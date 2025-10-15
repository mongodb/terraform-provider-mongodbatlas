package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccSTSAssumeRole_basic(t *testing.T) {
	acc.SkipInPAK(t, "skipping as this test is for AWS credentials only")
	acc.SkipInSA(t, "skipping as this test is for AWS credentials only")
	var (
		resourceName = "mongodbatlas_project.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSTSAssumeRole(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configProject(orgID, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"with_default_alerts_settings"},
			},
		},
	})
}

func TestAccServiceAccount_basic(t *testing.T) {
	acc.SkipInPAK(t, "skipping as this test is for SA only")
	var (
		resourceName = "data.mongodbatlas_organization.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDataSourceOrg(orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
				),
			},
		},
	})
}

func TestAccAccessToken_basic(t *testing.T) {
	acc.SkipTestForCI(t) // access token has a validity period of 1 hour, so it cannot be used in CI reliably
	acc.SkipInPAK(t, "skipping as this test is for Token credentials only")
	acc.SkipInSA(t, "skipping as this test is for Token credentials only")
	var (
		resourceName = "data.mongodbatlas_organization.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAccessToken(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDataSourceOrg(orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
				),
			},
		},
	})
}

func configProject(orgID, projectName string) string {
	// Use project in TPF and organization in SDKv2 so both providers are tested.
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name  			 = %[2]q
		}
		data "mongodbatlas_organization" "test" {
			org_id = %[1]q
		}
	`, orgID, projectName)
}

func configDataSourceOrg(orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_organization" "test" {
  			org_id = %[1]q
		}
	`, orgID)
}
