<!--forhugo
+++
title="5. Distributed Locking and Distributed Consensus"
+++
forhugo-->

# 5

## Distributed Locking and Distributed Consensus

> In a program, sometimes we need to lock a resource. Think about a physical device, like a printer: we only want one program to print at a time. Locking applies to lots of other kinds of resources too, often when we need to update multiple pieces of data in consistent ways in a multi-threaded context.

### How to do distributed locking

We need to be able to do locking in distributed systems as well. Read Martin Kleppmann’s article [How to do distributed locking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html).

### What are the two main reasons to do distributed locking?

It turns out that in terms of computer science, distributed locking is theoretically equivalent to reliably electing a single leader (like in database replication, for example). It is also logically the same as determining the exact order that a sequence of events on different machines in a distributed system occur. All of these are useful.

### RAFT

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

See [raft-otel](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/raft-otel) project.
This project is an opportunity to explore the RAFT distributed consensus protocol through the medium of distributed tracing.

## Further Optional Reading {#further-optional-reading}

Discuss any of these pieces at office hours with Laura.

[Dan Luu’s list of postmortems](https://github.com/danluu/post-mortems)
: includes a lot of interesting stories of real-world distributed systems failure

[Jeff Hodges' Notes on Distributed Systems for Young Bloods](https://www.somethingsimilar.com/2013/01/14/notes-on-distributed-systems-for-young-bloods/)
: is a very practical take on distributed systems topics.

[Alvaro Videla's blog post and talks about learning distributed systems](https://alvaro-videla.com/2015/12/learning-about-distributed-systems.html)
: are accessible and well organised.

[Marc Brooker’s blog](https://brooker.co.za/blog/)
: is full of interesting pieces, which are very approachable.

[The SRE Book](https://sre.google/sre-book/table-of-contents/)
: is available in full online and is worth reading - it addresses many aspects of operating distributed software systems.

Aphyr (Kyle Kingsbury)
: provides detailed [notes for a distributed systems class he teaches.](https://github.com/aphyr/distsys-class)

[A Distributed Systems Reading List](https://dancres.github.io/Pages/)
: will keep you reading excellent distributed systems papers for many many months.
