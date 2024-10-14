package codespec

import (
	"fmt"
	"regexp"
	"strings"
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
		firstChar := s[0]
		restOfString := s[1:]
		return fmt.Sprintf("%c_%s", firstChar, strings.ToLower(restOfString))
	})

	return strings.ToLower(insertedUnderscores)
}
