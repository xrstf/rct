// Copyright (c) 2015, xrstf | MIT licensed

package tp4

import (
	"errors"
	"image"
	"io/ioutil"
	"os"
	"strconv"
)

const FileSize = 52000

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeFile(file *os.File) (image.Image, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(content) != FileSize {
		return nil, errors.New("File size must be exactly " + strconv.Itoa(FileSize) + " bytes.")
	}

	width := 254
	height := 200
	pos := 400 // skip 400 header bytes
	idx := 0

	m := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if idx%127 == 0 {
				pos = pos + 2
			}

			m.Set(x, y, ColorPalette[content[pos]])

			pos++
			idx++
		}
	}

	return m, nil
}
