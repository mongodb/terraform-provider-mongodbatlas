package searchdeployment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFSearchDeploymentDSModel struct {
	ID          types.String `tfsdk:"id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	ProjectID   types.String `tfsdk:"project_id"`
	Specs       types.List   `tfsdk:"specs"`
	StateName   types.String `tfsdk:"state_name"`
}
