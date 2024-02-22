package ldapconfiguration_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
)

func TestAccLDAPConfigurationDS_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_ldap_configuration.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname       = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username       = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password       = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port           = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		authEnabled    = true
		projectName    = acc.RandomProjectName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyLDAPConfiguration,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hostname"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authentication_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "port"),
				),
			},
		},
	})
}
