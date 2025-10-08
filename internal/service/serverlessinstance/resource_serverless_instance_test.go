package serverlessinstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

const (
	resourceName         = "mongodbatlas_serverless_instance.test"
	dataSourceName       = "data.mongodbatlas_serverless_instance.test"
	dataSourcePluralName = "data.mongodbatlas_serverless_instances.test"
)

func TestAccServerlessInstance_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccServerlessInstance_withTags(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, false, nil, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, false, nil, []admin.ResourceTag{
					{
						Key:   "key 1",
						Value: "value 1",
					},
					{
						Key:   "key 2",
						Value: "value 2",
					},
				},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, false, nil, []admin.ResourceTag{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
}

func TestAccServerlessInstance_autoIndexing(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, false, conversion.Pointer(false), nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_indexing", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_indexing", "false"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.auto_indexing"),
				),
			},
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, false, conversion.Pointer(true), nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_indexing", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_indexing", "true"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.auto_indexing"),
				),
			},
		},
	})
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		projectID    = acc.ProjectIDExecution(tb)
		instanceName = acc.RandomClusterName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstance(projectID, instanceName, true, nil, nil),
				Check:  resource.ComposeAggregateTestCheckFunc(basicChecks(projectID, instanceName)...),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func basicChecks(projectID, instanceName string) []resource.TestCheckFunc {
	commonChecks := map[string]string{
		"name":                           instanceName,
		"project_id":                     projectID,
		"termination_protection_enabled": "false",
		"continuous_backup_enabled":      "true",
	}
	commonSetChecks := []string{"state_name", "create_date", "mongo_db_version"}
	pluralSetChecks := []string{
		"project_id",
		"results.#",
		"results.0.id",
		"results.0.name",
		"results.0.state_name",
		"results.0.continuous_backup_enabled",
		"results.0.termination_protection_enabled",
	}

	checks := acc.AddAttrChecks(resourceName, nil, commonChecks)
	checks = acc.AddAttrChecks(dataSourceName, checks, commonChecks)
	checks = acc.AddAttrSetChecks(resourceName, checks, commonSetChecks...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, commonSetChecks...)
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, pluralSetChecks...)

	checks = append(checks, checkExists(resourceName), checkConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName))
	return checks
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
		_, _, err := acc.ConnV2().ServerlessInstancesApi.GetServerlessInstance(context.Background(), ids["project_id"], ids["name"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("serverless instance (%s) does not exist", ids["name"])
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_serverless_instance" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		serverlessInstance, _, err := acc.ConnV2().ServerlessInstancesApi.GetServerlessInstance(context.Background(), ids["project_id"], ids["name"]).Execute()
		if err == nil && serverlessInstance != nil {
			return fmt.Errorf("serverless instance (%s) still exists", ids["name"])
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

		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["name"]), nil
	}
}

func checkConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if connectionStringPrivateEndpoint := rs.Primary.Attributes["connection_strings_private_endpoint_srv.#"]; connectionStringPrivateEndpoint == "" {
			return fmt.Errorf("expected connection_strings_private_endpoint_srv to be present")
		}

		return nil
	}
}
