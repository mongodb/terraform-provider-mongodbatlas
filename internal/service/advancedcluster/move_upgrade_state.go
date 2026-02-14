package advancedcluster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

// MoveState is used with moved block to upgrade from cluster to adv_cluster
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

// UpgradeState is used to upgrade from adv_cluster schema v1 (SDKv2) or v2 (TPF v2) to v3 (TPF v3).
func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {StateUpgrader: stateUpgrader},
		2: {StateUpgrader: stateUpgrader},
	}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_cluster" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}
	setStateResponse(ctx, &resp.Diagnostics, req.SourceRawState, &resp.TargetState)
}

func stateUpgrader(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	setStateResponse(ctx, &resp.Diagnostics, req.RawState, &resp.State)
	if resp.Diagnostics.HasError() {
		return
	}
	// Upgrade-specific: extract replication_specs from old raw JSON to preserve zone_name, disk_iops,
	// node_count, instance_size, disk_size_gb, etc. during schema upgrades (v1/v2 → v3).
	// Not done for the move path (cluster→advanced_cluster) because the source has a different schema.
	setReplicationSpecsFromOldJSON(ctx, &resp.Diagnostics, req.RawState.JSON, &resp.State)
}

// stateAttrs has the attributes needed from source schema.
// Filling these attributes in the destination will prevent plan changes when moving/upgrading state.
// Read will fill in the rest.
var stateAttrs = map[string]tftypes.Type{
	"project_id":             tftypes.String, // project_id and name to identify the cluster.
	"name":                   tftypes.String,
	"retain_backups_enabled": tftypes.Bool,   // TF specific so can't be got in Read.
	"backup_enabled":         tftypes.Bool,   // Optional-only in v3; was Optional+Computed in v1/v2.
	"mongo_db_major_version": tftypes.String, // Optional-only in v3, must be preserved from old state to avoid plan diff.
	"timeouts": tftypes.Object{ // TF specific so can't be got in Read.
		AttributeTypes: map[string]tftypes.Type{
			"create": tftypes.String,
			"update": tftypes.String,
			"delete": tftypes.String,
		},
	},
}

func setStateResponse(ctx context.Context, diags *diag.Diagnostics, stateIn *tfprotov6.RawState, stateOut *tfsdk.State) {
	rawStateValue, err := stateIn.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: stateAttrs,
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		diags.AddError("Unable to Unmarshal state", err.Error())
		return
	}
	var stateObj map[string]tftypes.Value
	if err := rawStateValue.As(&stateObj); err != nil {
		diags.AddError("Unable to Parse state", err.Error())
		return
	}
	projectID, name := getProjectIDNameFromStateObj(diags, stateObj)
	if diags.HasError() {
		return
	}
	model := newTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: projectID,
		Name:    name,
	}, diags)
	if diags.HasError() {
		return
	}
	model.AdvancedConfiguration = types.ObjectNull(advancedConfigurationObjType.AttrTypes)
	model.Timeouts = getTimeoutFromStateObj(stateObj)
	if diags.HasError() {
		return
	}
	setOptionalModelAttrs(stateObj, model)
	// Set tags and labels to null instead of empty so there is no plan change if there are no tags or labels when Read is called.
	model.Tags = types.MapNull(types.StringType)
	model.Labels = types.MapNull(types.StringType)
	diags.Append(stateOut.Set(ctx, model)...)
}

func getProjectIDNameFromStateObj(diags *diag.Diagnostics, stateObj map[string]tftypes.Value) (projectID, name *string) {
	projectID = schemafunc.GetAttrFromStateObj[string](stateObj, "project_id")
	name = schemafunc.GetAttrFromStateObj[string](stateObj, "name")
	if !conversion.IsStringPresent(projectID) || !conversion.IsStringPresent(name) {
		diags.AddError("Unable to read project_id or name from state", fmt.Sprintf("project_id: %s, name: %s",
			conversion.SafeString(projectID), conversion.SafeString(name)))
		return
	}
	return projectID, name
}

func getTimeoutFromStateObj(stateObj map[string]tftypes.Value) timeouts.Value {
	attrTypes := map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}
	nullObj := timeouts.Value{Object: types.ObjectNull(attrTypes)}
	timeoutState := schemafunc.GetAttrFromStateObj[map[string]tftypes.Value](stateObj, "timeouts")
	if timeoutState == nil {
		return nullObj
	}
	timeoutMap := make(map[string]attr.Value)
	for action := range attrTypes {
		actionTimeout := schemafunc.GetAttrFromStateObj[string](*timeoutState, action)
		if actionTimeout == nil {
			timeoutMap[action] = types.StringNull()
		} else {
			timeoutMap[action] = types.StringPointerValue(actionTimeout)
		}
	}
	obj, d := types.ObjectValue(attrTypes, timeoutMap)
	if d.HasError() {
		return nullObj
	}
	return timeouts.Value{Object: obj}
}

func setOptionalModelAttrs(stateObj map[string]tftypes.Value, model *TFModel) {
	if retainBackupsEnabled := schemafunc.GetAttrFromStateObj[bool](stateObj, "retain_backups_enabled"); retainBackupsEnabled != nil {
		model.RetainBackupsEnabled = types.BoolPointerValue(retainBackupsEnabled)
	}
	if v := schemafunc.GetAttrFromStateObj[bool](stateObj, "backup_enabled"); v != nil {
		model.BackupEnabled = types.BoolPointerValue(v)
	}
	if v := schemafunc.GetAttrFromStateObj[string](stateObj, "mongo_db_major_version"); v != nil {
		model.MongoDBMajorVersion = types.StringPointerValue(v)
	}
}

// setReplicationSpecsFromOldJSON extracts replication_specs from the old raw JSON state to preserve
// user-visible attributes during schema upgrades (v1/v2 → v3). Without this, Optional-only attributes
// are nullified by overrideReplicationSpecsWithPrevStateValue because the upgraded state has no
// replication_specs to compare against. This function parses the JSON directly because the nested
// structure differs between v1 (SDKv2 blocks with num_shards) and v2 (TPF objects without num_shards).
//
// Fields extracted: zone_name, num_shards, node_count, instance_size, disk_size_gb, disk_iops,
// backing_provider_name, analytics_specs, read_only_specs.
// Fields deliberately NOT extracted: ebs_volume_type, auto_scaling, analytics_auto_scaling
// (these are nullified by the nullification path and should remain null to avoid removal diffs).
func setReplicationSpecsFromOldJSON(ctx context.Context, diags *diag.Diagnostics, jsonData []byte, stateOut *tfsdk.State) {
	if len(jsonData) == 0 {
		return
	}
	oldSpecs := parseReplicationSpecsFromJSON(jsonData)
	if len(oldSpecs) == 0 {
		return
	}
	repSpecs := buildReplicationSpecsFromOldSpecs(ctx, diags, oldSpecs)
	if diags.HasError() {
		return
	}
	var model TFModel
	diags.Append(stateOut.Get(ctx, &model)...)
	if diags.HasError() {
		return
	}
	model.ReplicationSpecs = repSpecs
	diags.Append(stateOut.Set(ctx, model)...)
}

// oldReplicationSpec holds values extracted from the old JSON state for building v3 replication_specs.
type oldReplicationSpec struct {
	ZoneName      string
	RegionConfigs []oldRegionConfig
	NumShards     int
}

// oldRegionConfig holds extracted region_config values. auto_scaling and analytics_auto_scaling are
// deliberately NOT extracted — the nullification path (Path A) nullifies them, so they should remain
// null in the upgraded state to avoid removal diffs for users who don't configure them.
type oldRegionConfig struct {
	ElectableSpecs      *oldSpecsConfig
	AnalyticsSpecs      *oldSpecsConfig
	ReadOnlySpecs       *oldSpecsConfig
	BackingProviderName *string
}

// oldSpecsConfig holds extracted spec fields. ebs_volume_type is deliberately NOT extracted —
// the nullification path (Path A) nullifies it, so it should remain null to match that behavior.
type oldSpecsConfig struct {
	DiskIops     *int64
	DiskSizeGb   *float64
	NodeCount    *int64
	InstanceSize *string
}

// parseReplicationSpecsFromJSON extracts replication_specs from raw JSON state.
// Handles both v1 (SDKv2: electable_specs as list, has num_shards) and v2 (TPF: electable_specs as object).
func parseReplicationSpecsFromJSON(jsonData []byte) []oldReplicationSpec {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(jsonData, &raw); err != nil {
		return nil
	}
	specBytes, ok := raw["replication_specs"]
	if !ok {
		return nil
	}
	var rawSpecs []map[string]json.RawMessage
	if err := json.Unmarshal(specBytes, &rawSpecs); err != nil {
		return nil
	}
	result := make([]oldReplicationSpec, 0, len(rawSpecs))
	for _, rawSpec := range rawSpecs {
		spec := oldReplicationSpec{NumShards: 1}
		if v, err := unmarshalJSONField[string](rawSpec, "zone_name"); err == nil && v != nil {
			spec.ZoneName = *v
		}
		// num_shards only exists in v1 (SDKv2). In v3, each shard is a separate entry.
		if v, err := unmarshalJSONField[float64](rawSpec, "num_shards"); err == nil && v != nil && *v > 0 {
			spec.NumShards = int(*v)
		}
		spec.RegionConfigs = parseRegionConfigsFromJSON(rawSpec)
		result = append(result, spec)
	}
	return result
}

func parseRegionConfigsFromJSON(rawSpec map[string]json.RawMessage) []oldRegionConfig {
	rcBytes, ok := rawSpec["region_configs"]
	if !ok {
		return nil
	}
	var rawRCs []map[string]json.RawMessage
	if err := json.Unmarshal(rcBytes, &rawRCs); err != nil {
		return nil
	}
	result := make([]oldRegionConfig, 0, len(rawRCs))
	for _, rawRC := range rawRCs {
		rc := oldRegionConfig{
			ElectableSpecs: parseSpecsFromJSON(rawRC, "electable_specs"),
			AnalyticsSpecs: parseSpecsFromJSON(rawRC, "analytics_specs"),
			ReadOnlySpecs:  parseSpecsFromJSON(rawRC, "read_only_specs"),
		}
		if v, err := unmarshalJSONField[string](rawRC, "backing_provider_name"); err == nil && v != nil {
			rc.BackingProviderName = v
		}
		result = append(result, rc)
	}
	return result
}

// parseSpecsFromJSON extracts spec fields from a region_config JSON key (electable_specs,
// analytics_specs, or read_only_specs). Handles both v1 (SDKv2 TypeList MaxItems=1, i.e.,
// JSON array) and v2 (TPF SingleNestedAttribute, i.e., JSON object).
func parseSpecsFromJSON(rawRC map[string]json.RawMessage, key string) *oldSpecsConfig {
	esBytes, ok := rawRC[key]
	if !ok {
		return nil
	}
	// Try as object first (v2/v3 TPF format).
	var esObj map[string]json.RawMessage
	if err := json.Unmarshal(esBytes, &esObj); err == nil {
		return extractSpecsFields(esObj)
	}
	// Try as list (v1 SDKv2 format: TypeList MaxItems=1).
	var esList []map[string]json.RawMessage
	if err := json.Unmarshal(esBytes, &esList); err == nil && len(esList) > 0 {
		return extractSpecsFields(esList[0])
	}
	return nil
}

// extractSpecsFields extracts individual fields from a raw JSON specs object.
// Deliberately does NOT extract ebs_volume_type — Path A nullifies it, so Path B should too.
func extractSpecsFields(raw map[string]json.RawMessage) *oldSpecsConfig {
	specs := &oldSpecsConfig{}
	hasAnyField := false
	if v, err := unmarshalJSONField[float64](raw, "disk_iops"); err == nil && v != nil {
		iops := int64(*v)
		specs.DiskIops = &iops
		hasAnyField = true
	}
	if v, err := unmarshalJSONField[float64](raw, "disk_size_gb"); err == nil && v != nil {
		specs.DiskSizeGb = v
		hasAnyField = true
	}
	if v, err := unmarshalJSONField[float64](raw, "node_count"); err == nil && v != nil {
		nc := int64(*v)
		specs.NodeCount = &nc
		hasAnyField = true
	}
	if v, err := unmarshalJSONField[string](raw, "instance_size"); err == nil && v != nil {
		specs.InstanceSize = v
		hasAnyField = true
	}
	if !hasAnyField {
		return nil
	}
	return specs
}

func unmarshalJSONField[T any](raw map[string]json.RawMessage, key string) (*T, error) {
	data, ok := raw[key]
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// buildReplicationSpecsFromOldSpecs creates v3 TFReplicationSpecsModel entries from parsed old state.
// Expands num_shards: a single v1 entry with num_shards=N becomes N identical v3 entries.
func buildReplicationSpecsFromOldSpecs(ctx context.Context, diags *diag.Diagnostics, oldSpecs []oldReplicationSpec) types.List {
	var allSpecs []TFReplicationSpecsModel
	for _, old := range oldSpecs {
		regionConfigs := buildMinimalRegionConfigs(ctx, diags, old.RegionConfigs)
		if diags.HasError() {
			return types.ListNull(replicationSpecsObjType)
		}
		// Expand num_shards: in v1, one entry with num_shards=N means N shards sharing the same config.
		for range old.NumShards {
			allSpecs = append(allSpecs, TFReplicationSpecsModel{
				ZoneName:      types.StringValue(old.ZoneName),
				RegionConfigs: regionConfigs,
			})
		}
	}
	list, listDiags := types.ListValueFrom(ctx, replicationSpecsObjType, allSpecs)
	diags.Append(listDiags...)
	return list
}

func buildMinimalRegionConfigs(ctx context.Context, diags *diag.Diagnostics, oldRCs []oldRegionConfig) types.List {
	if len(oldRCs) == 0 {
		return types.ListNull(regionConfigsObjType)
	}
	rcs := make([]TFRegionConfigsModel, 0, len(oldRCs))
	for _, old := range oldRCs {
		esObj := buildSpecsObject(ctx, diags, old.ElectableSpecs)
		if diags.HasError() {
			return types.ListNull(regionConfigsObjType)
		}
		// For analytics_specs and read_only_specs: null out entirely if node_count is 0 or absent,
		// matching Path A's nullifySpecsIfNodeCountZero behavior.
		analyticsObj := buildSpecsObject(ctx, diags, specsIfNodeCountPositive(old.AnalyticsSpecs))
		if diags.HasError() {
			return types.ListNull(regionConfigsObjType)
		}
		readOnlyObj := buildSpecsObject(ctx, diags, specsIfNodeCountPositive(old.ReadOnlySpecs))
		if diags.HasError() {
			return types.ListNull(regionConfigsObjType)
		}
		rcs = append(rcs, TFRegionConfigsModel{
			ElectableSpecs:       esObj,
			AnalyticsSpecs:       analyticsObj,
			ReadOnlySpecs:        readOnlyObj,
			AutoScaling:          types.ObjectNull(autoScalingObjType.AttrTypes),
			AnalyticsAutoScaling: types.ObjectNull(autoScalingObjType.AttrTypes),
			BackingProviderName:  types.StringPointerValue(old.BackingProviderName),
			ProviderName:         types.StringNull(),
			RegionName:           types.StringNull(),
			Priority:             types.Int64Null(),
		})
	}
	list, listDiags := types.ListValueFrom(ctx, regionConfigsObjType, rcs)
	diags.Append(listDiags...)
	return list
}

// specsIfNodeCountPositive returns nil if the specs have no nodes (node_count absent or 0),
// causing the built model to have a null object. This matches Path A's nullifySpecsIfNodeCountZero.
func specsIfNodeCountPositive(specs *oldSpecsConfig) *oldSpecsConfig {
	if specs == nil || specs.NodeCount == nil || *specs.NodeCount == 0 {
		return nil
	}
	return specs
}

// buildSpecsObject creates a TFSpecsModel from parsed old state specs.
// Returns null object if old is nil. EbsVolumeType is always null (Path A nullifies it).
func buildSpecsObject(ctx context.Context, diags *diag.Diagnostics, old *oldSpecsConfig) types.Object {
	if old == nil {
		return types.ObjectNull(specsObjType.AttrTypes)
	}
	model := TFSpecsModel{
		DiskIops:      types.Int64PointerValue(old.DiskIops),
		DiskSizeGb:    types.Float64PointerValue(old.DiskSizeGb),
		NodeCount:     types.Int64PointerValue(old.NodeCount),
		InstanceSize:  types.StringPointerValue(old.InstanceSize),
		EbsVolumeType: types.StringNull(), // Deliberately null: Path A nullifies this.
	}
	obj, objDiags := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, model)
	diags.Append(objDiags...)
	return obj
}
