package validate

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

func InstanceSizeNameValidator() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if value == "M2" || value == "M5" {
			diagError := diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf(constant.DeprecationSharedTier, constant.ServerlessSharedEOLDate, "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
			}
			diags = append(diags, diagError)
		}
		return diags
	}
}
