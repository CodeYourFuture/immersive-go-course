+++
title="CYF+02 Sprint 3"
date="01 Jan 2024"    
versions=["1-1-0"]
weight=4
+++

# Provisional start date: 6 May 2024

## Consolidation

- [ ] [Multiple servers](../../projects/multiple-servers) Revisit this project and deploy it with Docker and GH Actions

## Managing state at scale

- [ ] Read about [Distributed Systems: Scaling Stateless Services](../../primers/distributed-software-systems-architecture/scaling-stateless-services)
- [ ] Do this HashiCorp tutorial. It will give you hands-on experience with seeing health checks can be used to manage failure. https://learn.hashicorp.com/tutorials/consul/introduction-chaos-engineering?in=consul/resiliency
- [ ] Demonstrate circuit breaking with this HashiCorp tutorial: https://learn.hashicorp.com/tutorials/consul/service-mesh-circuit-breaking?in=consul/resiliency
- [ ] See different kinds of load-balancing algorithms in use with this HashiCorp tutorial: https://learn.hashicorp.com/tutorials/consul/load-balancing-envoy?in=consul/resiliency
- [ ] Do this tutorial which demonstrates autoscaling with minikube (you will need to install minikube on your computer if you don’t have it). It will give you hands-on experience in configuring autoscaling, plus some exposure to Kubernetes configuration.  
      You may need to run this command first: minikube addons enable metrics-server
      https://learndevops.novalagung.com/kubernetes-minikube-deployment-service-horizontal-autoscale.html

## Command line familiarity

- [ ] [Linux Bash I](https://www.bogotobogo.com/Linux/linux_tips2_bash.php) & [Linux Bash II](https://www.bogotobogo.com/Linux/linux_bash_2.php) (There are other blog posts in the same series that could be useful to go through in spare time, would recommend running the commands that the blog post suggests for better understanding.)

## Troubleshooting

- [ ] [Troubleshooting Primer](../../primers/troubleshooting/)
- [ ] Troubleshooting project #3
    - This project is designed to get you familiar with upstream service failure.
    - You'll be given 2 instances today. \o/
      A www host (www=webserver), and a db host (db=database).
    - To log in: `ssh -i </path/to/the/ssh-private-key> <username>@<IP>`
        - You have sudo access on both hosts, please give a shout if that doesn't work. (You'll need it.)
    - The goal of the exercise is:
        - When you run `lynx -dump http://localhost/`, you will see a cute image of a cat on your terminal.
        - (Unlike in the previous exercise, today we don't care that the database password is exposed in `ps`, so don't worry about that.)
        - Along the way, we expect you'll be able to answer:
            - What kind of failures do you see when talking to an upstream service?
              (Do you know about upstream and downstream services? It's useful jargon.)
            - The webserver doesn't start. Can you explain why?
    - Some knowledge to get you started:
        - `httpurr` makes a return. Code has changed a bit from previous exercise, if you need it, source code is in `/httpurr-source/` on the www host.
        - Unlike previous time, we don't use docker this time.
        - The database is postgresql. (Unlike previous time, when it was mysql.) It's a lovely database.
    - While doing the exercise, I would recommend logging what you do.

## Product

You're starting to learn efficiencies and optimisations. How will you apply them to your product work?
