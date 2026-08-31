package cloudbackupschedule_test

import (
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
)

type testGetter map[string]any

func (g testGetter) Get(key string) any {
	return g[key]
}

func TestValidateCopySettingsModes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		getter  testGetter
		wantErr string
	}{
		{
			name: "frequencies only without flag",
			getter: testGetter{
				"copy_settings": []any{
					map[string]any{
						"frequencies": testFreqSet("DAILY"),
					},
				},
			},
		},
		{
			name: "copy policy items with flag",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"copy_policy_items": []any{
							map[string]any{"frequency_type": "daily"},
						},
					},
				},
			},
		},
		{
			name: "last N with flag",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"last_number_of_snapshots": 5,
					},
				},
			},
		},
		{
			name: "copy policy items ignore leftover frequencies",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"frequencies": testFreqSet("DAILY"),
						"copy_policy_items": []any{
							map[string]any{"frequency_type": "daily"},
						},
					},
				},
			},
		},
		{
			name: "last N ignores leftover frequencies",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"frequencies":              testFreqSet("DAILY"),
						"last_number_of_snapshots": 5,
					},
				},
			},
		},
		{
			name: "copy policy items and last N",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"copy_policy_items": []any{
							map[string]any{"frequency_type": "daily"},
						},
						"last_number_of_snapshots": 5,
					},
				},
			},
			wantErr: "only one of frequencies, copy_policy_items, or last_number_of_snapshots",
		},
		{
			name: "copy policy items with flag omitted",
			getter: testGetter{
				"copy_settings": []any{
					map[string]any{
						"copy_policy_items": []any{
							map[string]any{"frequency_type": "daily"},
						},
					},
				},
			},
			wantErr: "require copy_policy_items_enabled to be true",
		},
		{
			name: "last N with flag omitted",
			getter: testGetter{
				"copy_settings": []any{
					map[string]any{
						"last_number_of_snapshots": 5,
					},
				},
			},
			wantErr: "require copy_policy_items_enabled to be true",
		},
		{
			name: "last N with flag false",
			getter: testGetter{
				"copy_policy_items_enabled": false,
				"copy_settings": []any{
					map[string]any{
						"last_number_of_snapshots": 5,
					},
				},
			},
			wantErr: "require copy_policy_items_enabled to be true",
		},
		{
			name: "frequencies with flag true",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"copy_settings": []any{
					map[string]any{
						"frequencies": testFreqSet("DAILY"),
					},
				},
			},
			wantErr: "frequencies cannot be set when copy_policy_items_enabled is true",
		},
		{
			name: "delete copy snapshots without flag",
			getter: testGetter{
				"delete_copy_snapshots": true,
			},
			wantErr: "delete_copy_snapshots requires copy_policy_items_enabled to be true",
		},
		{
			name: "delete copy snapshots with flag",
			getter: testGetter{
				"copy_policy_items_enabled": true,
				"delete_copy_snapshots":     true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := cloudbackupschedule.ValidateCopySettingsModes(tc.getter)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}
