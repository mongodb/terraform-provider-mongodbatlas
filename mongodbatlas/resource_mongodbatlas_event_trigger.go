package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	"go.mongodb.org/realm/realm"
)

const (
	errorEventTriggersCreate  = "error creating MongoDB EventTriggers (%s): %s"
	errorEventTriggersUpdate  = "error updating MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersRead    = "error reading MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersDelete  = "error deleting MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersSetting = "error setting `%s` for EventTriggers(%s)%s: %s"
)

func resourceMongoDBAtlasEventTriggers() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasEventTriggersCreate,
		Read:   resourceMongoDBAtlasEventTriggersRead,
		Update: resourceMongoDBAtlasEventTriggersUpdate,
		Delete: resourceMongoDBAtlasEventTriggersDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasEventTriggerImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DATABASE", "AUTHENTICATION"}, false),
			},
			"function_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"config_operation_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"INSERT", "UPDATE", "REPLACE", "DELETE"}, false),
				},
			},
			"config_operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOGIN", "CREATE", "DELETE"}, false),
			},
			"config_providers": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_collection": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_match": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"config_full_document": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"config_schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_processors": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_eventbridge": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"config_region": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"function_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasEventTriggersCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*MongoDBClient).Realm

	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)
	eventTriggerReq := &realm.EventTriggerRequest{
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		FunctionID: d.Get("function_id").(string),
	}

	if v, ok := d.GetOk("disabled"); ok {
		eventTriggerReq.Disabled = pointy.Bool(v.(bool))
	}

	eventTriggerConfig := &realm.EventTriggerConfig{
		OperationTypes: []interface{}{cast.ToStringSlice(d.Get("config_operation_types"))},
		OperationType:  d.Get("config_operation_type").(string),
		Providers:      d.Get("config_providers").(string),
		Database:       d.Get("config_database").(string),
		Collection:     d.Get("config_collection").(string),
		ServiceID:      d.Get("config_service_id").(string),
	}

	if v, ok := d.GetOk("config_match"); ok {
		eventTriggerConfig.Match = expandTriggerConfigMatch(v.([]interface{}))
	}
	if v, ok := d.GetOk("config_full_document"); ok {
		eventTriggerConfig.FullDocument = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("config_schedule"); ok {
		eventTriggerConfig.Schedule = v.(string)
	}

	if v, ok := d.GetOk("event_processors"); ok {
		eventTriggerReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(v.([]interface{}))
	}

	eventTriggerReq.Config = eventTriggerConfig

	eventResp, _, err := conn.EventTriggers.Create(context.Background(), projectID, appID, eventTriggerReq)
	if err != nil {
		return fmt.Errorf(errorEventTriggersCreate, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": eventResp.ID,
	}))

	return resourceMongoDBAtlasEventTriggersRead(d, meta)
}

func resourceMongoDBAtlasEventTriggersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*MongoDBClient).Realm

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	resp, _, err := conn.EventTriggers.Get(context.Background(), projectID, appID, triggerID)
	if err != nil {
		return fmt.Errorf(errorEventTriggersRead, projectID, appID, err)
	}

	if err = d.Set("trigger_id", resp.ID); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "trigger_id", projectID, appID, err)
	}
	if err = d.Set("name", resp.Name); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "name", projectID, appID, err)
	}
	if err = d.Set("type", resp.Type); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "type", projectID, appID, err)
	}
	if err = d.Set("function_id", resp.FunctionID); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "function_id", projectID, appID, err)
	}
	if err = d.Set("function_name", resp.FunctionName); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "function_name", projectID, appID, err)
	}
	if err = d.Set("disabled", resp.Disabled); err != nil {
		return fmt.Errorf(errorEventTriggersSetting, "disabled", projectID, appID, err)
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get the client connection.
	conn := meta.(*MongoDBClient).Realm

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	eventReq := &realm.EventTriggerRequest{}
	eventTriggerConfig := &realm.EventTriggerConfig{}

	if d.HasChange("name") {
		eventReq.Name = d.Get("name").(string)
	}
	if d.HasChange("type") {
		eventReq.Type = d.Get("type").(string)
	}
	if d.HasChange("function_id") {
		eventReq.FunctionID = d.Get("function_id").(string)
	}
	if d.HasChange("disabled") {
		eventReq.Disabled = pointy.Bool(d.Get("disabled").(bool))
	}
	if d.HasChange("config_operation_types") {
		eventTriggerConfig.OperationTypes = []interface{}{cast.ToStringSlice(d.Get("config_operation_types"))}
	}
	if d.HasChange("config_operation_type") {
		eventTriggerConfig.OperationType = d.Get("config_operation_type").(string)
	}
	if d.HasChange("config_providers") {
		eventTriggerConfig.Providers = d.Get("config_providers").(string)
	}
	if d.HasChange("config_database") {
		eventTriggerConfig.Database = d.Get("config_database").(string)
	}
	if d.HasChange("config_collection") {
		eventTriggerConfig.Collection = d.Get("config_collection").(string)
	}
	if d.HasChange("config_service_id") {
		eventTriggerConfig.ServiceID = d.Get("config_service_id").(string)
	}
	if d.HasChange("config_match") {
		eventTriggerConfig.Match = expandTriggerConfigMatch(d.Get("config_match").([]interface{}))
	}
	if d.HasChange("config_full_document") {
		eventTriggerConfig.FullDocument = pointy.Bool(d.Get("config_full_document").(bool))
	}
	if d.HasChange("config_schedule") {
		eventTriggerConfig.Schedule = d.Get("config_schedule").(string)
	}
	if d.HasChange("event_processors") {
		eventReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(d.Get("event_processors").([]interface{}))
	}

	eventReq.Config = eventTriggerConfig

	_, _, err := conn.EventTriggers.Update(context.Background(), projectID, appID, triggerID, eventReq)
	if err != nil {
		return fmt.Errorf(errorEventTriggersUpdate, projectID, appID, err)
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersDelete(d *schema.ResourceData, meta interface{}) error {
	// Get the client connection.
	conn := meta.(*MongoDBClient).Realm
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	_, err := conn.EventTriggers.Delete(context.Background(), projectID, appID, triggerID)
	if err != nil {
		return fmt.Errorf(errorEventTriggersDelete, projectID, appID, err)
	}

	return nil
}

func expandTriggerConfigMatch(p []interface{}) map[string]interface{} {
	if len(p) == 0 {
		return nil
	}
	matchObj := make(map[string]interface{}, 1)

	match := p[0].(map[string]interface{})

	matchObj[match["key"].(string)] = match["value"].(string)

	return matchObj
}

func expandTriggerEventProcessorAWSEventBridge(p []interface{}) map[string]interface{} {
	if len(p) == 0 {
		return nil
	}

	aws := p[0].(map[string]interface{})
	event := aws["aws_eventbridge"].([]interface{})
	if len(event) == 0 {
		return nil
	}
	eventObj := event[0].(map[string]interface{})

	return map[string]interface{}{
		"AWS_EVENTBRIDGE": map[string]interface{}{
			"type": "AWS_EVENTBRIDGE",
			"config": map[string]interface{}{
				"account_id": eventObj["config_account_id"].(string),
				"region":     eventObj["config_region"].(string),
			},
		},
	}
}

func resourceMongoDBAtlasEventTriggerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*realm.Client)

	parts := strings.Split(d.Id(), "--")
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a MongoDB Event Trigger, use the format {project_id}-{app_id}-{trigger_id} ")
	}

	projectID := parts[0]
	appID := parts[1]
	triggerID := parts[2]

	_, _, err := conn.EventTriggers.Get(context.Background(), projectID, appID, triggerID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import event trigger %s in project %s, error: %s", triggerID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": triggerID,
	}))

	return []*schema.ResourceData{d}, nil
}
