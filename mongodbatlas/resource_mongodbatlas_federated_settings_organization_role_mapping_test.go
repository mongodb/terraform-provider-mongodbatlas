package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasFederatedSettingsOrganizationRoleMapping_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		federatedSettingsOrganizationRoleMapping matlas.FederatedSettingsOrganizationRoleMapping
		resourceName                             = "mongodbatlas_cloud_federated_settings_org_role_mapping.test"
		federationSettingsID                     = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                                    = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		groupID                                  = os.Getenv("MONGODB_ATLAS_FEDERATED_GROUP_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasFederatedSettingsOrganizationRoleMappingConfig(federationSettingsID, orgID, groupID),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingExists(resourceName, &federatedSettingsOrganizationRoleMapping),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "external_group_name", "newgroup"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasFederatedSettingsOrganizationRoleMapping_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName         = "mongodbatlas_cloud_federated_settings_org_role_mapping.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		groupID              = os.Getenv("MONGODB_ATLAS_FEDERATED_GROUP_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingDestroy,
		Steps: []resource.TestStep{

			{
				Config:            testAccMongoDBAtlasFederatedSettingsOrganizationRoleMappingConfig(federationSettingsID, orgID, groupID),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingImportStateIDFunc(resourceName),
				ImportState:       false,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingExists(resourceName string, snapshotExportJob *matlas.FederatedSettingsOrganizationRoleMapping) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		response, _, err := conn.FederatedSettings.GetRoleMapping(context.Background(), ids["federation_settings_id"], ids["org_id"], ids["role_mapping_id"])
		if err == nil {
			*snapshotExportJob = *response
			return nil
		}

		return fmt.Errorf("role mapping (%s) does not exist", ids["role_mapping_id"])
	}
}

func testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_federated_settings_org_role_mapping" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		snapshotExportBucket, _, err := conn.FederatedSettings.GetRoleMapping(context.Background(), ids["federation_settings_id"], ids["org_id"], ids["role_mapping_id"])
		if err == nil && snapshotExportBucket != nil {
			return fmt.Errorf("identity provider (%s) still exists", ids["idp_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["federation_settings_id"], ids["org_id"], ids["role_mapping_id"]), nil
	}
}

func testAccMongoDBAtlasFederatedSettingsOrganizationRoleMappingConfig(federationSettingsID, orgID, groupID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_federated_settings_org_role_mapping" "test" {
		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		external_group_name    = "newgroup"
	  
		organization_roles = ["ORG_OWNER", "ORG_MEMBER"]
		group_id           = "%[3]s"
		group_roles        = ["GROUP_OWNER", "GROUP_CLUSTER_MANAGER"]
	  
	  }`, federationSettingsID, orgID, groupID)
}
