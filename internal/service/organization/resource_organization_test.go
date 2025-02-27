package organization_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSOrganization_Basic(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = acc.RandomName()
		updatedName  = acc.RandomName()
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgOwnerID, name, description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "gen_ai_features_enabled", "true"),
				),
			},
			{
				Config: configBasic(orgOwnerID, updatedName, description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "gen_ai_features_enabled", "true"),
				),
			},
		},
	})
}

func TestAccConfigRSOrganization_BasicAccess(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		orgOwnerID  = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name        = acc.RandomName()
		description = "test Key for Acceptance tests"
		roleName    = "ORG_BILLING_ADMIN"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(orgOwnerID, name, description, roleName),
				ExpectError: regexp.MustCompile("API Key must have the ORG_OWNER role"),
			},
		},
	})
}

func TestAccConfigRSOrganization_Settings(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = acc.RandomName()
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"

		settingsConfig = `
		api_access_list_required = false
		multi_factor_auth_required = true
		gen_ai_features_enabled = false`
		settingsConfigUpdated = `
		api_access_list_required = false
		multi_factor_auth_required = true
		restrict_employee_access = true
		gen_ai_features_enabled = true`
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithSettings(orgOwnerID, name, description, roleName, settingsConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "true"),
				),
			},
			{
				Config: configWithSettings(orgOwnerID, name, description, roleName, settingsConfigUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_access_list_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "multi_factor_auth_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "restrict_employee_access", "true"),
				),
			},
			{
				Config: configBasic(orgOwnerID, "org-name-updated", description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccConfigRSOrganizationCreate_Errors(t *testing.T) {
	var (
		roleName    = "ORG_OWNER"
		unknownUser = "65def6160f722a1507105aaa"
	)
	acc.SkipTestForCI(t) // test will fail in CI since API_KEY_MUST_BE_ASSOCIATED_WITH_PAYING_ORG is returned
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(unknownUser, acc.RandomName(), "should fail since user is not found", roleName),
				ExpectError: regexp.MustCompile(`USER_NOT_FOUND`),
			},
		},
	})
}

func TestAccConfigDSOrganization_basic(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		datasourceName = "data.mongodbatlas_organization.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithDS(orgID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "restrict_employee_access"),
					resource.TestCheckResourceAttrSet(datasourceName, "multi_factor_auth_required"),
					resource.TestCheckResourceAttrSet(datasourceName, "api_access_list_required"),
					resource.TestCheckResourceAttr(datasourceName, "gen_ai_features_enabled", "true"),
				),
			},
		},
	})
}

func TestAccConfigDSOrganization_noAccessShouldFail(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configWithPluralDS() + acc.ConfigOrgMemberProvider(),
				ExpectError: regexp.MustCompile("error getting organization settings .*"),
			},
		},
	})
}

func TestAccConfigDSOrganizations_basic(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPluralDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.restrict_employee_access"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.multi_factor_auth_required"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.api_access_list_required"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.gen_ai_features_enabled", "true"),
				),
			},
		},
	})
}

func TestAccConfigDSOrganizations_withPagination(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithPagination(2, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
				),
			},
		},
	})
}

func configWithPluralDS() string {
	return `	
		data "mongodbatlas_organizations" "test" {
		}
	`
}

func configWithPagination(pageNum, itemPage int) string {
	return fmt.Sprintf(`
		data "mongodbatlas_organizations" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, pageNum, itemPage)
}

func configWithDS(orgID string) string {
	cfg := fmt.Sprintf(`
		
		data "mongodbatlas_organization" "test" {
			org_id = %[1]q
		}
	`, orgID)
	return cfg
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
		conn, err := getTestClientWithNewOrgCreds(rs)
		if err != nil {
			return err
		}

		orgs, _, err := conn.OrganizationsApi.ListOrganizations(context.Background()).Execute()
		if err == nil {
			for _, val := range orgs.GetResults() {
				if val.GetId() == ids["org_id"] {
					return nil
				}
			}
			return fmt.Errorf("Organization (%s) doesn't exist", ids["org_id"])
		}

		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		conn, err := getTestClientWithNewOrgCreds(rs)
		if err != nil {
			return err
		}

		orgs, _, err := conn.OrganizationsApi.ListOrganizations(context.Background()).Execute()
		if err == nil {
			for _, val := range orgs.GetResults() {
				if val.GetId() == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["org_id"])
				}
			}
			return nil
		}
	}

	return nil
}

func configBasic(orgOwnerID, name, description, roleNames string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = "%s"
		name = "%s"
		description = "%s"
		role_names = ["%s"]
	  }
	`, orgOwnerID, name, description, roleNames)
}

func configWithSettings(orgOwnerID, name, description, roleNames, settingsConfig string) string {
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

// getTestClientWithNewOrgCreds creates a new Atlas client with credentials for the newly created organization which
// is required to call relevant API methods for the new organization, for example ListOrganizations requires that the requesting API
// key must have the Organization Member role. So we cannot invoke API methods on the new organization with credentials configured in the provider.
func getTestClientWithNewOrgCreds(rs *terraform.ResourceState) (*admin.APIClient, error) {
	if rs.Primary.Attributes["public_key"] == "" {
		return nil, fmt.Errorf("no public_key is set")
	}

	if rs.Primary.Attributes["private_key"] == "" {
		return nil, fmt.Errorf("no private_key is set")
	}

	cfg := config.Config{
		PublicKey:  rs.Primary.Attributes["public_key"],
		PrivateKey: rs.Primary.Attributes["private_key"],
		BaseURL:    acc.MongoDBClient.Config.BaseURL,
	}

	ctx := context.Background()
	clients, _ := cfg.NewClient(ctx)
	return clients.(*config.MongoDBClient).AtlasV2, nil
}

func TestValidateAPIKeyIsOrgOwner(t *testing.T) {
	tests := []struct {
		name    string
		roles   []string
		wantErr bool
	}{
		{
			name:    "Contains OrgOwner",
			roles:   []string{"ORG_MEMBER", "ORG_OWNER", "ORG_READ_ONLY"},
			wantErr: false,
		},
		{
			name:    "Does Not Contain OrgOwner",
			roles:   []string{"ORG_MEMBER", "READ_ONLY"},
			wantErr: true,
		},
		{
			name:    "Empty Roles",
			roles:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := organization.ValidateAPIKeyIsOrgOwner(tt.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAPIKeyIsOrgOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
