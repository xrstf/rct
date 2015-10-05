// Copyright (c) 2015, xrstf | MIT licensed

package csg

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

type Palette map[byte]color.RGBA

func (p *Palette) Color(key byte) (color.RGBA, bool) {
	color, exists := (*p)[key]
	return color, exists
}

func (p *Palette) ToImage() image.Image {
	squareSize := 46 // size of one color on the palette
	maxPerRow := 12  // max number of colors in one row
	colors := len(*p)
	width := colors
	height := 1

	if colors > maxPerRow {
		width = maxPerRow
		height = (colors / maxPerRow) + 1
	}

	width *= squareSize // number of squares to pixel
	height *= squareSize

	// map order is random, so sort manually
	keys := make([]int, 0, colors)
	for k := range *p {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// create the image
	img := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for idx, key := range keys {
		paletteEntry := (*p)[byte(key)]

		// logical coordinates, not pixel values
		x := idx % maxPerRow
		y := idx / maxPerRow

		// draw the background
		xOffset := x * squareSize
		yOffset := y * squareSize

		for i := yOffset; i < yOffset+squareSize; i++ {
			for j := xOffset; j < xOffset+squareSize; j++ {
				img.Set(j, i, paletteEntry)
			}
		}

		// label the color
		text := fmt.Sprintf("0x%02X", byte(key))
		c := color.RGBA{
			255 - paletteEntry.R,
			255 - paletteEntry.G,
			255 - paletteEntry.B,
			255,
		}

		drawText(xOffset+3, yOffset+20, text, img, c)
	}

	return img
}

func drawText(x int, y int, text string, img *image.NRGBA, color color.RGBA) {
	charWidth := 8
	charHeight := 7

	for _, character := range []byte(text) {
		pattern := bitmaps[character]

		for i := 0; i < charHeight; i++ {
			for j := 0; j < charWidth; j++ {
				pixel := pattern[i*charWidth+j]

				if pixel > 0 {
					img.Set(x+j, y+i, color)
				}
			}
		}

		// we drew a character, so advance
		x += charWidth

		// spacing in between each character
		x++
	}
}

var bitmaps = map[byte][]uint8{
	'0': stringsToBitmap([]string{
		" ###### ",
		"###   ##",
		"####  ##",
		"## ## ##",
		"##  ####",
		"##   ###",
		" ###### ",
	}),

	'1': stringsToBitmap([]string{
		"   ##   ",
		" ####   ",
		"## ##   ",
		"   ##   ",
		"   ##   ",
		"   ##   ",
		" ###### ",
	}),

	'2': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"      ##",
		" ###### ",
		"##      ",
		"##      ",
		"########",
	}),

	'3': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"      ##",
		"  ##### ",
		"      ##",
		"##    ##",
		" ###### ",
	}),

	'4': stringsToBitmap([]string{
		"##    ##",
		"##    ##",
		"##    ##",
		"########",
		"      ##",
		"      ##",
		"      ##",
	}),

	'5': stringsToBitmap([]string{
		"####### ",
		"##      ",
		"##      ",
		"####### ",
		"      ##",
		"##    ##",
		" ###### ",
	}),

	'6': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"##      ",
		"####### ",
		"##    ##",
		"##    ##",
		" ###### ",
	}),

	'7': stringsToBitmap([]string{
		"########",
		"     ## ",
		"    ##  ",
		"   ##   ",
		"  ##    ",
		" ##     ",
		"##      ",
	}),

	'8': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"##    ##",
		" ###### ",
		"##    ##",
		"##    ##",
		" ###### ",
	}),

	'9': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"##    ##",
		" #######",
		"      ##",
		"##    ##",
		" ###### ",
	}),

	'A': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"##    ##",
		"########",
		"##    ##",
		"##    ##",
		"##    ##",
	}),

	'B': stringsToBitmap([]string{
		"####### ",
		"##    ##",
		"##    ##",
		"####### ",
		"##    ##",
		"##    ##",
		"####### ",
	}),

	'C': stringsToBitmap([]string{
		" ###### ",
		"##    ##",
		"##      ",
		"##      ",
		"##      ",
		"##    ##",
		" ###### ",
	}),

	'D': stringsToBitmap([]string{
		"####### ",
		"##    ##",
		"##    ##",
		"##    ##",
		"##    ##",
		"##    ##",
		"####### ",
	}),

	'E': stringsToBitmap([]string{
		"########",
		"##      ",
		"##      ",
		"#####   ",
		"##      ",
		"##      ",
		"########",
	}),

	'F': stringsToBitmap([]string{
		"########",
		"##      ",
		"##      ",
		"#####   ",
		"##      ",
		"##      ",
		"##      ",
	}),

	'x': stringsToBitmap([]string{
		"        ",
		"        ",
		"##    ##",
		" ##  ## ",
		"  ####  ",
		" ##  ## ",
		"##    ##",
	}),
}

func stringsToBitmap(strs []string) []uint8 {
	result := make([]uint8, len(strs)*len(strs[0]))
	idx := 0

	for _, str := range strs {
		for _, character := range []byte(str) {
			if character != ' ' {
				result[idx] = 1
			}

			idx++
		}
	}

	return result
}
