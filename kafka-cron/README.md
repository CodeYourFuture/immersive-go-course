# Distributed Cron Based on Kafka

In this project we're going to build a simple distributed `cron` system, based on the Apache Kafka distributed queue system.

> **Note:** This project requires us to have Docker with [Compose](https://docs.docker.com/compose/) installed. Compose is a tool for defining and running multi-container Docker applications. With Compose, we use a YAML file to configure the application’s services. Then, with a single command, we create and start all the services from the configuration.

Learning Objectives
 - How can we use a distributed queue in software architecture?
 - How can we deal with errors in a system based on distributed queues?
 - How can we instrument a complex application with metrics? How should we design alerting?

Timebox: 3 days

## Project

### Background on cron

We are going to implement a distributed version of the `cron` job scheduler (read about [cron](https://en.wikipedia.org/wiki/Cron) if you are not familiar with it). Cron jobs are defined by two attributes: the command to be executed, and either the schedule that the job should run on or a definition of the times that the job should execute. The schedule is defined according 
to the `crontab` format (you can find parsers for this format for Golang - the most widely used is [robfig/cron](https://github.com/robfig/cron)).

The `cron` tool common to Unix operating systems runs jobs on a schedule. Cron only works on a single hosts. We want to create a version of cron that can schedule jobs across multiple workers, running on different hosts.

### Background on Apache Kafka

Kafka is an open-source distributed queue. You can read about the core Kafka concepts in the [Kafka: a Distributed Messaging System for Log Processing paper](https://www.microsoft.com/en-us/research/wp-content/uploads/2017/09/Kafka.pdf).

After reading that paper you should understand:
 * How Kafka stores data
 * What producers, consumers, and brokers are
 * What topics and partitions are

### Part 1: Distributed cron with one queue

In this section of the project, we will start by creating a functional distributed cron system. We will build two separate programs:
 * A Kafka producer that reads configuration files for jobs and queues jobs for execution
 * A Kafka consumer that dequeues jobs from a queue and runs them

We'll need to be able to run Kafka. The easiest way is to use `docker-compose`. The [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose) project provides several starter configurations for running Kafka. The config for `zk-single-kafka-single.yml` will work for development purposes.

There is a [Golang Kafka client](https://docs.confluent.io/kafka-clients/go/current/overview.html#go-example-code) that we will use to interact with Kafka.

We may want to run other Docker containers later, so we may want to make our own copy of that configuration that we can add to. 

Our producer program needs to be able to do the following:
 * Read and parse a file with cron job definitions (we'll set up our own for this project, don't reuse the system cron config file because we will want to modify the format later)
 * Write a message to Kafka specifying the command to run, the intended start time of the job, and any other information that we think is necessary. It probably makes sense to encode this information as JSON (see [Go By Example: JSON](https://gobyexample.com/json) if you have never worked with JSON in Golang before)
 * We will also need to [create a Kafka topic](https://kafka.apache.org/documentation/#quickstart_createtopic). In a production environment we would probably use separate tooling to manage topics (perhaps Terraform), but for this project, we can [create our Kafka topic using code](https://github.com/confluentinc/examples/blob/7.3.0-post/clients/cloud/go/producer.go#L39).

Our consumer program needs to be able to do the following things:
 * Read job information from a Kafka queue (decoding JSON)
 * Execute the commands to run the jobs (assume this is a simple one-line command that you can `exec` for now)
 * Because the producer is writing jobs to the queue when they are ready to run, your consumer does not need to do any scheduling or to parse crontab format

 We want to run two consumers - therefore, when we create our topic, we should create two partitions of the topic. We will also need to specify a key for each Kafka message that we produce - Kafka assigns messages to partitions based on a hash of the message ID. We can use a package such as [google's UUID package](https://pkg.go.dev/github.com/google/UUID) to generate keys. 

We can build Docker containers for our producer and consumer and add these to our docker-compose configuration. We should create a Makefile or script to make this repeatable.

Test our implementation and observe both of our consumers running jobs scheduled by your producer. What happens if we only create one partition in our topic? What happens if we create three?

> **Note:** For the purposes of keeping this project scope tractable, we are ignoring two things. The first is security: simply run commands as the user that our consumer runs as. The second thing is that we are assuming the jobs to be run consist of commands available on the consumers. You may address these concerns later in an optional extension of the project if you have time.
 
### Part 2: Distributed cron with multiple queues

A new requirement: our distributed cron system needs to be able to schedule jobs to run in multiple clusters. Imagine that we want to support users who have data stored in specific clusters and they want to make sure their cron jobs are running near their data.

We don't really need to set up any clusters - just write our program as though you had multiple sets of consumer workers in different clusters.

 * Define a set of clusters in our program (two is fine, `cluster-a` and `cluster-b`) 
 * Each cluster should have its own Kafka topic
 * Update the job format so that jobs must specify what cluster to run in
 * Run separate consumers that are configured to read from each cluster-specific topic 
 
Test that our new program and Kafka configuration works as expected.

How would you do this sort of a migration in a running production environment, where you could not drop existing jobs?

### Part 3: Handling errors

What happens if there is a problem running a job? Maybe the right thing is retry it. 

This should be a configurable property of our cron jobs: update our program to add this to the job configurations and message format.

However: we don't want to risk retry jobs displacing first-time runs of other jobs. This is why some queue-based systems [use separate queues for retries](https://www.uber.com/en-IE/blog/reliable-reprocessing/).

We can create a second set of topics for jobs that fail the first time and need to be retried (we need one retry topic for each cluster). If a job fails, the consumer should write the job to the corresponding retry topic for the cluster (and decrement the remaining allowed attempts in the job definition). 

Run some instances of your consumer program that read from your retry queues (we can make this a command-line option in your consumer).

Define a job that fails and observe your retry consumers retrying and eventually discarding it.

Define a job that randomly fails some percent of the time, and observe your retry consumers retrying and eventually completing it.

### Part 4: Monitoring and Alerting

In software operations, we want to know what our software is doing and how it is performing.

One very useful technique is to have our program export metrics. Metrics are basically values that our program makes available (the industry standard is to export and scrape over HTTP).

Specialised programs, such as Prometheus, can then fetch metrics regularly from all the running instances of our program, store the history of these metrics, and do useful arithmetic on them (like computing rates, averages, and maximums). We can use this data to do troubleshooting and to alert if things go wrong.

Read the [Overview of Prometheus](https://prometheus.io/docs/introduction/overview/) if you are not familiar with Prometheus.

The [Prometheus Guide to Instrumenting a Go Application](https://prometheus.io/docs/guides/go-application/) describes how to add metrics to a Golang application.

First, consider:
 * What kinds of things may go wrong with our system? (it is useful to look at errors your code is handling, as inspiration)
 * What would users' expectations be of this system?
 * What metrics can we add that will tell us when the system is not working as intended?
 * What metrics can we add that might help us to troubleshoot the system and understand how it is operating? Read back through the first three parts of this exercise to try and identify the properties of the system that we might want to know about.

Asking these questions should guide us in designing the metrics that our consumers and producer should export.
Think about what kinds of problems can happen both in the infrastructure - Kafka, your consumers and producers - and in the submitted jobs. 

Add metrics to your programs. Verify that they work as expected using `curl` or your web browser.

Kafka doesn't export Prometheus metrics natively. However, we can use the official 
[Prometheus JMX exporter](https://github.com/prometheus/jmx_exporter) to expose its metrics.
Set this up (it is probably best as another container in your `docker-compose` configuration - we'll need to define a simple Dockerfile for it).

Next, we can add Prometheus, AlertManager, and Grafana, a common monitoring stack, to our `docker-compose` configuration. Here is an example configuration that we can adapt: https://dzlab.github.io/monitoring/2021/12/30/monitoring-stack-docker/. AlertManager is used for notifying operators of unexpected conditions, and Grafana is useful for building dashboards that allow us to troubleshoot and understand our system's operation.

If your computer is struggling to run such a complex `docker-compose` system in a performant fashion, you can cut down the number of Kafka topics and consumers that you are running to a minimum (just one producer and consumer/retry consumer pair are fine - don't run sets of these for multiple clusters if your computer is under too much load). 

We'll need to set up a Prometheus configuration to scrape our producers and consumers. Prometheus [configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) is quite complex but we can adapt this [example configuration](https://github.com/prometheus/prometheus/blob/main/documentation/examples/prometheus.yml).

Check that Prometheus is correctly scraping our metrics via http://localhost:9090/metrics.

Next, write an [AlertManager configuration](https://prometheus.io/docs/alerting/latest/alertmanager/) and set up at least one alert. Set up PagerDuty (it has a free trial period available) and get your system to fire an alert, and notify us via PagerDuty.

We can also build a Grafana dashboard to display our Prometheus metrics. The [Grafana Fundamentals](https://grafana.com/tutorials/grafana-fundamentals/) tutorial will walk you through how to do this (although we will need to use our own application and not their sample application).

## Extensions

### Kafka Chaos

Try running multiple Kafka brokers and Zookeeper servers with our producers and consumers (using another of the [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose)) configurations. Experiment with downing Kafka and Zookeeper containers.

How many containers being down can our system tolerate?

What happens to the Kafka system logs and the metrics that our binaries export? Did our alerts fire?

### Dealing with long-running jobs and load (challenging)

What does our system do if someone submits a very long-running job? Try testing this with the `sleep` command. 

If this is an issue for the stable operation of our system, or for running jobs in a timely fashion, what can we do about this?

If your system had problems, did our alerts fire?

How can we prevent our consumers getting overloaded if compute-intensive jobs are submitted?

### Security using Firecracker VMs (challenging)

In an earlier note it was mentioned that there are security issues with simply `exec`-ing code in this way. 

A better solution would be to use a [Firecracker VM](https://github.com/firecracker-microvm/firecracker/) to run the cron commands. Firecracker is an open-source virtualization technology that lets us start lightweight virtual machines very quickly and cheaply. It was developed at AWS to support services like AWS Lambda.

Here are some demos and examples of projects built with Firecracker:
 * https://stanislas.blog/2021/08/firecracker/
 * https://jvns.ca/blog/2021/01/23/firecracker--start-a-vm-in-less-than-a-second/

There is a [Firecracker SDK for Golang](https://github.com/firecracker-microvm/firecracker-go-sdk). If you have a significant amount of extra time available, updating the system to run commands in Firecracker VMs instead of exec-ing the commands provided would be a very good challenge.