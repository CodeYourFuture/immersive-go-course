# Troubleshooting Project #2

## Problem definition

When running a request to the server on the default port 80, the connection is aborted with an unexpected network read error. Below is an example of the error encountered when using lynx to make an HTTP request to http://localhost:

```zsh
[saadia@ip-172-31-30-87 ]$ lynx -dump http://localhost/

Looking up localhost
Making HTTP connection to localhost
Sending HTTP request.
HTTP request sent; waiting for response.
Alert!: Unexpected network read error; connection aborted.
Can't Access `http://localhost/'
Alert!: Unable to access document.
lynx: Can't access startfile
```

or using curl in verbose mode :

```zsh
curl -v http://localhost/
* Host localhost:80 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:80...
* Connected to localhost (::1) port 80
> GET / HTTP/1.1
> Host: localhost
> User-Agent: curl/8.5.0
> Accept: */*
>
* Recv failure: Connection reset by peer
* Closing connection
curl: (56) Recv failure: Connection reset by peer

```

**Expected response**

Image of a cat on the terminal.

## Understanding the system

We know from the instructions that :

- The webserver is called httpurr, with a corresponding systemd service (httpurr).
- It’s written in go, and you can see the source code in /httpurr/.
- The database is mysql, with a corresponding systemd service (httpurr-db).
- Both are running inside docker.

## My journey in troubleshooting #2

- At the beginning I was trying to find any error message related to my problem.
- I investigated the `journalctl` logs :

```zsh
[saadia@ip-172-31-30-87 ~]$ journalctl | grep lynx :
Apr 25 16:11:20 ip-172-31-91-156.ec2.internal sudo[14211]: ec2-user : TTY=pts/1 ; PWD=/home/ec2-user ; USER=root ; COMMAND=/usr/bin/yum install lynx
```

- Well it's not right , `journalctl` is used for systemd services logs but lynx is not a service
- I tried to use `strace  lynx -dump http://localhost/` to trace all the syscall made while sending this request.
- But I got a long logs difficult to read from them. I limited the trace to count syscall using `strace -c lynx -dump http://localhost/`

```zsh
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
  ...
  0.00    0.000000           0         1         1 access
  0.00    0.000000           0         2         1 connect
  0.00    0.000000           0         3         1 wait4
  0.00    0.000000           0         2         1 arch_prctl
  0.00    0.000000           0        46        18 openat
  0.00    0.000000           0        32         2 newfstatat
  ...
------ ----------- ----------- --------- --------- ----------------
100.00    0.000000           0       283        24 total
```

- I got 24 different error. Actually, it's difficult to try to understand every syscall error and try find the bug in the request!
- I went back to investigate some part of my system:
- webserver container if it's running correctly

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED       STATUS       PORTS                                         NAMES
2301db300652   httpurr   "/httpurr --dbhost=h…"   3 hours ago   Up 3 hours   0.0.0.0:80->80/tcp, :::80->80/tcp, 8080/tcp   httpurr
```

- List of available images :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker image ls
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
httpurr      latest    ca2e5b2d7687   3 hours ago   13.3MB
mysql        latest    e9387c13ed83   13 days ago   578MB
<none>       <none>    fb79d4b9cc13   2 weeks ago   13.3MB
<none>       <none>    c365962c2633   2 weeks ago   13.3MB
<none>       <none>    a9eee9e53ccb   3 weeks ago   12.6MB
<none>       <none>    f24b813437bf   3 weeks ago   12.6MB
<none>       <none>    b333d8b6dd34   3 weeks ago   12.6MB
mysql        <none>    6f343283ab56   7 weeks ago   632MB
```

- I could't see any issue.
- I explored the source code on `httpurr` :

```zsh
[saadia@ip-172-31-30-87 ~]$ cd /httpurr/
Dockerfile  README.md  catsdb  go.mod  go.sum  httpurr.go  sql  vendor
```

- I read trough the main package `httpturr.go`, nothing seems odd.
- I tried run the code locally,

```zsh
[saadia@ip-172-31-30-87 ~]$ go run .
-bash: go: command not found
```

- I needed to install golang : `sudo yum install golang -y`
- I run the code again :

```zsh
[saadia@ip-172-31-30-87 ~]$ go run .
slack-github.com/avandersteldt/httpurr/catsdb: cannot compile Go 1.22 code
```

- I need to investigate other files, it’s related to the compatibility of the go version with other modules

```zsh
[saadia@ip-172-31-30-87 ~]$ cat go.mod
module slack-github.com/avandersteldt/httpurr
go 1.22
require github.com/go-sql-driver/mysql v1.8.1
require filippo.io/edwards25519 v1.1.0 // indirect
```

- I googled compatibility `go 1.22` with `go-sql-driver/mysql v1.8.1 (latest)`
- From the `go-sql-driver/mysql` docs, it’s compatible : `Go 1.20 or higher. We aim to support the 3 latest versions of Go.`
- But I changed the go version on go.mod to `go 1.20` and tested it.

```zsh
[saadia@ip-172-31-30-87 ~]$  go run .
2024/05/14 16:34:01 HTTPurr starting...
```

- The server is running.
- I needed to change the version in the Dockerfile as well : `FROM golang:1.20 AS build`
- I went to a totally wrong path trying to build the image and create a new container and run into a `no space left on device` problem and had to do docker clean up to claim space `sudo docker system prune -a`. You can read about it on the section [wrong path](#wrong-path).
- I should have done just a restart to the `httpurr` service : `sudo systemctl restart httpurr`
- Also I discovered later that the compatibility information are on the `vendor/module.txt`
- I send a request again and I got the same error.
- I decided to check if the server is running properly in the container

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED              STATUS              PORTS                                         NAMES
7ed55ffa8670   httpurr   "/httpurr --dbhost=h…"   About a minute ago   Up About a minute   0.0.0.0:80->80/tcp, :::80->80/tcp, 8080/tcp   httpurr
[saadia@ip-172-31-30-87 ~]$ sudo docker container logs 7ed55ffa8670
2024/05/15 12:04:56 HTTPurr starting...
```

- The server is running.
- I checked the docker stats to check if the container stop ant any time :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker stats
CONTAINER ID   NAME         CPU %     MEM USAGE / LIMIT     MEM %     NET I/O         BLOCK I/O     PIDS
7ed55ffa8670   httpurr      0.00%     1.043MiB / 949.6MiB   0.11%     1.39kB / 300B   69.6kB / 0B   4
94ad2d8e784d   httpurr-db   317.08%   211.7MiB / 949.6MiB   22.29%    586B / 0B       0B / 0B       1
```

- I noticed that the `httpurr-db` container keeps starting and stopping. When was this created ? Actually I have deleted all containers and created just an image via docker build, Why I have also the containers created ?
- After googling a bit , I discover that we can configure systemd to automatically create and manage Docker containers as services.
- I remembered seeing that in the instruction :) actually we have `httpurr` service and `httpurr-db` service.
  Let’s get more details about those services.

1. httpurr service

```zsh
[saadia@ip-172-31-30-87 ~]$ cat /etc/systemd/system/httpurr.service
[Unit]
Description=HTTPurr
After=docker.service
Requires=docker.service
Wants=httpurr-db.service
[Service]
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker stop httpurr
ExecStartPre=-/usr/bin/docker rm httpurr
ExecStartPre=/usr/bin/docker build -t httpurr /httpurr
ExecStart=/usr/bin/docker run --name httpurr --user 1001:1001 --link httpurr-db --rm -p 80:80/tcp httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr --dbpass=miauw --dbname=httpurr
[Install]
WantedBy=multi-user.target
```

- So we know that `httpturr`
  - Systemd service corresponding to the webserver
  - Requires docker
  - Have `httpurr-db.service` as a dependency (but it's a weak independency as the `httpurr` service can start even if `httpurr-db` fails)
  - Build the image `httpurr` defined on the `/httpurr/Dockerfile`
  - Create the container `httpurr` which set to publish the container's port 80 to the host port 80

2. httpurr-db service

```zsh
[saadia@ip-172-31-30-87 ~]$ cat /etc/systemd/system/httpurr-db.service
[Unit]
Description=HTTPurr-DB
After=docker.service
Requires=docker.service
[Service]
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker stop httpurr-db
ExecStartPre=-/usr/bin/docker rm httpurr-db
ExecStartPre=/usr/bin/docker pull mysql:latest
ExecStart=/usr/bin/docker run --name httpurr-db --rm -v /database:/var/lib/mysql:ro --user 1002:1002 -p 3306:3306/tcp mysql:latest
[Install]
WantedBy=multi-user.target
```

- So we know that `httpurr-db`

  - Systemd service corresponding to the mysql database
  - Requires docker
  - Build mysql image pulled from the Docker Hub
  - Create the container `httpurr-db` which set :
  - to mount a volume from `/database` on the host to `/var/lib/mysql` on docker with `ro` permissions
  - to publish the container's port 3306 to the host port 3306

- After learning more about those services :
  1. I restarted and checked the status of `httpurr` service:

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo systemctl restart httpurr
[saadia@ip-172-31-30-87 ~]$ sudo systemctl status httpurr
● httpurr.service - HTTPurr
     Loaded: loaded (/etc/systemd/system/httpurr.service; enabled; preset: disabled)
     Active: active (running) since Wed 2024-05-15 13:13:29 UTC; 8s ago
    Process: 639846 ExecStartPre=/usr/bin/docker stop httpurr (code=exited, status=1/FAILURE)
    Process: 639851 ExecStartPre=/usr/bin/docker rm httpurr (code=exited, status=1/FAILURE)
    Process: 639861 ExecStartPre=/usr/bin/docker build -t httpurr /httpurr (code=exited, status=0/SUCCESS)
   Main PID: 640157 (docker)
      Tasks: 5 (limit: 1114)
     Memory: 20.2M
        CPU: 227ms
     CGroup: /system.slice/httpurr.service
             └─640157 /usr/bin/docker run --name httpurr --user 1001:1001 --link httpurr-db --rm -p80:80/tcp httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr -->
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #10 CACHED
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #11 [httpurr 3/4] RUN adduser -D nonroot
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #11 CACHED
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #12 exporting to image
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #12 exporting layers done
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #12 writing image sha256:5bb3674934e91fd9b837436644034ef7b6fae09ea81bad48fadbefb35d9679b2 done
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #12 naming to docker.io/library/httpurr done
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[639885]: #12 DONE 0.0s
May 15 13:13:29 ip-172-31-30-87.ec2.internal systemd[1]: Started httpurr.service - HTTPurr.
May 15 13:13:29 ip-172-31-30-87.ec2.internal docker[640157]: 2024/05/15 13:13:29 HTTPurr starting...
```

=> We confirmed that the server container is running without error.

2.  I checked the status of `httpurr-db` service:

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo systemctl status httpurr-db
× httpurr-db.service - HTTPurr-DB
     Loaded: loaded (/etc/systemd/system/httpurr-db.service; disabled; preset: disabled)
     Active: failed (Result: exit-code) since Wed 2024-05-15 13:13:58 UTC; 3min 57s ago
   Duration: 1.369s
    Process: 644165 ExecStartPre=/usr/bin/docker stop httpurr-db (code=exited, status=1/FAILURE)
    Process: 644169 ExecStartPre=/usr/bin/docker rm httpurr-db (code=exited, status=1/FAILURE)
    Process: 644174 ExecStartPre=/usr/bin/docker pull mysql:latest (code=exited, status=0/SUCCESS)
    Process: 644179 ExecStart=/usr/bin/docker run --name httpurr-db --rm -v /database:/var/lib/mysql:ro --user 1002:1002 -p 3306:3306/tcp mysql:latest (code=exited, statu>
   Main PID: 644179 (code=exited, status=1/FAILURE)
        CPU: 75ms
May 15 13:13:57 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Main process exited, code=exited, status=1/FAILURE
May 15 13:13:57 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Failed with result 'exit-code'.
May 15 13:13:58 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Scheduled restart job, restart counter is at 15.
May 15 13:13:58 ip-172-31-30-87.ec2.internal systemd[1]: Stopped httpurr-db.service - HTTPurr-DB.
May 15 13:13:58 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Start request repeated too quickly.
May 15 13:13:58 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Failed with result 'exit-code'.
May 15 13:13:58 ip-172-31-30-87.ec2.internal systemd[1]: Failed to start httpurr-db.service - HTTPurr-DB.
```

=> It's failing

- I have restarted the `httpurr-db` service an checked the status. It was active and after a bit I checked the status again and it failed.
- So it’s starting and stopping , same behaviour noticed using : docker stats.
- I used `journalctl` to get more detailed logs :

```zsh
[saadia@ip-172-31-30-87 ~]$ journalctl -u httpurr-db.service | grep 'May 15 13:'
// I should have used ctrl+G to get the navigate to the end of the logs
May 15 13:36:08 ip-172-31-30-87.ec2.internal systemd[1]: Started httpurr-db.service - HTTPurr-DB.
May 15 13:36:08 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15 13:36:08+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.4.0-1.el9 started.
May 15 13:36:08 ip-172-31-30-87.ec2.internal docker[646935]: ln: failed to create symbolic link '/var/lib/mysql/mysql.sock': Read-only file system
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:08.954313Z 0 [System] [MY-015015] [Server] MySQL Server - start.
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.255183Z 0 [Warning] [MY-010091] [Server] Can't create test file /var/lib/mysql/mysqld_tmp_file_case_insensitive_test.lower-test
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.255308Z 0 [System] [MY-010116] [Server] /usr/sbin/mysqld (mysqld 8.4.0) starting as process 1
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.258958Z 0 [Warning] [MY-010091] [Server] Can't create test file /var/lib/mysql/mysqld_tmp_file_case_insensitive_test.lower-test
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.258966Z 0 [Warning] [MY-010159] [Server] Setting lower_case_table_names=2 because file system for /var/lib/mysql/ is case insensitive
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.259283Z 0 [Warning] [MY-010122] [Server] One can only use the --user switch if running as root
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: mysqld: File './binlog.index' not found (OS errno 30 - Read-only file system)
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.262234Z 0 [ERROR] [MY-010119] [Server] Aborting
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.263723Z 0 [System] [MY-010910] [Server] /usr/sbin/mysqld: Shutdown complete (mysqld 8.4.0)  MySQL Community Server - GPL.
May 15 13:36:09 ip-172-31-30-87.ec2.internal docker[646935]: 2024-05-15T13:36:09.263994Z 0 [System] [MY-015016] [Server] MySQL Server - end.
May 15 13:36:09 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Main process exited, code=exited, status=1/FAILURE
May 15 13:36:09 ip-172-31-30-87.ec2.internal systemd[1]: httpurr-db.service: Failed with result 'exit-code'.
```

- I saw the first error encountered : `May 15 13:36:08 ip-172-31-30-87.ec2.internal docker[646935]: ln: failed to create symbolic link '/var/lib/mysql/mysql.sock': Read-only file system`
- This message indicates that there was an attempt to create a symbolic link on a file system `/var/lib/mysql` that is currently mounted as read-only.
- On the `/etc/systemd/system/httpurr-db.service` , volume configuration is set to read-only : `-v /database:/var/lib/mysql:ro`
- On `/etc/systemd/system/httpurr-db.service` I changed : `-v /database:/var/lib/mysql:ro` to `-v /database:/var/lib/mysql`
- I tried to restart the service and I got :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo systemctl restart httpurr-db
Warning: The unit file, source configuration file or drop-ins of httpurr-db.service changed on disk. Run 'systemctl daemon-reload' to reload units.
```

- I needed to reload the daemon because I changed the service configuration : `sudo systemctl daemon-reload` and then restarted the service `sudo systemctl restart httpurr-db`
- I checked the status of `httpurr-db`, it's active
- I tested the endpoint again , but got the same error.
- I hit a dead end and seek help from Daniel, he suggest to trace down the request flow from the host to the container.
- What we know so far :
  - The container holding the server is running and the server is running , we saw the msg : `“2024/05/15 14:06:43 HTTPurr starting…”` on `journalctl` logs after restarting the `httpurr` service.
  - The request didn’t hit the server , since we haven’t seen any printed msg from the handler function `http.HandleFunc("/", purrEndpoint)` on the server code `httpurr/httpurr.go`
  - The server set to be running on port `8080` locally (on the docker container)
- We checked what port is running in the `httpurr` container with the id `d0f5c6e1ecb3`, by opening a bash session inside the container :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker exec -it d0f5c6e1ecb3 sh
[saadia@ip-172-31-30-87 ~]$ ps
PID   USER     TIME  COMMAND
    1 1001      0:00 /httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr --dbpass=miauw --dbname=httpurr
   34 1001      0:00 sh
   40 1001      0:00 ps
[saadia@ip-172-31-30-87 ~]$ cat  /proc/1/net/tcp  // nothing running there
[saadia@ip-172-31-30-87 ~]$ cat  /proc/1/net/tcp6
sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
0: 00000000000000000000000000000000:1F90 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000  1001        0 1373189 1 00000000b446b2b1 100 0 0 10 0
```

- From there we saw that the local port where the server is running is effectively `8080` IPV6 (`1F90` is hex of `8080`)
- I checked if the docker container is correctly set to map `:80` to `:8080`.

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED       STATUS       PORTS                                                  NAMES
d0f5c6e1ecb3   httpurr        "/httpurr --dbhost=h…"   3 hours ago   Up 3 hours   0.0.0.0:80->80/tcp, :::80->80/tcp, 8080/tcp            httpurr
7b3e3e36a59b   mysql:latest   "docker-entrypoint.s…"   3 hours ago   Up 3 hours   0.0.0.0:3306->3306/tcp, :::3306->3306/tcp, 33060/tcp   httpurr-db
```

- No the server container set to receive requests from the host port `80` on the post `80`
- I changed the service configuration on `/etc/systemd/system/httpurr.service`
- I changed the mapping from `-p 80:80/tcp` to `-p 80:8080/tcp` on :

```zsh
ExecStart=/usr/bin/docker run --name httpurr --user 1001:1001 --link httpurr-db --rm -p 80:80/tcp httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr --dbpass=miauw --dbname=httpurr
```

- I reloaded the daemon : `sudo systemctl daemon-reload` and restarted the service : `sudo systemctl restart httpurr`
- I sent the request again : `lynx -dump  http://localhost`
- I got the right response :
  |\---/|
  | o*o |
  \_^*/
- I got a new image whenever I send a new request.

## Enhancing security

The problem : We want to hide the database password , that we can see using `ps`

- I created `.env` holding the `DB_PASSWORD` on `/etc/http-db-pass/` directory
- I changed the `httpurr` service file to read password from `.env`

```zsh
[saadia@ip-172-31-30-87 ~]$ cat /etc/systemd/system/httpurr.service
[Unit]
Description=HTTPurr
After=docker.service
Requires=docker.service
Wants=httpurr-db.service
[Service]
EnvironmentFile=/etc/httpurr-db-pass/.env
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker stop httpurr
ExecStartPre=-/usr/bin/docker rm httpurr
ExecStartPre=/usr/bin/docker build -t httpurr /httpurr
ExecStart=/usr/bin/docker run --name httpurr --user 1001:1001 --link httpurr-db --rm -p 80:80/tcp httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr --dbpass=${DB_PASSWORD} --dbname=httpurr
[Install]
WantedBy=multi-user.target
```

- I set the permissions to allow only the root user and user running the service to view the file :

```zsh
sudo chown root:root /etc/http-db-pass/.env
sudo chmod 600 /etc/http-db-pass/.env
```

- I reloaded to daemon and restarted the `httpurr` service.
- I sent the request again and checked the `ps` : `ps auxww | grep miauw`
- I can still see the password in the logs
- Then I tried to create the .env file in the container and get the password from there.
- I changed the `httpurr` service file to do that as following:

```zsh
[saadia@ip-172-31-30-87 ~]$ cat /etc/systemd/system/httpurr.service
[Unit]
Description=HTTPurr
After=docker.service
Requires=docker.service
Wants=httpurr-db.service
[Service]
EnvironmentFile=/etc/httpurr-db-pass/.env
TimeoutStartSec=0
Restart=always
ExecStartPre=-/usr/bin/docker stop httpurr
ExecStartPre=-/usr/bin/docker rm httpurr
ExecStartPre=/usr/bin/docker build -t httpurr /httpurr
ExecStart=/usr/bin/docker run --name httpurr --user 1001:1001 --link httpurr-db --rm -p 80:80/tcp httpurr --dbhost=httpurr-db --dbport=3306 --dbuser=httpurr --env-file /etc/httpurr-db-pass/.env --dbname=httpurr
[Install]
WantedBy=multi-user.target
```

- Also changed the database code `catsdb/table.go` to use the password from .env instead of the command line :

```zsh
db_pass = os.getenv('DB_PASSWORD')
```

- I need to reload the daemon and restart the server to be able to test it. But I couldn't because of the problem on my instance.

## Wrong path

After changing the go version in `go.mod` and `Dockerfile` I tried to rebuild the image manually.

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker build . -t httpur
ERROR: failed to solve: failed to copy: write /var/lib/docker/buildkit/content/ingest/0a308448261750b0bafb3ea55a96287438d761352dfe7f5a9e7d0b10c5cd85ac/data: no space left on device
```

- I tried to delete the previous image to get more space :

```zsh
sudo docker image rm ac31c528f11a : got the error.
Error response from daemon: conflict: unable to delete ac31c528f11a (must be forced) - image is being used by stopped container 3d659c3b293a
```

- I deleted the related the container :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker rm 3d659c3b293a
3d659c3b293a
```

- I tried to delete the image again, but I got the same error again with different container :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker image rm ac31c528f11a
Error response from daemon: conflict: unable to delete ac31c528f11a (must be forced) - image is being used by stopped container 96a5dd65cd0d
```

- I checked the list of all the containers (stopped or running):

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker ps -a
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS          PORTS                                                  NAMES
39c1e00923fa   mysql:latest   "docker-entrypoint.s…"   30 seconds ago   Up 24 seconds   0.0.0.0:3306->3306/tcp, :::3306->3306/tcp, 33060/tcp   httpurr-db
96a5dd65cd0d   httpurr        "/httpurr"               18 hours ago     Created                                                                jovial_jones
```

- I was just one container using the httpur image, and I deleted it :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker rm 96a5dd65cd0d
96a5dd65cd0d
```

- I removed the image again :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker image rm ac31c528f11a :
Untagged: httpurr:latest
Deleted: sha256:ac31c528f11a553e477f282a091cbc3d48df3c5e2b0d617a2b343640612d6ea8
```

- I tried to create a new image again, but got the same error :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker build . -t httpurr
ERROR: failed to solve: failed to copy: write /var/lib/docker/buildkit/content/ingest/0a308448261750b0bafb3ea55a96287438d761352dfe7f5a9e7d0b10c5cd85ac/data: no space left on device
```

- I deleted all unused images : `sudo docker image prune -a`
- I tried to build the image again , but got the same error (in different level):

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker build . -t httpurr
ERROR: failed to solve: failed to register layer: write /usr/local/go/pkg/tool/linux_amd64/pprof: no space left on device
```

- I checked the space in the instance:

```zsh
[saadia@ip-172-31-30-87 ~]$ df -h
Filesystem      Size  Used Avail Use% Mounted on
devtmpfs        4.0M     0  4.0M   0% /dev
tmpfs           475M     0  475M   0% /dev/shm
tmpfs           190M  3.0M  187M   2% /run
/dev/xvda1      8.0G  7.9G  107M  99% /
tmpfs           475M     0  475M   0% /tmp
/dev/xvda128     10M  1.3M  8.7M  13% /boot/efi
tmpfs            95M     0   95M   0% /run/user/1003
```

- The root filesystem is nearly full
- I tried to investigate which directories/files are using lot of space

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo du -h --max-depth=1 / | sort -hr
7.5G	/
5.2G	/var
2.1G	/usr
189M	/database
73M	/home
39M	/boot
...
```

- I checked the `/var` usage :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo du -h --max-depth=1 /var | sort -hr
5.3G	/var
5.2G	/var/lib
79M	/var/log
57M	/var/cache
...
```

- I checked the `/var/lib` usage :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo du -h --max-depth=1 /var/lib | sort -hr
4.9G	/var/lib/docker
4.9G	/var/lib
19M	/var/lib/rpm
17M	/var/lib/selinux
2.0M	/var/lib/dnf
1.3M	/var/lib/sss
432K	/var/lib/cloud
236K	/var/lib/systemd
144K	/var/lib/containerd
...
```

- I checked the `/var/lib` usage :

```zsh
[saadia@ip-172-31-30-87 ~]$ ssudo du -h --max-depth=1 /var/lib/docker | sort -hr
5.1G	/var/lib/docker
4.7G	/var/lib/docker/overlay2
477M	/var/lib/docker/buildkit
9.8M	/var/lib/docker/image
48K	/var/lib/docker/network
24K	/var/lib/docker/volumes
0	/var/lib/docker/tmp
0	/var/lib/docker/swarm
0	/var/lib/docker/runtimes
0	/var/lib/docker/plugins
0	/var/lib/docker/containers
```

- I saw that `/var/lib/docker/overlay2` use most of the space, I checked what is this directory: is one of the storage driver of docker , I should not touch it and looked for how to of clean up docker;
- I found that I can check docker disk usage and what space can I claim using :

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker system df
TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE
Images          1         0         577.9MB   577.9MB (100%)
Containers      0         0         0B        0B
Local Volumes   0         0         0B        0B
Build Cache     60        3         894.1MB   893.7MB
```

- I decided to reclaim at least the build cache storage space

```zsh
[saadia@ip-172-31-30-87 ~]$ sudo docker system prune -a
Total reclaimed space: 893.7MB
```

- I tried build the image again `sudo docker build . -t httpurr` and it was rebuilt without error.
  **Note :**
  All those steps maybe were not relevant as we can just restart the service and rebuild the image , or maybe they are and we may hit the same problem of the space of the disk.

## Answering questions :

1. **How is systemd configured? Where are the logs?**

- systemd is configured using unit files (e.g `httpurr.service`) , we can create or edit them to describe how the `systemd` start , stop, reload and manage the service.
- We can access to the logs using `journalctl` command.

2. **Can you describe what docker does?**

- Docker containerise applications and their dependencies by packing them into a lightweight portable containers.
  It separate the application from the infrastructures for a faster delivery.

3. **Can you describe the two ways that docker volumes are used in this setup?**

- In the `httpurr-db` service set up we have a `bind mount` volume `-v /database:/var/lib/mysql` which is used to store database data on the host machine, ensuring data is not lost when containers are restarted or removed.
- I couldn't find the other way

4. **Docker uses images. What are images? And how do you create your own?**

- Images are templates to create docker containers. Built on layers containing everything needed to run an application (code, dependencies, libraries ...)
- Images can be create using a `Dockerfile` containing all the instructions defining what it goes into this image.
