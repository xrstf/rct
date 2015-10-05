// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"image"
	"image/color"
)

type Bitmap struct {
	Width  uint16
	Height uint16
	Pixels []byte
}

func (self *Bitmap) ToImage(palette *Palette, remapping RemapSet) image.Image {
	width := int(self.Width)
	height := int(self.Height)
	img := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			key := self.Pixels[y*width+x]
			rgb := color.RGBA{}
			exists := true

			if key >= 0xCA && key <= 0xD5 {
				rgb = remapping.Second.Palette[key-0xCA]
			} else if key >= 0xF3 && key <= 0xFE {
				rgb = remapping.First.Palette[key-0xF3]
			} else {
				rgb, exists = palette.Color(key)
			}

			if exists {
				img.Set(x, y, rgb)
			}
		}
	}

	return img
}
