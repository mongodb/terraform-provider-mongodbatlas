package validator

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type IPValidator struct{}

func (v IPValidator) Description(_ context.Context) string {
	return "string value must be defined as a valid IP Address."
}

func (v IPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v IPValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	ip := net.ParseIP(value)
	if ip == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
	}
}

func ValidIP() validator.String {
	return IPValidator{}
}
