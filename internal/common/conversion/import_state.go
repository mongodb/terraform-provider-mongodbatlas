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

// ImportSplitN splits an import string by "/" into at most expectedParts parts,
// so the last part may itself contain slashes (for example an IP access list CIDR entry).
func ImportSplitN(importRaw string, expectedParts int) (ok bool, parts []string) {
	parts = strings.SplitN(importRaw, "/", expectedParts)
	if len(parts) != expectedParts {
		return false, nil
	}
	return true, parts
}

const LegacyImportFormatWarningSummary = "Deprecated import ID format"

// LegacyImportFormatWarningDetail returns the warning detail used when a legacy import format is accepted.
func LegacyImportFormatWarningDetail(legacyFormat, newFormat string) string {
	return fmt.Sprintf("The import ID format %q is deprecated and will be removed in the next major version. Use %q instead.", legacyFormat, newFormat)
}

func ImportStateProjectIDClusterName(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, attrNameProjectID, attrNameClusterName string) {
	var (
		projectID, clusterName string
		legacy                 bool
	)
	if ok, parts := ImportSplit(req.ID, 2); ok {
		projectID, clusterName = parts[0], parts[1]
	} else {
		legacy = true
		parts := strings.SplitN(req.ID, "-", 2)
		if len(parts) != 2 {
			resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("expected 2 parts with %s and %s, use the format {project_id}/{cluster_name} or the legacy format {project_id}-{cluster_name}: %s", attrNameProjectID, attrNameClusterName, req.ID))
			return
		}
		projectID, clusterName = parts[0], parts[1]
	}
	if legacy {
		resp.Diagnostics.AddWarning(LegacyImportFormatWarningSummary, LegacyImportFormatWarningDetail("{project_id}-{cluster_name}", "{project_id}/{cluster_name}"))
	}
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
