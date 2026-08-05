package clusteroverloadsimulation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	for _, name := range []string{"project_id", "cluster_name"} {
		attr := s.Attributes[name].(schema.StringAttribute)
		attr.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		s.Attributes[name] = attr
	}

	duration := s.Attributes["duration_seconds"].(schema.Int64Attribute)
	duration.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
	duration.Validators = []validator.Int64{int64validator.OneOf(900, 3600, 28800, 86400)}
	s.Attributes["duration_seconds"] = duration

	return s
}
