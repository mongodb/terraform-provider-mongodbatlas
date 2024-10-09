package codespec

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	camelCase             = regexp.MustCompile(`([a-z])[A-Z]`)
	unsupportedCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

func terraformAttrName(attrName string) string {
	if attrName == "" {
		return attrName
	}

	removedUnsupported := unsupportedCharacters.ReplaceAllString(attrName, "")

	insertedUnderscores := camelCase.ReplaceAllStringFunc(removedUnsupported, func(s string) string {
		firstRune, size := utf8.DecodeRuneInString(s)
		if firstRune == utf8.RuneError && size <= 1 {
			return s
		}

		return fmt.Sprintf("%s_%s", string(firstRune), strings.ToLower(s[size:]))
	})
	return strings.ToLower(insertedUnderscores)
}
