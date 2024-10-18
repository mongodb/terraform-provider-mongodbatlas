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
	return xstrings.ToCamelCase(string(snake)) // in xstrings v1.15.0 we can switch to using ToPascalCase for same functionality
}

func (snake SnakeCaseString) LowerCaseNoUnderscore() string {
	return strings.ReplaceAll(string(snake), "_", "")
}
