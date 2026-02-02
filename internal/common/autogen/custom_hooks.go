package autogen

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Set of hooks which can be implemented at the resource/data source struct level to override default behavior.

// SchemaExtensionHook allows extending the auto-generated schema with additional attributes.
// Implement this hook to add computed attributes like 'id' without modifying the generated schema file.
type SchemaExtensionHook interface {
	ExtendSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema
}

// PostStateSetHook allows setting additional state attributes after the main state has been set.
// This is useful for setting computed attributes that are derived from other state values.
type PostStateSetHook interface {
	PostStateSet(ctx context.Context, state *tfsdk.State, model any) diag.Diagnostics
}

// StateModelHook allows using an extended model for state operations.
// Implement this to add computed attributes like 'id' that aren't in the base auto-generated model.
// This hook works together with SchemaExtensionHook - the schema extension adds the attribute definition,
// and this hook provides the extended model that includes the field for that attribute.
type StateModelHook interface {
	// NewStateModel returns a pointer to a new extended model instance for reading state.
	// The extended model should embed the base TFModel and add additional fields.
	NewStateModel() any
	// PrepareForState prepares the model before setting state.
	// It receives the model (either base or extended) and returns the extended model with computed values populated.
	PrepareForState(model any) any
}

type PreReadAPICallHook interface {
	PreReadAPICall(callParams config.APICallParams) config.APICallParams
}

type PostReadAPICallHook interface {
	PostReadAPICall(HandleReadReq, APICallResult) APICallResult
}

type PreCreateAPICallHook interface {
	PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (config.APICallParams, []byte)
}

type PostCreateAPICallHook interface {
	PostCreateAPICall(APICallResult) APICallResult
}

type PreDeleteAPICallHook interface {
	PreDeleteAPICall(callParams config.APICallParams) config.APICallParams
}

type PostDeleteAPICallHook interface {
	PostDeleteAPICall(APICallResult) APICallResult
}

type PreUpdateAPICallHook interface {
	PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (config.APICallParams, []byte)
}

type PostUpdateAPICallHook interface {
	PostUpdateAPICall(APICallResult) APICallResult
}
