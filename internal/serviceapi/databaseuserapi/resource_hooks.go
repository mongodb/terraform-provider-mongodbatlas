package databaseuserapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// ResourceSchema adds manual-parity behavior that the config.yml override surface
// cannot express yet (override.Validators / override.PlanModifiers are not consumed
// by the generator today):
//   - password conflicts with the auth-type fields at plan time (manual uses ConflictsWith)
//   - project_id and username are identity fields: changes force replacement
//     (manual marks both RequiresReplace; without it a username change would PATCH a
//     different user's path instead of replacing the resource)

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	if passwordAttr, ok := baseSchema.Attributes["password"].(schema.StringAttribute); ok {
		passwordAttr.Validators = append(passwordAttr.Validators, stringvalidator.ConflictsWith(path.Expressions{
			path.MatchRelative().AtParent().AtName("x509_type"),
			path.MatchRelative().AtParent().AtName("ldap_auth_type"),
			path.MatchRelative().AtParent().AtName("aws_iam_type"),
		}...))
		baseSchema.Attributes["password"] = passwordAttr
	}
	for _, name := range []string{"project_id", "username"} {
		if attr, ok := baseSchema.Attributes[name].(schema.StringAttribute); ok {
			attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.RequiresReplace())
			baseSchema.Attributes[name] = attr
		}
	}
	return baseSchema
}
