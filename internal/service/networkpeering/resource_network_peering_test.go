package networkpeering_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName          = "mongodbatlas_network_peering.test"
	resourceNameContainer = "mongodbatlas_network_container.test"
	dataSourceName        = "data.mongodbatlas_network_peering.test"
	pluralDataSourceName  = "data.mongodbatlas_network_peerings.test"
)

func TestAccNetworkNetworkPeering_basicAWS(t *testing.T) {
	resource.ParallelTest(t, *basicAWSTestCase(t))
}

func TestAccNetworkRSNetworkPeering_basicAzure(t *testing.T) {
	var (
		projectID         = acc.ProjectIDExecution(t)
		directoryID       = os.Getenv("AZURE_DIRECTORY_ID")
		subscriptionID    = os.Getenv("AZURE_SUBSCRIPTION_ID")
		resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		vNetName          = os.Getenv("AZURE_VNET_NAME")
		providerName      = "AZURE"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckPeeringEnvAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configAzure(projectID, providerName, directoryID, subscriptionID, resourceGroupName, vNetName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "vnet_name", vNetName),
					resource.TestCheckResourceAttr(resourceName, "azure_directory_id", directoryID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"container_id"},
			},
		},
	})
}

func TestAccNetworkRSNetworkPeering_updateBasicAzure(t *testing.T) {
	var (
		projectID         = acc.ProjectIDExecution(t)
		directoryID       = os.Getenv("AZURE_DIRECTORY_ID")
		subscriptionID    = os.Getenv("AZURE_SUBSCRIPTION_ID")
		resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		vNetName          = os.Getenv("AZURE_VNET_NAME")
		updatedvNetName   = os.Getenv("AZURE_VNET_NAME_UPDATED")
		providerName      = "AZURE"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckPeeringEnvAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configAzure(projectID, providerName, directoryID, subscriptionID, resourceGroupName, vNetName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "vnet_name", vNetName),
					resource.TestCheckResourceAttr(resourceName, "azure_directory_id", directoryID),
				),
			},
			{
				Config: configAzure(projectID, providerName, directoryID, subscriptionID, resourceGroupName, updatedvNetName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "vnet_name", updatedvNetName),
					resource.TestCheckResourceAttr(resourceName, "azure_directory_id", directoryID),
				),
			},
		},
	})
}

func TestAccNetworkRSNetworkPeering_basicGCP(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		providerName = "GCP"
		gcpProjectID = os.Getenv("GCP_PROJECT_ID")
		networkName  = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPeeringEnvGCP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configGCP(projectID, providerName, gcpProjectID, networkName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),

					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "gcp_project_id", gcpProjectID),
					resource.TestCheckResourceAttr(resourceName, "network_name", networkName),

					// computed values that are obtain from associated container, checks for existing prefix convention to ensure they are gcp related values
					resource.TestCheckResourceAttrWith(resourceName, "atlas_gcp_project_id", acc.MatchesExpression("p-.*")),
					resource.TestCheckResourceAttrWith(resourceName, "atlas_vpc_name", acc.MatchesExpression("nt-.*")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkRSNetworkPeering_updateBasicGCP(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		projectID          = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		providerName       = "GCP"
		gcpProjectID       = os.Getenv("GCP_PROJECT_ID")
		networkName        = acc.RandomName()
		updatedNetworkName = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPeeringEnvGCP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configGCP(projectID, providerName, gcpProjectID, networkName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_name"),

					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "gcp_project_id", gcpProjectID),
					resource.TestCheckResourceAttr(resourceName, "network_name", networkName),
				),
			},
			{
				Config: configGCP(projectID, providerName, gcpProjectID, updatedNetworkName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
					resource.TestCheckResourceAttrSet(resourceName, "network_name"),

					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "gcp_project_id", gcpProjectID),
					resource.TestCheckResourceAttr(resourceName, "network_name", updatedNetworkName),
				),
			},
		},
	})
}

func TestAccNetworkRSNetworkPeering_AWSDifferentRegionName(t *testing.T) {
	var (
		vpcID           = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock    = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID    = os.Getenv("AWS_ACCOUNT_ID")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		containerRegion = "US_WEST_2"
		peerRegion      = conversion.MongoDBRegionToAWSRegion(os.Getenv("AWS_REGION"))
		providerName    = "AWS"
		projectName     = acc.RandomProjectName()
	)
	checks := commonChecksAWS(vpcID, providerName, awsAccountID, vpcCIDRBlock, peerRegion)
	checks = acc.AddAttrSetChecks(resourceNameContainer, checks, "project_id", "atlas_cidr_block", "provider_name", "region_name")
	mapChecksContainer := map[string]string{
		"atlas_cidr_block": "192.168.208.0/21",
		"provider_name":    providerName,
		"region_name":      containerRegion,
	}
	checks = acc.AddAttrChecks(resourceNameContainer, checks, mapChecksContainer)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
			acc.PreCheckPeeringEnvAWS(t)
			func() {
				if strings.EqualFold(containerRegion, peerRegion) {
					t.Fatalf("the `AWS_REGION` (%s) must be different region than %s", peerRegion, containerRegion)
				}
			}()
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configAWS(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, containerRegion, peerRegion),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func basicAWSTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		vpcID           = os.Getenv("AWS_VPC_ID")
		vpcCIDRBlock    = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID    = os.Getenv("AWS_ACCOUNT_ID")
		containerRegion = os.Getenv("AWS_REGION")
		peerRegion      = conversion.MongoDBRegionToAWSRegion(containerRegion)
		providerName    = "AWS"
		projectName     = acc.RandomProjectName()
	)
	checks := commonChecksAWS(vpcID, providerName, awsAccountID, vpcCIDRBlock, peerRegion)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb); acc.PreCheckPeeringEnvAWS(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyNetworkPeering,
		Steps: []resource.TestStep{
			{
				Config: configAWS(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, containerRegion, peerRegion),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accepter_region_name", "container_id"},
			},
		},
	}
}

func commonChecksAWS(vpcID, providerName, awsAccountID, vpcCIDRBlock, regionPeer string) []resource.TestCheckFunc {
	attributes := map[string]string{
		"vpc_id":                 vpcID,
		"provider_name":          providerName,
		"aws_account_id":         awsAccountID,
		"route_table_cidr_block": vpcCIDRBlock,
		"accepter_region_name":   regionPeer,
	}
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrChecks(dataSourceName, checks, attributes)
	checks = acc.AddAttrSetChecks(resourceName, checks, "project_id", "container_id", "accepter_region_name")
	checks = acc.AddAttrSetChecks(dataSourceName, checks, "project_id", "container_id")
	checks = acc.AddAttrSetChecks(pluralDataSourceName, checks, "results.#", "results.0.provider_name", "results.0.vpc_id", "results.0.aws_account_id", "results.0.accepter_region_name")
	return checks
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["peer_id"], ids["provider_name"]), nil
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		log.Printf("[DEBUG] projectID: %s", ids["project_id"])
		if _, _, err := acc.ConnV2().NetworkPeeringApi.GetPeeringConnection(context.Background(), ids["project_id"], ids["peer_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("peer(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["peer_id"])
	}
}

func configAWS(orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegionContainer, awsRegionPeer string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "my_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_network_container" "test" {
		project_id   		 = mongodbatlas_project.my_project.id
		atlas_cidr_block  	 = "192.168.208.0/21"
		provider_name		 = %[3]q
		region_name			 = %[7]q
	}

	resource "mongodbatlas_network_peering" "test" {
		accepter_region_name	= %[8]q
		project_id    			= mongodbatlas_project.my_project.id
		container_id           	= mongodbatlas_network_container.test.id
		provider_name           = %[3]q
		route_table_cidr_block  = %[6]q
		vpc_id					= %[4]q
		aws_account_id	        = %[5]q
	}

	data "mongodbatlas_network_peering" "test" {
		project_id = mongodbatlas_project.my_project.id
		peering_id = mongodbatlas_network_peering.test.peer_id
	}

	data "mongodbatlas_network_peerings" "test" {
		project_id = mongodbatlas_network_peering.test.project_id
	}
`, orgID, projectName, providerName, vpcID, awsAccountID, vpcCIDRBlock, awsRegionContainer, awsRegionPeer)
}

func configAzure(projectID, providerName, directoryID, subscriptionID, resourceGroupName, vNetName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		  = %[1]q
			atlas_cidr_block  = "192.168.208.0/21"
			provider_name		  = %[2]q
			region    			  = "US_EAST_2"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id   		  = %[1]q
			container_id          = mongodbatlas_network_container.test.container_id
			provider_name         = %[2]q
			azure_directory_id    = %[3]q
			azure_subscription_id = %[4]q
			resource_group_name   = %[5]q
			vnet_name	          = %[6]q
		}
	`, projectID, providerName, directoryID, subscriptionID, resourceGroupName, vNetName)
}

func configGCP(projectID, providerName, gcpProjectID, networkName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id       = %[1]q
			atlas_cidr_block = "192.168.192.0/18"
			provider_name    = %[2]q
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id     = %[1]q
			container_id   = mongodbatlas_network_container.test.container_id
			provider_name  = %[2]q
			gcp_project_id = %[3]q
			network_name   = %[4]q
		}
	`, projectID, providerName, gcpProjectID, networkName)
}
