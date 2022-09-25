# Preparation

## Prerequisite learning

Before you start this course, there's a few things we assume you've done:

- You're familiar with the essentials of writing code in JavaScript
- You have experience with JavaScript in the browser and in [Node][node]
- You've completed the [Tour of Go][tourofgo]

This is important because we don't cover the basic language features of Go: you need to be familiar with writing Go functions and methods, plus the basics of types in Go. You'll also need to to navigate [packages and documentation](https://pkg.go.dev/).

Remember: you can _always_ Google or ask for help if you get stuck.

## Set up and get to know your IDE

We're going to assume you're using VS Code in this course.

With Code Your Future so far, you've mostly used VS Code just as a text editor. It can also be a lot more powerful than that, but in order to do so, it needs to know details of the language you're writing. Set up the Go extension for VS Code by following [these instructions](https://code.visualstudio.com/docs/languages/go).

Have a read of the features listed on that page.

Some of the really useful ones:

1. Go to definition - when you call a function or use a variable, this will show you where it was defined. This can help to understand what code is doing and why, and even works when calling into things like the standard library.
2. Go to references - this will show you what bits of code use a variable or function. Say you're changing a function to add a new parameter, this can help you find all the places you'll need to modify.
3. Autocomplete - Go can guess what you're about to type, and save you time. But more importantly, it can tell you what exists - if you're looking to use something related to HTTP, and you think it's probably in the `http` package, you can type `http.` and see what's auto-completed for you - that could help you find the code you want without needing to switch to Google.

Write a bit of Go in VS Code and experiment with these features. A small investment now will save a lot of time in the future!

> :warning: **Opening the right directory** - When writing Go in VS Code, many of these features only work if you opened the folder directly containing the file named `go.mod`.
>
> When working in `immersive-go-course`, you need to open a new window in the directory the code you're working on lives in.
>
> If you opened VS Code in the root directory of `immersive-go-course`, you'll probably see a lot of red squiggly lines in your Go code and errors starting "could not import".
>
> You can have more than one copy of VS Code open at a time if you need to.

## Learn how to navigate Go documentation

The [Go standard library](https://pkg.go.dev/std) has lots of documentation and examples, such as [net/http](https://pkg.go.dev/net/http). To find documentation, you can use the search feature or Google something like `golang net/http`, which will generally help you find what you're looking for.

The website `pkg.go.dev` also hosts documentation for other go packages that you might use: again you can use the search feature or Google for it.

The structure is fairly similar between different packages. Let's take `fmt` as an example:

- At the start is a summary of the package, discussing what it is for and some important information.
- The Index lists all the [functions that the package exports](https://www.callicoder.com/golang-packages/#exported-vs-unexported-names)
- The Examples show you how to use the package (good place to start!)
- The Variables, Functions and Types sections contain specific documentation on what the package contains
- Lastly, Source Files is... the source code! Good Go libraries can be quite readable, so don't be scared to jump in. You'll probably learn a lot,

> 💡 The best way to get familiar with a new package, particularly if the documentation is a bit dense (like for the the `fmt` package), is to look at the Examples section. It will have some basic and advanced usage that you can often use straight away.

## Conventions used in projects

### Command line examples

In the projects you'll need to run some programs on the command line. To show this, and the output, we'll use an example like the following:

```console
> echo "Hello, world."
Hello, world.
```

What we mean is: "run the command to the right of the `>` sign": `echo "Hello, world."`

Everything after that line is _output_ from the command: `Hello, world.`

In the above example, the command is `echo` with the argument `"Hello, world."`. Here's a more complex example using `curl`:

```console
> curl -i 'http://localhost:8080/'
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:42:30 GMT
Content-Length: 42

<!DOCTYPE html><html><em>Hello, world</em>
```

It doesn't matter what this does: what's important is the input command and the expected output.

Input command:

```console
curl -i 'http://localhost:8080/'
```

Expected output:

```
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sun, 24 Jul 2022 09:42:30 GMT
Content-Length: 42

<!DOCTYPE html><html><em>Hello, world</em>
```

**Important**: the output from commands that you run will often not be identical to the example. Dates, times and counts will be different.

Sometimes we may put more than one command in the same snippet:

```console
> echo hello
hello
> echo goodbye
goodbye
```

Generally each time a line starts with a `>`, it's a new command (but _occasionally_ it may be output from a previous one!)
