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
	leafChanges := map[string]bool{}
	for _, change := range a.Changes {
		var leaf string
		parts := strings.Split(change, ".")
		if len(parts) == 1 {
			leaf = parts[0]
		} else {
			leaf = parts[len(parts)-1]
		}
		if strings.HasSuffix(leaf, "]") {
			leaf = strings.Split(leaf, "[")[0]
		}
		leafChanges[leaf] = true
	}
	return leafChanges
}

func (a *AttributeChanges) AttributeChanged(name string) bool {
	changes := a.LeafChanges()
	changed := changes[name]
	return changed
}

func (a *AttributeChanges) KeepUnknown(attributeEffectedMapping map[string][]string) []string {
	keepUnknown := []string{}
	for attrChanged, affectedAttributes := range attributeEffectedMapping {
		if a.AttributeChanged(attrChanged) {
			keepUnknown = append(keepUnknown, attrChanged)
			keepUnknown = append(keepUnknown, affectedAttributes...)
		}
	}
	return keepUnknown
}

func FindAttributeChanges(ctx context.Context, src, dest any) AttributeChanges {
	changes := FindChanges(ctx, src, dest)
	return AttributeChanges{changes}
}

// FindChanges TODO: Add description
func FindChanges(ctx context.Context, src, dest any) []string {
	valSrc, valDest := validateStructPointers(src, dest)
	typeDest := valDest.Type()
	changes := []string{} // Always return an empty list, as nested attributes might be added and then removed
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
	changes := []string{}
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
			changes = append(changes, indexPath)
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
