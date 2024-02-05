package privatelinkendpoint_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkRSPrivateLinkEndpointAWS_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		region       = "us-east-1"
		providerName = "AWS"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
		},
	})
}

func TestAccNetworkRSPrivateLinkEndpointAWS_import(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		region       = "us-east-1"
		providerName = "AWS"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
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
func TestAccNetworkRSPrivateLinkEndpointAzure_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		region       = "US_EAST_2"
		providerName = "AZURE"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
		},
	})
}

func TestAccNetworkRSPrivateLinkEndpointAzure_import(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		region       = "US_EAST_2"
		providerName = "AZURE"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
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

func TestAccNetworkRSPrivateLinkEndpointGCP_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		region       = "us-central1"
		providerName = "GCP"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
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

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s-%s", ids["project_id"], ids["private_link_id"], ids["provider_name"], ids["region"]), nil
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
		if _, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpointService(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("the MongoDB Private Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["private_link_id"], rs.Primary.Attributes["project_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().PrivateEndpointServicesApi.GetPrivateEndpointService(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"]).Execute()
		if err == nil {
			return fmt.Errorf("the MongoDB Private Endpoint(%s) still exists", ids["private_link_id"])
		}
	}
	return nil
}

func configBasic(orgID, projectName, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = mongodbatlas_project.test.id
			provider_name = %[3]q
			region        = %[4]q
		}
	`, orgID, projectName, providerName, region)
}
