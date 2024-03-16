package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

func fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	fmt.Println(reflect.TypeOf(err))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	resp, err := fetch("http://localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	if resp.StatusCode == 200 {
		fmt.Fprint(os.Stdout, sb)
	}

}
