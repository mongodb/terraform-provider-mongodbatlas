package serverlessinstance_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

const (
	resourceName         = "mongodbatlas_serverless_instance.test"
	dataSourceName       = "data.mongodbatlas_serverless_instance.test"
	dataSourcePluralName = "data.mongodbatlas_serverless_instances.test"
)

func TestAccServerlessInstance_basic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		instanceName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstanceBasic(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					checkConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName),
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "create_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mongo_db_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "continuous_backup_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "termination_protection_enabled"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.state_name"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.continuous_backup_enabled"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.termination_protection_enabled"),
				),
			},
		},
	})
}

func TestAccServerlessInstance_WithTags(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		instanceName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstanceWithTags(orgID, projectName, instanceName, []admin.ResourceTag{}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: acc.ConfigServerlessInstanceWithTags(orgID, projectName, instanceName, []admin.ResourceTag{
					{
						Key:   conversion.StringPtr("key 1"),
						Value: conversion.StringPtr("value 1"),
					},
					{
						Key:   conversion.StringPtr("key 2"),
						Value: conversion.StringPtr("value 2"),
					},
				},
				),
				Check: resource.ComposeTestCheckFunc(
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
				Config: acc.ConfigServerlessInstanceWithTags(orgID, projectName, instanceName, []admin.ResourceTag{
					{
						Key:   conversion.StringPtr("key 3"),
						Value: conversion.StringPtr("value 3"),
					},
				},
				),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccServerlessInstance_importBasic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		instanceName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigServerlessInstanceBasic(orgID, projectName, instanceName, true),
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
