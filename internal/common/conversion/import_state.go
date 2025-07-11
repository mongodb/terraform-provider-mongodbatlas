package conversion

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportSplit is a generalized function to split import strings by "/" with validation
func ImportSplit(importRaw string, expectedParts int) (ok bool, parts []string) {
	parts = strings.Split(importRaw, "/")
	if len(parts) != expectedParts {
		return false, nil
	}
	return true, parts
}

func ImportStateProjectIDClusterName(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, attrNameProjectID, attrNameClusterName string) {
	parts := strings.SplitN(req.ID, "-", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("expected 2 parts with %s and %s: %s", attrNameProjectID, attrNameClusterName, req.ID))
		return
	}
	projectID, clusterName := parts[0], parts[1]
	if err := ValidateProjectID(projectID); err != nil {
		resp.Diagnostics.AddError("invalid project_id in import ID", err.Error())
	}
	if err := ValidateClusterName(clusterName); err != nil {
		resp.Diagnostics.AddError("invalid cluster_name in import ID", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrNameProjectID), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attrNameClusterName), clusterName)...)
}

func ValidateProjectID(projectID string) error {
	re := regexp.MustCompile("^([a-f0-9]{24})$")
	if !re.MatchString(projectID) {
		return fmt.Errorf("project_id must be a 24 character hex string: %s", projectID)
	}
	return nil
}

func ValidateClusterName(clusterName string) error {
	re := regexp.MustCompile("^([a-zA-Z0-9][a-zA-Z0-9-]*)?[a-zA-Z0-9]+$")
	if !re.MatchString(clusterName) || len(clusterName) < 1 || len(clusterName) > 64 {
		return fmt.Errorf("cluster_name must be a string with length between 1 and 64, starting and ending with an alphanumeric character, and containing only alphanumeric characters and hyphens: %s", clusterName)
	}
	return nil
}
