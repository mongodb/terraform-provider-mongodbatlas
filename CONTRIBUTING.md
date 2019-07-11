# Contributing to Terraform - MongoDB Atlas Provider

### This guide is on Work in Progress.

**First:** if you're unsure or afraid of _anything_, ask for help! You can
submit a work in progress (WIP) pull request, or file an issue with the parts
you know. We'll do our best to guide you in the right direction, and let you
know if there are guidelines we will need to follow. We want people to be able
to participate without fear of doing the wrong thing.

Below are our expectations for contributors. Following these guidelines gives us
the best opportunity to work with you, by making sure we have the things we need
in order to make it happen. Doing your best to follow it will speed up our
ability to merge PRs and respond to issues.

<!-- TOC depthFrom:2 -->

- [Issues](#issues)
    - [Issue Reporting Checklists](#issue-reporting-checklists)
        - [Bug Reports](#bug-reports)
        - [Feature Requests](#feature-requests)
        - [Questions](#questions)
    - [Issue Lifecycle](#issue-lifecycle)
- [Pull Requests](#pull-requests)
    - [Pull Request Lifecycle](#pull-request-lifecycle)
    - [Checklists for Contribution](#checklists-for-contribution)
        - [Documentation Update](#documentation-update)
        - [Enhancement/Bugfix to a Resource](#enhancementbugfix-to-a-resource)
        - [New Resource](#new-resource)
    - [Common Review Items](#common-review-items)
        - [Go Coding Style](#go-coding-style)
        - [Acceptance Testing Guidelines](#acceptance-testing-guidelines)
    - [Writing Acceptance Tests](#writing-acceptance-tests)
        - [Acceptance Tests Often Cost Money to Run](#acceptance-tests-often-cost-money-to-run)
        - [Running an Acceptance Test](#running-an-acceptance-test)
        - [Writing an Acceptance Test](#writing-an-acceptance-test)

<!-- /TOC -->

## Issues

### Issue Reporting Checklists

We welcome issues of all kinds including feature requests, bug reports, and
general questions. Below you'll find checklists with guidelines for well-formed
issues of each type.

#### Bug Reports

 - [ ] __Test against latest release__: Make sure you test against the latest
   released version. It is possible we already fixed the bug you're experiencing.

 - [ ] __Search for possible duplicate reports__: It's helpful to keep bug
   reports consolidated to one thread, so do a quick search on existing bug
   reports to check if anybody else has reported the same thing. You can [scope
      searches by the label "bug" to help narrow things down.

 - [ ] __Include steps to reproduce__: Provide steps to reproduce the issue,
   along with your `.tf` files, with secrets removed, so we can try to
   reproduce it. Without this, it makes it much harder to fix the issue.

 - [ ] __For panics, include `crash.log`__: If you experienced a panic, please
   create a [gist](https://gist.github.com) of the *entire* generated crash log
   for us to look at. Double check no sensitive items were in the log.

#### Feature Requests

 - [ ] __Search for possible duplicate requests__: It's helpful to keep requests
   consolidated to one thread, so do a quick search on existing requests to
   check if anybody else has reported the same thing. You can [scope searches by
      the label "enhancement"](https://github.com/terraform-providers/) to help narrow things down.

 - [ ] __Include a use case description__: In addition to describing the
   behavior of the feature you'd like to see added, it's helpful to also lay
   out the reason why the feature would be important and how it would benefit
   Terraform users.

#### Questions

 - [ ] __Search for answers in Terraform documentation__: We're happy to answer
   questions in GitHub Issues, but it helps reduce issue churn and maintainer
   workload if you work to [find answers to common questions in the
   documentation](https://www.terraform.io/docs/providers/). Oftentimes Question issues result in documentation updates
   to help future users, so if you don't find an answer, you can give us
   pointers for where you'd expect to see it in the docs.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a Terraform collaborator.
   Categorization is done via GitHub labels. We generally use a two-label
   system of (1) issue/PR type, and (2) section of the codebase. Type is
   one of "bug", "enhancement", "documentation", or "question", and section
   is usually the MongoDB Atlas service name.

3. An initial triage process determines whether the issue is critical and must
    be addressed immediately, or can be left open for community discussion.

4. The issue is addressed in a pull request or commit. The issue number will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.

5. The issue is closed. Sometimes, valid issues will be closed because they are
   tracked elsewhere or non-actionable. The issue is still indexed and
   available for future viewers, or can be re-opened if necessary.

## Pull Requests

We appreciate direct contributions to the provider codebase. Here's what to
expect:

 * For pull requests that follow the guidelines, we will proceed to reviewing
  and merging, following the provider team's review schedule. There may be some
  internal or community discussion needed before we can complete this.
 * Pull requests that don't follow the guidelines will be commented with what
  they're missing. The person who submits the pull request or another community
  member will need to address those requests before they move forward.

### Pull Request Lifecycle

1. [Fork the GitHub repository](https://help.github.com/en/articles/fork-a-repo),
   modify the code, and [create a pull request](https://help.github.com/en/articles/creating-a-pull-request-from-a-fork).
   You are welcome to submit your pull request for commentary or review before
   it is fully completed by creating a [draft pull request](https://help.github.com/en/articles/about-pull-requests#draft-pull-requests)
   or adding `[WIP]` to the beginning of the pull request title.
   Please include specific questions or items you'd like feedback on.

1. Once you believe your pull request is ready to be reviewed, ensure the
   pull request is not a draft pull request by [marking it ready for review](https://help.github.com/en/articles/changing-the-stage-of-a-pull-request)
   or removing `[WIP]` from the pull request title if necessary, and a
   maintainer will review it. Follow [the checklists below](#checklists-for-contribution)
   to help ensure that your contribution can be easily reviewed and potentially
   merged.

1. One of Terraform's provider team members will look over your contribution and
   either approve it or provide comments letting you know if there is anything
   left to do. We do our best to keep up with the volume of PRs waiting for
   review, but it may take some time depending on the complexity of the work.

1. Once all outstanding comments and checklist items have been addressed, your
   contribution will be merged! Merged PRs will be included in the next
   Terraform release. The provider team takes care of updating the CHANGELOG as
   they merge.

1. In some cases, we might decide that a PR should be closed without merging.
   We'll make sure to provide clear reasoning when this happens.

### Checklists for Contribution

There are several different kinds of contribution, each of which has its own
standards for a speedy review. The following sections describe guidelines for
each type of contribution.

#### Documentation Update

The [Terraform MongoDB Atlas Provider's website source][website] is in this repository
along with the code and tests. Below are some common items that will get
flagged during documentation reviews:

- [ ] __Reasoning for Change__: Documentation updates should include an explanation for why the update is needed.
- [ ] __Prefer MongoDB Atlas Documentation__: Documentation about MongoDB Atlas service features and valid argument values that are likely to update over time should link to MongoDB Atlas service user guides and API references where possible.
- [ ] __Large Example Configurations__: Example Terraform configuration that includes multiple resource definitions should be added to the repository `examples` directory instead of an individual resource documentation page. Each directory under `examples` should be self-contained to call `terraform apply` without special configuration.
- [ ] __Terraform Configuration Language Features__: Individual resource documentation pages and examples should refrain from highlighting particular Terraform configuration language syntax workarounds or features such as `variable`, `local`, `count`, and built-in functions.

#### Enhancement/Bugfix to a Resource

Working on existing resources is a great way to get started as a Terraform
contributor because you can work within existing code and tests to get a feel
for what to do.

In addition to the below checklist, please see the [Common Review
Items](#common-review-items) sections for more specific coding and testing
guidelines.

 - [ ] __Acceptance test coverage of new behavior__: Existing resources each
   have a set of [acceptance tests][acctests] covering their functionality.
   These tests should exercise all the behavior of the resource. Whether you are
   adding something or fixing a bug, the idea is to have an acceptance test that
   fails if your code were to be removed. Sometimes it is sufficient to
   "enhance" an existing test by adding an assertion or tweaking the config
   that is used, but it's often better to add a new test. You can copy/paste an
   existing test and follow the conventions you see there, modifying the test
   to exercise the behavior of your code.
 - [ ] __Documentation updates__: If your code makes any changes that need to
   be documented, you should include those doc updates in the same PR. This
   includes things like new resource attributes or changes in default values.
   The [Terraform website][website] source is in this repo and includes
   instructions for getting a local copy of the site up and running if you'd
   like to preview your changes.
 - [ ] __Well-formed Code__: Do your best to follow existing conventions you
   see in the codebase, and ensure your code is formatted with `go fmt`. (The
   Travis CI build will fail if `go fmt` has not been run on incoming code.)
   The PR reviewers can help out on this front, and may provide comments with
   suggestions on how to improve the code.
 - [ ] __Vendor additions__: Create a separate PR if you are updating the vendor
   folder. This is to avoid conflicts as the vendor versions tend to be fast-
   moving targets. We will plan to merge the PR with this change first.

#### New Resource

Implementing a new resource is a good way to learn more about how Terraform
interacts with upstream APIs. There are plenty of examples to draw from in the
existing resources, but you still get to implement something completely new.

In addition to the below checklist, please see the [Common Review
Items](#common-review-items) sections for more specific coding and testing
guidelines.

 - [ ] __Minimal LOC__: It's difficult for both the reviewer and author to go
   through long feedback cycles on a big PR with many resources. We ask you to
   only submit **1 resource at a time**.
 - [ ] __Acceptance tests__: New resources should include acceptance tests
   covering their behavior. See [Writing Acceptance
   Tests](#writing-acceptance-tests) below for a detailed guide on how to
   approach these.
 - [ ] __Naming__: Resources should be named `mongodbatlas_<service>_<name>` where
   `service` is the MongoDB Atlas short service name and `name` is a short, preferably
   single word, description of the resource. Use `_` as a separator.
 - [ ] __Arguments_and_Attributes__: The HCL for arguments and attributes should
   mimic the types and structs presented by the MongoDB Atlas API. API arguments should be
   converted from `CamelCase` to `camel_case`.
 - [ ] __Documentation__: Each resource gets a page in the Terraform
   documentation. The [Terraform website][website] source is in this
   repo and includes instructions for getting a local copy of the site up and
   running if you'd like to preview your changes. For a resource, you'll want
   to add a new file in the appropriate place and add a link to the sidebar for
   that page.
 - [ ] __Well-formed Code__: Do your best to follow existing conventions you
   see in the codebase, and ensure your code is formatted with `go fmt`. (The
   Travis CI build will fail if `go fmt` has not been run on incoming code.)
   The PR reviewers can help out on this front, and may provide comments with
   suggestions on how to improve the code.
 - [ ] __Vendor updates__: Create a separate PR if you are adding to the vendor
   folder. This is to avoid conflicts as the vendor versions tend to be fast-
   moving targets. We will plan to merge the PR with this change first.


### Common Review Items

The Terraform MongoDB Atlas Provider follows common practices to ensure consistent and
reliable implementations across all resources in the project. While there may be
older resource and testing code that predates these guidelines, new submissions
are generally expected to adhere to these items to maintain Terraform Provider
quality. For any guidelines listed, contributors are encouraged to ask any
questions and community reviewers are encouraged to provide review suggestions
based on these guidelines to speed up the review and merge process.

#### Go Coding Style

The following Go language resources provide common coding preferences that may be referenced during review, if not automatically handled by the project's linting tools.

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### Acceptance Testing Guidelines

The below are required items that will be noted during submission review and prevent immediate merging:

- [ ] __Implements CheckDestroy__: Resource testing should include a `CheckDestroy` function (typically named `testAccCheckMongoDBAtlas{RESOURCE}Destroy`) that calls the API to verify that the Terraform resource has been deleted or disassociated as appropriate. More information about `CheckDestroy` functions can be found in the [Extending Terraform TestCase documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html#checkdestroy).
- [ ] __Implements Exists Check Function__: Resource testing should include a `TestCheckFunc` function (typically named `testAccCheckMongoDBAtlas{RESOURCE}Exists`) that calls the API to verify that the Terraform resource has been created or associated as appropriate. Preferably, this function will also accept a pointer to an API object representing the Terraform resource from the API response that can be set for potential usage in later `TestCheckFunc`. More information about these functions can be found in the [Extending Terraform Custom Check Functions documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html#checkdestroy).

- [ ] __Uses resource.Test__:  [`resource.Test()`](https://godoc.org/github.com/hashicorp/terraform/helper/resource#Test) where serialized testing is absolutely required.
- [ ] __Uses fmt.Sprintf()__: Test configurations preferably should to be separated into their own functions (typically named `testAccMongoDBAtlas{RESOURCE}Config{PURPOSE}`) that call [`fmt.Sprintf()`](https://golang.org/pkg/fmt/#Sprintf) for variable injection or a string `const` for completely static configurations. Test configurations should avoid `var` or other variable injection functionality such as [`text/template`](https://golang.org/pkg/text/template/).
- [ ] __Uses Randomized Infrastructure Naming__: Test configurations that utilize resources where a unique name is required should generate a random name. Typically this is created via `rName := acctest.RandomWithPrefix("tf-acc-test")` in the acceptance test function before generating the configuration.

For resources that support import, the additional item below is required that will be noted during submission review and prevent immediate merging:

- [ ] __Implements ImportState Testing__: Tests should include an additional `TestStep` configuration that verifies resource import via `ImportState: true` and `ImportStateVerify: true`. This `TestStep` should be added to all possible tests for the resource to ensure that all infrastructure configurations are properly imported into Terraform.

The below are style-based items that _may_ be noted during review and are recommended for simplicity, consistency, and quality assurance:

- [ ] __Uses Builtin Check Functions__: Tests should utilize already available check functions, e.g. `resource.TestCheckResourceAttr()`, to verify values in the Terraform state over creating custom `TestCheckFunc`. More information about these functions can be found in the [Extending Terraform Builtin Check Functions documentation](https://www.terraform.io/docs/extend/testing/acceptance-tests/teststep.html#builtin-check-functions).
- [ ] __Uses TestCheckResoureAttrPair() for Data Sources__: Tests should utilize [`resource.TestCheckResourceAttrPair()`](https://godoc.org/github.com/hashicorp/terraform/helper/resource#TestCheckResourceAttrPair) to verify values in the Terraform state for data sources attributes to compare them with their expected resource attributes.
- [ ] __Excludes Timeouts Configurations__: Test configurations should not include `timeouts {...}` configuration blocks except for explicit testing of customizable timeouts (typically very short timeouts with `ExpectError`).
- [ ] __Implements Default and Zero Value Validation__: The basic test for a resource (typically named `TestAccMongoDBAtlas{RESOURCE}_basic`) should utilize available check functions, e.g. `resource.TestCheckResourceAttr()`, to verify default and zero values in the Terraform state for all attributes. Empty/missing configuration blocks can be verified with `resource.TestCheckResourceAttr(resourceName, "{ATTRIBUTE}.#", "0")` and empty maps with `resource.TestCheckResourceAttr(resourceName, "{ATTRIBUTE}.%", "0")`

### Writing Acceptance Tests

Terraform includes an acceptance test harness that does most of the repetitive
work involved in testing a resource. For additional information about testing
Terraform Providers, see the [Extending Terraform documentation](https://www.terraform.io/docs/extend/testing/index.html).

#### Acceptance Tests Often Cost Money to Run

Because acceptance tests create real resources, they often cost money to run.
Because the resources only exist for a short period of time, the total amount
of money required is usually a relatively small. Nevertheless, we don't want
financial limitations to be a barrier to contribution, so if you are unable to
pay to run acceptance tests for your contribution, mention this in your
pull request. We will happily accept "best effort" implementations of
acceptance tests and run them for you on our side. This might mean that your PR
takes a bit longer to merge, but it most definitely is not a blocker for
contributions.

#### Running an Acceptance Test

Acceptance tests can be run using the `testacc` target in the Terraform
`Makefile`. The individual tests to run can be controlled using a regular
expression. Prior to running the tests provider configuration details such as
access keys must be made available as environment variables.

For example, to run an acceptance test against the MongoDB Atlas
provider, the following environment variables must be set:

```sh
# Using a profile
export MONGODB_ATLAS_PUBLIC_KEY=...
export MONGODB_ATLAS_PRIVATE_KEY=...
```

Tests can then be run by specifying the target provider and a regular
expression defining the tests to run:

```sh
$ TF_LOG=DEBUG make testacc TESTARGS='-run=TestAccResourceMongoDBAtlasCluster_basic'
==> Checking that code complies with gofmt requirements...
TF_ACC=1 go test ./mongodbatlas -v -run=TestAccResourceMongoDBAtlasCluster_basic -timeout 120m
=== RUN   TestAccResourceMongoDBAtlasCluster_basic
--- PASS: TestAccResourceMongoDBAtlasCluster_basic (26.56s)
PASS
ok  	github.com/terraform-providers/mongodbaltas	26.607s
```

#### Writing an Acceptance Test

Terraform has a framework for writing acceptance tests which minimises the
amount of boilerplate code necessary to use common testing patterns. The entry
point to the framework is the `resource.Test()` function.

Tests are divided into `TestStep`s. Each `TestStep` proceeds by applying some
Terraform configuration using the provider under test, and then verifying that
results are as expected by making assertions using the provider API. It is
common for a single test function to exercise both the creation of and updates
to a single resource. Most tests follow a similar structure.

1. Pre-flight checks are made to ensure that sufficient provider configuration
   is available to be able to proceed - for example in an acceptance test
   targeting MongoDB Atlas, `MONGODB_ATLAS_PUBLIC_KEY` and `MONGODB_ATLAS_PRIVATE_KEY` must be set prior
   to running acceptance tests. This is common to all tests exercising a single
   provider.

```go
func TestAccResourceMongoDBAtlasCluster_basic(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfig(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}
```

When executing the test, the following steps are taken for each `TestStep`:

1. The Terraform configuration required for the test is applied. This is
   responsible for configuring the resource under test, and any dependencies it
   may have. For example, to test the `mongodbatlas_cluster` resource, a valid configuration with the requisite fields is required. This results in configuration which looks like this:

    ```hcl
    resource "mongodbatlas_cluster" "test" {
      project_id   = "<PROJECT_ID>"
      name         = "<NAME>"
      disk_size_gb = 100
      num_shards   = 1
      
      replication_factor           = 3
      backup_enabled               = true
      auto_scaling_disk_gb_enabled = true
      
      //Provider Settings "block"
      provider_name               = "AWS"
      provider_disk_iops          = 300
      provider_encrypt_ebs_volume = false
      provider_instance_size_name = "M40"
      provider_region_name        = "US_EAST_1"
    }
    ```

1. Assertions are run using the provider API. These use the provider API
   directly rather than asserting against the resource state. For example, to
   verify that the `mongodbatlas_cluster` described above was created
   successfully, a test function like this is used:

    ```go
      func testAccCheckMongoDBAtlasClusterExists(resourceName string, cluster *matlas.Cluster) resource.TestCheckFunc {
        return func(s *terraform.State) error {
          conn := testAccProvider.Meta().(*matlas.Client)

          rs, ok := s.RootModule().Resources[resourceName]
          if !ok {
            return fmt.Errorf("not found: %s", resourceName)
          }
          if rs.Primary.ID == "" {
            return fmt.Errorf("no ID is set")
          }

          log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

          if clusterResp, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]); err == nil {
            *cluster = *clusterResp
            return nil
          }

          return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
        }
      }
    ```

   Notice that the only information used from the Terraform state is the ID of
   the resource. For computed properties, we instead assert that the value saved in the Terraform state was the
   expected value if possible. The testing framework provides helper functions
   for several common types of check - for example:

    ```go
    resource.TestCheckResourceAttrSet(resourceName, "project_id"),
    ```

2. The resources created by the test are destroyed. This step happens
   automatically, and is the equivalent of calling `terraform destroy`.

3. Assertions are made against the provider API to verify that the resources
   have indeed been removed. If these checks fail, the test fails and reports
   "dangling resources". The code to ensure that the `mongodbatlas_cluster` shown
   above has been destroyed looks like this:

    ```go
    func testAccCheckMongoDBAtlasClusterDestroy(s *terraform.State) error {
      conn := testAccProvider.Meta().(*matlas.Client)

      for _, rs := range s.RootModule().Resources {
        if rs.Type != "mongodbatlas_cluster" {
          continue
        }

        // Try to find the cluster
        _, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"])

        if err == nil {
          return fmt.Errorf("cluster (%s:%s) still exists", rs.Primary.Attributes["name"], rs.Primary.ID)
        }
      }

      return nil
    }
    ```

[website]: https://github.com/terraform-providers/
[acctests]: https://github.com/hashicorp/terraform#acceptance-tests
[ml]: https://groups.google.com/group/terraform-tool
