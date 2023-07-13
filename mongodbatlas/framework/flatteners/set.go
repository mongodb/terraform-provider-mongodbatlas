package flatteners

// StringSet accepts a `[]attr.Value` and returns a `basetypes.SetValue`. The
// return type automatically handles `SetNull` for empty results and coercing
// all element values to a string if there are any elements.
//
// nolint: contextcheck
// func StringSet(in []attr.Value) basetypes.SetValue {
// 	if len(in) == 0 {
// 		return types.SetNull(types.StringType)
// 	}
// 	return types.SetValueMust(types.StringType, in)
// }

// func StringList(in []attr.Value) basetypes.ListValue {
// 	if len(in) == 0 {
// 		return types.ListNull(types.StringType)
// 	}
// 	return types.ListValueMust(types.StringType, in)
// }
