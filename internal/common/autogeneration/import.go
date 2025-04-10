package autogeneration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const ExpectedErrorMsg = "expected format: %s"

// GenericImportOperation handles the import operation for Terraform resources.
// It splits the request ID string by "/" delimiter and maps each part to the corresponding attribute specified in idAttributes.
// Example usage:
//   - GenericImportOperation(ctx, []string{"project_id", "name"}, req, resp)
//   - example import ID would be "5c9d0a239ccf643e6a35ddasdf/myCluster"
func GenericImportOperation(ctx context.Context, idAttrs []string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idAttrsWithValue, err := ProcessImportID(req.ID, idAttrs)
	if err != nil {
		resp.Diagnostics.AddError("Error processing import ID", err.Error())
		return
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
