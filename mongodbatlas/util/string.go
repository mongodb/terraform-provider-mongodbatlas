package util

// IsStringPresent returns true if the string is non-empty.
func IsStringPresent(strPtr *string) bool {
	return strPtr != nil && len(*strPtr) > 0
}
