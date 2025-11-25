package schemafunc

import "github.com/hashicorp/terraform-plugin-go/tftypes"

func GetAttrFromStateObj[T any](rawState map[string]tftypes.Value, attrName string) *T {
	var ret *T
	if err := rawState[attrName].As(&ret); err != nil {
		return nil
	}
	return ret
}
