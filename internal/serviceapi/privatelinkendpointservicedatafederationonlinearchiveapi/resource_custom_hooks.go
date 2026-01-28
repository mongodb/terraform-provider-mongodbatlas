package privatelinkendpointservicedatafederationonlinearchiveapi

import (
	"encoding/json"
	"errors"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostCreateAPICallHook = (*rs)(nil)

type createAPIResponse struct {
	Results []map[string]any `json:"results"`
}

// PostCreateAPICall selects the created/updated endpoint from the paginated response
// by matching the endpointId and returns a new APICallResult with that element.
//
// The create API returns PaginatedPrivateNetworkEndpointIdEntryView, but the autogen
// CRUD expects the response body to match the resource schema (PrivateNetworkEndpointIdEntry).
func (r *rs) PostCreateAPICall(req autogen.HandleCreateReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	endpointID := req.Plan.(*TFModel).EndpointId.ValueString()

	var responseJSON createAPIResponse
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	for _, entry := range responseJSON.Results {
		if entry["endpointId"] == endpointID {
			marshaledEntry, err := json.Marshal(entry)
			return autogen.APICallResult{Body: marshaledEntry, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: errors.New("endpointId not found in create private endpoint ids response")}
}
