package backupcompliancepolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// In-package unit tests for the SDKv2 state upgrade transforms (see resource_hooks.go UpgradeState).

func TestUpgradeOnDemandPolicyItem(t *testing.T) {
	rawMap := map[string]any{"on_demand_policy_item": []any{map[string]any{"frequency_interval": float64(0), "retention_unit": "days"}}}
	upgradeOnDemandPolicyItem(rawMap)
	assert.Equal(t, map[string]any{"frequency_interval": float64(0), "frequency_type": "ondemand", "retention_unit": "days"}, rawMap["on_demand_policy_item"])

	rawMapEmpty := map[string]any{"on_demand_policy_item": []any{}}
	upgradeOnDemandPolicyItem(rawMapEmpty)
	assert.NotContains(t, rawMapEmpty, "on_demand_policy_item")

	rawMapAbsent := map[string]any{}
	upgradeOnDemandPolicyItem(rawMapAbsent)
	assert.NotContains(t, rawMapAbsent, "on_demand_policy_item")
}

func TestUpgradeScheduledPolicyItems(t *testing.T) {
	rawMap := map[string]any{
		"policy_item_hourly":  []any{map[string]any{"frequency_interval": float64(6), "retention_unit": "days"}},
		"policy_item_monthly": []any{map[string]any{"frequency_type": "monthly", "frequency_interval": float64(0)}, map[string]any{"frequency_interval": float64(1)}},
		"policy_item_daily":   []any{},
	}
	upgradeScheduledPolicyItems(rawMap)
	scheduled := rawMap["scheduled_policy_items"].([]any)
	require.Len(t, scheduled, 3)
	assert.Equal(t, "hourly", scheduled[0].(map[string]any)["frequency_type"])
	assert.Equal(t, "monthly", scheduled[1].(map[string]any)["frequency_type"])
	assert.Equal(t, "monthly", scheduled[2].(map[string]any)["frequency_type"])
	for _, removed := range []string{"policy_item_hourly", "policy_item_daily", "policy_item_monthly"} {
		assert.NotContains(t, rawMap, removed)
	}
}
