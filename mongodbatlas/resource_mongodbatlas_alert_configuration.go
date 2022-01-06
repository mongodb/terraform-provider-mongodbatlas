package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorCreateAlertConf  = "error creating Alert Configuration information: %s"
	errorReadAlertConf    = "error getting Alert Configuration information: %s"
	errorUpdateAlertConf  = "error updating Alert Configuration information: %s"
	errorDeleteAlertConf  = "error deleting Alert Configuration information: %s"
	errorAlertConfSetting = "error setting `%s` for Alert Configuration (%s): %s"
	errorImportAlertConf  = "couldn't import Alert Configuration (%s) in project %s, error: %s"
	pagerDuty             = "PAGER_DUTY"
	opsGenie              = "OPS_GENIE"
	victorOps             = "VICTOR_OPS"
)

func resourceMongoDBAtlasAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasAlertConfigurationCreate,
		ReadContext:   resourceMongoDBAtlasAlertConfigurationRead,
		UpdateContext: resourceMongoDBAtlasAlertConfigurationUpdate,
		DeleteContext: resourceMongoDBAtlasAlertConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasAlertConfigurationImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Required: true,
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
				Optional: true,
			},
			"matcher": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operator": {
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
			"metric_threshold": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"metric_threshold_config"},
				Deprecated:    "use metric_threshold_config instead",
			},
			"threshold": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"threshold_config"},
				Deprecated:    "use threshold_config instead",
			},
			"metric_threshold_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"metric_threshold"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"GREATER_THAN", "LESS_THAN"}, false),
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"units": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RAW",
								"BITS",
								"BYTES",
								"KILOBITS",
								"KILOBYTES",
								"MEGABITS",
								"MEGABYTES",
								"GIGABITS",
								"GIGABYTES",
								"TERABYTES",
								"PETABYTES",
								"MILLISECONDS",
								"SECONDS",
								"MINUTES",
								"HOURS",
								"DAYS"}, false),
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"threshold_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"threshold"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"units": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RAW",
								"BITS",
								"BYTES",
								"KILOBITS",
								"KILOBYTES",
								"MEGABITS",
								"MEGABYTES",
								"GIGABITS",
								"GIGABYTES",
								"TERABYTES",
								"PETABYTES",
								"MILLISECONDS",
								"SECONDS",
								"MINUTES",
								"HOURS",
								"DAYS"}, false),
						},
					},
				},
			},
			"notification": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"datadog_api_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"datadog_region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
						},
						"delay_min": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"flowdock_api_token": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"flow_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval_min": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"mobile_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ops_genie_api_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"ops_genie_region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
						},
						"org_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"sms_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"team_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{"EMAIL", "SMS", pagerDuty, "SLACK",
								"FLOWDOCK", "DATADOG", opsGenie, victorOps,
								"WEBHOOK", "USER", "TEAM", "GROUP", "ORG"}, false),
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"victor_ops_api_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"victor_ops_routing_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"roles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasAlertConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)

	req := &matlas.AlertConfiguration{
		EventTypeName:   d.Get("event_type").(string),
		Enabled:         pointy.Bool(d.Get("enabled").(bool)),
		Matchers:        expandAlertConfigurationMatchers(d),
		MetricThreshold: expandAlertConfigurationMetricThresholdConfig(d),
		Threshold:       expandAlertConfigurationThresholdConfig(d),
	}

	notifications, err := expandAlertConfigurationNotification(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreateAlertConf, err))
	}
	req.Notifications = notifications

	resp, _, err := conn.AlertConfigurations.Create(ctx, projectID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreateAlertConf, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         resp.ID,
		"project_id": projectID,
	}))

	return resourceMongoDBAtlasAlertConfigurationRead(ctx, d, meta)
}

func resourceMongoDBAtlasAlertConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	alert, resp, err := conn.AlertConfigurations.GetAnAlertConfig(context.Background(), ids["project_id"], ids["id"])
	if err != nil {
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorReadAlertConf, err))
	}

	if err := d.Set("alert_configuration_id", alert.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "alert_configuration_id", ids["id"], err))
	}

	if err := d.Set("event_type", alert.EventTypeName); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "event_type", ids["id"], err))
	}

	if err := d.Set("created", alert.Created); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "created", ids["id"], err))
	}

	if err := d.Set("updated", alert.Updated); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "updated", ids["id"], err))
	}

	if err := d.Set("notification", flattenAlertConfigurationNotifications(alert.Notifications)); err != nil {
		return diag.FromErr(fmt.Errorf(errorAlertConfSetting, "notification", ids["id"], err))
	}

	return nil
}

func resourceMongoDBAtlasAlertConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conn = meta.(*MongoDBClient).Atlas
		ids  = decodeStateID(d.Id())
		err  error
	)

	// In order to update an alert config it is necessary to send the original alert configuration request again, if not the
	// server returns an error 500
	req, _, err := conn.AlertConfigurations.GetAnAlertConfig(ctx, ids["project_id"], ids["id"])
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorReadAlertConf, err))
	}
	// Removing the computed attributes to recreate the original request
	req.GroupID = ""
	req.Created = ""
	req.Updated = ""

	// If matchers changes ensure we are sending the information like the Terraform state
	// because the req variable above doesn't have the "field_name" in each matcher item attribute
	// if sent as is, the server sends an error
	req.Matchers = expandAlertConfigurationMatchers(d)

	// Only changes the updated fields
	if d.HasChange("enabled") {
		req.Enabled = pointy.Bool(d.Get("enabled").(bool))
	}

	if d.HasChange("event_type") {
		req.EventTypeName = d.Get("event_type").(string)
	}

	if d.HasChange("metric_threshold") {
		req.MetricThreshold = expandAlertConfigurationMetricThreshold(d)
	}

	if d.HasChange("threshold") {
		req.Threshold = expandAlertConfigurationThreshold(d)
	}

	if d.HasChange("metric_threshold_config") {
		req.MetricThreshold = expandAlertConfigurationMetricThresholdConfig(d)
	}

	if d.HasChange("threshold_config") {
		req.Threshold = expandAlertConfigurationThresholdConfig(d)
	}

	if d.HasChange("notification") {
		notifications, err := expandAlertConfigurationNotification(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorUpdateAlertConf, err))
		}
		req.Notifications = notifications
	}

	// Cannot enable/disable ONLY via update (if only send enable as changed field server returns a 500 error) so have to use different method to change enabled.
	if reflect.DeepEqual(req, &matlas.AlertConfiguration{Enabled: pointy.Bool(true)}) ||
		reflect.DeepEqual(req, &matlas.AlertConfiguration{Enabled: pointy.Bool(false)}) {
		_, _, err = conn.AlertConfigurations.EnableAnAlertConfig(ctx, ids["project_id"], ids["id"], req.Enabled)
	} else {
		_, _, err = conn.AlertConfigurations.Update(ctx, ids["project_id"], ids["id"], req)
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorUpdateAlertConf, err))
	}

	return resourceMongoDBAtlasAlertConfigurationRead(ctx, d, meta)
}

func resourceMongoDBAtlasAlertConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())

	_, err := conn.AlertConfigurations.Delete(ctx, ids["project_id"], ids["id"])
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDeleteAlertConf, err))
	}

	return nil
}

func resourceMongoDBAtlasAlertConfigurationImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	parts := strings.SplitN(d.Id(), "-", 2)

	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a cluster, use the format {project_id}-{id}")
	}

	projectID := parts[0]
	id := parts[1]

	alert, _, err := conn.AlertConfigurations.GetAnAlertConfig(ctx, projectID, id)
	if err != nil {
		return nil, fmt.Errorf(errorImportAlertConf, id, projectID, err)
	}

	if err := d.Set("project_id", alert.GroupID); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "project_id", id, err)
	}

	if err := d.Set("event_type", alert.EventTypeName); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "event_type", id, err)
	}

	if err := d.Set("enabled", alert.Enabled); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "enabled", id, err)
	}

	if err := d.Set("matcher", flattenAlertConfigurationMatchers(alert.Matchers)); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "matcher", id, err)
	}

	if err := d.Set("metric_threshold_config", flattenAlertConfigurationMetricThresholdConfig(alert.MetricThreshold)); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "metric_threshold_config", id, err)
	}

	if err := d.Set("threshold_config", flattenAlertConfigurationThresholdConfig(alert.Threshold)); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "threshold_config", id, err)
	}

	if err := d.Set("notification", flattenAlertConfigurationNotifications(alert.Notifications)); err != nil {
		return nil, fmt.Errorf(errorAlertConfSetting, "notification", id, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         alert.ID,
		"project_id": projectID,
	}))

	return []*schema.ResourceData{d}, nil
}

func expandAlertConfigurationMatchers(d *schema.ResourceData) []matlas.Matcher {
	matchers := make([]matlas.Matcher, 0)

	if m, ok := d.GetOk("matcher"); ok {
		for _, value := range m.([]interface{}) {
			v := value.(map[string]interface{})

			matchers = append(matchers, matlas.Matcher{
				FieldName: v["field_name"].(string),
				Operator:  v["operator"].(string),
				Value:     v["value"].(string),
			})
		}
	}

	return matchers
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

func expandAlertConfigurationMetricThreshold(d *schema.ResourceData) *matlas.MetricThreshold {
	if value, ok := d.GetOk("metric_threshold"); ok {
		v := value.(map[string]interface{})

		return &matlas.MetricThreshold{
			MetricName: cast.ToString(v["metric_name"]),
			Operator:   cast.ToString(v["operator"]),
			Threshold:  cast.ToFloat64(v["threshold"]),
			Units:      cast.ToString(v["units"]),
			Mode:       cast.ToString(v["mode"]),
		}
	}

	return nil
}

func expandAlertConfigurationThreshold(d *schema.ResourceData) *matlas.Threshold {
	if value, ok := d.GetOk("threshold"); ok {
		v := value.(map[string]interface{})

		return &matlas.Threshold{
			Operator:  cast.ToString(v["operator"]),
			Units:     cast.ToString(v["units"]),
			Threshold: cast.ToFloat64(v["threshold"]),
		}
	}

	return nil
}

func expandAlertConfigurationMetricThresholdConfig(d *schema.ResourceData) *matlas.MetricThreshold {
	if value, ok := d.GetOk("metric_threshold_config"); ok {
		vL := value.([]interface{})

		if len(vL) > 0 {
			v := vL[0].(map[string]interface{})

			return &matlas.MetricThreshold{
				MetricName: cast.ToString(v["metric_name"]),
				Operator:   cast.ToString(v["operator"]),
				Threshold:  cast.ToFloat64(v["threshold"]),
				Units:      cast.ToString(v["units"]),
				Mode:       cast.ToString(v["mode"]),
			}
		}
	}

	// Deprecated, will be removed later
	if value, ok := d.GetOk("metric_threshold"); ok {
		v := value.(map[string]interface{})

		return &matlas.MetricThreshold{
			MetricName: cast.ToString(v["metric_name"]),
			Operator:   cast.ToString(v["operator"]),
			Threshold:  cast.ToFloat64(v["threshold"]),
			Units:      cast.ToString(v["units"]),
			Mode:       cast.ToString(v["mode"]),
		}
	}

	return nil
}

func expandAlertConfigurationThresholdConfig(d *schema.ResourceData) *matlas.Threshold {
	if value, ok := d.GetOk("threshold_config"); ok {
		vL := value.([]interface{})

		if len(vL) > 0 {
			v := vL[0].(map[string]interface{})

			return &matlas.Threshold{
				Operator:  cast.ToString(v["operator"]),
				Units:     cast.ToString(v["units"]),
				Threshold: cast.ToFloat64(v["threshold"]),
			}
		}
	}

	// Deprecated, will be removed later
	if value, ok := d.GetOk("threshold"); ok {
		v := value.(map[string]interface{})

		return &matlas.Threshold{
			Operator:  cast.ToString(v["operator"]),
			Units:     cast.ToString(v["units"]),
			Threshold: cast.ToFloat64(v["threshold"]),
		}
	}

	return nil
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
	notifications := make([]matlas.Notification, len(d.Get("notification").([]interface{})))

	for i, value := range d.Get("notification").([]interface{}) {
		v := value.(map[string]interface{})
		if v1, ok := v["interval_min"]; ok && v1.(int) > 0 {
			typeName := v["type_name"].(string)
			if strings.EqualFold(typeName, pagerDuty) || strings.EqualFold(typeName, opsGenie) || strings.EqualFold(typeName, victorOps) {
				return nil, fmt.Errorf(`'interval_min' doesn't need to be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'`)
			}
		}
		notifications[i] = matlas.Notification{
			APIToken:            cast.ToString(v["api_token"]),
			ChannelName:         cast.ToString(v["channel_name"]),
			DatadogAPIKey:       cast.ToString(v["datadog_api_key"]),
			DatadogRegion:       cast.ToString(v["datadog_region"]),
			DelayMin:            pointy.Int(v["delay_min"].(int)),
			EmailAddress:        cast.ToString(v["email_address"]),
			EmailEnabled:        pointy.Bool(v["email_enabled"].(bool)),
			FlowdockAPIToken:    cast.ToString(v["flowdock_api_token"]),
			FlowName:            cast.ToString(v["flow_name"]),
			IntervalMin:         cast.ToInt(v["interval_min"]),
			MobileNumber:        cast.ToString(v["mobile_number"]),
			OpsGenieAPIKey:      cast.ToString(v["ops_genie_api_key"]),
			OpsGenieRegion:      cast.ToString(v["ops_genie_region"]),
			OrgName:             cast.ToString(v["org_name"]),
			ServiceKey:          cast.ToString(v["service_key"]),
			SMSEnabled:          pointy.Bool(v["sms_enabled"].(bool)),
			TeamID:              cast.ToString(v["team_id"]),
			TypeName:            cast.ToString(v["type_name"]),
			Username:            cast.ToString(v["username"]),
			VictorOpsAPIKey:     cast.ToString(v["victor_ops_api_key"]),
			VictorOpsRoutingKey: cast.ToString(v["victor_ops_routing_key"]),
			Roles:               cast.ToStringSlice(v["roles"]),
		}
	}

	return notifications, nil
}

func flattenAlertConfigurationNotifications(notifications []matlas.Notification) []map[string]interface{} {
	nts := make([]map[string]interface{}, len(notifications))

	for i := range notifications {
		nts[i] = map[string]interface{}{
			"api_token":              notifications[i].APIToken,
			"channel_name":           notifications[i].ChannelName,
			"datadog_api_key":        notifications[i].DatadogAPIKey,
			"datadog_region":         notifications[i].DatadogRegion,
			"delay_min":              notifications[i].DelayMin,
			"email_address":          notifications[i].EmailAddress,
			"email_enabled":          notifications[i].EmailEnabled,
			"flowdock_api_token":     notifications[i].FlowdockAPIToken,
			"flow_name":              notifications[i].FlowName,
			"interval_min":           notifications[i].IntervalMin,
			"mobile_number":          notifications[i].MobileNumber,
			"ops_genie_api_key":      notifications[i].OpsGenieAPIKey,
			"ops_genie_region":       notifications[i].OpsGenieRegion,
			"org_name":               notifications[i].OrgName,
			"service_key":            notifications[i].ServiceKey,
			"sms_enabled":            notifications[i].SMSEnabled,
			"team_id":                notifications[i].TeamID,
			"team_name":              notifications[i].TeamName,
			"type_name":              notifications[i].TypeName,
			"username":               notifications[i].Username,
			"victor_ops_api_key":     notifications[i].VictorOpsAPIKey,
			"victor_ops_routing_key": notifications[i].VictorOpsRoutingKey,
		}

		// We need to validate it due to the datasource haven't the roles attribute
		if len(notifications[i].Roles) > 0 {
			nts[i]["roles"] = notifications[i].Roles
		}
	}

	return nts
}
