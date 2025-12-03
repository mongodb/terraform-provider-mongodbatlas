# Migrating Existing Modules to use Effective Fields

This example will demonstrate how to migrate an existing Terraform module to use the `use_effective_fields` attribute, enabling it to handle both auto-scaling and non-auto-scaling configurations without requiring `lifecycle.ignore_changes` blocks.

**Status:** This example is currently under development and will be available in a future release.

## Coming Soon

This directory will contain:
- An example of an existing module that uses `lifecycle.ignore_changes`
- Step-by-step migration guide to adopt `use_effective_fields`
- Before and after comparisons
- Best practices for migrating production modules

## In the Meantime

If you need to migrate an existing module now, please refer to the "Migration from Legacy Implementations" section in the [new-module example](../new-module/README.md#migration-from-legacy-implementations) for general guidance.
