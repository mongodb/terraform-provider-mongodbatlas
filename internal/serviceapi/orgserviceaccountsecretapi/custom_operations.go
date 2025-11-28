package orgserviceaccountsecretapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

type response struct {
	Secrets []map[string]any `tfsdk:"secrets"`
}

func (r *rs) PerformRead(ctx context.Context, req *autogen.HandleReadReq) ([]byte, *http.Response, error) {
	model := req.State.(*TFModel)

	bodyResp, apiResp, err := autogen.CallAPIWithoutBody(ctx, req.Client, req.CallParams)
	if err != nil {
		return nil, apiResp, fmt.Errorf("failed to read service account secret: %w", err)
	}

	var responseJSON response
	if err := json.Unmarshal(bodyResp, &responseJSON); err != nil {
		return nil, nil, err
	}

	id := model.Id.ValueString()
	for _, secret := range responseJSON.Secrets {
		if secret["id"] == id {
			marshaledSecret, err := json.Marshal(secret)
			if err != nil {
				return nil, nil, err
			}
			return marshaledSecret, apiResp, nil
		}
	}

	return nil, apiResp, nil
}
