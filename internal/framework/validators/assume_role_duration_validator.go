package validators

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// This validator is currently unused, added only for POC purpose
type assumeRoleDurationValidator struct {
	Max int
	Min int
}

func (v assumeRoleDurationValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", v.Min, v.Max)
}

func (v assumeRoleDurationValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between `%d` and `%d`", v.Min, v.Max)
}

func (v assumeRoleDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	configVal := req.ConfigValue.ValueString()
	duration, err := time.ParseDuration(configVal)

	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Duration",
			fmt.Sprintf("%q cannot be parsed as a duration: %v", configVal, err),
		)
		return
	}

	if duration.Minutes() < 15 || duration.Hours() > 12 {
		resp.Diagnostics.AddError(
			"Invalid Duration",
			fmt.Sprintf("duration %q must be between 15 minutes (15m) and 12 hours (12h), inclusive", configVal),
		)
	}
}
