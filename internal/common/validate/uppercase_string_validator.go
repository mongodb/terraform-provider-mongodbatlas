package validate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type UppercaseStringValidator struct{}

func (v UppercaseStringValidator) Description(_ context.Context) string {
	return "string value must be uppercase."
}

func (v UppercaseStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v UppercaseStringValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if !isUppercase(value) {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			value,
		))
	}
}

func ValidUppercaseString() validator.String {
	return UppercaseStringValidator{}
}
