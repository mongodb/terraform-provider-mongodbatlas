package cluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

type ProcessArgs struct {
	argsDefault           *admin.ClusterDescriptionProcessArgs20240805
	clusterAdvancedConfig *matlas.AdvancedConfiguration
	argsLegacy            *admin20240530.ClusterDescriptionProcessArgs
}

func flattenCloudProviderSnapshotBackupPolicy(ctx context.Context, d *schema.ResourceData, conn *matlas.Client, projectID, clusterName string) ([]map[string]any, error) {
	backupPolicy, res, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound ||
			strings.Contains(err.Error(), "BACKUP_CONFIG_NOT_FOUND") ||
			strings.Contains(err.Error(), "Not Found") ||
			strings.Contains(err.Error(), "404") {
			return []map[string]any{}, nil
		}

		return []map[string]any{}, fmt.Errorf(ErrorSnapshotBackupPolicyRead, clusterName, err)
	}

	return []map[string]any{
		{
			"cluster_id":               backupPolicy.ClusterID,
			"cluster_name":             backupPolicy.ClusterName,
			"next_snapshot":            backupPolicy.NextSnapshot,
			"reference_hour_of_day":    backupPolicy.ReferenceHourOfDay,
			"reference_minute_of_hour": backupPolicy.ReferenceMinuteOfHour,
			"restore_window_days":      backupPolicy.RestoreWindowDays,
			"update_snapshots":         cast.ToBool(backupPolicy.UpdateSnapshots),
			"policies":                 flattenPolicies(backupPolicy.Policies),
		},
	}, nil
}

func flattenPolicies(policies []matlas.Policy) []map[string]any {
	actionList := make([]map[string]any, 0)
	for _, v := range policies {
		actionList = append(actionList, map[string]any{
			"id":          v.ID,
			"policy_item": flattenPolicyItems(v.PolicyItems),
		})
	}

	return actionList
}

func flattenPolicyItems(items []matlas.PolicyItem) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, v := range items {
		policyItems = append(policyItems, map[string]any{
			"id":                 v.ID,
			"frequency_interval": v.FrequencyInterval,
			"frequency_type":     v.FrequencyType,
			"retention_unit":     v.RetentionUnit,
			"retention_value":    v.RetentionValue,
		})
	}

	return policyItems
}

func flattenProcessArgs(p *ProcessArgs) []map[string]any {
	flattenedProcessArgs := []map[string]any{
		{
			// default_read_concern and fail_index_key_too_long have been deprecated, hence using the older SDK
			"default_read_concern":                 p.argsLegacy.DefaultReadConcern,
			"fail_index_key_too_long":              cast.ToBool(p.argsLegacy.FailIndexKeyTooLong),
			"default_write_concern":                p.argsDefault.DefaultWriteConcern,
			"javascript_enabled":                   cast.ToBool(p.argsDefault.JavascriptEnabled),
			"no_table_scan":                        cast.ToBool(p.argsDefault.NoTableScan),
			"oplog_size_mb":                        p.argsDefault.OplogSizeMB,
			"oplog_min_retention_hours":            p.argsDefault.OplogMinRetentionHours,
			"sample_size_bi_connector":             p.argsDefault.SampleSizeBIConnector,
			"sample_refresh_interval_bi_connector": p.argsDefault.SampleRefreshIntervalBIConnector,
			"transaction_lifetime_limit_seconds":   p.argsDefault.TransactionLifetimeLimitSeconds,
			"minimum_enabled_tls_protocol":         p.argsDefault.MinimumEnabledTlsProtocol,
			"tls_cipher_config_mode":               p.argsDefault.TlsCipherConfigMode,
			"custom_openssl_cipher_config_tls12":   conversion.SliceFromPtr(p.argsDefault.CustomOpensslCipherConfigTls12),
		},
	}

	if p.argsDefault.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != nil {
		flattenedProcessArgs[0]["change_stream_options_pre_and_post_images_expire_after_seconds"] = p.argsDefault.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds
	} else {
		flattenedProcessArgs[0]["change_stream_options_pre_and_post_images_expire_after_seconds"] = -1 // default in schema, otherwise user gets drift detection
	}

	if p.clusterAdvancedConfig != nil { // For TENANT cluster type, advancedConfiguration field may not be returned from cluster APIs
		flattenedProcessArgs[0]["minimum_enabled_tls_protocol"] = p.clusterAdvancedConfig.MinimumEnabledTLSProtocol
		flattenedProcessArgs[0]["tls_cipher_config_mode"] = p.clusterAdvancedConfig.TLSCipherConfigMode
		flattenedProcessArgs[0]["custom_openssl_cipher_config_tls12"] = conversion.SliceFromPtr(p.clusterAdvancedConfig.CustomOpensslCipherConfigTLS12)
	}

	return flattenedProcessArgs
}

func flattenLabels(l []matlas.Label) []map[string]any {
	labels := make([]map[string]any, len(l))
	for i, v := range l {
		labels[i] = map[string]any{
			"key":   v.Key,
			"value": v.Value,
		}
	}
	return labels
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	var pos int

	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			list = append(list[:pos], list[pos+1:]...)

			if pos > 0 {
				pos--
			}

			continue
		}
		pos++
	}

	return list
}

func flattenTags(l *[]*matlas.Tag) []map[string]any {
	if l == nil {
		return []map[string]any{}
	}
	tags := make([]map[string]any, len(*l))
	for i, v := range *l {
		tags[i] = map[string]any{
			"key":   v.Key,
			"value": v.Value,
		}
	}
	return tags
}

func flattenConnectionStrings(connectionStrings *matlas.ConnectionStrings) []map[string]any {
	connections := make([]map[string]any, 0)

	connections = append(connections, map[string]any{
		"standard":         connectionStrings.Standard,
		"standard_srv":     connectionStrings.StandardSrv,
		"private":          connectionStrings.Private,
		"private_srv":      connectionStrings.PrivateSrv,
		"private_endpoint": flattenPrivateEndpoint(connectionStrings.PrivateEndpoint),
	})

	return connections
}

func flattenPrivateEndpoint(privateEndpoints []matlas.PrivateEndpoint) []map[string]any {
	endpoints := make([]map[string]any, 0)
	for _, endpoint := range privateEndpoints {
		endpoints = append(endpoints, map[string]any{
			"connection_string":                     endpoint.ConnectionString,
			"srv_connection_string":                 endpoint.SRVConnectionString,
			"srv_shard_optimized_connection_string": endpoint.SRVShardOptimizedConnectionString,
			"endpoints":                             flattenEndpoints(endpoint.Endpoints),
			"type":                                  endpoint.Type,
		})
	}
	return endpoints
}

func flattenEndpoints(listEndpoints []matlas.Endpoint) []map[string]any {
	endpoints := make([]map[string]any, 0)
	for _, endpoint := range listEndpoints {
		endpoints = append(endpoints, map[string]any{
			"region":        endpoint.Region,
			"provider_name": endpoint.ProviderName,
			"endpoint_id":   endpoint.EndpointID,
		})
	}
	return endpoints
}

func flattenBiConnectorConfig(biConnector *matlas.BiConnector) []any {
	return []any{
		map[string]any{
			"enabled":         *biConnector.Enabled,
			"read_preference": biConnector.ReadPreference,
		},
	}
}

func expandBiConnectorConfig(d *schema.ResourceData) (*matlas.BiConnector, error) {
	var biConnector matlas.BiConnector

	if v, ok := d.GetOk("bi_connector_config"); ok {
		biConn := v.([]any)
		if len(biConn) > 0 {
			biConnMap := biConn[0].(map[string]any)

			enabled := cast.ToBool(biConnMap["enabled"])

			biConnector = matlas.BiConnector{
				Enabled:        &enabled,
				ReadPreference: cast.ToString(biConnMap["read_preference"]),
			}
		}
	}

	return &biConnector, nil
}

func expandTagSliceFromSetSchema(d *schema.ResourceData) []*matlas.Tag {
	list := d.Get("tags").(*schema.Set)
	res := make([]*matlas.Tag, list.Len())
	for i, val := range list.List() {
		v := val.(map[string]any)
		res[i] = &matlas.Tag{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}
	return res
}

func expandClusterAdvancedConfiguration(d *schema.ResourceData) *matlas.AdvancedConfiguration {
	ac := d.Get("advanced_configuration")
	if aclist, ok1 := ac.([]any); ok1 && len(aclist) > 0 {
		p := aclist[0].(map[string]any)
		res := matlas.AdvancedConfiguration{}

		if _, ok := d.GetOkExists("advanced_configuration.0.minimum_enabled_tls_protocol"); ok {
			res.MinimumEnabledTLSProtocol = conversion.StringPtr(cast.ToString(p["minimum_enabled_tls_protocol"]))
		}

		if _, ok := d.GetOkExists("advanced_configuration.0.tls_cipher_config_mode"); ok {
			res.TLSCipherConfigMode = conversion.StringPtr(cast.ToString(p["tls_cipher_config_mode"]))
		}

		if _, ok := d.GetOkExists("advanced_configuration.0.custom_openssl_cipher_config_tls12"); ok {
			tmp := conversion.ExpandStringListFromSetSchema(d.Get("advanced_configuration.0.custom_openssl_cipher_config_tls12").(*schema.Set))
			res.CustomOpensslCipherConfigTLS12 = &tmp
		}
		return &res
	}
	return nil
}

func expandProcessArgs(d *schema.ResourceData, p map[string]any, mongodbMajorVersion *string) (admin20240530.ClusterDescriptionProcessArgs, admin.ClusterDescriptionProcessArgs20240805) {
	res20240530 := admin20240530.ClusterDescriptionProcessArgs{}
	res := admin.ClusterDescriptionProcessArgs20240805{}

	// default_read_concern and fail_index_key_too_long have been deprecated, hence using the older SDK
	if _, ok := d.GetOkExists("advanced_configuration.0.default_read_concern"); ok {
		res20240530.DefaultReadConcern = conversion.StringPtr(cast.ToString(p["default_read_concern"]))
	}
	if _, ok := d.GetOkExists("advanced_configuration.0.fail_index_key_too_long"); ok {
		res20240530.FailIndexKeyTooLong = conversion.Pointer(cast.ToBool(p["fail_index_key_too_long"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.default_write_concern"); ok {
		res.DefaultWriteConcern = conversion.StringPtr(cast.ToString(p["default_write_concern"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.javascript_enabled"); ok {
		res.JavascriptEnabled = conversion.Pointer(cast.ToBool(p["javascript_enabled"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.no_table_scan"); ok {
		res.NoTableScan = conversion.Pointer(cast.ToBool(p["no_table_scan"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.sample_size_bi_connector"); ok {
		res.SampleSizeBIConnector = conversion.Pointer(cast.ToInt(p["sample_size_bi_connector"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.sample_refresh_interval_bi_connector"); ok {
		res.SampleRefreshIntervalBIConnector = conversion.Pointer(cast.ToInt(p["sample_refresh_interval_bi_connector"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.oplog_size_mb"); ok {
		if sizeMB := cast.ToInt64(p["oplog_size_mb"]); sizeMB != 0 {
			res.OplogSizeMB = conversion.Pointer(cast.ToInt(p["oplog_size_mb"]))
		} else {
			log.Printf(advancedcluster.ErrorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
		}
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.oplog_min_retention_hours"); ok {
		if minRetentionHours := cast.ToFloat64(p["oplog_min_retention_hours"]); minRetentionHours >= 0 {
			res.OplogMinRetentionHours = conversion.Pointer(cast.ToFloat64(p["oplog_min_retention_hours"]))
		} else {
			log.Printf(advancedcluster.ErrorClusterSetting, `oplog_min_retention_hours`, "", cast.ToString(minRetentionHours))
		}
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.transaction_lifetime_limit_seconds"); ok {
		if transactionLifetimeLimitSeconds := cast.ToInt64(p["transaction_lifetime_limit_seconds"]); transactionLifetimeLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = conversion.Pointer(cast.ToInt64(p["transaction_lifetime_limit_seconds"]))
		} else {
			log.Printf(advancedcluster.ErrorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLifetimeLimitSeconds))
		}
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.change_stream_options_pre_and_post_images_expire_after_seconds"); ok && advancedcluster.IsChangeStreamOptionsMinRequiredMajorVersion(mongodbMajorVersion) {
		res.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds = conversion.Pointer(cast.ToInt(p["change_stream_options_pre_and_post_images_expire_after_seconds"]))
	}

	return res20240530, res
}

func expandLabelSliceFromSetSchema(d *schema.ResourceData) []matlas.Label {
	list := d.Get("labels").(*schema.Set)
	res := make([]matlas.Label, list.Len())

	for i, val := range list.List() {
		v := val.(map[string]any)
		res[i] = matlas.Label{
			Key:   v["key"].(string),
			Value: v["value"].(string),
		}
	}

	return res
}

func expandReplicationSpecs(d *schema.ResourceData) ([]matlas.ReplicationSpec, error) {
	rSpecs := make([]matlas.ReplicationSpec, 0)

	vRSpecs, okRSpecs := d.GetOk("replication_specs")
	vPRName, okPRName := d.GetOk("provider_region_name")

	if okRSpecs {
		for _, s := range vRSpecs.(*schema.Set).List() {
			spec := s.(map[string]any)

			replaceRegion := ""
			originalRegion := ""
			id := ""

			if okPRName && d.Get("provider_name").(string) == "GCP" && cast.ToString(d.Get("cluster_type")) == "REPLICASET" {
				if d.HasChange("provider_region_name") {
					replaceRegion = vPRName.(string)
					original, _ := d.GetChange("provider_region_name")
					originalRegion = original.(string)
				}
			}

			if d.HasChange("replication_specs") {
				// Get original and new object
				var oldSpecs map[string]any
				original, _ := d.GetChange("replication_specs")
				for _, s := range original.(*schema.Set).List() {
					oldSpecs = s.(map[string]any)
					if spec["zone_name"].(string) == cast.ToString(oldSpecs["zone_name"]) {
						id = oldSpecs["id"].(string)
						break
					}
				}
				// If there was an item before and after then use the same id assuming it's the same replication spec
				if id == "" && oldSpecs != nil && len(vRSpecs.(*schema.Set).List()) == 1 && len(original.(*schema.Set).List()) == 1 {
					id = oldSpecs["id"].(string)
				}
			}

			regionsConfig, err := expandRegionsConfig(spec["regions_config"].(*schema.Set).List(), originalRegion, replaceRegion)
			if err != nil {
				return rSpecs, err
			}

			rSpec := matlas.ReplicationSpec{
				ID:            id,
				NumShards:     conversion.Pointer(cast.ToInt64(spec["num_shards"])),
				ZoneName:      cast.ToString(spec["zone_name"]),
				RegionsConfig: regionsConfig,
			}
			rSpecs = append(rSpecs, rSpec)
		}
	}

	return rSpecs, nil
}

func flattenReplicationSpecs(rSpecs []matlas.ReplicationSpec) []map[string]any {
	specs := make([]map[string]any, 0)

	for _, rSpec := range rSpecs {
		spec := map[string]any{
			"id":             rSpec.ID,
			"num_shards":     rSpec.NumShards,
			"zone_name":      cast.ToString(rSpec.ZoneName),
			"regions_config": flattenRegionsConfig(rSpec.RegionsConfig),
		}
		specs = append(specs, spec)
	}

	return specs
}

func expandRegionsConfig(regions []any, originalRegion, replaceRegion string) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, r := range regions {
		region := r.(map[string]any)

		r, err := conversion.ValRegion(region["region_name"])
		if err != nil {
			return regionsConfig, err
		}

		if replaceRegion != "" && r == originalRegion {
			r, err = conversion.ValRegion(replaceRegion)
		}
		if err != nil {
			return regionsConfig, err
		}

		regionsConfig[r] = matlas.RegionsConfig{
			AnalyticsNodes: conversion.Pointer(cast.ToInt64(region["analytics_nodes"])),
			ElectableNodes: conversion.Pointer(cast.ToInt64(region["electable_nodes"])),
			Priority:       conversion.Pointer(cast.ToInt64(region["priority"])),
			ReadOnlyNodes:  conversion.Pointer(cast.ToInt64(region["read_only_nodes"])),
		}
	}

	return regionsConfig, nil
}

func flattenRegionsConfig(regionsConfig map[string]matlas.RegionsConfig) []map[string]any {
	regions := make([]map[string]any, 0)

	for regionName, regionConfig := range regionsConfig {
		region := map[string]any{
			"region_name":     regionName,
			"priority":        regionConfig.Priority,
			"analytics_nodes": regionConfig.AnalyticsNodes,
			"electable_nodes": regionConfig.ElectableNodes,
			"read_only_nodes": regionConfig.ReadOnlyNodes,
		}
		regions = append(regions, region)
	}

	return regions
}

func expandProviderSetting(d *schema.ResourceData) (*matlas.ProviderSettings, error) {
	var (
		region, _          = conversion.ValRegion(d.Get("provider_region_name"))
		minInstanceSize    = getInstanceSizeToInt(d.Get("provider_auto_scaling_compute_min_instance_size").(string))
		maxInstanceSize    = getInstanceSizeToInt(d.Get("provider_auto_scaling_compute_max_instance_size").(string))
		instanceSize       = getInstanceSizeToInt(d.Get("provider_instance_size_name").(string))
		compute            *matlas.Compute
		autoScalingEnabled = d.Get("auto_scaling_compute_enabled").(bool)
		providerName       = cast.ToString(d.Get("provider_name"))
	)

	if minInstanceSize != 0 && autoScalingEnabled {
		if instanceSize < minInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_min_instance_size` must be lower than `provider_instance_size_name`")
		}

		compute = &matlas.Compute{
			MinInstanceSize: d.Get("provider_auto_scaling_compute_min_instance_size").(string),
		}
	}

	if maxInstanceSize != 0 && autoScalingEnabled {
		if instanceSize > maxInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_max_instance_size` must be higher than `provider_instance_size_name`")
		}

		if compute == nil {
			compute = &matlas.Compute{}
		}
		compute.MaxInstanceSize = d.Get("provider_auto_scaling_compute_max_instance_size").(string)
	}

	providerSettings := &matlas.ProviderSettings{
		InstanceSizeName: cast.ToString(d.Get("provider_instance_size_name")),
		ProviderName:     providerName,
		RegionName:       region,
		VolumeType:       cast.ToString(d.Get("provider_volume_type")),
	}

	if d.HasChange("provider_disk_type_name") {
		_, newdiskTypeName := d.GetChange("provider_disk_type_name")
		diskTypeName := cast.ToString(newdiskTypeName)
		if diskTypeName != "" { // ensure disk type is not included in request if attribute is removed, prevents errors in NVME intances
			providerSettings.DiskTypeName = diskTypeName
		}
	}

	if providerName == "TENANT" {
		providerSettings.BackingProviderName = cast.ToString(d.Get("backing_provider_name"))
	}

	if autoScalingEnabled {
		providerSettings.AutoScaling = &matlas.AutoScaling{Compute: compute}
	}

	if d.Get("provider_name") == "AWS" {
		// Check if the Provider Disk IOS sets in the Terraform configuration and if the instance size name is not NVME.
		// If it didn't, the MongoDB Atlas server would set it to the default for the amount of storage.
		if v, ok := d.GetOk("provider_disk_iops"); ok && !strings.Contains(providerSettings.InstanceSizeName, "NVME") {
			providerSettings.DiskIOPS = conversion.Pointer(cast.ToInt64(v))
		}

		providerSettings.EncryptEBSVolume = conversion.Pointer(true)
	}

	return providerSettings, nil
}

func flattenProviderSettings(d *schema.ResourceData, settings *matlas.ProviderSettings, clusterName string) {
	if settings.ProviderName == "TENANT" {
		if err := d.Set("backing_provider_name", settings.BackingProviderName); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "backing_provider_name", clusterName, err)
		}
	}

	if settings.DiskIOPS != nil && *settings.DiskIOPS != 0 {
		if err := d.Set("provider_disk_iops", *settings.DiskIOPS); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "provider_disk_iops", clusterName, err)
		}
	}

	if err := d.Set("provider_disk_type_name", settings.DiskTypeName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_disk_type_name", clusterName, err)
	}

	if settings.EncryptEBSVolume != nil {
		if err := d.Set("provider_encrypt_ebs_volume_flag", *settings.EncryptEBSVolume); err != nil {
			log.Printf(advancedcluster.ErrorClusterSetting, "provider_encrypt_ebs_volume_flag", clusterName, err)
		}
	}

	if err := d.Set("provider_instance_size_name", settings.InstanceSizeName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_instance_size_name", clusterName, err)
	}

	if err := d.Set("provider_name", settings.ProviderName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_name", clusterName, err)
	}

	if err := d.Set("provider_region_name", settings.RegionName); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_region_name", clusterName, err)
	}

	if err := d.Set("provider_volume_type", settings.VolumeType); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "provider_volume_type", clusterName, err)
	}
}

func containsLabelOrKey(list []matlas.Label, item matlas.Label) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || v.Key == item.Key {
			return true
		}
	}

	return false
}
