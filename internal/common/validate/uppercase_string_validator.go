package validate

import (
	"context"
	"fmt"
	"strings"

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
	if value != strings.ToUpper(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("The provided string %q must be uppercase.", value),
			"",
		)
	}
}

func ValidUppercaseString() validator.String {
	return UppercaseStringValidator{}
}
