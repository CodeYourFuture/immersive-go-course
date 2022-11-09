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

From researching why du (disk usage) only shows 1.6G in the "/" directory whereas df(disk free) shows full disk that is mounted to "/".

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
