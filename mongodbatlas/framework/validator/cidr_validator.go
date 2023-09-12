package validator

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type CIDRValidator struct{}

func (v CIDRValidator) Description(_ context.Context) string {
	return "string value must be defined as a valid cidr."
}

func (v CIDRValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v CIDRValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
		return
	}

	if ipnet == nil || ipnet.String() != value {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
		return
	}
}

func ValidCIDR() validator.String {
	return CIDRValidator{}
}
