package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics map[string]interface{}

var (
	ScheduledCrons = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "scheduled_crons",
		Help: "Number of scheduled crons",
	})
	QueuedJobs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queued_jobs",
		Help: "Number of queued jobs",
	}, []string{"topic", "status"})
)

func InitMonitoring(port int) (Metrics, error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()
	return nil, nil
}
