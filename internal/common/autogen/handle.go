package autogen

import (
	"context"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorReadingAPIResponse    = "error reading API response"
	errorProcessingAPIResponse = "error processing API response"
	errorBuildingAPIRequest    = "error building API request"
)

func HandleCreate(ctx context.Context, resp *resource.CreateResponse, client *config.MongoDBClient, plan any, callParams *config.APICallParams) {
	reqBody, err := Marshal(plan, false)
	if err != nil {
		resp.Diagnostics.AddError(errorBuildingAPIRequest, err.Error())
		return
	}
	callParams.Body = reqBody
	apiResp, err := client.UntypedAPICall(ctx, callParams)
	if err != nil {
		resp.Diagnostics.AddError("error during create operation", err.Error())
		return
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(respBody, plan); err != nil {
		resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	if err := ResolveUnknowns(plan); err != nil {
		resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func HandleRead(ctx context.Context, resp *resource.ReadResponse, client *config.MongoDBClient, state any, callParams *config.APICallParams) {
	apiResp, err := client.UntypedAPICall(ctx, callParams)
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during get operation", err.Error())
		return
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}

	// Use the current state as the base model to set the response state
	if err := Unmarshal(respBody, state); err != nil {
		resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	if err := ResolveUnknowns(state); err != nil {
		resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
