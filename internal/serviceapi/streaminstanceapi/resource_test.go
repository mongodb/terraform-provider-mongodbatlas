package streaminstanceapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName  = "mongodbatlas_stream_instance_api.test"
	region        = "VIRGINIA_USA"
	cloudProvider = "AWS"
)

func TestAccStreamInstanceAPI_basic(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, instanceName, region, cloudProvider, "SP10"),
				Check:  checkBasic(projectID, instanceName, region, cloudProvider, "SP10"),
			},
			// Same as in the curated resource, update can't be tested immediately after creation because the Atlas API doesn't expose a field to know when the resource is ready.
			/*
				{
					Config: configBasic(projectID, instanceName, region, cloudProvider, "SP30"),
					Check:  checkBasic(projectID, instanceName, region, cloudProvider, "SP30"),
				},*/
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func configBasic(projectID, instanceName, region, cloudProvider, tier string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_instance_api" "test" {
			group_id         = %[1]q
			name             = %[2]q
			region           = %[3]q
			cloud_provider   = %[4]q
			
			data_process_region = {
				region         = %[3]q
				cloud_provider = %[4]q
			}

			stream_config = {
				tier = %[5]q
			}
		}
	`, projectID, instanceName, region, cloudProvider, tier)
}

func checkBasic(projectID, instanceName, region, cloudProvider, tier string) resource.TestCheckFunc {
	mapChecks := map[string]string{
		"group_id":                           projectID,
		"name":                               instanceName,
		"region":                             region,
		"cloud_provider":                     cloudProvider,
		"data_process_region.region":         region,
		"data_process_region.cloud_provider": cloudProvider,
		"stream_config.tier":                 tier,
	}
	checks := acc.AddAttrChecks(resourceName, nil, mapChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		name := rs.Primary.Attributes["name"]
		if groupID == "" || name == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().StreamsApi.GetStreamInstance(context.Background(), groupID, name).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("stream instance(%s/%s) does not exist", groupID, name)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_instance_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		name := rs.Primary.Attributes["name"]
		if groupID == "" || name == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnV2().StreamsApi.GetStreamInstance(context.Background(), groupID, name).Execute()
		if err == nil {
			return fmt.Errorf("stream instance (%s/%s) still exists", groupID, name)
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
		groupID := rs.Primary.Attributes["group_id"]
		name := rs.Primary.Attributes["name"]
		if groupID == "" || name == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", groupID, name), nil
	}
}
