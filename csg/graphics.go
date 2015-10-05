// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"errors"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/xrstf/rct/utils"
)

type Graphics struct {
	data []byte // avoid sharing the ByteSlice struct, so we have independent cursors per goroutine
}

func NewGraphics(dataFile *os.File) (*Graphics, error) {
	content, err := ioutil.ReadAll(dataFile)
	if err != nil {
		return nil, err
	}

	return &Graphics{content}, nil
}

func (self *Graphics) ExtractBitmap(index IndexStruct) (*Bitmap, error) {
	switch index.Type {
	case DirectBitmapType:
		return self.extractDirectBitmap(index)

	case CompactedBitmapType:
		return self.extractCompactedImage(index)

	default:
		return nil, errors.New("The given index struct is not pointing to a bitmap.")
	}
}

func (self *Graphics) extractDirectBitmap(index IndexStruct) (*Bitmap, error) {
	reader, err := self.reader(index)
	if err != nil {
		return nil, err
	}

	width := index.Width
	height := index.Height
	imgData := reader.ConsumeBytes(uint32(width * height))

	return &Bitmap{width, height, imgData}, nil
}

func (self *Graphics) extractCompactedImage(index IndexStruct) (*Bitmap, error) {
	reader, err := self.reader(index)
	if err != nil {
		return nil, err
	}

	width := index.Width
	height := index.Height
	imgData := make([]byte, width*height)

	// taken from https://github.com/LinusU/node-rct-graphics/blob/6efe3864a93ed/src/api.js#L19
	for y := uint16(0); y < height; y++ {
		isLast := false
		startOffset := reader.ReadUint16(index.StartAddress + (2 * uint32(y)))
		currentPos := index.StartAddress + uint32(startOffset)

		for !isLast {
			size := reader.ReadUint8(currentPos)
			offset := reader.ReadUint8(currentPos + 1)

			currentPos += 2

			isLast = size&0x80 > 0
			size = size & 0x7F

			start := y*width + uint16(offset)
			copy(imgData[start:(start+uint16(size))], reader.ReadBytes(currentPos, uint32(size)))

			currentPos += uint32(size)
		}
	}

	return &Bitmap{width, height, imgData}, nil
}

func (self *Graphics) ExtractPalette(index IndexStruct) (*Palette, error) {
	if index.Type != PaletteType {
		return nil, errors.New("The given index structure does not point to a palette.")
	}

	reader, err := self.reader(index)
	if err != nil {
		return nil, err
	}

	p := make(Palette)
	idx := byte(uint8(index.XOffset))

	for i := uint16(0); i < index.Width; i++ {
		blue := reader.ConsumeByte()
		green := reader.ConsumeByte()
		red := reader.ConsumeByte()

		p[idx] = color.RGBA{red, green, blue, 255}
		idx++
	}

	return &p, nil
}

func (self *Graphics) reader(index IndexStruct) (*utils.ByteSlice, error) {
	data := utils.NewByteSlice(self.data)
	err := data.Seek(index.StartAddress)
	if err != nil {
		return nil, err
	}

	return data, nil
}
