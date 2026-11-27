package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
)

// Test helpers exported only for package project_test.

func PopulateProjectsDataSourceModelForTest(ctx context.Context, connV2 *admin.APIClient, projectsRes *admin.PaginatedAtlasGroup) ([]*TFProjectDSModel, diag.Diagnostics) {
	var stateModel tfProjectsDSModel
	diags := populateProjectsDataSourceModel(ctx, connV2, &stateModel, projectsRes)
	return stateModel.Results, diags
}
