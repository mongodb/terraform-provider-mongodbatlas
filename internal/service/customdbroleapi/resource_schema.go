// Code generated by terraform-provider-mongodbatlas using `make generate-resource`. DO NOT EDIT.

package customdbroleapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"actions": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of the individual privilege actions that the role grants.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Human-readable label that identifies the privilege action.",
						},
						"resources": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "List of resources on which you grant the action.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cluster": schema.BoolAttribute{
										Required:            true,
										MarkdownDescription: "Flag that indicates whether to grant the action on the cluster resource. If `true`, MongoDB Cloud ignores the **actions.resources.collection** and **actions.resources.db** parameters.",
									},
									"collection": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Human-readable label that identifies the collection on which you grant the action to one MongoDB user. If you don't set this parameter, you grant the action to all collections in the database specified in the **actions.resources.db** parameter. If you set `\"actions.resources.cluster\" : true`, MongoDB Cloud ignores this parameter.",
									},
									"db": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Human-readable label that identifies the database on which you grant the action to one MongoDB user. If you set `\"actions.resources.cluster\" : true`, MongoDB Cloud ignores this parameter.",
									},
								},
							},
						},
					},
				},
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"inherited_roles": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of the built-in roles that this custom role inherits.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"db": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Human-readable label that identifies the database on which someone grants the action to one MongoDB user.",
						},
						"role": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Human-readable label that identifies the role inherited. Set this value to `admin` for every role except `read` or `readWrite`.",
						},
					},
				},
			},
			"role_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.",
			},
		},
	}
}

type TFModel struct {
	Actions        types.List   `tfsdk:"actions"`
	GroupId        types.String `tfsdk:"group_id" autogeneration:"omitjson"`
	InheritedRoles types.Set    `tfsdk:"inherited_roles"`
	RoleName       types.String `tfsdk:"role_name" autogeneration:"omitjsonupdate"`
}
type TFActionsModel struct {
	Action    types.String `tfsdk:"action"`
	Resources types.List   `tfsdk:"resources"`
}
type TFResourcesModel struct {
	Collection types.String `tfsdk:"collection"`
	Db         types.String `tfsdk:"db"`
	Cluster    types.Bool   `tfsdk:"cluster"`
}
type TFInheritedRolesModel struct {
	Db   types.String `tfsdk:"db"`
	Role types.String `tfsdk:"role"`
}
