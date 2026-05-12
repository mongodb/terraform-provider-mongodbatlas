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

// TestAccAPIResource_aiModelAPIKey_basic mirrors the typed
// mongodbatlas_ai_model_api_key acceptance test (CLOUDP-371552-dev-voyage)
// using the generic resource. Exercises the preview API channel and the
// "secret returned only on create" pattern (secret on create, maskedSecret
// on subsequent reads).
func TestAccAPIResource_aiModelAPIKey_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		name        = acc.RandomName()
		nameUpdated = name + "-updated"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyAIModelAPIKey(projectID),
		Steps: []resource.TestStep{
			{
				Config: configAIModelAPIKey(projectID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsAIModelAPIKey(resourceName, projectID),
					resource.TestCheckResourceAttr(resourceName, "preview", "true"),
					resource.TestCheckResourceAttr(resourceName, "body.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "output.apiKeyId"),
					resource.TestCheckResourceAttrSet(resourceName, "output.createdAt"),
					resource.TestCheckResourceAttrSet(resourceName, "output.status"),
					// secret is present after Create.
					resource.TestCheckResourceAttrSet(resourceName, "output.secret"),
					// data source surfaces the masked variant (read response).
					resource.TestCheckResourceAttrSet(dataSourceName, "output.maskedSecret"),
				),
			},
			{
				Config: configAIModelAPIKey(projectID, nameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsAIModelAPIKey(resourceName, projectID),
					resource.TestCheckResourceAttr(resourceName, "body.name", nameUpdated),
					// After refresh the API stops returning secret; maskedSecret remains.
					resource.TestCheckResourceAttrSet(resourceName, "output.maskedSecret"),
				),
			},
		},
	})
}

func configAIModelAPIKey(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_resource" "test" {
			path         = "/api/atlas/v2/groups/%[1]s/aiModelApiKeys"
			id_attribute = ["apiKeyId"]
			preview      = true

			body = {
				name = %[2]q
			}
		}

		data "mongodbatlas_api_resource" "test" {
			path    = mongodbatlas_api_resource.test.id
			preview = true
		}
	`, projectID, name)
}

func checkExistsAIModelAPIKey(rsName, projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("not found: %s", rsName)
		}
		apiKeyID := rs.Primary.Attributes["output.apiKeyId"]
		if apiKeyID == "" {
			return fmt.Errorf("checkExists: output.apiKeyId not set for %s", rsName)
		}
		if !aiModelAPIKeyExists(projectID, apiKeyID) {
			return fmt.Errorf("ai model api key (%s/%s) does not exist", projectID, apiKeyID)
		}
		return nil
	}
}

func checkDestroyAIModelAPIKey(projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_api_resource" {
				continue
			}
			apiKeyID := rs.Primary.Attributes["output.apiKeyId"]
			if apiKeyID == "" {
				continue
			}
			if aiModelAPIKeyExists(projectID, apiKeyID) {
				return fmt.Errorf("ai model api key (%s/%s) still exists", projectID, apiKeyID)
			}
		}
		return nil
	}
}

// aiModelAPIKeyExists uses UntypedAPICall because the API is in preview and
// not yet available in the SDK (matches the typed resource's test approach).
func aiModelAPIKeyExists(projectID, apiKeyID string) bool {
	params := config.APICallParams{
		Method:        http.MethodGet,
		VersionHeader: "application/vnd.atlas.preview+json",
		RelativePath:  "/api/atlas/v2/groups/{projectId}/aiModelApiKeys/{apiKeyId}",
		PathParams: map[string]string{
			"projectId": projectID,
			"apiKeyId":  apiKeyID,
		},
	}
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), params, nil)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err == nil
}
