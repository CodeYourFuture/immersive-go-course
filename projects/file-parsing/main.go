package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/go-node"
)

func main() {
	// go run . -fx examples/json.txt '.sort((l, r) => r.high_score - l.high_score)[0].name'
	// go run . -fx examples/json.txt '.sort((l, r) => l.high_score - r.high_score)[0].name'
	// go run . -csvq examples/data.csv 'SELECT * ORDER BY `high score` DESC LIMIT 1'
	// go run . -csvq examples/data.csv 'SELECT * ORDER BY `high score` ASC LIMIT 1'
	fx := flag.Bool("fx", false, "Implementing -fx")
	csvq := flag.Bool("csvq", false, "Implementing -csvq")
	flag.Parse()
	args := flag.Args()

	file, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read: %s", err)
		return
	}
	if *fx {
		vm := node.New(nil)
		jsCommand := fmt.Sprint(string(file), args[1])
		result := vm.Run(jsCommand)
		println(result.String())
	}

	if *csvq {
		r := csv.NewReader(strings.NewReader(string(file)))
		records, err := r.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't read: %s", err)
		}

		resultName := records[1][0]
		resultScore, _ := strconv.Atoi(records[1][1])
		for i := 2; i < len(records); i++ {
			score, _ := strconv.Atoi(records[i][1])
			if strings.Contains(args[1], "DESC") && score > resultScore ||
				strings.Contains(args[1], "ASC") && score < resultScore {
				resultScore = score
				resultName = records[i][0]
			}
		}
		fmt.Println("+--------+------------+\n|  name  | high score |\n+--------+------------+\n|", resultName, "|", resultScore, "       |\n+--------+------------+")
	}
}
