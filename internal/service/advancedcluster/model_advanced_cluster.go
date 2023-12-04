package advancedcluster

import (
	"context"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

type TfLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var TfLabelType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

type TfTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var TfTagType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

type TfBiConnectorConfigModel struct {
	ReadPreference types.String `tfsdk:"read_preference"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

var TfBiConnectorConfigType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"read_preference": types.StringType,
	"enabled":         types.BoolType,
}}

type TfAdvancedConfigurationModel struct {
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
	"fail_index_key_too_long":              types.BoolType,
	"javascript_enabled":                   types.BoolType,
	"no_table_scan":                        types.BoolType,
}}

type tfConnectionStringModel struct {
	Standard        types.String `tfsdk:"standard"`
	StandardSrv     types.String `tfsdk:"standard_srv"`
	Private         types.String `tfsdk:"private"`
	PrivateSrv      types.String `tfsdk:"private_srv"`
	PrivateEndpoint types.List   `tfsdk:"private_endpoint"`
}

var tfConnectionStringType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"standard":         types.StringType,
	"standard_srv":     types.StringType,
	"private":          types.StringType,
	"private_srv":      types.StringType,
	"private_endpoint": types.ListType{ElemType: TfPrivateEndpointType},
}}

type TfPrivateEndpointModel struct {
	ConnectionString                  types.String `tfsdk:"connection_string"`
	SrvConnectionString               types.String `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString types.String `tfsdk:"srv_shard_optimized_connection_string"`
	EndpointType                      types.String `tfsdk:"type"`
	Endpoints                         types.List   `tfsdk:"endpoints"`
}

var TfPrivateEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"connection_string":                     types.StringType,
	"endpoints":                             types.ListType{ElemType: TfEndpointType},
	"srv_connection_string":                 types.StringType,
	"srv_shard_optimized_connection_string": types.StringType,
	"type":                                  types.StringType,
}}

type TfEndpointModel struct {
	EndpointID   types.String `tfsdk:"endpoint_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Region       types.String `tfsdk:"region"`
}

var TfEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"endpoint_id":   types.StringType,
	"provider_name": types.StringType,
	"region":        types.StringType,
},
}

type tfReplicationSpecModel struct {
	RegionsConfigs types.Set    `tfsdk:"region_configs"`
	ContainerID    types.Map    `tfsdk:"container_id"`
	ID             types.String `tfsdk:"id"`
	ZoneName       types.String `tfsdk:"zone_name"`
	NumShards      types.Int64  `tfsdk:"num_shards"`
}

var tfReplicationSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":             types.StringType,
	"zone_name":      types.StringType,
	"num_shards":     types.Int64Type,
	"container_id":   types.MapType{ElemType: types.StringType},
	"region_configs": types.SetType{ElemType: tfRegionsConfigType},
},
}

type tfRegionsConfigModel struct {
	AnalyticsSpecs       types.List   `tfsdk:"analytics_specs"`
	AutoScaling          types.List   `tfsdk:"auto_scaling"`
	AnalyticsAutoScaling types.List   `tfsdk:"analytics_auto_scaling"`
	ReadOnlySpecs        types.List   `tfsdk:"read_only_specs"`
	ElectableSpecs       types.List   `tfsdk:"electable_specs"`
	BackingProviderName  types.String `tfsdk:"backing_provider_name"`
	ProviderName         types.String `tfsdk:"provider_name"`
	RegionName           types.String `tfsdk:"region_name"`
	Priority             types.Int64  `tfsdk:"priority"`
}

var tfRegionsConfigType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"backing_provider_name":  types.StringType,
	"priority":               types.Int64Type,
	"provider_name":          types.StringType,
	"region_name":            types.StringType,
	"analytics_specs":        types.ListType{ElemType: tfRegionsConfigSpecType},
	"electable_specs":        types.ListType{ElemType: tfRegionsConfigSpecType},
	"read_only_specs":        types.ListType{ElemType: tfRegionsConfigSpecType},
	"auto_scaling":           types.ListType{ElemType: tfRegionsConfigAutoScalingSpecType},
	"analytics_auto_scaling": types.ListType{ElemType: tfRegionsConfigAutoScalingSpecType},
}}

type tfRegionsConfigSpecsModel struct {
	DiskIOPS      types.Int64  `tfsdk:"disk_iops"`
	InstanceSize  types.String `tfsdk:"instance_size"`
	NodeCount     types.Int64  `tfsdk:"node_count"`
	EBSVolumeType types.String `tfsdk:"ebs_volume_type"`
}

var tfRegionsConfigSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_iops":       types.Int64Type,
	"ebs_volume_type": types.StringType,
	"instance_size":   types.StringType,
	"node_count":      types.Int64Type,
}}

type tfRegionsConfigAutoScalingSpecsModel struct {
	DiskGBEnabled           types.Bool   `tfsdk:"disk_gb_enabled"`
	ComputeScaleDownEnabled types.Bool   `tfsdk:"compute_scale_down_enabled"`
	ComputeMinInstanceSize  types.String `tfsdk:"compute_min_instance_size"`
	ComputeMaxInstanceSize  types.String `tfsdk:"compute_max_instance_size"`
	ComputeEnabled          types.Bool   `tfsdk:"compute_enabled"`
}

var tfRegionsConfigAutoScalingSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_gb_enabled":            types.BoolType,
	"compute_enabled":            types.BoolType,
	"compute_scale_down_enabled": types.BoolType,
	"compute_min_instance_size":  types.StringType,
	"compute_max_instance_size":  types.StringType,
}}

func newTfConnectionStringsModel(ctx context.Context, connString *matlas.ConnectionStrings) []*tfConnectionStringModel {
	res := []*tfConnectionStringModel{}

	if connString != nil {
		res = append(res, &tfConnectionStringModel{
			Standard:        types.StringValue(connString.Standard),
			StandardSrv:     types.StringValue(connString.StandardSrv),
			Private:         types.StringValue(connString.Private),
			PrivateSrv:      types.StringValue(connString.PrivateSrv),
			PrivateEndpoint: NewTfPrivateEndpointModel(ctx, connString.PrivateEndpoint),
		})
	}
	return res
}

func newTfRegionConfig(ctx context.Context, conn *matlas.Client, apiObject *matlas.AdvancedRegionConfig, projectID string) (tfRegionsConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return tfRegionsConfigModel{}, diags
	}

	tfRegionsConfig := tfRegionsConfigModel{
		BackingProviderName: conversion.StringNullIfEmpty(apiObject.BackingProviderName),
		ProviderName:        conversion.StringNullIfEmpty(apiObject.ProviderName),
		RegionName:          conversion.StringNullIfEmpty(apiObject.RegionName),
		Priority:            types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.Priority)),
	}

	tfRegionsConfig.AnalyticsSpecs, diags = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.AnalyticsSpecs))
	tfRegionsConfig.ElectableSpecs, diags = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.ElectableSpecs))
	tfRegionsConfig.ReadOnlySpecs, diags = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.ReadOnlySpecs))

	tfRegionsConfig.AnalyticsAutoScaling, diags = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, newTfRegionsConfigAutoScalingSpecsModel(apiObject.AnalyticsAutoScaling))
	tfRegionsConfig.AutoScaling, diags = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, newTfRegionsConfigAutoScalingSpecsModel(apiObject.AutoScaling))

	return tfRegionsConfig, diags
}

func newTfRegionsConfigSpecsModel(apiSpecs *matlas.Specs) []*tfRegionsConfigSpecsModel {
	res := make([]*tfRegionsConfigSpecsModel, 0)

	if apiSpecs != nil {
		res = append(res, &tfRegionsConfigSpecsModel{
			DiskIOPS:      types.Int64PointerValue(apiSpecs.DiskIOPS),
			InstanceSize:  conversion.StringNullIfEmpty(apiSpecs.InstanceSize),
			NodeCount:     types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiSpecs.NodeCount)),
			EBSVolumeType: conversion.StringNullIfEmpty(apiSpecs.EbsVolumeType),
		})
	}

	return res
}

func newTfRegionsConfigAutoScalingSpecsModel(apiSpecs *matlas.AdvancedAutoScaling) []*tfRegionsConfigAutoScalingSpecsModel {
	res := make([]*tfRegionsConfigAutoScalingSpecsModel, 0)

	if apiSpecs != nil && apiSpecs.Compute != nil {
		res = append(res, &tfRegionsConfigAutoScalingSpecsModel{
			DiskGBEnabled:           types.BoolPointerValue(apiSpecs.DiskGB.Enabled),
			ComputeEnabled:          types.BoolPointerValue(apiSpecs.Compute.Enabled),
			ComputeScaleDownEnabled: types.BoolPointerValue(apiSpecs.Compute.ScaleDownEnabled),
			ComputeMinInstanceSize:  conversion.StringNullIfEmpty(apiSpecs.Compute.MinInstanceSize),
			ComputeMaxInstanceSize:  conversion.StringNullIfEmpty(apiSpecs.Compute.MaxInstanceSize),
		})
	}

	return res
}

func NewTfBiConnectorConfigModel(biConnector *matlas.BiConnector) []*TfBiConnectorConfigModel {
	if biConnector == nil {
		return []*TfBiConnectorConfigModel{}
	}

	return []*TfBiConnectorConfigModel{
		{
			Enabled:        types.BoolPointerValue(biConnector.Enabled),
			ReadPreference: conversion.StringNullIfEmpty(biConnector.ReadPreference),
		},
	}
}

func NewTfTagsModel(tags *[]*matlas.Tag) []*TfTagModel {
	res := make([]*TfTagModel, len(*tags))

	for i, v := range *tags {
		res[i] = &TfTagModel{
			Key:   conversion.StringNullIfEmpty(v.Key),
			Value: conversion.StringNullIfEmpty(v.Value),
		}
	}

	return res
}

func NewTfLabelsModel(labels []matlas.Label) []TfLabelModel {
	out := make([]TfLabelModel, len(labels))
	for i, v := range labels {
		out[i] = TfLabelModel{
			Key:   conversion.StringNullIfEmpty(v.Key),
			Value: conversion.StringNullIfEmpty(v.Value),
		}
	}

	return out
}

func NewTfAdvancedConfigurationModel(p *matlas.ProcessArgs) []*TfAdvancedConfigurationModel {
	res := []*TfAdvancedConfigurationModel{
		{
			DefaultReadConcern:               conversion.StringNullIfEmpty(p.DefaultReadConcern),
			DefaultWriteConcern:              conversion.StringNullIfEmpty(p.DefaultWriteConcern),
			FailIndexKeyTooLong:              types.BoolPointerValue(p.FailIndexKeyTooLong),
			JavascriptEnabled:                types.BoolPointerValue(p.JavascriptEnabled),
			MinimumEnabledTLSProtocol:        conversion.StringNullIfEmpty(p.MinimumEnabledTLSProtocol),
			NoTableScan:                      types.BoolPointerValue(p.NoTableScan),
			OplogSizeMB:                      types.Int64PointerValue(p.OplogSizeMB),
			OplogMinRetentionHours:           conversion.FloatToTfInt64(p.OplogMinRetentionHours),
			SampleSizeBiConnector:            types.Int64PointerValue(p.SampleSizeBIConnector),
			SampleRefreshIntervalBiConnector: types.Int64PointerValue(p.SampleRefreshIntervalBIConnector),
			TransactionLifetimeLimitSeconds:  types.Int64PointerValue(p.TransactionLifetimeLimitSeconds),
		},
	}
	return res
}

func NewTfPrivateEndpointModel(ctx context.Context, privateEndpoints []matlas.PrivateEndpoint) types.List {
	res := make([]TfPrivateEndpointModel, len(privateEndpoints))

	for i, pe := range privateEndpoints {
		res[i] = TfPrivateEndpointModel{
			ConnectionString:                  conversion.StringNullIfEmpty(pe.ConnectionString),
			SrvConnectionString:               conversion.StringNullIfEmpty(pe.SRVConnectionString),
			SrvShardOptimizedConnectionString: conversion.StringNullIfEmpty(pe.SRVShardOptimizedConnectionString),
			EndpointType:                      conversion.StringNullIfEmpty(pe.Type),
			Endpoints:                         NewTFEndpointModel(ctx, pe.Endpoints),
		}
	}
	s, _ := types.ListValueFrom(ctx, TfPrivateEndpointType, res)
	return s
}

func NewTFEndpointModel(ctx context.Context, endpoints []matlas.Endpoint) types.List {
	res := make([]TfEndpointModel, len(endpoints))

	for i, e := range endpoints {
		res[i] = TfEndpointModel{
			Region:       conversion.StringNullIfEmpty(e.Region),
			ProviderName: conversion.StringNullIfEmpty(e.ProviderName),
			EndpointID:   conversion.StringNullIfEmpty(e.EndpointID),
		}
	}
	s, _ := types.ListValueFrom(ctx, TfEndpointType, res)
	return s
}
