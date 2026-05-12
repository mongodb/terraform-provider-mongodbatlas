package apiresource

import "regexp"

// TypedCounterpart describes a typed resource that supersedes a particular
// api_resource path. Used to (a) emit migration warnings at plan time and
// (b) document the migration path.
type TypedCounterpart struct {
	PathPattern   *regexp.Regexp
	Preview       *bool // nil = match either; pointer to bool = require equality
	TypedTypeName string
	DocsAnchor    string
}

// typedRegistry is a small, hand-curated table. Order is irrelevant — the
// first matching entry wins.
var typedRegistry = []TypedCounterpart{
	{
		PathPattern:   regexp.MustCompile(`^/api/atlas/v2/orgs/[^/]+/serviceAccounts/?$`),
		Preview:       nil,
		TypedTypeName: "mongodbatlas_service_account",
		DocsAnchor:    "service-account",
	},
}

// LookupTypedCounterpart returns the typed resource (if any) that supersedes
// the given api_resource path.
func LookupTypedCounterpart(path string, preview bool) (TypedCounterpart, bool) {
	for _, entry := range typedRegistry {
		if entry.Preview != nil && *entry.Preview != preview {
			continue
		}
		if entry.PathPattern.MatchString(path) {
			return entry, true
		}
	}
	return TypedCounterpart{}, false
}
