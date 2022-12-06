package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
	"github.com/berkeli/kafka-cron/types"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topicPrefix = kingpin.Flag("topicPrefix", "Topic prefix, e.g. jobs will create topics for each cluster in the format jobs-cluster-a").Default("jobs").String()
	cronPath    = kingpin.Flag("cronPath", "Path to cron file").Default("./cron.yaml").String()
)

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()
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
	cmds, err := ReadConfig(*cronPath)

	if err != nil {
		log.Panic(err)
	}

	c := cron.New()

	err = ScheduleJobs(c, producer, cmds)

	if err != nil {
		log.Panic(err)
	}

	c.Run()
}

func ReadConfig(path string) ([]types.Command, error) {
	var cmds struct {
		Cron []types.Command `json:"cron" yaml:"cron"`
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cmds)
	if err != nil {
		return nil, err
	}
	return cmds.Cron, nil
}

func ScheduleJobs(c *cron.Cron, prod sarama.SyncProducer, cmds []types.Command) error {
	var runErr error
	for _, cmd := range cmds {
		fmt.Println("scheduling: ", cmd.Description)
		_, err := c.AddFunc(cmd.Schedule, func() {
			log.Println("Publishing command: ", cmd.Description)
			msgString, err := json.Marshal(cmd)
			if err != nil {
				log.Println(fmt.Errorf("error marshalling command: %v", err))
				runErr = err
			}
			err = PublishMessages(prod, string(msgString), cmd.Clusters)

			if err != nil {
				log.Println(fmt.Errorf("error publishing command: %v", err))
				runErr = err
			}
		})
		if err != nil {
			return err
		}
	}
	return runErr
}

func PublishMessages(prod sarama.SyncProducer, msg string, clusters []string) error {
	for _, cluster := range clusters {
		msg := &sarama.ProducerMessage{
			Topic: fmt.Sprintf("%s-%s", *topicPrefix, cluster),
			Key:   sarama.StringEncoder(uuid.New().String()),
			Value: sarama.StringEncoder(msg),
		}
		_, _, err := prod.SendMessage(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
