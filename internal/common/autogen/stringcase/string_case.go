package stringcase

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/huandu/xstrings"
)

var (
	camelCase             = regexp.MustCompile(`([a-z])[A-Z]`)
	unsupportedCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

type SnakeCaseString string

func (snake SnakeCaseString) SnakeCase() string {
	return string(snake)
}

func (snake SnakeCaseString) PascalCase() string {
	return xstrings.ToPascalCase(string(snake))
}

func (snake SnakeCaseString) CamelCase() string {
	return xstrings.ToCamelCase(string(snake))
}

func (snake SnakeCaseString) LowerCaseNoUnderscore() string {
	return strings.ReplaceAll(string(snake), "_", "")
}

func FromCamelCase(input string) SnakeCaseString {
	if input == "" {
		return SnakeCaseString(input)
	}

	removedUnsupported := unsupportedCharacters.ReplaceAllString(input, "")

	insertedUnderscores := camelCase.ReplaceAllStringFunc(removedUnsupported, func(s string) string {
		firstChar := s[0]
		restOfString := s[1:]
		return fmt.Sprintf("%c_%s", firstChar, strings.ToLower(restOfString))
	})

	return SnakeCaseString(strings.ToLower(insertedUnderscores))
}

func ToCamelCase(str string) string {
	return xstrings.ToCamelCase(str)
}

func ToSnakeCase(str string) string {
	return xstrings.ToSnakeCase(str)
}

func Capitalize(str string) string {
	return capitalization(str, true)
}

func Uncapitalize(str string) string {
	return capitalization(str, false)
}

func capitalization(str string, capitalize bool) string {
	if str == "" {
		return str
	}

	first, size := utf8.DecodeRuneInString(str)
	if first == utf8.RuneError {
		return str
	}

	builder := &strings.Builder{}
	if capitalize {
		builder.WriteRune(unicode.ToUpper(first))
	} else {
		builder.WriteRune(unicode.ToLower(first))
	}
	builder.WriteString(str[size:])
	return builder.String()
}
