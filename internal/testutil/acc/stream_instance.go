package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func StreamInstanceConfig(projectID, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region = %[3]q
				cloud_provider = %[4]q
			}
		}
	`, projectID, instanceName, region, cloudProvider)
}

func StreamInstanceWithStreamConfigConfig(projectID, instanceName, region, cloudProvider, configTier string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region = %[3]q
				cloud_provider = %[4]q
			}
			stream_config = {
				tier = %[5]q
			}
		}
	`, projectID, instanceName, region, cloudProvider, configTier)
}

func CheckDestroyStreamInstance(state *terraform.State) error {
	if projectDestroyedErr := CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_stream_instance" {
			_, _, err := ConnV2().StreamsApi.GetStreamWorkspace(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"]).Execute()
			if err == nil {
				return fmt.Errorf("stream instance (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"])
			}
		}
	}
	return nil
}
