package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 50051, "The server port")
	probeLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_average_latency",
		Help: "probe latency for request",
	},
		[]string{"endpoint"},
	)
	probeTotalRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probe_total_requests",
	},
		[]string{"endpoint"},
	)
)

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
	logger *log.Logger
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	var (
		totalElapsedTimeInMS        float32
		totalSuccessfulRequestCount uint64
	)

	type httpResponse struct {
		resp          *http.Response
		error         error
		elapsedTimeMS float32
	}

	httpResponseCh := make(chan httpResponse)
	var wg sync.WaitGroup
	for i := uint64(0); i < in.NumOfRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			resp, err := http.Get(in.GetEndpoint())
			if resp != nil {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						s.logger.Println("error closing the resp body")
					}
				}(resp.Body)
			}
			elapsed := time.Since(start)
			elapsedMsecs := float32(elapsed / time.Millisecond)
			httpResponseCh <- httpResponse{
				resp:          resp,
				elapsedTimeMS: elapsedMsecs,
				error:         err,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(httpResponseCh)
	}()

	for res := range httpResponseCh {
		if res.error != nil {
			continue
		}

		if res.resp.StatusCode == http.StatusOK {
			totalElapsedTimeInMS += res.elapsedTimeMS
			totalSuccessfulRequestCount++
		}
	}

	probeLatency.WithLabelValues(in.GetEndpoint()).Set(float64(totalElapsedTimeInMS / float32(totalSuccessfulRequestCount)))
	probeTotalRequests.WithLabelValues(in.GetEndpoint()).Set(float64(in.NumOfRequests))

	return &pb.ProbeReply{
		AverageLatencyMsecs:             totalElapsedTimeInMS / float32(totalSuccessfulRequestCount),
		TotalRequestCounts:              in.NumOfRequests,
		TotalRequestsWith_2XXStatusCode: totalSuccessfulRequestCount,
	}, nil
}

func main() {
	flag.Parse()
	logger := log.Default()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{
		logger: logger,
	})
	prometheus.MustRegister(probeLatency, probeTotalRequests)
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Println("attempting to start the metrics endpoint :4000")
		if err := http.ListenAndServe(":4000", nil); err != nil {
			logger.Printf("failed to start the metrics endpoint: %v\n", err)
			return
		}
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
