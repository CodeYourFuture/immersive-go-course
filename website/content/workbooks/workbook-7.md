+++
title="CYF+02 Sprint 1"
date="01 Jan 2024"    
versions=["1-1-0"]
hero="pictures/depths/dddepth--085.webp"
weight=2
+++

# Provisional start date: 08 April 2024

In the CYF+ process your goal is to demonstrate the profile of an engineer at a high performing company. This series of challenges was designed by a panel of engineers from top companies to help you demonstrate that profile. Key skills include:

- **Problem solving** - You can solve problems in a systematic way, and you can explain your approach to others.
- **Systems thinking** - You can trace the flow of data through a system, and you can explain how the system works to others.
- **Communication** - You can communicate effectively with others, both verbally and in writing. This doesn't mean knowing all the answers; it means being able to research things you don't know, and being able to ask for help when you need it.
- **Collaboration** - You can work effectively with others, and you can learn from them.
- **Self-management** - You can manage your own time and priorities, and you can learn independently.

## Study

- [ ] [Distributed Systems: Reliable RPCs](../../primers/distributed-software-systems-architecture/reliable-rpcs)
- [ ] [Troubleshooting Primer](../../primers/troubleshooting/)
- [ ] [Linux process Intro](https://tldp.org/LDP/tlk/kernel/processes.html)
- [ ] [Linux Process and Signals](https://www.bogotobogo.com/Linux/linux_process_and_signals.php) - This doc has some commands to give you an insight on how to view processes and pass signals to a process; we recommend running those commands and documenting your learning.
- [ ] [PHP fastcgi](https://www.php.net/manual/en/install.fpm.php)
- [ ] [Systemctl](https://www.freedesktop.org/software/systemd/man/systemctl.html) and then [this](https://www.redhat.com/sysadmin/linux-systemctl-manage-services)
- [ ] [Nginx](https://nginx.org/en/docs/) and [How request processing works in Nginx](https://nginx.org/en/docs/http/request_processing.html)
- [ ] [Nginx with PHP FastCGI](https://www.nginx.com/resources/wiki/start/topics/examples/phpfcgi/)

## Projects

- [ ] [Servers & databases](../../projects/server-database)
- [ ] [Multiple Servers](../../projects/multiple-servers)
- [ ] [Docker & Cloud Deployment](../../projects/docker-cloud)
- [ ] [Troubleshooting project #1 - Fix Nginx]
    - This exercise is designed to help you learn how to setup Ngix proxy and configure it with PHP FastCGI to disaplay web pages.
    - Instruction: (please let Radha Kumari know when you are ready to do this exercise)
        - `ssh -i </path/to/the/ssh-private-key> <username>@<IP>`
        - `ssh -i </path/to/the/ssh-private-key> <username>@<IP>`
    - The goal of this exercise is - when you run "curl http://127.0.0.1/" in the terminal you get "Hello World" with HTTP response code 200. This exercise is a mix of both system and application troubleshooting
        - Run the curl command in verbose mode and see what the output looks like.
        - Identify what is the port that the curl defaults to and what runs on that port.
        - Find where Nginx configuration lives and what each configuration field means along with the location on the file system where it serves the content from.
        - What is php fast CGI and php-fpm (fastCGI process manager)?
        - Find out where php-fpm configuration files live and what each configuration field means.
    - This exercise shouldn't take more than 1-2 days. If you are stuck, please ask for help. While you are doing the exercise, I would really recommend logging somewhere the commands you run and noting what information that gives you and whether or not that was helpful in reaching the end goal. If you were given Sudo access, make sure to document when you needed to use that and why.

## Product

You will join [CYF Products](https://codeyourfuture.io/volunteers/) as a junior engineer. You will be deployed on a team delivering a product with real users, stakeholders and deadlines. Your challenge is to work with your team to deliver a product that meets the needs of your users, while managing the competing demands of stakeholders, deadlines, and your own learning priorities. Communication, organisation, and collaboration are key skills here.
