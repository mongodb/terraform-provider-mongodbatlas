package privatelinkendpointservicedatafederationonlinearchive_test

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
	resourceName   = "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	comment        = "Terraform Acceptance Test"
	AWSRegion      = "US_EAST_1"
	azureRegion    = "US_EAST_2"
	dataSourceName = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test"
	pluralDSName   = "data.mongodbatlas_privatelink_endpoint_service_data_federation_online_archives.test"
)

const singularDataSourceConfig = `
data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id  = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
  endpoint_id = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.endpoint_id
}`

const pluralDataSourceConfig = `
data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "test" {
  project_id = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.project_id
}`

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
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("region"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("customer_endpoint_dns_name"), knownvalue.StringExact("")),
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
				Check: checkAggr(projectID, endpointID, comment, "", "",
					checkEncodedID(resourceName, projectID, endpointID),
					checkEncodedID(dataSourceName, projectID, endpointID),
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

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_basicAzure(t *testing.T) {
	var (
		projectID                    = acc.ProjectIDExecution(t)
		subscriptionID               = os.Getenv("AZURE_SUBSCRIPTION_ID")
		clientID                     = os.Getenv("AZURE_CLIENT_ID")
		clientSecret                 = os.Getenv("AZURE_APP_SECRET")
		tenantID                     = os.Getenv("AZURE_TENANT_ID")
		privateLinkServiceResourceID = os.Getenv("MONGODB_ATLAS_DATA_FEDERATION_PRIVATE_LINK_SERVICE_RESOURCE_ID_AZURE")
		resourceGroupName            = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckPrivateEndpointAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		ExternalProviders:        acc.ExternalProvidersOnlyAzurerm(),
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicAzure(projectID, subscriptionID, clientID, clientSecret, tenantID, privateLinkServiceResourceID, resourceGroupName, comment, azureRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					acc.PluralResultCheck(
						pluralDSName,
						"provider_name",
						knownvalue.StringExact("AZURE"),
						map[string]knownvalue.Check{
							"region":                       knownvalue.StringExact(azureRegion),
							"endpoint_id":                  knownvalue.NotNull(),
							"customer_endpoint_ip_address": knownvalue.NotNull(),
						},
					),
				},
				Check: checkAggrAzure(projectID, comment, azureRegion,
					resource.TestCheckResourceAttr(pluralDSName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.#"),
					resource.TestCheckResourceAttrSet(pluralDSName, "id"),
				),
			},
			// TODO: Uncomment once CLOUDP-391704 is released — Azure comment update is not yet supported.
			// {
			// 	Config: dataSourceConfigBasicAzure(projectID, subscriptionID, clientID, clientSecret, tenantID, privateLinkServiceResourceID, resourceGroupName, "updated comment", azureRegion),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		checkExists(resourceName),
			// 		resource.TestCheckResourceAttr(resourceName, "comment", "updated comment"),
			// 	),
			// },
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
		return checkResourceOnlyAggr(projectID, endpointID,
			resource.TestCheckResourceAttr(resourceName, "comment", expectedComment),
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
		return checkResourceOnlyAggr(projectID, endpointID,
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

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_optionalAttrsOmittedAWS(t *testing.T) {
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
				Config: resourceConfigBasicAWSNoOptional(projectID, endpointID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("region"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("customer_endpoint_dns_name"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("comment"), knownvalue.StringExact("")),
				},
				Check: checkResourceOnlyAggr(projectID, endpointID,
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "region", ""),
					resource.TestCheckResourceAttr(resourceName, "customer_endpoint_dns_name", ""),
				),
			},
		},
	})
}

func TestAccNetworkPrivatelinkEndpointServiceDataFederationOnlineArchive_forceNewEndpointID(t *testing.T) {
	var (
		projectID         = acc.ProjectIDExecution(t)
		endpointID        = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID")
		endpointIDReplace = os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID_REPLACE")
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckPrivateEndpoint(t)
			if endpointIDReplace == "" {
				t.Fatal("`MONGODB_ATLAS_PRIVATE_ENDPOINT_ID_REPLACE` must be set for force-new endpoint_id acceptance testing")
			}
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasicAWS(projectID, endpointID, comment),
				Check: checkResourceOnlyAggr(projectID, endpointID,
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: resourceConfigBasicAWS(projectID, endpointIDReplace, comment),
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
							"region":                     knownvalue.StringExact(AWSRegion),
							"customer_endpoint_dns_name": knownvalue.StringExact(customerEndpointDNSName),
						},
					),
				},
				Check: checkAggr(projectID, endpointID, comment, AWSRegion, customerEndpointDNSName,
					checkEncodedID(dataSourceName, projectID, endpointID),
					resource.TestCheckResourceAttr(pluralDSName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(pluralDSName, "results.#"),
					resource.TestCheckResourceAttrSet(pluralDSName, "id"),
				),
			},
		},
	})
}

func checkAggr(projectID, endpointID, commentValue, region, customerEndpointDNSName string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attrsSet := []string{"type", "provider_name"}
	attrsMap := map[string]string{
		"project_id":                 projectID,
		"endpoint_id":                endpointID,
		"comment":                    commentValue,
		"region":                     region,
		"customer_endpoint_dns_name": customerEndpointDNSName,
	}
	extraChecks := extra
	extraChecks = append(extraChecks, checkExists(resourceName))
	return acc.CheckRSAndDS(resourceName, &dataSourceName, nil, attrsSet, attrsMap, extraChecks...)
}

func checkAggrAzure(projectID, commentValue, region string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	attrsSet := []string{"type", "endpoint_id", "customer_endpoint_ip_address"}
	attrsMap := map[string]string{
		"project_id":    projectID,
		"provider_name": "AZURE",
		"comment":       commentValue,
		"region":        region,
	}
	extraChecks := extra
	extraChecks = append(extraChecks, checkExists(resourceName))
	return acc.CheckRSAndDS(resourceName, &dataSourceName, nil, attrsSet, attrsMap, extraChecks...)
}

func checkResourceOnlyAggr(projectID, endpointID string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		checkEncodedID(resourceName, projectID, endpointID),
		resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
		resource.TestCheckResourceAttr(resourceName, "endpoint_id", endpointID),
		resource.TestCheckResourceAttrSet(resourceName, "type"),
		resource.TestCheckResourceAttrSet(resourceName, "provider_name"),
	}
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
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
		if rs.Type != "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" {
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
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id    = %[1]q
	  endpoint_id   = %[2]q
	  provider_name = "AWS"
	  comment       = %[3]q
	}
	`, projectID, endpointID, comment)
}

func resourceConfigBasicAWSLowercase(projectID, endpointID, comment string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id				= %[1]q
	  endpoint_id				= %[2]q
	  provider_name				= "aws"
	  comment					= %[3]q
	}
	`, projectID, endpointID, comment)
}

func resourceConfigBasicAWSWithTimeouts(projectID, endpointID, createTimeout, deleteTimeout string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
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

func resourceConfigBasicAWSNoOptional(projectID, endpointID string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id    = %q
  endpoint_id   = %q
  provider_name = "AWS"
}
`, projectID, endpointID)
}

func dataSourceConfigBasicAWS(projectID, endpointID, comment string) string {
	return buildConfigWithDataSource(projectID, endpointID, comment, "", false)
}

func dataSourceConfigOptionalFieldsAWS(projectID, endpointID, comment, customerEndpointDNSName string) string {
	return buildConfigWithDataSource(projectID, endpointID, comment, customerEndpointDNSName, true)
}

func dataSourceConfigBasicAzure(projectID, subscriptionID, clientID, clientSecret, tenantID, privateLinkServiceResourceID, resourceGroupName, comment, region string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "azurerm_resource_group" "test" {
	  name     = %[8]q
	  location = "East US 2"
	}

	resource "azurerm_virtual_network" "test" {
	  name                = "vnet-df-pe-test"
	  address_space       = ["10.0.0.0/16"]
	  location            = azurerm_resource_group.test.location
	  resource_group_name = azurerm_resource_group.test.name
	}

	resource "azurerm_subnet" "test" {
	  name                              = "snet-df-pe-test"
	  resource_group_name               = azurerm_resource_group.test.name
	  virtual_network_name              = azurerm_virtual_network.test.name
	  address_prefixes                  = ["10.0.1.0/24"]
	  private_endpoint_network_policies = "Disabled"
	}

	resource "azurerm_private_endpoint" "test" {
	  name                = "pe-df-test"
	  location            = azurerm_resource_group.test.location
	  resource_group_name = azurerm_resource_group.test.name
	  subnet_id           = azurerm_subnet.test.id

	  private_service_connection {
	    name                           = "atlas-df-connection"
	    private_connection_resource_id = %[7]q
	    is_manual_connection           = true
	    request_message                = "Terraform Acceptance Test"
	  }
	}

	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id                   = %[2]q
	  endpoint_id                  = azurerm_private_endpoint.test.id
	  provider_name                = "AZURE"
	  customer_endpoint_ip_address = azurerm_private_endpoint.test.private_service_connection[0].private_ip_address
	  comment                      = %[9]q
	  region                       = %[10]q
	}

	%[11]s
	%[12]s
	`,
		acc.ConfigAzurermProvider(subscriptionID, clientID, clientSecret, tenantID),
		projectID,
		subscriptionID,
		clientID,
		clientSecret,
		tenantID,
		privateLinkServiceResourceID,
		resourceGroupName,
		comment,
		region,
		singularDataSourceConfig,
		pluralDataSourceConfig,
	)
}

func buildConfigWithDataSource(projectID, endpointID, comment, customerEndpointDNSName string, includeOptionalFields bool) string {
	optionalFields := ""
	if includeOptionalFields {
		optionalFields = fmt.Sprintf(`
	  region                     = %q
	  customer_endpoint_dns_name = %q`, AWSRegion, customerEndpointDNSName)
	}

	return fmt.Sprintf(`
	resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
	  project_id    = %[1]q
	  endpoint_id   = %[2]q
	  provider_name = "AWS"
	  comment       = %[3]q
	  %[4]s
	}

	%[5]s
	%[6]s
	`, projectID, endpointID, comment, optionalFields, singularDataSourceConfig, pluralDataSourceConfig)
}
