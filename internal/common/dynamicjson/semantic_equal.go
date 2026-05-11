package dynamicjson

import "github.com/hashicorp/terraform-plugin-framework/types"

// SemanticallyEqual reports whether two Dynamic values represent the same
// JSON payload after canonicalisation (sorted keys, normalised numbers).
// Null/Unknown both serialize to JSON null, so they compare equal — callers
// that need to distinguish those must guard before calling.
func SemanticallyEqual(a, b types.Dynamic) bool {
	aj, err := ToJSON(a)
	if err != nil {
		return false
	}
	bj, err := ToJSON(b)
	if err != nil {
		return false
	}
	if len(aj) != len(bj) {
		return false
	}
	for i := range aj {
		if aj[i] != bj[i] {
			return false
		}
	}
	return true
}
