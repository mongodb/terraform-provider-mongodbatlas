// Package deletenotfound provides shared handling of "already gone" responses from Atlas API delete calls.
// Deleting a resource that is already gone achieves the goal of the delete, so a 404 (or, for some APIs,
// a specific error code) is treated as success (see https://developer.hashicorp.com/terraform/plugin/framework/resources/delete).
package deletenotfound

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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

// IsNotFoundOrErrorCode is the equivalent of IsNotFound for APIs that answer a delete of a missing
// resource with a non-404 status and a specific Atlas error code (e.g. 400 API_KEY_NOT_FOUND).
// Besides a 404, err carrying errorCode is also treated as not-found. The error code is read from
// the JSON body of the error, so errors from any Atlas SDK version are matched. The HTTP status is
// intentionally not checked: the error code alone identifies the case, and the status could change.
func IsNotFoundOrErrorCode(ctx context.Context, resp *http.Response, err error, errorCode string) bool {
	if IsNotFound(ctx, resp, err) {
		return true
	}
	if !hasErrorCode(err, errorCode) {
		return false
	}
	tflog.Info(ctx, "delete returned "+errorCode+", resource already no longer exists, considering the delete as successful")

	return true
}

// hasErrorCode reports whether err has a JSON body with the given Atlas error code. It relies on
// conversion.ErrBody instead of a specific SDK type because the SDK error type (GenericOpenAPIError)
// is different in each SDK version, while its Body method is the same.
func hasErrorCode(err error, errorCode string) bool {
	var errBody conversion.ErrBody
	if !errors.As(err, &errBody) {
		return false
	}
	var apiError struct {
		ErrorCode string `json:"errorCode"`
	}
	return json.Unmarshal(errBody.Body(), &apiError) == nil && apiError.ErrorCode == errorCode
}

// IsNotFoundLegacy is the equivalent of IsNotFound for legacy SDKs (go.mongodb.org/atlas and
// go.mongodb.org/realm) whose calls return a wrapper instead of the raw *http.Response.
// The wrapper can be nil when the request fails before a response exists (e.g. a transport
// failure), in which case the error is not treated as not-found.
func IsNotFoundLegacy(ctx context.Context, resp *mongodbatlas.Response, err error) bool {
	if resp == nil {
		return false
	}
	return IsNotFound(ctx, resp.Response, err)
}
