package autogen

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errBuildingAPIRequest    = "building API request"
	errCallingAPI            = "calling API"
	errUnmarshallingResponse = "unmarshalling API response"
	errResolvingResponse     = "resolving API response"
	errWaitingForChanges     = "waiting for changes"
	opCreate                 = "Create"
	opRead                   = "Read"
	opUpdate                 = "Update"
	opDelete                 = "Delete"
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
	d := &req.Resp.Diagnostics
	reqBody, err := Marshal(req.Plan, false)
	if err != nil {
		addError(d, opCreate, errBuildingAPIRequest, err)
		return
	}
	respBody, err := callAPIWithBody(ctx, req.Client, req.CallParams, reqBody)
	if err != nil {
		addError(d, opCreate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(respBody, req.Plan); err != nil {
		addError(d, opCreate, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		addError(d, opCreate, errResolvingResponse, err)
		return
	}
	if err := handleWaitCreateUpdate(ctx, req.Wait, req.Client, req.Plan); err != nil {
		addError(d, opCreate, errWaitingForChanges, err)
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
	d := &req.Resp.Diagnostics
	respBody, apiResp, err := callAPIWithoutBody(ctx, req.Client, req.CallParams)
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			req.Resp.State.RemoveResource(ctx)
			return
		}
		addError(d, opRead, errCallingAPI, err)
		return
	}

	// Use the current state as the base model to set the response state
	if err := Unmarshal(respBody, req.State); err != nil {
		addError(d, opRead, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.State); err != nil {
		addError(d, opRead, errResolvingResponse, err)
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
	d := &req.Resp.Diagnostics
	reqBody, err := Marshal(req.Plan, true)
	if err != nil {
		addError(d, opUpdate, errBuildingAPIRequest, err)
		return
	}
	respBody, err := callAPIWithBody(ctx, req.Client, req.CallParams, reqBody)
	if err != nil {
		addError(d, opUpdate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(respBody, req.Plan); err != nil {
		addError(d, opUpdate, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		addError(d, opUpdate, errResolvingResponse, err)
		return
	}
	if err := handleWaitCreateUpdate(ctx, req.Wait, req.Client, req.Plan); err != nil {
		addError(d, opUpdate, errWaitingForChanges, err)
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
	d := &req.Resp.Diagnostics
	if _, _, err := callAPIWithoutBody(ctx, req.Client, req.CallParams); err != nil {
		addError(d, opDelete, errCallingAPI, err)
		return
	}
	if err := handleWaitDelete(ctx, req.Wait, req.Client); err != nil {
		addError(d, opDelete, errWaitingForChanges, err)
	}
}

// handleWaitCreateUpdate waits until a long-running operation is done if needed.
// It also updates the model with the latest JSON response from the API.
func handleWaitCreateUpdate(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any) error {
	if wait == nil {
		return nil
	}
	respBody, err := waitForChanges(ctx, wait, client)
	if err != nil {
		return err
	}
	if len(respBody) == 0 {
		return nil
	}
	if err := Unmarshal(respBody, model); err != nil {
		return err
	}
	return ResolveUnknowns(model)
}

// handleWaitDelete waits until a long-running operation to delete a resource if neeed.
func handleWaitDelete(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) error {
	if wait == nil {
		return nil
	}
	if _, err := waitForChanges(ctx, wait, client); err != nil {
		return err
	}
	return nil
}

func addError(d *diag.Diagnostics, opName, errSummary string, err error) {
	d.AddError(fmt.Sprintf("Error %s in %s", errSummary, opName), err.Error())
}

// callAPIWithBody makes a request to the API with the given request body and returns the response body.
// It is used for POST, PUT, and PATCH requests where a request body is required.
func callAPIWithBody(ctx context.Context, client *config.MongoDBClient, callParams *config.APICallParams, reqBody []byte) ([]byte, error) {
	apiResp, err := client.UntypedAPICall(ctx, callParams, reqBody)
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
	apiResp, err := client.UntypedAPICall(ctx, callParams, nil)
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

// waitForChanges waits until a long-running operation is done.
// It returns the latest JSON response from the API so it can be used to update the response state.
func waitForChanges(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) ([]byte, error) {
	// time.Sleep(time.Duration(wait.TimeoutSeconds) * time.Second) // TODO: TimeoutSeconds is temporarily used to allow time to destroy the resource until autogen long-running operations are supported in CLOUDP-314960
	if len(wait.TargetStates) == 0 {
		return nil, fmt.Errorf("no target states, this is an error in the provider, pending states: %v", wait.PendingStates)
	}
	stateConf := retry.StateChangeConf{
		Target:     wait.TargetStates,
		Pending:    wait.PendingStates,
		Timeout:    time.Duration(wait.TimeoutSeconds) * time.Second,
		MinTimeout: time.Duration(wait.MinTimeoutSeconds) * time.Second,
		Delay:      time.Duration(wait.DelaySeconds) * time.Second,
		Refresh:    refreshFunc(ctx, wait, client),
	}
	respBody, err := stateConf.WaitForStateContext(ctx)
	if err != nil || respBody == nil {
		return nil, err
	}
	return respBody.([]byte), err
}

// refreshFunc retries until a target state or error happens.
// It uses a special state value of "DELETED" when the when API returns 404 or empty object
func refreshFunc(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) retry.StateRefreshFunc {
	return func() (result any, state string, err error) {
		respBody, httpResp, err := callAPIWithoutBody(ctx, client, wait.CallParams)
		if validate.StatusNotFound(httpResp) {
			return nil, retrystrategy.RetryStrategyDeletedState, nil
		}
		if err != nil {
			return nil, "", err
		}
		var objJSON map[string]any
		// TODO: check search deployment when deleted with new API. StatusServiceUnavailable not checked.
		if err := json.Unmarshal(respBody, &objJSON); err != nil {
			return nil, "", err
		}
		stateValAny, found := objJSON[wait.StateAttribute]
		if !found { // if state attribute is not found we assume that the result is deleted, it could also happen that the state attribute name is incorrectly specified.
			// TODO: also check id not ID or any attribute at all comes, so we can keep !found error
			return nil, retrystrategy.RetryStrategyDeletedState, nil
		}
		stateValStr, ok := stateValAny.(string)
		if !ok {
			return nil, "", fmt.Errorf("wait state value is not a string, state name: %s, state value: %s", wait.StateAttribute, stateValAny)
		}
		return respBody, stateValStr, nil
	}
}
