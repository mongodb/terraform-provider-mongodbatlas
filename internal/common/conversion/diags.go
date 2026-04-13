package conversion

import (
	"encoding/json"
	"errors"
	"fmt"

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

// SDKv2DiagnosticsToError converts SDKv2 diagnostics to a single error joining all error-level entries.
// Returns nil when there are no errors.
func SDKv2DiagnosticsToError(diags sdkv2diag.Diagnostics) error {
	var errs []error
	for _, d := range diags {
		if d.Severity != sdkv2diag.Error {
			continue
		}
		if d.Detail != "" {
			errs = append(errs, fmt.Errorf("%s: %s", d.Summary, d.Detail))
		} else {
			errs = append(errs, errors.New(d.Summary))
		}
	}
	return errors.Join(errs...)
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
