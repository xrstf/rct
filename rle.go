package rct

import (
	"bytes"
	"io/ioutil"
	"os"
)

type RLEDecoder interface {
	Decode(*os.File) ([]byte, error)
}

func NewRLEDecoder() RLEDecoder {
	return &dummy{}
}

type dummy struct{}

func (d *dummy) Decode(file *os.File) ([]byte, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	size := len(content)
	end := 0
	result := make([]byte, 0, SaveStateSize)

	for pos := 0; pos < (size - 4); { // -4 to ignore the trailing checksum
		cmd := content[pos]

		pos = pos + 1

		// we have (cmd+1) literal bytes to copy from the input
		if cmd < 128 {
			end = pos + int(cmd) + 1

			result = append(result, content[pos:end]...)

			pos = end
		} else {
			next := content[pos]

			result = append(result, bytes.Repeat([]byte{next}, int(byte(1)-cmd))...)

			pos = pos + 1
		}
	}

	return result, nil
}
