package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceListOptions() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"include_count": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func readListOptions(listOptionsArr []interface{}) *matlas.ListOptions {
	var listOptions map[string]interface{}

	if len(listOptionsArr) > 0 {
		listOptions = listOptionsArr[0].(map[string]interface{})
	} else {
		listOptions = map[string]interface{}{
			"page_num":       0,
			"items_per_page": 100,
			"include_count":  false,
		}
	}

	return &matlas.ListOptions{
		PageNum:      listOptions["page_num"].(int),
		ItemsPerPage: listOptions["items_per_page"].(int),
		IncludeCount: listOptions["include_count"].(bool),
	}
}

func dataSourceMongoDBAtlasAlertConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAlertConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"list_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceListOptions(),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceMongoDBAtlasAlertConfiguration(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"output_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"resource_hcl", "resource_import"}, false),
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAlertConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	listOptions := d.Get("list_options").([]interface{})

	alerts, _, err := conn.AlertConfigurations.List(ctx, projectID, readListOptions(listOptions))

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorReadAlertConf, err))
	}

	results := flattenAlertConfigurations(ctx, conn, alerts, d)

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "results", projectID, err))
	}

	if err := d.Set("list_options", listOptions); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "list_options", projectID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
	}))

	return nil
}

func flattenAlertConfigurations(ctx context.Context, conn *matlas.Client, alerts []matlas.AlertConfiguration, d *schema.ResourceData) []map[string]interface{} {
	var outputTypes []string

	results := make([]map[string]interface{}, 0)

	if output := d.Get("output_type"); output != nil {
		for _, o := range output.([]interface{}) {
			outputTypes = append(outputTypes, o.(string))
		}
	}

	for _, alert := range alerts {
		results = append(results, map[string]interface{}{
			"alert_configuration_id":  alert.ID,
			"event_type":              alert.EventTypeName,
			"created":                 alert.Created,
			"updated":                 alert.Updated,
			"enabled":                 alert.Enabled,
			"matcher":                 flattenAlertConfigurationMatchers(alert.Matchers),
			"metric_threshold_config": flattenAlertConfigurationMetricThresholdConfig(alert.MetricThreshold),
			"threshold_config":        flattenAlertConfigurationThresholdConfig(alert.Threshold),
			"notification":            flattenAlertConfigurationNotifications(d, alert.Notifications),
			"output":                  computeOutput(&alert, outputTypes),
		})
	}

	return results
}
