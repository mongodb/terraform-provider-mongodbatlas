# MongoDB Atlas Advanced Cluster - Effective Fields Examples

This directory contains examples demonstrating how to use the `use_effective_fields` attribute in Terraform modules for MongoDB Atlas clusters. The `use_effective_fields` feature enables modules to handle both auto-scaling and non-auto-scaling configurations without requiring `lifecycle.ignore_changes` blocks.

## Directory Structure

This directory contains two subdirectories, each addressing a different use case:

### 1. [new_module](./new_module/)
**Creating a New Module from Scratch**

This example demonstrates how to build a reusable Terraform module from the ground up that uses `use_effective_fields`. It includes:
- Complete module implementation with `use_effective_fields = true`
- Examples of using the module with auto-scaling enabled
- Examples of using the module without auto-scaling
- Module interface documentation and best practices

**Use this example when:**
- You're creating a new Terraform module for MongoDB Atlas clusters
- You want to understand the complete architecture of an effective fields module
- You need a reference implementation for module design patterns

### 2. [existing_module](./existing_module/)
**Migrating an Existing Module**

This example demonstrates how to migrate an existing Terraform module that uses `lifecycle.ignore_changes` to adopt the `use_effective_fields` approach.

**Use this example when:**
- You have an existing module that uses `lifecycle.ignore_changes`
- You want to migrate to the `use_effective_fields` approach
- You need guidance on updating production modules safely

## What is use_effective_fields?

The `use_effective_fields` attribute addresses a key challenge in creating reusable Terraform modules for MongoDB Atlas clusters with auto-scaling:

**The Problem:** When Atlas auto-scales a cluster, it changes instance sizes and other specifications. Without special handling, these changes cause Terraform plan drift, forcing module authors to use `lifecycle.ignore_changes` blocks. However, lifecycle blocks cannot be conditional, requiring separate module implementations for auto-scaling and non-auto-scaling scenarios.

**The Solution:** Setting `use_effective_fields = true` tells Terraform to read the user-configured values for planning and ignore Atlas-managed changes from auto-scaling. This enables a single module to handle both scenarios without lifecycle blocks.

## Benefits

1. **Single module implementation**: One codebase supports both auto-scaling and non-auto-scaling use cases
2. **No lifecycle blocks needed**: Cleaner, more maintainable code
3. **Operational visibility**: Modules can expose both configured and actual (effective) values
4. **Forward compatibility**: Aligns with provider v3.x where this becomes default behavior

## Getting Started

1. If you're creating a new module, start with the [new_module](./new_module/) example
2. If you're migrating an existing module, check the [existing_module](./existing_module/) directory

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
