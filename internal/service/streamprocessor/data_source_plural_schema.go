package streamprocessor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"instance_name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable label that identifies the stream instance.",
				MarkdownDescription: "Human-readable label that identifies the stream instance.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: DSAttributes(false),
				},
				Description:         "Returns all Stream Processors within the specified stream instance.\n\nTo use this resource, the requesting API Key must have the Project Owner\n\nrole or Project Stream Processing Owner role.",
				MarkdownDescription: "Returns all Stream Processors within the specified stream instance.\n\nTo use this resource, the requesting API Key must have the Project Owner\n\nrole or Project Stream Processing Owner role.",
			},
		},
	}
}

type TFStreamProcessorsDSModel struct {
	ProjectID    types.String               `tfsdk:"project_id"`
	InstanceName types.String               `tfsdk:"instance_name"`
	Results      []TFStreamProcessorDSModel `tfsdk:"results"`
}
