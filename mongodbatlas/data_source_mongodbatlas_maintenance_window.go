package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasMaintenanceWindowRead,
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
		},
	}
}

func dataSourceMongoDBAtlasMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	maintenance, _, err := conn.MaintenanceWindows.Get(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorMaintenanceRead, projectID, err)
	}

	if err := d.Set("day_of_week", maintenance.DayOfWeek); err != nil {
		return fmt.Errorf(errorMaintenanceRead, projectID, err)
	}
	if err := d.Set("hour_of_day", maintenance.HourOfDay); err != nil {
		return fmt.Errorf(errorMaintenanceRead, projectID, err)
	}
	if err := d.Set("number_of_deferrals", maintenance.NumberOfDeferrals); err != nil {
		return fmt.Errorf(errorMaintenanceRead, projectID, err)
	}
	if err := d.Set("start_asap", cast.ToBool(maintenance.StartASAP)); err != nil {
		return fmt.Errorf(errorMaintenanceRead, projectID, err)
	}

	d.SetId(projectID)
	return nil
}
