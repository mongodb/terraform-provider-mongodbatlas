package mcpconfig_test

import (
	"context"
	"fmt"
	"os"
	"strings"
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

func (e ipAccessListEntry) attrMap() map[string]string {
	result := map[string]string{}
	if e.cidr != "" {
		result["cidr_block"] = e.cidr
	}
	if e.ip != "" {
		result["ip_address"] = e.ip
	}
	return result
}

func (e ipAccessListEntry) hclStr() string {
	if e.cidr != "" {
		return fmt.Sprintf("{cidr_block = %q}", e.cidr)
	}
	if e.ip != "" {
		return fmt.Sprintf("{ip_address = %q}", e.ip)
	}
	return ""
}

func TestAccMcpConfig_basic(t *testing.T) {
	var (
		orgID  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name1  = acc.RandomName()
		name2  = fmt.Sprintf("%s-updated", name1)
		withDS = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name1, []string{"ORG_READ_ONLY"}, nil, !withDS),
				Check:  checkBasic([]string{"ORG_READ_ONLY"}, nil, !withDS),
			},
			{
				Config: configBasic(orgID, name2, []string{"ORG_MEMBER", "ORG_READ_ONLY"}, nil, !withDS),
				Check:  checkBasic([]string{"ORG_MEMBER", "ORG_READ_ONLY"}, nil, !withDS),
			},
			{
				Config: configBasic(orgID, name2, []string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.113.0"}}, !withDS),
				Check:  checkBasic([]string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.113.0"}}, !withDS),
			},
			{ // Change name keeping ip_access_list the same, plans ip + cidr, hook removes cidr.
				Config: configBasic(orgID, name1, []string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.113.0"}}, !withDS),
				Check:  checkBasic([]string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.113.0"}}, !withDS),
			},
			{
				Config: configBasic(orgID, name1, []string{"ORG_MEMBER"}, []ipAccessListEntry{{cidr: "203.0.113.0/24"}}, !withDS),
				Check:  checkBasic([]string{"ORG_MEMBER"}, []ipAccessListEntry{{cidr: "203.0.113.0/24"}}, !withDS),
			},
			{
				Config: configBasic(orgID, name1, []string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.111.0"}, {cidr: "203.0.112.0/32"}, {cidr: "203.0.113.0/24"}}, !withDS),
				Check:  checkBasic([]string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.111.0"}, {cidr: "203.0.112.0/32"}, {cidr: "203.0.113.0/24"}}, !withDS),
			},
			{ // No-op step: the plural list endpoint returns 404 while the backing service accounts are being updated, so read it only once the resource is stable.
				Config: configBasic(orgID, name1, []string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.111.0"}, {cidr: "203.0.112.0/32"}, {cidr: "203.0.113.0/24"}}, withDS),
				Check:  checkBasic([]string{"ORG_MEMBER"}, []ipAccessListEntry{{ip: "203.0.111.0"}, {cidr: "203.0.112.0/32"}, {cidr: "203.0.113.0/24"}}, withDS),
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

func configBasic(orgID, name string, roles []string, entries []ipAccessListEntry, withDS bool) string {
	rolesHCL := hcl.StringSliceToHCL(roles)
	pluralDSHCL := ""
	ipAccessListHCL := ""
	if len(entries) > 0 {
		var entryBlocks []string
		for _, e := range entries {
			entryBlocks = append(entryBlocks, e.hclStr())
		}
		ipAccessListHCL = fmt.Sprintf("ip_access_list = [%s]", strings.Join(entryBlocks, ", "))
	}
	// The list operation returns 404 while the backing MCP Config service accounts are being
	// created or updated, so only request it in a step that does not mutate the resource.
	if withDS {
		pluralDSHCL = fmt.Sprintf(`
			data "mongodbatlas_mcp_configs" "test" {
				org_id     = %[1]q
				depends_on = [mongodbatlas_mcp_config.test]
			}
		`, orgID)
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

		%[5]s
	`, orgID, name, rolesHCL, ipAccessListHCL, pluralDSHCL)
}

func checkBasic(roles []string, entries []ipAccessListEntry, withDS bool) resource.TestCheckFunc {
	commonAttrsSet := []string{"mcp_config_id", "client_id", "egress_client_id"}
	attrsMap := map[string]string{
		"roles.#":          fmt.Sprintf("%d", len(roles)),
		"ip_access_list.#": fmt.Sprintf("%d", len(entries)),
	}
	checks := []resource.TestCheckFunc{
		acc.CheckRSAndDS(resourceName, new(dataSourceName), nil, commonAttrsSet, attrsMap, checkExists(resourceName)),
	}
	if withDS {
		checks = append(checks, resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)))
	}
	for _, e := range entries {
		attrMap := e.attrMap()
		checks = append(checks,
			resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ip_access_list.*", attrMap),
			resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "ip_access_list.*", attrMap),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
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
	for name, rs := range s.RootModule().Resources {
		if name != resourceName {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if orgID == "" || mcpConfigID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		if _, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetOrgMcpConfig(context.Background(), orgID, mcpConfigID).Execute(); err == nil {
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
