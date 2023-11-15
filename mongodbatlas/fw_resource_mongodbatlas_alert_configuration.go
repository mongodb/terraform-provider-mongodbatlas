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
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
)

const (
	alertConfigurationResourceName = "alert_configuration"
	errorCreateAlertConf           = "error creating Alert Configuration information: %s"
	errorReadAlertConf             = "error getting Alert Configuration information: %s"
	errorUpdateAlertConf           = "error updating Alert Configuration information: %s"
	pagerDuty                      = "PAGER_DUTY"
	opsGenie                       = "OPS_GENIE"
	victorOps                      = "VICTOR_OPS"
	encodedIDKeyAlertID            = "id"
	encodedIDKeyProjectID          = "project_id"
)

var _ resource.ResourceWithConfigure = &AlertConfigurationRS{}
var _ resource.ResourceWithImportState = &AlertConfigurationRS{}

func NewAlertConfigurationRS() resource.Resource {
	return &AlertConfigurationRS{
		RSCommon: RSCommon{
			resourceName: alertConfigurationResourceName,
		},
	}
}

type AlertConfigurationRS struct {
	RSCommon
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
	NotifierID               types.String `tfsdk:"notifier_id"`
	TypeName                 types.String `tfsdk:"type_name"`
	ChannelName              types.String `tfsdk:"channel_name"`
	VictorOpsAPIKey          types.String `tfsdk:"victor_ops_api_key"`
	Roles                    []string     `tfsdk:"roles"`
	IntervalMin              types.Int64  `tfsdk:"interval_min"`
	DelayMin                 types.Int64  `tfsdk:"delay_min"`
	SMSEnabled               types.Bool   `tfsdk:"sms_enabled"`
	EmailEnabled             types.Bool   `tfsdk:"email_enabled"`
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
							Required: true,
						},
						"operator": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
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
							Required: true,
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
							Computed: true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
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
						"notifier_id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"type_name": schema.StringAttribute{
							Required: true,
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
	connV2 := r.client.AtlasV2

	var alertConfigPlan tfAlertConfigurationRSModel

	diags := req.Plan.Get(ctx, &alertConfigPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := alertConfigPlan.ProjectID.ValueString()

	apiReq := &admin.GroupAlertsConfig{
		EventTypeName:   alertConfigPlan.EventType.ValueStringPointer(),
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

	apiResp, _, err := connV2.AlertConfigurationsApi.CreateAlertConfiguration(ctx, projectID, apiReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorCreateAlertConf, err.Error())
		return
	}

	encodedID := encodeStateID(map[string]string{
		encodedIDKeyAlertID:   *apiResp.Id,
		encodedIDKeyProjectID: projectID,
	})
	alertConfigPlan.ID = types.StringValue(encodedID)

	newAlertConfigurationState := newTFAlertConfigurationModel(apiResp, &alertConfigPlan)

	// set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, newAlertConfigurationState)...)
}

func (r *AlertConfigurationRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	connV2 := r.client.AtlasV2

	var alertConfigState tfAlertConfigurationRSModel

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(alertConfigState.ID.ValueString())

	alert, getResp, err := connV2.AlertConfigurationsApi.GetAlertConfiguration(context.Background(), ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID]).Execute()
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
	connV2 := r.client.AtlasV2

	var alertConfigState, alertConfigPlan tfAlertConfigurationRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &alertConfigPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(alertConfigState.ID.ValueString())

	// In order to update an alert config it is necessary to send the original alert configuration request again, if not the
	// server returns an error 500
	apiReq, _, err := connV2.AlertConfigurationsApi.GetAlertConfiguration(ctx, ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID]).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}
	// Removing the computed attributes to recreate the original request
	// //TODO: fix or remove comments
	// s := ""
	// apiReq.GroupId = &s
	apiReq.Created = nil
	apiReq.Updated = nil

	// Only changes the updated fields
	if !alertConfigPlan.Enabled.Equal(alertConfigState.Enabled) {
		apiReq.Enabled = alertConfigPlan.Enabled.ValueBoolPointer()
	}

	if !alertConfigPlan.EventType.Equal(alertConfigState.EventType) {
		apiReq.EventTypeName = alertConfigPlan.EventType.ValueStringPointer()
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

	var updatedAlertConfigResp *admin.GroupAlertsConfig

	// Cannot enable/disable ONLY via update (if only send enable as changed field server returns a 500 error) so have to use different method to change enabled.
	if reflect.DeepEqual(apiReq, &admin.GroupAlertsConfig{Enabled: pointy.Bool(true)}) ||
		reflect.DeepEqual(apiReq, &admin.GroupAlertsConfig{Enabled: pointy.Bool(false)}) {
		// this code seems unreachable, as notifications are always being set
		updatedAlertConfigResp, _, err = connV2.AlertConfigurationsApi.ToggleAlertConfiguration(context.Background(), ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID], &admin.AlertsToggle{Enabled: apiReq.Enabled}).Execute()
	} else {
		updatedAlertConfigResp, _, err = connV2.AlertConfigurationsApi.UpdateAlertConfiguration(context.Background(), ids[encodedIDKeyProjectID], ids[encodedIDKeyAlertID], apiReq).Execute()
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
		resp.Diagnostics.AddError("import format error", "to import an alert configuration, use the format {project_id}-{alert_configuration_id}")
		return
	}

	projectID := parts[0]
	alertConfigurationID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), encodeStateID(map[string]string{
		"id":         alertConfigurationID,
		"project_id": projectID,
	}))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}

func newNotificationList(tfNotificationSlice []tfNotificationModel) ([]admin.AlertsNotificationRootForGroup, error) {
	notifications := make([]admin.AlertsNotificationRootForGroup, 0)

	for i := range tfNotificationSlice {

		if !tfNotificationSlice[i].IntervalMin.IsNull() && tfNotificationSlice[i].IntervalMin.ValueInt64() > 0 {
			typeName := tfNotificationSlice[i].TypeName.ValueString()
			if strings.EqualFold(typeName, pagerDuty) || strings.EqualFold(typeName, opsGenie) || strings.EqualFold(typeName, victorOps) {
				return nil, fmt.Errorf(`'interval_min' doesn't need to be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'`)
			}
		}
	}

	for i := range tfNotificationSlice {
		notifierId := tfNotificationSlice[i].NotifierID.ValueStringPointer()
		//AlertConfigurationsApi expects nil or a valid notifier id
		if *notifierId == "" {
			notifierId = nil
		}
		notification := admin.AlertsNotificationRootForGroup{
			ApiToken:                 tfNotificationSlice[i].APIToken.ValueStringPointer(),
			ChannelName:              tfNotificationSlice[i].ChannelName.ValueStringPointer(),
			DatadogApiKey:            tfNotificationSlice[i].DatadogAPIKey.ValueStringPointer(),
			DatadogRegion:            tfNotificationSlice[i].DatadogRegion.ValueStringPointer(),
			DelayMin:                 pointy.Int(int(tfNotificationSlice[i].DelayMin.ValueInt64())),
			EmailAddress:             tfNotificationSlice[i].EmailAddress.ValueStringPointer(),
			EmailEnabled:             tfNotificationSlice[i].EmailEnabled.ValueBoolPointer(),
			IntervalMin:              pointy.Int(int(*tfNotificationSlice[i].IntervalMin.ValueInt64Pointer())),
			MobileNumber:             tfNotificationSlice[i].MobileNumber.ValueStringPointer(),
			OpsGenieApiKey:           tfNotificationSlice[i].OpsGenieAPIKey.ValueStringPointer(),
			OpsGenieRegion:           tfNotificationSlice[i].OpsGenieRegion.ValueStringPointer(),
			ServiceKey:               tfNotificationSlice[i].ServiceKey.ValueStringPointer(),
			SmsEnabled:               tfNotificationSlice[i].SMSEnabled.ValueBoolPointer(),
			TeamId:                   tfNotificationSlice[i].TeamID.ValueStringPointer(),
			TypeName:                 tfNotificationSlice[i].TypeName.ValueStringPointer(),
			Username:                 tfNotificationSlice[i].Username.ValueStringPointer(),
			VictorOpsApiKey:          tfNotificationSlice[i].VictorOpsAPIKey.ValueStringPointer(),
			VictorOpsRoutingKey:      tfNotificationSlice[i].VictorOpsRoutingKey.ValueStringPointer(),
			Roles:                    tfNotificationSlice[i].Roles,
			NotifierId:               notifierId,
			MicrosoftTeamsWebhookUrl: tfNotificationSlice[i].MicrosoftTeamsWebhookURL.ValueStringPointer(),
			WebhookSecret:            tfNotificationSlice[i].WebhookSecret.ValueStringPointer(),
			WebhookUrl:               tfNotificationSlice[i].WebhookURL.ValueStringPointer(),
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func newThreshold(tfThresholdConfigSlice []tfThresholdConfigModel) *admin.GreaterThanRawThreshold {
	if len(tfThresholdConfigSlice) < 1 {
		return nil
	}

	v := tfThresholdConfigSlice[0]
	return &admin.GreaterThanRawThreshold{
		Operator:  v.Operator.ValueStringPointer(),
		Units:     v.Units.ValueStringPointer(),
		Threshold: pointy.Int(int(v.Threshold.ValueFloat64())),
	}
}

func newMetricThreshold(tfMetricThresholdConfigSlice []tfMetricThresholdConfigModel) *admin.ServerlessMetricThreshold {
	if len(tfMetricThresholdConfigSlice) < 1 {
		return nil
	}
	v := tfMetricThresholdConfigSlice[0]
	return &admin.ServerlessMetricThreshold{
		MetricName: v.MetricName.ValueString(),
		Operator:   v.Operator.ValueStringPointer(),
		Threshold:  v.Threshold.ValueFloat64Pointer(),
		Units:      v.Units.ValueStringPointer(),
		Mode:       v.Mode.ValueStringPointer(),
	}
}

func newMatcherList(tfMatcherSlice []tfMatcherModel) []map[string]interface{} {
	matchers := make([]map[string]interface{}, 0)

	for i := range tfMatcherSlice {
		matcher := map[string]interface{}{
			"fieldName": tfMatcherSlice[i].FieldName.ValueString(),
			"operator":  tfMatcherSlice[i].Operator.ValueString(),
			"value":     tfMatcherSlice[i].Value.ValueString(),
		}
		matchers = append(matchers, matcher)
	}
	return matchers
}

func newTFAlertConfigurationModel(apiRespConfig *admin.GroupAlertsConfig, currState *tfAlertConfigurationRSModel) tfAlertConfigurationRSModel {
	return tfAlertConfigurationRSModel{
		ID:                    currState.ID,
		ProjectID:             currState.ProjectID,
		AlertConfigurationID:  types.StringValue(*apiRespConfig.Id),
		EventType:             types.StringValue(*apiRespConfig.EventTypeName),
		Created:               types.StringPointerValue(util.TimePtrToStringPtr(apiRespConfig.Created)),
		Updated:               types.StringPointerValue(util.TimePtrToStringPtr(apiRespConfig.Updated)),
		Enabled:               types.BoolPointerValue(apiRespConfig.Enabled),
		MetricThresholdConfig: newTFMetricThresholdConfigModel(apiRespConfig.MetricThreshold, currState.MetricThresholdConfig),
		ThresholdConfig:       newTFThresholdConfigModel(apiRespConfig.Threshold, currState.ThresholdConfig),
		Notification:          newTFNotificationModelList(apiRespConfig.Notifications, currState.Notification),
		Matcher:               newTFMatcherModelList(apiRespConfig.Matchers, currState.Matcher),
	}
}

func newTFNotificationModelList(n []admin.AlertsNotificationRootForGroup, currStateNotifications []tfNotificationModel) []tfNotificationModel {
	notifications := make([]tfNotificationModel, len(n))

	if len(n) != len(currStateNotifications) { // notifications were modified elsewhere from terraform, or import statement is being called
		for i := range n {
			value := n[i]
			notifications[i] = tfNotificationModel{
				TeamName:       conversion.StringPtrNullIfEmpty(value.TeamName),
				Roles:          value.Roles,
				ChannelName:    conversion.StringPtrNullIfEmpty(value.ChannelName),
				DatadogRegion:  conversion.StringPtrNullIfEmpty(value.DatadogRegion),
				DelayMin:       types.Int64PointerValue(util.IntPtrToInt64Ptr(value.DelayMin)),
				EmailAddress:   conversion.StringPtrNullIfEmpty(value.EmailAddress),
				IntervalMin:    types.Int64PointerValue(util.IntPtrToInt64Ptr(value.IntervalMin)),
				MobileNumber:   conversion.StringPtrNullIfEmpty(value.MobileNumber),
				OpsGenieRegion: conversion.StringPtrNullIfEmpty(value.OpsGenieRegion),
				TeamID:         conversion.StringPtrNullIfEmpty(value.TeamId),
				NotifierID:     types.StringPointerValue(value.NotifierId),
				TypeName:       conversion.StringPtrNullIfEmpty(value.TypeName),
				Username:       conversion.StringPtrNullIfEmpty(value.Username),
				EmailEnabled:   types.BoolValue(value.EmailEnabled != nil && *value.EmailEnabled),
				SMSEnabled:     types.BoolValue(value.SmsEnabled != nil && *value.SmsEnabled),
			}
		}
		return notifications
	}

	for i := range n {
		value := n[i]
		currState := currStateNotifications[i]
		newState := tfNotificationModel{
			TeamName: conversion.StringPtrNullIfEmpty(value.TeamName),
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
			newState.ChannelName = conversion.StringPtrNullIfEmpty(value.ChannelName)
		}
		if !currState.DatadogRegion.IsNull() {
			newState.DatadogRegion = conversion.StringPtrNullIfEmpty(value.DatadogRegion)
		}
		if !currState.EmailAddress.IsNull() {
			newState.EmailAddress = conversion.StringPtrNullIfEmpty(value.EmailAddress)
		}
		if !currState.MobileNumber.IsNull() {
			newState.MobileNumber = conversion.StringPtrNullIfEmpty(value.MobileNumber)
		}
		if !currState.OpsGenieRegion.IsNull() {
			newState.OpsGenieRegion = conversion.StringPtrNullIfEmpty(value.OpsGenieRegion)
		}
		if !currState.TeamID.IsNull() {
			newState.TeamID = conversion.StringPtrNullIfEmpty(value.TeamId)
		}
		if !currState.TypeName.IsNull() {
			newState.TypeName = conversion.StringPtrNullIfEmpty(value.TypeName)
		}
		if !currState.Username.IsNull() {
			newState.Username = conversion.StringPtrNullIfEmpty(value.Username)
		}

		newState.NotifierID = types.StringPointerValue(value.NotifierId)
		newState.IntervalMin = types.Int64PointerValue(util.IntPtrToInt64Ptr(value.IntervalMin))
		newState.DelayMin = types.Int64PointerValue(util.IntPtrToInt64Ptr(value.DelayMin))
		newState.EmailEnabled = types.BoolValue(value.EmailEnabled != nil && *value.EmailEnabled)
		newState.SMSEnabled = types.BoolValue(value.SmsEnabled != nil && *value.SmsEnabled)

		notifications[i] = newState
	}

	return notifications
}

func newTFMetricThresholdConfigModel(t *admin.ServerlessMetricThreshold, currStateSlice []tfMetricThresholdConfigModel) []tfMetricThresholdConfigModel {
	if t == nil {
		return []tfMetricThresholdConfigModel{}
	}
	if len(currStateSlice) == 0 { // metric threshold was created elsewhere from terraform, or import statement is being called
		return []tfMetricThresholdConfigModel{
			{
				MetricName: conversion.StringNullIfEmpty(t.MetricName),
				Operator:   conversion.StringNullIfEmpty(*t.Operator),
				Threshold:  types.Float64Value(*t.Threshold),
				Units:      conversion.StringNullIfEmpty(*t.Units),
				Mode:       conversion.StringNullIfEmpty(*t.Mode),
			},
		}
	}
	currState := currStateSlice[0]
	newState := tfMetricThresholdConfigModel{}
	if !currState.MetricName.IsNull() {
		newState.MetricName = conversion.StringNullIfEmpty(t.MetricName)
	}
	if !currState.Operator.IsNull() {
		newState.Operator = conversion.StringNullIfEmpty(*t.Operator)
	}
	if !currState.Units.IsNull() {
		newState.Units = conversion.StringNullIfEmpty(*t.Units)
	}
	if !currState.Mode.IsNull() {
		newState.Mode = conversion.StringNullIfEmpty(*t.Mode)
	}
	newState.Threshold = types.Float64Value(*t.Threshold)
	return []tfMetricThresholdConfigModel{newState}
}

func newTFThresholdConfigModel(t *admin.GreaterThanRawThreshold, currStateSlice []tfThresholdConfigModel) []tfThresholdConfigModel {
	if t == nil {
		return []tfThresholdConfigModel{}
	}

	if len(currStateSlice) == 0 { // threshold was created elsewhere from terraform, or import statement is being called
		return []tfThresholdConfigModel{
			{
				Operator:  conversion.StringNullIfEmpty(*t.Operator),
				Threshold: types.Float64Value(float64(*t.Threshold)), // int in new SDK but keeping float64 for backward compatibility
				Units:     conversion.StringNullIfEmpty(*t.Units),
			},
		}
	}
	currState := currStateSlice[0]
	newState := tfThresholdConfigModel{}
	if !currState.Operator.IsNull() {
		newState.Operator = conversion.StringNullIfEmpty(*t.Operator)
	}
	if !currState.Units.IsNull() {
		newState.Units = conversion.StringNullIfEmpty(*t.Units)
	}
	newState.Threshold = types.Float64Value(float64(*t.Threshold))

	return []tfThresholdConfigModel{newState}
}

func newTFMatcherModelList(m []map[string]any, currStateSlice []tfMatcherModel) []tfMatcherModel {
	matchers := make([]tfMatcherModel, len(m))
	if len(m) != len(currStateSlice) { // matchers were modified elsewhere from terraform, or import statement is being called
		for i, matcher := range m {
			fieldName, _ := matcher["fieldName"].(string)
			operator, _ := matcher["operator"].(string)
			value, _ := matcher["value"].(string)
			matchers[i] = tfMatcherModel{
				FieldName: conversion.StringNullIfEmpty(fieldName),
				Operator:  conversion.StringNullIfEmpty(operator),
				Value:     conversion.StringNullIfEmpty(value),
			}
		}
		return matchers
	}
	for i, matcher := range m {
		currState := currStateSlice[i]
		newState := tfMatcherModel{}
		if !currState.FieldName.IsNull() {
			fieldName, _ := matcher["fieldName"].(string)
			newState.FieldName = conversion.StringNullIfEmpty(fieldName)
		}
		if !currState.Operator.IsNull() {
			operator, _ := matcher["operator"].(string)
			newState.Operator = conversion.StringNullIfEmpty(operator)
		}
		if !currState.Value.IsNull() {
			value, _ := matcher["value"].(string)
			newState.Value = conversion.StringNullIfEmpty(value)
		}
		matchers[i] = newState
	}
	return matchers
}
