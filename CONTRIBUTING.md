# Contributing to Code Your Future Immersive Go Course

This repository is a permissively licensed, open source work. We welcome contributions, both as corrections/improvements and additions, and recommend filing an issue for discussion before investing too much work.

## Goals

This repository primarily exists to help junior software engineers get experience with the Go programming language, and learn about topics that would be useful as a software engineer working in production.

It is intended to be followed without needing support (i.e. should explain things, ideally via references to high quality existing materials), but requiring real learning and research by those taking it (i.e. it shouldn't spoon-feed every step and every answer). This is important, as working out how to do things, and thinking about design, are important software engineering skills.

## How we teach

### Learning by doing

We primarily teach via exercises and projects. We describe outcomes, show examples, and give references (either terms to Google, or links to useful material).

### Don't reinvent the wheel

Where we think there isn't a useful / accurate / well-targeted existing resource, we may write explanations ourselves, but this should be rare - many excellent people have come before us!

### Make it challenging

We don't explain every step along the way - working out the steps is an important part of learning. Deciding exactly how many steps to explain is a tricky balance.

### Fundamentals first

As far as possible and reasonable, we teach foundational topics _before_ teaching abstractions. The goal is that learners can re-apply their knowledge to new situations by re-combining the fundamentals, rather than simply being able to use a specific tool.

- **Yes**: a project that teaches setting up & running code on a cloud VM (such as [EC2](https://aws.amazon.com/ec2/)) before a project that hides the details (such as [Lambda](https://aws.amazon.com/lambda/))
- **Yes**: a project that introduces configuring cloud infrastructure before a project using an [IaC](https://en.m.wikipedia.org/wiki/Infrastructure_as_code) tool such as [Terraform](https://www.terraform.io/).
- **No**: introducing concurrency with a framework like [ants](https://github.com/panjf2000/ants) without building familiarity with the basics of channels.

### Self-contained

Each project here should be self-contained; projects should not extend existing projects and assume they have been built. However, projects can reference previous projects if there is incidental overlap in what is being built.

- **Yes**: (after instructions for the build) "You might have spotted that this is the same server as in the server-database project."
- **No**: "Implement this server as per the instructions in the server-database project."
- **Also no**: "Take your server from the server-database project and make it also respond to gRPC requests"

The course itself should be self-contained in that it should not depend on the existence or implementation of any specific Code Your Future resources. This is so that the course does not add a maintenance burden to CYF as it develops: for example, the person managing AWS resources for learners at CYF should _not_ have to update instructions in this course. Instead, this course should link to canonical resources inside CYF.

It's also fine to rely on prerequisite knowledge as long as there is a link to resources which would explain that prerequisite at the correct level.

Rule of thumb: a person should be able to complete a project _without_ being part of CYF, and without having done any of the other projects.

## How we work

Work on this repository is tracked in the [Immersive Go Course GitHub Project](https://github.com/orgs/CodeYourFuture/projects/20).

Work should progress Backlog -> To Do -> In Progress -> Done. Every task should have a Target Sprint capturing which CYF sprint it contributes to. If you're not sure, ask in the GitHub issue.

For the initial implementation of the course:

- Tasks are organised into Development Sprints of one week
- The work in the current Development Sprint is on the [Current Sprint board](https://github.com/orgs/CodeYourFuture/projects/20/views/2)
- The work in the current Development Sprint should be prioritised over everything else
- Every issue and PR should have an associated **milestone** relating to the CYF Course Sprint

## Project structure

Each project has a name, e.g. `http-auth`.

Every project contains at least the following elements, and may have others as needed:

### On `main` in a directory named for the project

1. A README.md containing:
   1. Learning objectives. A list of things the person doing the project will understand, know, and be able to do after they have completed the project.
   1. A concise description of the high-level idea of the project - what's it for?
   1. A decscription of the outcomes desired from the project, and a rough sequence of steps.
   1. Ideas for extensions to the project.
1. Any scaffolding or supporting files that may be useful (e.g. files to be served, starter code if getting started isn't the focus of the project).

### On `main` in the README.md

A reference to the project, in order (or at least after prerequisite projects), and a quick description of it.

### On a branch named `impl/${project_name}`

1. An IMPLEMENTATION.md file, talking through a whole implementation, pointing out the important pieces, and explaining how to run it.
1. A fully working implementation of the project, intended to be readable and intelligible someone doing the project for the first time:
   1. The code should always be correct from a security and safety point of view - this should be model code, don't take short-cuts, and if you need to, explain why, and why it's bad.
   1. Always choose clarity of code as the primary motivator. Favour clarity over efficiency, conciseness, cuteness, cleverness, etc. These are learning examples.
   1. Comment your code extensively. Assume the reader doesn't understand what your code does or why: explain what it's doing, how it's doing it, and most importantly, why it's doing it. If it's important that it's doing something a certain way (e.g. avoiding SQL injection attacks), _explain this_.
   1. From the IMPLEMENTATION.md, and the implementation, it should be clear how to run the code. If dependencies are needed, they should be made clear in one of those locations.

Don't copy the README.md over from the `main` branch. That way, we avoid tricky rebasing issues for everyone.

## How to add a project

If you haven't already discussed the project idea with the team, we recommend [filing an issue](https://github.com/CodeYourFuture/immersive-go-course/issues/new) to discuss it.

Start with learning objectives - in your README.md, write down what you're hoping to teach through the project.

Fill in the rest of your README.md, and any other files needed, and make a pull request against `main`.

Put together a model implementation as described above, and make a pull request against `main`. After review, one of the maintainers will push it to the appropriately named branch.

## Getting help

Feel free to [file an issue](https://github.com/CodeYourFuture/immersive-go-course/issues/new), or if you're on the Code Your Future slack, ping one of [the named authors](https://github.com/CodeYourFuture/immersive-go-course#authors) to get directed to the right channel.

## Code of Conduct

Contributors and users of this repository are expected to adhere to [the Code Your Future Code of Conduct](https://codeyourfuture.io/about/code-of-conduct/).
