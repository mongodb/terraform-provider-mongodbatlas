package orginvitation_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrgInvitation_basic(t *testing.T) {
	var (
		dataSourceName = "mongodbatlas_org_invitation.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = acc.RandomEmail()
		initialRole    = []string{"ORG_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyOrgInvitation,
		Steps: []resource.TestStep{
			{
				Config: configDS(orgID, name, initialRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "invitation_id"),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(dataSourceName, "username", name),
					resource.TestCheckResourceAttr(dataSourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, "roles.*", initialRole[0]),
				),
			},
		},
	})
}

func configDS(orgID, username string, roles []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_invitation" "test" {
			org_id   = %[1]q
			username = %[2]q
			roles  	 = ["%[3]s"]
		}

		data "mongodbatlas_org_invitation" "test" {
			org_id        = mongodbatlas_org_invitation.test.org_id
			username      = mongodbatlas_org_invitation.test.username
			invitation_id = mongodbatlas_org_invitation.test.invitation_id
		}`, orgID, username,
		strings.Join(roles, `", "`),
	)
}
