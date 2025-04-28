package resourcepolicyapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_resource_policy_api.test"

func TestAccResourcePolicyAPI_basic(t *testing.T) {
	var (
		orgID      = os.Getenv("MONGODB_ATLAS_ORG_ID")
		policyName = "test-policy-autogen"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, policyName),
				Check:  checkBasic(),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func configBasic(orgID, policyName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_resource_policy" "test" {
			org_id = %[1]q
			name   = %[2]q
			description = "some description"
			
			policies = [
				{
					body = <<EOF
					forbid (
						principal,
						action == cloud::Action::"cluster.createEdit",
						resource
					) when {
						context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"aws"])
					};
					EOF
				},
				{
					body = <<EOF
					forbid (
						principal,
						action == cloud::Action::"project.edit",
						resource
					) when {
						context.project.ipAccessList.contains(ip("0.0.0.0/0"))
					};
					EOF
				}
			]
		}
	`, orgID, policyName)
}

func checkBasic() resource.TestCheckFunc {
	// adds checks for computed attributes not defined in config
	setAttrsChecks := []string{"id", "policies.0.id", "policies.1.id", "version", "created_date", "created_by_user.id", "created_by_user.name"}
	checks := acc.AddAttrSetChecks(resourceName, nil, setAttrsChecks...)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["id"]
		if orgID == "" || id == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().ResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("resource policy (%s/%s) does not exist", orgID, id)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_resource_policy_api" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["id"]
		if orgID == "" || id == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().ResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute()
		if err == nil {
			return fmt.Errorf("resource policy (%s/%s) still exists", orgID, id)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["id"]
		if orgID == "" || id == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", orgID, id), nil
	}
}
