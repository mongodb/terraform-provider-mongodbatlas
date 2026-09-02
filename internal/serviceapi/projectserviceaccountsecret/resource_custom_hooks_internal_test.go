package projectserviceaccountsecret

import (
	"errors"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourcePostReadAPICall(t *testing.T) {
	listBody := []byte(`{"secrets":[{"id":"secret-1","maskedSecretValue":"mdb_sa_sk_masked"}]}`)

	t.Run("secret found returns its body", func(t *testing.T) {
		result := resourcePostReadAPICall("secret-1", autogen.APICallResult{Body: listBody})
		require.NoError(t, result.Err)
		assert.JSONEq(t, `{"id":"secret-1","maskedSecretValue":"mdb_sa_sk_masked"}`, string(result.Body))
	})

	t.Run("secret absent signals not found", func(t *testing.T) {
		result := resourcePostReadAPICall("other-secret", autogen.APICallResult{Body: listBody})
		require.ErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("existing error passes through", func(t *testing.T) {
		err := errors.New("some API error")
		result := resourcePostReadAPICall("secret-1", autogen.APICallResult{Err: err})
		require.ErrorIs(t, result.Err, err)
		assert.NotErrorIs(t, result.Err, autogen.ErrNotFound)
	})

	t.Run("decode error is not classified as not found", func(t *testing.T) {
		result := resourcePostReadAPICall("secret-1", autogen.APICallResult{Body: []byte(`not-json`)})
		require.Error(t, result.Err)
		assert.NotErrorIs(t, result.Err, autogen.ErrNotFound)
	})
}
