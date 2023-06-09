# Releasing

## Prerequisites

- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)

# Generating the CHANGELOG.md
Run: 
```bash 
github_changelog_generator -u mongodb -p terraform-provider-mongodbatlas --enhancement-label "**Enhancements**" --bugs-label "**Bug Fixes**"
```
