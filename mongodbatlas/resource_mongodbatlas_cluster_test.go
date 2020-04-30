package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
	"github.com/mwielbut/pointy"
)

func TestAccResourceMongoDBAtlasCluster_basicAWS(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAWS(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigAWS(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_basicAdvancedConf(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              pointy.Bool(true),
					JavascriptEnabled:                pointy.Bool(true),
					MinimumEnabledTLSProtocol:        "TLS1_2",
					NoTableScan:                      pointy.Bool(true),
					OplogSizeMB:                      pointy.Int64(1000),
					SampleRefreshIntervalBIConnector: pointy.Int64(310),
					SampleSizeBIConnector:            pointy.Int64(110),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.fail_index_key_too_long", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.no_table_scan", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttrSet(resourceName, "advanced_configuration.%"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, "true", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              pointy.Bool(false),
					JavascriptEnabled:                pointy.Bool(false),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      pointy.Bool(false),
					OplogSizeMB:                      pointy.Int64(990),
					SampleRefreshIntervalBIConnector: pointy.Int64(0),
					SampleSizeBIConnector:            pointy.Int64(0),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.javascript_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.oplog_size_mb", "990"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_size_bi_connector", "0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_refresh_interval_bi_connector", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "advanced_configuration.%"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_basicAzure(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAzure(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigAzure(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_basicGCP(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvGCP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigGCP(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigGCP(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_MultiRegion(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-multi-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_Global(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigGlobal(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.1.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "80"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "GEOSHARDED"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.1.regions_config.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_AWSWithLabels(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("testAcc-%s-%s-%s", "AWS", "M10", acctest.RandString(1))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, name, "false", "M10", "EU_CENTRAL_1", []matlas.Label{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, name, "false", "M10", "EU_CENTRAL_1",
					[]matlas.Label{
						{
							Key:   "key 4",
							Value: "value 4",
						},
						{
							Key:   "key 3",
							Value: "value 3",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "3"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, name, "false", "M10", "EU_CENTRAL_1",
					[]matlas.Label{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 5",
							Value: "value 5",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_withPrivateEndpointLink(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	vpcID := os.Getenv("AWS_VPC_ID")
	subnetID := os.Getenv("AWS_SUBNET_ID")
	securityGroupID := os.Getenv("AWS_SECURITY_GROUP_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkAwsEnv(t); checkPeeringEnvAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigWithPrivateEndpointLink(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_withAzureNetworkPeering(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	directoryID := os.Getenv("AZURE_DIRECTORY_ID")
	subcrptionID := os.Getenv("AZURE_SUBCRIPTION_ID")
	resourceGroupName := os.Getenv("AZURE_RESOURCE_GROUP_NAME")
	vNetName := os.Getenv("AZURE_VNET_NAME")
	providerName := "AZURE"
	region := os.Getenv("AZURE_REGION")

	atlasCidrBlock := "192.168.208.0/21"
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAzureWithNetworkPeering(projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_withGCPNetworkPeering(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	gcpRegion := os.Getenv("GCP_REGION_NAME")
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	providerName := "GCP"
	gcpPeeringName := fmt.Sprintf("test-acc-%s", acctest.RandString(3))
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(3))
	gcpClusterRegion := os.Getenv("GCP_CLUSTER_REGION_NAME")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvGCP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigGCPWithNetworkPeering(gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_withAzureAndContainerID(t *testing.T) {
	resourceName := "mongodbatlas_cluster.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	providerName := "AZURE"
	region := os.Getenv("AZURE_REGION")
	directoryID := os.Getenv("AZURE_DIRECTORY_ID")
	subcrptionID := os.Getenv("AZURE_SUBCRIPTION_ID")
	resourceGroupName := os.Getenv("AZURE_RESOURCE_GROUP_NAME")
	vNetName := os.Getenv("AZURE_VNET_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkPeeringEnvAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAzureWithContainerID(projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_withAWSAndContainerID(t *testing.T) {
	resourceName := "mongodbatlas_cluster.test"

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	providerName := "AWS"
	awsRegion := os.Getenv("AWS_REGION")
	vpcCIDRBlock := os.Getenv("AWS_VPC_CIDR_BLOCK")
	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAWSWithContainerID(awsAccessKey, awsSecretKey, projectID, clusterName, providerName, awsRegion, vpcCIDRBlock, awsAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	importStateID := fmt.Sprintf("%s-%s", projectID, clusterName)

	resourceName := "mongodbatlas_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAWS(projectID, clusterName, "true"),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           importStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_tenant(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.tenant"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigTenant(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigTenantUpdated(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
		},
	})

}

func testAccCheckMongoDBAtlasClusterExists(resourceName string, cluster *matlas.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if clusterResp, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]); err == nil {
			*cluster = *clusterResp
			return nil
		}

		return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasClusterAttributes(cluster *matlas.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cluster.Name != name {
			return fmt.Errorf("bad name: %s", cluster.Name)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" {
			continue
		}

		// Try to find the cluster
		_, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf("cluster (%s:%s) still exists", rs.Primary.Attributes["name"], rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasClusterConfigAWS(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 100
			num_shards   = 1
			replication_factor           = 3
			provider_backup_enabled      = %[3]s
			pit_enabled 				 = %[3]s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops 			    = 300
			provider_encrypt_ebs_volume = false
			provider_instance_size_name = "M30"
			provider_region_name        = "EU_CENTRAL_1"
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, backupEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 10
			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
			provider_instance_size_name = "M10"
			provider_region_name        = "EU_CENTRAL_1"

			advanced_configuration = {
				fail_index_key_too_long              = %t
				javascript_enabled                   = %t
				minimum_enabled_tls_protocol         = "%s"
				no_table_scan                        = %t
				oplog_size_mb                        = %d
				sample_size_bi_connector			       = %d
				sample_refresh_interval_bi_connector = %d
			}
		}
	`, projectID, name, backupEnabled,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector)
}

func testAccMongoDBAtlasClusterConfigAzure(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			num_shards   = 1

			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "AZURE"
			provider_disk_type_name     = "P6"
			provider_instance_size_name = "M30"
			provider_region_name        = "US_EAST_2"
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigGCP(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 40
			num_shards   = 1

			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "GCP"
			provider_instance_size_name = "M30"
			provider_region_name        = "US_EAST_4"
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id     = "%s"
			name           = "%s"
			disk_size_gb   = 100
			num_shards     = 1
			backup_enabled = %s
			cluster_type   = "REPLICASET"
			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 300
			provider_instance_size_name = "M10"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "EU_CENTRAL_1"
					electable_nodes = 2
					priority        = 6
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_WEST_1"
					electable_nodes = 2
					priority        = 5
					read_only_nodes = 2
				}
			}
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigGlobal(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id              = "%s"
			name                    = "%s"
			disk_size_gb            = 80
			num_shards              = 1
			backup_enabled          = %s
			provider_backup_enabled = true
			cluster_type            = "GEOSHARDED"

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 240
			provider_instance_size_name = "M30"

			replication_specs {
				zone_name  = "Zone 1"
				num_shards = 2
				regions_config {
				region_name     = "EU_CENTRAL_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
			}

			replication_specs {
				zone_name  = "Zone 2"
				num_shards = 2
				regions_config {
				region_name     = "US_EAST_2"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
			}
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigTenant(projectID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cluster" "tenant" {
		project_id = "%s"
		name       = "%s"

		provider_name         = "TENANT"
		backing_provider_name = "AWS"
		provider_region_name  = "EU_CENTRAL_1"

		provider_instance_size_name  = "M2"
		disk_size_gb                 = 2
		auto_scaling_disk_gb_enabled = false
	  }
	`, projectID, name)
}

func testAccMongoDBAtlasClusterConfigTenantUpdated(projectID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cluster" "tenant" {
		project_id = "%s"
		name       = "%s"

		provider_name        = "AWS"
		provider_region_name = "EU_CENTRAL_1"

		provider_instance_size_name  = "M10"
		disk_size_gb                 = 10
		provider_disk_iops           = 100
		auto_scaling_disk_gb_enabled = true
	  }
	`, projectID, name)
}

func testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, name, backupEnabled, tier, region string, labels []matlas.Label) string {
	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			num_shards   = 1
			disk_size_gb = 10
			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = false
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
			%s
		}
	`, projectID, name, backupEnabled, tier, region, labelsConf)
}

func testAccMongoDBAtlasClusterConfigWithPrivateEndpointLink(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "${lower(replace("%[5]s", "_", "-"))}"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%[6]s"
			service_name       = "${mongodbatlas_private_endpoint.test.endpoint_service_name}"
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%[7]s"]
			security_group_ids = ["%[8]s"]
		}

		resource "mongodbatlas_private_endpoint_interface_link" "test" {
			project_id            = "${mongodbatlas_private_endpoint.test.project_id}"
			private_link_id       = "${mongodbatlas_private_endpoint.test.private_link_id}"
			interface_endpoint_id = "${aws_vpc_endpoint.ptfe_service.id}"
		}

		resource "mongodbatlas_cluster" "test" {
		  project_id             = "%[3]s"
		  name                   = "%[9]s"
		  disk_size_gb           = 5

		  //Provider Settings "block"
		  provider_name               = "AWS"
		  provider_region_name        = "${upper(replace("%[5]s", "-", "_"))}"
		  provider_instance_size_name = "M10"
		  provider_backup_enabled     = true // enable cloud provider snapshots
		  provider_disk_iops          = 100
		  provider_encrypt_ebs_volume = false
		  depends_on = ["mongodbatlas_private_endpoint_interface_link.test"]
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName)
}

func testAccMongoDBAtlasClusterConfigAzureWithNetworkPeering(projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id       = "%[1]s"
			atlas_cidr_block = "%[8]s"
			provider_name    = "%[2]s"
			region           = "%[9]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id            = "%[1]s"
			atlas_cidr_block      = "192.168.0.0/21"
			container_id          = mongodbatlas_network_container.test.container_id
			provider_name         = "%[2]s"
			azure_directory_id    = "%[3]s"
			azure_subscription_id = "%[4]s"
			resource_group_name   = "%[5]s"
			vnet_name             = "%[6]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[1]s"
			name         = "%[7]s"

			replication_factor           = 3
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "%[2]s"
			provider_disk_type_name     = "P6"
			provider_instance_size_name = "M10"
			provider_region_name        = "%[9]s"

			depends_on = ["mongodbatlas_network_peering.test"]
		}
	`, projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region)
}

func testAccMongoDBAtlasClusterConfigGCPWithNetworkPeering(gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion string) string {
	return fmt.Sprintf(`
		provider "google" {
			project     = "%[1]s"
			region      = "%[2]s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id       = "%[3]s"
			atlas_cidr_block = "192.168.192.0/18"
			provider_name    = "%[4]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id     = "%[3]s"
			container_id   = "${mongodbatlas_network_container.test.container_id}"
			provider_name  = "%[4]s"
			gcp_project_id = "%[1]s"
			network_name   = "default"
		}

		data "google_compute_network" "default" {
			name = "default"
		}

		resource "google_compute_network_peering" "gcp_peering" {
			name         = "%[5]s"
			network      = "${data.google_compute_network.default.self_link}"
			peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[3]s"
			name         = "%[6]s"
			disk_size_gb = 5
			num_shards   = 1

			replication_factor           = 3
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "%[4]s"
			provider_instance_size_name = "M10"
			provider_region_name        = "%[7]s"

			depends_on = ["google_compute_network_peering.gcp_peering"]
		}
	`, gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion)
}

func testAccMongoDBAtlasClusterConfigAzureWithContainerID(projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[1]s"
			name         = "%[2]s"

			replication_factor           = 3
			auto_scaling_disk_gb_enabled = false
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "%[3]s"
			provider_instance_size_name = "M10"
			provider_region_name        = "%[4]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id            = "%[1]s"
			atlas_cidr_block      = "192.168.0.0/21"
			container_id          = "${mongodbatlas_cluster.test.container_id}"
			provider_name         = "%[3]s"
			azure_directory_id    = "%[5]s"
			azure_subscription_id = "%[6]s"
			resource_group_name   = "%[7]s"
			vnet_name             = "%[8]s"
		}
	`, projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName)
}

func testAccMongoDBAtlasClusterConfigAWSWithContainerID(awsAccessKey, awsSecretKey, projectID, clusterName, providerName, region, vpcCIDRBlock, awsAccountID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = lower(replace("%[6]s", "_", "-"))
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[3]s"
			name         = "%[4]s"
			disk_size_gb = 5

			replication_factor           = 3
			auto_scaling_disk_gb_enabled = false
			mongo_db_major_version       = "4.0"

			//Provider Settings "block"
			provider_name               = "%[5]s"
			provider_instance_size_name = "M10"
			provider_region_name        = "%[6]s"
		}

		resource "aws_default_vpc" "default" {
			tags = {
				Name = "Default VPC"
			}
		}

		resource "mongodbatlas_network_peering" "mongo_peer" {
			accepter_region_name   = lower(replace("%[6]s", "_", "-"))
			project_id             = "%[3]s"
			container_id           = mongodbatlas_cluster.test.container_id
			provider_name          = "%[5]s"
			route_table_cidr_block = "%[7]s"
			vpc_id                 = aws_default_vpc.default.id
			aws_account_id         = "%[8]s"
		}

		resource "aws_vpc_peering_connection_accepter" "aws_peer" {
			vpc_peering_connection_id = mongodbatlas_network_peering.mongo_peer.connection_id
			auto_accept               = true

			tags = {
				Side = "Accepter"
			}
		}
	`, awsAccessKey, awsSecretKey, projectID, clusterName, providerName, region, vpcCIDRBlock, awsAccountID)
}
