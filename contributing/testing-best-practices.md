
## Testing Practices

### File structure

- Unit and Acceptances tests are in the same `_test.go` file. They are not in the same package as the code tests, e.g. `advancedcluster` tests are in `advancedcluster_test` package so coupling is minimized.
- Migration tests are in `_migration_test.go` files.
- Helper methods must have their own tests, e.g. `common_advanced_cluster_test.go` has tests for `common_advanced_cluster.go`.
- `internal/testutils/acc` contains helper test methods for Acceptance and Migration tests.

### Unit tests

- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock) and [Mockery](https://github.com/vektra/mockery) are used for test doubles in Atlas Go SDK unit tests.
