package customplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}
var _ resource.ResourceWithModifyPlan = &rs{}

type BaseResourcePlanModify interface {
	resource.Resource
	resource.ResourceWithConfigure
	resource.ResourceWithImportState
	resource.ResourceWithModifyPlan
}

func WrappedResource(base BaseResourcePlanModify) func() resource.Resource {
	return func() resource.Resource {
		return &rs{base: base}
	}
}

type rs struct {
	base BaseResourcePlanModify
}

func (r *rs) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.base.Metadata(ctx, req, resp)
}

func (r *rs) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.base.Configure(ctx, req, resp)
}

func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.base.ModifyPlan(ctx, req, resp)
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.base.Schema(ctx, req, resp)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.base.Create(ctx, req, resp)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.base.Read(ctx, req, resp)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.base.Update(ctx, req, resp)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.base.Delete(ctx, req, resp)
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.base.ImportState(ctx, req, resp)
}

func TestWrappingAdvancedClusterTPF(t *testing.T) {
	var (
		resources = []func() resource.Resource{
			WrappedResource(advancedclustertpf.Resource().(BaseResourcePlanModify)),
		}
		mockConfig   = unit.MockConfigAdvancedClusterTPF.WithResources(resources)
		baseConfig   = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterReplicasetOneRegion)
		resourceName = baseConfig.ResourceName
		testCases    = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_mongo_db_major_version_changed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("mongo_db_version")),
				},
			},
			{
				ConfigFilename: "main_backup_enabled.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("mongo_db_version"), knownvalue.StringExact("8.0.5")),
				},
			},
		}
	)
	baseConfig.TestdataPrefix = unit.PackagePath("advancedclustertpf")
	unit.RunPlanCheckTests(t, baseConfig, testCases)
}
