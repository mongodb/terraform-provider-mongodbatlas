package sdkv2config_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/sdkv2config"
)

type fakeView struct {
	raw     cty.Value
	changes map[string]bool
}

func (f fakeView) GetRawConfig() cty.Value { return f.raw }

func (f fakeView) HasChange(key string) bool { return f.changes[key] }

func view(raw cty.Value, changed ...string) fakeView {
	changes := make(map[string]bool, len(changed))
	for _, name := range changed {
		changes[name] = true
	}
	return fakeView{raw: raw, changes: changes}
}

func obj(attrs map[string]cty.Value) cty.Value {
	return cty.ObjectVal(attrs)
}

func TestString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		raw     cty.Value
		want    types.String
		changed []string
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.StringVal("x")}),
			want: types.StringNull(),
		},
		{
			name: "json null",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.NullVal(cty.String)}),
			want: types.StringNull(),
		},
		{
			name: "unknown",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.UnknownVal(cty.String)}),
			want: types.StringUnknown(),
		},
		{
			name: "empty string is a value",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.StringVal("")}),
			want: types.StringValue(""),
		},
		{
			name: "set string",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.StringVal("a@b.com")}),
			want: types.StringValue("a@b.com"),
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.True}),
			want: types.StringNull(),
		},
		{
			name:    "null and HasChange",
			raw:     obj(map[string]cty.Value{"operations_contact": cty.NullVal(cty.String)}),
			want:    types.StringNull(),
			changed: []string{"operations_contact"},
		},
		{
			name: "unknown parent yields unknown",
			raw:  cty.UnknownVal(cty.EmptyObject),
			want: types.StringUnknown(),
		},
		{
			name: "null parent yields null",
			raw:  cty.NullVal(cty.EmptyObject),
			want: types.StringNull(),
		},
		{
			name: "non-object parent yields null",
			raw:  cty.StringVal("nope"),
			want: types.StringNull(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.String(view(tc.raw, tc.changed...), "operations_contact")
			assertEqualString(t, tc.want, got)
			assert.Equal(t, len(tc.changed) > 0, got.HasChange())
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw     cty.Value
		name    string
		changed []string
		want    types.Bool
		wantVal bool
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.True}),
			want: types.BoolNull(),
		},
		{
			name: "unknown",
			raw:  obj(map[string]cty.Value{"auto_export_enabled": cty.UnknownVal(cty.Bool)}),
			want: types.BoolUnknown(),
		},
		{
			name:    "false is a value",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.False}),
			want:    types.BoolValue(false),
			wantVal: false,
		},
		{
			name:    "true",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.True}),
			want:    types.BoolValue(true),
			wantVal: true,
		},
		{
			name:    "null and changed PATCHes false via ValueBool",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.NullVal(cty.Bool)}),
			want:    types.BoolNull(),
			changed: []string{"auto_export_enabled"},
			wantVal: false,
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"auto_export_enabled": cty.StringVal("true")}),
			want: types.BoolNull(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.Bool(view(tc.raw, tc.changed...), "auto_export_enabled")
			assertEqualBool(t, tc.want, got)
			assert.Equal(t, len(tc.changed) > 0, got.HasChange())
			if !got.IsUnknown() {
				assert.Equal(t, tc.wantVal, got.ValueBool())
			}
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw     cty.Value
		name    string
		want    types.Int64
		wantVal int64
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.NumberIntVal(1)}),
			want: types.Int64Null(),
		},
		{
			name: "unknown omitted Optional+Computed",
			raw:  obj(map[string]cty.Value{"reference_hour_of_day": cty.UnknownVal(cty.Number)}),
			want: types.Int64Unknown(),
		},
		{
			name:    "zero is a value",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberIntVal(0)}),
			want:    types.Int64Value(0),
			wantVal: 0,
		},
		{
			name:    "set int",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberIntVal(7)}),
			want:    types.Int64Value(7),
			wantVal: 7,
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"reference_hour_of_day": cty.StringVal("7")}),
			want: types.Int64Null(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.Int64(view(tc.raw), "reference_hour_of_day")
			assertEqualInt64(t, tc.want, got)
			assert.False(t, got.HasChange())
			if !got.IsNull() && !got.IsUnknown() {
				assert.Equal(t, tc.wantVal, got.ValueInt64())
			}
		})
	}
}

func TestCollectionEmpty(t *testing.T) {
	t.Parallel()
	listTy := cty.List(cty.EmptyObject)
	tests := []struct {
		raw  cty.Value
		name string
		want bool
	}{
		{
			name: "null root is not an HCL omit",
			raw:  cty.NullVal(cty.EmptyObject),
			want: false,
		},
		{
			name: "unknown root is not an HCL omit",
			raw:  cty.UnknownVal(cty.EmptyObject),
			want: false,
		},
		{
			name: "attr null",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.NullVal(listTy)}),
			want: true,
		},
		{
			name: "attr unknown",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.UnknownVal(listTy)}),
			want: true,
		},
		{
			name: "length 0",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.ListValEmpty(cty.EmptyObject)}),
			want: true,
		},
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.ListValEmpty(cty.EmptyObject)}),
			want: true,
		},
		{
			name: "known length greater than 0",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.ListVal([]cty.Value{cty.EmptyObjectVal})}),
			want: false,
		},
		{
			name: "empty set",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.SetValEmpty(cty.String)}),
			want: true,
		},
		{
			name: "non-empty set",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.SetVal([]cty.Value{cty.StringVal("DAILY")})}),
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, sdkv2config.CollectionEmpty(view(tc.raw), "copy_settings"))
		})
	}
}

func TestNestedCollectionLen(t *testing.T) {
	t.Parallel()
	freqSet := cty.Set(cty.String)
	entryTy := cty.Object(map[string]cty.Type{"frequencies": freqSet})
	listTy := cty.List(entryTy)

	entryWith := func(freqs cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{"frequencies": freqs})
	}
	listOf := func(entries ...cty.Value) cty.Value {
		return obj(map[string]cty.Value{"copy_settings": cty.ListVal(entries)})
	}

	tests := []struct {
		raw   cty.Value
		name  string
		index int
		want  int
	}{
		{
			name: "omitted list",
			raw:  obj(map[string]cty.Value{"other": cty.ListValEmpty(entryTy)}),
			want: 0,
		},
		{
			name: "unknown list",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.UnknownVal(listTy)}),
			want: 0,
		},
		{
			name: "empty list",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.ListValEmpty(entryTy)}),
			want: 0,
		},
		{
			name:  "index out of range",
			raw:   listOf(entryWith(cty.SetVal([]cty.Value{cty.StringVal("DAILY")}))),
			index: 1,
			want:  0,
		},
		{
			name: "nested attr omitted",
			raw:  listOf(cty.EmptyObjectVal),
			want: 0,
		},
		{
			name: "nested attr unknown",
			raw:  listOf(entryWith(cty.UnknownVal(freqSet))),
			want: 0,
		},
		{
			name: "nested attr empty set",
			raw:  listOf(entryWith(cty.SetValEmpty(cty.String))),
			want: 0,
		},
		{
			name: "two frequencies",
			raw:  listOf(entryWith(cty.SetVal([]cty.Value{cty.StringVal("DAILY"), cty.StringVal("WEEKLY")}))),
			want: 2,
		},
		{
			name:  "second entry",
			raw:   listOf(entryWith(cty.SetValEmpty(cty.String)), entryWith(cty.SetVal([]cty.Value{cty.StringVal("DAILY")}))),
			index: 1,
			want:  1,
		},
		{
			name: "parent is a set so Index is not used",
			raw:  obj(map[string]cty.Value{"copy_settings": cty.SetVal([]cty.Value{cty.EmptyObjectVal})}),
			want: 0,
		},
		{
			name: "null root",
			raw:  cty.NullVal(cty.EmptyObject),
			want: 0,
		},
		{
			name: "unknown root",
			raw:  cty.UnknownVal(cty.EmptyObject),
			want: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.NestedCollectionLen(view(tc.raw), "copy_settings", tc.index, "frequencies")
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStringDoesNotPanicOnMissingAttr(t *testing.T) {
	t.Parallel()
	require.NotPanics(t, func() {
		_ = sdkv2config.String(view(obj(map[string]cty.Value{"x": cty.StringVal("y")})), "missing")
	})
}

func assertEqualString(t *testing.T, want types.String, got sdkv2config.StringValue) {
	t.Helper()
	assert.Equal(t, want.IsNull(), got.IsNull())
	assert.Equal(t, want.IsUnknown(), got.IsUnknown())
	if !want.IsNull() && !want.IsUnknown() {
		assert.Equal(t, want.ValueString(), got.ValueString())
	}
}

func assertEqualBool(t *testing.T, want types.Bool, got sdkv2config.BoolValue) {
	t.Helper()
	assert.Equal(t, want.IsNull(), got.IsNull())
	assert.Equal(t, want.IsUnknown(), got.IsUnknown())
	if !want.IsNull() && !want.IsUnknown() {
		assert.Equal(t, want.ValueBool(), got.ValueBool())
	}
}

func assertEqualInt64(t *testing.T, want types.Int64, got sdkv2config.Int64Value) {
	t.Helper()
	assert.Equal(t, want.IsNull(), got.IsNull())
	assert.Equal(t, want.IsUnknown(), got.IsUnknown())
	if !want.IsNull() && !want.IsUnknown() {
		assert.Equal(t, want.ValueInt64(), got.ValueInt64())
	}
}
