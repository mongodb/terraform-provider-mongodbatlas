package backupcompliancepolicy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorBackupPolicyUpdate          = "error updating a Backup Compliance Policy: %s: %s"
	errorBackupPolicyRead            = "error getting a Backup Compliance Policy for the project(%s): %s"
	errorBackupPolicySetting         = "error setting `%s` for Backup Compliance Policy : %s: %s"
	errorSnapshotBackupPolicySetting = "error setting `%s` for Cloud Provider Snapshot Backup Policy(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		UpdateContext: resourceUpdate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_user_first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorized_user_last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"copy_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encryption_at_rest_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"restore_window_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"on_demand_policy_item": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Computed: true,
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
			"pit_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	backupPolicy := matlas.BackupCompliancePolicy{}
	backupPolicyItem := matlas.ScheduledPolicyItem{}
	var backupPoliciesItem []matlas.ScheduledPolicyItem

	backupCompliancePolicyReq := &matlas.BackupCompliancePolicy{}

	backupCompliancePolicyReq.ProjectID = projectID

	backupCompliancePolicyReq.AuthorizedEmail = d.Get("authorized_email").(string)

	backupCompliancePolicyReq.AuthorizedUserFirstName = d.Get("authorized_user_first_name").(string)

	backupCompliancePolicyReq.AuthorizedUserLastName = d.Get("authorized_user_last_name").(string)

	backupCompliancePolicyReq.CopyProtectionEnabled = pointy.Bool(d.Get("copy_protection_enabled").(bool))

	backupCompliancePolicyReq.EncryptionAtRestEnabled = pointy.Bool(d.Get("encryption_at_rest_enabled").(bool))

	backupCompliancePolicyReq.PitEnabled = pointy.Bool(d.Get("pit_enabled").(bool))

	backupCompliancePolicyReq.RestoreWindowDays = pointy.Int64(cast.ToInt64(d.Get("restore_window_days")))

	backupCompliancePolicyReq.OnDemandPolicyItem = *expandDemandBackupPolicyItem(d)

	if v, ok := d.GetOk("policy_item_hourly"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPolicyItem.FrequencyType = cloudbackupschedule.Hourly
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPolicyItem.FrequencyType = cloudbackupschedule.Daily
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPolicyItem.FrequencyType = cloudbackupschedule.Weekly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPolicyItem.FrequencyType = cloudbackupschedule.Monthly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}

	backupPolicy.ScheduledPolicyItems = backupPoliciesItem
	if len(backupPoliciesItem) > 0 {
		backupCompliancePolicyReq.ScheduledPolicyItems = backupPoliciesItem
	}

	// there is not an entry point to create a backup compliance policy until it will use the update entry point
	_, _, err := conn.BackupCompliancePolicy.Update(ctx, projectID, backupCompliancePolicyReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicyUpdate, projectID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]

	policy, resp, err := connV2.CloudBackupsApi.GetDataProtectionSettings(context.Background(), projectID).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorBackupPolicyRead, projectID, err))
	}

	if err := d.Set("authorized_email", policy.GetAuthorizedEmail()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "authorized_email", projectID, err))
	}

	if err := d.Set("authorized_user_first_name", policy.GetAuthorizedUserFirstName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "authorized_user_first_name", projectID, err))
	}

	if err := d.Set("authorized_user_last_name", policy.GetAuthorizedUserLastName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "authorized_user_last_name", projectID, err))
	}

	if err := d.Set("restore_window_days", policy.GetRestoreWindowDays()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "restore_window_days", projectID, err))
	}

	if err := d.Set("updated_date", conversion.TimePtrToStringPtr(policy.UpdatedDate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "updated_date", projectID, err))
	}

	if err := d.Set("updated_user", policy.GetUpdatedUser()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "updated_user", projectID, err))
	}

	if err := d.Set("state", policy.GetState()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "state", projectID, err))
	}

	if err := d.Set("copy_protection_enabled", policy.GetCopyProtectionEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "copy_protection_enabled", projectID, err))
	}

	if err := d.Set("encryption_at_rest_enabled", policy.GetEncryptionAtRestEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "encryption_at_rest_enabled", projectID, err))
	}

	if err := d.Set("pit_enabled", policy.GetPitEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "pit_enabled", projectID, err))
	}

	if err := d.Set("on_demand_policy_item", flattenOnDemandBackupPolicyItem(policy.GetOnDemandPolicyItem())); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "scheduled_policy_items", projectID, err))
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

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]

	backupPolicy := matlas.BackupCompliancePolicy{}
	backupPolicyItem := matlas.ScheduledPolicyItem{}
	var backupPoliciesItem []matlas.ScheduledPolicyItem

	backupCompliancePolicyUpdate := &matlas.BackupCompliancePolicy{}

	backupCompliancePolicyUpdate.ProjectID = projectID

	backupCompliancePolicyUpdate.AuthorizedEmail = d.Get("authorized_email").(string)

	backupCompliancePolicyUpdate.AuthorizedUserFirstName = d.Get("authorized_user_first_name").(string)

	backupCompliancePolicyUpdate.AuthorizedUserLastName = d.Get("authorized_user_last_name").(string)

	if d.HasChange("copy_protection_enabled") {
		backupCompliancePolicyUpdate.CopyProtectionEnabled = pointy.Bool(d.Get("copy_protection_enabled").(bool))
	}

	if d.HasChange("encryption_at_rest_enabled") {
		backupCompliancePolicyUpdate.CopyProtectionEnabled = pointy.Bool(d.Get("copy_protection_enabled").(bool))
	}

	if d.HasChange("pit_enabled") {
		backupCompliancePolicyUpdate.PitEnabled = pointy.Bool(d.Get("pit_enabled").(bool))
	}

	if d.HasChange("restore_window_days") {
		backupCompliancePolicyUpdate.RestoreWindowDays = pointy.Int64(cast.ToInt64(d.Get("restore_window_days")))
	}

	backupCompliancePolicyUpdate.OnDemandPolicyItem = *expandDemandBackupPolicyItem(d)

	if v, ok := d.GetOk("policy_item_hourly"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPolicyItem.FrequencyType = cloudbackupschedule.Hourly
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPolicyItem.FrequencyType = cloudbackupschedule.Daily
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPolicyItem.FrequencyType = cloudbackupschedule.Weekly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPolicyItem.FrequencyType = cloudbackupschedule.Monthly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}

	backupPolicy.ScheduledPolicyItems = backupPoliciesItem
	if len(backupPoliciesItem) > 0 {
		backupCompliancePolicyUpdate.ScheduledPolicyItems = backupPoliciesItem
	}

	_, _, err := conn.BackupCompliancePolicy.Update(context.Background(), projectID, backupCompliancePolicyUpdate)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicyUpdate, projectID, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// There is no resource to delete a backup compliance policy, it can only be updated.
	log.Printf("[WARN] Note: Deleting a Backup Compliance Policy resource in Terraform does not remove the policy from your Atlas Project. " +
		"To disable a Backup Compliance Policy, the security or legal representative specified for the Backup Compliance Policy must contact " +
		"MongoDB Support and complete an extensive verification process. ")

	d.SetId("")
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 1 {
		return nil, errors.New("import format error: to import a Backup Compliance Policy use the format {project_id}")
	}
	projectID := parts[0]

	_, _, err := connV2.CloudBackupsApi.GetDataProtectionSettings(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorBackupPolicyRead, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorBackupPolicySetting, "project_id", projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenOnDemandBackupPolicyItem(item admin.BackupComplianceOnDemandPolicyItem) []map[string]any {
	return []map[string]any{
		{
			"id":                 item.GetId(),
			"frequency_interval": item.GetFrequencyInterval(),
			"frequency_type":     item.GetFrequencyType(),
			"retention_unit":     item.GetRetentionUnit(),
			"retention_value":    item.GetRetentionValue(),
		},
	}
}

func expandDemandBackupPolicyItem(d *schema.ResourceData) *matlas.PolicyItem {
	var onDemand matlas.PolicyItem

	if v, ok := d.GetOk("on_demand_policy_item"); ok {
		demandItem := v.([]any)
		if len(demandItem) > 0 {
			demandItemMap := demandItem[0].(map[string]any)

			onDemand = matlas.PolicyItem{
				ID:                demandItemMap["id"].(string),
				FrequencyInterval: demandItemMap["frequency_interval"].(int),
				FrequencyType:     "ondemand",
				RetentionUnit:     demandItemMap["retention_unit"].(string),
				RetentionValue:    demandItemMap["retention_value"].(int),
			}
		}
	}

	return &onDemand
}

func flattenBackupPolicyItems(items []admin.BackupComplianceScheduledPolicyItem, frequencyType string) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for i := range items {
		item := &items[i]
		if frequencyType == item.FrequencyType {
			policyItems = append(policyItems, map[string]any{
				"id":                 item.GetId(),
				"frequency_interval": item.GetFrequencyInterval(),
				"frequency_type":     item.GetFrequencyType(),
				"retention_unit":     item.GetRetentionUnit(),
				"retention_value":    item.GetRetentionValue(),
			})
		}
	}
	return policyItems
}
