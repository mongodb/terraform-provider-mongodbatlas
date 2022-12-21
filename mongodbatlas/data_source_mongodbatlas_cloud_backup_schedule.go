package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: the schema is the same as dataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicy
// see documentation at https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/
func dataSourceMongoDBAtlasCloudBackupSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudBackupScheduleRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"copy_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequencies": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_spec_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_copy_oplogs": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"next_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference_hour_of_day": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"reference_minute_of_hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restore_window_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_export_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"use_org_and_group_names_in_export_prefix": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"export": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_bucket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_item_hourly": {
				Type:     schema.TypeList,
				Computed: true,
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
				Computed: true,
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
				Computed: true,
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
							Computed: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"policy_item_monthly": {
				Type:     schema.TypeList,
				Computed: true,
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
							Computed: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Almost the same as dataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead
// just do not save the update_snapshots because is not specified in the DS
func dataSourceMongoDBAtlasCloudBackupScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
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

	if err := d.Set("next_snapshot", backupPolicy.NextSnapshot); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "next_snapshot", clusterName, err)
	}
	if err := d.Set("use_org_and_group_names_in_export_prefix", backupPolicy.UseOrgAndGroupNamesInExportPrefix); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "use_org_and_group_names_in_export_prefix", clusterName, err)
	}
	if err := d.Set("auto_export_enabled", backupPolicy.AutoExportEnabled); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "auto_export_enabled", clusterName, err)
	}
	if err := d.Set("id_policy", backupPolicy.Policies[0].ID); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "id_policy", clusterName, err)
	}
	if err := d.Set("export", flattenExport(backupPolicy)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "auto_export_enabled", clusterName, err)
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

	if err := d.Set("copy_settings", flattenCopySettings(backupPolicy.CopySettings)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "copy_settings", clusterName, err)
	}

	if err := d.Set("export", flattenExport(backupPolicy)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "export", clusterName, err)
	}
	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}
