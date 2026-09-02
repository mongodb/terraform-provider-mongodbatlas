package cloudbackupschedule

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCopySettingsWithEmptyFrequencies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantUnchanged map[string]any
		name          string
		copySettings  []any
		wantFreqLens  []int
		wantChanged   bool
	}{
		{
			name: "leftover frequencies",
			copySettings: []any{
				map[string]any{
					"cloud_provider":           "AWS",
					"region_name":              "US_WEST_1",
					"zone_id":                  "zone-1",
					"should_copy_oplogs":       true,
					"copy_policy_items":        []any{map[string]any{"frequency_type": "daily"}},
					"last_number_of_snapshots": 0,
					"frequencies":              testFreqSet("HOURLY", "DAILY"),
				},
			},
			wantChanged:  true,
			wantFreqLens: []int{0},
			wantUnchanged: map[string]any{
				"cloud_provider":           "AWS",
				"region_name":              "US_WEST_1",
				"zone_id":                  "zone-1",
				"should_copy_oplogs":       true,
				"copy_policy_items":        []any{map[string]any{"frequency_type": "daily"}},
				"last_number_of_snapshots": 0,
			},
		},
		{
			name: "missing frequencies",
			copySettings: []any{
				map[string]any{
					"cloud_provider": "AWS",
				},
			},
			wantChanged:  false,
			wantFreqLens: []int{0},
		},
		{
			name: "two entries only index 1 has leftover",
			copySettings: []any{
				map[string]any{
					"cloud_provider": "AWS",
					"region_name":    "US_EAST_1",
					"frequencies":    testFreqSet(),
				},
				map[string]any{
					"cloud_provider": "AWS",
					"region_name":    "US_WEST_1",
					"frequencies":    testFreqSet("DAILY"),
				},
			},
			wantChanged:  true,
			wantFreqLens: []int{0, 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			origFreqLens := make([]int, len(tc.copySettings))
			for i, raw := range tc.copySettings {
				entry := raw.(map[string]any)
				origFreqLens[i] = setLen(entry["frequencies"])
			}

			rewritten, changed := copySettingsWithEmptyFrequencies(tc.copySettings)
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}
			if len(rewritten) != len(tc.copySettings) {
				t.Fatalf("rewritten len = %d, want %d", len(rewritten), len(tc.copySettings))
			}
			for i, raw := range rewritten {
				entry, ok := raw.(map[string]any)
				if !ok {
					t.Fatalf("rewritten[%d] is %T, want map[string]any", i, raw)
				}
				gotLen := setLen(entry["frequencies"])
				if gotLen != tc.wantFreqLens[i] {
					t.Fatalf("rewritten[%d].frequencies len = %d, want %d", i, gotLen, tc.wantFreqLens[i])
				}
				orig := tc.copySettings[i].(map[string]any)
				if setLen(orig["frequencies"]) != origFreqLens[i] {
					t.Fatalf("input[%d].frequencies was mutated", i)
				}
			}
			if tc.wantUnchanged != nil {
				entry := rewritten[0].(map[string]any)
				for k, want := range tc.wantUnchanged {
					if got := entry[k]; !equalValue(got, want) {
						t.Fatalf("rewritten[0].%s = %#v, want %#v", k, got, want)
					}
				}
			}
		})
	}
}

func TestResourceCustomizeDiffClearsOmittedCopySettings(t *testing.T) {
	t.Parallel()
	r := Resource()
	block := r.CoreConfigSchema()
	ty := block.ImpliedType()
	copySettingsType := ty.AttributeTypes()["copy_settings"]
	priorAttrs := map[string]string{
		"id":                                 "id",
		"project_id":                         "p",
		"cluster_name":                       "c",
		"copy_policy_items_enabled":          "false",
		"copy_settings.#":                    "1",
		"copy_settings.0.cloud_provider":     "AWS",
		"copy_settings.0.region_name":        "US_WEST_1",
		"copy_settings.0.zone_id":            "zone",
		"copy_settings.0.should_copy_oplogs": "true",
		"copy_settings.0.frequencies.#":      "1",
		"copy_settings.0.frequencies.0":      "DAILY",
	}
	tests := []struct {
		override cty.Value
		name     string
	}{
		{name: "null", override: cty.NullVal(copySettingsType)},
		{name: "unknown", override: cty.UnknownVal(copySettingsType)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := objectWithNullDefaults(ty, map[string]cty.Value{
				"project_id":                cty.StringVal("p"),
				"cluster_name":              cty.StringVal("c"),
				"copy_policy_items_enabled": cty.False,
				"copy_settings":             tc.override,
			})
			state := &terraform.InstanceState{
				ID:         "id",
				Attributes: priorAttrs,
				RawConfig:  raw,
			}
			diff, err := r.SimpleDiff(context.Background(), state, terraform.NewResourceConfigShimmed(raw, block), nil)
			if err != nil {
				t.Fatalf("SimpleDiff: %v", err)
			}
			if diff == nil {
				t.Fatal("expected a diff to clear omitted copy_settings, got nil")
				return
			}
			got := diff.Attributes["copy_settings.#"]
			if got == nil || got.New != "0" || got.NewComputed {
				t.Fatalf("copy_settings.# diff = %#v, want New=0", got)
			}
		})
	}
}

func TestCopySettingsRawConfigEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw  cty.Value
		name string
		want bool
	}{
		{
			name: "null",
			raw:  cty.NullVal(cty.List(cty.EmptyObject)),
			want: true,
		},
		{
			name: "empty",
			raw:  cty.ListValEmpty(cty.EmptyObject),
			want: true,
		},
		{
			name: "unknown",
			raw:  cty.UnknownVal(cty.List(cty.EmptyObject)),
			want: true,
		},
		{
			name: "non-empty",
			raw:  cty.ListVal([]cty.Value{cty.EmptyObjectVal}),
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := copySettingsRawConfigEmpty(tc.raw); got != tc.want {
				t.Fatalf("copySettingsRawConfigEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
}

func testFreqSet(items ...string) *schema.Set {
	raw := make([]any, len(items))
	for i, item := range items {
		raw[i] = item
	}
	return schema.NewSet(schema.HashString, raw)
}

func objectWithNullDefaults(ty cty.Type, overrides map[string]cty.Value) cty.Value {
	vals := make(map[string]cty.Value, len(ty.AttributeTypes()))
	for name, aty := range ty.AttributeTypes() {
		if v, ok := overrides[name]; ok {
			vals[name] = v
			continue
		}
		vals[name] = cty.NullVal(aty)
	}
	return cty.ObjectVal(vals)
}

func setLen(v any) int {
	s, ok := v.(*schema.Set)
	if !ok || s == nil {
		return 0
	}
	return s.Len()
}

func equalValue(got, want any) bool {
	switch w := want.(type) {
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return false
		}
		for i := range w {
			if !equalValue(g[i], w[i]) {
				return false
			}
		}
		return true
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok || len(g) != len(w) {
			return false
		}
		for k, wv := range w {
			if !equalValue(g[k], wv) {
				return false
			}
		}
		return true
	default:
		return got == want
	}
}
