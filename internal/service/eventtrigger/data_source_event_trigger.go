package eventtrigger

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasEventTriggerRead,
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
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config_operation_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config_operation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config_database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_collection": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_match": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_full_document": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config_full_document_before": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config_schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_schedule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_processors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_eventbridge": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"config_region": {
										Type:     schema.TypeString,
										Computed: true,
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
			"unordered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasEventTriggerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn, err := meta.(*config.MongoDBClient).GetAppServicesClient(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)
	triggerID := conversion.GetEncodedID(d.Get("trigger_id").(string), "trigger_id")

	eventResp, _, err := conn.EventTriggers.Get(ctx, projectID, appID, triggerID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersRead, projectID, triggerID, err))
	}

	if err = d.Set("name", eventResp.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "name", projectID, appID, err))
	}
	if err = d.Set("type", eventResp.Type); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "type", projectID, appID, err))
	}
	if err = d.Set("function_id", eventResp.FunctionID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "function_id", projectID, appID, err))
	}
	if err = d.Set("function_name", eventResp.FunctionName); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "function_name", projectID, appID, err))
	}
	if err = d.Set("disabled", eventResp.Disabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "disabled", projectID, appID, err))
	}
	if err = d.Set("config_operation_types", eventResp.Config.OperationTypes); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_operation_types", projectID, appID, err))
	}
	if err = d.Set("config_operation_type", eventResp.Config.OperationType); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_operation_type", projectID, appID, err))
	}
	if err = d.Set("config_providers", eventResp.Config.Providers); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_providers", projectID, appID, err))
	}
	if err = d.Set("config_database", eventResp.Config.Database); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_database", projectID, appID, err))
	}
	if err = d.Set("config_collection", eventResp.Config.Collection); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_collection", projectID, appID, err))
	}
	if err = d.Set("config_service_id", eventResp.Config.ServiceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_service_id", projectID, appID, err))
	}
	if err = d.Set("config_match", matchToString(eventResp.Config.Match)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_match", projectID, appID, err))
	}
	if err = d.Set("config_project", matchToString(eventResp.Config.Project)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_project", projectID, appID, err))
	}
	if err = d.Set("config_full_document", eventResp.Config.FullDocument); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_full_document", projectID, appID, err))
	}
	if err = d.Set("config_full_document_before", eventResp.Config.FullDocumentBeforeChange); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_full_document_before", projectID, appID, err))
	}
	if err = d.Set("config_schedule", eventResp.Config.Schedule); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_schedule", projectID, appID, err))
	}
	if err = d.Set("config_schedule_type", eventResp.Config.ScheduleType); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "config_schedule_type", projectID, appID, err))
	}
	if err = d.Set("unordered", eventResp.Config.Unordered); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "unordered", projectID, appID, err))
	}
	if err = d.Set("event_processors", flattenTriggerEventProcessorAWSEventBridge(eventResp.EventProcessors)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEventTriggersSetting, "event_processors", projectID, appID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"app_id":     appID,
		"trigger_id": eventResp.ID,
	}))

	return nil
}
