package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/berkeli/kafka-cron/types"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("jobs-cluster-a").String()
	retryTopic = kingpin.Flag("retryTopic", "Retry topic name").Default("jobs-cluster-a-retries").String()
	partition  = kingpin.Flag("partition", "Partition number").Default("0").Int32()
)

func main() {

	InitMonitoring(2112)

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}

	config = sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, *partition, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	chDone := make(chan bool)
	wgWorkers := sync.WaitGroup{}
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				var cmd types.Command
				err := json.Unmarshal(msg.Value, &cmd)
				if err != nil {
					log.Println(err)
				}
				wgWorkers.Add(1)
				// TODO: Add a worker pool with semaphore? What if jobs are dependent on each other?
				processCommand(producer, cmd, &wgWorkers)
			case <-signals:
				chDone <- true
				return
			}
		}
	}()
	<-chDone
	log.Println("Interrupt is detected, shutting down gracefully...")
	wgWorkers.Wait()
}

func processCommand(producer sarama.SyncProducer, cmd types.Command, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Starting a job for: ", cmd.Description)
	//metrics
	timer := prometheus.NewTimer(JobDuration.WithLabelValues(*topic, cmd.Description))
	defer timer.ObserveDuration()
	JobStatus.WithLabelValues("new", *topic, cmd.Description).Inc()

	out, err := executeCommand(cmd.Command)
	if err != nil {

		log.Printf("Command: %s resulted in error: %s\n", cmd.Command, err)
		if cmd.MaxRetries > 0 {
			JobStatus.WithLabelValues("retry", *topic, cmd.Description).Inc()
			cmd.MaxRetries--
			log.Printf("Retrying command: %s, %d retries left\n", cmd.Command, cmd.MaxRetries)
			cmdBytes, err := json.Marshal(cmd)

			if err != nil {
				log.Println(err)
			}
			_, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Topic: *retryTopic,
				Key:   sarama.StringEncoder(uuid.New().String()),
				Value: sarama.ByteEncoder(cmdBytes),
			})

			if err != nil {
				log.Println(err)
			}

		} else {
			JobStatus.WithLabelValues("failed", *topic, cmd.Description).Inc()
			log.Printf("Command: %s failed, no more retries left\n", cmd.Command)
		}
	}
	JobStatus.WithLabelValues("success", *topic, cmd.Description).Inc()
	log.Println(out)
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}
