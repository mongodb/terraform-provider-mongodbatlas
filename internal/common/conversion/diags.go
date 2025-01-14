package conversion

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func FromTPFDiagsToSDKV2Diags(diagsTpf []diag.Diagnostic) sdkv2diag.Diagnostics {
	results := []sdkv2diag.Diagnostic{}
	for _, tpfDiag := range diagsTpf {
		sdkV2Sev := sdkv2diag.Warning
		if tpfDiag.Severity() == diag.SeverityError {
			sdkV2Sev = sdkv2diag.Error
		}
		results = append(results, sdkv2diag.Diagnostic{
			Severity: sdkV2Sev,
			Summary:  tpfDiag.Summary(),
			Detail:   tpfDiag.Detail(),
		})
	}
	return results
}
