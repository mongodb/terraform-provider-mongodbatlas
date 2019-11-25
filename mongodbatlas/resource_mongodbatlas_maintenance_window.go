package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorMaintenanceCreate = "error create MongoDB Mainteanace Window (%s): %s"
	errorMaintenanceUpdate = "error update MongoDB Mainteanace Window (%s): %s"
	errorMaintenanceRead   = "error reading MongoDB Mainteanace Window (%s): %s"
	errorMaintenanceDelete = "error delete MongoDB Mainteanace Window (%s): %s"
	errorMaintenanceDefer  = "error defer MongoDB Mainteanace Window (%s): %s"
)

func resourceMongoDBAtlasMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasMaintenanceWindowCreate,
		Read:   resourceMongoDBAtlasMaintenanceWindowRead,
		Update: resourceMongoDBAtlasMaintenanceWindowUpdate,
		Delete: resourceMongoDBAtlasMaintenanceWindowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"day_of_week": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"start_asap"},

				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 23 {
						errs = append(errs, fmt.Errorf("%q value should be between 0 and 23, got: %d", key, v))
					}
					return
				},
			},
			"start_asap": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"day_of_week", "hour_of_day", "number_of_deferrals", "number_of_deferrals"},
			},
			"number_of_deferrals": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"start_asap"},
			},
			"defer": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"start_asap"},
			},
		},
	}
}

func resourceMongoDBAtlasMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	if deferValue := d.Get("defer").(bool); deferValue {
		_, err := conn.MaintenanceWindows.Defer(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf(errorMaintenanceDefer, projectID, err)
		}
	}

	maintenanceWindowReq := &matlas.MaintenanceWindow{}
	if dayOfWeek, ok := d.GetOk("day_of_week"); ok {
		maintenanceWindowReq.DayOfWeek = cast.ToInt(dayOfWeek)
	}
	if hourOfDay, ok := d.GetOkExists("hour_of_day"); ok {
		maintenanceWindowReq.HourOfDay = pointy.Int(cast.ToInt(hourOfDay))
	}
	if startASAP, ok := d.GetOkExists("start_asap"); ok {
		maintenanceWindowReq.StartASAP = pointy.Bool(cast.ToBool(startASAP))
	}
	if numberOfDeferrals, ok := d.GetOk("number_of_deferrals"); ok {
		maintenanceWindowReq.NumberOfDeferrals = cast.ToInt(numberOfDeferrals)
	}

	_, err := conn.MaintenanceWindows.Update(context.Background(), projectID, maintenanceWindowReq)
	if err != nil {
		return fmt.Errorf(errorMaintenanceCreate, projectID, err)
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasMaintenanceWindowRead(d, meta)
}

func resourceMongoDBAtlasMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	maintenanceWindow, _, err := conn.MaintenanceWindows.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorMaintenanceRead, d.Id(), err)
	}

	if err := d.Set("day_of_week", maintenanceWindow.DayOfWeek); err != nil {
		return fmt.Errorf(errorMaintenanceRead, d.Id(), err)
	}
	if err := d.Set("hour_of_day", maintenanceWindow.HourOfDay); err != nil {
		return fmt.Errorf(errorMaintenanceRead, d.Id(), err)
	}
	if err := d.Set("number_of_deferrals", maintenanceWindow.NumberOfDeferrals); err != nil {
		return fmt.Errorf(errorMaintenanceRead, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	maintenanceWindowReq := &matlas.MaintenanceWindow{}

	if d.HasChange("defer") {
		_, err := conn.MaintenanceWindows.Defer(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf(errorMaintenanceDefer, d.Id(), err)
		}
	}

	if d.HasChange("day_of_week") {
		maintenanceWindowReq.DayOfWeek = cast.ToInt(d.Get("day_of_week"))
	}
	if d.HasChange("hour_of_day") {
		maintenanceWindowReq.HourOfDay = pointy.Int(cast.ToInt(d.Get("hour_of_day")))
	}
	if d.HasChange("number_of_deferrals") {
		maintenanceWindowReq.NumberOfDeferrals = cast.ToInt(d.Get("number_of_deferrals"))
	}

	_, err := conn.MaintenanceWindows.Update(context.Background(), d.Id(), maintenanceWindowReq)
	if err != nil {
		return fmt.Errorf(errorMaintenanceUpdate, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
	//Get the client connection.
	conn := meta.(*matlas.Client)

	_, err := conn.MaintenanceWindows.Reset(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorMaintenanceDelete, d.Id(), err)
	}

	return nil
}
