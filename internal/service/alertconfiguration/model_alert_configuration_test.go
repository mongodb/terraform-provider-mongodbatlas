package alertconfiguration_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/alertconfiguration"
	"go.mongodb.org/atlas-sdk/v20231115001/admin"
)

type sdkToTFNotificationModelTestCase struct {
	name                      string
	SDKResp                   *[]admin.AlertsNotificationRootForGroup
	currentStateNotifications []alertconfiguration.TfNotificationModel
	expectedTFModel           []alertconfiguration.TfNotificationModel
}

type sdkToTFMetricThresholdModelTestCase struct {
	name                        string
	SDKResp                     *admin.ServerlessMetricThreshold
	currentStateMetricThreshold []alertconfiguration.TfMetricThresholdConfigModel
	expectedTFModel             []alertconfiguration.TfMetricThresholdConfigModel
}

type sdkToTFMatcherModelTestCase struct {
	name                string
	SDKResp             []map[string]interface{}
	currentStateMatcher []alertconfiguration.TfMatcherModel
	expectedTFModel     []alertconfiguration.TfMatcherModel
}

type sdkToTFThresholdConfigModelTestCase struct {
	name                        string
	SDKResp                     *admin.GreaterThanRawThreshold
	currentStateThresholdConfig []alertconfiguration.TfThresholdConfigModel
	expectedTFModel             []alertconfiguration.TfThresholdConfigModel
}

type sdkToTFAlertConfigurationModelTestCase struct {
	name                           string
	SDKResp                        *admin.GroupAlertsConfig
	currentStateAlertConfiguration *alertconfiguration.TfAlertConfigurationRSModel
	expectedTFModel                alertconfiguration.TfAlertConfigurationRSModel
}

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

func TestNotificationSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFNotificationModelTestCase{
		{
			name: "Complete SDK response",
			SDKResp: &[]admin.AlertsNotificationRootForGroup{
				{
					TypeName:     admin.PtrString(group),
					IntervalMin:  admin.PtrInt(intervalMin),
					DelayMin:     admin.PtrInt(delayMin),
					SmsEnabled:   admin.PtrBool(disabled),
					EmailEnabled: admin.PtrBool(enabled),
					Roles: []string{
						"GROUP_DATA_ACCESS_READ_ONLY",
						"GROUP_CLUSTER_MANAGER",
						"GROUP_DATA_ACCESS_ADMIN",
					},
				},
			},
			currentStateNotifications: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(previousIntervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					Roles: []string{
						"GROUP_DATA_ACCESS_READ_ONLY",
						"GROUP_CLUSTER_MANAGER",
						"GROUP_DATA_ACCESS_ADMIN",
					},
				},
			},
			expectedTFModel: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(intervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					Roles: []string{
						"GROUP_DATA_ACCESS_READ_ONLY",
						"GROUP_CLUSTER_MANAGER",
						"GROUP_DATA_ACCESS_ADMIN",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFNotificationModelList(*tc.SDKResp, tc.currentStateNotifications)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestMetricThresholdSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFMetricThresholdModelTestCase{
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
					Units:      types.StringValue(units),
					Mode:       types.StringValue(mode),
				},
			},
			expectedTFModel: []alertconfiguration.TfMetricThresholdConfigModel{
				{
					Threshold:  types.Float64Value(threshold),
					MetricName: types.StringValue("ASSERT_REGULAR"),
					Operator:   types.StringValue(operator),
					Units:      types.StringValue(units),
					Mode:       types.StringValue(mode),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFMetricThresholdConfigModel(tc.SDKResp, tc.currentStateMetricThreshold)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestThresholdConfigSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFThresholdConfigModelTestCase{
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
					Operator:  types.StringValue("LESS_THAN"),
					Units:     types.StringValue("MINUTES"),
				},
			},
			expectedTFModel: []alertconfiguration.TfThresholdConfigModel{
				{
					Threshold: types.Float64Value(1.0),
					Operator:  types.StringValue("LESS_THAN"),
					Units:     types.StringValue("HOURS"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := alertconfiguration.NewTFThresholdConfigModel(tc.SDKResp, tc.currentStateThresholdConfig)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestMatcherSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFMatcherModelTestCase{
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
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestAlertConfigurationSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFAlertConfigurationModelTestCase{
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
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

type tfToSDKNotificationModelTestCase struct {
	name           string
	expectedSDKReq *[]admin.AlertsNotificationRootForGroup
	tfModel        []alertconfiguration.TfNotificationModel
}

type tfToSDKMetricThresholdModelTestCase struct {
	name           string
	expectedSDKReq *admin.ServerlessMetricThreshold
	tfModel        []alertconfiguration.TfMetricThresholdConfigModel
}

type tfToSDKMatcherModelTestCase struct {
	name           string
	expectedSDKReq []map[string]interface{}
	tfModel        []alertconfiguration.TfMatcherModel
}

type tfToSDKThresholdModelTestCase struct {
	name           string
	expectedSDKReq *admin.GreaterThanRawThreshold
	tfModel        []alertconfiguration.TfThresholdConfigModel
}

func TestNotificationTFModelToSDK(t *testing.T) {
	testCases := []tfToSDKNotificationModelTestCase{
		{
			name: "Complete TF model",
			tfModel: []alertconfiguration.TfNotificationModel{
				{
					TypeName:     types.StringValue(group),
					IntervalMin:  types.Int64Value(int64(intervalMin)),
					DelayMin:     types.Int64Value(int64(delayMin)),
					SMSEnabled:   types.BoolValue(disabled),
					EmailEnabled: types.BoolValue(enabled),
					Roles: []string{
						"GROUP_DATA_ACCESS_READ_ONLY",
						"GROUP_CLUSTER_MANAGER",
						"GROUP_DATA_ACCESS_ADMIN",
					},
				},
			},
			expectedSDKReq: &[]admin.AlertsNotificationRootForGroup{
				{
					TypeName:     admin.PtrString(group),
					IntervalMin:  admin.PtrInt(intervalMin),
					DelayMin:     admin.PtrInt(delayMin),
					SmsEnabled:   admin.PtrBool(disabled),
					EmailEnabled: admin.PtrBool(enabled),
					Roles: []string{
						"GROUP_DATA_ACCESS_READ_ONLY",
						"GROUP_CLUSTER_MANAGER",
						"GROUP_DATA_ACCESS_ADMIN",
					},
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
	testCases := []tfToSDKThresholdModelTestCase{
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
	testCases := []tfToSDKMetricThresholdModelTestCase{
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
	testCases := []tfToSDKMatcherModelTestCase{
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
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}
