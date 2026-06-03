package organization3_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceType = "mongodbatlas_organization3"

func TestAccOrganization3_withRotationBlock(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.SkipUnlessHasOrgOwner(t)

	orgOwnerID := os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
	name := acc.RandomName()
	var firstCurrentSecretID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization3,
		Steps: []resource.TestStep{
			{
				Config: configWithRotation(name, orgOwnerID, 720, 87600, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "1"),
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret_rotation.current_secret.secret_id"),
					saveAttr(resourceName(name), "client_secret_rotation.current_secret.secret_id", &firstCurrentSecretID),
				),
			},
			{
				Config: configWithRotation(name, orgOwnerID, 720, 87600, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "2"),
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.old_secret.secret_id", firstCurrentSecretID),
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret_rotation.current_secret.secret_id"),
				),
			},
		},
	})
}

func TestAccOrganization3_rotationDeletesOldSecret(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.SkipUnlessHasOrgOwner(t)

	orgOwnerID := os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
	name := acc.RandomName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization3,
		Steps: []resource.TestStep{
			{Config: configWithRotation(name, orgOwnerID, 720, 87600, nil)},
			{Config: configWithRotation(name, orgOwnerID, 720, 87600, nil)},
			{
				Config: configWithRotation(name, orgOwnerID, 720, 87600, nil),
				Check:  checkAtMostTwoSecrets(resourceName(name)),
			},
		},
	})
}

func TestAccOrganization3_forceSecretVersion(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.SkipUnlessHasOrgOwner(t)

	orgOwnerID := os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
	name := acc.RandomName()
	var firstCurrentSecretID string
	forcedVersion := int64(2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization3,
		Steps: []resource.TestStep{
			{
				Config: configWithRotation(name, orgOwnerID, 720, 360, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "1"),
					saveAttr(resourceName(name), "client_secret_rotation.current_secret.secret_id", &firstCurrentSecretID),
				),
			},
			{
				Config: configWithRotation(name, orgOwnerID, 720, 360, &forcedVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "2"),
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.old_secret.secret_id", firstCurrentSecretID),
				),
			},
		},
	})
}

func resourceName(name string) string {
	return fmt.Sprintf("%s.%s", resourceType, name)
}

func configProvider() string {
	return `provider "mongodbatlas" {}`
}

func configWithRotation(name, orgOwnerID string, expiresAfter, rotateBefore int64, secretVersion *int64) string {
	versionAttr := ""
	if secretVersion != nil {
		versionAttr = fmt.Sprintf("\n    secret_version = %d", *secretVersion)
	}
	return fmt.Sprintf(`
%s

resource %q %q {
  name         = %q
  org_owner_id = %q

  client_secret_rotation = {
    expires_after_hours        = %d
    rotate_before_expiry_hours = %d
%s
  }
}
`, configProvider(), resourceType, name, name, orgOwnerID, expiresAfter, rotateBefore, versionAttr)
}

func checkDestroyOrganization3(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		_, resp, err := acc.MongoDBClient.AtlasV2.OrganizationsApi.GetOrg(context.Background(), orgID).Execute()
		if err == nil {
			return fmt.Errorf("organization3 org %q still exists", orgID)
		}
		if validate.StatusNotFound(resp) {
			return nil
		}
	}
	return nil
}

func checkAtMostTwoSecrets(resourceAddress string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceAddress)
		}
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		sa, _, err := acc.MongoDBClient.AtlasV2.ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
		if err != nil {
			return err
		}
		if len(sa.GetSecrets()) > 2 {
			return fmt.Errorf("expected at most 2 secrets, got %d", len(sa.GetSecrets()))
		}
		return nil
	}
}

func saveAttr(resourceAddress, attr string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceAddress)
		}
		value, ok := rs.Primary.Attributes[attr]
		if !ok {
			return fmt.Errorf("attribute %q not found on %s", attr, resourceAddress)
		}
		*target = value
		return nil
	}
}
