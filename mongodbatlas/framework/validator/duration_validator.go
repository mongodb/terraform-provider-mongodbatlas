package validator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type durationValidator struct {
	MinMinutes int
	MaxMinutes int
}

func (v durationValidator) Description(_ context.Context) string {
	ds := "string value must be defined as a valid duration, and must be between %d and %d minutes, inclusive. Valid time units are ns, us (or µs), ms, s, h, or m."
	return fmt.Sprintf(ds, v.MinMinutes, v.MaxMinutes)
}

func (v durationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

<<<<<<< HEAD
=======
//nolint:gocritic //we have to comply with validator interface, cannot pass req as a pointer
>>>>>>> 6feaad9a (feat: new framework provider, main and acceptance tests to use mux server with existing sdk v2 provider (#1366))
func (v durationValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
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
	return durationValidator{
		MinMinutes: minMinutes,
		MaxMinutes: maxMinutes,
	}
}
