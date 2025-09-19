package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccSTSAssumeRole_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSTSAssumeRole(t); acc.PreCheckRegularCredsAreEmpty(t) },
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
	var (
		resourceName = "data.mongodbatlas_projects.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckServiceAccount(t); acc.PreCheckRegularCredsAreEmpty(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDataSourceProject(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func configProject(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name  			 = %[2]q
		}
	`, orgID, projectName)
}

func configDataSourceProject() string {
	return `
		data "mongodbatlas_projects" "test" {
		}
	`
}
