// Package deletenotfound provides shared handling of 404 responses from Atlas API delete calls.
// Deleting a resource that is already gone achieves the goal of the delete, so a 404 is treated
// as success (see https://developer.hashicorp.com/terraform/plugin/framework/resources/delete).
package deletenotfound

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mongodb.org/atlas/mongodbatlas"

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

// IsNotFound is the equivalent of IsNotFound for legacy SDKs (go.mongodb.org/atlas and
// go.mongodb.org/realm) whose calls return a wrapper instead of the raw *http.Response.
// The wrapper can be nil when the request fails before a response exists (e.g. a transport
// failure), in which case the error is not treated as not-found.
func IsNotFoundLegacy(ctx context.Context, resp *mongodbatlas.Response, err error) bool {
	if resp == nil {
		return false
	}
	return IsNotFound(ctx, resp.Response, err)
}
