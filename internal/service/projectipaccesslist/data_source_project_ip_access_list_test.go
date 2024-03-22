package projectipaccesslist_test

import (
	"fmt"
)

func testAccDataMongoDBAtlasProjectIPAccessListConfigSettingIPAddress(orgID, projectName, ipAddress, comment string) string {
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

		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			ip_address = mongodbatlas_project_ip_access_list.test.ip_address
		}
	`, orgID, projectName, ipAddress, comment)
}

func testAccDataMongoDBAtlasProjectIPAccessListConfigSettingCIDRBlock(orgID, projectName, cidrBlock, comment string) string {
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
		data "mongodbatlas_project_ip_access_list" "test" {
			project_id = mongodbatlas_project_ip_access_list.test.project_id
			cidr_block = mongodbatlas_project_ip_access_list.test.cidr_block
		}
	`, orgID, projectName, cidrBlock, comment)
}
