# Releasing

## Steps

### Remove deprecated attributes

**Note**: Only applies if the right most version digit is 0 (considered a major or minor version in [semantic versioning](https://semver.org/)).

- If some deprecated attributes need to be removed in the following release, create a Jira ticket and merge the corresponding PR before starting the release workflow.
You can search in the code for the constansts in [deprecation.go](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/common/constant/deprecation.go) to find them.

### Revise jira release

Before triggering a release, view the corresponding [unreleased jira page](https://jira.mongodb.org/projects/CLOUDP?selectedItem=com.atlassian.jira.jira-projects-plugin:release-page&status=unreleased&contains=terraform) to ensure there are no pending tickets. In case there are pending tickets, verify with the team if the expectation is to have them included within the current release. After release workflow is successful the version will be marked as released automatically.

### Make sure that the acceptance tests are successful

While QA acceptance tests are run in the release process automatically, we check [workflows/test-suite.yml](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/test-suite.yml) and see if the latest run of the Test Suite action is successful (it runs every day at midnight UTC time). This can help detect failures before proceeding with the next steps.

### Verify upgrade guide is defined (if required)

- A doc ./docs/guides/X.Y.Z-upgrade-guide.md must be defined containing a summary of the most significant features, breaking changes, and additional information that can be helpful. The expectation is that this file is created during relevant pull requests (breaking changes, significant features), and not before the release process.

- We keep [Guides](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/docs/guides) only for 12 months. Add header `subcategory: "Older Guides"` to previous versions.

### Trigger release workflow

- Using our [Release GitHub Action](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/release.yml) run a new workflow using `master` and the following inputs:
  - Version number: `vX.Y.Z`
  - Skip QA acceptance tests: Should be left empty. Only used in case failing tests have been encountered in QA and the team agrees the release can still de done, or successful run of QA acceptance tests has already been done with the most recent changes.
  - Using an existing tag: Should be left empty (default `false` creates a new tag from `master`). This should be set to `true` only if you want to re-use an existing tag for the release. This can be helpful for rerunning a failed release process in which the tag has already been created.

#### How to create a pre-release (not part of regular process)
Pre-releases are not needed for a regular release process, but they can be generated for exceptional cases (e.g. sharing with external teams for additional testing). The process is the same as a regular release except for the format of the version number which must be `vX.Y.Z-pre` with additional numbers at the end if needed. When a pre release is triggered, steps related to updating the changelog header or updating the jira release version are skipped.

### Post-trigger checks
- You will see the release in the [GitHub Release page](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) once the [release action](.github/workflows/release.yml) has completed. HashiCorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the HashiCorp Terraform Registry (more details [here](https://developer.hashicorp.com/terraform/registry/providers/publishing#webhooks)).
- **CDKTF Update - Only for major release, i.e. the left most version digit increment (see this [comment](https://github.com/cdktf/cdktf-repository-manager/pull/202#issuecomment-1602562201))**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

## FAQ

**What happens if a release execution fails to create the tag but generated automatic commits into master?**

All steps before creating the tag are idempotent, meaning you can run the process again and no additional commits will be generated in the second run.

**What happens if a release execution creates a tag but fails during acceptance tests or creating the release (go releaser step)?**

Once a tag has been created in a previous execution, you can make use of the input `Using an existing tag` to run a new release process. 

Depending on the nature of the failure, you may want to introduce new changes into the current release. In this case you must:
- Delete the existing tag.
- Incorporate any new fixes into master.
- Manually trigger the [Generate Changelog workflow](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/generate-changelog.yml) to remove the current header and including any new entries that have been merged. This will run automatically if you have merged PRs after deleting the tag.
- Trigger a new release process, this will create a tag that includes your latest fixes.

**What happens if all the release process works except the last step to release the version in Jira?**

In this case there is no need to run the full release process again. Once the problem is found and fixed, the [Jira Release Version action](.github/workflows/jira-release-version.yml) can be run manually.
