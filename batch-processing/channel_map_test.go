package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestIdentity(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{1, 2, 3}

	c1 := make(chan int, len(input))
	for _, v := range input {
		c1 <- v
	}
	close(c1)

	c2 := Map(
		c1,
		func(v int) int {
			return v
		},
	)

	output := make([]int, 0, len(input))
	for v := range c2 {
		output = append(output, v)
	}
	sort.Ints(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}

func TestDoNothing(t *testing.T) {
	input := []int{}
	expected := []int{}

	c1 := make(chan int, len(input))
	close(c1)

	c2 := Map(
		c1,
		func(v int) int {
			return v
		},
	)

	output := make([]int, 0, len(input))
	for v := range c2 {
		output = append(output, v)
	}
	sort.Ints(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}

func TestStruct(t *testing.T) {
	type Task[T any] struct {
		Result T
		Error  error
	}

	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	c1 := make(chan Task[int], len(input))
	for _, v := range input {
		c1 <- Task[int]{
			Result: v,
			Error:  nil,
		}
	}
	close(c1)

	c2 := Map(
		c1,
		func(t Task[int]) Task[string] {
			return Task[string]{fmt.Sprintf("%d", t.Result), nil}
		},
	)

	output := make([]string, 0, len(expected))
	for t := range c2 {
		output = append(output, t.Result)
	}
	sort.Strings(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}

func TestError(t *testing.T) {
	type Task[T any] struct {
		Result T
		Error  error
	}

	input := []string{"1"}
	expected := []Task[int]{
		{
			Result: 0,
			Error:  errors.New("Something went wrong"),
		},
	}

	c1 := make(chan Task[string], len(input))
	for _, s := range input {
		c1 <- Task[string]{
			Result: s,
			Error:  nil,
		}
	}
	close(c1)

	c2 := Map(
		c1,
		func(t Task[string]) Task[int] {
			// Imaginary error happens here
			return Task[int]{
				Result: 0,
				Error:  errors.New("Something went wrong"),
			}
		},
	)

	output := make([]Task[int], 0, len(input))
	for t := range c2 {
		output = append(output, t)
	}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}

func TestChain(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{4, 6, 8}

	c1 := make(chan int, len(input))
	for _, v := range input {
		c1 <- v
	}
	close(c1)

	c2 := Map(
		c1,
		func(v int) int {
			return v + 1
		},
	)

	c3 := Map(
		c2,
		func(v int) int {
			return v * 2
		},
	)

	output := make([]int, 0, len(expected))
	for v := range c3 {
		output = append(output, v)
	}
	sort.Ints(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}

func TestIdentityUnbuffered(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{1, 2, 3}

	c1 := make(chan int)

	c2 := Map(
		c1,
		func(v int) int {
			return v
		},
	)

	for _, v := range input {
		c1 <- v
	}
	close(c1)

	output := make([]int, 0, len(input))
	for v := range c2 {
		output = append(output, v)
	}
	sort.Ints(output)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Expected %v but got %v", expected, output)
	}
}
