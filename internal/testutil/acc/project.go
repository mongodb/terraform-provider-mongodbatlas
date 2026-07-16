package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20250312022/admin"
)

func CheckDestroyProject(s *terraform.State) error {
	return checkDestroyProject(ConnV2(), s)
}

func CheckDestroyProjectGov(s *terraform.State) error {
	return checkDestroyProject(ConnV2UsingGov(), s)
}

func checkDestroyProject(conn *admin.APIClient, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project" {
			continue
		}
		projectRes, _, _ := conn.ProjectsAPI.GetGroupByName(context.Background(), rs.Primary.ID).Execute()
		if projectRes != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func ImportStateProjectIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
