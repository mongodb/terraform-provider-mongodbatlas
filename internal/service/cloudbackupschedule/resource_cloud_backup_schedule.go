package cloudbackupschedule

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	Hourly                             = "hourly"
	Daily                              = "daily"
	Weekly                             = "weekly"
	Monthly                            = "monthly"
	Yearly                             = "yearly"
	errorSnapshotBackupScheduleCreate  = "error creating a Cloud Backup Schedule: %s"
	errorSnapshotBackupScheduleUpdate  = "error updating a Cloud Backup Schedule: %s"
	errorSnapshotBackupScheduleRead    = "error getting a Cloud Backup Schedule for the cluster(%s): %s"
	errorSnapshotBackupScheduleSetting = "error setting `%s` for Cloud Backup Schedule(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
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
			"auto_export_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_org_and_group_names_in_export_prefix": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"copy_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"frequencies": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"should_copy_oplogs": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"export": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_bucket_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"frequency_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
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
			"reference_hour_of_day": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	// When a new cluster is created with the backup feature enabled,
	// MongoDB Atlas automatically generates a default backup policy for that cluster.
	// As a result, we need to first delete the default policies to avoid having
	// the infrastructure differs from the TF configuration file.
	if _, _, err := connV2.CloudBackupsApi.DeleteClusterBackupSchedule(ctx, projectID, clusterName).Execute(); err != nil {
		diagWarning := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error deleting default backup schedule",
			Detail:   fmt.Sprintf("error deleting default MongoDB Cloud Backup Schedule (%s): %s", clusterName, err),
		}
		diags = append(diags, diagWarning)
	}

	if err := cloudBackupScheduleCreateOrUpdate(ctx, connV2, d, projectID, clusterName); err != nil {
		diags = append(diags, diag.Errorf(errorSnapshotBackupScheduleCreate, err)...)
		return diags
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	var backupSchedule *admin.DiskBackupSnapshotSchedule20240805
	var resp *http.Response
	var err error

	backupSchedule, resp, err = connV2.CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusterName).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	diags := setSchemaFields(d, backupSchedule)
	if diags.HasError() {
		return diags
	}

	return nil
}

func setSchemaFields(d *schema.ResourceData, backupSchedule *admin.DiskBackupSnapshotSchedule20240805) diag.Diagnostics {
	clusterName := backupSchedule.GetClusterName()
	if err := d.Set("cluster_id", backupSchedule.GetClusterId()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "cluster_id", clusterName, err)
	}

	if err := d.Set("reference_hour_of_day", backupSchedule.GetReferenceHourOfDay()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_hour_of_day", clusterName, err)
	}

	if err := d.Set("reference_minute_of_hour", backupSchedule.GetReferenceMinuteOfHour()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_minute_of_hour", clusterName, err)
	}

	if err := d.Set("restore_window_days", backupSchedule.GetRestoreWindowDays()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "restore_window_days", clusterName, err)
	}

	if err := d.Set("next_snapshot", conversion.TimePtrToStringPtr(backupSchedule.NextSnapshot)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "next_snapshot", clusterName, err)
	}

	if err := d.Set("id_policy", backupSchedule.GetPolicies()[0].GetId()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "id_policy", clusterName, err)
	}

	if err := d.Set("use_org_and_group_names_in_export_prefix", backupSchedule.GetUseOrgAndGroupNamesInExportPrefix()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "use_org_and_group_names_in_export_prefix", clusterName, err)
	}

	if err := d.Set("policy_item_hourly", FlattenPolicyItem(backupSchedule.GetPolicies()[0].GetPolicyItems(), Hourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_hourly", clusterName, err)
	}

	if err := d.Set("policy_item_daily", FlattenPolicyItem(backupSchedule.GetPolicies()[0].GetPolicyItems(), Daily)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_daily", clusterName, err)
	}

	if err := d.Set("policy_item_weekly", FlattenPolicyItem(backupSchedule.GetPolicies()[0].GetPolicyItems(), Weekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_weekly", clusterName, err)
	}

	if err := d.Set("policy_item_monthly", FlattenPolicyItem(backupSchedule.GetPolicies()[0].GetPolicyItems(), Monthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_monthly", clusterName, err)
	}

	if err := d.Set("policy_item_yearly", FlattenPolicyItem(backupSchedule.GetPolicies()[0].GetPolicyItems(), Yearly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_yearly", clusterName, err)
	}

	if err := d.Set("copy_settings", FlattenCopySettings(backupSchedule.GetCopySettings())); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "copy_settings", clusterName, err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if restoreWindowDays, ok := d.GetOk("restore_window_days"); ok {
		if cast.ToInt64(restoreWindowDays) <= 0 {
			return diag.Errorf("`restore_window_days` cannot be <= 0")
		}
	}

	err := cloudBackupScheduleCreateOrUpdate(ctx, connV2, d, projectID, clusterName)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleUpdate, err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	_, _, err := connV2.CloudBackupsApi.DeleteClusterBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.Errorf("error deleting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}

	d.SetId("")

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Cloud Backup Schedule use the format {project_id}-{cluster_name}")
	}

	projectID := parts[0]
	clusterName := parts[1]

	_, _, err := connV2.CloudBackupsApi.GetBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "project_id", clusterName, err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		return nil, fmt.Errorf(errorSnapshotBackupScheduleSetting, "cluster_name", clusterName, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return []*schema.ResourceData{d}, nil
}

func cloudBackupScheduleCreateOrUpdate(ctx context.Context, connV2 *admin.APIClient, d *schema.ResourceData, projectID, clusterName string) error {
	var err error
	copySettings := d.Get("copy_settings")

	req := &admin.DiskBackupSnapshotSchedule20240805{}

	var policiesItem []admin.DiskBackupApiPolicyItem
	if v, ok := d.GetOk("policy_item_hourly"); ok {
		policiesItem = append(policiesItem, *ExpandPolicyItems(v.([]any), Hourly)...)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		policiesItem = append(policiesItem, *ExpandPolicyItems(v.([]any), Daily)...)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		policiesItem = append(policiesItem, *ExpandPolicyItems(v.([]any), Weekly)...)
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		policiesItem = append(policiesItem, *ExpandPolicyItems(v.([]any), Monthly)...)
	}
	if v, ok := d.GetOk("policy_item_yearly"); ok {
		policiesItem = append(policiesItem, *ExpandPolicyItems(v.([]any), Yearly)...)
	}

	if d.HasChange("auto_export_enabled") {
		req.AutoExportEnabled = conversion.Pointer(d.Get("auto_export_enabled").(bool))
	}

	if v, ok := d.GetOk("export"); ok {
		req.Export = expandAutoExportPolicy(v.([]any))
	}

	if d.HasChange("use_org_and_group_names_in_export_prefix") {
		req.UseOrgAndGroupNamesInExportPrefix = conversion.Pointer(d.Get("use_org_and_group_names_in_export_prefix").(bool))
	}

	if v, ok := d.GetOkExists("reference_hour_of_day"); ok {
		req.ReferenceHourOfDay = conversion.Pointer(v.(int))
	}
	if v, ok := d.GetOkExists("reference_minute_of_hour"); ok {
		req.ReferenceMinuteOfHour = conversion.Pointer(v.(int))
	}
	if v, ok := d.GetOkExists("restore_window_days"); ok {
		req.RestoreWindowDays = conversion.Pointer(v.(int))
	}

	value := conversion.Pointer(d.Get("update_snapshots").(bool))
	if *value {
		req.UpdateSnapshots = value
	}

	resp, _, err := connV2.CloudBackupsApi.GetBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf("error getting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}
	if isCopySettingsNonEmptyOrChanged(d) {
		req.CopySettings = ExpandCopySettings(copySettings.([]any))
	}

	req.Policies = getRequestPolicies(policiesItem, resp.GetPolicies())

	_, _, err = connV2.CloudBackupsApi.UpdateBackupSchedule(context.Background(), projectID, clusterName, req).Execute()
	if err != nil {
		return err
	}

	return nil
}

func ExpandCopySetting(tfMap map[string]any) *admin.DiskBackupCopySetting20240805 {
	if tfMap == nil {
		return nil
	}

	frequencies := conversion.ExpandStringList(tfMap["frequencies"].(*schema.Set).List())
	copySetting := &admin.DiskBackupCopySetting20240805{
		CloudProvider:    conversion.Pointer(tfMap["cloud_provider"].(string)),
		Frequencies:      &frequencies,
		RegionName:       conversion.Pointer(tfMap["region_name"].(string)),
		ZoneId:           tfMap["zone_id"].(string),
		ShouldCopyOplogs: conversion.Pointer(tfMap["should_copy_oplogs"].(bool)),
	}
	return copySetting
}

func ExpandCopySettings(tfList []any) *[]admin.DiskBackupCopySetting20240805 {
	copySettings := make([]admin.DiskBackupCopySetting20240805, 0)

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			continue
		}
		apiObject := ExpandCopySetting(tfMap)
		copySettings = append(copySettings, *apiObject)
	}
	return &copySettings
}

func expandAutoExportPolicy(items []any) *admin.AutoExportPolicy {
	itemObj := items[0].(map[string]any)
	return &admin.AutoExportPolicy{
		ExportBucketId: conversion.StringPtr(itemObj["export_bucket_id"].(string)),
		FrequencyType:  conversion.StringPtr(itemObj["frequency_type"].(string)),
	}
}

func ExpandPolicyItems(items []any, frequencyType string) *[]admin.DiskBackupApiPolicyItem {
	results := make([]admin.DiskBackupApiPolicyItem, len(items))

	for i, s := range items {
		itemObj := s.(map[string]any)
		results[i] = expandPolicyItem(itemObj, frequencyType)
	}
	return &results
}

func expandPolicyItem(itemObj map[string]any, frequencyType string) admin.DiskBackupApiPolicyItem {
	return admin.DiskBackupApiPolicyItem{
		Id:                policyItemID(itemObj),
		RetentionUnit:     itemObj["retention_unit"].(string),
		RetentionValue:    itemObj["retention_value"].(int),
		FrequencyInterval: itemObj["frequency_interval"].(int),
		FrequencyType:     frequencyType,
	}
}

func policyItemID(policyState map[string]any) *string {
	// if the policyItem ID is present then it's an update operation
	if val, ok := policyState["id"]; ok {
		if id, ok := val.(string); ok && id != "" {
			return &id
		}
	}
	return nil
}

func isCopySettingsNonEmptyOrChanged(d *schema.ResourceData) bool {
	copySettings := d.Get("copy_settings")
	return copySettings != nil && (conversion.HasElementsSliceOrMap(copySettings) || d.HasChange("copy_settings"))
}

func getRequestPolicies(policiesItem []admin.DiskBackupApiPolicyItem, respPolicies []admin.AdvancedDiskBackupSnapshotSchedulePolicy) *[]admin.AdvancedDiskBackupSnapshotSchedulePolicy {
	if len(policiesItem) > 0 {
		policy := admin.AdvancedDiskBackupSnapshotSchedulePolicy{
			PolicyItems: &policiesItem,
		}
		if len(respPolicies) == 1 {
			policy.Id = respPolicies[0].Id
		}
		return &[]admin.AdvancedDiskBackupSnapshotSchedulePolicy{policy}
	}
	return nil
}
