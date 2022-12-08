package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics map[string]interface{}

var (
	ScheduledCrons = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scheduled_crons",
		Help: "Number of scheduled crons",
	})
	QueuedJobs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "queued_jobs",
		Help: "Number of queued jobs",
	}, []string{"topic", "status"})
)

func InitMonitoring(port int) (Metrics, error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
	return nil, nil
}
