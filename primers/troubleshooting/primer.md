<!--forhugo
+++
title="Troubleshooting Primer"   
+++
forhugo-->

# Troubleshooting Primer

## About this document {#about-this-document}

This document is a crash course in troubleshooting. It is aimed at people with some knowledge of computing and programming, but without significant professional experience in operating distributed software systems. This document does not aim to prepare readers to be oncall or be an incident responder. It aims primarily to describe the skills needed to make progress in day-to-day software operations work, which often involves a lot of troubleshooting in development and test environments.

## Learning objectives:

- Explain troubleshooting and how it differs from debugging
- Name common troubleshooting methods
- Experience working through some example scenarios.
- Use commonly used tools to troubleshoot

## Troubleshooting Versus Debugging {#troubleshooting-versus-debugging}

Troubleshooting is related to debugging, but isn’t the same thing.

Debugging relates to a codebase into which the debugger has full visibility. In general, the scope of the debugging is limited to a single program; in rare cases it might extend to libraries or dependent systems. In general, debugging techniques involve the use of debuggers, logging, and code inspection.

In troubleshooting we investigate and resolve problems in complex systems. The troubleshooter may not know which subsystem the fault originates with. The troubleshooter may not know the full set of relationships and dependencies in the system exhibiting the fault, and they may not be able to inspect all of the source code (or familiarity with the codebases involved).

We can use debugging to fix a system of limited complexity that we know completely. When many programs interact in a complex system that we can’t know completely, we must use troubleshooting.

Examples of troubleshooting:

- finding the cause of unexpectedly high load in a distributed system
- finding why packets are being dropped in a network between two hosts
- determining why a program is failing to run correctly on a host.

## General Troubleshooting Methods {#general-troubleshooting-methods}

Because troubleshooting involves a wide variety of systems, some of which may be unknown, we cannot create a comprehensive troubleshooting curriculum. Troubleshooting is not just learning a bank of answers but learning how to ask, and answer, questions methodically. Troubleshooters use documentation, tooling, logic, and discussion to analyse system behaviours and close in on problems. Being able to recognize gaps in your knowledge and to fill them is a key troubleshooting skill. Filling in these gaps becomes easier the wider your knowledge becomes over time.

There are some approaches to troubleshooting that are generally useful:

- Defining the problem
- Understanding the request path
- Bisecting the problem space
- Generating hypotheses and proving or disproving them - this step is iterative, in other words, you keep doing this until you have found the problem

### Defining the problem {#defining-the-problem}

To troubleshoot a problem, you must be able to understand it and describe it. This is a starting point for discussion with technical and non-technical stakeholders.

Begin your notes as you begin your investigation. Take lots of notes, and record each step you take in your journey towards understanding. If things get difficult, it’s valuable to have a record of what your theories were, how you tested them, and what happened. The first step is to make your first definition of the problem, which you will update very often.

This step is necessary because problem reports are often vague and poorly-defined. It’s unusual for _reports_ of a problem to accurately _diagnose_ the problem. More often, they describe symptoms. Where reports do attribute causes, this is often without systematic reasoning and might be completely wrong.

For example, a user might report that your website is slow. You verify this using your monitoring systems, which display graphs of page load times showing an increase in the past day. Based on your knowledge of the system, you know that slow requests can happen when the database is having issues. Using tooling (such as a [monitoring system](https://www.rabbitmq.com/monitoring.html), or[ logging slow queries](https://severalnines.com/blog/how-identify-postgresql-performance-issues-slow-queries/)) you can verify whether the database is indeed having performance issues, and from this you have your base definition of the problem. You can share this, and use it to build upon. We still do not know \_why \_the database latency is high, but it is a specific problem that we can investigate.

Often we do not know the actual cause of the problem, we just see the impact on other systems. It is like a rock thrown into a lake. The rock is out of sight, but we can still see its impact, as ripples on the water. The ripples help us guess the size and location of the rock in the lake.

- What do we already know about the issue?
- Can we describe the problem to others?
- Can we find proof that validates our assumptions?
- Can we reproduce the issue?
- What are the effects of the issue?
- Can we see the effects in a graph somewhere?
- Is it affecting all or only some operations?

It can be difficult to untangle cause and effect in computer systems. A database seeming slow may be an effect \_caused \_by additional load exhausting database resources. Another commonly-observed example is systems experiencing spikes in the number of incoming requests due to retries when there are problems with an application. If an application is very slow or serving errors, clients and users will naturally retry, creating unusually high load on the system. The additional load may cause further problems; but some other problem has triggered the application’s incorrect behaviour.

### Understanding the request path {#understanding-the-request-path}

Often the first challenge in any troubleshooting exercise is to figure out what is supposed to be happening in the system that is exhibiting a fault.

First figure out what is _meant_ to be happening, and then determine what is _actually_ happening.

Like when you order food, and the driver turns up at the wrong address. Where in the flow did things divert from what was expected?

- Was the wrong address provided?
- Was the address read incorrectly?
- Are there two identically named streets?
- Has GPS broken?

Making a mental model of the request path helps you navigate through the issue and narrow down the component(s) in the request path that might be having issues. It can be helpful to sketch an actual diagram on paper to help you build your mental model.

Depending on the situation, there may be documentation that can help you understand the intended operation of your system and build that mental model. Open-source software often comes with extensive manuals. In-house developed systems may also have documentation. However, often you need to use troubleshooting tools and experimentation to understand what is happening (there is a short catalogue of useful tools in a later section).

### Finding the cause of the problem: generating and validating hypotheses {#finding-the-cause-of-the-problem-generating-and-validating-hypotheses}

Once we determine and observe that there is a problem, we have supporting evidence to say that the [happy path](https://en.wikipedia.org/wiki/Happy_path) is not being taken. However, we still do not understand the cause, or how to fix it.

Once we believe that we know what might be causing the problem, we now have a hypothesis. We want to use our skills, or the skills of others, to verify that what we think is the problem, is actually the problem. With this we can mitigate or fix the problem.

#### Looking for causes, or “what happened earlier” {#looking-for-causes-or-“what-happened-earlier”}

When a problem has a distinct ‘starting time’, it can be worth checking for changes or events that occurred at that time. Not all problems are a consequence of recent events or changes - sometimes latent problems in systems that have existed for years can be triggered. However, changes (such as deployments and configuration updates) are the most common cause of problems in software systems, so recent changes and events should not be ignored.

For example, we have a graph that lines up with us getting alerted at 13:00, but we can see that the graph started getting worse at 12:00. If you leave a tap running, it isn’t a problem at first, but eventually the sink overflows. The alert is when the sink overflows, the _cause_ is leaving the tap running.

So when we have a better understanding of when things started to get worse, we need to see if anything changed around the same time.

- Did a deploy go out?
- Was a feature flag changed?
- Was there a sudden increase in traffic?
- Did a server vanish?

We can see that a deployment went out at 12:00. This feels like a hypothesis to dig into.

#### Examining possible causes {#examining-possible-causes}

So now we can find the changes that were included in that deployment. Either ourselves or a subject matter expert can help confirm that the changes might be related to the problem.

#### Testing solutions {#testing-solutions}

If we believe that the recent change might be the problem, then we may be able to revert the change. In stateless systems reverting recent changes is generally low-risk. Reverting changes to stateful systems must only be done with great care: an older codebase may not be able to read data written by a newer version of the system and data loss or corruption may result. However, in most production software systems, the great majority of changes that are made are safe to revert, and it is often the quickest way to fix a problem.

If reverting a recent change makes the problem go away, then it is a very strong signal that those changes were indeed the cause of the issue (although coincidences do occur). It does not explicitly prove or disprove our hypothesis that the recent change was the problem.

The next step is to examine those changes more closely and determine exactly how they caused the alert, which would definitively prove the hypothesis.

It is often easier to disprove a hypothesis than to prove it. For example, if one of the changes in the recent deployment (that we just rolled back) introduced a dependency on a new service, and we hypothesise that this new service’s performance might be the cause of the problem, we can inspect the monitoring dashboard for that microservice. If the monitoring shows that the new service is performing well (low latency and low error rates) then the hypothesis would be disproved. We can then generate new hypotheses and try to [falsify](https://en.wikipedia.org/wiki/Falsifiability) those.

### Finding problems by iteratively reducing the possibility space {#finding-problems-by-iteratively-reducing-the-possibility-space}

It is not always possible to track down problems by taking the fast path of looking for recent breaking changes. In this case, we need to use our knowledge of the system (including information that we gain by using debugging tools) to zero in on parts of our system that are not behaving as they should be.

_All swans are white vs no swan is black_

It is more efficient to find a way to disprove your hypothesis or falsify your proposition, if you can. This is because you only need to disprove something once to discard it, but you may _apparently_ verify a hypothesis many times in many different ways and still be wrong.

Every time we disprove a hypothesis we reduce the scope of our problem, which is helpful. We close in on something.

For example, let us imagine that we are troubleshooting an issue in a system that sends notifications to users at a specific times. We know that the system works as follows:

1. The user configures a notification using a web service
2. The notification configuration is written into a database
3. A scheduler reads from that database and then sends a push notification to the user’s device

In this situation, the problem could be in any of these steps. Perhaps:

1. The web service silently failed to write to the database for some reason
2. The database lost data due to a hardware failure
3. The scheduler had a bug and failed to operate correctly
4. The push notification could not be sent to the user for some reason; perhaps network-related

A good place to start would be by checking the scheduler program’s logs to determine whether it did attempt to send the notification or not. In order to do this, you would need some specifics about the request, such as a user ID, in order to identify the correct log lines.

Finding the log line (or not finding it) will tell you if the scheduler attempted to send the notification. If the scheduler did attempt to send the notification, then it eliminates database write failures, data loss, and some types of scheduler bug from your search, and indicates that you should focus on what happens on the part of the request path involving the scheduler sending the notification. You may find error details in the scheduler log lines.

If you do not find a log line for the specific request in question - and you do see log lines relating to other notifications around that same time - then it should signal you to examine the first three steps in the process and iterate. Checking whether the notification configuration is present in the database would further reduce the space of possible problems.

### USE Method {#use-method}

[Brendan Gregg’s USE](https://www.brendangregg.com/usemethod.html) (Utilisation, Saturation, Errors) method is particularly helpful for troubleshooting performance problems. Slowness and performance problems can be more difficult to debug than full breakage of some component, because when something is completely down (unresponsive and not serving requests) it’s generally more obvious than something being slow.

Performance problems are generally a result either of some component in the system being highly utilised or saturated, or of errors somewhere that are going unnoticed.

#### Bottlenecks {#bottlenecks}

Imagine a wine bottle and a wide tumbler of water, both holding the same amount of water. Turn the bottle and the glass upside down. The water in the glass falls at once. The water in the bottle empties more slowly. It sits above the narrow neck of the bottle; the speed of the pour is limited by the capacity of the bottle neck.

When a component - either physical, such as CPU or network, or logical, such as locks or cloud API quota - is too highly utilised, other parts of the system will end up waiting for the heavily-loaded component. This occurs because of queuing: when a system is under very heavy load, requests cannot usually be served quickly and on average, will have to wait. The heavily loaded component is known as a bottleneck.

The performance problem then tends to spread beyond the original bottleneck. The clients of the bottleneck system serve requests more slowly (because they are waiting for the bottleneck system), and this in turn affects clients of those systems. This is why heavy utilisation and saturation are very important signals in your systems.

## Scenarios {#scenarios}

### Slack Client Crashes {#slack-client-crashes}

Let’s look at [Slack's Secret STDERR Messages](https://brendangregg.com/blog/2021-08-27/slack-crashes-secret-stderr.html) by Brendan Gregg. This is a great piece not only because of Gregg’s expertise but because of how clearly Gregg describes his process.

Here, Gregg is attempting to troubleshoot why his Slack client seems to be crashing. He doesn’t know much about the inner workings of the software, and he doesn’t have the code, so he has to treat it as a black box. However, in this case, he does have a fairly clear probable location of the fault: the Slack program itself.

He knows that the program is exiting, so he starts there. He has a false start looking for a [core dump](https://linux-audit.com/understand-and-configure-core-dumps-work-on-linux/) to attach a debugger to, but this doesn’t work out so he moves on to try a tool called [exitsnoop](https://manpages.ubuntu.com/manpages/jammy/en/man8/exitsnoop-bpfcc.8.html). Exitsnoop is an eBBF-based tool that traces process termination. Gregg finds that the Slack client is exiting because it receives a SIGABRT signal (which is generally handled by termination).

Gregg still doesn’t know why the program is receiving this signal. He tries some more tracing tools - trying to get a stack trace - but draws a blank. He moves on to looking for logs instead.

He has no idea where Slack logs to, so he uses the lsof tool. Specifically, he runs

```
lsof -p `pgrep -n slack` | grep -i log
```

This command line does the following:

1. pgrep -n slack: find the process ID of the most recently started slack process.
2. `pgrep -n slack`: the use of backticks in a command line means ‘run the commands in the backticks first, and then substitute the result
3. lsof -p `pgrep -n slack`: this runs lsof -p &lt;PID of the most recently running slack process>. lsof with the ‘-p’ flag lists all the open files that belong to the given PID.
4. grep -i log: this searches for the word ‘log’ in the text that the command is given. The ‘-i’ flag just makes it case insensitive.
5. |: this is the pipe symbol. It takes the output of the commands to the left, and send it as input to the command on the right.

The overall command line searches for any file that the most recently started slack process has open, with the word ‘log’ (case insensitive) in the name.

This sort of command line is a result of the [UNIX tools philosophy](https://www.linuxtopia.org/online_books/gnu_linux_tools_guide/the-unix-tools-philosophy.html): commandline tools should be flexible and composable using mechanisms like the pipe (|). This also means that in a Linux/UNIX system, there are often many ways to achieve the same result. Likewise, Linux largely makes system state available to administrators. System state is all the state of the running kernel and its resources - such as which files are being held open by which process, which TCP/IP sockets are listening, which processes are running, and so forth. This state is generally exposed via the [/proc filesystem](https://www.kernel.org/doc/html/latest/filesystems/proc.html), as well as a variety of commandline tools.

Gregg doesn’t find any log files, but he realises that the logs might still exist, having been opened by another slack process. He tries a program called [pstree](https://man7.org/linux/man-pages/man1/pstree.1.html) which displays the entire tree of slack processes. It turns out that there are a number of slack processes, and he tries lsof again with the oldest. This time he finds log files.

Gregg is an expert troubleshooter, but we see that he attempts and abandons several methods for understanding the reason for this program’s failure. This is fairly typical of troubleshooting, unfortunately. He is also using his knowledge of core Linux concepts - processes, open files - to increase his knowledge about how the slack program works.

Once more, Gregg finds that the slack program logs don’t reveal the reason for the program’s crashing. However, he notices that the [stderr stream](https://www.howtogeek.com/435903/what-are-stdin-stdout-and-stderr-on-linux/) isn’t present in the log files.

Gregg knows of another tool, shellsnoop, that he uses to read the stderr stream.

Here he finds an error:

```
/snap/slack/42/usr/lib/x86_64-linux-gnu/gdk-pixbuf-2.0/2.10.0/loaders/libpixbufloader-png.so: cannot open shared object file: No such file or directory (gdk-pixbuf-error-quark, 5)
```

This error log indicates that the slack process tried to dynamically load a [shared code module](https://programs.wiki/wiki/linux-system-programming-static-and-dynamic-libraries.html) that didn’t exist. Resolving that issue resolves the failures.

Gregg concludes by pointing out that the tool opensnoop could have led him straight to this answer; but of course, that’s much easier to know in hindsight.

### Broken Load Balancers {#broken-load-balancers}

Let’s look at the [broken loadbalancer scenario](https://hostedgraphite1.wordpress.com/2018/03/01/spooky-action-at-a-distance-how-an-aws-outage-ate-our-load-balancer/) from Hosted Graphite.

On a Friday night, Hosted Graphite engineers start receiving reports of trouble from their monitoring indicating that their website was intermittently becoming unavailable. Their periodic checks to the website and graphs rendering API were both intermittently timing out.

They begin to gather more information to find out what might be happening. They note an

overall drop in all traffic across all their ingestion endpoints which suggested there might

be a network connectivity issue at play. When looking through the impact from canaries vs HTTP API endpoint traffic, they notice that canaries ingestion traffic is affected only in certain AWS regions but the HTTP API ingestion is affected regardless of the location. To add further confusion, some of their internal services also start to report timeouts.

There are conflicting information but all of the events indicate an AWS connectivity issue

and they decide to follow it through.

Digging further, they realise that the internal services having issues are relying on S3 (another AWS service) and that their AWS dependent integrations are also severely impacted. At this point AWS status page is reporting connectivity issues both in us-east-1 and us-west-2 region which is even more confusing as they cannot comprehend "_how_" an AWS outage could affect how they serve their website when it’s neither hosted on AWS nor physically located anywhere near the affected regions.

**One of the hardest problems during any incident is differentiating cause(s) from symptoms.**

So they start looking into the only service they are using which was hosted on AWS, [Route53 health checks](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/welcome-health-checks.html) for their own (self-hosted) load balancers. These Route53 health checks were configured to ensure that traffic was only routed to healthy load balancers to server production traffic and unhealthy ones were removed from the DNS entry. The health check logs indicate failures from many locations. They don’t know if this was a cause or a symptom, so they disable the route53 health checks, to either confirm or rule out that theory. Unfortunately, disabling the health checks didn’t resolve the issue so they continue digging for clues.

It is at this point, the AWS incident gets resolved and they notice that their traffic rate starts to recover with it, further confirming that the AWS outage was the trigger for this incident. Now they know what happened but not why.

They find two visualisations from their load balancing logs which help them paint a clear

picture of what happened. The first graph shows a sharp rise in the active connections through their load balancing tier, followed by a flat line exactly during the time period of most impact and then a decline towards the end of the incident. This explains the SSL handshake woes they noticed as any new connection won't be accepted by their load balancers once

the maximum connection limit was reached.

This still doesn't explain where these connections were originating from as they didn't see any increase in the number of requests received. This is where the second visualisation comes into picture. This graph shows the average time it took for hosts from different ISPs to send a full request over time and the top row represents requests from AWS hosts. These requests were taking up to 5 seconds to make a full request while other ISPs remained largely unaffected. At this point they finally crack the case.

The AWS hosts from the affected regions were experiencing connectivity issues which significantly slowed down their connections to the load balancers. As a result, these hosts were hogging all the available connections until they hit a connection limit in their load balancers, causing hosts in other locations to be unable to create a new connection.

## Tools {#tools}

There is no exhaustive list of troubleshooting tools. We use whatever is both available and best-suited for the problem at hand. Sometimes that’s a general-purpose tool – like Linux OS-level tooling or TCP packet dumps – and sometimes it’s something system-specific, like application-specific counters, or logging.

In many cases, Google is your friend. It is a starting point for understanding what an error message may mean and what may have caused the error. It is also a good way to find new observability tools (e.g. searching for things like ‘linux how to debug full disk’ will generally throw up some tutorials).

However: do not spend too long trying to find information in this way. Google is your friend, but it is not your only friend. It is easy to fall into a rabbit hole and lose your entire day to Googling, so give yourself some time limits. Set an alarm for 90 minutes. Write up your journey with the problem so far and take it to a more senior engineer. They will appreciate the work you have already done on describing and exploring the problem.

System-specific monitoring can help you: does the system you are investigating export metrics (statistics exposed by a software program for purposes of analysis and alerting)? Does it have a status page, or a command-line tool that you can use to understand its state?

Examples:

- [https://www.haproxy.com/blog/exploring-the-haproxy-stats-page/](https://www.haproxy.com/blog/exploring-the-haproxy-stats-page/)
- [https://vitess.io/docs/13.0/user-guides/configuration-basic/monitoring/](https://vitess.io/docs/13.0/user-guides/configuration-basic/monitoring/)
- You many have system-specific dashboards built in Grafana, Prometheus, a SaaS provider like Honeycomb or Datadog, or another monitoring tool

Loadbalancers and datastore statistics are particularly useful - these can often help you to determine whether problems exist upstream (towards backends or servers) or downstream (towards frontends, or clients) of your loadbalancers/datastores. Understanding this can help you narrow down the scope of your investigations.

Logs are also a great source of information, although they can also be a source of plentiful red herrings. Many organisations use a tool such as Splunk or ELK to centralise logs for searching and analysis.

There is a fairly long list of tools below. You don’t need to be an expert in all of these, but it is worth knowing that they exist, what they do in broad strokes, and where to get more information (man pages, google). Try running all of these

You should be familiar with basic Linux tooling such as:

- [perf](https://www.brendangregg.com/perf.html)
- [strace](https://man7.org/linux/man-pages/man1/strace.1.html)
- ltrace
- top, htop
- sar
- netstat
- lsof
- kill
- df, du, iotop
- ps, pstree
- the[ /proc/](https://www.kernel.org/doc/html/latest/filesystems/proc.html) filesystem
- dmesg, location of system logfiles - generally /var/syslog/, journalctl
- tools like cat, less, grep, sed, and awk are invaluable for working with lengthy logs or text output from tools
- jq is useful for parsing and formatting JSON

For debugging network or connectivity issues, you should know tools like:

- dig (for DNS)
- traceroute
- tcpdump and wireshark
- netcat

[curl](https://curl.se/docs/manpage.html) is invaluable for reproducing and understanding problems with HTTP servers, including issues with certificates.

[Man pages](https://www.kernel.org/doc/man-pages/) can be super useful when you are looking for more information and available options for any of the common Linux tools.

[eBPF ](https://ebpf.io/)is a newer technology that lets you insert traces into the Linux OS, making everything visible. A lot of observability tools use eBPF under the hood, or you can write your own eBPF programs.

# Related Reading {#related-reading}

[https://github.com/iovisor/bcc/blob/master/docs/tutorial.md](https://github.com/iovisor/bcc/blob/master/docs/tutorial.md)

[https://netflixtechblog.com/linux-performance-analysis-in-60-000-milliseconds-accc10403c55](https://netflixtechblog.com/linux-performance-analysis-in-60-000-milliseconds-accc10403c55)

[https://linuxtect.com/linux-strace-command-tutorial/](https://linuxtect.com/linux-strace-command-tutorial/)

[https://www.brendangregg.com/perf.html](https://www.brendangregg.com/perf.html)

[https://danielmiessler.com/study/tcpdump/](https://danielmiessler.com/study/tcpdump/)

[https://www.lifewire.com/wireshark-tutorial-4143298](https://www.lifewire.com/wireshark-tutorial-4143298)
