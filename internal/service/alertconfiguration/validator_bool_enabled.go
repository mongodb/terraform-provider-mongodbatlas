package alertconfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var validEmailSMSEnabledTypes = []string{"ORG", "GROUP", "USER"}

type boolEnabledValidator struct {
	fieldName string
}

func (v boolEnabledValidator) Description(_ context.Context) string {
	return fmt.Sprintf("'%s' is only valid if 'type_name' is set to 'ORG', 'GROUP', or 'USER'", v.fieldName)
}

func (v boolEnabledValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v boolEnabledValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || !req.ConfigValue.ValueBool() {
		return
	}
	notificationPath := req.Path.ParentPath()
	var notification TfNotificationModel
	diags := req.Config.GetAttribute(ctx, notificationPath, &notification)
	if diags.HasError() {
		return
	}
	if notification.TypeName.IsNull() || notification.TypeName.IsUnknown() {
		return
	}
	typeNameValue := notification.TypeName.ValueString()
	for _, validType := range validEmailSMSEnabledTypes {
		if strings.EqualFold(typeNameValue, validType) {
			return
		}
	}
	resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		req.Path,
		v.Description(ctx),
		"true",
	))
}

func validEmailEnabled() validator.Bool {
	return boolEnabledValidator{fieldName: "email_enabled"}
}

func validSMSEnabled() validator.Bool {
	return boolEnabledValidator{fieldName: "sms_enabled"}
}
