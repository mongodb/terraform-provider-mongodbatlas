package advancedclustertpf

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TODO: We temporarily use mongodbatlas_database_user instead of mongodbatlas_cluster to set up the initial environment
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			StateMover: stateMoverTemporaryTPFDatabaseUser,
		},
		{
			StateMover: stateMoverTemporaryV2,
		},
		{
			StateMover: stateMoverCluster,
		},
	}
}

func stateMoverTemporaryTPFDatabaseUser(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if !isSource(req, "database_user") {
		return
	}
	rawStateValue, err := req.SourceRawState.Unmarshal(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                 tftypes.String,
			"project_id":         tftypes.String,
			"auth_database_name": tftypes.String,
			"username":           tftypes.String,
			"password":           tftypes.String,
			"x509_type":          tftypes.String,
			"oidc_auth_type":     tftypes.String,
			"ldap_auth_type":     tftypes.String,
			"aws_iam_type":       tftypes.String,
			"roles": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"collection_name": tftypes.String,
					"database_name":   tftypes.String,
					"role_name":       tftypes.String,
				},
			}},
			"labels": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"key":   tftypes.String,
					"value": tftypes.String,
				},
			}},
			"scopes": tftypes.Set{ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"name": tftypes.String,
					"type": tftypes.String,
				},
			}},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Unmarshal Source State", err.Error())
		return
	}
	var rawState map[string]tftypes.Value
	if err := rawStateValue.As(&rawState); err != nil {
		resp.Diagnostics.AddError("Unable to Convert Source State", err.Error())
		return
	}
	var projectID *string // TODO: take username as the cluster name
	if err := rawState["project_id"].As(&projectID); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("project_id"), "Unable to Convert Source State", err.Error())
		return
	}
	var clusterName *string // TODO: take username as the cluster name
	if err := rawState["username"].As(&clusterName); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Unable to Convert Source State", err.Error())
		return
	}
	setMoveState(ctx, *projectID, *clusterName, resp)
}

func stateMoverTemporaryV2(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
}

func stateMoverCluster(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
}

func isSource(req resource.MoveStateRequest, resourceName string) bool {
	return req.SourceTypeName == "mongodbatlas_"+resourceName &&
		req.SourceSchemaVersion == 0 &&
		strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas")
}

func setMoveState(ctx context.Context, projectID, clusterName string, resp *resource.MoveStateResponse) {
	// TODO: timeout should be read from source if provided
	timeout := timeouts.Value{
		Object: types.ObjectValueMust(
			map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			},
			map[string]attr.Value{
				"create": types.StringValue("30m"),
				"update": types.StringValue("30m"),
				"delete": types.StringValue("30m"),
			}),
	}
	// TODO: we need to have a good state (all attributes known or null) but not need to be the final ones as Reas is called after
	tfNewModel, shouldReturn := mockedSDK(ctx, &resp.Diagnostics, timeout)
	if shouldReturn {
		return
	}
	// TODO: setting attributed needed by Read, confirm if ClusterID is needed
	tfNewModel.ProjectID = types.StringValue(projectID)
	tfNewModel.Name = types.StringValue(clusterName)
	resp.Diagnostics.Append(resp.TargetState.Set(ctx, tfNewModel)...)
}
