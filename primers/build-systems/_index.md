+++
title="Build Systems"
+++

## About this primer {#about-this-document}

This document outlines a short course in build systems - tools used in software engineering to turn source code (and other input files) into derived artifacts (often bundles ready for deployment, or test results). It is aimed at students who have some knowledge of computing and programming and are likely to have used a build tool (perhaps without identifying that this is what they were doing), but without significant professional experience in building, customising, or modifying build systems.

This primer references associated projects which are intended to be completed alongisde to give practical experience with the topics.

### Learning outcomes for this primer:

#### General topics

- [ ] Break a build process down into logical build actions.
- [ ] Identify inputs and outputs of build actions which build and test Java.
- [ ] Sort a dependency graph topologically.
- [ ] Identify the critical path of a dependency graph annotated with individual action execution times.
- [ ] Identify which actions in a dependency graph can be executed in parallel with each other.
- [ ] Compare and contrast the parallelism of building Java and C++.
- [ ] Determine the inputs to an exhaustive cache key for actions building and testing Java.
- [ ] Identify undeclared dependencies in an under-specified build action.
- [ ] Build a cross-machine action cache for a high-trust environment.
- [ ] Explain the trade-offs involved in requiring strict dependency declarations and allowing transitive dependency use.
- [ ] Compare and contrast declarative and imperative build systems.
- [ ] Explain how remote execution can help enforce cache key correctness and avoid undeclared dependencies.
- [ ] Analyse situations when Java header extraction optimises or pessimises overall build time.
- [ ] Explain the trade-offs involved in applicative and monadic builds.

#### `make` as a build tool

- [ ] Write a Makefile to build and test Java.
- [ ] Write a Makefile to build and test C++.
- [ ] Use `make` to analyse a dependency graph for a series of actions which build and test Java.
- [ ] Use `make` to analyse a dependency graph for a series of actions which build and test C++.

#### `bazel` as a build tool

- [ ] Write WORKSPACE.bazel and BUILD.bazel files to build and test Java using Bazel's built-in Java rules.
- [ ] Use `bazel` to analyse a dependency graph for a series of actions which build and test Java.
- [ ] Describe the purposes and differences of Bazel's target and action graphs.
- [ ] Implement a custom (simpified) `java_library`, `java_binary`, and `java_test` rules in Starlark to run in Bazel.
- [ ] Explain how Bazel determines which actions need to run in a clean build.
- [ ] Explain how Bazel determines which actions need to run in an incremental build.
- [ ] Identify undeclared dependencies in a simplified Java ruleset rule in Bazel.
- [ ] Use a toolchain to convert `javac` and `java` from undeclared to a declared dependencies in a simplified Java ruleset in in Bazel.
- [ ] Deploy a remote execution cluster to a cloud computing environment and configure Bazel to use it.
- [ ] Demonstrate that using a toolchain to declare dependencies produces more consistent results across remote execution environments.
- [ ] Reduce the critical path of Java compilation by using [Turbine](https://github.com/google/turbine) to extract headers in environments with large amounts of available parallelism.
