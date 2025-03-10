package advancedclustertpf

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

func LastPart(p path.Path) string {
	parts := strings.Split(p.String(), ".")
	return parts[len(parts)-1]
}

func IsListIndex(p path.Path) bool {
	lastPart := LastPart(p)
	if IsMapIndex(p) {
		return false
	}
	return strings.HasSuffix(lastPart, "]")
}

func IsMapIndex(p path.Path) bool {
	lastPart := LastPart(p)
	return strings.HasSuffix(lastPart, "\"]")
}

func hasPrefix(p path.Path, prefix path.Path) bool {
	prefixString := prefix.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func AttributeNameEquals(p path.Path, name string) bool {
	noBrackets := StripSquareBrackets(p)
	return noBrackets == name || strings.HasSuffix(noBrackets, fmt.Sprintf(".%s", name))
}

func AsAddedIndex(p path.Path) string {
	parentString := p.ParentPath().ParentPath().String()
	lastPart := LastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[+", 1)
	if parentString == "" {
		return indexWithSign
	}
	return parentString + "." + indexWithSign
}

func AsRemovedIndex(p path.Path) string {
	parentString := p.ParentPath().ParentPath().String()
	lastPart := LastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[-", 1)
	if parentString == "" {
		return indexWithSign
	}
	return parentString + "." + indexWithSign
}

func StripSquareBrackets(p path.Path) string {
	if IsListIndex(p) {
		return p.ParentPath().String()
	}
	if IsMapIndex(p) {
		return p.ParentPath().String()
	}
	return p.String()
}
