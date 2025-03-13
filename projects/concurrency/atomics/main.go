package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var x atomic.Int32

func increment(wg *sync.WaitGroup) {
	x.Add(1)
	wg.Done()
}

func main() {
	var w sync.WaitGroup
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go increment(&w)
	}
	w.Wait()
	fmt.Println("final value of x", x.Load())
}
