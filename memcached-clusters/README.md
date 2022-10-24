# Memcached Clusters: Replicated and Sharded

This project should be done after reading Section 2 of the [Distributed Systems Primer](https://docs.google.com/document/d/1WoOTLTdtDqnL3fv3YVfI32kfySHqh7y1UfLizBJ3LXY/edit?usp=sharing). 

Timebox: 2 days

Learning objectives:
 * Understand the differences between sharded and replicated datastores
 * Use mcrouter, a widely-used proxy that can be used to create both sharded and replicated memcached clusters

## Read about mcrouter

Read the [FaceBook Engineering mcrouter blog post](https://engineering.fb.com/2014/09/15/web/introducing-mcrouter-a-memcached-protocol-router-for-scaling-memcached-deployments/).

## Running A Replicated Memcached Cluster

Make sure you have Docker and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine.

In the `replicated` directory of this project run:

```console
> docker-compose up -d
```

This will run an instance of the mcrouter proxy plus three separate memcached instances.
mcrouter is listening on prt 11211, and the three memcached instances are on ports 11212, 11213, and 11214.
Look at the `docker-compose.yml` file to understand what is happening.

You can run commands against mcrouter and memcached by sending commands to them.
Use the program `nc` (netcat) to do this. Try running these commands:

```console
> printf "set mykey 0 60 4\r\ndata\r\n" | nc localhost 11211
> printf "get mykey\r\n" | nc localhost 11211
```

The `set` command will store an association between a key called `mykey` and the value 4, with a TTL of 60 seconds. 
You should see the value 4 returned when you run the `get` command.
The 60 second TTL means that the value will disappear after 60 seconds - TTLs are discussed further in the [Distributed Systems Primer](https://docs.google.com/document/d/1WoOTLTdtDqnL3fv3YVfI32kfySHqh7y1UfLizBJ3LXY/edit?usp=sharing). 

Wait a minute and run the `get` command again: you should see the value 4 is no longer returned.

You can also try the `delete` command.

You have been running commands against the mcrouter proxy.
Set the value once more and then run the `get` command directly against the three memcached instances (on ports 11212, 11213, and 11214).

You should see that the value you set is returned from each individual cache. This is down to the replication scheme specified in the mcrouter command line, 
which you can see in the `docker-compose.yml` file. Read about [replicated mcrouter configuration](https://github.com/facebook/mcrouter/wiki/Replicated-pools-setup) 
and try to understand what the configuration is doing.

You can bring down your cluster by running
```console
> docker-compose down.
```

## Running A Sharded Memcached Cluster

Now, in the `sharded` directory of this project run:

```console
> docker-compose up -d
```

Use the `set` and `get` commands again to understand the behaviour of this mcrouter configuration.
How does the behaviour differ from the replicated cluster?

Again, you can look at the `docker-compose.yml` file. Read about [sharded mcrouter configuration](https://github.com/facebook/mcrouter/wiki/Sharded-pools-setup) 
and try to understand what the configuration is doing.

Again, you can bring down your cluster by running
```console
> docker-compose down.
```

## Build a Go program that can determine the topology of a memcached cluster

Now that you have seen how sharded and replicated mcrouter clusters work, write a Go program that can tell the difference between them.

Your program should take command-line flags like this:
```console
> go run cache-tester --mcrouter=11211 --memcacheds=11212,11213,11214
```

Your program should do `set` and `get`  commands against mcrouter and the memcached instances and print out whether the caches are 
operating in sharded or replicated mode. Test it against both configurations.

You can use the [Go memcache client](https://pkg.go.dev/github.com/bradfitz/gomemcache/memcache) to run commands against mcrouter and memcache from Go.
