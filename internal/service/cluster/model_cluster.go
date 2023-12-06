package cluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
)

const (
	ErrorClusterSetting            = "error setting `%s` for MongoDB Cluster (%s): %s"
	ErrorAdvancedConfRead          = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
	ErrorClusterAdvancedSetting    = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	ErrorAdvancedClusterListStatus = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
)

var (
	dsTagsSchema = schema.Schema{
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
	rsTagsSchema = schema.Schema{
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

func flattenProcessArgs(p *matlas.ProcessArgs) []any {
	return []any{
		map[string]any{
			"default_read_concern":                 p.DefaultReadConcern,
			"default_write_concern":                p.DefaultWriteConcern,
			"fail_index_key_too_long":              cast.ToBool(p.FailIndexKeyTooLong),
			"javascript_enabled":                   cast.ToBool(p.JavascriptEnabled),
			"minimum_enabled_tls_protocol":         p.MinimumEnabledTLSProtocol,
			"no_table_scan":                        cast.ToBool(p.NoTableScan),
			"oplog_size_mb":                        p.OplogSizeMB,
			"oplog_min_retention_hours":            p.OplogMinRetentionHours,
			"sample_size_bi_connector":             p.SampleSizeBIConnector,
			"sample_refresh_interval_bi_connector": p.SampleRefreshIntervalBIConnector,
			"transaction_lifetime_limit_seconds":   p.TransactionLifetimeLimitSeconds,
		},
	}
}

func expandProcessArgs(d *schema.ResourceData, p map[string]any) *matlas.ProcessArgs {
	res := &matlas.ProcessArgs{}

	if _, ok := d.GetOkExists("advanced_configuration.0.default_read_concern"); ok {
		res.DefaultReadConcern = cast.ToString(p["default_read_concern"])
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.default_write_concern"); ok {
		res.DefaultWriteConcern = cast.ToString(p["default_write_concern"])
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.fail_index_key_too_long"); ok {
		res.FailIndexKeyTooLong = pointy.Bool(cast.ToBool(p["fail_index_key_too_long"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.javascript_enabled"); ok {
		res.JavascriptEnabled = pointy.Bool(cast.ToBool(p["javascript_enabled"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.minimum_enabled_tls_protocol"); ok {
		res.MinimumEnabledTLSProtocol = cast.ToString(p["minimum_enabled_tls_protocol"])
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.no_table_scan"); ok {
		res.NoTableScan = pointy.Bool(cast.ToBool(p["no_table_scan"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.sample_size_bi_connector"); ok {
		res.SampleSizeBIConnector = pointy.Int64(cast.ToInt64(p["sample_size_bi_connector"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.sample_refresh_interval_bi_connector"); ok {
		res.SampleRefreshIntervalBIConnector = pointy.Int64(cast.ToInt64(p["sample_refresh_interval_bi_connector"]))
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.oplog_size_mb"); ok {
		if sizeMB := cast.ToInt64(p["oplog_size_mb"]); sizeMB != 0 {
			res.OplogSizeMB = pointy.Int64(cast.ToInt64(p["oplog_size_mb"]))
		} else {
			log.Printf(ErrorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
		}
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.oplog_min_retention_hours"); ok {
		if minRetentionHours := cast.ToFloat64(p["oplog_min_retention_hours"]); minRetentionHours >= 0 {
			res.OplogMinRetentionHours = pointy.Float64(cast.ToFloat64(p["oplog_min_retention_hours"]))
		} else {
			log.Printf(ErrorClusterSetting, `oplog_min_retention_hours`, "", cast.ToString(minRetentionHours))
		}
	}

	if _, ok := d.GetOkExists("advanced_configuration.0.transaction_lifetime_limit_seconds"); ok {
		if transactionLifetimeLimitSeconds := cast.ToInt64(p["transaction_lifetime_limit_seconds"]); transactionLifetimeLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = pointy.Int64(cast.ToInt64(p["transaction_lifetime_limit_seconds"]))
		} else {
			log.Printf(ErrorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLifetimeLimitSeconds))
		}
	}

	return res
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
