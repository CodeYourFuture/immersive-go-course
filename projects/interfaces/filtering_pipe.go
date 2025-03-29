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
        f.Writer.Write([]byte{c})
    }
    return 0, nil
}
