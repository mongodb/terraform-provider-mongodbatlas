# MongoDB Atlas Provider -- MongoDB Atlas Programmatic API Key Examples 

In the Terraform MongoDB Atlas Provider v1.8.0 release, we introduced support for MongoDB Atlas Programmatic API Keys (PAK). While this was an exciting development, we quickly realized that some of our customers were facing challenges in optimally leveraging this resource. In response to this, we initiated several enhancements as part of the v1.10.0 release to refine the user experience with this resource. These enhancements encompassed both code revisions and documentation updates.

The most notable improvement among these as part of v1.10.0 release was the deprecation of the `api_keys` parameter from the `mongodbatlas_project` resource. We also extended the functionality of the `mongodbatlas_project_api_key` resource by incorporating the `project_assignment` parameter. This enhancement removes the necessity for users to create multiple `mongodbatlas_project` resource blocks just to assign keys. This update streamlines the process of assigning an API Key to multiple projects, making it less cumbersome and more manageable.

To further facilitate the transition, we've included in this atlas-api-key folder, three reference examples for Programmatic API Key (PAK) usage:

* "Create and Assign PAK Together" — this demonstrates how to create a PAK and assign it simultaneously.

* "Create and Assign PAK to Multiple Projects" — this shows how to create a PAK and assign it to several projects at once.

* "Create and Assign PAK Separately" (Deprecated) — this is an older method of creating and assigning PAKs, now deprecated but still available for reference.

Lastly, in MongoDB Atlas, all PAKs are Organization API keys. Once created, a PAK is linked at the organization level with an 'Organization Member' role. However, these Organization API keys can also be assigned to one or more projects within the organization. When a PAK is assigned to a specific project, it essentially takes on the 'Project Owner' role for that particular project. This enables the key to perform operations at the project level, in addition to the organization level. The flexibility of PAKs provides a powerful mechanism for fine-grained access and control, once their functioning is clearly understood. 

Our hope is that these examples will provide clear guidance and help ease your transition to this new PAK workflow in our Terraform Provider.

