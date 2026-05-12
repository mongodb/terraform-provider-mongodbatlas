package apiresource

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TypedCounterpart describes a typed resource that supersedes a particular
// api_resource path. Used to (a) emit migration warnings at plan time and
// (b) document the migration path.
type TypedCounterpart struct {
	PathPattern   *regexp.Regexp
	Preview       *bool // nil = match either; pointer to bool = require equality
	TypedTypeName string
	DocsAnchor    string
}

// typedRegistry is a small, hand-curated table. The first matching entry wins,
// so more-specific patterns must be listed before less-specific ones.
var typedRegistry = []TypedCounterpart{
	{
		PathPattern:   regexp.MustCompile(`^/api/atlas/v2/orgs/[^/]+/serviceAccounts/?$`),
		TypedTypeName: "mongodbatlas_service_account",
		DocsAnchor:    "service-account",
	},
}

// LookupTypedCounterpart returns the typed resource (if any) that supersedes
// the given api_resource path.
func LookupTypedCounterpart(path string, preview bool) (TypedCounterpart, bool) {
	return lookupIn(typedRegistry, path, preview)
}

func lookupIn(registry []TypedCounterpart, path string, preview bool) (TypedCounterpart, bool) {
	for _, entry := range registry {
		if entry.Preview != nil && *entry.Preview != preview {
			continue
		}
		if entry.PathPattern.MatchString(path) {
			return entry, true
		}
	}
	return TypedCounterpart{}, false
}

// emitTypedCounterpartWarning consults the registry and appends a warning
// diagnostic to diags if a typed resource supersedes the given path.
func emitTypedCounterpartWarning(_ context.Context, path string, preview bool, diags *diag.Diagnostics) {
	entry, ok := LookupTypedCounterpart(path, preview)
	if !ok {
		return
	}
	diags.AddWarning(
		"Typed resource available",
		fmt.Sprintf(
			"A typed resource is available for this endpoint: %s. "+
				"Migrate with a `moved` block to get typed schema, validation, "+
				"and IDE completions.",
			entry.TypedTypeName,
		),
	)
	// entry.DocsAnchor is intentionally unused for now: it identifies the
	// section of the migration guide a future warning revision should link
	// to once that guide is published.
	_ = entry.DocsAnchor
}
