package schemafunc

import (
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
)

// Returns true if both slices contain the same elements regardless of order.
func StringSlicesEqualIgnoringOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sortedA := make([]string, len(a))
	copy(sortedA, a)
	sort.Strings(sortedA)
	sortedB := make([]string, len(b))
	copy(sortedB, b)
	sort.Strings(sortedB)
	return slices.Equal(sortedA, sortedB)
}

// DiffSuppressFunc that suppresses diffs on a TypeList of strings when the only difference is element order.
func EqualStringListsIgnoringOrder(k, oldValue, newValue string, d *schema.ResourceData) bool {
	baseKey := k[:strings.LastIndex(k, ".")]
	oldVal, newVal := d.GetChange(baseKey)
	return StringSlicesEqualIgnoringOrder(cast.ToStringSlice(oldVal), cast.ToStringSlice(newVal))
}
