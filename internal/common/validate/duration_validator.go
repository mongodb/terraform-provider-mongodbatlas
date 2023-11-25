package validate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type DurationValidator struct {
	MinMinutes int
	MaxMinutes int
}

func (v DurationValidator) Description(_ context.Context) string {
	ds := "string value must be defined as a valid duration, and must be between %d and %d minutes, inclusive. Valid time units are ns, us (or µs), ms, s, h, or m."
	return fmt.Sprintf(ds, v.MinMinutes, v.MaxMinutes)
}

func (v DurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v DurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	duration, err := time.ParseDuration(value)

	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
		return
	}

	if duration.Minutes() < float64(v.MinMinutes) || duration.Minutes() > float64(v.MaxMinutes) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
	}
}

func ValidDurationBetween(minMinutes, maxMinutes int) validator.String {
	return DurationValidator{
		MinMinutes: minMinutes,
		MaxMinutes: maxMinutes,
	}
}
