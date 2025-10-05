package eventtrigger

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/realm/realm"
)

func PluralDataSource() *schema.Resource {
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasEventTriggersRead(d *schema.ResourceData, meta any) error {
	// Get client connection.
	ctx := context.Background()
	conn, err := meta.(*config.MongoDBClient).Realm.Get(ctx)
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	appID := d.Get("app_id").(string)

	eventTriggers, _, err := conn.EventTriggers.List(ctx, projectID, appID)
	if err != nil {
		return fmt.Errorf("error getting event triggers information: %s", err)
	}

	if err := d.Set("results", flattenEventTriggers(eventTriggers)); err != nil {
		return fmt.Errorf("error setting `result` for event triggers: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenEventTriggers(eventTriggers []realm.EventTrigger) []map[string]any {
	var triggersMap []map[string]any

	if len(eventTriggers) > 0 {
		triggersMap = make([]map[string]any, len(eventTriggers))

		for i := range eventTriggers {
			triggersMap[i] = map[string]any{
				"trigger_id":                  eventTriggers[i].ID,
				"name":                        eventTriggers[i].Name,
				"type":                        eventTriggers[i].Type,
				"function_id":                 eventTriggers[i].FunctionID,
				"function_name":               eventTriggers[i].FunctionName,
				"disabled":                    eventTriggers[i].Disabled,
				"config_operation_types":      eventTriggers[i].Config.OperationTypes,
				"config_operation_type":       eventTriggers[i].Config.OperationType,
				"config_providers":            eventTriggers[i].Config.Providers,
				"config_database":             eventTriggers[i].Config.Database,
				"config_collection":           eventTriggers[i].Config.Collection,
				"config_service_id":           eventTriggers[i].Config.ServiceID,
				"config_match":                matchToString(eventTriggers[i].Config.Match),
				"config_project":              matchToString(eventTriggers[i].Config.Project),
				"config_full_document":        eventTriggers[i].Config.FullDocument,
				"config_full_document_before": eventTriggers[i].Config.FullDocumentBeforeChange,
				"config_schedule":             eventTriggers[i].Config.Schedule,
				"config_schedule_type":        eventTriggers[i].Config.ScheduleType,
				"event_processors":            flattenTriggerEventProcessorAWSEventBridge(eventTriggers[i].EventProcessors),
				"unordered":                   eventTriggers[i].Config.Unordered,
			}
		}
	}

	return triggersMap
}
