package aimodelapikey_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceType         = "mongodbatlas_ai_model_api_key"
	resourceName         = resourceType + ".test"
	dataSourceName       = "data." + resourceType + ".test"
	dataSourcePluralName = "data." + resourceType + "s.test"
)

func TestAccAIModelAPIKey_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		name        = acc.RandomName()
		nameUpdated = name + "-updated"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name),
				Check:  checkBasic(projectID, name),
			},
			{
				Config: configBasic(projectID, nameUpdated),
				Check:  checkBasic(projectID, nameUpdated),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "api_key_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"secret"}, // secret is not imported.
			},
		},
	})
}

func configBasic(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ai_model_api_key" "test" {
			project_id = %[1]q
			name       = %[2]q
		}

		data "mongodbatlas_ai_model_api_key" "test" {
			project_id = %[1]q
			api_key_id = mongodbatlas_ai_model_api_key.test.api_key_id
		}

		data "mongodbatlas_ai_model_api_keys" "test" {
			project_id = %[1]q
			depends_on = [mongodbatlas_ai_model_api_key.test]
		}
	`, projectID, name)
}

func checkBasic(projectID, name string) resource.TestCheckFunc {
	commonAttrsSet := []string{"api_key_id", "created_at", "created_by", "masked_secret", "status"}
	commonAttrsMap := map[string]string{
		"project_id": projectID,
		"name":       name,
	}
	return resource.ComposeAggregateTestCheckFunc(
		acc.CheckRSAndDS(resourceName, admin.PtrString(dataSourceName), admin.PtrString(dataSourcePluralName), commonAttrsSet, commonAttrsMap, checkExists(resourceName)),
		// TODO: secret update check will fail until CLOUDP-373517 is done.
		// resource.TestCheckResourceAttrSet(resourceName, "secret"), // secret only in resource
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
	)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]
		if !apiKeyExists(rs) {
			return fmt.Errorf("api key (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		if apiKeyExists(rs) {
			return fmt.Errorf("api key (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs := s.RootModule().Resources[resourceName]
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["api_key_id"]), nil
	}
}

// apiKeyExists checks if an API key exists.
// Uses UntypedAPICall because the API is in preview and not yet available in the SDK.
// TODO: Use SDK before merging to master in CLOUDP-372674.
func apiKeyExists(rs *terraform.ResourceState) bool {
	callParams := config.APICallParams{
		VersionHeader: "application/vnd.atlas.preview+json",
		RelativePath:  "/api/atlas/v2/groups/{projectId}/aiModelApiKeys/{apiKeyId}",
		PathParams: map[string]string{
			"projectId": rs.Primary.Attributes["project_id"],
			"apiKeyId":  rs.Primary.Attributes["api_key_id"],
		},
		Method: "GET",
	}
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), callParams, nil)
	if resp != nil {
		resp.Body.Close()
	}
	return err == nil
}
