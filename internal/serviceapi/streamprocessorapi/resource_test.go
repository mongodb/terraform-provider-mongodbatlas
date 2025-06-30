package streamprocessorapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_stream_processor_api.test"
)

func TestAccStreamProcessorAPI_basic(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		instanceName  = acc.RandomName()
		processorName = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, instanceName, processorName),
				Check:  checkBasic(projectID, instanceName, processorName),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"stats"},
				ImportStateVerifyIdentifierAttribute: "name", // id is not used because _id is returned in Atlas which is not a legal name for a Terraform attribute.
			},
		},
	})
}

func configBasic(projectID, instanceName, processorName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance_api" "test" {
			group_id = %[1]q
			name     = %[2]q
			data_process_region = {
				region         = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
			stream_config = {
				tier = "SP10"
			}
		}

		resource "mongodbatlas_stream_connection" "test" {
			project_id      = mongodbatlas_stream_instance_api.test.group_id
			instance_name   = mongodbatlas_stream_instance_api.test.name
			connection_name = "sample_stream_solar"
			type            = "Sample"
			depends_on      = [mongodbatlas_stream_instance_api.test]
		}

		resource "mongodbatlas_stream_processor_api" "test" {
			group_id      = mongodbatlas_stream_instance_api.test.group_id
			tenant_name   = mongodbatlas_stream_instance_api.test.name
			name          = %[3]q

			pipeline = [
				jsonencode({
					"$source" = {
						"connectionName" = "sample_stream_solar"
					}
				}),
				jsonencode({
					"$emit" = {
						"connectionName" = "__testLog"
					}
				})
			]

		  depends_on = [mongodbatlas_stream_connection.test]
		}
	`, projectID, instanceName, processorName)
}

func checkBasic(projectID, instanceName, processorName string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"group_id":    projectID,
		"tenant_name": instanceName,
		"name":        processorName,
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		tenantName := rs.Primary.Attributes["tenant_name"]
		processorName := rs.Primary.Attributes["name"]
		if groupID == "" || tenantName == "" || processorName == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), groupID, tenantName, processorName).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("stream processor(%s/%s/%s) does not exist", groupID, tenantName, processorName)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_processor_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		tenantName := rs.Primary.Attributes["tenant_name"]
		processorName := rs.Primary.Attributes["name"]
		if groupID == "" || tenantName == "" || processorName == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), groupID, tenantName, processorName).Execute()
		if err == nil {
			return fmt.Errorf("stream processor (%s/%s/%s) still exists", groupID, tenantName, processorName)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		tenantName := rs.Primary.Attributes["tenant_name"]
		processorName := rs.Primary.Attributes["name"]
		if groupID == "" || tenantName == "" || processorName == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", groupID, tenantName, processorName), nil
	}
}
