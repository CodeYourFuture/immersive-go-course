package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/bradfitz/gomemcache/memcache"
)


func main() {

	var HOST = "localhost"
	var mrPort = flag.String("mcrouter", "11211", "port for mc router")
	var mcPort = flag.String("memcacheds", "11212", "server ports for memcacheds \n more than one port is required")

	flag.Parse()

	mcPorts := strings.Split(*mcPort, ",")

	if len(mcPorts) <= 1 {
		fmt.Fprint(os.Stderr, "more than one cache is required\n")
		flag.Usage()
		os.Exit(1)
	}
	routerServer := fmt.Sprintf("%s:%s", HOST, *mrPort)

	routerClient := memcache.New(routerServer)
	err := routerClient.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	if err != nil {
		panic(err)
	}

	for _, cache := range mcPorts {
		cache = fmt.Sprintf("%s:%s", HOST, cache)
		mc := memcache.New(cache)
		_, err := mc.Get("foo")

		if err != nil {
			if errors.Is(err, memcache.ErrCacheMiss) {
				fmt.Println("Cache is sharded")
				os.Exit(0)
			} else {
				panic(err)
			}
		}

	}
	fmt.Println("Cache is replicated")
}
