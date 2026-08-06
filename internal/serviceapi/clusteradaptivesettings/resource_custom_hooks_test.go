package clusteradaptivesettings_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

func TestAdaptiveSettingsRequestHooks(t *testing.T) {
	t.Parallel()

	resourceInstance := clusteradaptivesettings.Resource()
	createHook, ok := resourceInstance.(autogen.PreCreateAPICallHook)
	require.True(t, ok)
	updateHook, ok := resourceInstance.(autogen.PreUpdateAPICallHook)
	require.True(t, ok)

	testCases := map[string]struct {
		request  string
		expected string
	}{
		"omitted overrides reset all": {
			request:  `{}`,
			expected: `{"adaptiveSettingsOverrides":null}`,
		},
		"null overrides reset all": {
			request:  `{"adaptiveSettingsOverrides":null}`,
			expected: `{"adaptiveSettingsOverrides":null}`,
		},
		"empty overrides reset every known key": {
			request:  `{"adaptiveSettingsOverrides":{}}`,
			expected: `{"adaptiveSettingsOverrides":{"OVERLOAD_PROTECTION":null,"SEARCH_OVERLOAD_PROTECTION":null}}`,
		},
		"configured override resets omitted known key": {
			request:  `{"adaptiveSettingsOverrides":{"OVERLOAD_PROTECTION":true}}`,
			expected: `{"adaptiveSettingsOverrides":{"OVERLOAD_PROTECTION":true,"SEARCH_OVERLOAD_PROTECTION":null}}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			callParams := config.APICallParams{Method: "PATCH"}
			createParams, createBody := createHook.PreCreateAPICall(callParams, []byte(testCase.request))
			require.Equal(t, callParams, createParams)
			require.JSONEq(t, testCase.expected, string(createBody))

			updateParams, updateBody := updateHook.PreUpdateAPICall(callParams, []byte(testCase.request))
			require.Equal(t, callParams, updateParams)
			require.JSONEq(t, testCase.expected, string(updateBody))
		})
	}
}
