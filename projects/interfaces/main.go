package main

import (
	"bytes"
	"fmt"
)

// RAM: variables / data structure / buffer
// Stack: part of RAM - fast - allocation: Lifo - free: after finish the func
// Heap: part of RAM - slower than stack - allocation: dynamic - free: manually or garbage collector
// Array: fix size - allocate in compile time - store in stack
// Slice: getting size dynamically - allocate in run time(execute) - header will be in stack and data will be in heap specially when created by make()

// arr:=[3]int{1,2,3}  arr in stack
// s:=arr[1:4]  s in stack (because pointing to stack)
// x:=make([]int,3)  x in heap

// buffer: it is a structure data to store data as a block of bytes or string between I/O instead of reading byte to byte from hard disk
// one thread can read from buffer, another thread can write to the buffer
// input buffer: collect data from keyboard / hard disk / network and write on buffer
// output buffer: save before sending to screen / hard disk / network
// package bytes: working with buffer - using slice


func main() {
	b:=bytes.NewBufferString("Hi")
	c:=b.Bytes()
	b.Write([]byte(" Hello"))
	d:=b.Bytes()
	x:=make([]byte,4)
	// x:=[]byte("test")  
	l,_:=b.Read(x) // It cut from start of my buffer and write it in an empty slice (or rewrite it)
	
	fmt.Println(b,*b,c,d,x,l)
	b.Reset()
	fmt.Println(b)
}

