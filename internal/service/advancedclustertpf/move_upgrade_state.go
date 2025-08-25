package advancedclustertpf

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

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

// UpgradeState is used to upgrade from adv_cluster schema v1 (SDKv2) to v2 (TPF)
func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {StateUpgrader: stateUpgraderFromV1},
	}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_cluster" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}
	// Use always new sharding config when moving from cluster to adv_cluster
	setStateResponse(ctx, &resp.Diagnostics, req.SourceRawState, &resp.TargetState, false)
}

func stateUpgraderFromV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	// Use same sharding config as in SDKv2 when upgrading to TPF
	setStateResponse(ctx, &resp.Diagnostics, req.RawState, &resp.State, true)
}

// stateAttrs has the attributes needed from source schema.
// Filling these attributes in the destination will prevent plan changes when moving/upgrading state.
// Read will fill in the rest.
var stateAttrs = map[string]tftypes.Type{
	"project_id":             tftypes.String, // project_id and name to identify the cluster
	"name":                   tftypes.String,
	"retain_backups_enabled": tftypes.Bool,   // TF specific so can't be got in Read
	"mongo_db_major_version": tftypes.String, // Has special logic in overrideAttributesWithPrevStateValue that needs the previous state
	"timeouts": tftypes.Object{ // TF specific so can't be got in Read
		AttributeTypes: map[string]tftypes.Type{
			"create": tftypes.String,
			"update": tftypes.String,
			"delete": tftypes.String,
		},
	},
	"replication_specs": tftypes.List{ // Needed to send num_shards to Read so it can decide if it's using the legacy schema.
		ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"num_shards": tftypes.Number,
			},
		},
	},
}

func setStateResponse(ctx context.Context, diags *diag.Diagnostics, stateIn *tfprotov6.RawState, stateOut *tfsdk.State, allowOldShardingConfig bool) {
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
	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: projectID,
		Name:    name,
	}, diags, ExtraAPIInfo{})
	if diags.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, &ProcessArgs{
		ArgsDefault: nil,
		// ArgsLegacy:            nil,
		ClusterAdvancedConfig: nil,
	}, diags)
	model.Timeouts = getTimeoutFromStateObj(stateObj)
	if diags.HasError() {
		return
	}
	setOptionalModelAttrs(stateObj, model)
	if allowOldShardingConfig {
		setReplicationSpecNumShardsAttr(ctx, stateObj, model)
	}
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
	if mongoDBMajorVersion := schemafunc.GetAttrFromStateObj[string](stateObj, "mongo_db_major_version"); mongoDBMajorVersion != nil {
		model.MongoDBMajorVersion = types.StringPointerValue(mongoDBMajorVersion)
	}
}

func setReplicationSpecNumShardsAttr(ctx context.Context, stateObj map[string]tftypes.Value, model *TFModel) {
	specsVal := schemafunc.GetAttrFromStateObj[[]tftypes.Value](stateObj, "replication_specs")
	if specsVal == nil {
		return
	}
	var specModels []TFReplicationSpecsModel
	for _, specVal := range *specsVal {
		var specObj map[string]tftypes.Value
		if err := specVal.As(&specObj); err != nil {
			continue
		}
		if specModel := replicationSpecModelWithNumShards(specObj["num_shards"]); specModel != nil {
			specModels = append(specModels, *specModel)
		}
	}
	if len(specModels) > 0 {
		model.ReplicationSpecs, _ = types.ListValueFrom(ctx, ReplicationSpecsObjType, specModels)
	}
}

func replicationSpecModelWithNumShards(numShardsVal tftypes.Value) *TFReplicationSpecsModel {
	var numShardsFloat *big.Float
	if err := numShardsVal.As(&numShardsFloat); err != nil || numShardsFloat == nil {
		return nil
	}
	return &TFReplicationSpecsModel{
		RegionConfigs: types.ListNull(RegionConfigsObjType),
		ContainerId:   types.MapNull(types.StringType),
		ExternalId:    types.StringNull(),
		ZoneId:        types.StringNull(),
		ZoneName:      types.StringNull(),
	}
}
