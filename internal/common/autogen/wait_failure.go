package autogen

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WaitFailure is the last poll snapshot when a wait ends badly.
// TimeoutErr is set only on SDK timeout; it is nil on unexpected state.
type WaitFailure struct {
	LastJSON   map[string]any
	Model      any
	TimeoutErr error
	LastState  string
}

const tfsdkTag = "tfsdk"

// ResourceID joins id attribute values from the in-memory model in import order.
// Non-string fields are skipped. Empty or unknown strings still occupy their slot.
func ResourceID(model any, idAttrs []string) string {
	if model == nil || len(idAttrs) == 0 {
		return ""
	}
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ""
	}

	byTag := make(map[string]reflect.Value, val.NumField())
	typ := val.Type()
	for i := range val.NumField() {
		tag := typ.Field(i).Tag.Get(tfsdkTag)
		if tag == "" {
			continue
		}
		byTag[tag] = val.Field(i)
	}

	parts := make([]string, 0, len(idAttrs))
	for _, attr := range idAttrs {
		field, ok := byTag[attr]
		if !ok {
			continue
		}
		str, ok := field.Interface().(types.String)
		if !ok {
			continue
		}
		parts = append(parts, str.ValueString())
	}
	return strings.Join(parts, "/")
}

func waitFailureIDPrefix(model any, idAttrs []string) string {
	id := ResourceID(model, idAttrs)
	if id == "" {
		return ""
	}
	return "id " + id + ": "
}

// DefaultFormatWaitFailure builds the wait error with an import-style id prefix.
// When TimeoutErr is set it wraps that error with %w so *retry.TimeoutError stays in the chain.
func DefaultFormatWaitFailure(wait *WaitReq, req WaitFailure) error {
	var idAttrs []string
	if wait != nil {
		idAttrs = wait.IDAttributes
	}
	prefix := waitFailureIDPrefix(req.Model, idAttrs)

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
	if wait != nil && wait.ErrorDescriptionProperty != "" && req.LastJSON != nil {
		if msg, ok := req.LastJSON[wait.ErrorDescriptionProperty].(string); ok && msg != "" {
			return msg
		}
	}
	return fmt.Sprintf("operation failed with state %q", req.LastState)
}

// IsWaitContinueState reports whether the polled state should keep waiting.
func IsWaitContinueState(pending, target []string, state string) bool {
	return slices.Contains(pending, state) || slices.Contains(target, state)
}
