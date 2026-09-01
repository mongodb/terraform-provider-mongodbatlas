package cloudbackupschedule_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
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
			name: "already empty frequencies",
			copySettings: []any{
				map[string]any{
					"cloud_provider": "AWS",
					"frequencies":    testFreqSet(),
				},
			},
			wantChanged:  false,
			wantFreqLens: []int{0},
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

			rewritten, changed := cloudbackupschedule.CopySettingsWithEmptyFrequencies(tc.copySettings)
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
			if tc.name == "two entries only index 1 has leftover" {
				first := tc.copySettings[0].(map[string]any)
				rewrittenFirst := rewritten[0].(map[string]any)
				if setLen(first["frequencies"]) != 0 {
					t.Fatalf("input[0].frequencies was mutated, len = %d", setLen(first["frequencies"]))
				}
				if rewrittenFirst["region_name"] != "US_EAST_1" {
					t.Fatalf("rewritten[0].region_name = %v, want US_EAST_1", rewrittenFirst["region_name"])
				}
				if rewritten[1].(map[string]any)["region_name"] != "US_WEST_1" {
					t.Fatalf("rewritten[1].region_name = %v, want US_WEST_1", rewritten[1].(map[string]any)["region_name"])
				}
			}
		})
	}
}

type testDiffSetter struct {
	vals   map[string]any
	setNew map[string]any
}

func (s *testDiffSetter) Get(key string) any {
	return s.vals[key]
}

func (s *testDiffSetter) SetNew(key string, value any) error {
	if s.setNew == nil {
		s.setNew = map[string]any{}
	}
	s.setNew[key] = value
	return nil
}

func TestClearFrequenciesWhenCopyPolicyEnabled(t *testing.T) {
	t.Parallel()
	leftover := []any{
		map[string]any{
			"cloud_provider": "AWS",
			"frequencies":    testFreqSet("DAILY"),
			"copy_policy_items": []any{
				map[string]any{"frequency_type": "daily"},
			},
		},
	}
	emptyFreq := []any{
		map[string]any{
			"cloud_provider": "AWS",
			"frequencies":    testFreqSet(),
			"copy_policy_items": []any{
				map[string]any{"frequency_type": "daily"},
			},
		},
	}

	tests := []struct {
		vals        map[string]any
		name        string
		wantFreqLen int
		wantSetNew  bool
	}{
		{
			name: "flag false leftover frequencies",
			vals: map[string]any{
				"copy_policy_items_enabled": false,
				"copy_settings":             leftover,
			},
		},
		{
			name: "flag true leftover frequencies",
			vals: map[string]any{
				"copy_policy_items_enabled": true,
				"copy_settings":             leftover,
			},
			wantSetNew:  true,
			wantFreqLen: 0,
		},
		{
			name: "flag true frequencies already empty",
			vals: map[string]any{
				"copy_policy_items_enabled": true,
				"copy_settings":             emptyFreq,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := &testDiffSetter{vals: tc.vals}
			if err := cloudbackupschedule.ClearFrequenciesWhenCopyPolicyEnabled(d); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, ok := d.setNew["copy_settings"]
			if tc.wantSetNew != ok {
				t.Fatalf("SetNew called = %v, want %v (setNew=%#v)", ok, tc.wantSetNew, d.setNew)
			}
			if !tc.wantSetNew {
				return
			}
			list, ok := got.([]any)
			if !ok || len(list) != 1 {
				t.Fatalf("SetNew copy_settings = %#v, want 1-entry list", got)
			}
			entry := list[0].(map[string]any)
			if setLen(entry["frequencies"]) != tc.wantFreqLen {
				t.Fatalf("SetNew frequencies len = %d, want %d", setLen(entry["frequencies"]), tc.wantFreqLen)
			}
		})
	}
}

func TestResourceCustomizeDiffClearsOmittedCopySettings(t *testing.T) {
	t.Parallel()
	r := cloudbackupschedule.Resource()
	block := r.CoreConfigSchema()
	ty := block.ImpliedType()
	copySettingsType := ty.AttributeTypes()["copy_settings"]
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
				ID: "id",
				Attributes: map[string]string{
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
				},
				RawConfig: raw,
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

func copySettingsListVal(copySettingsType cty.Type) cty.Value {
	entryTy := copySettingsType.ElementType()
	entry := objectWithNullDefaults(entryTy, map[string]cty.Value{
		"cloud_provider":     cty.StringVal("AWS"),
		"region_name":        cty.StringVal("US_WEST_1"),
		"zone_id":            cty.StringVal("zone"),
		"should_copy_oplogs": cty.True,
		"frequencies":        cty.SetVal([]cty.Value{cty.StringVal("DAILY")}),
	})
	return cty.ListVal([]cty.Value{entry})
}

func resourceObject(ty cty.Type, copySettings cty.Value) cty.Value {
	return objectWithNullDefaults(ty, map[string]cty.Value{
		"project_id":                cty.StringVal("p"),
		"cluster_name":              cty.StringVal("c"),
		"copy_policy_items_enabled": cty.False,
		"copy_settings":             copySettings,
	})
}

func priorCopySettingsFlatmap() map[string]string {
	return map[string]string{
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
}

func describeCtyVal(v cty.Value) string {
	if v == cty.NilVal {
		return "nil"
	}
	length := "n/a"
	if !v.IsNull() && (v.IsKnown() || v.Type().IsTupleType()) {
		length = fmt.Sprintf("%d", v.LengthInt())
	}
	return fmt.Sprintf("null=%v known=%v whollyKnown=%v type=%s len=%s", v.IsNull(), v.IsKnown(), v.IsWhollyKnown(), v.Type().FriendlyName(), length)
}

func copySettingsFromResource(rawConfig cty.Value) cty.Value {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return cty.NullVal(cty.DynamicPseudoType)
	}
	return rawConfig.GetAttr("copy_settings")
}

func copySettingsGetLen(v any) int {
	list, ok := v.([]any)
	if !ok {
		return 0
	}
	return len(list)
}

func logCopySettingsCustomizeDiff(t *testing.T, d *schema.ResourceDiff) {
	t.Helper()
	rawConfig := d.GetRawConfig()
	rawPlan := d.GetRawPlan()
	copySettingsRaw := copySettingsFromResource(rawConfig)
	copySettingsPlan := copySettingsFromResource(rawPlan)
	t.Logf("CustomizeDiff id=%q", d.Id())
	t.Logf("  RawConfig.copy_settings: %s empty=%v", describeCtyVal(copySettingsRaw), cloudbackupschedule.CopySettingsRawConfigEmpty(copySettingsRaw))
	t.Logf("  RawPlan.copy_settings: %s", describeCtyVal(copySettingsPlan))
	t.Logf("  Get(copy_settings) len=%d HasChange=%v NewValueKnown=%v", copySettingsGetLen(d.Get("copy_settings")), d.HasChange("copy_settings"), d.NewValueKnown("copy_settings"))
}

func runSimpleDiffCopySettingsDump(t *testing.T, r *schema.Resource, state *terraform.InstanceState, shimVal cty.Value, priorFlatmap map[string]string) {
	t.Helper()
	block := r.CoreConfigSchema()
	if state.RawState != cty.NilVal && !state.RawState.IsNull() {
		t.Logf("prior RawState.copy_settings: %s", describeCtyVal(copySettingsFromResource(state.RawState)))
	} else if priorFlatmap["copy_settings.#"] != "" {
		t.Logf("prior flatmap copy_settings.#=%s", priorFlatmap["copy_settings.#"])
	}
	cfg := terraform.NewResourceConfigShimmed(shimVal, block)
	diff, err := r.SimpleDiff(context.Background(), state, cfg, nil)
	if err != nil {
		t.Fatalf("SimpleDiff: %v", err)
	}
	if diff == nil {
		t.Log("SimpleDiff returned nil; PlanResourceChange would return req.PriorState")
		if priorFlatmap["copy_settings.#"] != "" {
			t.Logf("prior flatmap copy_settings.#=%s", priorFlatmap["copy_settings.#"])
		}
		return
	}
	got := diff.Attributes["copy_settings.#"]
	t.Logf("SimpleDiff copy_settings.# diff: %#v", got)
	attrs := make(map[string]string, len(priorFlatmap))
	for k, v := range priorFlatmap {
		attrs[k] = v
	}
	applied, applyErr := diff.Apply(attrs, block)
	if applyErr != nil {
		t.Fatalf("diff.Apply: %v", applyErr)
	}
	t.Logf("after diff.Apply copy_settings.#=%s", applied["copy_settings.#"])
}

func TestDebugCopySettingsOmitPlan(t *testing.T) {
	r := cloudbackupschedule.Resource()
	block := r.CoreConfigSchema()
	ty := block.ImpliedType()
	copySettingsType := ty.AttributeTypes()["copy_settings"]
	copiesPresent := copySettingsListVal(copySettingsType)
	priorFlatmap := priorCopySettingsFlatmap()
	priorStateCty := resourceObject(ty, copiesPresent)
	proposedWithCopies := priorStateCty

	origCustomizeDiff := r.CustomizeDiff
	r.CustomizeDiff = func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		beforeLen := copySettingsGetLen(d.Get("copy_settings"))
		logCopySettingsCustomizeDiff(t, d)
		err := origCustomizeDiff(ctx, d, meta)
		afterLen := copySettingsGetLen(d.Get("copy_settings"))
		t.Logf("  CustomizeDiff after SetNew: beforeLen=%d afterLen=%d err=%v", beforeLen, afterLen, err)
		return err
	}

	t.Run("create_block_present", func(t *testing.T) {
		config := resourceObject(ty, copiesPresent)
		state := &terraform.InstanceState{
			RawConfig: config,
			RawPlan:   config,
		}
		runSimpleDiffCopySettingsDump(t, r, state, config, nil)
	})

	omitShapes := []struct {
		val  cty.Value
		name string
	}{
		{name: "null", val: cty.NullVal(copySettingsType)},
		{name: "empty", val: cty.ListValEmpty(copySettingsType.ElementType())},
		{name: "unknown", val: cty.UnknownVal(copySettingsType)},
	}
	for _, omit := range omitShapes {
		omitConfig := resourceObject(ty, omit.val)
		t.Run("omit_same_raw_"+omit.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID:         "id",
				Attributes: priorFlatmap,
				RawConfig:  omitConfig,
			}
			runSimpleDiffCopySettingsDump(t, r, state, omitConfig, priorFlatmap)
		})
		t.Run("omit_plan_rpc_"+omit.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID:         "id",
				Attributes: priorFlatmap,
				RawConfig:  omitConfig,
				RawState:   priorStateCty,
				RawPlan:    proposedWithCopies,
			}
			runSimpleDiffCopySettingsDump(t, r, state, proposedWithCopies, priorFlatmap)
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
			if got := cloudbackupschedule.CopySettingsRawConfigEmpty(tc.raw); got != tc.want {
				t.Fatalf("CopySettingsRawConfigEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
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
