package streamconnectionapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_stream_connection_api.test"

func TestAccStreamConnectionAPI_basic(t *testing.T) {
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName           = fmt.Sprintf("kafka-conn-api-%s", acc.RandomName())
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configureKafka(
					projectID,
					workspaceName,
					connectionName,
					"localhost:9092,localhost:9092",
					"earliest",
					"user",
					"rawpassword",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
					resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", "localhost:9092,localhost:9092"),
					resource.TestCheckResourceAttr(resourceName, "config.auto.offset.reset", "earliest"),
					resource.TestCheckResourceAttr(resourceName, "authentication.mechanism", "PLAIN"),
					resource.TestCheckResourceAttr(resourceName, "authentication.username", "user"),
					resource.TestCheckResourceAttr(resourceName, "security.protocol", "SASL_PLAINTEXT"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "connection_name",
				ImportStateVerifyIgnore:              []string{"authentication.password", "delete_on_create_timeout"},
			},
		},
	})
}

func configureKafka(projectID, workspaceName, connectionName, bootstrapServers, autoOffsetReset, username, password string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection_api" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			connection_name = %[3]q
			type = "Kafka"
			bootstrap_servers = %[4]q
			config = {
				"auto.offset.reset" = %[5]q
			}
			authentication = {
				mechanism = "PLAIN"
				username = %[6]q
				password = %[7]q
			}
			security = {
				protocol = "SASL_PLAINTEXT"
			}
			networking = {
				access = {
					type = "PUBLIC"
				}
			}
		}
	`, projectID, workspaceName, connectionName, bootstrapServers, autoOffsetReset, username, password)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		connectionName := rs.Primary.Attributes["connection_name"]
		if projectID == "" || workspaceName == "" || connectionName == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}

		if _, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, workspaceName, connectionName).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("stream connection (%s/%s/%s) does not exist", projectID, workspaceName, connectionName)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_connection_api" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		connectionName := rs.Primary.Attributes["connection_name"]
		if projectID == "" || workspaceName == "" || connectionName == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, workspaceName, connectionName).Execute()
		if err == nil {
			return fmt.Errorf("stream connection (%s/%s/%s) still exists", projectID, workspaceName, connectionName)
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

		projectID := rs.Primary.Attributes["project_id"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		connectionName := rs.Primary.Attributes["connection_name"]
		if projectID == "" || workspaceName == "" || connectionName == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", workspaceName, projectID, connectionName), nil
	}
}
