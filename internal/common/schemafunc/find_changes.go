package schemafunc

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AttributeChanges struct {
	Changes []string
}

func (a *AttributeChanges) LeafChanges() map[string]bool {
	return a.leafChanges(true)
}

func (a *AttributeChanges) AttributeChanged(name string) bool {
	changes := a.LeafChanges()
	changed := changes[name]
	return changed
}

func (a *AttributeChanges) KeepUnknown(attributeEffectedMapping map[string][]string) []string {
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
func (a *AttributeChanges) ListIndexChanged(name string, index int) bool {
	leafChanges := a.leafChanges(false)
	indexPath := fmt.Sprintf("%s[%d]", name, index)
	return leafChanges[indexPath]
}

// NestedListLenChanges accepts a fullPath, e.g., "replication_specs[0].region_configs" and returns true if the length of the nested list has changed
func (a *AttributeChanges) NestedListLenChanges(fullPath string) bool {
	addPrefix := fmt.Sprintf("%s[+", fullPath)
	removePrefix := fmt.Sprintf("%s[-", fullPath)
	for _, change := range a.Changes {
		if strings.HasPrefix(change, addPrefix) || strings.HasPrefix(change, removePrefix) {
			return true
		}
	}
	return false
}

func (a *AttributeChanges) ListLenChanges(name string) bool {
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

func (a *AttributeChanges) leafChanges(removeIndex bool) map[string]bool {
	leafChanges := map[string]bool{}
	for _, change := range a.Changes {
		var leaf string
		parts := strings.Split(change, ".")
		leaf = parts[len(parts)-1]
		if removeIndex && strings.HasSuffix(leaf, "]") {
			leaf = strings.Split(leaf, "[")[0]
		}
		leafChanges[leaf] = true
	}
	return leafChanges
}

// FindAttributeChanges: Iterates through TFModel of state+plan and returns AttributeChanges for querying changed attributes
// The implementation is similar to KeepUnknown, no support for types.Set or types.Tuple yet
func FindAttributeChanges(ctx context.Context, src, dest any) AttributeChanges {
	changes := FindChanges(ctx, src, dest)
	return AttributeChanges{changes}
}

func FindChanges(ctx context.Context, src, dest any) []string {
	valSrc, valDest := validateStructPointers(src, dest)
	typeDest := valDest.Type()
	changes := []string{} // Always return an empty list, as nested attributes might be added and then removed, which make the test cases fail on nil vs []
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name, tfName := fieldNameTFName(&fieldDest)
		nestedSrc := valSrc.FieldByName(name).Interface()
		nestedDest := valDest.FieldByName(name).Interface()
		compareSrc := nestedSrc.(attr.Value)
		compareDest := nestedDest.(attr.Value)
		if compareDest.IsNull() || compareDest.IsUnknown() || compareDest.Equal(compareSrc) {
			continue
		}
		changes = append(changes, tfName)
		objValueSrc, okSrc := nestedSrc.(types.Object)
		objValueDest, okDest := nestedDest.(types.Object)
		if okSrc && okDest {
			moreChanges := findChangesInObject(ctx, objValueSrc, objValueDest, []string{tfName})
			if len(moreChanges) == 0 {
				changes = slices.Delete(changes, len(changes)-1, len(changes))
			}
			changes = append(changes, moreChanges...)
			continue
		}
		listValueSrc, okSrc := nestedSrc.(types.List)
		listValueDest, okDest := nestedDest.(types.List)
		if okSrc && okDest {
			moreChanges := findChangesInList(ctx, listValueSrc, listValueDest, []string{tfName})
			if len(moreChanges) == 0 {
				changes = slices.Delete(changes, len(changes)-1, len(changes))
			}
			changes = append(changes, moreChanges...)
			continue
		}
	}
	return changes
}

func findChangesInObject(ctx context.Context, src, dest types.Object, parentPath []string) []string {
	var changes []string
	attributesSrc := src.Attributes()
	attributesDest := dest.Attributes()
	for name, attr := range attributesDest {
		path := slices.Clone(parentPath)
		path = append(path, name)
		if attr.IsNull() || attr.IsUnknown() || attr.Equal(attributesSrc[name]) {
			continue
		}
		changes = append(changes, strings.Join(path, "."))
		tfListDest, isList := attr.(types.List)
		tfObjectDest, isObject := attr.(types.Object)
		if isObject {
			moreChanges := findChangesInObject(ctx, attributesSrc[name].(types.Object), tfObjectDest, path)
			if len(moreChanges) == 0 {
				changes = slices.Delete(changes, len(changes)-1, len(changes))
			}
			changes = append(changes, moreChanges...)
		}
		if isList {
			moreChanges := findChangesInList(ctx, attributesSrc[name].(types.List), tfListDest, path)
			if len(moreChanges) == 0 {
				changes = slices.Delete(changes, len(changes)-1, len(changes))
			}
			changes = append(changes, moreChanges...)
		}
	}
	return changes
}

func findChangesInList(ctx context.Context, src, dest types.List, parentPath []string) []string {
	changes := []string{}
	srcElements := src.Elements()
	destElements := dest.Elements()
	if dest.IsNull() {
		return changes
	}
	maxCount := max(len(srcElements), len(destElements))
	for i := range maxCount {
		srcObj, srcOk := lookupIndex(srcElements, i)
		destObj, destOk := lookupIndex(destElements, i)
		path := slices.Clone(parentPath)
		indexPath := fmt.Sprintf("%s[%d]", strings.Join(path, "."), i)
		switch {
		case srcOk && destOk:
			indexChanges := findChangesInObject(ctx, srcObj, destObj, []string{})
			if len(indexChanges) == 0 {
				continue
			}
			changes = append(changes, indexPath)
			for _, change := range indexChanges {
				changes = append(changes, fmt.Sprintf("%s.%s", indexPath, change))
			}
		case srcOk && !destOk: // removed from list
			changes = append(changes, fmt.Sprintf("%s[-%d]", strings.Join(parentPath, "."), i))
		default: // added to list
			changes = append(changes, fmt.Sprintf("%s[+%d]", strings.Join(parentPath, "."), i))
		}
	}
	return changes
}
func lookupIndex(elements []attr.Value, index int) (types.Object, bool) {
	if index >= len(elements) {
		return types.ObjectNull(nil), false
	}
	obj, ok := elements[index].(types.Object)
	return obj, ok
}
