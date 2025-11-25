package alertconfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
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
	// Path format: notification[0].interval_min
	notificationPath := req.Path.ParentPath()
	var notification TfNotificationModel
	diags := req.Config.GetAttribute(ctx, notificationPath, &notification)
	if diags.HasError() {
		return
	}
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
