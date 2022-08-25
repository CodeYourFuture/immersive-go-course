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
