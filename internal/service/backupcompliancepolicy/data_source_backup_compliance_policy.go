package backupcompliancepolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasBackupCompliancePolicyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_user_first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_user_last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"copy_protection_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_at_rest_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"on_demand_policy_item": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
			"pit_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restore_window_days": {
				Type:     schema.TypeInt,
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
		},
	}
}

func dataSourceMongoDBAtlasBackupCompliancePolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	backupPolicy, resp, err := conn.BackupCompliancePolicy.Get(ctx, projectID)
	if resp != nil && resp.StatusCode == http.StatusNotFound || backupPolicy.ProjectID == "" {
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf(cluster.ErrorSnapshotBackupPolicyRead, projectID, err))
	}

	if err := d.Set("authorized_email", backupPolicy.AuthorizedEmail); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_email", projectID, err))
	}

	if err := d.Set("authorized_user_first_name", backupPolicy.AuthorizedUserFirstName); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_user_first_name", projectID, err))
	}

	if err := d.Set("authorized_user_last_name", backupPolicy.AuthorizedUserLastName); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_user_last_name", projectID, err))
	}

	if err := d.Set("state", backupPolicy.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "state", projectID, err))
	}

	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "restore_window_days", projectID, err))
	}

	if err := d.Set("pit_enabled", backupPolicy.PitEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "pit_enabled", projectID, err))
	}

	if err := d.Set("copy_protection_enabled", backupPolicy.CopyProtectionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "copy_protection_enabled", projectID, err))
	}

	if err := d.Set("encryption_at_rest_enabled", backupPolicy.EncryptionAtRestEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "encryption_at_rest_enabled", projectID, err))
	}

	if err := d.Set("updated_date", backupPolicy.UpdatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "updated_date", projectID, err))
	}

	if err := d.Set("updated_user", backupPolicy.UpdatedUser); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "updated_user", projectID, err))
	}

	if err := d.Set("on_demand_policy_item", flattenOnDemandBackupPolicyItem(backupPolicy.OnDemandPolicyItem)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "policies", projectID, err))
	}

	if err := d.Set("policy_item_hourly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, cloudbackupschedule.Hourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_hourly", projectID, err)
	}

	if err := d.Set("policy_item_daily", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, cloudbackupschedule.Daily)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_daily", projectID, err)
	}

	if err := d.Set("policy_item_weekly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, cloudbackupschedule.Weekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_weekly", projectID, err)
	}

	if err := d.Set("policy_item_monthly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, cloudbackupschedule.Monthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_monthly", projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}

func flattenBackupPolicyItems(items []matlas.ScheduledPolicyItem, frequencyType string) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, v := range items {
		if frequencyType == v.FrequencyType {
			policyItems = append(policyItems, map[string]any{
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
