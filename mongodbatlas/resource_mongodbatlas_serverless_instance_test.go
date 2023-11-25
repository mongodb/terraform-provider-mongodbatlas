package mongodbatlas_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccServerlessInstance_basic(t *testing.T) {
	var (
		serverlessInstance      matlas.Cluster
		resourceName            = "mongodbatlas_serverless_instance.test"
		instanceName            = acctest.RandomWithPrefix("test-acc-serverless")
		orgID                   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName             = acctest.RandomWithPrefix("test-acc-serverless")
		datasourceName          = "data.mongodbatlas_serverless_instance.test"
		datasourceInstancesName = "data.mongodbatlas_serverless_instances.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfig(orgID, projectName, instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "state_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "create_date"),
					resource.TestCheckResourceAttrSet(datasourceName, "mongo_db_version"),
					resource.TestCheckResourceAttrSet(datasourceName, "continuous_backup_enabled"),
					resource.TestCheckResourceAttrSet(datasourceName, "termination_protection_enabled"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.0.id"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.0.state_name"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.0.continuous_backup_enabled"),
					resource.TestCheckResourceAttrSet(datasourceInstancesName, "results.0.termination_protection_enabled"),
					testAccCheckConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName),
				),
			},
		},
	})
}

func TestAccServerlessInstance_WithTags(t *testing.T) {
	var (
		serverlessInstance      matlas.Cluster
		resourceName            = "mongodbatlas_serverless_instance.test"
		instanceName            = acctest.RandomWithPrefix("test-acc-serverless")
		orgID                   = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName             = acctest.RandomWithPrefix("test-acc-serverless")
		dataSourceName          = "data.mongodbatlas_serverless_instance.test"
		dataSourceInstancesName = "data.mongodbatlas_serverless_instances.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfigWithTags(orgID, projectName, instanceName, []matlas.Tag{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceInstancesName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfigWithTags(orgID, projectName, instanceName, []matlas.Tag{
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceInstancesName, "results.0.tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceInstancesName, "results.0.tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceInstancesName, "results.0.tags.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfigWithTags(orgID, projectName, instanceName, []matlas.Tag{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName, &serverlessInstance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceInstancesName, "results.0.tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceInstancesName, "results.0.tags.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
}

func TestAccServerlessInstance_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_serverless_instance.test"
		instanceName = acctest.RandomWithPrefix("test-acc-serverless")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-serverless")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceConfig(orgID, projectName, instanceName, true),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasServerlessInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasServerlessInstanceExists(resourceName string, serverlessInstance *matlas.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := config.DecodeStateID(rs.Primary.ID)

		serverlessResponse, _, err := conn.ServerlessInstances.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			*serverlessInstance = *serverlessResponse
			return nil
		}

		return fmt.Errorf("serverless instance (%s) does not exist", ids["name"])
	}
}

func testAccCheckMongoDBAtlasServerlessInstanceDestroy(state *terraform.State) error {
	conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_serverless_instance" {
			continue
		}

		ids := config.DecodeStateID(rs.Primary.ID)

		serverlessInstance, _, err := conn.ServerlessInstances.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil && serverlessInstance != nil {
			return fmt.Errorf("serverless instance (%s) still exists", ids["name"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasServerlessInstanceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := config.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["name"]), nil
	}
}

func testAccCheckConnectionStringPrivateEndpointIsPresentWithNoElement(resourceName string) resource.TestCheckFunc {
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

func testAccMongoDBAtlasServerlessInstanceConfig(orgID, projectName, name string, ignoreConnectionStrings bool) string {
	lifecycle := ""

	if ignoreConnectionStrings {
		lifecycle = `

		lifecycle {
			ignore_changes = [connection_strings_private_endpoint_srv]
		}
		`
	}

	return fmt.Sprintf(serverlessConfig, orgID, projectName, name, lifecycle)
}

const serverlessConfig = `
	resource "mongodbatlas_project" "test" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q

		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true
		%[4]s
	}

	data "mongodbatlas_serverless_instance" "test" {
		name        = mongodbatlas_serverless_instance.test.name
		project_id  = mongodbatlas_serverless_instance.test.project_id
	}

	data "mongodbatlas_serverless_instances" "test" {
		project_id         = mongodbatlas_serverless_instance.test.project_id
	}
`

func testAccMongoDBAtlasServerlessInstanceConfigWithTags(orgID, projectName, name string, tags []matlas.Tag) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}
	return fmt.Sprintf(serverlessConfig, orgID, projectName, name, tagsConf)
}
