package teamprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func NewTFModel(ctx context.Context, apiResp *admin.TeamRole, projectID string) (*TFModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if apiResp == nil {
		return nil, diags
	}

	roleNames := conversion.TFSetValueOrNull(ctx, apiResp.RoleNames, types.StringType)

	return &TFModel{
		ProjectId: 	types.StringValue(projectID),
		TeamId:    	types.StringPointerValue(apiResp.TeamId),
		RoleNames:  roleNames,

	}, diags
}



func NewAtlasReq(ctx context.Context, plan *TFModel) (*[]admin.TeamRole, diag.Diagnostics) {
    var roleNames []string
	if !plan.RoleNames.IsNull() && !plan.RoleNames.IsUnknown() {
		roleNames = conversion.TypesSetToString(ctx, plan.RoleNames)
	}

	teamRole := admin.TeamRole{
		TeamId: plan.TeamId.ValueStringPointer(),
		RoleNames: &roleNames,

	}
	return &[]admin.TeamRole{teamRole}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.TeamRole, diag.Diagnostics) {
	var roleNames []string
	if !plan.RoleNames.IsNull() && !plan.RoleNames.IsUnknown() {
		roleNames = conversion.TypesSetToString(ctx, plan.RoleNames)
	}

	return &admin.TeamRole{
		RoleNames: &roleNames,
		TeamId: plan.TeamId.ValueStringPointer(),
	}, nil
}



