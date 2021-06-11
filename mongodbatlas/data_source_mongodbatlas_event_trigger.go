package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMongoDBAtlasEventTrigger() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasEventTriggerRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasEventTriggerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*MongoDBClient).Realm
	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)
	triggerID := getEncodedID(d.Get("trigger_id").(string), "trigger_id")

	eventResp, _, err := conn.EventTriggers.Get(context.Background(), projectID, appID, triggerID)
	if err != nil {
		return fmt.Errorf(errorEventTriggersRead, projectID, triggerID, err)
	}

	if err := d.Set("name", eventResp.Name); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "name", projectID, triggerID, err)
	}
	if err := d.Set("type", eventResp.Type); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "type", projectID, triggerID, err)
	}
	if err := d.Set("function_id", eventResp.FunctionID); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "function_id", projectID, triggerID, err)
	}
	if err := d.Set("function_name", eventResp.FunctionName); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "function_name", projectID, triggerID, err)
	}
	if err := d.Set("disabled", eventResp.Disabled); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "disabled", projectID, triggerID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": eventResp.ID,
	}))

	return nil
}
