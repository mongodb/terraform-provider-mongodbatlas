package backupcompliancepolicy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
)

const (
	errorBackupPolicyUpdate          = "error updating a Backup Compliance Policy: %s: %s"
	errorBackupPolicyDelete          = "error disabling the Backup Compliance Policy: %s: %s"
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
				Optional: true,
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
			"policy_item_yearly": {
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	err := updateOrCreateDataProtectionSetting(ctx, d, connV2, projectID)

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

	policy, resp, err := connV2.CloudBackupsApi.GetCompliancePolicy(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
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

	if err := d.Set("on_demand_policy_item", flattenOnDemandBackupPolicyItem(policy.OnDemandPolicyItem)); err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicySetting, "on_demand_policy_item", projectID, err))
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

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]

	err := updateOrCreateDataProtectionSetting(ctx, d, connV2, projectID)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicyUpdate, projectID, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	_, err := connV2.CloudBackupsApi.DisableCompliancePolicy(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorBackupPolicyDelete, projectID, err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 1 {
		return nil, errors.New("import format error: to import a Backup Compliance Policy use the format {project_id}")
	}
	projectID := parts[0]

	_, _, err := connV2.CloudBackupsApi.GetCompliancePolicy(ctx, projectID).Execute()
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

func flattenOnDemandBackupPolicyItem(item *admin.BackupComplianceOnDemandPolicyItem) []map[string]any {
	if item == nil {
		return nil
	}
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

func expandDemandBackupPolicyItem(d *schema.ResourceData) *admin.BackupComplianceOnDemandPolicyItem {
	if v, ok := d.GetOk("on_demand_policy_item"); ok {
		demandItem := v.([]any)
		if len(demandItem) > 0 {
			demandItemMap := demandItem[0].(map[string]any)
			return &admin.BackupComplianceOnDemandPolicyItem{
				Id:                conversion.StringPtr(demandItemMap["id"].(string)),
				FrequencyInterval: demandItemMap["frequency_interval"].(int),
				FrequencyType:     "ondemand",
				RetentionUnit:     demandItemMap["retention_unit"].(string),
				RetentionValue:    demandItemMap["retention_value"].(int),
			}
		}
	}
	return nil
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

func updateOrCreateDataProtectionSetting(ctx context.Context, d *schema.ResourceData, connV2 *admin.APIClient, projectID string) error {
	dataProtectionSettings := &admin.DataProtectionSettings20231001{
		ProjectId:               conversion.StringPtr(projectID),
		AuthorizedEmail:         d.Get("authorized_email").(string),
		AuthorizedUserFirstName: d.Get("authorized_user_first_name").(string),
		AuthorizedUserLastName:  d.Get("authorized_user_last_name").(string),
		CopyProtectionEnabled:   conversion.Pointer(d.Get("copy_protection_enabled").(bool)),
		EncryptionAtRestEnabled: conversion.Pointer(d.Get("encryption_at_rest_enabled").(bool)),
		PitEnabled:              conversion.Pointer(d.Get("pit_enabled").(bool)),
		RestoreWindowDays:       conversion.Pointer(cast.ToInt(d.Get("restore_window_days"))),
		OnDemandPolicyItem:      expandDemandBackupPolicyItem(d),
	}

	var backupPoliciesItem []admin.BackupComplianceScheduledPolicyItem
	if v, ok := d.GetOk("policy_item_hourly"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPoliciesItem = append(backupPoliciesItem, admin.BackupComplianceScheduledPolicyItem{
			FrequencyType:     cloudbackupschedule.Hourly,
			RetentionUnit:     itemObj["retention_unit"].(string),
			FrequencyInterval: itemObj["frequency_interval"].(int),
			RetentionValue:    itemObj["retention_value"].(int),
		})
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		item := v.([]any)
		itemObj := item[0].(map[string]any)
		backupPoliciesItem = append(backupPoliciesItem, admin.BackupComplianceScheduledPolicyItem{
			FrequencyType:     cloudbackupschedule.Daily,
			RetentionUnit:     itemObj["retention_unit"].(string),
			FrequencyInterval: itemObj["frequency_interval"].(int),
			RetentionValue:    itemObj["retention_value"].(int),
		})
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPoliciesItem = append(backupPoliciesItem, admin.BackupComplianceScheduledPolicyItem{
				FrequencyType:     cloudbackupschedule.Weekly,
				RetentionUnit:     itemObj["retention_unit"].(string),
				FrequencyInterval: itemObj["frequency_interval"].(int),
				RetentionValue:    itemObj["retention_value"].(int),
			})
		}
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPoliciesItem = append(backupPoliciesItem, admin.BackupComplianceScheduledPolicyItem{
				FrequencyType:     cloudbackupschedule.Monthly,
				RetentionUnit:     itemObj["retention_unit"].(string),
				FrequencyInterval: itemObj["frequency_interval"].(int),
				RetentionValue:    itemObj["retention_value"].(int),
			})
		}
	}
	if v, ok := d.GetOk("policy_item_yearly"); ok {
		items := v.([]any)
		for _, s := range items {
			itemObj := s.(map[string]any)
			backupPoliciesItem = append(backupPoliciesItem, admin.BackupComplianceScheduledPolicyItem{
				FrequencyType:     cloudbackupschedule.Yearly,
				RetentionUnit:     itemObj["retention_unit"].(string),
				FrequencyInterval: itemObj["frequency_interval"].(int),
				RetentionValue:    itemObj["retention_value"].(int),
			})
		}
	}
	if len(backupPoliciesItem) > 0 {
		dataProtectionSettings.ScheduledPolicyItems = &backupPoliciesItem
	}

	params := admin.UpdateCompliancePolicyApiParams{
		GroupId:                        projectID,
		DataProtectionSettings20231001: dataProtectionSettings,
		OverwriteBackupPolicies:        conversion.Pointer(false),
	}
	_, _, err := connV2.CloudBackupsApi.UpdateCompliancePolicyWithParams(ctx, &params).Execute()
	return err
}
