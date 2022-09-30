package main

import "sync"

// Map runs a function on all values from an input channel and publishes them to a new output channel that it returns.
// After the input channel closes, Map will block until work is complete and then close the channel.
func Map[In any, Out any](in <-chan In, work func(in In) Out) chan Out {
	// We'll push out results on this channel
	out := make(chan Out)
	// This WaitGroup tracks work-in-progress which will be allowed to complete before we close the output channel
	var wg sync.WaitGroup
	// Kick off a goroutine so we can return the output channel (and don't block)
	go func() {
		// Iterate through every value on the input channel
		for inV := range in {
			// Track work in progress
			wg.Add(1)
			// Run the work in another goroutine so again we don't block
			go func(inV In) {
				// Do the work and put the value onto the output channel
				out <- work(inV)
				// Work is complete
				wg.Done()
			}(inV)
		}
		// Make sure all tracked work is complete. We will only reach here when the input channel closes because
		// the `for x := range c` for tracks the closed-state of the channel it is ranging over.
		wg.Wait()
		// We're done! Close the channel.
		close(out)
	}()
	// Return the output channel so the caller can do something with it (like Map again!)
	return out
}
