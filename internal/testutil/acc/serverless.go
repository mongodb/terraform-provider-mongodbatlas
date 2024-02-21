package acc

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func ConfigServerlessInstanceBasic(orgID, projectName, name string, ignoreConnectionStrings bool) string {
	lifecycle := ""

	if ignoreConnectionStrings {
		lifecycle = `

		lifecycle {
			ignore_changes = [connection_strings_private_endpoint_srv]
		}
		`
	}

	return fmt.Sprintf(serverlessConfig, orgID, projectName, name, lifecycle)
}

func ConfigServerlessInstanceWithTags(orgID, projectName, name string, tags []admin.ResourceTag) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = %q
				value = %q
			}
		`, label.GetKey(), label.GetValue())
	}
	return fmt.Sprintf(serverlessConfig, orgID, projectName, name, tagsConf)
}

const serverlessConfig = `
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q

		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true
		%[4]s
	}

	data "mongodbatlas_serverless_instance" "test" {
		name        = mongodbatlas_serverless_instance.test.name
		project_id  = mongodbatlas_serverless_instance.test.project_id
	}

	data "mongodbatlas_serverless_instances" "test" {
		project_id         = mongodbatlas_serverless_instance.test.project_id
	}
`
