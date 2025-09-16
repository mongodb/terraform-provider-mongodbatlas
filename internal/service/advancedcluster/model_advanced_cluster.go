package advancedcluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"slices"
	"strconv"
	"strings"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

const minVersionForChangeStreamOptions = 6.0
const minVersionForDefaultMaxTimeMS = 8.0

type OldShardConfigMeta struct {
	ID       string
	NumShard int
}

var (
	DSTagsSchema = schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
	RSTagsSchema = schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
)

func SchemaAdvancedConfigDS() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_read_concern": {
					Type:       schema.TypeString,
					Computed:   true,
					Deprecated: DeprecationMsgOldSchema,
				},
				"default_write_concern": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"fail_index_key_too_long": {
					Type:       schema.TypeBool,
					Computed:   true,
					Deprecated: DeprecationMsgOldSchema,
				},
				"javascript_enabled": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"minimum_enabled_tls_protocol": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"no_table_scan": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"oplog_size_mb": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"sample_size_bi_connector": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"sample_refresh_interval_bi_connector": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"oplog_min_retention_hours": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"change_stream_options_pre_and_post_images_expire_after_seconds": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"default_max_time_ms": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"tls_cipher_config_mode": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"custom_openssl_cipher_config_tls12": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func SchemaConnectionStrings() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"standard": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"standard_srv": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"private": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"private_srv": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"private_endpoint": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"connection_string": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"endpoints": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"endpoint_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"provider_name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"region": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
							},
							"srv_connection_string": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"srv_shard_optimized_connection_string": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"type": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func SchemaAdvancedConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_read_concern": {
					Type:       schema.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: DeprecationMsgOldSchema,
				},
				"default_write_concern": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"fail_index_key_too_long": {
					Type:       schema.TypeBool,
					Optional:   true,
					Computed:   true,
					Deprecated: DeprecationMsgOldSchema,
				},
				"javascript_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"minimum_enabled_tls_protocol": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"no_table_scan": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"oplog_size_mb": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"oplog_min_retention_hours": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"sample_size_bi_connector": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"sample_refresh_interval_bi_connector": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"change_stream_options_pre_and_post_images_expire_after_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  -1,
				},
				"default_max_time_ms": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"custom_openssl_cipher_config_tls12": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"tls_cipher_config_mode": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func HashFunctionForKeyValuePair(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)
	buf.WriteString(m["key"].(string))
	buf.WriteString(m["value"].(string))
	return HashCodeString(buf.String())
}

// HashCodeString hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func HashCodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func IsSharedTier(instanceSize string) bool {
	return instanceSize == "M0" || instanceSize == "M2" || instanceSize == "M5"
}

// GetDiskSizeGBFromReplicationSpec obtains the diskSizeGB value by looking into the electable spec of the first replication spec.
// Independent storage size scaling is not supported (CLOUDP-201331), meaning all electable/analytics/readOnly configs in all replication specs are the same.
func GetDiskSizeGBFromReplicationSpec(cluster *admin.ClusterDescription20240805) float64 {
	specs := cluster.GetReplicationSpecs()
	if len(specs) < 1 {
		return 0
	}
	configs := specs[0].GetRegionConfigs()
	if len(configs) < 1 {
		return 0
	}
	return configs[0].ElectableSpecs.GetDiskSizeGB()
}

func WaitStateTransitionClusterUpgrade(ctx context.Context, name, projectID string,
	client admin.ClustersApi, pendingStates, desiredStates []string, timeout time.Duration) (*admin.ClusterDescription20240805, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     desiredStates,
		Refresh:    advancedclustertpf.ResourceRefreshFunc(ctx, name, projectID, client),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	if cluster, ok := result.(*admin.ClusterDescription20240805); ok && cluster != nil {
		return cluster, nil
	}

	return nil, errors.New("did not obtain valid result when waiting for cluster upgrade state transition")
}

func ResourceClusterListAdvancedRefreshFunc(ctx context.Context, projectID string, clustersAPI admin.ClustersApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		clusters, resp, err := clustersAPI.ListClusters(ctx, projectID).Execute()

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && clusters == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}
			if validate.StatusServiceUnavailable(resp) {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		for i := range clusters.GetResults() {
			cluster := clusters.GetResults()[i]
			if cluster.GetStateName() != "IDLE" {
				return cluster, "PENDING", nil
			}
		}
		return clusters, "IDLE", nil
	}
}

func FormatMongoDBMajorVersion(val any) string {
	return advancedclustertpf.FormatMongoDBMajorVersion(val.(string))
}

// CheckRegionConfigsPriorityOrder will be deleted in CLOUDP-275825
func CheckRegionConfigsPriorityOrder(regionConfigs []admin.ReplicationSpec20240805) error {
	for _, spec := range regionConfigs {
		configs := spec.GetRegionConfigs()
		for i := range len(configs) - 1 {
			if configs[i].GetPriority() < configs[i+1].GetPriority() {
				return errors.New("priority values in region_configs must be in descending order")
			}
		}
	}
	return nil
}

// CheckRegionConfigsPriorityOrderOld will be deleted in CLOUDP-275825
func CheckRegionConfigsPriorityOrderOld(regionConfigs []admin20240530.ReplicationSpec) error {
	for _, spec := range regionConfigs {
		configs := spec.GetRegionConfigs()
		for i := range len(configs) - 1 {
			if configs[i].GetPriority() < configs[i+1].GetPriority() {
				return errors.New("priority values in region_configs must be in descending order")
			}
		}
	}
	return nil
}

func FlattenPinnedFCV(cluster *admin.ClusterDescription20240805) []map[string]string {
	if cluster.FeatureCompatibilityVersionExpirationDate == nil { // pinned_fcv is defined in state only if featureCompatibilityVersionExpirationDate is present in cluster response
		return nil
	}
	nestedObj := map[string]string{}
	nestedObj["version"] = cluster.GetFeatureCompatibilityVersion()
	nestedObj["expiration_date"] = conversion.TimeToString(cluster.GetFeatureCompatibilityVersionExpirationDate())
	return []map[string]string{nestedObj}
}

func FlattenAdvancedReplicationSpecsOldShardingConfig(ctx context.Context, apiObjects []admin.ReplicationSpec20240805, zoneNameToOldReplicationSpecMeta map[string]OldShardConfigMeta, tfMapObjects []any,
	d *schema.ResourceData, connV2 *admin.APIClient) ([]map[string]any, error) {
	replicationSpecFlattener := func(ctx context.Context, sdkModel *admin.ReplicationSpec20240805, tfModel map[string]any, resourceData *schema.ResourceData, client *admin.APIClient) (map[string]any, error) {
		return flattenAdvancedReplicationSpecOldShardingConfig(ctx, sdkModel, zoneNameToOldReplicationSpecMeta, tfModel, resourceData, connV2)
	}
	compressedAPIObjects := compressAPIObjectList(apiObjects)
	return flattenAdvancedReplicationSpecsLogic(ctx, compressedAPIObjects, tfMapObjects, d,
		doesAdvancedReplicationSpecMatchAPIOldShardConfig, replicationSpecFlattener, connV2)
}

// compressAPIObjectList returns an array of ReplicationSpec20240805. The input array is reduced from all shards to only one shard per zoneName
func compressAPIObjectList(apiObjects []admin.ReplicationSpec20240805) []admin.ReplicationSpec20240805 {
	var compressedAPIObjectList []admin.ReplicationSpec20240805
	wasZoneNameUsed := populateZoneNameMap(apiObjects)
	for _, apiObject := range apiObjects {
		if !wasZoneNameUsed[apiObject.GetZoneName()] {
			compressedAPIObjectList = append(compressedAPIObjectList, apiObject)
			wasZoneNameUsed[apiObject.GetZoneName()] = true
		}
	}
	return compressedAPIObjectList
}

// populateZoneNameMap returns a map of zoneNames and initializes all keys to false.
func populateZoneNameMap(apiObjects []admin.ReplicationSpec20240805) map[string]bool {
	zoneNames := make(map[string]bool)
	for _, apiObject := range apiObjects {
		if _, exists := zoneNames[apiObject.GetZoneName()]; !exists {
			zoneNames[apiObject.GetZoneName()] = false
		}
	}
	return zoneNames
}

type ReplicationSpecSDKModel interface {
	admin20240530.ReplicationSpec | admin.ReplicationSpec20240805
}

func flattenAdvancedReplicationSpecsLogic[T ReplicationSpecSDKModel](
	ctx context.Context, apiObjects []T, tfMapObjects []any, d *schema.ResourceData,
	tfModelWithSDKMatcher func(map[string]any, *T) bool,
	flattenRepSpec func(context.Context, *T, map[string]any, *schema.ResourceData, *admin.APIClient) (map[string]any, error),
	connV2 *admin.APIClient) ([]map[string]any, error) {
	if len(apiObjects) == 0 {
		return nil, nil
	}

	tfList := make([]map[string]any, len(apiObjects))
	wasAPIObjectUsed := make([]bool, len(apiObjects))

	for i := range len(tfList) {
		var tfMapObject map[string]any

		if len(tfMapObjects) > i {
			tfMapObject = tfMapObjects[i].(map[string]any)
		}
		for j := 0; j < len(apiObjects); j++ {
			if wasAPIObjectUsed[j] || !tfModelWithSDKMatcher(tfMapObject, &apiObjects[j]) {
				continue
			}

			advancedReplicationSpec, err := flattenRepSpec(ctx, &apiObjects[j], tfMapObject, d, connV2)

			if err != nil {
				return nil, err
			}

			tfList[i] = advancedReplicationSpec
			wasAPIObjectUsed[j] = true
			break
		}
	}

	for i, tfo := range tfList {
		var tfMapObject map[string]any

		if tfo != nil {
			continue
		}

		if len(tfMapObjects) > i {
			tfMapObject = tfMapObjects[i].(map[string]any)
		}

		j := slices.IndexFunc(wasAPIObjectUsed, func(isUsed bool) bool { return !isUsed })
		advancedReplicationSpec, err := flattenRepSpec(ctx, &apiObjects[j], tfMapObject, d, connV2)

		if err != nil {
			return nil, err
		}

		tfList[i] = advancedReplicationSpec
		wasAPIObjectUsed[j] = true
	}

	return tfList, nil
}

func doesAdvancedReplicationSpecMatchAPIOldShardConfig(tfObject map[string]any, apiObject *admin.ReplicationSpec20240805) bool {
	return tfObject["zone_name"] == apiObject.GetZoneName()
}

func flattenAdvancedReplicationSpecRegionConfigs(ctx context.Context, apiObjects []admin.CloudRegionConfig20240805, tfMapObjects []any,
	d *schema.ResourceData, connV2 *admin.APIClient) (tfResult []map[string]any, containersIDs map[string]string, err error) {
	if len(apiObjects) == 0 {
		return nil, nil, nil
	}

	var tfList []map[string]any
	containerIDs := make(map[string]string)

	for i := range apiObjects {
		apiObject := apiObjects[i]
		if len(tfMapObjects) > i {
			tfMapObject := tfMapObjects[i].(map[string]any)
			tfList = append(tfList, flattenAdvancedReplicationSpecRegionConfig(&apiObject, tfMapObject))
		} else {
			tfList = append(tfList, flattenAdvancedReplicationSpecRegionConfig(&apiObject, nil))
		}

		if apiObject.GetProviderName() != "TENANT" {
			params := &admin.ListGroupContainersApiParams{
				GroupId:      d.Get("project_id").(string),
				ProviderName: apiObject.ProviderName,
			}
			containers, _, err := connV2.NetworkPeeringApi.ListGroupContainersWithParams(ctx, params).Execute()
			if err != nil {
				return nil, nil, err
			}
			if result := advancedclustertpf.GetAdvancedClusterContainerID(containers.GetResults(), &apiObject); result != "" {
				// Will print as "providerName:regionName" = "containerId" in terraform show
				containerIDs[fmt.Sprintf("%s:%s", apiObject.GetProviderName(), apiObject.GetRegionName())] = result
			}
		}
	}
	return tfList, containerIDs, nil
}

func flattenAdvancedReplicationSpecRegionConfig(apiObject *admin.CloudRegionConfig20240805, tfMapObject map[string]any) map[string]any {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]any{}
	if tfMapObject != nil {
		if v, ok := tfMapObject["analytics_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["analytics_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.AnalyticsSpecs, apiObject.GetProviderName(), tfMapObject["analytics_specs"].([]any))
		}
		if v, ok := tfMapObject["electable_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["electable_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(hwSpecToDedicatedHwSpec(apiObject.ElectableSpecs), apiObject.GetProviderName(), tfMapObject["electable_specs"].([]any))
		}
		if v, ok := tfMapObject["read_only_specs"]; ok && len(v.([]any)) > 0 {
			tfMap["read_only_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ReadOnlySpecs, apiObject.GetProviderName(), tfMapObject["read_only_specs"].([]any))
		}
		if v, ok := tfMapObject["auto_scaling"]; ok && len(v.([]any)) > 0 {
			tfMap["auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AutoScaling)
		}
		if v, ok := tfMapObject["analytics_auto_scaling"]; ok && len(v.([]any)) > 0 {
			tfMap["analytics_auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AnalyticsAutoScaling)
		}
	} else {
		tfMap["analytics_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.AnalyticsSpecs, apiObject.GetProviderName(), nil)
		tfMap["electable_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(hwSpecToDedicatedHwSpec(apiObject.ElectableSpecs), apiObject.GetProviderName(), nil)
		tfMap["read_only_specs"] = flattenAdvancedReplicationSpecRegionConfigSpec(apiObject.ReadOnlySpecs, apiObject.GetProviderName(), nil)
		tfMap["auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AutoScaling)
		tfMap["analytics_auto_scaling"] = flattenAdvancedReplicationSpecAutoScaling(apiObject.AnalyticsAutoScaling)
	}

	tfMap["region_name"] = apiObject.GetRegionName()
	tfMap["provider_name"] = apiObject.GetProviderName()
	tfMap["backing_provider_name"] = apiObject.GetBackingProviderName()
	tfMap["priority"] = apiObject.GetPriority()

	return tfMap
}

func hwSpecToDedicatedHwSpec(apiObject *admin.HardwareSpec20240805) *admin.DedicatedHardwareSpec20240805 {
	if apiObject == nil {
		return nil
	}
	return &admin.DedicatedHardwareSpec20240805{
		NodeCount:     apiObject.NodeCount,
		DiskIOPS:      apiObject.DiskIOPS,
		EbsVolumeType: apiObject.EbsVolumeType,
		InstanceSize:  apiObject.InstanceSize,
		DiskSizeGB:    apiObject.DiskSizeGB,
	}
}

func flattenAdvancedReplicationSpecRegionConfigSpec(apiObject *admin.DedicatedHardwareSpec20240805, providerName string, tfMapObjects []any) []map[string]any {
	if apiObject == nil {
		return nil
	}
	var tfList []map[string]any

	tfMap := map[string]any{}

	if len(tfMapObjects) > 0 {
		tfMapObject := tfMapObjects[0].(map[string]any)

		if providerName == constant.AWS || providerName == constant.AZURE {
			if cast.ToInt64(apiObject.GetDiskIOPS()) > 0 {
				tfMap["disk_iops"] = apiObject.GetDiskIOPS()
			}
		}
		if providerName == constant.AWS {
			if v, ok := tfMapObject["ebs_volume_type"]; ok && v.(string) != "" {
				tfMap["ebs_volume_type"] = apiObject.GetEbsVolumeType()
			}
		}
		if _, ok := tfMapObject["disk_size_gb"]; ok {
			tfMap["disk_size_gb"] = apiObject.GetDiskSizeGB()
		}
		if _, ok := tfMapObject["node_count"]; ok {
			tfMap["node_count"] = apiObject.GetNodeCount()
		}
		if v, ok := tfMapObject["instance_size"]; ok && v.(string) != "" {
			tfMap["instance_size"] = apiObject.GetInstanceSize()
			tfList = append(tfList, tfMap)
		}
	} else {
		tfMap["disk_size_gb"] = apiObject.GetDiskSizeGB()
		tfMap["disk_iops"] = apiObject.GetDiskIOPS()
		tfMap["ebs_volume_type"] = apiObject.GetEbsVolumeType()
		tfMap["node_count"] = apiObject.GetNodeCount()
		tfMap["instance_size"] = apiObject.GetInstanceSize()
		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenAdvancedReplicationSpecAutoScaling(apiObject *admin.AdvancedAutoScalingSettings) []map[string]any {
	if apiObject == nil {
		return nil
	}
	var tfList []map[string]any
	tfMap := map[string]any{}
	if apiObject.DiskGB != nil {
		tfMap["disk_gb_enabled"] = apiObject.DiskGB.GetEnabled()
	}
	if apiObject.Compute != nil {
		tfMap["compute_enabled"] = apiObject.Compute.GetEnabled()
		tfMap["compute_scale_down_enabled"] = apiObject.Compute.GetScaleDownEnabled()
		tfMap["compute_min_instance_size"] = apiObject.Compute.GetMinInstanceSize()
		tfMap["compute_max_instance_size"] = apiObject.Compute.GetMaxInstanceSize()
	}
	tfList = append(tfList, tfMap)
	return tfList
}

func isMinRequiredMajorVersion(input *string, minVersion float64) bool {
	if input == nil || *input == "" {
		return true
	}
	parts := strings.SplitN(*input, ".", 2)
	if len(parts) == 0 {
		return false
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return false
	}

	return value >= minVersion
}

func IsChangeStreamOptionsMinRequiredMajorVersion(input *string) bool {
	return isMinRequiredMajorVersion(input, minVersionForChangeStreamOptions)
}

func IsDefaultMaxTimeMinRequiredMajorVersion(input *string) bool {
	return isMinRequiredMajorVersion(input, minVersionForDefaultMaxTimeMS)
}

func flattenAdvancedReplicationSpecOldShardingConfig(ctx context.Context, apiObject *admin.ReplicationSpec20240805, zoneNameToOldShardConfigMeta map[string]OldShardConfigMeta, tfMapObject map[string]any,
	d *schema.ResourceData, connV2 *admin.APIClient) (map[string]any, error) {
	if apiObject == nil {
		return nil, nil
	}

	tfMap := map[string]any{}
	if oldShardConfigData, ok := zoneNameToOldShardConfigMeta[apiObject.GetZoneName()]; ok {
		tfMap["num_shards"] = oldShardConfigData.NumShard
		tfMap["id"] = oldShardConfigData.ID
	}
	if tfMapObject != nil {
		object, containerIDs, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.GetRegionConfigs(), tfMapObject["region_configs"].([]any), d, connV2)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIDs
	} else {
		object, containerIDs, err := flattenAdvancedReplicationSpecRegionConfigs(ctx, apiObject.GetRegionConfigs(), nil, d, connV2)
		if err != nil {
			return nil, err
		}
		tfMap["region_configs"] = object
		tfMap["container_id"] = containerIDs
	}
	tfMap["zone_name"] = apiObject.GetZoneName()
	tfMap["zone_id"] = apiObject.GetZoneId()

	return tfMap, nil
}
