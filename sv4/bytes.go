package sv4

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func rol32(x uint32, shift uint) uint32 {
	return ((x << shift) | (x >> (32 - shift)))
}

func uint32ToBytes(value int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))

	return buf
}

// Reads a single byte at pos.
func (s *SaveState) readByte(pos uint32) byte {
	return s.data[pos]
}

// Reads num bytes starting (including) from pos.
func (s *SaveState) readBytes(pos uint32, num uint32) []byte {
	return s.data[pos:(pos + num)]
}

// Read an unsigned 8-bit integer at pos.
func (s *SaveState) readUint8(pos uint32) uint8 {
	return uint8(s.readByte(pos))
}

// Read an unsigned 16-bit integer at pos.
func (s *SaveState) readUint16(pos uint32) uint16 {
	return binary.LittleEndian.Uint16(s.readBytes(pos, 2))
}

// Read an signed 16-bit integer at pos.
func (s *SaveState) readInt16(pos uint32) int16 {
	b := s.readBytes(pos, 2)
	buf := bytes.NewReader(b)
	result := int16(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}

// Read an unsigned 32-bit integer at pos.
func (s *SaveState) readUint32(pos uint32) uint32 {
	return binary.LittleEndian.Uint32(s.readBytes(pos, 4))
}

// Read an signed 32-bit integer at pos.
func (s *SaveState) readInt32(pos uint32) int32 {
	b := s.readBytes(pos, 4)
	buf := bytes.NewReader(b)
	result := int32(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}
