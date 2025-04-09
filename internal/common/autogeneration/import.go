package autogeneration

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// GenericImportOperation handles the import operation for Terraform resources.
// It splits the request ID string by "/" delimiter and maps each part to the corresponding attribute specified in idAttributes.
// Example usage:
//   - GenericImportOperation(ctx, []string{"project_id", "name"}, req, resp)
//   - example import ID would be "5c9d0a239ccf643e6a35ddasdf/myCluster"
func GenericImportOperation(ctx context.Context, idAttributes []string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	attrValues, err := processImportID(req.ID, idAttributes)
	if err != nil {
		resp.Diagnostics.AddError("error processing import ID", err.Error())
		return
	}
	for attr, value := range attrValues {
		resp.State.SetAttribute(ctx, path.Root(attr), value)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, resp.State)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func processImportID(importID string, idAttributes []string) (map[string]string, error) {
	parts := strings.Split(importID, "/")
	if len(parts) != len(idAttributes) {
		return nil, errors.New("Expected format: " + strings.Join(idAttributes, "/"))
	}

	result := make(map[string]string)
	for i, part := range parts {
		result[idAttributes[i]] = part
	}

	return result, nil
}
