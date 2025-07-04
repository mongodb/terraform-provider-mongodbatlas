package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

var (
	//go:embed testdata/dummy-ca.pem
	DummyCACert               string
	networkingTypeVPC         = "VPC"
	networkingTypePublic      = "PUBLIC"
	networkingTypePrivatelink = "PRIVATE_LINK"
	kafkaNetworkingVPC        = fmt.Sprintf(`networking = {
			access = {
				type = %[1]q
			}
		}`, networkingTypeVPC)
	kafkaNetworkingPublic = fmt.Sprintf(`networking = {
			access = {
				type = %[1]q
			}
		}`, networkingTypePublic)

	resourceName   = "mongodbatlas_stream_connection.test"
	dataSourceName = "data.mongodbatlas_stream_connection.test"
)

func TestAccStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	testCase := testCaseKafkaPlaintext(t)
	resource.ParallelTest(t, *testCase)
}

func testCaseKafkaPlaintext(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureKafka(projectID, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", "", false),
				Check:  checkKafkaAttributes(resourceName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, true),
			},
			{
				Config: configureKafka(projectID, instanceName, "user2", "otherpassword", "localhost:9093", "latest", kafkaNetworkingPublic, false),
				Check:  checkKafkaAttributes(resourceName, instanceName, "user2", "otherpassword", "localhost:9093", "latest", networkingTypePublic, false, true),
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
		projectID            = acc.ProjectIDExecution(t)
		instanceName         = acc.RandomName()
		vpcID                = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock         = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID         = os.Getenv("AWS_ACCOUNT_ID")
		containerRegion      = os.Getenv("AWS_REGION")
		peerRegion           = conversion.MongoDBRegionToAWSRegion(containerRegion)
		providerName         = "AWS"
		networkPeeringConfig = configNetworkPeeringAWS(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, containerRegion, peerRegion)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckPeeringEnvAWS(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: networkPeeringConfig + configureKafka(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingPublic, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributes(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true, true),
				),
			},
			{
				Config: networkPeeringConfig + configureKafka(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingVPC, true),
				Check:  checkKafkaAttributes(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypeVPC, true, true),
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
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureKafka(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingPublic, true),
				Check:  checkKafkaAttributes(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true, true),
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
	testCase := testCaseCluster(t)
	resource.ParallelTest(t, *testCase)
}

func testCaseCluster(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, false)
		instanceName           = acc.RandomName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureCluster(projectID, instanceName, clusterName),
				Check:  checkClusterAttributes(resourceName, clusterName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccStreamRSStreamConnection_sample(t *testing.T) {
	var (
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
				Config: configureSampleStream(projectID, instanceName, sampleName),
				Check:  checkSampleStreamAttributes(resourceName, instanceName, sampleName),
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

func TestAccStreamStreamConnection_https(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
		url          = "https://example.com"
		updatedURL   = "https://example2.com"
		headerStr    = `headers = {
			Authorization = "Bearer token"
			key1 = "value1"
		}`
		updatedHeaderStr = `headers = {
			updatedKey = "updatedValue"
		}`
		emptyHeaders string
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureHTTPS(projectID, instanceName, url, headerStr),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(instanceName, url),
					resource.TestCheckResourceAttr(resourceName, "headers.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "headers.Authorization", "Bearer token"),
					resource.TestCheckResourceAttr(resourceName, "headers.key1", "value1"),
				),
			},
			{
				Config: configureHTTPS(projectID, instanceName, updatedURL, updatedHeaderStr),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(instanceName, updatedURL),
					resource.TestCheckResourceAttr(resourceName, "headers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "headers.updatedKey", "updatedValue"),
				),
			},
			{
				Config: configureHTTPS(projectID, instanceName, updatedURL, emptyHeaders),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(instanceName, updatedURL),
					resource.TestCheckResourceAttr(resourceName, "headers.%", "0"),
				),
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

func TestAccStreamPrivatelinkEndpoint_streamConnection(t *testing.T) {
	acc.SkipTestForCI(t) // requires Confluent Cloud resources
	var (
		projectID                  = acc.ProjectIDExecution(t)
		instanceName               = acc.RandomName()
		vendor                     = "CONFLUENT"
		provider                   = "AWS"
		region                     = "us-east-1"
		awsAccountID               = os.Getenv("AWS_ACCOUNT_ID")
		networkID                  = os.Getenv("CONFLUENT_CLOUD_NETWORK_ID")
		privatelinkAccessID        = os.Getenv("CONFLUENT_CLOUD_PRIVATELINK_ACCESS_ID")
		privatelinkConfig          = acc.GetCompleteConfluentConfig(true, true, projectID, provider, region, vendor, awsAccountID, networkID, privatelinkAccessID)
		kafkaNetworkingPrivatelink = fmt.Sprintf(`networking = {
			access = {
				type = %[1]q
				connection_id = mongodbatlas_stream_privatelink_endpoint.test.id
			}
		}`, networkingTypePrivatelink)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					%[1]s
					%[2]s
				`, privatelinkConfig, configureKafka(projectID, instanceName, "user", "rawpassword", "localhost:9092", "earliest", kafkaNetworkingPrivatelink, true)),
				Check: checkKafkaAttributes(resourceName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", networkingTypePrivatelink, true, true),
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

func TestAccStreamRSStreamConnection_AWSLambda(t *testing.T) {
	var (
		projectID      = os.Getenv("MONGODB_ATLAS_ASP_PROJECT_EAR_PE_ID") // test-acc-tf-p-keep-ear-AWS-private-endpoint project has aws integration
		instanceName   = acc.RandomName()
		connectionName = acc.RandomName()
		roleArn        = os.Getenv("MONGODB_ATLAS_ASP_PROJECT_AWS_ROLE_ARN")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureAWSLambda(projectID, instanceName, connectionName, roleArn),
				Check:  checkAWSLambdaAttributes(resourceName, instanceName, connectionName, roleArn),
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

func configureKafka(projectID, instanceName, username, password, bootstrapServers, configValue, networkingConfig string, useSSL bool) string {
	projectAndStreamInstanceConfig := acc.StreamInstanceConfig(projectID, instanceName, "VIRGINIA_USA", "AWS")
	securityConfig := `
		security = {
			protocol = "SASL_PLAINTEXT"
		}`

	if useSSL {
		securityConfig = fmt.Sprintf(`
		security = {
		    broker_public_certificate = %q
		    protocol = "SASL_SSL"
		}`, DummyCACert)
	}
	return fmt.Sprintf(`
		%[1]s
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = mongodbatlas_stream_instance.test.instance_name
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

func configureSampleStream(projectID, instanceName, sampleName string) string {
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

func checkSampleStreamAttributes(
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

func checkHTTPSAttributes(instanceName, url string) resource.TestCheckFunc {
	setChecks := []string{"project_id"}
	mapChecks := map[string]string{
		"instance_name":   instanceName,
		"connection_name": "ConnectionNameHttps",
		"type":            "Https",
		"url":             url,
	}
	extra := []resource.TestCheckFunc{checkStreamConnectionExists()}
	return acc.CheckRSAndDS(resourceName, conversion.StringPtr(dataSourceName), nil, setChecks, mapChecks, extra...)
}

func checkKafkaAttributes(
	resourceName, instanceName, username, password, bootstrapServers, configValue, networkingType string, usesSSL, checkPassword bool) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
		resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
		resource.TestCheckResourceAttr(resourceName, "authentication.mechanism", "PLAIN"),
		resource.TestCheckResourceAttr(resourceName, "authentication.username", username),
		resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", bootstrapServers),
		resource.TestCheckResourceAttr(resourceName, "config.auto.offset.reset", configValue),
	}
	if mig.IsProviderVersionAtLeast("1.25.0") {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "networking.access.type", networkingType))
	}
	if checkPassword {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "authentication.password", password))
	}
	if !usesSSL {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "security.protocol", "SASL_PLAINTEXT"))
	} else {
		resourceChecks = append(resourceChecks,
			resource.TestCheckResourceAttr(resourceName, "security.protocol", "SASL_SSL"),
			resource.TestCheckResourceAttrSet(resourceName, "security.broker_public_certificate"),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func configureCluster(projectID, instanceName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
			stream_config = {
				tier = "SP10"
			}
		}
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = "ConnectionNameCluster"
		 	type = "Cluster"
		 	cluster_name = %[3]q
			db_role_to_execute = {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
	`, projectID, instanceName, clusterName)
}

func configureHTTPS(projectID, instanceName, url, headers string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "test" {
			project_id = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
			stream_config = {
				tier = "SP10"
			}
		}
			
		resource "mongodbatlas_stream_connection" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
			connection_name = "ConnectionNameHttps"
			type = "Https"
			url = %[3]q
			%[4]s	
		}

		data "mongodbatlas_stream_connection" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
			connection_name = mongodbatlas_stream_connection.test.connection_name
		}
	`, projectID, instanceName, url, headers)
}

func checkClusterAttributes(resourceName, clusterName string) resource.TestCheckFunc {
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

func configNetworkPeeringAWS(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegionContainer, awsRegionPeer string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_network_container" "test" {
		project_id   		 = %[1]q
		atlas_cidr_block  	 = "192.168.208.0/21"
		provider_name		 = %[2]q
		region_name			 = %[6]q
	}

	resource "mongodbatlas_network_peering" "test" {
		accepter_region_name	= %[7]q
		project_id    			= %[1]q
		container_id           	= mongodbatlas_network_container.test.id
		provider_name           = %[2]q
		route_table_cidr_block  = %[5]q
		vpc_id					= %[3]q
		aws_account_id	        = %[4]q
	}
`, projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegionContainer, awsRegionPeer)
}

func configureAWSLambda(projectID, instanceName, connectionName, roleArn string) string {
	streamInstanceConfig := acc.StreamInstanceConfig(projectID, instanceName, "VIRGINIA_USA", "AWS")

	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = %[2]q
		 	type = "AWSLambda"
            aws = {
				role_arn = %[3]q
			}
		}
	`, streamInstanceConfig, connectionName, roleArn)
}

func checkAWSLambdaAttributes(
	resourceName, instanceName, connectionName, roleArn string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "AWSLambda"),
		resource.TestCheckResourceAttr(resourceName, "aws.role_arn", roleArn),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}
