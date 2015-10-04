// Copyright (c) 2015, xrstf | MIT licensed

package main

import (
	"image/png"
	"log"
	"os"

	"github.com/xrstf/rct/csg"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("No filenames given.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	idxDecoder := csg.NewIndexDecoder()

	index, err := idxDecoder.DecodeFile(file)
	if err != nil {
		log.Fatal(err)
	}

	file, err = os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	graphics, err := csg.NewGraphics(file)
	if err != nil {
		log.Fatal(err)
	}

	bitmap, err := graphics.ExtractCompactedImage(index[61996])
	palette, err := graphics.ExtractPalette(index[2024])

	out, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}

	encoder := png.Encoder{png.BestCompression}
	err = encoder.Encode(out, bitmap.ToImage(palette))
	if err != nil {
		log.Fatal(err)
	}
}
