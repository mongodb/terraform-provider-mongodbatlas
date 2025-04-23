package autogen

import (
	"context"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type WaitForChangesReq struct {
	CallParams        *config.APICallParams
	StateAttribute    string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}

// WaitForChanges waits until a long-running operation is done.
// TODO: This is an basic implementation, it will be replaced in CLOUDP-314960.
func WaitForChanges(ctx context.Context, req *WaitForChangesReq) error {
	time.Sleep(time.Duration(req.TimeoutSeconds) * time.Second) // TODO: TimeoutSeconds is temporarily used to allow time to destroy the resource until autogen long-running operations are supported in CLOUDP-314960
	return nil
}
