<!--forhugo
+++
title="Projects"
+++
forhugo-->

## Requirements

Before you start this course, make sure you've read and followed all of the instructions in the [prep](../prep/README.md) document. This will get you set up and explain how to work through the projects.

Remember: you can _always_ Google or ask for help if you get stuck.

## Projects

This course is structured into self-contained projects that you can work through at your own pace.

Each project has its own directory with a README.md file that has instructions. If you want to take a look at one way of completing an exercise, there's some code waiting on a branch prefixed `impl/` (short for "implementation") and an associated [Pull Request](https://github.com/CodeYourFuture/immersive-go-course/pulls) for you to look at. Try not to copy!

Most exercises finish with a list of optional extension tasks. It's highly recommended that you try them out. Note that often the extensions are open-ended and under-specified - make sure to think about them with a curious mind: Why are they useful? What trade-offs do they have?

1. [Output and Error Handling](./output-and-error-handling)
   <br>An introduction to how to handle errors in Go, and how to present information to users of programs run on the command line.
1. [CLI & Files](./cli-files)
   <br>An introduction to building things with Go by replicating the unix tools `cat` and `ls`.
1. [File Parsing](./file-parsing)
   <br>Practice parsing different formats of files, both standard and custom.
1. [Servers & HTTP requests](./http-auth)
   <br>Learn about long-running processes, HTTP and `curl`
1. [Servers & Databases](./server-database)
   <br>Build a server that takes data from a database and serves it in `json` format
1. [Multiple Servers](./multiple-servers)
   <br>Build and run file & API servers behind nginx in a simple multi-server architecture
1. [Docker & Cloud Deployment](./docker-cloud/)
   <br>Use containers to reproducibly deploy applications into the cloud
1. [gRPC](./grpc-client-server)
   <br>Learn about RPCs and how they differ from REST, and start thinking about observability
1. [Batch Processing](./batch-processing/)
   <br>Build an image processing pipeline with cloud storage
1. [Buggy App](./buggy-app/)
   <br>Run, debug, and fix a buggy application
1. [Memcache](./memcached-clusters)
   <br>Explore sharding and replication of state
1. [Kafka Cron](./kafka-cron)
   <br>Build a distributed multi-server application handling variable load, with Kafka as a task queue
1. [RAFT and OTel](./raft-otel)
   <br>Build a complex distributed system for with strong consistency, and instrument it with tracing
