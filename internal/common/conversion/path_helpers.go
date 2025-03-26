package conversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	noBrackets := TrimLastIndex(p)
	return noBrackets == name || strings.HasSuffix(noBrackets, fmt.Sprintf(".%s", name))
}

func AttributeName(p path.Path) string {
	noBrackets := TrimLastIndex(p)
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

func TrimLastIndex(p path.Path) string {
	if IsIndexValue(p) {
		return p.ParentPath().String()
	}
	return p.String()
}

func TrimLastIndexPath(p path.Path) path.Path {
	for {
		if IsIndexValue(p) {
			p = p.ParentPath()
		} else {
			return p
		}
	}
}

func ParentPathWithIndex(p path.Path, attributeName string, diags *diag.Diagnostics) path.Path {
	for {
		p = p.ParentPath()
		if p.Equal(path.Empty()) {
			diags.AddError("Parent path not found", fmt.Sprintf("Parent attribute %s not found in path %s", attributeName, p.String()))
			return p
		}
		if AttributeNameEquals(p, attributeName) {
			return p
		}
	}
}

func ParentPathNoIndex(p path.Path, attributeName string, diags *diag.Diagnostics) path.Path {
	parent := ParentPathWithIndex(p, attributeName, diags)
	if diags.HasError() {
		return parent
	}
	return TrimLastIndexPath(parent)
}

func HasPathParent(p path.Path, parentAttributeName string) bool {
	for {
		p = p.ParentPath()
		if p.Equal(path.Empty()) {
			return false
		}
		if AttributeNameEquals(p, parentAttributeName) {
			return true
		}
	}
}
