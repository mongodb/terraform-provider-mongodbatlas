package cloudbackupschedule

import (
	"context"
	"errors"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/sdkv2config"
)

const (
	errCopySettingsModes        = "copy_settings entry must set only one of frequencies, copy_policy_items, or last_number_of_snapshots"
	errCopyPolicyRequiresFlag   = "copy_policy_items and last_number_of_snapshots require copy_policy_items_enabled to be true"
	errFrequenciesWithFlag      = "frequencies cannot be set when copy_policy_items_enabled is true"
	errDeleteCopyRequiresFlag   = "delete_copy_snapshots requires copy_policy_items_enabled to be true"
	errUpdateCopyRequiresFlag   = "update_copy_snapshots requires copy_policy_items_enabled to be true"
	errAutoCopyWithCopySettings = "copy_settings cannot be set when auto_copy_settings_enabled is true"
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
	// Skip force-empty when auto_copy_settings_enabled is true so Atlas-managed copies are not overwritten or deleted.
	if sdkv2config.CollectionEmpty(d, "copy_settings") {
		v := sdkv2config.Bool(d, "auto_copy_settings_enabled")
		if !shouldForceEmptyCopySettings(v.Set() && v.ValueBool()) {
			// Atlas owns the copy region list. Keep the last GET so an unchanged config still plans
			// empty, but mark computed on the apply that flips the flag on: the PATCH omits
			// copySettings and the read then writes whatever Atlas returns, which a known planned
			// value from the removed manual entries cannot predict.
			if v.HasChange() {
				return d.SetNewComputed("copy_settings")
			}
			return nil
		}
		return d.SetNew("copy_settings", []any{})
	}
	return clearFrequenciesWhenCopyPolicyEnabled(d)
}

// shouldForceEmptyCopySettings reports whether the plan or PATCH must send an empty copy_settings
// list. Atlas owns the copy region list when auto_copy_settings_enabled is true, so the force-empty
// delete-on-omit only applies while the flag is off.
func shouldForceEmptyCopySettings(autoCopyEnabled bool) bool {
	return !autoCopyEnabled
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
	// Raw config only: a copy_settings block present in prior state but omitted in HCL must not trip this.
	if sdkv2config.Bool(d, "auto_copy_settings_enabled").ValueBool() && !sdkv2config.CollectionEmpty(d, "copy_settings") {
		return errors.New(errAutoCopyWithCopySettings)
	}

	enabled, _ := d.Get("copy_policy_items_enabled").(bool)
	updateCopy, _ := d.Get("update_copy_snapshots").(bool)
	deleteCopy, _ := d.Get("delete_copy_snapshots").(bool)
	if updateCopy && !enabled {
		return errors.New(errUpdateCopyRequiresFlag)
	}
	if deleteCopy && !enabled {
		return errors.New(errDeleteCopyRequiresFlag)
	}

	// With the flag on, the raw config cannot carry copy_settings entries (the check above rejects
	// them), so anything d.Get returns here is Atlas-written state. Skip entry validation so those
	// Atlas-written modes do not fail a plan that omits the list, for example the migration from a
	// manual copy_settings config to auto copy settings.
	autoCopy := sdkv2config.Bool(d, "auto_copy_settings_enabled")
	if autoCopy.Set() && autoCopy.ValueBool() && sdkv2config.CollectionEmpty(d, "copy_settings") {
		return nil
	}

	// Validate the raw config, not d.Get: copy_settings is Optional+Computed, so d.Get still
	// carries prior state when HCL omits the list. Reading all three modes from raw config
	// keeps them consistent and lets an in-place removal of every copy target plan cleanly
	// instead of tripping on the stale state entry.
	copySettings, _ := d.Get("copy_settings").([]any)
	for i, raw := range copySettings {
		if _, ok := raw.(map[string]any); !ok {
			continue
		}
		hasFreq := sdkv2config.NestedCollectionLen(d, "copy_settings", i, "frequencies") > 0
		hasItems := sdkv2config.NestedCollectionLen(d, "copy_settings", i, "copy_policy_items") > 0
		hasLastN := sdkv2config.NestedInt(d, "copy_settings", i, "last_number_of_snapshots") > 0

		modeCount := 0
		if hasFreq {
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
		// Frequencies are ignored when another mode is set so a frequencies config can switch in one apply.
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
