package codespec

import (
	"strings"

	"github.com/huandu/xstrings"
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
