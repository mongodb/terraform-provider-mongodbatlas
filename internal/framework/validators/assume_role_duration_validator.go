package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringLengthBetweenValidator struct {
	Max int
	Min int
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", v.Min, v.Max)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between `%d` and `%d`", v.Min, v.Max)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringLengthBetweenValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	// strLen := len(req.ConfigValue.ValueString())

	// if strLen < v.Min || strLen > v.Max {
	// 	resp.Diagnostics.AddAttributeError(
	// 		// req.AttributePath,
	// 		"Invalid String Length",
	// 		fmt.Sprintf("String length must be between %d and %d, got: %d.", v.Min, v.Max, strLen),
	// 	)

	// 	return
	// }
}

// schema.Schema{
//     Attributes: map[string]schema.Attribute{
//         "root_list_attribute": schema.ListNestedAttribute{
//             NestedObject: schema.NestedAttributeObject{
//                 Attributes: map[string]schema.Attribute{
//                     "nested_list_attribute": schema.ListNestedAttribute{
//                         NestedObject: schema.NestedAttributeObject{
//                             Attributes: map[string]schema.Attribute{
//                                 "deeply_nested_string_attribute": schema.StringAttribute{
//                                     Required: true,
//                                 },
//                             },
//                         },
//                         Required: true,
//                         Validators: []validator.List{
//                             exampleValidatorThatAcceptsExpressions(
//                                 path.MatchRelative().AtParent().AtName("nested_string_attribute"),
//                             ),
//                         },
//                     },
//                     "nested_string_attribute": schema.StringAttribute{
//                         Required: true,
//                     },
//                 },
//             },
//             Required: true,
//         },
//     },
// }
