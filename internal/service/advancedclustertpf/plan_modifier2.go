package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func UseStateForUnknownBasedOnChanges(ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, schema SimplifiedSchema, state, plan *TFModel) {
	attributeChanges := d.AttributeChanges()
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)

}
