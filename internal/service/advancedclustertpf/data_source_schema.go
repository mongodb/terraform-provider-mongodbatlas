package advancedclustertpf

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TODO: see if resource model can be used instead
type ModelDS struct {
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
}
