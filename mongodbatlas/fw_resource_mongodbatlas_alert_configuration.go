package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	conversion "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/conversion"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	alertConfigurationResourceName = "alert_configuration"
	errorCreateAlertConf           = "error creating Alert Configuration information: %s"
	errorReadAlertConf             = "error getting Alert Configuration information: %s"
	errorUpdateAlertConf           = "error updating Alert Configuration information: %s"
	errorAlertConfSetting          = "error setting `%s` for Alert Configuration (%s): %s"
	pagerDuty                      = "PAGER_DUTY"
	opsGenie                       = "OPS_GENIE"
	victorOps                      = "VICTOR_OPS"
	encodedIDKeyAlertID            = "id"
	encodedIDKeyProjectID          = "project_id"
)

var _ resource.Resource = &AlertConfigurationRS{}
var _ resource.ResourceWithImportState = &AlertConfigurationRS{}

func NewAlertConfigurationRS() resource.Resource {
	return &AlertConfigurationRS{}
}

type AlertConfigurationRS struct {
	client *MongoDBClient
}

type tfAlertConfigurationRSModel struct {
	ID                    types.String                   `tfsdk:"id"`
	ProjectID             types.String                   `tfsdk:"project_id"`
	AlertConfigurationID  types.String                   `tfsdk:"alert_configuration_id"`
	EventType             types.String                   `tfsdk:"event_type"`
	Created               types.String                   `tfsdk:"created"`
	Updated               types.String                   `tfsdk:"updated"`
	Matcher               []tfMatcherModel               `tfsdk:"matcher"`
	MetricThresholdConfig []tfMetricThresholdConfigModel `tfsdk:"metric_threshold_config"`
	ThresholdConfig       []tfThresholdConfigModel       `tfsdk:"threshold_config"`
	Notification          []tfNotificationModel          `tfsdk:"notification"`
	Enabled               types.Bool                     `tfsdk:"enabled"`
}

type tfMatcherModel struct {
	FieldName types.String `tfsdk:"field_name"`
	Operator  types.String `tfsdk:"operator"`
	Value     types.String `tfsdk:"value"`
}

type tfMetricThresholdConfigModel struct {
	Threshold  types.Float64 `tfsdk:"threshold"`
	MetricName types.String  `tfsdk:"metric_name"`
	Operator   types.String  `tfsdk:"operator"`
	Units      types.String  `tfsdk:"units"`
	Mode       types.String  `tfsdk:"mode"`
}

type tfThresholdConfigModel struct {
	Threshold types.Float64 `tfsdk:"threshold"`
	Operator  types.String  `tfsdk:"operator"`
	Units     types.String  `tfsdk:"units"`
}

type tfNotificationModel struct {
	OpsGenieRegion           types.String `tfsdk:"ops_genie_region"`
	Username                 types.String `tfsdk:"username"`
	APIToken                 types.String `tfsdk:"api_token"`
	DatadogRegion            types.String `tfsdk:"datadog_region"`
	ServiceKey               types.String `tfsdk:"service_key"`
	EmailAddress             types.String `tfsdk:"email_address"`
	WebhookSecret            types.String `tfsdk:"webhook_secret"`
	MicrosoftTeamsWebhookURL types.String `tfsdk:"microsoft_teams_webhook_url"`
	MobileNumber             types.String `tfsdk:"mobile_number"`
	VictorOpsRoutingKey      types.String `tfsdk:"victor_ops_routing_key"`
	DatadogAPIKey            types.String `tfsdk:"datadog_api_key"`
	WebhookURL               types.String `tfsdk:"webhook_url"`
	OpsGenieAPIKey           types.String `tfsdk:"ops_genie_api_key"`
	TeamID                   types.String `tfsdk:"team_id"`
	TeamName                 types.String `tfsdk:"team_name"`
	TypeName                 types.String `tfsdk:"type_name"`
	ChannelName              types.String `tfsdk:"channel_name"`
	VictorOpsAPIKey          types.String `tfsdk:"victor_ops_api_key"`
	Roles                    []string     `tfsdk:"roles"`
	IntervalMin              types.Int64  `tfsdk:"interval_min"`
	DelayMin                 types.Int64  `tfsdk:"delay_min"`
	SMSEnabled               types.Bool   `tfsdk:"sms_enabled"`
	EmailEnabled             types.Bool   `tfsdk:"email_enabled"`
}

func (r *AlertConfigurationRS) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, alertConfigurationResourceName)
}

func (r *AlertConfigurationRS) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := ConfigureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.client = client
}

func (r *AlertConfigurationRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alert_configuration_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"event_type": schema.StringAttribute{
				Required: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"matcher": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"field_name": schema.StringAttribute{
							Optional: true,
						},
						"operator": schema.StringAttribute{
							Optional: true,
						},
						"value": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"metric_threshold_config": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"metric_name": schema.StringAttribute{
							Optional: true,
						},
						"operator": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("GREATER_THAN", "LESS_THAN"),
							},
						},
						"threshold": schema.Float64Attribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Float64{
								float64planmodifier.UseStateForUnknown(),
							},
						},
						"units": schema.StringAttribute{
							Optional: true,
						},
						"mode": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"threshold_config": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"operator": schema.StringAttribute{
							Optional: true,
						},
						"threshold": schema.Float64Attribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Float64{
								float64planmodifier.UseStateForUnknown(),
							},
						},
						"units": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
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
									"DAYS"),
							},
						},
					},
				},
			},
			"notification": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"api_token": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"channel_name": schema.StringAttribute{
							Optional: true,
						},
						"datadog_api_key": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"datadog_region": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("US", "EU"),
							},
						},
						"delay_min": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"email_address": schema.StringAttribute{
							Optional: true,
						},
						"email_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"interval_min": schema.Int64Attribute{
							Optional: true,
						},
						"mobile_number": schema.StringAttribute{
							Optional: true,
						},
						"ops_genie_api_key": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"ops_genie_region": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("US", "EU"),
							},
						},
						"service_key": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"sms_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"team_id": schema.StringAttribute{
							Optional: true,
						},
						"team_name": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type_name": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("EMAIL", "SMS", pagerDuty, "SLACK",
									"DATADOG", opsGenie, victorOps,
									"WEBHOOK", "USER", "TEAM", "GROUP", "ORG", "MICROSOFT_TEAMS"),
							},
						},
						"username": schema.StringAttribute{
							Optional: true,
						},
						"victor_ops_api_key": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"victor_ops_routing_key": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"roles": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
						"microsoft_teams_webhook_url": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"webhook_secret": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
						"webhook_url": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
						},
					},
				},
			},
		},
	}
}

func (r *AlertConfigurationRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.client.Atlas

	var alertConfigPlan tfAlertConfigurationRSModel

	diags := req.Plan.Get(ctx, &alertConfigPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := alertConfigPlan.ProjectID.ValueString()

	apiReq := &matlas.AlertConfiguration{
		EventTypeName:   alertConfigPlan.EventType.ValueString(),
		Enabled:         alertConfigPlan.Enabled.ValueBoolPointer(),
		Matchers:        newMatcherList(alertConfigPlan.Matcher),
		MetricThreshold: newMetricThreshold(alertConfigPlan.MetricThresholdConfig),
		Threshold:       newThreshold(alertConfigPlan.ThresholdConfig),
	}

	notifications, err := newNotificationList(alertConfigPlan.Notification)
	if err != nil {
		resp.Diagnostics.AddError(errorCreateAlertConf, err.Error())
		return
	}
	apiReq.Notifications = notifications

	apiResp, _, err := conn.AlertConfigurations.Create(ctx, projectID, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(errorCreateAlertConf, err.Error())
		return
	}

	encodedID := encodeStateID(map[string]string{
		encodedIDKeyAlertID:   apiResp.ID,
		encodedIDKeyProjectID: projectID,
	})
	alertConfigPlan.ID = types.StringValue(encodedID)

	newAlertConfigurationState := newTFAlertConfigurationModel(apiResp, &alertConfigPlan)

	// set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, newAlertConfigurationState)...)
}

func (r *AlertConfigurationRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.client.Atlas

	var alertConfigState tfAlertConfigurationRSModel

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(alertConfigState.ID.ValueString())

	alert, getResp, err := conn.AlertConfigurations.GetAnAlertConfig(context.Background(), ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID])
	if err != nil {
		// deleted in the backend case
		if getResp != nil && getResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}

	newAlertConfigurationState := newTFAlertConfigurationModel(alert, &alertConfigState)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newAlertConfigurationState)...)
}

func (r *AlertConfigurationRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.client.Atlas

	var alertConfigState, alertConfigPlan tfAlertConfigurationRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &alertConfigPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(alertConfigState.ID.ValueString())

	// In order to update an alert config it is necessary to send the original alert configuration request again, if not the
	// server returns an error 500
	apiReq, _, err := conn.AlertConfigurations.GetAnAlertConfig(ctx, ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID])
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}
	// Removing the computed attributes to recreate the original request
	apiReq.GroupID = ""
	apiReq.Created = ""
	apiReq.Updated = ""

	// Only changes the updated fields
	if !alertConfigPlan.Enabled.Equal(alertConfigState.Enabled) {
		apiReq.Enabled = alertConfigPlan.Enabled.ValueBoolPointer()
	}

	if !alertConfigPlan.EventType.Equal(alertConfigState.EventType) {
		apiReq.EventTypeName = alertConfigPlan.EventType.ValueString()
	}

	if !reflect.DeepEqual(alertConfigPlan.MetricThresholdConfig, alertConfigState.MetricThresholdConfig) {
		apiReq.MetricThreshold = newMetricThreshold(alertConfigPlan.MetricThresholdConfig)
	}

	if !reflect.DeepEqual(alertConfigPlan.ThresholdConfig, alertConfigState.ThresholdConfig) {
		apiReq.Threshold = newThreshold(alertConfigPlan.ThresholdConfig)
	}

	if !reflect.DeepEqual(alertConfigPlan.Matcher, alertConfigState.Matcher) {
		apiReq.Matchers = newMatcherList(alertConfigPlan.Matcher)
	}

	// Always refresh structure to handle service keys being obfuscated coming back from read API call
	notifications, err := newNotificationList(alertConfigPlan.Notification)
	if err != nil {
		resp.Diagnostics.AddError(errorUpdateAlertConf, err.Error())
		return
	}
	apiReq.Notifications = notifications

	var updatedAlertConfigResp *matlas.AlertConfiguration

	// Cannot enable/disable ONLY via update (if only send enable as changed field server returns a 500 error) so have to use different method to change enabled.
	if reflect.DeepEqual(apiReq, &matlas.AlertConfiguration{Enabled: pointy.Bool(true)}) ||
		reflect.DeepEqual(apiReq, &matlas.AlertConfiguration{Enabled: pointy.Bool(false)}) {
		// this code seems unreachable, as notifications are always being set
		updatedAlertConfigResp, _, err = conn.AlertConfigurations.EnableAnAlertConfig(ctx, ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID], apiReq.Enabled)
	} else {
		updatedAlertConfigResp, _, err = conn.AlertConfigurations.Update(ctx, ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID], apiReq)
	}

	if err != nil {
		resp.Diagnostics.AddError(errorUpdateAlertConf, err.Error())
		return
	}

	newAlertConfigurationState := newTFAlertConfigurationModel(updatedAlertConfigResp, &alertConfigPlan)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newAlertConfigurationState)...)
}

func (r *AlertConfigurationRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.client.Atlas

	var alertConfigState tfAlertConfigurationRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(alertConfigState.ID.ValueString())

	_, err := conn.AlertConfigurations.Delete(ctx, ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID])
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
	}
}

func (r *AlertConfigurationRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "-", 2)

	if len(parts) != 2 {
		resp.Diagnostics.AddError("import format error: to import a cluster, use the format {project_id}-{id}", "")
	}

	projectID := parts[0]
	alertConfigurationID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), encodeStateID(map[string]string{
		"id":         alertConfigurationID,
		"project_id": projectID,
	}))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}

func newNotificationList(tfNotificationSlice []tfNotificationModel) ([]matlas.Notification, error) {
	notifications := make([]matlas.Notification, len(tfNotificationSlice))
	if len(tfNotificationSlice) == 0 {
		return notifications, nil
	}

	for i := range tfNotificationSlice {
		value := tfNotificationSlice[i]

		if value.IntervalMin.ValueInt64() > 0 {
			typeName := value.TypeName.ValueString()
			if strings.EqualFold(typeName, pagerDuty) || strings.EqualFold(typeName, opsGenie) || strings.EqualFold(typeName, victorOps) {
				return nil, fmt.Errorf(`'interval_min' doesn't need to be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'`)
			}
		}

		notifications[i] = matlas.Notification{
			APIToken:                 value.APIToken.ValueString(),
			ChannelName:              value.ChannelName.ValueString(),
			DatadogAPIKey:            value.DatadogAPIKey.ValueString(),
			DatadogRegion:            value.DatadogRegion.ValueString(),
			DelayMin:                 pointy.Int(int(value.DelayMin.ValueInt64())),
			EmailAddress:             value.EmailAddress.ValueString(),
			EmailEnabled:             value.EmailEnabled.ValueBoolPointer(),
			IntervalMin:              int(value.IntervalMin.ValueInt64()),
			MobileNumber:             value.MobileNumber.ValueString(),
			OpsGenieAPIKey:           value.OpsGenieAPIKey.ValueString(),
			OpsGenieRegion:           value.OpsGenieRegion.ValueString(),
			ServiceKey:               value.ServiceKey.ValueString(),
			SMSEnabled:               value.SMSEnabled.ValueBoolPointer(),
			TeamID:                   value.TeamID.ValueString(),
			TypeName:                 value.TypeName.ValueString(),
			Username:                 value.Username.ValueString(),
			VictorOpsAPIKey:          value.VictorOpsAPIKey.ValueString(),
			VictorOpsRoutingKey:      value.VictorOpsRoutingKey.ValueString(),
			Roles:                    value.Roles,
			MicrosoftTeamsWebhookURL: value.MicrosoftTeamsWebhookURL.ValueString(),
			WebhookSecret:            value.WebhookSecret.ValueString(),
			WebhookURL:               value.WebhookURL.ValueString(),
		}
	}

	return notifications, nil
}

func newThreshold(tfThresholdConfigSlice []tfThresholdConfigModel) *matlas.Threshold {
	if len(tfThresholdConfigSlice) < 1 {
		return nil
	}

	v := tfThresholdConfigSlice[0]
	return &matlas.Threshold{
		Operator:  v.Operator.ValueString(),
		Units:     v.Units.ValueString(),
		Threshold: v.Threshold.ValueFloat64(),
	}
}

func newMetricThreshold(tfMetricThresholdConfigSlice []tfMetricThresholdConfigModel) *matlas.MetricThreshold {
	if len(tfMetricThresholdConfigSlice) < 1 {
		return nil
	}
	v := tfMetricThresholdConfigSlice[0]
	return &matlas.MetricThreshold{
		MetricName: v.MetricName.ValueString(),
		Operator:   v.Operator.ValueString(),
		Threshold:  v.Threshold.ValueFloat64(),
		Units:      v.Units.ValueString(),
		Mode:       v.Mode.ValueString(),
	}
}

func newMatcherList(tfMatcherSlice []tfMatcherModel) []matlas.Matcher {
	matchers := make([]matlas.Matcher, len(tfMatcherSlice))

	for i, m := range tfMatcherSlice {
		matchers[i] = matlas.Matcher{
			FieldName: m.FieldName.ValueString(),
			Operator:  m.Operator.ValueString(),
			Value:     m.Value.ValueString(),
		}
	}

	return matchers
}

func newTFAlertConfigurationModel(apiRespConfig *matlas.AlertConfiguration, currState *tfAlertConfigurationRSModel) tfAlertConfigurationRSModel {
	return tfAlertConfigurationRSModel{
		ID:                    currState.ID,
		ProjectID:             currState.ProjectID,
		AlertConfigurationID:  types.StringValue(apiRespConfig.ID),
		EventType:             types.StringValue(apiRespConfig.EventTypeName),
		Created:               types.StringValue(apiRespConfig.Created),
		Updated:               types.StringValue(apiRespConfig.Updated),
		Enabled:               types.BoolPointerValue(apiRespConfig.Enabled),
		MetricThresholdConfig: newTFMetricThresholdConfigModel(apiRespConfig.MetricThreshold, currState.MetricThresholdConfig),
		ThresholdConfig:       newTFThresholdConfigModel(apiRespConfig.Threshold, currState.ThresholdConfig),
		Notification:          newTFNotificationModelList(apiRespConfig.Notifications, currState.Notification),
		Matcher:               newTFMatcherModelList(apiRespConfig.Matchers, currState.Matcher),
	}
}

func newTFNotificationModelList(matlasSlice []matlas.Notification, currStateNotifications []tfNotificationModel) []tfNotificationModel {
	notifications := make([]tfNotificationModel, len(matlasSlice))

	if len(matlasSlice) != len(currStateNotifications) { // notifications were modified elsewhere from terraform, or import statement is being called
		for i := range matlasSlice {
			value := matlasSlice[i]
			notifications[i] = tfNotificationModel{
				TeamName:       types.StringValue(value.TeamName),
				Roles:          value.Roles,
				ChannelName:    conversion.StringNullIfEmpty(value.ChannelName),
				DatadogRegion:  conversion.StringNullIfEmpty(value.DatadogRegion),
				DelayMin:       types.Int64Value(int64(*value.DelayMin)),
				EmailAddress:   conversion.StringNullIfEmpty(value.EmailAddress),
				IntervalMin:    types.Int64Value(int64(value.IntervalMin)),
				MobileNumber:   conversion.StringNullIfEmpty(value.MobileNumber),
				OpsGenieRegion: conversion.StringNullIfEmpty(value.OpsGenieRegion),
				TeamID:         conversion.StringNullIfEmpty(value.TeamID),
				TypeName:       conversion.StringNullIfEmpty(value.TypeName),
				Username:       conversion.StringNullIfEmpty(value.Username),
				EmailEnabled:   types.BoolValue(value.EmailEnabled != nil && *value.EmailEnabled),
				SMSEnabled:     types.BoolValue(value.SMSEnabled != nil && *value.SMSEnabled),
			}
		}
		return notifications
	}

	for i := range matlasSlice {
		value := matlasSlice[i]
		currState := currStateNotifications[i]
		newState := tfNotificationModel{
			TeamName: types.StringValue(value.TeamName),
			Roles:    value.Roles,
		}

		// sentive attributes do not use value returned from API
		newState.APIToken = conversion.StringNullIfEmpty(currState.APIToken.ValueString())
		newState.DatadogAPIKey = conversion.StringNullIfEmpty(currState.DatadogAPIKey.ValueString())
		newState.OpsGenieAPIKey = conversion.StringNullIfEmpty(currState.OpsGenieAPIKey.ValueString())
		newState.ServiceKey = conversion.StringNullIfEmpty(currState.ServiceKey.ValueString())
		newState.VictorOpsAPIKey = conversion.StringNullIfEmpty(currState.VictorOpsAPIKey.ValueString())
		newState.VictorOpsRoutingKey = conversion.StringNullIfEmpty(currState.VictorOpsRoutingKey.ValueString())
		newState.WebhookURL = conversion.StringNullIfEmpty(currState.WebhookURL.ValueString())
		newState.WebhookSecret = conversion.StringNullIfEmpty(currState.WebhookSecret.ValueString())
		newState.MicrosoftTeamsWebhookURL = conversion.StringNullIfEmpty(currState.MicrosoftTeamsWebhookURL.ValueString())

		// for optional attributes that are not computed we must check if they were previously defined in state
		if !currState.ChannelName.IsNull() {
			newState.ChannelName = conversion.StringNullIfEmpty(value.ChannelName)
		}
		if !currState.DatadogRegion.IsNull() {
			newState.DatadogRegion = conversion.StringNullIfEmpty(value.DatadogRegion)
		}
		if !currState.EmailAddress.IsNull() {
			newState.EmailAddress = conversion.StringNullIfEmpty(value.EmailAddress)
		}
		if !currState.IntervalMin.IsNull() {
			newState.IntervalMin = types.Int64Value(int64(value.IntervalMin))
		}
		if !currState.MobileNumber.IsNull() {
			newState.MobileNumber = conversion.StringNullIfEmpty(value.MobileNumber)
		}
		if !currState.OpsGenieRegion.IsNull() {
			newState.OpsGenieRegion = conversion.StringNullIfEmpty(value.OpsGenieRegion)
		}
		if !currState.TeamID.IsNull() {
			newState.TeamID = conversion.StringNullIfEmpty(value.TeamID)
		}
		if !currState.TypeName.IsNull() {
			newState.TypeName = conversion.StringNullIfEmpty(value.TypeName)
		}
		if !currState.Username.IsNull() {
			newState.Username = conversion.StringNullIfEmpty(value.Username)
		}

		newState.DelayMin = types.Int64Value(int64(*value.DelayMin))
		newState.EmailEnabled = types.BoolValue(value.EmailEnabled != nil && *value.EmailEnabled)
		newState.SMSEnabled = types.BoolValue(value.SMSEnabled != nil && *value.SMSEnabled)

		notifications[i] = newState
	}

	return notifications
}

func newTFMetricThresholdConfigModel(matlasMetricThreshold *matlas.MetricThreshold, currStateSlice []tfMetricThresholdConfigModel) []tfMetricThresholdConfigModel {
	if matlasMetricThreshold == nil {
		return []tfMetricThresholdConfigModel{}
	}
	if len(currStateSlice) == 0 { // metric threshold was created elsewhere from terraform, or import statement is being called
		return []tfMetricThresholdConfigModel{
			{
				MetricName: conversion.StringNullIfEmpty(matlasMetricThreshold.MetricName),
				Operator:   conversion.StringNullIfEmpty(matlasMetricThreshold.Operator),
				Threshold:  types.Float64Value(matlasMetricThreshold.Threshold),
				Units:      conversion.StringNullIfEmpty(matlasMetricThreshold.Units),
				Mode:       conversion.StringNullIfEmpty(matlasMetricThreshold.Mode),
			},
		}
	}
	currState := currStateSlice[0]
	newState := tfMetricThresholdConfigModel{}
	if !currState.MetricName.IsNull() {
		newState.MetricName = conversion.StringNullIfEmpty(matlasMetricThreshold.MetricName)
	}
	if !currState.Operator.IsNull() {
		newState.Operator = conversion.StringNullIfEmpty(matlasMetricThreshold.Operator)
	}
	if !currState.Units.IsNull() {
		newState.Units = conversion.StringNullIfEmpty(matlasMetricThreshold.Units)
	}
	if !currState.Mode.IsNull() {
		newState.Mode = conversion.StringNullIfEmpty(matlasMetricThreshold.Mode)
	}
	newState.Threshold = types.Float64Value(matlasMetricThreshold.Threshold)
	return []tfMetricThresholdConfigModel{newState}
}

func newTFThresholdConfigModel(atlasThreshold *matlas.Threshold, currStateSlice []tfThresholdConfigModel) []tfThresholdConfigModel {
	if atlasThreshold == nil {
		return []tfThresholdConfigModel{}
	}

	if len(currStateSlice) == 0 { // threshold was created elsewhere from terraform, or import statement is being called
		return []tfThresholdConfigModel{
			{
				Operator:  conversion.StringNullIfEmpty(atlasThreshold.Operator),
				Threshold: types.Float64Value(atlasThreshold.Threshold),
				Units:     conversion.StringNullIfEmpty(atlasThreshold.Units),
			},
		}
	}
	currState := currStateSlice[0]
	newState := tfThresholdConfigModel{}
	if !currState.Operator.IsNull() {
		newState.Operator = conversion.StringNullIfEmpty(atlasThreshold.Operator)
	}
	if !currState.Units.IsNull() {
		newState.Units = conversion.StringNullIfEmpty(atlasThreshold.Units)
	}
	newState.Threshold = types.Float64Value(atlasThreshold.Threshold)

	return []tfThresholdConfigModel{newState}
}

func newTFMatcherModelList(matlasSlice []matlas.Matcher, currStateSlice []tfMatcherModel) []tfMatcherModel {
	matchers := make([]tfMatcherModel, len(matlasSlice))
	if len(matlasSlice) != len(currStateSlice) { // matchers were modified elsewhere from terraform, or import statement is being called
		for i, value := range matlasSlice {
			matchers[i] = tfMatcherModel{
				FieldName: conversion.StringNullIfEmpty(value.FieldName),
				Operator:  conversion.StringNullIfEmpty(value.Operator),
				Value:     conversion.StringNullIfEmpty(value.Value),
			}
		}
		return matchers
	}
	for i, value := range matlasSlice {
		currState := currStateSlice[i]
		newState := tfMatcherModel{}
		if !currState.FieldName.IsNull() {
			newState.FieldName = conversion.StringNullIfEmpty(value.FieldName)
		}
		if !currState.Operator.IsNull() {
			newState.Operator = conversion.StringNullIfEmpty(value.Operator)
		}
		if !currState.Value.IsNull() {
			newState.Value = conversion.StringNullIfEmpty(value.Value)
		}
		matchers[i] = newState
	}
	return matchers
}
