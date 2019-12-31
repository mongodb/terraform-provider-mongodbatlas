package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorCreateAlertConf  = "error creating Alert Configuration information: %s"
	errorReadAlertConf    = "error getting Alert Configuration information: %s"
	errorDeleteAlertConf  = "error deleting Alert Configuration information: %s"
	errorAlertConfSetting = "error setting `%s` for Alert Configuration (%s): %s"
	errorImportAlertConf  = "couldn't import Alert Configuration (%s) in project %s, error: %s"
)

func resourceMongoDBAtlasAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasAlertConfigurationCreate,
		Read:   resourceMongoDBAtlasAlertConfigurationRead,
		Update: resourceMongoDBAtlasAlertConfigurationUpdate,
		Delete: resourceMongoDBAtlasAlertConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasAlertConfigurationImportState,
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
			"group_id": {
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
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"TYPE_NAME", "HOSTNAME", "PORT", "HOSTNAME_AND_PORT", "REPLICA_SET_NAME", "REPLICA_SET_NAME", "SHARD_NAME", "CLUSTER_NAME", "CLUSTER_NAME", "SHARD_NAME"}, false),
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"EQUALS", "NOT_EQUALS", "CONTAINS", "NOT_CONTAINS", "STARTS_WITH", "ENDS_WITH", "REGEX"}, false),
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"PRIMARY", "SECONDARY", "STANDALONE", "CONFIG", "MONGOS"}, false),
						},
					},
				},
			},
			"metric_threshold": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"TYPE_NAME",
								"HOSTNAME",
								"PORT",
								"HOSTNAME_AND_PORT",
								"REPLICA_SET_NAME",
								"REPLICA_SET_NAME",
								"SHARD_NAME",
								"CLUSTER_NAME",
								"CLUSTER_NAME",
								"SHARD_NAME"}, false),
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
			"notification": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"datadog_api_key": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Optional: true,
						},
						"sms_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"DATADOG",
								"EMAIL",
								"FLOWDOCK",
								"GROUP",
								"OPS_GENIE",
								"ORG",
								"PAGER_DUTY",
								"SLACK",
								"SMS",
								"TEAM",
								"USER",
								"VICTOR_OPS",
								"WEBHOOK"}, false),
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"victor_ops_api_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"victor_ops_routing_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasAlertConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)

	req := &matlas.AlertConfiguration{
		EventTypeName:   d.Get("event_type").(string),
		Enabled:         pointy.Bool(d.Get("enabled").(bool)),
		Matchers:        expandAlertConfigurationMatchers(d),
		MetricThreshold: expandAlertConfigurationMetricThreshold(d),
		Notifications:   expandAlertConfigurationNotification(d),
	}

	resp, _, err := conn.AlertConfigurations.Create(context.Background(), projectID, req)
	if err != nil {
		return fmt.Errorf(errorCreateAlertConf, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"id":         resp.ID,
		"project_id": projectID,
	}))

	return resourceMongoDBAtlasAlertConfigurationRead(d, meta)
}

func resourceMongoDBAtlasAlertConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	alert, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), ids["project_id"], ids["id"])
	if err != nil {
		return fmt.Errorf(errorReadAlertConf, err)
	}

	if err := d.Set("alert_configuration_id", alert.ID); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "alert_configuration_id", ids["id"], err)
	}
	if err := d.Set("group_id", alert.GroupID); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "group_id", ids["id"], err)
	}
	if err := d.Set("created", alert.Created); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "created", ids["id"], err)
	}
	if err := d.Set("updated", alert.Updated); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "updated", ids["id"], err)
	}
	if err := d.Set("notification", flattenAlertConfigurationNotifications(alert.Notifications)); err != nil {
		return fmt.Errorf(errorAlertConfSetting, "notification", ids["id"], err)
	}
	return nil
}

func resourceMongoDBAtlasAlertConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	var err error

	// TO update is nesessary to send the original create alert configuration request if not the server returns
	// error 500
	req, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), ids["project_id"], ids["id"])
	if err != nil {
		return fmt.Errorf(errorReadAlertConf, err)
	}
	// Removing the computed attributest to create the original create request
	req.GroupID = ""
	req.Created = ""
	req.Updated = ""

	// If something change, always the matches need to be sent like
	// the Terraform state to due that the req variable got above doesn't have
	// the "field_name" in each matcher item attribute so if the request is sent like
	// above we will get an error
	req.Matchers = expandAlertConfigurationMatchers(d)

	// Only changes the updated fields
	if d.HasChange("enabled") {
		req.Enabled = pointy.Bool(d.Get("enabled").(bool))
	}
	if d.HasChange("event_type_name") {
		req.EventTypeName = d.Get("event_type_name").(string)
	}
	if d.HasChange("metric_threshold") {
		req.MetricThreshold = expandAlertConfigurationMetricThreshold(d)
	}
	if d.HasChange("notification") {
		req.Notifications = expandAlertConfigurationNotification(d)
	}

	// it's necessary if just enabled attr is seated, due if only the enabled attr is sent the server returns error 500,
	// so we need to use the enabled/disabled entry point to set it
	if reflect.DeepEqual(req, &matlas.AlertConfiguration{Enabled: pointy.Bool(true)}) ||
		reflect.DeepEqual(req, &matlas.AlertConfiguration{Enabled: pointy.Bool(false)}) {
		_, _, err = conn.AlertConfigurations.EnableAnAlert(context.Background(), ids["project_id"], ids["id"], req.Enabled)
	} else {
		_, _, err = conn.AlertConfigurations.Update(context.Background(), ids["project_id"], ids["id"], req)
	}
	if err != nil {
		return fmt.Errorf(errorReadAlertConf, err)
	}

	return resourceMongoDBAtlasAlertConfigurationRead(d, meta)
}

func resourceMongoDBAtlasAlertConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())

	_, err := conn.AlertConfigurations.Delete(context.Background(), ids["project_id"], ids["id"])
	if err != nil {
		return fmt.Errorf(errorDeleteAlertConf, err)
	}
	return nil
}

func resourceMongoDBAtlasAlertConfigurationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a cluster, use the format {project_id}-{id}")
	}

	projectID := parts[0]
	id := parts[1]

	alert, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), projectID, id)
	if err != nil {
		return nil, fmt.Errorf(errorImportAlertConf, id, projectID, err)
	}

	if err := d.Set("project_id", alert.GroupID); err != nil {
		log.Printf(errorAlertConfSetting, "project_id", id, err)
	}
	if err := d.Set("event_type", alert.EventTypeName); err != nil {
		log.Printf(errorAlertConfSetting, "event_type", id, err)
	}
	if err := d.Set("enabled", alert.Enabled); err != nil {
		log.Printf(errorAlertConfSetting, "enabled", id, err)
	}
	if err := d.Set("notification", flattenAlertConfigurationNotifications(alert.Notifications)); err != nil {
		log.Printf(errorAlertConfSetting, "notification", id, err)
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
	return &matlas.MetricThreshold{}
}

func flattenAlertConfigurationMetricThreshold(m *matlas.MetricThreshold) map[string]interface{} {
	return map[string]interface{}{
		"metric_name": m.MetricName,
		"operator":    m.Operator,
		"threshold":   cast.ToString(m.Threshold),
		"units":       m.Units,
		"mode":        m.Mode,
	}
}

func expandAlertConfigurationNotification(d *schema.ResourceData) []matlas.Notification {
	notifications := make([]matlas.Notification, len(d.Get("notification").([]interface{})))
	for i, value := range d.Get("notification").([]interface{}) {
		v := value.(map[string]interface{})
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
		}
	}
	return notifications
}

func flattenAlertConfigurationNotifications(notifications []matlas.Notification) []map[string]interface{} {
	nts := make([]map[string]interface{}, len(notifications))

	for i, n := range notifications {
		nts[i] = map[string]interface{}{
			"api_token":              n.APIToken,
			"channel_name":           n.ChannelName,
			"datadog_api_key":        n.DatadogRegion,
			"datadog_region":         n.DatadogRegion,
			"delay_min":              n.DelayMin,
			"email_address":          n.EmailAddress,
			"email_enabled":          n.EmailEnabled,
			"flowdock_api_token":     n.FlowdockAPIToken,
			"flow_name":              n.FlowName,
			"interval_min":           n.IntervalMin,
			"mobile_number":          n.MobileNumber,
			"ops_genie_api_key":      n.OpsGenieAPIKey,
			"ops_genie_region":       n.OpsGenieRegion,
			"org_name":               n.OrgName,
			"service_key":            n.ServiceKey,
			"sms_enabled":            n.SMSEnabled,
			"team_id":                n.TeamID,
			"type_name":              n.TypeName,
			"username":               n.Username,
			"victor_ops_api_key":     n.VictorOpsAPIKey,
			"victor_ops_routing_key": n.VictorOpsRoutingKey,
		}
	}
	return nts
}
