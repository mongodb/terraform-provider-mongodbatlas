package mongodbatlas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource mongodbatlas_example "test" {
					configurable_attribute = "config_attr_val"
				}`,
			},
		},
	})
}
