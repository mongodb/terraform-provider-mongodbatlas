# Releasing

## Prerequisites

- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)

## Steps

### Make sure that the acceptance tests are successful
Check [workflows/acceptance-tests.yml](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/acceptance-tests.yml) and see if the latest run of the Acceptance Test action is successful (it runs every day at 4 AM Dublin Time). If tests are failing, you should investigate the failure before proceeding with the next steps.

### Pre-release the provider 
We pre-release the provider to make for testing purpose. **A Pre-release is not published to the Hashicorp Terraform Registry**.

- Create and push the pre-release tag (`X.Y.Z-pre`) to master
```bash
git tag [YOUR_TAG]-pre
git push origin [YOUR_TAG]-pre
```

- You will see the release in the [GitHub Release page](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) once the [release action](.github/workflows/release.yml) has completed.

### Generate the CHANGELOG.md 
We use a tool called [github changelog generator](https://github.com/github-changelog-generator/github-changelog-generator) to automatically update our changelog. It provides options for downloading a CLI or using a docker image with interactive mode to update the CHANGELOG.md file locally.

- Update `since_tag` and `future-release` in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator)
- **There is a bug with `github_changelog_generator` ([#971](https://github.com/github-changelog-generator/github-changelog-generator/issues/971))**: Make sure to update the `future-tag` with the pre-release tag. Once you generate the changelog, update `future-tag` with the final release tag in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator). Then, manually update the generated changelog to remove references to the pre-release tag
- Run the following command: 
    ```bash 
    docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u mongodb -p terraform-provider-mongodbatlas -t <GH_TOKEN> --breaking-labels "breaking-change" --enhancement-label "**Enhancements**" --bugs-label "**Bug Fixes**"  --issues-label "**Closed Issues**" --pr-label "**Internal Improvements**"
    ```
    To obtain your github personal access token you can use the following guide: [Authorizing a personal access token for use with SAML single sign-on](https://docs.github.com/en/enterprise-cloud@latest/authentication/authenticating-with-saml-single-sign-on/authorizing-a-personal-access-token-for-use-with-saml-single-sign-on)
- Create a JIRA ticket and open a PR against the **master** branch. Make any manual adjustments if needed taking into account date format and format parameter names and resources/data source names if they begin with `mongodbatlas`.
- Include the PM as a PR reviewer
- Contact Documentation team in Slack to review the PR
- Example: [#1478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1478). 
- If the right most version digit is 0 then create a new doc in /website/docs/guides/X.Y.0-upgrade-guide.html

### Release the provider
- Follow the same steps in the pre-release but provide the final release tag (example `v1.9.0`). This will trigger the release action that will release the provider to the GitHub Release page. Harshicorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the Hashicorp Terraform Registry.
- **CDKTF Update - Only for major release, i.e. the left most version digit increment (see this [comment](https://github.com/cdktf/cdktf-repository-manager/pull/202#issuecomment-1602562201))**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

