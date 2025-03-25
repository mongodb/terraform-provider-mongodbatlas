package customplanmodifier

import (
	"fmt"
	"strings"
)

type AttributeChanges []string

func (a AttributeChanges) LeafChanges() map[string]struct{} {
	return a.leafChanges(true)
}

func (a AttributeChanges) AttributeChanged(name string) bool {
	changes := a.LeafChanges()
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

// ListIndexChanged returns true if the list at the given index has changed, false if it was added or removed
func (a AttributeChanges) ListIndexChanged(name string, index int) bool {
	leafChanges := a.leafChanges(false)
	indexPath := fmt.Sprintf("%s[%d]", name, index)
	_, found := leafChanges[indexPath]
	return found
}

// NestedListLenChanges accepts a fullPath, e.g., "replication_specs[0].region_configs" and returns true if the length of the nested list has changed
func (a AttributeChanges) NestedListLenChanges(fullPath string) bool {
	addPrefix := fmt.Sprintf("%s[+", fullPath)
	removePrefix := fmt.Sprintf("%s[-", fullPath)
	for _, change := range a {
		if strings.HasPrefix(change, addPrefix) || strings.HasPrefix(change, removePrefix) {
			return true
		}
	}
	return false
}

func (a AttributeChanges) ListLenChanges(name string) bool {
	leafChanges := a.leafChanges(false)
	addPrefix := fmt.Sprintf("%s[+", name)
	removePrefix := fmt.Sprintf("%s[-", name)
	for change := range leafChanges {
		if strings.HasPrefix(change, addPrefix) || strings.HasPrefix(change, removePrefix) {
			return true
		}
	}
	return false
}

func (a AttributeChanges) leafChanges(removeIndex bool) map[string]struct{} {
	leafChanges := make(map[string]struct{})
	for _, change := range a {
		var leaf string
		parts := strings.Split(change, ".")
		leaf = parts[len(parts)-1]
		if removeIndex && strings.HasSuffix(leaf, "]") {
			leaf = strings.Split(leaf, "[")[0]
		}
		leafChanges[leaf] = struct{}{}
	}
	return leafChanges
}
