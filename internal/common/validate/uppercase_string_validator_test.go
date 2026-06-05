package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestValidUppercaseString(t *testing.T) {
	tests := []struct {
		name       string
		value      types.String
		wantHasErr bool
	}{
		{
			name:       "uppercase value",
			value:      types.StringValue("US_EAST_1"),
			wantHasErr: false,
		},
		{
			name:       "lowercase value",
			value:      types.StringValue("us_east_1"),
			wantHasErr: true,
		},
		{
			name:       "empty value",
			value:      types.StringValue(""),
			wantHasErr: false,
		},
		{
			name:       "null value",
			value:      types.StringNull(),
			wantHasErr: false,
		},
		{
			name:       "unknown value",
			value:      types.StringUnknown(),
			wantHasErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				ConfigValue: tc.value,
			}
			resp := validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			validate.ValidUppercaseString().ValidateString(t.Context(), req, &resp)
			if resp.Diagnostics.HasError() != tc.wantHasErr {
				t.Fatalf("unexpected diagnostics error state, got=%v want=%v", resp.Diagnostics.HasError(), tc.wantHasErr)
			}
		})
	}
}
