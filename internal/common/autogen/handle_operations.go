package autogen

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorReadingAPIResponse    = "error reading API response"
	errorProcessingAPIResponse = "error processing API response"
	errorBuildingAPIRequest    = "error building API request"
)

type WaitReq struct {
	CallParams        *config.APICallParams
	StateAttribute    string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}

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
	respBody, err := callAPIWithBody(ctx, req.Client, req.CallParams, reqBody)
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
	if err := handleWait(ctx, req.Wait, req.Client, req.Plan); err != nil {
		req.Resp.Diagnostics.AddError("error waiting for changes", err.Error())
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
	respBody, apiResp, err := callAPIWithoutBody(ctx, req.Client, req.CallParams)
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			req.Resp.State.RemoveResource(ctx)
			return
		}
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
	respBody, err := callAPIWithBody(ctx, req.Client, req.CallParams, reqBody)
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
	if err := handleWait(ctx, req.Wait, req.Client, req.Plan); err != nil {
		req.Resp.Diagnostics.AddError("error waiting for changes", err.Error())
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
	if _, _, err := callAPIWithoutBody(ctx, req.Client, req.CallParams); err != nil {
		req.Resp.Diagnostics.AddError(errorReadingAPIResponse, err.Error())
		return
	}
	// don't consider wait errors as it can happen that the resource is already deleted
	_ = handleWait(ctx, req.Wait, req.Client, req.State)
}

// handleWait waits until a long-running operation is done if needed.
// It also updates the model with the latest JSON response from the API.
func handleWait(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any) error {
	if wait == nil {
		return nil
	}
	respBody, err := waitForChanges(ctx, wait, client)
	if err != nil {
		return err
	}
	if respBody == nil {
		return nil
	}
	if err := Unmarshal(respBody, model); err != nil {
		return err
	}
	return ResolveUnknowns(model)
}

// waitForChanges waits until a long-running operation is done.
// It returns the latest JSON response from the API so it can be used to update the response state.
// TODO: This is a basic implementation, it will be replaced in CLOUDP-314960.
func waitForChanges(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) ([]byte, error) {
	time.Sleep(time.Duration(wait.TimeoutSeconds) * time.Second) // TODO: TimeoutSeconds is temporarily used to allow time to destroy the resource until autogen long-running operations are supported in CLOUDP-314960
	respBody, _, err := callAPIWithoutBody(ctx, client, wait.CallParams)
	return respBody, err
}

// callAPIWithBody makes a request to the API with the given request body and returns the response body.
// It is used for POST, PUT, and PATCH requests where a request body is required.
func callAPIWithBody(ctx context.Context, client *config.MongoDBClient, callParams *config.APICallParams, reqBody []byte) ([]byte, error) {
	callParams.Body = reqBody
	apiResp, err := client.UntypedAPICall(ctx, callParams)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// callAPIWithoutBody makes a request to the API without a request body and returns the response body.
// It is used for GET or DELETE requests where no request body is required.
func callAPIWithoutBody(ctx context.Context, client *config.MongoDBClient, callParams *config.APICallParams) ([]byte, *http.Response, error) {
	callParams.Body = nil
	apiResp, err := client.UntypedAPICall(ctx, callParams)
	if err != nil {
		return nil, apiResp, err
	}
	respBody, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		return nil, apiResp, err
	}
	return respBody, apiResp, nil
}
