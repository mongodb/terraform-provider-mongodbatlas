package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckDestroyDeleteOrgMcpConfigs(s *terraform.State) error {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	var errs []error
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "mongodbatlas_project_mcp_config.") {
			continue
		}
		if mcpConfigID := rs.Primary.Attributes["mcp_config_id"]; mcpConfigID != "" {
			if _, err := ConnPreview().RemoteMCPConfigurationsAPI.DeleteOrgMcpConfig(context.Background(), orgID, mcpConfigID).Cascading(true).Execute(); err != nil {
				errs = append(errs, fmt.Errorf("failed to delete org mcp config %s: %w", mcpConfigID, err))
			}
		}
	}
	return errors.Join(errs...)
}
