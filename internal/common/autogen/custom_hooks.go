package autogen

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Set of hooks which can be implemented at the resource/data source struct level to override default behavior.

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

type SchemaExtensionHook interface {
	ExtendSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema
}
