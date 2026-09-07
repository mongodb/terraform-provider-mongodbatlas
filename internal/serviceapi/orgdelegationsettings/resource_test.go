package orgdelegationsettings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceType = "mongodbatlas_org_delegation_settings"
	resourceName = resourceType + ".this"
)

func TestAccOrgDelegationSettings_basic(t *testing.T) {
	acc.SkipInUnitTest(t) // baseline capture below runs before resource.Test skips
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	// Delete is a no-op for this singleton resource, so the current settings are
	// captured before the test and restored afterwards to leave the org unchanged.
	baseline, err := getDelegationSettings(orgID)
	if err != nil {
		t.Fatalf("failed to capture current delegation settings for org %s: %v", orgID, err)
	}
	t.Cleanup(func() {
		if err := updateDelegationSettings(orgID, baseline); err != nil {
			t.Errorf("failed to restore delegation settings for org %s: %v", orgID, err)
		}
	})

	// Serial execution because delegation settings are a singleton at the org level.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, "READ_ONLY", "READ_WRITE", 3600, 86400),
				Check:  checkBasic(orgID, "READ_ONLY", "READ_WRITE", 3600, 86400),
			},
			{
				Config: configBasic(orgID, "READ_WRITE", "DISALLOWED", 7200, 43200),
				Check:  checkBasic(orgID, "READ_WRITE", "DISALLOWED", 7200, 43200),
			},
			{
				// Access policies are computed, so omitting them keeps the values returned by Atlas.
				// Token lifetimes are optional-only with send_null_as_null_on_update, so omitting them
				// sends null and resets them to the system default (returned as null by Atlas).
				Config: configOmitted(orgID),
				Check:  checkReset(orgID, "READ_WRITE", "DISALLOWED"),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "org_id",
			},
		},
	})
}

func configBasic(orgID, mcpAccess, partnerAccess string, idleLifetime, maxLifetime int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_delegation_settings" "this" {
			org_id                        = %[1]q
			delegated_mcp_access          = %[2]q
			delegated_partner_access      = %[3]q
			idle_refresh_token_lifetime   = %[4]d
			maximum_refresh_token_lifetime = %[5]d
		}
	`, orgID, mcpAccess, partnerAccess, idleLifetime, maxLifetime)
}

func configOmitted(orgID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_delegation_settings" "this" {
			org_id = %[1]q
		}
	`, orgID)
}

func checkBasic(orgID, mcpAccess, partnerAccess string, idleLifetime, maxLifetime int) resource.TestCheckFunc {
	attrChecks := map[string]string{
		"org_id":                         orgID,
		"delegated_mcp_access":           mcpAccess,
		"delegated_partner_access":       partnerAccess,
		"idle_refresh_token_lifetime":    fmt.Sprintf("%d", idleLifetime),
		"maximum_refresh_token_lifetime": fmt.Sprintf("%d", maxLifetime),
	}
	checks := acc.AddAttrChecks(resourceName, nil, attrChecks)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkReset(orgID, mcpAccess, partnerAccess string) resource.TestCheckFunc {
	attrChecks := map[string]string{
		"org_id":                   orgID,
		"delegated_mcp_access":     mcpAccess,
		"delegated_partner_access": partnerAccess,
	}
	checks := acc.AddAttrChecks(resourceName, nil, attrChecks)
	checks = append(checks,
		resource.TestCheckNoResourceAttr(resourceName, "idle_refresh_token_lifetime"),
		resource.TestCheckNoResourceAttr(resourceName, "maximum_refresh_token_lifetime"),
		checkExists(resourceName),
	)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		if orgID == "" {
			return fmt.Errorf("no org_id is set")
		}
		if _, err := getDelegationSettings(orgID); err != nil {
			return fmt.Errorf("delegation settings for org(%s) do not exist: %w", orgID, err)
		}
		return nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["org_id"], nil
	}
}

type delegationSettings struct {
	DelegatedMcpAccess          *string `json:"delegatedMcpAccess"`
	DelegatedPartnerAccess      *string `json:"delegatedPartnerAccess"`
	IdleRefreshTokenLifetime    *int64  `json:"idleRefreshTokenLifetime"`
	MaximumRefreshTokenLifetime *int64  `json:"maximumRefreshTokenLifetime"`
}

func getDelegationSettings(orgID string) (*delegationSettings, error) {
	resp, err := delegationSettingsAPICall(orgID, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching delegation settings for org %s", resp.StatusCode, orgID)
	}
	var result delegationSettings
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func updateDelegationSettings(orgID string, settings *delegationSettings) error {
	body, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	resp, err := delegationSettingsAPICall(orgID, http.MethodPatch, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d updating delegation settings for org %s", resp.StatusCode, orgID)
	}
	return nil
}

func delegationSettingsAPICall(orgID, method string, body []byte) (*http.Response, error) {
	callParams := config.APICallParams{
		VersionHeader: "application/vnd.atlas.2025-03-12+json",
		RelativePath:  "/api/atlas/v2/orgs/{orgId}/delegationSettings",
		PathParams:    map[string]string{"orgId": orgID},
		Method:        method,
	}
	return acc.MongoDBClient.UntypedAPICall(context.Background(), callParams, body)
}
