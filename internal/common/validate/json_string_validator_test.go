package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestStringIsJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "Valid JSON",
			json: `{
				"test": "value"
			}`,
			wantErr: false,
		},
		{
			name:    "invalid value",
			json:    "12312321",
			wantErr: true,
		},
		{
			name: "missing comma",
			json: `{
				"test" "value"
			}`,
			wantErr: true,
		},
		{
			name:    "empty",
			json:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		val := tt.json
		wantErr := tt.wantErr
		cidrValidator := validate.JSONStringValidator{}

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
