package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func StreamInstanceConfig(orgID, projectName, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_stream_instance" "test" {
			project_id = mongodbatlas_project.test.id
			instance_name = %[3]q
			data_process_region = {
				region = %[4]q
				cloud_provider = %[5]q
			}
		}
	`, orgID, projectName, instanceName, region, cloudProvider)
}

func CheckDestroyStreamInstance(state *terraform.State) error {
	if projectDestroyedErr := CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	connV2 := TestMongoDBClient.(*config.MongoDBClient).AtlasV2
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_stream_instance" {
			_, _, err := connV2.StreamsApi.GetStreamInstance(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"]).Execute()
			if err == nil {
				return fmt.Errorf("stream instance (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"])
			}
		}
	}
	return nil
}
