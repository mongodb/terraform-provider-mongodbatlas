package organization_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSOrganization_Basic(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-organization-%s", acctest.RandString(5))
		updatedName  = fmt.Sprintf("test-acc-organization-updated-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, updatedName, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccConfigRSOrganization_BasicAccess(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		orgOwnerID  = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name        = fmt.Sprintf("test-acc-organization-%s", acctest.RandString(5))
		description = "test Key for Acceptance tests"
		roleName    = "ORG_BILLING_ADMIN"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
				ExpectError: regexp.MustCompile("API Key must have the ORG_OWNER role"),
			},
		},
	})
}

func TestAccConfigRSOrganization_Settings(t *testing.T) {
	acc.SkipTestForCI(t)
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-organization-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"

		settingsConfig = `
		api_access_list_required = false
		multi_factor_auth_required = true`
		settingsConfigUpdated = `
		api_access_list_required = false
		multi_factor_auth_required = true
		restrict_employee_access = true`
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigWithSettings(orgOwnerID, name, description, roleName, settingsConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "true"),
				),
			},
			{
				Config: testAccMongoDBAtlasOrganizationConfigWithSettings(orgOwnerID, name, description, roleName, settingsConfigUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "true"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasOrganizationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if rs.Primary.Attributes["public_key"] == "" {
			return fmt.Errorf("no public_key is set")
		}

		if rs.Primary.Attributes["private_key"] == "" {
			return fmt.Errorf("no private_key is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		cfg := config.Config{
			PublicKey:  rs.Primary.Attributes["public_key"],
			PrivateKey: rs.Primary.Attributes["private_key"],
			BaseURL:    acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Config.BaseURL,
		}

		ctx := context.Background()
		clients, _ := cfg.NewClient(ctx)
		conn := clients.(*config.MongoDBClient).AtlasV2

		orgs, _, err := conn.OrganizationsApi.ListOrganizations(ctx).Execute()
		if err == nil {
			for _, val := range *orgs.Results {
				if *val.Id == ids["org_id"] {
					return nil
				}
			}
			return fmt.Errorf("Organization (%s) doesn't exist", ids["org_id"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasOrganizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		if rs.Primary.Attributes["public_key"] == "" {
			return fmt.Errorf("no public_key is set")
		}

		if rs.Primary.Attributes["private_key"] == "" {
			return fmt.Errorf("no private_key is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		cfg := config.Config{
			PublicKey:  rs.Primary.Attributes["public_key"],
			PrivateKey: rs.Primary.Attributes["private_key"],
			BaseURL:    acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Config.BaseURL,
		}

		ctx := context.Background()
		clients, _ := cfg.NewClient(ctx)
		conn := clients.(*config.MongoDBClient).AtlasV2

		orgs, _, err := conn.OrganizationsApi.ListOrganizations(context.Background()).Execute()
		if err == nil {
			for _, val := range *orgs.Results {
				if *val.Id == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["org_id"])
				}
			}
			return nil
		}
	}

	return nil
}

func testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleNames string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = "%s"
		name = "%s"
		description = "%s"
		role_names = ["%s"]
	  }
	`, orgOwnerID, name, description, roleNames)
}

func testAccMongoDBAtlasOrganizationConfigWithSettings(orgOwnerID, name, description, roleNames, settingsConfig string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = "%s"
		name = "%s"
		description = "%s"
		role_names = ["%s"]
		%s
	  }
	`, orgOwnerID, name, description, roleNames, settingsConfig)
}
