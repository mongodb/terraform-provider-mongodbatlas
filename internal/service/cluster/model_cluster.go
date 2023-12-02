package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfConnectionStringDSModel struct {
	Standard          types.String `tfsdk:"standard"`
	StandardSrv       types.String `tfsdk:"standard_srv"`
	AwsPrivateLink    types.String `tfsdk:"aws_private_link"`
	AwsPrivateLinkSrv types.String `tfsdk:"aws_private_link_srv"`
	Private           types.String `tfsdk:"private"`
	PrivateSrv        types.String `tfsdk:"private_srv"`
	// PrivateEndpoint []tfPrivateEndpointModel `tfsdk:"private_endpoint"`
	PrivateEndpoint types.List `tfsdk:"private_endpoint"`
}

type tfLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type tfBiConnectorConfigModel struct {
	ReadPreference types.String `tfsdk:"read_preference"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

type tfAdvancedConfigurationModel struct {
	DefaultReadConcern               types.String `tfsdk:"default_read_concern"`
	DefaultWriteConcern              types.String `tfsdk:"default_write_concern"`
	MinimumEnabledTLSProtocol        types.String `tfsdk:"minimum_enabled_tls_protocol"`
	OplogSizeMB                      types.Int64  `tfsdk:"oplog_size_mb"`
	OplogMinRetentionHours           types.Int64  `tfsdk:"oplog_min_retention_hours"`
	SampleSizeBiConnector            types.Int64  `tfsdk:"sample_size_bi_connector"`
	SampleRefreshIntervalBiConnector types.Int64  `tfsdk:"sample_refresh_interval_bi_connector"`
	TransactionLifetimeLimitSeconds  types.Int64  `tfsdk:"transaction_lifetime_limit_seconds"`
	FailIndexKeyTooLong              types.Bool   `tfsdk:"fail_index_key_too_long"`
	JavascriptEnabled                types.Bool   `tfsdk:"javascript_enabled"`
	NoTableScan                      types.Bool   `tfsdk:"no_table_scan"`
}

var tfAdvancedConfigurationType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"default_read_concern":                 types.StringType,
	"default_write_concern":                types.StringType,
	"minimum_enabled_tls_protocol":         types.StringType,
	"oplog_size_mb":                        types.Int64Type,
	"oplog_min_retention_hours":            types.Int64Type,
	"sample_size_bi_connector":             types.Int64Type,
	"sample_refresh_interval_bi_connector": types.Int64Type,
	"transaction_lifetime_limit_seconds":   types.Int64Type,

	"fail_index_key_too_long": types.BoolType,
	"javascript_enabled":      types.BoolType,
	"no_table_scan":           types.BoolType,
}}

type tfTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
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

type tfPrivateEndpointModel struct {
	ConnectionString                  types.String `tfsdk:"connection_string"`
	SrvConnectionString               types.String `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString types.String `tfsdk:"srv_shard_optimized_connection_string"`
	EndpointType                      types.String `tfsdk:"type"`
	// Endpoints                         []tfEndpointModel `tfsdk:"endpoints"`
	Endpoints types.List `tfsdk:"endpoints"`
}

var tfPrivateEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"connection_string":                     types.StringType,
	"endpoints":                             types.ListType{ElemType: tfEndpointType},
	"srv_connection_string":                 types.StringType,
	"srv_shard_optimized_connection_string": types.StringType,
	"type":                                  types.StringType,
}}

type tfEndpointModel struct {
	EndpointID   types.String `tfsdk:"endpoint_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Region       types.String `tfsdk:"region"`
}

var tfEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"endpoint_id":   types.StringType,
	"provider_name": types.StringType,
	"region":        types.StringType,
},
}
