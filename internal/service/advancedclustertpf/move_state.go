package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{
		{
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
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
				tfNewModel, shouldReturn := mockedSDK(ctx, &resp.Diagnostics, timeout)
				if shouldReturn {
					return
				}
				resp.Diagnostics.Append(resp.TargetState.Set(ctx, tfNewModel)...)
			},
		},
	}
}
