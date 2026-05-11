package dynamicreshape_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicreshape"
	"github.com/stretchr/testify/assert"
)

func TestReshape_KeepsConfiguredShapeDropsExtras(t *testing.T) {
	configured := map[string]any{
		"name":  "x",
		"roles": []any{"A"},
	}
	response := map[string]any{
		"name":      "x",
		"roles":     []any{"A"},
		"createdAt": "2026-01-01",
		"clientId":  "abc",
	}
	got := dynamicreshape.Reshape(configured, response, dynamicreshape.Options{})
	assert.Equal(t, configured, got)
}

func TestReshape_ConfiguredKeyMissingFromResponse(t *testing.T) {
	// Keys configured but absent from the response are write-only request
	// fields. Preserve the configured value rather than null'ing it out.
	configured := map[string]any{"name": "x", "description": "y"}
	response := map[string]any{"name": "x"}
	got := dynamicreshape.Reshape(configured, response, dynamicreshape.Options{})
	assert.Equal(t, map[string]any{"name": "x", "description": "y"}, got)
}

func TestReshape_ConfiguredKeyPresentButNull(t *testing.T) {
	// When the response explicitly returns null for a configured key, treat
	// that as real drift to null (not "absent from response").
	configured := map[string]any{"name": "x", "description": "y"}
	response := map[string]any{"name": "x", "description": nil}
	got := dynamicreshape.Reshape(configured, response, dynamicreshape.Options{})
	assert.Equal(t, map[string]any{"name": "x", "description": nil}, got)
}

func TestReshape_SensitivePathsExcluded(t *testing.T) {
	configured := map[string]any{"name": "x", "headers": map[string]any{"value": "secret"}}
	response := map[string]any{"name": "x", "headers": map[string]any{"value": "redacted"}}
	opts := dynamicreshape.Options{SensitivePaths: map[string]struct{}{"headers.value": {}}}
	got := dynamicreshape.Reshape(configured, response, opts)
	assert.Equal(t, map[string]any{"name": "x", "headers": map[string]any{}}, got)
}

func TestReshape_ListByID(t *testing.T) {
	configured := []any{
		map[string]any{"name": "h1", "value": "a"},
		map[string]any{"name": "h2", "value": "b"},
	}
	response := []any{
		map[string]any{"name": "h2", "value": "B"},
		map[string]any{"name": "h1", "value": "A"},
	}
	opts := dynamicreshape.Options{ListIDKeys: map[string]string{"": "name"}}
	got := dynamicreshape.Reshape(configured, response, opts)
	assert.Equal(t, []any{
		map[string]any{"name": "h1", "value": "A"},
		map[string]any{"name": "h2", "value": "B"},
	}, got)
}

func TestCollectSensitivePaths(t *testing.T) {
	got := dynamicreshape.CollectSensitivePaths(map[string]any{
		"a": "x",
		"b": map[string]any{"c": "y"},
	})
	assert.Contains(t, got, "a")
	assert.Contains(t, got, "b")
	assert.Contains(t, got, "b.c")
}
