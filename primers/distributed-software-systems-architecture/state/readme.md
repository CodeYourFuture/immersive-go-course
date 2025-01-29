<!--forhugo
+++
title="2. State"
+++
forhugo-->

# 2

## State {#section-2-state}

**State*ful* and state*less***

Components in distributed systems are often divided into stateful and stateless services. Stateless services don’t store any state between serving requests. Stateful services, such as databases and caches, do store state between requests. State _stays_.

When we need to scale up a stateless service we simply run more instances and loadbalance between them. Scaling up a stateful service is different: we need to deal with splitting up or replicating our data, and keeping things in sync.

### Caching {#caching}

Caches are an example of a stateful service. A [cache service](https://aws.amazon.com/caching/) is a high-performance storage service that stores only a subset of data, generally the data most recently accessed by your service. A cache can generally serve data faster than the primary data store (for example, a database), because a cache stores a small set of data in RAM whereas a database stores all of your data on disk, which is slower.

However, caches are not _durable_ - when an instance of your cache service restarts, it does not hold any data. Until the cache has filled with a useful _working set_ of data, all requests will be _cache misses_, meaning that they need to be served from the primary datastore. The _hit rate_ is the percentage of cacheable data that can be served from the cache. Hit rate is an important aspect of cache performance which should be monitored.

There are different strategies for getting data into your cache. The most common is _lazy loading_: your application always tries to load data from the cache first when it needs it. If it isn’t in the cache, the data is fetched from the primary data store and copied into the cache. You can also use a _write through_ strategy: every time your application writes data, it writes it to the cache at the same time as writing it to the datastore. However, when using this strategy, you have to deal with cases where the data isn’t in the cache (for instance, via lazy loading). Read about the [pros and cons](https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/Strategies.html#Strategies.WriteThrough) of these cache-loading strategies.

#### Why a cache service? {#why-a-cache-service}

Cache services, such as [memcached](https://www.memcached.org/) or [redis](https://redis.io/), are very often used in web application architectures to improve read performance, as well as to increase system throughput (the number of requests that can be served given a particular hardware configuration). You can also cache data in your application layer – but this adds complexity, because in order to be effective, requests affecting the same set of data must be routed to the same instance of your application each time. This means that your loadbalancer has to support the use of _[sticky sessions](https://www.linode.com/docs/guides/configuring-load-balancer-sticky-session/)._

In-application caches are less effective when your application is restarted frequently, because restarting your application means that all cached data is lost and the cache will be _cold_: it does not have a useful working set of data\_. \_In many organisations, web applications get restarted very frequently. There are two main reasons for this.

1. Deployment of new code. Many organisations using modern [CI/CD (Continuous Integration/Continuous Delivery)](https://en.wikipedia.org/wiki/CI/CD) deploy new code on every change, or if not on every change, many times a day.
2. The use of _[autoscaling](https://en.wikipedia.org/wiki/Autoscaling)_. Autoscaling is the use of automation to adjust the number of running instances in a software service. It is very common to use autoscaling with stateless services, particularly when running on a cloud provider. Autoscaling is an easy way to scale your service up to serve a peak load, while not paying for unused resources at other times. However, autoscaling also means that the lifespan of instances of your service will be reduced.

The use of a separate cache service means that the stateless web application layer can be deployed frequently and can autoscale according to workload while still using a cache effectively. Cache services typically are not deployed frequently and are less likely to use autoscaling.

#### Hazards of using cache services {#hazards-of-using-cache-services}

Operations involving cache services must be done carefully. Any operation that causes all your caches to restart will leave your application with an entirely _cold_ cache - a cache with no data in it. All reads that your cache would normally have served will go to your primary datastore. There are multiple possible outcomes in this scenario.

In the first case, your primary datastore can handle the increased read load and the only consequence for your application will be an increase in duration of requests served. Your application will also consume additional system resources: it will hold more open transactions to your primary datastore than usual, and it will have more requests in flight than usual. Over time your cache service fills with a useful working set of data and everything returns to normal. Some systems offer a feature called _[cache warmup](https://github.com/facebook/mcrouter/wiki/Cold-cache-warm-up-setup)_ to speed this process up.

In the second case, even though your datastore can handle the load, your application cannot handle the increase in in-flight requests. It may run out of resources such as worker threads, datastore connections, CPU, or RAM - this depends on the specific web application framework you are using. Applications in this state will serve errors or fail to respond to requests in a timely fashion.

In the third case, your primary datastore is unable to handle the entire load, and some requests to your datastore begin to fail or to time out. As discussed in the previous section, you should be setting a deadline for your datastore requests. In the absence of deadlines, your application will wait indefinitely for your overloaded datastore, and is likely to run out of system resources and then grind to a halt, until it is restarted. Setting deadlines allows your application to serve what requests it can without grinding to a halt. After some time your cache will fill and your application should return to normal.

Read about an incident at Slack involving failure of a caching layer: [Slack’s Incident on 2-22-22](https://slack.engineering/slacks-incident-on-2-22-22/).

#### Cache Invalidation {#cache-invalidation}

There are two hard things in computer science: cache invalidation, naming things, and off-by-one errors.

Cache invalidation means removing or updating an entry in a cache. The reason that cache invalidation is hard is that caches are optimised towards fast reads, without synchronising with the primary datastore. If a cached item is updated then the application must either tolerate stale cached data or update all cache replicas that reference the updated item.

Caches generally support specifying a Time-To-Live (TTL). After the TTL passes, or expires, the item is removed from the cache and must be fetched again from the main datastore. This lets you put an upper limit on how stale cached data may be. It is useful to add some ‘jitter’ when specifying TTLs. This means varying the TTL duration slightly - for example, instead of always using a 60 second TTL, you might randomly choose a duration in the range 54 to 66 seconds, up to ten percent higher or lower. This reduces the likelihood of load spikes on backend systems as a result of coordinated waves of cache evictions.

##### Immutable data {#immutable-data}

One of the simplest strategies to manage cache invalidation problems is to avoid updating data where possible. Instead, we can create new versions of the data. For example, if we are building an application that has profile pictures for users, if the user updates their profile picture we create a new profile picture with a new unique ID, and update the user data to refer to the new profile picture rather than the old one. Doing this means that there is no need to invalidate the old cached profile picture - we just stop referring to it. Eventually, the old picture will be removed from the cache, as the cache removes entries that have not been accessed recently.

In this strategy the profile picture data is immutable - meaning that it never changes. However, the user data does change, meaning that cached copies must be invalidated, or must have a short TTL so that users will see new profile pictures in a reasonable period of time.

Read more about [Ways to Maintain Cache Consistency](https://redis.com/blog/three-ways-to-maintain-cache-consistency/) in this article.

#### Scaling Caches {#scaling-caches}

##### Replicated Caches {#replicated-caches}

In very large systems, it may not be possible to serve all cache reads from one instance of a cache. Caches can run out of memory, connections, network capacity, or CPU. And with only one cache instance, if we lose or even update the instance, we will lose all of our cache capacity if we lose that single instance, or if we need to update it.

We can solve these problems by running multiple cache instances. We could split the caches up according to the type of data to be cached. For example, in a collaborative document-editing system, we might have one cache for Users, one cache for Documents, one cache for Comments, and so on. This works, but we may still have more requests for one or more of these data types than a single cache instance can handle.

A way to solve this problem is to replicate the caches. In a replicated cache setup, we would have two or more caches serving the same data - for instance, we might choose to run three replicas of the Users cache. Reads for Users data can go to any of the three replicas. However, when Users data is updated, we must either:

- invalidate all three instances of the Users cache
- tolerate stale data until the TTL expires (as described above)

The need to invalidate all instances of the Users cache adds cost to every write operation: more latency, because we must wait for the slowest cache instance to respond; as well as higher use of bandwidth and other computing resources, because we have to do work on every cache instance that stores the data being invalidated. This cost increases the more replicated instances we run.

There is an additional complication: what happens if we write to the Users database table, but cannot connect to one or more of the cache servers? In practice, this can happen and it means that inconsistency between datastores and caches is always a possibility in distributed systems.

There is further discussion of this problem below in the section on [CAP theorems](#the-cap-theorem).

##### Sharded Caches {#sharded-caches}

Another approach to scaling cache systems is to shard the data instead of replicating it. This means that your data is split across multiple machines, instead of each instance of your cache storing all of the cached items. This is a good choice when the working set of recently accessed data is too large for any one machine. Each machine hosts a _partition_ of the data set. In this setup there is usually a router service that is responsible for proxying cache requests to the correct instance: the instance that stores the data to be accessed.

Data sharding can be vertical or horizontal. In vertical sharding, we store different fields on different machines (or sets of machines). In horizontal sharding, we split the data up by rows.

In the case of horizontal sharding, we can shard algorithmically - meaning that the shard to store a given row is determined by a function - or dynamically, meaning that the system maintains a lookup table that tracks where data ranges are stored.

Read [Introduction to Distributed Data Storage by Quentin Truong](https://towardsdatascience.com/introduction-to-distributed-data-storage-2ee03e02a11d) for more on sharding.

A simple algorithmic sharding implementation might route requests to one of _N_ routers using a modulo operation: `cache_shard = id % num_shards`. The problem is that whenever a shard is added or removed, most of the keys will now be routed to a different cache server. This is equivalent to restarting all of our caches at once and starting cold. As discussed above, this is bad for performance, and potentially could cause an outage.

This problem is usually solved using a technique called [consistent hashing](https://en.wikipedia.org/wiki/Consistent_hashing). A consistent hash function maps a key to a data partition, and it has the very useful property that when the number of partitions changes, most keys do not get remapped to a different partition. This means that we can add and remove cache servers safely without risking reducing our cache hit rate by too much. Consistent hashing is a very widely used technique for load balancing across stateful services.

You can read about [how Pinterest scaled its Cache Infrastructure](https://medium.com/pinterest-engineering/scaling-cache-infrastructure-at-pinterest-422d6d294ece) - it’s a useful example of a real architecture. They use an open-source system called [mcrouter](https://github.com/facebook/mcrouter/wiki) to horizontally scale their use of memcached. Mcrouter is one of the most widely-used distributed cache solutions in industry.

#### Questions: {#questions}

- What is a cache hit rate and why is it important?
- What do we mean by ‘cold’ or ‘warm’ when we discuss caches?
- Why do we use consistent hashing when we shard data?
- When should we consider sharding a cache rather than replicating it?
- Why do we need cache invalidation?
- What is a TTL used for in caching?
- Why should we cache in a standalone cache service, as opposed to within our application?

### Databases {#databases}

Most web applications we build include at least one database. Databases don’t have to be distributed systems: you can run a database on one machine. However, there are reasons that you might want to run a distributed database across two or more machines. The logic is similar to the rationale for running more than one cache server, as described above.

1. Reliability: you don’t want to have an interruption in service if your database server experiences a hardware failure, power failure, or network failure.
2. Capacity: you might run a distributed database to handle more load than a single database instance can serve.

To do either of these things we must _replicate_ the data, which means to copy the data to at least one other database instance.

#### The CAP Theorem {#the-cap-theorem}

Before discussing distributed datastores further, we must introduce the CAP Theorem.

The [CAP Theorem](https://en.wikipedia.org/wiki/CAP_theorem) is a computer science concept that states that any distributed data store can provide at most two of the following three properties:

- Consistency: every read should receive the most recently written data or else an error
- Availability: every request receives a response, but there is no guarantee that it is based on the most recently written data
- (Network) Partition tolerance: the system continues to operate even if network connectivity between some or all of the computers running the distributed data store is disrupted

This seems complicated, but let’s break it down.

###### Network Partition {#network-partition}

A network partition means that your computers cannot all communicate over the network. Imagine that you are running your service in two datacenters, and the fiber optic cables between them are dug up (there should be redundancy in these paths, of course, but accidents do happen that can take out multiple connections). This is a network partition. Your servers won’t be able to communicate with servers in the opposite datacenter. Configuration problems, or simply too much network traffic can also cause network partitions.

Network partitions do occur, and there is no way that you can avoid this unpleasant fact of life. So in practice, distributed datastores have to choose between being consistent when network failures occur, or being available. It’s not a matter of choosing any two properties: you must choose either consistency and network partition tolerance or availability and network partition tolerance.

###### Consistency and availability {#consistency-and-availability}

Choosing consistency means that your application won’t be available on one side of the partition (or it might be read-only and serving old, stale data). Choosing availability means that your application will remain available, including for writes, but when the network partition is healed, you need a way of merging the writes that happened on different sides of the partition.

The CAP Theorem is a somewhat simplified model of distributed datastores, and it has been [criticised](https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html) on that basis. However, it remains a reasonably useful model for learning about distributed datastore concepts. In general, replicated traditional datastores like MySQL choose consistency in the event of network partitions; [NoSQL datastores](https://www.mongodb.com/nosql-explained) such as Cassandra and Riak tend to choose availability.

(However, the behaviour of systems in the event of network failure can also vary based on how that specific datastore is configured, in terms of how .)

#### Leader/Follower Datastore Replication {#leader-follower-datastore-replication}

The most straightforward distributed datastore is the leader/follower datastore, also known as primary/secondary or single leader replication. The aim is to increase the availability and durability of the data: in other words, if we lose one of the datastore machines, we should not lose data and we should still be able to run the system.

###### Synchronous Replication {#synchronous-replication}

To do something synchronously means to do something at the same time as something else.

In synchronous replication, the datastore client sends all writes to the leader. The leader has one or more followers. For every write, the leader sends those writes to its followers, and waits for all its followers to acknowledge completion of the write before the leader sends an acknowledgement to the client. Think of a race where all the runners are tied together with a rope. Leader and followers commit the data as part of the same database operation. Reads can go either to the leader, or to a follower, depending on how the datastore is configured. Reading from followers means that the datastore can serve a higher read load, which can be useful in applications that are serving a lot of traffic.

There is one problem with synchronous replication: availability. Not only must the leader be available, but _all of the followers_ must be available as well. This is a problem: the system’s availability for writes will actually be lower than a single machine, because the chances of one of _N_ machines being unavailable are by definition higher than the chances of one machine being unavailable, because we have to add up the probabilities of downtime.

For example, if one machine has 99.9% uptime and 0.1% downtime, and we have three machines, then we would expect the availability for all three together to be closer to 99.7% (i.e. 100% - 3 x 0.1%). Adding more replicas makes this problem worse. Response time is also a problem with synchronous replication, because the leader has to wait for all followers. This means that the system cannot commit writes faster than the slowest follower.

<p id="gdcalert2" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: inline image link here (to images/image2.png). Store image on your image server and adjust path/filename/extension if necessary. </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert3">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>

![alt_text](images/image2.png "image_tooltip")

###### Asynchronous and Semisynchronous Replication {#asynchronous-and-semisynchronous-replication}

We can solve these performance and availability problems by not requiring the leader to write to _all_ the followers before returning to the client. Asynchronous replication means that the leader can commit the transaction before replicating to any followers. Semisynchronous replication means that the leader must confirm the write at a subset of its followers before committing the transaction.

<p id="gdcalert3" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: inline image link here (to images/image3.png). Store image on your image server and adjust path/filename/extension if necessary. </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert4">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>

![alt_text](images/image3.png "image_tooltip")

Many databases provide configuration settings to control how many replicas need to acknowledge a write before the write can be considered complete. [MySQL](https://dev.mysql.com/doc/refman/8.0/en/replication.html), for example, defaults to asynchronous replication. MySQL can also support fully synchronous replication as described in the previous section, or semisynchronous replication, which means that at least one follower has committed every write before the leader considers the write complete.

However, asynchronous and semisynchronous replication are no free lunch. If using these forms of replication, followers can have data that is _behind_ the data that the leader stores. Data on followers can be different to the leader and can also vary between followers. This means that a client might write a value to the datastore and then read back an older value for that same field from a follower. A client might also read the same data twice from two different followers and see an older value on the second read. This can result in display of incorrect data, or application bugs.

###### Transaction ID {#transaction-id}

This problem can be solved by use of a transaction ID. Basically, the client has to track transaction IDs for their writes and reads and send that value with their reads to followers. Followers don’t return data until their replication has caught up with the transaction ID that the client specifies. This involves changing your application code, which is inconvenient. One solution to this problem is using a proxy (another program that relays requests and responses) between your application and your database. ProxySQL is a proxy that is designed to work with replicated MySQL databases, and it [supports tracking transaction IDs ](https://proxysql.com/blog/proxysql-gtid-causal-reads/)in this way.

#### Time and Ordering in Distributed Systems {#time-and-ordering-in-distributed-systems}

The kind of transaction IDs being discussed here are _logical timestamps_. Logical timestamps aren’t the same as system time, or time as we would read it from a clock. They are an ascending sequence of IDs that can be compared. Given two logical timestamps, we can tell if one event happened before the other event or not: this is what lets us solve our problems with data consistency above.

Logical timestamps create an ordering based on communications in distributed systems.

When the database leader sends transaction data, it sends a message along with it. The message is a transaction ID, or logical timestamp. A database replica can order its data in time – tell which data is current – using the transaction IDs.

Why not just use system time to figure out ordering in distributed systems? The reason is that in a distributed system we normally have no guarantees that different computers’ system clocks are in sync. We can use [Network Time Protocol ](https://en.wikipedia.org/wiki/Network_Time_Protocol)(NTP) to synchronise within a few milliseconds, but it is not perfect - and challenges can arise with respect to leap seconds or other problems. Some specialised systems do use special hardware clocks to keep systems synchronised (such as [Google’s Spanner datastore](https://cloud.google.com/spanner/docs/true-time-external-consistency)), but this is unusual.

Further Reading on Time, Database Replication and Failover:

- [There is No Now: Problems with simultaneity in distributed systems](https://queue.acm.org/detail.cfm?id=2745385)
- [Database Replication Explained. Part 1 — Single Leader Replication | by Zixuan Zhang | Towards Data Science](https://towardsdatascience.com/database-replication-explained-5c76a200d8f3)
- [DB Replication (II): Failure recovery fundamentals](https://www.brainstobytes.com/db-replication-i-introduction-to-database-replicationdb-replication-ii-failure-recovery-fundamentals/)

#### Questions: {#questions}

- What are the challenges associated with recovering from a failed database replica?

### Sharding {#sharding}

As we saw with caches, replication can help us to scale our read load for databases. However, at a certain scale, sharding becomes necessary. The general considerations are quite similar to sharding a caching layer. Read [Understanding Database Sharding](https://www.digitalocean.com/community/tutorials/understanding-database-sharding) for more detail.

A common sharded datastore used in industry is Vitess, originally developed at YouTube but available as an open-source project. It is an [orchestration layer](https://www.databricks.com/glossary/orchestration) designed for operating a massively sharded fleet of MySQL database instances and doing so reliably. It is worth reading some of the Vitess documentation:

- [Vitess - What is Vitess](https://vitess.io/docs/19.0/overview/whatisvitess/)
- [Vitess - Architecture](https://vitess.io/docs/19.0/overview/architecture/)
- [Vitess - History](https://vitess.io/docs/19.0/overview/history/)
- [Vitess - Scalability Philosophy](https://vitess.io/docs/19.0/overview/scalability-philosophy/)

#### Questions: {#questions}

- Why are smaller data shards recommended?
- What kind of changes might you need to make to your application before moving to a sharded datastore?
- What is a cell? What happens if a cell fails?

### Project work for this section {#project-work-for-this-section}

[https://github.com/CodeYourFuture/immersive-go-course/tree/main/memcached-clusters](https://github.com/CodeYourFuture/immersive-go-course/tree/main/memcached-clusters)

#### Stateful Services: Additional Reading {#stateful-services-additional-reading}

Discuss these papers in office hours with Laura Nolan.

Some of these are academic papers. If you aren’t used to reading academic writing, take some time to read [How to read and understand a scientific article](https://violentmetaphors.files.wordpress.com/2018/01/how-to-read-and-understand-a-scientific-article.pdf) first.

[Dynamo: Amazon’s Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf): This paper from 2007 describes how Dynamo, Amazon’s scalable K/V store which they originally developed to support shopping carts, works. The core of the system is consistent hashing, as described earlier in this section.

- What kind of sharding is this, algorithmic or dynamic?

[Amazon DynamoDB: A Scalable, Predictably Performant, and Fully Managed NoSQL Database Service](https://www.usenix.org/system/files/atc22-elhemali.pdf): This much more recent paper shows DynamoDB’s evolution into a Cloud service, with quite different demands to the original shopping cart use case.

- Describe two ways in which this architecture differs from the original architecture described in the 2007 paper above, and why.

[The Google File System](https://static.googleusercontent.com/media/research.google.com/en//archive/gfs-sosp2003.pdf): Classic paper from 2003 describing Google’s first distributed file system architecture.

- GFS uses a different kind of sharding model to DynamoDB - is it algorithmic or dynamic? Why do you think GFS’s designers chose that model?

[Bigtable: A Distributed Storage System for Structured Data](https://static.googleusercontent.com/media/research.google.com/en//archive/bigtable-osdi06.pdf): Paper which builds on the Google File System (GFS) paper about how Google built their K/V store, Bigtable, on top of GFS.

- What constraints did GFS’s properties impose on the BigTable design?

[After the Retrospective: The 2017 Amazon S3 Outage](https://www.gremlin.com/blog/the-2017-amazon-s-3-outage/): You can read this analysis of an AWS S3 storage service outage.

- Draw a rough diagram of the S3 architecture showing what happened.

##
