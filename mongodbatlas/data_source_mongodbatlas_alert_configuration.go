package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasAlertConfigurationRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_configuration_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"matcher": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metric_threshold": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"units": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"notification": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datadog_api_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datadog_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flowdock_api_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mobile_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ops_genie_api_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ops_genie_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sms_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"victor_ops_api_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"victor_ops_routing_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAlertConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	alertID := d.Get("alert_configuration_id").(string)

	alert, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), projectID, alertID)
	if err != nil {
		return fmt.Errorf(errorReadAlertConf, err)
	}

	if err := d.Set("group_id", alert.GroupID); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "group_id", projectID, err)
	}
	if err := d.Set("event_type", alert.EventTypeName); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "event_type", projectID, err)
	}
	if err := d.Set("created", alert.Created); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "created", projectID, err)
	}
	if err := d.Set("updated", alert.Updated); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "updated", projectID, err)
	}
	if err := d.Set("matcher", flattenAlertConfigurationMatchers(alert.Matchers)); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "matcher", projectID, err)
	}
	if err := d.Set("metric_threshold", flattenAlertConfigurationMetricThreshold(alert.MetricThreshold)); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "metric_threshold", projectID, err)
	}
	if err := d.Set("notification", flattenAlertConfigurationNotifications(alert.Notifications)); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "notification", projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         alert.ID,
		"project_id": projectID,
	}))

	return nil
}
