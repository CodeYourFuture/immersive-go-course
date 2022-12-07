package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/berkeli/kafka-cron/types"
	"github.com/berkeli/kafka-cron/utils"
)

func Test_PublishMessages(t *testing.T) {
	prod := mocks.NewSyncProducer(t, nil)

	clusters := []string{"cluster-a", "cluster-b"}
	wantMsg := "hello world"
	prefix := "jobs"
	topicPrefix = &prefix

	prod.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
		bytes, err := msg.Value.Encode()
		if err != nil {
			return err
		}

		gotMsg := string(bytes)

		if gotMsg != wantMsg {
			return fmt.Errorf("got message %s, want %s", gotMsg, wantMsg)
		}

		if msg.Topic != "jobs-cluster-a" {
			return fmt.Errorf("got topic %s, want %s", msg.Topic, "jobs-cluster-a")
		}

		return nil
	})

	prod.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
		bytes, err := msg.Value.Encode()
		if err != nil {
			return err
		}

		gotMsg := string(bytes)

		if gotMsg != wantMsg {
			return fmt.Errorf("got message %s, want %s", gotMsg, wantMsg)
		}

		if msg.Topic != "jobs-cluster-b" {
			return fmt.Errorf("got topic %s, want %s", msg.Topic, "jobs-cluster-b")
		}

		return nil
	})

	PublishMessages(prod, wantMsg, clusters)
}

func Test_CommandPublisher(t *testing.T) {
	t.Run("should publish command", func(t *testing.T) {
		cmd := types.Command{
			Clusters:    []string{"cluster-a"},
			Description: "hello world",
			Command:     "echo hello world",
			MaxRetries:  3,
			Schedule:    "*/1 * * * *",
		}

		prod := mocks.NewSyncProducer(t, nil)

		prod.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {

			gotCmd := types.Command{}

			bytes, err := msg.Value.Encode()

			if err != nil {
				return err
			}

			err = json.Unmarshal(bytes, &gotCmd)

			if err != nil {
				return err
			}

			err = utils.AssertCommandsEqual(t, cmd, gotCmd)

			if err != nil {
				return err
			}

			return nil
		})

		cmdPub := CommandPublisher{
			Command:   cmd,
			publisher: prod,
		}

		cmdPub.Run()
	})
}
