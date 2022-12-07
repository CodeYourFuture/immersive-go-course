package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics map[string]interface{}

var (
	JobStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "job_status",
		Help: "Status of jobs",
	}, []string{"status", "topic", "description"})
	JobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "job_duration_seconds",
		Help: "Duration of jobs in seconds",
	}, []string{"topic", "description"})
)

func InitMonitoring(port int) (Metrics, error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()
	return nil, nil
}