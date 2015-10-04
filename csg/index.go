// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"io/ioutil"
	"os"

	"github.com/xrstf/rct/utils"
)

type ElementType byte

const (
	DirectBitmapType    ElementType = 0x1
	CompactedBitmapType ElementType = 0x5
	PaletteType         ElementType = 0x8
)

type IndexStruct struct {
	StartAddress uint32
	Width        uint16
	Height       uint16
	XOffset      int16
	YOffset      int16
	Type         ElementType
}

type IndexDecoder struct{}

func NewIndexDecoder() *IndexDecoder {
	return &IndexDecoder{}
}

func (d *IndexDecoder) DecodeFile(file *os.File) ([]IndexStruct, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return d.Decode(content)
}

func (d *IndexDecoder) Decode(index []byte) ([]IndexStruct, error) {
	elements := len(index) / 16
	result := make([]IndexStruct, 0, elements)
	input := utils.NewByteSlice(index)

	for input.At() < input.Size() {
		result = append(result, IndexStruct{
			input.ConsumeUint32(),
			input.ConsumeUint16(),
			input.ConsumeUint16(),
			input.ConsumeInt16(),
			input.ConsumeInt16(),
			ElementType(input.ConsumeByte() & 0x0F),
		})

		// padding
		input.Skip(3)
	}

	return result, nil
}
