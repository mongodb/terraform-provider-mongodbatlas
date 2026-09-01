package cloudbackupschedule_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
)

type testCopySettingsUpdateData struct {
	get       map[string]any
	hasChange map[string]bool
	getChange map[string][2]any
	rawConfig cty.Value
}

func (d testCopySettingsUpdateData) Get(key string) any {
	return d.get[key]
}

func (d testCopySettingsUpdateData) HasChange(key string) bool {
	return d.hasChange[key]
}

func (d testCopySettingsUpdateData) GetChange(key string) (old, newVal any) {
	ch := d.getChange[key]
	return ch[0], ch[1]
}

func (d testCopySettingsUpdateData) GetRawConfig() cty.Value {
	return d.rawConfig
}

func TestCopySettingsForUpdate(t *testing.T) {
	t.Parallel()
	leftover := []any{
		map[string]any{
			"cloud_provider": "AWS",
			"region_name":    "US_WEST_1",
			"zone_id":        "zone",
			"frequencies":    testFreqSet("DAILY"),
		},
	}
	copySettingsListTy := cty.List(cty.Object(map[string]cty.Type{
		"cloud_provider": cty.String,
	}))
	omitRawConfig := cty.ObjectVal(map[string]cty.Value{
		"copy_settings": cty.NullVal(copySettingsListTy),
	})
	tests := []struct {
		data    testCopySettingsUpdateData
		name    string
		wantLen int
		wantOK  bool
	}{
		{
			name: "omit raw config empty despite stale get",
			data: testCopySettingsUpdateData{
				get:       map[string]any{"copy_settings": leftover, "copy_settings.#": 1},
				rawConfig: omitRawConfig,
			},
			wantOK:  true,
			wantLen: 0,
		},
		{
			name: "omit uses count change when raw config unavailable",
			data: testCopySettingsUpdateData{
				get:       map[string]any{"copy_settings": leftover, "copy_settings.#": 0},
				hasChange: map[string]bool{"copy_settings.#": true},
				getChange: map[string][2]any{"copy_settings": {leftover, []any{}}},
			},
			wantOK:  true,
			wantLen: 0,
		},
		{
			name: "unchanged non-empty",
			data: testCopySettingsUpdateData{
				get:       map[string]any{"copy_settings": leftover},
				hasChange: map[string]bool{"copy_settings": false},
				rawConfig: cty.NullVal(cty.Object(map[string]cty.Type{"copy_settings": copySettingsListTy})),
			},
			wantOK:  true,
			wantLen: 1,
		},
		{
			name: "changed rewrite",
			data: testCopySettingsUpdateData{
				get:       map[string]any{"copy_settings": leftover},
				hasChange: map[string]bool{"copy_settings": true},
				getChange: map[string][2]any{
					"copy_settings": {
						leftover,
						[]any{
							map[string]any{
								"cloud_provider": "AWS",
								"region_name":    "US_WEST_1",
								"zone_id":        "zone",
								"frequencies":    testFreqSet(),
							},
						},
					},
				},
			},
			wantOK:  true,
			wantLen: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, ok := cloudbackupschedule.CopySettingsForUpdate(tc.data)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if len(got) != tc.wantLen {
				t.Fatalf("len = %d, want %d (got=%#v)", len(got), tc.wantLen, got)
			}
		})
	}
}
