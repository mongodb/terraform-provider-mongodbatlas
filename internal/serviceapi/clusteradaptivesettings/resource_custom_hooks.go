package clusteradaptivesettings

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = &rs{}

// PostReadAPICall normalizes the response when Atlas omits cleared overrides.
// After a whole-map reset (null), Atlas omits adaptiveSettingsOverrides from
// the response; the shared decoder preserves absent and null fields.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err == nil {
		req.State.(*TFModel).AdaptiveSettingsOverrides = jsontypes.NewNormalizedNull()
	}
	return result
}
