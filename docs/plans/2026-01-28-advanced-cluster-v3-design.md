# Advanced Cluster v3.0.0 Design

## Overview

This document describes the breaking changes for the `mongodbatlas_advanced_cluster` resource in Terraform Provider v3.0.0. The primary change is simplifying how `replication_specs` are handled by making them Optional only (not Computed) and introducing a new `effective_replication_specs` attribute for data sources.

## Goals

1. Simplify the `replication_specs` schema by removing Optional+Computed complexity
2. Remove the `use_effective_fields` flag - the new behavior is always on
3. Provide clear visibility into configured vs actual values via `effective_replication_specs` on data sources
4. Remove plan modifier complexity that was required for Optional+Computed handling

## Schema Changes

### Resource Schema (`internal/service/advancedcluster/schema.go`)

**Removals:**
- `use_effective_fields` attribute - deleted entirely
- `UseEffectiveFieldsValidator` - no longer needed

**Modifications:**
- `replication_specs`: Change from `Optional + Computed` to `Optional` only
- All nested children (`region_configs`, `electable_specs`, `read_only_specs`, `analytics_specs`, `auto_scaling`, etc.): Change from `Optional + Computed` to `Optional` only
- Schema `Version`: Change from `2` to `3`

**TFModel changes:**
- Remove `UseEffectiveFields` field

### Data Source Schema (singular and plural)

**Removals:**
- `use_effective_fields` attribute
- `effective_electable_specs`, `effective_read_only_specs`, `effective_analytics_specs` from `region_configs`

**Additions:**
- `effective_replication_specs` (Computed) - identical structure to `replication_specs`, mirrors the full tree without the old `effective_*_specs` nested attributes

## API Integration

### Resource Implementation (`internal/service/advancedcluster/resource.go`)

**Changes to API calls:**

1. **Create operation** - Add `.UseEffectiveFieldsReplicationSpecs(true)` to the create request
2. **Read operation** - Add `.UseEffectiveFieldsReplicationSpecs(true)` to the get request
3. **Update operation** - Add `.UseEffectiveFieldsReplicationSpecs(true)` to the update request
4. **Delete operation** - Add if applicable

**Removals:**
- All references to `UseEffectiveInstanceFields`
- `waitParams.UseEffectiveFields` field and usage
- `ClusterWaitParams.UseEffectiveFields` field

### Data Source Implementation

- Add `.UseEffectiveFieldsReplicationSpecs(true)` to all cluster fetch operations
- Remove `UseEffectiveInstanceFields` usage
- Update `GetClusterDetails()` signature to remove the `useEffectiveFields` parameter

### Model Conversion (`internal/service/advancedcluster/model_ClusterDescription20240805.go`)

- Add conversion logic for `effective_replication_specs` from API response
- Remove conversion logic for `effective_electable_specs`, `effective_read_only_specs`, `effective_analytics_specs`

## Plan Modifier

**Action:** Delete the entire file `internal/service/advancedcluster/plan_modifier.go`

Since all `replication_specs` children are now Optional only (not Computed), Terraform handles plan/state comparison naturally. No plan modifier logic is needed.

**Functions to remove:**
- `handleModifyPlan()`
- `isReadOnlySpecsDeleted()`
- `isAnalyticsSpecsDeleted()`
- `adjustRegionConfigsChildren()`
- All helper functions

**Resource changes:**
- Remove `ModifyPlan` method from resource, or make it a no-op if framework requires it

## State Upgrade (`internal/service/advancedcluster/move_upgrade_state.go`)

**Additions:**
- Add `stateUpgraderFromV2` function for v2 → v3 migration
- Simple upgrade: copy existing state fields to new schema structure

**Modifications:**
- Update `UpgradeState()` to include v2 → v3 upgrade path
- Keep existing v1 → v2 upgrade for users upgrading from very old versions (v1 → v2 → v3 chain)

## Tests

### Test File Changes (`internal/service/advancedcluster/effective_fields_test.go`)

**Rename/Repurpose:** Convert to test `effective_replication_specs` behavior

**Remove tests for:**
- `use_effective_fields` flag toggle
- Tenant/Flex restriction validation
- Spec removal validation during toggle

**Convert/Keep tests for:**
- Basic cluster creation verifying `effective_replication_specs` is populated correctly
- Auto-scaling scenarios where configured values differ from effective values
- Data source returns both `replication_specs` (configured) and `effective_replication_specs` (actual)

**New test scenarios:**
- Verify `effective_replication_specs` structure matches `replication_specs`
- Verify auto-scaling changes reflected in `effective_replication_specs` but not in `replication_specs`

## Documentation & Examples

### Documentation Updates

**`docs/resources/advanced_cluster.md`:**
- Remove all references to `use_effective_fields`
- Remove warnings about toggling the flag
- Update auto-scaling example
- Add note that all `replication_specs` children must be explicitly defined

**`docs/data-sources/advanced_cluster.md` and `docs/data-sources/advanced_clusters.md`:**
- Remove `use_effective_fields` attribute documentation
- Remove `effective_electable_specs`, `effective_read_only_specs`, `effective_analytics_specs` documentation
- Add `effective_replication_specs` attribute documentation with full structure

### New Migration Guide

**Create `docs/guides/3.0.0-advanced-cluster-migration.md`:**

Migration steps:
1. Remove `use_effective_fields` from resource config
2. Update data source references from `effective_*_specs` to `effective_replication_specs`

### Breaking Behavior Change Warning

**Important: Previously removed specs will now be deleted**

In Terraform Provider 2.x, if you **removed** `read_only_specs` or `analytics_specs` from your configuration, the plan modifier preserved them - they continued running even though they were no longer in your config.

In 3.0.0, since all `replication_specs` children are Optional only (not Computed), Terraform will now detect this mismatch and **plan to delete** those nodes that exist in your cluster but are missing from your configuration.

**Before upgrading:** Review your actual cluster state in Atlas. If your cluster has read-only or analytics nodes that you previously removed from your Terraform config:
1. Add them back to your configuration before upgrading to preserve them, or
2. Accept that Terraform will remove them after upgrading

### Examples

**Delete:** `examples/mongodbatlas_advanced_cluster/effective_fields/` (entire directory)

**Create:** `examples/mongodbatlas_advanced_cluster/effective_replication_specs/`
- Show resource without `use_effective_fields`
- Show data source accessing both `replication_specs` and `effective_replication_specs`
- Demonstrate auto-scaling scenario where values differ
- Include README.md with warning about the behavior change

## Files Summary

### Files to Modify

| File | Changes |
|------|---------|
| `internal/service/advancedcluster/schema.go` | Remove `use_effective_fields`, change Optional+Computed to Optional, Version 2→3, add `effective_replication_specs` to DS schemas |
| `internal/service/advancedcluster/resource.go` | Add `UseEffectiveFieldsReplicationSpecs(true)`, remove `UseEffectiveInstanceFields` usage |
| `internal/service/advancedcluster/data_source.go` | Add `UseEffectiveFieldsReplicationSpecs(true)`, remove old effective fields |
| `internal/service/advancedcluster/plural_data_source.go` | Same as above |
| `internal/service/advancedcluster/model_ClusterDescription20240805.go` | Add `effective_replication_specs` conversion, remove `effective_*_specs` |
| `internal/service/advancedcluster/move_upgrade_state.go` | Add v2→v3 state upgrade |
| `internal/service/advancedcluster/effective_fields_test.go` | Convert to test new behavior |
| `docs/resources/advanced_cluster.md` | Remove `use_effective_fields` docs, update examples |
| `docs/data-sources/advanced_cluster.md` | Update for `effective_replication_specs` |
| `docs/data-sources/advanced_clusters.md` | Same as above |

### Files to Delete

| File | Reason |
|------|--------|
| `internal/service/advancedcluster/plan_modifier.go` | No longer needed |
| `examples/mongodbatlas_advanced_cluster/effective_fields/` | Old approach, replaced |

### Files to Create

| File | Purpose |
|------|---------|
| `docs/guides/3.0.0-advanced-cluster-migration.md` | Migration guide with breaking change warning |
| `examples/mongodbatlas_advanced_cluster/effective_replication_specs/` | New example with README warning |

## Implementation Notes

- The new `UseEffectiveFieldsReplicationSpecs(true)` API flag is separate from the old `UseEffectiveInstanceFields`
- The flag is always set to `true` internally and not exposed in the Terraform schema
- `effective_replication_specs` has identical structure to `replication_specs` (no nested `effective_*_specs` attributes)
- Build verification: `make build && golangci-lint run`
