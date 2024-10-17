//nolint:gocritic
package flexcluster

// TODO: `ctx` parameter and `diags` return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
// func NewTFModel(ctx context.Context, apiResp *admin.FlexCluster) (*TFModel, diag.Diagnostics) {
// 	// complexAttr, diagnostics := types.ListValueFrom(ctx, InnerObjectType, newTFComplexAttrModel(apiResp.ComplexAttr))
// 	// if diagnostics.HasError() {
// 	// 	return nil, diagnostics
// 	// }
// 	return &TFModel{}, nil
// }

// TODO: If SDK defined different models for create and update separate functions will need to be defined.
// TODO: `ctx` parameter and `diags` in return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
// func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.FlexCluster, diag.Diagnostics) {
//     // var tfList []complexArgumentData
// 	// resp.Diagnostics.Append(plan.ComplexArgument.ElementsAs(ctx, &tfList, false)...)
// 	// if resp.Diagnostics.HasError() {
// 	// 	return nil, diagnostics
// 	// }
// 	return &admin.FlexCluster{}, nil
// }
