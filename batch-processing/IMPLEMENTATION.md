# batch-processing

See https://github.com/CodeYourFuture/immersive-go-course/issues/26 for context.

## Plan

The planned architecture of this:

- Set up the file structure
- Take each row in the input CSV
- Download the associated file
- Run imagemagick to monochrome it
- Upload it to S3

To build this concurrently, we'll use goroutines and channels. Something like this:

- Range over CSV file, spawning goroutines for each line to download the file, with a "done" channel to return the file path (or error?)
- Range over the done downloading channel, spawning goroutines for the resizing, with a done channel
- Range over done channel to upload

Read this as inspiration: https://go.dev/blog/pipelines

This would be unbounded in terms of parallelism (? CPU), so we could tweak this to run a pool of workers with a configurable size. The done channel would be shared: each message would be the next task to complete. The done channel would be buffered to match the size of the list of files.

Will approach this in (at least) two iterations: the first, quite simple, and then the second. Who knows if that's where we'll end up!

## Sketch

Started by sketching out the flow using channels. Turns out I don't know channels very well!

After some iteration, we have a thing working that feels idiomatic. There's minimal use of `sync.WaitGroup` and it's isolated to the `Map` function. I've used `timer.Sleep` with some randomness to simulate interleaving to make sure it handles that OK.

A question will be: should the students also be required to use channels? I think yes. If so, how to teach it? Possible: make it work (linear, blocking) then make it run in parallel (through a concurrent design) using channels.
