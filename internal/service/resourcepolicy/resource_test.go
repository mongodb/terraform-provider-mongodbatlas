package resourcepolicy_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceType = "mongodbatlas_resource_policy"
	resourceID   = fmt.Sprintf("%s.test", resourceType)
)

func TestAccResourcePolicy_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_RP_ORG_ID")
		policyName  = "test-policy"
		updatedName = "updated-policy"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckRPBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, &policyName),
				Check:  checksResourcePolicy(orgID, &policyName, 1),
			},
			{
				Config: configBasic(orgID, &updatedName),
				Check:  checksResourcePolicy(orgID, &updatedName, 1),
			},
			{
				Config:            configBasic(orgID, &updatedName),
				ResourceName:      resourceID,
				ImportStateIdFunc: checkImportStateIDFunc(resourceID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	},
	)
}

func checksResourcePolicy(orgID string, name *string, policyCount int) resource.TestCheckFunc {
	attrMap := map[string]string{
		"org_id":     orgID,
		"policies.#": fmt.Sprintf("%d", policyCount),
	}
	attrSet := []string{
		"created_by_user.id",
		"created_by_user.name",
		"created_date",
		"last_updated_by_user.id",
		"last_updated_by_user.name",
		"last_updated_date",
		"id",
		"version",
	}
	if name != nil {
		attrMap["name"] = *name
	}
	checks := []resource.TestCheckFunc{checkExists()}
	checks = acc.AddAttrChecks(resourceID, checks, attrMap)
	checks = acc.AddAttrSetChecks(resourceID, checks, attrSet...)
	for i := 0; i < policyCount; i++ {
		checks = acc.AddAttrSetChecks(resourceID, checks, fmt.Sprintf("policies.%d.body", i), fmt.Sprintf("policies.%d.id", i))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configBasic(orgID string, policyName *string) string {
	var name string
	if policyName != nil {
		name = fmt.Sprintf("  name = %q", *policyName)
	}
	return acc.ConfigRPProvider() + fmt.Sprintf(`
resource "mongodbatlas_resource_policy" "test" {
	org_id = %[1]q
	%[2]s
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
   }
 ]
}
`, orgID, name)
}

func checkImportStateIDFunc(resourceID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceID)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["id"]), nil
	}
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == resourceType {
				orgID := rs.Primary.Attributes["org_id"]
				id := rs.Primary.Attributes["id"]
				_, _, err := acc.ConnV2UsingResourcePolicy().AtlasResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute()
				if err != nil {
					return fmt.Errorf("resource policy (%s:%s) not found", orgID, id)
				}
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceType {
			orgID := rs.Primary.Attributes["org_id"]
			id := rs.Primary.Attributes["id"]
			_, _, err := acc.ConnV2UsingResourcePolicy().AtlasResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute()
			if err == nil {
				return fmt.Errorf("resource policy (%s:%s) still exists", orgID, id)
			}
		}
	}
	return nil
}
