+++
title="CYF+02 Sprint 2"
date="01 Jan 2024"    
versions=["1-1-0"]
hero="dddepth--085.webp"
weight=3
+++

# Provisional start date: 22 Apr 2024

## Study

- [ ] [Distributed Systems: State](../../primers/distributed-software-systems-architecture/state)
- [ ] [Troubleshooting Primer](../../primers/troubleshooting/)

## Projects

- [ ] [Batch processing](../../projects/batch-processing)
- [ ] [Buggy app](../../projects/buggy-app)
- [ ] [Memcached Clusters](../../projects/memcached-clusters)
- [ ] Troubleshooting project #2
    - This project is designed to get you familiar with docker.
    - To log in: `ssh -i </path/to/the/ssh-private-key> <username>@<IP>`
        - You have sudo access on the host, please give a shout if that doesn't work. (You'll need it.)
        - There's a database (mysql) on the machine. If you have a mysql client, you can use `--host=127.0.0.1 --port=3306 --user=root --password` to log in.
          (Ask any of the instructors for the password.)
    - The goal of the exercise is:
        - when you run `lynx -dump http://localhost/`, you will see a cute image of a cat on your terminal.
          This shouldn't take more than a day or 2.
        - secondary, once you've reached that, we'll ask you to hide the password that you can see in `ps auxww | grep miauw`.
          (A password visible in `ps` is bad for security, so we want that to be somewhat hidden.)
        - Along the way, we expect you'll be able to answer:
            - How is systemd configured? Where are the logs?
            - Can you describe what docker does?
            - Can you describe the two ways that docker volumes are used in this setup?
            - Docker uses images. What are images? And how do you create your own?
    - Some knowledge to get you started:
        - The webserver is called `httpurr`, with a corresponding systemd service (`httpurr`).
          It's written in go, and you can see the source code in `/httpurr/`.
        - The database is mysql, with a corresponding systemd service (`httpurr-db`).
        - Both are run inside docker.
    - While doing this exercise, I would recommend logging what you do.

## Product

Your product work is ramping up, and so is the complexity of your study. How will you manage your time?
