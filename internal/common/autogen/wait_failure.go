package autogen

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// WaitFailure is the last poll snapshot when a wait ends badly.
// TimeoutErr is set only on SDK timeout; it is nil on unexpected state.
type WaitFailure struct {
	LastJSON   map[string]any
	Model      any
	TimeoutErr error
	LastState  string
}

func waitFailureIDPrefix(wait *WaitReq, model any) string {
	if wait == nil || wait.FormatID == nil {
		return ""
	}
	id := wait.FormatID(model)
	if id == "" {
		return ""
	}
	return id + ": "
}

// DefaultFormatWaitFailure builds the wait error with named id attributes.
// When TimeoutErr is set it wraps that error with %w so *retry.TimeoutError stays in the chain.
func DefaultFormatWaitFailure(wait *WaitReq, req WaitFailure) error {
	prefix := waitFailureIDPrefix(wait, req.Model)

	if req.TimeoutErr != nil {
		if prefix == "" {
			return req.TimeoutErr
		}
		return fmt.Errorf("%s%w", prefix, req.TimeoutErr)
	}

	msg := unexpectedStateMessage(wait, req)
	return errors.New(prefix + msg)
}

func unexpectedStateMessage(wait *WaitReq, req WaitFailure) string {
	msg := fmt.Sprintf("Operation failed with state %q", req.LastState)
	if wait != nil && len(wait.TargetStates) > 0 {
		msg = fmt.Sprintf("%s, wanted target %q", msg, strings.Join(wait.TargetStates, ", "))
	}
	if wait != nil && wait.ErrorDescriptionProperty != "" && req.LastJSON != nil {
		if desc, ok := req.LastJSON[wait.ErrorDescriptionProperty].(string); ok && desc != "" {
			return msg + ". " + desc
		}
	}
	return msg
}

// IsWaitContinueState reports whether the polled state should keep waiting.
func IsWaitContinueState(pending, target []string, state string) bool {
	return slices.Contains(pending, state) || slices.Contains(target, state)
}
