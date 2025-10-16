package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

const (
	dataSourceConfig = `
data "mongodbatlas_stream_connection" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
		connection_name = mongodbatlas_stream_connection.test.connection_name
}
`

	dataSourcePluralConfig = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
}
`
	dataSourcePluralConfigWithPage = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
		page_num = 2 # no specific reason for 2, just to test pagination
		items_per_page = 1
	}
	`
)

const (
	dataSourceConfigMigration = `
data "mongodbatlas_stream_connection" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		connection_name = mongodbatlas_stream_connection.test.connection_name
}
`

	dataSourcePluralConfigMigration = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
}
`
	dataSourcePluralConfigWithPageMigration = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		page_num = 2 # no specific reason for 2, just to test pagination
		items_per_page = 1
	}
	`
)

var (
	dataSourcesConfig         = dataSourceConfig + dataSourcePluralConfig
	dataSourcesWithPagination = dataSourceConfig + dataSourcePluralConfigWithPage
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

	resourceName         = "mongodbatlas_stream_connection.test"
	dataSourceName       = "data.mongodbatlas_stream_connection.test"
	pluralDataSourceName = "data.mongodbatlas_stream_connections.test"
)

func TestAccStreamRSStreamConnection_kafkaPlaintext(t *testing.T) {
	testCase := testCaseKafkaPlaintext(t)
	resource.ParallelTest(t, *testCase)
}

func testCaseKafkaPlaintext(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = "kafka-conn-plaintext"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfig + configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092,localhost:9092", "earliest", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributesAcceptance(resourceName, instanceName, connectionName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, true),
					checkKafkaAttributesAcceptance(dataSourceName, instanceName, connectionName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecksAcceptance(pluralDataSourceName, nil, nil),
				),
			},
			{
				Config: dataSourcesWithPagination + configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("PLAIN", "user2", "otherpassword", "", "", "", "", "", ""), "localhost:9093", "latest", kafkaNetworkingPublic, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributesAcceptance(resourceName, instanceName, connectionName, "user2", "otherpassword", "localhost:9093", "latest", networkingTypePublic, false, true),
					checkKafkaAttributesAcceptance(dataSourceName, instanceName, connectionName, "user2", "otherpassword", "localhost:9093", "latest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecksAcceptance(pluralDataSourceName, conversion.Pointer(2), conversion.Pointer(1)),
				),
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

func testCaseKafkaPlaintextMigration(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = "kafka-conn-plaintext"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigMigration + dataSourcePluralConfigMigration + configureKafkaMigration(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092,localhost:9092", "earliest", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributesMigration(resourceName, instanceName, connectionName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, true),
					checkKafkaAttributesMigration(dataSourceName, instanceName, connectionName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecksMigration(pluralDataSourceName, nil, nil),
				),
			},
			{
				Config: dataSourceConfigMigration + dataSourcePluralConfigWithPageMigration + configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("PLAIN", "user2", "otherpassword", "", "", "", "", "", ""), "localhost:9093", "latest", kafkaNetworkingPublic, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributesMigration(resourceName, instanceName, connectionName, "user2", "otherpassword", "localhost:9093", "latest", networkingTypePublic, false, true),
					checkKafkaAttributesMigration(dataSourceName, instanceName, connectionName, "user2", "otherpassword", "localhost:9093", "latest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecksMigration(pluralDataSourceName, conversion.Pointer(2), conversion.Pointer(1)),
				),
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

func TestAccStreamRSStreamConnection_kafkaOAuthBearer(t *testing.T) {
	t.Helper()
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = "kafka-conn-oauthbearer"
	)

	testCase := &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfig + configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("OAUTHBEARER", "", "", tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtentions, method), "localhost:9092,localhost:9092", "earliest", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaOAuthAttributes(resourceName, connectionName, tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtentions, method, "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, true),
					checkKafkaOAuthAttributes(dataSourceName, connectionName, tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtentions, method, "localhost:9092,localhost:9092", "earliest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecks(pluralDataSourceName, nil, nil),
				),
			},
			{
				Config: dataSourcesWithPagination + configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName, getKafkaAuthenticationConfig("OAUTHBEARER", "", "", tokenEndpointURL, "clientId2", "clientSecret", scope, saslOauthbearerExtentions, method), "localhost:9093", "latest", kafkaNetworkingPublic, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaOAuthAttributes(resourceName, connectionName, tokenEndpointURL, "clientId2", "clientSecret", scope, saslOauthbearerExtentions, method, "localhost:9093", "latest", networkingTypePublic, false, true),
					checkKafkaOAuthAttributes(dataSourceName, connectionName, tokenEndpointURL, "clientId2", "clientSecret", scope, saslOauthbearerExtentions, method, "localhost:9093", "latest", networkingTypePublic, false, false),
					streamConnectionsAttributeChecks(pluralDataSourceName, conversion.Pointer(2), conversion.Pointer(1)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication.client_secret"},
			},
		},
	}
	resource.ParallelTest(t, *testCase)
}

func TestAccStreamRSStreamConnection_kafkaNetworkingVPC(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		vpcID                   = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock            = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID            = os.Getenv("AWS_ACCOUNT_ID")
		peerRegion              = os.Getenv("AWS_REGION")
		containerRegion         = conversion.AWSRegionToMongoDBRegion(peerRegion)
		providerName            = "AWS"
		networkPeeringConfig    = configNetworkPeeringAWS(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, containerRegion, peerRegion)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckPeeringEnvAWS(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: networkPeeringConfig + configureKafka("mongodbatlas_network_peering.test.project_id", instanceName, "kafka-conn-vpc", getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092", "earliest", kafkaNetworkingVPC, true),
				Check:  checkKafkaAttributesAcceptance(resourceName, instanceName, "kafka-conn-vpc", "user", "rawpassword", "localhost:9092", "earliest", networkingTypeVPC, true, true),
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
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		vpcID                   = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock            = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID            = os.Getenv("AWS_ACCOUNT_ID")
		peerRegion              = os.Getenv("AWS_REGION")
		containerRegion         = conversion.AWSRegionToMongoDBRegion(peerRegion)
		providerName            = "AWS"
		networkPeeringConfig    = configNetworkPeeringAWS(projectID, providerName, vpcID, awsAccountID, vpcCIDRBlock, containerRegion, peerRegion)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf("%s\n%s", configureKafka(fmt.Sprintf("%q", projectID), instanceName, "kafka-conn-ssl", getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092", "earliest", kafkaNetworkingPublic, true), dataSourceConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkKafkaAttributesAcceptance(resourceName, instanceName, "kafka-conn-ssl", "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true, true),
					checkKafkaAttributesAcceptance(dataSourceName, instanceName, "kafka-conn-ssl", "user", "rawpassword", "localhost:9092", "earliest", networkingTypePublic, true, false),
				),
			},
			// cannot change networking access type once set
			{
				Config:      networkPeeringConfig + configureKafka("mongodbatlas_network_peering.test.project_id", instanceName, "kafka-conn-ssl", getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092", "earliest", kafkaNetworkingVPC, true),
				ExpectError: regexp.MustCompile("STREAM_NETWORKING_ACCESS_TYPE_CANNOT_BE_MODIFIED"),
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
		_, instanceName        = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName         = "conn-cluster"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfig + configureCluster(projectID, instanceName, connectionName, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkClusterAttributesAcceptance(resourceName, clusterName),
					checkClusterAttributesAcceptance(dataSourceName, clusterName),
				),
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

func testCaseClusterMigration(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, clusterName = acc.ClusterNameExecution(t, false)
		_, instanceName        = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName         = "conn-cluster-mig"
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigMigration + dataSourcePluralConfigMigration + configureClusterMigration(projectID, instanceName, connectionName, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkClusterAttributesMigration(resourceName, clusterName),
					checkClusterAttributesMigration(dataSourceName, clusterName),
				),
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
		projectID, _ = acc.ProjectIDExecutionWithStreamInstance(t)
		instanceName = acc.RandomStreamInstanceName() // The execution stream instance use sample stream, so we need to create this in a different instance
		sampleName   = "sample_stream_solar"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourcesConfig + configureSampleStream(projectID, instanceName, sampleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkSampleStreamAttributes(resourceName, instanceName, sampleName),
					checkSampleStreamAttributes(dataSourceName, instanceName, sampleName),
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

func TestAccStreamStreamConnection_https(t *testing.T) {
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		url                      = "https://example.com"
		updatedURL               = "https://example2.com"
		headerStr                = `headers = {
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
				Config: configureHTTPS(projectID, workspaceName, url, headerStr),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(workspaceName, url),
					resource.TestCheckResourceAttr(resourceName, "headers.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "headers.Authorization", "Bearer token"),
					resource.TestCheckResourceAttr(resourceName, "headers.key1", "value1"),
				),
			},
			{
				Config: configureHTTPS(projectID, workspaceName, updatedURL, updatedHeaderStr),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(workspaceName, updatedURL),
					resource.TestCheckResourceAttr(resourceName, "headers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "headers.updatedKey", "updatedValue"),
				),
			},
			{
				Config: configureHTTPS(projectID, workspaceName, updatedURL, emptyHeaders),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHTTPSAttributes(workspaceName, updatedURL),
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
		projectID, instanceName    = acc.ProjectIDExecutionWithStreamInstance(t)
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
				`, privatelinkConfig, configureKafka(fmt.Sprintf("%q", projectID), instanceName, "kafka-conn-privatelink", getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092", "earliest", kafkaNetworkingPrivatelink, true)),
				Check: checkKafkaAttributesAcceptance(resourceName, instanceName, "kafka-conn-privatelink", "user", "rawpassword", "localhost:9092", "earliest", networkingTypePrivatelink, true, true),
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
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		awsIAMRoleName          = acc.RandomIAMRole()
		connectionName          = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureAWSLambda(projectID, instanceName, connectionName, awsIAMRoleName),
				Check:  checkAWSLambdaAttributes(resourceName, instanceName, connectionName),
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

func TestAccStreamRSStreamConnection_instanceName(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureKafkaWithInstanceName(projectID, instanceName, connectionName, "user", "password", "localhost:9092"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamConnectionExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
					resource.TestCheckNoResourceAttr(resourceName, "workspace_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// When the new resource is imported, it will contain workspace_name instead of instance_name. This is expected so we will ignore it
				ImportStateVerifyIgnore: []string{"authentication.password", "instance_name", "workspace_name"},
			},
		},
	})
}

func TestAccStreamRSStreamConnection_conflictingFields(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = "conflict-test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config:      configureKafkaWithInstanceAndWorkspaceName(projectID, instanceName, connectionName, "user", "password", "localhost:9092"),
				ExpectError: regexp.MustCompile("Attribute \"workspace_name\" cannot be specified when \"instance_name\" is\n.*specified"),
			},
		},
	})
}

func getKafkaAuthenticationConfig(mechanism, username, password, tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtensions, method string) string {
	if mechanism == "PLAIN" {
		return fmt.Sprintf(`authentication = {
			mechanism = %[1]q
			username = %[2]q
			password = %[3]q
		}`, mechanism, username, password)
	}
	return fmt.Sprintf(`authentication = {
			mechanism = %[1]q
			method = %[2]q
			token_endpoint_url = %[3]q
			client_id = %[4]q
			client_secret = %[5]q
			scope = %[6]q
			sasl_oauthbearer_extensions = %[7]q
		}`, mechanism, method, tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtensions)
}

func configureKafka(projectRef, workspaceName, connectionName, authenticationConfig, bootstrapServers, configValue, networkingConfig string, useSSL bool) string {
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
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]s
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Kafka"
		 	%[4]s
		    bootstrap_servers = %[5]q
		    config = {
		    	"auto.offset.reset": %[6]q
		    }
		    %[7]s
			%[8]s
		}
	`, projectRef, workspaceName, connectionName, authenticationConfig, bootstrapServers, configValue, networkingConfig, securityConfig)
}

// configureKafkaForMigration uses instance_name for compatibility with older provider versions
func configureKafkaMigration(projectRef, instanceName, connectionName, authenticationConfig, bootstrapServers, configValue, networkingConfig string, useSSL bool) string {
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
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]s
			instance_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Kafka"
		 	%[4]s
		    bootstrap_servers = %[5]q
		    config = {
		    	"auto.offset.reset": %[6]q
		    }
		    %[7]s
			%[8]s
		}
	`, projectRef, instanceName, connectionName, authenticationConfig, bootstrapServers, configValue, networkingConfig, securityConfig)
}

func configureSampleStream(projectID, workspaceName, sampleName string) string {
	streamInstanceConfig := acc.StreamInstanceConfig(projectID, workspaceName, "VIRGINIA_USA", "AWS")

	return fmt.Sprintf(`
		%[1]s
		
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = mongodbatlas_stream_instance.test.project_id
			workspace_name = mongodbatlas_stream_instance.test.instance_name
		 	connection_name = %[2]q
		 	type = "Sample"
		}
	`, streamInstanceConfig, sampleName)
}

// configureKafkaWithInstanceName tests that the deprecated isntance_name field is still functional
func configureKafkaWithInstanceName(projectID, instanceName, connectionName, username, password, bootstrapServers string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			instance_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Kafka"
		 	authentication = {
		    	mechanism = "PLAIN"
		    	username = %[4]q
		    	password = %[5]q
		    }
		    bootstrap_servers = %[6]q
		    config = {
		    	"auto.offset.reset": "earliest"
		    }
		    security = {
				protocol = "SASL_PLAINTEXT"
			}
		}
	`, projectID, instanceName, connectionName, username, password, bootstrapServers)
}

func configureKafkaWithInstanceAndWorkspaceName(projectID, instanceName, connectionName, username, password, bootstrapServers string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			instance_name = %[2]q
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Kafka"
		 	authentication = {
		    	mechanism = "PLAIN"
		    	username = %[4]q
		    	password = %[5]q
		    }
		    bootstrap_servers = %[6]q
		    config = {
		    	"auto.offset.reset": "earliest"
		    }
		    security = {
				protocol = "SASL_PLAINTEXT"
			}
		}
	`, projectID, instanceName, connectionName, username, password, bootstrapServers)
}

func checkSampleStreamAttributes(
	resourceName, workspaceName, sampleName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", sampleName),
		resource.TestCheckResourceAttr(resourceName, "type", "Sample"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func checkHTTPSAttributes(workspaceName, url string) resource.TestCheckFunc {
	setChecks := []string{"project_id"}
	mapChecks := map[string]string{
		"workspace_name":  workspaceName,
		"connection_name": "ConnectionNameHttps",
		"type":            "Https",
		"url":             url,
	}
	extra := []resource.TestCheckFunc{checkStreamConnectionExists()}
	return acc.CheckRSAndDS(resourceName, conversion.StringPtr(dataSourceName), nil, setChecks, mapChecks, extra...)
}

func checkKafkaAttributes(
	resourceName, connectionName, username, password, bootstrapServers, configValue, networkingType string, usesSSL, checkPassword bool) resource.TestCheckFunc {
	authAttrs := map[string]string{
		"authentication.mechanism": "PLAIN",
		"authentication.username":  username,
	}
	if checkPassword {
		authAttrs["authentication.password"] = password
	}
	return checkKafkaConnectionAttributes(resourceName, connectionName, bootstrapServers, configValue, networkingType, usesSSL, authAttrs)
}

func checkKafkaOAuthAttributes(
	resourceName, connectionName, tokenEndpointURL, clientID, clientSecret, scope, saslOauthbearerExtensions, method, bootstrapServers, configValue, networkingType string, usesSSL, checkClientSecret bool) resource.TestCheckFunc {
	authAttrs := map[string]string{
		"authentication.mechanism":                   "OAUTHBEARER",
		"authentication.method":                      method,
		"authentication.token_endpoint_url":          tokenEndpointURL,
		"authentication.client_id":                   clientID,
		"authentication.scope":                       scope,
		"authentication.sasl_oauthbearer_extensions": saslOauthbearerExtensions,
	}
	if checkClientSecret {
		authAttrs["authentication.client_secret"] = clientSecret
	}
	return checkKafkaConnectionAttributes(resourceName, connectionName, bootstrapServers, configValue, networkingType, usesSSL, authAttrs)
}

func checkKafkaConnectionAttributes(resourceName, connectionName, bootstrapServers, configValue, networkingType string, usesSSL bool, authAttrs map[string]string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
		resource.TestCheckResourceAttr(resourceName, "bootstrap_servers", bootstrapServers),
		resource.TestCheckResourceAttr(resourceName, "config.auto.offset.reset", configValue),
	}

	resourceChecks = acc.AddAttrChecks(resourceName, resourceChecks, authAttrs)

	if mig.IsProviderVersionAtLeast("1.25.0") {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "networking.access.type", networkingType))
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

func checkKafkaAttributesMigration(
	resourceName, instanceName, connectionName, username, password, bootstrapServers, configValue, networkingType string, usesSSL, checkPassword bool) resource.TestCheckFunc {
	commonTests := checkKafkaAttributes(resourceName, connectionName, username, password, bootstrapServers, configValue, networkingType, usesSSL, checkPassword)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName))
}

func checkKafkaAttributesAcceptance(
	resourceName, workspaceName, connectionName, username, password, bootstrapServers, configValue, networkingType string, usesSSL, checkPassword bool) resource.TestCheckFunc {
	commonTests := checkKafkaAttributes(resourceName, connectionName, username, password, bootstrapServers, configValue, networkingType, usesSSL, checkPassword)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName))
}

func configureCluster(projectID, workspaceName, connectionName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Cluster"
		 	cluster_name = %[4]q
			db_role_to_execute = {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
	`, projectID, workspaceName, connectionName, clusterName)
}

// configureClusterMigration uses instance_name for compatibility with older provider versions
func configureClusterMigration(projectID, instanceName, connectionName, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			instance_name = %[2]q
		 	connection_name = %[3]q
		 	type = "Cluster"
		 	cluster_name = %[4]q
			db_role_to_execute = {
				role = "atlasAdmin"
				type = "BUILT_IN"
			}
		}
	`, projectID, instanceName, connectionName, clusterName)
}

func configureHTTPS(projectID, workspaceName, url, headers string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			connection_name = "ConnectionNameHttps"
			type = "Https"
			url = %[3]q
			%[4]s	
		}

		data "mongodbatlas_stream_connection" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			connection_name = mongodbatlas_stream_connection.test.connection_name
		}
	`, projectID, workspaceName, url, headers)
}

func checkClusterAttributes(resourceName, clusterName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_name"),
		resource.TestCheckResourceAttr(resourceName, "type", "Cluster"),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.role", "atlasAdmin"),
		resource.TestCheckResourceAttr(resourceName, "db_role_to_execute.type", "BUILT_IN"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func checkClusterAttributesAcceptance(resourceName, clusterName string) resource.TestCheckFunc {
	commonTests := checkClusterAttributes(resourceName, clusterName)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttrSet(resourceName, "workspace_name"))
}

func checkClusterAttributesMigration(resourceName, clusterName string) resource.TestCheckFunc {
	commonTests := checkClusterAttributes(resourceName, clusterName)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttrSet(resourceName, "instance_name"))
}

func checkStreamConnectionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func checkStreamConnectionExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_stream_connection" {
				continue
			}
			projectID := rs.Primary.Attributes["project_id"]
			workspaceName := rs.Primary.Attributes["workspace_name"]
			if workspaceName == "" {
				workspaceName = rs.Primary.Attributes["instance_name"]
			}
			connectionName := rs.Primary.Attributes["connection_name"]
			_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, workspaceName, connectionName).Execute()
			if err != nil {
				return fmt.Errorf("stream connection (%s:%s:%s) does not exist", projectID, workspaceName, connectionName)
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
		workspaceName := rs.Primary.Attributes["workspace_name"]
		if workspaceName == "" {
			workspaceName = rs.Primary.Attributes["instance_name"]
		}
		connectionName := rs.Primary.Attributes["connection_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamConnection(context.Background(), projectID, workspaceName, connectionName).Execute()
		if err == nil {
			return fmt.Errorf("stream connection (%s:%s:%s) still exists", projectID, workspaceName, connectionName)
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

func configureAWSLambda(projectID, instanceName, connectionName, awsIamRoleName string) string {
	config := fmt.Sprintf(`
		resource "aws_iam_role" "test_role" {
			name = %[4]q
	  
			assume_role_policy = jsonencode({
				"Version" : "2012-10-17",
				"Statement" : [
					{
						"Effect" : "Allow",
						"Principal" : {
							"AWS" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
						},
						"Action" : "sts:AssumeRole",
						"Condition" : {
							"StringEquals" : {
								"sts:ExternalId" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
							}
						}
					}
				]
			})
		}

		resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
			project_id    = %[1]q
			provider_name = "AWS"
		}
	  
		resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	  
			aws {
				iam_assumed_role_arn = aws_iam_role.test_role.arn
			}
		}

		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "AWSLambda"
            aws = {
				role_arn = mongodbatlas_cloud_provider_access_authorization.auth_role.aws[0].iam_assumed_role_arn
			}
		}
	`, projectID, instanceName, connectionName, awsIamRoleName)
	return config
}

func checkAWSLambdaAttributes(resourceName, workspaceName, connectionName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "AWSLambda"),
		resource.TestCheckResourceAttrSet(resourceName, "aws.role_arn"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func streamConnectionsAttributeChecks(resourceName string, pageNum, itemsPerPage *int) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttrSet(resourceName, "total_count"),
		resource.TestCheckResourceAttrSet(resourceName, "results.#"),
	}
	if pageNum != nil {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "page_num", fmt.Sprint(*pageNum)))
	}
	if itemsPerPage != nil {
		resourceChecks = append(resourceChecks, resource.TestCheckResourceAttr(resourceName, "items_per_page", fmt.Sprint(*itemsPerPage)))
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func streamConnectionsAttributeChecksAcceptance(resourceName string, pageNum, itemsPerPage *int) resource.TestCheckFunc {
	commonTests := streamConnectionsAttributeChecks(resourceName, pageNum, itemsPerPage)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttrSet(resourceName, "workspace_name"))
}

func streamConnectionsAttributeChecksMigration(resourceName string, pageNum, itemsPerPage *int) resource.TestCheckFunc {
	commonTests := streamConnectionsAttributeChecks(resourceName, pageNum, itemsPerPage)
	return resource.ComposeAggregateTestCheckFunc(commonTests, resource.TestCheckResourceAttrSet(resourceName, "instance_name"))
}
