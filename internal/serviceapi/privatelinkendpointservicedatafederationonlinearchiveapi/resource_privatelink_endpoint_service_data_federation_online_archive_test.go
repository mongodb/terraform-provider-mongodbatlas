package privatelinkendpointservicedatafederationonlinearchiveapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test"
	comment      = "Terraform Acceptance Test"
)

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basic(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkEncodedID(resourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importLegacyStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_on_create_timeout",
				},
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importNormalizedStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_on_create_timeout",
				},
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_updateComment(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		endpointID      = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		commentValue    = "Terraform Acceptance Test"
		commentUpdated2 = "Terraform Acceptance Test Updated Again"
	)
	checkWithComment := func(expectedComment string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			checkExists(resourceName),
			checkEncodedID(resourceName, projectID, endpointID),
			resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
			resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
			resource.TestCheckResourceAttr(resourceName, "comment", expectedComment),
			resource.TestCheckResourceAttrSet(resourceName, "type"),
			resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, ""),
				Check:  checkWithComment(""),
			},
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, commentValue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkWithComment(commentValue),
			},
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, commentUpdated2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkWithComment(commentUpdated2),
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_optionalStringEmptyState(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		endpointID = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, comment),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("region"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("customer_endpoint_dns_name"), knownvalue.StringExact("")),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkEncodedID(resourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importNormalizedStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_on_create_timeout",
				},
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_forceNewEndpointID(t *testing.T) {
	acc.SkipTestForCI(t)

	var (
		projectID   = acc.ProjectIDExecution(t)
		endpointID  = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		endpointID2 = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID_2")
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckPrivateEndpoint(t)
			if endpointID2 == "" {
				t.Fatal("`MONGODB_ATLAS_PRIVATE_ENDPOINT_ID_2` must be set for force-new endpoint_id acceptance testing")
			}
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkEncodedID(resourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: resourceConfigBasicAWS(projectID, endpointID2, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func importLegacyStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["endpoint_id"]), nil
	}
}

func importNormalizedStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s/%s", ids["project_id"], ids["endpoint_id"]), nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataFederationApi.GetPrivateEndpointId(context.Background(), ids["project_id"], ids["endpoint_id"]).Execute()
		if err == nil {
			return fmt.Errorf("Private endpoint service data federation online archive still exists")
		}
	}
	return nil
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Private endpoint service data federation online archive not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Private endpoint service data federation online archive ID not set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataFederationApi.GetPrivateEndpointId(context.Background(), ids["project_id"], ids["endpoint_id"]).Execute()
		if err != nil {
			return err
		}
		return nil
	}
}

func checkEncodedID(resourceName, expectedProjectID, expectedEndpointID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("id is empty")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		if ids["project_id"] != expectedProjectID || ids["endpoint_id"] != expectedEndpointID {
			return fmt.Errorf("unexpected decoded ID map: %+v", ids)
		}
		return nil
	}
}

func resourceConfigBasicAWS(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "AWS"
	  comment					= %[3]q
	}
	`, projectID, endpointID, comment)
}
