# Migration Example: Project Teams → Team Project Assignment

This example demonstrates migrating from the deprecated `mongodbatlas_project.teams` attribute to `mongodbatlas_team_project_assignment`.

## Basic usage (direct resource usage)

- [basic](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_team_project_assignment/basic): Direct resource usage without modules.

## Module-based examples

- [module_maintainer](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_team_project_assignment/module_maintainer): How module maintainers should migrate.
- [module_user](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_team_project_assignment/module_user): How module consumers should migrate.
