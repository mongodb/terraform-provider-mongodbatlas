package ldapconfiguration_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"github.com/spf13/cast"
)

func TestAccMigrationLDAPConfiguration_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname    = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username    = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password    = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port        = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		authEnabled = true
		projectName = acc.RandomProjectName()
		config      = configBasic(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port))
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckLDAP(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "authentication_enabled", strconv.FormatBool(authEnabled)),
					resource.TestCheckResourceAttr(resourceName, "port", port),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
