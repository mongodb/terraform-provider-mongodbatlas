package privateendpointconnectionstring_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceType         = "mongodbatlas_private_endpoint_connection_string"
	resourceName         = resourceType + ".test"
	dataSourceName       = "data." + resourceType + ".test"
	pluralDataSourceName = "data." + resourceType + "s.test"
	firstEndpointName    = "mongodbatlas_privatelink_endpoint_service.first"
	secondEndpointName   = "mongodbatlas_privatelink_endpoint_service.second"
	// awsRegion is the region of the AWS_VPC_ID, AWS_SUBNET_ID and AWS_SECURITY_GROUP_ID test network.
	awsRegion = "us-east-1"
	// atlasRegion matches awsRegion because the cluster must span the region of the selected private endpoints.
	atlasRegion = "US_EAST_1"
)

var (
	// dataSourceAttrs are the attributes that the singular data source shares with the resource.
	dataSourceAttrs = []string{
		"project_id", "cluster_name", "connection_string_id", "name", "private_endpoint_ids.#", "optimized_mode_enabled",
		"connection_string", "srv_connection_string", "status", "management_type", "slug", "created_at", "updated_at",
	}
	errOptimizedModeUpdate = regexp.MustCompile(`optimized_mode_enabled cannot be updated`)
)

// TestAccPrivateEndpointConnectionString_basic covers the resource and both data sources on a new replica set with two
// AWS private endpoints in its region: create, update of name and private_endpoint_ids in place, the create-only
// optimized_mode_enabled, and import.
func TestAccPrivateEndpointConnectionString_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{
			ReplicationSpecs: []acc.ReplicationSpecRequest{{Region: atlasRegion}},
		})
		infra       = configInfra(&clusterInfo, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET_ID"), os.Getenv("AWS_SECURITY_GROUP_ID"))
		name        = acc.RandomName()
		nameUpdated = name + "-updated"
		// connection_string_id, slug and the SRV connection string must not change when the string is updated in place.
		sameID   = statecheck.CompareValue(compare.ValuesSame())
		sameSlug = statecheck.CompareValue(compare.ValuesSame())
		sameSRV  = statecheck.CompareValue(compare.ValuesSame())
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAwsEnvPrivateLinkEndpointService(t); acc.SerialSleep(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				// Opting in before the cluster exists means Atlas creates no IMPORTED string for it.
				PreConfig: func() { enableConnectionStrings(t, clusterInfo.ProjectID) },
				Config:    configBasic(infra, &clusterInfo, name, firstEndpointName, ""),
				Check:     checkBasic(&clusterInfo, name, firstEndpointName),
				ConfigStateChecks: []statecheck.StateCheck{
					sameID.AddStateValue(resourceName, tfjsonpath.New("connection_string_id")),
					sameSlug.AddStateValue(resourceName, tfjsonpath.New("slug")),
					sameSRV.AddStateValue(resourceName, tfjsonpath.New("srv_connection_string")),
					checkPluralResult(name),
				},
			},
			{
				Config: configBasic(infra, &clusterInfo, nameUpdated, secondEndpointName, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkBasic(&clusterInfo, nameUpdated, secondEndpointName),
				ConfigStateChecks: []statecheck.StateCheck{
					sameID.AddStateValue(resourceName, tfjsonpath.New("connection_string_id")),
					sameSlug.AddStateValue(resourceName, tfjsonpath.New("slug")),
					sameSRV.AddStateValue(resourceName, tfjsonpath.New("srv_connection_string")),
					checkPluralResult(nameUpdated),
				},
			},
			{
				// The PATCH body does not accept optimizedModeEnabled, so a change must fail at plan time.
				Config:      configBasic(infra, &clusterInfo, nameUpdated, secondEndpointName, "optimized_mode_enabled = true"),
				PlanOnly:    true,
				ExpectError: errOptimizedModeUpdate,
			},
			{
				// Setting the default explicitly matches the state, so the plan is empty.
				Config:   configBasic(infra, &clusterInfo, nameUpdated, secondEndpointName, "optimized_mode_enabled = false"),
				PlanOnly: true,
			},
			{
				Config:                               configBasic(infra, &clusterInfo, nameUpdated, secondEndpointName, ""),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "connection_string_id",
				ImportStateVerifyIgnore:              []string{"delete_on_create_timeout"},
			},
		},
	})
}

// configInfra returns the cluster and two AWS private endpoints in the cluster region, registered in Atlas.
func configInfra(clusterInfo *acc.ClusterInfo, vpcID, subnetID, securityGroupID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region = %[1]q
			%[2]s
		}

		%[3]s

		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = %[4]q
			provider_name = "AWS"
			region        = %[1]q
		}

		resource "aws_vpc_endpoint" "test" {
			count              = 2
			vpc_id             = %[5]q
			subnet_ids         = [%[6]q]
			security_group_ids = [%[7]q]
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
		}

		resource "mongodbatlas_privatelink_endpoint_service" "first" {
			project_id          = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id     = mongodbatlas_privatelink_endpoint.test.private_link_id
			endpoint_service_id = aws_vpc_endpoint.test[0].id
			provider_name       = "AWS"
		}

		# Registered after the first endpoint so the two registrations do not change the private link at the same time.
		resource "mongodbatlas_privatelink_endpoint_service" "second" {
			project_id          = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id     = mongodbatlas_privatelink_endpoint.test.private_link_id
			endpoint_service_id = aws_vpc_endpoint.test[1].id
			provider_name       = "AWS"
			depends_on          = [mongodbatlas_privatelink_endpoint_service.first]
		}
	`, awsRegion, acc.ConfigDefaultTags(), clusterInfo.TerraformStr, clusterInfo.ProjectID, vpcID, subnetID, securityGroupID)
}

func configBasic(infra string, clusterInfo *acc.ClusterInfo, name, endpointName, extraAttrs string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_private_endpoint_connection_string" "test" {
			project_id           = %[2]q
			cluster_name         = %[3]s
			name                 = %[4]q
			private_endpoint_ids = [%[5]s.endpoint_service_id]
			%[6]s
		}

		data "mongodbatlas_private_endpoint_connection_string" "test" {
			project_id           = mongodbatlas_private_endpoint_connection_string.test.project_id
			cluster_name         = mongodbatlas_private_endpoint_connection_string.test.cluster_name
			connection_string_id = mongodbatlas_private_endpoint_connection_string.test.connection_string_id
		}

		data "mongodbatlas_private_endpoint_connection_strings" "test" {
			project_id      = mongodbatlas_private_endpoint_connection_string.test.project_id
			cluster_name    = mongodbatlas_private_endpoint_connection_string.test.cluster_name
			management_type = "SELF_MANAGED"
			status          = "AVAILABLE"
			depends_on      = [mongodbatlas_private_endpoint_connection_string.test]
		}
	`, infra, clusterInfo.ProjectID, clusterInfo.TerraformNameRef, name, endpointName, extraAttrs)
}

func checkBasic(clusterInfo *acc.ClusterInfo, name, endpointName string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		checkExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "project_id", clusterInfo.ProjectID),
		resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "private_endpoint_ids.#", "1"),
		resource.TestCheckTypeSetElemAttrPair(resourceName, "private_endpoint_ids.*", endpointName, "endpoint_service_id"),
		resource.TestCheckResourceAttr(resourceName, "optimized_mode_enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "status", "AVAILABLE"),
		resource.TestCheckNoResourceAttr(resourceName, "error_message"),
		resource.TestCheckResourceAttr(resourceName, "management_type", "SELF_MANAGED"),
		resource.TestMatchResourceAttr(resourceName, "connection_string_id", regexp.MustCompile(`^[a-f0-9]{24}$`)),
		resource.TestMatchResourceAttr(resourceName, "connection_string", regexp.MustCompile(`^mongodb://`)),
		resource.TestMatchResourceAttr(resourceName, "srv_connection_string", regexp.MustCompile(`^mongodb\+srv://`)),
		checkSRVContainsSlug(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "created_at"),
		resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
		resource.TestCheckResourceAttr(resourceName, "delete_on_create_timeout", "true"),
		resource.TestCheckTypeSetElemAttrPair(dataSourceName, "private_endpoint_ids.*", endpointName, "endpoint_service_id"),
	}
	for _, attr := range dataSourceAttrs {
		checks = append(checks, resource.TestCheckResourceAttrPair(dataSourceName, attr, resourceName, attr))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// checkPluralResult finds the string by name because the plural data source does not guarantee an order.
func checkPluralResult(name string) statecheck.StateCheck {
	return acc.PluralResultCheck(pluralDataSourceName, "name", knownvalue.StringExact(name), map[string]knownvalue.Check{
		"connection_string_id":   knownvalue.StringRegexp(regexp.MustCompile(`^[a-f0-9]{24}$`)),
		"management_type":        knownvalue.StringExact("SELF_MANAGED"),
		"status":                 knownvalue.StringExact("AVAILABLE"),
		"optimized_mode_enabled": knownvalue.Bool(false),
		"private_endpoint_ids":   knownvalue.SetSizeExact(1),
		"srv_connection_string":  knownvalue.StringRegexp(regexp.MustCompile(`^mongodb\+srv://`)),
	})
}

// checkSRVContainsSlug verifies that the SRV hostname is the one Atlas generated for this string, which embeds its slug.
func checkSRVContainsSlug(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		slug, srv := rs.Primary.Attributes["slug"], rs.Primary.Attributes["srv_connection_string"]
		if slug == "" || !strings.Contains(srv, slug) {
			return fmt.Errorf("srv_connection_string %q does not contain slug %q", srv, slug)
		}
		return nil
	}
}

func checkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if _, err := getConnectionString(rs); err != nil {
			return fmt.Errorf("private endpoint connection string does not exist: %w", err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		if _, err := getConnectionString(rs); err == nil {
			return fmt.Errorf("private endpoint connection string %s still exists", rs.Primary.Attributes["connection_string_id"])
		}
	}
	return nil
}

// getConnectionString uses the preview SDK because the API is only available in the preview version.
func getConnectionString(rs *terraform.ResourceState) (*adminpreview.PrivateEndpointConnectionStringResponse, error) {
	result, _, err := acc.ConnPreview().ClustersAPI.GetEndpointConnectionString(
		context.Background(),
		rs.Primary.Attributes["project_id"],
		rs.Primary.Attributes["cluster_name"],
		rs.Primary.Attributes["connection_string_id"],
	).Execute()
	return result, err
}

// enableConnectionStrings opts the project in to private endpoint connection strings, which Atlas requires before it
// accepts a self-managed string. The project is deleted after the test run, so the opt-in is not reverted.
func enableConnectionStrings(tb testing.TB, projectID string) {
	tb.Helper()
	settings := &adminpreview.GroupSettingsPreview{IsPrivateEndpointConnectionStringsEnabled: new(true)}
	result, _, err := acc.ConnPreview().ProjectsAPI.UpdateGroupSettings(context.Background(), projectID, settings).Execute()
	require.NoError(tb, err, "failed to enable private endpoint connection strings in project %s", projectID)
	require.True(tb, result.GetIsPrivateEndpointConnectionStringsEnabled(), "private endpoint connection strings not enabled in project %s", projectID)
}

func importStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["cluster_name"],
			rs.Primary.Attributes["connection_string_id"],
		), nil
	}
}
