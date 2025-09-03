package networkcontainer_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_network_container.test"
	dataSourceName       = "data.mongodbatlas_network_container.test"
	dataSourcePluralName = "data.mongodbatlas_network_containers.test"
)

func TestAccNetworkContainer_basicAWS(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		randInt          = acctest.RandIntRange(0, 255)
		cidrBlock        = fmt.Sprintf("10.8.%d.0/24", randInt)
		randIntUpdated   = (randInt + 1) % 256
		cidrBlockUpdated = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, cidrBlock, constant.AWS, "US_EAST_1"),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
			{
				Config: configBasic(projectID, cidrBlockUpdated, constant.AWS, "US_EAST_2"),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccNetworkContainer_basicAzure(t *testing.T) {
	var (
		randInt          = acctest.RandIntRange(0, 255)
		cidrBlock        = fmt.Sprintf("10.8.%d.0/24", randInt)
		randIntUpdated   = (randInt + 1) % 256
		cidrBlockUpdated = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
		projectID        = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, cidrBlock, constant.AZURE, "US_EAST_2"),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AZURE)...),
			},
			{
				Config: configBasic(projectID, cidrBlockUpdated, constant.AZURE, "US_EAST_2"),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AZURE)...),
			},
		},
	})
}

func TestAccNetworkContainer_basicGCP(t *testing.T) {
	var (
		randInt          = acctest.RandIntRange(0, 255)
		gcpCidrBlock     = fmt.Sprintf("10.%d.0.0/18", randInt)
		randIntUpdated   = (randInt + 1) % 256
		cidrBlockUpdated = fmt.Sprintf("10.%d.0.0/18", randIntUpdated)
		projectID        = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, gcpCidrBlock, constant.GCP, ""),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.GCP)...),
			},
			{
				Config: configBasic(projectID, cidrBlockUpdated, constant.GCP, ""),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.GCP)...),
			},
		},
	})
}

func TestAccNetworkContainer_withRegionsGCP(t *testing.T) {
	var (
		projectID               = acc.ProjectIDExecution(t)
		randInt                 = acctest.RandIntRange(0, 255)
		gcpWithRegionsCidrBlock = fmt.Sprintf("10.%d.0.0/21", randInt)
		regions                 = "[\"US_EAST_4\", \"US_WEST_3\"]"
	)
	// No ParallelTest because only one GCP network container is allowed per project
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, gcpWithRegionsCidrBlock, constant.GCP, regions),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.GCP)...),
			},
		},
	})
}

func TestAccNetworkContainer_updateIndividualFields(t *testing.T) {
	var (
		projectID        = acc.ProjectIDExecution(t)
		randInt          = acctest.RandIntRange(0, 255)
		cidrBlock        = fmt.Sprintf("10.8.%d.0/24", randInt)
		randIntUpdated   = (randInt + 1) % 256
		cidrBlockUpdated = fmt.Sprintf("10.8.%d.0/24", randIntUpdated)
		region           = "EU_WEST_1"
		regionUpdated    = "EU_WEST_2"
	)

	// Serial so it doesn't conflict with TestAccNetworkContainer_basicAWS
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, cidrBlock, constant.AWS, region),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
			{
				Config: configBasic(projectID, cidrBlockUpdated, constant.AWS, region),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
			{
				Config: configBasic(projectID, cidrBlockUpdated, constant.AWS, regionUpdated),
				Check:  resource.ComposeAggregateTestCheckFunc(commonChecks(constant.AWS)...),
			},
		},
	})
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]), nil
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
		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])
		if _, _, err := acc.ConnV2().NetworkPeeringApi.GetGroupContainer(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("container(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_network_container" {
			continue
		}

		_, _, err := acc.ConnV2().NetworkPeeringApi.GetGroupContainer(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"]).Execute()

		if err == nil {
			return fmt.Errorf("container (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["container_id"])
		}
	}
	return nil
}

func commonChecks(providerName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
		resource.TestCheckResourceAttrSet(resourceName, "provisioned"),

		resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
		resource.TestCheckResourceAttr(dataSourceName, "provider_name", providerName),
		resource.TestCheckResourceAttrSet(dataSourceName, "provisioned"),

		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.id"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.atlas_cidr_block"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.provider_name"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.provisioned"),
	}
}

func configBasic(projectID, cidrBlock, providerName, region string) string {
	var regionStr string
	if region != "" {
		switch providerName {
		case constant.AWS:
			regionStr = fmt.Sprintf("region_name	= %q", region)
		case constant.AZURE:
			regionStr = fmt.Sprintf("region	= %q", region)
		case constant.GCP:
			regionStr = fmt.Sprintf("regions = %s", region)
		}
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id   		 = %[1]q
			atlas_cidr_block = %[2]q
			provider_name		 = %[3]q
			%[4]s
		}

		data "mongodbatlas_network_container" "test" {
			project_id   		= mongodbatlas_network_container.test.project_id
			container_id		= mongodbatlas_network_container.test.id
		}

		data "mongodbatlas_network_containers" "test" {
			project_id = mongodbatlas_network_container.test.project_id
			provider_name		 = %[3]q
		}
	`, projectID, cidrBlock, providerName, regionStr)
}
