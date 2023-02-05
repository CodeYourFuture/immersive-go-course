<!--forhugo
+++
title="Output and Error Handling"
+++
forhugo-->

In this project you're going to get familiar with conventions around output and exit codes, as well as error handling, and how to apply these in the [Go programming language][go].

Timebox: 3 days

## Objectives:

- Know when to write to standard out and standard error
- Exit programs with conventional exit codes
- Know when to propagate errors, wrap errors, and terminate due to errors.

## Project

Most programs can run into problems. Sometimes these problems are recoverable, and other times they can't be recovered from.

We are going to write a program which may encounter several kinds of error, and handle them appropriately. We will also make sure we tell the user of the program information they need, in ways they can usefully consume it.

### The program

In this project, we have been supplied with a server - its code lives in the `server` subdirectory of this project. You can run it by `cd`ing into that directory, and running `go run .`. The server is an HTTP server, which listens on port 8080 and responds in a few different ways:
* If you make an HTTP GET request to it, it will respond with the current weather. When this happens, you should display it to the user on the terminal.
* Sometimes the server simulates being overloaded by too many requests, and responds with a status code 429. When this happens, the client should wait the amount of time indicated in the `Retry-After` response header, and attempt the request again. You can learn about [the `Retry-After` header on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After). Note that it has two formats, and may contain either a number of seconds or a timestamp.
* Other times, it may drop a connection before responding. When this happens, you should assume the server is non-responsive (and that you making more requests to it may make things worse), and give up your request, telling the user something irrecoverable went wrong.

Have a read of the server code and make sure you understand what it's doing, and what kinds of responses you may need to handle.

We are not expected to change the server code as part of this project - it is intentionally buggy because we sometimes need to handle bad responses. We may, however, want to make edits to it while you're developing your client to help you better manually test out your code (randomness is hard to test against!). Our final client code should work against the server code as it was given to us with no modifications.

We're going to focus in this project on how we handle errors, and how we present output to the user.

### Standard out and standard error

Typically, terminal programs have two places they can write output to: standard out (also known as "standard output"), and standard error. As far as your program is concerned, these are "files" you can write to, but in reality they are both by default connected to your terminal, so if you write to them, what you write will end up displayed in your terminal.

You can also request that one or both of these "files" be redirected somewhere else, e.g. in your terminal you can run:

```console
% echo hello > /tmp/the-output-from-echo
```

`echo`'s job is to write something to standard out, but if you run this command, you won't see "hello" output to the terminal, instead, it will get written to the file `/tmp/the-output-from-echo`. `echo`'s standard out was redirected. If you `cat /tmp/the-output-from-echo` you'll see `hello` was written in that file.

You can redirect standard error by writing `2>` instead of `>`:

```console
% echo hello 2> /tmp/error-from-echo
```

In this example, you'll still see `hello` on your terminal (because you didn't redirect standard out anywhere), and you'll see `/tmp/error-from-echo` was created, but is empty, because `echo` didn't write anything to standard error.

You can redirect both if you want by using both redirect instructions:

```console
% echo hello > /tmp/file-1 2> /tmp/file-2
```

The main reason we have these two different locations is that often times we want the output of our program to be passed to some other program or file for further processing. Standard error exists as a place you can write information which a user may be interested in (e.g. information about something going wrong, or progress messages explaining what's happening), but which you don't want to pass on for that further processing.

For example, imagine your program writes out a series of scores, one per line, and you were going to write those scores to a file which another program may analyse. If the output file had "Waiting for HTTP request..." or "Server error" printed in it, that would be annoying to process later on, but as a user, you my want to know why your program appears to be hanging or failing.

Another example is when something goes wrong - in your score-recording program, you may want to (reasonably) assume that anything it outputs is a number. But if something goes wrong (say, your network connection was down so the scores couldn't be fetched), reporting that on standard error means you won't accidentally try to add the error string "Network was down" to some other number.

#### Standard out and standard error in Go

In go, standard out and standard error can be accessed as `os.Stdout` and `os.Stderr`.

You can write to them by writing code like `fmt.Fprint(os.Stdout, "Hello")` or `fmt.Fprint(os.Stderr, "Something went wrong")`. (The "F" before "printf" stands for "file" - we're saying "print some string to a file I'll specify as the first argument". A lot of times in Unix systems, we like to pretend anything we read or write is a file).

More often, we'll write `fmt.Print("Hello")` - this is the same as writing `fmt.Fprint(os.Stdout, "Hello")` (if you look at [the Go standard library source code](https://cs.opensource.google/go/go/+/refs/tags/go1.19.5:src/fmt/print.go;l=251-253), you can see it's literally the same), but it's worth remembering we can choose to write to other locations, like standard error, if it's more appropriate.

#### When to write to standard out/error

As a rule, the intended output of your program should be written to standard out, and anything that isn't the intended output of your program should be written to standard error.

Some things you may write to standard error:
* Progress messages explaining what the program is doing.
* Error information about something that went wrong.

Thinking about our program we're going to write, that means we're likely to write:
* The current weather to standard out - it's what our program is for.
* A message saying that we've been asked to wait and retry later to standard error - it's a progress message, not the intended output of our program.
* Error information if the server seems broken to standard error - it's not the intended output of our program, it's diagnostic information.

### Exit codes

By convention, most programs exit with an exit code of `0` when they successfully did what they expected, and any number that isn't `0` when they didn't. Often a specific program will attach specific meaning to specific non-zero exit codes (e.g. `1` may mean "You didn't specify all the flags I needed" and `2` may mean "A remote server couldn't give me information I needed"), but there are no general conventions for specific non-`0` exit codes across different programs.

By default, your program will exit with exit code `0` unless you tell it to do otherwise, or it crashes.

In Go, you can choose what code to exit your program with by calling [`os.Exit`](https://pkg.go.dev/os#Exit). After calling `os.Exit`, your program stops and can't do anything else.

From a terminal, you can check the exist code of a process you just ran by inspecting the environment variable `$?`. In this example, `pwd` is a successful command, and `cd /doesnotexist` is an unsuccessful one:

```console
% pwd
/tmp
% echo $?
0
% cd /doesnotexist
cd: no such file or directory: /doestnotexist
% echo $?
1
```

Many common utility programs document their exit codes in their manpages, which you can access by running `man [command-name]` (if you find yourself stuck in a manpage, pressing `q` will exit from it). For example:

```console
% man cat
```

shows the somewhat vague but still useful:

```
EXIT STATUS
     The cat utility exits 0 on success, and >0 if an error occurs.
```

Or the same for `curl`:

```console
% man curl
```

shows lots of very specific information:

<details><summary>Click to expand the full output of the exit codes section of <code>man curl</code></summary>

```
EXIT CODES
       There are a bunch of different error codes and their corresponding error messages that may appear under error conditions. At the time of this writing, the exit codes are:

       0      Success. The operation completed successfully according to the instructions.

       1      Unsupported protocol. This build of curl has no support for this protocol.

       2      Failed to initialize.

       3      URL malformed. The syntax was not correct.

       4      A feature or option that was needed to perform the desired request was not enabled or was explicitly disabled at build-time. To make curl able to do this, you probably need another build of libcurl.

       5      Could not resolve proxy. The given proxy host could not be resolved.

       6      Could not resolve host. The given remote host could not be resolved.

       7      Failed to connect to host.

       8      Weird server reply. The server sent data curl could not parse.

       9      FTP access denied. The server denied login or denied access to the particular resource or directory you wanted to reach. Most often you tried to change to a directory that does not exist on the server.

       10     FTP accept failed. While waiting for the server to connect back when an active FTP session is used, an error code was sent over the control connection or similar.

       11     FTP weird PASS reply. Curl could not parse the reply sent to the PASS request.

       12     During an active FTP session while waiting for the server to connect back to curl, the timeout expired.

       13     FTP weird PASV reply, Curl could not parse the reply sent to the PASV request.

       14     FTP weird 227 format. Curl could not parse the 227-line the server sent.

       15     FTP cannot use host. Could not resolve the host IP we got in the 227-line.

       16     HTTP/2 error. A problem was detected in the HTTP2 framing layer. This is somewhat generic and can be one out of several problems, see the error message for details.

       17     FTP could not set binary. Could not change transfer method to binary.

       18     Partial file. Only a part of the file was transferred.

       19     FTP could not download/access the given file, the RETR (or similar) command failed.

       21     FTP quote error. A quote command returned error from the server.

       22     HTTP page not retrieved. The requested URL was not found or returned another error with the HTTP error code being 400 or above. This return code only appears if --fail is used.

       23     Write error. Curl could not write data to a local filesystem or similar.

       25     FTP could not STOR file. The server denied the STOR operation, used for FTP uploading.

       26     Read error. Various reading problems.

       27     Out of memory. A memory allocation request failed.

       28     Operation timeout. The specified time-out period was reached according to the conditions.

       30     FTP PORT failed. The PORT command failed. Not all FTP servers support the PORT command, try doing a transfer using PASV instead!

       31     FTP could not use REST. The REST command failed. This command is used for resumed FTP transfers.

       33     HTTP range error. The range "command" did not work.

       34     HTTP post error. Internal post-request generation error.

       35     SSL connect error. The SSL handshaking failed.

       36     Bad download resume. Could not continue an earlier aborted download.

       37     FILE could not read file. Failed to open the file. Permissions?

       38     LDAP cannot bind. LDAP bind operation failed.

       39     LDAP search failed.

       41     Function not found. A required LDAP function was not found.

       42     Aborted by callback. An application told curl to abort the operation.

       43     Internal error. A function was called with a bad parameter.

       45     Interface error. A specified outgoing interface could not be used.

       47     Too many redirects. When following redirects, curl hit the maximum amount.

       48     Unknown option specified to libcurl. This indicates that you passed a weird option to curl that was passed on to libcurl and rejected. Read up in the manual!

       49     Malformed telnet option.

       52     The server did not reply anything, which here is considered an error.

       53     SSL crypto engine not found.

       54     Cannot set SSL crypto engine as default.

       55     Failed sending network data.

       56     Failure in receiving network data.

       58     Problem with the local certificate.

       59     Could not use specified SSL cipher.

       60     Peer certificate cannot be authenticated with known CA certificates.

       61     Unrecognized transfer encoding.

       63     Maximum file size exceeded.

       64     Requested FTP SSL level failed.

       65     Sending the data requires a rewind that failed.

       66     Failed to initialise SSL Engine.

       67     The user name, password, or similar was not accepted and curl failed to log in.

       68     File not found on TFTP server.

       69     Permission problem on TFTP server.

       70     Out of disk space on TFTP server.

       71     Illegal TFTP operation.

       72     Unknown TFTP transfer ID.

       73     File already exists (TFTP).

       74     No such user (TFTP).

       77     Problem reading the SSL CA cert (path? access rights?).

       78     The resource referenced in the URL does not exist.

       79     An unspecified error occurred during the SSH session.

       80     Failed to shut down the SSL connection.

       82     Could not load CRL file, missing or wrong format.

       83     Issuer check failed.

       84     The FTP PRET command failed.

       85     Mismatch of RTSP CSeq numbers.

       86     Mismatch of RTSP Session Identifiers.

       87     Unable to parse FTP file list.

       88     FTP chunk callback reported error.

       89     No connection available, the session will be queued.

       90     SSL public key does not matched pinned public key.

       91     Invalid SSL certificate status.

       92     Stream error in HTTP/2 framing layer.

       93     An API function was called from inside a callback.

       94     An authentication function returned an error.

       95     A problem was detected in the HTTP/3 layer. This is somewhat generic and can be one out of several problems, see the error message for details.

       96     QUIC connection error. This error may be caused by an SSL library error. QUIC is the protocol used for HTTP/3 transfers.

       XX     More error codes will appear here in future releases. The existing ones are meant to never change.
```
</details>

### Handling errors in your code

A lot of time when writing code, we need to handle the possibility that an error has occurred.

This may be an explicit error returned from a function (e.g. when you make a `GET` request, [`Get`](https://pkg.go.dev/net/http#Client.Get) returns an `error` (which may be `nil`, but will be non-`nil` if, for instance, the server couldn't be connected to) alongside the response).

Alternatively, this may be something which we detect, but which other code didn't tell us was an error. For instance, if we make a `GET` request to a server which returns a 429 status code, the `error` will be `nil`, but by looking at [`Response.StatusCode`](https://pkg.go.dev/net/http#Response) we can see that something went wrong which we may need to handle.

When encountering or detecting errors, there are typically four options for how to handle them:
1. Propagating the error to the calling function (and possibly wrapping it with some extra contextual information).
1. Working around the error to recover from it.
1. Terminating the program completely.
1. Ignoring the error - sometimes an error actually doesn't matter.

When we should do each of these isn't always obvious, but here are some guidelines:

#### Propagating the error to the calling function

This is generally our default behaviour. If an error has happened, and we don't know how to handle it, we should early-return from our function, handing the error to the caller.

This means that generally any time we write a function, and it calls another function which may return an error, our function will probably also possibly return an error.

Often times, we want to wrap the error to provide more context. For instance, say we have the following code:

```go
package main

import "os"

func main() {
    password, err := readPassword()
    // ...
}

func readPassword() (string, error) {
    password, err := os.ReadFile(".some-file")
    if err != nil {
        return "", err
    }
    return string(password), nil
}
```

If the password file doesn't exist, the error message `open .some-file: no such file or directory` is less useful than an error message like `failed to read password file: open .some-file: no such file or directory`. By wrapping the error with more contextual information, you help the person seeing the error understand _what_ went wrong, _why_ it failed, and what they need to do to fix the situation.

Accordingly, we may write `readPassword` instead like:

```go
func readPassword() (string, error) {
    password, err := os.ReadFile(".some-file")
    if err != nil {
        return "", fmt.Errorf("failed to read password file: %w", err)
    }
    return string(password), nil
}
```

You can learn more about creating and wrapping errors in Go by reading [this article](https://earthly.dev/blog/golang-errors/).

#### Working around the error to recover from it

Sometimes, an error may be expected, or may be recoverable. For instance, suppose we have some expensive computation we want to do, but which may have already been done and saved to a file. We may try to read the file, but if we encounter an error that it doesn't exist, we may know how to compute the answer we need instead.

This kind of behaviour will depend entirely on the problem domain you're solving, and there isn't really a general rule for when it's appropriate.

#### Terminating the program completely

A lot of the time, when we run into errors, there's nothing we can do about them.

If we're running a server, and the error happened when processing one request, normally we don't want to terminate our program - we just want to respond to that request saying an error happened, but keep trying to process other requests. If you terminate your program because of input from users, you are implicitly giving users the ability to reduce your service's capacity (at least temporarily). This is dangerous.

Other times, for instance when first starting up a server, or when writing a program that isn't a server but just does a one-off task, it may make sense to terminate our program, and exit (with a non-`0` status code).

_Where_ we do this, however, is worth thinking about. We generally don't want to call `os.Exit` from anywhere _except_ our `main` function.

There are a few reasons for this:
1. If we call `os.Exit`, there's no way any code can handle that or recover. Let's say we started calling `os.Exit` in some other function - it's possible we'll end up in the future calling that function from a request handler, and we'll end up terminating the whole server just because one request couldn't be handled. This will probably cause an outage, because no one will be able to talk to our server any more.
1. When writing unit tests, we generally don't want our program to exit. But if you call `os.Exit` inside a unit test, it will stop running. In general, we never want to call `os.Exit` from any code which is called from a test.

While you may write an `os.Exit` call in some other function, thinking it's only ever called from places it's ok to call `os.Exit`, code changes a lot over time, and you may find yourself or someone else adding calls to functions from other places without realising that your function isn't safe to be called everywhere. The easiest way to avoid this is to use the rule: only ever call `os.Exit` from your `main` function - everything else should propagate any errors they encounter.

#### Ignoring the error

Sometimes an error actually doesn't matter at all, and can just be ignored. This is rare, and you should be wary if you think this is the case. [This paper](https://www.usenix.org/system/files/conference/osdi14/osdi14-paper-yuan.pdf) describes that a very large percentage of real-life bugs are caused by ignoring or poorly handling errors - it's worth a read.

### Back to our program

Recall the server we've been supplied with for telling the weather.

Our task is to write a client, in Go, which makes HTTP requests to that server and tells the user about the weather.

We should focus in this project on making sure:
1. If the server replies with a retryable error, we will retry it appropriately. For a 429 response code, this means reading the `Retry-After` response header, calling `time.Sleep` until the appropriate time has passed, and trying again.
   * If we're going to sleep for more than 1 second, we should notify the user that things may be a bit slow because we're doing a retry.
   * If the server tells us we should sleep for more than 5 seconds, we should give up and tell the user we can't get them the weather.
   * If we can't determine how long to sleep for, consider what the best thing to do is - you should decide whether we should sleep for some amount of time (and if so what) and then retry, or give up. Make sure to write down why you decided what you dod.
1. If the server terminates our connection, we will give up and tell the user that we can't get them the weather.

Make sure all error messages are clear and useful to the user, that we're properly printing to standard out or standard error when appropriate, and that our program always exits with an appropriate exit code.

[go]: https://go.dev/
