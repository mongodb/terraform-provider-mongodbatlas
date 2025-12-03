# MongoDB Atlas Advanced Cluster - Effective Fields Examples

This directory contains examples demonstrating how to use the `use_effective_fields` attribute in Terraform modules for MongoDB Atlas clusters. This feature enhances auto-scaling workflows by eliminating the need for `lifecycle.ignore_changes` blocks and providing visibility into Atlas-managed changes.

## Directory Structure

### 1. [new_module](./new_module/)
**Building a New Module with Effective Fields**

Demonstrates how to build a reusable Terraform module from scratch using `use_effective_fields`. Includes complete module implementation with examples for both auto-scaling and non-auto-scaling configurations.

**Use this example when:**
- Creating a new Terraform module for MongoDB Atlas clusters
- Learning the architecture and best practices for effective fields modules
- Need a reference implementation for module design patterns

### 2. [existing_module](./existing_module/)
**Migrating an Existing Module to Effective Fields**

Demonstrates how to migrate an existing module from `lifecycle.ignore_changes` to `use_effective_fields`, showing both v1 (with lifecycle.ignore_changes) and v2 (with use_effective_fields) approaches side-by-side.

**Use this example when:**
- Migrating an existing module that uses `lifecycle.ignore_changes`
- Understanding the upgrade path and benefits of effective fields
- Providing a seamless upgrade experience for module users

## Understanding use_effective_fields

### The Challenge

When auto-scaling is enabled, Atlas automatically adjusts instance sizes and disk capacity based on workload. Without `use_effective_fields`, these Atlas-managed changes create plan drift, requiring `lifecycle.ignore_changes` blocks to prevent Terraform from reverting the changes. This approach has limitations:

- **Configuration drift**: Actual cluster state diverges from Terraform configuration
- **Maintenance overhead**: Careful management of ignore_changes blocks required
- **Limited visibility**: Actual scaled values cannot be easily inspected
- **Module inflexibility**: Lifecycle blocks cannot be conditional, requiring separate module implementations

### How use_effective_fields Works

The `use_effective_fields` attribute changes how the provider handles specification attributes:

**With `use_effective_fields = true`:**
- **Clear separation of concerns**:
  - Specification attributes (`electable_specs`, `analytics_specs`, `read_only_specs`) remain exactly as defined in your configuration
  - Atlas-computed values are available separately in effective specs (`effective_electable_specs`, `effective_analytics_specs`, `effective_read_only_specs`)
- **No plan drift**: Atlas auto-scaling changes do not affect your Terraform configuration
- **Visibility**: Use data sources to read effective specs showing actual provisioned values
- **Module flexibility**: Single module implementation works for both auto-scaling and non-auto-scaling scenarios

**Key principle:** Your configuration stays clean and represents your intent, while effective specs show the reality of what Atlas has provisioned.

## Benefits for Module Authors

1. **Single module implementation**: One codebase supports both auto-scaling and non-auto-scaling use cases
2. **No lifecycle blocks needed**: Cleaner, more maintainable code
3. **Operational visibility**: Expose both configured and actual (effective) values to module users
4. **Forward compatibility**: Aligns with provider v3.x where effective fields becomes default behavior

## Getting Started

1. **Creating a new module**: Start with [new_module](./new_module/)
2. **Migrating an existing module**: Review [existing_module](./existing_module/)

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields) - Complete documentation
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) - Full resource documentation
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster) - Data source for reading effective specs
