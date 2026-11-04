package alertconfiguration_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
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
			want: `{"notifications":[{"typeName":"GROUP","username":"user"}]}`,
		},
		"non-empty values are kept": {
			body: `{"notifications":[{"typeName":"GROUP","channelName":"general"}]}`,
			want: `{"notifications":[{"typeName":"GROUP","channelName":"general"}]}`,
		},
		"body without notifications or thresholds is untouched": {
			body: `{"enabled":true}`,
			want: `{"enabled":true}`,
		},
		"sensitive values returned by the API are dropped": {
			body: `{"notifications":[{"typeName":"PAGER_DUTY","serviceKey":"echoed-from-integration","intervalMin":10}]}`,
			want: `{"notifications":[{"typeName":"PAGER_DUTY","intervalMin":10}]}`,
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
				"notifications": []any{map[string]any{"typeName": "GROUP"}},
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

func TestPreImport_rewritesLegacyDashFormat(t *testing.T) {
	hook, ok := alertconfiguration.Resource().(autogen.PreImportHook)
	require.True(t, ok)

	tests := map[string]struct {
		id   string
		want string
	}{
		"legacy dash format is rewritten": {
			id:   "64e3d2d9a25b07293b8f5e21-651dd9336afac13e1c112222",
			want: "64e3d2d9a25b07293b8f5e21/651dd9336afac13e1c112222",
		},
		"slash format passes through": {
			id:   "64e3d2d9a25b07293b8f5e21/651dd9336afac13e1c112222",
			want: "64e3d2d9a25b07293b8f5e21/651dd9336afac13e1c112222",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := hook.PreImport(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
