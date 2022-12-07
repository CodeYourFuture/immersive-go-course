package utils

import (
	"fmt"
	"testing"

	"github.com/berkeli/kafka-cron/types"
)

func AssertCommandsEqual(t *testing.T, want, got types.Command) error {
	t.Helper()

	if want.Command != got.Command {
		return fmt.Errorf("expected command %s, got %s", want.Command, got.Command)
	}

	if want.Description != got.Description {
		return fmt.Errorf("expected description %s, got %s", want.Description, got.Description)
	}

	if want.Schedule != got.Schedule {
		return fmt.Errorf("expected schedule %s, got %s", want.Schedule, got.Schedule)
	}

	if want.MaxRetries != got.MaxRetries {
		return fmt.Errorf("expected max retries %d, got %d", want.MaxRetries, got.MaxRetries)
	}

	if len(want.Clusters) != len(got.Clusters) {
		return fmt.Errorf("expected %d clusters, got %d", len(want.Clusters), len(got.Clusters))
	}

	return nil
}
