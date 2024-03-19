# Implementation of interfaces project

## How to run

The primary artifacts of this project are a struct implementation, and a series of tests.

To run the tests for this project, `cd` to this directory and run `go test ./...`. You should see all the tests pass.

## Comparing the tests

There are two different implementations of this test suite.

In [`buffer_table_test.go`](./buffer_table_test.go) there is a table-driven test. In [`buffer_test.go`](./buffer_test.go) there is a non-table-driven test.

Often we use [table-driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests) in Go because they make it easy to do similar things many times (factoring out commonality to be [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)), and make it really easy to add extra test-cases.

In this example, it's not obvious that table-driven tests are better than non-table-driven tests. Take a bit of time to look at both tests and think about the trade-offs involved here:
* Which set of tests would you prefer to write?
* Which set of tests would you prefer to read?
* Which set of tests would you prefer to maintain?

Think about this yourself before reading the next section.

### Advantages of non-table-driven tests

In this example, the table-driven tests are harder to read and verify they're obviously correct.

In the non-table-driven tests, each test reads top-to-bottom quite clearly.

In the table-driven tests, you need to understand quite a lot about the `operations` field. If you want to verify the exact behaviours, you need to read several interface implementations, and understand how they interact.

We don't often write tests for our tests, so in tests, we strongly value being able to see that they're obviously correct just by reading them. One could argue that the non-table-driven test does this better, in this example.

### Advantages of table-driven tests

#### Adding more tests

In this example, the table-driven tests are more re-usable. If we wanted to add lots more tests with different permutations ("write then read then write again then write again then read"), it would be easier to add these extra tests in the table-driven way. Making it easy to write more tests is generally a good thing.

#### Changing all the tests

Imagine we changed the implementation of our buffer (e.g. added an `error` return to some function, or changed how we construct a buffer). In the table-driven tests, we only have a few actual uses of the buffer API - probably one per `operation` implementation.

In the non-table-driven test version, we need to make that change in lots of places. By avoiding repeating ourselves and factoring away the commonality, we reduced the number of places we need to change our tests if something changes in our implementation.

### Summary

Neither approach is obviously _better_. They have different trade-offs. This is a common situation when writing software.

If we were testing something a little simpler, like testing the `FilteringPipe` type, the tests would have fewer differences between them (e.g. just differing in what's written, rather than having lots of different operations), a table-driven test would probably be a more obvious winner.

We need to choose what we're optimising for: Are we optimising for making it easy to add lots more complex ordering tests? Are we optimising for making it easy to read and verify the test works? How do we imagine needing to change the code in the future, and what will make _that_ easy?

## Implementation notes on `OurByteBuffer`

#### Struct vs Pointer

When implementing a method on a struct, we can choose to implement it on the struct type itself (`func (b OurByteBuffer) someMethod()`) or on a pointer to the struct type (`func (b *OurByteBuffer) someMethod()`).

We chose to implement the methods on a pointer to the struct, rather than just the struct.

This is because we have information which we need to be preserved across the different methods. If we implemented the methods on the struct, each call to `Write` would get a new copy of `OurByteBuffer`, so when we update `b.bytes` it would only get updated _in that copy_.
