package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type optionsValidator struct{}

func OptionsValidator() validator.Object {
	return optionsValidator{}
}

func (optionsValidator) Description(_ context.Context) string {
	return "must not be empty when configured"
}

func (v optionsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v optionsValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	options, diags := req.ConfigValue.ToObjectValue(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags.Errors()...)
		return
	}
	for _, value := range options.Attributes() {
		if !value.IsNull() && !value.IsUnknown() {
			return
		}
	}
	resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		req.Path,
		v.Description(ctx),
		req.ConfigValue.String(),
	))
}
