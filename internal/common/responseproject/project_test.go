package responseproject_test

import (
	"reflect"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/responseproject"
)

func TestProject(t *testing.T) {
	resp := map[string]any{
		"apiKeyId":  "k-123",
		"name":      "voyage-prod",
		"createdAt": "2026-05-12T00:00:00Z",
		"secret":    "sk-abc",
		"profile": map[string]any{
			"region": "us-east-1",
			"owner":  "alice",
		},
		"secrets": []any{
			map[string]any{"id": "s0", "value": "v0"},
			map[string]any{"id": "s1", "value": "v1"},
		},
	}

	tests := []struct {
		want  map[string]any
		name  string
		paths []string
	}{
		{
			name:  "empty paths returns empty map",
			paths: nil,
			want:  map[string]any{},
		},
		{
			name:  "top-level scalar",
			paths: []string{"apiKeyId"},
			want:  map[string]any{"apiKeyId": "k-123"},
		},
		{
			name:  "multiple top-level",
			paths: []string{"apiKeyId", "name"},
			want: map[string]any{
				"apiKeyId": "k-123",
				"name":     "voyage-prod",
			},
		},
		{
			name:  "nested object path",
			paths: []string{"profile.region"},
			want: map[string]any{
				"profile": map[string]any{"region": "us-east-1"},
			},
		},
		{
			name:  "two nested fields share parent",
			paths: []string{"profile.region", "profile.owner"},
			want: map[string]any{
				"profile": map[string]any{"region": "us-east-1", "owner": "alice"},
			},
		},
		{
			name:  "list element by index",
			paths: []string{"secrets.0.value"},
			want: map[string]any{
				"secrets": []any{map[string]any{"value": "v0"}},
			},
		},
		{
			name:  "non-contiguous list index pads with nil",
			paths: []string{"secrets.1.id"},
			want: map[string]any{
				"secrets": []any{nil, map[string]any{"id": "s1"}},
			},
		},
		{
			name:  "missing path silently skipped",
			paths: []string{"apiKeyId", "doesNotExist", "profile.missing"},
			want:  map[string]any{"apiKeyId": "k-123"},
		},
		{
			name:  "out-of-bounds list index silently skipped",
			paths: []string{"secrets.5.value"},
			want:  map[string]any{},
		},
		{
			name:  "empty string path is skipped",
			paths: []string{"", "apiKeyId"},
			want:  map[string]any{"apiKeyId": "k-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := responseproject.Project(resp, tt.paths)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project()\n got = %#v\nwant = %#v", got, tt.want)
			}
		})
	}
}

func TestProject_NilResponse(t *testing.T) {
	got := responseproject.Project(nil, []string{"foo"})
	if !reflect.DeepEqual(got, map[string]any{}) {
		t.Errorf("Project(nil, ...) = %#v, want empty map", got)
	}
}

func TestPathsOverlap(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want []string
	}{
		{name: "no overlap", a: []string{"x"}, b: []string{"y"}, want: nil},
		{name: "one overlap", a: []string{"x", "y"}, b: []string{"y", "z"}, want: []string{"y"}},
		{name: "multiple overlaps", a: []string{"x", "y", "z"}, b: []string{"y", "z"}, want: []string{"y", "z"}},
		{name: "empty a", a: nil, b: []string{"y"}, want: nil},
		{name: "empty b", a: []string{"x"}, b: nil, want: nil},
		{name: "nested path equality", a: []string{"profile.region"}, b: []string{"profile.region"}, want: []string{"profile.region"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := responseproject.PathsOverlap(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PathsOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}
