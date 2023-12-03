package alertconfiguration_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/alertconfiguration"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

const (
	group               string  = "GROUP"
	previousIntervalMin int     = 10
	intervalMin         int     = 5
	delayMin            int     = 0
	enabled             bool    = true
	disabled            bool    = false
	previousOperator    string  = "MORE_THAN"
	operator            string  = "LESS_THAN"
	threshold           float64 = 99.0
	units               string  = "RAW"
	mode                string  = "AVERAGE"
)

var (
	roles = []string{"GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"}
)

func TestNotificationSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                      string
		SDKResp                   *[]admin.AlertsNotificationRootForGroup
		currentStateNotifications []alertconfiguration.TfNotificationModel
		expectedTFModel           []alertconfiguration.TfNotificationModel
	}{
		{
			name: "Complete SDK response",
			SDKResp: &[]admin.AlertsNotificationRootForGroup{
				{
					TypeName:     admin.PtrString(group),
					IntervalMin:  admin.PtrInt(intervalMin),
					DelayMin:     admin.PtrInt(delayMin),
					SmsEnabled:   admin.PtrBool(disabled),
					EmailEnabled: admin.PtrBool(enabled),
					ChannelName:  admin.PtrString("#channel"),
					Roles:        roles,
					ApiToken:     admin.PtrString("newApiToken"),
				},
			},
			currentStateNotifications: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(previousIntervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					APIToken:     types.StringValue("apiToken"),
					Roles:        roles,
				},
			},
			expectedTFModel: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(intervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					ChannelName:  types.StringNull(),
					APIToken:     types.StringValue("apiToken"),
					Roles:        roles,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFNotificationModelList(*tc.SDKResp, tc.currentStateNotifications)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestMetricThresholdSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                        string
		SDKResp                     *admin.ServerlessMetricThreshold
		currentStateMetricThreshold []alertconfiguration.TfMetricThresholdConfigModel
		expectedTFModel             []alertconfiguration.TfMetricThresholdConfigModel
	}{
		{
			name: "Complete SDK response",
			SDKResp: &admin.ServerlessMetricThreshold{
				MetricName: "ASSERT_REGULAR",
				Operator:   admin.PtrString(operator),
				Threshold:  admin.PtrFloat64(threshold),
				Units:      admin.PtrString(units),
				Mode:       admin.PtrString(mode),
			},
			currentStateMetricThreshold: []alertconfiguration.TfMetricThresholdConfigModel{
				{
					Threshold:  types.Float64Value(threshold),
					MetricName: types.StringValue("ASSERT_REGULAR"),
					Operator:   types.StringValue(previousOperator),
					Units:      types.StringNull(),
					Mode:       types.StringValue(mode),
				},
			},
			expectedTFModel: []alertconfiguration.TfMetricThresholdConfigModel{
				{
					Threshold:  types.Float64Value(threshold),
					MetricName: types.StringValue("ASSERT_REGULAR"),
					Operator:   types.StringValue(operator),
					Units:      types.StringNull(),
					Mode:       types.StringValue(mode),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFMetricThresholdConfigModel(tc.SDKResp, tc.currentStateMetricThreshold)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestThresholdConfigSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                        string
		SDKResp                     *admin.GreaterThanRawThreshold
		currentStateThresholdConfig []alertconfiguration.TfThresholdConfigModel
		expectedTFModel             []alertconfiguration.TfThresholdConfigModel
	}{
		{
			name: "Complete SDK response",
			SDKResp: &admin.GreaterThanRawThreshold{
				Threshold: admin.Int64PtrToIntPtr(admin.PtrInt64(1.0)),
				Operator:  admin.PtrString("LESS_THAN"),
				Units:     admin.PtrString("HOURS"),
			},
			currentStateThresholdConfig: []alertconfiguration.TfThresholdConfigModel{
				{
					Threshold: types.Float64Value(1.0),
					Operator:  types.StringNull(),
					Units:     types.StringValue("MINUTES"),
				},
			},
			expectedTFModel: []alertconfiguration.TfThresholdConfigModel{
				{
					Threshold: types.Float64Value(1.0),
					Operator:  types.StringNull(),
					Units:     types.StringValue("HOURS"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFThresholdConfigModel(tc.SDKResp, tc.currentStateThresholdConfig)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestMatcherSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                string
		SDKResp             []map[string]interface{}
		currentStateMatcher []alertconfiguration.TfMatcherModel
		expectedTFModel     []alertconfiguration.TfMatcherModel
	}{
		{
			name: "Complete SDK response",
			SDKResp: []map[string]interface{}{{
				"fieldName": "HOSTNAME",
				"operator":  "EQUALS",
				"value":     "PRIMARY",
			},
			},
			currentStateMatcher: []alertconfiguration.TfMatcherModel{
				{
					FieldName: types.StringValue("HOSTNAME"),
					Operator:  types.StringValue("EQUALS"),
					Value:     types.StringValue("SECONDARY"),
				},
			},
			expectedTFModel: []alertconfiguration.TfMatcherModel{
				{
					FieldName: types.StringValue("HOSTNAME"),
					Operator:  types.StringValue("EQUALS"),
					Value:     types.StringValue("PRIMARY"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFMatcherModelList(tc.SDKResp, tc.currentStateMatcher)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestAlertConfigurationSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                           string
		SDKResp                        *admin.GroupAlertsConfig
		currentStateAlertConfiguration *alertconfiguration.TfAlertConfigurationRSModel
		expectedTFModel                alertconfiguration.TfAlertConfigurationRSModel
	}{
		{
			name: "Complete SKD response",
			SDKResp: &admin.GroupAlertsConfig{
				Enabled:       admin.PtrBool(true),
				EventTypeName: admin.PtrString("EventType"),
				GroupId:       admin.PtrString("projectId"),
				Id:            admin.PtrString("alertConfigurationId"),
			},
			currentStateAlertConfiguration: &alertconfiguration.TfAlertConfigurationRSModel{
				ID:                    types.StringValue("id"),
				ProjectID:             types.StringValue("projectId"),
				AlertConfigurationID:  types.StringValue("alertConfigurationId"),
				EventType:             types.StringValue("EventType"),
				Matcher:               []alertconfiguration.TfMatcherModel{},
				MetricThresholdConfig: []alertconfiguration.TfMetricThresholdConfigModel{},
				ThresholdConfig:       []alertconfiguration.TfThresholdConfigModel{},
				Notification:          []alertconfiguration.TfNotificationModel{},
				Enabled:               types.BoolValue(true),
			},
			expectedTFModel: alertconfiguration.TfAlertConfigurationRSModel{
				ID:                    types.StringValue("id"),
				ProjectID:             types.StringValue("projectId"),
				AlertConfigurationID:  types.StringValue("alertConfigurationId"),
				EventType:             types.StringValue("EventType"),
				Matcher:               []alertconfiguration.TfMatcherModel{},
				MetricThresholdConfig: []alertconfiguration.TfMetricThresholdConfigModel{},
				ThresholdConfig:       []alertconfiguration.TfThresholdConfigModel{},
				Notification:          []alertconfiguration.TfNotificationModel{},
				Enabled:               types.BoolValue(true),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFAlertConfigurationModel(tc.SDKResp, tc.currentStateAlertConfiguration)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestNotificationTFModelToSDK(t *testing.T) {
	testCases := []struct {
		name           string
		expectedSDKReq *[]admin.AlertsNotificationRootForGroup
		tfModel        []alertconfiguration.TfNotificationModel
	}{
		{
			name: "Complete TF model",
			tfModel: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(intervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					Roles:        roles,
				},
			},
			expectedSDKReq: &[]admin.AlertsNotificationRootForGroup{
				{
					TypeName:     admin.PtrString(group),
					IntervalMin:  admin.PtrInt(intervalMin),
					DelayMin:     admin.PtrInt(delayMin),
					SmsEnabled:   admin.PtrBool(disabled),
					EmailEnabled: admin.PtrBool(enabled),
					Roles:        roles,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult, _ := alertconfiguration.NewNotificationList(tc.tfModel)
			if !reflect.DeepEqual(apiReqResult, *tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func TestThresholdTFModelToSDK(t *testing.T) {
	testCases := []struct {
		name           string
		expectedSDKReq *admin.GreaterThanRawThreshold
		tfModel        []alertconfiguration.TfThresholdConfigModel
	}{
		{
			name:           "Empty TF model",
			tfModel:        []alertconfiguration.TfThresholdConfigModel{},
			expectedSDKReq: nil,
		},
		{
			name: "Complete TF model",
			tfModel: []alertconfiguration.TfThresholdConfigModel{
				{
					Threshold: types.Float64Value(1.0),
					Operator:  types.StringValue("LESS_THAN"),
					Units:     types.StringValue("MINUTES"),
				},
			},
			expectedSDKReq: &admin.GreaterThanRawThreshold{
				Threshold: admin.Int64PtrToIntPtr(admin.PtrInt64(1.0)),
				Operator:  admin.PtrString("LESS_THAN"),
				Units:     admin.PtrString("MINUTES"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult := alertconfiguration.NewThreshold(tc.tfModel)
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func TestMetricThresholdTFModelToSDK(t *testing.T) {
	testCases := []struct {
		name           string
		expectedSDKReq *admin.ServerlessMetricThreshold
		tfModel        []alertconfiguration.TfMetricThresholdConfigModel
	}{
		{
			name:           "Empty TF model",
			tfModel:        []alertconfiguration.TfMetricThresholdConfigModel{},
			expectedSDKReq: nil,
		},
		{
			name: "Complete TF model",
			tfModel: []alertconfiguration.TfMetricThresholdConfigModel{
				{
					Threshold:  types.Float64Value(threshold),
					MetricName: types.StringValue("ASSERT_REGULAR"),
					Operator:   types.StringValue(operator),
					Units:      types.StringValue(units),
					Mode:       types.StringValue(mode),
				},
			},
			expectedSDKReq: &admin.ServerlessMetricThreshold{
				MetricName: "ASSERT_REGULAR",
				Operator:   admin.PtrString(operator),
				Threshold:  admin.PtrFloat64(threshold),
				Units:      admin.PtrString(units),
				Mode:       admin.PtrString(mode),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult := alertconfiguration.NewMetricThreshold(tc.tfModel)
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func TestMatcherTFModelToSDK(t *testing.T) {
	testCases := []struct {
		name           string
		expectedSDKReq []map[string]interface{}
		tfModel        []alertconfiguration.TfMatcherModel
	}{
		{
			name:           "Empty TF model",
			tfModel:        []alertconfiguration.TfMatcherModel{},
			expectedSDKReq: make([]map[string]interface{}, 0),
		},
		{
			name: "Complete TF model",
			tfModel: []alertconfiguration.TfMatcherModel{
				{
					FieldName: types.StringValue("HOSTNAME"),
					Operator:  types.StringValue("EQUALS"),
					Value:     types.StringValue("PRIMARY"),
				},
			},
			expectedSDKReq: []map[string]interface{}{{
				"fieldName": "HOSTNAME",
				"operator":  "EQUALS",
				"value":     "PRIMARY",
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult := alertconfiguration.NewMatcherList(tc.tfModel)
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedSDKReq, apiReqResult)
			}
		})
	}
}

func TestAlertConfigurationSdkToTFDSModel(t *testing.T) {
	testCases := []struct {
		name            string
		apiRespConfig   *admin.GroupAlertsConfig
		projectID       string
		expectedTFModel alertconfiguration.TFAlertConfigurationDSModel
	}{
		{
			name: "Complete SDK model",
			apiRespConfig: &admin.GroupAlertsConfig{
				Enabled:       admin.PtrBool(true),
				EventTypeName: admin.PtrString("EventType"),
				GroupId:       admin.PtrString("projectId"),
				Id:            admin.PtrString("alertConfigurationId"),
			},
			projectID: "123",
			expectedTFModel: alertconfiguration.TFAlertConfigurationDSModel{
				ID: types.StringValue(conversion.EncodeStateID(map[string]string{
					"id":         "alertConfigurationId",
					"project_id": "123",
				})),
				ProjectID:             types.StringValue("123"),
				AlertConfigurationID:  types.StringValue("alertConfigurationId"),
				EventType:             types.StringValue("EventType"),
				Enabled:               types.BoolValue(true),
				Matcher:               []alertconfiguration.TfMatcherModel{},
				MetricThresholdConfig: []alertconfiguration.TfMetricThresholdConfigModel{},
				ThresholdConfig:       []alertconfiguration.TfThresholdConfigModel{},
				Notification:          []alertconfiguration.TfNotificationModel{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTfAlertConfigurationDSModel(tc.apiRespConfig, tc.projectID)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestAlertConfigurationSdkToDSModelList(t *testing.T) {
	testCases := []struct {
		name            string
		projectID       string
		definedOutputs  []string
		alerts          []admin.GroupAlertsConfig
		expectedTfModel []alertconfiguration.TFAlertConfigurationDSModel
	}{
		{
			name: "Complete SDK model",
			alerts: []admin.GroupAlertsConfig{
				{
					Enabled:       admin.PtrBool(true),
					EventTypeName: admin.PtrString("EventType"),
					GroupId:       admin.PtrString("projectId"),
					Id:            admin.PtrString("alertConfigurationId"),
				},
			},
			projectID:      "projectId",
			definedOutputs: []string{"resource_hcl"},
			expectedTfModel: []alertconfiguration.TFAlertConfigurationDSModel{
				{
					ID: types.StringValue(conversion.EncodeStateID(map[string]string{
						"id":         "alertConfigurationId",
						"project_id": "projectId",
					})),
					ProjectID:             types.StringValue("projectId"),
					AlertConfigurationID:  types.StringValue("alertConfigurationId"),
					EventType:             types.StringValue("EventType"),
					Enabled:               types.BoolValue(true),
					Matcher:               []alertconfiguration.TfMatcherModel{},
					MetricThresholdConfig: []alertconfiguration.TfMetricThresholdConfigModel{},
					ThresholdConfig:       []alertconfiguration.TfThresholdConfigModel{},
					Notification:          []alertconfiguration.TfNotificationModel{},
					Output: []alertconfiguration.TfAlertConfigurationOutputModel{
						{
							Type:  types.StringValue("resource_hcl"),
							Label: types.StringValue("EventType_0"),
							Value: types.StringValue("resource \"mongodbatlas_alert_configuration\" \"EventType_0\" {\n  project_id = \"projectId\"\n  event_type = \"EventType\"\n  enabled    = true\n}\n"),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFAlertConfigurationDSModelList(tc.alerts, tc.projectID, tc.definedOutputs)
			if !reflect.DeepEqual(resultModel, tc.expectedTfModel) {
				t.Errorf("created terraform model did not match expected output.\n Expected: %s \n Result: %s", tc.expectedTfModel, resultModel)
			}
		})
	}
}
