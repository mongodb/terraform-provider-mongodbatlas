package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
)

// PopulateProjectsDataSourceResultsForTest exposes populateProjectsDataSourceModel to package project_test.
func PopulateProjectsDataSourceResultsForTest(ctx context.Context, connV2 *admin.APIClient, projectsRes *admin.PaginatedAtlasGroup) ([]*TFProjectDSModel, diag.Diagnostics) {
	var stateModel tfProjectsDSModel
	diags := populateProjectsDataSourceModel(ctx, connV2, &stateModel, projectsRes)
	return stateModel.Results, diags
}
