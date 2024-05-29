# Notes

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


### Handling errors in your code

When encountering or detecting errors, there are typically four options for handling them:

1. **Propagating** the error to the calling function. (And possibly wrapping it with some extra contextual information.)
1. **Working around** the error to recover from it.
1. **Terminating** the program completely.
1. **Ignoring** the error. (Sometimes an error actually doesn't matter.)
