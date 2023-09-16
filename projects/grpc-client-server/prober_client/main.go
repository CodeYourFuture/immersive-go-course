// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/url"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr            = flag.String("addr", "localhost:50051", "the address to connect to")
	num_of_requests = flag.Uint64("num_of_requests", 1, "number of requests")
	endpoint        = flag.String("endpoint", "", "endpoint to probe")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProberClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	e := *endpoint
	if e == "" {
		log.Fatal("endpoint cannot be empty")
	}

	u, err := url.Parse(e)
	if err != nil || u.Scheme == "" || u.Host == "" {
		log.Fatalf("invalid url: %v", u.String())
	}

	_, err = net.LookupIP(u.Host)
	if err != nil {
		log.Fatalf("error looking up ip: %v", err)
	}

	r, err := c.DoProbes(ctx, &pb.ProbeRequest{Endpoint: u.String(), NumOfRequests: *num_of_requests})
	if err != nil {
		log.Fatalf("could not probe: %v", err)
	}
	log.Printf("Response Time: %f", r.GetAverageLatencyMsecs())
	log.Printf("Successful Requests: %d", r.GetTotalRequestsWith_2XXStatusCode())
	log.Printf("Total Request Count: %d", r.GetTotalRequestCounts())
}
