package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
	"github.com/zclconf/go-cty/cty"
	"go.mongodb.org/atlas-sdk/v20230201006/admin"
)

var _ datasource.DataSource = &AlertConfigurationDS{}
var _ datasource.DataSourceWithConfigure = &AlertConfigurationDS{}

type tfAlertConfigurationDSModel struct {
	ID                    types.String                      `tfsdk:"id"`
	ProjectID             types.String                      `tfsdk:"project_id"`
	AlertConfigurationID  types.String                      `tfsdk:"alert_configuration_id"`
	EventType             types.String                      `tfsdk:"event_type"`
	Created               types.String                      `tfsdk:"created"`
	Updated               types.String                      `tfsdk:"updated"`
	Matcher               []tfMatcherModel                  `tfsdk:"matcher"`
	MetricThresholdConfig []tfMetricThresholdConfigModel    `tfsdk:"metric_threshold_config"`
	ThresholdConfig       []tfThresholdConfigModel          `tfsdk:"threshold_config"`
	Notification          []tfNotificationModel             `tfsdk:"notification"`
	Output                []tfAlertConfigurationOutputModel `tfsdk:"output"`
	Enabled               types.Bool                        `tfsdk:"enabled"`
}

type tfAlertConfigurationOutputModel struct {
	Type  types.String `tfsdk:"type"`
	Label types.String `tfsdk:"label"`
	Value types.String `tfsdk:"value"`
}

func NewAlertConfigurationDS() datasource.DataSource {
	return &AlertConfigurationDS{
		DSCommon: DSCommon{
			dataSourceName: alertConfigurationResourceName,
		},
	}
}

type AlertConfigurationDS struct {
	DSCommon
}

var alertConfigDSSchemaBlocks = map[string]schema.Block{
	"output": schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("resource_hcl", "resource_import"),
					},
				},
				"label": schema.StringAttribute{
					Optional: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	},
}

var alertConfigDSSchemaAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Computed: true,
	},
	"project_id": schema.StringAttribute{
		Required: true,
	},
	"alert_configuration_id": schema.StringAttribute{
		Required: true,
	},
	"event_type": schema.StringAttribute{
		Computed: true,
	},
	"created": schema.StringAttribute{
		Computed: true,
	},
	"updated": schema.StringAttribute{
		Computed: true,
	},
	"enabled": schema.BoolAttribute{
		Computed: true,
	},
	"matcher": schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"field_name": schema.StringAttribute{
					Computed: true,
				},
				"operator": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	},
	"metric_threshold_config": schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"metric_name": schema.StringAttribute{
					Computed: true,
				},
				"operator": schema.StringAttribute{
					Computed: true,
				},
				"threshold": schema.Float64Attribute{
					Computed: true,
				},
				"units": schema.StringAttribute{
					Computed: true,
				},
				"mode": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	},
	"threshold_config": schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"operator": schema.StringAttribute{
					Computed: true,
				},
				"threshold": schema.Float64Attribute{
					Computed: true,
				},
				"units": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	},
	"notification": schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"api_token": schema.StringAttribute{
					Computed:  true,
					Sensitive: true,
				},
				"channel_name": schema.StringAttribute{
					Computed: true,
				},
				"datadog_api_key": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"datadog_region": schema.StringAttribute{
					Computed: true,
				},
				"delay_min": schema.Int64Attribute{
					Computed: true,
				},
				"email_address": schema.StringAttribute{
					Computed: true,
				},
				"email_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"interval_min": schema.Int64Attribute{
					Computed: true,
				},
				"mobile_number": schema.StringAttribute{
					Computed: true,
				},
				"ops_genie_api_key": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"ops_genie_region": schema.StringAttribute{
					Computed: true,
				},
				"service_key": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"sms_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"team_id": schema.StringAttribute{
					Computed: true,
				},
				"team_name": schema.StringAttribute{
					Computed: true,
				},
				"notifier_id": schema.StringAttribute{
					Computed: true,
				},
				"type_name": schema.StringAttribute{
					Computed: true,
				},
				"username": schema.StringAttribute{
					Computed: true,
				},
				"victor_ops_api_key": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"victor_ops_routing_key": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"roles": schema.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},
				"microsoft_teams_webhook_url": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"webhook_secret": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
				"webhook_url": schema.StringAttribute{
					Sensitive: true,
					Computed:  true,
				},
			},
		},
	},
}

func (d *AlertConfigurationDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: alertConfigDSSchemaAttributes,
		Blocks:     alertConfigDSSchemaBlocks,
	}
}

func (d *AlertConfigurationDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var alertConfigurationConfig tfAlertConfigurationDSModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &alertConfigurationConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := alertConfigurationConfig.ProjectID.ValueString()

	// this is very hard to follow as the data source currently receieves the alert_configuration resource id in alert_configuration_id attribute
	alertID := getEncodedID(alertConfigurationConfig.AlertConfigurationID.ValueString(), encodedIDKeyAlertID)
	outputs := alertConfigurationConfig.Output

	connV2 := d.client.AtlasV2
	alert, _, err := connV2.AlertConfigurationsApi.GetAlertConfiguration(ctx, projectID, alertID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}

	resultAlertConfigModel := newTFAlertConfigurationDSModel(alert, projectID)
	resultAlertConfigModel.Output = computeAlertConfigurationOutput(alert, outputs, *alert.EventTypeName)

	// setting initial value for backwards compatibility, but setting the alert_configuration resource id here is not consistent with the resource
	resultAlertConfigModel.AlertConfigurationID = alertConfigurationConfig.AlertConfigurationID

	resp.Diagnostics.Append(resp.State.Set(ctx, &resultAlertConfigModel)...)
}

func computeAlertConfigurationOutput(alert *admin.GroupAlertsConfig, definedOutputs []tfAlertConfigurationOutputModel, defaultLabel string) []tfAlertConfigurationOutputModel {
	resultOutputs := make([]tfAlertConfigurationOutputModel, len(definedOutputs))
	for i, defined := range definedOutputs {
		resultOutput := tfAlertConfigurationOutputModel{}
		resultOutput.Type = defined.Type
		if defined.Label.IsNull() {
			resultOutput.Label = types.StringValue(defaultLabel)
		} else {
			resultOutput.Label = defined.Label
		}
		if outputValue := outputAlertConfiguration(alert, resultOutput.Type.ValueString(), resultOutput.Label.ValueString()); outputValue != "" {
			resultOutput.Value = types.StringValue(outputValue)
		}
		resultOutputs[i] = resultOutput
	}
	return resultOutputs
}

func newTFAlertConfigurationDSModel(apiRespConfig *admin.GroupAlertsConfig, projectID string) tfAlertConfigurationDSModel {
	return tfAlertConfigurationDSModel{
		ID: types.StringValue(encodeStateID(map[string]string{
			encodedIDKeyAlertID:   *apiRespConfig.Id,
			encodedIDKeyProjectID: projectID,
		})),
		ProjectID:             types.StringValue(projectID),
		AlertConfigurationID:  types.StringValue(*apiRespConfig.Id),
		EventType:             types.StringValue(*apiRespConfig.EventTypeName),
		Created:               types.StringPointerValue(util.TimePtrToStringPtr(apiRespConfig.Created)),
		Updated:               types.StringPointerValue(util.TimePtrToStringPtr(apiRespConfig.Updated)),
		Enabled:               types.BoolPointerValue(apiRespConfig.Enabled),
		MetricThresholdConfig: newTFMetricThresholdConfigModelV2(apiRespConfig.MetricThreshold, []tfMetricThresholdConfigModel{}),
		ThresholdConfig:       newTFThresholdConfigModelV2(apiRespConfig.Threshold, []tfThresholdConfigModel{}),
		Notification:          newTFNotificationModelListV2(apiRespConfig.Notifications, []tfNotificationModel{}),
		Matcher:               newTFMatcherModelListV2(apiRespConfig.Matchers, []tfMatcherModel{}),
	}
}

func outputAlertConfiguration(alert *admin.GroupAlertsConfig, outputType, resourceLabel string) string {
	if outputType == "resource_hcl" {
		return outputAlertConfigurationResourceHcl(resourceLabel, alert)
	}
	if outputType == "resource_import" {
		return outputAlertConfigurationResourceImport(resourceLabel, alert)
	}

	return ""
}

func outputAlertConfigurationResourceHcl(label string, alert *admin.GroupAlertsConfig) string {
	f := hclwrite.NewEmptyFile()
	root := f.Body()
	resource := root.AppendNewBlock("resource", []string{"mongodbatlas_alert_configuration", label}).Body()

	resource.SetAttributeValue("project_id", cty.StringVal(*alert.GroupId))
	resource.SetAttributeValue("event_type", cty.StringVal(*alert.EventTypeName))

	if alert.Enabled != nil {
		resource.SetAttributeValue("enabled", cty.BoolVal(*alert.Enabled))
	}

	for _, matcher := range alert.Matchers {
		appendBlockWithCtyValues(resource, "matcher", []string{}, convertMatcherToCtyValues(matcher))
	}

	if alert.MetricThreshold != nil {
		appendBlockWithCtyValues(resource, "metric_threshold_config", []string{}, convertMetricThresholdToCtyValues(*alert.MetricThreshold))
	}

	if alert.Threshold != nil {
		appendBlockWithCtyValues(resource, "threshold_config", []string{}, convertThresholdToCtyValues(alert.Threshold))
	}

	for i := 0; i < len(alert.Notifications); i++ {
		appendBlockWithCtyValues(resource, "notification", []string{}, convertNotificationToCtyValues(&alert.Notifications[i]))
	}

	return string(f.Bytes())
}

func outputAlertConfigurationResourceImport(label string, alert *admin.GroupAlertsConfig) string {
	return fmt.Sprintf("terraform import mongodbatlas_alert_configuration.%s %s-%s\n", label, *alert.GroupId, *alert.Id)
}

func convertMatcherToCtyValues(matcher map[string]interface{}) map[string]cty.Value {
	fieldName, _ := matcher["fieldName"].(string)
	operator, _ := matcher["operator"].(string)
	value, _ := matcher["value"].(string)
	return map[string]cty.Value{
		"field_name": cty.StringVal(fieldName),
		"operator":   cty.StringVal(operator),
		"value":      cty.StringVal(value),
	}
}

func convertMetricThresholdToCtyValues(metric admin.ServerlessMetricThreshold) map[string]cty.Value {
	var t float64
	if metric.Threshold != nil {
		t = *metric.Threshold
	}
	return map[string]cty.Value{
		"metric_name": ctyStringPtrVal(metric.MetricName),
		"operator":    ctyStringPtrVal(metric.Operator),
		"threshold":   cty.NumberFloatVal(t),
		"units":       ctyStringPtrVal(metric.Units),
		"mode":        ctyStringPtrVal(metric.Mode),
	}
}

func convertThresholdToCtyValues(threshold *admin.GreaterThanRawThreshold) map[string]cty.Value {
	var t int
	if threshold.Threshold != nil {
		t = *threshold.Threshold
	}
	return map[string]cty.Value{
		"operator":  ctyStringPtrVal(threshold.Operator),
		"units":     ctyStringPtrVal(threshold.Units),
		"threshold": cty.NumberFloatVal(float64(t)), // int in new SDK but keeping float64 for backward compatibility
	}
}

func convertNotificationToCtyValues(notification *admin.AlertsNotificationRootForGroup) map[string]cty.Value {
	values := map[string]cty.Value{}

	if util.IsStringPresent(notification.ChannelName) {
		values["channel_name"] = cty.StringVal(*notification.ChannelName)
	}

	if util.IsStringPresent(notification.DatadogRegion) {
		values["datadog_region"] = cty.StringVal(*notification.DatadogRegion)
	}

	if util.IsStringPresent(notification.EmailAddress) {
		values["email_address"] = cty.StringVal(*notification.EmailAddress)
	}

	if notification.IntervalMin != nil && *notification.IntervalMin > 0 {
		values["interval_min"] = cty.NumberIntVal(int64(*notification.IntervalMin))
	}

	if util.IsStringPresent(notification.MobileNumber) {
		values["mobile_number"] = cty.StringVal(*notification.MobileNumber)
	}

	if util.IsStringPresent(notification.OpsGenieRegion) {
		values["ops_genie_region"] = cty.StringVal(*notification.OpsGenieRegion)
	}

	if util.IsStringPresent(notification.TeamId) {
		values["team_id"] = cty.StringVal(*notification.TeamId)
	}

	if util.IsStringPresent(notification.TeamName) {
		values["team_name"] = cty.StringVal(*notification.TeamName)
	}

	if util.IsStringPresent(notification.NotifierId) {
		values["notifier_id"] = cty.StringVal(*notification.NotifierId)
	}

	if util.IsStringPresent(notification.TypeName) {
		values["type_name"] = cty.StringVal(*notification.TypeName)
	}

	if util.IsStringPresent(notification.Username) {
		values["username"] = cty.StringVal(*notification.Username)
	}

	if notification.DelayMin != nil && *notification.DelayMin > 0 {
		values["delay_min"] = cty.NumberIntVal(int64(*notification.DelayMin))
	}

	if notification.EmailEnabled != nil && *notification.EmailEnabled {
		values["email_enabled"] = cty.BoolVal(*notification.EmailEnabled)
	}

	if notification.SmsEnabled != nil && *notification.SmsEnabled {
		values["sms_enabled"] = cty.BoolVal(*notification.SmsEnabled)
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

func ctyStringPtrVal(ptr *string) cty.Value {
	if ptr == nil {
		return cty.StringVal("")
	}
	return cty.StringVal(*ptr)
}
