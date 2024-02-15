package advancedcluster

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	ErrorClusterSetting            = "error setting `%s` for MongoDB Cluster (%s): %s"
	ErrorAdvancedConfRead          = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
	ErrorClusterAdvancedSetting    = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	ErrorAdvancedClusterListStatus = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
)

var (
	defaultLabel = matlas.Label{Key: "Infrastructure Tool", Value: "MongoDB Atlas Terraform Provider"}
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

func RemoveLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
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

func ContainsLabelOrKey(list []matlas.Label, item matlas.Label) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || v.Key == item.Key {
			return true
		}
	}

	return false
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

func UpgradeCluster(ctx context.Context, conn *matlas.Client, request *matlas.Cluster, projectID, name string, timeout time.Duration) (*matlas.Cluster, *matlas.Response, error) {
	request.Name = name

	cluster, resp, err := conn.Clusters.Upgrade(ctx, projectID, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    ResourceClusterRefreshFunc(ctx, name, projectID, ServiceFromClient(conn)),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func ResourceClusterRefreshFunc(ctx context.Context, name, projectID string, client ClusterService) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		} else if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func ResourceClusterListAdvancedRefreshFunc(ctx context.Context, projectID string, client ClusterService) retry.StateRefreshFunc {
	return func() (any, string, error) {
		clusters, resp, err := client.List(ctx, projectID, nil)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && clusters == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		for i := range clusters.Results {
			if clusters.Results[i].StateName != "IDLE" {
				return clusters.Results[i], "PENDING", nil
			}
		}

		return clusters, "IDLE", nil
	}
}

func FormatMongoDBMajorVersion(val any) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}

	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
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

func flattenProcessArgs(p *admin.ClusterDescriptionProcessArgs) []map[string]any {
	if p == nil {
		return nil
	}
	return []map[string]any{
		{
			"default_read_concern":                 p.GetDefaultReadConcern(),
			"default_write_concern":                p.GetDefaultWriteConcern(),
			"fail_index_key_too_long":              p.GetFailIndexKeyTooLong(),
			"javascript_enabled":                   p.GetJavascriptEnabled(),
			"minimum_enabled_tls_protocol":         p.GetMinimumEnabledTlsProtocol(),
			"no_table_scan":                        p.GetNoTableScan(),
			"oplog_size_mb":                        p.GetOplogSizeMB(),
			"oplog_min_retention_hours":            p.GetOplogMinRetentionHours(),
			"sample_size_bi_connector":             p.GetSampleSizeBIConnector(),
			"sample_refresh_interval_bi_connector": p.GetSampleRefreshIntervalBIConnector(),
			"transaction_lifetime_limit_seconds":   p.GetTransactionLifetimeLimitSeconds(),
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

func SchemaAdvancedConfigDS() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_read_concern": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"default_write_concern": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"fail_index_key_too_long": {
					Type:     schema.TypeBool,
					Computed: true,
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
					Type:     schema.TypeInt,
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func SchemaConnectionStrings() *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeList,
		Computed:   true,
		ConfigMode: schema.SchemaConfigModeAttr,
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
					Type:       schema.TypeList,
					Computed:   true,
					ConfigMode: schema.SchemaConfigModeAttr,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"connection_string": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"endpoints": {
								Type:       schema.TypeList,
								Computed:   true,
								ConfigMode: schema.SchemaConfigModeAttr,
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
		Type:       schema.TypeList,
		Optional:   true,
		Computed:   true,
		ConfigMode: schema.SchemaConfigModeAttr,
		MaxItems:   1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_read_concern": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"default_write_concern": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"fail_index_key_too_long": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
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
					Type:     schema.TypeInt,
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
			},
		},
	}
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

func StringIsUppercase() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if value != strings.ToUpper(value) {
			diagError := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("The provided string '%q' must be uppercase.", value),
			}
			diags = append(diags, diagError)
		}
		return diags
	}
}
