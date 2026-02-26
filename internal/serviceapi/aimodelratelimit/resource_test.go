package aimodelratelimit_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceType            = "mongodbatlas_ai_model_rate_limit"
	resourceName            = resourceType + ".this"
	dataSourceName          = "data." + resourceType + ".this"
	dataSourcePluralName    = "data." + resourceType + "s.this"
	orgDataSourceType       = "mongodbatlas_ai_model_org_rate_limit"
	orgDataSourceName       = "data." + orgDataSourceType + ".this"
	orgDataSourcePluralName = "data." + orgDataSourceType + "s.this"
	modelGroupName          = "embed_large"
)

func TestAccAIModelRateLimit_basic(t *testing.T) {
	var (
		orgID     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID = acc.ProjectIDExecution(t)
	)
	// Serial test execution to avoid conflicts with other tests that use the same project and model group names.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		// TODO: CLOUDP-374704, CLOUDP-372674 - Implement CheckDestroy checking default limits when SDK can be used.
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectID, 100, 1000),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectID, 200, 2000),
				Check:  checkBasic(),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccAIModelRateLimit_invalidValues(t *testing.T) {
	projectID := acc.ProjectIDExecution(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configInvalid(projectID, "nonexistent_model_group", 100, 1000),
				ExpectError: regexp.MustCompile("RESOURCE_NOT_FOUND"),
			},
			{
				Config:      configInvalid(projectID, modelGroupName, 0, 1000),
				ExpectError: regexp.MustCompile("BAD_REQUEST"),
			},
			{
				Config:      configInvalid(projectID, modelGroupName, -1, 1000),
				ExpectError: regexp.MustCompile("BAD_REQUEST"),
			},
			{
				Config:      configInvalid(projectID, modelGroupName, 100, 0),
				ExpectError: regexp.MustCompile("BAD_REQUEST"),
			},
			{
				Config:      configInvalid(projectID, modelGroupName, 100, -1),
				ExpectError: regexp.MustCompile("BAD_REQUEST"),
			},
		},
	})
}

func configBasic(orgID, projectID string, requestsPerMinute, tokensPerMinute int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ai_model_rate_limit" "this" {
			project_id                 = %[2]q
			model_group_name           = %[3]q
			requests_per_minute_limit  = %[4]d
			tokens_per_minute_limit    = %[5]d
		}

		data "mongodbatlas_ai_model_rate_limit" "this" {
			project_id       = %[2]q
			model_group_name = mongodbatlas_ai_model_rate_limit.this.model_group_name
		}

		data "mongodbatlas_ai_model_rate_limits" "this" {
			project_id = %[2]q
			depends_on = [mongodbatlas_ai_model_rate_limit.this]
		}

		data "mongodbatlas_ai_model_org_rate_limit" "this" {
			org_id           = %[1]q
			model_group_name = mongodbatlas_ai_model_rate_limit.this.model_group_name
		}

		data "mongodbatlas_ai_model_org_rate_limits" "this" {
			org_id     = %[1]q
			depends_on = [mongodbatlas_ai_model_rate_limit.this]
		}
	`, orgID, projectID, modelGroupName, requestsPerMinute, tokensPerMinute)
}

func configInvalid(projectID, modelGroupName string, requestsPerMinute, tokensPerMinute int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ai_model_rate_limit" "this" {
			project_id                 = %[1]q
			model_group_name           = %[2]q
			requests_per_minute_limit  = %[3]d
			tokens_per_minute_limit    = %[4]d
		}
	`, projectID, modelGroupName, requestsPerMinute, tokensPerMinute)
}

func checkBasic() resource.TestCheckFunc {
	attrsSet := []string{"model_group_name", "requests_per_minute_limit", "tokens_per_minute_limit"}
	return resource.ComposeAggregateTestCheckFunc(
		acc.CheckRSAndDS(resourceName, new(dataSourceName), new(dataSourcePluralName), attrsSet, nil, checkExists(resourceName)),
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
		acc.CheckRSAndDS(orgDataSourceName, nil, new(orgDataSourcePluralName), attrsSet, nil),
		resource.TestCheckResourceAttrWith(orgDataSourcePluralName, "results.#", acc.IntGreatThan(0)),
	)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]
		if !rateLimitExists(rs) {
			return fmt.Errorf("rate limit (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs := s.RootModule().Resources[resourceName]
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["model_group_name"]), nil
	}
}

// rateLimitExists checks if a rate limit exists.
// Uses UntypedAPICall because the API is in preview and not yet available in the SDK.
// TODO: CLOUDP-374704 - Use SDK before merging to master in CLOUDP-372674.
func rateLimitExists(rs *terraform.ResourceState) bool {
	callParams := config.APICallParams{
		VersionHeader: "application/vnd.atlas.preview+json",
		RelativePath:  "/api/atlas/v2/groups/{groupId}/aiModelRateLimits/{modelGroupName}",
		PathParams: map[string]string{
			"groupId":        rs.Primary.Attributes["project_id"],
			"modelGroupName": rs.Primary.Attributes["model_group_name"],
		},
		Method: "GET",
	}
	resp, err := acc.MongoDBClient.UntypedAPICall(context.Background(), callParams, nil)
	if resp != nil {
		resp.Body.Close()
	}
	return err == nil
}
