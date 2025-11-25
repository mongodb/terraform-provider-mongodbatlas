package stringcase

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var unsupportedCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

// ToSnakeCase Multiple consecutive uppercase letters are treated as part of the same word except for the last one.
// Example: "MongoDBMajorVersion" -> "mongo_db_major_version"
func ToSnakeCase(str string) string {
	if str == "" {
		return str
	}

	str = unsupportedCharacters.ReplaceAllString(str, "")

	builder := &strings.Builder{}
	runes := []rune(str)
	length := len(runes)

	prevIsUpper := unicode.IsUpper(runes[0])
	builder.WriteRune(unicode.ToLower(runes[0]))

	for i := 1; i < length; i++ {
		current := runes[i]
		currentIsUpper := unicode.IsUpper(runes[i])

		// Write an underscore before uppercase letter if:
		// - Previous char was lowercase, so this is the first uppercase.
		// - Next char is lowercase, so this is the last uppercase in a sequence.
		if currentIsUpper {
			if !prevIsUpper || (i+1 != length && unicode.IsLower(runes[i+1])) {
				builder.WriteByte('_')
			}
			current = unicode.ToLower(current)
		}

		builder.WriteRune(current)
		prevIsUpper = currentIsUpper
	}

	return builder.String()
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
