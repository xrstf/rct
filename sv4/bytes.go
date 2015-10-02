// Copyright (c) 2015, xrstf | MIT licensed

package sv4

import "encoding/binary"

func rol32(x uint32, shift uint) uint32 {
	return ((x << shift) | (x >> (32 - shift)))
}

func uint32ToBytes(value int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))

	return buf
}

// convenience wrappers

func (s *SaveState) readByte(pos uint32) byte {
	return s.data.ReadByte(pos)
}

func (s *SaveState) readBytes(pos uint32, num uint32) []byte {
	return s.data.ReadBytes(pos, num)
}

func (s *SaveState) readUint8(pos uint32) uint8 {
	return s.data.ReadUint8(pos)
}

func (s *SaveState) readUint16(pos uint32) uint16 {
	return s.data.ReadUint16(pos)
}

func (s *SaveState) readInt16(pos uint32) int16 {
	return s.data.ReadInt16(pos)
}

func (s *SaveState) readUint32(pos uint32) uint32 {
	return s.data.ReadUint32(pos)
}

func (s *SaveState) readInt32(pos uint32) int32 {
	return s.data.ReadInt32(pos)
}
