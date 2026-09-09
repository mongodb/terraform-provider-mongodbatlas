package searchdeploymentapi_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/searchdeploymentapi"
)

func TestPostReadAPICall(t *testing.T) {
	hook, ok := searchdeploymentapi.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok, "resource must implement autogen.PostReadAPICallHook")

	t.Run("empty JSON body signals not found", func(t *testing.T) {
		result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{Body: []byte("{}")})
		require.ErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("empty body signals not found", func(t *testing.T) {
		result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{Body: nil})
		require.ErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("real body passes through", func(t *testing.T) {
		body := []byte(`{"stateName":"IDLE"}`)
		result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{Body: body})
		require.NoError(t, result.Err)
		assert.Equal(t, body, result.Body)
	})

	t.Run("existing error passes through", func(t *testing.T) {
		err := errors.New("some API error")
		result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{Err: err})
		require.ErrorIs(t, result.Err, err)
		assert.NotErrorIs(t, result.Err, autogen.ErrNotFound)
	})
}
