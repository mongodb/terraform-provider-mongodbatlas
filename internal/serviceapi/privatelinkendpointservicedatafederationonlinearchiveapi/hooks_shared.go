package privatelinkendpointservicedatafederationonlinearchiveapi

func normalizeOptionalStringFields(obj map[string]any) {
	setEmptyStringIfMissing(obj, "comment")
	setEmptyStringIfMissing(obj, "region")
	setEmptyStringIfMissing(obj, "customerEndpointDNSName")
}

func setEmptyStringIfMissing(obj map[string]any, responseKey string) {
	if val, exists := obj[responseKey]; !exists || val == nil {
		obj[responseKey] = ""
	}
}
