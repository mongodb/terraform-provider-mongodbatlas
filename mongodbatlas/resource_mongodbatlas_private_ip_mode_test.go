package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourcePrivateIPMode_basic(t *testing.T) {
	var (
		privateIPMode matlas.PrivateIPMode
		resourceName  = "mongodbatlas_private_ip_mode.test"
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateIPModeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateIPModeConfig(projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateIPModeExists(resourceName, &privateIPMode),
					testAccCheckPrivateIPModeAttributes(&privateIPMode, pointy.Bool(false)),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPrivateIPModeExists(resourceName string, privateIPMode *matlas.PrivateIPMode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no project ID is set")
		}

		if privateIPModeResp, _, err := conn.PrivateIPMode.Get(context.Background(), rs.Primary.Attributes["project_id"]); err == nil {
			*privateIPMode = *privateIPModeResp
			enabled, _ := strconv.ParseBool(rs.Primary.Attributes["enabled"])
			privateIPMode.Enabled = pointy.Bool(enabled)

			return nil
		}

		return fmt.Errorf("privateIPMode(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["enabled"])
	}
}

func testAccCheckPrivateIPModeAttributes(privateIPMode *matlas.PrivateIPMode, enabled *bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *privateIPMode.Enabled != *enabled {
			return fmt.Errorf("bad enabled: %t", *privateIPMode.Enabled)
		}

		return nil
	}
}

func testAccCheckPrivateIPModeDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_private_ip_mode" {
			continue
		}

		privateIPMode, _, err := conn.PrivateIPMode.Get(context.Background(), rs.Primary.Attributes["project_id"])

		if err == nil && *privateIPMode.Enabled {
			return fmt.Errorf("privateIPMode from project (%s) still enabled", rs.Primary.Attributes["project_id"])
		}
	}

	return nil
}

func testAccMongoDBAtlasPrivateIPModeConfig(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_ip_mode" "test" {
			project_id  = "%s"
			enabled		= false
		}
	`, projectID)
}
