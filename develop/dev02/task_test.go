package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Unpack(test.input)

			if (err != nil) != test.wantErr {
				t.Errorf("Unexpected error status: got %v, want %v", err, test.wantErr)
				return
			}

			if result != test.expected {
				t.Errorf("Unpack(%s) = %s, want %s", test.input, result, test.expected)
			}
		})
	}
}
