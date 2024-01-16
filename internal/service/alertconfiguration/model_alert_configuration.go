package alertconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20231115003/admin"
)

func NewNotificationList(tfNotificationSlice []TfNotificationModel) ([]admin.AlertsNotificationRootForGroup, error) {
	notifications := make([]admin.AlertsNotificationRootForGroup, 0)

	for i := range tfNotificationSlice {
		if !tfNotificationSlice[i].IntervalMin.IsNull() && tfNotificationSlice[i].IntervalMin.ValueInt64() > 0 {
			typeName := tfNotificationSlice[i].TypeName.ValueString()
			if strings.EqualFold(typeName, pagerDuty) || strings.EqualFold(typeName, opsGenie) || strings.EqualFold(typeName, victorOps) {
				return nil, fmt.Errorf(`'interval_min' must not be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'`)
			}
		}
	}

	for i := range tfNotificationSlice {
		n := &tfNotificationSlice[i]
		notification := admin.AlertsNotificationRootForGroup{
			ApiToken:                 n.APIToken.ValueStringPointer(),
			ChannelName:              n.ChannelName.ValueStringPointer(),
			DatadogApiKey:            n.DatadogAPIKey.ValueStringPointer(),
			DatadogRegion:            n.DatadogRegion.ValueStringPointer(),
			DelayMin:                 pointy.Int(int(n.DelayMin.ValueInt64())),
			EmailAddress:             n.EmailAddress.ValueStringPointer(),
			EmailEnabled:             n.EmailEnabled.ValueBoolPointer(),
			IntervalMin:              conversion.Int64PtrToIntPtr(n.IntervalMin.ValueInt64Pointer()),
			MobileNumber:             n.MobileNumber.ValueStringPointer(),
			OpsGenieApiKey:           n.OpsGenieAPIKey.ValueStringPointer(),
			OpsGenieRegion:           n.OpsGenieRegion.ValueStringPointer(),
			ServiceKey:               n.ServiceKey.ValueStringPointer(),
			SmsEnabled:               n.SMSEnabled.ValueBoolPointer(),
			TeamId:                   n.TeamID.ValueStringPointer(),
			TypeName:                 n.TypeName.ValueStringPointer(),
			Username:                 n.Username.ValueStringPointer(),
			VictorOpsApiKey:          n.VictorOpsAPIKey.ValueStringPointer(),
			VictorOpsRoutingKey:      n.VictorOpsRoutingKey.ValueStringPointer(),
			Roles:                    conversion.NonEmptyToPtr(n.Roles),
			MicrosoftTeamsWebhookUrl: n.MicrosoftTeamsWebhookURL.ValueStringPointer(),
			WebhookSecret:            n.WebhookSecret.ValueStringPointer(),
			WebhookUrl:               n.WebhookURL.ValueStringPointer(),
		}
		if !n.NotifierID.IsUnknown() {
			notification.NotifierId = n.NotifierID.ValueStringPointer()
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func NewThreshold(tfThresholdConfigSlice []TfThresholdConfigModel) *admin.GreaterThanRawThreshold {
	if len(tfThresholdConfigSlice) == 0 {
		return nil
	}

	v := tfThresholdConfigSlice[0]
	return &admin.GreaterThanRawThreshold{
		Operator:  v.Operator.ValueStringPointer(),
		Units:     v.Units.ValueStringPointer(),
		Threshold: pointy.Int(int(v.Threshold.ValueFloat64())),
	}
}

func NewMetricThreshold(tfMetricThresholdConfigSlice []TfMetricThresholdConfigModel) *admin.ServerlessMetricThreshold {
	if len(tfMetricThresholdConfigSlice) == 0 {
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

func NewMatcherList(tfMatcherSlice []TfMatcherModel) []map[string]interface{} {
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

func NewTFAlertConfigurationModel(apiRespConfig *admin.GroupAlertsConfig, currState *TfAlertConfigurationRSModel) TfAlertConfigurationRSModel {
	return TfAlertConfigurationRSModel{
		ID:                    currState.ID,
		ProjectID:             currState.ProjectID,
		AlertConfigurationID:  types.StringValue(conversion.SafeString(apiRespConfig.Id)),
		EventType:             types.StringValue(conversion.SafeString(apiRespConfig.EventTypeName)),
		Created:               types.StringPointerValue(conversion.TimePtrToStringPtr(apiRespConfig.Created)),
		Updated:               types.StringPointerValue(conversion.TimePtrToStringPtr(apiRespConfig.Updated)),
		Enabled:               types.BoolPointerValue(apiRespConfig.Enabled),
		MetricThresholdConfig: NewTFMetricThresholdConfigModel(apiRespConfig.MetricThreshold, currState.MetricThresholdConfig),
		ThresholdConfig:       NewTFThresholdConfigModel(apiRespConfig.Threshold, currState.ThresholdConfig),
		Notification:          NewTFNotificationModelList(apiRespConfig.GetNotifications(), currState.Notification),
		Matcher:               NewTFMatcherModelList(apiRespConfig.GetMatchers(), currState.Matcher),
	}
}

func NewTFNotificationModelList(n []admin.AlertsNotificationRootForGroup, currStateNotifications []TfNotificationModel) []TfNotificationModel {
	notifications := make([]TfNotificationModel, len(n))

	if len(n) != len(currStateNotifications) { // notifications were modified elsewhere from terraform, or import statement is being called
		for i := range n {
			value := n[i]
			notifications[i] = TfNotificationModel{
				TeamName:       conversion.StringPtrNullIfEmpty(value.TeamName),
				Roles:          value.GetRoles(),
				ChannelName:    conversion.StringPtrNullIfEmpty(value.ChannelName),
				DatadogRegion:  conversion.StringPtrNullIfEmpty(value.DatadogRegion),
				DelayMin:       types.Int64PointerValue(conversion.IntPtrToInt64Ptr(value.DelayMin)),
				EmailAddress:   conversion.StringPtrNullIfEmpty(value.EmailAddress),
				IntervalMin:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(value.IntervalMin)),
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
		newState := TfNotificationModel{
			TeamName: conversion.StringPtrNullIfEmpty(value.TeamName),
			Roles:    value.GetRoles(),
			// sentive attributes do not use value returned from API
			APIToken:                 conversion.StringNullIfEmpty(currState.APIToken.ValueString()),
			DatadogAPIKey:            conversion.StringNullIfEmpty(currState.DatadogAPIKey.ValueString()),
			OpsGenieAPIKey:           conversion.StringNullIfEmpty(currState.OpsGenieAPIKey.ValueString()),
			ServiceKey:               conversion.StringNullIfEmpty(currState.ServiceKey.ValueString()),
			VictorOpsAPIKey:          conversion.StringNullIfEmpty(currState.VictorOpsAPIKey.ValueString()),
			VictorOpsRoutingKey:      conversion.StringNullIfEmpty(currState.VictorOpsRoutingKey.ValueString()),
			WebhookURL:               conversion.StringNullIfEmpty(currState.WebhookURL.ValueString()),
			WebhookSecret:            conversion.StringNullIfEmpty(currState.WebhookSecret.ValueString()),
			MicrosoftTeamsWebhookURL: conversion.StringNullIfEmpty(currState.MicrosoftTeamsWebhookURL.ValueString()),
			NotifierID:               types.StringPointerValue(value.NotifierId),
			IntervalMin:              types.Int64PointerValue(conversion.IntPtrToInt64Ptr(value.IntervalMin)),
			DelayMin:                 types.Int64PointerValue(conversion.IntPtrToInt64Ptr(value.DelayMin)),
			EmailEnabled:             types.BoolValue(value.EmailEnabled != nil && *value.EmailEnabled),
			SMSEnabled:               types.BoolValue(value.SmsEnabled != nil && *value.SmsEnabled),
		}

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
		notifications[i] = newState
	}
	return notifications
}

func NewTFMetricThresholdConfigModel(t *admin.ServerlessMetricThreshold, currStateSlice []TfMetricThresholdConfigModel) []TfMetricThresholdConfigModel {
	if t == nil {
		return []TfMetricThresholdConfigModel{}
	}
	if len(currStateSlice) == 0 { // metric threshold was created elsewhere from terraform, or import statement is being called
		return []TfMetricThresholdConfigModel{
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
	newState := TfMetricThresholdConfigModel{
		Threshold: types.Float64Value(*t.Threshold),
	}
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
	return []TfMetricThresholdConfigModel{newState}
}

func NewTFThresholdConfigModel(t *admin.GreaterThanRawThreshold, currStateSlice []TfThresholdConfigModel) []TfThresholdConfigModel {
	if t == nil {
		return []TfThresholdConfigModel{}
	}

	if len(currStateSlice) == 0 { // threshold was created elsewhere from terraform, or import statement is being called
		return []TfThresholdConfigModel{
			{
				Operator:  conversion.StringNullIfEmpty(*t.Operator),
				Threshold: types.Float64Value(float64(*t.Threshold)), // int in new SDK but keeping float64 for backward compatibility
				Units:     conversion.StringNullIfEmpty(*t.Units),
			},
		}
	}
	currState := currStateSlice[0]
	newState := TfThresholdConfigModel{}
	if !currState.Operator.IsNull() {
		newState.Operator = conversion.StringNullIfEmpty(*t.Operator)
	}
	if !currState.Units.IsNull() {
		newState.Units = conversion.StringNullIfEmpty(*t.Units)
	}
	newState.Threshold = types.Float64Value(float64(*t.Threshold))

	return []TfThresholdConfigModel{newState}
}

func NewTFMatcherModelList(m []map[string]any, currStateSlice []TfMatcherModel) []TfMatcherModel {
	matchers := make([]TfMatcherModel, len(m))
	if len(m) != len(currStateSlice) { // matchers were modified elsewhere from terraform, or import statement is being called
		for i, matcher := range m {
			fieldName, _ := matcher["fieldName"].(string)
			operator, _ := matcher["operator"].(string)
			value, _ := matcher["value"].(string)
			matchers[i] = TfMatcherModel{
				FieldName: conversion.StringNullIfEmpty(fieldName),
				Operator:  conversion.StringNullIfEmpty(operator),
				Value:     conversion.StringNullIfEmpty(value),
			}
		}
		return matchers
	}
	for i, matcher := range m {
		currState := currStateSlice[i]
		newState := TfMatcherModel{}
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

func NewTfAlertConfigurationDSModel(apiRespConfig *admin.GroupAlertsConfig, projectID string) TFAlertConfigurationDSModel {
	return TFAlertConfigurationDSModel{
		ID: types.StringValue(conversion.EncodeStateID(map[string]string{
			EncodedIDKeyAlertID:   *apiRespConfig.Id,
			EncodedIDKeyProjectID: projectID,
		})),
		ProjectID:             types.StringValue(projectID),
		AlertConfigurationID:  types.StringValue(*apiRespConfig.Id),
		EventType:             types.StringValue(*apiRespConfig.EventTypeName),
		Created:               types.StringPointerValue(conversion.TimePtrToStringPtr(apiRespConfig.Created)),
		Updated:               types.StringPointerValue(conversion.TimePtrToStringPtr(apiRespConfig.Updated)),
		Enabled:               types.BoolPointerValue(apiRespConfig.Enabled),
		MetricThresholdConfig: NewTFMetricThresholdConfigModel(apiRespConfig.MetricThreshold, []TfMetricThresholdConfigModel{}),
		ThresholdConfig:       NewTFThresholdConfigModel(apiRespConfig.Threshold, []TfThresholdConfigModel{}),
		Notification:          NewTFNotificationModelList(apiRespConfig.GetNotifications(), []TfNotificationModel{}),
		Matcher:               NewTFMatcherModelList(apiRespConfig.GetMatchers(), []TfMatcherModel{}),
	}
}

func NewTFAlertConfigurationDSModelList(alerts []admin.GroupAlertsConfig, projectID string, definedOutputs []string) []TFAlertConfigurationDSModel {
	outputConfigurations := make([]TfAlertConfigurationOutputModel, len(definedOutputs))
	for i, output := range definedOutputs {
		outputConfigurations[i] = TfAlertConfigurationOutputModel{
			Type: types.StringValue(output),
		}
	}

	results := make([]TFAlertConfigurationDSModel, len(alerts))

	for i := 0; i < len(alerts); i++ {
		alert := alerts[i]
		label := fmt.Sprintf("%s_%d", *alert.EventTypeName, i)
		resultAlertConfigModel := NewTfAlertConfigurationDSModel(&alerts[i], projectID)
		computedOutputs := computeAlertConfigurationOutput(&alert, outputConfigurations, label)
		resultAlertConfigModel.Output = computedOutputs
		results[i] = resultAlertConfigModel
	}

	return results
}
