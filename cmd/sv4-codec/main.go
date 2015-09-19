package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xrstf/rct"
)

func main() {
	file, err := os.Open("test.sv4") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := rct.NewRLEDecoder()
	result, err := decoder.Decode(file)

	fmt.Print(string(result))
}
