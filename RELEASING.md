# Releasing

## Steps

### Revise jira release

Before triggering a release, view the corresponding [unreleased jira page](https://jira.mongodb.org/projects/CLOUDP?selectedItem=com.atlassian.jira.jira-projects-plugin:release-page&status=unreleased&contains=terraform) to ensure there are no pending tickets. In case there are pending tickets, verify with the team if the expectation is to have them included within the current release. After release workflow is successful the version will be marked as released automatically.

### Make sure that the acceptance tests are successful

While QA acceptance tests are run in the release process automatically, we check [workflows/test-suite.yml](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/test-suite.yml) and see if the latest run of the Test Suite action is successful (it runs every day at midnight UTC time). This can help detect failures before proceeding with the next steps.

### Verify upgrade guide is defined (if required)

- A document (./docs/guides/X.0.0-upgrade-guide.md) must be provided for each major version, summarizing the most significant features, breaking changes, and other helpful information. For minor version releases, this can be created if there are notable changes that warrant it.

- The expectation is that this file is created during relevant pull requests (breaking changes, significant features), and not before the release process.

### Trigger release workflow

- Using our [Release GitHub Action](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/release.yml) run a new workflow using `master` and the following inputs:
  - Version number: `vX.Y.Z`
  - Base branch: Leave empty for standard releases from `master`. For backport releases, see [Backport Releases](#backport-releases) section.
  - Skip QA acceptance tests: Should be left empty. Only used in case failing tests have been encountered in QA and the team agrees the release can still be done, or successful run of QA acceptance tests has already been done with the most recent changes.
  - Using an existing tag: Should be left empty (default `false` creates a new tag from `master`). This should be set to `true` only if you want to re-use an existing tag for the release. This can be helpful for rerunning a failed release process in which the tag has already been created.

#### Backport Releases

Backport releases are used to release fixes or features to previous major versions. The same [Release GitHub Action](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/release.yml) is used for both standard and backport releases.

To create a backport release (e.g., v1.x.x from the `v1-lts` branch):
1. Ensure all changes are merged to the appropriate LTS branch (e.g., `v1-lts` for v1.x releases).
2. Trigger the release workflow with the following inputs:
   - Version number: `v1.Y.Z` (must match the major version of the base branch)
   - Base branch: `v1-lts` (or the appropriate LTS branch name)
   - Other inputs follow the same guidelines as standard releases.

**Important notes for backport releases:**
- The Jira release version step is automatically skipped for backport releases.
- Automatic commits (CHANGELOG.md header update, example links in docs, SSDLC report) are made to the LTS branch, not `master`.
- After a backport release, CHANGELOG.md changes from the LTS branch should be manually merged to `master` branch to keep the main changelog up to date.
- Version validation ensures that version numbers match the base branch (e.g., only `v1.x.x` versions are allowed when releasing from `v1-lts`).

#### How to create a pre-release (not part of regular process)
Pre-releases are not needed for a regular release process, but they can be generated for exceptional cases (e.g. sharing with external teams for additional testing). The process is the same as a regular release except for the format of the version number which must be `vX.Y.Z-pre` with additional numbers at the end if needed. When a pre release is triggered, steps related to updating the changelog header or updating the jira release version are skipped.

### Post-trigger checks
- You will see the release in the [GitHub Release page](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) once the [release action](.github/workflows/release.yml) has completed. HashiCorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the HashiCorp Terraform Registry (more details [here](https://developer.hashicorp.com/terraform/registry/providers/publishing#webhooks)).
- **CDKTF Update - Only for major release, i.e. the left most version digit increment (see this [comment](https://github.com/cdktf/cdktf-repository-manager/pull/202#issuecomment-1602562201))**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

## Post-release considerations for a new major version

When releasing a new major version (e.g., v3.0.0), you need to create an LTS branch for the previous major version to enable backport releases. Follow these steps:

### 1. Create the LTS branch

Create an LTS branch from the last release tag of the previous major version. This can be done at any time after the major release when backport support is needed.

```bash
# Create from the last v2.x.x release tag
git checkout v2.X.X  # Replace with actual last v2 release tag (e.g., v2.23.0)
git checkout -b v2-lts
git push origin v2-lts
git branch -u origin/v2-lts v2-lts
```

### 2. Update GitHub branch protection rules

Add branch protection rules for the new LTS branch (`v2-lts`) matching the rules configured for `master` and existing LTS branches. This is done in the repository settings under "Branches" > "Branch protection rules".

### 3. Update workflow files

Update the following workflow files to include the new LTS branch in their branch selectors:

**`.github/workflows/release.yml`** - Add the new branch to the `base_branch` options:
```yaml
base_branch:
  type: choice
  options:
    - 'master'
    - 'v2-lts'  # Add new LTS branch
    - 'v1-lts'
```

**`.github/workflows/generate-changelog.yml`** - Add to both `branches` and `base_branch` options:
```yaml
branches: [master, v2-lts, v1-lts]  # Add new LTS branch
# ...
base_branch:
  options:
    - 'master'
    - 'v2-lts'  # Add new LTS branch
    - 'v1-lts'
```

**`.github/workflows/code-health.yml`** - Add to branches list:
```yaml
branches:
  - master
  - v2-lts  # Add new LTS branch
  - v1-lts
```

**`.github/workflows/examples.yml`** - Add to branches list:
```yaml
branches:
  - master
  - v2-lts  # Add new LTS branch
  - v1-lts
```

## FAQ

**What happens if a release execution fails to create the tag but generated automatic commits into master?**

All steps before creating the tag are idempotent, meaning you can run the process again and no additional commits will be generated in the second run. This applies to both standard releases and backport releases.

**What happens if a release execution creates a tag but fails during acceptance tests or creating the release (go releaser step)?**

Once a tag has been created in a previous execution, you can make use of the input `Using an existing tag` to run a new release process. 

Depending on the nature of the failure, you may want to introduce new changes into the current release. In this case you must:
- Delete the existing tag.
- Incorporate any new fixes into the target branch (`master` for standard releases, or the appropriate LTS branch like `v1-lts` for backport releases).
- Manually trigger the [Generate Changelog workflow](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/generate-changelog.yml) to remove the current header and including any new entries that have been merged. This will run automatically if you have merged PRs after deleting the tag.
- Trigger a new release process, this will create a tag that includes your latest fixes.

**What happens if all the release process works except the last step to release the version in Jira?**

In this case there is no need to run the full release process again. Once the problem is found and fixed, the [Jira Release Version action](.github/workflows/jira-release-version.yml) can be run manually. Note: Jira release is automatically skipped for backport releases.

**How do I create a backport release?**

Use the same [Release GitHub Action](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/release.yml) as standard releases, but set the `Base branch` input to the appropriate LTS branch (e.g., `v1-lts` for v1.x releases). The workflow will automatically:
- Validate that the version number matches the base branch major version.
- Skip the Jira release version step.
- Perform all automatic commits (CHANGELOG.md header, example links, SSDLC report) on the LTS branch.
- Create the tag from the LTS branch.

After a backport release completes, remember to manually merge the CHANGELOG.md changes from the LTS branch to `master` to keep the main changelog complete.
