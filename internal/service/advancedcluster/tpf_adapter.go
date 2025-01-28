package advancedcluster

import (
	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	v2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ConvertV2DiagsToFrameworkDiags(v2Diags v2diag.Diagnostics) *frameworkdiag.Diagnostics {
	var diags frameworkdiag.Diagnostics
	for _, v2Diag := range v2Diags {
		diags.AddError(v2Diag.Summary, v2Diag.Detail)
	}
	return &diags
}

func ConvertFrameworkDiagsToV2Diags(diags frameworkdiag.Diagnostics) v2diag.Diagnostics {
	var v2Diags v2diag.Diagnostics
	for _, diag := range diags {
		v2Diags = append(v2Diags, v2diag.Diagnostic{
			Severity: v2diag.Severity(diag.Severity()),
			Summary:  diag.Summary(),
			Detail:   diag.Detail(),
		})
	}
	return v2Diags
}
