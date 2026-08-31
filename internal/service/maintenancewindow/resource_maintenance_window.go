package maintenancewindow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	admin "github.com/mongodb/atlas-sdk-go/admin"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"hour_of_day"},
			},
			"hour_of_day": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"day_of_week"},
			},
			"start_asap": {
				Type:     schema.TypeBool,
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
			"time_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected_hours": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_hour_of_day": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start_hour_of_day": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"wave_assignment": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasPreview
	projectID := d.Get("project_id").(string)

	if deferValue := d.Get("defer").(bool); deferValue {
		_, err := connV2.MaintenanceWindowsAPI.DeferMaintenanceWindow(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, projectID, err))
		}
	}

	params := new(admin.GroupMaintenanceWindowPreviewUpdateRequest)

	if !d.GetRawConfig().GetAttr("day_of_week").IsNull() {
		params.DayOfWeek = new(d.Get("day_of_week").(int))
	}
	if !d.GetRawConfig().GetAttr("hour_of_day").IsNull() {
		params.HourOfDay = new(d.Get("hour_of_day").(int))
	}

	if autoDeferOnceEnabled, ok := d.GetOk("auto_defer_once_enabled"); ok {
		params.AutoDeferOnceEnabled = new(autoDeferOnceEnabled.(bool))
	}

	if !d.GetRawConfig().GetAttr("wave_assignment").IsNull() {
		wave := d.Get("wave_assignment").(int)
		params.WaveAssignment = &wave
	}

	params.ProtectedHours = newProtectedHours(d)
	_, err := connV2.MaintenanceWindowsAPI.UpdateMaintenanceWindow(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceCreate, projectID, err))
	}

	if autoDeferValue := d.Get("auto_defer").(bool); autoDeferValue {
		_, err := connV2.MaintenanceWindowsAPI.ToggleMaintenanceAutoDefer(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, projectID, err))
		}
	}

	d.SetId(projectID)

	return resourceRead(ctx, d, meta)
}

func newProtectedHours(d *schema.ResourceData) *admin.ProtectedHours {
	if protectedHours, ok := d.Get("protected_hours").([]any); ok && len(protectedHours) > 0 {
		item := protectedHours[0].(map[string]any)

		return &admin.ProtectedHours{
			EndHourOfDay:   conversion.IntPtr(item["end_hour_of_day"].(int)),
			StartHourOfDay: conversion.IntPtr(item["start_hour_of_day"].(int)),
		}
	}

	return nil
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasPreview
	projectID := d.Id()

	maintenanceWindow, resp, err := connV2.MaintenanceWindowsAPI.GetMaintenanceWindow(ctx, projectID).Execute()
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

	if err := d.Set("time_zone_id", maintenanceWindow.GetTimeZoneId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}

	if maintenanceWindow.ProtectedHours != nil {
		if err := d.Set("protected_hours", flattenProtectedHours(maintenanceWindow.GetProtectedHours())); err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
		}
	}

	if err := d.Set("wave_assignment", maintenanceWindow.GetWaveAssignment()); err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceRead, projectID, err))
	}
	return nil
}

func flattenProtectedHours(protectedHours admin.ProtectedHours) []map[string]int {
	res := make([]map[string]int, 0)
	res = append(res, map[string]int{
		"end_hour_of_day":   protectedHours.GetEndHourOfDay(),
		"start_hour_of_day": protectedHours.GetStartHourOfDay(),
	})
	return res
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasPreview
	projectID := d.Id()

	if d.HasChange("defer") {
		_, err := connV2.MaintenanceWindowsAPI.DeferMaintenanceWindow(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceDefer, projectID, err))
		}
	}

	params := new(admin.GroupMaintenanceWindowPreviewUpdateRequest)
	// TODO(CLOUDP-440317): omitting day_of_week/hour_of_day leaves them out of the PATCH rather than
	// clearing them, so removing an existing schedule to go wave-only does not converge. Send them as
	// explicit null (SetDayOfWeekNil/SetHourOfDayNil) once the API honors null unset (CLOUDP-439562).
	if !d.GetRawConfig().GetAttr("day_of_week").IsNull() {
		params.DayOfWeek = new(d.Get("day_of_week").(int))
	}
	if !d.GetRawConfig().GetAttr("hour_of_day").IsNull() {
		params.HourOfDay = new(d.Get("hour_of_day").(int))
	}

	if d.HasChange("auto_defer_once_enabled") {
		params.AutoDeferOnceEnabled = new(d.Get("auto_defer_once_enabled").(bool))
	}

	if oldPAny, newPAny := d.GetChange("protected_hours"); d.HasChange("protected_hours") {
		oldP := oldPAny.([]any)
		newP := newPAny.([]any)

		if len(oldP) == 1 && len(newP) == 0 {
			params.ProtectedHours = &admin.ProtectedHours{
				StartHourOfDay: nil,
				EndHourOfDay:   nil,
			}
		} else {
			params.ProtectedHours = newProtectedHours(d)
		}
	}

	if d.HasChange("wave_assignment") {
		// SDKv2 GetOk() cannot distinguish an explicit value (including 0) from an unset field,
		// since TypeInt treats 0 and absent identically. GetRawConfig() allows to distinguish between the two.
		if d.GetRawConfig().GetAttr("wave_assignment").IsNull() {
			params.SetWaveAssignmentNil()
		} else {
			params.SetWaveAssignment(d.Get("wave_assignment").(int))
		}
	}

	_, err := connV2.MaintenanceWindowsAPI.UpdateMaintenanceWindow(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceUpdate, projectID, err))
	}

	if d.HasChange("auto_defer") {
		_, err := connV2.MaintenanceWindowsAPI.ToggleMaintenanceAutoDefer(ctx, projectID).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorMaintenanceAutoDefer, projectID, err))
		}
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasPreview
	projectID := d.Id()

	_, err := connV2.MaintenanceWindowsAPI.ResetMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorMaintenanceDelete, projectID, err))
	}
	return nil
}
