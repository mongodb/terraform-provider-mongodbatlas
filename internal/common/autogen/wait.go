package autogen

import (
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type WaitReq struct {
	CallParams        *config.APICallParams
	StateAttribute    string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}

// waitForChanges waits until a long-running operation is done.
// It returns the latest JSON response from the API so it can be used to update the response state.
// TODO: This is a basic implementation, it will be replaced in CLOUDP-314960.
func waitForChanges(req *WaitReq) ([]byte, error) {
	time.Sleep(time.Duration(req.TimeoutSeconds) * time.Second) // TODO: TimeoutSeconds is temporarily used to allow time to destroy the resource until autogen long-running operations are supported in CLOUDP-314960
	return nil, nil
}
