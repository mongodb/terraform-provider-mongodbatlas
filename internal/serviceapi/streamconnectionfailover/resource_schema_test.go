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

// toAny re-keys an attribute map so the same assertions can run against the resource and the data
// source schemas, which use distinct (but structurally identical) attribute map types.
func toAny[V any](attrs map[string]V) map[string]any {
	out := make(map[string]any, len(attrs))
	for name, attr := range attrs {
		out[name] = attr
	}
	return out
}

// assertAwsRoleArn checks the MSK IAM authentication block: `role_arn` is exposed, while the
// `links` API boilerplate and the S3-only `test_bucket` stay ignored (see tools/codegen/config.yml).
func assertAwsRoleArn(t *testing.T, label string, authAttrs map[string]any) {
	t.Helper()
	var awsAttrs map[string]any
	switch aws := authAttrs["aws"].(type) {
	case schema.SingleNestedAttribute:
		awsAttrs = toAny(aws.Attributes)
	case dsschema.SingleNestedAttribute:
		awsAttrs = toAny(aws.Attributes)
	default:
		t.Fatalf("%s: authentication.aws not found or not a SingleNestedAttribute", label)
	}
	if _, ok := awsAttrs["role_arn"]; !ok {
		t.Errorf("%s: authentication.aws.role_arn not found", label)
	}
	if _, ok := awsAttrs["links"]; ok {
		t.Errorf("%s: authentication.aws.links must be ignored (see tools/codegen/config.yml)", label)
	}
	if _, ok := awsAttrs["test_bucket"]; ok {
		t.Errorf("%s: authentication.aws.test_bucket must be ignored: it is S3-only, and failover supports only Kafka and Cluster", label)
	}
}

// TestSchemaExposesAuthenticationAws guards that `authentication.aws.role_arn` is exposed by the
// resource and by both data sources, so MSK IAM auth can be configured and read back.
func TestSchemaExposesAuthenticationAws(t *testing.T) {
	ctx := context.Background()

	resAuth, ok := streamconnectionfailover.ResourceSchema(ctx).Attributes["authentication"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("resource: `authentication` not found or not a SingleNestedAttribute")
	}
	assertAwsRoleArn(t, "resource", toAny(resAuth.Attributes))

	dsAuth, ok := streamconnectionfailover.DataSourceSchema(ctx).Attributes["authentication"].(dsschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("data source: `authentication` not found or not a SingleNestedAttribute")
	}
	assertAwsRoleArn(t, "data source", toAny(dsAuth.Attributes))

	results, ok := streamconnectionfailover.PluralDataSourceSchema(ctx).Attributes["results"].(dsschema.ListNestedAttribute)
	if !ok {
		t.Fatal("plural data source: `results` not found or not a ListNestedAttribute")
	}
	pluralAuth, ok := results.NestedObject.Attributes["authentication"].(dsschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("plural data source: `results.authentication` not found or not a SingleNestedAttribute")
	}
	assertAwsRoleArn(t, "plural data source", toAny(pluralAuth.Attributes))
}

// TestSchemaAwsRoleArnIsRequired guards that an empty `aws = {}` block is rejected. `aws` stays
// optional so PLAIN/SCRAM/OAUTHBEARER remain valid, but with `test_bucket` ignored role_arn is the
// only field that applies to a failover connection, so it is required inside the block.
func TestSchemaAwsRoleArnIsRequired(t *testing.T) {
	resAuth, ok := streamconnectionfailover.ResourceSchema(context.Background()).Attributes["authentication"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("resource: `authentication` not found or not a SingleNestedAttribute")
	}
	aws, ok := resAuth.Attributes["aws"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("resource: `authentication.aws` not found or not a SingleNestedAttribute")
	}
	roleArn, ok := aws.Attributes["role_arn"].(schema.StringAttribute)
	if !ok {
		t.Fatal("resource: `authentication.aws.role_arn` not found or not a StringAttribute")
	}
	if !roleArn.IsRequired() {
		t.Error("resource: authentication.aws.role_arn must be required, so that `aws = {}` is rejected")
	}
}
