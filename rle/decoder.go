// Copyright (c) 2015, xrstf | MIT licensed

// Package rle implements encoder and decoder for RLE encoded data.
//
// RLE (run-length encoding) is a very simple compression algorithm that shrinks multiple
// consecutive identical bytes by placing the count and one example.
//
// In the Rollercoaster Tycoon world, this encoding is sometimes referred to as
// "sawyer coding".
//
// For more information, see http://tid.rctspace.com/RLE.html.
package rle

import (
	"bytes"
	"io/ioutil"
	"os"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeFile(file *os.File, checksumLen int) ([]byte, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	length := len(content)

	return d.Decode(content[0:(length - checksumLen)])
}

func (d *Decoder) Decode(encoded []byte) ([]byte, error) {
	size := len(encoded)
	end := 0
	result := make([]byte, 0)

	for pos := 0; pos < size; {
		cmd := encoded[pos]

		pos = pos + 1

		// we have (cmd+1) literal bytes to copy from the input
		if cmd < 128 {
			end = pos + int(cmd) + 1

			result = append(result, encoded[pos:end]...)

			pos = end
		} else {
			next := encoded[pos]

			result = append(result, bytes.Repeat([]byte{next}, int(byte(1)-cmd))...)

			pos = pos + 1
		}
	}

	return result, nil
}
