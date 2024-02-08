package maintenancewindow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
	// Get the client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	if deferValue := d.Get("defer").(bool); deferValue {
		_, err := conn.MaintenanceWindows.Defer(ctx, projectID)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, projectID, err))
		}
	}

	maintenanceWindowReq := &matlas.MaintenanceWindow{}

	if dayOfWeek, ok := d.GetOk("day_of_week"); ok {
		maintenanceWindowReq.DayOfWeek = cast.ToInt(dayOfWeek)
	}

	if hourOfDay, ok := d.GetOk("hour_of_day"); ok {
		maintenanceWindowReq.HourOfDay = pointy.Int(cast.ToInt(hourOfDay))
	}

	if autoDeferOnceEnabled, ok := d.GetOk("auto_defer_once_enabled"); ok {
		maintenanceWindowReq.AutoDeferOnceEnabled = pointy.Bool(autoDeferOnceEnabled.(bool))
	}

	_, err := conn.MaintenanceWindows.Update(ctx, projectID, maintenanceWindowReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceCreate, projectID, err))
	}

	if autoDeferValue := d.Get("auto_defer").(bool); autoDeferValue {
		_, err := conn.MaintenanceWindows.AutoDefer(ctx, projectID)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, projectID, err))
		}
	}

	d.SetId(projectID)

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	maintenanceWindow, resp, err := connV2.MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), d.Id()).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	if err := d.Set("day_of_week", maintenanceWindow.GetDayOfWeek()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	if err := d.Set("hour_of_day", maintenanceWindow.GetHourOfDay()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	if err := d.Set("number_of_deferrals", maintenanceWindow.GetNumberOfDeferrals()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	if err := d.Set("start_asap", maintenanceWindow.GetStartASAP()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	if err := d.Set("auto_defer_once_enabled", maintenanceWindow.GetAutoDeferOnceEnabled()); err != nil {
		return diag.Errorf(errorMaintenanceRead, d.Id(), err)
	}

	if err := d.Set("project_id", d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, d.Id(), err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get the client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	maintenanceWindowReq := &matlas.MaintenanceWindow{}

	if d.HasChange("defer") {
		_, err := conn.MaintenanceWindows.Defer(ctx, d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, d.Id(), err))
		}
	}

	if d.HasChange("day_of_week") {
		maintenanceWindowReq.DayOfWeek = cast.ToInt(d.Get("day_of_week"))
	}

	if d.HasChange("hour_of_day") {
		maintenanceWindowReq.HourOfDay = pointy.Int(cast.ToInt(d.Get("hour_of_day")))
	}

	if d.HasChange("auto_defer_once_enabled") {
		maintenanceWindowReq.AutoDeferOnceEnabled = pointy.Bool(d.Get("auto_defer_once_enabled").(bool))
	}

	_, err := conn.MaintenanceWindows.Update(ctx, d.Id(), maintenanceWindowReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceUpdate, d.Id(), err))
	}

	if d.HasChange("auto_defer") {
		_, err := conn.MaintenanceWindows.AutoDefer(ctx, d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, d.Id(), err))
		}
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get the client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	_, err := conn.MaintenanceWindows.Reset(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceDelete, d.Id(), err))
	}

	return nil
}
