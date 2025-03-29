package interfaces

import (
	"io"
	"strconv"
)

type FilteringPipe struct {
    Writer      io.Writer
}

func(f *FilteringPipe) Write(p []byte) (n int, err error) {
    for _, c := range p {
        if _, err := strconv.Atoi(string(c)); err != nil {
            continue
        }
        if _, err := f.Writer.Write([]byte{c}); err != nil {
            return 0, err
        }
    }
    return len(p), nil
}
