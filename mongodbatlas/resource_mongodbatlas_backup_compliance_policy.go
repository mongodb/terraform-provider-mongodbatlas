package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorBackupPolicyUpdate  = "error updating a Backup Compliance Policy: %s: %s"
	errorBackupPolicyRead    = "error getting a Backup Compliance Policy for the project(%s): %s"
	errorBackupPolicySetting = "error setting `%s` for Backup Compliance Policy : %s: %s"
)

func resourceMongoDBAtlasBackupCompliancePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasBackupCompliancePolicyCreate,
		UpdateContext: resourceMongoDBAtlasBackupCompliancePolicyUpdate,
		ReadContext:   resourceMongoDBAtlasBackupCompliancePolicyRead,
		DeleteContext: resourceMongoDBAtlasBackupCompliancePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasBackupCompliancePolicyImportState,
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
			"copy_protection_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"encryption_at_rest_enabled": {
				Type:     schema.TypeBool,
				Required: true,
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
				Required: true,
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

func resourceMongoDBAtlasBackupCompliancePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	backupPolicy := matlas.BackupCompliancePolicy{}
	backupPolicyItem := matlas.ScheduledPolicyItem{}
	var backupPoliciesItem []matlas.ScheduledPolicyItem

	backupCompliancePolicyReq := &matlas.BackupCompliancePolicy{}

	backupCompliancePolicyReq.ProjectID = projectID

	backupCompliancePolicyReq.AuthorizedEmail = d.Get("authorized_email").(string)

	backupCompliancePolicyReq.CopyProtectionEnabled = pointy.Bool(d.Get("copy_protection_enabled").(bool))

	backupCompliancePolicyReq.EncryptionAtRestEnabled = pointy.Bool(d.Get("encryption_at_rest_enabled").(bool))

	backupCompliancePolicyReq.PitEnabled = pointy.Bool(d.Get("pit_enabled").(bool))

	backupCompliancePolicyReq.RestoreWindowDays = pointy.Int64(cast.ToInt64(d.Get("restore_window_days")))

	backupCompliancePolicyReq.OnDemandPolicyItem = *expandDemandBackupPolicyItem(d)

	if v, ok := d.GetOk("policy_item_hourly"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		backupPolicyItem.FrequencyType = snapshotScheduleHourly
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		backupPolicyItem.FrequencyType = snapshotScheduleDaily
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		items := v.([]interface{})
		for _, s := range items {
			itemObj := s.(map[string]interface{})
			backupPolicyItem.FrequencyType = snapshotScheduleWeekly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		items := v.([]interface{})
		for _, s := range items {
			itemObj := s.(map[string]interface{})
			backupPolicyItem.FrequencyType = snapshotScheduleMonthly
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

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return resourceMongoDBAtlasBackupCompliancePolicyRead(ctx, d, meta)
}

func resourceMongoDBAtlasBackupCompliancePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]

	backupPolicy, resp, err := conn.BackupCompliancePolicy.Get(context.Background(), projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorBackupPolicyRead, projectID, err))
	}

	if err := d.Set("authorized_email", backupPolicy.AuthorizedEmail); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "authorized_email", projectID, err))
	}

	if err := d.Set("restore_window_days", backupPolicy.RestoreWindowDays); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "restore_window_days", projectID, err))
	}

	if err := d.Set("updated_date", backupPolicy.UpdatedDate); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "updated_date", projectID, err))
	}

	if err := d.Set("updated_user", backupPolicy.UpdatedUser); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "updated_user", projectID, err))
	}

	if err := d.Set("state", backupPolicy.State); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "state", projectID, err))
	}

	if err := d.Set("copy_protection_enabled", backupPolicy.CopyProtectionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "copy_protection_enabled", projectID, err))
	}

	if err := d.Set("encryption_at_rest_enabled", backupPolicy.EncryptionAtRestEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "encryption_at_rest_enabled", projectID, err))
	}

	if err := d.Set("pit_enabled", backupPolicy.PitEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "pit_enabled", projectID, err))
	}

	if err := d.Set("on_demand_policy_item", flattenOnDemandBackupPolicyItem(backupPolicy.OnDemandPolicyItem)); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "scheduled_policy_items", projectID, err))
	}

	if err := d.Set("policy_item_hourly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, snapshotScheduleHourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_hourly", projectID, err)
	}

	if err := d.Set("policy_item_daily", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, snapshotScheduleDaily)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_daily", projectID, err)
	}

	if err := d.Set("policy_item_weekly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, snapshotScheduleWeekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_weekly", projectID, err)
	}

	if err := d.Set("policy_item_monthly", flattenBackupPolicyItems(backupPolicy.ScheduledPolicyItems, snapshotScheduleMonthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupPolicySetting, "policy_item_monthly", projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}

func resourceMongoDBAtlasBackupCompliancePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]

	backupPolicy := matlas.BackupCompliancePolicy{}
	backupPolicyItem := matlas.ScheduledPolicyItem{}
	var backupPoliciesItem []matlas.ScheduledPolicyItem

	backupCompliancePolicyUpdate := &matlas.BackupCompliancePolicy{}

	backupCompliancePolicyUpdate.ProjectID = projectID

	backupCompliancePolicyUpdate.AuthorizedEmail = d.Get("authorized_email").(string)

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
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		backupPolicyItem.FrequencyType = snapshotScheduleHourly
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]interface{})
		itemObj := item[0].(map[string]interface{})
		backupPolicyItem.FrequencyType = snapshotScheduleDaily
		backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
		backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
		backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
		backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		items := v.([]interface{})
		for _, s := range items {
			itemObj := s.(map[string]interface{})
			backupPolicyItem.FrequencyType = snapshotScheduleWeekly
			backupPolicyItem.RetentionUnit = itemObj["retention_unit"].(string)
			backupPolicyItem.FrequencyInterval = itemObj["frequency_interval"].(int)
			backupPolicyItem.RetentionValue = itemObj["retention_value"].(int)
			backupPoliciesItem = append(backupPoliciesItem, backupPolicyItem)
		}
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		items := v.([]interface{})
		for _, s := range items {
			itemObj := s.(map[string]interface{})
			backupPolicyItem.FrequencyType = snapshotScheduleMonthly
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

	return resourceMongoDBAtlasBackupCompliancePolicyRead(ctx, d, meta)
}

func resourceMongoDBAtlasBackupCompliancePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no resource to delete a backup compliance policy, it can only be updated.
	log.Printf("[WARN] Note: Deleting a Backup Compliance Policy resource in Terraform does not remove the policy from your Atlas Project. " +
		"To disable a Backup Compliance Policy, the security or legal representative specified for the Backup Compliance Policy must contact " +
		"MongoDB Support and complete an extensive verification process. ")

	d.SetId("")
	return nil
}

func resourceMongoDBAtlasBackupCompliancePolicyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 1 {
		return nil, errors.New("import format error: to import a Backup Compliance Policy use the format {project_id}")
	}

	projectID := parts[0]

	_, _, err := conn.BackupCompliancePolicy.Get(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(errorBackupPolicyRead, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorBackupPolicySetting, "project_id", projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenOnDemandBackupPolicyItem(item matlas.PolicyItem) []map[string]interface{} {
	policyItems := make([]map[string]interface{}, 0)

	policyItems = append(policyItems, map[string]interface{}{
		"id":                 item.ID,
		"frequency_interval": item.FrequencyInterval,
		"frequency_type":     item.FrequencyType,
		"retention_unit":     item.RetentionUnit,
		"retention_value":    item.RetentionValue,
	})

	return policyItems
}

func expandDemandBackupPolicyItem(d *schema.ResourceData) *matlas.PolicyItem {
	var onDemand matlas.PolicyItem

	if v, ok := d.GetOk("on_demand_policy_item"); ok {
		demandItem := v.([]interface{})
		if len(demandItem) > 0 {
			demandItemMap := demandItem[0].(map[string]interface{})

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
