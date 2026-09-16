package serviceaccount_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/serviceaccount"
)

func TestResourceSchema(t *testing.T) {
	hook, ok := serviceaccount.Resource().(autogen.ResourceSchemaHook)
	require.True(t, ok, "serviceaccount.Resource() must implement autogen.ResourceSchemaHook")

	s := hook.ResourceSchema(t.Context(), serviceaccount.ResourceSchema(t.Context()))
	attr, ok := s.Attributes["secret_expires_after_hours"].(schema.Int64Attribute)
	require.True(t, ok, "secret_expires_after_hours must stay an Int64Attribute")

	var hasCreateOnly, hasCreateRequired bool
	for _, modifier := range attr.PlanModifiers {
		if reflect.TypeOf(modifier) == reflect.TypeOf(customplanmodifier.CreateOnly()) {
			hasCreateOnly = true
		}
		if reflect.TypeOf(modifier) == reflect.TypeOf(customplanmodifier.RequestOnlyRequiredOnCreate()) {
			hasCreateRequired = true
		}
	}
	assert.False(t, hasCreateRequired, "omitting secret_expires_after_hours must be valid on create")
	assert.True(t, hasCreateOnly, "changing secret_expires_after_hours after creation must still fail")
}

func TestPreCreateAPICall(t *testing.T) {
	hook, ok := serviceaccount.Resource().(autogen.PreCreateAPICallHook)
	require.True(t, ok, "serviceaccount.Resource() must implement autogen.PreCreateAPICallHook")

	call := func(body string) []byte {
		_, updated := hook.PreCreateAPICall(config.APICallParams{}, []byte(body))
		return updated
	}

	t.Run("body without hours requests no initial secret", func(t *testing.T) {
		assert.JSONEq(t, `{"name":"sa","roles":["ORG_READ_ONLY"],"withoutInitialSecret":true}`,
			string(call(`{"name":"sa","roles":["ORG_READ_ONLY"]}`)))
	})

	t.Run("body with hours is unchanged", func(t *testing.T) {
		body := `{"name":"sa","secretExpiresAfterHours":24}`
		assert.JSONEq(t, body, string(call(body)))
	})

	t.Run("invalid JSON is returned unchanged", func(t *testing.T) {
		body := `not-json`
		assert.Equal(t, body, string(call(body)))
	})
}
