package clusteradaptivesettings_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

func TestPostReadAPICall(t *testing.T) {
	t.Parallel()

	t.Run("clears stale overrides on success", func(t *testing.T) {
		t.Parallel()
		hook := clusteradaptivesettings.Resource().(autogen.PostReadAPICallHook)
		state := &clusteradaptivesettings.TFModel{
			AdaptiveSettingsOverrides: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`),
		}
		result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(`{}`)})
		require.NoError(t, result.Err)
		require.Equal(t, jsontypes.NewNormalizedNull(), state.AdaptiveSettingsOverrides)
	})

	t.Run("preserves state on error", func(t *testing.T) {
		t.Parallel()
		hook := clusteradaptivesettings.Resource().(autogen.PostReadAPICallHook)
		overrides := jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`)
		state := &clusteradaptivesettings.TFModel{
			AdaptiveSettingsOverrides: overrides,
		}
		testErr := errors.New("test error")
		result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Err: testErr})
		require.Equal(t, testErr, result.Err)
		require.Equal(t, overrides, state.AdaptiveSettingsOverrides)
	})
}
