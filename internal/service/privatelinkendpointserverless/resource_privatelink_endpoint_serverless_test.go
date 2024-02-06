package privatelinkendpointserverless_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccServerlessPrivateLinkEndpoint_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint_serverless.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName = "serverlessplink"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
				),
			},
		},
	})
}

func TestAccServerlessPrivateLinkEndpoint_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint_serverless.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-serverless")
		instanceName = "serverlessimport"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
				),
			},
			{
				Config:                  configBasic(orgID, projectName, instanceName, false),
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFuncBasic(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_strings_private_endpoint_srv"},
			},
		},
	})
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_serverless" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		privateLink, _, err := acc.ConnV2().ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"]).Execute()
		if err == nil && privateLink != nil {
			return fmt.Errorf("endpoint_id (%s) still exists", ids["endpoint_id"])
		}
	}
	return nil
}

func configBasic(orgID, projectName, instanceName string, ignoreConnectionStrings bool) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id    = mongodbatlas_serverless_instance.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	}

	%s
	`, acc.ConfigServerlessInstanceBasic(orgID, projectName, instanceName, ignoreConnectionStrings))
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ServerlessPrivateEndpointsApi.GetServerlessPrivateEndpoint(context.Background(), ids["project_id"], ids["instance_name"], ids["endpoint_id"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("endpoint_id (%s) does not exist", ids["endpoint_id"])
	}
}

func importStateIDFuncBasic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["instance_name"], ids["endpoint_id"]), nil
	}
}
