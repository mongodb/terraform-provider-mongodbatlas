package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestValidCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantErr bool
	}{
		{
			name:    "Valid Value",
			cidr:    "192.0.0.0/28",
			wantErr: false,
		},
		{
			name:    "invalid value",
			cidr:    "12312321",
			wantErr: true,
		},
		{
			name:    "missing slash",
			cidr:    "192.0.0.8",
			wantErr: true,
		},
		{
			name:    "empty",
			cidr:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		val := tt.cidr
		wantErr := tt.wantErr
		cidrValidator := validate.CIDRValidator{}

		validatorRequest := validator.StringRequest{
			ConfigValue: types.StringValue(val),
		}

		validatorResponse := validator.StringResponse{
			Diagnostics: diag.Diagnostics{},
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cidrValidator.ValidateString(t.Context(), validatorRequest, &validatorResponse)

			if validatorResponse.Diagnostics.HasError() && !wantErr {
				t.Errorf("URL() error = %v, wantErr %v", validatorResponse.Diagnostics.Errors(), wantErr)
			}
		})
	}
}
