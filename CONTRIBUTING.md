Contributing
---------------------------

# Contributing

## Workflow

MongoDB welcomes community contributions! If you’re interested in making a contribution to  Terraform MongoDB Atlas provider, please follow the steps below before you start writing any code:

1. Sign the [contributor's agreement](http://www.mongodb.com/contributor). This will allow us to review and accept contributions.
1. Read the [Terraform contribution guidelines](https://www.terraform.io/docs/extend/community/contributing.html) and the [README](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/README.md) in this repo
1. Reach out by filing an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues) to discuss your proposed contribution, be it a bug fix or feature/other improvements.  

After the above 3 steps are completed and we've agreed on a path forward:
1. Fork the repository on GitHub
1. Create a branch with a name that briefly describes your submission
1. Implement your feature, improvement or bug fix, ensuring it adheres to the [Terraform Plugin Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html)
1. Ensure you follow the [Terraform Plugin Testing requirements](https://www.terraform.io/docs/extend/testing/index.html).
1. Add comments around your new code that explain what's happening
1. Commit and push your changes to your branch then submit a pull request against the current release branch, not master.  The naming scheme of the branch is `release-staging-v#.#.#`. Note: There will only be one release branch at a time.  
1. A repo maintainer will review the your pull request, and may either request additional changes or merge the pull request.

## PR Title Format
We use [*Conventional Commits*](https://www.conventionalcommits.org/):
- `fix: description of the PR`: a commit of the type fix patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- `chore: description of the PR`: the commit includes a technical or preventative maintenance task that is necessary for managing the product or the repository, but it is not tied to any specific feature or user story (this correlates with PATCH in Semantic Versioning).
- `doc: description of the PR`: The commit adds, updates, or revises documentation that is stored in the repository (this correlates with PATCH in Semantic Versioning).
- `test: description of the PR`: The commit enhances, adds to, revised, or otherwise changes the suite of automated tests for the product (this correlates with PATCH in Semantic Versioning).
- `security: description of the PR`: The commit improves the security of the product or resolves a security issue that has been reported (this correlates with PATCH in Semantic Versioning).
- `refactor: description of the PR`: The commit refactors existing code in the product, but does not alter or change existing behavior in the product (this correlates with Minor in Semantic Versioning).
- `perf: description of the PR`: The commit improves the performance of algorithms or general execution time of the product, but does not fundamentally change an existing feature (this correlates with Minor in Semantic Versioning).
- `ci: description of the PR`: The commit makes changes to continuous integration or continuous delivery scripts or configuration files (this correlates with Minor in Semantic Versioning).
- `revert: description of the PR`: The commit reverts one or more commits that were previously included in the product, but were accidentally merged or serious issues were discovered that required their removal from the main branch (this correlates with Minor in Semantic Versioning).
- `style: description of the PR`: The commit updates or reformats the style of the source code, but does not otherwise change the product implementation (this correlates with Minor in Semantic Versioning).
- `feat: description of the PR`: a commit of the type feat introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- `deprecate: description of the PR`: The commit deprecates existing functionality, but does not remove it from the product (this correlates with MINOR in Semantic Versioning).
- `BREAKING CHANGE`: a commit that has a footer BREAKING CHANGE:, or appends a ! after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
Examples:
  - `fix!: description of the ticket`
  - If the PR has `BREAKING CHANGE`: in its description is a breaking change
- `remove!: description of the PR`: The commit removes a feature from the product. Typically features are deprecated first for a period of time before being removed. Removing a feature is a breaking change (correlating with MAJOR in Semantic Versioning).

## Terraform Plugin Framework Migration
Certain resources are being migrated to the new Terraform Plugin Framework from Terraform Plugin SDKv2.
Below conventions are followed for resources implemented with Terraform Plugin Framework:
- **File names:** All resource/data source and test files for resources implemented with Terraform Plugin Framework must be prefixed with `fw_`. For example, `fw_mongodbatlas_resource_project.go`, `fw_mongodbatlas_resource_project_test.go`and `fw_mongodbatlas_datasource_project.go`.
- **Data models:** Data models associated to a resource/data source schema must follow naming format - `tf<resourceName><RS|DS>Model`, where `RS|DS` demonstrates whether the model belongs to a Terraform resource (`RS`) or Data source (`DS`).
For example, `tfProjectRSModel` or `tfProjectDSModel`.


## Documentation Best Practises

1. In our documentation, when a resource field allows a maximum of only one item, we do not format that field as an array. Instead, we create a subsection specifically for this field. Within this new subsection, we enumerate all the attributes of the field. Let's illustrate this with an example: [cloud_backup_schedule.html.markdown](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/r/cloud_backup_schedule.html.markdown?plain=1#L207)
