package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasBackupCompliancePolicy() *schema.Resource {
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
			"scheduled_policy_items": {
				Type:     schema.TypeList,
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
	}
}

func dataSourceMongoDBAtlasBackupCompliancePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	backupPolicy, resp, err := conn.BackupCompliancePolicy.Get(ctx, projectID)
	if resp != nil && resp.StatusCode == http.StatusNotFound || backupPolicy.ProjectID == "" {
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicyRead, projectID, err))
	}

	if err := d.Set("authorized_email", backupPolicy.AuthorizedEmail); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_email", projectID, err))
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

	if err := d.Set("scheduled_policy_items", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "policies", projectID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}

func flattenBackupPolicyItems(items []matlas.ScheduledPolicyItem) []map[string]interface{} {
	policyItems := make([]map[string]interface{}, 0)
	for _, v := range items {
		policyItems = append(policyItems, map[string]interface{}{
			"id":                 v.ID,
			"frequency_interval": v.FrequencyInterval,
			"frequency_type":     v.FrequencyType,
			"retention_unit":     v.RetentionUnit,
			"retention_value":    v.RetentionValue,
		})
	}

	return policyItems
}
