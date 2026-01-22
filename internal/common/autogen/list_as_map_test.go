package autogen_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
)

func TestModifyJSONFromListToMap(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		var input any
		output := autogen.ModifyJSONFromListToMap(input)
		assert.Nil(t, output)
	})

	t.Run("non list input returned as-is", func(t *testing.T) {
		input := map[string]any{"foo": "bar"}
		output := autogen.ModifyJSONFromListToMap(input)
		assert.Equal(t, input, output)
	})

	t.Run("empty list returns empty map", func(t *testing.T) {
		input := []any{}
		output := autogen.ModifyJSONFromListToMap(input)
		assert.Equal(t, map[string]any{}, output)
	})

	t.Run("list with key/value objects returns map", func(t *testing.T) {
		input := []any{
			map[string]any{"key": "key1", "value": "val1"},
			map[string]any{"key": "key2", "value": 2},
		}

		output := autogen.ModifyJSONFromListToMap(input)

		expected := map[string]any{
			"key1": "val1",
			"key2": 2,
		}
		assert.Equal(t, expected, output)
	})
}

func TestModifyJSONFromMapToList(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		var input any
		output := autogen.ModifyJSONFromMapToList(input)
		assert.Nil(t, output)
	})

	t.Run("non map input returned as-is", func(t *testing.T) {
		input := []any{"foo"}
		output := autogen.ModifyJSONFromMapToList(input)
		assert.Equal(t, input, output)
	})

	t.Run("empty map returns empty list", func(t *testing.T) {
		input := map[string]any{}
		output := autogen.ModifyJSONFromMapToList(input)
		assert.Equal(t, []any{}, output)
	})

	t.Run("map with values returns list", func(t *testing.T) {
		input := map[string]any{
			"b": 2,
			"a": "one",
		}

		output := autogen.ModifyJSONFromMapToList(input)

		expected := []any{
			map[string]any{"key": "a", "value": "one"},
			map[string]any{"key": "b", "value": 2},
		}
		assert.Equal(t, expected, output)
	})
}
