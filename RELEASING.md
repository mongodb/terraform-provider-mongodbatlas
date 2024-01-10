# Releasing

## Prerequisites

- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)

## Steps

### Make sure that the acceptance tests are successful
While QA acceptance tests are run in the release process automatically, it is advised to check [workflows/test-suite.yml](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/test-suite.yml) and see if the latest run of the Test Suite action is successful (it runs every day at midnight UTC time). This can help detect failures before proceeding with the next steps.

### Pre-release the provider 
We pre-release the provider to make for testing purpose. **A Pre-release is not published to the Hashicorp Terraform Registry**.

- Using our [Release GitHub Action](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/release.yml) run a new workflow using `master` and the following inputs:
  - Version number: vX.Y.Z-pre
  - Skip QA acceptance tests: Should be left empty. Only used in case failing tests have been encountered in QA and the team agrees the release can still de done, or successful run of QA acceptance tests has already been done with the most recent changes.

- You will see the release in the [GitHub Release page](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) once the [release action](.github/workflows/release.yml) has completed.

**Note**: If a failure is encountered during the go releaser step you must manually delete the created tag and then retry running the action.

### Create PR updating Changelog and Upgrade Guide

- Create a JIRA ticket and open a PR against the **master** branch. Make any manual adjustments if needed taking into account date format and format parameter names and resources/data source names if they begin with `mongodbatlas`.
- Include the PM as a PR reviewer
- Contact Documentation team in Slack to review the PR
- Example: [#1478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1478).

The PR includes the following changes:

#### Generate the CHANGELOG.md 
We use a tool called [github changelog generator](https://github.com/github-changelog-generator/github-changelog-generator) to automatically update our changelog. It provides options for downloading a CLI or using a docker image with interactive mode to update the CHANGELOG.md file locally.

- Update `since_tag` and `future-release` in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator)
- **There is a bug with `github_changelog_generator` ([#971](https://github.com/github-changelog-generator/github-changelog-generator/issues/971))**: Make sure to update the `future-tag` with the pre-release tag. Once you generate the changelog, update `future-tag` with the final release tag in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator). Then, manually update the generated changelog to remove references to the pre-release tag
- Run the following command: 
    ```bash 
    docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u mongodb -p terraform-provider-mongodbatlas -t <GH_TOKEN> --breaking-labels "breaking-change" --enhancement-label "**Enhancements**" --bugs-label "**Bug Fixes**"  --issues-label "**Closed Issues**" --pr-label "**Internal Improvements**" --max-issues 1000
    ```
    To obtain your github personal access token you can use the following guide: [Authorizing a personal access token for use with SAML single sign-on](https://docs.github.com/en/enterprise-cloud@latest/authentication/authenticating-with-saml-single-sign-on/authorizing-a-personal-access-token-for-use-with-saml-single-sign-on)
 
#### Define the Upgrade Guide

**Note**: Only applies if the right most version digit is 0 (considered a major or minor version in [semantic versioning](https://semver.org/)).

- Create a new doc in /website/docs/guides/X.Y.0-upgrade-guide.html. This will contain a summary of the most significant features and breaking changes. Additional information that can be helpful to users can be defined here.

### Release the provider
- Follow the same steps in the pre-release but provide the final release tag (example `v1.9.0`). Harshicorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the Hashicorp Terraform Registry.
- **CDKTF Update - Only for major release, i.e. the left most version digit increment (see this [comment](https://github.com/cdktf/cdktf-repository-manager/pull/202#issuecomment-1602562201))**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

