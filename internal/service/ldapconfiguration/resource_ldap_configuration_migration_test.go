package ldapconfiguration_test

import (
	"os"
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
		CheckDestroy: acc.CheckDestroyLDAPConfiguration,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
