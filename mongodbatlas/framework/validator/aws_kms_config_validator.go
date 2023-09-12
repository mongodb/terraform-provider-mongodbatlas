package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type awsKmsConfigValidator struct{}

func (v awsKmsConfigValidator) Description(_ context.Context) string {
	return "for credentials: `access_key_id` and `secret_access_key` are allowed but not `role_id`." +
		" For roles: `access_key_id` and `secret_access_key` are not allowed but `role_id` is allowed"
}

func (v awsKmsConfigValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v awsKmsConfigValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, response *validator.ObjectResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	obj, diag := req.ConfigValue.ToObjectValue(ctx)
	if diag.HasError() {
		response.Diagnostics.Append(diag.Errors()...)
		return
	}

	attrMap := obj.Attributes()
	ak, akOk := attrMap["access_key_id"]
	sa, saOk := attrMap["secret_access_key"]
	r, rOk := attrMap["role_id"]
	accessKeyDefined := akOk && !ak.IsNull()
	secretAccessKeyDefined := saOk && !sa.IsNull()
	roleIDDefined := rOk && !r.IsNull()

	if (accessKeyDefined && secretAccessKeyDefined && roleIDDefined) ||
		(accessKeyDefined && roleIDDefined) ||
		(secretAccessKeyDefined && roleIDDefined) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.String(),
		))
	}
}

func AwsKmsConfig() validator.Object {
	return awsKmsConfigValidator{}
}
