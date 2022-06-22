# CLI & Files

In this project you're going to get familiar with the [Go programming language][go] by rebuilding two very common tools that programmers and computer administrators use all the time: [cat][cat] and [ls][ls].

Timebox: 4 days

Objectives:

- Install and use [cobra][cobra]
- Use go build/go install/go get etc
- Understand what a process is & the basics of process - lifecycle
- Accept arguments on the CLI
- Open, read (and close) files from CLI arguments
- Reading directories for files

## Project

You're going to build a command-line application that reads data from your computer's file system.

Let's say you have a few files containing poems in a directory:

- `dew.txt`
- `for_you.txt`
- `rain.txt`

The `ls` tool will list the files in that directory:

```
> ls
dew.txt
for_you.txt
rain.txt
```

While the `cat` command will output the contents of a file to you:

```
> cat dew.txt
“A World of Dew” by Kobayashi Issa

A world of dew,
And within every dewdrop
A world of struggle.
```

As you can see, `cat` takes an argument, namely the file to be written out. You'll learn how to do this too.

You're going to rebuild these two tools in Go:

```bash
> go-ls
dew.txt
for_you.txt
rain.txt

> go-cat dew.txt
“A World of Dew” by Kobayashi Issa

A world of dew,
And within every dewdrop
A world of struggle.
```

### go-ls

These steps will leave some work for you to do. If something's not clear, search around for to find the answer. If you're stuck for half an hour at most, ask for help. Remember to use Git to store your progress, committing often in small increments with useful descriptions.

A full implementation of this is available on the `impl/cat-ls` branch.

The `go-ls` directory is provided for you with some example files in the `assets/` folder.

In your command line/terminal, make sure your working directory is go-ls: `cd go-ls`

Create `go.mod`, `main.go`, `cmd/root.go`. The `touch` command creates files: `touch go.mod`.

```go
// go.mod
module vinery/cli-files/go-ls

go 1.18
```

```go
// main.go
package main

import (
	"vinery/cli-files/go-ls/cmd"
)

func main() {
	cmd.Execute()
}
```

```go
// cmd/root.go
package cmd

func Execute() {}
```

We're going to use the [Cobra][cobra] package to make these commands. It does a lot of magic for you.

Install the Cobra package using the `go get` command: `go get -u github.com/spf13/cobra@latest`

The Cobra [user guide](https://github.com/spf13/cobra/blob/master/user_guide.md) will show you how to make a root command that prints hello in `cmd/root.go`.

To use your command, install and run it: `go install .`

Then run the following so that your command line knows where to look for the executable code that Go is generating: `export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))`

You should now be able to run `go-ls`.

Now, when you change your code, install and run it: `go install . && go-ls`

Now you've got something working, we'll speed up the steps. You can do it!

The `ls` command reads a directory, generating a list of files or sub-directories. Go comes with packages for interacting with files built-in, include a package called [os][os], which contains a function called `ReadDir` which will do lots of the work of `ls` for you.

See if you can implement basic `ls` with `os.ReadDir`. It should read the list of files in the current, "working" directory:

```
> go install .
> cd assets
> go-ls
dew.txt
for_you.txt
rain.txt
```

The real `ls` allows you pass a directory to be read: `ls assets`.

Extend your `go-ls` to allow the command to take arguments (look for `cobra.ArbitraryArgs`) and then, when passed an argument such as `go-ls assets`, read from the passed directory.

Make that this directory path can be relative: `go-ls ..` and `go-ls ../go-ls` should both work.

Handle the error `Error: fdopendir go.mod: not a directory` when passing `go-ls` a file argument: `go-ls go.mod`. Think hard about why this is happening before you try to fix it.

Update `go-ls` to match `ls` in terms of how it handles files (hint: `os.Stat`) — it should just output the name of the file.

Make `go-ls -h` include a helpful description.

If you smash through this, here's some fun/tricky extensions:

- Write some tests for `go-ls`
- Extend `go-ls` to support some more features of the real `ls` (for example, `ls -m assets`)

[go]: https://go.dev/
[cat]: https://en.m.wikipedia.org/wiki/Cat_(Unix)
[ls]: https://en.m.wikipedia.org/wiki/Ls
[cobra]: https://github.com/spf13/cobra#overview
[os]: https://pkg.go.dev/os
