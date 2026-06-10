package streamconnection_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

// azureServicePrincipalMu serializes Azure Blob Storage tests that share the same
// service principal ID, preventing DUPLICATE_AZURE_SERVICE_PRINCIPAL errors
var azureServicePrincipalMu sync.Mutex

const (
	dataSourceConfig = `
data "mongodbatlas_stream_connection" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
		connection_name = mongodbatlas_stream_connection.test.connection_name
		depends_on = [mongodbatlas_stream_connection.test]
}
`

	dataSourcePluralConfig = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
		depends_on = [mongodbatlas_stream_connection.test]
}
`
	dataSourcePluralConfigWithPage = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		workspace_name = mongodbatlas_stream_connection.test.workspace_name
		page_num = 2 # no specific reason for 2, just to test pagination
		items_per_page = 1
		depends_on = [mongodbatlas_stream_connection.test]
}
`

	dataSourceConfigMigration = `
data "mongodbatlas_stream_connection" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		connection_name = mongodbatlas_stream_connection.test.connection_name
		depends_on = [mongodbatlas_stream_connection.test]
}
`

	dataSourcePluralConfigMigration = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		depends_on = [mongodbatlas_stream_connection.test]
}
`

	dataSourcePluralConfigWithPageMigration = `
data "mongodbatlas_stream_connections" "test" {
		project_id = mongodbatlas_stream_connection.test.project_id
		instance_name = mongodbatlas_stream_connection.test.instance_name
		page_num = 2 # no specific reason for 2, just to test pagination
		items_per_page = 1
		depends_on = [mongodbatlas_stream_connection.test]
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
	resource.Test(t, *testCase)
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
					streamConnectionsAttributeChecksAcceptance(pluralDataSourceName, new(2), new(1)),
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
		connectionName          = "kafka-conn-plaintext-mig"
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
					streamConnectionsAttributeChecksMigration(pluralDataSourceName, new(2), new(1)),
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
					streamConnectionsAttributeChecks(pluralDataSourceName, new(2), new(1)),
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
	resource.Test(t, *testCase)
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
	resource.Test(t, resource.TestCase{
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
			{
				Config:      networkPeeringConfig + configureKafka("mongodbatlas_network_peering.test.project_id", instanceName, "kafka-conn-ssl", getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""), "localhost:9092", "earliest", kafkaNetworkingVPC, true),
				ExpectError: regexp.MustCompile("STREAM_NETWORKING_CANNOT_BE_MODIFIED"),
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
	resource.Test(t, *testCase)
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
	resource.Test(t, resource.TestCase{
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
	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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
	resource.Test(t, resource.TestCase{
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

func TestAccStreamRSStreamConnection_GCPPubSub(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = acc.RandomName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureGCPPubSub(projectID, instanceName, connectionName),
				Check:  checkGCPPubSubAttributes(resourceName, instanceName, connectionName),
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

func TestAccStreamRSStreamConnection_GCPPubSubPrivateLink(t *testing.T) {
	acc.SkipTestForCI(t) // requires a GCP cluster in the same region for privatelink provisioning, too slow for CI
	var (
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomStreamInstanceName()
		clusterName    = acc.RandomClusterName()
		connectionName = acc.RandomName()
		region         = "us-east4"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureGCPPubSubPrivateLink(projectID, instanceName, clusterName, connectionName, region),
				Check:  checkGCPPubSubPrivateLinkAttributes(resourceName, instanceName, connectionName),
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

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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

func TestAccStreamRSStreamConnection_SchemaRegistry(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = acc.RandomName()
		schemaRegistryURLs      = []string{"https://schemaregistry.example.com", "https://schemaregistry2.example.com"}
		username                = "user"
		password                = "password"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig + configureSchemaRegistry(projectID, instanceName, connectionName, "CONFLUENT", "USER_INFO", username, password, schemaRegistryURLs),
				Check:  checkSchemaRegistryAttributes(resourceName, instanceName, connectionName, "CONFLUENT", "USER_INFO", username),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamConnectionImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"schema_registry_authentication.password",
				},
			},
		},
	})
}

func configureSchemaRegistry(projectID, workspaceName, connectionName, provider, authType, username, password string, urls []string) string {
	quotedURLs := make([]string, len(urls))
	for i, url := range urls {
		quotedURLs[i] = fmt.Sprintf("%q", url)
	}
	urlsStr := strings.Join(quotedURLs, ", ")
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "SchemaRegistry"
            schema_registry_provider = %[4]q
			schema_registry_urls = [%[5]s]
		 	schema_registry_authentication = {
				type = %[6]q
				username = %[7]q
				password = %[8]q
			}
		}
	`, projectID, workspaceName, connectionName, provider, urlsStr, authType, username, password)
}

func checkSchemaRegistryAttributes(resourceName, workspaceName, connectionName, provider, authType, username string) resource.TestCheckFunc {
	// check map similar to http way of doing it
	setChecks := []string{"project_id"}
	mapChecks := map[string]string{
		"workspace_name":                          workspaceName,
		"connection_name":                         connectionName,
		"type":                                    "SchemaRegistry",
		"schema_registry_provider":                provider,
		"schema_registry_urls.#":                  "2",
		"schema_registry_authentication.type":     authType,
		"schema_registry_authentication.username": username,
	}
	extra := []resource.TestCheckFunc{checkStreamConnectionExists()}
	return acc.CheckRSAndDS(resourceName, conversion.StringPtr(dataSourceName), nil, setChecks, mapChecks, extra...)
}

func TestAccStreamRSStreamConnection_SchemaRegistrySASLInherit(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = acc.RandomName()
		schemaRegistryURLs      = []string{"https://schemaregistry.example.com"}
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig + configureSchemaRegistrySASLInherit(projectID, instanceName, connectionName, "CONFLUENT", schemaRegistryURLs),
				Check:  checkSchemaRegistrySASLInheritAttributes(resourceName, instanceName, connectionName, "CONFLUENT"),
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

func configureSchemaRegistrySASLInherit(projectID, workspaceName, connectionName, provider string, urls []string) string {
	quotedURLs := make([]string, len(urls))
	for i, url := range urls {
		quotedURLs[i] = fmt.Sprintf("%q", url)
	}
	urlsStr := strings.Join(quotedURLs, ", ")
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
		    project_id = %[1]q
			workspace_name = %[2]q
		 	connection_name = %[3]q
		 	type = "SchemaRegistry"
            schema_registry_provider = %[4]q
			schema_registry_urls = [%[5]s]
		 	schema_registry_authentication = {
				type = "SASL_INHERIT"
			}
		}
	`, projectID, workspaceName, connectionName, provider, urlsStr)
}

func checkSchemaRegistrySASLInheritAttributes(resourceName, workspaceName, connectionName, provider string) resource.TestCheckFunc {
	setChecks := []string{"project_id"}
	mapChecks := map[string]string{
		"workspace_name":                      workspaceName,
		"connection_name":                     connectionName,
		"type":                                "SchemaRegistry",
		"schema_registry_provider":            provider,
		"schema_registry_urls.#":              "1",
		"schema_registry_authentication.type": "SASL_INHERIT",
	}
	extra := []resource.TestCheckFunc{checkStreamConnectionExists()}
	return acc.CheckRSAndDS(resourceName, conversion.StringPtr(dataSourceName), nil, setChecks, mapChecks, extra...)
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

func configureGCPPubSub(projectID, instanceName, connectionName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_access_setup" "gcp_setup" {
			project_id    = %[1]q
			provider_name = "GCP"
		}

		resource "mongodbatlas_cloud_provider_access_authorization" "gcp_auth" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.gcp_setup.role_id
		}

		resource "mongodbatlas_stream_connection" "test" {
			project_id      = %[1]q
			workspace_name  = %[2]q
			connection_name = %[3]q
			type            = "GCPPubSub"
			gcp = {
				service_account_id = mongodbatlas_cloud_provider_access_setup.gcp_setup.gcp_config[0].service_account_for_atlas
			}
			depends_on = [mongodbatlas_cloud_provider_access_authorization.gcp_auth]
		}
	`, projectID, instanceName, connectionName)
}

func checkGCPPubSubAttributes(resourceName, workspaceName, connectionName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "GCPPubSub"),
		resource.TestCheckResourceAttrSet(resourceName, "gcp.service_account_id"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func configureGCPPubSubPrivateLink(projectID, instanceName, clusterName, connectionName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[3]q
			cluster_type = "REPLICASET"
			replication_specs = [{
				region_configs = [{
					priority      = 7
					provider_name = "GCP"
					region_name   = "US_EAST_4"
					electable_specs = {
						instance_size = "M10"
						node_count    = 3
					}
				}]
			}]
		}

		resource "mongodbatlas_stream_workspace" "test" {
			project_id     = %[1]q
			workspace_name = %[2]q
			data_process_region = {
				region         = "US_EAST4"
				cloud_provider = "GCP"
			}
		}

		resource "mongodbatlas_cloud_provider_access_setup" "gcp_setup" {
			project_id    = %[1]q
			provider_name = "GCP"
		}

		resource "mongodbatlas_cloud_provider_access_authorization" "gcp_auth" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.gcp_setup.role_id
		}

		resource "mongodbatlas_stream_privatelink_endpoint" "test" {
			project_id    = %[1]q
			provider_name = "GCP"
			vendor        = "PUBSUB"
			region        = %[5]q
			depends_on    = [mongodbatlas_advanced_cluster.test, mongodbatlas_cloud_provider_access_authorization.gcp_auth]
		}

		resource "mongodbatlas_stream_connection" "test" {
			project_id      = %[1]q
			workspace_name  = mongodbatlas_stream_workspace.test.workspace_name
			connection_name = %[4]q
			type            = "GCPPubSub"
			gcp = {
				service_account_id = mongodbatlas_cloud_provider_access_setup.gcp_setup.gcp_config[0].service_account_for_atlas
			}
			networking = {
				access = {
					type          = "PRIVATE_LINK"
					connection_id = mongodbatlas_stream_privatelink_endpoint.test.id
				}
			}
		}
	`, projectID, instanceName, clusterName, connectionName, region)
}

func checkGCPPubSubPrivateLinkAttributes(resourceName, workspaceName, connectionName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "GCPPubSub"),
		resource.TestCheckResourceAttrSet(resourceName, "gcp.service_account_id"),
		resource.TestCheckResourceAttr(resourceName, "networking.access.type", "PRIVATE_LINK"),
		resource.TestCheckResourceAttrSet(resourceName, "networking.access.connection_id"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
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

func TestAccStreamRSStreamConnection_AzureBlobStorage(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = acc.RandomName()
		clientID                = os.Getenv("AZURE_CLIENT_ID")
		clientSecret            = os.Getenv("AZURE_APP_SECRET")
		subscriptionID          = os.Getenv("AZURE_SUBSCRIPTION_ID")
		tenantID                = os.Getenv("AZURE_TENANT_ID")
		atlasAzureAppID         = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID      = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		resourceGroupName       = acc.RandomName()
		storageAccountName      = "tfacctest" + acctest.RandString(10)
		storageContainerName    = acc.RandomBucketName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckAzureEnvWithServicePrincipal(t)
			azureServicePrincipalMu.Lock()
			t.Cleanup(azureServicePrincipalMu.Unlock)
		},
		ExternalProviders:        acc.ExternalProvidersOnlyAzurerm(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig + configureAzureBlobStorage(projectID, instanceName, connectionName, clientID, clientSecret, subscriptionID, tenantID, atlasAzureAppID, servicePrincipalID, resourceGroupName, storageAccountName, storageContainerName, networkingTypePublic),
				Check:  checkAzureBlobStorageAttributes(resourceName, dataSourceName),
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

func TestAccStreamRSStreamConnection_AzureBlobStoragePrivateLink(t *testing.T) {
	var (
		projectID            = acc.ProjectIDExecution(t)
		instanceName         = acc.RandomStreamInstanceName()
		clusterName          = acc.RandomClusterName()
		connectionName       = acc.RandomName()
		clientID             = os.Getenv("AZURE_CLIENT_ID")
		clientSecret         = os.Getenv("AZURE_APP_SECRET")
		subscriptionID       = os.Getenv("AZURE_SUBSCRIPTION_ID")
		tenantID             = os.Getenv("AZURE_TENANT_ID")
		atlasAzureAppID      = os.Getenv("AZURE_ATLAS_APP_ID")
		servicePrincipalID   = os.Getenv("AZURE_SERVICE_PRINCIPAL_ID")
		resourceGroupName    = acc.RandomName()
		storageAccountName   = "tfacctest" + acctest.RandString(10)
		storageContainerName = acc.RandomBucketName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckAzureEnvWithServicePrincipal(t)
			azureServicePrincipalMu.Lock()
			t.Cleanup(azureServicePrincipalMu.Unlock)
		},
		ExternalProviders:        acc.ExternalProvidersOnlyAzurerm(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig + configureAzureBlobStoragePrivateLink(projectID, instanceName, clusterName, connectionName, clientID, clientSecret, subscriptionID, tenantID, atlasAzureAppID, servicePrincipalID, resourceGroupName, storageAccountName, storageContainerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAzureBlobStoragePrivateLinkAttributes(resourceName, instanceName, connectionName, servicePrincipalID, storageAccountName),
					checkAzureBlobStoragePrivateLinkAttributes(dataSourceName, instanceName, connectionName, servicePrincipalID, storageAccountName),
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
func configAzureBlobStorageStreamConnection(projectID, workspaceName, connectionName, networkingType string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_connection" "test" {
			project_id      = %[1]q
			workspace_name  = %[2]q
			connection_name = %[3]q
			type            = "AzureBlobStorage"
			azure = {
				service_principal_id = mongodbatlas_cloud_provider_access_setup.azure_setup.azure_config[0].service_principal_id
				storage_account_name = azurerm_storage_account.blob_storage.name
				region               = azurerm_resource_group.blob_rg.location
			}
			networking = {
				access = {
					type = %[4]q
				}
			}
			depends_on = [
				mongodbatlas_cloud_provider_access_authorization.azure_auth,
				azurerm_role_assignment.blob_contributor,
			]
		}
	`, projectID, workspaceName, connectionName, networkingType)
}

func configureAzureBlobStorage(projectID, workspaceName, connectionName, clientID, clientSecret, subscriptionID, tenantID, atlasAzureAppID, servicePrincipalID, resourceGroupName, storageAccountName, storageContainerName, networkingType string) string {
	return acc.ConfigAzurermProvider(subscriptionID, clientID, clientSecret, tenantID) +
		acc.ConfigAzureCloudProviderAccess(projectID, atlasAzureAppID, servicePrincipalID, tenantID) +
		acc.ConfigAzureStorageResources("blob", resourceGroupName, storageAccountName, storageContainerName, servicePrincipalID) +
		configAzureBlobStorageStreamConnection(projectID, workspaceName, connectionName, networkingType)
}

func checkAzureBlobStorageAttributes(resourceNames ...string) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range resourceNames {
		checks = append(checks,
			checkStreamConnectionExists(),
			resource.TestCheckResourceAttrSet(name, "project_id"),
			resource.TestCheckResourceAttrSet(name, "workspace_name"),
			resource.TestCheckResourceAttrSet(name, "connection_name"),
			resource.TestCheckResourceAttr(name, "type", "AzureBlobStorage"),
			resource.TestCheckResourceAttrSet(name, "azure.service_principal_id"),
			resource.TestCheckResourceAttrSet(name, "azure.storage_account_name"),
			resource.TestCheckResourceAttrSet(name, "azure.region"),
			resource.TestCheckResourceAttr(name, "networking.access.type", networkingTypePublic),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configureAzureBlobStoragePrivateLink(projectID, workspaceName, clusterName, connectionName, clientID, clientSecret, subscriptionID, tenantID, atlasAzureAppID, servicePrincipalID, resourceGroupName, storageAccountName, storageContainerName string) string {
	return acc.ConfigAzurermProvider(subscriptionID, clientID, clientSecret, tenantID) +
		acc.ConfigAzureCloudProviderAccess(projectID, atlasAzureAppID, servicePrincipalID, tenantID) +
		acc.ConfigAzureStorageResources("blob", resourceGroupName, storageAccountName, storageContainerName, servicePrincipalID) +
		configAzureBlobStoragePrivateLinkResources(projectID, workspaceName, clusterName, connectionName)
}

func configAzureBlobStoragePrivateLinkResources(projectID, workspaceName, clusterName, connectionName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[4]q
			cluster_type = "REPLICASET"
			replication_specs = [{
				region_configs = [{
					priority      = 7
					provider_name = "AZURE"
					region_name   = "US_EAST_2"
					electable_specs = {
						instance_size = "M10"
						node_count    = 3
					}
				}]
			}]
		}

		resource "mongodbatlas_stream_workspace" "test" {
			project_id     = %[1]q
			workspace_name = %[2]q
			data_process_region = {
				region         = "eastus2"
				cloud_provider = "AZURE"
			}
		}

		resource "mongodbatlas_stream_privatelink_endpoint" "test" {
			project_id          = %[1]q
			provider_name       = "AZURE"
			vendor              = "AZURE_BLOB_STORAGE"
			region              = azurerm_resource_group.blob_rg.location
			service_endpoint_id = azurerm_storage_account.blob_storage.id
			dns_domain          = "${azurerm_storage_account.blob_storage.name}.blob.core.windows.net"
			depends_on          = [mongodbatlas_advanced_cluster.test]
		}

		resource "mongodbatlas_stream_connection" "test" {
			project_id      = %[1]q
			workspace_name  = mongodbatlas_stream_workspace.test.workspace_name
			connection_name = %[3]q
			type            = "AzureBlobStorage"
			azure = {
				service_principal_id = mongodbatlas_cloud_provider_access_setup.azure_setup.azure_config[0].service_principal_id
				storage_account_name = azurerm_storage_account.blob_storage.name
				region               = azurerm_resource_group.blob_rg.location
			}
			networking = {
				access = {
					type          = "PRIVATE_LINK"
					connection_id = mongodbatlas_stream_privatelink_endpoint.test.id
				}
			}
			depends_on = [
				mongodbatlas_cloud_provider_access_authorization.azure_auth,
				azurerm_role_assignment.blob_contributor,
			]
		}
	`, projectID, workspaceName, connectionName, clusterName)
}

func checkAzureBlobStoragePrivateLinkAttributes(resourceName, workspaceName, connectionName, servicePrincipalID, storageAccountName string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkStreamConnectionExists(),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
		resource.TestCheckResourceAttr(resourceName, "type", "AzureBlobStorage"),
		resource.TestCheckResourceAttr(resourceName, "azure.service_principal_id", servicePrincipalID),
		resource.TestCheckResourceAttr(resourceName, "azure.storage_account_name", storageAccountName),
		resource.TestCheckResourceAttrSet(resourceName, "azure.region"),
		resource.TestCheckResourceAttr(resourceName, "networking.access.type", networkingTypePrivatelink),
		resource.TestCheckResourceAttrSet(resourceName, "networking.access.connection_id"),
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

// TestAccStreamRSStreamConnection_withFailoverConnections exercises the full CRUD lifecycle for
// failover_connections: create with one failover, update it, then remove it.
// The failover connection name must match the primary connection name (server requirement).
func TestAccStreamRSStreamConnection_withFailoverConnections(t *testing.T) {
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		connectionName          = "kafka-with-failover"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: configureKafkaWithFailovers(projectID, instanceName, connectionName, []failoverKafkaConfig{
					{name: connectionName, region: "DUBLIN_IRL", bootstrapServers: "failover1:9092", username: "fcuser", password: "fcpass"},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamConnectionExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "type", "Kafka"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.type", "Kafka"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.region", "DUBLIN_IRL"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.bootstrap_servers", "failover1:9092"),
					resource.TestCheckResourceAttrSet(resourceName, "failover_connections.0.id"),
				),
			},
			{
				Config: configureKafkaWithFailovers(projectID, instanceName, connectionName, []failoverKafkaConfig{
					{name: connectionName, region: "DUBLIN_IRL", bootstrapServers: "failover1-updated:9093", username: "fcuser", password: "fcpass"},
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamConnectionExists(),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.region", "DUBLIN_IRL"),
					resource.TestCheckResourceAttr(resourceName, "failover_connections.0.bootstrap_servers", "failover1-updated:9093"),
				),
			},
			{
				Config: configureKafka(fmt.Sprintf("%q", projectID), instanceName, connectionName,
					getKafkaAuthenticationConfig("PLAIN", "user", "rawpassword", "", "", "", "", "", ""),
					"localhost:9092", "earliest", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkStreamConnectionExists(),
					resource.TestCheckNoResourceAttr(resourceName, "failover_connections.0.name"),
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
	})
}

type failoverKafkaConfig struct {
	name             string
	region           string
	bootstrapServers string
	username         string
	password         string
}

func configureKafkaWithFailovers(projectID, workspaceName, connectionName string, failovers []failoverKafkaConfig) string {
	var fcBlocks strings.Builder
	for _, fc := range failovers {
		fmt.Fprintf(&fcBlocks, `
		{
			name              = %q
			type              = "Kafka"
			region            = %q
			bootstrap_servers = %q
			authentication = {
				mechanism = "PLAIN"
				username  = %q
				password  = %q
			}
			security = {
				protocol = "SASL_PLAINTEXT"
			}
			config = {
				"auto.offset.reset" = "earliest"
			}
		},`, fc.name, fc.region, fc.bootstrapServers, fc.username, fc.password)
	}

	failoverAttr := ""
	if fcBlocks.Len() > 0 {
		failoverAttr = fmt.Sprintf("failover_connections = [%s\n\t\t]", fcBlocks.String())
	}

	return fmt.Sprintf(`
	resource "mongodbatlas_stream_connection" "test" {
		project_id        = %q
		workspace_name    = %q
		connection_name   = %q
		type              = "Kafka"
		bootstrap_servers = "localhost:9092"
		authentication = {
			mechanism = "PLAIN"
			username  = "user"
			password  = "rawpassword"
		}
		config = {
			"auto.offset.reset" = "earliest"
		}
		security = {
			protocol = "SASL_PLAINTEXT"
		}
		%s
	}
`, projectID, workspaceName, connectionName, failoverAttr)
}
