package ldapverify_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMigLDAPVerify_basic(t *testing.T) {
	var (
		hostname     = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username     = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password     = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port         = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		projectID, _ = acc.ClusterNameExecution(t)
		config       = configBasic(projectID, hostname, username, password, cast.ToInt(port))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckLDAP(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "port", port),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
