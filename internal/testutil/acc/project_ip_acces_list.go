package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func CheckProjectIPAccessListExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		connV2 := TestMongoDBClient.(*config.MongoDBClient).AtlasV2

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		_, _, err := connV2.ProjectIPAccessListApi.GetProjectIpList(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err != nil {
			return fmt.Errorf("project ip access list entry (%s) does not exist", ids["entry"])
		}

		return nil
	}
}

func CheckDestroyProjectIPAccessList(s *terraform.State) error {
	connV2 := TestMongoDBClient.(*config.MongoDBClient).AtlasV2

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_access_list" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		_, _, err := connV2.ProjectIPAccessListApi.GetProjectIpList(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err == nil {
			return fmt.Errorf("project ip access list entry (%s) still exists", ids["entry"])
		}
	}

	return nil
}

func ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project.test.id
			ip_address = %[3]q
			comment    = %[4]q
		}
	`, orgID, projectName, ipAddress, comment)
}

func ConfigProjectIPAccessListWithCIDRBlock(orgID, projectName, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}

		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project.test.id
			cidr_block = %[3]q
			comment    = %[4]q
		}
	`, orgID, projectName, cidrBlock, comment)
}

func ConfigProjectIPAccessListWithAWSSecurityGroup(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment string) string {
	return fmt.Sprintf(`

		resource "mongodbatlas_project" "my_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = mongodbatlas_project.my_project.id
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = "%[3]s"
			region_name			  = "%[7]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			accepter_region_name	  = "us-east-1"
			project_id    			    = mongodbatlas_project.my_project.id
			container_id            = mongodbatlas_network_container.test.container_id
			provider_name           = "%[3]s"
			route_table_cidr_block  = "%[6]s"
			vpc_id					        = "%[4]s"
			aws_account_id	        = "%[5]s"
		}

		resource "mongodbatlas_project_ip_access_list" "test" {
			project_id         = mongodbatlas_project.my_project.id
			aws_security_group = "%[8]s"
			comment            = "%[9]s"

			depends_on = ["mongodbatlas_network_peering.test"]
		}
	`, orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegion, awsSGroup, comment)
}

func ConfigProjectIPAccessListWithMultiple(projectName, orgID string, accessList []map[string]string, isUpdate bool) string {
	cfg := fmt.Sprintf(`
			resource "mongodbatlas_project" "test" {
				name   = %[1]q
				org_id = %[2]q
			}`, projectName, orgID)

	for i, entry := range accessList {
		comment := entry["comment"]

		if isUpdate {
			comment = entry["comment"] + " update"
		}

		if cidr, ok := entry["cidr_block"]; ok {
			cfg += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
				project_id   = mongodbatlas_project.test.id
				cidr_block = %[2]q
				comment    = %[3]q
			}
		`, i, cidr, comment)
		} else {
			cfg += fmt.Sprintf(`
			resource "mongodbatlas_project_ip_access_list" "test_%[1]d" {
				project_id   = mongodbatlas_project.test.id
				ip_address = %[2]q
				comment    = %[3]q
			}
		`, i, entry["ip_address"], comment)
		}
	}
	return cfg
}

func ImportStateProjecIPAccessListtIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["entry"]), nil
	}
}
