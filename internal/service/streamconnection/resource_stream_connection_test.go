package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

//go:embed testdata/dummy-ca.pem
var DummyCACert string

func TestAccStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-stream")
		instanceName = acctest.RandomWithPrefix("test-acc-name")
		resourceName = "mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false, true),
			},
			{
				Config: kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user2", "otherpassword", "localhost:9093", "latest", false),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, orgID, projectName, instanceName, "user2", "otherpassword", "localhost:9093", "latest", false, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication.password"},
			},
		},
	})
}

func TestAccStreamRSStreamConnection_kafkaSSL(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-stream")
		instanceName = acctest.RandomWithPrefix("test-acc-name")
		resourceName = "mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication.password"},
			},
		},
	})
}

func TestAccStreamRSStreamConnection_cluster(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo  = acc.GetClusterInfo(orgID)
		instanceName = acctest.RandomWithPrefix("test-acc-name")
		resourceName = "mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: clusterStreamConnectionConfig(clusterInfo.ProjectIDStr, instanceName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr),
				Check:  clusterStreamConnectionAttributeChecks(resourceName, clusterInfo.ClusterName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func kafkaStreamConnectionConfig(orgID, projectName, instanceName, username, password, bootstrapServers, configValue string, useSSL bool) string {
	projectAndStreamInstanceConfig := acc.StreamInstanceConfig(orgID, projectName, instanceName, "VIRGINIA_USA", "AWS")
	securityConfig := `
		security = {
			protocol = "PLAINTEXT"
		}`

	if useSSL {
		securityConfig = fmt.Sprintf(`
		security = {
		    broker_public_certificate = %q
		    protocol = "SSL"
		}`, DummyCACert)
	}

	return fmt.Sprintf(`
		%[1]s
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_project.test.id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = "ConnectionNameKafka"
		 	type = "Kafka"
		 	authentication = {
		    	mechanism = "PLAIN"
		    	username = %[2]q
		    	password = %[3]q
		    }
		    bootstrap_servers = %[4]q
		    config = {
		    	"auto.offset.reset": %[5]q
		    }
		    %[6]s
		}
	`, projectAndStreamInstanceConfig, username, password, bootstrapServers, configValue, securityConfig)
}

func kafkaStreamConnectionAttributeChecks(
	resourceName, orgID, projectName, instanceName, username, password, bootstrapServers, configValue string, usesSSL, checkPassword bool) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
		resource.TestCheckResourceAttr(resourceName, "authentication.mechanism", "PLAIN"),
		resource.TestCheckResourceAttr(resourceName, "authentication.username", username),
		resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", bootstrapServers),
		resource.TestCheckResourceAttr(resourceName, "config.auto.offset.reset", configValue),
	}
	if checkPassword {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "authentication.password", password))
	}
	if !usesSSL {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "security.protocol", "PLAINTEXT"))
	} else {
		resourceChecks = append(resourceChecks,
			resource.TestCheckResourceAttr(resourceName, "security.protocol", "SSL"),
			resource.TestCheckResourceAttrSet(resourceName, "security.broker_public_certificate"),
		)
	}
	return resource.ComposeTestCheckFunc(resourceChecks...)
}

func clusterStreamConnectionConfig(projectIDStr, instanceName, clusterNameStr, clusterTerraformStr string) string {
	return clusterTerraformStr + fmt.Sprintf(`
		
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]s
			instance_name = %[2]q
			data_process_region = {
				region = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
		}
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = "ConnectionNameKafka"
		 	type = "Cluster"
		 	cluster_name = %[3]s
		}
	`, projectIDStr, instanceName, clusterNameStr)
}

func clusterStreamConnectionAttributeChecks(resourceName, clusterName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Cluster"),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
	}
	return resource.ComposeTestCheckFunc(resourceChecks...)
}

func checkStreamConnectionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["instance_name"], rs.Primary.Attributes["project_id"], rs.Primary.Attributes["connection_name"]), nil
	}
}

func checkStreamConnectionExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_connection" {
				continue
			}
			projectID := rs.Primary.Attributes["project_id"]
			instanceName := rs.Primary.Attributes["instance_name"]
			connectionName := rs.Primary.Attributes["connection_name"]
			_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, instanceName, connectionName).Execute()
			if err != nil {
				return fmt.Errorf("stream connection (%s:%s:%s) does not exist", projectID, instanceName, connectionName)
			}
		}
		return nil
	}
}

func CheckDestroyStreamConnection(state *terraform.State) error {
	if instanceDestroyedErr := acc.CheckDestroyStreamInstance(state); instanceDestroyedErr != nil {
		return instanceDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_connection" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		instanceName := rs.Primary.Attributes["instance_name"]
		connectionName := rs.Primary.Attributes["connection_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, instanceName, connectionName).Execute()
		if err == nil {
			return fmt.Errorf("stream connection (%s:%s:%s) still exists", projectID, instanceName, connectionName)
		}
	}
	return nil
}
