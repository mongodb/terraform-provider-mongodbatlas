package maintenancewindow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
			"day_of_week": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hour_of_day": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"start_asap": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"number_of_deferrals": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_defer_once_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"time_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected_hours": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_hour_of_day": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_hour_of_day": {
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

	maintenance, _, err := connV2.MaintenanceWindowsApi.GetMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("day_of_week", maintenance.GetDayOfWeek()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("hour_of_day", maintenance.GetHourOfDay()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("number_of_deferrals", maintenance.GetNumberOfDeferrals()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("start_asap", maintenance.GetStartASAP()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("auto_defer_once_enabled", maintenance.GetAutoDeferOnceEnabled()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("time_zone_id", maintenance.GetTimeZoneId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if maintenance.ProtectedHours != nil {
		if err := d.Set("protected_hours", flattenProtectedHours(maintenance.GetProtectedHours())); err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
		}
	}

	d.SetId(projectID)

	return nil
}
