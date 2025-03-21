package eventtrigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb-labs/go-client-mongodb-atlas-app-services/appservices"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorEventTriggersCreate  = "error creating MongoDB EventTriggers (%s): %s"
	errorEventTriggersUpdate  = "error updating MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersRead    = "error reading MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersDelete  = "error deleting MongoDB EventTriggers (%s)%s: %s"
	errorEventTriggersSetting = "error setting `%s` for EventTriggers(%s)%s: %s"
)

func Resource() *schema.Resource {
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
				ForceNew: true,
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
				DiffSuppressFunc: func(k, oldAttr, newAttr string, d *schema.ResourceData) bool {
					var j, j2 any
					if err := json.Unmarshal([]byte(oldAttr), &j); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
						return false
					}
					if err := json.Unmarshal([]byte(newAttr), &j2); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
						return false
					}
					if !reflect.DeepEqual(&j, &j2) {
						return false
					}

					return true
				},
			},
			"config_project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, oldAttr, newAttr string, d *schema.ResourceData) bool {
					var j, j2 any
					if err := json.Unmarshal([]byte(oldAttr), &j); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
						return false
					}
					if err := json.Unmarshal([]byte(newAttr), &j2); err != nil {
						log.Printf("[ERROR] json.Unmarshal %v", err)
						return false
					}
					if !reflect.DeepEqual(&j, &j2) {
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
			"unordered": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasEventTriggersCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)
	typeTrigger := d.Get("type").(string)
	eventTriggerReq := &appservices.EventTriggerRequest{
		Name:       d.Get("name").(string),
		Type:       typeTrigger,
		FunctionID: d.Get("function_id").(string),
	}

	if v, ok := d.GetOk("disabled"); ok {
		eventTriggerReq.Disabled = conversion.Pointer(v.(bool))
	}

	eventTriggerConfig := &appservices.EventTriggerConfig{}

	cots, okCots := d.GetOk("config_operation_types")
	cot, okCot := d.GetOk("config_operation_type")
	pro, okP := d.GetOk("config_providers")
	data, okD := d.GetOk("config_database")
	coll, okC := d.GetOk("config_collection")
	si, okSI := d.GetOk("config_service_id")
	sche, oksch := d.GetOk("config_schedule")

	if typeTrigger == "DATABASE" {
		if !okCots || !okD || !okSI {
			return diag.FromErr(fmt.Errorf("`config_operation_types`, `config_database`,`config_service_id` must be provided if type is DATABASE"))
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
		eventTriggerConfig.FullDocument = conversion.Pointer(v.(bool))
	}
	if v, ok := d.GetOk("config_full_document_before"); ok {
		eventTriggerConfig.FullDocumentBeforeChange = conversion.Pointer(v.(bool))
	}

	if oksch {
		eventTriggerConfig.Schedule = sche.(string)
	}

	if v, ok := d.GetOk("event_processors"); ok {
		eventTriggerReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(v.([]any))
	}

	if v, ok := d.GetOk("unordered"); ok {
		eventTriggerConfig.Unordered = conversion.Pointer(v.(bool))
	}

	eventTriggerReq.Config = eventTriggerConfig

	eventResp, _, err := conn.EventTriggers.Create(context.Background(), projectID, appID, eventTriggerReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersCreate, projectID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": eventResp.ID,
	}))

	return resourceMongoDBAtlasEventTriggersRead(ctx, d, meta)
}

func resourceMongoDBAtlasEventTriggersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := conversion.DecodeStateID(d.Id())
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
	if err = d.Set("unordered", resp.Config.Unordered); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "unordered", projectID, appID, err))
	}
	if err = d.Set("event_processors", flattenTriggerEventProcessorAWSEventBridge(resp.EventProcessors)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "event_processors", projectID, appID, err))
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]
	typeTrigger := d.Get("type").(string)

	eventReq := &appservices.EventTriggerRequest{
		Name:       d.Get("name").(string),
		Type:       typeTrigger,
		FunctionID: d.Get("function_id").(string),
		Disabled:   conversion.Pointer(d.Get("disabled").(bool)),
	}
	eventTriggerConfig := &appservices.EventTriggerConfig{}

	if typeTrigger == "DATABASE" {
		eventTriggerConfig.OperationTypes = cast.ToStringSlice(d.Get("config_operation_types"))
		eventTriggerConfig.Database = d.Get("config_database").(string)
		eventTriggerConfig.Collection = d.Get("config_collection").(string)
		eventTriggerConfig.ServiceID = d.Get("config_service_id").(string)
		eventTriggerConfig.Match = cast.ToStringMap(d.Get("config_match").(string))
		eventTriggerConfig.Project = cast.ToStringMap(d.Get("config_project").(string))
		eventTriggerConfig.FullDocument = conversion.Pointer(d.Get("config_full_document").(bool))
		eventTriggerConfig.FullDocumentBeforeChange = conversion.Pointer(d.Get("config_full_document_before").(bool))
		eventTriggerConfig.Unordered = conversion.Pointer(d.Get("unordered").(bool))
	}
	if typeTrigger == "AUTHENTICATION" {
		eventTriggerConfig.OperationType = d.Get("config_operation_type").(string)
		eventTriggerConfig.Providers = cast.ToStringSlice(d.Get("config_providers"))
	}

	if typeTrigger == "SCHEDULED" {
		eventTriggerConfig.Schedule = d.Get("config_schedule").(string)
	}
	eventReq.EventProcessors = expandTriggerEventProcessorAWSEventBridge(d.Get("event_processors").([]any))

	eventReq.Config = eventTriggerConfig

	_, _, err = conn.EventTriggers.Update(ctx, projectID, appID, triggerID, eventReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersUpdate, projectID, appID, err))
	}

	return nil
}

func resourceMongoDBAtlasEventTriggersDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get the client connection.
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	ids := conversion.DecodeStateID(d.Id())

	projectID := ids["project_id"]
	appID := ids["app_id"]
	triggerID := ids["trigger_id"]

	_, err = conn.EventTriggers.Delete(ctx, projectID, appID, triggerID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersDelete, projectID, appID, err))
	}

	return nil
}

func expandTriggerEventProcessorAWSEventBridge(p []any) map[string]any {
	if len(p) == 0 {
		return nil
	}

	aws := p[0].(map[string]any)
	event := aws["aws_eventbridge"].([]any)
	if len(event) == 0 {
		return nil
	}
	eventObj := event[0].(map[string]any)

	return map[string]any{
		"AWS_EVENTBRIDGE": map[string]any{
			"type": "AWS_EVENTBRIDGE",
			"config": map[string]any{
				"account_id": eventObj["config_account_id"].(string),
				"region":     eventObj["config_region"].(string),
			},
		},
	}
}

func flattenTriggerEventProcessorAWSEventBridge(eventProcessor map[string]any) []map[string]any {
	results := make([]map[string]any, 0)
	if eventProcessor != nil && eventProcessor["AWS_EVENTBRIDGE"] != nil {
		event := eventProcessor["AWS_EVENTBRIDGE"].(map[string]any)
		cfg := event["config"].(map[string]any)
		mapEvent := map[string]any{
			"aws_eventbridge": []map[string]any{
				{
					"config_account_id": cfg["account_id"].(string),
					"config_region":     cfg["region"].(string),
				},
			},
		}
		results = append(results, mapEvent)
	}

	return results
}

func resourceMongoDBAtlasEventTriggerImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(d.Id(), "--")
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a MongoDB Event Trigger, use the format {project_id}--{app_id}--{trigger_id} ")
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

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": triggerID,
	}))

	return []*schema.ResourceData{d}, nil
}

func matchToString(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("[ERROR] %v ", err)
	}
	return string(b)
}
