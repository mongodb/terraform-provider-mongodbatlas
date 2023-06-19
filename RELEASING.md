# Releasing

## Prerequisites

- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)

## Steps
### Generate the CHANGELOG.md 
- Run

    ```bash 
    github_changelog_generator -u mongodb -p terraform-provider-mongodbatlas --enhancement-label "**Enhancements**" --bugs-label "**Bug Fixes**"
    ```
-  Open a PR against the **master** branch

### Pre-release the provider 
We pre-release the provider to make for testing purpose. **A Pre-release is not published to the Hashicorp Terraform Registry**.

- Open the GitHub repository release page and click draft a new release
- Fill the pre-release tag and select the correct target branch

    <img width="370" alt="image2" src="https://github.com/mongodb/terraform-provider-mongodbatlas/assets/5663078/e710c0ff-dc00-44c2-9eb6-146cd791d47e">
- Select the **master** branch
- Generate Release Notes: Click Generate release notes button to populate release notes
- Set publishing to Pre-release
    
    <img width="477" alt="image3" src="https://github.com/mongodb/terraform-provider-mongodbatlas/assets/5663078/30d2db83-6b2d-4eb2-9da6-93fc34d64c09">

- **There is a bug in the GitHub release page**: after binaries get created, GitHub  flips backthe  status of release as Draft so you have to set it to Pre-Release again.

### Release the provider
- Follow the same steps in the pre-release but provide the final release tag (example `v1.9.0`). This will trigger the release action that will release the provider to the GitHub Release page. Harshicorp has a process in place that will retrieve the latest release from the GitHub repository and add the binaries to the Hashicorp Terraform Registry.
- **CDKTF Update**: Once the provider has been released, we need to update the provider version in our CDKTF. Raise a PR against [cdktf/cdktf-repository-manager](https://github.com/cdktf/cdktf-repository-manager).
  - Example PR: [#183](https://github.com/cdktf/cdktf-repository-manager/pull/183)

