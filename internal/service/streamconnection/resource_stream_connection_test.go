package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

//go:embed testdata/dummy-ca.pem
var DummyCACert string
var (
	networkingNameVPC    = "vpc-096df0a10902ad9a9"
	networkingTypeVPC    = "VPC"
	networkingTypePublic = "PUBLIC"
	kafkaNetworkingVPC   = fmt.Sprintf(`networking = {
			access = {
				type = %[1]q
				name = %[2]q
			}
		}`, networkingTypeVPC, networkingNameVPC)
	kafkaNetworkingPublic = fmt.Sprintf(`networking = {
			access = {
				type = %[1]q
			}
		}`, networkingTypePublic)
)

func TestAccStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	testCase := testCaseKafkaPlaintext(t)
	resource.ParallelTest(t, *testCase)
}

func testCaseKafkaPlaintext(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", "", false),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false, true),
			},
			{
				Config: kafkaStreamConnectionConfig(projectID, instanceName, "user2", "otherpassword", "localhost:9093", "latest", kafkaNetworkingPublic, false),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user2", "otherpassword", "localhost:9093", "latest", false, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication.password"},
			},
		},
	}
}

func TestAccStreamRSStreamConnection_kafkaNetworkingVPC(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingPublic, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
					resource.TestCheckResourceAttr(resourceName, "networking.access.type", networkingTypePublic),
				),
			},
			{
				Config: kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingVPC, true),
				Check: resource.ComposeAggregateTestCheckFunc(kafkaStreamConnectionAttributeChecks(
					resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
					resource.TestCheckResourceAttr(resourceName, "networking.access.type", networkingTypeVPC),
					resource.TestCheckResourceAttr(resourceName, "networking.access.name", networkingNameVPC)),
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
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: kafkaStreamConnectionConfig(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true),
				Check:  kafkaStreamConnectionAttributeChecks(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, true),
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
		resourceName           = "mongodbatlas_stream_connection.test"
		projectID, clusterName = acc.ClusterNameExecution(t)
		instanceName           = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: clusterStreamConnectionConfig(projectID, instanceName, clusterName),
				Check:  clusterStreamConnectionAttributeChecks(resourceName, clusterName),
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

func TestAccStreamRSStreamConnection_sample(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_connection.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
		sampleName   = "sample_stream_solar"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: sampleStreamConnectionConfig(projectID, instanceName, sampleName),
				Check:  sampleStreamConnectionAttributeChecks(resourceName, instanceName, sampleName),
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

func kafkaStreamConnectionConfig(projectID, instanceName, username, password, bootstrapServers, configValue, networkingConfig string, useSSL bool) string {
	projectAndStreamInstanceConfig := acc.StreamInstanceConfig(projectID, instanceName, "VIRGINIA_USA", "AWS")
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
		    project_id = mongodbatlas_stream_instance.test.project_id
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
			%[7]s
		}
	`, projectAndStreamInstanceConfig, username, password, bootstrapServers, configValue, networkingConfig, securityConfig)
}

func sampleStreamConnectionConfig(projectID, instanceName, sampleName string) string {
	streamInstanceConfig := acc.StreamInstanceConfig(projectID, instanceName, "VIRGINIA_USA", "AWS")

	return fmt.Sprintf(`
		%[1]s
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = %[2]q
		 	type = "Sample"
		}
	`, streamInstanceConfig, sampleName)
}

func sampleStreamConnectionAttributeChecks(
	resourceName, instanceName, sampleName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", sampleName),
		resource.TestCheckResourceAttr(resourceName, "type", "Sample"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func kafkaStreamConnectionAttributeChecks(
	resourceName, instanceName, username, password, bootstrapServers, configValue string, usesSSL, checkPassword bool) resource.TestCheckFunc {
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
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func clusterStreamConnectionConfig(projectID, instanceName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]q
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
		 	cluster_name = %[3]q
			db_role_to_execute = {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
	`, projectID, instanceName, clusterName)
}

func clusterStreamConnectionAttributeChecks(resourceName, clusterName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Cluster"),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.role", "atlasAdmin"),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.type", "BUILT_IN"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
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
