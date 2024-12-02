package dsschema_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
)

func TestPaginatedDSSchema(t *testing.T) {
	expectedSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"page_num": schema.Int64Attribute{
				Optional: true,
			},
			"items_per_page": schema.Int64Attribute{
				Optional: true,
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
			},
			"results": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"instance_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}

	resultSchema := dsschema.PaginatedDSSchema(
		map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
		},
		map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"instance_name": schema.StringAttribute{
				Computed: true,
			},
		})

	if !reflect.DeepEqual(resultSchema, expectedSchema) {
		t.Errorf("created schema does not matched expected")
	}
}
