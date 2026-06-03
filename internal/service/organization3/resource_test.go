package organization3_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceType = "mongodbatlas_organization3"

// TestAccOrganization3_rotationLifecycle exercises create, practitioner-forced rotation,
// stable re-apply, and ModifyPlan-scheduled rotation from a widened rotate_before_expiry_hours policy.
// Short expires_after_hours (8h) keeps Atlas overlap windows testable without long sleeps.
func TestAccOrganization3_rotationLifecycle(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.SkipUnlessHasOrgOwner(t)

	orgOwnerID := os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
	name := acc.RandomName()
	addr := resourceName(name)
	var firstCurrentSecretID, secondCurrentSecretID string
	forcedVersion := int64(2)
	expiresAfter := int64(8) // Atlas secret lifetime; also default rotate_before = 4h
	rotateBefore := int64(8) // renew window covers full lifetime so ModifyPlan schedules immediately

	createConfig := configRotationBlock(name, orgOwnerID, rotationBlockConfig{
		expiresAfterHours: expiresAfter,
	})
	afterForceConfig := configRotationBlock(name, orgOwnerID, rotationBlockConfig{
		expiresAfterHours: expiresAfter,
	}) // same as create; secret_version removed after forced rotation
	widenRotateBeforeConfig := configRotationBlock(name, orgOwnerID, rotationBlockConfig{
		expiresAfterHours:       expiresAfter,
		rotateBeforeExpiryHours: &rotateBefore,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrganization3,
		Steps: []resource.TestStep{
			// Step 1: Create org + SA with rotation block; initial secret is version 1.
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "client_secret_rotation.secret_version", "1"),
					resource.TestCheckResourceAttrSet(addr, "client_secret_rotation.current_secret.secret_id"),
					saveAttr(addr, "client_secret_rotation.current_secret.secret_id", &firstCurrentSecretID),
				),
			},
			// Step 2: Re-apply unchanged config; expect no drift (ModifyPlan must not schedule yet).
			{
				Config: createConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Practitioner forces rotation via secret_version = 2 (no wait for expires_at).
			{
				Config: configRotationBlock(name, orgOwnerID, rotationBlockConfig{
					expiresAfterHours: expiresAfter,
					secretVersion:     &forcedVersion,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "client_secret_rotation.secret_version", "2"),
					resource.TestCheckResourceAttr(addr, "client_secret_rotation.old_secret.secret_id", firstCurrentSecretID),
					resource.TestCheckResourceAttrSet(addr, "client_secret_rotation.current_secret.secret_id"),
					saveAttr(addr, "client_secret_rotation.current_secret.secret_id", &secondCurrentSecretID),
				),
			},
			// Step 4: Drop secret_version from config; state stays at 2 with two secrets, no further change.
			{
				Config: afterForceConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 5: Widen rotate_before_expiry_hours to 8 so renewAt is in the past; ModifyPlan → version 3.
			// At state version 2 the provider deletes old_secret before POST (deletion policy).
			{
				Config: widenRotateBeforeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "client_secret_rotation.secret_version", "3"),
					resource.TestCheckResourceAttr(addr, "client_secret_rotation.old_secret.secret_id", secondCurrentSecretID),
					resource.TestCheckResourceAttrSet(addr, "client_secret_rotation.current_secret.secret_id"),
				),
			},
		},
	})
}

func resourceName(name string) string {
	return fmt.Sprintf("%s.%s", resourceType, name)
}

type rotationBlockConfig struct {
	expiresAfterHours       int64
	rotateBeforeExpiryHours *int64
	secretVersion           *int64
}

func configRotationBlock(name, orgOwnerID string, cfg rotationBlockConfig) string {
	extra := ""
	if cfg.rotateBeforeExpiryHours != nil {
		extra += fmt.Sprintf("\n    rotate_before_expiry_hours = %d", *cfg.rotateBeforeExpiryHours)
	}
	if cfg.secretVersion != nil {
		extra += fmt.Sprintf("\n    secret_version = %d", *cfg.secretVersion)
	}
	return fmt.Sprintf(`
resource %q %q {
  name         = %q
  org_owner_id = %q

  client_secret_rotation = {
    expires_after_hours = %d%s
  }
}
`, resourceType, name, name, orgOwnerID, cfg.expiresAfterHours, extra)
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
