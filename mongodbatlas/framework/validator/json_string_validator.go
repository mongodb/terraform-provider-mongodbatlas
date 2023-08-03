package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
)

type jsonStringValidator struct{}

func (v jsonStringValidator) Description(_ context.Context) string {
	return "string value must be a valid JSON"
}

func (v jsonStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

//nolint:gocritic //we have to comply with validator interface, cannot pass req as a pointer
func (v jsonStringValidator) ValidateString(ctx context.Context, req validator.StringRequest, response *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	if _, err := structure.NormalizeJsonString(req.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
	}
}

func StringIsJSON() validator.String {
	return jsonStringValidator{}
}
