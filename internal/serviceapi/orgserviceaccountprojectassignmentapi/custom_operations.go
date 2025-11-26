package orgserviceaccountprojectassignmentapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

type response struct {
	Results []projectServiceAccount `json:"results"`
}

type projectServiceAccount struct {
	ClientId string   `json:"clientId"`
	Roles    []string `json:"roles"`
}

func (r *rs) PerformRead(ctx context.Context, req *autogen.HandleReadReq) ([]byte, *http.Response, error) {
	model := req.State.(*TFModel)

	bodyResp, apiResp, err := autogen.CallAPIWithoutBody(ctx, req.Client, req.CallParams)
	if err != nil {
		return nil, apiResp, fmt.Errorf("failed to read service accounts of project: %w", err)
	}

	var responseJSON response
	if err := json.Unmarshal(bodyResp, &responseJSON); err != nil {
		return nil, nil, err
	}

	for _, serviceAccount := range responseJSON.Results {
		if serviceAccount.ClientId == model.ClientId.ValueString() {
			marshaledServiceAccount, err := json.Marshal(serviceAccount)
			if err != nil {
				return nil, nil, err
			}
			return marshaledServiceAccount, apiResp, nil
		}
	}

	return nil, apiResp, nil
}
