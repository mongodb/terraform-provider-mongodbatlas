package streamprocessor

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &StreamProccesorDS{}
var _ datasource.DataSourceWithConfigure = &StreamProccesorDS{}

func PluralDataSource() datasource.DataSource {
	return &streamProcessorsDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", StreamProcessorName),
		},
	}
}

type streamProcessorsDS struct {
	config.DSCommon
}

func (d *streamProcessorsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"instance_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: DSAttributes(false),
				},
				MarkdownDescription: "Returns all Stream Processors within the specified stream instance.\n\nTo use this resource, the requesting API Key must have the Project Owner\n\nrole or Project Stream Processing Owner role.",
			},
		},
	}
	conversion.UpdateSchemaDescription(&ds1)

	requiredFields := []string{"project_id", "instance_name"}
	desc := "Returns all Stream Processors within the specified stream instance.\n\nTo use this resource, the requesting API Key must have the Project Owner\n\nrole or Project Stream Processing Owner role."
	ds2 := conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, nil, desc, false)
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
}

type TFStreamProcessorsDSModel struct {
	ProjectID    types.String               `tfsdk:"project_id"`
	InstanceName types.String               `tfsdk:"instance_name"`
	Results      []TFStreamProcessorDSModel `tfsdk:"results"`
}
