<!--forhugo
+++
title="Interfaces"
+++
forhugo-->

In this project you're going to get familiar with interfaces in Go.

Timebox: 2 days

## Objectives:

- Describe when interfaces are useful
- Write tests in Go
- Implement existing interfaces on your own struct

## Project

Read [this blog post about interfaces in Go](https://www.alexedwards.net/blog/interfaces-explained).

Interfaces are used to build abstractions. Abstractions in programming are when we can just think about _what_ something does, and not _how_ it does it.

For instance, when programming we never really think about how memory works, we just focus on the fact that we can store values.

Interfaces allow us to write code which relies on other types having some behaviour, regardless of how it has that behaviour. For instance, we can write code which relies on being able to read some bytes, without needing to worry about whether those bytes came from a file, a response to an HTTP request, or some variable we stored earlier in our code.

Two common interfaces in Go are [`io.Reader`](https://pkg.go.dev/io#Reader) and [`io.Writer`](https://pkg.go.dev/io#Writer). We can write a function which accepts an `io.Reader` as a parameter, and that function can know it can read bytes from that parameter, without worrying about where they're coming from.

In this exercise, we're going to implement two [structs](https://go.dev/tour/moretypes/2) in Go, and have them implement some interfaces.

Make a new Go program, and we'll begin:

### Testing `bytes.Buffer`

The Go standard library has a type called [`bytes.Buffer`](https://pkg.go.dev/bytes#Buffer). Have a read of its documentation.

A `bytes.Buffer` allows you store some bytes, and retrieve them later. You can do this in a number of ways:
* You can store some bytes when you create one, by using `bytes.NewBufferString`.
* You can add bytes to the buffer by using `bytes.Buffer.Write` - this is the function on the `io.Writer` interface, so `bytes.Buffer` implementes `io.Writer`.
* You can read all of the bytes in the buffer using `bytes.Buffer.Bytes()`.
* You can read some bytes from the buffer into a slice by using `bytes.Buffer.Read` - this is the function on the `io.Reader` interface, so `bytes.Buffer` implements `io.Reader`.
There are also other ways of accessing and manipulating data in `bytes.Buffer`, but we won't worry about them for now.

The first task is to implement unit tests for this type.

You should write unit tests which show at least the following (you can write more if you want!):
* If you make a buffer named `b` containing some bytes, calling `b.Bytes()` returns the same bytes you created it with.
* If you write some extra bytes to that buffer using `b.Write()`, a call to `b.Bytes()` returns both the initial bytes and the extra bytes.
* If you call `b.Read()` with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.
* If you call `b.Read()` with a slice smaller than the contents of the buffer, some of the bytes are read. If you call it again, the next bytes are read.

### Implementing our own `bytes.Buffer`

Now that we have some tests for the behaviour of this struct, let's implement our own copy.

Make a new struct, we'll call it `OurByteBuffer`.

Change the tests we have to refer to our type instead of `bytes.Buffer`.

Implement functions on `OurByteBuffer` so that all the tests pass.

(Hint: We need to implement our functions on _pointers_ to `OurByteBuffer`, not to _copies_ of it).

#### Limitations

When implementing a type we're going to use in different places, it's important to think about the limitations of our type.

Some examples of limitations for `OurByteBuffer` may be:
- Is there a maximum about of data an `OurByteBuffer` can store?
- What operations are safe or unsafe to perform concurrently on an `OurByteBuffer` from different threads?
- Are there any important performance characteristics? e.g. is it much faster or slower to `Write` one large amount of data to it than write the same amount of data one byte at a time?
- Are there any important memory characteristics? e.g. does an `OurByteBuffer` always retain all data that was stored in it, or does it free some of its memory after it's been read?

It can be useful to note these limitations in a doc comment for the type, or method, they apply to. This will help someone using your type to use it the best way. Make sure you write these docs.

When you need to pick (or implement) a type which implements an interface, it's important to think about these limitations.

### Implementing a custom filter

Implement a new struct called `FilteringPipe`.

When it's constructed, it should be constructed with an `io.Writer` which we'll store as a field in the struct.

Implement `io.Writer` for `FilteringPipe`, which will write whatever is written to it to the `io.Writer` it was constructed with, _except_ it should skip any numbers.

So if we call:

```go
filteringPipe := NewFilteringPipe(someWriter)
filteringPipe.Write([]byte("start=1, end=10"))
```

we'll end up writing `start=, end=` to `someWriter`.

Make sure to write some tests for this type, too. These should probably be [table-driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests).

## Thinking about trade-offs

When you're finished, have a read of the [sample implementation in the `impl/interfaces` branch](https://github.com/CodeYourFuture/immersive-go-course/tree/impl/interfaces/projects/interfaces). Read through the [IMPLEMENTATION.md](https://github.com/CodeYourFuture/immersive-go-course/blob/impl/interfaces/projects/interfaces/IMPLEMENTATION.md) file. It talks a lot about two different implementations of the tests for this project. Think about the trade-offs involved in both approaches.
