package organization2_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization2"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceType = "mongodbatlas_organization2"

func TestMain(m *testing.M) {
	os.Setenv("MONGODB_ATLAS_ORGANIZATION2_POC_STORE", filepath.Join(os.TempDir(), "mongodbatlas-organization2-acc-store.json"))
	os.Exit(acc.Run(m))
}

func TestAccOrganization2_noRotationBlock(t *testing.T) {
	acc.SkipInUnitTest(t)
	t.Cleanup(organization2.ResetStoreForTest)

	name := fmt.Sprintf("acc-no-rotation-%d", time.Now().UnixNano())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheckPoC(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization2,
		Steps: []resource.TestStep{
			{
				Config: configNoRotation(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "name", name),
					resource.TestCheckResourceAttrSet(resourceName(name), "org_id"),
					resource.TestCheckResourceAttrSet(resourceName(name), "client_id"),
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret"),
					resource.TestCheckNoResourceAttr(resourceName(name), "client_secret_rotation"),
				),
			},
		},
	})
}

func TestAccOrganization2_withRotationBlock(t *testing.T) {
	acc.SkipInUnitTest(t)
	t.Cleanup(organization2.ResetStoreForTest)

	name := fmt.Sprintf("acc-rotation-%d", time.Now().UnixNano())
	var firstCurrentSecretID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheckPoC(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization2,
		Steps: []resource.TestStep{
			{
				Config: configWithRotation(name, "2s", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "1"),
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret_rotation.current_secret_id"),
					saveAttr(resourceName(name), "client_secret_rotation.current_secret_id", &firstCurrentSecretID),
				),
			},
			{
				Config: configWithRotation(name, "2s", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "1"),
				),
			},
			{
				PreConfig: func() { time.Sleep(3 * time.Second) },
				Config:    configWithRotation(name, "2s", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "2"),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							resourceName(name),
							"client_secret_rotation.old_secret_id",
							firstCurrentSecretID,
						)(s)
					},
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret_rotation.current_secret_id"),
				),
			},
		},
	})
}

func TestAccOrganization2_forceSecretVersion(t *testing.T) {
	acc.SkipInUnitTest(t)
	t.Cleanup(organization2.ResetStoreForTest)

	name := fmt.Sprintf("acc-force-version-%d", time.Now().UnixNano())
	var firstCurrentSecretID string
	forcedVersion := int64(2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheckPoC(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization2,
		Steps: []resource.TestStep{
			{
				Config: configWithRotation(name, "240h", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "1"),
					saveAttr(resourceName(name), "client_secret_rotation.current_secret_id", &firstCurrentSecretID),
				),
			},
			{
				Config: configWithRotation(name, "240h", &forcedVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName(name), "client_secret_rotation.secret_version", "2"),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(
							resourceName(name),
							"client_secret_rotation.old_secret_id",
							firstCurrentSecretID,
						)(s)
					},
					resource.TestCheckResourceAttrSet(resourceName(name), "client_secret_rotation.current_secret_id"),
				),
			},
		},
	})
}

func preCheckPoC(t *testing.T) {
	t.Helper()
}

func resourceName(name string) string {
	return fmt.Sprintf("%s.%s", resourceType, name)
}

func configProvider() string {
	return `provider "mongodbatlas" {}`
}

func configNoRotation(name string) string {
	return fmt.Sprintf(`
%s

resource %q %q {
  name = %q
}
`, configProvider(), resourceType, name, name)
}

func configWithRotation(name, interval string, secretVersion *int64) string {
	versionAttr := ""
	if secretVersion != nil {
		versionAttr = fmt.Sprintf("\n    secret_version = %d", *secretVersion)
	}
	return fmt.Sprintf(`
%s

resource %q %q {
  name = %q

  client_secret_rotation = {
    interval = %q
%s
  }
}
`, configProvider(), resourceType, name, name, interval, versionAttr)
}

func checkDestroyOrganization2(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		if organization2.HasStoreEntry(rs.Primary.Attributes["name"]) {
			return fmt.Errorf("organization2 %q still exists in mock store", rs.Primary.Attributes["name"])
		}
	}
	return nil
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
