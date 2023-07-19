# Releasing

## Prerequisites

- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)

## Steps

### Pre-release the provider 
We pre-release the provider to make for testing purpose. **A Pre-release is not published to the Hashicorp Terraform Registry**.

- Open the GitHub repository release page and click draft a new release
- Fill the pre-release tag and select `master` as the target branch

    <img width="370" alt="image2" src="https://github.com/mongodb/terraform-provider-mongodbatlas/assets/5663078/e710c0ff-dc00-44c2-9eb6-146cd791d47e">
- Generate Release Notes: Click Generate release notes button to populate release notes
- Set publishing to Pre-release
    
    <img width="477" alt="image3" src="https://github.com/mongodb/terraform-provider-mongodbatlas/assets/5663078/30d2db83-6b2d-4eb2-9da6-93fc34d64c09">

- **There is a bug in the GitHub release page**: after binaries get created, GitHub  flips backthe  status of release as Draft so you have to set it to Pre-Release (or Latest, if publishing the final version) again.

### Generate the CHANGELOG.md 
- Update `since_tag` and `future-release` in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator)
- **There is a bug with `github_changelog_generator` ([#971](https://github.com/github-changelog-generator/github-changelog-generator/issues/971))**: Make sure to update the `future-tag` with the pre-release tag. Once you generate the changelog, update `future-tag` with the final release tag in [.github_changelog_generator](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/.github_changelog_generator). Then, manually update the generated changelog to remove references to the pre-release tag
- Run

    ```bash 
    github_changelog_generator -u mongodb -p terraform-provider-mongodbatlas --enhancement-label "**Enhancements**" --bugs-label "**Bug Fixes**"  --issues-label "**Closed Issues**" --pr-label "**Internal Improvements**"
    ```
-  Open a PR against the **master** branch
-  Example: [#1308](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1308)

### Release the provider
- Follow the same steps in the pre-release but provide the final release tag (example `v1.9.0`). This will trigger the release action that will release the provider to the GitHub Release page. Harshicorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the Hashicorp Terraform Registry.
- **CDKTF Update - Only for major release, i.e. the left most version digit increment (see this [comment](https://github.com/cdktf/cdktf-repository-manager/pull/202#issuecomment-1602562201))**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

