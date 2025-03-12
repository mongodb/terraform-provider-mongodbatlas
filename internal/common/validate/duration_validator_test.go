package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestValidDurationBetween(t *testing.T) {
	tests := []struct {
		name       string
		minutes    string
		maxMinutes int
		minMinutes int
		wantErr    bool
	}{
		{
			name:       "valid minutes",
			minutes:    "11m",
			minMinutes: 10,
			maxMinutes: 12,
			wantErr:    false,
		},
		{
			name:       "out of range",
			minutes:    "11h45m",
			minMinutes: 10,
			maxMinutes: 12,
			wantErr:    true,
		},
		{
			name:       "unvalid minutes",
			minutes:    "1m",
			minMinutes: 10,
			maxMinutes: 12,
			wantErr:    true,
		},
		{
			name:       "max minutes smaller than min minutes",
			minutes:    "11",
			minMinutes: 10,
			maxMinutes: 1,
			wantErr:    true,
		},
		{
			name:       "negative number",
			minutes:    "-11",
			minMinutes: 10,
			maxMinutes: 1,
			wantErr:    true,
		},
		{
			name:       "empty",
			minutes:    "",
			minMinutes: 10,
			maxMinutes: 12,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		wantErr := tt.wantErr
		cidrValidator := validate.DurationValidator{
			MinMinutes: tt.minMinutes,
			MaxMinutes: tt.maxMinutes,
		}

		val := tt.minutes
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
