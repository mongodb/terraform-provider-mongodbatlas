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

API Specification schema of GET response:
oneOf:
    - description: Other alerts which don't have extra details beside of basic one.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Incident that triggered this alert.
            oneOf:
                - enum:
                    - CREDIT_CARD_ABOUT_TO_EXPIRE
                  title: Billing Event Types
                  type: string
                - enum:
                    - CPS_SNAPSHOT_STARTED
                    - CPS_SNAPSHOT_SUCCESSFUL
                    - CPS_SNAPSHOT_FAILED
                    - CPS_CONCURRENT_SNAPSHOT_FAILED_WILL_RETRY
                    - CPS_SNAPSHOT_FALLBACK_SUCCESSFUL
                    - CPS_SNAPSHOT_FALLBACK_FAILED
                    - CPS_COPY_SNAPSHOT_STARTED
                    - CPS_COPY_SNAPSHOT_FAILED
                    - CPS_COPY_SNAPSHOT_FAILED_WILL_RETRY
                    - CPS_COPY_SNAPSHOT_SUCCESSFUL
                    - CPS_RESTORE_SUCCESSFUL
                    - CPS_EXPORT_SUCCESSFUL
                    - CPS_RESTORE_FAILED
                    - CPS_EXPORT_FAILED
                    - CPS_AUTO_EXPORT_FAILED
                    - CPS_SNAPSHOT_DOWNLOAD_REQUEST_FAILED
                    - CPS_OPLOG_CAUGHT_UP
                  title: Cps Backup Event Types
                  type: string
                - enum:
                    - CPS_DATA_PROTECTION_ENABLE_REQUESTED
                    - CPS_DATA_PROTECTION_ENABLED
                    - CPS_DATA_PROTECTION_UPDATE_REQUESTED
                    - CPS_DATA_PROTECTION_UPDATED
                    - CPS_DATA_PROTECTION_DISABLE_REQUESTED
                    - CPS_DATA_PROTECTION_DISABLED
                    - CPS_DATA_PROTECTION_APPROVED_FOR_DISABLEMENT
                  title: Data Protection Event Types
                  type: string
                - enum:
                    - FTS_INDEX_DELETION_FAILED
                    - FTS_INDEX_BUILD_COMPLETE
                    - FTS_INDEX_BUILD_FAILED
                    - FTS_INDEXES_RESTORE_FAILED
                    - FTS_INDEXES_SYNONYM_MAPPING_INVALID
                  title: FTS Index Audit Types
                  type: string
                - enum:
                    - USERS_WITHOUT_MULTI_FACTOR_AUTH
                    - ENCRYPTION_AT_REST_KMS_NETWORK_ACCESS_DENIED
                    - ENCRYPTION_AT_REST_CONFIG_NO_LONGER_VALID
                  title: Group Event Types
                  type: string
                - enum:
                    - CLUSTER_INSTANCE_STOP_START
                    - CLUSTER_INSTANCE_RESYNC_REQUESTED
                    - CLUSTER_INSTANCE_UPDATE_REQUESTED
                    - SAMPLE_DATASET_LOAD_REQUESTED
                    - TENANT_UPGRADE_TO_SERVERLESS_SUCCESSFUL
                    - TENANT_UPGRADE_TO_SERVERLESS_FAILED
                    - NETWORK_PERMISSION_ENTRY_ADDED
                    - NETWORK_PERMISSION_ENTRY_REMOVED
                    - NETWORK_PERMISSION_ENTRY_UPDATED
                  title: NDS Audit Types
                  type: string
                - enum:
                    - MAINTENANCE_IN_ADVANCED
                    - MAINTENANCE_AUTO_DEFERRED
                    - MAINTENANCE_STARTED
                    - MAINTENANCE_NO_LONGER_NEEDED
                  title: NDS Maintenance Window Audit Types
                  type: string
                - enum:
                    - ONLINE_ARCHIVE_INSUFFICIENT_INDEXES_CHECK
                    - ONLINE_ARCHIVE_MAX_CONSECUTIVE_OFFLOAD_WINDOWS_CHECK
                  title: Online Archive Event Types
                  type: string
                - enum:
                    - JOINED_GROUP
                    - REMOVED_FROM_GROUP
                    - USER_ROLES_CHANGED_AUDIT
                  title: User Event Types
                  type: string
                - enum:
                    - TAGS_MODIFIED
                    - CLUSTER_TAGS_MODIFIED
                    - GROUP_TAGS_MODIFIED
                  title: Resource Event Types
                  type: string
                - enum:
                    - STREAM_PROCESSOR_STATE_IS_FAILED
                    - OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD
                  title: Stream Processor Event Types
                  type: string
                - enum:
                    - COMPUTE_AUTO_SCALE_INITIATED_BASE
                    - COMPUTE_AUTO_SCALE_INITIATED_ANALYTICS
                    - COMPUTE_AUTO_SCALE_SCALE_DOWN_FAIL_BASE
                    - COMPUTE_AUTO_SCALE_SCALE_DOWN_FAIL_ANALYTICS
                    - COMPUTE_AUTO_SCALE_MAX_INSTANCE_SIZE_FAIL_BASE
                    - COMPUTE_AUTO_SCALE_MAX_INSTANCE_SIZE_FAIL_ANALYTICS
                    - COMPUTE_AUTO_SCALE_OPLOG_FAIL_BASE
                    - COMPUTE_AUTO_SCALE_OPLOG_FAIL_ANALYTICS
                    - DISK_AUTO_SCALE_INITIATED
                    - DISK_AUTO_SCALE_MAX_DISK_SIZE_FAIL
                    - DISK_AUTO_SCALE_OPLOG_FAIL
                    - PREDICTIVE_COMPUTE_AUTO_SCALE_INITIATED_BASE
                    - PREDICTIVE_COMPUTE_AUTO_SCALE_MAX_INSTANCE_SIZE_FAIL_BASE
                    - PREDICTIVE_COMPUTE_AUTO_SCALE_OPLOG_FAIL_BASE
                  title: NDS Auto Scaling Audit Types
                  type: string
                - enum:
                    - RESOURCE_POLICY_VIOLATED
                  title: Atlas Resource Policy Audit Types
                  type: string
            type: object
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Any Other Alert Configurations
      type: object
    - description: App Services metric alert configuration allows to select which app service conditions and events trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - URL_CONFIRMATION
                - SUCCESSFUL_DEPLOY
                - DEPLOYMENT_FAILURE
                - REQUEST_RATE_LIMIT
                - LOG_FORWARDER_FAILURE
                - SYNC_FAILURE
                - TRIGGER_FAILURE
                - TRIGGER_AUTO_RESUMED
                - DEPLOYMENT_MODEL_CHANGE_SUCCESS
                - DEPLOYMENT_MODEL_CHANGE_FAILURE
            example: DEPLOYMENT_FAILURE
            title: App Services Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an app service metric against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - APPLICATION_ID
                        example: APPLICATION_ID
                        title: App Services Metric Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: App Services Alert Configuration
      type: object
    - description: App Services metric alert configuration allows to select which app service metrics trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - OUTSIDE_REALM_METRIC_THRESHOLD
            example: OUTSIDE_REALM_METRIC_THRESHOLD
            title: App Services Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an app service metric against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - APPLICATION_ID
                        example: APPLICATION_ID
                        title: App Services Metric Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        metricThreshold:
            description: Threshold for the metric that, when exceeded, triggers an alert. The metric threshold pertains to event types which reflects changes of measurements and metrics in the app services.
            discriminator:
                mapping:
                    REALM_AUTH_LOGIN_FAIL: '#/components/schemas/RawMetricThresholdView'
                    REALM_ENDPOINTS_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_ENDPOINTS_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_ENDPOINTS_FAILED_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    REALM_ENDPOINTS_RESPONSE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_GQL_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_GQL_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_GQL_FAILED_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    REALM_GQL_RESPONSE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_OVERALL_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_OVERALL_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_OVERALL_FAILED_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    REALM_SDK_FNS_RESPONSE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SDK_FUNCTIONS_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SDK_FUNCTIONS_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_SDK_MQL_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SDK_MQL_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_SDK_MQL_RESPONSE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SDKFNS_FAILED_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    REALM_SYNC_CLIENT_BOOTSTRAP_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SYNC_CLIENT_CHANGESETS_INVALID: '#/components/schemas/DataMetricThresholdView'
                    REALM_SYNC_CLIENT_READS_FAILED: '#/components/schemas/DataMetricThresholdView'
                    REALM_SYNC_CLIENT_UPLOADS_FAILED: '#/components/schemas/DataMetricThresholdView'
                    REALM_SYNC_CURRENT_OPLOG_LAG_MS_SUM: '#/components/schemas/TimeMetricThresholdView'
                    REALM_SYNC_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_SYNC_FAILED_REQUESTS: '#/components/schemas/DataMetricThresholdView'
                    REALM_SYNC_NUM_UNSYNCABLE_DOCS_PERCENT: '#/components/schemas/RawMetricThresholdView'
                    REALM_SYNC_SESSIONS_ENDED: '#/components/schemas/DataMetricThresholdView'
                    REALM_TRIGGERS_COMPUTE_MS: '#/components/schemas/TimeMetricThresholdView'
                    REALM_TRIGGERS_CURRENT_OPLOG_LAG_MS_SUM: '#/components/schemas/TimeMetricThresholdView'
                    REALM_TRIGGERS_EGRESS_BYTES: '#/components/schemas/DataMetricThresholdView'
                    REALM_TRIGGERS_FAILED_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    REALM_TRIGGERS_RESPONSE_MS: '#/components/schemas/TimeMetricThresholdView'
                propertyName: metricName
            properties:
                metricName:
                    description: Human-readable label that identifies the metric against which MongoDB Cloud checks the configured **metricThreshold.threshold**.
                    type: string
                mode:
                    description: MongoDB Cloud computes the current metric value as an average.
                    enum:
                        - AVERAGE
                    type: string
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: double
                    type: number
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            required:
                - metricName
            title: App Services Metric Threshold
            type: object
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: App Services Metric Alert Configuration
      type: object
    - description: Billing threshold alert configuration allows to select thresholds for bills and invoices which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - PENDING_INVOICE_OVER_THRESHOLD
                - DAILY_BILL_OVER_THRESHOLD
            example: PENDING_INVOICE_OVER_THRESHOLD
            title: Billing Event Type
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: A Limit that triggers an alert when greater than a number.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    default: RAW
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - RAW
                    title: Raw Metric Units
                    type: string
            title: Greater Than Raw Threshold
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Billing Threshold Alert Configuration
      type: object
    - description: Cluster alert configuration allows to select which conditions of mongod cluster which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - CLUSTER_MONGOS_IS_MISSING
                - CLUSTER_AGENT_IN_CRASH_LOOP
            example: CLUSTER_MONGOS_IS_MISSING
            title: Cluster Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an cluster against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - CLUSTER_NAME
                        example: CLUSTER_NAME
                        title: Cluster Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Cluster Alert Configuration
      type: object
    - description: Cps Backup threshold alert configuration allows to select thresholds for conditions of CPS backup or oplogs anomalies which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - CPS_SNAPSHOT_BEHIND
                - CPS_PREV_SNAPSHOT_OLD
                - CPS_OPLOG_BEHIND
            example: CPS_SNAPSHOT_BEHIND
            title: Cps Backup Event Type
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: A Limit that triggers an alert when greater than a time period.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    default: HOURS
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - NANOSECONDS
                        - MILLISECONDS
                        - MILLION_MINUTES
                        - SECONDS
                        - MINUTES
                        - HOURS
                        - DAYS
                    title: Time Metric Units
                    type: string
            title: Greater Than Time Threshold
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Cps Backup Threshold Alert Configuration
      type: object
    - description: Encryption key alert configuration allows to select thresholds  which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - AWS_ENCRYPTION_KEY_NEEDS_ROTATION
                - AZURE_ENCRYPTION_KEY_NEEDS_ROTATION
                - GCP_ENCRYPTION_KEY_NEEDS_ROTATION
                - AWS_ENCRYPTION_KEY_INVALID
                - AZURE_ENCRYPTION_KEY_INVALID
                - GCP_ENCRYPTION_KEY_INVALID
            example: AWS_ENCRYPTION_KEY_NEEDS_ROTATION
            title: Encryption Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: Threshold value that triggers an alert.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - DAYS
                    type: string
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Encryption Key Alert Configuration
      type: object
    - description: Host alert configuration allows to select which mongod host events trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - HOST_DOWN
                - HOST_HAS_INDEX_SUGGESTIONS
                - HOST_MONGOT_CRASHING_OOM
                - HOST_MONGOT_STOP_REPLICATION
                - HOST_NOT_ENOUGH_DISK_SPACE
                - SSH_KEY_NDS_HOST_ACCESS_REQUESTED
                - SSH_KEY_NDS_HOST_ACCESS_REFRESHED
                - PUSH_BASED_LOG_EXPORT_STOPPED
                - PUSH_BASED_LOG_EXPORT_DROPPED_LOG
                - HOST_VERSION_BEHIND
                - VERSION_BEHIND
                - HOST_EXPOSED
                - HOST_SSL_CERTIFICATE_STALE
                - HOST_SECURITY_CHECKUP_NOT_MET
            example: HOST_DOWN
            title: Host Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an host against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - TYPE_NAME
                            - HOSTNAME
                            - PORT
                            - HOSTNAME_AND_PORT
                            - REPLICA_SET_NAME
                        example: HOSTNAME
                        title: Host Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        enum:
                            - STANDALONE
                            - PRIMARY
                            - SECONDARY
                            - ARBITER
                            - MONGOS
                            - CONFIG
                        example: STANDALONE
                        title: Matcher Host Types
                        type: string
                required:
                    - fieldName
                    - operator
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Host Alert Configuration
      type: object
    - description: Host metric alert configuration allows to select which mongod host metrics trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - OUTSIDE_METRIC_THRESHOLD
            example: OUTSIDE_METRIC_THRESHOLD
            title: Host Metric Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an host against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - TYPE_NAME
                            - HOSTNAME
                            - PORT
                            - HOSTNAME_AND_PORT
                            - REPLICA_SET_NAME
                        example: HOSTNAME
                        title: Host Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        enum:
                            - STANDALONE
                            - PRIMARY
                            - SECONDARY
                            - ARBITER
                            - MONGOS
                            - CONFIG
                        example: STANDALONE
                        title: Matcher Host Types
                        type: string
                required:
                    - fieldName
                    - operator
                title: Matchers
                type: object
            type: array
        metricThreshold:
            description: Threshold for the metric that, when exceeded, triggers an alert. The metric threshold pertains to event types which reflects changes of measurements and metrics about mongod host.
            discriminator:
                mapping:
                    ASSERT_MSG: '#/components/schemas/RawMetricThresholdView'
                    ASSERT_REGULAR: '#/components/schemas/RawMetricThresholdView'
                    ASSERT_USER: '#/components/schemas/RawMetricThresholdView'
                    ASSERT_WARNING: '#/components/schemas/RawMetricThresholdView'
                    AVG_COMMAND_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    AVG_READ_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    AVG_WRITE_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    BACKGROUND_FLUSH_AVG: '#/components/schemas/TimeMetricThresholdView'
                    CACHE_BYTES_READ_INTO: '#/components/schemas/DataMetricThresholdView'
                    CACHE_BYTES_WRITTEN_FROM: '#/components/schemas/DataMetricThresholdView'
                    CACHE_USAGE_DIRTY: '#/components/schemas/DataMetricThresholdView'
                    CACHE_USAGE_USED: '#/components/schemas/DataMetricThresholdView'
                    COMPUTED_MEMORY: '#/components/schemas/DataMetricThresholdView'
                    CONNECTIONS: '#/components/schemas/RawMetricThresholdView'
                    CONNECTIONS_MAX: '#/components/schemas/RawMetricThresholdView'
                    CONNECTIONS_PERCENT: '#/components/schemas/RawMetricThresholdView'
                    CURSORS_TOTAL_CLIENT_CURSORS_SIZE: '#/components/schemas/RawMetricThresholdView'
                    CURSORS_TOTAL_OPEN: '#/components/schemas/RawMetricThresholdView'
                    CURSORS_TOTAL_TIMED_OUT: '#/components/schemas/RawMetricThresholdView'
                    DB_DATA_SIZE_TOTAL: '#/components/schemas/DataMetricThresholdView'
                    DB_DATA_SIZE_TOTAL_WO_SYSTEM: '#/components/schemas/DataMetricThresholdView'
                    DB_INDEX_SIZE_TOTAL: '#/components/schemas/DataMetricThresholdView'
                    DB_STORAGE_TOTAL: '#/components/schemas/DataMetricThresholdView'
                    DISK_PARTITION_QUEUE_DEPTH_DATA: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_QUEUE_DEPTH_INDEX: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_QUEUE_DEPTH_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_READ_IOPS_DATA: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_READ_IOPS_INDEX: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_READ_IOPS_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_READ_LATENCY_DATA: '#/components/schemas/TimeMetricThresholdView'
                    DISK_PARTITION_READ_LATENCY_INDEX: '#/components/schemas/TimeMetricThresholdView'
                    DISK_PARTITION_READ_LATENCY_JOURNAL: '#/components/schemas/TimeMetricThresholdView'
                    DISK_PARTITION_SPACE_USED_DATA: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_SPACE_USED_INDEX: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_SPACE_USED_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_WRITE_IOPS_DATA: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_WRITE_IOPS_INDEX: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_WRITE_IOPS_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    DISK_PARTITION_WRITE_LATENCY_DATA: '#/components/schemas/TimeMetricThresholdView'
                    DISK_PARTITION_WRITE_LATENCY_INDEX: '#/components/schemas/TimeMetricThresholdView'
                    DISK_PARTITION_WRITE_LATENCY_JOURNAL: '#/components/schemas/TimeMetricThresholdView'
                    DOCUMENT_DELETED: '#/components/schemas/RawMetricThresholdView'
                    DOCUMENT_INSERTED: '#/components/schemas/RawMetricThresholdView'
                    DOCUMENT_RETURNED: '#/components/schemas/RawMetricThresholdView'
                    DOCUMENT_UPDATED: '#/components/schemas/RawMetricThresholdView'
                    EXTRA_INFO_PAGE_FAULTS: '#/components/schemas/RawMetricThresholdView'
                    FTS_DISK_UTILIZATION: '#/components/schemas/DataMetricThresholdView'
                    FTS_JVM_CURRENT_MEMORY: '#/components/schemas/DataMetricThresholdView'
                    FTS_JVM_MAX_MEMORY: '#/components/schemas/DataMetricThresholdView'
                    FTS_MEMORY_MAPPED: '#/components/schemas/DataMetricThresholdView'
                    FTS_MEMORY_RESIDENT: '#/components/schemas/DataMetricThresholdView'
                    FTS_MEMORY_VIRTUAL: '#/components/schemas/DataMetricThresholdView'
                    FTS_PROCESS_CPU_KERNEL: '#/components/schemas/RawMetricThresholdView'
                    FTS_PROCESS_CPU_USER: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_ACCESSES_NOT_IN_MEMORY: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_LOCK_CURRENT_QUEUE_READERS: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_LOCK_CURRENT_QUEUE_TOTAL: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_LOCK_CURRENT_QUEUE_WRITERS: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_LOCK_PERCENTAGE: '#/components/schemas/RawMetricThresholdView'
                    GLOBAL_PAGE_FAULT_EXCEPTIONS_THROWN: '#/components/schemas/RawMetricThresholdView'
                    INDEX_COUNTERS_BTREE_ACCESSES: '#/components/schemas/RawMetricThresholdView'
                    INDEX_COUNTERS_BTREE_HITS: '#/components/schemas/RawMetricThresholdView'
                    INDEX_COUNTERS_BTREE_MISS_RATIO: '#/components/schemas/RawMetricThresholdView'
                    INDEX_COUNTERS_BTREE_MISSES: '#/components/schemas/RawMetricThresholdView'
                    JOURNALING_COMMITS_IN_WRITE_LOCK: '#/components/schemas/RawMetricThresholdView'
                    JOURNALING_MB: '#/components/schemas/DataMetricThresholdView'
                    JOURNALING_WRITE_DATA_FILES_MB: '#/components/schemas/DataMetricThresholdView'
                    LOGICAL_SIZE: '#/components/schemas/DataMetricThresholdView'
                    MAX_DISK_PARTITION_QUEUE_DEPTH_DATA: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_QUEUE_DEPTH_INDEX: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_QUEUE_DEPTH_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_READ_IOPS_DATA: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_READ_IOPS_INDEX: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_READ_IOPS_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_READ_LATENCY_DATA: '#/components/schemas/TimeMetricThresholdView'
                    MAX_DISK_PARTITION_READ_LATENCY_INDEX: '#/components/schemas/TimeMetricThresholdView'
                    MAX_DISK_PARTITION_READ_LATENCY_JOURNAL: '#/components/schemas/TimeMetricThresholdView'
                    MAX_DISK_PARTITION_SPACE_USED_DATA: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_SPACE_USED_INDEX: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_SPACE_USED_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_IOPS_DATA: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_IOPS_INDEX: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_IOPS_JOURNAL: '#/components/schemas/RawMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_LATENCY_DATA: '#/components/schemas/TimeMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_LATENCY_INDEX: '#/components/schemas/TimeMetricThresholdView'
                    MAX_DISK_PARTITION_WRITE_LATENCY_JOURNAL: '#/components/schemas/TimeMetricThresholdView'
                    MAX_NORMALIZED_SYSTEM_CPU_STEAL: '#/components/schemas/RawMetricThresholdView'
                    MAX_NORMALIZED_SYSTEM_CPU_USER: '#/components/schemas/RawMetricThresholdView'
                    MAX_SWAP_USAGE_FREE: '#/components/schemas/DataMetricThresholdView'
                    MAX_SWAP_USAGE_USED: '#/components/schemas/DataMetricThresholdView'
                    MAX_SYSTEM_MEMORY_AVAILABLE: '#/components/schemas/DataMetricThresholdView'
                    MAX_SYSTEM_MEMORY_PERCENT_USED: '#/components/schemas/RawMetricThresholdView'
                    MAX_SYSTEM_MEMORY_USED: '#/components/schemas/DataMetricThresholdView'
                    MAX_SYSTEM_NETWORK_IN: '#/components/schemas/DataMetricThresholdView'
                    MAX_SYSTEM_NETWORK_OUT: '#/components/schemas/DataMetricThresholdView'
                    MEMORY_MAPPED: '#/components/schemas/DataMetricThresholdView'
                    MEMORY_RESIDENT: '#/components/schemas/DataMetricThresholdView'
                    MEMORY_VIRTUAL: '#/components/schemas/DataMetricThresholdView'
                    MUNIN_CPU_IOWAIT: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_IRQ: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_NICE: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_SOFTIRQ: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_STEAL: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_SYSTEM: '#/components/schemas/RawMetricThresholdView'
                    MUNIN_CPU_USER: '#/components/schemas/RawMetricThresholdView'
                    NETWORK_BYTES_IN: '#/components/schemas/DataMetricThresholdView'
                    NETWORK_BYTES_OUT: '#/components/schemas/DataMetricThresholdView'
                    NETWORK_NUM_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    NORMALIZED_FTS_PROCESS_CPU_KERNEL: '#/components/schemas/RawMetricThresholdView'
                    NORMALIZED_FTS_PROCESS_CPU_USER: '#/components/schemas/RawMetricThresholdView'
                    NORMALIZED_SYSTEM_CPU_STEAL: '#/components/schemas/RawMetricThresholdView'
                    NORMALIZED_SYSTEM_CPU_USER: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_CMD: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_DELETE: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_GETMORE: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_INSERT: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_QUERY: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_REPL_CMD: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_REPL_DELETE: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_REPL_INSERT: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_REPL_UPDATE: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_TTL_DELETED: '#/components/schemas/RawMetricThresholdView'
                    OPCOUNTER_UPDATE: '#/components/schemas/RawMetricThresholdView'
                    OPERATION_THROTTLING_REJECTED_OPERATIONS: '#/components/schemas/RawMetricThresholdView'
                    OPERATIONS_QUERIES_KILLED: '#/components/schemas/RawMetricThresholdView'
                    OPERATIONS_SCAN_AND_ORDER: '#/components/schemas/RawMetricThresholdView'
                    OPLOG_MASTER_LAG_TIME_DIFF: '#/components/schemas/TimeMetricThresholdView'
                    OPLOG_MASTER_TIME: '#/components/schemas/TimeMetricThresholdView'
                    OPLOG_MASTER_TIME_ESTIMATED_TTL: '#/components/schemas/TimeMetricThresholdView'
                    OPLOG_RATE_GB_PER_HOUR: '#/components/schemas/DataMetricThresholdView'
                    OPLOG_SLAVE_LAG_MASTER_TIME: '#/components/schemas/TimeMetricThresholdView'
                    QUERY_EXECUTOR_SCANNED: '#/components/schemas/RawMetricThresholdView'
                    QUERY_EXECUTOR_SCANNED_OBJECTS: '#/components/schemas/RawMetricThresholdView'
                    QUERY_SPILL_TO_DISK_DURING_SORT: '#/components/schemas/RawMetricThresholdView'
                    QUERY_TARGETING_SCANNED_OBJECTS_PER_RETURNED: '#/components/schemas/RawMetricThresholdView'
                    QUERY_TARGETING_SCANNED_PER_RETURNED: '#/components/schemas/RawMetricThresholdView'
                    RESTARTS_IN_LAST_HOUR: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_INDEX_SIZE: '#/components/schemas/DataMetricThresholdView'
                    SEARCH_MAX_NUMBER_OF_LUCENE_DOCS: '#/components/schemas/NumberMetricThresholdView'
                    SEARCH_NUMBER_OF_FIELDS_IN_INDEX: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_NUMBER_OF_QUERIES_ERROR: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_NUMBER_OF_QUERIES_SUCCESS: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_NUMBER_OF_QUERIES_TOTAL: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_OPCOUNTER_DELETE: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_OPCOUNTER_GETMORE: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_OPCOUNTER_INSERT: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_OPCOUNTER_UPDATE: '#/components/schemas/RawMetricThresholdView'
                    SEARCH_REPLICATION_LAG: '#/components/schemas/TimeMetricThresholdView'
                    SWAP_USAGE_FREE: '#/components/schemas/DataMetricThresholdView'
                    SWAP_USAGE_USED: '#/components/schemas/DataMetricThresholdView'
                    SYSTEM_MEMORY_AVAILABLE: '#/components/schemas/DataMetricThresholdView'
                    SYSTEM_MEMORY_PERCENT_USED: '#/components/schemas/RawMetricThresholdView'
                    SYSTEM_MEMORY_USED: '#/components/schemas/DataMetricThresholdView'
                    SYSTEM_NETWORK_IN: '#/components/schemas/DataMetricThresholdView'
                    SYSTEM_NETWORK_OUT: '#/components/schemas/DataMetricThresholdView'
                    TICKETS_AVAILABLE_READS: '#/components/schemas/RawMetricThresholdView'
                    TICKETS_AVAILABLE_WRITES: '#/components/schemas/RawMetricThresholdView'
                propertyName: metricName
            properties:
                metricName:
                    description: Human-readable label that identifies the metric against which MongoDB Cloud checks the configured **metricThreshold.threshold**.
                    type: string
                mode:
                    description: MongoDB Cloud computes the current metric value as an average.
                    enum:
                        - AVERAGE
                    type: string
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: double
                    type: number
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            required:
                - metricName
            title: Host Metric Threshold
            type: object
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Host Metric Alert Configuration
      type: object
    - description: NDS X509 User Authentication alert configuration allows to select thresholds for expiration of client, CA certificates and CRL which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - NDS_X509_USER_AUTHENTICATION_CUSTOMER_CA_EXPIRATION_CHECK
                - NDS_X509_USER_AUTHENTICATION_CUSTOMER_CRL_EXPIRATION_CHECK
                - NDS_X509_USER_AUTHENTICATION_MANAGED_USER_CERTS_EXPIRATION_CHECK
            example: NDS_X509_USER_AUTHENTICATION_CUSTOMER_CA_EXPIRATION_CHECK
            title: NDS x509 User Auth Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: Threshold value that triggers an alert.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - DAYS
                    type: string
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: NDS X509 User Authentication Alert Configuration
      type: object
    - description: Replica Set alert configuration allows to select which conditions of mongod replica set trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - NO_PRIMARY
                - PRIMARY_ELECTED
            example: NO_PRIMARY
            title: ReplicaSet Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an replica set against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - REPLICA_SET_NAME
                            - SHARD_NAME
                            - CLUSTER_NAME
                        example: REPLICA_SET_NAME
                        title: Replica Set Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: A Limit that triggers an alert when  exceeded. The resource returns this parameter when **eventTypeName** has not been set to `OUTSIDE_METRIC_THRESHOLD`.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - <
                        - '>'
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            title: Threshold
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Replica Set Alert Configuration
      type: object
    - description: Replica Set threshold alert configuration allows to select thresholds for conditions of mongod replica set which trigger alerts and how users are notified.
      discriminator:
        mapping:
            REPLICATION_OPLOG_WINDOW_RUNNING_OUT: '#/components/schemas/LessThanTimeThresholdAlertConfigViewForNdsGroup'
            TOO_MANY_ELECTIONS: '#/components/schemas/GreaterThanRawThresholdAlertConfigViewForNdsGroup'
        propertyName: eventTypeName
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - TOO_MANY_ELECTIONS
                - REPLICATION_OPLOG_WINDOW_RUNNING_OUT
                - TOO_FEW_HEALTHY_MEMBERS
                - TOO_MANY_UNHEALTHY_MEMBERS
            example: TOO_MANY_ELECTIONS
            title: ReplicaSet Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.
            items:
                description: Rules to apply when comparing an replica set against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - REPLICA_SET_NAME
                            - SHARD_NAME
                            - CLUSTER_NAME
                        example: REPLICA_SET_NAME
                        title: Replica Set Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: A Limit that triggers an alert when  exceeded. The resource returns this parameter when **eventTypeName** has not been set to `OUTSIDE_METRIC_THRESHOLD`.
            properties:
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - <
                        - '>'
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: int32
                    type: integer
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            title: Threshold
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Replica Set Threshold Alert Configuration
      type: object
    - description: Serverless metric alert configuration allows to select which serverless database metrics trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - OUTSIDE_SERVERLESS_METRIC_THRESHOLD
            example: OUTSIDE_SERVERLESS_METRIC_THRESHOLD
            title: Serverless Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        metricThreshold:
            description: Threshold for the metric that, when exceeded, triggers an alert. The metric threshold pertains to event types which reflects changes of measurements and metrics about the serverless database.
            discriminator:
                mapping:
                    SERVERLESS_AVG_COMMAND_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    SERVERLESS_AVG_READ_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    SERVERLESS_AVG_WRITE_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    SERVERLESS_CONNECTIONS: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_CONNECTIONS_PERCENT: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_DATA_SIZE_TOTAL: '#/components/schemas/DataMetricThresholdView'
                    SERVERLESS_NETWORK_BYTES_IN: '#/components/schemas/DataMetricThresholdView'
                    SERVERLESS_NETWORK_BYTES_OUT: '#/components/schemas/DataMetricThresholdView'
                    SERVERLESS_NETWORK_NUM_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_CMD: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_DELETE: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_GETMORE: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_INSERT: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_QUERY: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_OPCOUNTER_UPDATE: '#/components/schemas/RawMetricThresholdView'
                    SERVERLESS_TOTAL_READ_UNITS: '#/components/schemas/RPUMetricThresholdView'
                    SERVERLESS_TOTAL_WRITE_UNITS: '#/components/schemas/RPUMetricThresholdView'
                propertyName: metricName
            properties:
                metricName:
                    description: Human-readable label that identifies the metric against which MongoDB Cloud checks the configured **metricThreshold.threshold**.
                    type: string
                mode:
                    description: MongoDB Cloud computes the current metric value as an average.
                    enum:
                        - AVERAGE
                    type: string
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: double
                    type: number
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            required:
                - metricName
            title: Serverless Metric Threshold
            type: object
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Serverless Alert Configuration
      type: object
    - description: Flex metric alert configuration allows to select which Flex database metrics trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - OUTSIDE_FLEX_METRIC_THRESHOLD
            example: OUTSIDE_FLEX_METRIC_THRESHOLD
            title: Flex Metric Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: Matching conditions for target resources.
            items:
                description: Rules to apply when comparing an target instance against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        metricThreshold:
            description: Threshold for the metric that, when exceeded, triggers an alert. The metric threshold pertains to event types which reflects changes of measurements and metrics about the serverless database.
            discriminator:
                mapping:
                    FLEX_AVG_COMMAND_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    FLEX_AVG_READ_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    FLEX_AVG_WRITE_EXECUTION_TIME: '#/components/schemas/TimeMetricThresholdView'
                    FLEX_CONNECTIONS: '#/components/schemas/RawMetricThresholdView'
                    FLEX_CONNECTIONS_PERCENT: '#/components/schemas/RawMetricThresholdView'
                    FLEX_DATA_SIZE_TOTAL: '#/components/schemas/DataMetricThresholdView'
                    FLEX_NETWORK_BYTES_IN: '#/components/schemas/DataMetricThresholdView'
                    FLEX_NETWORK_BYTES_OUT: '#/components/schemas/DataMetricThresholdView'
                    FLEX_NETWORK_NUM_REQUESTS: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_CMD: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_DELETE: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_GETMORE: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_INSERT: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_QUERY: '#/components/schemas/RawMetricThresholdView'
                    FLEX_OPCOUNTER_UPDATE: '#/components/schemas/RawMetricThresholdView'
                propertyName: metricName
            properties:
                metricName:
                    description: Human-readable label that identifies the metric against which MongoDB Cloud checks the configured **metricThreshold.threshold**.
                    type: string
                mode:
                    description: MongoDB Cloud computes the current metric value as an average.
                    enum:
                        - AVERAGE
                    type: string
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: double
                    type: number
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            required:
                - metricName
            title: Flex Cluster Metric Threshold
            type: object
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Flex Alert Configuration
      type: object
    - description: Host metric alert configuration allows to select which Atlas streams processors trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - STREAM_PROCESSOR_STATE_IS_FAILED
            example: STREAM_PROCESSOR_STATE_IS_FAILED
            title: Stream Processor Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration.
            items:
                description: Rules to apply when comparing a stream processing instance or stream processor against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - INSTANCE_NAME
                            - PROCESSOR_NAME
                        example: INSTANCE_NAME
                        title: Streams Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Stream Processor Alert Configuration
      type: object
    - description: Stream Processor threshold alert configuration allows to select thresholds on metrics which trigger alerts and how users are notified.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
        eventTypeName:
            description: Event type that triggers an alert.
            enum:
                - OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD
            example: OUTSIDE_STREAM_PROCESSOR_METRIC_THRESHOLD
            title: Stream Processor Event Types
            type: string
        groupId:
            description: Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        id:
            description: Unique 24-hexadecimal digit string that identifies this alert configuration.
            example: 32b6e34b3d91647abb20e7b8
            pattern: ^([a-f0-9]{24})$
            readOnly: true
            type: string
        links:
            description: List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
            externalDocs:
                description: Web Linking Specification (RFC 5988)
                url: https://datatracker.ietf.org/doc/html/rfc5988
            items:
                properties:
                    href:
                        description: Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: https://cloud.mongodb.com/api/atlas
                        type: string
                    rel:
                        description: Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.
                        example: self
                        type: string
                type: object
            readOnly: true
            type: array
        matchers:
            description: List of rules that determine whether MongoDB Cloud checks an object for the alert configuration.
            items:
                description: Rules to apply when comparing a stream processing instance or stream processor against this alert configuration.
                properties:
                    fieldName:
                        description: Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
                        enum:
                            - INSTANCE_NAME
                            - PROCESSOR_NAME
                        example: INSTANCE_NAME
                        title: Streams Matcher Fields
                        type: string
                    operator:
                        description: Comparison operator to apply when checking the current metric value against **matcher[n].value**.
                        enum:
                            - EQUALS
                            - CONTAINS
                            - STARTS_WITH
                            - ENDS_WITH
                            - NOT_EQUALS
                            - NOT_CONTAINS
                            - REGEX
                        type: string
                    value:
                        description: Value to match or exceed using the specified **matchers.operator**.
                        example: event-replica-set
                        type: string
                required:
                    - fieldName
                    - operator
                    - value
                title: Matchers
                type: object
            type: array
        notifications:
            description: List that contains the targets that MongoDB Cloud sends notifications.
            items:
                description: One target that MongoDB Cloud sends notifications when an alert triggers.
                oneOf:
                    - description: Datadog notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        datadogApiKey:
                            description: |-
                                Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************a23c'
                            type: string
                        datadogRegion:
                            default: US
                            description: 'Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `"notifications.[n].typeName" : "DATADOG"`.'
                            enum:
                                - US
                                - EU
                                - US3
                                - US5
                                - AP1
                                - US1_FED
                            externalDocs:
                                description: Datadog regions
                                url: https://docs.datadoghq.com/getting_started/site/
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - DATADOG
                            type: string
                      required:
                        - typeName
                      title: Datadog Notification
                      type: object
                    - description: Email notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailAddress:
                            description: |-
                                Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "EMAIL"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:

                                - specific MongoDB Cloud users (`"notifications.[n].typeName" : "USER"`)
                                - MongoDB Cloud users with specific project roles (`"notifications.[n].typeName" : "GROUP"`)
                                - MongoDB Cloud users with specific organization roles (`"notifications.[n].typeName" : "ORG"`)
                                - MongoDB Cloud teams (`"notifications.[n].typeName" : "TEAM"`)

                                To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
                            format: email
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - EMAIL
                            type: string
                      required:
                        - typeName
                      title: Email Notification
                      type: object
                    - description: Group notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more project roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Project Roles
                                url: https://dochub.mongodb.org/core/atlas-proj-roles
                            items:
                                description: One or more project roles that receive the configured alert.
                                enum:
                                    - GROUP_BACKUP_MANAGER
                                    - GROUP_CLUSTER_MANAGER
                                    - GROUP_DATA_ACCESS_ADMIN
                                    - GROUP_DATA_ACCESS_READ_ONLY
                                    - GROUP_DATA_ACCESS_READ_WRITE
                                    - GROUP_DATABASE_ACCESS_ADMIN
                                    - GROUP_OBSERVABILITY_VIEWER
                                    - GROUP_OWNER
                                    - GROUP_READ_ONLY
                                    - GROUP_SEARCH_INDEX_EDITOR
                                    - GROUP_STREAM_PROCESSING_OWNER
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - GROUP
                            type: string
                      required:
                        - typeName
                      title: Group Notification
                      type: object
                    - description: HipChat notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notificationToken:
                            description: |-
                                HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '************************************1234'
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roomName:
                            description: 'HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "HIP_CHAT"`".'
                            example: test room
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - HIP_CHAT
                            type: string
                      required:
                        - typeName
                      title: HipChat Notification
                      type: object
                    - description: Microsoft Teams notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        microsoftTeamsWebhookUrl:
                            description: |-
                                Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `"notifications.[n].typeName" : "MICROSOFT_TEAMS"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - MICROSOFT_TEAMS
                            type: string
                      required:
                        - typeName
                      title: Microsoft Teams Notification
                      type: object
                    - description: OpsGenie notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        opsGenieApiKey:
                            description: |-
                                API Key that MongoDB Cloud needs to send this notification via Opsgenie. The resource requires this parameter when `"notifications.[n].typeName" : "OPS_GENIE"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************a111'
                            type: string
                        opsGenieRegion:
                            default: US
                            description: Opsgenie region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - OPS_GENIE
                            type: string
                      required:
                        - typeName
                      title: OpsGenie Notification
                      type: object
                    - description: Org notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        roles:
                            description: 'List that contains the one or more organization roles that receive the configured alert. This parameter is available when `"notifications.[n].typeName" : "GROUP"` or `"notifications.[n].typeName" : "ORG"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.'
                            externalDocs:
                                description: Organization Roles
                                url: https://dochub.mongodb.org/core/atlas-org-roles
                            items:
                                description: One or more organization roles that receive the configured alert.
                                enum:
                                    - ORG_OWNER
                                    - ORG_MEMBER
                                    - ORG_GROUP_CREATOR
                                    - ORG_BILLING_ADMIN
                                    - ORG_BILLING_READ_ONLY
                                    - ORG_READ_ONLY
                                type: string
                            type: array
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - ORG
                            type: string
                      required:
                        - typeName
                      title: Org Notification
                      type: object
                    - description: PagerDuty notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        region:
                            default: US
                            description: PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
                            enum:
                                - US
                                - EU
                            type: string
                        serviceKey:
                            description: |-
                                PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `"notifications.[n].typeName" : "PAGER_DUTY"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '****************************7890'
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - PAGER_DUTY
                            type: string
                      required:
                        - typeName
                      title: PagerDuty Notification
                      type: object
                    - description: Slack notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        apiToken:
                            description: "Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token. \n\n**NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:\n\n* View or edit the alert through the Atlas UI.\n\n* Query the alert for the notification through the Atlas Administration API."
                            example: '**********************************************************************abcd'
                            type: string
                        channelName:
                            description: 'Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SLACK"`.'
                            example: alerts
                            type: string
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SLACK
                            type: string
                      required:
                        - typeName
                      title: Slack Notification
                      type: object
                    - description: SMS notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        mobileNumber:
                            description: 'Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `"notifications.[n].typeName" : "SMS"`.'
                            example: "1233337892"
                            type: string
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - SMS
                            type: string
                      required:
                        - typeName
                      title: SMS Notification
                      type: object
                    - description: Team notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        teamId:
                            description: 'Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        teamName:
                            description: 'Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `"notifications.[n].typeName" : "TEAM"`.'
                            example: Atlas
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - TEAM
                            type: string
                      required:
                        - typeName
                      title: Team Notification
                      type: object
                    - description: User notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        emailEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        smsEnabled:
                            description: |-
                                Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:

                                - `"notifications.[n].typeName" : "ORG"`
                                - `"notifications.[n].typeName" : "GROUP"`
                                - `"notifications.[n].typeName" : "USER"`
                            type: boolean
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - USER
                            type: string
                        username:
                            description: 'MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `"notifications.[n].typeName" : "USER"`.'
                            format: email
                            type: string
                      required:
                        - typeName
                      title: User Notification
                      type: object
                    - description: VictorOps notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - VICTOR_OPS
                            type: string
                        victorOpsApiKey:
                            description: |-
                                API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.

                                **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:

                                * View or edit the alert through the Atlas UI.

                                * Query the alert for the notification through the Atlas Administration API.
                            example: '********************************9abc'
                            type: string
                        victorOpsRoutingKey:
                            description: 'Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `"notifications.[n].typeName" : "VICTOR_OPS"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.'
                            example: test routing
                            type: string
                      required:
                        - typeName
                      title: VictorOps Notification
                      type: object
                    - description: Webhook notification configuration for MongoDB Cloud to send information when an event triggers an alert condition.
                      properties:
                        delayMin:
                            description: Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
                            format: int32
                            type: integer
                        integrationId:
                            description: The id of the associated integration, the credentials of which to use for requests.
                            example: 32b6e34b3d91647abb20e7b8
                            type: string
                        intervalMin:
                            description: |-
                                Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.

                                PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
                            format: int32
                            minimum: 5
                            type: integer
                        notifierId:
                            description: The notifierId is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
                            example: 32b6e34b3d91647abb20e7b8
                            pattern: ^([a-f0-9]{24})$
                            type: string
                        typeName:
                            description: Human-readable label that displays the alert notification type.
                            enum:
                                - WEBHOOK
                            type: string
                        webhookSecret:
                            description: |-
                                Authentication secret for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookSecret` to a non-empty string
                                * You set a default webhookSecret either on the Integrations page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
                            externalDocs:
                                description: Integrations page
                                url: https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations
                            format: password
                            type: string
                        webhookUrl:
                            description: |-
                                Target URL for a webhook-based alert.

                                Atlas returns this value if you set `"notifications.[n].typeName" :"WEBHOOK"` and either:
                                * You set `notification.[n].webhookURL` to a non-empty string
                                * You set a default webhookUrl either on the [Integrations](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/#std-label-third-party-integrations) page, or with the [Integrations API](#tag/Third-Party-Service-Integrations/operation/createIntegration)

                                **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
                            example: https://webhook.com/****
                            type: string
                      required:
                        - typeName
                      title: Webhook Notification
                      type: object
                type: object
            type: array
        threshold:
            description: Threshold for the metric that, when exceeded, triggers an alert. The metric threshold pertains to event types which reflects changes of measurements and metrics in stream processors.
            discriminator:
                mapping:
                    STREAM_PROCESSOR_CHANGE_STREAM_LAG: '#/components/schemas/TimeMetricThresholdView'
                    STREAM_PROCESSOR_DLQ_MESSAGE_COUNT: '#/components/schemas/RawMetricThresholdView'
                    STREAM_PROCESSOR_KAFKA_LAG: '#/components/schemas/RawMetricThresholdView'
                    STREAM_PROCESSOR_OUTPUT_MESSAGE_COUNT: '#/components/schemas/RawMetricThresholdView'
                propertyName: metricName
            properties:
                metricName:
                    description: Human-readable label that identifies the metric against which MongoDB Cloud checks the configured **metricThreshold.threshold**.
                    type: string
                mode:
                    description: MongoDB Cloud computes the current metric value as an average.
                    enum:
                        - AVERAGE
                    type: string
                operator:
                    description: Comparison operator to apply when checking the current metric value.
                    enum:
                        - LESS_THAN
                        - GREATER_THAN
                    type: string
                threshold:
                    description: Value of metric that, when exceeded, triggers an alert.
                    format: double
                    type: number
                units:
                    description: Element used to express the quantity. This can be an element of time, storage capacity, and the like.
                    enum:
                        - bits
                        - Kbits
                        - Mbits
                        - Gbits
                        - bytes
                        - KB
                        - MB
                        - GB
                        - TB
                        - PB
                        - nsec
                        - msec
                        - sec
                        - min
                        - hours
                        - million minutes
                        - days
                        - requests
                        - 1000 requests
                        - GB seconds
                        - GB hours
                        - GB days
                        - RPU
                        - thousand RPU
                        - million RPU
                        - WPU
                        - thousand WPU
                        - million WPU
                        - count
                        - thousand
                        - million
                        - billion
                    type: string
            required:
                - metricName
            title: Stream Processor Metric Threshold
            type: object
        updated:
            description: Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
      required:
        - eventTypeName
        - notifications
      title: Stream Processor Metric Alert Configuration
      type: object
type: object

