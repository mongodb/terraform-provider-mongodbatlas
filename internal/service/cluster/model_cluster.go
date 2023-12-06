package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfConnectionStringDSModel struct {
	Standard          types.String `tfsdk:"standard"`
	StandardSrv       types.String `tfsdk:"standard_srv"`
	AwsPrivateLink    types.Map    `tfsdk:"aws_private_link"`
	AwsPrivateLinkSrv types.Map    `tfsdk:"aws_private_link_srv"`
	Private           types.String `tfsdk:"private"`
	PrivateSrv        types.String `tfsdk:"private_srv"`
	PrivateEndpoint   types.List   `tfsdk:"private_endpoint"`
}

type tfReplicationSpecModel struct {
	ID            types.String          `tfsdk:"id"`
	ZoneName      types.String          `tfsdk:"zone_name"`
	RegionsConfig []tfRegionConfigModel `tfsdk:"regions_config"`
	NumShards     types.Int64           `tfsdk:"num_shards"`
}

type tfRegionConfigModel struct {
	RegionName     types.String `tfsdk:"region_name"`
	ElectableNodes types.Int64  `tfsdk:"electable_nodes"`
	Priority       types.Int64  `tfsdk:"priority"`
	ReadOnlyNodes  types.Int64  `tfsdk:"read_only_nodes"`
	AnalyticsNodes types.Int64  `tfsdk:"analytics_nodes"`
}

type tfSnapshotPolicyModel struct {
	ID         types.String `tfsdk:"id"`
	PolicyItem types.List   `tfsdk:"policy_item"`
}

type tfSnapshotBackupPolicyModel struct {
	ClusterID             types.String `tfsdk:"cluster_id"`
	ClusterName           types.String `tfsdk:"cluster_name"`
	NextSnapshot          types.String `tfsdk:"next_snapshot"`
	Policies              types.List   `tfsdk:"policies"`
	ReferenceHourOfDay    types.Int64  `tfsdk:"reference_hour_of_day"`
	ReferenceMinuteOfHour types.Int64  `tfsdk:"reference_minute_of_hour"`
	RestoreWindowDays     types.Int64  `tfsdk:"restore_window_days"`
	UpdateSnapshots       types.Bool   `tfsdk:"update_snapshots"`
}

var tfSnapshotPolicyType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":          types.StringType,
	"policy_item": types.ListType{ElemType: tfSnapshotPolicyItemType},
}}

type tfSnapshotPolicyItemModel struct {
	ID                types.String `tfsdk:"id"`
	FrequencyType     types.String `tfsdk:"frequency_type"`
	RetentionUnit     types.String `tfsdk:"retention_unit"`
	FrequencyInterval types.Int64  `tfsdk:"frequency_interval"`
	RetentionValue    types.Int64  `tfsdk:"retention_value"`
}

var tfSnapshotPolicyItemType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":                 types.StringType,
	"frequency_type":     types.StringType,
	"retention_unit":     types.StringType,
	"frequency_interval": types.Int64Type,
	"retention_value":    types.Int64Type,
}}
