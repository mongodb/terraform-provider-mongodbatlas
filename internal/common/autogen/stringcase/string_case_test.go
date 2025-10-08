package stringcase_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Non capitalized to capitalized",
			input:    "toCaps",
			expected: "ToCaps",
		},
		{
			name:     "Already capitalized does nothing",
			input:    "ToCaps",
			expected: "ToCaps",
		},
		{
			name:     "Single char capitalizes",
			input:    "c",
			expected: "C",
		},
		{
			name:     "Non-capitalizable does nothing",
			input:    "_",
			expected: "_",
		},
		{
			name:     "Empty does nothing",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid utf8 does nothing",
			input:    string([]byte{0xFF}),
			expected: string([]byte{0xFF}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := stringcase.Capitalize(tt.input); actual != tt.expected {
				t.Errorf("Capitalize() returned %v, expected %v", actual, tt.expected)
			}
		})
	}
}

func TestUncapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Capitalized to non capitalized",
			input:    "FromCaps",
			expected: "fromCaps",
		},
		{
			name:     "Non capitalized does nothing",
			input:    "fromCaps",
			expected: "fromCaps",
		},
		{
			name:     "Single char uncapitalizes",
			input:    "C",
			expected: "c",
		},
		{
			name:     "Non-uncapitalizable does nothing",
			input:    "_",
			expected: "_",
		},
		{
			name:     "Empty does nothing",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid utf8 does nothing",
			input:    string([]byte{0xFF}),
			expected: string([]byte{0xFF}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := stringcase.Uncapitalize(tt.input); actual != tt.expected {
				t.Errorf("Uncapitalize() returned %v, expected %v", actual, tt.expected)
			}
		})
	}
}
