# v2: Migration Phase

This configuration demonstrates the migration approach by running both old and new data sources side-by-side.

## What this shows

- **Attribute mapping**: Direct comparison between old and new attribute structures
- **Validation**: Outputs that verify the new data sources return equivalent data
- **Migration readiness**: Checks to confirm you're ready to move to v3

## Key comparisons

### Single User Reads
- `email_address` → `username`
- Complex role filtering → Structured `roles.org_roles` and `roles.project_role_assignments`

### User Lists
- `results[*].email_address` → `users[*].username`
- `results` → `users`

## Usage

1. Apply this configuration: `terraform apply`
2. Review the comparison outputs to verify data consistency
3. Check `migration_validation.ready_for_v3` is `true`
4. Once validated, proceed to v3

## Expected outputs

The outputs will show side-by-side comparisons of:
- Email retrieval methods
- Role access patterns  
- User list structures
- Count validations

If `migration_validation.ready_for_v3` is `true`, you can safely proceed to the final v3 configuration.
