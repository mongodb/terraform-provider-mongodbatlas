package acc

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func ConfigServerlessInstance(projectID, name string, ignoreConnectionStrings bool, autoIndexing *bool, tags []admin.ResourceTag) string {
	var extra string

	if ignoreConnectionStrings {
		extra += `
			lifecycle {
				ignore_changes = [connection_strings_private_endpoint_srv]
			}
		`
	}
	if autoIndexing != nil {
		extra += fmt.Sprintf(`
			auto_indexing = %t
		`, *autoIndexing)
	}
	for _, label := range tags {
		extra += fmt.Sprintf(`
			tags {
				key   = %q
				value = %q
			}
		`, label.GetKey(), label.GetValue())
	}

	return fmt.Sprintf(serverlessConfig, projectID, name, extra)
}

const serverlessConfig = `
	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = %[1]q
		name         = %[2]q

		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true
		%[3]s
	}

	data "mongodbatlas_serverless_instance" "test" {
		name        = mongodbatlas_serverless_instance.test.name
		project_id  = mongodbatlas_serverless_instance.test.project_id
	}

	data "mongodbatlas_serverless_instances" "test" {
		project_id         = mongodbatlas_serverless_instance.test.project_id
	}
`
