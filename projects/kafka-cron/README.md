<!--forhugo
+++
title="Distributed Cron Based on Kafka"
+++
forhugo-->

In this project we're going to build a simple distributed `cron` system, based on the Apache Kafka distributed queue system.

> **Note:** This project requires us to have Docker with [Compose](https://docs.docker.com/compose/) installed. Compose is a tool for defining and running multi-container Docker applications. With Compose, we use a YAML file to configure the application’s services. Then, with a single command, we create and start all the services from the configuration.

## Learning Objectives

- Use a distributed queue in software architecture
- Deal with errors in a system based on distributed queues
- Instrument a complex application with metrics
- Design alerting in a complex application

Timebox: 5 days

## Project

### Background on cron

We are going to implement a distributed version of the `cron` job scheduler (read about [cron](https://en.wikipedia.org/wiki/Cron) if you are not familiar with it). Cron jobs are defined by two attributes: the command to be executed, and either the schedule that the job should run on or a definition of the times that the job should execute. The schedule is defined according
to the `crontab` format (you can find parsers for this format for Golang - the most widely used is [robfig/cron](https://github.com/robfig/cron)).

The `cron` tool common to Unix operating systems runs jobs on a schedule. Cron only works on a single hosts. We want to create a version of cron that can schedule jobs across multiple workers, running on different hosts.

### Background on Apache Kafka

Kafka is an open-source distributed queue. You can read about the core Kafka concepts in the [Kafka: a Distributed Messaging System for Log Processing paper](https://www.microsoft.com/en-us/research/wp-content/uploads/2017/09/Kafka.pdf).

After reading that paper you should understand:

- How Kafka stores data
- What producers, consumers, and brokers are
- What topics and partitions are

### Part 1: Distributed cron with one queue

In this section of the project, we will start by creating a functional distributed cron system. We will build two separate programs:

- A Kafka producer that reads configuration files for jobs and queues tasks for execution
- A Kafka consumer that dequeues jobs from a queue and runs them

Kafka itself is just a queue that lets you communicate in a structured and asynchronous way between producers and consumers. Therefore, all the scheduling logic for managing recurring jobs must be part of your producer (although it is recommended to reuse a suitable library to assist with parsing crontabs and scheduling). Every time a job is due to be run, your producer creates a new message and writes it to Kafka, for a consumer to dequeue and run.

We'll need to be able to run Kafka. The easiest way is to use `docker-compose`. The [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose) project provides several starter configurations for running Kafka. The config for `zk-single-kafka-single.yml` will work for development purposes.

There is a [Golang Kafka client](https://docs.confluent.io/kafka-clients/go/current/overview.html#go-example-code) that we will use to interact with Kafka.

We may want to run other Docker containers later, so we may want to make our own copy of that configuration that we can add to.

Our producer program needs to be able to do the following:

- Read and parse a file with cron job definitions (we'll set up our own for this project, don't reuse the system cron config file because we will want to modify the format later)
- Write a message to Kafka specifying the command to run, the intended start time of the job, and any other information that we think is necessary. It probably makes sense to encode this information as JSON (see [Go By Example: JSON](https://gobyexample.com/json) if you have never worked with JSON in Golang before)
- We will also need to [create a Kafka topic](https://kafka.apache.org/documentation/#quickstart_createtopic). In a production environment we would probably use separate tooling to manage topics (perhaps Terraform), but for this project, we can [create our Kafka topic using code](https://github.com/confluentinc/examples/blob/7.3.0-post/clients/cloud/go/producer.go#L39).

Our consumer program needs to be able to do the following things:

- Read job information from a Kafka queue (decoding JSON)
- Execute the commands to run the jobs (assume this is a simple one-line command that you can `exec` for now)
- Because the producer is writing jobs to the queue when they are ready to run, your consumer does not need to do any scheduling or to parse crontab format

We want to run two consumers - therefore, when we create our topic, we should create two partitions of the topic. We will also need to specify a key for each Kafka message that we produce - Kafka assigns messages to partitions based on a hash of the message ID. We can use a package such as [google's UUID package](https://pkg.go.dev/github.com/google/UUID) to generate keys.

We can build Docker containers for our producer and consumer and add these to our docker-compose configuration. We should create a Makefile or script to make this repeatable.

Test our implementation and observe both of our consumers running jobs scheduled by your producer. What happens if we only create one partition in our topic? What happens if we create three?

> **Note:** For the purposes of keeping this project scope tractable, we are ignoring two things. The first is security: simply run commands as the user that our consumer runs as. The second thing is that we are assuming the jobs to be run consist of commands available on the consumers. You may address these concerns later in an optional extension of the project if you have time.

### Part 2: Distributed cron with multiple queues

A new requirement: our distributed cron system needs to be able to schedule jobs to run in multiple clusters. Imagine that we want to support users who have data stored in specific cells/AZs and they want to make sure their cron jobs are running near their data.

We don't really need to set up any cells for this - just write our program as though you had multiple sets of consumer workers. 
You *don't* need to set up multiple Kafka clusters for this - this extension is just about having multiple sets of consumer jobs, which we notionally call clusters.

- Define a set of clusters in our program (two is fine, `cluster-a` and `cluster-b`)
- Each cluster should have its own Kafka topic 
- Update the job format so that jobs must specify what cluster to run in
- Run separate consumers that are configured to read from each cluster-specific topic

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

- What kinds of things may go wrong with our system? (it is useful to look at errors your code is handling, as inspiration)
- What would users' expectations be of this system?
- What metrics can we add that will tell us when the system is not working as intended?
- What metrics can we add that might help us to troubleshoot the system and understand how it is operating? Read back through the first three parts of this exercise to try and identify the properties of the system that we might want to know about.

Asking these questions should guide us in designing the metrics that our consumers and producer should export.
Think about what kinds of problems can happen both in the infrastructure - Kafka, your consumers and producers - and in the submitted jobs.

Add metrics to your programs. Verify that they work as expected using `curl` or your web browser.

#### Running the Prometheus JMX Exporter to get Kafka metrics

Kafka doesn't export Prometheus metrics natively. However, we can use the official
[Prometheus JMX exporter](https://github.com/prometheus/jmx_exporter) to expose its metrics.

> **Note:** Kafka is a Java program. We don't need to know much about Java programs in order to run them, but it's useful to know that Java programs run in a host process called a Java Virtual Machine (JVM). The JVM also allows for injecting extra code called Java agents, which can modify how a program is run.

The Prometheus JMX exporter can run as a Java agent (alongside a Java program such as Kafka) or else as a standalone HTTP server, which collects metrics from a JVM running elsewhere and re-exports them as Prometheus metrics. If you're using [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose) as suggested above then your image contains the `jmx_prometheus_javaagent` already.

You need to create a `config.yaml`. A config file that will collect all metrics is:

```
rules:
- pattern: ".*"
```

Now, update the Kafka service in your `docker-compose.yml`. Add a volume - for example:

```
    volumes:
      - ./kafka-jmx-config.yaml:/kafka-jmx-config.yaml
```

Finally, you need to add a new line in your `environment` section for your Kafka server in your `docker-compose.yml`:

```
KAFKA_OPTS: -javaagent:/usr/share/java/cp-base-new/jmx_prometheus_javaagent-0.14.0.jar=8999:/kafka-jmx-config.yaml
```

The version of the `jmx_prometheus_javaagent` jar might change in a later version of the `cp-kafka` image, so if you have any issues running the software, this would be the first thing to check. You can't just map a newer version of the agent as a volume as this is likely to cause runtime errors due to multiple version of the agent on the Java classpath.

Now you should be able to see JVM and Kafka metrics on http://localhost:8999. Check this using `curl` or your web browser.

#### Running Prometheus, Alertmanager, and Grafana

Next, we can add Prometheus, AlertManager, and Grafana, a common monitoring stack, to our `docker-compose` configuration. Here is an example configuration that we can adapt: https://dzlab.github.io/monitoring/2021/12/30/monitoring-stack-docker/. AlertManager is used for notifying operators of unexpected conditions, and Grafana is useful for building dashboards that allow us to troubleshoot and understand our system's operation.

If your computer is struggling to run such a complex `docker-compose` system in a performant fashion, you can cut down the number of Kafka topics and consumers that you are running to a minimum (just one producer and consumer/retry consumer pair are fine - don't run sets of these for multiple clusters if your computer is under too much load).

We'll need to set up a Prometheus configuration to scrape our producers and consumers. Prometheus [configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) is quite complex but we can adapt this [example configuration](https://github.com/prometheus/prometheus/blob/main/documentation/examples/prometheus.yml).

For example, to scrape your Kafka metrics, you can add ths to the Prometheus configuration:

```
scrape_configs:
  - job_name: "kafka"
    static_configs:
      - targets: ["kafka1:8999"]
```

Once you have adapted the sample Prometheus configuration to scrape metrics from your running producer and consumer(s) and from the JMX exporter that is exporting the Kafka metrics, you should check that Prometheus is correctly scraping all those metrics. If you haven't changed the default port, you can access Prometheus's status page at http://localhost:9090/.

You can now try out some queries in the Prometheus UI.

For example, let's say that our consumers are exporting a metric `job_runtime` that describes how long it takes to run jobs. And let's say the metric is labelled with the name of the queue the consumer is reading from.

Because this metric is describing a population of observed latencies, the best metric type to use is a [histogram](https://prometheus.io/docs/practices/histograms/).

We can query this as follows:

```
histogram_quantile(0.9, sum by (queue, le)(rate(job_runtime[10m])))
```

This will give you the 90th percentile job runtime (i.e. the runtime where 90% of jobs complete this fast or faster) over the past 10 minutes (the `rate` function does this for histogram queries - it's a little counterintuitive).

For some more PromQL examples, see the [Prometheus Query Examples page](https://prometheus.io/docs/prometheus/latest/querying/examples/).

#### Alertmanager

Next, write an [AlertManager configuration](https://prometheus.io/docs/alerting/latest/alertmanager/) and set up at least one alert.

For instance:

- We could alert on the age of jobs being unqueued - if this rises too high (more than a few seconds) then users' jobs aren't being executed in a timely fashion. We should use a percentile for this calculation.
- We could also alert on failure to queue jobs, and failure to read from the queue.
- We expect to see fetch requests against all of our topics. If we don't, it may mean that our consumers are not running, or are otherwise broken. We could set up alerts on the `kafka_server_BrokerTopicMetrics_Count{name="TotalFetchRequestsPerSec"}` metric to check this.

For critical alerts in a production environment we would usually use PagerDuty or a similar tool, but for our purposes the easiest way to configure an alert is to use email.
This article describes how to send [Alertmanager email using GMail](https://www.robustperception.io/sending-email-with-the-alertmanager-via-gmail/) as an email server.

> **Note:** If you do this, be careful not to check your `GMAIL_AUTH_TOKEN` into GitHub - we should never check ANY token into source control. Instead, we can check in a template file and use a tool such as [heredoc](https://tldp.org/LDP/abs/html/here-docs.html) to substitute the value of an environment variable (our token) into the final generated Alertmanager configuration (and include this step in a build script/Makefile). It is also advisable a throwaway GMail account for this purpose, for additional security - just in case.

We can also build a Grafana dashboard to display our Prometheus metrics. The [Grafana Fundamentals](https://grafana.com/tutorials/grafana-fundamentals/) tutorial will walk you through how to do this (although we will need to use our own application and not their sample application).

## Extensions

#### Comprehensive Alerting Design and Runbooks

You should have at least one alert defined. However, for a production system, we need a comprehensive set of alerts that we can rely on to
tell us when our system is not meeting user expectations. Try to implement the smallest set of alerts that covers all cases.
Use [symptom-based alerting](https://docs.google.com/document/d/199PqyG3UsyXlwieHaqbGiWVa8eMWi8zzAn0YfcApr8Q/preview) and avoid cause-based alerting.
Write a short README about how you designed your alerts.

Now, for each alert, write a playbook that describes how to handle that type of alert. Information to include:

- A summary of the relevant system architecture (you can include a diagram, either as an image or using [mermaid.js](https://github.blog/2022-02-14-include-diagrams-markdown-files-mermaid/)).
- What the likely user impact is of the alert (e.g. "all scheduled tasks will fail" or "tasks will be slow to execute").
- What kinds of things might cause this alert?
- How would the engineer receiving that alert narrow down the possible causes and troubleshoot?
- How should the engineer address each possible cause that you can foresee?

A useful way to proceed is to think about all of the entities in your system: services, topics, and so on.
What would happen if each of these disappeared?
Now consider all the places where communication occurs in your system. What would happen if each of these communication paths failed, or if a software bug caused wrong messages to be sent? Don't forget that your monitoring system itself is a communication link to your production systems.

Your alert definitions should include a link to your playbooks on GitHub.

### Kafka Chaos

Try running multiple Kafka brokers and Zookeeper servers with our producers and consumers (using another of the [conduktor/kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose)) configurations. Experiment with downing Kafka and Zookeeper containers.

How many containers being down can our system tolerate?

What happens to the Kafka system logs and the metrics that our binaries export? Did our alerts fire? If not, consider how they could be improved - remember, the point of them is to tell us when something's wrong!

### Porting your system from docker-compose to minikube

In this project, and in several previous projects, you have used `docker-compose` to deploy our code, alongside dependencies (such as Kafka and
Zookeeper here). `docker-compose` is an extremely convenient tool for running a multi-part software stack locally (it also works well for running 
integration tests in a Continuous Integration workflow as part of your development process). However, in most deployments, we want to be able
to run our code and its dependencies across more than one host, in order to scale horizontally and to be robust to single-node failures. For this,
`kubernetes` is a better tool. `Kubernetes`, like `docker-compose`, is a platform for running containerised applications, but where `docker-compose` 
is focused on running a set of related containers on a single host, `kubernetes` is optimized for running services across many hosts.

Here are some introductions to Kubernetes: 
 * [What is Kubernetes](https://www.digitalocean.com/community/tutorials/an-introduction-to-kubernetes)
 * [Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/) - explore the linked pages here

You may have already used `minikube` in one of the previous projects. `Minikube` is a local version of Kubernetes, which we can use to learn 
(rather than needing to incur the expense of cloud services such as [EKS](https://aws.amazon.com/eks/) for learning on.)

Get your local `minikube setup` working by following [minikube start](https://minikube.sigs.k8s.io/docs/start/).
Follow the steps to install the program and run the `hello-minikube` application.

Once you have done this, you will need to convert your `docker-compose.yml` files to `kubernetes` deployment files.
There is a tool, [kompose](https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/) which may assist you.

> **Note:** `kompose` *will not* give you perfect results, however. In particular, `kompose` will not correctly convert cases where you are using a `docker-compose` `volume` to map a configuration file into your running container. `kompose` will just create a `persistentVolumeClaim` with empty data. The best thing to do is to remove those and use a [ConfigMap] for the config file content and map that in as a volume. You'll have to do this by hand. Note that there are some `persistentVolumeClaims` for application data that are appropriate - don't remove these, only the ones that are substituted for config. 

Test that your system works as it did on `docker-compose`.

Learn your way around the `kubernetes` command-line tool, `kubectl` (see the [kubectl documentation](https://kubernetes.io/docs/reference/kubectl/)).
* How can you resize your service, i.e. change the number of running instances (pods)?
* How can you remove one instance of your service (a single pod)?
* How can you see the logs for your service?
* How can you see the log of Kubernetes' management operations?

### Dealing with long-running jobs and load (challenging)

What does our system do if someone submits a very long-running job? Try testing this with the `sleep` command.

If this is an issue for the stable operation of our system, or for running jobs in a timely fashion, what can we do about this?

If your system had problems, did our alerts fire?

How can we prevent our consumers getting overloaded if compute-intensive jobs are submitted?
