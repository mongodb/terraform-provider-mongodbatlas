package apiresource_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	updateResourceName = "mongodbatlas_api_update.failover"
	streamResourceName = "mongodbatlas_stream_instance.demo"
)

// TestAccAPIUpdate_streamWorkspace_failoverRegions exercises the
// mongodbatlas_api_update resource against the failoverRegions field on the
// stream workspace PATCH endpoint.
//
// failoverRegions is a perfect demo target for `api_update`:
//   - It is accepted by Atlas in production today (verified by reading MMS:
//     ApiStreamsTenantUpdateRequestView + ApiStreamsResource.updateTenant).
//   - It is `@Schema(hidden = true)` in the MMS DTO, so it does NOT appear in
//     the public OpenAPI spec. The typed stream_instance / stream_workspace
//     resources cannot expose it.
//   - It demonstrates exactly the hybrid-coexistence story: typed resource
//     owns the workspace lifecycle, api_update reaches a hidden field the
//     typed resource cannot see.
//
// Body shape (matches ApiStreamsTenantUpdateRequestView.failoverRegions):
//
//	{"failoverRegions": [{"cloudProvider": "AWS", "region": "DUBLIN_IRL"}]}
//
// The region value must be an ApiStreamsAWSRegionView enum name (DUBLIN_IRL,
// VIRGINIA_USA, etc.), not an AWS region code. cloud-dev's allowlist
// (StreamsConstants._allowedAWSRegions) only permits VIRGINIA_USA and
// DUBLIN_IRL on AppEnv.DEV — and the failover region must differ from the
// workspace primary, so DUBLIN_IRL is the only valid choice here.
//
// Note: failoverRegions is mutually exclusive with dataProcessRegion AND with
// processorStatus in the same PATCH call. We only send failoverRegions here.
func TestAccAPIUpdate_streamWorkspace_failoverRegions(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamInstance(projectID),
		Steps: []resource.TestStep{
			{
				// Step 1: create workspace + initial PATCH with one failover region.
				// Region value uses the ApiStreamsAWSRegionView enum name (e.g.
				// DUBLIN_IRL = EU_WEST_1), not the AWS region code. On cloud-dev,
				// only VIRGINIA_USA (workspace primary) and DUBLIN_IRL (allowed
				// AWS failover target) pass the StreamsConstants region allowlist
				// check — see docs-context/api-update-demo-target-investigation.md
				// for the trail.
				Config: configStreamWorkspaceWithFailover(projectID, instanceName, "DUBLIN_IRL"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamInstanceExists(streamResourceName, projectID),
					resource.TestCheckResourceAttr(updateResourceName, "preview", "true"),
					resource.TestCheckResourceAttr(updateResourceName, "body.failoverRegions.0.cloudProvider", "AWS"),
					resource.TestCheckResourceAttr(updateResourceName, "body.failoverRegions.0.region", "DUBLIN_IRL"),
					resource.TestCheckResourceAttrSet(updateResourceName, "id"),
				),
			},
			{
				// Step 2: drop the api_update block — typed resource stays,
				// failoverRegions persists on the workspace (no-op delete).
				Config: configStreamWorkspaceOnly(projectID, instanceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamInstanceExists(streamResourceName, projectID),
				),
			},
		},
	})
}

func configStreamWorkspaceWithFailover(projectID, name, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "demo" {
			project_id    = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region         = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
		}

		resource "mongodbatlas_api_update" "failover" {
			path    = "/api/atlas/v2/groups/${mongodbatlas_stream_instance.demo.project_id}/streams/${mongodbatlas_stream_instance.demo.instance_name}"
			preview = true

			body = {
				failoverRegions = [
					{
						cloudProvider = "AWS"
						region        = %[3]q
					},
				]
			}
		}
	`, projectID, name, region)
}

func configStreamWorkspaceOnly(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "demo" {
			project_id    = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region         = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
		}
	`, projectID, name)
}

func checkStreamInstanceExists(rsName, projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("not found: %s", rsName)
		}
		name := rs.Primary.Attributes["instance_name"]
		if name == "" {
			return fmt.Errorf("checkExists: instance_name not set for %s", rsName)
		}
		if !streamWorkspaceExists(projectID, name) {
			return fmt.Errorf("stream workspace (%s/%s) does not exist", projectID, name)
		}
		return nil
	}
}

func checkDestroyStreamInstance(projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_instance" {
				continue
			}
			name := rs.Primary.Attributes["instance_name"]
			if name == "" {
				continue
			}
			if streamWorkspaceExists(projectID, name) {
				return fmt.Errorf("stream workspace (%s/%s) still exists", projectID, name)
			}
		}
		return nil
	}
}

func streamWorkspaceExists(projectID, name string) bool {
	params := config.APICallParams{
		Method:        http.MethodGet,
		VersionHeader: "application/vnd.atlas.2023-02-01+json",
		RelativePath:  "/api/atlas/v2/groups/{projectId}/streams/{tenantName}",
		PathParams: map[string]string{
			"projectId":  projectID,
			"tenantName": name,
		},
	}
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), params, nil)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err == nil
}
