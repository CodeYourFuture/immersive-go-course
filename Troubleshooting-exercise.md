~ df -h
Report file system disk space usage

Logs:
Filesystem Size Used Avail Use% Mounted on
devtmpfs 474M 0 474M 0% /dev
tmpfs 483M 0 483M 0% /dev/shm
tmpfs 483M 460K 483M 1% /run
tmpfs 483M 0 483M 0% /sys/fs/cgroup
/dev/xvda1 8.0G 8.0G 1.4M 100% /
tmpfs 97M 0 97M 0% /run/user/1000
tmpfs 97M 0 97M 0% /run/user/1001

I can see that /dev/xvda1 is full

~ du -sch /
du - short for disk usage is used to estimate file space usage.
The du command can be used to track the files and directories which are consuming excessive amount of space on hard disk drive.

Logs:
1.6G /
1.6G total

~ sudo du -sch /

Logs:
1.8G /
1.8G total

Admin mode shows a bit more usage but still much less than df is showing.

From researching why du (disk usage) only shows 1.6G in the "/" directory whereas df(disk free) shows full disk that is mounted to "/". But it can't see these unlinked files, and will thus return less space used than the filesystem will

First hypothesis:
There is a process running which is holding open a file which has been deleted.
df thinks that the partition is full since the inodes are all taken up by that process.

~ lsof | grep '(deleted)'
I want to see processes that use files, which were deleted

Logs:
Nothing

I need root permissions.

~ sudo lsof | grep '(deleted)'

Logs:
systemd-j 1741 root txt REG 202,1 325536 199663 /usr/lib/systemd/systemd-journald (deleted)
dbus-daem 2582 dbus txt REG 202,1 227352 8676635 /usr/bin/dbus-daemon (deleted)
systemd-l 2583 root txt REG 202,1 606928 199665 /usr/lib/systemd/systemd-logind (deleted)
findme 3513 root 3w REG 202,1 6583025664 2186 /root/$TMP (deleted)
sleep     25933     root    3w      REG   202,1 6583025664       2186 /root/$TMP (deleted)

~ sudo lsof +L1

Logs:
COMMAND PID USER FD TYPE DEVICE SIZE/OFF NLINK NODE NAME
systemd-j 1741 root txt REG 202,1 325536 0 199663 /usr/lib/systemd/systemd-journald (deleted)
dbus-daem 2582 dbus txt REG 202,1 227352 0 8676635 /usr/bin/dbus-daemon (deleted)
systemd-l 2583 root txt REG 202,1 606928 0 199665 /usr/lib/systemd/systemd-logind (deleted)
findme 3513 root 3w REG 202,1 6583025664 0 2186 /root/$TMP (deleted)
sleep     26433 root    3w   REG  202,1 6583025664     0    2186 /root/$TMP (deleted)

lsof +L1 lists: open files that are completely unlinked. There are no remaining open file descriptors that refer to it.

So think I am confident to think that some of these processes are causing the issue and need to be stopped.

Command "findme" and "sleep" need more investigation.
Name: Shows the name of the mount point and file system on which the file resides, in this case /root/$TMP, where TMP I would assume stand for temporary, so it can is probably isn't very crucial for our system.

In computing, "sleep" is a command in Unix, Unix-like and other operating systems that suspends program execution for a specified time. So it seems like FINDME is the one we need to work with.

~ lsof -p 3513

I want to look into command "findme" in details.

Logs:
COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
findme 3513 root cwd unknown /proc/3513/cwd (readlink: Permission denied)
findme 3513 root rtd unknown /proc/3513/root (readlink: Permission denied)
findme 3513 root txt unknown /proc/3513/exe (readlink: Permission denied)
findme 3513 root NOFD /proc/3513/fd (opendir: Permission denied)

~ ps -ef | grep 3513 | grep -v grep

List the process for confirming the open process.

Logs:
root 3513 1 0 Nov09 ? 00:00:01 /bin/sh /sbin/findme
root 6921 3513 0 11:10 ? 00:00:00 sleep 10

I can see the is a path to the file findme.

~ cd /sbin/findme
~ file findme

Logs:
findme: POSIX shell script, ASCII text executable

It is a shell script file. Lets have a look at more stats and what is inside.

~ stat findme

Logs:
File: ‘findme’
Size: 147 Blocks: 8 IO Block: 4096 regular file
Device: ca01h/51713d Inode: 12584491 Links: 1
Access: (0755/-rwxr-xr-x) Uid: ( 0/ root) Gid: ( 0/ root)
Access: 2022-10-13 12:59:50.451194465 +0000
Modify: 2022-10-13 12:59:50.451194465 +0000
Change: 2022-11-09 13:15:17.890346507 +0000
Birth: -

I can see that it has been changed recently.

~ tail findme

Logs:
#!/bin/sh
set -e
// displays and sets the names and values of shell and Linux environment variables

TMP="$(mktemp)"
// create temporary file

exec 3>"\$TMP"
// links filed descriptor 3 to "\$TMP" file

dd bs="1M" count="9000" if="/dev/zero" of="\$TMP" || :
// is used to clean a drive or device before forensically copying data. The point in dd having seperate bs and count argument is that bs controls how much is written at a time. Count is multiplied by bs. So we filling our "\$TMP" with zeros??? Also it looks like 9000M is more that the space we have.

rm -f "\$TMP"
while true; do sleep 10; done
// It seems like we will run out of space before the dd command does it's job. Also we are saying "while true" what is true in this case? Is it an infinite loop??????
I think it results in continuous output of progress every 10 seconds until it is true ?

After some research and discussion with Berkeli, I found out that "|| :" will ignore error or overfilled space. So go to next line, remove file and enter infinite loop.

~ pgrep -f findme
I want to find PID of just this process
Logs: 3513

I can make sure that running the following command I only kill this process
~ kill -9 $(pgrep -f findme)

~ pgrep -f findme
No PID found, process is removed

~ df -h
Logs:
df -h
Filesystem Size Used Avail Use% Mounted on
devtmpfs 474M 0 474M 0% /dev
tmpfs 483M 0 483M 0% /dev/shm
tmpfs 483M 408K 483M 1% /run
tmpfs 483M 0 483M 0% /sys/fs/cgroup
/dev/xvda1 8.0G 1.9G 6.2G 24% /
tmpfs 97M 0 97M 0% /run/user/1001

File system now has space available.

If the process was not findme and it was apache or service we can not just kill:
I fill up the disk again:

~ sh /usr/sbin/findme &
Throws expected error that no space left.

Check disc is full:
~ df -h

Logs:
Filesystem Size Used Avail Use% Mounted on
devtmpfs 474M 0 474M 0% /dev
tmpfs 483M 0 483M 0% /dev/shm
tmpfs 483M 496K 483M 1% /run
tmpfs 483M 0 483M 0% /sys/fs/cgroup
/dev/xvda1 8.0G 8.0G 636K 100% /
tmpfs 97M 0 97M 0% /run/user/1001
tmpfs 97M 0 97M 0% /run/user/1000

~ sudo lsof +L1
Logs:
COMMAND PID USER FD TYPE DEVICE SIZE/OFF NLINK NODE NAME
systemd-j 1741 root txt REG 202,1 325536 0 199663 /usr/lib/systemd/systemd-journald (deleted)
dbus-daem 2582 dbus txt REG 202,1 227352 0 8676635 /usr/bin/dbus-daemon (deleted)
systemd-l 2583 root txt REG 202,1 606928 0 199665 /usr/lib/systemd/systemd-logind (deleted)
sh 14545 margarita 0u CHR 136,0 0t0 0 3 /dev/pts/0 (deleted)
sh 14545 margarita 1u CHR 136,0 0t0 0 3 /dev/pts/0 (deleted)
sh 14545 margarita 2u CHR 136,0 0t0 0 3 /dev/pts/0 (deleted)
sh 14545 margarita 3w REG 202,1 6600851456 0 12584487 /home/margarita/$TMP (deleted)
sleep     15227 margarita    0u   CHR  136,0        0t0     0        3 /dev/pts/0 (deleted)
sleep     15227 margarita    1u   CHR  136,0        0t0     0        3 /dev/pts/0 (deleted)
sleep     15227 margarita    2u   CHR  136,0        0t0     0        3 /dev/pts/0 (deleted)
sleep     15227 margarita    3w   REG  202,1 6600851456     0 12584487 /home/margarita/$TMP (deleted)

We know that $TMP file is full, it showed error which was ignored and the program deleted the file $TMP but as process is in infinite loop, the file can't actually be deleted unless we kill the process.

So I need to research how can I make this file smaller and keep process running. We know that program wants to remove this file, so by making it smaller we won't lose any essential data cos it is not our intention to keep anything in it.

https://unix.stackexchange.com/questions/88808/most-efficient-method-to-empty-the-contents-of-a-file

People suggest different methods but some say turncate is most efficient.

Truncating a file is much faster and easier than deleting the file , recreating it, and setting the correct permissions and ownership . Also, if the file is opened by a process, removing the file may cause the program that uses it to malfunction.

Sounds just like what we need!!!! But we deleted file in our bash script so we cannot access it directly.
We know that that file descriptor will still be there because it has to due to process still running.

[margarita@ip-172-31-31-103 ~]$ ls -lh /proc/14545/fd
total 0
lrwx------ 1 margarita margarita 64 Nov 10 13:43 0 -> /dev/pts/0 (deleted)
lrwx------ 1 margarita margarita 64 Nov 10 13:43 1 -> /dev/pts/0 (deleted)
lrwx------ 1 margarita margarita 64 Nov 10 13:43 2 -> /dev/pts/0 (deleted)
lr-x------ 1 margarita margarita 64 Nov 10 13:43 255 -> /usr/sbin/findme
l-wx------ 1 margarita margarita 64 Nov 10 13:43 3 -> /home/margarita/$TMP (deleted)

/proc/[pid]/fd/
This is a subdirectory containing one entry for each file
which the process has open, named by its file descriptor,
and which is a symbolic link to the actual file. Thus, 0
is standard input, 1 standard output, 2 standard error,
and so on.

To be able to truncate a file, you need to have write permissions on the file. Usually, you would use sudo for this, but the elevated root privileges do not apply to the redirection.
https://linuxize.com/post/truncate-files-in-linux/#:~:text=Truncating%20a%20file%20is%20much,that%20uses%20it%20to%20malfunction.

We can see that this is file descriptor 3 pointing to our file.

~ truncate -s 0 /proc/14545/fd/3
~ df -h
Filesystem Size Used Avail Use% Mounted on
devtmpfs 474M 0 474M 0% /dev
tmpfs 483M 0 483M 0% /dev/shm
tmpfs 483M 440K 483M 1% /run
tmpfs 483M 0 483M 0% /sys/fs/cgroup
/dev/xvda1 8.0G 1.9G 6.2G 24% /
tmpfs 97M 0 97M 0% /run/user/1001
