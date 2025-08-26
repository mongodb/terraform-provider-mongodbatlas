package config_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestUserAgentExtra_ToHeaderValue(t *testing.T) {
	testCases := map[string]struct {
		extra    config.UserAgentExtra
		old      string
		expected string
	}{
		"all fields": {
			extra: config.UserAgentExtra{
				Name:      "name1",
				Operation: "op1",
			},
			old:      "base/1.0",
			expected: "base/1.0 Name/name1 Operation/op1",
		},
		"some fields empty": {
			extra: config.UserAgentExtra{
				Name:      "name2",
				Operation: "",
			},
			old:      "",
			expected: "Name/name2",
		},
		"none": {
			extra:    config.UserAgentExtra{},
			old:      "",
			expected: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.extra.ToHeaderValue(t.Context(), tc.old)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestUserAgentExtra_Combine(t *testing.T) {
	testCases := map[string]struct {
		base     config.UserAgentExtra
		other    config.UserAgentExtra
		expected config.UserAgentExtra
	}{
		"other overwrites non-empty": {
			base:     config.UserAgentExtra{Name: "B", Operation: "C"},
			other:    config.UserAgentExtra{Name: "Y", Operation: "Z"},
			expected: config.UserAgentExtra{Name: "Y", Operation: "Z"},
		},
		"other empty": {
			base:     config.UserAgentExtra{Name: "B", Operation: "C"},
			other:    config.UserAgentExtra{},
			expected: config.UserAgentExtra{Name: "B", Operation: "C"},
		},
		"mixed": {
			base:     config.UserAgentExtra{Name: "B", Operation: "O"},
			other:    config.UserAgentExtra{Name: "Y"},
			expected: config.UserAgentExtra{Name: "Y", Operation: "O"},
		},
		"extras combine base set": {
			base:     config.UserAgentExtra{Extras: map[string]string{"A": "ok"}},
			other:    config.UserAgentExtra{},
			expected: config.UserAgentExtra{Extras: map[string]string{"A": "ok"}},
		},
		"extras combine other set": {
			base:     config.UserAgentExtra{},
			other:    config.UserAgentExtra{Extras: map[string]string{"A": "ok"}},
			expected: config.UserAgentExtra{Extras: map[string]string{"A": "ok"}},
		},
		"extras combine both set": {
			base:     config.UserAgentExtra{Extras: map[string]string{"A": "ok"}},
			other:    config.UserAgentExtra{Extras: map[string]string{"B": "yes"}},
			expected: config.UserAgentExtra{Extras: map[string]string{"A": "ok", "B": "yes"}},
		},
		"all attributes set": {
			other: config.UserAgentExtra{
				Extras:        map[string]string{"B": "yes"},
				ModuleName:    "module-name",
				ModuleVersion: "1.2.3",
				Name:          "some-name",
				Operation:     "my-operation",
			},
			expected: config.UserAgentExtra{
				Extras:        map[string]string{"B": "yes"},
				ModuleName:    "module-name",
				ModuleVersion: "1.2.3",
				Name:          "some-name",
				Operation:     "my-operation",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.base.Combine(tc.other)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestAddUserAgentExtra(t *testing.T) {
	base := config.UserAgentExtra{Name: "B"}
	other := config.UserAgentExtra{Name: "Y", Operation: "O"}
	ctx := config.AddUserAgentExtra(t.Context(), base)
	ctx2 := config.AddUserAgentExtra(ctx, other)
	ua := config.ReadUserAgentExtra(ctx2)
	// The combined should have Type from base, Name from other, ScriptLocation from other
	assert.Equal(t, "Y", ua.Name)
	assert.Equal(t, "O", ua.Operation)
	assert.Empty(t, ua.Operation)
}
