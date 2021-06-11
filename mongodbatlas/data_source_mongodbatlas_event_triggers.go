package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"go.mongodb.org/realm/realm"
)

func dataSourceMongoDBAtlasEventTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasEventTriggersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasEventTriggersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*MongoDBClient).Realm

	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)

	eventTriggers, _, err := conn.EventTriggers.List(context.Background(), projectID, appID)
	if err != nil {
		return fmt.Errorf("error getting event triggers information: %s", err)
	}

	if err := d.Set("results", flattenEventTriggers(eventTriggers)); err != nil {
		return fmt.Errorf("error setting `result` for event triggers: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenEventTriggers(eventTriggers []realm.EventTrigger) []map[string]interface{} {
	var triggersMap []map[string]interface{}

	if len(eventTriggers) > 0 {
		triggersMap = make([]map[string]interface{}, len(eventTriggers))

		for i := range eventTriggers {
			triggersMap[i] = map[string]interface{}{
				"trigger_id":    eventTriggers[i].ID,
				"name":          eventTriggers[i].Name,
				"type":          eventTriggers[i].Type,
				"function_id":   eventTriggers[i].FunctionID,
				"function_name": eventTriggers[i].FunctionName,
				"disabled":      eventTriggers[i].Disabled,
			}
		}
	}

	return triggersMap
}
