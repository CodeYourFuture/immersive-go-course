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

## Instructions

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

### Steps (notes)

- make sure the working directory is go-ls: `cd go-ls`
- create `go.mod`, `main.go`, `cmd/root.go`

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

- `go get -u github.com/spf13/cobra@latest`
- follow cobra [user guide](https://github.com/spf13/cobra/blob/master/user_guide.md) to make a root command that prints hello in `cmd/root.go`
- implement basic ls with `os.ReadDir`
- allow the command to take arguments with `cobra.ArbitraryArgs`
- when passed an argument such as `go-ls assets`, read from the passed directory
- ensure that this directory path can be relative: `go-ls ..` and `go-ls ../go-ls` should both work
- handle the error (e.g. `Error: fdopendir go.mod: not a directory` when passing `go-ls` a file argument: `go-ls go.mod`)
- update `go-ls` to match `ls` in terms of how it handles files (hint: `os.Stat`)
- make `go-ls -h` include a helpful description
- bonus: write some tests for `go-ls`

[go]: https://go.dev/
[cat]: https://en.m.wikipedia.org/wiki/Cat_(Unix)
[ls]: https://en.m.wikipedia.org/wiki/Ls
[cobra]: https://github.com/spf13/cobra#overview
