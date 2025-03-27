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
	if !isIndexValue(p) {
		return ""
	}
	lastPart := lastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[+", 1)
	everythingExceptLast, _ := strings.CutSuffix(p.String(), lastPart)
	return everythingExceptLast + indexWithSign
}

// AsRemovedIndex returns empty string if the path is not an index otherwise it adds `-` before the index
func AsRemovedIndex(p path.Path) string {
	if !isIndexValue(p) {
		return ""
	}
	lastPart := lastPart(p)
	lastPartWithRemoveIndex := strings.Replace(lastPart, "[", "[-", 1)
	everythingExceptLast, _ := strings.CutSuffix(p.String(), lastPart)
	return everythingExceptLast + lastPartWithRemoveIndex
}

func TrimLastIndex(p path.Path) string {
	if isIndexValue(p) {
		return p.ParentPath().String()
	}
	return p.String()
}

func TrimLastIndexPath(p path.Path) path.Path {
	for {
		if isIndexValue(p) {
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
