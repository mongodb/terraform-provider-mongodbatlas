// Package dynamicreshape walks a configured body alongside an API response
// and produces a state body that mirrors the configured shape while taking
// values from the response.
//
// Keys configured but absent from the response are preserved at their
// configured value (handles write-only request fields the API does not echo
// back — overwriting with null would surface as phantom drift on every
// refresh). Keys present only in the response are dropped (noise suppressed).
// Paths listed in SensitivePaths are excluded entirely so write-only secrets
// never leak into state.
package dynamicreshape

import (
	"reflect"
	"strings"
)

// Options controls reshape behavior.
type Options struct {
	// SensitivePaths is the set of dotted paths to exclude from the reshape
	// output (matches keys derived from sensitive_body). A path "a.b" excludes
	// the b leaf of object a. Lists are not indexed in this MVP.
	SensitivePaths map[string]struct{}
	// ListIDKeys maps a dotted list path to the field name identifying items
	// (e.g. "headers" -> "name"). When set, the engine matches items by ID
	// instead of by position.
	ListIDKeys map[string]string
}

// Reshape returns a new value shaped like configured, populated from response.
// Both inputs are plain Go trees (map[string]any / []any / scalars / nil).
func Reshape(configured, response any, opts Options) any {
	return reshapeAt(configured, response, opts, "")
}

func reshapeAt(configured, response any, opts Options, path string) any {
	if _, sensitive := opts.SensitivePaths[path]; sensitive {
		return nil
	}
	switch c := configured.(type) {
	case map[string]any:
		return reshapeObject(c, response, opts, path)
	case []any:
		return reshapeList(c, response, opts, path)
	}
	// Scalar: just pass through the response value. nil means "absent in response".
	return response
}

func reshapeObject(configured map[string]any, response any, opts Options, path string) any {
	resp, _ := response.(map[string]any)
	out := make(map[string]any, len(configured))
	for k, cv := range configured {
		child := joinPath(path, k)
		if _, sensitive := opts.SensitivePaths[child]; sensitive {
			continue
		}
		var rv any
		var present bool
		if resp != nil {
			rv, present = resp[k]
		}
		// Keys configured but absent from the response are typically write-only
		// request fields the API does not echo back. Preserve the configured
		// value instead of overwriting with null — null would create a phantom
		// "drift to null" diff on every refresh.
		if !present {
			out[k] = cv
			continue
		}
		out[k] = reshapeAt(cv, rv, opts, child)
	}
	return out
}

func reshapeList(configured []any, response any, opts Options, path string) any {
	resp, _ := response.([]any)
	if resp == nil {
		// Whole list missing in response → null per slot.
		out := make([]any, len(configured))
		for i, cv := range configured {
			out[i] = reshapeAt(cv, nil, opts, path)
		}
		return out
	}
	if idKey, ok := opts.ListIDKeys[path]; ok {
		return reshapeListByID(configured, resp, opts, path, idKey)
	}
	// Scalar arrays containing the same elements as configured are treated as
	// order-insensitive: preserve the configured order so server-side reorders
	// of set-like fields (roles, labels, etc.) do not cause spurious drift.
	if isScalarArray(configured) && isScalarArray(resp) && sameElementMultiset(configured, resp) {
		out := make([]any, len(configured))
		copy(out, configured)
		return out
	}
	n := min(len(resp), len(configured))
	out := make([]any, len(configured))
	for i, cv := range configured {
		var rv any
		if i < n {
			rv = resp[i]
		}
		out[i] = reshapeAt(cv, rv, opts, path)
	}
	return out
}

func isScalarArray(arr []any) bool {
	for _, e := range arr {
		switch e.(type) {
		case map[string]any, []any:
			return false
		}
	}
	return true
}

func sameElementMultiset(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[any]int, len(a))
	for _, v := range a {
		counts[normalizeID(v)]++
	}
	for _, v := range b {
		k := normalizeID(v)
		c, ok := counts[k]
		if !ok || c == 0 {
			return false
		}
		counts[k] = c - 1
	}
	return true
}

func reshapeListByID(configured, response []any, opts Options, path, idKey string) any {
	respByID := make(map[any]map[string]any, len(response))
	for _, item := range response {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if id, present := m[idKey]; present {
			respByID[normalizeID(id)] = m
		}
	}
	out := make([]any, len(configured))
	for i, cv := range configured {
		cm, ok := cv.(map[string]any)
		if !ok {
			out[i] = reshapeAt(cv, nil, opts, path)
			continue
		}
		var match any
		if id, present := cm[idKey]; present {
			match = respByID[normalizeID(id)]
		}
		out[i] = reshapeAt(cm, match, opts, path)
	}
	return out
}

func normalizeID(v any) any {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.String {
		return rv.String()
	}
	// json.Number / numerics: compare by string representation.
	type stringer interface{ String() string }
	if s, ok := v.(stringer); ok {
		return s.String()
	}
	return v
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// CollectSensitivePaths walks a sensitive_body tree and returns the set of
// dotted paths that should be excluded from reshape output.
func CollectSensitivePaths(body any) map[string]struct{} {
	out := make(map[string]struct{})
	walkSensitive(body, "", out)
	return out
}

func walkSensitive(v any, path string, out map[string]struct{}) {
	if path != "" {
		out[path] = struct{}{}
	}
	if m, ok := v.(map[string]any); ok {
		for k, child := range m {
			walkSensitive(child, joinPath(path, k), out)
		}
	}
}

// HasPrefix reports whether p starts with prefix as a path segment.
// Exported only for tests.
func HasPrefix(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, prefix+".")
}
