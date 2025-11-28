package autogen

import (
	"context"
	"fmt"
	"net/http"
)

var _ ResourceAPIOperations = &DefaultResourceAPIOperations{}

// ResourceAPIOperations defines the interface for API operations in autogen resources.
type ResourceAPIOperations interface {
	PerformRead(ctx context.Context, req *HandleReadReq) ([]byte, *http.Response, error)
	PerformCreate(ctx context.Context, req *HandleCreateReq) ([]byte, *http.Response, error)
	PerformUpdate(ctx context.Context, req *HandleUpdateReq) ([]byte, *http.Response, error)
	PerformDelete(ctx context.Context, req *HandleDeleteReq) error
}

// DefaultResourceAPIOperations is used as an embedded struct for all autogen resources and provides default implementations for all API operations.
type DefaultResourceAPIOperations struct{}

func (d *DefaultResourceAPIOperations) PerformRead(ctx context.Context, req *HandleReadReq) ([]byte, *http.Response, error) {
	return CallAPIWithoutBody(ctx, req.Client, req.CallParams)
}

func (d *DefaultResourceAPIOperations) PerformCreate(ctx context.Context, req *HandleCreateReq) ([]byte, *http.Response, error) {
	bodyReq, err := Marshal(req.Plan, false)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errBuildingAPIRequest, err)
	}
	return CallAPIWithBody(ctx, req.Client, req.CallParams, bodyReq)
}

func (d *DefaultResourceAPIOperations) PerformUpdate(ctx context.Context, req *HandleUpdateReq) ([]byte, *http.Response, error) {
	bodyReq, err := Marshal(req.Plan, true)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errBuildingAPIRequest, err)
	}
	return CallAPIWithBody(ctx, req.Client, req.CallParams, bodyReq)
}

func (d *DefaultResourceAPIOperations) PerformDelete(ctx context.Context, req *HandleDeleteReq) error {
	var err error
	var bodyResp []byte
	var apiResp *http.Response
	if req.StaticRequestBody == "" {
		bodyResp, apiResp, err = CallAPIWithoutBody(ctx, req.Client, req.CallParams)
	} else {
		bodyResp, apiResp, err = CallAPIWithBody(ctx, req.Client, req.CallParams, []byte(req.StaticRequestBody))
	}
	if NotFound(bodyResp, apiResp) {
		return nil // Already deleted
	}
	return err
}
