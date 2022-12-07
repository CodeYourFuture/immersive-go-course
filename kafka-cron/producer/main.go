package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	err = ScheduleJobs(producer, cmds)

	if err != nil {
		log.Panic(err)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
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

func ScheduleJobs(prod sarama.SyncProducer, cmds []types.Command) error {
	c := cron.New()
	for _, cmd := range cmds {
		fmt.Println("scheduling: ", cmd.Description)

		sch, err := cron.ParseStandard(cmd.Schedule)

		if err != nil {
			return err
		}

		job := CommandPublisher{
			Command:   cmd,
			publisher: prod,
		}

		c.Schedule(sch, &job)
		ScheduledCrons.Inc()
	}
	c.Start()
	return nil
}

func PublishMessages(prod sarama.SyncProducer, msg string, clusters []string) error {
	for _, cluster := range clusters {
		topic := fmt.Sprintf("%s-%s", *topicPrefix, cluster)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(uuid.New().String()),
			Value: sarama.StringEncoder(msg),
		}
		_, _, err := prod.SendMessage(msg)
		if err != nil {
			QueuedJobs.WithLabelValues(topic, "error").Inc()
			return err
		}
		QueuedJobs.WithLabelValues(topic, "success").Inc()
	}
	return nil
}

type CommandPublisher struct {
	types.Command
	publisher sarama.SyncProducer
}

func (c *CommandPublisher) Run() {
	fmt.Println("Running command: ", c.Description)
	msgString, err := json.Marshal(&c)
	if err != nil {
		log.Println(fmt.Errorf("error marshalling command: %v", err))
	}
	err = PublishMessages(c.publisher, string(msgString), c.Clusters)

	if err != nil {
		log.Println(fmt.Errorf("error publishing command: %v", err))
	}
}
