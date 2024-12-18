package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

// AlignStateReplicationSpecsChanged ensures len(state.ReplicationSpecs) == len(plan.ReplicationSpecs) & each element in both state and plan is A) likely the same or B) use an empty replication spec
// If length already match, we expect either zone_name || regionNames to be the same, if both are updated, we assume the state cannot be used
// If length missmatch (add/remove) we require both zone_name and regionNames to match to allow using an existing state
// If an element doesn't match we set the Replication spec to the "empty" value to avoid using the `id` from state
func AlignStateReplicationSpecsChanged(ctx context.Context, state, plan *admin.ClusterDescription20240805) bool {
	stateSpecs := state.GetReplicationSpecs()
	planSpecs := plan.GetReplicationSpecs()
	var alignedSpecs []admin.ReplicationSpec20240805
	if len(stateSpecs) == len(planSpecs) {
		alignedSpecs = alignSpecs(ctx, &stateSpecs, &planSpecs, specsMatchPartial)
	} else {
		alignedSpecs = alignSpecs(ctx, &stateSpecs, &planSpecs, specsMatchFull)
	}
	state.ReplicationSpecs = &alignedSpecs
	return !reflect.DeepEqual(stateSpecs, alignedSpecs)
}

func alignSpecs(ctx context.Context, state, plan *[]admin.ReplicationSpec20240805, match func(admin.ReplicationSpec20240805, admin.ReplicationSpec20240805) bool) []admin.ReplicationSpec20240805 {
	remainingStateSpecs := make(map[int]admin.ReplicationSpec20240805)
	for i, stateSpec := range *state {
		remainingStateSpecs[i] = stateSpec
	}
	alignedSpecs := make([]admin.ReplicationSpec20240805, len(*plan))
	for i, planSpec := range *plan {
		for j := range *state {
			stateSpec, ok := remainingStateSpecs[j]
			if ok && match(stateSpec, planSpec) {
				alignedSpecs[i] = remainingStateSpecs[j]
				delete(remainingStateSpecs, j)
				break
			}
		}
	}
	for index, stateSpec := range remainingStateSpecs {
		tflog.Info(ctx, fmt.Sprintf("Replication spec %d in state does not match any spec in config, zone_name=%s, regions=%s, assuming it has been deleted", index, stateSpec.GetZoneName(), strings.Join(regionNames(stateSpec), ", ")))
	}
	for i, newStateSpec := range alignedSpecs {
		if update.IsZeroValues(&newStateSpec) {
			planSpec := (*plan)[i]
			tflog.Info(ctx, fmt.Sprintf("Couldn't match replication spec %d in config, assuming it is a new spec, zone_name=%s, regions=%s", i, planSpec.GetZoneName(), strings.Join(regionNames(planSpec), ", ")))
		}
	}
	return alignedSpecs
}

func specsMatchFull(state, plan admin.ReplicationSpec20240805) bool {
	return state.GetZoneName() == plan.GetZoneName() && reflect.DeepEqual(regionNames(state), regionNames(plan))
}

func specsMatchPartial(state, plan admin.ReplicationSpec20240805) bool {
	return state.GetZoneName() == plan.GetZoneName() || reflect.DeepEqual(regionNames(state), regionNames(plan))
}

func regionNames(spec admin.ReplicationSpec20240805) []string {
	names := make([]string, len(spec.GetRegionConfigs()))
	for i, config := range spec.GetRegionConfigs() {
		names[i] = config.GetRegionName()
	}
	return names
}
