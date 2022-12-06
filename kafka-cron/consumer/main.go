package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"os/signal"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/berkeli/kafka-cron/types"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("jobs-cluster-a").String()
	retryTopic = kingpin.Flag("retryTopic", "Retry topic name").Default("jobs-cluster-a-failed").String()
	partition  = kingpin.Flag("partition", "Partition number").Default("0").Int32()
)

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
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
				log.Println("Starting a job for: ", cmd.Description)
				out, err := executeCommand(cmd.Command)
				if err != nil {
					log.Printf("Command: %s resulted in error: %s\n", cmd.Command, err)
					if cmd.MaxRetries > 0 {
						cmd.MaxRetries--
						log.Printf("Retrying command: %s, %d retries left\n", cmd.Command, cmd.MaxRetries)
						
					} else {
						log.Printf("Command: %s failed, no more retries left\n", cmd.Command)
					}
				}
				log.Println(out)
			case <-signals:
				log.Println("Interrupt is detected")
				chDone <- true
			}
		}
	}()
	<-chDone
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}
