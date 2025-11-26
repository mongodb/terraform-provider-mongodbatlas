package orgserviceaccountsecretapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

type response struct {
	Secrets []map[string]any `tfsdk:"secrets"`
}

func (r *rs) PerformRead(ctx context.Context, req *autogen.HandleReadReq) ([]byte, *http.Response, error) {
	model := req.State.(*TFModel)

	bodyResp, apiResp, err := autogen.CallAPIWithoutBody(ctx, req.Client, req.CallParams)
	if err != nil {
		tflog.Error(ctx, "Failed to read service account secret", map[string]any{
			"error": err.Error(),
		})
		return nil, apiResp, err
	}
	if autogen.NotFound(bodyResp, apiResp) {
		return bodyResp, apiResp, nil
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
