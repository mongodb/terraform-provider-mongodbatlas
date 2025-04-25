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

type HandleCreateReq struct {
	Resp       *resource.CreateResponse
	Client     *config.MongoDBClient
	Plan       any
	CallParams *config.APICallParams
	Wait       *WaitReq
}

func HandleCreate(ctx context.Context, req HandleCreateReq) {
	reqBody, err := Marshal(req.Plan, false)
	if err != nil {
		req.Resp.Diagnostics.AddError(errorBuildingAPIRequest, err.Error())
		return
	}
	req.CallParams.Body = reqBody
	apiResp, err := req.Client.UntypedAPICall(ctx, req.CallParams)
	if err != nil {
		req.Resp.Diagnostics.AddError("error during create operation", err.Error())
		return
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		req.Resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(respBody, req.Plan); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	req.Resp.Diagnostics.Append(req.Resp.State.Set(ctx, req.Plan)...)
}

type HandleReadReq struct {
	Resp       *resource.ReadResponse
	Client     *config.MongoDBClient
	State      any
	CallParams *config.APICallParams
}

func HandleRead(ctx context.Context, req HandleReadReq) {
	apiResp, err := req.Client.UntypedAPICall(ctx, req.CallParams)
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			req.Resp.State.RemoveResource(ctx)
			return
		}
		req.Resp.Diagnostics.AddError("error during get operation", err.Error())
		return
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		req.Resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}

	// Use the current state as the base model to set the response state
	if err := Unmarshal(respBody, req.State); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	if err := ResolveUnknowns(req.State); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}

	req.Resp.Diagnostics.Append(req.Resp.State.Set(ctx, req.State)...)
}

type HandleUpdateReq struct {
	Resp       *resource.UpdateResponse
	Client     *config.MongoDBClient
	Plan       any
	CallParams *config.APICallParams
	Wait       *WaitReq
}

func HandleUpdate(ctx context.Context, req HandleUpdateReq) {
	reqBody, err := Marshal(req.Plan, true)
	if err != nil {
		req.Resp.Diagnostics.AddError(errorBuildingAPIRequest, err.Error())
		return
	}
	req.CallParams.Body = reqBody
	apiResp, err := req.Client.UntypedAPICall(ctx, req.CallParams)
	if err != nil {
		req.Resp.Diagnostics.AddError("error during update operation", err.Error())
		return
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		req.Resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(respBody, req.Plan); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		req.Resp.Diagnostics.AddError(errorProcessingAPIResponse, err.Error())
		return
	}
	req.Resp.Diagnostics.Append(req.Resp.State.Set(ctx, req.Plan)...)
}

type HandleDeleteReq struct {
	Resp       *resource.DeleteResponse
	Client     *config.MongoDBClient
	State      any
	CallParams *config.APICallParams
	Wait       *WaitReq
}

func HandleDelete(ctx context.Context, req HandleDeleteReq) {
	if _, err := req.Client.UntypedAPICall(ctx, req.CallParams); err != nil {
		req.Resp.Diagnostics.AddError("error during delete", err.Error())
		return
	}
	if req.Wait != nil {
		if _, err := waitForChanges(req.Wait); err != nil {
			req.Resp.Diagnostics.AddError("error waiting for changes", err.Error())
			return
		}
	}
}
