package clusteroverloadsimulation_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteroverloadsimulation"
)

func TestResourceSchemaRequiresReplacement(t *testing.T) {
	var resp resource.SchemaResponse
	clusteroverloadsimulation.Resource().Schema(t.Context(), resource.SchemaRequest{}, &resp)

	for _, name := range []string{"project_id", "cluster_name"} {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("attribute %q not found or not a StringAttribute", name)
		}
		if !requiresReplacement(t.Context(), attr.PlanModifiers) {
			t.Errorf("attribute %q must require replacement", name)
		}
	}

	duration, ok := resp.Schema.Attributes["duration_seconds"].(schema.Int64Attribute)
	if !ok {
		t.Fatal("attribute \"duration_seconds\" not found or not an Int64Attribute")
	}
	if !requiresReplacement(t.Context(), duration.PlanModifiers) {
		t.Error("attribute \"duration_seconds\" must require replacement")
	}
	if len(duration.Validators) != 1 {
		t.Errorf("attribute \"duration_seconds\" must have one validator, got %d", len(duration.Validators))
	}
}

func requiresReplacement[T interface{ Description(context.Context) string }](ctx context.Context, modifiers []T) bool {
	for _, modifier := range modifiers {
		if strings.Contains(modifier.Description(ctx), "recreate") {
			return true
		}
	}
	return false
}
