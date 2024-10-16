package codespec

import (
	"regexp"
	"strings"
)

type SnakeCaseString string

func (snake SnakeCaseString) SnakeCase() string {
	return string(snake)
}

func (snake SnakeCaseString) PascalCase() string {
	words := strings.Split(string(snake), "_")
	var pascalCase string
	for i := range words {
		if words[i] != "" {
			pascalCase += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return pascalCase
}

func (snake SnakeCaseString) LowerCaseNoUnderscore() string {
	return strings.ReplaceAll(string(snake), "_", "")
}

func ValidateSnakeCase(s string) bool {
	snakeCaseRegex := `^[a-z]+(_[a-z]+)*$`
	matched, err := regexp.MatchString(snakeCaseRegex, s)
	if err != nil {
		return false
	}
	return matched
}
