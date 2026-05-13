// Package responseproject projects an API response map down to a
// practitioner-declared allow-list of dotted paths.
//
// The generic api_resource / api_update resources store the response in two
// computed attributes — `output` (plain) and `output_sensitive` (Sensitive).
// Each attribute is populated by Project against its own path list. Anything
// not listed is dropped before state is written, which keeps the state file
// free of API metadata the practitioner didn't ask for and prevents secrets
// from leaking into plan/apply stdout via the non-sensitive `output`.
//
// Path syntax (v1):
//   - Dotted keys: "apiKeyId", "data.profile.email"
//   - Numeric list indices: "secrets.0.value", "headers.1.name"
//
// Wildcards (e.g. "headers.*.value") are intentionally not supported yet —
// when added they will be additive.
package responseproject

import (
	"strconv"
	"strings"
)

// Project returns a new map containing only the values reachable via the
// supplied paths. Paths whose targets are absent are silently skipped (the
// API legitimately stops returning create-only fields like the AI Model API
// Key `secret` after the first read, and that case must not be an error).
//
// If paths is empty the result is an empty map — the caller decides whether
// to translate that into a null Dynamic attribute.
func Project(response map[string]any, paths []string) map[string]any {
	out := map[string]any{}
	if response == nil || len(paths) == 0 {
		return out
	}
	for _, p := range paths {
		segs := splitPath(p)
		if len(segs) == 0 {
			continue
		}
		val, ok := lookup(response, segs)
		if !ok {
			continue
		}
		insert(out, segs, val)
	}
	return out
}

// PathsOverlap reports whether any path appears in both a and b. Used by
// ValidateConfig to reject configurations that put the same response path
// in both response_export_values and response_export_values_sensitive.
func PathsOverlap(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(a))
	for _, p := range a {
		set[p] = struct{}{}
	}
	var overlap []string
	for _, p := range b {
		if _, ok := set[p]; ok {
			overlap = append(overlap, p)
		}
	}
	return overlap
}

func splitPath(p string) []string {
	p = strings.TrimSpace(p)
	if p == "" {
		return nil
	}
	return strings.Split(p, ".")
}

// lookup walks segs through tree. Map segments traverse by key; numeric
// segments traverse []any by index. Returns (value, true) on hit.
func lookup(tree any, segs []string) (any, bool) {
	cur := tree
	for _, seg := range segs {
		switch node := cur.(type) {
		case map[string]any:
			v, ok := node[seg]
			if !ok {
				return nil, false
			}
			cur = v
		case []any:
			i, err := strconv.Atoi(seg)
			if err != nil || i < 0 || i >= len(node) {
				return nil, false
			}
			cur = node[i]
		default:
			return nil, false
		}
	}
	return cur, true
}

// insert places val at the location identified by segs inside dst, creating
// intermediate containers as needed. List indices grow the slice with nil
// padding so non-contiguous projections (e.g. only "items.2.name") still
// produce a valid tree.
func insert(dst map[string]any, segs []string, val any) {
	if len(segs) == 0 {
		return
	}
	if len(segs) == 1 {
		dst[segs[0]] = val
		return
	}
	key := segs[0]
	rest := segs[1:]
	// Decide whether the next container should be a list or a map based on
	// the next segment: if it's a non-negative integer, treat as list.
	if idx, isIdx := parseIndex(rest[0]); isIdx {
		var list []any
		if existing, ok := dst[key].([]any); ok {
			list = existing
		}
		list = growList(list, idx+1)
		insertList(list, rest, val)
		dst[key] = list
		return
	}
	var sub map[string]any
	if existing, ok := dst[key].(map[string]any); ok {
		sub = existing
	} else {
		sub = map[string]any{}
	}
	insert(sub, rest, val)
	dst[key] = sub
}

func insertList(list []any, segs []string, val any) {
	idx, _ := parseIndex(segs[0])
	if len(segs) == 1 {
		list[idx] = val
		return
	}
	rest := segs[1:]
	if nextIdx, isIdx := parseIndex(rest[0]); isIdx {
		var nested []any
		if existing, ok := list[idx].([]any); ok {
			nested = existing
		}
		nested = growList(nested, nextIdx+1)
		insertList(nested, rest, val)
		list[idx] = nested
		return
	}
	var sub map[string]any
	if existing, ok := list[idx].(map[string]any); ok {
		sub = existing
	} else {
		sub = map[string]any{}
	}
	insert(sub, rest, val)
	list[idx] = sub
}

func parseIndex(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil || i < 0 {
		return 0, false
	}
	return i, true
}

func growList(list []any, n int) []any {
	for len(list) < n {
		list = append(list, nil)
	}
	return list
}
