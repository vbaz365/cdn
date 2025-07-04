package utils

import "testing"

func TestCommonPrefixLen(t *testing.T) {
	tests := []struct {
		name     string
		node     uint64
		inserted uint64
		expected uint8
	}{
		{
			name:     "Common prefix 8 bits",
			node:     0b1111000000000000000000000000000000000000000000000000000000000000,
			inserted: 0b1111000011110000000000000000000000000000000000000000000000000000,
			expected: 8,
		},
		{
			name:     "Identical inputs full 64",
			node:     0xFFFFFFFFFFFFFFFF,
			inserted: 0xFFFFFFFFFFFFFFFF,
			expected: 64,
		},
		{
			name:     "No common prefix",
			node:     0,
			inserted: 0x8000000000000000,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommonPrefixLen(tt.node, tt.inserted)
			if got != tt.expected {
				t.Errorf("CommonPrefixLen(%064b, %064b) = %d; expected %d", tt.node, tt.inserted, got, tt.expected)
			}
		})
	}
}

func TestGetFirstBit(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected uint8
	}{
		{name: "MSB is 1", input: 0x8000000000000000, expected: 1},
		{name: "MSB is 0", input: 0x7FFFFFFFFFFFFFFF, expected: 0},
		{name: "All zeros", input: 0, expected: 0},
		{name: "All ones", input: 0xFFFFFFFFFFFFFFFF, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFirstBit(tt.input)
			if got != tt.expected {
				t.Errorf("GetFirstBit(%064b) = %d; expected %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetBitsAfterIndex(t *testing.T) {
	tests := []struct {
		name     string
		high     uint64
		low      uint64
		index    uint8
		expected uint64
	}{
		{
			name:     "index less than 64",
			high:     0xFFFF0000FFFF0000,
			low:      0x0000000FFFFFFFFF,
			index:    32,
			expected: 0xFFFF00000000000F,
		},
		{
			name:     "index equal to 64",
			high:     0x123456789ABCDEF0,
			low:      0x0FEDCBA987654321,
			index:    64,
			expected: 0x0FEDCBA987654321,
		},
		{
			name:     "index zero",
			high:     0x123456789ABCDEF0,
			low:      0x0FEDCBA987654321,
			index:    0,
			expected: 0x123456789ABCDEF0,
		},
		{
			name:     "index greater than 64",
			high:     0x123456789ABCDEF0,
			low:      0b1111000011110000111100001111000011110000111100001111000011110000,
			index:    65,
			expected: 0b1110000111100001111000011110000111100001111000011110000111100000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBitsAfterIndex(tt.high, tt.low, tt.index)
			if got != tt.expected {
				t.Errorf("GetBitsAfterIndex(%064b, %064b, %d) = %064b; expected %064b",
					tt.high, tt.low, tt.index, got, tt.expected)
			}
		})
	}
}

func TestCombineSmallCompressions(t *testing.T) {
	tests := []struct {
		name           string
		first          uint64
		firstLen       uint8
		second         uint64
		secondLen      uint8
		expectedVal    uint64
		expectedLength uint8
	}{
		{
			name:           "Combine simple 4 and 4 bits",
			first:          0b1111,
			firstLen:       4,
			second:         0b1100,
			secondLen:      4,
			expectedVal:    0b1111 | (0b1100 >> 4),
			expectedLength: 8,
		},
		{
			name:           "Combine zero lengths",
			first:          0,
			firstLen:       0,
			second:         0,
			secondLen:      0,
			expectedVal:    0,
			expectedLength: 0,
		},
		{
			name:           "Combine with zero first",
			first:          0,
			firstLen:       0,
			second:         0b101010,
			secondLen:      6,
			expectedVal:    0b101010,
			expectedLength: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotLen := CombineSmallCompressions(tt.first, tt.firstLen, tt.second, tt.secondLen)
			if gotVal != tt.expectedVal || gotLen != tt.expectedLength {
				t.Errorf("CombineSmallCompressions(%b, %d, %b, %d) = (%b, %d); expected (%b, %d)",
					tt.first, tt.firstLen, tt.second, tt.secondLen, gotVal, gotLen, tt.expectedVal, tt.expectedLength)
			}
		})
	}
}

func TestSplitCompression(t *testing.T) {
	tests := []struct {
		name         string
		input        uint64
		splitIndex   uint8
		expectedHigh uint64
		expectedLow  uint64
	}{
		{
			name:         "Split at 8",
			input:        0b1111000011110000111100001111000011110000111100001111000011110000,
			splitIndex:   8,
			expectedHigh: 0b1111000000000000000000000000000000000000000000000000000000000000,
			expectedLow:  0b1111000011110000111100001111000011110000111100001111000000000000,
		},
		{
			name:         "Split at 0",
			input:        0xFFFFFFFFFFFFFFFF,
			splitIndex:   0,
			expectedHigh: 0,
			expectedLow:  0xFFFFFFFFFFFFFFFF,
		},
		{
			name:         "Split at 64",
			input:        0xFFFFFFFFFFFFFFFF,
			splitIndex:   64,
			expectedHigh: 0xFFFFFFFFFFFFFFFF,
			expectedLow:  0,
		},
		{
			name:         "Split at 32",
			input:        0xFFFF0000FFFF0000,
			splitIndex:   32,
			expectedHigh: 0xFFFF000000000000,
			expectedLow:  0xFFFF000000000000,
		},
		{
			name:         "Split at 12",
			input:        0x0FF0F00FFF000000,
			splitIndex:   12,
			expectedHigh: 0x0FF0000000000000,
			expectedLow:  0x0F00FFF000000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right := SplitCompression(tt.input, tt.splitIndex)
			if left != tt.expectedHigh {
				t.Errorf("Left mismatch: got %064b, expected %064b", left, tt.expectedHigh)
			}
			if right != tt.expectedLow {
				t.Errorf("Right mismatch: got %064b, expected %064b", right, tt.expectedLow)
			}
		})
	}
}
