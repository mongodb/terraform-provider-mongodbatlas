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

func AssertDiagsOK(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if diags.HasError() {
		t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
	}
}
