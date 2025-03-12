package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestValidAwsKmsConfig(t *testing.T) {
	enabled := true
	validType1 := map[string]attr.Type{
		"enabled":                types.BoolType,
		"customer_master_key_id": types.StringType,
		"region":                 types.StringType,
		"role_id":                types.StringType,
	}
	validValue1 := map[string]attr.Value{
		"enabled":                types.BoolValue(enabled),
		"customer_master_key_id": types.StringValue("testCustomerMasterKeyID"),
		"region":                 types.StringValue("testRegion"),
		"role_id":                types.StringValue("testRoleID"),
	}
	validType2 := map[string]attr.Type{
		"enabled":                types.BoolType,
		"access_key_id":          types.StringType,
		"secret_access_key":      types.StringType,
		"customer_master_key_id": types.StringType,
		"region":                 types.StringType,
	}
	validValue2 := map[string]attr.Value{
		"enabled":                types.BoolValue(enabled),
		"access_key_id":          types.StringValue("testAccessKey"),
		"secret_access_key":      types.StringValue("testSecretAccessKey"),
		"customer_master_key_id": types.StringValue("testCustomerMasterKeyID"),
		"region":                 types.StringValue("testRegion"),
	}
	inValidType := map[string]attr.Type{
		"enabled":                types.BoolType,
		"access_key_id":          types.StringType,
		"secret_access_key":      types.StringType,
		"customer_master_key_id": types.StringType,
		"region":                 types.StringType,
		"role_id":                types.StringType,
	}
	inValidValue := map[string]attr.Value{
		"enabled":                types.BoolValue(enabled),
		"access_key_id":          types.StringValue("testAccessKey"),
		"secret_access_key":      types.StringValue("testSecretAccessKey"),
		"customer_master_key_id": types.StringValue("testCustomerMasterKeyID"),
		"region":                 types.StringValue("testRegion"),
		"role_id":                types.StringValue("testRoleID"),
	}

	tests := []struct {
		awsKmsConfigValue map[string]attr.Value
		awsKmsConfigType  map[string]attr.Type
		name              string
		wantErr           bool
	}{
		{
			name:              "Valid Value 1",
			awsKmsConfigValue: validValue1,
			awsKmsConfigType:  validType1,
			wantErr:           false,
		},
		{
			name:              "Valid Value 2",
			awsKmsConfigValue: validValue2,
			awsKmsConfigType:  validType2,
			wantErr:           false,
		},
		{
			name:              "Invalid Value",
			awsKmsConfigValue: inValidValue,
			awsKmsConfigType:  inValidType,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		wantErr := tt.wantErr

		AwsKmsConfigValidator := validate.AwsKmsConfigValidator{}
		validatorRequest := validator.ObjectRequest{
			ConfigValue: types.ObjectValueMust(tt.awsKmsConfigType, tt.awsKmsConfigValue),
		}

		validatorResponse := validator.ObjectResponse{
			Diagnostics: diag.Diagnostics{},
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AwsKmsConfigValidator.ValidateObject(t.Context(), validatorRequest, &validatorResponse)

			if validatorResponse.Diagnostics.HasError() && !wantErr {
				t.Errorf("error = %v, wantErr %v", validatorResponse.Diagnostics.Errors(), wantErr)
			}
		})
	}
}
