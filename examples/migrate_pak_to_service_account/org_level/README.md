# Combined Example: Organization-Level PAK → Service Account

This combined example is organized into step subfolders (v1–v3):

- v1/: Initial state with:
  - `mongodbatlas_api_key` (organization-level PAK),
  - `mongodbatlas_api_key_project_assignment` (project assignment), and
  - `mongodbatlas_access_list_api_key` (IP access list entry).
- v2/: Migration step showcasing the intermediate state:
  - add `mongodbatlas_service_account` (organization-level SA),
  - add `mongodbatlas_service_account_project_assignment` (project assignment),
  - add `mongodbatlas_service_account_access_list_entry` (IP access list entry),
  - keep existing PAK resources alongside SA resources for testing.
- v3/: Cleaned-up final configuration after v2 is applied:
  - remove all PAK resources (`mongodbatlas_api_key`, `mongodbatlas_api_key_project_assignment`, `mongodbatlas_access_list_api_key`),
  - keep only SA resources (`mongodbatlas_service_account`, `mongodbatlas_service_account_project_assignment`, `mongodbatlas_service_account_access_list_entry`).

Navigate into each version folder to see the step-specific configuration.

