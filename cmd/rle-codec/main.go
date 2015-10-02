// Copyright (c) 2015, xrstf | MIT licensed

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xrstf/rct/rle"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No filename given.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	decoder := rle.NewDecoder()
	result, _ := decoder.DecodeFile(file, 4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(result))
}
