package x509authenticationdatabaseuser_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGenericX509AuthDBUser_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		config    = configBasic(projectID, username)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			mig.PreCheckBasic(t)
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigGenericX509AuthDBUser_withCustomerX509(t *testing.T) {
	var (
		cas         = os.Getenv("CA_CERT")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid CANNOT_GENERATE_CERT_IF_ADVANCED_X509
		config      = configWithCustomerX509(orgID, projectName, cas)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckCert(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "customer_x509_cas"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
