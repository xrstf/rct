// Copyright (c) 2015, xrstf | MIT licensed

package rle

import "errors"

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

type encoderMode byte

const (
	modeUnknown encoderMode = iota
	modeConsecutive
	modeIndividuals
)

type encoderStack struct {
	mode  encoderMode
	stack []byte
}

func newEncoderStack(mode encoderMode, initial byte) *encoderStack {
	return &encoderStack{mode, []byte{initial}}
}

func (s *encoderStack) size() int {
	return len(s.stack)
}

func (s *encoderStack) end() byte {
	return s.stack[len(s.stack)-1]
}

func (s *encoderStack) push(b byte) {
	s.stack = append(s.stack, b)
}

func (s *encoderStack) pop() byte {
	l := len(s.stack)
	last := s.stack[l-1]

	s.stack = s.stack[:l-1]

	return last
}

func (s *encoderStack) flush() []byte {
	size := uint8(len(s.stack) - 1)

	if s.mode == modeConsecutive {
		controlByte := byte(-size) | 0x80

		return []byte{controlByte, s.stack[0]}
	}

	controlByte := byte(size)

	return append([]byte{controlByte}, s.stack...)
}

// Encodes (compresses) a byte slices
func (d *Encoder) Encode(raw []byte) ([]byte, error) {
	size := len(raw)

	if size == 0 {
		return nil, errors.New("Cannot encode zero bytes.")
	}

	prev := raw[0]
	stack := newEncoderStack(modeUnknown, prev)
	result := make([]byte, 0)

	for i := 1; i < size; i++ {
		current := raw[i]

		if stack.mode == modeUnknown { // this basically means that the stack only has one byte yet
			stack.push(current)

			if current == prev {
				stack.mode = modeConsecutive
			} else {
				stack.mode = modeIndividuals
			}
		} else if stack.mode == modeConsecutive {
			if current == prev { // continue our streak
				stack.push(current)
			} else { // bail out
				result = append(result, stack.flush()...)
				stack = newEncoderStack(modeUnknown, current)
			}
		} else {
			if current == prev { // we found two identical bytes, switch to consecutive mode
				// remove the start of the streak from the stack
				stack.pop()

				result = append(result, stack.flush()...)
				stack = newEncoderStack(modeConsecutive, prev)
			}

			stack.push(current)
		}

		// respect RCT1's internal limit
		if stack.size() >= 125 {
			result = append(result, stack.flush()...)

			// start a new stack if we aren't on the end yet
			if i < (size - 1) {
				i++
				stack = newEncoderStack(modeUnknown, raw[i])
			}
		}

		prev = raw[i]
	}

	result = append(result, stack.flush()...)

	return result, nil
}
