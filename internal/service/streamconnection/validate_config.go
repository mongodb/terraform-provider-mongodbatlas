package streamconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const awsMSKIAMMechanism = "AWS_MSK_IAM"

type kafkaIAMAuthenticationValidator struct{}

func (kafkaIAMAuthenticationValidator) Description(context.Context) string {
	return "AWS_MSK_IAM authentication requires authentication.aws.role_arn."
}

func (v kafkaIAMAuthenticationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (kafkaIAMAuthenticationValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var authentication types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("authentication"), &authentication)...)
	if resp.Diagnostics.HasError() || authentication.IsNull() || authentication.IsUnknown() {
		return
	}

	authenticationModel := &TFConnectionAuthenticationModel{}
	resp.Diagnostics.Append(authentication.As(ctx, authenticationModel, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() || authenticationModel.Mechanism.IsNull() || authenticationModel.Mechanism.IsUnknown() || authenticationModel.Mechanism.ValueString() != awsMSKIAMMechanism {
		return
	}

	roleARNPath := path.Root("authentication").AtName("aws").AtName("role_arn")
	if authenticationModel.AWS.IsNull() || authenticationModel.AWS.IsUnknown() {
		resp.Diagnostics.AddAttributeError(roleARNPath, "AWS MSK IAM role ARN is required", "authentication.aws.role_arn must be set when authentication.mechanism is AWS_MSK_IAM.")
		return
	}

	awsModel := &TFAWSModel{}
	resp.Diagnostics.Append(authenticationModel.AWS.As(ctx, awsModel, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() || awsModel.RoleArn.IsUnknown() {
		return
	}
	if awsModel.RoleArn.IsNull() || awsModel.RoleArn.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(roleARNPath, "AWS MSK IAM role ARN is required", "authentication.aws.role_arn must be set when authentication.mechanism is AWS_MSK_IAM.")
	}
}
