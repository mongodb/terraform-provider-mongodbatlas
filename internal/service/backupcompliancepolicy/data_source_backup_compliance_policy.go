package backupcompliancepolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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
						"retention_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retention_unit": {
							Type:     schema.TypeString,
							Computed: true,
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
			"policy_item_yearly": {
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	policy, resp, err := connV2.CloudBackupsApi.GetCompliancePolicy(ctx, projectID).Execute()
	if validate.StatusNotFound(resp) || policy.GetProjectId() == "" {
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf(cluster.ErrorSnapshotBackupPolicyRead, projectID, err))
	}

	if err := d.Set("authorized_email", policy.GetAuthorizedEmail()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_email", projectID, err))
	}

	if err := d.Set("authorized_user_first_name", policy.GetAuthorizedUserFirstName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_user_first_name", projectID, err))
	}

	if err := d.Set("authorized_user_last_name", policy.GetAuthorizedUserLastName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "authorized_user_last_name", projectID, err))
	}

	if err := d.Set("state", policy.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "state", projectID, err))
	}

	if err := d.Set("restore_window_days", policy.GetRestoreWindowDays()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "restore_window_days", projectID, err))
	}

	if err := d.Set("pit_enabled", policy.GetPitEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "pit_enabled", projectID, err))
	}

	if err := d.Set("copy_protection_enabled", policy.GetCopyProtectionEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "copy_protection_enabled", projectID, err))
	}

	if err := d.Set("encryption_at_rest_enabled", policy.GetEncryptionAtRestEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "encryption_at_rest_enabled", projectID, err))
	}

	if err := d.Set("updated_date", conversion.TimePtrToStringPtr(policy.UpdatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "updated_date", projectID, err))
	}

	if err := d.Set("updated_user", policy.GetUpdatedUser()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "updated_user", projectID, err))
	}

	if err := d.Set("on_demand_policy_item", flattenOnDemandBackupPolicyItem(policy.OnDemandPolicyItem)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSnapshotBackupPolicySetting, "on_demand_policy_item", projectID, err))
	}

	if err := d.Set("policy_item_hourly", flattenBackupPolicyItems(policy.GetScheduledPolicyItems(), cloudbackupschedule.Hourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_hourly", projectID, err)
	}

	if err := d.Set("policy_item_daily", flattenBackupPolicyItems(policy.GetScheduledPolicyItems(), cloudbackupschedule.Daily)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_daily", projectID, err)
	}

	if err := d.Set("policy_item_weekly", flattenBackupPolicyItems(policy.GetScheduledPolicyItems(), cloudbackupschedule.Weekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_weekly", projectID, err)
	}

	if err := d.Set("policy_item_monthly", flattenBackupPolicyItems(policy.GetScheduledPolicyItems(), cloudbackupschedule.Monthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_monthly", projectID, err)
	}

	if err := d.Set("policy_item_yearly", flattenBackupPolicyItems(policy.GetScheduledPolicyItems(), cloudbackupschedule.Yearly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_yearly", projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}
