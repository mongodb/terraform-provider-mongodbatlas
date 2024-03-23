package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func CheckProjectIPAccessListExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := ConnV2().ProjectIPAccessListApi.GetProjectIpList(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err != nil {
			return fmt.Errorf("project ip access list entry (%s) does not exist", ids["entry"])
		}
		return nil
	}
}

func CheckDestroyProjectIPAccessList(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_access_list" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := ConnV2().ProjectIPAccessListApi.GetProjectIpList(context.Background(), ids["project_id"], ids["entry"]).Execute()
		if err == nil {
			return fmt.Errorf("project ip access list entry (%s) still exists", ids["entry"])
		}
	}
	return nil
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
