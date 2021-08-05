package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorSnapshotBackupScheduleCreate  = "error creating a Cloud Backup Schedule: %s"
	errorSnapshotBackupScheduleUpdate  = "error updating a Cloud Backup Schedule: %s"
	errorSnapshotBackupScheduleRead    = "error getting a Cloud Backup Schedule for the cluster(%s): %s"
	errorSnapshotBackupScheduleSetting = "error setting `%s` for Cloud Backup Schedule(%s): %s"
	snapshotScheduleHourly             = "hourly"
	snapshotScheduleDaily              = "daily"
	snapshotScheduleWeekly             = "weekly"
	snapshotScheduleMonthly            = "monthly"
)

// https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/
// same as resourceMongoDBAtlasCloudProviderSnapshotBackupPolicy
func resourceMongoDBAtlasCloudBackupSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasCloudBackupScheduleCreate,
		ReadContext:   resourceMongoDBAtlasCloudBackupScheduleRead,
		UpdateContext: resourceMongoDBAtlasCloudBackupScheduleUpdate,
		DeleteContext: resourceMongoDBAtlasCloudBackupScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasCloudBackupScheduleImportState,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_item_hourly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"policy_item_daily": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"policy_item_weekly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"policy_item_monthly": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			// Optionals
			"reference_hour_of_day": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 23 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 23, got: %d", key, v))
					}
					return
				},
			},
			"reference_minute_of_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 59 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 59, got: %d", key, v))
					}
					return
				},
			},
			"restore_window_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"update_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Only computed
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasCloudBackupScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	err := cloudBackupScheduleCreateOrUpdate(ctx, conn, d, projectID, clusterName)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasCloudBackupScheduleRead(ctx, d, meta)
}

func resourceMongoDBAtlasCloudBackupScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if err := d.Set("cluster_id", backupPolicy.ClusterID); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "cluster_id", clusterName, err)
	}

	if err := d.Set("reference_hour_of_day", backupPolicy.ReferenceHourOfDay); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_hour_of_day", clusterName, err)
	}

	if err := d.Set("reference_minute_of_hour", backupPolicy.ReferenceMinuteOfHour); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_minute_of_hour", clusterName, err)
	}

	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "restore_window_days", clusterName, err)
	}

	if err := d.Set("update_snapshots", backupPolicy.UpdateSnapshots); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "update_snapshots", clusterName, err)
	}

	if err := d.Set("next_snapshot", backupPolicy.NextSnapshot); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "next_snapshot", clusterName, err)
	}

	if err := d.Set("id_policy", backupPolicy.Policies[0].ID); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "id_policy", clusterName, err)
	}

	if err := d.Set("policy_item_hourly", flattenPolicyItem(backupPolicy.Policies[0].PolicyItems, snapshotScheduleHourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_hourly", clusterName, err)
	}

	if err := d.Set("policy_item_daily", flattenPolicyItem(backupPolicy.Policies[0].PolicyItems, snapshotScheduleDaily)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_daily", clusterName, err)
	}

	if err := d.Set("policy_item_weekly", flattenPolicyItem(backupPolicy.Policies[0].PolicyItems, snapshotScheduleWeekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_weekly", clusterName, err)
	}

	if err := d.Set("policy_item_monthly", flattenPolicyItem(backupPolicy.Policies[0].PolicyItems, snapshotScheduleMonthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_monthly", clusterName, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudBackupScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if restoreWindowDays, ok := d.GetOk("restore_window_days"); ok {
		if cast.ToInt64(restoreWindowDays) <= 0 {
			return diag.Errorf("`restore_window_days` cannot be <= 0")
		}
	}

	err := cloudBackupScheduleCreateOrUpdate(ctx, conn, d, projectID, clusterName)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleUpdate, err)
	}

	return resourceMongoDBAtlasCloudBackupScheduleRead(ctx, d, meta)
}

func resourceMongoDBAtlasCloudBackupScheduleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Delete(ctx, projectID, clusterName)
	if err != nil {
		return diag.Errorf("error deleting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}

	return nil
}

func resourceMongoDBAtlasCloudBackupScheduleImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Cloud Backup Schedule use the format {project_id}-{cluster_name}")
	}

	projectID := parts[0]
	clusterName := parts[1]

	_, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "project_id", clusterName, err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "cluster_name", clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return []*schema.ResourceData{d}, nil
}

func cloudBackupScheduleCreateOrUpdate(ctx context.Context, conn *matlas.Client, d *schema.ResourceData, projectID, clusterName string) error {
	req := &matlas.CloudProviderSnapshotBackupPolicy{}

	// Delete policies items
	resp, _, err := conn.CloudProviderSnapshotBackupPolicies.Delete(ctx, projectID, clusterName)
	if err != nil {
		return fmt.Errorf("error deleting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}

	policy := matlas.Policy{}
	policyItem := matlas.PolicyItem{}
	var policiesItem []matlas.PolicyItem

	if v, ok := d.GetOk("policy_item_hourly"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		policyItem.FrequencyType = snapshotScheduleHourly
		policyItem.RetentionUnit = itemObj["retention_unit"].(string)
		policyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		policyItem.RetentionValue = itemObj["retention_value"].(int)
		policiesItem = append(policiesItem, policyItem)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		policyItem.FrequencyType = snapshotScheduleDaily
		policyItem.RetentionUnit = itemObj["retention_unit"].(string)
		policyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		policyItem.RetentionValue = itemObj["retention_value"].(int)
		policiesItem = append(policiesItem, policyItem)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		policyItem.FrequencyType = snapshotScheduleWeekly
		policyItem.RetentionUnit = itemObj["retention_unit"].(string)
		policyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		policyItem.RetentionValue = itemObj["retention_value"].(int)
		policiesItem = append(policiesItem, policyItem)
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		policyItem.FrequencyType = snapshotScheduleMonthly
		policyItem.RetentionUnit = itemObj["retention_unit"].(string)
		policyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		policyItem.RetentionValue = itemObj["retention_value"].(int)
		policiesItem = append(policiesItem, policyItem)
	}

	policy.ID = resp.Policies[0].ID
	policy.PolicyItems = policiesItem
	if len(policiesItem) > 0 {
		req.Policies = []matlas.Policy{policy}
	}

	req.ReferenceHourOfDay = pointy.Int64(cast.ToInt64(d.Get("reference_hour_of_day")))
	req.ReferenceMinuteOfHour = pointy.Int64(cast.ToInt64(d.Get("reference_minute_of_hour")))
	req.RestoreWindowDays = pointy.Int64(cast.ToInt64(d.Get("restore_window_days")))
	value := pointy.Bool(d.Get("update_snapshots").(bool))
	if *value {
		req.UpdateSnapshots = value
	}

	_, _, err = conn.CloudProviderSnapshotBackupPolicies.Update(context.Background(), projectID, clusterName, req)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupScheduleUpdate, err)
	}

	return nil
}

func flattenPolicyItem(items []matlas.PolicyItem, frequencyType string) []map[string]interface{} {
	policyItems := make([]map[string]interface{}, 0)
	for _, v := range items {
		if frequencyType == v.FrequencyType {
			policyItems = append(policyItems, map[string]interface{}{
				"id":                 v.ID,
				"frequency_interval": v.FrequencyInterval,
				"frequency_type":     v.FrequencyType,
				"retention_unit":     v.RetentionUnit,
				"retention_value":    v.RetentionValue,
			})
		}
	}

	return policyItems
}
