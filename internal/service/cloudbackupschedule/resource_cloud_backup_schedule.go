package cloudbackupschedule

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
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
	DeprecationOldSchemaAction         = "Please refer to our examples, documentation, and 1.18.0 migration guide for more details at https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide.html.markdown"
)

var DeprecationMsgOldSchema = fmt.Sprintf("%s %s", fmt.Sprintf("%s %s ", constant.DeprecationParamWithReplacement, "`copy_settings.#.zone_id`"), DeprecationOldSchemaAction)

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
				Computed: true,
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
						"replication_spec_id": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: DeprecationMsgOldSchema,
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
				Computed: true,
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
	connV220231115 := meta.(*config.MongoDBClient).AtlasV220231115
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	// When a new cluster is created with the backup feature enabled,
	// MongoDB Atlas automatically generates a default backup policy for that cluster.
	// As a result, we need to first delete the default policies to avoid having
	// the infrastructure differs from the TF configuration file.
	if _, _, err := connV2.CloudBackupsApi.DeleteAllBackupSchedules(ctx, projectID, clusterName).Execute(); err != nil {
		diagWarning := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error deleting default backup schedule",
			Detail:   fmt.Sprintf("error deleting default MongoDB Cloud Backup Schedule (%s): %s", clusterName, err),
		}
		diags = append(diags, diagWarning)
	}

	if err := cloudBackupScheduleCreateOrUpdate(ctx, connV220231115, connV2, d, projectID, clusterName, true); err != nil {
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
	connV220231115 := meta.(*config.MongoDBClient).AtlasV220231115
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	var backupSchedule *admin.DiskBackupSnapshotSchedule20250101
	var backupScheduleOldSDK *admin20231115.DiskBackupSnapshotSchedule
	var copySettings []map[string]any
	var resp *http.Response
	var err error

	useOldAPI, err := shouldUseOldAPI(d, false)
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	if useOldAPI {
		backupScheduleOldSDK, resp, err = connV220231115.CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusterName).Execute()
		// TODO handle asymmetric error, ask user to set zone_id for asymmetric cluster ????
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
			return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
		}

		copySettings = flattenCopySettingsOldSDK(backupScheduleOldSDK.GetCopySettings())
		backupSchedule = convertBackupScheduleToLatestExcludeCopySettings(backupScheduleOldSDK)
	} else {
		backupSchedule, resp, err = connV2.CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusterName).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
			return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
		}
		copySettings = flattenCopySettings(backupSchedule.GetCopySettings())
	}

	diags := setSchemaFieldsExceptCopySettings(d, backupSchedule)
	if diags.HasError() {
		return diags
	}

	if err := d.Set("copy_settings", copySettings); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "copy_settings", clusterName, err)
	}

	return nil
}

func setSchemaFieldsExceptCopySettings(d *schema.ResourceData, backupPolicy *admin.DiskBackupSnapshotSchedule20250101) diag.Diagnostics {
	clusterName := backupPolicy.GetClusterName()
	if err := d.Set("cluster_id", backupPolicy.GetClusterId()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "cluster_id", clusterName, err)
	}

	if err := d.Set("reference_hour_of_day", backupPolicy.GetReferenceHourOfDay()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_hour_of_day", clusterName, err)
	}

	if err := d.Set("reference_minute_of_hour", backupPolicy.GetReferenceMinuteOfHour()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "reference_minute_of_hour", clusterName, err)
	}

	if err := d.Set("restore_window_days", backupPolicy.GetRestoreWindowDays()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "restore_window_days", clusterName, err)
	}

	if err := d.Set("next_snapshot", conversion.TimePtrToStringPtr(backupPolicy.NextSnapshot)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "next_snapshot", clusterName, err)
	}

	if err := d.Set("id_policy", backupPolicy.GetPolicies()[0].GetId()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "id_policy", clusterName, err)
	}

	if err := d.Set("export", flattenExport(backupPolicy)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "export", clusterName, err)
	}

	if err := d.Set("auto_export_enabled", backupPolicy.GetAutoExportEnabled()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "auto_export_enabled", clusterName, err)
	}

	if err := d.Set("use_org_and_group_names_in_export_prefix", backupPolicy.GetUseOrgAndGroupNamesInExportPrefix()); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "use_org_and_group_names_in_export_prefix", clusterName, err)
	}

	if err := d.Set("policy_item_hourly", flattenPolicyItem(backupPolicy.GetPolicies()[0].GetPolicyItems(), Hourly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_hourly", clusterName, err)
	}

	if err := d.Set("policy_item_daily", flattenPolicyItem(backupPolicy.GetPolicies()[0].GetPolicyItems(), Daily)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_daily", clusterName, err)
	}

	if err := d.Set("policy_item_weekly", flattenPolicyItem(backupPolicy.GetPolicies()[0].GetPolicyItems(), Weekly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_weekly", clusterName, err)
	}

	if err := d.Set("policy_item_monthly", flattenPolicyItem(backupPolicy.GetPolicies()[0].GetPolicyItems(), Monthly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_monthly", clusterName, err)
	}

	if err := d.Set("policy_item_yearly", flattenPolicyItem(backupPolicy.GetPolicies()[0].GetPolicyItems(), Yearly)); err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleSetting, "policy_item_yearly", clusterName, err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV220231115 := meta.(*config.MongoDBClient).AtlasV220231115
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	if restoreWindowDays, ok := d.GetOk("restore_window_days"); ok {
		if cast.ToInt64(restoreWindowDays) <= 0 {
			return diag.Errorf("`restore_window_days` cannot be <= 0")
		}
	}

	err := cloudBackupScheduleCreateOrUpdate(ctx, connV220231115, connV2, d, projectID, clusterName, false)
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

	_, _, err := connV2.CloudBackupsApi.DeleteAllBackupSchedules(ctx, projectID, clusterName).Execute()
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

func cloudBackupScheduleCreateOrUpdate(ctx context.Context, connV220231115 *admin20231115.APIClient, connV2 *admin.APIClient, d *schema.ResourceData, projectID, clusterName string, isCreate bool) error {
	var err error
	copySettings := d.Get("copy_settings")

	useOldAPI, err := shouldUseOldAPI(d, isCreate)
	if err != nil {
		return err
	}

	req := &admin.DiskBackupSnapshotSchedule20250101{}

	var policiesItem []admin.DiskBackupApiPolicyItem
	if v, ok := d.GetOk("policy_item_hourly"); ok {
		policiesItem = append(policiesItem, *expandPolicyItems(v.([]any), Hourly)...)
	}
	if v, ok := d.GetOk("policy_item_daily"); ok {
		policiesItem = append(policiesItem, *expandPolicyItems(v.([]any), Daily)...)
	}
	if v, ok := d.GetOk("policy_item_weekly"); ok {
		policiesItem = append(policiesItem, *expandPolicyItems(v.([]any), Weekly)...)
	}
	if v, ok := d.GetOk("policy_item_monthly"); ok {
		policiesItem = append(policiesItem, *expandPolicyItems(v.([]any), Monthly)...)
	}
	if v, ok := d.GetOk("policy_item_yearly"); ok {
		policiesItem = append(policiesItem, *expandPolicyItems(v.([]any), Yearly)...)
	}

	if d.HasChange("auto_export_enabled") {
		req.AutoExportEnabled = conversion.Pointer(d.Get("auto_export_enabled").(bool))
	}

	if v, ok := d.GetOk("export"); ok {
		req.Export = expandAutoExportPolicy(v.([]any), d)
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

	if useOldAPI {
		resp, _, err := connV220231115.CloudBackupsApi.GetBackupSchedule(ctx, projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("error getting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
			// TODO handle asymmetric error, ask user to set zone_id for asymmetric cluster
		}
		var copySettingsOldSDK *[]admin20231115.DiskBackupCopySetting
		if isCopySettingsNonEmptyOrChanged(d) {
			copySettingsOldSDK = expandCopySettingsOldSDK(copySettings.([]any))
		}

		policiesOldSDK := getRequestPoliciesOldSDK(convertPolicyItemsToOldSDK(&policiesItem), resp.GetPolicies())

		reqOld := convertBackupScheduleReqToOldSDK(req, copySettingsOldSDK, policiesOldSDK)
		_, _, err = connV220231115.CloudBackupsApi.UpdateBackupSchedule(context.Background(), projectID, clusterName, reqOld).Execute()
		if err != nil {
			return err
			// TODO handle asymmetric error, ask user to set zone_id for asymmetric cluster
		}

		return nil
	}

	resp, _, err := connV2.CloudBackupsApi.GetBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf("error getting MongoDB Cloud Backup Schedule (%s): %s", clusterName, err)
	}
	if isCopySettingsNonEmptyOrChanged(d) {
		req.CopySettings = expandCopySettings(copySettings.([]any))
	}

	req.Policies = getRequestPolicies(policiesItem, resp.GetPolicies())

	_, _, err = connV2.CloudBackupsApi.UpdateBackupSchedule(context.Background(), projectID, clusterName, req).Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenPolicyItem(items []admin.DiskBackupApiPolicyItem, frequencyType string) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, v := range items {
		if frequencyType == v.GetFrequencyType() {
			policyItems = append(policyItems, map[string]any{
				"id":                 v.GetId(),
				"frequency_interval": v.GetFrequencyInterval(),
				"frequency_type":     v.GetFrequencyType(),
				"retention_unit":     v.GetRetentionUnit(),
				"retention_value":    v.GetRetentionValue(),
			})
		}
	}
	return policyItems
}

func flattenExport(roles *admin.DiskBackupSnapshotSchedule20250101) []map[string]any {
	exportList := make([]map[string]any, 0)
	emptyStruct := admin.DiskBackupSnapshotSchedule20250101{}
	if emptyStruct.GetExport() != roles.GetExport() {
		exportList = append(exportList, map[string]any{
			"frequency_type":   roles.Export.GetFrequencyType(),
			"export_bucket_id": roles.Export.GetExportBucketId(),
		})
	}
	return exportList
}

func flattenCopySettingsOldSDK(copySettingList []admin20231115.DiskBackupCopySetting) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		copySettings = append(copySettings, map[string]any{
			"cloud_provider":      v.GetCloudProvider(),
			"frequencies":         v.GetFrequencies(),
			"region_name":         v.GetRegionName(),
			"replication_spec_id": v.GetReplicationSpecId(),
			"should_copy_oplogs":  v.GetShouldCopyOplogs(),
		})
	}
	return copySettings
}

func flattenCopySettings(copySettingList []admin.DiskBackupCopySetting20250101) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		copySettings = append(copySettings, map[string]any{
			"cloud_provider":      v.GetCloudProvider(),
			"frequencies":         v.GetFrequencies(),
			"region_name":         v.GetRegionName(),
			"zone_id":             v.GetZoneId(),
			"replication_spec_id": nil,
			"should_copy_oplogs":  v.GetShouldCopyOplogs(),
		})
	}
	return copySettings
}

func expandCopySetting(tfMap map[string]any) *admin.DiskBackupCopySetting20250101 {
	if tfMap == nil {
		return nil
	}

	frequencies := conversion.ExpandStringList(tfMap["frequencies"].(*schema.Set).List())
	copySetting := &admin.DiskBackupCopySetting20250101{
		CloudProvider:    conversion.Pointer(tfMap["cloud_provider"].(string)),
		Frequencies:      &frequencies,
		RegionName:       conversion.Pointer(tfMap["region_name"].(string)),
		ZoneId:           tfMap["zone_id"].(string),
		ShouldCopyOplogs: conversion.Pointer(tfMap["should_copy_oplogs"].(bool)),
	}
	return copySetting
}

func expandCopySettings(tfList []any) *[]admin.DiskBackupCopySetting20250101 {
	copySettings := make([]admin.DiskBackupCopySetting20250101, 0)

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			continue
		}
		apiObject := expandCopySetting(tfMap)
		copySettings = append(copySettings, *apiObject)
	}
	return &copySettings
}

func expandCopySettingsOldSDK(tfList []any) *[]admin20231115.DiskBackupCopySetting {
	copySettings := make([]admin20231115.DiskBackupCopySetting, 0)

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			continue
		}
		apiObject := expandCopySettingOldSDK(tfMap)
		copySettings = append(copySettings, *apiObject)
	}
	return &copySettings
}

func expandCopySettingOldSDK(tfMap map[string]any) *admin20231115.DiskBackupCopySetting {
	if tfMap == nil {
		return nil
	}

	frequencies := conversion.ExpandStringList(tfMap["frequencies"].(*schema.Set).List())
	copySetting := &admin20231115.DiskBackupCopySetting{
		CloudProvider:     conversion.Pointer(tfMap["cloud_provider"].(string)),
		Frequencies:       &frequencies,
		RegionName:        conversion.Pointer(tfMap["region_name"].(string)),
		ReplicationSpecId: conversion.Pointer(tfMap["replication_spec_id"].(string)),
		ShouldCopyOplogs:  conversion.Pointer(tfMap["should_copy_oplogs"].(bool)),
	}
	return copySetting
}

func expandPolicyItems(items []any, frequencyType string) *[]admin.DiskBackupApiPolicyItem {
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

func shouldUseOldAPI(d *schema.ResourceData, isCreate bool) (bool, error) {
	copySettings := d.Get("copy_settings")
	if isCopySettingsNonEmptyOrChanged(d) {
		return checkCopySettings(copySettings.([]any), isCreate)
	}
	return false, nil
}

func isCopySettingsNonEmptyOrChanged(d *schema.ResourceData) bool {
	copySettings := d.Get("copy_settings")
	return copySettings != nil && (conversion.HasElementsSliceOrMap(copySettings) || d.HasChange("copy_settings"))
}

// checkCopySettings verifies that all elements in tfList use either `replication_spec_id` or `zone_id`
// Returns an error if any element has both `replication_spec_id` and `zone_id` set during create
func checkCopySettings(tfList []any, isCreate bool) (bool, error) {
	allHaveRepID := true

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			return false, fmt.Errorf("element is not a valid map[string]any")
		}

		repSpecID, repOk := tfMap["replication_spec_id"].(string)
		zoneID, zoneOk := tfMap["zone_id"].(string)

		if repOk && repSpecID != "" && zoneOk && zoneID != "" {
			if isCreate {
				return false, fmt.Errorf("both 'replication_spec_id' and 'zone_id' cannot be set unless resource is being updated")
			} else {
				return false, nil
			}
		}

		if (repOk && repSpecID != "" && zoneOk && zoneID != "") || (!repOk && !zoneOk) {
			return false, fmt.Errorf("each element must have either 'replication_spec_id' or 'zone_id' set, but not both")
		}

		if !repOk || repSpecID == "" {
			allHaveRepID = false
		}
	}

	if allHaveRepID {
		return true, nil
	} else {
		return false, nil
	}
}
