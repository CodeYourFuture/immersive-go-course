<!--forhugo
+++
title="4. Asynchronous Work and Pipelines" 
+++
forhugo-->

# 4

## Asynchronous Work and Pipelines {#asynchronous-work-and-pipelines}

Not all work that we want to do with computers involves serving a request in near-real-time and responding to a user. Sometimes we need to do _asynchronous_ tasks like:

- Periodic work, such as a nightly data export, or computing monthly reports
- Work scheduled for later, such as scheduling reminders to users
- Long-running work, such as scheduling a build or a set of tests
- Running a continuous statistics computation based on incoming data

Batch processes may not be that large and may just run as a scheduled [cron](https://en.wikipedia.org/wiki/Cron) job.

### MapReduce {#mapreduce}

However, not all computations fit on one machine. One way of running large batch computations across a fleet of computers is the MapReduce paradigm. Read this article which [describes how MapReduce works](https://medium.com/edureka/mapreduce-tutorial-3d9535ddbe7c).

- How does MapReduce help us to scale big computations?

Read this book chapter about [Data Processing Pipelines](https://sre.google/sre-book/data-processing-pipelines/).

- Give two reasons why data processing pipelines can be fragile

Optionally, you can follow this short tutorial to implement a distributed word count application, and run it locally on a Glow cluster. You will get hands-on experience with MapReduce.

- [https://blog.gopheracademy.com/advent-2015/glow-map-reduce-for-golang/](https://blog.gopheracademy.com/advent-2015/glow-map-reduce-for-golang/)

### Queues {#queues}

Queues are a frequently-seen component of large software systems that involve potentially heavyweight or long-running requests. A queue can act as a form of buffer, smoothing out spikes of load so that the system can deal with work when it has the resources to do so. Read about the [Queue-Based Load-Leveling Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/queue-based-load-leveling).

- How can results of tasks be communicated back to users in a queue-based system?

Kafka is a commonly-used open-source distributed queue. Read [Apache Kafka in a Nutshell](https://medium.com/swlh/apache-kafka-in-a-nutshell-5782b01d9ffb).

- What are the components of the Kafka architecture?
- How are topics different from partitions?

### Project work for this section {#project-work-for-this-section}

- [https://github.com/CodeYourFuture/immersive-go-course/tree/main/kafka-cron](https://github.com/CodeYourFuture/immersive-go-course/tree/main/kafka-cron)
