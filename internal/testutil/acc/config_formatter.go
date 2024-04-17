package acc

import (
	"fmt"
	"sort"
	"strings"
)

func FormatToHCLMap(m map[string]string, indent, varName string) string {
	if m == nil {
		return ""
	}
	lines := []string{
		fmt.Sprintf("%s%s = {", indent, varName),
	}
	indentKeyValues := indent + "\t"

	for _, k := range sortStringMapKeys(m) {
		v := m[k]
		lines = append(lines, fmt.Sprintf("%s%s = %[3]q", indentKeyValues, k, v))
	}
	lines = append(lines, fmt.Sprintf("%s}", indent))
	return strings.Join(lines, "\n")
}

func FormatToHCLLifecycleIgnore(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	ignoreParts := []string{}
	for _, ignoreKey := range keys {
		ignoreParts = append(ignoreParts, fmt.Sprintf("\t\t\t%s,", ignoreKey))
	}
	lines := []string{
		"\tlifecycle {",
		"\t\tignore_changes = [",
		strings.Join(ignoreParts, "\n"),
		"\t\t]",
		"\t}",
	}
	return strings.Join(lines, "\n")
}

// make test deterministic
func sortStringMapKeys(m map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
