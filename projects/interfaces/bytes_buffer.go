package interfaces

import (
	"io"
	"slices"
)

type OurByteBuffer struct {
    Available       []byte
}

// Return number of bytes available in buffer.
func (b *OurByteBuffer) Len() (bytes int) {
   return len(b.Available) 
}

func(b *OurByteBuffer) Write (c []byte) {
   b.Available = append(b.Available, c...) 
}
func(b *OurByteBuffer) Bytes () []byte {
    return b.Available
}
func(b *OurByteBuffer) Read (c []byte) (int, error) {
    // Check for and empty buffer slice.
    if len(b.Available) == 0 {
        return 0, io.EOF
    }
    // Copy elements and set to zero the obsolet elements from the buffer new slice.
    elemsCopied := copy(c, b.Available)
    b.Available = slices.Delete(b.Available, 0, len(c))

    return elemsCopied, nil
}
