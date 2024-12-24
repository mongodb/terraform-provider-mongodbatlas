package conversion

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	legacyDiag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

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

func AddLegacyDiags(diags *diag.Diagnostics, legacyDiags legacyDiag.Diagnostics) {
	for _, diag := range legacyDiags {
		if diag.Severity == legacyDiag.Error {
			diags.AddError(diag.Summary, diag.Detail)
		} else {
			diags.AddWarning(diag.Summary, diag.Detail)
		}
	}
}
