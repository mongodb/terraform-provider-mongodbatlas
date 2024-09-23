package projectinvitation_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigProjectInvitation_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		name         = acc.RandomEmail()
		roles        = []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_CLUSTER_MANAGER"}
		config       = configBasic(orgID, projectName, name, roles)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", roles[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", roles[1]),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
