package privatelinkendpointserverless_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_privatelink_endpoint_serverless.test"
)

func TestAccServerlessPrivateLinkEndpoint_basic(t *testing.T) {
	acc.SkipTestForCI(t)
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID    = acc.ProjectIDExecution(tb)
		instanceName = acc.RandomClusterName()
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, instanceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
				),
			},
			{
				Config:                  configBasic(projectID, instanceName, false),
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_strings_private_endpoint_srv"},
			},
		},
	}
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

func configBasic(projectID, instanceName string, ignoreConnectionStrings bool) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id    = mongodbatlas_serverless_instance.test.project_id
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	}

	%s
	`, acc.ConfigServerlessInstance(projectID, instanceName, ignoreConnectionStrings, nil, nil))
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

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s--%s--%s", ids["project_id"], ids["instance_name"], ids["endpoint_id"]), nil
	}
}
