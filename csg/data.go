// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"errors"
	"image"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/xrstf/rct/utils"
)

type Palette map[byte]color.RGBA

type Bitmap struct {
	Width  uint16
	Height uint16
	Pixels []byte
}

func (self *Bitmap) ToImage(palette *Palette) image.Image {
	width := int(self.Width)
	height := int(self.Height)
	img := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			key := self.Pixels[y*width+x]

			color, exists := (*palette)[key]
			if exists {
				img.Set(x, y, color)
			}
		}
	}

	return img
}

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

func (self *Graphics) ExtractDirectBitmap(index IndexStruct) (*Bitmap, error) {
	if index.Type != DirectBitmapType {
		return nil, errors.New("The given index structure does not point to a direct bitmap.")
	}

	data := utils.NewByteSlice(self.data)
	err := data.Seek(index.StartAddress)
	if err != nil {
		return nil, err
	}

	width := index.Width
	height := index.Height
	imgData := data.ConsumeBytes(uint32(width * height))

	return &Bitmap{width, height, imgData}, nil
}

func (self *Graphics) ExtractCompactedImage(index IndexStruct) (*Bitmap, error) {
	if index.Type != CompactedBitmapType {
		return nil, errors.New("The given index structure does not point to a direct bitmap.")
	}

	data := utils.NewByteSlice(self.data)
	err := data.Seek(index.StartAddress)
	if err != nil {
		return nil, err
	}

	width := index.Width
	height := index.Height
	imgData := make([]byte, width*height)

	// taken from https://github.com/LinusU/node-rct-graphics/blob/6efe3864a93ed/src/api.js#L19
	for y := uint16(0); y < height; y++ {
		isLast := false
		startOffset := data.ReadUint16(index.StartAddress + (2 * uint32(y)))
		currentPos := index.StartAddress + uint32(startOffset)

		for !isLast {
			size := data.ReadUint8(currentPos)
			offset := data.ReadUint8(currentPos + 1)

			currentPos += 2

			isLast = size&0x80 > 0
			size = size & 0x7F

			start := y*width + uint16(offset)
			copy(imgData[start:(start+uint16(size))], data.ReadBytes(currentPos, uint32(size)))

			currentPos += uint32(size)
		}
	}

	return &Bitmap{width, height, imgData}, nil
}

func (self *Graphics) ExtractPalette(index IndexStruct) (*Palette, error) {
	if index.Type != PaletteType {
		return nil, errors.New("The given index structure does not point to a palette.")
	}

	data := utils.NewByteSlice(self.data)
	err := data.Seek(index.StartAddress)
	if err != nil {
		return nil, err
	}

	p := make(Palette)
	idx := byte(uint8(index.XOffset))

	for i := uint16(0); i < index.Width; i++ {
		blue := data.ConsumeByte()
		green := data.ConsumeByte()
		red := data.ConsumeByte()

		p[idx] = color.RGBA{red, green, blue, 255}
		idx++
	}

	return &p, nil
}
