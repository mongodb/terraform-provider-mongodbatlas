package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20241023002/admin"
)

// NewTFModel temporarily creates a fake model to be able to run the tests
func NewTFModel(ctx context.Context, apiResp *admin.ClusterDescription20240805, diags *diag.Diagnostics) *TFModel {
	model := &TFModel{
		BiConnector: NewObjectTF(ctx, diags, BiConnectorObjType, TFBiConnectorModel{
			ReadPreference: types.StringValue("PREFERENCE"),
			Enabled:        types.BoolValue(false),
		}),
		ConnectionStrings: NewObjectTF(ctx, diags, ConnectionStringsObjType, TFConnectionStringsModel{
			AwsPrivateLink:    types.MapNull(types.StringType),
			AwsPrivateLinkSrv: types.MapNull(types.StringType),
			Private:           types.StringValue("Private"),
			PrivateEndpoint:   types.ListNull(types.StringType),
			PrivateSrv:        types.StringValue("PrivateSrv"),
			Standard:          types.StringValue("Standard"),
			StandardSrv:       types.StringValue("StandardSrv"),
		}),
		MongoDBEmployeeAccessGrant: NewObjectTF(ctx, diags, MongoDbemployeeAccessGrantObjType, TFMongoDbemployeeAccessGrantModel{
			ExpirationTime: types.StringValue("2024-10-31T00:00:00Z"),
			GrantType:      types.StringValue("TYPE"),
		}),
		ReplicationSpecs: types.ListNull(ReplicationSpecsObjType),
		Tags:             types.MapNull(types.StringType),
		Labels:           types.ListNull(LabelsObjType),
	}

	if diags.HasError() {
		return nil
	}
	return model
}

func NewObjectTF(ctx context.Context, diags *diag.Diagnostics, objectType types.ObjectType, model any) types.Object {
	object, moreDiags := types.ObjectValueFrom(ctx, objectType.AttributeTypes(), &model)
	diags.Append(moreDiags...)
	if diags.HasError() {
		return types.ObjectNull(objectType.AttributeTypes())
	}
	return object
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.ClusterDescription20240805, diag.Diagnostics) {
	return &admin.ClusterDescription20240805{}, nil
}
