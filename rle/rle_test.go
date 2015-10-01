// Copyright (c) 2015, xrstf | MIT licensed

package rle

import (
	"fmt"
	"strings"
	"testing"
)

type testcase struct {
	input    string
	expected string
}

func TestEncodingDecoding(t *testing.T) {
	testcases := []testcase{
		{"00", "00 00"},                                                      // ind(1, [0])
		{"FF", "00 FF"},                                                      // ind(1, [255])
		{"00 00", "FF 00"},                                                   // con(2, 0)
		{"00 00 00", "FE 00"},                                                // con(3, 0)
		{"01 02 03", "02 01 02 03"},                                          // ind(3, [1,2,3])
		{"01 01 01 02 02 03", "FE 01 FF 02 00 03"},                           // con(3, 1) + con(2, 2) + ind(1, [3])
		{"01 01 02 03 03", "FF 01 00 02 FF 03"},                              // con(2, 1) + ind(1, [2]) + con(2, 3)
		{strings.TrimSpace(strings.Repeat("00 ", 300)), "84 00 84 00 CF 00"}, // con(125, 0) +  con(125, 0) +  con(50, 0)
	}

	encoder := NewEncoder()
	decoder := NewDecoder()

	for _, test := range testcases {
		encoded, err := encoder.Encode(strToByteSlice(test.input))

		if err != nil {
			t.Error(err)
		}

		assertByteSliceEqual(t, test.input, encoded, test.expected)

		// are we compatible to our own output?
		decoded, err := decoder.Decode(encoded)

		if err != nil {
			t.Error(err)
		}

		assertDecodedByteSliceEqual(t, encoded, decoded, test.input)
	}
}

func strToByteSlice(str string) []byte {
	result := make([]byte, 0)

	fmt.Sscanf(strings.Replace(str, " ", "", -1), "%x", &result)

	return result
}

func assertByteSliceEqual(t *testing.T, input string, result []byte, expected string) {
	output := fmt.Sprintf("% X", result)

	if output != expected {
		t.Errorf("Encoded result does not meet the expectation.\nInput...: %s\nExpected: %s\nActual..: %s\n\n", input, expected, output)
	}
}

func assertDecodedByteSliceEqual(t *testing.T, input []byte, result []byte, expected string) {
	in := fmt.Sprintf("% X", input)
	out := fmt.Sprintf("% X", result)

	if out != expected {
		t.Errorf("Re-decoding the encoder output failed.\nInput...: %s\nExpected: %s\nActual..: %s\n\n", in, expected, out)
	}
}
