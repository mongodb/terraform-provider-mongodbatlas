package autogen

import "slices"

func ModifyJSONFromListToMap(attrObjJSON any) any {
	if attrObjJSON == nil {
		return nil
	}

	list, ok := attrObjJSON.([]any)
	if !ok {
		// If it's not a list, leave it unchanged.
		return attrObjJSON
	}

	// For an empty list, return an empty map.
	if len(list) == 0 {
		return map[string]any{}
	}

	result := make(map[string]any, len(list))
	for _, item := range list {
		obj, ok := item.(map[string]any)
		if !ok {
			// If any element is not an object, fall back to original value.
			return attrObjJSON
		}

		keyRaw, hasKey := obj["key"]
		value, hasValue := obj["value"]
		if !hasKey || !hasValue {
			// Skip items without both key and value.
			continue
		}

		key, ok := keyRaw.(string)
		if !ok {
			// Keys must be strings; skip otherwise.
			continue
		}

		result[key] = value
	}

	return result
}

func ModifyJSONFromMapToList(val any) any {
	if val == nil {
		return nil
	}

	obj, ok := val.(map[string]any)
	if !ok {
		// If it's not a map, leave it unchanged.
		return val
	}

	// For an empty map, return an empty list.
	if len(obj) == 0 {
		return []any{}
	}

	// To ensure deterministic output (useful for tests), sort keys.
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	result := make([]any, 0, len(obj))
	for _, k := range keys {
		result = append(result, map[string]any{
			"key":   k,
			"value": obj[k],
		})
	}

	return result
}
