package streamconnection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const failoverResourceName = "mongodbatlas_stream_connection_failover.test"

// TestAccStreamRSStreamConnectionFailover exercises the CRUD lifecycle of a failover stream
// connection: create, update its bootstrap servers, then import. The workspace must have failover
// regions enabled, and the failover connection shares the primary connection's name.
func TestAccStreamRSStreamConnectionFailover(t *testing.T) {
	projectID, workspaceName := acc.ProjectIDExecutionWithStreamInstanceWithFailover(t)
	connectionName := "kafka-failover-primary"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamConnectionFailover,
		Steps: []resource.TestStep{
			{
				Config: configFailover(projectID, workspaceName, connectionName, "DUBLIN_IRL", "failover1:9092"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFailoverConnectionExists(),
					resource.TestCheckResourceAttr(failoverResourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(failoverResourceName, "region", "DUBLIN_IRL"),
					resource.TestCheckResourceAttr(failoverResourceName, "type", "Kafka"),
					resource.TestCheckResourceAttr(failoverResourceName, "bootstrap_servers", "failover1:9092"),
					resource.TestCheckResourceAttrSet(failoverResourceName, "id"),
				),
			},
			{
				Config: configFailover(projectID, workspaceName, connectionName, "DUBLIN_IRL", "failover1-updated:9093"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFailoverConnectionExists(),
					resource.TestCheckResourceAttr(failoverResourceName, "bootstrap_servers", "failover1-updated:9093"),
				),
			},
			{
				ResourceName:      failoverResourceName,
				ImportStateIdFunc: failoverImportStateIDFunc(failoverResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// Secrets are never returned by the API, so they cannot be verified on import.
				ImportStateVerifyIgnore: []string{"authentication.password", "authentication.client_secret", "schema_registry_authentication.password"},
			},
		},
	})
}

func configFailover(projectID, workspaceName, connectionName, region, bootstrap string) string {
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
			config   = { "auto.offset.reset" = "earliest" }
			security = { protocol = "SASL_PLAINTEXT" }
		}
	`, projectID, workspaceName, connectionName, region, bootstrap)
}

func failoverWorkspaceName(a map[string]string) string {
	if ws := a["workspace_name"]; ws != "" {
		return ws
	}
	return a["instance_name"]
}

func failoverImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		a := rs.Primary.Attributes
		return fmt.Sprintf("%s-%s-%s-%s", failoverWorkspaceName(a), a["project_id"], a["connection_name"], a["id"]), nil
	}
}

func checkFailoverConnectionExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_connection_failover" {
				continue
			}
			a := rs.Primary.Attributes
			_, _, err := acc.ConnV2().StreamsApi.GetStreamFailoverConnection(context.Background(), a["project_id"], failoverWorkspaceName(a), a["connection_name"], a["id"]).Execute()
			if err != nil {
				return fmt.Errorf("failover connection (%s:%s:%s) does not exist", a["connection_name"], a["region"], a["id"])
			}
		}
		return nil
	}
}

func checkDestroyStreamConnectionFailover(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_connection_failover" {
			continue
		}
		a := rs.Primary.Attributes
		_, _, err := acc.ConnV2().StreamsApi.GetStreamFailoverConnection(context.Background(), a["project_id"], failoverWorkspaceName(a), a["connection_name"], a["id"]).Execute()
		if err == nil {
			return fmt.Errorf("failover connection (%s:%s) still exists", a["connection_name"], a["id"])
		}
	}
	return nil
}
