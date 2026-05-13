package apiresource //nolint:testpackage // in-package test: later tasks access unexported helpers

import (
	"regexp"
	"testing"
)

func TestLookupTypedCounterpart(t *testing.T) {
	tests := map[string]struct {
		path         string
		wantTypeName string
		preview      bool
		wantOK       bool
	}{
		"service account org-scoped GA": {
			path:         "/api/atlas/v2/orgs/abc123/serviceAccounts",
			preview:      false,
			wantOK:       true,
			wantTypeName: "mongodbatlas_service_account",
		},
		"service account with trailing slash": {
			path:         "/api/atlas/v2/orgs/abc123/serviceAccounts/",
			preview:      false,
			wantOK:       true,
			wantTypeName: "mongodbatlas_service_account",
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
			if ok && got.TypedTypeName != tc.wantTypeName {
				t.Fatalf("TypedTypeName = %q, want %q", got.TypedTypeName, tc.wantTypeName)
			}
		})
	}
}

func TestLookupTypedCounterpart_PreviewFilter(t *testing.T) {
	syntheticRegistry := []TypedCounterpart{
		{
			PathPattern:   regexp.MustCompile(`^/preview/resource$`),
			Preview:       new(true),
			TypedTypeName: "preview_only",
		},
		{
			PathPattern:   regexp.MustCompile(`^/ga/resource$`),
			Preview:       new(false),
			TypedTypeName: "ga_only",
		},
		{
			PathPattern:   regexp.MustCompile(`^/any/resource$`),
			TypedTypeName: "any_channel",
		},
	}

	tests := map[string]struct {
		path         string
		wantTypeName string
		preview      bool
		wantOK       bool
	}{
		"Preview:true entry, caller preview=true → match": {
			path:         "/preview/resource",
			preview:      true,
			wantOK:       true,
			wantTypeName: "preview_only",
		},
		"Preview:true entry, caller preview=false → no match": {
			path:    "/preview/resource",
			preview: false,
			wantOK:  false,
		},
		"Preview:false entry, caller preview=true → no match": {
			path:    "/ga/resource",
			preview: true,
			wantOK:  false,
		},
		"Preview:false entry, caller preview=false → match": {
			path:         "/ga/resource",
			preview:      false,
			wantOK:       true,
			wantTypeName: "ga_only",
		},
		"Preview:nil entry, caller preview=true → match (any)": {
			path:         "/any/resource",
			preview:      true,
			wantOK:       true,
			wantTypeName: "any_channel",
		},
		"Preview:nil entry, caller preview=false → match (any)": {
			path:         "/any/resource",
			preview:      false,
			wantOK:       true,
			wantTypeName: "any_channel",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := lookupIn(syntheticRegistry, tc.path, tc.preview)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got.TypedTypeName != tc.wantTypeName {
				t.Fatalf("TypedTypeName = %q, want %q", got.TypedTypeName, tc.wantTypeName)
			}
		})
	}
}
