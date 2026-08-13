package streamconnectionfailover_test

import (
	"context"
	"strings"
	"testing"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/streamconnectionfailover"
)

// TestResourceSchemaImmutableFields guards the ResourceSchema hook: `type` and `region` are immutable
// on the failover connection PATCH, so both must force replacement. This asserts the RequiresReplace
// plan modifier survives schema generation.
func TestResourceSchemaImmutableFields(t *testing.T) {
	var resp resource.SchemaResponse
	streamconnectionfailover.Resource().Schema(context.Background(), resource.SchemaRequest{}, &resp)

	for _, name := range []string{"type", "region"} {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("attribute %q not found or not a StringAttribute", name)
		}
		requiresReplace := false
		for _, pm := range attr.PlanModifiers {
			if strings.Contains(pm.Description(context.Background()), "recreate") {
				requiresReplace = true
				break
			}
		}
		if !requiresReplace {
			t.Errorf("attribute %q must have a RequiresReplace plan modifier (it is immutable on update)", name)
		}
	}
}

// authSecrets are the Kafka authentication fields that carry `format: password` in the API spec.
var authSecrets = []string{"client_secret", "password", "ssl_key", "ssl_key_password"}

// assertSecretsSensitive checks every field in authSecrets is marked sensitive.
func assertSecretsSensitive[T interface{ IsSensitive() bool }](t *testing.T, label string, attrs map[string]T) {
	t.Helper()
	for _, name := range authSecrets {
		attr, ok := attrs[name]
		if !ok {
			t.Errorf("%s: authentication.%s not found", label, name)
			continue
		}
		if !attr.IsSensitive() {
			t.Errorf("%s: authentication.%s must be sensitive", label, name)
		}
	}
}

// TestSchemaAuthSecretsAreSensitive guards that the Kafka authentication secrets are sensitive in the
// resource and in both data sources.
func TestSchemaAuthSecretsAreSensitive(t *testing.T) {
	ctx := context.Background()

	resAuth, ok := streamconnectionfailover.ResourceSchema(ctx).Attributes["authentication"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("resource: `authentication` not found or not a SingleNestedAttribute")
	}
	assertSecretsSensitive(t, "resource", resAuth.Attributes)

	dsAuth, ok := streamconnectionfailover.DataSourceSchema(ctx).Attributes["authentication"].(dsschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("data source: `authentication` not found or not a SingleNestedAttribute")
	}
	assertSecretsSensitive(t, "data source", dsAuth.Attributes)

	results, ok := streamconnectionfailover.PluralDataSourceSchema(ctx).Attributes["results"].(dsschema.ListNestedAttribute)
	if !ok {
		t.Fatal("plural data source: `results` not found or not a ListNestedAttribute")
	}
	pluralAuth, ok := results.NestedObject.Attributes["authentication"].(dsschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("plural data source: `results.authentication` not found or not a SingleNestedAttribute")
	}
	assertSecretsSensitive(t, "plural data source", pluralAuth.Attributes)
}
