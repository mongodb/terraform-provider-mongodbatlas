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

// TestAccAPIUpdate_streamWorkspace_processorStatus exercises the
// mongodbatlas_api_update resource against the regional-failover private
// preview field on a stream workspace.
//
// Flow:
//  1. Typed mongodbatlas_stream_instance creates the workspace (GA fields).
//  2. mongodbatlas_api_update PATCHes processorStatus on the same workspace
//     using preview = true (preview content-type unlocks the field).
//  3. We change mode GRACEFUL -> FORCED in HCL: plan shows a single PATCH,
//     the typed resource is untouched.
//  4. We drop the api_update block: the typed stream_instance survives
//     destroy and Atlas keeps the field at its last-applied value.
func TestAccAPIUpdate_streamWorkspace_processorStatus(t *testing.T) {
	// Skipped: the processorStatus PATCH requires server-side state
	// (failoverRegions on the workspace + regional_failover_config.enabled on
	// processors) that the current public OpenAPI spec does not expose. We
	// cannot construct a valid end-to-end scenario from typed or generic
	// resources today. See docs-context/api-update-demo-target-investigation.md
	// for the full Glean trail (TD doc, MMS PR, and validation rules).
	//
	// The resource itself is fully exercised by:
	//   - TestUpdateResource_ValidateConfig_PreviewVersionHeaderMutex
	//   - manual run against this endpoint surfacing a 400 with the correct
	//     content-type and body (proves auth, preview negotiation, body
	//     marshalling, and error surfacing).
	t.Skip("blocked on streams failover prerequisites — see docs-context/api-update-demo-target-investigation.md")

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
				// Step 1: create workspace + initial patch (GRACEFUL).
				Config: configStreamWorkspaceWithFailover(projectID, instanceName, "GRACEFUL"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamInstanceExists(streamResourceName, projectID),
					resource.TestCheckResourceAttr(updateResourceName, "preview", "true"),
					resource.TestCheckResourceAttr(updateResourceName, "body.processorStatus.mode", "GRACEFUL"),
					resource.TestCheckResourceAttr(updateResourceName, "body.processorStatus.status", "PROCESSORS_STARTED"),
					resource.TestCheckResourceAttrSet(updateResourceName, "id"),
				),
			},
			{
				// Step 2: change mode -> FORCED. Only api_update should plan.
				Config: configStreamWorkspaceWithFailover(projectID, instanceName, "FORCED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamInstanceExists(streamResourceName, projectID),
					resource.TestCheckResourceAttr(updateResourceName, "body.processorStatus.mode", "FORCED"),
				),
			},
			{
				// Step 3: drop the api_update block — typed resource stays.
				Config: configStreamWorkspaceOnly(projectID, instanceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamInstanceExists(streamResourceName, projectID),
				),
			},
		},
	})
}

func configStreamWorkspaceWithFailover(projectID, name, mode string) string {
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
				processorStatus = {
					mode   = %[3]q
					region = "us-east-1"
					status = "PROCESSORS_STARTED"
				}
			}
		}
	`, projectID, name, mode)
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
