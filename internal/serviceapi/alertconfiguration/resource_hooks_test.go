package alertconfiguration_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/alertconfiguration"
)

func decodeBodyJSON(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var obj map[string]any
	require.NoError(t, json.Unmarshal(body, &obj))
	return obj
}

func TestPostReadAPICall_thresholdPrecedenceAndNormalization(t *testing.T) {
	hook, ok := alertconfiguration.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)

	tests := map[string]struct{ body, want string }{
		"metric threshold wins when both are returned": {
			body: `{"notifications":[],"metricThreshold":{"threshold":1},"threshold":{"threshold":2}}`,
			want: `{"notifications":[],"metricThreshold":{"threshold":1}}`,
		},
		"threshold is kept when metric threshold is absent": {
			body: `{"threshold":{"threshold":2}}`,
			want: `{"threshold":{"threshold":2}}`,
		},
		"empty strings are dropped from notifications": {
			body: `{"notifications":[{"typeName":"GROUP","channelName":"","username":"user"}]}`,
			want: `{"notifications":[{"typeName":"GROUP","username":"user","emailEnabled":false,"smsEnabled":false}]}`,
		},
		"non-empty values are kept": {
			body: `{"notifications":[{"typeName":"GROUP","channelName":"general"}]}`,
			want: `{"notifications":[{"typeName":"GROUP","channelName":"general","emailEnabled":false,"smsEnabled":false}]}`,
		},
		"body without notifications or thresholds is untouched": {
			body: `{"enabled":true}`,
			want: `{"enabled":true}`,
		},
		"sensitive values returned by the API are dropped": {
			body: `{"notifications":[{"typeName":"PAGER_DUTY","serviceKey":"echoed-from-integration","intervalMin":10}]}`,
			want: `{"notifications":[{"typeName":"PAGER_DUTY","intervalMin":10,"emailEnabled":false,"smsEnabled":false}]}`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{Body: []byte(tc.body)})
			require.NoError(t, result.Err)
			assert.Equal(t, decodeBodyJSON(t, []byte(tc.want)), decodeBodyJSON(t, result.Body))
		})
	}
}

func TestPostCreateAndUpdateAPICall_normalization(t *testing.T) {
	for _, name := range []string{"PostCreateAPICall", "PostUpdateAPICall"} {
		t.Run(name, func(t *testing.T) {
			body := []byte(`{"notifications":[{"typeName":"GROUP","teamName":""}]}`)
			var result autogen.APICallResult
			if name == "PostCreateAPICall" {
				hook, ok := alertconfiguration.Resource().(autogen.PostCreateAPICallHook)
				require.True(t, ok)
				result = hook.PostCreateAPICall(autogen.APICallResult{Body: body})
			} else {
				hook, ok := alertconfiguration.Resource().(autogen.PostUpdateAPICallHook)
				require.True(t, ok)
				result = hook.PostUpdateAPICall(autogen.APICallResult{Body: body})
			}
			require.NoError(t, result.Err)
			assert.Equal(t, map[string]any{
				"notifications": []any{map[string]any{"typeName": "GROUP", "emailEnabled": false, "smsEnabled": false}},
			}, decodeBodyJSON(t, result.Body))
		})
	}
}

func TestPostHooks_errorAndEmptyBodyUnchanged(t *testing.T) {
	res := alertconfiguration.Resource()
	createHook, _ := res.(autogen.PostCreateAPICallHook)
	updateHook, _ := res.(autogen.PostUpdateAPICallHook)
	readHook, _ := res.(autogen.PostReadAPICallHook)
	for name, result := range map[string]autogen.APICallResult{
		"error": {Err: assert.AnError},
		"empty": {Body: []byte("")},
	} {
		t.Run(name, func(t *testing.T) {
			for hookName, got := range map[string]autogen.APICallResult{
				"create": createHook.PostCreateAPICall(result),
				"update": updateHook.PostUpdateAPICall(result),
				"read":   readHook.PostReadAPICall(autogen.HandleReadReq{}, result),
			} {
				assert.Equal(t, result.Err, got.Err, hookName)
			}
		})
	}
}

func TestPreCreateAPICall_fillsRequiredZeroValues(t *testing.T) {
	body := `{"enabled":true,"eventTypeName":"OUTSIDE_METRIC_THRESHOLD","notifications":[{"typeName":"EMAIL","emailAddress":"a@b.com"}],"metricThreshold":{"metricName":"ASSERT_REGULAR","operator":"LESS_THAN","units":"RAW","mode":"AVERAGE"}}`
	params, got := alertconfiguration.Resource().(autogen.PreCreateAPICallHook).PreCreateAPICall(config.APICallParams{Method: "POST"}, []byte(body))
	require.Equal(t, "POST", params.Method)
	obj := decodeBodyJSON(t, got)
	notification := obj["notifications"].([]any)[0].(map[string]any)
	// The API requires the intervalMin and threshold keys (MISSING_ATTRIBUTE otherwise); the manual sent
	// 0 for the unknown optional+computed attributes (ValueXPointer returns a pointer to 0 on unknown).
	assert.InDelta(t, 0, notification["intervalMin"], 0)
	assert.InDelta(t, 0, obj["metricThreshold"].(map[string]any)["threshold"], 0)
	// delayMin 0, smsEnabled false and emailEnabled false are filled (the manual sent them for unknown
	// optional+computed attributes; an absent emailEnabled defaults to true in the API).
	assert.InDelta(t, 0, notification["delayMin"], 0)
	assert.Equal(t, false, notification["smsEnabled"])
	assert.Equal(t, false, notification["emailEnabled"])
	// Attributes the manual omitted (null pointers) stay omitted.
	assert.NotContains(t, notification, "notifierId")
	assert.NotContains(t, notification, "integrationId")
	// Configured values are kept.
	assert.Equal(t, true, obj["enabled"])
	assert.Equal(t, "a@b.com", notification["emailAddress"])
}

func TestPreCreateAPICall_keepsConfiguredIntervalAndThreshold(t *testing.T) {
	body := `{"notifications":[{"typeName":"GROUP","intervalMin":5}],"threshold":{"operator":"GREATER_THAN","units":"HOURS","threshold":48}}`
	_, got := alertconfiguration.Resource().(autogen.PreCreateAPICallHook).PreCreateAPICall(config.APICallParams{}, []byte(body))
	obj := decodeBodyJSON(t, got)
	notification := obj["notifications"].([]any)[0].(map[string]any)
	assert.InDelta(t, 5, notification["intervalMin"], 0)
	assert.InDelta(t, 48, obj["threshold"].(map[string]any)["threshold"], 0)
}

func TestUpgradeManualStateJSON_unwrapsThresholdLists(t *testing.T) {
	raw := `{
		"id": "aWQ=:encoded",
		"project_id": "64e3d2d9a25b07293b8f5e21",
		"event_type": "OUTSIDE_METRIC_THRESHOLD",
		"matcher": [{"field_name":"HOSTNAME_AND_PORT","operator":"EQUALS","value":"SECONDARY"}],
		"notification": [{"type_name":"GROUP","interval_min":5}],
		"metric_threshold_config": [{"threshold":99,"metric_name":"ASSERT_REGULAR","operator":"LESS_THAN","units":"RAW","mode":"AVERAGE"}],
		"threshold_config": []
	}`
	rawWithoutMatchers := `{"matcher": [], "notification": [{"type_name":"GROUP"}]}`
	got, err := alertconfiguration.UpgradeManualStateJSONForTest([]byte(raw))
	require.NoError(t, err)
	obj := map[string]any{}
	require.NoError(t, json.Unmarshal(got, &obj))
	// list-of-one becomes the object; empty list is removed (null object); lists and the rest are kept.
	metricThreshold, ok := obj["metric_threshold_config"].(map[string]any)
	require.True(t, ok, "metric_threshold_config must be an object after the upgrade")
	assert.InDelta(t, 99, metricThreshold["threshold"], 0.001)
	_, ok = obj["threshold_config"]
	assert.False(t, ok, "empty threshold_config list must be removed")
	assert.Len(t, obj["matcher"], 1)
	assert.Len(t, obj["notification"], 1)
	assert.Equal(t, "aWQ=:encoded", obj["id"])

	got, err = alertconfiguration.UpgradeManualStateJSONForTest([]byte(rawWithoutMatchers))
	require.NoError(t, err)
	obj = map[string]any{}
	require.NoError(t, json.Unmarshal(got, &obj))
	_, ok = obj["matcher"]
	assert.False(t, ok, "empty matcher list must be removed")
}
