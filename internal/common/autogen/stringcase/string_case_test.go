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
				t.Errorf("Capitalize(%q) returned %q, expected %q", tt.input, actual, tt.expected)
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
				t.Errorf("Uncapitalize(%q) returned %q, expected %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single letter",
			input:    "a",
			expected: "a",
		},
		{
			name:     "Single uppercase letter",
			input:    "A",
			expected: "a",
		},
		{
			name:     "Simple camelCase",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "Simple PascalCase",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "All lowercase",
			input:    "word",
			expected: "word",
		},
		{
			name:     "All uppercase",
			input:    "WORD",
			expected: "word",
		},
		{
			name:     "Consecutive uppercase at start, middle and end",
			input:    "THISIsANExampleWORD",
			expected: "this_is_an_example_word",
		},
		{
			name:     "Numbers do not split words",
			input:    "Example123Word456WithNUMBERS789",
			expected: "example123_word456_with_numbers789",
		},
		{
			name:     "Already snake_case",
			input:    "already_snake_case",
			expected: "already_snake_case",
		},
		{
			name:     "Unsupported characters are removed",
			input:    "Example#!Unsup-.ported%&Chars",
			expected: "example_unsupported_chars",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := stringcase.ToSnakeCase(tt.input); actual != tt.expected {
				t.Errorf("ToSnakeCase(%q) returned %q, expected %q", tt.input, actual, tt.expected)
			}
		})
	}
}
