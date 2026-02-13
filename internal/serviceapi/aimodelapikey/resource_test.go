package aimodelapikey_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceType            = "mongodbatlas_ai_model_api_key"
	resourceName            = resourceType + ".this"
	dataSourceName          = "data." + resourceType + ".this"
	dataSourcePluralName    = "data." + resourceType + "s.this"
	orgDataSourceType       = "mongodbatlas_ai_model_org_api_key"
	orgDataSourceName       = "data." + orgDataSourceType + ".this"
	orgDataSourcePluralName = "data." + orgDataSourceType + "s.this"
)

func TestAccAIModelAPIKey_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
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
				Config: configBasic(orgID, projectID, name),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectID, nameUpdated),
				Check:  checkBasic(),
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

func configBasic(orgID, projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ai_model_api_key" "this" {
			project_id = %[2]q
			name       = %[3]q
		}

		data "mongodbatlas_ai_model_api_key" "this" {
			project_id = %[2]q
			api_key_id = mongodbatlas_ai_model_api_key.this.api_key_id
		}

		data "mongodbatlas_ai_model_api_keys" "this" {
			project_id = %[2]q
			depends_on = [mongodbatlas_ai_model_api_key.this]
		}

		data "mongodbatlas_ai_model_org_api_key" "this" {
			org_id     = %[1]q
			api_key_id = mongodbatlas_ai_model_api_key.this.api_key_id
		}

		data "mongodbatlas_ai_model_org_api_keys" "this" {
			org_id     = %[1]q
			depends_on = [mongodbatlas_ai_model_api_key.this]
		}
	`, orgID, projectID, name)
}

func checkBasic() resource.TestCheckFunc {
	attrsSet := []string{"api_key_id", "created_at", "created_by", "masked_secret", "status", "name", "project_id"}
	return resource.ComposeAggregateTestCheckFunc(
		acc.CheckRSAndDS(resourceName, admin.PtrString(dataSourceName), admin.PtrString(dataSourcePluralName), attrsSet, nil, checkExists(resourceName)),
		resource.TestCheckResourceAttrSet(resourceName, "secret"), // secret only in resource
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
		acc.CheckRSAndDS(orgDataSourceName, nil, admin.PtrString(orgDataSourcePluralName), attrsSet, nil),
		resource.TestCheckResourceAttrWith(orgDataSourcePluralName, "results.#", acc.IntGreatThan(0)),
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
