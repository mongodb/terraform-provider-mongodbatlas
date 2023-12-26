package advancedcluster

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func (r *advancedClusterRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conn := r.Client.Atlas

	projectID, name, err := splitSClusterAdvancedImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to IMPORT cluster. An error occurred when attempting to read resource ID", err.Error())
		return
	}

	u, _, err := conn.AdvancedClusters.Get(ctx, *projectID, *name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to IMPORT cluster. An error occurred when getting cluster details from Atlas.",
			fmt.Sprintf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err))
		return
	}
	id := conversion.EncodeStateID(map[string]string{
		"cluster_id":   u.ID,
		"project_id":   u.GroupID,
		"cluster_name": u.Name,
	})

	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))
	if resp.Diagnostics.HasError() {
		return
	}
}

func splitSClusterAdvancedImportID(id string) (projectID, clusterName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a advanced cluster, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	clusterName = &parts[2]

	return
}
