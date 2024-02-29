+++
title="CYF+02 Prep Workbook"
date="01 Jan 2024"    
versions=["1-1-0"]
weight=1
+++

# Prep

The CYF+ Immersive Engineering Programme is an intensive three month course. There's a lot to get through and you'll need to hit the ground running. To prepare for this exciting opportunity, you will:

## Computers

Read some of [How Computers Really Work by Matthew Justice](https://www.howcomputersreallywork.com/).

You _must_ read chapters 1, 2, 7, 9, 10, 11, 12.

Any other chapters may be interesting (feel free to read them! We particularly recommend 8), but aren't necessary.

Note that this book isn't free - someone at CYF may be able to lend you a copy if you need.

## Linux

- [ ] Play [the Bandit](https://overthewire.org/wargames/bandit/) - you must be able to complete up to level 20 (repeatedly)
- [ ] Take this [Unix 101 course](https://www.opsschool.org/unix_101.html) (and then try 102)
- [ ] Print out this [Linux-Cheat-Sheet](https://www.loggly.com/wp-content/uploads/2015/05/Linux-Cheat-Sheet-Sponsored-By-Loggly.pdf)

## Go

- [ ] Learn the basics of the Go programming language: [Get Started - The Go Programming Language](https://go.dev/learn/)
- [ ] Read this, you might find it useful for working through your first projects: [How to use the fmt package in Golang](https://www.educative.io/answers/how-to-use-the-fmt-package-in-golang)
- [ ] And optionally [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/)
- [ ] Read the [Pointers chapter](https://www.golang-book.com/books/intro/8) of [An Introduction to Programming in Go](https://www.golang-book.com/books/intro), and do the problems listed in the chapter.
- [ ] Read about [stack and heap memory](https://courses.engr.illinois.edu/cs225/fa2022/resources/stack-heap/) - note that the examples here are C++ not Go, but you should be able to follow them. One important difference between these languages is that in Go, you don't need to worry about whether a value lives on the heap or the stack (in fact, Go is allowed to choose where your variable lives, and even move it around) - this means that it's ok to return a pointer to a stack variable in Go. You also don't need to call `delete` yourself in Go - Go is [a garbage collected language](https://en.wikipedia.org/wiki/Garbage_collection_(computer_science)) so it will delete variables when they're no longer needed.

## Tooling
- [ ] [Set up and get to know your IDE](https://github.com/CodeYourFuture/immersive-go-course/tree/main/prep#set-up-and-get-to-know-your-ide)
- [ ] [Learnt how to navigate Go documentation](https://github.com/CodeYourFuture/immersive-go-course/tree/main/prep#learn-how-to-navigate-go-documentation)

## Projects

Complete the [prep](https://github.com/CodeYourFuture/immersive-go-course/blob/main/prep/README.md) and first _five_ [projects](https://github.com/CodeYourFuture/immersive-go-course/blob/main/projects/README.md) from the [Immersive Go](https://github.com/CodeYourFuture/immersive-go-course) course

- [ ] [Prep](https://github.com/CodeYourFuture/immersive-go-course/tree/main/prep)
- [ ] [Output and Error Handling](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/output-and-error-handling)
- [ ] [CLI Files](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/cli-files)
- [ ] [File Parsing](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/file-parsing)
- [ ] [Concurrency](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/concurrency)
- [ ] [Servers & HTTP Requests](https://github.com/CodeYourFuture/immersive-go-course/tree/main/projects/http-auth)

## Reading

Dip in to some longer books, but don't feel you need to read the whole lot!

- [ ] [Google - Site Reliability Engineering](https://sre.google/sre-book/table-of-contents/)
- [ ] [The Unix Tools Philosophy](https://www.linuxtopia.org/online_books/gnu_linux_tools_guide/the-unix-tools-philosophy.html)
- [ ] [The Phoenix Project](https://smile.amazon.co.uk/Phoenix-Project-Helping-Business-Anniversary/dp/B00VBEBRK6/)
