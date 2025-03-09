package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func UseStateForUnknown2(ctx context.Context, diags *diag.Diagnostics, d *DiffHelper, state, plan *TFModel) {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, d.AttributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs2(ctx, diags, state, plan, d)
	}
}

func useStateForUnknownsReplicationSpecs2(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, d *DiffHelper) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan, d.AttributeChanges)
	keepUnknownsUnchangedSpec = append(keepUnknownsUnchangedSpec, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	if diags.HasError() {
		return
	}
	for i := range planRepSpecsTF {
		if i < len(stateRepSpecsTF) {
			keepUnknowns := keepUnknownsUnchangedSpec
			if d.AttributeChanges.ListIndexChanged("replication_specs", i) {
				keepUnknowns = determineKeepUnknownsChangedReplicationSpec(keepUnknownsUnchangedSpec, d.AttributeChanges, fmt.Sprintf("replication_specs[%d]", i))
			}
			d.UseStateForUnknown(ctx, diags, keepUnknowns, path.Root("replication_specs").AtListIndex(i))
		}
	}
}
