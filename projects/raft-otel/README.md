<!--forhugo
+++
title="RAFT implementation"
+++
forhugo-->

# RAFT implementation with Distributed Tracing

In this project we're going to build (or reuse) an implementation of RAFT, a distributed consensus algorithm, and
we are going to use distributed tracing to understand its behaviour.

Learning Objectives

- How is distributed tracing different from logging and from metrics?
- How can we use distributed tracing to get a detailed understanding of complex application behaviour?
- How can costs of distributed tracing be minimized?

Timebox: 5 days

## Project

### Background on Distributed Consensus and RAFT

In a program, sometimes we need to lock a resource. Think about a physical device, like a printer: we only want one program to print at a time.
Locking applies to lots of other kinds of resources too, often when we need to update multiple pieces of data in consistent ways in a multi-threaded context.

We need to be able to do [locking in distributed systems](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html) as well.

It turns out that in terms of computer science, distributed locking is theoretically equivalent to reliably electing a single leader (like in database replication, for example). It is also logically the same as determining the exact order that a sequence of events on different machines in a distributed system occur. All of these are useful. All of these are versions of the same problem: distributed consensus. Distributed consensus in computer science means reaching agreement on state among multiple processes on different machines in the presence of unreliable components (such as networks).

Distributed consensus gives us the strongest possible guarantees of data consistency. If we run a distributed consensus protocol on a cluster of machines, as long as a majority of those machines remain available then:

- we can never lose data that has been written
- we will never see 'stale' data that has been modified by another process

Contrast this with a system such as memcached: if a node becomes unavailable we can lose data that has been written. If multiple cache nodes store the same value,
the values can become out of sync if a node becomes temporarily unavailable and misses an update. Replicated databases can also lose data or become inconsistent.
This is not possible with a distributed consensus system. The downside of these strong guarantees is that they require more network round-trips and therefore have higher latency and lower throughput.
Because of this, they are normally used only for critical data.

Distributed consensus is a key component of widely-used software infrastructure such as Zookeeper (used by Kafka), etcd (used by Kubernetes),
Vitess (which can use Consul, etcd or Zookeeper) and the Consul service catalog.

Typically, these services use distributed consensus to store key information such as:

- which database instances are the leader for a given shard (Vitess)
- which brokers are members of a Kafka cluster, and which brokers are leaders of each partition (Kafka)
- cluster state and configuration (Kubernetes)

What these usecases have in common is that if there is inconsistency in how the system views this kind of state, it has far-reaching consequences.
For example, if two database instances both act as leaders for a shard in Vitess, we would have a 'split-brain' situation, where both instances
might commit different sets of transactions. These kinds of situations are very difficult to resolve without data loss.

RAFT is a specific protocol that implements distributed consensus. It has been designed specifically to be easier
to understand than earlier protocols (such as Paxos). However, it remains quite a complex protocol. The [RAFT website](https://raft.github.io) provides
a set of resources for understanding how RAFT works. Read the [RAFT paper](https://raft.github.io/raft.pdf), and
watch one of the introductory RAFT talks linked from the [RAFT website](https://raft.github.io).

### Background on Distributed Tracing and Open Telemetry

In [other projects](https://github.com/CodeYourFuture/immersive-go-course/tree/main/kafka-cron) we have added metrics to our programs and collected and
aggregated those metrics using the Prometheus monitoring tool. Metrics are a widely-used methodology for understanding the behaviour of our systems at
a statistical level: what percentage of requests are being completed successfully, what is the 90th percentile latency, what is our current cache hit
rate or queue length. These kinds of queries are very useful for telling us whether our systems seem to be healthy overall or not, and, in some cases,
may provide useful insights into problems or inefficiencies.

However, one thing that metrics are not normally very good for is understanding how user experience for a system may vary between different types of
requests, why particular requests are outliers in terms of latency, and how a single user request flows through backend services - many complex web services
may involve dozens of backend services or datastores. It may be possible to answer some of these questions using logs analysis.
However, there is a better solution, designed just for this problem: distributed tracing.

Distributed tracing has two key concepts: traces and spans. A trace represents a whole request or transaction. Traces are uniquely identified by trace IDs.
Traces are made up of a set of spans, each tagged with the trace ID of the trace it belongs to. Each span is a unit of work: a remote procedure call or
web request to a specific service, a method execution, or perhaps the time that a message spends in a queue. Spans can have child spans.
There are specific tools that are designed to collect and store distributed traces, and to perform useful queries against them.

One of the key aspects of distributed tracing is that when services call other services the trace ID is propagated to those calls (in HTTP-based systems
this is done using a special HTTP [traceparent header](https://uptrace.dev/opentelemetry/opentelemetry-traceparent.html)) so that the
overall trace may be assembled. This is necessary because each service in a complex chain of calls independently posts its spans to the distributed
trace collector. The collector uses the trace ID to assemble the spans together, like a jigsaw puzzle, so that we can see a holistic view of an
entire operation.

[OpenTelemetry](https://opentelemetry.io/) (also known as OTel) is the main industry standard for distributed tracing. It governs the format of traces and spans, and how traces and spans are collected. It is worth spending some time exploring the [OTel documentation](https://opentelemetry.io/docs/), particularly the Concepts
section. The [awesome-opentelemetry repo](https://github.com/magsther/awesome-opentelemetry) is another very comprehensive set of
resources.

Distributed tracing is a useful technique for understanding how complex systems such as RAFT are operating. The goal of this project is to use distributed
tracing to trace and to visualise operations in an implementation of RAFT, which is a nontrivial distributed system.

### Part 1: Get RAFT working

To begin with, we will need a running implementation of RAFT. Eli Bendersky has written a
[detailed set of blog posts describing a RAFT implementation in Go](https://eli.thegreenplace.net/2020/implementing-raft-part-0-introduction/).
Read these blog posts carefully.

You can either

1.  try to write your own RAFT implementation, building up the functionality in the stages described by Bendersky, or
2.  use Bendersky's code, after having thoroughly read it and understood it
3.  write your own implementation, using Bendersky's as a reference if you get stuck

**Note:** Note that Bendersky's code as it stands does not include a main package for actually running standalone RAFT servers.
It instead includes a test harness that simulates running a cluster in one process.
We want to run a real cluster sending real RPCs between the members, so that we can see traces composed of spans from different instances.
In order to do this we will need to make some straightforward changes to Bendersky's codebase:

- add a `main` package (with code similar to the test harness setup code)
- pass in the list of cluster members (i.e. a list of host:port pairs) as an argument to the program (in a real-world application we would likely use some form of service discovery)
- change the code in `server.go` to send and receive gRPC calls (in a similar fashion to the [gRPC Client-Server project](https://github.com/CodeYourFuture/immersive-go-course/tree/main/grpc-client-server))

Reading code written by others is a useful skill to have, so if you opt to create your own implementation, you should still review Bendersky's code.
Does it differ from yours in any significant respect?

By the end of Part 1 we should have a running RAFT cluster with 5 instances.
We may choose to run our RAFT cluster locally using `docker-compose`, `minikube`, or any other appropriate tool.

If you are using Bendersky's code, you will notice that his RAFT implementation implements a method `func (cm *ConsensusModule) Submit(command interface{}) bool`.
This simply appends the `command` provided to a log.

**Note:** What Bendersky is doing here is rather sketchily demonstrating a theoretical [Finite State Machine](https://en.wikipedia.org/wiki/Finite-state_machine) (FSM).
FSMs are a computer science concept (that you don't need to know any details about for this project): the idea that you can implement a program's core state as an abstract machine, with a specific set of states
and a specific set of transitions between those states. Externally-provided commands are the trigger to move between states.

Try modifying your code to instead implement Get, Set, and CompareAndSet commands - a highly-consistent key-value store - and make these available via gRPC methods.
For instance, a Set command may have two fields, a key and a value, and a Get command may have one field, a key. If your RAFT cluster accepts a Set command setting X = 10, then a Set command setting X = 20, and then gets a Get command for X, it should return 20.
There is already a `storage` module that you can use.

**Note:** [CompareAndSwap](https://en.wikipedia.org/wiki/Compare-and-swap) (also called CompareAndSet, or CAS)
is a very useful pattern for concurrent systems that lets you update a key to a given value only if that key already has a specific value.
This is useful for implementing sequences of operations where we read a value, perform some computation that modifies that value, and then write that value back -
but without potentially overwriting any changes to that value that other processes might have performed.

Next, write a client that uses your RAFT cluster to perform Sets, Gets, and CompareAndSets.
Bendersky's code also doesn't do anything if you send a `Submit` to any server other than the leader. It may be useful to have your program return the
address (host:port) of the leader instead, as part of your gRPC reply, along with an indication that the operation was not attempted. Your client can then
retry the operation against the leader.

Please timebox this part of the project to no more than two days, in order to leave time for the other sections.
At the start of the third day, if your implementation is not complete, begin modifying Bendersky's. You can come back and complete your
own implementation if you have time at the end of the sprint.

### Part 2: Add distributed tracing

In this section we will add distributed tracing support to the RAFT implementation from Part 1.
We will use [Honeycomb](https://www.honeycomb.io/), a SaaS distributed tracing provider.

Honeycomb provides a useful guide for their own
[OpenTelemetry Distribution for Go](https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/).

Add tracing to all parts of your application which you might wish to trace. Consider what operations are
interesting - in general, anything that may feasibly take a long time, such as a RPC, writing to storage, or
taking a lock may be a candidate for a span.

Run your system and view traces in Honeycomb. Run through the [Honeycomb sandbox tour](https://play.honeycomb.io/sandbox/environments/analyze-debug-tour)
and then explore your own data in the same way.

Do distributed operations (such as leader elections) surface in Honeycomb as coherent traces?

Create a [Board](https://docs.honeycomb.io/working-with-your-data/boards/) in honeycomb with some useful visualisations.

### Part 3: Debugging latency and failures using distributed tracing

Bendersky's implementation of RAFT allows you to simulate unreliable RPCs. Enable this (if you wrote your own implementation then you will
need to add this kind of chaos capability as a feature first - for example, adding an environment variable which, if set, will drop some percent of requests completely, and add random latency to others).
Now, use Honeycomb to observe the dropped RPCs and delays that
the simulation injects. Did these show up on your board from Part 2?

Can you add further kinds of chaos? What about storage failures?

### Part 4: Comparison with logging and metrics

Consider if you didn't have distributed tracing in this project, what logging you may add, and what metrics you may record.

What kind of analysis does each form of observability make easier? What's harder to do with each? We tend to do all three in real systems - when might each be useful?

### Part 5: Reducing costs of distributed tracing

In this exercise we are running a small system and using Honeycomb's free tier, so cost is not a consideration.
However, in real production systems, distributed tracing can create a large volume of traces and spans. This can be
costly in terms of network, storage, or SaaS bills.

For this reason, many distributed tracing users use sampling or ratelimiting to control the number of traces that are
collected. Read about [OTel Sampling and Ratelimiting](https://uptrace.dev/opentelemetry/sampling.html).

Modify your solution to support sampling a specific percentage of requests, and to limit the total number of traces sent to no more than 20 per minute per cluster member.
(Hint: You typically want to get all or no spans for a whole trace, rather than dropping spans independently of the trace they're in).

Consider also whether some requests may be more important to trace than others. What may make a request more or less interesting than others?
Note that certain kinds of sampling strategies are not possible to implement at the client. For instance, it is not possible to sample only failed
requests without first collecting all the spans: because we only know whether the request succeeds at the end of the operation.
You can read about some of these concerns in [this article about head-based and tail-based sampling](https://uptrace.dev/opentelemetry/sampling.html).
