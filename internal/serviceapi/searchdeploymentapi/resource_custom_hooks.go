package searchdeploymentapi

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = (*rs)(nil)

// PostReadAPICall maps the search deployment API quirk of returning an ok status code with an empty
// JSON body for a missing resource to autogen.ErrNotFound. This also applies while polling the
// delete wait, where the empty body indicates the deployment is gone.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err == nil && autogen.IsEmptyJSON(result.Body) {
		result.Err = autogen.ErrNotFound
	}
	return result
}
