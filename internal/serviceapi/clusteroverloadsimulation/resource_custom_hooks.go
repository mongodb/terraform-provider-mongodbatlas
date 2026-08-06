package clusteroverloadsimulation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	for _, name := range []string{"project_id", "cluster_name"} {
		attr, ok := s.Attributes[name].(schema.StringAttribute)
		if !ok {
			continue
		}
		attr.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		s.Attributes[name] = attr
	}

	if duration, ok := s.Attributes["duration_seconds"].(schema.Int64Attribute); ok {
		duration.PlanModifiers = []planmodifier.Int64{int64planmodifier.RequiresReplace()}
		s.Attributes["duration_seconds"] = duration
	}

	return s
}
