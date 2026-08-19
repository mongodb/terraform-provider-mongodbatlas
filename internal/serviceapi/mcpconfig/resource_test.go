package mcpconfig_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

const resourceName = "mongodbatlas_mcp_config.test"
const dataSourceName = "data.mongodbatlas_mcp_config.test"
const dataSourcePluralName = "data.mongodbatlas_mcp_configs.test"

type ipAccessListEntry struct {
	cidr string
	ip   string
}

func (e ipAccessListEntry) hclStr() string {
	if e.cidr != "" {
		return fmt.Sprintf("cidr_block = %q", e.cidr)
	}
	if e.ip != "" {
		return fmt.Sprintf("ip_address = %q", e.ip)
	}
	return ""
}

func (e ipAccessListEntry) attrMap() map[string]string {
	if e.cidr == "" && e.ip == "" {
		return nil
	}
	result := map[string]string{"ip_access_list.#": "1"}
	if e.cidr != "" {
		result["ip_access_list.0.cidr_block"] = e.cidr
	}
	if e.ip != "" {
		result["ip_access_list.0.ip_address"] = e.ip
	}
	return result
}

func TestAccMcpConfig_basic(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name1 = acc.RandomName()
		name2 = fmt.Sprintf("%s-updated", name1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name1, []string{"ORG_READ_ONLY"}, ipAccessListEntry{}),
				Check:  checkBasic(ipAccessListEntry{}),
			},
			{
				Config: configBasic(orgID, name2, []string{"ORG_MEMBER"}, ipAccessListEntry{}),
				Check:  checkBasic(ipAccessListEntry{}),
			},
			{
				Config: configBasic(orgID, name2, []string{"ORG_MEMBER"}, ipAccessListEntry{ip: "203.0.113.10"}),
				Check:  checkBasic(ipAccessListEntry{ip: "203.0.113.10"}),
			},
			{
				Config: configBasic(orgID, name2, []string{"ORG_MEMBER"}, ipAccessListEntry{cidr: "203.0.113.0/24"}),
				Check:  checkBasic(ipAccessListEntry{cidr: "203.0.113.0/24"}),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "mcp_config_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func configBasic(orgID, name string, roles []string, entry ipAccessListEntry) string {
	rolesHCL := hcl.StringSliceToHCL(roles)
	ipAccessListHCL := ""
	if entryHCL := entry.hclStr(); entryHCL != "" {
		ipAccessListHCL = fmt.Sprintf(`
			ip_access_list = [
				{
					%s
				}
			]`, entryHCL)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_mcp_config" "test" {
			org_id          = %[1]q
			mcp_config_name = %[2]q
			roles           = %[3]s
			%[4]s
		}

		data "mongodbatlas_mcp_config" "test" {
			org_id        = %[1]q
			mcp_config_id = mongodbatlas_mcp_config.test.mcp_config_id
		}

		data "mongodbatlas_mcp_configs" "test" {
			org_id     = %[1]q
			depends_on = [mongodbatlas_mcp_config.test]
		}
	`, orgID, name, rolesHCL, ipAccessListHCL)
}

func checkBasic(entry ipAccessListEntry) resource.TestCheckFunc {
	commonAttrsSet := []string{"mcp_config_id", "client_id", "egress_client_id"}
	checks := acc.CheckRSAndDS(resourceName, new(dataSourceName), new(dataSourcePluralName), commonAttrsSet, entry.attrMap(), checkExists(resourceName))
	return checks
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if orgID == "" || mcpConfigID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetOrgMcpConfig(context.Background(), orgID, mcpConfigID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("mcp config (%s/%s) does not exist: %w", orgID, mcpConfigID, err)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_mcp_config" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if orgID == "" || mcpConfigID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetOrgMcpConfig(context.Background(), orgID, mcpConfigID).Execute()
		if err == nil {
			return fmt.Errorf("mcp config (%s/%s) still exists", orgID, mcpConfigID)
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
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if orgID == "" || mcpConfigID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", orgID, mcpConfigID), nil
	}
}
