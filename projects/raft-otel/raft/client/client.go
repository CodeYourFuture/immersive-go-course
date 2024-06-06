package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"raft/raft_proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 7600

func main() {
	addr := flag.String("dns", "raft", "dns address for raft cluster")

	if addr == nil || *addr == "" {
		fmt.Printf("Must supply dns address of cluster\n")
		os.Exit(1)
	}

	time.Sleep(time.Second * 5) // wait for raft servers to come up

	ips, err := net.LookupIP(*addr)
	if err != nil {
		fmt.Printf("Could not get IPs: %v\n", err)
		os.Exit(1)
	}

	clients := make([]raft_proto.RaftKVServiceClient, 0)

	for _, ip := range ips {
		fmt.Printf("Connecting to %s\n", ip.String())
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip.String(), port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("%v", err)
		}
		client := raft_proto.NewRaftKVServiceClient(conn)
		clients = append(clients, client)
	}

	for {
		for _, c := range clients {
			n := time.Now().Second()
			res, err := c.Set(context.TODO(), &raft_proto.SetRequest{Keyname: "cursec", Value: fmt.Sprintf("%d", n)})
			fmt.Printf("Called set cursec %d, got %v, %v\n", n, res, err)

			time.Sleep(1 * time.Second) // allow consensus to happen

			getres, err := c.Get(context.TODO(), &raft_proto.GetRequest{Keyname: "cursec"})
			fmt.Printf("Called get cursec, got %v, %v\n", getres, err)
		}
		time.Sleep(5 * time.Second)
	}
}
