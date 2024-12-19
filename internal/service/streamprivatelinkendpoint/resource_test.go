package streamprivatelinkendpoint_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceType         = "mongodbatlas_stream_privatelink_endpoint"
	resourceName         = "mongodbatlas_stream_privatelink_endpoint.test"
	dataSourceName       = "data.mongodbatlas_stream_privatelink_endpoint.test"
	dataSourcePluralName = "data.mongodbatlas_stream_privatelink_endpoints.test"
)

func TestAccStreamPrivatelinkEndpoint_basic(t *testing.T) {
	tc := basicTestCase(t)
	// Tests include testing of plural data source and so cannot be run in parallel
	resource.Test(t, *tc)
}

func TestAccStreamPrivatelinkEndpoint_failedUpdate(t *testing.T) {
	tc := failedUpdateTestCase(t)
	// Tests include testing of plural data source and so cannot be run in parallel
	resource.Test(t, *tc)
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	var (
		// need specific projectID because feature is currently under a Feature flag
		projectID    = os.Getenv("MONGODB_ATLAS_STREAM_AWS_PL_PROJECT_ID")
		provider     = "AWS"
		region       = "us-east-1"
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyConfluent(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configConfluentDedicatedCluster(provider, region, awsAccountID) + configBasic(projectID, provider, region, vendor, true),
				Check:  checksStreamPrivatelinkEndpoint(projectID, provider, region, vendor, false),
			},
			{
				Config:            configBasic(projectID, provider, region, vendor, true),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func failedUpdateTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	var (
		// need specific projectID because feature is currently under a Feature flag
		projectID    = os.Getenv("MONGODB_ATLAS_STREAM_AWS_PL_PROJECT_ID")
		provider     = "AWS"
		region       = "us-east-1"
		vendor       = "CONFLUENT"
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configConfluentDedicatedCluster(provider, region, awsAccountID) + configBasic(projectID, provider, region, vendor, false),
				Check:  checksStreamPrivatelinkEndpoint(projectID, provider, region, vendor, false),
			},
			{
				Config:      configConfluentDedicatedCluster(provider, region, awsAccountID) + configBasic(projectID, provider, region, vendor, true),
				ExpectError: regexp.MustCompile(`Operation not supported`),
			},
		},
	}
}

func configBasic(projectID, provider, region, vendor string, withDNSSubdomains bool) string {
	dnsSubDomainConfig := ""
	if withDNSSubdomains {
		dnsSubDomainConfig = `dns_sub_domain = local.dns_sub_domain_entries`
	}

	return fmt.Sprintf(`
	locals {
		dns_sub_domain_entries = [
    		for zone in confluent_network.private-link.zones :
    		"${zone}.${confluent_network.private-link.dns_domain}"
  		]
	}	

	resource "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id          = %[1]q
		dns_domain          = confluent_network.private-link.dns_domain
		provider_name       = %[2]q
		region              = %[3]q
		vendor              = %[4]q
		service_endpoint_id = confluent_network.private-link.aws[0].private_link_endpoint_service
		%[5]s
		depends_on = [
			confluent_kafka_cluster.dedicated
		]
	}

	data "mongodbatlas_stream_privatelink_endpoint" "test" {
		project_id = %[1]q
		id         = mongodbatlas_stream_privatelink_endpoint.test.id
		depends_on = [
    		mongodbatlas_stream_privatelink_endpoint.test
  		]
	}

	data "mongodbatlas_stream_privatelink_endpoints" "test" {
		project_id = %[1]q
		depends_on = [
    		mongodbatlas_stream_privatelink_endpoint.test
  		]
	}`, projectID, provider, region, vendor, dnsSubDomainConfig)
}

func configConfluentDedicatedCluster(provider, region, awsAccountID string) string {
	return fmt.Sprintf(`
	%[1]s

	data "confluent_environment" "default_environment" {
  		display_name = "default"
	}

	resource "confluent_network" "private-link" {
		display_name     = "terraform-test-private-link-network"
		cloud            = %[2]q
		region           = %[3]q
		connection_types = ["PRIVATELINK"]
		zones            = ["use1-az1", "use1-az4", "use1-az6"]
		environment {
			id = data.confluent_environment.default_environment.id
		}
		dns_config {
			resolution = "PRIVATE"
		}
	}

	resource "confluent_private_link_access" "aws" {
		display_name = "terraform-test-aws-private-link-access"
		aws {
			account = %[4]q
		}
		environment {
			id = data.confluent_environment.default_environment.id
		}
		network {
			id = confluent_network.private-link.id
		}
	}

	resource "confluent_kafka_cluster" "dedicated" {
		display_name = "terraform-test-cluster"
		availability = "MULTI_ZONE"
		cloud        = confluent_network.private-link.cloud
		region       = confluent_network.private-link.region
		dedicated {
			cku = 2
		}
		environment {
			id = data.confluent_environment.default_environment.id
		}
		network {
			id = confluent_network.private-link.id
		}
	}`, acc.ConfigConfluentProvider(), provider, region, awsAccountID)
}

func checksStreamPrivatelinkEndpoint(projectID, provider, region, vendor string, dnsSubdomainsCheck bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists()}
	attrMap := map[string]string{
		"project_id":    projectID,
		"provider_name": provider,
		"region":        region,
		"vendor":        vendor,
	}
	pluralMap := map[string]string{
		"project_id": projectID,
		"results.#":  "1",
	}
	attrSet := []string{
		"id",
		"interface_endpoint_id",
		"state",
		"dns_domain",
		"service_endpoint_id",
	}
	if dnsSubdomainsCheck {
		attrSet = append(attrSet, "dns_sub_domain.0")
	}
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, pluralMap)
	return acc.CheckRSAndDS(resourceName, &dataSourceName, &dataSourcePluralName, attrSet, attrMap, checks...)
}

func checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == resourceType {
				projectID := rs.Primary.Attributes["project_id"]
				connectionID := rs.Primary.Attributes["id"]
				_, _, err := acc.ConnV2().StreamsApi.GetPrivateLinkConnection(context.Background(), projectID, connectionID).Execute()
				if err != nil {
					return fmt.Errorf("Privatelink Connection (%s:%s) not found", projectID, connectionID)
				}
			}
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceType {
			projectID := rs.Primary.Attributes["project_id"]
			connectionID := rs.Primary.Attributes["id"]
			_, _, err := acc.ConnV2().StreamsApi.GetPrivateLinkConnection(context.Background(), projectID, connectionID).Execute()
			if err == nil {
				return fmt.Errorf("Privatelink Connection (%s:%s) still exists", projectID, id)
			}
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]), nil
	}
}
