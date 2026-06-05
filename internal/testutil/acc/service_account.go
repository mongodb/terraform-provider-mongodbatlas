package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// CheckDestroyDeleteProjectSAs call this from Project SA tests checkDestroy functions to delete the service accounts at
// the org level (project_service_account DELETE only removes the project assignment)
func CheckDestroyDeleteProjectSAs(s *terraform.State) error {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	var errs []error
	for name, rs := range s.RootModule().Resources {
		// Prefix includes `.` to avoid matching other resources such as project_service_account_access_list_entry.
		if !strings.HasPrefix(name, "mongodbatlas_project_service_account.") {
			continue
		}
		if clientID := rs.Primary.Attributes["client_id"]; clientID != "" {
			if _, err := ConnV2().ServiceAccountsApi.DeleteOrgServiceAccount(context.Background(), clientID, orgID).Execute(); err != nil {
				errs = append(errs, fmt.Errorf("failed to delete service account %s: %w", clientID, err))
			}
		}
	}
	return errors.Join(errs...)
}
