package conversion

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
