package streaminstance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamRSStreamInstance_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_stream_instance.test"
		dataSourceName = "data.mongodbatlas_stream_instance.test"
		pluralDSName   = "data.mongodbatlas_stream_instances.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstanceWithDataSourcesConfig(projectID, instanceName, region, cloudProvider),
				Check:  streamInstanceChecksWithDataSources(resourceName, dataSourceName, pluralDSName, instanceName, "SP30"),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStreamRSStreamInstance_withStreamConfig(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_stream_instance.test"
		dataSourceName = "data.mongodbatlas_stream_instance.test"
		pluralDSName   = "data.mongodbatlas_stream_instances.test"
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamInstanceWithDataSourcesAndStreamConfigConfig(projectID, instanceName, region, cloudProvider, "SP10"),
				Check:  streamInstanceChecksWithDataSources(resourceName, dataSourceName, pluralDSName, instanceName, "SP10"),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamInstanceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func streamInstanceAttributeChecks(resourceName, instanceName, region, cloudProvider string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkSearchInstanceExists(),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "instance_name", instanceName),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.region", region),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.cloud_provider", cloudProvider),
		resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func checkStreamInstanceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"]), nil
	}
}

func checkSearchInstanceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "mongodbatlas_stream_instance" {
				_, _, err := acc.ConnV2().StreamsApi.GetStreamWorkspace(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"]).Execute()
				if err != nil {
					return fmt.Errorf("stream instance (%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["instance_name"])
				}
			}
		}
		return nil
	}
}

func streamInstanceChecksWithDataSources(resourceName, dataSourceName, pluralDataSourceName, instanceName, tier string) resource.TestCheckFunc {
	attrsMap := map[string]string{
		"instance_name":                      instanceName,
		"data_process_region.region":         region,
		"data_process_region.cloud_provider": cloudProvider,
		"stream_config.tier":                 tier,
	}
	attrsSet := []string{"id", "project_id", "hostnames.#"}
	return resource.ComposeAggregateTestCheckFunc(
		checkSearchInstanceExists(),
		acc.CheckRSAndDS(resourceName, &dataSourceName, nil, attrsSet, attrsMap),
		resource.TestCheckResourceAttrSet(pluralDataSourceName, "project_id"),
		resource.TestCheckResourceAttrSet(pluralDataSourceName, "total_count"),
		resource.TestCheckResourceAttrWith(pluralDataSourceName, "results.#", acc.IntGreatThan(0)),
		resource.TestCheckTypeSetElemNestedAttrs(pluralDataSourceName, "results.*", map[string]string{
			"instance_name": instanceName,
		}),
	)
}

func streamInstanceWithDataSourcesConfig(projectID, instanceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instance" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		}

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
		}
	`, acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider))
}

func streamInstanceWithDataSourcesAndStreamConfigConfig(projectID, instanceName, region, cloudProvider, tier string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_instance" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
			instance_name = mongodbatlas_stream_instance.test.instance_name
		}

		data "mongodbatlas_stream_instances" "test" {
			project_id = mongodbatlas_stream_instance.test.project_id
		}
	`, acc.StreamInstanceWithStreamConfigConfig(projectID, instanceName, region, cloudProvider, tier))
}
