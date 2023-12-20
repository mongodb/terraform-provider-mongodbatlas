package main

import (
	"strings"
	"unicode"
)

// toPascalCase converts camel case to pascal case.
func ToPascalCase(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}

// toSnakeCase converts camel case to snake case.
func ToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
