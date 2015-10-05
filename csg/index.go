// Copyright (c) 2015, xrstf | MIT licensed

//go:generate stringer -type=ElementType -output=index_strings.go

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

type Index struct {
	Elements []IndexStruct
}

func (self *Index) Bitmaps() []IndexStruct {
	return self.filter(DirectBitmapType | CompactedBitmapType)
}

func (self *Index) Palettes() []IndexStruct {
	return self.filter(PaletteType)
}

func (self *Index) filter(t ElementType) []IndexStruct {
	result := make([]IndexStruct, 0)

	for _, element := range self.Elements {
		if element.Type&t > 0 {
			result = append(result, element)
		}
	}

	return result
}

type IndexDecoder struct{}

func NewIndexDecoder() *IndexDecoder {
	return &IndexDecoder{}
}

func (d *IndexDecoder) DecodeFile(file *os.File) (Index, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return Index{}, err
	}

	return d.Decode(content)
}

func (d *IndexDecoder) Decode(data []byte) (Index, error) {
	elements := len(data) / 16
	result := make([]IndexStruct, 0, elements)
	input := utils.NewByteSlice(data)

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

	return Index{result}, nil
}
