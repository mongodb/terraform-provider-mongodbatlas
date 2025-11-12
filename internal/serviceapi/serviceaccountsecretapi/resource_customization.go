package serviceaccountsecretapi

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

// reqHandle.State is SA Secret state
// set it to a SA state
func (r *rs) PreReadHandler(reqHandle autogen.HandleReadReq) autogen.HandleReadReq {
	reqHandle.CallParams.RelativePath += "/secrets"
	reqHandle.State = nil
	return reqHandle
}

// responseState is a SA state
// convert to SA Secret state
func (r *rs) PostReadHandler(responseState any, state *TFModel) {

	fmt.Println("dsafdasf")
}
