package main

import (
	"flag"
	"log"
	"strings"
)

var (
	REPLICATED = "The provided memcached cluster is replicated"
	SHARDED    = "The provided memcached cluster is sharded"
)

type arrayFlag []string

func (i *arrayFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlag) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}

func main() {
	var (
		routerPort string
		nodePorts  arrayFlag
	)
	flag.StringVar(&routerPort, "mcrouter", "11211", "port of mcrouter, defaults to 11211")
	flag.Var(&nodePorts, "memcacheds", "comma separated list of ports for memcached instances in the cluster")
	flag.Parse()

	mcrouter := NewCacheService(routerPort)

	memcacheds := make(map[string]ICacheService)

	for _, port := range nodePorts {
		memcacheds[port] = NewCacheService(port)
	}

	log.Println(checkIfSharded(mcrouter, memcacheds))
}

func checkIfSharded(router ICacheService, nodes map[string]ICacheService) string {
	router.Set("foo", "bar")

	for _, node := range nodes {
		val, err := node.Get("foo")
		if err != nil || val != "bar" {
			return SHARDED
		}
	}

	return REPLICATED
}
