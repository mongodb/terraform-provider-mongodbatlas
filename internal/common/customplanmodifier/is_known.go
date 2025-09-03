package customplanmodifier

import "github.com/hashicorp/terraform-plugin-framework/attr"

// IsKnown returns true if the attribute is known (not null or unknown). Note that !IsKnown is not the same as IsUnknown because null is !IsKnown but not IsUnknown.
func IsKnown(attribute attr.Value) bool {
	return !attribute.IsNull() && !attribute.IsUnknown()
}
