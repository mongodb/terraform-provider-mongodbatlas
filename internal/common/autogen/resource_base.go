package autogen

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ ResourceAPIOperations = &RSAutogen{}

// RSAutogen is used as an embedded struct for all autogen resources.
type RSAutogen struct {
	config.RSCommon
}

// ResourceAPIOperations defines the interface for API operations in autogen resources.
// RSAutogen provides default implementations which can be overridden by individual resources.
type ResourceAPIOperations interface {
	PerformRead(ctx context.Context, req *HandleReadReq) ([]byte, *http.Response, error)
	PerformCreate(ctx context.Context, req *HandleCreateReq) ([]byte, *http.Response, error)
	PerformUpdate(ctx context.Context, req *HandleUpdateReq) ([]byte, *http.Response, error)
	PerformDelete(ctx context.Context, req *HandleDeleteReq) error
}

func (r *RSAutogen) PerformRead(ctx context.Context, req *HandleReadReq) ([]byte, *http.Response, error) {
	return CallAPIWithoutBody(ctx, r.Client, req.CallParams)
}

func (r *RSAutogen) PerformCreate(ctx context.Context, req *HandleCreateReq) ([]byte, *http.Response, error) {
	bodyReq, err := Marshal(req.Plan, false)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errBuildingAPIRequest, err)
	}
	return CallAPIWithBody(ctx, r.Client, req.CallParams, bodyReq)
}

func (r *RSAutogen) PerformUpdate(ctx context.Context, req *HandleUpdateReq) ([]byte, *http.Response, error) {
	bodyReq, err := Marshal(req.Plan, true)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errBuildingAPIRequest, err)
	}
	return CallAPIWithBody(ctx, r.Client, req.CallParams, bodyReq)
}

func (r *RSAutogen) PerformDelete(ctx context.Context, req *HandleDeleteReq) error {
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
