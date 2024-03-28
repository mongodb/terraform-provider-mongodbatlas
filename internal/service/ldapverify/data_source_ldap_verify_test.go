package ldapverify_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
)

func TestAccLDAPVerifyDS_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_ldap_verify.test"
		hostname       = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username       = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password       = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port           = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		projectID, _   = acc.ClusterNameExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, hostname, username, password, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "request_id"),
					resource.TestCheckResourceAttr(dataSourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(dataSourceName, "bind_username", username),
					resource.TestCheckResourceAttr(dataSourceName, "port", port),
				),
			},
		},
	})
}
