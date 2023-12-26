package advancedcluster

import (
	"context"
	"fmt"
	"log"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

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

type tfReplicationSpecDSModel struct {
	RegionsConfigs types.Set    `tfsdk:"region_configs"`
	ContainerID    types.Map    `tfsdk:"container_id"`
	ID             types.String `tfsdk:"id"`
	ZoneName       types.String `tfsdk:"zone_name"`
	NumShards      types.Int64  `tfsdk:"num_shards"`
}

var tfReplicationSpecDSType = types.ObjectType{AttrTypes: map[string]attr.Type{
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
	InstanceSize  types.String `tfsdk:"instance_size"`
	EBSVolumeType types.String `tfsdk:"ebs_volume_type"`
	DiskIOPS      types.Int64  `tfsdk:"disk_iops"`
	NodeCount     types.Int64  `tfsdk:"node_count"`
}

var tfRegionsConfigSpecType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"disk_iops":       types.Int64Type,
	"ebs_volume_type": types.StringType,
	"instance_size":   types.StringType,
	"node_count":      types.Int64Type,
}}

type tfRegionsConfigAutoScalingSpecsModel struct {
	ComputeMinInstanceSize  types.String `tfsdk:"compute_min_instance_size"`
	ComputeMaxInstanceSize  types.String `tfsdk:"compute_max_instance_size"`
	DiskGBEnabled           types.Bool   `tfsdk:"disk_gb_enabled"`
	ComputeScaleDownEnabled types.Bool   `tfsdk:"compute_scale_down_enabled"`
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
			Standard:        conversion.StringNullIfEmpty(connString.Standard),
			StandardSrv:     conversion.StringNullIfEmpty(connString.StandardSrv),
			Private:         conversion.StringNullIfEmpty(connString.Private),
			PrivateSrv:      conversion.StringNullIfEmpty(connString.PrivateSrv),
			PrivateEndpoint: newTfPrivateEndpointModel(ctx, connString.PrivateEndpoint),
		})
	}
	return res
}

func newTfRegionConfig(ctx context.Context, conn *matlas.Client, apiObject *matlas.AdvancedRegionConfig, projectID string) (tfRegionsConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	if apiObject == nil {
		return tfRegionsConfigModel{}, diags
	}

	providerName := apiObject.ProviderName
	tfRegionsConfig := tfRegionsConfigModel{
		BackingProviderName: conversion.StringNullIfEmpty(apiObject.BackingProviderName),
		ProviderName:        conversion.StringNullIfEmpty(providerName),
		RegionName:          conversion.StringNullIfEmpty(apiObject.RegionName),
		Priority:            types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.Priority)),
	}

	tfRegionsConfig.AnalyticsSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.AnalyticsSpecs, providerName))
	diags.Append(d...)
	tfRegionsConfig.ElectableSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.ElectableSpecs, providerName))
	diags.Append(d...)
	tfRegionsConfig.ReadOnlySpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, newTfRegionsConfigSpecsModel(apiObject.ReadOnlySpecs, providerName))
	diags.Append(d...)
	tfRegionsConfig.AnalyticsAutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, newTfRegionsConfigAutoScalingSpecsModel(apiObject.AnalyticsAutoScaling))
	diags.Append(d...)
	tfRegionsConfig.AutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, newTfRegionsConfigAutoScalingSpecsModel(apiObject.AutoScaling))
	diags.Append(d...)

	return tfRegionsConfig, diags
}

func newTfRegionsConfigSpecsModel(apiSpecs *matlas.Specs, providerName string) []*tfRegionsConfigSpecsModel {
	res := make([]*tfRegionsConfigSpecsModel, 0)

	if apiSpecs != nil {
		tmp := &tfRegionsConfigSpecsModel{
			InstanceSize: conversion.StringNullIfEmpty(apiSpecs.InstanceSize),
			NodeCount:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiSpecs.NodeCount)),
		}
		if providerName == "AWS" {
			tmp.DiskIOPS = types.Int64PointerValue(apiSpecs.DiskIOPS)
			tmp.EBSVolumeType = conversion.StringNullIfEmpty(apiSpecs.InstanceSize)
		}
		res = append(res, tmp)
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

func newTfBiConnectorConfigModel(biConnector *matlas.BiConnector) []*TfBiConnectorConfigModel {
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

func newTfTagsModel(tags *[]*matlas.Tag) []*TfTagModel {
	res := make([]*TfTagModel, len(*tags))

	for i, v := range *tags {
		res[i] = &TfTagModel{
			Key:   conversion.StringNullIfEmpty(v.Key),
			Value: conversion.StringNullIfEmpty(v.Value),
		}
	}

	return res
}

func newTfLabelsModel(labels []matlas.Label) []TfLabelModel {
	out := make([]TfLabelModel, len(labels))

	for i, v := range labels {
		out[i] = TfLabelModel{
			Key:   conversion.StringNullIfEmpty(v.Key),
			Value: conversion.StringNullIfEmpty(v.Value),
		}
	}

	return out
}

func newTfAdvancedConfigurationModel(p *matlas.ProcessArgs) []*TfAdvancedConfigurationModel {
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

func newTfPrivateEndpointModel(ctx context.Context, privateEndpoints []matlas.PrivateEndpoint) types.List {
	res := make([]TfPrivateEndpointModel, len(privateEndpoints))

	for i, pe := range privateEndpoints {
		res[i] = TfPrivateEndpointModel{
			ConnectionString:                  conversion.StringNullIfEmpty(pe.ConnectionString),
			SrvConnectionString:               conversion.StringNullIfEmpty(pe.SRVConnectionString),
			SrvShardOptimizedConnectionString: conversion.StringNullIfEmpty(pe.SRVShardOptimizedConnectionString),
			EndpointType:                      conversion.StringNullIfEmpty(pe.Type),
			Endpoints:                         newTFEndpointModel(ctx, pe.Endpoints),
		}
	}
	s, _ := types.ListValueFrom(ctx, TfPrivateEndpointType, res)
	return s
}

func newTFEndpointModel(ctx context.Context, endpoints []matlas.Endpoint) types.List {
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

func newTags(ctx context.Context, tfSet basetypes.SetValue) []*matlas.Tag {
	if tfSet.IsNull() || len(tfSet.Elements()) == 0 {
		return nil
	}
	var tfArr []TfTagModel
	tfSet.ElementsAs(ctx, &tfArr, true)

	res := make([]*matlas.Tag, len(tfArr))
	for i, v := range tfArr {
		res[i] = &matlas.Tag{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return res
}

func newLabels(ctx context.Context, tfSet basetypes.SetValue) []matlas.Label {
	if tfSet.IsNull() || len(tfSet.Elements()) == 0 {
		return nil
	}

	var tfArr []TfLabelModel
	tfSet.ElementsAs(ctx, &tfArr, true)

	res := make([]matlas.Label, len(tfArr))

	for i, v := range tfArr {
		res[i] = matlas.Label{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}

	return res
}

func newBiConnectorConfig(ctx context.Context, tfList basetypes.ListValue) *matlas.BiConnector {
	// res := matlas.BiConnector{}
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		// if isUpdate {
		// 	return &res
		// }
		return nil
	}

	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfArr []TfBiConnectorConfigModel
	tfList.ElementsAs(ctx, &tfArr, true)

	tfBiConnector := tfArr[0]

	biConnector := matlas.BiConnector{
		Enabled:        tfBiConnector.Enabled.ValueBoolPointer(),
		ReadPreference: tfBiConnector.ReadPreference.ValueString(),
	}

	return &biConnector
}

func newReplicationSpecs(ctx context.Context, tfList basetypes.ListValue) []*matlas.AdvancedReplicationSpec {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfRepSpecs []tfReplicationSpecRSModel
	tfList.ElementsAs(ctx, &tfRepSpecs, true)

	var repSpecs []*matlas.AdvancedReplicationSpec

	for i := range tfRepSpecs {
		rs := newReplicationSpec(ctx, &tfRepSpecs[i])
		repSpecs = append(repSpecs, rs)
	}
	return repSpecs
}

func newReplicationSpec(ctx context.Context, tfRepSpec *tfReplicationSpecRSModel) *matlas.AdvancedReplicationSpec {
	if tfRepSpec == nil {
		return nil
	}

	zoneName := tfRepSpec.ZoneName.ValueString()
	if !conversion.IsStringPresent(&zoneName) {
		zoneName = DefaultZoneName
	}
	res := &matlas.AdvancedReplicationSpec{
		NumShards:     int(tfRepSpec.NumShards.ValueInt64()),
		ZoneName:      zoneName,
		RegionConfigs: newRegionConfigs(ctx, tfRepSpec.RegionsConfigs),
	}

	if v := tfRepSpec.ID; !v.IsUnknown() {
		res.ID = v.ValueString()
	}
	return res
}

func newRegionConfigs(ctx context.Context, tfList basetypes.ListValue) []*matlas.AdvancedRegionConfig {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfRegionConfigs []tfRegionsConfigModel
	tfList.ElementsAs(ctx, &tfRegionConfigs, true)

	var regionConfigs []*matlas.AdvancedRegionConfig

	for i := range tfRegionConfigs {
		rc := newRegionConfig(ctx, &tfRegionConfigs[i])

		regionConfigs = append(regionConfigs, rc)
	}

	return regionConfigs
}

func newRegionConfig(ctx context.Context, tfRegionConfig *tfRegionsConfigModel) *matlas.AdvancedRegionConfig {
	if tfRegionConfig == nil {
		return nil
	}

	providerName := tfRegionConfig.ProviderName.ValueString()
	apiObject := &matlas.AdvancedRegionConfig{
		Priority:     conversion.Int64PtrToIntPtr(tfRegionConfig.Priority.ValueInt64Pointer()),
		ProviderName: providerName,
		RegionName:   tfRegionConfig.RegionName.ValueString(),
	}

	if v := tfRegionConfig.AnalyticsSpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AnalyticsSpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.ElectableSpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.ElectableSpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.ReadOnlySpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.ReadOnlySpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.AutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AutoScaling = newRegionConfigAutoScalingSpec(ctx, v)
	}
	if v := tfRegionConfig.AnalyticsAutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AnalyticsAutoScaling = newRegionConfigAutoScalingSpec(ctx, v)
	}
	if v := tfRegionConfig.BackingProviderName; !v.IsNull() && v.ValueString() != defaultString {
		apiObject.BackingProviderName = v.ValueString()
	}

	return apiObject
}

func newRegionConfigAutoScalingSpec(ctx context.Context, tfList basetypes.ListValue) *matlas.AdvancedAutoScaling {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var specs []tfRegionsConfigAutoScalingSpecsModel
	tfList.ElementsAs(ctx, &specs, true)

	spec := specs[0]
	advancedAutoScaling := &matlas.AdvancedAutoScaling{}
	diskGB := &matlas.DiskGB{}
	compute := &matlas.Compute{}

	if v := spec.DiskGBEnabled; !v.IsUnknown() {
		diskGB.Enabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeEnabled; !v.IsUnknown() {
		compute.Enabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeScaleDownEnabled; !v.IsUnknown() {
		compute.ScaleDownEnabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeMinInstanceSize; !v.IsUnknown() {
		value := compute.ScaleDownEnabled
		if *value {
			compute.MinInstanceSize = v.ValueString()
		}
	}
	if v := spec.ComputeMaxInstanceSize; !v.IsUnknown() {
		value := compute.Enabled
		if *value {
			compute.MaxInstanceSize = v.ValueString()
		}
	}

	advancedAutoScaling.DiskGB = diskGB
	advancedAutoScaling.Compute = compute

	return advancedAutoScaling
}

func newRegionConfigSpec(ctx context.Context, tfList basetypes.ListValue, providerName string) *matlas.Specs {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var specs []tfRegionsConfigSpecsModel
	tfList.ElementsAs(ctx, &specs, true)

	spec := specs[0]
	apiObject := &matlas.Specs{}

	if providerName == "AWS" {
		if v := spec.DiskIOPS; !v.IsNull() && v.ValueInt64() > 0 {
			apiObject.DiskIOPS = v.ValueInt64Pointer()
		}
		if v := spec.EBSVolumeType; !v.IsNull() && v.ValueString() != defaultString {
			apiObject.EbsVolumeType = v.ValueString()
		}
	}

	if v := spec.InstanceSize; !v.IsNull() {
		apiObject.InstanceSize = v.ValueString()
	}
	if v := spec.NodeCount; !v.IsNull() && v.ValueInt64() > 0 {
		apiObject.NodeCount = conversion.Int64PtrToIntPtr(v.ValueInt64Pointer())
	}
	return apiObject
}

func newAdvancedConfiguration(ctx context.Context, tfList basetypes.ListValue) *matlas.ProcessArgs {
	res := &matlas.ProcessArgs{}

	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		// if isUpdate {
		// 	return res
		// }
		return nil
	}
	// else if isUpdate && (tfList.IsNull() || len(tfList.Elements()) == 0) { // if during update user removed the advanced_configuration block
	// 	return &matlas.ProcessArgs{}
	// }

	var tfAdvancedConfigArr []TfAdvancedConfigurationModel
	tfList.ElementsAs(ctx, &tfAdvancedConfigArr, true)

	if len(tfAdvancedConfigArr) == 0 {
		return nil
	}
	tfModel := tfAdvancedConfigArr[0]

	if v := tfModel.DefaultReadConcern; !v.IsUnknown() {
		res.DefaultReadConcern = v.ValueString()
	}
	if v := tfModel.DefaultWriteConcern; !v.IsUnknown() {
		res.DefaultWriteConcern = v.ValueString()
	}

	if v := tfModel.FailIndexKeyTooLong; !v.IsUnknown() {
		res.FailIndexKeyTooLong = v.ValueBoolPointer()
	}

	if v := tfModel.JavascriptEnabled; !v.IsUnknown() {
		res.JavascriptEnabled = v.ValueBoolPointer()
	}

	if v := tfModel.MinimumEnabledTLSProtocol; !v.IsUnknown() {
		res.MinimumEnabledTLSProtocol = v.ValueString()
	}

	if v := tfModel.NoTableScan; !v.IsUnknown() {
		res.NoTableScan = v.ValueBoolPointer()
	}

	if v := tfModel.SampleSizeBiConnector; !v.IsUnknown() {
		res.SampleSizeBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.SampleRefreshIntervalBiConnector; !v.IsUnknown() {
		res.SampleRefreshIntervalBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.OplogSizeMB; !v.IsUnknown() {
		if sizeMB := v.ValueInt64(); sizeMB != 0 {
			res.OplogSizeMB = v.ValueInt64Pointer()
		} else {
			log.Printf(ErrorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
		}
	}

	if v := tfModel.OplogMinRetentionHours; !v.IsNull() {
		if minRetentionHours := v.ValueInt64(); minRetentionHours >= 0 {
			res.OplogMinRetentionHours = pointy.Float64(cast.ToFloat64(v.ValueInt64()))
		} else {
			log.Printf(ErrorClusterSetting, `oplog_min_retention_hours`, "", cast.ToString(minRetentionHours))
		}
	}

	if v := tfModel.TransactionLifetimeLimitSeconds; !v.IsUnknown() {
		if transactionLimitSeconds := v.ValueInt64(); transactionLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = v.ValueInt64Pointer()
		} else {
			log.Printf(ErrorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLimitSeconds))
		}
	}

	return res
}

func newTfReplicationSpecRSModel(ctx context.Context, apiObject *matlas.AdvancedReplicationSpec, configSpec *tfReplicationSpecRSModel,
	conn *matlas.Client, projectID string) (*tfReplicationSpecRSModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return nil, diags
	}

	tfMap := tfReplicationSpecRSModel{}
	tfMap.NumShards = types.Int64Value(cast.ToInt64(apiObject.NumShards))
	tfMap.ID = types.StringValue(apiObject.ID)
	if configSpec != nil {
		object, containerIds, diags := newTfRegionsConfigsRSModel(ctx, apiObject.RegionConfigs, configSpec.RegionsConfigs, conn, projectID)
		if diags.HasError() {
			return nil, diags
		}
		l, diags := types.ListValueFrom(ctx, tfRegionsConfigType, object)
		if diags.HasError() {
			return nil, diags
		}
		tfMap.RegionsConfigs = l
		tfMap.ContainerID = containerIds
	} else {
		object, containerIds, diags := newTfRegionsConfigsRSModel(ctx, apiObject.RegionConfigs, types.ListNull(tfRegionsConfigType), conn, projectID)
		if diags.HasError() {
			return nil, diags
		}
		l, diags := types.ListValueFrom(ctx, tfRegionsConfigType, object)
		if diags.HasError() {
			return nil, diags
		}
		tfMap.RegionsConfigs = l
		tfMap.ContainerID = containerIds
	}
	tfMap.ZoneName = types.StringValue(apiObject.ZoneName)

	return &tfMap, diags
}

func newTfRegionsConfigsRSModel(ctx context.Context, apiObjects []*matlas.AdvancedRegionConfig, configRegionConfigsList types.List,
	conn *matlas.Client, projectID string) (tfResult []tfRegionsConfigModel, containersIDs types.Map, diags1 diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiObjects) == 0 {
		return nil, types.MapNull(types.StringType), diags
	}

	var configRegionConfigs []*tfRegionsConfigModel
	containerIDsMap := map[string]attr.Value{}

	if !configRegionConfigsList.IsNull() { // create return to state - filter by config, read/tf plan - filter by config, update - filter by config, import - return everything from API
		configRegionConfigsList.ElementsAs(ctx, &configRegionConfigs, true)
	}

	var tfList []tfRegionsConfigModel

	for i, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		if len(configRegionConfigs) > i {
			tfMapObject := configRegionConfigs[i]
			rc, diags := newTfRegionsConfigRSModel(ctx, apiObject, tfMapObject)
			if diags.HasError() {
				break
			}

			tfList = append(tfList, *rc)
		} else {
			rc, diags := newTfRegionsConfigRSModel(ctx, apiObject, nil)
			if diags.HasError() {
				break
			}

			tfList = append(tfList, *rc)
		}

		if apiObject.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, projectID,
				&matlas.ContainersListOptions{ProviderName: apiObject.ProviderName})
			if err != nil {
				diags.AddError("error when getting containers list from Atlas", err.Error())
				return nil, types.MapNull(types.StringType), diags
			}
			if result := getAdvancedClusterContainerID(containers, apiObject); result != "" {
				// Will print as "providerName:regionName" = "containerId" in terraform show
				key := fmt.Sprintf("%s:%s", apiObject.ProviderName, apiObject.RegionName)
				containerIDsMap[key] = types.StringValue(result)
			}
		}
	}
	tfContainersIDsMap, _ := types.MapValue(types.StringType, containerIDsMap)

	return tfList, tfContainersIDsMap, diags
}

func newTfRegionsConfigRSModel(ctx context.Context, apiObject *matlas.AdvancedRegionConfig, configRegionConfig *tfRegionsConfigModel) (*tfRegionsConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	if apiObject == nil {
		return nil, diags
	}

	tfMap := tfRegionsConfigModel{}
	if configRegionConfig != nil {
		if v := configRegionConfig.AnalyticsSpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AnalyticsSpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.AnalyticsSpecs, apiObject.ProviderName, configRegionConfig.AnalyticsSpecs)
		} else {
			//  tfMap.AnalyticsSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
			tfMap.AnalyticsSpecs = types.ListNull(tfRegionsConfigSpecType)

		}
		diags.Append(d...)
		if v := configRegionConfig.ElectableSpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.ElectableSpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.ElectableSpecs, apiObject.ProviderName, configRegionConfig.ElectableSpecs)
		} else {
			// tfMap.ElectableSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
			tfMap.ElectableSpecs = types.ListNull(tfRegionsConfigSpecType)

		}
		diags.Append(d...)
		if v := configRegionConfig.ReadOnlySpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.ReadOnlySpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.ReadOnlySpecs, apiObject.ProviderName, configRegionConfig.ReadOnlySpecs)
		} else {
			// tfMap.ReadOnlySpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
			tfMap.ReadOnlySpecs = types.ListNull(tfRegionsConfigSpecType)

		}
		diags.Append(d...)
		if v := configRegionConfig.AutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AutoScaling, d = newTfRegionsConfigAutoScalingSpecsRSModel(ctx, apiObject.AutoScaling)
		} else {
			// tfMap.AutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, []tfRegionsConfigAutoScalingSpecsModel{})
			tfMap.AutoScaling = types.ListNull(tfRegionsConfigAutoScalingSpecType)

		}
		diags.Append(d...)
		if v := configRegionConfig.AnalyticsAutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AnalyticsAutoScaling, d = newTfRegionsConfigAutoScalingSpecsRSModel(ctx, apiObject.AnalyticsAutoScaling)
		} else {
			// tfMap.AnalyticsAutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, []tfRegionsConfigAutoScalingSpecsModel{})
			tfMap.AnalyticsAutoScaling = types.ListNull(tfRegionsConfigAutoScalingSpecType)

		}
		diags.Append(d...)
	} else {
		nilSpecList := types.ListNull(tfRegionsConfigSpecType)
		tfMap.AnalyticsSpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.AnalyticsSpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.ElectableSpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.ElectableSpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.ReadOnlySpecs, d = newTfRegionsConfigSpecRSModel(ctx, apiObject.ReadOnlySpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.AutoScaling, d = newTfRegionsConfigAutoScalingSpecsRSModel(ctx, apiObject.AutoScaling)
		diags.Append(d...)
		tfMap.AnalyticsAutoScaling, d = newTfRegionsConfigAutoScalingSpecsRSModel(ctx, apiObject.AnalyticsAutoScaling)
		diags.Append(d...)
	}

	tfMap.RegionName = types.StringValue(apiObject.RegionName)
	tfMap.ProviderName = types.StringValue(apiObject.ProviderName)
	tfMap.BackingProviderName = conversion.StringNullIfEmpty(apiObject.BackingProviderName)
	tfMap.Priority = types.Int64Value(cast.ToInt64(apiObject.Priority))

	return &tfMap, diags
}

func newTfRegionsConfigSpecRSModel(ctx context.Context, apiObject *matlas.Specs, providerName string, tfMapObjects types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ListNull(tfRegionsConfigSpecType), diags
	}

	var configRegionConfigSpecs []*tfRegionsConfigSpecsModel

	if !tfMapObjects.IsNull() { // create return to state - filter by config, read/tf plan - filter by config, update - filter by config, import - return everything from API
		tfMapObjects.ElementsAs(ctx, &configRegionConfigSpecs, true)
	}

	var tfList []tfRegionsConfigSpecsModel

	tfMap := tfRegionsConfigSpecsModel{}

	if len(configRegionConfigSpecs) > 0 {
		tfMapObject := configRegionConfigSpecs[0]

		if providerName == "AWS" {
			if cast.ToInt64(apiObject.DiskIOPS) > 0 {
				tfMap.DiskIOPS = types.Int64PointerValue(apiObject.DiskIOPS)
			} else {
				tfMap.DiskIOPS = types.Int64Null()
			}
			// if v := tfMapObject.EBSVolumeType; !v.IsNull() && v.ValueString() != "" {
			// 	tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
			// }
			if v := tfMapObject.EBSVolumeType; !v.IsNull() {
				tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
			}

		}
		if v := tfMapObject.NodeCount; !v.IsNull() {
			tfMap.NodeCount = types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.NodeCount))
		}
		if v := tfMapObject.InstanceSize; !v.IsNull() && v.ValueString() != "" {
			tfMap.InstanceSize = types.StringValue(apiObject.InstanceSize)
		}

		// if tfMap.DiskIOPS.IsNull() {
		// 	tfMap.DiskIOPS = types.Int64Value(defaultInt)
		// }
		// if tfMap.NodeCount.IsNull() {
		// 	tfMap.NodeCount = types.Int64Value(defaultInt)
		// }
		// if tfMap.EBSVolumeType.IsNull() {
		// 	tfMap.EBSVolumeType = types.StringValue(defaultString)
		// }
		tfList = append(tfList, tfMap)
	} else {
		tfMap.DiskIOPS = types.Int64PointerValue(apiObject.DiskIOPS)
		tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
		tfMap.NodeCount = types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.NodeCount))
		tfMap.InstanceSize = types.StringValue(apiObject.InstanceSize)
		// if tfMap.DiskIOPS.IsNull() {
		// 	tfMap.DiskIOPS = types.Int64Value(defaultInt)
		// }
		// if tfMap.NodeCount.IsNull() {
		// 	tfMap.NodeCount = types.Int64Value(defaultInt)
		// }
		// if tfMap.EBSVolumeType.IsNull() {
		// 	tfMap.EBSVolumeType = types.StringValue(defaultString)
		// }
		tfList = append(tfList, tfMap)
	}

	return types.ListValueFrom(ctx, tfRegionsConfigSpecType, tfList)
}

func newTfRegionsConfigAutoScalingSpecsRSModel(ctx context.Context, apiObject *matlas.AdvancedAutoScaling) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ListNull(tfRegionsConfigAutoScalingSpecType), diags
	}

	var tfList []tfRegionsConfigAutoScalingSpecsModel

	tfMap := tfRegionsConfigAutoScalingSpecsModel{}
	if apiObject.DiskGB != nil {
		tfMap.DiskGBEnabled = types.BoolPointerValue(apiObject.DiskGB.Enabled)
	}
	if apiObject.Compute != nil {
		tfMap.ComputeEnabled = types.BoolPointerValue(apiObject.Compute.Enabled)
		tfMap.ComputeScaleDownEnabled = types.BoolPointerValue(apiObject.Compute.ScaleDownEnabled)
		tfMap.ComputeMinInstanceSize = types.StringValue(apiObject.Compute.MinInstanceSize)
		tfMap.ComputeMaxInstanceSize = types.StringValue(apiObject.Compute.MaxInstanceSize)
	}

	tfList = append(tfList, tfMap)

	return types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, tfList)
}
