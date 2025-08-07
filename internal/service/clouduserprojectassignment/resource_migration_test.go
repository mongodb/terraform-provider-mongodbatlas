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
		roles       = []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"}
		rolesStr    = `"` + strings.Join(roles, `", "`) + `"`
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            originalConfigFirst(username, projectName, orgID, rolesStr),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   userProjectAssignmentConfigSecond(username, projectName, orgID, rolesStr),
				Check:                    checksSecond(username, projectName, roles),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceInvitationName, plancheck.ResourceActionDestroy),
					},
				},
				Config: removeProjectInvitationConfigThird(username, projectName, orgID, rolesStr),
			},
			mig.TestStepCheckEmptyPlan(removeProjectInvitationConfigThird(username, projectName, orgID, rolesStr)),
		},
	})
}

func originalConfigFirst(username, projectName, orgID, roles string) string {
	return fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]q]
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
`, username, roles, projectName, orgID)
}

func userProjectAssignmentConfigSecond(username, projectName, orgID, roles string) string {
	return fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]q]
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
		`, username, roles, projectName, orgID)
}

func removeProjectInvitationConfigThird(username, projectName, orgID, roles string) string {
	return fmt.Sprintf(`
		locals {
			username = %[1]q
			roles    = [%[2]q]
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
		`, username, roles, projectName, orgID)
}

func checksSecond(username, projectName string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceUserProjectAssignmentName, "username", username),
		resource.TestCheckResourceAttrSet(resourceUserProjectAssignmentName, "project_id"),
		resource.TestCheckResourceAttr(resourceUserProjectAssignmentName, "roles.#", fmt.Sprintf("%d", len(roles))),
	}

	for i, role := range roles {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(resourceUserProjectAssignmentName, fmt.Sprintf("roles.%d", i), role))
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
