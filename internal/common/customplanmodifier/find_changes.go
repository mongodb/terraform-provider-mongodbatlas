package customplanmodifier

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

type AttributeChanges []string

func (a AttributeChanges) AttributeChanged(name string) bool {
	changes := a.allAttributeNameChanges()
	_, found := changes[name]
	return found
}

func (a AttributeChanges) KeepUnknown(attributeEffectedMapping map[string][]string) []string {
	var keepUnknown []string
	for attrChanged, affectedAttributes := range attributeEffectedMapping {
		if a.AttributeChanged(attrChanged) {
			keepUnknown = append(keepUnknown, attrChanged)
			keepUnknown = append(keepUnknown, affectedAttributes...)
		}
	}
	return keepUnknown
}

func (a AttributeChanges) PathChanged(path path.Path) bool {
	return slices.Contains(a, path.String())
}

// ListLenChanged accepts a fullPath, e.g., "replication_specs[0].region_configs" and returns true if the length of the nested list has changed
func (a AttributeChanges) ListLenChanged(p path.Path) bool {
	addPrefix := asAddPrefix(p.String())
	removePrefix := asRemovePrefix(p.String())
	for _, change := range a {
		if strings.HasPrefix(change, addPrefix) || strings.HasPrefix(change, removePrefix) {
			return true
		}
	}
	return false
}

func (a AttributeChanges) allAttributeNameChanges() map[string]struct{} {
	nameChanges := make(map[string]struct{})
	for _, change := range a {
		parts := strings.Split(change, ".")
		attributeName := parts[len(parts)-1]
		nameChanges[attributeName] = struct{}{}
	}
	return nameChanges
}

// asAddPrefix must match conversion.AsAddedIndex
func asAddPrefix(p string) string {
	return fmt.Sprintf("%s[+", p)
}

// asRemovePrefix must match conversion.AsRemovedIndex
func asRemovePrefix(p string) string {
	return fmt.Sprintf("%s[-", p)
}
