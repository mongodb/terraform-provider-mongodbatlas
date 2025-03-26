package customplanmodifier

import (
	"fmt"
	"strings"
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

// ListIndexChanged returns true if the list at the given index has changed, false if it was added or removed
func (a AttributeChanges) ListIndexChanged(name string, index int) bool {
	attributeNameChanges := a.allAttributeNameChanges()
	indexPath := fmt.Sprintf("%s[%d]", name, index)
	_, found := attributeNameChanges[indexPath]
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
	attributeNameChanges := a.allAttributeNameChanges()
	addPrefix := fmt.Sprintf("%s[+", name)
	removePrefix := fmt.Sprintf("%s[-", name)
	for change := range attributeNameChanges {
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
