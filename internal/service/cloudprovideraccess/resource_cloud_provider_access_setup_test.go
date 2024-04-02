package cloudprovideraccess_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudprovideraccess"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccCloudProviderAccessSetupAWS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicSetupTestCase(t))
}

func TestAccCloudProviderAccessSetupAzure_basic(t *testing.T) {
	var (
		resourceName       = "mongodbatlas_cloud_provider_access_setup.test"
		dataSourceName     = "data.mongodbatlas_cloud_provider_access_setup.test"
		atlasAzureAppID    = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		tenantID           = os.Getenv("AZURE_TENANT_ID")
		projectID          = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCloudProviderAccessAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSetupAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.atlas_azure_app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.service_principal_id"),
					resource.TestCheckResourceAttrSet(resourceName, "azure_config.0.tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.atlas_azure_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.service_principal_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "azure_config.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated_date"),
				),
			},
		},
	},
	)
}

func basicSetupTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		resourceName   = "mongodbatlas_cloud_provider_access_setup.test"
		dataSourceName = "data.mongodbatlas_cloud_provider_access_setup.test"
		projectID      = acc.ProjectIDExecution(tb)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSetupAWS(projectID),
				Check: resource.ComposeTestCheckFunc(
					// same as regular cloud resource
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_config.0.atlas_assumed_role_external_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "aws_config.0.atlas_aws_account_arn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_date"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func configSetupAWS(projectID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = %[1]q
		provider_name = "AWS"
	 }

	 data "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_cloud_provider_access_setup.test.project_id
		provider_name = mongodbatlas_cloud_provider_access_setup.test.provider_name
		role_id =  mongodbatlas_cloud_provider_access_setup.test.role_id
	 }

	`, projectID)
}

func configSetupAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = %[1]q
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[2]q
			service_principal_id = %[3]q
			tenant_id = %[4]q
		}
	 }

	 data "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = mongodbatlas_cloud_provider_access_setup.test.project_id
		provider_name = "AWS"
		role_id =  mongodbatlas_cloud_provider_access_setup.test.role_id
	 }
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		providerName := ids["provider_name"]
		id := ids["id"]
		roles, _, err := acc.Conn().CloudProviderAccess.ListRoles(context.Background(), ids["project_id"])
		if err != nil {
			return fmt.Errorf(cloudprovideraccess.ErrorCloudProviderGetRead, err)
		}
		if providerName == "AWS" {
			for i := range roles.AWSIAMRoles {
				if roles.AWSIAMRoles[i].RoleID == id && roles.AWSIAMRoles[i].ProviderName == providerName {
					return nil
				}
			}
		}
		if providerName == "AZURE" {
			for i := range roles.AzureServicePrincipals {
				if *roles.AzureServicePrincipals[i].AzureID == id && roles.AzureServicePrincipals[i].ProviderName == providerName {
					return nil
				}
			}
		}
		return fmt.Errorf("error cloud Provider Access (%s) does not exist", ids["project_id"])
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["provider_name"], ids["id"]), nil
	}
}
