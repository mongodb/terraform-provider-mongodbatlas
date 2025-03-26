package resourcepolicy_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceType                      = "mongodbatlas_resource_policy"
	resourceID                        = fmt.Sprintf("%s.test", resourceType)
	dataSourceID                      = "data.mongodbatlas_resource_policy.test"
	dataSourcePluralID                = "data.mongodbatlas_resource_policies.test"
	invalidPolicyUnknownCloudProvider = `
	forbid (
	principal,
	action == cloud::Action::"cluster.createEdit",
	resource
	) when {
	context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"aws222"])
	};`
	invalidPolicyMissingComma = `
	forbid (
	principal,
	action == cloud::Action::"cluster.createEdit"
	resource
	) when {
	context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"aws"])
	};`
	validPolicyForbidAwsCloudProvider = `
	forbid (
	principal,
	action == cloud::Action::"cluster.createEdit",
	resource
	) when {
	context.cluster.cloudProviders.containsAny([cloud::cloudProvider::"aws"])
	};`
	validPolicyProjectForbidIPAccessAnywhere = `
	forbid (
		principal,
		action == cloud::Action::"project.edit",
		resource
	) 
		when {
		context.project.ipAccessList.contains(ip("0.0.0.0/0"))
	};`
)

func TestAccResourcePolicy_basic(t *testing.T) {
	tc := basicTestCase(t)
	resource.Test(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		policyName  = "test-policy"
		updatedName = "updated-policy"
	)
	return &resource.TestCase{ // Need sequential execution for assertions to be deterministic (plural data source)
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, policyName),
				Check:  checksResourcePolicy(orgID, policyName, 1),
			},
			{
				Config: configBasic(orgID, updatedName),
				Check:  checksResourcePolicy(orgID, updatedName, 1),
			},
			{
				Config:            configBasic(orgID, updatedName),
				ResourceName:      resourceID,
				ImportStateIdFunc: checkImportStateIDFunc(resourceID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccResourcePolicy_multipleNestedPolicies(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithPolicyBodies(orgID, "test-policy-multiple", validPolicyForbidAwsCloudProvider, validPolicyProjectForbidIPAccessAnywhere),
				Check:  checksResourcePolicy(orgID, "test-policy-multiple", 2),
			},
			{
				Config:            configWithPolicyBodies(orgID, "test-policy-multiple", validPolicyForbidAwsCloudProvider, validPolicyProjectForbidIPAccessAnywhere),
				ResourceName:      resourceID,
				ImportStateIdFunc: checkImportStateIDFunc(resourceID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	},
	)
}

func TestAccResourcePolicy_invalidConfig(t *testing.T) {
	var (
		orgID      = os.Getenv("MONGODB_ATLAS_ORG_ID")
		policyName = "test-policy-invalid"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configWithPolicyBodies(orgID, policyName, invalidPolicyMissingComma),
				ExpectError: regexp.MustCompile("unexpected token `resource`"),
			},
			{
				Config:      configWithPolicyBodies(orgID, policyName, invalidPolicyUnknownCloudProvider),
				ExpectError: regexp.MustCompile(`entity id aws222 does not exist in the context of this organization`),
			},
			{
				Config:      configWithPolicyBodies(orgID, policyName, validPolicyForbidAwsCloudProvider, invalidPolicyUnknownCloudProvider),
				ExpectError: regexp.MustCompile(`entity id aws222 does not exist in the context of this organization`),
			},
		},
	},
	)
}

func checksResourcePolicy(orgID, name string, policyCount int) resource.TestCheckFunc {
	attrMap := map[string]string{
		"org_id":     orgID,
		"policies.#": fmt.Sprintf("%d", policyCount),
		"name":       name,
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
	pluralMap := map[string]string{
		"org_id":    orgID,
		"results.#": "1",
	}
	checks := []resource.TestCheckFunc{checkExists()}
	checks = acc.AddAttrChecks(dataSourcePluralID, checks, pluralMap)
	for i := range policyCount {
		checks = acc.AddAttrSetChecks(resourceID, checks, fmt.Sprintf("policies.%d.body", i), fmt.Sprintf("policies.%d.id", i))
		checks = acc.AddAttrSetChecks(dataSourceID, checks, fmt.Sprintf("policies.%d.body", i), fmt.Sprintf("policies.%d.id", i))
		checks = acc.AddAttrSetChecks(dataSourcePluralID, checks, fmt.Sprintf("results.0.policies.%d.body", i), fmt.Sprintf("results.0.policies.%d.id", i))
	}
	// cannot use dataSourcePluralID as it doesn't have the `results` attribute
	return acc.CheckRSAndDS(resourceID, &dataSourceID, nil, attrSet, attrMap, resource.ComposeAggregateTestCheckFunc(checks...))
}

func configBasic(orgID, policyName string) string {
	return configWithPolicyBodies(orgID, policyName, validPolicyForbidAwsCloudProvider)
}

func configWithPolicyBodies(orgID, policyName string, description string, bodies ...string) string {
	policies := ""
	for _, body := range bodies {
		policies += fmt.Sprintf(`
		{
			body = <<EOF
			%s
			EOF
		},
		`, body)
	}
	return fmt.Sprintf(`
resource "mongodbatlas_resource_policy" "test" {
	org_id = %[1]q
	name   = %[2]q
	
	policies = [
%s
	]
}
data "mongodbatlas_resource_policy" "test" {
	org_id = mongodbatlas_resource_policy.test.org_id
	id = mongodbatlas_resource_policy.test.id
}
data "mongodbatlas_resource_policies" "test" {
	org_id = mongodbatlas_resource_policy.test.org_id
}
	`, orgID, policyName, policies, description)
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
				_, _, err := acc.ConnV2().ResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute()
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
			_, _, err := acc.ConnV2().ResourcePoliciesApi.GetAtlasResourcePolicy(context.Background(), orgID, id).Execute()
			if err == nil {
				return fmt.Errorf("resource policy (%s:%s) still exists", orgID, id)
			}
		}
	}
	return nil
}
