
## Testing Best Practices

### Types of test

- Unit tests: In Terraform terminology they refer to tests that [validate a resource schema](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#unit-testing). That is done automatically [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/provider/provider_test.go) for all resources and data sources using Terraform Framework Plugin. Here we’re referring to the broader concept of testing a resource or unit without calling the external systems like the Atlas Go SDK.
- Acceptance (acc) tests: In Terraform terminology they refer to the use of real Terraform configurations to exercise the code in plan, apply, refresh, and destroy life cycles (real infrastructure resources are created as part of the test).
- Migration (mig) tests: These tests are designed to ensure that after an upgrade to a new Atlas provider version, user configs do not result in unexpected plan changes. Migration tests are a subset of Acceptance tests.

### File structure

- Unit and Acceptances tests are in the same `_test.go` file. They are not in the same package as the code tests, e.g. `advancedcluster` tests are in `advancedcluster_test` package so coupling is minimized.
- Migration tests are in `_migration_test.go` files.
- Helper methods must have their own tests, e.g. `common_advanced_cluster_test.go` has tests for `common_advanced_cluster.go`.
- `internal/testutils/acc` contains helper test methods for Acceptance tests.
- `internal/testutils/mig` contains helper test methods specifically for Migration tests.
- `internal/testutils/replay` contains helper test methods for [Hoverfly](https://docs.hoverfly.io/en/latest/). Hoverfly is used to capture and replay HTTP traffic with Atlas Cloud to speed up local development process.

### Unit tests
- Unit tests must not create Terraform resources or use external systems like [Atlas Go SDK](https://github.com/mongodb/atlas-sdk-go).
- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock) and [Mockery](https://github.com/vektra/mockery) are used for test doubles in Atlas Go SDK unit tests.

### Acceptance tests

### Migration tests

- Migration tests are also acceptance tests so most of the info above also applies here
