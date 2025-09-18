package advancedcluster

import (
	"bytes"
	"context"
	"hash/crc32"
	"strconv"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

const minVersionForChangeStreamOptions = 6.0
const V20240530 = "(v20240530)"

const (
	ErrorAdvancedConfRead          = "error reading Advanced Configuration Option %s for MongoDB Cluster (%s): %s"
	ErrorClusterAdvancedSetting    = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	ErrorFlexClusterSetting        = "error setting `%s` for MongoDB Flex Cluster (%s): %s"
	ErrorAdvancedClusterListStatus = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
	ErrorDefaultMaxTimeMinVersion  = "`advanced_configuration.default_max_time_ms` can only be configured if the mongo_db_major_version is 8.0 or higher"
)

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

func WarningIfFCVExpiredOrUnpinnedExternally(d *schema.ResourceData, cluster *admin.ClusterDescription20240805) diag.Diagnostics {
	pinnedFCVBlock, _ := d.Get("pinned_fcv").([]any)
	fcvPresentInState := len(pinnedFCVBlock) > 0
	diagsTpf := advancedclustertpf.GenerateFCVPinningWarningForRead(fcvPresentInState, cluster.FeatureCompatibilityVersionExpirationDate)
	return conversion.FromTPFDiagsToSDKV2Diags(diagsTpf)
}

// GetDiskSizeGBFromReplicationSpec obtains the diskSizeGB value by looking into the electable spec of the first replication spec.
// Independent storage size scaling is not supported (CLOUDP-201331), meaning all electable/analytics/readOnly configs in all replication specs are the same.

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

func FlattenPinnedFCV(cluster *admin.ClusterDescription20240805) []map[string]string {
	if cluster.FeatureCompatibilityVersionExpirationDate == nil { // pinned_fcv is defined in state only if featureCompatibilityVersionExpirationDate is present in cluster response
		return nil
	}
	nestedObj := map[string]string{}
	nestedObj["version"] = cluster.GetFeatureCompatibilityVersion()
	nestedObj["expiration_date"] = conversion.TimeToString(cluster.GetFeatureCompatibilityVersionExpirationDate())
	return []map[string]string{nestedObj}
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
