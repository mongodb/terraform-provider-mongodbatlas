package streamconnectionfailover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_stream_connection_failover.test"

// TestAccStreamConnectionFailover exercises the CRUD lifecycle of a failover stream connection:
// create, update its bootstrap servers, then import. The workspace must have failover regions
// enabled, and the failover connection shares the primary connection's name.
func TestAccStreamConnectionFailover(t *testing.T) {
	projectID, workspaceName := acc.ProjectIDExecutionWithStreamInstanceWithFailover(t)
	connectionName := "kafka-failover-primary"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configFailover(projectID, workspaceName, connectionName, "DUBLIN_IRL", "failover1:9092", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(),
					resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "region", "DUBLIN_IRL"),
					resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", "failover1:9092"),
					resource.TestCheckResourceAttr(resourceName, "config.auto.offset.reset", "earliest"),
					resource.TestCheckResourceAttrSet(resourceName, "failover_connection_id"),
				),
			},
			{
				// Remove the `config` block (same bootstrap) to exercise unsetting an optional field on
				// update. If the PATCH omits it and Atlas keeps the old value, this step fails with
				// "Provider produced inconsistent result after apply".
				Config: configFailover(projectID, workspaceName, connectionName, "DUBLIN_IRL", "failover1:9092", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(),
					resource.TestCheckNoResourceAttr(resourceName, "config.%"),
				),
			},
			{
				Config: configFailover(projectID, workspaceName, connectionName, "DUBLIN_IRL", "failover1-updated:9093", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", "failover1-updated:9093"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// This resource has no `id` attribute (id is aliased to failover_connection_id), so the
				// import-verify identity must point at the real identifier.
				ImportStateVerifyIdentifierAttribute: "failover_connection_id",
				// Secrets are never returned by the API, so they cannot be verified on import.
				ImportStateVerifyIgnore: []string{"authentication.password", "authentication.client_secret", "schema_registry_authentication.password"},
			},
		},
	})
}

func configFailover(projectID, workspaceName, connectionName, region, bootstrap string, includeFailoverConfig bool) string {
	failoverConfig := ""
	if includeFailoverConfig {
		failoverConfig = `config   = { "auto.offset.reset" = "earliest" }`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "primary" {
			project_id        = %[1]q
			workspace_name    = %[2]q
			connection_name   = %[3]q
			type              = "Kafka"
			bootstrap_servers = "localhost:9092"
			authentication = {
				mechanism = "PLAIN"
				username  = "user"
				password  = "rawpassword"
			}
			config   = { "auto.offset.reset" = "earliest" }
			security = { protocol = "SASL_PLAINTEXT" }
		}

		resource "mongodbatlas_stream_connection_failover" "test" {
			project_id        = %[1]q
			workspace_name    = %[2]q
			connection_name   = mongodbatlas_stream_connection.primary.connection_name
			region            = %[4]q
			type              = "Kafka"
			bootstrap_servers = %[5]q
			authentication = {
				mechanism = "PLAIN"
				username  = "fcuser"
				password  = "fcpass"
			}
			%[6]s
			security = { protocol = "SASL_PLAINTEXT" }
		}
	`, projectID, workspaceName, connectionName, region, bootstrap, failoverConfig)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf("%s/%s/%s/%s", a["project_id"], a["workspace_name"], a["connection_name"], a["failover_connection_id"]), nil
	}
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_connection_failover" {
				continue
			}
			a := rs.Primary.Attributes
			_, _, err := acc.ConnV2().StreamsApi.GetStreamFailoverConnection(context.Background(), a["project_id"], a["workspace_name"], a["connection_name"], a["failover_connection_id"]).Execute()
			if err != nil {
				return fmt.Errorf("failover connection (%s:%s:%s) does not exist", a["connection_name"], a["region"], a["failover_connection_id"])
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_connection_failover" {
			continue
		}
		a := rs.Primary.Attributes
		_, _, err := acc.ConnV2().StreamsApi.GetStreamFailoverConnection(context.Background(), a["project_id"], a["workspace_name"], a["connection_name"], a["failover_connection_id"]).Execute()
		if err == nil {
			return fmt.Errorf("failover connection (%s:%s) still exists", a["connection_name"], a["failover_connection_id"])
		}
	}
	return nil
}
