package configs

// Hardcoded, for demo purposes.

import (
	"fmt"
)

func GetClusters() []string {
	clusters := []string{"cluster-a", "cluster-b"}
	return clusters
}

func GetTopicName(cluster string) string {
	return fmt.Sprintf("%s-cronjobs", cluster)
}

func GetRetryTopicName(cluster string) string {
	return fmt.Sprintf("%s-cronjobs-retry", cluster)
}

func GetDLQName() string {
	return "dead-letter-queue"
}
