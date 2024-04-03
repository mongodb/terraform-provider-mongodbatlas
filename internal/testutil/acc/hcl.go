package acc

import (
	"fmt"
	"strings"
)

func MapToHcl(m map[string]string, indent, varName string) string {
	if m == nil {
		return ""
	}
	lines := []string{
		fmt.Sprintf("%s = {", varName),
	}
	indentKeyValues := indent + "\t"
	for k, v := range m {
		lines = append(lines, fmt.Sprintf("%s%s = %[3]q", indentKeyValues, k, v))
	}
	lines = append(lines, fmt.Sprintf("%s}", indent))
	return strings.Join(lines, "\n")
}
