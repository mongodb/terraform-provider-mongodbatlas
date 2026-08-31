package serviceaccountsecret_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/serviceaccountsecret"
)

func TestPostReadAPICall(t *testing.T) {
	t.Run("resource", func(t *testing.T) {
		hook := postReadHook(t, serviceaccountsecret.Resource())
		assertSecretLookup(t, func(secretID string, result autogen.APICallResult) autogen.APICallResult {
			state := &serviceaccountsecret.TFModel{SecretId: types.StringValue(secretID)}
			return hook.PostReadAPICall(autogen.HandleReadReq{State: state}, result)
		})
	})

	t.Run("data source", func(t *testing.T) {
		hook := postReadHook(t, serviceaccountsecret.DataSource())
		assertSecretLookup(t, func(secretID string, result autogen.APICallResult) autogen.APICallResult {
			state := &serviceaccountsecret.TFDSModel{SecretId: types.StringValue(secretID)}
			return hook.PostReadAPICall(autogen.HandleReadReq{State: state}, result)
		})
	})
}

func postReadHook(t *testing.T, from any) autogen.PostReadAPICallHook {
	t.Helper()
	hook, ok := from.(autogen.PostReadAPICallHook)
	require.True(t, ok, "%T must implement autogen.PostReadAPICallHook", from)
	return hook
}

// assertSecretLookup covers the secret lookup shared by the resource and data source hook.
func assertSecretLookup(t *testing.T, postRead func(secretID string, result autogen.APICallResult) autogen.APICallResult) {
	t.Helper()
	listBody := []byte(`{"secrets":[{"id":"secret-1","maskedSecretValue":"mdb_sa_sk_masked"}]}`)

	t.Run("secret found returns its body", func(t *testing.T) {
		result := postRead("secret-1", autogen.APICallResult{Body: listBody})
		require.NoError(t, result.Err)
		assert.JSONEq(t, `{"id":"secret-1","maskedSecretValue":"mdb_sa_sk_masked"}`, string(result.Body))
	})

	t.Run("secret absent signals not found", func(t *testing.T) {
		result := postRead("other-secret", autogen.APICallResult{Body: listBody})
		require.ErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("existing error passes through", func(t *testing.T) {
		err := errors.New("some API error")
		result := postRead("secret-1", autogen.APICallResult{Err: err})
		require.ErrorIs(t, result.Err, err)
		assert.NotErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("decode error is not classified as not found", func(t *testing.T) {
		result := postRead("secret-1", autogen.APICallResult{Body: []byte(`not-json`)})
		require.Error(t, result.Err)
		assert.NotErrorIs(t, result.Err, autogen.ErrNotFound)
	})
}
