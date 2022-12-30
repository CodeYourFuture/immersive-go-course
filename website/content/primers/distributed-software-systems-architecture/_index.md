+++
title="The Distributed Software Systems Architecture Primer"
author="Laura Nolan"
date="28 Dec 2022 12:22:11 BST"    
+++

## About this primer {#about-this-document}

This document outlines a short course in distributed systems architecture patterns and common issues and considerations. It is aimed at students with some knowledge of computing and programming, but without significant professional experience in building and operating distributed software systems.

### Learning outcomes for this primer:

- Explain the reasons that we build distributed systems
- Demonstrate how to build your systems to deal with common types of failure, such as slow backends and network partitions
- Describe the pros and cons of asynchronous, semisynchronous, and synchronous database replication
- Describe the difference between asynchronous work systems (like pipelines) and serving systems

There are five sections in this primer, one per sprint.

## Why Distributed Systems? {#why-distributed-systems}

All client-server systems are distributed systems. Any computer system that involves communication between multiple physical computers is a distributed system. Any system that separates data storage from web serving, or which uses cloud services via APIs, is a distributed system.

A lot of distributed systems in operation go far beyond these simple architectures. Organisations build distributed systems to make their systems more reliable in the face of failure. A distributed system can serve a user from any one of several physical computers, perhaps on different networking and power domains.

Some workloads are too high to be served from a single machine. We use distributed systems techniques to spread workloads across many machines This is called horizontal scalability. In distributed systems you can serve requests closer to your users, which is faster, and a better user experience (it is remarkable how different many web applications feel when accessed from Australia or South Africa).

Almost all computer systems that we build today are distributed systems, whether large or small.
