<!--forhugo
+++
title="Output and Error Handling"
+++
forhugo-->

In this project you're going to get familiar with conventions around output and exit codes. You will learn about error handling, and how to apply these in the [Go programming language](https://go.dev/).

Timebox: 3 days

## Learning Objectives:

- Write to standard out and standard error
- Exit programs with conventional exit codes
- Explain when to propagate errors
- Decide when to wrap errors, and terminate due to errors.

## Project

Most programs can run into problems. Sometimes these problems are recoverable, and sometimes not.

We will write a program which may encounter several kinds of error. We will handle these errors. We will tell the user about these errors, and make that information easy to consume.

### The program

In this project, we have been supplied with a server. Our server code lives in the `server` subdirectory of this project. Run it by `cd`ing into that directory, and running `go run`. The server is an HTTP server, which listens on port 8080 and responds in a few different ways:

- If you make an HTTP [`GET`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/GET) request to it, it will respond with the current weather. When this happens, you should display it to the user on the terminal.
- Sometimes this server will overload and respond with a status code 429. When this happens, the client should:
  1. wait the amount of time indicated in the `Retry-After`[^1] response header, and
  2. attempt the request again.
- Sometimes, this server will drop a connection before responding. When this happens:

  _You should assume:_

  1. The server is non-responsive.
  2. Making more requests to it could make things worse.

  _The client should:_

  1.  Give up its request.
  2.  Tell the user something irrecoverable went wrong.

[^1]: Learn about [the `Retry-After` header on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After). This header has two formats: delay in seconds, or a timestamp.)

Have a read of the server code. Make sure you understand what it's doing, and what kinds of responses you may need to handle.

We won't propose changes to the server code as part of this project. This server is intentionally buggy because as part of the exercise we sometimes need to handle bad responses. We may, however, want to make edits to it while we're developing our client to help us better manually test out yur code (randomness is hard to test against!).

Our final client code should work against the original server code.

We're going to focus on how we handle errors, and how we present output to the user.

### Standard out and standard error

Terminal programs can write output to two places: standard out and standard error. Standard out is sometimes called standard output.

Think of these as "files" you can write to. These files are connected to your terminal by default. What you write to these files will display in your terminal.

You can connect, or redirect, these "files" somewhere else. Run this command:

```console
> echo hello > /tmp/the-output-from-echo
```

`echo`'s job is to write something to standard out. But if you run this command, you won't see "hello" output to the terminal. Instead, it will get written to the file `/tmp/the-output-from-echo`. You redirected echo's standard out to the file you named. If you `cat /tmp/the-output-from-echo` you'll see "hello" was written in that file.

You can redirect standard error by writing `2>` instead of `>`:

```console
> echo hello 2> /tmp/error-from-echo
```

In this example, you'll still see "hello" on your terminal. This is because you didn't redirect standard out anywhere. Also notice that `/tmp/error-from-echo` was created, but is empty. This is because `echo` didn't write anything to standard error.

You can redirect both by using both redirect instructions:

```console
> echo hello > /tmp/file-1 2> /tmp/file-2
```

Often, we want to pass the output of our program to some other program or file for further processing. Standard error is a place you can write information which a user may want, but which you don't want to pass forward.

The user might want to know about a problem or progress messages explaining what's happening. These are important messages, but they aren't helpful for your forward process.

Imagine your program writes out a series of scores, one per line. Next, you write those scores to a file for another program to analyse.

You have two different use cases: As a user, you want to know why your program appears to be hanging or failing. As a consumer of the output, you only want the scores. If the output file had "Waiting for HTTP request… or "Server error" printed in it, that would be annoying to process.

What about when something goes wrong? Say, your network connection goes down and you cannot fetch your scores. In your score-analysing program, you may want to assume that anything it analyses is a number. You might need to add numbers together. Reporting your problem on standard error means you won't try to add the error string "Network was down" to a number.

#### Standard out and standard error in Go

In Go, we access standard out and standard error with `os.Stdout` and `os.Stderr`.

Write to them by writing code like `fmt.Fprint(os.Stdout, "Hello")` or `fmt.Fprint(os.Stderr, "Something went wrong")`. The "F” before "printf” stands for "file”. We're saying "print some string to a file I'll specify as the first argument”. In Unix systems, we often like to pretend anything we read or write is a file.

More often, we'll write `fmt.Print("Hello")`. This is the same as writing `fmt.Fprint(os.Stdout, "Hello")`. If you look at [the Go standard library source code](https://cs.opensource.google/go/go/+/refs/tags/go1.19.5:src/fmt/print.go;l=251-253), you can see it's literally the same. But we can choose to write to other locations, and _sometimes_ we _should_. This is why we are thinking about `Stdout` and `Stderr` separately now.

#### When to write to standard out/error

As a rule, write the _intended_ output of your program to standard out. Write anything that isn't the intended output of your program to standard error.

Some things you may write to standard error:

- Progress messages explaining what the program is doing.
- Error information about something that went wrong.

Thinking about our program we're going to write, that means we're likely to write:

##### Standard out

- The current weather. It's what our program is for.

##### Standard error

- A message saying that we've been asked to wait and retry later. It's a progress message, not the intended output of our program.
- Error information if the server seems broken. It's not the intended output of our program; it's diagnostic information.

### Exit codes

By convention, most programs exit with an exit code of `0` when they complete successfully. Programs exit with any number that isn't `0` when they fail.

Often a specific program will attach specific meaning to specific non-zero exit codes. `1` may mean "You didn't specify all the flags I needed”. `2` may mean "A remote server couldn't give me information I needed”. But these meanings belong to that program. There are no general conventions for specific non-0 exit codes across different programs.

By default, your program will exit with exit code `0` unless you tell it to do otherwise, or it crashes.

In Go, you can choose what code to exit your program with by calling [`os.Exit`](https://pkg.go.dev/os#Exit). After calling `os.Exit`, your program stops and can't do anything else.

From a terminal, you can check the exit code of a process you just ran by inspecting the environment variable `$?`. In this example, `pwd` is a successful command, and `cd /doesnotexist` is an unsuccessful one:

```console
> pwd
/tmp
> echo $?
0
> cd /doesnotexist
cd: no such file or directory: /doestnotexist
> echo $?
1
```

Many common utility programs document their exit codes in their manpages. You can access these by running `man [command-name]`[^2]. For example:
[^2]: If you find yourself stuck in a manpage, pressing `q` will exit from it.

```console
> man cat
```

shows the somewhat vague but still useful:

```
EXIT STATUS
     The cat utility exits 0 on success, and >0 if an error occurs.
```

Or the same for `curl`:

```console
> man curl
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

When writing code, we often need to handle the possibility that an error has occurred.

This may be an explicit error returned from a function. When you make a [`GET`](https://pkg.go.dev/net/http#Client.Get) request, `Get` returns an error alongside the response. This error may be `nil`, but will be non-`nil` if, for instance, the server was down.

Alternatively, this may be something which we detect, but which other code didn't tell us was an error. If we make a `GET` request to a server which returns a `429` status code, the error will be `nil`. By looking at [`Response.StatusCode`](https://pkg.go.dev/net/http#Response) we can see that something went wrong which we may need to handle.

When encountering or detecting errors, there are typically four options for handling them:

1. **Propagating** the error to the calling function. (And possibly wrapping it with some extra contextual information.)
1. **Working around** the error to recover from it.
1. **Terminating** the program completely.
1. **Ignoring** the error. (Sometimes an error actually doesn't matter.)

When we should do each of these isn't always obvious, but here are some guidelines:

#### 1. Propagating the error to the calling function

Propagation is our default behaviour. If we don't know how to handle an error, we should early-return from our function, handing the error to the caller.

> This means that when function Alf calls function Betty, and Betty goes wrong, Betty hands her error to Alf.

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

If the password file doesn't exist, the error message

```
open .some-file: no such file or directory
```

is less useful than an error message like

```
failed to read password file: open .some-file: no such file or directory
```

#### Context is valuable

By wrapping the error with more contextual information, we help the person seeing the error understand:

- _what_ went wrong
- _why_ it failed
- and what they need to _do_ to fix the situation

With this in mind, we may write `readPassword` instead like:

```go
func readPassword() (string, error) {
    password, err := os.ReadFile(".some-file")
    if err != nil {
        return "", fmt.Errorf("failed to read password file: %w", err)
    }
    return string(password), nil
}
```

Learn more about creating and wrapping errors in [Effective Error Messages in Go](https://earthly.dev/blog/golang-errors/).

#### 2. Working around the error to recover from it

Sometimes, an error may be expected, or may be recoverable. Suppose we have some expensive computation we want to do, but which may have already been done and saved to a file. We may try to read the file. If we encounter an error that the file doesn't exist, we may know how to compute the answer we need instead.

This kind of behaviour will depend on the problem domain you're solving. There isn't a general rule for when to write workarounds.

#### 3. Terminating the program completely

> When user input can terminate your program, user input can reduce your service's capacity. This is dangerous.

When we run into errors, there's often nothing we can do about them.

Let's imagine we're running a server. Our error happened when processing one request. We don't want to end our program! We want to respond to that request saying an error happened. But we also want to keep trying to process other requests.

In other situations it may make sense to end our program and exit (with a non-`0` status code). Some examples are when starting up a server, or writing a program that does a one-off task. In these cases we aren't responding to user inputs, or there's only one user of our program.

Where we exit is worth thinking about. We generally don't want to call `os.Exit` from anywhere except our main function.

There are a few reasons for this:

1. If we call `os.Exit`, there's no way any code can handle that or recover. Let's say we started calling `os.Exit` in some other function. It's possible we'll end up in the future calling that function from a request handler. If so, we'll end up terminating the whole server because we couldn't handle one request. This will probably cause an outage, because no one will be able to talk to our server any more.
1. When writing unit tests, we generally don't want our program to exit. If you call `os.Exit` inside a unit test, it will stop running. In general, don't call `os.Exit` from any code called from a test.

**Only call `os.Exit` from your main function.**

Code changes a lot over time. Writing an `os.Exit` call creates a function that isn't safe to call from other places. Your future self, or someone else, could call your function without realising this. The easiest way to avoid this is to use the rule: **only ever call `os.Exit` from your main function**. Everything else should propagate any errors they encounter.

#### 4. Ignoring the error

It may seem to you that you can ignore errors. Be suspicious of this idea.

It is true that occasionally an error actually doesn't matter at all. But this is rare, and you should be wary if you think this is the case. Programmers cause many real-life bugs by ignoring or poorly handling errors.

Read [Simple Testing Can Prevent Most Critical Failures](https://www.usenix.org/system/files/conference/osdi14/osdi14-paper-yuan.pdf) to learn more about this. It's worth a read.

---

### Back to our program

Recall the server that gives us the weather.

Our task is to write a client, in Go, which makes HTTP requests to that server and tells the user about the weather.

We should focus in this project on handling errors and retries.

If the server replies with a retryable error, we will retry it appropriately. For a `429` response code, this means:

1. Reading the Retry-After response header.
1. Calling `time.Sleep` until the appropriate time has passed.
1. Trying again.

Make sure all error messages are clear and useful to the user, that we're properly printing to standard out or standard error when appropriate, and that our program always exits with an appropriate exit code.

#### Creating the program

We'll create the program in the same directory as this README.md file. In a terminal, `cd` to this directory, and run `go mod init github.com/CodeYourFuture/immersive-go-course/projects/output-and-error-handling`. This will create a `go.mod` file (which indicates a Go program lives here). Then we will need to manually create a `main.go` file and start writing our program in it.

#### time.Sleep

##### If we're going to sleep for more than 1 second:

We should notify the user that things may be a bit slow because we're doing a retry.

##### If the server tells us we should sleep for more than 5 seconds:

We should give up and tell the user we can't get them the weather.

##### If we can't determine how long to sleep for:

You should decide whether we should sleep for some amount of time (and if so what) and then retry, or give up. Make sure to write down the reasons for your decision.

##### When to stop

If the server terminates our connection, we will give up and tell the user that we can't get them the weather.

### Make sure that:

1. All error messages are clear and useful to the user.
2. We're properly printing to standard out or standard error when appropriate.
3. Our program always exits with an appropriate exit code.
