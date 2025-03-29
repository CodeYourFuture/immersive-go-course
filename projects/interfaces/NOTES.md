I want to point something that I discovered while making this project, and is the deletion of elements within a slice.

So, while I was figuring out how to achieve the reader method of my own bytes.Buffer struct, I wanted to be able to update
the status of the buffer it was read from (my bytes.Buffer struct) so every time it was read from it, the amount of bytes 
remaining were the original amount - the amount of bytes read from it.

Doing some research I found out the slices.Delete function which was just the tool I needed to update the status of the bufferand at the same time to zero out the obsolete elements, that you may think are totally invisible but the truth is that they are still using memory.

Link: https://go.dev/blog/generic-slice-functions
