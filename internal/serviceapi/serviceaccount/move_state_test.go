//nolint:testpackage // accesses unexported parseAPIResourceID + stateMover
package serviceaccount

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestParseAPIResourceID(t *testing.T) {
	tests := map[string]struct {
		id           string
		wantOrgID    string
		wantClientID string
		wantErr      bool
	}{
		"happy path": {
			id:           "/api/atlas/v2/orgs/653f8c1234567890abcdef12/serviceAccounts/mdb_sa_id_xyz",
			wantOrgID:    "653f8c1234567890abcdef12",
			wantClientID: "mdb_sa_id_xyz",
		},
		"with trailing slash": {
			id:           "/api/atlas/v2/orgs/abc/serviceAccounts/clientID/",
			wantOrgID:    "abc",
			wantClientID: "clientID",
		},
		"missing client id (list URL)": {
			id:      "/api/atlas/v2/orgs/abc/serviceAccounts",
			wantErr: true,
		},
		"completely unrelated": {
			id:      "/api/atlas/v2/groups/abc/clusters/xyz",
			wantErr: true,
		},
		"empty": {
			id:      "",
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			orgID, clientID, err := parseAPIResourceID(tc.id)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			if orgID != tc.wantOrgID || clientID != tc.wantClientID {
				t.Fatalf("got (%q,%q), want (%q,%q)", orgID, clientID, tc.wantOrgID, tc.wantClientID)
			}
		})
	}
}

func TestStateMover(t *testing.T) {
	ctx := context.Background()
	rs := ResourceSchema(ctx)

	t.Run("wrong source type is no-op", func(t *testing.T) {
		req := resource.MoveStateRequest{
			SourceTypeName:        "mongodbatlas_other",
			SourceProviderAddress: "registry.terraform.io/mongodb/mongodbatlas",
			SourceRawState:        rawJSON(`{"id":"/api/atlas/v2/orgs/abc/serviceAccounts/xyz"}`),
		}
		resp := &resource.MoveStateResponse{TargetState: nullTargetState(ctx, tfsdk.State{Schema: rs})}
		stateMover(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
	})

	t.Run("happy path sets org_id and client_id", func(t *testing.T) {
		req := resource.MoveStateRequest{
			SourceTypeName:        apiResourceSourceTypeName,
			SourceProviderAddress: "registry.terraform.io/mongodb/mongodbatlas",
			SourceRawState:        rawJSON(`{"id":"/api/atlas/v2/orgs/myOrg/serviceAccounts/mdb_xyz"}`),
		}
		resp := &resource.MoveStateResponse{TargetState: nullTargetState(ctx, tfsdk.State{Schema: rs})}
		stateMover(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics.Errors())
		}
		var got TFModel
		if diags := resp.TargetState.Get(ctx, &got); diags.HasError() {
			t.Fatalf("get target state: %v", diags)
		}
		if got.OrgId.ValueString() != "myOrg" {
			t.Fatalf("org_id = %q, want %q", got.OrgId.ValueString(), "myOrg")
		}
		if got.ClientId.ValueString() != "mdb_xyz" {
			t.Fatalf("client_id = %q, want %q", got.ClientId.ValueString(), "mdb_xyz")
		}
	})

	t.Run("malformed id reports error", func(t *testing.T) {
		req := resource.MoveStateRequest{
			SourceTypeName:        apiResourceSourceTypeName,
			SourceProviderAddress: "registry.terraform.io/mongodb/mongodbatlas",
			SourceRawState:        rawJSON(`{"id":"/garbage"}`),
		}
		resp := &resource.MoveStateResponse{TargetState: nullTargetState(ctx, tfsdk.State{Schema: rs})}
		stateMover(ctx, req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected error diagnostic, got none")
		}
	})

	t.Run("foreign provider is no-op", func(t *testing.T) {
		req := resource.MoveStateRequest{
			SourceTypeName:        apiResourceSourceTypeName,
			SourceProviderAddress: "registry.terraform.io/other/provider",
			SourceRawState:        rawJSON(`{"id":"/api/atlas/v2/orgs/abc/serviceAccounts/xyz"}`),
		}
		resp := &resource.MoveStateResponse{TargetState: nullTargetState(ctx, tfsdk.State{Schema: rs})}
		stateMover(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
		}
	})
}

func rawJSON(s string) *tfprotov6.RawState {
	return &tfprotov6.RawState{JSON: []byte(s)}
}

// nullTargetState creates a properly initialized null tfsdk.State for the given schema,
// matching how the framework initializes TargetState before calling StateMover.
func nullTargetState(ctx context.Context, s tfsdk.State) tfsdk.State {
	s.Raw = tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)
	return s
}
