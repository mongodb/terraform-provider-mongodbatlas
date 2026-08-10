package autogen_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type contextHookError struct{}

func (contextHookError) PreCreateAPICall(context.Context, config.APICallParams, []byte) (config.APICallParams, []byte, error) {
	return config.APICallParams{}, nil, errors.New("prepare create request")
}

func (contextHookError) PreUpdateAPICall(context.Context, config.APICallParams, []byte) (config.APICallParams, []byte, error) {
	return config.APICallParams{}, nil, errors.New("prepare update request")
}

func TestContextHookErrorsStopAPICall(t *testing.T) {
	t.Parallel()

	model := struct {
		Name types.String `tfsdk:"name"`
	}{Name: types.StringValue("name")}
	callParams := &config.APICallParams{}

	createResp := &resource.CreateResponse{}
	autogen.HandleCreate(t.Context(), autogen.HandleCreateReq{
		Hooks:      contextHookError{},
		Resp:       createResp,
		Plan:       &model,
		CallParams: callParams,
	})
	require.True(t, createResp.Diagnostics.HasError())
	require.Contains(t, createResp.Diagnostics[0].Detail(), "prepare create request")

	updateResp := &resource.UpdateResponse{}
	autogen.HandleUpdate(t.Context(), autogen.HandleUpdateReq{
		Hooks:      contextHookError{},
		Resp:       updateResp,
		Plan:       &model,
		CallParams: callParams,
	})
	require.True(t, updateResp.Diagnostics.HasError())
	require.Contains(t, updateResp.Diagnostics[0].Detail(), "prepare update request")
}
