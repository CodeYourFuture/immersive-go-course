package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilteringPipe(t *testing.T) {
	for name, tc := range map[string]struct {
		inputs []string
		output string
	}{
		"no_numbers_in_input": {
			inputs: []string{"hello"},
			output: "hello",
		},
		"just_numbers": {
			inputs: []string{"123"},
			output: "",
		},
		"mixed_numbers_and_letters": {
			inputs: []string{"start=1, end=10"},
			output: "start=, end=",
		},
		"multiple_writes": {
			inputs: []string{"start=", "1, end=10"},
			output: "start=, end=",
		},
	} {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString("")

			fp := NewFilteringPipe(buf)

			for _, input := range tc.inputs {
				n, err := fp.Write([]byte(input))
				require.NoError(t, err)
				require.Equal(t, len(input), n)
			}
			require.Equal(t, tc.output, buf.String())
		})
	}
}
