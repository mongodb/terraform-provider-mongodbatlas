package utils

import "os"

const (
	AttrNotSetError     = "attribute %s must be set"
	ProviderConfigError = "error in configuring the provider."
)

func MultiEnvDefaultFunc(ks []string, def interface{}) interface{} {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return def
}
