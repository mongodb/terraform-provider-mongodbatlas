package autogen

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	errProcessingImportID = "processing import ID"
	opImport              = "Import"
	ExpectedErrorMsg      = "expected format: %s"
)

// HandleImport handles the import operation for Terraform resources.
// It splits the request ID string by "/" delimiter and maps each part to the corresponding attribute specified in idAttributes.
// When a fallback hook is provided, it is used if default parsing fails.
// Example usage:
//   - HandleImport(ctx, []string{"project_id", "name"}, req, resp, r)
//   - HandleImport(ctx, []string{"project_id", "name"}, req, resp, nil)
//   - example import ID would be "5c9d0a239ccf643e6a35ddasdf/myCluster"
func HandleImport(ctx context.Context, idAttrs []string, req resource.ImportStateRequest, resp *resource.ImportStateResponse, hook any) {
	d := &resp.Diagnostics
	idAttrsWithValue, defaultErr := ProcessImportID(req.ID, idAttrs)
	if defaultErr != nil {
		if fallbackHook, ok := hook.(ImportStateFallbackHook); ok {
			legacyValues, fallbackErr := fallbackHook.ParseLegacyImportID(req.ID)
			if fallbackErr == nil {
				idAttrsWithValue = legacyValues
			} else {
				addError(d, opImport, errProcessingImportID, fmt.Errorf("%v; fallback parse failed: %w", defaultErr, fallbackErr))
				return
			}
		} else {
			addError(d, opImport, errProcessingImportID, defaultErr)
			return
		}
	}
	for attrName, value := range idAttrsWithValue {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrName), value)...)
	}
}

// ProcessImportID is exported for testing purposes only and is not intended for direct usage.
func ProcessImportID(importID string, idAttrs []string) (map[string]string, error) {
	parts := strings.Split(importID, "/")
	if len(parts) != len(idAttrs) {
		return nil, fmt.Errorf(ExpectedErrorMsg, strings.Join(idAttrs, "/"))
	}

	result := make(map[string]string)
	for i, part := range parts {
		result[idAttrs[i]] = part
	}

	return result, nil
}
