package conversion

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func FromTPFDiagsToSDKV2Diags(diagsTpf []diag.Diagnostic) sdkv2diag.Diagnostics {
	var results []sdkv2diag.Diagnostic
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

type ErrBody interface {
	Body() []byte
}

// AddJSONBodyErrorToDiagnostics tries to get the JSON body from the error and add it to the diagnostics.
// For example, admin.GenericOpenAPIError has the Body() []byte method.
func AddJSONBodyErrorToDiagnostics(msgPrefix string, err error, diags *diag.Diagnostics) {
	errGeneric, ok := err.(ErrBody)
	if !ok {
		diags.AddError(msgPrefix, err.Error())
		return
	}
	var respJSON map[string]any
	errMarshall := json.Unmarshal(errGeneric.Body(), &respJSON)
	if errMarshall != nil {
		diags.AddError(msgPrefix, err.Error())
		return
	}
	errorBytes, errMarshall := json.MarshalIndent(respJSON, "", "  ")
	if errMarshall != nil {
		diags.AddError(msgPrefix, err.Error())
		return
	}
	errorJSON := string(errorBytes)
	diags.AddError(msgPrefix, errorJSON)
}

func FormatDiags(diags *diag.Diagnostics) string {
	if diags == nil {
		return ""
	}
	summary := []string{}

	addDiag := func(d diag.Diagnostic) {
		summary = append(summary, d.Summary(), "\t detail: "+d.Detail())
	}

	for _, errorDiag := range diags.Errors() {
		addDiag(errorDiag)
	}
	if diags.WarningsCount() > 0 {
		if len(summary) > 0 {
			summary = append(summary, "\nWarnings:")
		} else {
			summary = append(summary, "Warnings:")
		}
		for _, warningDiag := range diags.Warnings() {
			addDiag(warningDiag)
		}
	}
	return strings.Join(summary, "\n")
}
