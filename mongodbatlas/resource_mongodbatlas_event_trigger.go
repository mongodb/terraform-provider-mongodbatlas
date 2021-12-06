package mongodbatlas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		CreateContext: resourceMongoDBAtlasEventTriggersCreate,
		ReadContext:   resourceMongoDBAtlasEventTriggersRead,
		UpdateContext: resourceMongoDBAtlasEventTriggersUpdate,
		DeleteContext: resourceMongoDBAtlasEventTriggersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasEventTriggerImportState,
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
				ValidateFunc: validation.StringInSlice([]string{"DATABASE", "AUTHENTICATION", "SCHEDULED"}, false),
			},
			"function_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"event_processors"},
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"config_operation_types": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"INSERT", "UPDATE", "REPLACE", "DELETE"}, false),
				},
			},
			"config_operation_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOGIN", "CREATE", "DELETE"}, false),
			},
			"config_providers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"anon-user", "local-userpass", "api-key", "custom-token", "custom-function", "oauth2-facebook", "oauth2-google", "oauth2-apple"}, false),
				},
			},
			"config_database": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_collection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_match": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					var j, j2 interface{}
					if err := json.Unmarshal([]byte(old), &j); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
					}
					if err := json.Unmarshal([]byte(new), &j2); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
					}
					if diff := deep.Equal(&j, &j2); diff != nil {
						log.Printf("[DEBUG] deep equal not passed: %v", diff)
						return false
					}

					return true
				},
			},
			"config_project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					var j, j2 interface{}
					if err := json.Unmarshal([]byte(old), &j); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
					}
					if err := json.Unmarshal([]byte(new), &j2); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
					}
					if diff := deep.Equal(&j, &j2); diff != nil {
						log.Printf("[DEBUG] deep equal not passed: %v", diff)
						return false
					}

					return true
				},
			},
			"config_full_document": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"config_full_document_before": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"config_schedule": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_schedule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_processors": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"function_id"},
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

func resourceMongoDBAtlasEventTriggersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)
	typeTrigger := d.Get("type").(string)
	eventTriggerReq := &realm.EventTriggerRequest{
		Name:       d.Get("name").(string),
		Type:       typeTrigger,
		FunctionID: d.Get("function_id").(string),
	}

	if v, ok := d.GetOk("disabled"); ok {
		eventTriggerReq.Disabled = pointy.Bool(v.(bool))
	}

	eventTriggerConfig := &realm.EventTriggerConfig{}

	cots, okCots := d.GetOk("config_operation_types")
	cot, okCot := d.GetOk("config_operation_type")
	pro, okP := d.GetOk("config_providers")
	data, okD := d.GetOk("config_database")
	coll, okC := d.GetOk("config_collection")
	si, okSI := d.GetOk("config_service_id")
	sche, oksch := d.GetOk("config_schedule")

	if typeTrigger == "DATABASE" {
		if !okCots || !okD || !okC || !okSI {
			return diag.FromErr(fmt.Errorf("`config_operation_types`, `config_database`,`config_collection`,`config_service_id` must be provided if type is DATABASE"))
		}
	}
	if typeTrigger == "AUTHENTICATION" {
		if !okCot || !okP {
			return diag.FromErr(fmt.Errorf("`config_operation_type`, `config_providers` must be provided if type is AUTHENTICATION"))
		}
	}
	if typeTrigger == "SCHEDULED" {
		if !oksch {
			return diag.FromErr(fmt.Errorf("`config_schedule` must be provided if type is SCHEDULED"))
		}
	}

	if okCots {
		eventTriggerConfig.OperationTypes = cast.ToStringSlice(cots)
	}
	if okCot {
		eventTriggerConfig.OperationType = cot.(string)
	}
	if okP {
		eventTriggerConfig.Providers = cast.ToStringSlice(pro)
	}
	if okD {
		eventTriggerConfig.Database = data.(string)
	}
	if okC {
		eventTriggerConfig.Collection = coll.(string)
	}
	if okSI {
		eventTriggerConfig.ServiceID = si.(string)
	}
	if v, ok := d.GetOk("config_match"); ok {
		eventTriggerConfig.Match = cast.ToStringMap(v)
	}
	if v, ok := d.GetOk("config_project"); ok {
		eventTriggerConfig.Project = cast.ToStringMap(v)
	}
	if v, ok := d.GetOk("config_full_document"); ok {
		eventTriggerConfig.FullDocument = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("config_full_document_before"); ok {
		eventTriggerConfig.FullDocumentBeforeChange = pointy.Bool(v.(bool))
	}
	if oksch {
		eventTriggerConfig.Schedule = sche.(string)
	}

	if v, ok := d.GetOk("event_processors"); ok {
		eventTriggerReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(v.([]interface{}))
	}

	eventTriggerReq.Config = eventTriggerConfig

	eventResp, _, err := conn.EventTriggers.Create(context.Background(), projectID, appID, eventTriggerReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersCreate, projectID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": eventResp.ID,
	}))

	return resourceMongoDBAtlasEventTriggersRead(ctx, d, meta)
}

func resourceMongoDBAtlasEventTriggersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	resp, recodes, err := conn.EventTriggers.Get(ctx, projectID, appID, triggerID)
	if err != nil {
		if recodes != nil && recodes.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorEventTriggersRead, projectID, appID, err))
	}

	if err = d.Set("trigger_id", resp.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "trigger_id", projectID, appID, err))
	}
	if err = d.Set("name", resp.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "name", projectID, appID, err))
	}
	if err = d.Set("type", resp.Type); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "type", projectID, appID, err))
	}
	functionID := resp.FunctionID
	if resp.FunctionID == "000000000000000000000000" {
		functionID = ""
	}
	if err = d.Set("function_id", functionID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "function_id", projectID, appID, err))
	}
	if err = d.Set("function_name", resp.FunctionName); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "function_name", projectID, appID, err))
	}
	if err = d.Set("disabled", resp.Disabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "disabled", projectID, appID, err))
	}
	if err = d.Set("config_operation_types", resp.Config.OperationTypes); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_operation_types", projectID, appID, err))
	}
	if err = d.Set("config_operation_type", resp.Config.OperationType); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_operation_type", projectID, appID, err))
	}
	if err = d.Set("config_providers", resp.Config.Providers); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_providers", projectID, appID, err))
	}
	if err = d.Set("config_database", resp.Config.Database); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_database", projectID, appID, err))
	}
	if err = d.Set("config_collection", resp.Config.Collection); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_collection", projectID, appID, err))
	}
	if err = d.Set("config_service_id", resp.Config.ServiceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_service_id", projectID, appID, err))
	}
	if err = d.Set("config_match", matchToString(resp.Config.Match)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_match", projectID, appID, err))
	}
	if err = d.Set("config_project", matchToString(resp.Config.Project)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_project", projectID, appID, err))
	}
	if err = d.Set("config_full_document", resp.Config.FullDocument); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_full_document", projectID, appID, err))
	}

	if err = d.Set("config_full_document_before", resp.Config.FullDocumentBeforeChange); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_full_document", projectID, appID, err))
	}

	if err = d.Set("config_schedule", resp.Config.Schedule); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_schedule", projectID, appID, err))
	}
	if err = d.Set("config_schedule_type", resp.Config.ScheduleType); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_schedule_type", projectID, appID, err))
	}
	if err = d.Set("event_processors", flattenTriggerEventProcessorAWSEventBridge(resp.EventProcessors)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "event_processors", projectID, appID, err))
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection.
	conn, err := meta.(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]
	typeTrigger := d.Get("type").(string)

	eventReq := &realm.EventTriggerRequest{
		Name:       d.Get("name").(string),
		Type:       typeTrigger,
		FunctionID: d.Get("function_id").(string),
	}
	eventTriggerConfig := &realm.EventTriggerConfig{}

	if d.HasChange("disabled") {
		eventReq.Disabled = pointy.Bool(d.Get("disabled").(bool))
	}
	if typeTrigger == "DATABASE" {
		eventTriggerConfig.OperationTypes = cast.ToStringSlice(d.Get("config_operation_types"))
		eventTriggerConfig.Database = d.Get("config_database").(string)
		eventTriggerConfig.Collection = d.Get("config_collection").(string)
		eventTriggerConfig.ServiceID = d.Get("config_service_id").(string)
		eventTriggerConfig.Match = cast.ToStringMap(d.Get("config_match").(string))
		eventTriggerConfig.Project = cast.ToStringMap(d.Get("config_project").(string))
		eventTriggerConfig.FullDocument = pointy.Bool(d.Get("config_full_document").(bool))
		eventTriggerConfig.FullDocumentBeforeChange = pointy.Bool(d.Get("config_full_document_before").(bool))
	}
	if typeTrigger == "AUTHENTICATION" {
		eventTriggerConfig.OperationType = d.Get("config_operation_type").(string)
		eventTriggerConfig.Providers = cast.ToStringSlice(d.Get("config_providers"))
	}

	if typeTrigger == "SCHEDULED" {
		eventTriggerConfig.Schedule = d.Get("config_schedule").(string)
	}
	eventReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(d.Get("event_processors").([]interface{}))

	eventReq.Config = eventTriggerConfig

	_, _, err = conn.EventTriggers.Update(ctx, projectID, appID, triggerID, eventReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersUpdate, projectID, appID, err))
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection.
	conn, err := meta.(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	ids := decodeStateID(d.Id())

	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	_, err = conn.EventTriggers.Delete(ctx, projectID, appID, triggerID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersDelete, projectID, appID, err))
	}

	return nil
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

func flattenTriggerEventProcessorAWSEventBridge(eventProcessor map[string]interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	if eventProcessor != nil && eventProcessor["AWS_EVENTBRIDGE"] != nil {
		event := eventProcessor["AWS_EVENTBRIDGE"].(map[string]interface{})
		config := event["config"].(map[string]interface{})
		mapEvent := map[string]interface{}{
			"aws_eventbridge": []map[string]interface{}{
				{
					"config_account_id": config["account_id"].(string),
					"config_region":     config["region"].(string),
				},
			},
		}
		results = append(results, mapEvent)
	}

	return results
}

func resourceMongoDBAtlasEventTriggerImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn, err := meta.(*MongoDBClient).GetRealmClient(ctx)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(d.Id(), "--")
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a MongoDB Event Trigger, use the format {project_id}-{app_id}-{trigger_id} ")
	}

	projectID := parts[0]
	appID := parts[1]
	triggerID := parts[2]

	_, _, err = conn.EventTriggers.Get(ctx, projectID, appID, triggerID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import event trigger %s in project %s, error: %s", triggerID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorEventTriggersSetting, "project_id", projectID, appID, err)
	}
	if err := d.Set("app_id", appID); err != nil {
		return nil, fmt.Errorf(errorEventTriggersSetting, "app_id", projectID, appID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": triggerID,
	}))

	return []*schema.ResourceData{d}, nil
}

func matchToString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("[ERROR] %v ", err)
	}
	return string(b)
}
