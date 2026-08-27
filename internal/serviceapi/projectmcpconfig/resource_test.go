package projectmcpconfig_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

const resourceName = "mongodbatlas_project_mcp_config.test"
const dataSourceName = "data.mongodbatlas_project_mcp_config.test"
const dataSourcePluralName = "data.mongodbatlas_project_mcp_configs.test"

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

func ipAccessListAttrMap(entries []ipAccessListEntry) map[string]string {
	if len(entries) == 0 {
		return nil
	}
	result := map[string]string{"ip_access_list.#": fmt.Sprintf("%d", len(entries))}
	for i, e := range entries {
		if e.cidr != "" {
			result[fmt.Sprintf("ip_access_list.%d.cidr_block", i)] = e.cidr
		}
		if e.ip != "" {
			result[fmt.Sprintf("ip_access_list.%d.ip_address", i)] = e.ip
		}
	}
	return result
}

func TestAccProjectMcpConfig_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		name1     = acc.RandomName()
		name2     = fmt.Sprintf("%s-updated", name1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name1, []string{"GROUP_READ_ONLY"}, nil),
				Check:  checkBasic([]string{"GROUP_READ_ONLY"}, nil),
			},
			{
				Config: configBasic(projectID, name2, []string{"GROUP_OWNER"}, nil),
				Check:  checkBasic([]string{"GROUP_OWNER"}, nil),
			},
			{
				Config: configBasic(projectID, name2, []string{"GROUP_OWNER", "GROUP_READ_ONLY"}, nil),
				Check:  checkBasic([]string{"GROUP_OWNER", "GROUP_READ_ONLY"}, nil),
			},
			{
				Config: configBasic(projectID, name2, []string{"GROUP_OWNER"}, []ipAccessListEntry{{ip: "203.0.113.10"}}),
				Check:  checkBasic([]string{"GROUP_OWNER"}, []ipAccessListEntry{{ip: "203.0.113.10"}}),
			},
			{
				Config: configBasic(projectID, name2, []string{"GROUP_OWNER"}, []ipAccessListEntry{{cidr: "203.0.113.0/24"}}),
				Check:  checkBasic([]string{"GROUP_OWNER"}, []ipAccessListEntry{{cidr: "203.0.113.0/24"}}),
			},
			{
				Config: configBasic(projectID, name2, []string{"GROUP_OWNER"}, []ipAccessListEntry{{ip: "203.0.113.10"}, {cidr: "203.0.113.0/24"}}),
				Check:  checkBasic([]string{"GROUP_OWNER"}, []ipAccessListEntry{{ip: "203.0.113.10"}, {cidr: "203.0.113.0/24"}}),
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

func configBasic(projectID, name string, roles []string, entries []ipAccessListEntry) string {
	rolesHCL := hcl.StringSliceToHCL(roles)
	ipAccessListHCL := ""
	if len(entries) > 0 {
		var entryBlocks []string
		for _, e := range entries {
			entryBlocks = append(entryBlocks, fmt.Sprintf("{\n\t\t\t\t\t%s\n\t\t\t\t}", e.hclStr()))
		}
		ipAccessListHCL = fmt.Sprintf(`
			ip_access_list = [
				%s
			]`, strings.Join(entryBlocks, ",\n\t\t\t\t"))
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project_mcp_config" "test" {
			project_id      = %[1]q
			mcp_config_name = %[2]q
			roles           = %[3]s
			%[4]s
		}

		data "mongodbatlas_project_mcp_config" "test" {
			project_id    = %[1]q
			mcp_config_id = mongodbatlas_project_mcp_config.test.mcp_config_id
		}

		data "mongodbatlas_project_mcp_configs" "test" {
			project_id = %[1]q
			depends_on = [mongodbatlas_project_mcp_config.test]
		}
	`, projectID, name, rolesHCL, ipAccessListHCL)
}

func checkBasic(roles []string, entries []ipAccessListEntry) resource.TestCheckFunc {
	commonAttrsSet := []string{"mcp_config_id", "client_id", "egress_client_id"}
	attrsMap := ipAccessListAttrMap(entries)
	if attrsMap == nil {
		attrsMap = map[string]string{}
	}
	attrsMap["roles.#"] = fmt.Sprintf("%d", len(roles))
	checks := acc.CheckRSAndDS(resourceName, new(dataSourceName), new(dataSourcePluralName), commonAttrsSet, attrsMap, checkExists(resourceName))
	return checks
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if projectID == "" || mcpConfigID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetGroupMcpConfig(context.Background(), projectID, mcpConfigID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("mcp config (%s/%s) does not exist: %w", projectID, mcpConfigID, err)
	}
}

func checkDestroy(s *terraform.State) error {
	var errs []error
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "mongodbatlas_project_mcp_config.") {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if projectID == "" || mcpConfigID == "" {
			errs = append(errs, fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName))
			continue
		}

		if _, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetGroupMcpConfig(context.Background(), projectID, mcpConfigID).Execute(); err == nil {
			errs = append(errs, fmt.Errorf("mcp config (%s/%s) still exists", projectID, mcpConfigID))
		}
	}
	if err := acc.CheckDestroyDeleteOrgMcpConfigs(s); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		if projectID == "" || mcpConfigID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, mcpConfigID), nil
	}
}
