package conversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func IsListIndex(p path.Path) bool {
	lastPart := lastPart(p)
	if IsMapIndex(p) || IsSetIndex(p) {
		return false
	}
	return strings.HasSuffix(lastPart, "]")
}

func IsMapIndex(p path.Path) bool {
	lastPart := lastPart(p)
	return strings.HasSuffix(lastPart, "\"]")
}

func IsSetIndex(p path.Path) bool {
	lastPart := lastPart(p)
	return strings.Contains(lastPart, "[Value(")
}

func HasAncestor(p, ancestor path.Path) bool {
	prefixString := ancestor.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func AttributeName(p path.Path) string {
	noIndex := trimLastIndex(p)
	parts := strings.Split(noIndex, ".")
	return parts[len(parts)-1]
}

// AsAddedIndex returns "" if the path is not an index otherwise it adds `+` before the index
func AsAddedIndex(p path.Path) string {
	if !isIndexValue(p) {
		return ""
	}
	lastPart := lastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[+", 1)
	everythingExceptLast, _ := strings.CutSuffix(p.String(), lastPart)
	return everythingExceptLast + indexWithSign
}

// AsRemovedIndex returns "" if the path is not an index otherwise it adds `-` before the index
func AsRemovedIndex(p path.Path) string {
	if !isIndexValue(p) {
		return ""
	}
	lastPart := lastPart(p)
	lastPartWithRemoveIndex := strings.Replace(lastPart, "[", "[-", 1)
	everythingExceptLast, _ := strings.CutSuffix(p.String(), lastPart)
	return everythingExceptLast + lastPartWithRemoveIndex
}

func AncestorPathWithIndex(p path.Path, attributeName string, diags *diag.Diagnostics) path.Path {
	for {
		p = p.ParentPath()
		if p.Equal(path.Empty()) {
			diags.AddError("Parent path not found", fmt.Sprintf("Parent attribute %s not found in path %s", attributeName, p.String()))
			return p
		}
		if attributeNameEquals(p, attributeName) {
			return p
		}
	}
}

func AncestorPathNoIndex(p path.Path, attributeName string, diags *diag.Diagnostics) path.Path {
	parent := AncestorPathWithIndex(p, attributeName, diags)
	if diags.HasError() {
		return parent
	}
	return trimLastIndexPath(parent)
}

func AncestorPaths(p path.Path) []path.Path {
	ancestors := []path.Path{}
	for {
		ancestor := p.ParentPath()
		if ancestor.Equal(path.Empty()) {
			return ancestors
		}
		ancestors = append(ancestors, ancestor)
		p = ancestor
	}
}

func lastPart(p path.Path) string {
	parts := strings.Split(p.String(), ".")
	return parts[len(parts)-1]
}

func isIndexValue(p path.Path) bool {
	return IsMapIndex(p) || IsListIndex(p) || IsSetIndex(p)
}

func attributeNameEquals(p path.Path, name string) bool {
	return AttributeName(p) == name
}

func trimLastIndex(p path.Path) string {
	return trimLastIndexPath(p).String()
}

func trimLastIndexPath(p path.Path) path.Path {
	if isIndexValue(p) {
		return p.ParentPath()
	}
	return p
}
