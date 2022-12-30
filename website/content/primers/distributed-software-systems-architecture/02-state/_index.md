<!-- Output copied to clipboard! -->

<!-----

You have some errors, warnings, or alerts. If you are using reckless mode, turn it off to see inline alerts.
* ERRORs: 0
* WARNINGs: 0
* ALERTS: 3

Conversion time: 2.297 seconds.


Using this Markdown file:

1. Paste this output into your source file.
2. See the notes and action items below regarding this conversion run.
3. Check the rendered output (headings, lists, code blocks, tables) for proper
   formatting and use a linkchecker before you publish this page.

Conversion notes:

* Docs to Markdown version 1.0β34
* Fri Dec 30 2022 09:24:19 GMT-0800 (PST)
* Source doc: Copy of Distributed Systems Primer
* This document has images: check for >>>>>  gd2md-html alert:  inline image link in generated source and store images to your server. NOTE: Images in exported zip file from Google Docs may not appear in  the same order as they do in your doc. Please check the images!


WARNING:
You have 2 H1 headings. You may want to use the "H1 -> H2" option to demote all headings by one level.

----->

<p style="color: red; font-weight: bold">>>>>>  gd2md-html alert:  ERRORs: 0; WARNINGs: 1; ALERTS: 3.</p>
<ul style="color: red; font-weight: bold"><li>See top comment block for details on ERRORs and WARNINGs. <li>In the converted Markdown or HTML, search for inline alerts that start with >>>>>  gd2md-html alert:  for specific instances that need correction.</ul>

<p style="color: red; font-weight: bold">Links to alert messages:</p><a href="#gdcalert1">alert1</a>
<a href="#gdcalert2">alert2</a>
<a href="#gdcalert3">alert3</a>

<p style="color: red; font-weight: bold">>>>>> PLEASE check and correct alert issues and delete this message and the inline alerts.<hr></p>

# The Distributed Software Systems Architecture Primer

_Last updated: 29 Nov 2022 by Laura Nolan_

[TOC]

##

## About this document {#about-this-document}

This document outlines a short course in distributed systems architecture patterns and common issues and considerations. It is aimed at students with some knowledge of computing and programming, but without significant professional experience in building and operating distributed software systems.

Learning outcomes for this primer:

- You will be able to explain the reasons that we build distributed systems
- You will understand how to build your systems to deal with common types of failure, such as slow backends and network partitions
- You will be able to describe the pros and cons of asynchronous, semisynchronous, and synchronous database replication
- You will be able to describe the difference between asynchronous work systems (like pipelines) and serving systems

There are five sections in this document, one per sprint.

## Why Distributed Systems? {#why-distributed-systems}

All client-server systems are distributed systems. Any computer system that involves communication between multiple physical computers is a distributed system. Any system that separates data storage from web serving, or which uses cloud services via APIs, is a distributed system.

A lot of distributed systems in operation go far beyond these simple architectures. Organisations build distributed systems to make their systems more reliable in the face of failure. A distributed system can serve a user from any one of several physical computers, perhaps on different networking and power domains.

Some workloads are too high to be served from a single machine. We use distributed systems techniques to spread workloads across many machines This is called horizontal scalability. In distributed systems you can serve requests closer to your users, which is faster, and a better user experience (it is remarkable how different many web applications feel when accessed from Australia or South Africa).

Almost all computer systems that we build today are distributed systems, whether large or small.

##

## Section 1: Reliable RPCs {#section-1-reliable-rpcs}

### About Remote Procedure Calls (RPCs) {#about-remote-procedure-calls-rpcs}

The Remote Procedure Call is the fundamental building block of distributed systems: they are how you can run a procedure - i.e. a piece of code - on another machine. Without RPCs we cannot build an efficient distributed system: they are essential to the existence of GMail, Netflix, Spotify, Facebook, and all other large-scale distributed software systems. Keeping RPCs flowing reliably is a significant part of most infrastructure engineering jobs.

RPCs aren’t the same as [HTTP requests](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview). HTTP is a text-based protocol, oriented towards resources or entities. RPCs are oriented towards calling a specific function. RPCs are highly structured; normally the request and response are encoded in efficient binary formats which are not human-readable. However these structured formats can be used to transmit data between computers using much less bandwidth than JSON or other text formats. Structured RPCs are also more efficient to parse than text, meaning that CPU usage will be lower. It depends on the logic of the specific application, but parsing of requests and encoding of responses often uses a very significant percentage of the computing resources needed by a web service.

HTTP APIs are usually limited to CRUD (Create, Read, Update, Delete) operations; RPCs can perform any kind of operation. If you want to efficiently integrate two systems that you control, RPCs are a good choice. If you want to provide an API for use by developers outside of your organisation then HTTP APIs are generally a better choice, because they are simpler to develop against and all programming languages provide good HTTP support.

<p id="gdcalert1" ><span style="color: red; font-weight: bold">>>>>>  gd2md-html alert: inline image link here (to images/image1.png). Store image on your image server and adjust path/filename/extension if necessary. </span><br>(<a href="#">Back to top</a>)(<a href="#gdcalert2">Next alert</a>)<br><span style="color: red; font-weight: bold">>>>>> </span></p>

![alt_text](images/image1.png "image_tooltip")

This diagram is from [Nelson and Birell’s classic paper](http://birrell.org/andrew/papers/ImplementingRPC.pdf) on implementing RPCs. Both client and server use _stubs_ to define the structure of the request and the response.

There is also an RPCRuntime library, which takes care of the mechanics of locating the server machine, transmitting the request, and waiting for the response.

Today, two of the most popular RPC frameworks are [Thrift](https://thrift.apache.org/) and gRPC. Thrift was developed at Meta (Facebook) and is popular in the Java ecosystem. gRPC was developed at Google and supports a variety of languages, but is most associated with the Go programming language.

The [gRPC quickstart guide for Go](https://grpc.io/docs/languages/go/quickstart/) explains how to define a simple gRPC service definition (this is the API that your program serves), compile it - which generates the user stubs and service stubs shown in the diagram above - and use it in your program.

### RPCs: What Could Go Wrong? {#rpcs-what-could-go-wrong}

The idea of RPCs is not complicated. The real challenge of distributed systems is dealing with the problems that unreliable networks create. Since the 1990s programmers have talked about the [Eight Fallacies of Distributed Computing](https://web.archive.org/web/20171107014323/http://blog.fogcreek.com/eight-fallacies-of-distributed-computing-tech-talk/) - in other words, things that most programmers wrongly assume are true.

**The most important fallacies are:**

1. The network is reliable;
2. Latency is zero;
3. Bandwidth is infinite;
4. The network is secure.

### Fallacy 1: The network is reliable {#fallacy-1-the-network-is-reliable}

#### Unreliable networks, errors, and retries {#unreliable-networks-errors-and-retries}

All networks that really exist, that are configured and accessed by human beings, served from machines sitting in a building, connected by cables and powered by electricity, are unreliable, because any of these things can break. In reality, your network is not reliable.

Sometimes communication from a client will fail to reach the server, or the server’s response will fail to return to the client, even though the work requested was completed.

RPCs can fail for other reasons. Sometimes servers will reject requests because they are overloaded and do not have the capacity to serve additional requests.

Graceful shutdowns of servers can cause client errors. All servers have a lifecycle: they are launched, they begin serving, and at some point, they will stop serving and be replaced. In a distributed system, we usually have multiple servers - so we can gracefully replace servers by sending requests elsewhere. Servers that are about to shut down can send error codes to clients attempting to begin new requests (the HTTP/2 [GOAWAY ](https://httpwg.org/specs/rfc7540.html#GOAWAY)code is an example of this). The client is expected to try another available server.

All client-server code needs to deal with errors. A common method is to _retry_ a request that failed, by re-sending it to a different server. Retries are a useful strategy but there are pitfalls. These considerations can be grouped under: Status codes, Idempotency, and Limits.

##### Status codes {#status-codes}

The first question to ask when retrying a request is whether it is a request that could possibly succeed. Errors returned by the gRPC library to a client include a [status code](https://grpc.github.io/grpc/core/md_doc_statuscodes.html) - these are somewhat similar to HTTP status codes. The status code can be generated either by the server or by the client (for example, if the client can’t connect to the server the client will generate an `UNAVAILABLE` status code and return it to the calling code). Some of these codes indicate that a request is never likely to succeed - for example,` INVALID_ARGUMENT` means that the request is malformed in some way - while others, like `UNAVAILABLE` or `RESOURCE_EXHAUSTED `might succeed if retried.

##### Idempotency {#idempotency}

The second question to ask is whether it is safe to retry the request or not. Some kinds of requests are[ idempotent](https://betterprogramming.pub/architecting-distributed-systems-the-importance-of-idempotence-138722a6b88e): this means that the result will be the same, whether you make the request once or multiple times. For example, setting an attribute to a specific value is idempotent, but adding a number to the existing value is not. Non-idempotent requests can’t be safely retried. This means that it’s best to design your systems with idempotent APIs wherever possible. Services can use unique request IDs to deduplicate requests, changing a non-idempotent API into one that may be safely retried.

##### Limits {#limits}

It is very important to limit the number of retries that a client makes to a service. If a service becomes overwhelmed with load and each of its clients repeatedly retries failing requests then it becomes a vicious cycle. The service will not be able to recover without intervention. Groups of continuously-retrying clients can create an effect similar to a [Distributed Denial of Service](https://www.cloudflare.com/en-gb/learning/ddos/what-is-a-ddos-attack/) (DDoS) attack.

A maximum of three retries is a reasonable rule of thumb for interactive workloads, in other words, when a human is waiting for a response. By the time your system has retried three times the user is likely to try again or to give up, so there is no point in making further attempts. However, some systems are not interactive, such as [batch jobs](https://en.wikipedia.org/wiki/Batch_processing) or [data pipelines](https://www.snowflake.com/guides/data-pipeline). In these systems it is reasonable to attempt more than three retries, after a long delay.

There are lots of ways to limit retries: by layer, by duration between attempts, and with circuit breakers.

###### Layers {#layers}

Systems may implement retrying at multiple levels. For example, it is not uncommon to see client retries combined with retries at a loadbalancing layer. This has a multiplicative effect: if both levels retry failed requests three times, then this can result in nine retries. It is best to limit the use of retries to one layer of your stack and to be consistent about where they are performed.

###### Duration Between Attempts {#duration-between-attempts}

Clients should increase the duration of their waits between retries exponentially: double it for the second retry attempt, again for the third, and so on. This technique is called _exponential backoff_. In cases where a backend service is briefly unavailable, the use of exponential backoff by its clients means that the backend service is less likely to be overwhelmed with load when it comes back up.

Jitter means a short variable wait before retrying. It is good practice to use jitter when retrying. The jitter, or wait time, is randomly chosen, between a minimum and maximum value. The idea of introducing such jitter is to avoid overwhelming backend servers with a coordinated ‘echo’ of an original spike of errors, which could cause a further overload.

Use jitter and exponential backoff together r when implementing retries. Many programming languages include libraries that implement exponential backoff with jitter - for example, see [GoLang’s go-retry library](https://pkg.go.dev/gopkg.in/retry.v1). Some loadbalancers, such as [Envoy Proxy](https://www.envoyproxy.io/docs/envoy/latest/faq/load_balancing/transient_failures), also allow you to configure these retry strategies at the loadbalancer layer.

###### Circuit Breakers {#circuit-breakers}

A very useful pattern for client-server requests is the [Circuit Breaker](https://martinfowler.com/bliki/CircuitBreaker.html). Use circuit breakers to limit potentially-dangerous retries. According to Martin Fowler: “The basic idea behind the circuit breaker is very simple. You wrap a protected function call in a circuit breaker object, which monitors for failures. Once the failures reach a certain threshold, the circuit breaker trips, and all further calls to the circuit breaker return with an error, without the protected call being made at all. Usually you'll also want some kind of monitor alert if the circuit breaker trips.”

### Fallacy 2: Latency is zero {#fallacy-2-latency-is-zero}

#### Variable latency and deadlines {#variable-latency-and-deadlines}

All real networks are connected over real physical distances. The distance between nodes is like the travel time between cities on Google Maps: the route, mode of transport, amount of traffic, and distance, all affect the travel time. Latency means the time taken for a request to traverse the network. Network latency is travel time.

Latency can be highly variable, particularly over long distances. Processing time at the server can also vary a lot, which matters because latency as perceived by clients is the sum of network latency and application processing time at the server.

High latency has implications for clients. It is always a bad idea to wait indefinitely for a response from a server. It is entirely possible to end up waiting forever. For example, if the server ‘hangs up’ and your TCP connection does not have keepalives enabled, you will wait indefinitely. When this happens, no useful work gets done and it can result in problems that are difficult to track down, because a client that waits indefinitely generates no errors.

If some abnormal condition means that many clients in a system are experiencing long waits for some service (even a small and trivial one) it can result in widespread problems as everything slows down waiting for the slowest service. Again, this is generally difficult to troubleshoot. When everything is slow, the source of that slowness is not obvious - it can be quite hard to pinpoint the reason.

Instead, define a specific deadline for every RPC request and deal with it as an error when the deadline is past. This is a much better approach because it is usually very easy to find errors in monitoring and in logs. When our systems are not operating correctly, it is very important to be able to quickly find the problem and solve it.

Set deadlines much higher than typical response latencies. In general, the deadline should be higher than your typical 99.9% response time - this means that only one in 1000 requests would be slower than the deadline in normal operation.

You should be able to find out what your observed response latencies are either by using your monitoring system or by analysing application logs. Your organisation may run an open-source monitoring tool such as [Prometheus ](https://prometheus.io/)or [Graphite](https://grafana.com/oss/graphite/), or you may use a SaaS monitoring system such as [New Relic](https://newrelic.com/) or [DataDog](https://www.datadoghq.com/). All of these systems allow you to [instrument your application](https://alex.dzyoba.com/blog/go-prometheus-service/) so that you can understand runtime behaviour, such as response latencies and error rates.

It is very important to monitor your observed error rates and latencies and to investigate if you see a significant number of timeouts. The threshold you set here depends on the degree of reliability required for your application, but a typical threshold for a high-traffic application would be more than 0.01% of requests timing out, or failing for any reason.

You should also monitor whether your observed request latencies are approaching the deadline you have set. For example, you might configure your monitoring system to alert if the 99th [percentile](https://en.wikipedia.org/wiki/Percentile) latency (i.e. the slowest 1% of requests) exceeds 80% of your configured request deadline. The reason for this is that request latencies may drift higher as changes are made to systems over time. If request latencies become so slow that many requests begin to take longer than the deadline, then your system is doing wasted work. Clients will retry failed requests, and your system may become overloaded as a result. AWS’s account of their [DynamoDB Service Disruption](https://aws.amazon.com/message/5467D2/) in 2015 is a good description of this phenomenon.

### Fallacy 3: Bandwidth is infinite {#fallacy-3-bandwidth-is-infinite}

#### Limited bandwidth {#limited-bandwidth}

All systems that actually exist in the world have limits. Computers do not have infinite memory. Networks do not have infinite bandwidth. All the nodes in your system are physically real and all real things are finite, or limited.

Network bandwidth is not unlimited. Passing very large requests or responses between machines may be unreliable. As we saw above, it is helpful to set deadlines on requests. However, it is possible that the time to process a very large request might exceed an RPC deadline. Some services also attempt to buffer incoming requests before beginning to process them. Because no computer has infinite memory, it is common (and good practice) to limit the size of buffered requests. However, very large requests may exceed the configured buffer size and fail. In most cases, very large requests bodies are not legitimate requests. Where large amounts of data must be transferred, it is worth considering a specialised protocol (such as [tus, which has a Go implementation](https://pkg.go.dev/github.com/tus/tusd?utm_source=godoc)).

### Fallacy 4: The network is secure {#fallacy-4-the-network-is-secure}

Your network exists in the world and interacts with other networks, systems, and actors, with their own objectives. Many of these objectives conflict with your own, and some are hostile. Networks can never be assumed to be secure. You should assume that all network links are being eavesdropped, and therefore by default encrypt all data on the wire. This is modern best practice. Applications should use methods such as [mTLS ](https://www.cloudflare.com/en-gb/learning/access-management/what-is-mutual-tls/)to secure communications between clients and servers.

Here’s how you do [mTLS in Go](https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go), although in a real production environment you also need to manage the creation, distribution and renewal of keys. Here’s how [Zendesk does mTLS key management](https://zendesk.engineering/implementing-mtls-and-securing-apache-kafka-at-zendesk-10f309db208d).

### Questions {#questions}

- ​​What is latency? Why is it important?
- Why is idempotency a desirable property in a web service?
- When should I use RPCs, as opposed to simple HTTP requests?

### Project work for this section {#project-work-for-this-section}

See:

- [https://github.com/CodeYourFuture/immersive-go-course/tree/main/http-auth](https://github.com/CodeYourFuture/immersive-go-course/tree/main/http-auth) [should have been done in prework]
- [https://github.com/CodeYourFuture/immersive-go-course/tree/main/grpc-client-server](https://github.com/CodeYourFuture/immersive-go-course/tree/main/grpc-client-server)

##

## Section 2: State {#section-2-state}

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

- [https://vitess.io/docs/14.0/overview/whatisvitess/](https://vitess.io/docs/14.0/overview/whatisvitess/)
- [https://vitess.io/docs/14.0/overview/architecture/](https://vitess.io/docs/14.0/overview/architecture/)
- [https://vitess.io/docs/14.0/overview/history/](https://vitess.io/docs/14.0/overview/history/)
- [https://vitess.io/docs/14.0/overview/scalability-philosophy/](https://vitess.io/docs/14.0/overview/scalability-philosophy/)

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

## Section 3: Scaling Stateless Services {#section-3-scaling-stateless-services}

### Microservices or monoliths {#microservices-or-monoliths}

A monolith is a single large program that encapsulates all the logic for your application. A microservice architecture, on the other hand, splits application functionality across a number of smaller programs, which are composed together to form your application. Both have advantages and drawbacks. ​​

Read [Microservices versus Monoliths](https://www.atlassian.com/microservices/microservices-architecture/microservices-vs-monolith) for a discussion of microservices and monoliths. Optionally, watch - a comedy about the extremes of microservice-based architectures. The middle ground between a single monolith and many tiny microservices is to break the application into a more moderate number of services, each of which have high cohesion (or relatedness). Some call this approach [‘macroservices’](https://www.geeksforgeeks.org/software-engineering-coupling-and-cohesion/).

### Horizontal Scaling and Load Balancing {#horizontal-scaling-and-load-balancing}

We have seen how stateless and stateful services can be scaled horizontally - i.e. run across many separate machines to provide a scalable service that can serve effectively unlimited load - with the correct architecture.

Load balancers are an essential component of horizontal scaling. For stateless services, the load balancers are general-purpose proxies (like Nginx, Envoy Proxy, or HAProxy). For scaling datastores, proxies are typically specialised, like mcrouter or the Vitess vtgates.

### How load balancers work {#how-load-balancers-work}

It is worth understanding a little of how load balancers work. Load balancers today are typically based on software, although hardware load balancer appliances are still available. The biggest split is between Layer 4 load balancers, which focus on balancing load at the level of TCP/IP packets, and Layer 7 load balancers, which understand HTTP.

Read [Introduction to modern network load balancing and proxying by Matt Klein](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236).

- Name 5 functions of load balancers
- Why might you use a Layer 7 loadbalancer instead of a Layer 4?
- When might you use a Layer 4 loadbalancer?
- Give a reason to use a sidecar proxy architecture rather than a middle proxy
- Why would you use Direct Server Return?
- What is a Service Mesh? What are the advantages and disadvantages of a service mesh?
- What is a VIP? What is Anycast?

#### Round Robin and other loadbalancing algorithms {#round-robin-and-other-loadbalancing-algorithms}

Many load balancers use Round Robin to allocate connections or requests to backend servers. This means that they assign the first connection or request to the first backend, the second to the second, and so on until they loop back around to the first backend. This has the virtue of being simple to understand, doesn’t need a lot of complex state or feedback to be managed, and it’s difficult for anything to go very wrong with it.

Read about other [loadbalancing algorithms.](https://www.cloudflare.com/en-gb/learning/performance/types-of-load-balancing-algorithms/)

Now consider the weighted response time loadbalancing algorithm, which sends the most requests to the server that responds fastest. What can go wrong here?

If one server is misconfigured and happens to be very very rapidly serving errors or empty responses, then loadbalancers configured to use weighted response time algorithms would send more traffic to the faulty server.

##### DNS {#dns}

There is one more approach to load balancing that is worth knowing about: [DNS Load Balancing](https://www.cloudflare.com/en-gb/learning/performance/what-is-dns-load-balancing/). DNS-based load balancing is often used to route users to a specific set of servers that is closest to their location, in order to minimize latency (page load time).

### Performance: Edge Serving and CDNs {#performance-edge-serving-and-cdns}

Your users may be anywhere on Earth, but quite often, your serving infrastructure (web applications, datastores, and so on) is all in one region. This means that users in other continents may find your application slow. A network round-trip between Sydney in Australia and most locations in the US or Europe takes around 200 milliseconds. 200ms is not overly long, but the problem is that serving a user request may take several round trips.

##### Round trips {#round-trips}

First, the user may need to look up the DNS name for your site (this may be cached nearby). Next, they need to open a TCP connection to one of your servers, which requires a network round trip. Finally, the user must perform a [SSL handshake](https://zoompf.com/blog/2014/12/optimizing-tls-handshake/) with your server, which also requires one or two network round trips, depending on configuration of client and server (a recent session may be resumed in one round trip if both client and server support TLS 1.3).

All of this takes place before any data may be returned to the user, and, unless there is already an open TCP connection between the user and the website, involves an absolute minimum of <span style="text-decoration:underline;">three network round trips</span> before the first byte of data can be received by the client.

##### SSL Termination at the Edge {#ssl-termination-at-the-edge}

SSL termination _at the edge_ is the solution to this issue. This involves running some form of proxy much nearer to the user which can perform the SSL handshake with the user. If the user need only make network round trips of 10 milliseconds to a local Point of Presence (PoP), as opposed to 200 milliseconds to serving infrastructure in a different continent, then a standard TCP connection initiation and SSL handshake will take only around 60 milliseconds, as opposed to 1.2 seconds.

Of course, the edge proxies must still relay requests to your serving infrastructure, which remains 200 milliseconds of network latency away. However, the edge proxies will have persistent encrypted network connections to your servers: there is no per-request overhead for connection setup or handshakes.

##### Content Delivery Networks {#content-delivery-networks}

Termination at the edge is a service often performed by [Content Delivery Networks](https://en.wikipedia.org/wiki/Content_delivery_network) (CDNs). CDNs can also be used to cache static assets such as CSS or images close to your users, in order to reduce site load time as well as reducing load on your origin servers. You can also run your own compute infrastructure close to your users. However, this is an area of computing where being big is an advantage: it is hard to beat the number of edge locations that large providers like Cloudflare, Fastly, and AWS operate.

[Edge Regions in AWS](https://www.lastweekinaws.com/blog/what-is-an-edge-location-in-aws-a-simple-explanation/) is worth reading to get an idea of the scale of Amazon’s edge presence.

### QUIC {#quic}

It is worth being aware of [QUI](https://en.wikipedia.org/wiki/QUIC)C, an emerging network protocol that is designed to be faster for modern Internet applications than TCP/IP. While it is by no means ubiquitous yet, it is certainly an area that the largest Internet companies are investing in. HTTP/3, the next major version of the HTTP protocol, uses QUIC.

Read about [HTTP over QUIC](https://blog.cloudflare.com/http3-the-past-present-and-future/) in the context of the development of the HTTP protocol.

### Autoscaling {#autoscaling}

Aside from load balancing, the other major component of successful horizontal scaling is autoscaling. Autoscaling means to scale the number of instances of your service up and down according to the load that the service is experiencing. This can be more cost-effective than sizing your service for expected peak loads.

On AWS, for example, you can create an Autoscaling Group (ASG) which acts as a container for your running EC2 instances. ASGs can be [configured](https://docs.aws.amazon.com/autoscaling/ec2/userguide/scale-your-group.html) to scale up or scale down the number of instances based on a schedule, based on[ predicted load](https://docs.aws.amazon.com/autoscaling/ec2/userguide/ec2-auto-scaling-predictive-scaling.html) (based on the past two weeks of history), or based on [current metrics](https://docs.aws.amazon.com/autoscaling/ec2/userguide/as-scale-based-on-demand.html). Kubernetes [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) is a tool similar to ASGs scaling policies for Kubernetes.

##### CPU utilisation {#cpu-utilisation}

CPU utilisation is very commonly used as a scaling signal. It is a good general proxy for how ‘busy’ a particular instance is, and it’s a ‘generic’ metric that your platform understands, as opposed to an application-specific metric.

The CPU utilisation target should not be set too high. Your service will take some time to scale up when load increases. You might experience failure of a subset of your infrastructure (think of a bad code push to a subset of your instances, or a failure of an AWS Availability Zone). Your service probably cannot serve reliably at close to 100% CPU utilisation. 40% utilisation is a common target.

### Autoscaling and Long-Running Connections {#autoscaling-and-long-running-connections}

One case where autoscaling does not help you to manage load is when your application is based on long-running connections, such as websockets or gRPC streaming. Managing these at scale can be challenging. Read the article [Load balancing and scaling long-lived connections in Kubernetes](https://learnk8s.io/kubernetes-long-lived-connections).

- Why doesn’t autoscaling work to redistribute load in systems with long-lived connections?
- How can we make these kinds of systems robust?

### Project work for this section {#project-work-for-this-section}

- [https://github.com/CodeYourFuture/immersive-go-course/tree/main/multiple-servers](https://github.com/CodeYourFuture/immersive-go-course/tree/main/multiple-servers)
- Optional: do this HashiCorp tutorial. It will give you hands-on experience with seeing health checks can be used to manage failure. [https://learn.hashicorp.com/tutorials/consul/introduction-chaos-engineering?in=consul/resiliency](https://learn.hashicorp.com/tutorials/consul/introduction-chaos-engineering?in=consul/resiliency)
- Optional: demonstrate circuit breaking. [https://learn.hashicorp.com/tutorials/consul/service-mesh-circuit-breaking?in=consul/resiliency](https://learn.hashicorp.com/tutorials/consul/service-mesh-circuit-breaking?in=consul/resiliency)
- Optional: see different kinds of load-balancing algorithms in use: [https://learn.hashicorp.com/tutorials/consul/load-balancing-envoy?in=consul/resiliency](https://learn.hashicorp.com/tutorials/consul/load-balancing-envoy?in=consul/resiliency)
- Optional: do this tutorial which demonstrates autoscaling with minikube (you will need to install minikube on your computer if you don’t have it). It will give you hands-on experience in configuring autoscaling, plus some exposure to Kubernetes configuration.
  - You may need to run this first: `minikube addons enable metrics-server`
  - [https://devops.novalagung.com/kubernetes-minikube-deployment-service-horizontal-autoscale.html](https://devops.novalagung.com/kubernetes-minikube-deployment-service-horizontal-autoscale.html)

##

## Section 4: Asynchronous Work and Pipelines {#section-4-asynchronous-work-and-pipelines}

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

##

## Section 5: Distributed Locking and Distributed Consensus {#section-5-distributed-locking-and-distributed-consensus}

In a program, sometimes we need to lock a resource. Think about a physical device, like a printer: we only want one program to print at a time. Locking applies to lots of other kinds of resources too, often when we need to update multiple pieces of data in consistent ways in a multi-threaded context.

We need to be able to do locking in distributed systems as well. Read Martin Kleppmann’s article [How to do distributed locking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html).

- What are the two main reasons to do distributed locking?

It turns out that in terms of computer science, distributed locking is theoretically equivalent to reliably electing a single leader (like in database replication, for example). It is also logically the same as determining the exact order that a sequence of events on different machines in a distributed system occur. All of these are useful.

An algorithm that you should know about is the RAFT distributed consensus protocol: read [https://raft.github.io/raft.pdf](https://raft.github.io/raft.pdf).

- Under what circumstances is a RAFT cluster available?
- What does the leader do?
- Is there always a leader?
- How is a leader elected?
- What happens if one replica falls behind the leader, i.e. the leader has committed transactions that the replica doesn’t have?
- What happens if a server is replaced?

Read about the operational characteristics of distributed consensus algorithms in [Managing Critical State.](https://sre.google/sre-book/managing-critical-state/)

- What are the scaling limitations of distributed consensus algorithms?
- How can we scale read-heavy workloads?

### Project work for this section {#project-work-for-this-section}

Follow along with this long tutorial and try to get a working Raft implementation running. This is a real programming challenge.

- [https://eli.thegreenplace.net/2020/implementing-raft-part-0-introduction/](https://eli.thegreenplace.net/2020/implementing-raft-part-0-introduction/)

##

## Further Optional Reading {#further-optional-reading}

Discuss any of these pieces at office hours with Laura.

Dan Luu’s list of postmortems includes a lot of interesting stories of real-world distributed systems failure: [https://github.com/danluu/post-mortems](https://github.com/danluu/post-mortems)

[Jeff Hodges' Notes on Distributed Systems for Young Bloods](https://www.somethingsimilar.com/2013/01/14/notes-on-distributed-systems-for-young-bloods/) is a very practical take on distributed systems topics.

[Alvaro Videla's blog post and talks about learning distributed systems](https://alvaro-videla.com/2015/12/learning-about-distributed-systems.html) are accessible and well organised.

[Marc Brooker’s blog](https://brooker.co.za/blog/) is full of interesting pieces, which are very approachable.

[The SRE Book](https://sre.google/sre-book/table-of-contents/) is available in full online and is worth reading - it addresses many aspects of operating distributed software systems.

Aphyr (Kyle Kingsbury) provides detailed [notes for a distributed systems class he teaches.](https://github.com/aphyr/distsys-class)

[A Distributed Systems Reading List](https://dancres.github.io/Pages/) will keep you reading excellent distributed systems papers for many many months.

# Glossary of Abbreviations {#glossary-of-abbreviations}

API : Application Programming Interface

CAP : Consistency Availability Partition tolerance

CDN : Content Delivery Network

CRUD : Create Read Update Delete

CPU : Central Processing Unit

DDoS : Distributed Denial of Service

DNS : Domain Name System

gRPC : google Remote Procedure Call

HTTP : Hypertext Transfer Protocol

HTTP2 : Hypertext Transfer Protocol 2

ID : Identity

LB : Load Balancer

mTLS : mutual Transport Layer Security

QUIC : Quick UDP Internet Connections

RPC : Remote Procedure Call

SSL : Secure Socket Layer

TCP : Transmission Control Protocol

TLS : Transport Layer Security

TTL : Time To Live

UDP : User Datagram Protocol

2PC : 2 Phase Commit

3PC : 3 Phase Commit
