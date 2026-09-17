package clusteradaptivesettings

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = &rs{}

// PostReadAPICall clears the override state when Atlas omits the field from
// the response. Atlas omits adaptiveSettingsOverrides after a whole-map reset
// (null), e.g. when someone clears overrides via the Atlas UI. The shared
// decoder skips absent fields, so without this hook the old value stays in
// state and Terraform misses the drift.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err == nil {
		req.State.(*TFModel).AdaptiveSettingsOverrides = jsontypes.NewNormalizedNull()
	}
	return result
}
