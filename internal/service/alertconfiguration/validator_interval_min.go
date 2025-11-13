package alertconfiguration

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IntervalMinValidator validates that interval_min is not set for notification types
// that don't support it (PAGER_DUTY, OPS_GENIE, VICTOR_OPS).
//
// API-level validation cannot catch this error because when the notification type
// doesn't support interval_min, the nested_type in the API payload omits the field,
// so it's never sent to the API. Client-side validation ensures users get immediate
// feedback during Terraform's plan phase.
type IntervalMinValidator struct{}

func (v IntervalMinValidator) Description(_ context.Context) string {
	return "'interval_min' must not be set if type_name is 'PAGER_DUTY', 'OPS_GENIE' or 'VICTOR_OPS'"
}

func (v IntervalMinValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v IntervalMinValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	// If the value is unknown/null/0, there is nothing to validate.
	if req.ConfigValue.ValueInt64() <= 0 {
		return
	}

	// Parse the path to find which notification index we're validating
	// Path format: notification[0].interval_min
	pathStr := req.Path.String()
	notificationIndex := -1

	// Extract the index from the path (e.g., "notification[0]" -> 0)
	if idxStart := strings.Index(pathStr, "["); idxStart != -1 {
		if idxEnd := strings.Index(pathStr[idxStart:], "]"); idxEnd != -1 {
			idxStr := pathStr[idxStart+1 : idxStart+idxEnd]
			if idx, err := strconv.Atoi(idxStr); err == nil {
				notificationIndex = idx
			}
		}
	}

	// If we couldn't parse the index, skip validation
	if notificationIndex < 0 {
		return
	}

	// Get the entire notification list from config
	var notifications []TfNotificationModel
	diags := req.Config.GetAttribute(ctx, path.Root("notification"), &notifications)
	if diags.HasError() {
		// If we can't read notifications, skip validation (might be unknown during plan)
		return
	}

	// Check if we have the notification at the parsed index
	if notificationIndex >= len(notifications) {
		return
	}

	notification := notifications[notificationIndex]
	typeNameValue := notification.TypeName.ValueString()
	// Check if the type_name is one of the unsupported types
	if strings.EqualFold(typeNameValue, pagerDuty) ||
		strings.EqualFold(typeNameValue, opsGenie) ||
		strings.EqualFold(typeNameValue, victorOps) {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", req.ConfigValue.ValueInt64()),
		))
	}
}

func ValidIntervalMin() validator.Int64 {
	return IntervalMinValidator{}
}
