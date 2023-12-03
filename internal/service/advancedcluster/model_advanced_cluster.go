package advancedcluster

import (
	"context"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

type TfLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type TfTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type TfBiConnectorConfigModel struct {
	ReadPreference types.String `tfsdk:"read_preference"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

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

type tfConnectionStringModel struct {
	Standard        types.String `tfsdk:"standard"`
	StandardSrv     types.String `tfsdk:"standard_srv"`
	Private         types.String `tfsdk:"private"`
	PrivateSrv      types.String `tfsdk:"private_srv"`
	PrivateEndpoint types.List   `tfsdk:"private_endpoint"`
}

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
	DiskIOPS      types.String `tfsdk:"disk_iops"`
	InstanceSize  types.String `tfsdk:"instance_size"`
	NodeCount     types.String `tfsdk:"node_count"`
	EBSVolumeType types.Int64  `tfsdk:"ebs_volume_type"`
}

var tfRegionsConfigSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_iops":       types.StringType,
	"ebs_volume_type": types.Int64Type,
	"instance_size":   types.StringType,
	"node_count":      types.StringType,
}}

type tfRegionsConfigAutoScalingSpecsModel struct {
	DiskGBEnabled           types.String `tfsdk:"disk_gb_enabled"`
	ComputeScaleDownEnabled types.String `tfsdk:"compute_scale_down_enabled"`
	ComputeMinInstanceSize  types.String `tfsdk:"compute_min_instance_size"`
	ComputeMaxInstanceSize  types.String `tfsdk:"compute_max_instance_size"`
	ComputeEnabled          types.Int64  `tfsdk:"compute_enabled"`
}

var tfRegionsConfigAutoScalingSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_gb_enabled":            types.StringType,
	"compute_enabled":            types.Int64Type,
	"compute_scale_down_enabled": types.StringType,
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
			PrivateEndpoint: NewTFPrivateEndpointModel(ctx, connString.PrivateEndpoint),
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

func NewTFTagsModel(tags *[]*matlas.Tag) []*TfTagModel {
	res := make([]*TfTagModel, len(*tags))

	for i, v := range *tags {
		res[i] = &TfTagModel{
			Key:   conversion.StringNullIfEmpty(v.Key),
			Value: conversion.StringNullIfEmpty(v.Value),
		}
	}

	return res
}

func NewTFLabelsModel(labels []matlas.Label) []TfLabelModel {
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

func NewTFPrivateEndpointModel(ctx context.Context, privateEndpoints []matlas.PrivateEndpoint) types.List {
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
