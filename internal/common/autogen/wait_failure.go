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

type idAttrPair struct {
	name  string
	value string
}

func resourceIDPairs(model any, idAttrs []string) []idAttrPair {
	if model == nil || len(idAttrs) == 0 {
		return nil
	}
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
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

	pairs := make([]idAttrPair, 0, len(idAttrs))
	for _, attr := range idAttrs {
		field, ok := byTag[attr]
		if !ok {
			continue
		}
		str, ok := field.Interface().(types.String)
		if !ok {
			continue
		}
		pairs = append(pairs, idAttrPair{name: attr, value: str.ValueString()})
	}
	return pairs
}

func waitFailureIDPrefix(model any, idAttrs []string) string {
	pairs := resourceIDPairs(model, idAttrs)
	if len(pairs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(pairs))
	for _, p := range pairs {
		parts = append(parts, fmt.Sprintf("%s=%q", p.name, p.value))
	}
	return strings.Join(parts, ", ") + ": "
}

// DefaultFormatWaitFailure builds the wait error with named id attributes.
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
