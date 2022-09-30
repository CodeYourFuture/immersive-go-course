package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestIdentity(t *testing.T) {
	expected := []int{1, 2, 3}
	c1 := make(chan int, len(expected))
	for _, v := range expected {
		c1 <- v
	}
	close(c1)

	c2 := Map(
		c1,
		func(v int) int {
			return v
		},
	)

	output := make([]int, 0, len(expected))
	for v := range c2 {
		output = append(output, v)
	}
	sort.Ints(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}
