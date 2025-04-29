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
	bodyReq, err := Marshal(req.Plan, false)
	if err != nil {
		addError(d, opCreate, errBuildingAPIRequest, err)
		return
	}
	bodyResp, err := callAPIWithBody(ctx, req.Client, req.CallParams, bodyReq)
	if err != nil {
		addError(d, opCreate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(bodyResp, req.Plan); err != nil {
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
	bodyResp, apiResp, err := callAPIWithoutBody(ctx, req.Client, req.CallParams)
	if notFound(bodyResp, apiResp) {
		req.Resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		addError(d, opRead, errCallingAPI, err)
		return
	}

	// Use the current state as the base model to set the response state
	if err := Unmarshal(bodyResp, req.State); err != nil {
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
	bodyReq, err := Marshal(req.Plan, true)
	if err != nil {
		addError(d, opUpdate, errBuildingAPIRequest, err)
		return
	}
	bodyResp, err := callAPIWithBody(ctx, req.Client, req.CallParams, bodyReq)
	if err != nil {
		addError(d, opUpdate, errCallingAPI, err)
		return
	}

	// Use the plan as the base model to set the response state
	if err := Unmarshal(bodyResp, req.Plan); err != nil {
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
	bodyResp, err := waitForChanges(ctx, wait, client)
	if err != nil || isEmptyJSON(bodyResp) {
		return err
	}
	if err := Unmarshal(bodyResp, model); err != nil {
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
func callAPIWithBody(ctx context.Context, client *config.MongoDBClient, callParams *config.APICallParams, bodyReq []byte) ([]byte, error) {
	apiResp, err := client.UntypedAPICall(ctx, callParams, bodyReq)
	if err != nil {
		return nil, err
	}
	bodyResp, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		return nil, err
	}
	return bodyResp, nil
}

// callAPIWithoutBody makes a request to the API without a request body and returns the response body.
// It is used for GET or DELETE requests where no request body is required.
func callAPIWithoutBody(ctx context.Context, client *config.MongoDBClient, callParams *config.APICallParams) ([]byte, *http.Response, error) {
	apiResp, err := client.UntypedAPICall(ctx, callParams, nil)
	if err != nil {
		return nil, apiResp, err
	}
	bodyResp, err := io.ReadAll(apiResp.Body)
	apiResp.Body.Close()
	if err != nil {
		return nil, apiResp, err
	}
	return bodyResp, apiResp, nil
}

// waitForChanges waits until a long-running operation is done.
// It returns the latest JSON response from the API so it can be used to update the response state.
func waitForChanges(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) ([]byte, error) {
	if len(wait.TargetStates) == 0 {
		return nil, nil // nothing to do if no target states
	}
	stateConf := retry.StateChangeConf{
		Target:     wait.TargetStates,
		Pending:    wait.PendingStates,
		Timeout:    time.Duration(wait.TimeoutSeconds) * time.Second,
		MinTimeout: time.Duration(wait.MinTimeoutSeconds) * time.Second,
		Delay:      time.Duration(wait.DelaySeconds) * time.Second,
		Refresh:    refreshFunc(ctx, wait, client),
	}
	bodyResp, err := stateConf.WaitForStateContext(ctx)
	if err != nil || bodyResp == nil {
		return nil, err
	}
	return bodyResp.([]byte), err
}

// refreshFunc retries until a target state or error happens.
// It uses a special state value of "DELETED" when the when API returns 404 or empty object
func refreshFunc(ctx context.Context, wait *WaitReq, client *config.MongoDBClient) retry.StateRefreshFunc {
	return func() (result any, state string, err error) {
		bodyResp, httpResp, err := callAPIWithoutBody(ctx, client, wait.CallParams)
		if notFound(bodyResp, httpResp) {
			return emptyJSON, retrystrategy.RetryStrategyDeletedState, nil
		}
		if err != nil {
			return nil, "", err
		}
		var objJSON map[string]any
		// TODO: check search deployment when deleted with new API. StatusServiceUnavailable not checked.
		if err := json.Unmarshal(bodyResp, &objJSON); err != nil {
			return nil, "", err
		}
		stateValAny, found := objJSON[wait.StateAttribute]
		if !found {
			return nil, "", fmt.Errorf("wait state attribute not found: %s", wait.StateAttribute)
		}
		stateValStr, ok := stateValAny.(string)
		if !ok {
			return nil, "", fmt.Errorf("wait state attribute value is not a string, attribute name: %s, state value: %s", wait.StateAttribute, stateValAny)
		}
		return bodyResp, stateValStr, nil
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
