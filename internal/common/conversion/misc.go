package conversion

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

func ValRegion(reg any, opt ...string) (string, error) {
	region, err := cast.ToStringE(reg)
	if err != nil {
		return "", err
	}

	if region == "" {
		return "", fmt.Errorf("region must be set")
	}

	/*
		We need to check if the option will be similar to network_peering word
		 (this comes in from the same resource) because network_pering resource
		 has not the standard region name pattern "US_EAST_1",
		 instead it needs the following one: "us-east-1".
	*/
	if len(opt) > 0 && strings.EqualFold("network_peering", opt[0]) {
		return strings.ToLower(strings.ReplaceAll(region, "_", "-")), nil
	}

	return strings.ReplaceAll(region, "-", "_"), nil
}

func IsEmpty(val any) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case *bool, *float64, *int64:
		return v == nil
	case string:
		return v == ""
	case *string:
		return v == nil || *v == ""
	}
	return false
}
