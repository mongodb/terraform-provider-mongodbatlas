package mongodbatlas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	"github.com/zclconf/go-cty/cty"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAlertConfigurationRead,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"threshold": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metric_threshold_config": {
				Type:     schema.TypeList,
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
			"threshold_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"notification": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datadog_api_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
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
						"interval_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mobile_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ops_genie_api_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"ops_genie_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"sms_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"team_name": {
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
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"victor_ops_routing_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"roles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"microsoft_teams_webhook_url": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"webhook_secret": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"webhook_url": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
					},
				},
			},
			"output": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"resource_hcl", "resource_import"}, false),
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAlertConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	alertID := getEncodedID(d.Get("alert_configuration_id").(string), "id")

	alert, _, err := conn.AlertConfigurations.GetAnAlertConfig(ctx, projectID, alertID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorReadAlertConf, err))
	}

	if err := d.Set("event_type", alert.EventTypeName); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "event_type", projectID, err))
	}

	if err := d.Set("created", alert.Created); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "created", projectID, err))
	}

	if err := d.Set("updated", alert.Updated); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "updated", projectID, err))
	}

	if err := d.Set("matcher", flattenAlertConfigurationMatchers(alert.Matchers)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "matcher", projectID, err))
	}

	if err := d.Set("metric_threshold", flattenAlertConfigurationMetricThreshold(alert.MetricThreshold)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "metric_threshold", projectID, err))
	}

	if err := d.Set("threshold", flattenAlertConfigurationThreshold(alert.Threshold)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "threshold", projectID, err))
	}

	if err := d.Set("metric_threshold_config", flattenAlertConfigurationMetricThresholdConfig(alert.MetricThreshold)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "metric_threshold_config", projectID, err))
	}

	if err := d.Set("threshold_config", flattenAlertConfigurationThresholdConfig(alert.Threshold)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "threshold_config", projectID, err))
	}

	if err := d.Set("notification", flattenAlertConfigurationNotifications(d, alert.Notifications)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "notification", projectID, err))
	}

	if dOutput := d.Get("output"); dOutput != nil {
		if err := d.Set("output", computeAlertConfigurationOutput(alert, dOutput.([]interface{}), alert.EventTypeName)); err != nil {
			return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "output", projectID, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         alert.ID,
		"project_id": projectID,
	}))

	return nil
}

func computeAlertConfigurationOutput(alert *matlas.AlertConfiguration, outputConfigurations []interface{}, defaultLabel string) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)

	for i := 0; i < len(outputConfigurations); i++ {
		config := outputConfigurations[i].(map[string]interface{})
		var o = map[string]interface{}{
			"type": config["type"],
		}

		if label, ok := o["label"]; ok {
			o["label"] = label
		} else {
			o["label"] = defaultLabel
		}

		if outputValue := outputAlertConfiguration(alert, o["type"].(string), o["label"].(string)); outputValue != "" {
			o["value"] = outputValue
		}

		output = append(output, o)
	}

	return output
}

func outputAlertConfiguration(alert *matlas.AlertConfiguration, outputType, resourceLabel string) string {
	if outputType == "resource_hcl" {
		return outputAlertConfigurationResourceHcl(resourceLabel, alert)
	}
	if outputType == "resource_import" {
		return outputAlertConfigurationResourceImport(resourceLabel, alert)
	}

	return ""
}

func outputAlertConfigurationResourceHcl(label string, alert *matlas.AlertConfiguration) string {
	f := hclwrite.NewEmptyFile()
	root := f.Body()
	resource := root.AppendNewBlock("resource", []string{"mongodbatlas_alert_configuration", label}).Body()

	resource.SetAttributeValue("project_id", cty.StringVal(alert.GroupID))
	resource.SetAttributeValue("event_type", cty.StringVal(alert.EventTypeName))

	if alert.Enabled != nil {
		resource.SetAttributeValue("enabled", cty.BoolVal(*alert.Enabled))
	}

	for _, matcher := range alert.Matchers {
		values := convertMatcherToCtyValues(matcher)

		appendBlockWithCtyValues(resource, "matcher", []string{}, values)
	}

	if alert.MetricThreshold != nil {
		values := convertMetricThresholdToCtyValues(*alert.MetricThreshold)

		appendBlockWithCtyValues(resource, "metric_threshold_config", []string{}, values)
	}

	if alert.Threshold != nil {
		values := convertThresholdToCtyValues(*alert.Threshold)

		appendBlockWithCtyValues(resource, "threshold_config", []string{}, values)
	}

	for i := 0; i < len(alert.Notifications); i++ {
		values := convertNotificationToCtyValues(&alert.Notifications[i])

		appendBlockWithCtyValues(resource, "notification", []string{}, values)
	}

	return string(f.Bytes())
}

func outputAlertConfigurationResourceImport(label string, alert *matlas.AlertConfiguration) string {
	return fmt.Sprintf("terraform import mongodbatlas_alert_configuration.%s %s-%s\n", label, alert.GroupID, alert.ID)
}

func convertMatcherToCtyValues(matcher matlas.Matcher) map[string]cty.Value {
	return map[string]cty.Value{
		"field_name": cty.StringVal(matcher.FieldName),
		"operator":   cty.StringVal(matcher.Operator),
		"value":      cty.StringVal(matcher.Value),
	}
}

func convertMetricThresholdToCtyValues(metric matlas.MetricThreshold) map[string]cty.Value {
	return map[string]cty.Value{
		"metric_name": cty.StringVal(metric.MetricName),
		"operator":    cty.StringVal(metric.Operator),
		"threshold":   cty.NumberFloatVal(metric.Threshold),
		"units":       cty.StringVal(metric.Units),
		"mode":        cty.StringVal(metric.Mode),
	}
}

func convertThresholdToCtyValues(threshold matlas.Threshold) map[string]cty.Value {
	return map[string]cty.Value{
		"operator":  cty.StringVal(threshold.Operator),
		"units":     cty.StringVal(threshold.Units),
		"threshold": cty.NumberFloatVal(threshold.Threshold),
	}
}

func convertNotificationToCtyValues(notification *matlas.Notification) map[string]cty.Value {
	values := map[string]cty.Value{}

	if notification.ChannelName != "" {
		values["channel_name"] = cty.StringVal(notification.ChannelName)
	}

	if notification.DatadogRegion != "" {
		values["datadog_region"] = cty.StringVal(notification.DatadogRegion)
	}

	if notification.EmailAddress != "" {
		values["email_address"] = cty.StringVal(notification.EmailAddress)
	}

	if notification.FlowName != "" {
		values["flow_name"] = cty.StringVal(notification.FlowName)
	}

	if notification.IntervalMin > 0 {
		values["interval_min"] = cty.NumberIntVal(int64(notification.IntervalMin))
	}

	if notification.MobileNumber != "" {
		values["mobile_number"] = cty.StringVal(notification.MobileNumber)
	}

	if notification.OpsGenieRegion != "" {
		values["ops_genie_region"] = cty.StringVal(notification.OpsGenieRegion)
	}

	if notification.OrgName != "" {
		values["org_name"] = cty.StringVal(notification.OrgName)
	}

	if notification.TeamID != "" {
		values["team_id"] = cty.StringVal(notification.TeamID)
	}

	if notification.TeamName != "" {
		values["team_name"] = cty.StringVal(notification.TeamName)
	}

	if notification.TypeName != "" {
		values["type_name"] = cty.StringVal(notification.TypeName)
	}

	if notification.Username != "" {
		values["username"] = cty.StringVal(notification.Username)
	}

	if notification.DelayMin != nil && *notification.DelayMin > 0 {
		values["delay_min"] = cty.NumberIntVal(int64(*notification.DelayMin))
	}

	if notification.EmailEnabled != nil && *notification.EmailEnabled {
		values["email_enabled"] = cty.BoolVal(*notification.EmailEnabled)
	}

	if notification.SMSEnabled != nil && *notification.SMSEnabled {
		values["sms_enabled"] = cty.BoolVal(*notification.SMSEnabled)
	}

	if len(notification.Roles) > 0 {
		roles := make([]cty.Value, 0)

		for _, r := range notification.Roles {
			if r != "" {
				roles = append(roles, cty.StringVal(r))
			}
		}

		values["roles"] = cty.TupleVal(roles)
	}

	return values
}

func flattenAlertConfigurationMetricThreshold(m *matlas.MetricThreshold) map[string]interface{} {
	if m != nil {
		return map[string]interface{}{
			"metric_name": m.MetricName,
			"operator":    m.Operator,
			"threshold":   cast.ToString(m.Threshold),
			"units":       m.Units,
			"mode":        m.Mode,
		}
	}

	return map[string]interface{}{}
}

func flattenAlertConfigurationThreshold(m *matlas.Threshold) map[string]interface{} {
	if m != nil {
		return map[string]interface{}{
			"operator":  m.Operator,
			"units":     m.Units,
			"threshold": cast.ToString(m.Threshold),
		}
	}

	return map[string]interface{}{}
}

func flattenAlertConfigurationMetricThresholdConfig(m *matlas.MetricThreshold) []interface{} {
	if m != nil {
		return []interface{}{map[string]interface{}{
			"metric_name": m.MetricName,
			"operator":    m.Operator,
			"threshold":   m.Threshold,
			"units":       m.Units,
			"mode":        m.Mode,
		}}
	}

	return []interface{}{}
}

func flattenAlertConfigurationThresholdConfig(m *matlas.Threshold) []interface{} {
	if m != nil {
		return []interface{}{map[string]interface{}{
			"operator":  m.Operator,
			"units":     m.Units,
			"threshold": m.Threshold,
		}}
	}

	return []interface{}{}
}

func expandAlertConfigurationNotification(d *schema.ResourceData) ([]matlas.Notification, error) {
	notificationCount := 0

	if notifications, ok := d.GetOk("notification"); ok {
		notificationCount = len(notifications.([]interface{}))
	}

	notifications := make([]matlas.Notification, notificationCount)

	if notificationCount == 0 {
		return notifications, nil
	}

	for i, value := range d.Get("notification").([]interface{}) {
		v := value.(map[string]interface{})
		if v1, ok := v["interval_min"]; ok && v1.(int) > 0 {
			typeName := v["type_name"].(string)
			if strings.EqualFold(typeName, pagerDuty) || strings.EqualFold(typeName, opsGenie) || strings.EqualFold(typeName, victorOps) {
				return nil, fmt.Errorf(`'interval_min' doesn't need to be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'`)
			}
		}
		notifications[i] = matlas.Notification{
			APIToken:                 cast.ToString(v["api_token"]),
			ChannelName:              cast.ToString(v["channel_name"]),
			DatadogAPIKey:            cast.ToString(v["datadog_api_key"]),
			DatadogRegion:            cast.ToString(v["datadog_region"]),
			DelayMin:                 pointy.Int(v["delay_min"].(int)),
			EmailAddress:             cast.ToString(v["email_address"]),
			EmailEnabled:             pointy.Bool(v["email_enabled"].(bool)),
			IntervalMin:              cast.ToInt(v["interval_min"]),
			MobileNumber:             cast.ToString(v["mobile_number"]),
			OpsGenieAPIKey:           cast.ToString(v["ops_genie_api_key"]),
			OpsGenieRegion:           cast.ToString(v["ops_genie_region"]),
			ServiceKey:               cast.ToString(v["service_key"]),
			SMSEnabled:               pointy.Bool(v["sms_enabled"].(bool)),
			TeamID:                   cast.ToString(v["team_id"]),
			TypeName:                 cast.ToString(v["type_name"]),
			Username:                 cast.ToString(v["username"]),
			VictorOpsAPIKey:          cast.ToString(v["victor_ops_api_key"]),
			VictorOpsRoutingKey:      cast.ToString(v["victor_ops_routing_key"]),
			Roles:                    cast.ToStringSlice(v["roles"]),
			MicrosoftTeamsWebhookURL: cast.ToString(v["microsoft_teams_webhook_url"]),
			WebhookSecret:            cast.ToString(v["webhook_secret"]),
			WebhookURL:               cast.ToString(v["webhook_url"]),
		}
	}

	return notifications, nil
}

func flattenAlertConfigurationNotifications(d *schema.ResourceData, notifications []matlas.Notification) []map[string]interface{} {
	notificationsSchema, err := expandAlertConfigurationNotification(d)
	if err != nil {
		return nil
	}

	if len(notificationsSchema) > 0 {
		for i := range notificationsSchema {
			notifications[i].APIToken = notificationsSchema[i].APIToken
			notifications[i].DatadogAPIKey = notificationsSchema[i].DatadogAPIKey
			notifications[i].OpsGenieAPIKey = notificationsSchema[i].OpsGenieAPIKey
			notifications[i].ServiceKey = notificationsSchema[i].ServiceKey
			notifications[i].VictorOpsAPIKey = notificationsSchema[i].VictorOpsAPIKey
			notifications[i].VictorOpsRoutingKey = notificationsSchema[i].VictorOpsRoutingKey
			notifications[i].WebhookURL = notificationsSchema[i].WebhookURL
			notifications[i].WebhookSecret = notificationsSchema[i].WebhookSecret
			notifications[i].SMSEnabled = notificationsSchema[i].SMSEnabled
			notifications[i].EmailEnabled = notificationsSchema[i].EmailEnabled
			notifications[i].MicrosoftTeamsWebhookURL = notificationsSchema[i].MicrosoftTeamsWebhookURL
		}
	}

	nts := make([]map[string]interface{}, len(notifications))

	for i := range notifications {
		nts[i] = map[string]interface{}{
			"api_token":                   notifications[i].APIToken,
			"channel_name":                notifications[i].ChannelName,
			"datadog_api_key":             notifications[i].DatadogAPIKey,
			"datadog_region":              notifications[i].DatadogRegion,
			"delay_min":                   notifications[i].DelayMin,
			"email_address":               notifications[i].EmailAddress,
			"email_enabled":               notifications[i].EmailEnabled,
			"interval_min":                notifications[i].IntervalMin,
			"mobile_number":               notifications[i].MobileNumber,
			"ops_genie_api_key":           notifications[i].OpsGenieAPIKey,
			"ops_genie_region":            notifications[i].OpsGenieRegion,
			"service_key":                 notifications[i].ServiceKey,
			"sms_enabled":                 notifications[i].SMSEnabled,
			"team_id":                     notifications[i].TeamID,
			"team_name":                   notifications[i].TeamName,
			"type_name":                   notifications[i].TypeName,
			"username":                    notifications[i].Username,
			"victor_ops_api_key":          notifications[i].VictorOpsAPIKey,
			"victor_ops_routing_key":      notifications[i].VictorOpsRoutingKey,
			"microsoft_teams_webhook_url": notifications[i].MicrosoftTeamsWebhookURL,
			"webhook_secret":              notifications[i].WebhookSecret,
			"webhook_url":                 notifications[i].WebhookURL,
		}

		// We need to validate it due to the datasource haven't the roles attribute
		if len(notifications[i].Roles) > 0 {
			nts[i]["roles"] = notifications[i].Roles
		}
	}

	return nts
}

func flattenAlertConfigurationMatchers(matchers []matlas.Matcher) []map[string]interface{} {
	mts := make([]map[string]interface{}, len(matchers))

	for i, m := range matchers {
		mts[i] = map[string]interface{}{
			"field_name": m.FieldName,
			"operator":   m.Operator,
			"value":      m.Value,
		}
	}

	return mts
}
