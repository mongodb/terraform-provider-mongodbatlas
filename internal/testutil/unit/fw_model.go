package unit

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TFObjectValue(t *testing.T, objType types.ObjectType, attributes any) types.Object {
	t.Helper()
	object, diags := types.ObjectValueFrom(t.Context(), objType.AttrTypes, attributes)
	AssertDiagsOK(t, diags)
	return object
}

func TFListValue(t *testing.T, elementType types.ObjectType, tfList any) types.List {
	t.Helper()
	list, diags := types.ListValueFrom(t.Context(), elementType, tfList)
	AssertDiagsOK(t, diags)
	return list
}

func AssertDiagsOK(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if diags.HasError() {
		t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
	}
}
