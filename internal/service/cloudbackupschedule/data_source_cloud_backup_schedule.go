package cloudbackupschedule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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
			"copy_policy_items_enabled": {
				Type:     schema.TypeBool,
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_copy_oplogs": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"copy_policy_items": copyPolicyItemsSchema(true),
						"last_number_of_snapshots": {
							Type:     schema.TypeInt,
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
			"auto_copy_settings_enabled": {
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
	clusterName := d.Get("cluster_name").(string)

	backupSchedule, _, err := connV2.CloudBackupsAPI.GetBackupSchedule(ctx, projectID, clusterName).Execute()
	if err != nil {
		return diag.Errorf(errorSnapshotBackupScheduleRead, clusterName, err)
	}

	diags := setSchemaFields(d, backupSchedule)
	if diags.HasError() {
		return diags
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}
