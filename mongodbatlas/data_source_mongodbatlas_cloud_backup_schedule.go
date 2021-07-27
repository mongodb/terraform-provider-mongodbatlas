package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: the schema is the same as dataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicy
// see documentation at https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/
func dataSourceMongoDBAtlasCloudBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasCloudBackupScheduleRead,
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
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_item": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"frequency_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"frequency_type": {
										Type:     schema.TypeString,
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
				},
			},
		},
	}
}

// Almost the same as dataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicyRead
// just do not save the update_snapshots because is not specified in the DS
func dataSourceMongoDBAtlasCloudBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	backupPolicy, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
	if err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	if err := d.Set("cluster_id", backupPolicy.ClusterID); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "cluster_id", clusterName, err)
	}

	if err := d.Set("reference_hour_of_day", backupPolicy.ReferenceHourOfDay); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "reference_hour_of_day", clusterName, err)
	}

	if err := d.Set("reference_minute_of_hour", backupPolicy.ReferenceMinuteOfHour); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "reference_minute_of_hour", clusterName, err)
	}

	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "restore_window_days", clusterName, err)
	}

	if err := d.Set("next_snapshot", backupPolicy.NextSnapshot); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "next_snapshot", clusterName, err)
	}

	if err := d.Set("policies", flattenPolicies(backupPolicy.Policies)); err != nil {
		return fmt.Errorf(errorSnapshotBackupPolicySetting, "policies", clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}
