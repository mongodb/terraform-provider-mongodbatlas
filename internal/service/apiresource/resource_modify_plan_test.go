//nolint:testpackage // accesses unexported helper emitTypedCounterpartWarning
package apiresource

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestEmitTypedCounterpartWarning(t *testing.T) {
	tests := map[string]struct {
		path        string
		preview     bool
		wantWarning bool
	}{
		"matching SA path emits warning": {
			path:        "/api/atlas/v2/orgs/abc/serviceAccounts",
			preview:     false,
			wantWarning: true,
		},
		"non-matching path silent": {
			path:        "/api/atlas/v2/groups/abc/clusters",
			preview:     false,
			wantWarning: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var diags diag.Diagnostics
			emitTypedCounterpartWarning(context.Background(), tc.path, tc.preview, &diags)
			gotWarning := diags.WarningsCount() > 0
			if gotWarning != tc.wantWarning {
				t.Fatalf("warning emitted=%v, want=%v; diags=%v", gotWarning, tc.wantWarning, diags)
			}
			if tc.wantWarning {
				w := diags.Warnings()[0]
				if !strings.Contains(w.Detail(), "mongodbatlas_service_account") {
					t.Fatalf("warning detail missing typed name: %q", w.Detail())
				}
			}
		})
	}
}
