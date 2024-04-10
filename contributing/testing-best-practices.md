
## Testing Best Practices

### Types of test

- Unit tests: In Terraform terminology they refer to tests that [validate a resource schema](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#unit-testing). That is done automatically [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/provider/provider_test.go) for all resources and data sources using Terraform Framework Plugin. Here we’re referring to the broader concept of testing a resource or unit without calling the external systems like the Atlas Go SDK.
- Acceptance (acc) tests: In Terraform terminology they refer to the use of real Terraform configurations to exercise the code in plan, apply, refresh, and destroy life cycles (real infrastructure resources are created as part of the test).
- Migration (mig) tests: These tests are designed to ensure that after an upgrade to a new Atlas provider version, user configs do not result in unexpected plan changes. Migration tests are a subset of Acceptance tests.

### File structure

- Unit and Acceptances tests are in the same `_test.go` file. They are not in the same package as the code tests, e.g. `advancedcluster` tests are in `advancedcluster_test` package so coupling is minimized.
- Migration tests are in `_migration_test.go` files.
- All resources need a `main_test.go` file to handle resource reuse lifecycle, e.g. [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/service/advancedcluster/main_test.go).
- Helper methods must have their own tests, e.g. `common_advanced_cluster_test.go` has tests for `common_advanced_cluster.go`.
- `internal/testutils/acc` contains helper test functions for Acceptance tests.
- `internal/testutils/mig` contains helper test functions specifically for Migration tests.
- `internal/testutils/replay` contains helper test functions for [Hoverfly](https://docs.hoverfly.io/en/latest/). Hoverfly is used to capture and replay HTTP traffic with Atlas Cloud to speed up local development process.


### Local development

- Many tests support use of environment variables `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_CLUSTER_NAME` to resuse an exisiting project or cluster when runnning tests. This significantly reduces run duration for those tests.
- Go test cache can be used without any special setup.

### Unit tests

- Unit tests must not create Terraform resources or use external systems like [Atlas Go SDK](https://github.com/mongodb/atlas-sdk-go)
- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock) is used for test doubles
- Altlas Go SDK mocked interfaces are generated in [mockadmin](https://github.com/mongodb/atlas-sdk-go/tree/main/mockadmin) package using [Mockery](https://github.com/vektra/mockery)

### Acceptance tests

- There must be at least one `basic acceptance test` for each resource
- `Basic import tests` are done as the last step in the `basic acceptance tests`, not as a different test. Exceptions apply for more specific import tests, e.g. testing with incorrect IDs.
- Data sources are tested in the same tests as the resources. There are not specific test files or tests for data sources as a resource must typically be created. (There are very few exceptions to this, e.g. when there is only data sources but not resource)
- Main way to reduce use of projects is `ProjectIDExecution`. This function returns the ID of a project that is created for the current execution of tests for a resource.
  - As the project is shared for all acceptance tests for a resource, sometimes tests can affect each other if they're using global resources to the project (e.g. network peering, maintenance window or LDAP config). In that case:
    - Run the tests in serial (`resource.Test` instead of `resource.ParallelTest`) if the tests are fast and saving resources is prefered.
    - Don’t use `ProjectIDExecution` and create a project for each test if a faster test execution is prefered even if more resources are needed.
- Main way to reduce use of clusters is `ClusterNameExecution`. This function returns the name of a cluster that is created for the current execution of tests for a resource. Similar precautions to project reuse must be taken here. If a global resource to cluster is being tested (e.g. cluster global config) then it's prefered to run tests in serial or create their own clusters.
- Plural data sources can be challeging to test when tests run in parallel or they share projects and/or clusters. In that case:
  - Don't check for a specific total count as other tests can also create resources. As an example you can check that there are results using `acc.IntGreatThan(0)`.
  - Don't assume results are in a certain order, use functions like `resource.TestCheckTypeSetElemNestedAttrs` or `resource.TestCheckTypeSetElemAttr`.

### Migration tests

- Migration tests are also acceptance tests so most of the info above also applies here, e.g. use of `ProjectIDExecution` or `ClusterNameExecution`.
- There must be at least one `basic migration test` for each resource that leverages on the `basic acceptance tests` using helper test functions such as `CreateAndRunTest`.
