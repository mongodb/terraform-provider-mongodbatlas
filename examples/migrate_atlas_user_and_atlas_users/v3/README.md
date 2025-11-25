# v3: Final State

This is the clean, final configuration using only the new data sources.

## What changed from v1

### Simpler attribute access
- `email_address` → `username`
- Complex role filtering → Direct access via `roles.org_roles` and `roles.project_role_assignments`
- `results[*]` → `users[*]`

### Cleaner code
- No complex list comprehensions for basic role access
- Structured role data instead of flat lists
- More intuitive attribute names

### Better performance
- Organization context required for user reads (more efficient API calls)
- Structured data reduces client-side filtering

## Key improvements

1. **Structured roles**: Organization and project roles are clearly separated
2. **Direct access**: No need to filter consolidated role lists
3. **Consistent naming**: `username` instead of `email_address`
4. **Better organization**: User lists come from their natural containers (org/project/team)

## Usage

This configuration represents the target state after migration. All references to deprecated data sources have been removed and replaced with their modern equivalents.
