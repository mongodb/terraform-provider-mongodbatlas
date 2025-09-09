package clouduserprojectassignment_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

const (
	resourceInvitationName            = "mongodbatlas_project_invitation.mig_test"
	resourceProjectName               = "mongodbatlas_project.mig_test"
	resourceUserProjectAssignmentName = "mongodbatlas_cloud_user_project_assignment.user_mig_test"
)

func TestMigCloudUserProjectAssignmentRS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCloudUserProjectAssignmentRS_migrationJourney(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username    = acc.RandomEmail()
		projectName = fmt.Sprintf("mig_user_project_%s", acc.RandomName())
		roles       = []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            legacyProjectInvitationConfig(username, projectName, orgID, roles),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   userProjectAssignmentConfigSecond(username, projectName, orgID, roles),
				Check:                    checksSecond(username, roles),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceInvitationName, plancheck.ResourceActionDestroy),
					},
				},
				Config: removeProjectInvitationConfigThird(username, projectName, orgID, roles),
			},
			mig.TestStepCheckEmptyPlan(removeProjectInvitationConfigThird(username, projectName, orgID, roles)),
		},
	})
}

func legacyProjectInvitationConfig(username, projectName, orgID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	config := fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]s]
		}

		resource "mongodbatlas_project" "mig_test" {
			name 	= %[3]q
			org_id  = %[4]q
		}

		resource "mongodbatlas_project_invitation" "mig_test" {
			project_id  = mongodbatlas_project.mig_test.id
			username    = local.username
			roles       = local.roles
		}
	`, username, rolesStr, projectName, orgID)
	return config
}

func userProjectAssignmentConfigSecond(username, projectName, orgID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]s]
		}

		resource "mongodbatlas_project" "mig_test" {
			name 	= %[3]q
			org_id  = %[4]q
		}

		resource "mongodbatlas_project_invitation" "mig_test" {
			project_id  = mongodbatlas_project.mig_test.id
			username    = local.username
			roles       = local.roles
		}

		resource "mongodbatlas_cloud_user_project_assignment" "user_mig_test" {
			project_id = mongodbatlas_project.mig_test.id
			username   = local.username
			roles      = local.roles
		}
		`, username, rolesStr, projectName, orgID)
}

func removeProjectInvitationConfigThird(username, projectName, orgID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]s]
		}

		resource "mongodbatlas_project" "mig_test" {
			name 	= %[3]q
			org_id  = %[4]q
		}

		resource "mongodbatlas_cloud_user_project_assignment" "user_mig_test" {
			project_id = mongodbatlas_project.mig_test.id
			username   = local.username
			roles      = local.roles
		}
		`, username, rolesStr, projectName, orgID)
}

func checksSecond(username string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceUserProjectAssignmentName, "username", username),
		resource.TestCheckResourceAttrSet(resourceUserProjectAssignmentName, "project_id"),
		resource.TestCheckResourceAttr(resourceUserProjectAssignmentName, "roles.#", fmt.Sprintf("%d", len(roles))),
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
