package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StringIsUppercase() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if value != strings.ToUpper(value) {
			diagError := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("The provided string '%q' must be uppercase.", value),
			}
			diags = append(diags, diagError)
		}
		return diags
	}
}
