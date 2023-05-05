# Deployment

This section will guide you through the process of deploying your code to the cloud using GitHub Actions. By the end of this module, you'll understand the concept of deployment pipelines, various environments, and testing along the way.

## Introduction

When you develop a software application, the code is initially written on your local computer. To make this code accessible to others you need to **deploy** it. This process of getting code to users is often referred to as a **deployment pipeline**.

### What is Deployment?

The meanings can vary greatly depending on what you are trying to do. Imagine a very simple use case, for example deploying a simple front-end website using HTML, CSS and Javascript. In this instance, deployment may be as simple as uploading the files somewhere publically available and ensuring that your URL is pointing at it.

If we want to also deploy a Node.js API, things become a bit more complicated. Yes, we need to upload our front-end files as before but remember our Node.js API is always running, listening to requests and responding. We would need our back end to actually run somewhere in the cloud, so we would provision a Cloud Compute platform (like AWS EC2 or Heroku). We also need to make sure the front end and the back end are configured to talk to each other, which will probably involve secrets/api keys, which means we may also need to talk to a Secrets Manager.

Now imagine we have an email server to send emails, a remote database and database backups, scripts that runs automatically at certain times of the day, slack integration, the list goes on. A complex deployment can involve a huge array of services, all configured to talk to each other securely.

In this section we'll run you through the key concepts of a simple deployment.

## Deploying From Local

Before we start, let's remind ourselves why any of this is necessary. We have written some code on our computer and our goal is to put it somewhere that users can access it. So can we just upload it from our computer directly to the cloud? The answer is yes. You can set up an account with a cloud provider and upload files directly using tools like FTP or SCP.

![Illustration of Local Computer deploying directly to the Cloud](https://i.imgur.com/BtPPm7K.png)

In the image above you can see the simplest possible process to deploy. Run a command on your local machine and deploy your code to the cloud for users. This approach is very easy to do but it has some drawbacks.

### Multiple developers

Imagine we have two developers working on the same codebase. They both have the ability to deploy code directly to the cloud from their computer. How do they know their code base is up to date, do they have the latest changes from their team-mate, if they both deploy code at the same time, will one over-write the other? Who pushed their changes first? Did they definitely do a Git Pull before they deployed?

![enter image description here](https://i.imgur.com/wOYLGA6.png)

We run a real risk of our cloud deployment becoming out of sync with GitHub the changes from both developers. We have no **Single Source of Truth**.

In the diagram we have illustrated that both developers are also pushing their changes to GitHub. Now we also know that GitHub does a really good job of tracking changes to a codebase and ensuring everything stays in-sync and unconflicted. Is there a way to push our code to GitHub and allow GitHub to deploy it for us?

## GitHub Actions

If we put GitHub in the middle of our deployment we no longer need to worry about the risk of things getting out of Sync. Our ideal setup would be to just push our changes to GitHub and have it automatically deployed to the Cloud every time a new commit is pushed.

![enter image description here](https://i.imgur.com/nBEi24n.png)

[GitHub Actions](https://github.com/features/actions) is a powerful automation tool that allows you to deploy code directly within your repository. It enables you to automate, customize, and execute your software development workflows right in your repo.

### Deployment Pipeline

A deployment pipeline is a series of automated steps that take the code from your computer to a cloud. These steps can include building, testing, and deploying the code to different environments. Let's take a look at a sample **workflow file**. GitHub reads this file every time you push your code and follows the steps one by one. Don't worry about memorising this now, you'll have a deployment challenge at the end of this section. For now, just take a look at the workflow file and try to guess what it will do.

```yaml
# .github/workflows/deployment.yml

name: Deployment Pipeline
on:
  push:
    branches:
      - main # only run this if the developer pushes to the main branch

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps: aws deploy # example deployment command
```

## Environments

In the previous digram, you can see that we've replaced the word 'Cloud' with the word 'Production'. Why is this? It turns out that it's common to have multiple versions of your app in the cloud at the same time. The "main version", the one that users see, is usually known as your 'Production' environment. However there are good reasons to put other versions of your app in the cloud.

The most common environment that you will see is known as the 'Staging' environment. The purpose of Staging is to allow the team to preview any changes before releasing it to users. Here our workflow might look something like this:

![enter image description here](https://i.imgur.com/UhNdaRF.png)

Two developers are both pushing code to GitHub. Every time new code comes in, GitHub will automatically deploy it to the Staging Environment. Once we are confident that we are happy with the changes, we can then in turn upload it to the Production environment.

We have looked at three different environments:

- **Local Environment**: This is the code running on your computer
- **Production Environment**: The code that's seen by users
- **Staging Environment**: A preview of the production environment that the team can use to determine that there are no bugs before releasing it to users.

Let's look at one more common type of environment. First let's take a deeper look at GitHub

### Branch Environments

The diagram below is the same as the previous diagram, except GitHub has been expanded out so we can see what's happening in a bit more detail.

![enter image description here](https://i.imgur.com/RfYlONR.png)

As we have learnt, the developers both push their code to GitHub, and GitHub automatically deploys it to the Cloud. Usually to the Staging environment.

We know that Github has the concept of branches, Let's take a look deeper into the internals of GitHub to see what is happening when two developers both push code to the same repository. In the diagram below we can see the main branch (in blue) and two other branches (in black). It is common for each developer to work in their own branch and then merge their branch into the main branch once their feature is complete.

![enter image description here](https://i.imgur.com/hjoI82B.png)

Let's now update this diagram to show where **Branch Environments** fit in.

![enter image description here](https://i.imgur.com/s7zdGoz.png)

Every time a developer pushes code to a branch, GitHub will automatically deploy the entire app into a cloud environment, with just the code for that branch. It means that developers can preview their changes before they merge to the main branch.

## Deployment Strategies

There are different strategies for how often to deploy code. Frequency of releases and stages of testing vary from team to team. Below is a common framework that will seem familiar to a lot of software development teams.

![enter image description here](https://i.imgur.com/YccRIOP.png)

### Challenge

Follow the [Create an example workflow task](https://docs.github.com/en/actions/learn-github-actions/understanding-github-actions) in the GitHub Actions tutorial.
