package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// Retries a "still exists" check after delete.
func WaitUntilGone(stillExists func() bool) bool {
	for range 5 {
		if !stillExists() {
			return true
		}
		time.Sleep(3 * time.Second)
	}
	return false
}

func CheckDestroyDeleteOrgMcpConfigs(s *terraform.State) error {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	var errs []error
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "mongodbatlas_project_mcp_config.") {
			continue
		}
		if mcpConfigID := rs.Primary.Attributes["mcp_config_id"]; mcpConfigID != "" {
			resp, err := ConnPreview().RemoteMCPConfigurationsAPI.DeleteOrgMcpConfig(context.Background(), orgID, mcpConfigID).Cascading(true).Execute()
			if err != nil && !validate.StatusNotFound(resp) {
				errs = append(errs, fmt.Errorf("failed to delete org mcp config %s: %w", mcpConfigID, err))
			}
		}
	}
	return errors.Join(errs...)
}
