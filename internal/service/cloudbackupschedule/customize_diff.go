package cloudbackupschedule

import (
	"context"
	"errors"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/sdkv2config"
)

const (
	errCopySettingsModes      = "copy_settings entry must set only one of frequencies, copy_policy_items, or last_number_of_snapshots"
	errCopyPolicyRequiresFlag = "copy_policy_items and last_number_of_snapshots require copy_policy_items_enabled to be true"
	errFrequenciesWithFlag    = "frequencies cannot be set when copy_policy_items_enabled is true"
	errDeleteCopyRequiresFlag = "delete_copy_snapshots requires copy_policy_items_enabled to be true"
	errUpdateCopyRequiresFlag = "update_copy_snapshots requires copy_policy_items_enabled to be true"
)

// resourceCustomizeDiff rewrites copy_settings during plan because the list and nested
// frequencies are Optional+Computed. Without it, omitted HCL would keep prior copies in the
// plan (delete-on-omit breaks), and switching to copy_policy_items or last_number_of_snapshots
// would keep leftover frequencies from state (SDKv2 normalizeNullValues restores the planned set
// after apply). SetNew only works on the top-level copy_settings list, so we rewrite the whole
// block: force [] when raw config omits the list, and clear frequency sets when the flag is on.
func resourceCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if err := validateCopySettingsModes(d); err != nil {
		return err
	}
	// Optional+Computed would keep last state when HCL omits copy_settings. Force an empty list so delete-on-omit stays.
	if sdkv2config.CollectionEmpty(d, "copy_settings") {
		return d.SetNew("copy_settings", []any{})
	}
	return clearFrequenciesWhenCopyPolicyEnabled(d)
}

func clearFrequenciesWhenCopyPolicyEnabled(d *schema.ResourceDiff) error {
	enabled, _ := d.Get("copy_policy_items_enabled").(bool)
	if !enabled {
		return nil
	}
	copySettings, _ := d.Get("copy_settings").([]any)
	rewritten, changed := copySettingsWithEmptyFrequencies(copySettings)
	if !changed {
		return nil
	}
	return d.SetNew("copy_settings", rewritten)
}

func copySettingsWithEmptyFrequencies(copySettings []any) ([]any, bool) {
	changed := false
	rewritten := make([]any, len(copySettings))
	for i, raw := range copySettings {
		entry, ok := raw.(map[string]any)
		if !ok {
			rewritten[i] = raw
			continue
		}
		copied := maps.Clone(entry)
		if collectionLen(copied["frequencies"]) > 0 {
			copied["frequencies"] = schema.NewSet(schema.HashString, []any{})
			changed = true
		}
		rewritten[i] = copied
	}
	return rewritten, changed
}

func validateCopySettingsModes(d *schema.ResourceDiff) error {
	enabled, _ := d.Get("copy_policy_items_enabled").(bool)
	updateCopy, _ := d.Get("update_copy_snapshots").(bool)
	deleteCopy, _ := d.Get("delete_copy_snapshots").(bool)
	if updateCopy && !enabled {
		return errors.New(errUpdateCopyRequiresFlag)
	}
	if deleteCopy && !enabled {
		return errors.New(errDeleteCopyRequiresFlag)
	}

	copySettings, _ := d.Get("copy_settings").([]any)
	for i, raw := range copySettings {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		hasFreq := collectionLen(entry["frequencies"]) > 0
		hasItems := collectionLen(entry["copy_policy_items"]) > 0
		lastN, _ := entry["last_number_of_snapshots"].(int)
		hasLastN := lastN > 0

		modeCount := 0
		// Raw config only: d.Get frequencies can still be leftover Optional+Computed state.
		if sdkv2config.NestedCollectionLen(d, "copy_settings", i, "frequencies") > 0 {
			modeCount++
		}
		if hasItems {
			modeCount++
		}
		if hasLastN {
			modeCount++
		}
		if modeCount > 1 {
			return errors.New(errCopySettingsModes)
		}
		if (hasItems || hasLastN) && !enabled {
			return errors.New(errCopyPolicyRequiresFlag)
		}
		// Leftover Optional+Computed frequencies are ignored when another mode is set so a frequencies config can switch in one apply.
		if enabled && hasFreq && !hasItems && !hasLastN {
			return errors.New(errFrequenciesWithFlag)
		}
	}
	return nil
}

func collectionLen(v any) int {
	switch x := v.(type) {
	case *schema.Set:
		if x == nil {
			return 0
		}
		return x.Len()
	case []any:
		return len(x)
	default:
		return 0
	}
}
