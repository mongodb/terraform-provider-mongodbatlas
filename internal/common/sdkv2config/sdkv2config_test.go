package sdkv2config_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
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
		raw     cty.Value
		name    string
		want    string
		changed []string
		null    bool
		unknown bool
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.StringVal("x")}),
			null: true,
		},
		{
			name: "json null",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.NullVal(cty.String)}),
			null: true,
		},
		{
			name:    "json null and HasChange is Removed",
			raw:     obj(map[string]cty.Value{"operations_contact": cty.NullVal(cty.String)}),
			null:    true,
			changed: []string{"operations_contact"},
		},
		{
			name:    "unknown",
			raw:     obj(map[string]cty.Value{"operations_contact": cty.UnknownVal(cty.String)}),
			unknown: true,
		},
		{
			name: "empty string is a value",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.StringVal("")}),
		},
		{
			name: "set string",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.StringVal("a@b.com")}),
			want: "a@b.com",
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"operations_contact": cty.True}),
			null: true,
		},
		{
			name:    "unknown parent yields unknown",
			raw:     cty.UnknownVal(cty.EmptyObject),
			unknown: true,
		},
		{
			name: "null parent yields null",
			raw:  cty.NullVal(cty.EmptyObject),
			null: true,
		},
		{
			name: "non-object parent yields null",
			raw:  cty.StringVal("nope"),
			null: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.String(view(tc.raw, tc.changed...), "operations_contact")
			assert.Equal(t, tc.null, got.IsNull())
			assert.Equal(t, tc.unknown, got.IsUnknown())
			assert.Equal(t, len(tc.changed) > 0, got.HasChange())
			assert.Equal(t, tc.want, got.ValueString())
			assert.Equal(t, tc.null && len(tc.changed) > 0, got.Removed())
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw     cty.Value
		name    string
		changed []string
		null    bool
		unknown bool
		wantVal bool
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.True}),
			null: true,
		},
		{
			name: "json null",
			raw:  obj(map[string]cty.Value{"auto_export_enabled": cty.NullVal(cty.Bool)}),
			null: true,
		},
		{
			name:    "json null and changed PATCHes false via ValueBool",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.NullVal(cty.Bool)}),
			null:    true,
			changed: []string{"auto_export_enabled"},
		},
		{
			name:    "unknown",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.UnknownVal(cty.Bool)}),
			unknown: true,
		},
		{
			name: "false is a value",
			raw:  obj(map[string]cty.Value{"auto_export_enabled": cty.False}),
		},
		{
			name:    "true",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.True}),
			wantVal: true,
		},
		{
			name:    "set and changed",
			raw:     obj(map[string]cty.Value{"auto_export_enabled": cty.True}),
			wantVal: true,
			changed: []string{"auto_export_enabled"},
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"auto_export_enabled": cty.StringVal("true")}),
			null: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.Bool(view(tc.raw, tc.changed...), "auto_export_enabled")
			assert.Equal(t, tc.null, got.IsNull())
			assert.Equal(t, tc.unknown, got.IsUnknown())
			assert.Equal(t, len(tc.changed) > 0, got.HasChange())
			assert.Equal(t, tc.wantVal, got.ValueBool())
			assert.Equal(t, !tc.unknown && !tc.null, got.Set())
			assert.Equal(t, tc.null && len(tc.changed) > 0, got.Removed())
			assert.Equal(t, !tc.unknown && (!tc.null || len(tc.changed) > 0), got.SetOrRemoved())
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw     cty.Value
		name    string
		wantVal int64
		null    bool
		unknown bool
	}{
		{
			name: "missing attr",
			raw:  obj(map[string]cty.Value{"other": cty.NumberIntVal(1)}),
			null: true,
		},
		{
			name:    "unknown omitted Optional+Computed",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.UnknownVal(cty.Number)}),
			unknown: true,
		},
		{
			name: "zero is a value",
			raw:  obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberIntVal(0)}),
		},
		{
			name:    "set int",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberIntVal(7)}),
			wantVal: 7,
		},
		{
			name: "wrong type",
			raw:  obj(map[string]cty.Value{"reference_hour_of_day": cty.StringVal("7")}),
			null: true,
		},
		{
			name:    "fractional number is unknown, not truncated",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberFloatVal(1.5)}),
			unknown: true,
		},
		{
			name:    "out-of-range number is unknown, not clamped",
			raw:     obj(map[string]cty.Value{"reference_hour_of_day": cty.NumberFloatVal(1e30)}),
			unknown: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := sdkv2config.Int64(view(tc.raw), "reference_hour_of_day")
			assert.Equal(t, tc.null, got.IsNull())
			assert.Equal(t, tc.unknown, got.IsUnknown())
			assert.False(t, got.HasChange())
			assert.Equal(t, tc.wantVal, got.ValueInt64())
			assert.Equal(t, !tc.unknown && !tc.null, got.Set())
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
