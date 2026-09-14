package advancedcluster

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312025/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

type MajorVersionOperator int

const (
	EqualOrHigher MajorVersionOperator = iota
	Higher
	EqualOrLower
	Lower
)

func MajorVersionCompatible(input *string, version float64, operator MajorVersionOperator) *bool {
	if !conversion.IsStringPresent(input) {
		return nil
	}
	value, err := strconv.ParseFloat(*input, 64)
	if err != nil {
		return nil
	}
	var result bool
	switch operator {
	case EqualOrHigher:
		result = value >= version
	case Higher:
		result = value > version
	case EqualOrLower:
		result = value <= version
	case Lower:
		result = value < version
	default:
		return nil
	}
	return &result
}

func containerIDKey(providerName, regionName string) string {
	return fmt.Sprintf("%s:%s", providerName, regionName)
}

// based on flattenAdvancedReplicationSpecRegionConfigs in model_advanced_cluster.go
func resolveContainerIDs(ctx context.Context, projectID string, cluster *admin.ClusterDescription20240805, api admin.NetworkPeeringAPI) (map[string]string, error) {
	containerIDs := map[string]string{}
	responseCache := map[string]*admin.PaginatedCloudProviderContainer{}
	for _, spec := range cluster.GetReplicationSpecs() {
		for _, regionConfig := range spec.GetRegionConfigs() {
			providerName := regionConfig.GetProviderName()
			if providerName == constant.TENANT {
				continue
			}
			params := &admin.ListGroupContainersApiParams{
				GroupId:      projectID,
				ProviderName: &providerName,
			}
			key := containerIDKey(providerName, regionConfig.GetRegionName())
			if _, ok := containerIDs[key]; ok {
				continue
			}
			var containersResponse *admin.PaginatedCloudProviderContainer
			var err error
			if response, ok := responseCache[providerName]; ok {
				containersResponse = response
			} else {
				containersResponse, _, err = api.ListGroupContainersWithParams(ctx, params).Execute()
				if err != nil {
					return nil, err
				}
				responseCache[providerName] = containersResponse
			}
			if results := getAdvancedClusterContainerID(containersResponse.GetResults(), &regionConfig); results != "" {
				containerIDs[key] = results
			} else {
				return nil, fmt.Errorf("container id not found for %s", key)
			}
		}
	}
	return containerIDs, nil
}

func OverrideAttributesWithPrevStateValue(ctx context.Context, modelIn, modelOut *TFModel, diags *diag.Diagnostics) {
	if modelIn == nil || modelOut == nil || diags == nil {
		return
	}
	overrideUnreportedInactiveSpecs(ctx, modelIn, modelOut, diags)
	beforeVersion := conversion.NilForUnknown(modelIn.MongoDBMajorVersion, modelIn.MongoDBMajorVersion.ValueStringPointer())
	afterVersion := conversion.NilForUnknown(modelOut.MongoDBMajorVersion, modelOut.MongoDBMajorVersion.ValueStringPointer())
	if beforeVersion != nil {
		warnIfMajorVersionChanged(*beforeVersion, afterVersion, diags)
		modelOut.MongoDBMajorVersion = types.StringPointerValue(beforeVersion)
	}
	overrideMapStringWithPrevStateValue(&modelIn.Labels, &modelOut.Labels)
	overrideMapStringWithPrevStateValue(&modelIn.Tags, &modelOut.Tags)

	// Copy Terraform-only attributes which are not returned by Atlas API.
	// These fields are Optional-only or Optional+Computed (because of a default value),
	// so no need for more complex logic as they can't be Unknown in the plan.
	modelOut.Timeouts = modelIn.Timeouts
	modelOut.DeleteOnCreateTimeout = modelIn.DeleteOnCreateTimeout
	modelOut.RetainBackupsEnabled = modelIn.RetainBackupsEnabled
	modelOut.UseEffectiveFields = modelIn.UseEffectiveFields
}

// overrideUnreportedInactiveSpecs keeps previously known hardware spec attributes that Atlas didn't
// report back for a spec describing no nodes (node_count = 0).
//
// A spec with no nodes has no hardware provisioned for it, so Atlas may omit its hardware attributes
// instead of echoing what was sent. This is most visible with use_effective_fields enabled, where
// the specs returned describe the effective (actually provisioned) hardware.
//
// Writing those omissions to state as null contradicts a value the practitioner set in the config,
// which Terraform reports as "Provider produced inconsistent result after apply", and which would
// otherwise show up as a permanent diff on every subsequent plan. modelIn is the plan in Create and
// Update and the prior state in Read, so in every case it holds the value to fall back to.
func overrideUnreportedInactiveSpecs(ctx context.Context, modelIn, modelOut *TFModel, diags *diag.Diagnostics) {
	// Flex clusters and unresolved configs can have no usable replication specs to reconcile.
	for _, repSpecs := range []types.List{modelIn.ReplicationSpecs, modelOut.ReplicationSpecs} {
		if !isKnown(repSpecs) {
			return
		}
	}
	inRepSpecs := TFModelList[TFReplicationSpecsModel](ctx, diags, modelIn.ReplicationSpecs)
	outRepSpecs := TFModelList[TFReplicationSpecsModel](ctx, diags, modelOut.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	anyChanged := false
	for i := range minLen(outRepSpecs, inRepSpecs) {
		inRegionConfigs := TFModelList[TFRegionConfigsModel](ctx, diags, inRepSpecs[i].RegionConfigs)
		outRegionConfigs := TFModelList[TFRegionConfigsModel](ctx, diags, outRepSpecs[i].RegionConfigs)
		if diags.HasError() {
			return
		}
		regionConfigsChanged := false
		for j := range minLen(outRegionConfigs, inRegionConfigs) {
			specs := []struct{ in, out *types.Object }{
				{&inRegionConfigs[j].AnalyticsSpecs, &outRegionConfigs[j].AnalyticsSpecs},
				{&inRegionConfigs[j].ElectableSpecs, &outRegionConfigs[j].ElectableSpecs},
				{&inRegionConfigs[j].ReadOnlySpecs, &outRegionConfigs[j].ReadOnlySpecs},
			}
			for _, spec := range specs {
				if overrideUnreportedSpecAttrs(ctx, diags, spec.in, spec.out) {
					regionConfigsChanged = true
				}
				if diags.HasError() {
					return
				}
			}
		}
		if !regionConfigsChanged {
			continue
		}
		listRegionConfigs, diagsLocal := types.ListValueFrom(ctx, regionConfigsObjType, outRegionConfigs)
		diags.Append(diagsLocal...)
		if diags.HasError() {
			return
		}
		outRepSpecs[i].RegionConfigs = listRegionConfigs
		anyChanged = true
	}
	if !anyChanged {
		return
	}
	listRepSpecs, diagsLocal := types.ListValueFrom(ctx, replicationSpecsObjType, outRepSpecs)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return
	}
	modelOut.ReplicationSpecs = listRepSpecs
}

// overrideUnreportedSpecAttrs fills the attributes of specOut that Atlas didn't report, using specIn.
// It returns true if specOut was changed. Only specs describing no nodes are considered: when a spec
// has nodes, Atlas reports the hardware actually backing them and that value must win.
func overrideUnreportedSpecAttrs(ctx context.Context, diags *diag.Diagnostics, specIn, specOut *types.Object) bool {
	if specIn.IsNull() || specOut.IsNull() {
		return false
	}
	in := TFModelObject[TFSpecsModel](ctx, *specIn)
	out := TFModelObject[TFSpecsModel](ctx, *specOut)
	if in == nil || out == nil {
		return false
	}
	nodeCount := out.NodeCount // prefer what Atlas reported, fall back to the previous model
	if !isKnown(nodeCount) {
		nodeCount = in.NodeCount
	}
	if nodeCount.ValueInt64() != 0 {
		return false
	}
	changed := copyAttrIfUnreported(&in.DiskSizeGb, &out.DiskSizeGb)
	changed = copyAttrIfUnreported(&in.DiskIops, &out.DiskIops) || changed
	changed = copyAttrIfUnreported(&in.EbsVolumeType, &out.EbsVolumeType) || changed
	changed = copyAttrIfUnreported(&in.InstanceSize, &out.InstanceSize) || changed
	changed = copyAttrIfUnreported(&in.NodeCount, &out.NodeCount) || changed
	if !changed {
		return false
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, out)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return false
	}
	*specOut = objType
	return true
}

// copyAttrIfUnreported copies src into dest when Atlas reported no value (dest is null) and src holds
// a usable one, returning true if the copy happened. src can come from the plan, where computed
// attributes are unknown, and unknown must never reach the state, so src is only used when known.
func copyAttrIfUnreported[T attr.Value](src, dest *T) bool {
	if !(*dest).IsNull() || !isKnown(*src) {
		return false
	}
	*dest = *src
	return true
}

func warnIfMajorVersionChanged(before string, after *string, diags *diag.Diagnostics) {
	if after == nil || before == *after || majorComponent(before) == majorComponent(*after) {
		return
	}
	diags.AddWarning(
		"MongoDB major version modified outside of Terraform",
		fmt.Sprintf("Atlas reports mongo_db_major_version as %q but your Terraform state has %q. "+
			"Your cluster's major version may have been modified outside of Terraform. "+
			"Consider setting mongo_db_major_version = %q in your configuration and applying the changes. "+
			"This warning will continue until you update your configuration. "+
			"In an upcoming major version of the provider, this drift will result in a non-empty plan.", *after, before, *after),
	)
}

func majorComponent(version string) string {
	major, _, _ := strings.Cut(version, ".")
	return major
}

func overrideMapStringWithPrevStateValue(mapIn, mapOut *types.Map) {
	if mapIn == nil || mapOut == nil || len(mapOut.Elements()) > 0 {
		return
	}
	if mapIn.IsNull() {
		*mapOut = types.MapNull(types.StringType)
	} else {
		*mapOut = types.MapValueMust(types.StringType, nil)
	}
}
