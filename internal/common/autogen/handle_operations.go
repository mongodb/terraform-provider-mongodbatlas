package autogen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/cleanup"
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
	CallParams        func(model any) *config.APICallParams
	StateProperty     string
	PendingStates     []string
	TargetStates      []string
	Timeout           time.Duration
	MinTimeoutSeconds int
	DelaySeconds      int
}
type HandleCreateReq struct {
	CreateAPICallHooks    CreateAPICallHooks
	ReadAPICallHooks      ReadAPICallHooks
	Resp                  *resource.CreateResponse
	Client                *config.MongoDBClient
	Plan                  any
	CallParams            *config.APICallParams
	DeleteReq             func(model any) *HandleDeleteReq
	Wait                  *WaitReq
	DeleteOnCreateTimeout bool
}

func HandleCreate(ctx context.Context, req HandleCreateReq) {
	d := &req.Resp.Diagnostics
	bodyReq, err := Marshal(req.Plan, false)
	if err != nil {
		addError(d, opCreate, errBuildingAPIRequest, err)
		return
	}

	modifiedParams, modifiedBody := req.CreateAPICallHooks.PreCreateAPICall(*req.CallParams, bodyReq)
	callResult := req.CreateAPICallHooks.PostCreateAPICall(callAPI(ctx, req.Client, modifiedParams, modifiedBody))
	if callResult.Err != nil {
		addError(d, opCreate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(callResult.Body, req.Plan); err != nil {
		addError(d, opCreate, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		addError(d, opCreate, errResolvingResponse, err)
		return
	}

	errWait := handleWaitCreateUpdate(ctx, req.Wait, req.Client, req.Plan, req.ReadAPICallHooks)
	if req.DeleteReq != nil {
		// Handle timeout with cleanup if delete_on_create_timeout is enabled.
		errWait = cleanup.HandleCreateTimeout(req.DeleteOnCreateTimeout, errWait, func(ctxCleanup context.Context) error {
			deleteReq := req.DeleteReq(req.Plan)
			return callDelete(ctxCleanup, deleteReq)
		})
	}
	if errWait != nil {
		addError(d, opCreate, errWaitingForChanges, errWait)
		return
	}
	req.Resp.Diagnostics.Append(req.Resp.State.Set(ctx, req.Plan)...)
}

type HandleReadReq struct {
	ReadAPICallHooks ReadAPICallHooks
	State            any
	RespState        *tfsdk.State
	Client           *config.MongoDBClient
	CallParams       *config.APICallParams
	RespDiags        diag.Diagnostics
}

// HandleRead handles the read operation for a resource.
func HandleRead(ctx context.Context, req HandleReadReq) {
	handleReadCore(
		ctx,
		req,
		func() { req.RespState.RemoveResource(ctx) }, // Resource: silently remove from state
	)
}

// HandleDataSourceRead handles the read operation for a data source.
func HandleDataSourceRead(ctx context.Context, req HandleReadReq) {
	handleReadCore(
		ctx,
		req,
		func() { req.RespDiags.AddError("Resource not found", "The requested resource does not exist") }, // Data source: return error
	)
}

// handleReadCore contains the shared read logic for both resources and data sources.
// The onNotFound callback handles the not-found scenario differently:
//   - Resource: silently removes from state (standard Terraform refresh behavior)
//   - Data source: returns an error (resource must exist)
//
// The setState callback sets the response state with the unmarshalled model.
func handleReadCore(
	ctx context.Context,
	req HandleReadReq,
	onNotFound func(),
) {
	modifiedParams := req.ReadAPICallHooks.PreReadAPICall(*req.CallParams)
	callResult := callAPIWithoutBody(ctx, req.Client, modifiedParams)
	callResult = req.ReadAPICallHooks.PostReadAPICall(req, callResult)
	if notFound(callResult.Body, callResult.Resp) {
		onNotFound()
		return
	}
	if callResult.Err != nil {
		addError(&req.RespDiags, opRead, errCallingAPI, callResult.Err)
		return
	}

	// Use the current state as the base model to set the response state
	if err := Unmarshal(callResult.Body, req.State); err != nil {
		addError(&req.RespDiags, opRead, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.State); err != nil {
		addError(&req.RespDiags, opRead, errResolvingResponse, err)
		return
	}
	req.RespDiags.Append(req.RespState.Set(ctx, req.State)...)
}

type HandleUpdateReq struct {
	UpdateAPICallHooks UpdateAPICallHooks
	ReadAPICallHooks   ReadAPICallHooks
	Resp               *resource.UpdateResponse
	Client             *config.MongoDBClient
	Plan               any
	CallParams         *config.APICallParams
	Wait               *WaitReq
}

func HandleUpdate(ctx context.Context, req HandleUpdateReq) {
	d := &req.Resp.Diagnostics
	bodyReq, err := Marshal(req.Plan, true)
	if err != nil {
		addError(d, opUpdate, errBuildingAPIRequest, err)
		return
	}
	modifiedParams, modifiedBody := req.UpdateAPICallHooks.PreUpdateAPICall(*req.CallParams, bodyReq)
	callResult := req.UpdateAPICallHooks.PostUpdateAPICall(callAPI(ctx, req.Client, modifiedParams, modifiedBody))
	if callResult.Err != nil {
		addError(d, opUpdate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(callResult.Body, req.Plan); err != nil {
		addError(d, opUpdate, errUnmarshallingResponse, err)
		return
	}
	if err := ResolveUnknowns(req.Plan); err != nil {
		addError(d, opUpdate, errResolvingResponse, err)
		return
	}
	if err := handleWaitCreateUpdate(ctx, req.Wait, req.Client, req.Plan, req.ReadAPICallHooks); err != nil {
		addError(d, opUpdate, errWaitingForChanges, err)
		return
	}
	req.Resp.Diagnostics.Append(req.Resp.State.Set(ctx, req.Plan)...)
}

type HandleDeleteReq struct {
	DeleteAPICallHooks DeleteAPICallHooks
	ReadAPICallHooks   ReadAPICallHooks
	Diags              *diag.Diagnostics
	Client             *config.MongoDBClient
	State              any
	CallParams         *config.APICallParams
	Wait               *WaitReq
	StaticRequestBody  string
}

func HandleDelete(ctx context.Context, req HandleDeleteReq) {
	if err := callDelete(ctx, &req); err != nil {
		addError(req.Diags, opDelete, errCallingAPI, err)
		return
	}
	if errWait := handleWaitDelete(ctx, req.Wait, req.Client, req.State, req.ReadAPICallHooks); errWait != nil {
		addError(req.Diags, opDelete, errWaitingForChanges, errWait)
	}
}

// handleWaitCreateUpdate waits until a long-running operation is done if needed.
// It also updates the model with the latest JSON response from the API.
func handleWaitCreateUpdate(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any, hooks ReadAPICallHooks) error {
	if wait == nil {
		return nil
	}
	bodyResp, err := waitForChanges(ctx, wait, client, model, hooks)
	if err != nil || isEmptyJSON(bodyResp) {
		return err
	}
	if err := Unmarshal(bodyResp, model); err != nil {
		return err
	}
	return ResolveUnknowns(model)
}

// handleWaitDelete waits until a long-running operation to delete a resource if neeed.
func handleWaitDelete(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any, hooks ReadAPICallHooks) error {
	if wait == nil {
		return nil
	}
	if _, err := waitForChanges(ctx, wait, client, model, hooks); err != nil {
		return err
	}
	return nil
}

func addError(d *diag.Diagnostics, opName, errSummary string, err error) {
	d.AddError(fmt.Sprintf("Error %s in %s", errSummary, opName), err.Error())
}

type APICallResult struct {
	Err  error
	Resp *http.Response
	Body []byte
}

// callAPI makes a request to the API with the given request body and returns the response body.
// It is used for POST, PUT, PATCH and DELETE with static content.
func callAPI(ctx context.Context, client *config.MongoDBClient, callParams config.APICallParams, bodyReq []byte) APICallResult {
	apiResp, err := client.UntypedAPICall(ctx, callParams, bodyReq)
	if err != nil {
		return APICallResult{Body: nil, Resp: apiResp, Err: err}
	}
	bodyResp, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		return APICallResult{Body: nil, Resp: apiResp, Err: err}
	}
	return APICallResult{Body: bodyResp, Resp: apiResp}
}

// callAPIWithoutBody makes a request to the API without a request body and returns the response body.
// It is used for GET or DELETE requests where no request body is required.
func callAPIWithoutBody(ctx context.Context, client *config.MongoDBClient, callParams config.APICallParams) APICallResult {
	return callAPI(ctx, client, callParams, nil)
}

// callDelete makes a DELETE request to the API, supporting both requests with and without a body.
// Returns nil if the resource is not found (already deleted).
func callDelete(ctx context.Context, req *HandleDeleteReq) error {
	var callResult APICallResult
	modifiedParams := req.DeleteAPICallHooks.PreDeleteAPICall(*req.CallParams)
	if req.StaticRequestBody == "" {
		callResult = callAPIWithoutBody(ctx, req.Client, modifiedParams)
	} else {
		callResult = callAPI(ctx, req.Client, modifiedParams, []byte(req.StaticRequestBody))
	}
	callResult = req.DeleteAPICallHooks.PostDeleteAPICall(callResult)

	if notFound(callResult.Body, callResult.Resp) { // Resource is already deleted, don't fail.
		return nil
	}
	return callResult.Err
}

// waitForChanges waits until a long-running operation is done.
// It returns the latest JSON response from the API so it can be used to update the response state.
func waitForChanges(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any, hooks ReadAPICallHooks) ([]byte, error) {
	if len(wait.TargetStates) == 0 {
		return nil, fmt.Errorf("wait must have at least one target state, pending states: %v", wait.PendingStates)
	}
	stateConf := retry.StateChangeConf{
		Target:     wait.TargetStates,
		Pending:    wait.PendingStates,
		Timeout:    wait.Timeout,
		MinTimeout: time.Duration(wait.MinTimeoutSeconds) * time.Second,
		Delay:      time.Duration(wait.DelaySeconds) * time.Second,
		Refresh:    refreshFunc(ctx, wait, client, model, hooks),
	}
	bodyResp, err := stateConf.WaitForStateContext(ctx)
	if err != nil || bodyResp == nil {
		return nil, err
	}
	return bodyResp.([]byte), err
}

// refreshFunc retries until a target state or error happens.
// It uses a special state value of "DELETED" when the API returns 404 or empty object
func refreshFunc(ctx context.Context, wait *WaitReq, client *config.MongoDBClient, model any, hooks ReadAPICallHooks) retry.StateRefreshFunc {
	return func() (result any, state string, err error) {
		callParams := wait.CallParams(model)
		modifiedParams := hooks.PreReadAPICall(*callParams)
		callResult := callAPIWithoutBody(ctx, client, modifiedParams)
		// TODO fix?
		callResult = hooks.PostReadAPICall(HandleReadReq{
			ReadAPICallHooks: hooks,
			Client:           client,
			State:            model,
			CallParams:       callParams,
		}, callResult)
		if notFound(callResult.Body, callResult.Resp) {
			// if "artificial" states continue to grow we can evaluate using a prefix to clearly separate states coming from API and those defined by refreshFunc
			return emptyJSON, retrystrategy.RetryStrategyDeletedState, nil
		}
		if callResult.Err != nil {
			return nil, "", err
		}
		var objJSON map[string]any
		if err := json.Unmarshal(callResult.Body, &objJSON); err != nil {
			return nil, "", err
		}
		stateValAny, found := objJSON[wait.StateProperty]
		if !found {
			return nil, "", fmt.Errorf("wait state attribute not found: %s", wait.StateProperty)
		}
		stateValStr, ok := stateValAny.(string)
		if !ok {
			return nil, "", fmt.Errorf("wait state attribute value is not a string, attribute name: %s, value: %s", wait.StateProperty, stateValAny)
		}
		return callResult.Body, stateValStr, nil
	}
}

// notFound returns if the resource is not found (API response is 404 or response body is empty JSON).
// That is because some resources like search_deployment can return an ok status code with empty json when resource doesn't exist.
func notFound(bodyResp []byte, apiResp *http.Response) bool {
	return validate.StatusNotFound(apiResp) || isEmptyJSON(bodyResp)
}

func isEmptyJSON(raw []byte) bool {
	return len(raw) == 0 || bytes.Equal(raw, emptyJSON)
}

var emptyJSON = []byte("{}")
