package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"raft"
	"strings"
	"syscall"
	"time"
)

const port = 7600

func main() {
	addr := flag.String("dns", "raft", "dns address for raft cluster")
	if_addr := flag.String("if", "eth0", "use IPV4 address of this interface") // eth0 works on docker, may vary for other platforms

	if addr == nil || *addr == "" {
		fmt.Printf("Must supply dns address of cluster\n")
		os.Exit(1)
	}

	id := getOwnAddr(*if_addr)
	fmt.Printf("My address/node ID is %s\n", id)

	ready := make(chan interface{})
	storage := raft.NewMapStorage()
	commitChan := make(chan raft.CommitEntry)
	server := raft.NewServer(id, id, storage, ready, commitChan, port)
	server.Serve(raft.NewKV())

	ips, err := net.LookupIP(*addr)
	if err != nil {
		fmt.Printf("Could not get IPs: %v\n", err)
		os.Exit(1)
	}

	// Connect to all peers with appropriate waits
	// TODO: we only do this once, on startup - we really should periodically check to see if the DNS listing for peers has changed
	for _, ip := range ips {
		// if not own IP
		if !ownAddr(ip, id) {
			peerAddr := fmt.Sprintf("%s:%d", ip.String(), port)

			connected := false
			for rt := 0; rt <= 3 && !connected; rt++ {
				fmt.Printf("Connecting to peer %s\n", peerAddr)
				err = server.ConnectToPeer(peerAddr, peerAddr)
				if err == nil {
					connected = true
				} else { // probably just not started up yet, retry
					fmt.Printf("Error connecting to peer: %+v", err)
					time.Sleep(time.Duration(rt+1) * time.Second)
				}
			}
			if err != nil {
				fmt.Printf("Exhausted retries connecting to peer %s", peerAddr)
				os.Exit(1)
			}
		}
	}

	close(ready) // start raft server, peers are connected

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulShutdown
	server.DisconnectAll()
	server.Shutdown()
}

func getOwnAddr(intf string) string {
	ifs, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Could not get intf: %v\n", err)
		os.Exit(1)
	}

	for _, cif := range ifs {
		if cif.Name == intf {
			ads, _ := cif.Addrs()
			for _, addr := range ads {
				if isIPV4(addr.String()) {
					ip := getIP(addr.String())
					return ip.String()
				}

			}
		}
	}

	fmt.Printf("Could not find intf: %s\n", intf)
	os.Exit(1)
	return ""
}

func isIPV4(addr string) bool {
	parts := strings.Split(addr, "::")
	return len(parts) == 1
}

func getIP(addr string) net.IP {
	parts := strings.Split(addr, "/")
	return net.ParseIP(parts[0])
}

func ownAddr(ip net.IP, myAddr string) bool {
	res := ip.String() == myAddr
	return res
}
