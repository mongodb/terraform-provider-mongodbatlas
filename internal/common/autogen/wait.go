package autogen

import "time"

type WaitForChangesReq struct {
	StateAttribute    string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}

// WaitForChanges waits until a long-running operation is done.
// TODO: this is an basic implementation, it will be replaced in CLOUDP-314960
func WaitForChanges(req *WaitForChangesReq) error {
	time.Sleep(time.Duration(req.TimeoutSeconds) * time.Second) // TODO: TimeoutSeconds is temporarily used to allow time to destroy the resource until autogen long-running operations are supported in CLOUDP-314960
	return nil
}
