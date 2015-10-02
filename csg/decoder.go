// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xrstf/rct/utils"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeFile(file *os.File) ([]byte, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return d.Decode(content)
}

func (d *Decoder) Decode(encoded []byte) ([]byte, error) {
	result := make([]byte, 0)
	input := utils.NewByteSlice(encoded)

	for input.At() < input.Size() {
		startAddress := input.ConsumeUint32()
		width := input.ConsumeInt16()
		height := input.ConsumeInt16()
		xoffset := input.ConsumeInt16()
		yoffset := input.ConsumeInt16()
		flags := input.ConsumeByte()

		// padding
		input.Skip(3)

		fmt.Printf("start @ % 8X, width = % 3d, height = % 3d, xoffset = % 4d, yoffset = % 4d, flags = %X\n", startAddress, width, height, xoffset, yoffset, flags & 0x0F)
	}

	return result, nil
}
