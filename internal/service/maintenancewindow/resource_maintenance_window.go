package maintenancewindow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	errorMaintenanceCreate    = "error creating the MongoDB Atlas Maintenance Window (%s): %s"
	errorMaintenanceUpdate    = "error updating the MongoDB Atlas Maintenance Window (%s): %s"
	errorMaintenanceRead      = "error reading the MongoDB Atlas Maintenance Window (%s): %s"
	errorMaintenanceDelete    = "error deleting the MongoDB Atlas Maintenance Window (%s): %s"
	errorMaintenanceDefer     = "error deferring the MongoDB Atlas Maintenance Window (%s): %s"
	errorMaintenanceAutoDefer = "error auto deferring the MongoDB Atlas Maintenance Window (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"day_of_week": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 7 {
						errs = append(errs, fmt.Errorf("%q value should be between 1 and 7, got: %d", key, v))
					}
					return
				},
			},
			"hour_of_day": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"start_asap"},
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 23 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 23, got: %d", key, v))
					}
					return
				},
			},
			"start_asap": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"number_of_deferrals": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"defer": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_defer": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_defer_once_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	if deferValue := d.Get("defer").(bool); deferValue {
		_, err := connV2.MaintenanceWindowsApi.DeferMaintenanceWindow(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, projectID, err))
		}
	}

	params := new(admin.GroupMaintenanceWindow)

	params.DayOfWeek = cast.ToInt(d.Get("day_of_week"))

	hourOfDay := d.Get("hour_of_day")
	params.HourOfDay = conversion.Pointer(cast.ToInt(hourOfDay)) // during creation of maintenance window hourOfDay needs to be set in PATCH to avoid errors, 0 value is sent when absent

	if autoDeferOnceEnabled, ok := d.GetOk("auto_defer_once_enabled"); ok {
		params.AutoDeferOnceEnabled = conversion.Pointer(autoDeferOnceEnabled.(bool))
	}

	_, _, err := connV2.MaintenanceWindowsApi.UpdateMaintenanceWindow(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceCreate, projectID, err))
	}

	if autoDeferValue := d.Get("auto_defer").(bool); autoDeferValue {
		_, err := connV2.MaintenanceWindowsApi.ToggleMaintenanceAutoDefer(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, projectID, err))
		}
	}

	d.SetId(projectID)

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()

	maintenanceWindow, resp, err := connV2.MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if err := d.Set("day_of_week", maintenanceWindow.GetDayOfWeek()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if err := d.Set("hour_of_day", maintenanceWindow.GetHourOfDay()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if err := d.Set("number_of_deferrals", maintenanceWindow.GetNumberOfDeferrals()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if err := d.Set("start_asap", maintenanceWindow.GetStartASAP()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if err := d.Set("auto_defer_once_enabled", maintenanceWindow.GetAutoDeferOnceEnabled()); err != nil {
		return diag.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()

	if d.HasChange("defer") {
		_, err := connV2.MaintenanceWindowsApi.DeferMaintenanceWindow(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, projectID, err))
		}
	}

	params := new(admin.GroupMaintenanceWindow)
	params.DayOfWeek = cast.ToInt(d.Get("day_of_week"))

	if d.HasChange("hour_of_day") {
		params.HourOfDay = conversion.Pointer(cast.ToInt(d.Get("hour_of_day")))
	}

	if d.HasChange("auto_defer_once_enabled") {
		params.AutoDeferOnceEnabled = conversion.Pointer(d.Get("auto_defer_once_enabled").(bool))
	}

	_, _, err := connV2.MaintenanceWindowsApi.UpdateMaintenanceWindow(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceUpdate, projectID, err))
	}

	if d.HasChange("auto_defer") {
		_, err := connV2.MaintenanceWindowsApi.ToggleMaintenanceAutoDefer(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, projectID, err))
		}
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Id()

	_, err := connV2.MaintenanceWindowsApi.ResetMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceDelete, projectID, err))
	}
	return nil
}
