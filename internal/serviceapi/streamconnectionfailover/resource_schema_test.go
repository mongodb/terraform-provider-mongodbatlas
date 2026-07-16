package streamconnectionfailover_test

import (
	"context"
	"strings"
	"testing"

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
