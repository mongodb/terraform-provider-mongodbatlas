---
layout: "mongodbatlas"
page_title: "Upgrade Guide for Terraform MongoDB Atlas Provider PAK Resource in v1.10.0"
sidebar_current: "docs-mongodbatlas-guides-PAK-upgrade-guide"
description: |-
MongoDB Atlas Provider : Upgrade and Information Guide
---

# MongoDB Atlas Provider: PAK Upgrade Guide in v1.10.0
In Terraform MongoDB Atlas Provider v1.10.0, some improvements were introduced which mainly focus on the MongoDB Atlas Programmatic API Keys (PAK) handling. This guide will help you to transition smoothly from the previous version which this resource was first released (v1.8.0) to the new version (v1.10.0).

## Changes Overview
* `api_keys` parameter is deprecated from the mongodbatlas_project resource.
* The `mongodbatlas_project_api_key` resource is extended to include a `project_assignment` parameter.

## Upgrade Steps
1. Replace `api_keys` in `mongodbatlas_project`: The `api_keys` parameter in the `mongodbatlas_project` resource is deprecated. Remove any instances of `api_keys` in your current Terraform scripts.

2. Use `project_assignment` in `mongodbatlas_project_api_key`: Instead of creating multiple mongodbatlas_project resources for each API Key assignment, you can now assign a PAK to multiple projects in a single resource block using the new `project_assignment` parameter in `mongodbatlas_project_api_key`.

## Examples
We provide three examples in the atlas-api-key folder (under examples) to help you understand the new changes:

* "Create and Assign PAK Together": This example demonstrates how to create a Programmatic API Key and assign it to a project within the same resource block. This is a good place to start if you're used to creating and assigning keys at the same time.

* "Create and Assign PAK to Multiple Projects": This example shows how to create a PAK and assign it to several projects at once using the new `project_assignment` parameter. This is useful if you have multiple projects that require the same key.

* "Create and Assign PAK Separately" (Deprecated): This is the older method of creating and assigning PAKs, which is now deprecated. This example remains for reference and to help you understand the changes that have been made.

Before making any changes, please ensure you have thoroughly read and understood these examples, and that your current Terraform scripts align with the new PAK workflow.

Remember, your scripts will still work with deprecated features for now, but it's best to upgrade as soon as possible to benefit from the latest enhancements. Code removal is planned for v1.12.0 at which point prior PAK workflow will no longer function.

If you have any questions or face any issues during the migration, feel free to reach out to us by creating a GitHub Issue or PR in our repo. Thank you.  
