
# Testing Best Practices

## Types of test

- Unit tests: In Terraform terminology they refer to tests that [validate a resource schema](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas#unit-testing). That is done automatically [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/provider/provider_test.go) for all resources and data sources using Terraform Framework Plugin. Additionally, we have general unit testing for testing a resource or unit without calling external systems like MongoDB Atlas.
- Acceptance (acc) tests: In Terraform terminology they refer to the use of real Terraform configurations to exercise the code in plan, apply, refresh, and destroy life cycles (real infrastructure resources are created as part of the test), more info [here](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests).
- Migration (mig) tests: These tests are designed to ensure that after an upgrade to a new Atlas provider version, user configs do not result in unexpected plan changes, more info [here](https://developer.hashicorp.com/terraform/plugin/framework/migrating/testing). Migration tests are a subset of Acceptance tests.

## Test Organization
- A resource and associated data sources are implemented in a folder that is also a Go package, e.g. `advancedcluster` implementation is in [`internal/service/advancedcluster`](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/internal/service/advancedcluster)
- We enforce "black box" testing, tests must be in a separate "_test" package, e.g. `advancedcluster` tests are in `advancedcluster_test` package.
- Acceptance and general unit tests are in corresponding  `_test.go` file as the resource or data source source file.  If business logic is extracted into a separate file, unit testing for that logic will be including in its associated `_test.go` file, e.g. [state_transition_search_deployment_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/searchdeployment/state_transition_search_deployment_test.go).
- Migration tests are in `_migration_test.go` files.
- When functions are in their own file because they are shared by resource and data sources, a test file can be created to test them, e.g. [model_alert_configuration_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/alertconfiguration/model_alert_configuration_test.go) has tests for [model_alert_configuration](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/alertconfiguration/model_alert_configuration.go).
- All resource folders must have a `main_test.go` file to handle resource reuse lifecycle, e.g. [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/f3ff5bb678c1b07c16cc467471f483e483565427/internal/service/advancedcluster/main_test.go).
- `internal/testutils/acc` contains helper test functions for Acceptance tests.
- `internal/testutils/mig` contains helper test functions specifically for Migration tests.
- `internal/testutils/replay` contains helper test functions for [Hoverfly](https://docs.hoverfly.io/en/latest/). Hoverfly is used to capture and replay HTTP traffic with  MongoDB Atlas.

## Unit tests

- Unit tests must not create Terraform resources or use external systems, e.g unit tests using [Atlas Go SDK](https://github.com/mongodb/atlas-sdk-go) must not call MongoDB Atlas.
- Mock of specific interfaces like `admin.ProjectsApi` is prefered to use the whole `client`, e.g. [resource_project.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/project/resource_project.go#L792)
- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock) is used for test doubles.
- Altlas Go SDK mocked interfaces are generated in [mockadmin](https://github.com/mongodb/atlas-sdk-go/tree/main/mockadmin) package using [Mockery](https://github.com/vektra/mockery), example of use in [resource_project_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/project/resource_project_test.go#L114).

## Acceptance tests

- There must be at least one `basic acceptance test` for each resource, e.g.: [TestAccSearchIndex_basic](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/searchindex/resource_search_index_test.go#L14). They test the happy path with minimum resource configuration.
- `Basic import tests` are done as the last step in the `basic acceptance tests`, not as a different test, e.g. [basicTestCase](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/searchindex/resource_search_index_test.go#L211). Exceptions apply for more specific import tests, e.g. testing with incorrect IDs. [Import tests](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/import#resource-acceptance-testing-implementation) verify that the [Terraform Import](https://developer.hashicorp.com/terraform/cli/import) functionality is working fine.
- Data sources are tested in the same tests as the resources, e.g. [commonChecks](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/searchindex/resource_search_index_test.go#L262-L263).
- Helper functions such as `resource.TestCheckTypeSetElemNestedAttrs` or `resource.TestCheckTypeSetElemAttr` can be used to check resource and data source attributes more easily, e.g. [resource_serverless_instance_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/serverlessinstance/resource_serverless_instance_test.go#L61).

### Cloud Gov tests

1. Use [`PreCheck: PreCheckGovBasic`](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/CLOUDP-250271_cloud_gov/internal/testutil/acc/pre_check.go#L98)
2. Use the [`acc.ConfigGovProvider`](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/CLOUDP-250271_cloud_gov/internal/testutil/acc/provider.go#L61) together with your normal terraform config
3. Modify the `checkExist` and `CheckDestroy` to use `acc.ConnV2UsingGov`
4. Follow naming convention:
   1. `TestAccGovProject_withProjectOwner`, note prefix: `TestAccGov`
   2. `TestMigGovProject_regionUsageRestrictionsDefault`, note prefix: `TestMigGov`

## Migration tests

- There must be at least one `basic migration test` for each resource that leverages on the `basic acceptance tests` using helper test functions such as `CreateAndRunTest`, e.g. [TestMigServerlessInstance_basic](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/serverlessinstance/resource_serverless_instance_migration_test.go#L10).

## Local testing

These enviroment variables can be used in local to speed up development process.

Enviroment Variable | Description
--- | ---
`MONGODB_ATLAS_PROJECT_ID` | Re-use an existing project reducing test run duration for resources supporting this variable
`MONGODB_ATLAS_CLUSTER_NAME` | Re-use an existing cluster reducing significantly test run duration for resources supporting this variable
`REPLAY_MODE` | Use [Hoverfly](https://docs.hoverfly.io/en/latest/), more info about possible variable values [here](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/contributing/development-setup.md#replaying-http-requests-with-hoverfly)

## Shared resources

Acceptance and migration tests can reuse projects and clusters in order to be more efficient in resource utilization.

- A project can be reused using `ProjectIDExecution`. It returns the ID of a project that is created for the current execution of tests for a resource, e.g. [TestAccConfigRSDatabaseUser_basic](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/databaseuser/resource_database_user_test.go#L24).
  - As the project is shared for all tests for a resource, sometimes tests can affect each other if they're using global resources to the project (e.g. network peering, maintenance window or LDAP config). In that case:
    - Run the tests in serial (`resource.Test` instead of `resource.ParallelTest`) if the tests are fast and saving resources is prefered, e.g. [TestAccConfigRSProjectAPIKey_multiple](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/projectapikey/resource_project_api_key_test.go#L149).
    - Don’t use `ProjectIDExecution` and create a project for each test if a faster test execution is prefered even if more resources are needed, e.g. [TestAccFederatedDatabaseInstance_basic](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/federateddatabaseinstance/resource_federated_database_instance_test.go#L19).
- A cluster can be reused using `ClusterNameExecution`. This function returns the project id (created with `ProjectIDExecution`) and the name of a cluster that is created for the current execution of tests for a resource, e.g. [TestAccSearchIndex_withSearchType](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/searchindex/resource_search_index_test.go#L20). Similar precautions to project reuse must be taken here. If a global resource to cluster is being tested (e.g. cluster global config) then it's prefered to run tests in serial or create their own clusters.
- Plural data sources can be challenging to test when tests run in parallel or they share projects and/or clusters.
  - Avoid checking for a specific total count as other tests can also create resources, e.g. [resource_network_container_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/networkcontainer/resource_network_container_test.go#L214).
  - Don't assume results are in a certain order, e.g. [resource_network_container_test.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/66c44e62c9afe04ffe8be0dbccaec682bab830e6/internal/service/networkcontainer/resource_network_container_test.go#L215).
