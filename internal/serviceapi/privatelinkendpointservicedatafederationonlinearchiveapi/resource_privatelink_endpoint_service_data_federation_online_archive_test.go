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
	resourceName   = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test"
	comment        = "Terraform Acceptance Test"
	AWSregion      = "US_EAST_1"
	dataSourceName = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test"
	pluralDSName   = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives_api.test"
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
				Config: dataSourceConfigBasicAWS(projectID, endpointID, comment),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("region"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("customer_endpoint_dns_name"), knownvalue.StringExact("")),
					acc.PluralResultCheck(
						pluralDSName,
						"endpoint_id",
						knownvalue.StringExact(endpointID),
						map[string]knownvalue.Check{
							"region":                     knownvalue.StringExact(""),
							"customer_endpoint_dns_name": knownvalue.StringExact(""),
						},
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					checkEncodedID(resourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
					resource.TestCheckResourceAttr(dataSourceName, "region", ""),
					resource.TestCheckResourceAttr(dataSourceName, "customer_endpoint_dns_name", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "provider_name"),
					checkDataSourceEncodedID(dataSourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(pluralDSName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.#"),
					resource.TestCheckResourceAttrSet(pluralDSName, "id"),
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

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_providerNameLowercase(t *testing.T) {
	var (
		projectID      = acc.ProjectIDExecution(t)
		endpointID     = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		commentValue   = "Terraform Acceptance Test"
		commentUpdated = "Terraform Acceptance Test Lowercase Updated"
	)

	checkWithComment := func(expectedComment string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			checkExists(resourceName),
			checkEncodedID(resourceName, projectID, endpointID),
			resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
			resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
			resource.TestCheckResourceAttr(resourceName, "provider_name", "aws"),
			resource.TestCheckResourceAttr(resourceName, "comment", expectedComment),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWSLowercase(projectID, endpointID, commentValue),
				Check:  checkWithComment(commentValue),
			},
			{
				Config: resourceConfigBasicAWSLowercase(projectID, endpointID, commentUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkWithComment(commentUpdated),
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

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_timeoutsBlock(t *testing.T) {
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
				Config: resourceConfigBasicAWSWithTimeouts(projectID, endpointID, "600s", "900s"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "600s"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.delete", "900s"),
				),
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchiveDS_OptionalFieldsAWS(t *testing.T) {
	var (
		projectID               = acc.ProjectIDExecution(t)
		endpointID              = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		customerEndpointDNSName = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_DNS_NAME")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpoint(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigOptionalFieldsAWS(projectID, endpointID, comment, customerEndpointDNSName),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(
						pluralDSName,
						"endpoint_id",
						knownvalue.StringExact(endpointID),
						map[string]knownvalue.Check{
							"region":                     knownvalue.StringExact(AWSregion),
							"customer_endpoint_dns_name": knownvalue.StringExact(customerEndpointDNSName),
						},
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "endpoint_id", endpointID),
					resource.TestCheckResourceAttr(dataSourceName, "comment", comment),
					resource.TestCheckResourceAttr(dataSourceName, "region", AWSregion),
					resource.TestCheckResourceAttr(dataSourceName, "customer_endpoint_dns_name", customerEndpointDNSName),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "provider_name"),
					checkDataSourceEncodedID(dataSourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(pluralDSName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.#"),
					resource.TestCheckResourceAttrSet(pluralDSName, "id"),
				),
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

func checkDataSourceEncodedID(resourceName, expectedProjectID, expectedEndpointID string) resource.TestCheckFunc {
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

func resourceConfigBasicAWSLowercase(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "aws"
	  comment					= %[3]q
	}
	`, projectID, endpointID, comment)
}

func resourceConfigBasicAWSWithTimeouts(projectID, endpointID, createTimeout, deleteTimeout string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id    = %[1]q
	  endpoint_id   = %[2]q
	  provider_name = "AWS"

	  timeouts {
	    create = %[3]q
	    delete = %[4]q
	  }
	}
	`, projectID, endpointID, createTimeout, deleteTimeout)
}

func dataSourceConfigBasicAWS(projectID, endpointID, comment string) string {
	return buildConfigWithDataSource(projectID, endpointID, comment, "", false)
}

func dataSourceConfigOptionalFieldsAWS(projectID, endpointID, comment, customerEndpointDNSName string) string {
	return buildConfigWithDataSource(projectID, endpointID, comment, customerEndpointDNSName, true)
}

func buildConfigWithDataSource(projectID, endpointID, comment, customerEndpointDNSName string, includeOptionalFields bool) string {
	optionalFields := ""
	if includeOptionalFields {
		optionalFields = fmt.Sprintf(`
	  region						= %q
	  customer_endpoint_dns_name	= %q`, AWSregion, customerEndpointDNSName)
	}

	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
	  project_id					= %[1]q
	  endpoint_id					= %[2]q
	  provider_name					= "AWS"
	  comment						= %[3]q
	  %[4]s
	}

	%[5]s
	%[6]s
	`, projectID, endpointID, comment, optionalFields, singularDataSourceConfig, pluralDataSourceConfig)
}

const singularDataSourceConfig = `
data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api" "test" {
  project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test.project_id
  endpoint_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test.endpoint_id
}`

const pluralDataSourceConfig = `
data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives_api" "test" {
  project_id				= mongodbatlas_privatelink_endpoint_service_data_federation_online_archive_api.test.project_id
}`
