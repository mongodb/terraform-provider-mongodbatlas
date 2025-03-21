package conversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

func LastPart(p path.Path) string {
	parts := strings.Split(p.String(), ".")
	return parts[len(parts)-1]
}

func IsIndexValue(p path.Path) bool {
	return IsMapIndex(p) || IsListIndex(p) || IsSetIndex(p)
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

func IsSetIndex(p path.Path) bool {
	lastPart := LastPart(p)
	return strings.Contains(lastPart, "[Value(")
}

func HasPrefix(p, prefix path.Path) bool {
	prefixString := prefix.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func AttributeNameEquals(p path.Path, name string) bool {
	noBrackets := StripSquareBrackets(p)
	return noBrackets == name || strings.HasSuffix(noBrackets, fmt.Sprintf(".%s", name))
}

func AttributeName(p path.Path) string {
	noBrackets := StripSquareBrackets(p)
	parts := strings.Split(noBrackets, ".")
	return parts[len(parts)-1]
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
	if IsIndexValue(p) {
		return p.ParentPath().String()
	}
	return p.String()
}
