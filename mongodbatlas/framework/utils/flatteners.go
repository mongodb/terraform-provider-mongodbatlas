package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ArrToSetValue(in []attr.Value) basetypes.SetValue {
	if len(in) == 0 {
		return types.SetNull(types.StringType)
	}
	return types.SetValueMust(types.StringType, in)
}

func ArrToListValue(in []attr.Value) basetypes.ListValue {
	if len(in) == 0 {
		return types.ListNull(types.StringType)
	}
	return types.ListValueMust(types.StringType, in)
}

func ToStringValue(in string) basetypes.StringValue {
	if in == "" {
		return types.StringNull()
	}
	return types.StringValue(in)
}
