// Package deletenotfound provides shared handling of 404 responses from Atlas API delete calls.
// Deleting a resource that is already gone achieves the goal of the delete, so a 404 is treated
// as success (see https://developer.hashicorp.com/terraform/plugin/framework/resources/delete).
package deletenotfound

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// IsNotFound reports whether err is a 404 response from an Atlas API delete call, meaning the
// resource already no longer exists and the delete can be considered successful. The tolerated
// case is logged so it can be observed in production. A nil response with a non-nil error
// (e.g. a transport failure) is never treated as not-found.
func IsNotFound(ctx context.Context, resp *http.Response, err error) bool {
	if err == nil || !validate.StatusNotFound(resp) {
		return false
	}
	tflog.Info(ctx, "delete returned 404, resource already no longer exists, considering the delete as successful")
	return true
}
