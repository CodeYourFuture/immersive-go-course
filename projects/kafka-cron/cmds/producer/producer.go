package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"kafka-cron/configs"
	"kafka-cron/message"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type cronjob struct {
	crontab string
	command string
	name    string
	cluster string
	retries int
}

func main() {
	var bootstrapServers, configPath string
	var partitions int

	flag.StringVar(&bootstrapServers, "bootstrap.servers", "127.0.0.1:9092", "bootstrap servers")
	flag.StringVar(&configPath, "config", "./data/cronjobs.txt", "path to cronjob spec file")
	flag.IntVar(&partitions, "partitions", 1, "number of partitions per topic")
	flag.Parse()

	// Create Producer instance
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Create topics if needed
	for _, cluster := range configs.GetClusters() {
		CreateTopic(p, configs.GetTopicName(cluster), partitions)
		CreateTopic(p, configs.GetRetryTopicName(cluster), partitions)
	}
	CreateTopic(p, configs.GetDLQName(), partitions)

	// Parse cronjobs and create schedule
	cronjobs := parseCronjobs(configPath)
	fmt.Printf("cronjobs: %v\n", cronjobs)

	var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	c := cron.New(cron.WithParser(secondParser), cron.WithChain())
	for _, j := range cronjobs {
		j := j // for closure
		c.AddFunc(j.crontab, func() {
			queueJob(p, j.command, j.command, j.cluster, j.retries)
		})
		fmt.Printf("cronjobs: started cron for %+v\n", j)
	}
	go c.Start()

	fmt.Printf("cronjobs: cron entries are %+v\n", c.Entries())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

func queueJob(p *kafka.Producer, command string, name string, cluster string, retries int) {
	msg := message.CronjobMessage{
		Command:  strings.TrimSpace(command),
		Exectime: time.Now().Format(time.RFC3339),
		Name:     strings.TrimSpace(name),
		Retries:  retries,
	}

	jmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	key, _ := uuid.New().MarshalText()
	topic := configs.GetTopicName(cluster)
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

func parseCronjobs(path string) []cronjob {
	result := make([]cronjob, 0)

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		val := scanner.Text()
		fields := strings.Split(val, ",")

		if len(fields) != 5 {
			log.Fatal(fmt.Errorf("wrong format %s %d fields, expect 5", val, len(fields)))
		}

		retries, err := strconv.Atoi(strings.TrimSpace(fields[4]))
		if err != nil {
			log.Fatal(fmt.Errorf("wrong format could not parse %s as int (num retries)", fields[4]))
		}

		cron := cronjob{fields[0], fields[1], fields[2], strings.TrimSpace(fields[3]), retries}
		result = append(result, cron)
	}
	return result
}

// CreateTopic creates a topic using the Admin Client API
func CreateTopic(p *kafka.Producer, topic string, partitions int) {

	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		os.Exit(1)
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Printf("ParseDuration(60s): %s", err)
		os.Exit(1)
	}
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Admin Client request error: %v\n", err)
		os.Exit(1)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Printf("%v\n", result)
	}
	a.Close()

}
