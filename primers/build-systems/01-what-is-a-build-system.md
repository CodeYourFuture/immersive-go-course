+++
title="1. What is a build system?"
+++

# 1. What is a build system?

## What is building software?

"Building software" can refer to many things. One example is a software engineer may describe themself as "building software" when they write code.

In this primer we're specifically concerned with one idea of building software: Converting some input file (typically source code) into some useful output (often executables or bundles ready to run or deploy, or test results). Most programming languages have some kind of build process, but different languages require more or less building.

If you have written complex JavaScript, you have probably seen `package.json` files, and interacted with `npm` or `yarn`. These are tools used to build JavaScript. Here is an example `package.json` file:

```json
{
    "name": "fancy-project",
    "type": "module",
    "scripts": {
        "build": "mkdir -p dist && webpack",
        "test": "jest test"
    },
    "dependencies": {
        "dotenv": "^16.4.5"
    },
    "devDependencies": {
        "jest": "^29.7.0",
        "webpack": "^5.91.0",
        "webpack-cli": "^5.1.4"
    }
}
```

Some of the actions these tools are used to do are:
* Fetch all of the dependencies a project needs into a `node_modules` folder.
* Run all of the tests in a project.
* Combine a source file and everything it imports into a single file, and minify that file.

These kind of actions are all related to building the software. Building software is concerned with analysing dependencies between different pieces of code, and running some commands to take source code and convert it to something useful.

## What is a build tool?

A "build system" or "build tool" is a tool which is aware of the relationships between different pieces of code, and performs actions to transform code from one form to another. A good build tool is **correct** (it always produces the correct output when you run it) and **fast** (i.e. it performs its actions as quickly as possible). We will use the terms "build tool" and "build system" interchangeably.

## Why do we automate builds?

A lot of the time we could perform the actions a build tool performs ourselves manually.

To fetch dependencies, we could look in a `package.json` file and, for each dependency, work out what URL can be used to fetch it, download that URL, and unpack the files into the correctly named directory in the `node_modules` directory. And then do the same for each of its dependencies too.

To run tests, we could manually find the `jest` tool in our `node_modules` directory, and run it to run our tests.

We automate these processes with build tools for a few reasons:
1. It avoids us needing to know things. What URL can a dependency be found at? What should its directory in `node_modules` be named? While we could find these things out, the build tool knows them, so we don't need to.
2. It avoids us needing to work out the order we need to do things. The build tool knows which actions need to happen before which other actions, and will make sure they're done in the right order.
3. It avoids us needing to think about what's already been done, and what needs to be re-done. Imagine we had already manually downloaded `dotenv` and `jest`. Then we changed the version of `dotenv` in the `package.json`. We can delete and re-download both `dotenv` and `jest` (which ensures we're doing the _correct_ thing), but this would be slower than it could be. If the version of `jest` hasn't changed, maybe we don't need to delete it and download it again. By just leaving the downloaded `jest` as-is, we can be _faster_. But manually analysing what we can skip and what we need to re-do is complicated and error-prone (what if changing the version of `dotenv` actually _does_ mean we need to re-download `jest`? Skipping it would mean our build was not _correct_!)
