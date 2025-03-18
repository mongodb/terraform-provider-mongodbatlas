package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestValidIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{
			name:    "Valid IP v4",
			ip:      "192.0.2.1",
			wantErr: false,
		},
		{
			name:    "Valid IP v6",
			ip:      "2001:db8::68",
			wantErr: false,
		},
		{
			name:    "Valid IP v6",
			ip:      "::ffff:192.0.2.1",
			wantErr: false,
		},
		{
			name:    "invalid IP",
			ip:      "12312321",
			wantErr: true,
		},
		{
			name:    "empty",
			ip:      "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		val := tt.ip
		wantErr := tt.wantErr
		cidrValidator := validate.IPValidator{}

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
