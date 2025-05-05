Resource name: mongodbatlas_alert_configuration

Resource schema:
```
package alertconfiguration

import (
	"context"
	"reflect"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

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

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	alertConfigurationResourceName = "alert_configuration"
	errorCreateAlertConf           = "error creating Alert Configuration information: %s"
	errorReadAlertConf             = "error getting Alert Configuration information: %s"
	errorUpdateAlertConf           = "error updating Alert Configuration information: %s"
	pagerDuty                      = "PAGER_DUTY"
	opsGenie                       = "OPS_GENIE"
	victorOps                      = "VICTOR_OPS"
	EncodedIDKeyAlertID            = "id"
	EncodedIDKeyProjectID          = "project_id"
)

var _ resource.ResourceWithConfigure = &alertConfigurationRS{}
var _ resource.ResourceWithImportState = &alertConfigurationRS{}

func Resource() resource.Resource {
	return &alertConfigurationRS{
		RSCommon: config.RSCommon{
			ResourceName: alertConfigurationResourceName,
		},
	}
}

type alertConfigurationRS struct {
	config.RSCommon
}

type TfAlertConfigurationRSModel struct {
	ID                    types.String                   `tfsdk:"id"`
	ProjectID             types.String                   `tfsdk:"project_id"`
	AlertConfigurationID  types.String                   `tfsdk:"alert_configuration_id"`
	EventType             types.String                   `tfsdk:"event_type"`
	Created               types.String                   `tfsdk:"created"`
	Updated               types.String                   `tfsdk:"updated"`
	Matcher               []TfMatcherModel               `tfsdk:"matcher"`
	MetricThresholdConfig []TfMetricThresholdConfigModel `tfsdk:"metric_threshold_config"`
	ThresholdConfig       []TfThresholdConfigModel       `tfsdk:"threshold_config"`
	Notification          []TfNotificationModel          `tfsdk:"notification"`
	Enabled               types.Bool                     `tfsdk:"enabled"`
}

type TfMatcherModel struct {
	FieldName types.String `tfsdk:"field_name"`
	Operator  types.String `tfsdk:"operator"`
	Value     types.String `tfsdk:"value"`
}

type TfMetricThresholdConfigModel struct {
	Threshold  types.Float64 `tfsdk:"threshold"`
	MetricName types.String  `tfsdk:"metric_name"`
	Operator   types.String  `tfsdk:"operator"`
	Units      types.String  `tfsdk:"units"`
	Mode       types.String  `tfsdk:"mode"`
}

type TfThresholdConfigModel struct {
	Threshold types.Float64 `tfsdk:"threshold"`
	Operator  types.String  `tfsdk:"operator"`
	Units     types.String  `tfsdk:"units"`
}

type TfNotificationModel struct {
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
	IntegrationID            types.String `tfsdk:"integration_id"`
	TypeName                 types.String `tfsdk:"type_name"`
	ChannelName              types.String `tfsdk:"channel_name"`
	VictorOpsAPIKey          types.String `tfsdk:"victor_ops_api_key"`
	Roles                    []string     `tfsdk:"roles"`
	IntervalMin              types.Int64  `tfsdk:"interval_min"`
	DelayMin                 types.Int64  `tfsdk:"delay_min"`
	SMSEnabled               types.Bool   `tfsdk:"sms_enabled"`
	EmailEnabled             types.Bool   `tfsdk:"email_enabled"`
}

func (r *alertConfigurationRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
						"integration_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
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

func (r *alertConfigurationRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	connV2 := r.Client.AtlasV2

	var alertConfigPlan TfAlertConfigurationRSModel

	diags := req.Plan.Get(ctx, &alertConfigPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := alertConfigPlan.ProjectID.ValueString()

	apiReq := &admin.GroupAlertsConfig{
		EventTypeName:   alertConfigPlan.EventType.ValueStringPointer(),
		Enabled:         alertConfigPlan.Enabled.ValueBoolPointer(),
		Matchers:        NewMatcherList(alertConfigPlan.Matcher),
		MetricThreshold: NewMetricThreshold(alertConfigPlan.MetricThresholdConfig),
		Threshold:       NewThreshold(alertConfigPlan.ThresholdConfig),
	}

	notifications, err := NewNotificationList(alertConfigPlan.Notification)
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

	encodedID := conversion.EncodeStateID(map[string]string{
		EncodedIDKeyAlertID:   conversion.SafeString(apiResp.Id),
		EncodedIDKeyProjectID: projectID,
	})
	alertConfigPlan.ID = types.StringValue(encodedID)

	newAlertConfigurationState := NewTFAlertConfigurationModel(apiResp, &alertConfigPlan)

	// set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, newAlertConfigurationState)...)
}

func (r *alertConfigurationRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	connV2 := r.Client.AtlasV2

	var alertConfigState TfAlertConfigurationRSModel

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := conversion.DecodeStateID(alertConfigState.ID.ValueString())

	alert, getResp, err := connV2.AlertConfigurationsApi.GetAlertConfiguration(context.Background(), ids[EncodedIDKeyProjectID], ids[EncodedIDKeyAlertID]).Execute()
	if err != nil {
		// deleted in the backend case
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}

	newAlertConfigurationState := NewTFAlertConfigurationModel(alert, &alertConfigState)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newAlertConfigurationState)...)
}

func (r *alertConfigurationRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	connV2 := r.Client.AtlasV2

	var alertConfigState, alertConfigPlan TfAlertConfigurationRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &alertConfigPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := conversion.DecodeStateID(alertConfigState.ID.ValueString())

	// In order to update an alert config it is necessary to send the original alert configuration request again, if not the
	// server returns an error 500
	apiReq, _, err := connV2.AlertConfigurationsApi.GetAlertConfiguration(ctx, ids[EncodedIDKeyProjectID], ids[EncodedIDKeyAlertID]).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
		return
	}
	// Removing the computed attributes to recreate the original request
	apiReq.GroupId = nil
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
		apiReq.MetricThreshold = NewMetricThreshold(alertConfigPlan.MetricThresholdConfig)
	}

	if !reflect.DeepEqual(alertConfigPlan.ThresholdConfig, alertConfigState.ThresholdConfig) {
		apiReq.Threshold = NewThreshold(alertConfigPlan.ThresholdConfig)
	}

	if !reflect.DeepEqual(alertConfigPlan.Matcher, alertConfigState.Matcher) {
		apiReq.Matchers = NewMatcherList(alertConfigPlan.Matcher)
	}

	// Always refresh structure to handle service keys being obfuscated coming back from read API call
	notifications, err := NewNotificationList(alertConfigPlan.Notification)
	if err != nil {
		resp.Diagnostics.AddError(errorUpdateAlertConf, err.Error())
		return
	}
	apiReq.Notifications = notifications

	var updatedAlertConfigResp *admin.GroupAlertsConfig

	// Cannot enable/disable ONLY via update (if only send enable as changed field server returns a 500 error) so have to use different method to change enabled.
	if reflect.DeepEqual(apiReq, &admin.GroupAlertsConfig{Enabled: conversion.Pointer(true)}) ||
		reflect.DeepEqual(apiReq, &admin.GroupAlertsConfig{Enabled: conversion.Pointer(false)}) {
		// this code seems unreachable, as notifications are always being set
		updatedAlertConfigResp, _, err = connV2.AlertConfigurationsApi.ToggleAlertConfiguration(
			context.Background(), ids[EncodedIDKeyProjectID], ids[EncodedIDKeyAlertID], &admin.AlertsToggle{Enabled: apiReq.Enabled}).Execute()
	} else {
		updatedAlertConfigResp, _, err = connV2.AlertConfigurationsApi.UpdateAlertConfiguration(context.Background(), ids[EncodedIDKeyProjectID], ids[EncodedIDKeyAlertID], apiReq).Execute()
	}

	if err != nil {
		resp.Diagnostics.AddError(errorUpdateAlertConf, err.Error())
		return
	}

	newAlertConfigurationState := NewTFAlertConfigurationModel(updatedAlertConfigResp, &alertConfigPlan)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newAlertConfigurationState)...)
}

func (r *alertConfigurationRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	connV2 := r.Client.AtlasV2

	var alertConfigState TfAlertConfigurationRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &alertConfigState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := conversion.DecodeStateID(alertConfigState.ID.ValueString())

	_, err := connV2.AlertConfigurationsApi.DeleteAlertConfiguration(ctx, ids[EncodedIDKeyProjectID], ids[EncodedIDKeyAlertID]).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorReadAlertConf, err.Error())
	}
}

func (r *alertConfigurationRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "-", 2)

	if len(parts) != 2 {
		resp.Diagnostics.AddError("import format error", "to import an alert configuration, use the format {project_id}-{alert_configuration_id}")
		return
	}

	projectID := parts[0]
	alertConfigurationID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), conversion.EncodeStateID(map[string]string{
		"id":         alertConfigurationID,
		"project_id": projectID,
	}))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
}
```
