// Copyright (c) 2015, xrstf | MIT licensed

package main

import (
	"image/png"
	"log"
	"os"

	"github.com/xrstf/rct/tp4"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No filename given.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	decoder := tp4.NewDecoder()
	img, err := decoder.DecodeFile(file)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}

	encoder := png.Encoder{png.BestCompression}
	err = encoder.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
