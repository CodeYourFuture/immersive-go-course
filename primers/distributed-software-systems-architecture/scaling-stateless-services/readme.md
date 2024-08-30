<!--forhugo
+++
title="3. Scaling Stateless Services"
+++
forhugo-->

# 3

## Scaling Stateless Services {#section-3-scaling-stateless-services}

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
