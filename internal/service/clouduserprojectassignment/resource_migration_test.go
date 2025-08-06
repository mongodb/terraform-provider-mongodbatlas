package clouduserprojectassignment_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudUserProjectAssignmentRS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCloudUserProjectAssignmentRS_migrationJourney(t *testing.T) {

}

func originalConfigFirst(t *testing.T) string {
	t.Helper()
	return `
	locals {
		username ="user_mig_test@email.com"
		roles    = [ "GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY" ]
	}

	resource "mongodbatlas_project" "mig_test" {
		name = "migration_user_project"
		org_id       = var.org_id
	}

	resource "mongodbatlas_project_invitation" "mig_test" {
		project_id  = mongodbatlas_project.mig_test.id
		username    = local.username
		roles       = local.roles
	}
`
}

func userProjectAssignmentConfigSecond(t *testing.T) string {
	t.Helper()
	return `
		locals {
		username ="user_mig_test@email.com"
		roles    = [ "GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY" ]
		}

		resource "mongodbatlas_project" "mig_test" {
			name = "migration_user_project"
			org_id       = var.org_id
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
		`
}
