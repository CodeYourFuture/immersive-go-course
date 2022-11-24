# Distributed Cron Based on Kafka

In this project you're going to build a simple distributed `cron` system, based on the 
Apache Kafka distributed queue system. 

> **Note:** This project requires you to have Docker with [Compose](https://docs.docker.com/compose/) installed. Compose is a tool for defining and running multi-container Docker applications. With Compose, we use a YAML file to configure the application’s services. Then, with a single command, we create and start all the services from the configuration.

Learning Objectives
 - How can we use a distributed queue in software architecture?
 - How do distributed queues scale?
 -  How can we deal with errors in a system based on distributed queues?
 - How can we instrument a complex application with metrics? How should we design alerting?

Timebox: 3 days

## Project

### Background on cron

We are going to implement a distributed version of the `cron` job scheduler (read about [cron](https://en.wikipedia.org/wiki/Cron) if you are not familiar with it). Cron jobs are defined by two attributes: the command to be executed and the schedule that the schedule, or a definition of the times that the job should execute. The schedule is defined according 
to the `crontab` format (you can find parsers for this format for Golang - the most widely used is [robfig/cron](https://github.com/robfig/cron)).

The `cron` tool common to Unix operating systems runs jobs on a schedule. Cron only works on a single hosts. We want to create a version of cron that can schedule jobs across multiple workers, running on different hosts.

### Background on Apache Kafka

Kafka is an open-source distributed queue. You can read about the core Kafka concepts in the [Kafka: a Distributed Messaging System for Log Processing paper](https://www.microsoft.com/en-us/research/wp-content/uploads/2017/09/Kafka.pdf).

After reading that paper you should understand:
 * How does Kafka store data? 
 * What are producers, consumers, and brokers?
 * How do brokers manage caching?
 * What is a topic, partition?

### Part 1: Distributed cron with one queue

In this section of the project, we will start by creating a functional distributed cron system. You will build two separate programs:
 * A Kafka producer that reads configuration files for jobs and queues jobs for execution
 * A Kafka consumer that dequeues jobs from a queue and runs them

You'll need to be able to run Kafka. The easiest way is to use `docker-compose`. The [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose) project provides several starter configurations for running Kafka. The config for `zk-single-kafka-single.yml` will work for development purposes.

There is a [Golang Kafka client](https://docs.confluent.io/kafka-clients/go/current/overview.html#go-example-code) that you will use to interact with Kafka.

You may want to run other Docker containers later, so you may want to make your own copy of that configuration that you can add to. 

Your producer program needs to be able to do the following:
 * Read and parse a file with cron job definitions
 * Write a message to Kafka specifying the command to run, the intended start time of the job, and any other information that you think is necessary. It probably makes sense to encode this information as JSON (see [Go By Example: JSON](https://gobyexample.com/json) if you have never worked with JSON in Golang before)
 * You will also need to [create a Kafka topic](https://kafka.apache.org/documentation/#quickstart_createtopic). In a production environment we would probably use separate tooling to manage topics (perhaps Terraform), but for this project, you can [create your Kafka topic using code](https://github.com/confluentinc/examples/blob/7.3.0-post/clients/cloud/go/producer.go#L39).

Your consumer program needs to be able to do the following things:
 * Read job information from a Kafka queue (decoding JSON)
 * Execute the commands to run the jobs (assume this is a simple one-line command that you can `exec` for now)
 * Because the consumer is writing jobs to the queue when they are ready to run, your consumer does not need to do any scheduling or to parse crontab format

 We want to run two consumers - therefore, when you create your topic, you should create two partitions of the topic. 

You can build Docker containers for your producer and consumer and add these to your docker-compose configuration. You should create a Makefile or script to make this repeatable.

Test your implementation and observe both of your consumers running jobs scheduled by your producer. What happens if you only create one partition in your topic? What happens if you create three?

> **Note:** For the purposes of keeping this project scope tractable, we are ignoring two things. The first is security: simply run commands as the user that your consumer runs as. The second thing is that we are assuming the jobs to be run consist of commands available on the consumers. You may address these concerns later in an optional extension of the project if you have time.

### Part 2: Status

You should now have a working program that runs scheduled jobs across multiple consumers.
However, this system lacks a way for operators to understand what is happening to scheduled jobs.

You can fix this by:
 * Running a database - any database you are familiar with is fine. Add it to your docker-compose configuration
 * Your producer should write job execution information to this database, and the consumers should update the database when they start and finish the job. 
 * It will be easiest to generate a unique ID for each job and add this to your job definition
 * Write a CLI that can display details of running and recently completed jobs. Did jobs succeed? How long did they take? 

### Part 3: Distributed cron with multiple queues

A new requirement: your distributed cron system needs to be able to schedule jobs to run in multiple clusters. 
 * Define a set of clusters in your program (two is fine, `cluster-a` and `cluster-b`)
 * Create a topic for each cluster, as well as a catchall `any-cluster` queue
 * Update the job format so that jobs can specify what cluster to run in. If they do not specify, then put them in `any-cluster`
 * Run separate consumers that are configured to read from each cluster specific queue and the `any-cluster` queue
 
Test that your new program and Kafka configuration works as expected.

### Part 4: Handling errors

What happens if there is a problem running a job? Maybe the right thing is retry it. 
This should be a configurable property of your cron jobs: update your program to add this to the job configurations and message format.

However: you don't want to risk retry jobs displacing first-time runs of other jobs. This is why some queue-based systems [use separate queues for retries](https://www.uber.com/en-IE/blog/reliable-reprocessing/).

You can create a second set of topics for jobs that fail the first time and need to be retried (we need one for each cluster and for `any-cluster`). If a job fails, the consumer should write the job to the corresponding retry queue (and decrement the remaining allowed attempts in the job definition). If there are no more allowed attempts, then discard the job.

Run some instances of your consumer program that read from your retry queues. Define a job that fails and observe your retry consumers retrying and eventually discarding it.

But what about if a job cannot be parsed, or is otherwise invalid? In that case, it cannot be run. It should be written to special 'dead-letter queue' which is set aside for this purpose. This is useful for debugging systems. 

Implement a dead-letter queue, and introduce an invalid message into your system (you can write Golang code to do this or try one of the [Kafka command-line tools](https://medium.com/@TimvanBaarsen/apache-kafka-cli-commands-cheat-sheet-a6f06eac01b)). Observe your invalid message being written to the dead-letter queue.

### Part 5: Monitoring and Alerting

In software operations, we want to know what our software is doing and how it is performing.
One very useful technique is to have our program export metrics. Metrics are basically values that your 
program makes available (the industry standard is to export and scrape over HTTP). 

Specialised programs, such as Prometheus, can then fetch metrics regularly
from all the running instances of your program, store the history of these metrics, and do useful arithmetic on them
(like computing rates, averages, and maximums). We can use this data to do troubleshooting and to alert if things 
go wrong.

Read the [Overview of Prometheus](https://prometheus.io/docs/introduction/overview/) if you are not familiar with Prometheus.

The [Prometheus Guide to Instrumenting a Go Application](https://prometheus.io/docs/guides/go-application/) describes how to add metrics to a Golang application.

First, consider:
 * What kinds of things may go wrong with your system?
 * What metrics can we add that will tell us when the system is not working as intended?
 * What metrics can we add that might help us to troubleshoot the system and understand how it is operating?

Asking these questions should guide you in designing the metrics that your consumers and producer should export. Add these metrics to your programs. Verify that they work as expected using `curl` or your web browser.

Kafka doesn't export Prometheus metrics natively. However, you can use the official 
[Prometheus JMX exporter](https://github.com/prometheus/jmx_exporter) to expose its metrics.
Set this up (it is probably best as another container in your `docker-compose` configuration - you'll need to define a simple Dockerfile for it).

Next, you can add Prometheus, AlertManager, and Grafana, a common monitoring stack, to your `docker-compose` configuration. Here is an example configuration that you can adapt: https://dzlab.github.io/monitoring/2021/12/30/monitoring-stack-docker/. AlertManager is used for notifying operators of unexpected conditions, and Grafana is useful for building dashboards that allow you to troubleshoot and understand your systems operation.

If your system is struggling to run such a complex `docker-compose` system in a performant fashion, you can cut down the number of Kafka topics and consumers that you are running to a minimum (just one consumer and one retry consumer are fine - don't run sets of these for multiple clusters if your system is under too much load). 

You'll need to set up a Prometheus configuration to scrape your producers and consumers. Prometheus [configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) is quite complex but you can adapt this [example configuration](https://github.com/prometheus/prometheus/blob/main/documentation/examples/prometheus.yml).

Check that Prometheus is correctly scraping your metrics via http://localhost:9090/metrics.

Next, write an [AlertManager configuration] and set up at least one alert. Set up PagerDuty (it has a free trial period available) and get your system to fire an alert, and notify you via PagerDuty.

You can also build a Grafana dashboard to display your Prometheus metrics. 
The [Grafana Fundamentals](https://grafana.com/tutorials/grafana-fundamentals/) tutorial will walk you through how to do this (although you will need to use your own application and not their sample application).

### Part 6: (Optional) Kafka Chaos

Try running multiple Kafka brokers and Zookeeper servers (using another of the [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose)) configurations. Experiment with downing Kafka and Zookeeper containers.

How many containers being down can your system tolerate?
What happens to the system logs and your metrics? Did you get alerts?

### Part 7: (Optional) Expiring data

In Part 2, we added a datastore to track job execution. 

In a production system this might grow very large over time.
Do you need to keep job execution data indefinitely?
Some datastores allow you to set an expiry time for data as it is written, and automatically deal with expiry. 
In other systems, you may need to deal with deletion of older data in some other way - via a periodic deletion job that operates based on a per-row timestamp, or using [time-series tables](https://docs.aws.amazon.com/redshift/latest/dg/c_best-practices-time-series-tables.html). 

Implement a strategy for your project to manage periodic deletion of data in your datastore. Consider how to add monitoring and alerting to make sure it continues to work.

### Part 8: (Optional) Security using Firecracker VMs

In an earlier note it was mentioned that there are security issues with simply `exec`-ing code in this way. 

A better solution would be to use a [Firecracker VM](https://github.com/firecracker-microvm/firecracker/) to run the cron commands. Firecracker is an open-source virtualization technology that lets you start 
lightweight virtual machines very quickly and cheaply. It was developed at AWS to support 
services like AWS Lambda.  

Here are some demos and examples of projects built with Firecracker:
 * https://stanislas.blog/2021/08/firecracker/
 * https://jvns.ca/blog/2021/01/23/firecracker--start-a-vm-in-less-than-a-second/

There is a [Firecracker SDK for Golang](https://github.com/firecracker-microvm/firecracker-go-sdk). If you have a significant amount of extra time available, reimplementing the system to run cron jobs in Firecracker VMs would be a good challenge.