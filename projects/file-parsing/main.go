package main

import (
	"fmt"
	"go-file-parsing/parser"
	"os"
)

func main() {
	file, _ := os.Open("examples/json.txt")
	defer file.Close()
	players, _ := parser.ParseJson(file)
	fmt.Println(players)

	file2, _ := os.Open("examples/repeated-json.txt")
	defer file2.Close()
	players, _ = parser.ParseRepeatedJson(file2)
	fmt.Println(players)

	file3, _ := os.Open("examples/data.csv")
	defer file3.Close()
	players, _ = parser.ParseCsv(file3)
	fmt.Println(players)

	file4, _ := os.Open("examples/custom-binary-be.bin")
	defer file4.Close()
	players, _ = parser.ParseCustomBinaryBigEndian(file4)
	fmt.Println(players)

	file5, _ := os.Open("examples/custom-binary-le.bin")
	defer file5.Close()
	players, _ = parser.ParseCustomBinaryLittleEndian(file5)
	fmt.Println(players)
}
