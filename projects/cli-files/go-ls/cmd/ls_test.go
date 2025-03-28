package cmd

import (
	"os/exec"
	"testing"
)

func TestGoLsCommand(t *testing.T) {
   t.Run("Testing against original ls command", func(t *testing.T) {
        cases := []struct {
            lsCommand       string
            golsCommand     string 
        }{
            { "ls", "go-ls" },
            { "ls ..", "go-ls .." },
            { "ls ../assets", "go-ls ../assets" },
            { "ls ../cmd", "go-ls ../cmd "},
        }

        for _, test := range cases {
            lsOutput, _ := exec.Command(test.lsCommand).Output()
            golsOutput, _ := exec.Command(test.golsCommand).Output()
            
            if string(lsOutput) != string(golsOutput) {
                t.Errorf("expected output: %s got: %s", lsOutput, golsOutput)
            }


        }
   }) 
}
