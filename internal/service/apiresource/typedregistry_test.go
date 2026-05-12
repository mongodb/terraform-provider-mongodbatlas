package apiresource //nolint:testpackage // in-package test: later tasks access unexported helpers

import "testing"

func TestLookupTypedCounterpart(t *testing.T) {
	tests := map[string]struct {
		path    string
		wantTyp string
		preview bool
		wantOK  bool
	}{
		"service account org-scoped GA": {
			path:    "/api/atlas/v2/orgs/abc123/serviceAccounts",
			preview: false,
			wantOK:  true,
			wantTyp: "mongodbatlas_service_account",
		},
		"service account with trailing slash": {
			path:    "/api/atlas/v2/orgs/abc123/serviceAccounts/",
			preview: false,
			wantOK:  true,
			wantTyp: "mongodbatlas_service_account",
		},
		"service account child path does not match": {
			path:    "/api/atlas/v2/orgs/abc123/serviceAccounts/mdbsa-xyz",
			preview: false,
			wantOK:  false,
		},
		"unrelated path": {
			path:    "/api/atlas/v2/groups/abc/clusters",
			preview: false,
			wantOK:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := LookupTypedCounterpart(tc.path, tc.preview)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got.TypedTypeName != tc.wantTyp {
				t.Fatalf("TypedTypeName = %q, want %q", got.TypedTypeName, tc.wantTyp)
			}
		})
	}
}
