package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"kafka-cron/configs"
	"kafka-cron/message"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	var bootstrapServers, clusterName string
	var retry bool
	flag.StringVar(&bootstrapServers, "bootstrap.servers", "127.0.0.1:9092", "bootstrap servers")
	flag.StringVar(&clusterName, "cluster", "cluster-a", "cluster name")
	flag.BoolVar(&retry, "retry", false, "retry worker")
	flag.Parse()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          "foo",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := configs.GetTopicName(clusterName)
	if retry {
		topic = configs.GetRetryTopicName(clusterName)
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s %v\n", topic, err)
		os.Exit(1)
	} else {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}

	// Create Producer instance for retrying
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			recordKey := string(msg.Key)
			recordValue := msg.Value
			data := message.CronjobMessage{}
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			fmt.Printf("Consumed record with key %s and value %s\n", recordKey, recordValue)

			args := strings.Split(data.Command, " ")
			cmd := exec.Command(args[0], args[1:]...)
			err = cmd.Run()

			if err != nil {
				fmt.Printf("Failed to run command %s - re-enqueuing on retry queue: %v\n", data.Command, err)
				if data.Retries > 1 {
					queueRetryJob(p, data, msg.Key, clusterName)
				} else {
					// todo dlq
					fmt.Printf("To the dead letter queue\n")
					queueDLQ(p, data, msg.Key)
				}
			} else {
				fmt.Printf("Command ran OK %s\n", data.Command)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func queueRetryJob(p *kafka.Producer, msg message.CronjobMessage, key []byte, cluster string) {
	msg.Retries = msg.Retries - 1

	jmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	topic := configs.GetRetryTopicName(cluster)
	fmt.Printf("queuing job: %s on topic %s\n", string(jmsg), topic)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          jmsg},
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func queueDLQ(p *kafka.Producer, msg message.CronjobMessage, key []byte) {
	msg.Retries = 0

	jmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	topic := configs.GetDLQName()
	fmt.Printf("queuing job: %s on topic %s\n", string(jmsg), topic)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          jmsg},
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
