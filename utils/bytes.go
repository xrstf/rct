// Copyright (c) 2015, xrstf | MIT licensed

package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type ByteSlice struct {
	data   []byte
	cursor uint32
}

func NewByteSlice(bytes []byte) *ByteSlice {
	return &ByteSlice{bytes, 0}
}

func (self *ByteSlice) At() uint32 {
	return self.cursor
}

func (self *ByteSlice) Size() uint32 {
	return uint32(len(self.data))
}

func (self *ByteSlice) Skip(num uint32) error {
	return self.Seek(self.cursor + num)
}

func (self *ByteSlice) Seek(pos uint32) error {
	if pos > uint32(len(self.data)) {
		return errors.New(fmt.Sprintf("Target position %d is out of range (slice has %d bytes)!", pos, len(self.data)))
	}

	self.cursor = pos

	return nil
}

// Reads a single byte at pos.
func (self *ByteSlice) ReadByte(pos uint32) byte {
	return self.data[pos]
}

// Reads num bytes starting (including) from pos.
func (self *ByteSlice) ReadBytes(pos uint32, num uint32) []byte {
	return self.data[pos:(pos + num)]
}

// Read an unsigned 8-bit integer at pos.
func (self *ByteSlice) ReadUint8(pos uint32) uint8 {
	return uint8(self.ReadByte(pos))
}

// Read an unsigned 16-bit integer at pos.
func (self *ByteSlice) ReadUint16(pos uint32) uint16 {
	return binary.LittleEndian.Uint16(self.ReadBytes(pos, 2))
}

// Read an signed 16-bit integer at pos.
func (self *ByteSlice) ReadInt16(pos uint32) int16 {
	b := self.ReadBytes(pos, 2)
	buf := bytes.NewReader(b)
	result := int16(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}

// Read an unsigned 32-bit integer at pos.
func (self *ByteSlice) ReadUint32(pos uint32) uint32 {
	return binary.LittleEndian.Uint32(self.ReadBytes(pos, 4))
}

// Read an signed 32-bit integer at pos.
func (self *ByteSlice) ReadInt32(pos uint32) int32 {
	b := self.ReadBytes(pos, 4)
	buf := bytes.NewReader(b)
	result := int32(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}

// The following functions are useful when processing a byte slice in a linear fashion.

func (self *ByteSlice) advance(steps uint32) {
	self.cursor = self.cursor + steps
}

func (self *ByteSlice) ConsumeByte() byte {
	return self.ConsumeBytes(1)[0]
}

func (self *ByteSlice) ConsumeBytes(num uint32) []byte {
	defer self.advance(num)
	return self.ReadBytes(self.cursor, num)
}

func (self *ByteSlice) ConsumeUint8() uint8 {
	defer self.advance(1)
	return self.ReadUint8(self.cursor)
}

func (self *ByteSlice) ConsumeUint16() uint16 {
	defer self.advance(2)
	return self.ReadUint16(self.cursor)
}

func (self *ByteSlice) ConsumeInt16() int16 {
	defer self.advance(2)
	return self.ReadInt16(self.cursor)
}

func (self *ByteSlice) ConsumeUint32() uint32 {
	defer self.advance(4)
	return self.ReadUint32(self.cursor)
}

func (self *ByteSlice) ConsumeInt32() int32 {
	defer self.advance(4)
	return self.ReadInt32(self.cursor)
}
