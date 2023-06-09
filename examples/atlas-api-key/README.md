# MongoDB Atlas Provider -- MongoDB Atlas Programmatic API Key Examples 

In the Terraform MongoDB Atlas Provider v1.8.0 release, we introduced support for MongoDB Atlas Programmatic API Keys (PAK). While this was an exciting development, we quickly realized that some of our customers were facing challenges in optimally leveraging this resource. In response to this, we initiated several enhancements as part of the v1.10.0 release to refine the user experience with this resource. These enhancements encompassed both code revisions and documentation updates.

The most notable improvement among these as part of v1.10.0 release was the deprecation of the `api_keys` parameter from the `mongodbatlas_project` resource. We also extended the functionality of the `mongodbatlas_project_api_key` resource by incorporating the `project_assignment` parameter. This enhancement removes the necessity for users to create multiple `mongodbatlas_project` resource blocks just to assign keys. This update streamlines the process of assigning an API Key to multiple projects, making it less cumbersome and more manageable.

To further facilitate the transition, we've included in this atlas-api-key folder, three reference examples for Programmatic API Key (PAK) usage:

* "Create and Assign PAK Together" — this demonstrates how to create a PAK and assign it simultaneously.

* "Create and Assign PAK to Multiple Projects" — this shows how to create a PAK and assign it to several projects at once.

* "Create and Assign PAK Separately" (Deprecated) — this is an older method of creating and assigning PAKs, now deprecated but still available for reference.

Our hope is that these examples will provide clear guidance and help ease your transition to this new PAK workflow in our Terraform Provider.

