# v2: Final State

This is the clean, final configuration using only the new `mongodbatlas_team_project_assignment` resource.

## What changed from v1

### Resource purpose
- **Old**: Managed through deprecated `mongodbatlas_project.teams` attribute
- **New**: Uses `mongodbatlas_team_project_assignment` resource for team-to-project assignments

### Data source support
- **Old**: Had data source for reading teams attribute
- **New**: Has data source for reading team assignments

## Usage patterns

This configuration demonstrates:
- Basic team assignment to project
- Data source usage for reading assignments
- Output examples showing how to print team assignments in various formats

## Migration complete

At this point, you have successfully migrated from the deprecated `mongodbatlas_project.teams` attribute to the modern `mongodbatlas_team_project_assignment` resource. All references to the old attribute have been removed and replaced with the new resource.