# Preparation

## Prerequisite learning

Before you start this course, there's a few things we assume you've done:

- You're familiar with the essentials of writing code in JavaScript
- You have experience with JavaScript in the browser and in [Node](https://nodejs.org/en/)
- You've completed the [Tour of Go](https://go.dev/tour/welcome/1)

This is important because we don't cover the basic language features of Go: you need to be familiar with writing Go functions and methods, plus the basics of types in Go. You'll also need to know how to navigate [packages and documentation](https://pkg.go.dev/), and we have a [short guide on how to do that](#learn-how-to-navigate-go-documentation).

Remember: you can _always_ Google or ask for help if you get stuck.

## Set up and get to know your IDE

We're going to assume you're using [Visual Studio Code](https://code.visualstudio.com/) in this course.

VS Code is a lot more powerful that a simple text editor: it can help you write code, spot mistakes, and help you fix them. But it needs to know details of the programming language, so we need to install an extension to support go.

Set up the Go extension for VS Code by following [these instructions](https://code.visualstudio.com/docs/languages/go). Have a read of the features listed on that page.

Some of the really useful ones:

1. Go to definition - when you call a function or use a variable, this will show you where it was defined. This can help to understand what code is doing and why, and even works when calling into things like the standard library.
2. Go to references - this will show you what bits of code use a variable or function. Say you're changing a function to add a new parameter, this can help you find all the places you'll need to modify.
3. Autocomplete - Go can guess what you're about to type, and save you time. But more importantly, it can tell you what exists - if you're looking to use something related to HTTP, and you think it's probably in the `http` package, you can type `http.` and see what's auto-completed for you - that could help you find the code you want without needing to switch to Google.

Write a bit of Go in VS Code and experiment with these features. A small investment now will save a lot of time in the future!

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

## Working on the projects

To work on the projects:

- [Fork this repository](https://github.com/CodeYourFuture/immersive-go-course/fork) — this [creates your own copy](https://docs.github.com/en/get-started/quickstart/fork-a-repo)
- Work through [the projects in order](../README.md) — they are designed to conceptually build on each other and increase in complexity
- For each project, open a [pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests) against your own fork with your implementation of the project. Don't put code for different projects into the same PR!
- Remember to make [small, incremental commits](https://dustinspecker.com/posts/making-smaller-git-commits/) and to only include code for a _single_ project in each commit

In general, work through each project like this:

1. Make it **work**
1. Make it **right**
1. Optional: Make it **pretty** (or fast)

Making it "right" means ensuring the structure is clear, reasonably simple, and documented, with error handling and tests in place. But it does't mean perfect!

In other words, resist the temptation to optimise or refactor individual pieces of the project, for example to use fancy techniques like go-routines and channels, until you have project working end-to-end. Sometimes it's only when you have the whole working project in front of you that you can see how the individual pieces should _really_ work, and optimising too early can waste time.

### Opening the right directory

When writing Go in VS Code, many of the helper features only work if you opened the folder directly containing the file named `go.mod`.

So when working in `immersive-go-course`, you _need to open a new window for project directory with the code you're working on_.

If you opened VS Code in the main ("root") directory of `immersive-go-course`, you'll probably see a lot of red squiggly lines in your Go code and errors starting "could not import".

You can have more than one copy of VS Code open at a time if you need to.
