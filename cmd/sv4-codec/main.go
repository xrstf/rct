package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xrstf/rct"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No filename given.")
	}

	// for {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	decoder := rct.NewRLEDecoder()
	result, _ := decoder.Decode(file)

	saveState, _ := rct.ParseSaveState(result)

	fmt.Printf("save state = %+v\n", saveState)

	file.Close()
	// 	<-time.After(1 * time.Second)
	// }
}
