package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/berkeli/kafka-cron/types"
	"github.com/berkeli/kafka-cron/utils"
	"github.com/stretchr/testify/require"
)

func Test_ExecuteCommand(t *testing.T) {
	t.Run("should execute hello world", func(t *testing.T) {
		want := "hello world\n"

		got, err := executeCommand("echo hello world")

		require.NoError(t, err)

		require.Equal(t, want, got)
	})

	t.Run("should throw error with invalid command", func(t *testing.T) {
		got, err := executeCommand("invalid command")

		require.Error(t, err)

		require.ErrorContains(t, err, "exit status 127")

		require.Equal(t, "", got)
	})
}

func Test_ProcessCommand(t *testing.T) {
	t.Run("should process command", func(t *testing.T) {
		producer := mocks.NewSyncProducer(t, nil)
		cmd := types.Command{
			Clusters:    []string{"cluster-a"},
			Description: "hello world",
			Command:     "echo hello world",
			MaxRetries:  3,
			Schedule:    "*/1 * * * *",
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)

		processCommand(producer, cmd, wg)
	})

	t.Run("erroneous command should result in retry", func(t *testing.T) {
		producer := mocks.NewSyncProducer(t, nil)
		expRetryTopic := "retry-topic"
		retryTopic = &expRetryTopic

		cmd := types.Command{
			Clusters:    []string{"cluster-a"},
			Description: "hello world",
			Command:     "invalid command",
			MaxRetries:  3,
			Schedule:    "*/1 * * * *",
		}

		producer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(val *sarama.ProducerMessage) error {
			if expRetryTopic != val.Topic {
				return fmt.Errorf("Expected topic %s, got %s", expRetryTopic, val.Topic)
			}

			var gotCmd types.Command

			bytes, err := val.Value.Encode()

			if err != nil {
				return err
			}

			err = json.Unmarshal(bytes, &gotCmd)

			if err != nil {
				return err
			}

			wantCmd := cmd
			wantCmd.MaxRetries = 2

			if err := utils.AssertCommandsEqual(t, wantCmd, gotCmd); err != nil {
				return err
			}

			return nil
		})

		wg := &sync.WaitGroup{}
		wg.Add(1)

		processCommand(producer, cmd, wg)
	})
}
