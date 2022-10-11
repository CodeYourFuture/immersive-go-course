# GRPC Client and Server Communication 

The goals of this project are: 
 * to use gRPC requests, and understand how they differ from REST
 * to use deadlines with gRPCs
 * to build a simple prober and understand how this might be used in production
 * optionally, gain experience with instrumenting an application and running Prometheus

## What is gRPC?

RPCs aren’t the same as HTTP requests. HTTP is a text-based protocol, oriented towards resources or entities. RPCs are oriented towards calling a specific function. RPCs are highly structured; normally the request and response are encoded in efficient binary formats which are not human-readable. HTTP APIs are usually limited to CRUD (Create, Read, Update, Delete) operations and RPCs can perform any kind of operation. If you want to efficiently integrate two systems that you control, RPCs are a good choice. If you want to provide an API for use by developers outside of your organisation then HTTP APIs are generally a better choice, because they are simpler to develop against and all programming languages provide good HTTP support. gRPC, which we will use in this exercise, is a RPC implementation that is fairly common in industry. 

Read the [gRPC Introduction](https://grpc.io/docs/what-is-grpc/introduction/) and 
[gRPC Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/) for an overview of gRPC.
todo stuff from primer

## Get started with a HelloWorld gRPC example

Start by building a simple gRPC ‘hello world’ server and cli client in golang. See [gRPC Quickstart](https://grpc.io/docs/languages/go/quickstart/) for detailed instructions.

## Change HelloWorld to Prober

Once you have the 'hello world' example from the gRPC Quickstart tutorial working, we will move on to implementing a prober.

Rename `greeter_client` and `greeter_server` to `prober_client` and `prober_server`.   
Change the `helloworld` subdirectory to be called `prober`,` and helloworld.proto` to `prober.proto`.
Delete the genereted .go files in the same directory as `prober.proto`.

Now update the the prober.proto file so that the main content section looks like this:

```
package prober;

// The greeting service definition.
service Prober {
  // Sends a greeting
  rpc DoProbes (ProbeRequest) returns (ProbeReply) {}
}

// The request message containing the user's name.
message ProbeRequest {
  string name = 1;
}

// The response message containing the greetings
message ProbeReply {
  string message = 1;
}
```

Regenerate the generated protobuf code: 
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative prober/prober.proto
```

Observe the new generated files:
```
ls prober
prober.pb.go		prober.proto		prober_grpc.pb.go
```

Read through `prober_grpc.pb.go` - this is the interface you will use in your code. This is how gRPC works: a `proto` format 
gets generated into language-specific code that you can use to interact with gRPCs. If you're working with multiple programming languages
that need to interact through RPCs, you can do this by generating from the same protocol buffer definition with the language-specific
tooling. 

Now, update the client and server `main.go` files to use the changed package, service and message names.
Get these running again as they did originally. The functionality will still be the same for now, just the hello world interaction.

```
go run prober_server/main.go
2022/10/09 20:20:55 server listening at [::]:50051
2022/10/09 20:21:24 Received: world
````

```
go run prober_client/main.go 
2022/10/09 20:21:24 Greeting: Hello world
```

While we haven't yet changed what the code does, we've gained some experience with the protocol buffer format, learned 
how to generate Go code from protocol buffer definitions, and call that code from Go programs.

Note: in a real production system you can't rename protocol buffers as we've done here: you would run into problems because your old clients can't talk to your new servers, and vice-versa. Since we gradually release code, we would have to support both services until all old clients were turned down.

## Implement prober logic

Let's modify the prober service slightly. Instead of the simple HelloWorld greeting, we are going to implement a service that will probe a HTTP endpoint N times and return the average time to fetch that endpoint to the client. 

This tool could be useful in a real operational context. Imagine that we want to verify that our site is available and has acceptable latency from many different locations around the world. We might run several instances of our prober server in
different regions. 

Change your prober request and response to look like this and regenerate your generated go code.
```
// The request message
message ProbeRequest {
  string endpoint = 1;    // Endpoint to probe
  int32 repetitions = 2;  // How many probes to execute
}

// The response message 
message ProbeReply {
  float average = 1;    // Average response time (in milliseconds)
}
```

Update your client to read the endpoint and number of repetitions from the [command line](https://gobyexample.com/command-line-arguments).
Then update your server to execute the probe: do a HTTP fetch of `endpoint` the specified number of times.
Use the standard [`net/http` package](https://pkg.go.dev/net/http).

Time each fetch using the `time` package:

```
  	start := time.Now()
	resp, _:= http.Get(endpoint)
    elapsed := time.Since(start)
```

Add up all the elapsed times, divide by the number of repetitions, and return the average to the client.
The client should print out the average value received.
You can do arithmetic operations like addition and division on `time.Duration` values in Go.
However, under the hood `time.Duration` is an `int64`. To divide by an `int32` and return a `float` will need to do some [type conversions](https://go.dev/tour/basics/13).

## Add a client timeout

Maybe the site you are probing is very slow, or the number of repetitions is very high.
Either way, you don't want your program to wait forever. 

On the client side, add a [timeout](https://pkg.go.dev/context#WithTimeout) to stop waiting after 1 second].
```

Run your client against some website - how many repetitions do you need to see your client timeout?

## Extra Challenges: Errors, Continuous Probing, Serve Prometheus Metrics

These sections are optional - do them for an extra challenge if time permits.

### Part 1:
How do we know if the HTTP fetch succeeded at the server? Add a check to make sure it did.
How should we deal with errors, i.e. if the endpoint isn't found?
Modify your code and proto format to handle these cases.

### Part 2: Run probes on a schedule
Instead of getting probe times via the gRPC client we might consider running the endpoint probes on a schedule. Change the
prober_server to accept a frequency argument instead of the old `repetitions` parameter, 
and continue running probes every N seconds forever. 

You will need to use a [goroutine](https://gobyexample.com/goroutines) to keep running probes in the background.
Now, add a new client operation to get the value of the last probe run.

### Part 3: Add Prometheus Metrics
If you still have time, then let's learn something about how to monitor applications.

In software operations, we want to know what our software is doing and how it is performing.
One very useful technique is to have your program export metrics. Metrics are basically values that your 
program makes available (via HTTP). Specialised programs, such as Prometheus, can then fetch metrics regularly
from all the running instances of your program, store the history of these metrics, and do useful arithmetic on them
(like computing rates, averages, and maximums). You can use this data to do troubleshooting and to alert if things 
go wrong.

Read the [Overview of Prometheus](https://prometheus.io/docs/introduction/overview/).

Now add Prometheus metrics to your prober server. Every time you execute a probe, update a `gauge` metric that tracks the latency.
Add a `label` specifying the endpoint being probed. 
The [Prometheus Guide to Instrumenting a Go Application](https://prometheus.io/docs/guides/go-application/) has all the information you need to do this.

Once you've run your program, use your client to execute periodic probes against some endpoint. 
Now use the `curl` program or your browser to view the metrics. 
You should see a number of built-in Go metrics, plus your new gauge.

If you use your client to start probing a second endpoint, you should see a second labelled metric appear.

### Part 3: Scrape Prometheus Metrics
The final step is to set up the Prometheus application to periodically pull metrics from your `prober_server`.

Run Prometheus in Docker locally - see [Prometheus Installation](https://prometheus.io/docs/prometheus/latest/installation/).
You'll need to set up a simple configuration to scrape your `prober_server`. Prometheus [configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) is quite complex but you can adapt this [example configuration](https://github.com/prometheus/prometheus/blob/main/documentation/examples/prometheus.yml).

Next, find your custom gauge metric from your `prober_server` in http://localhost:9090/metrics.
Graph it in http://localhost:9090/graph.