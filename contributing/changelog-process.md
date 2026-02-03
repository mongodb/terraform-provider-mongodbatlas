# Changelog Process

- [Script for creating changelog entry files](#script-for-creating-changelog-entry-files)
- [Changelog format](#changelog-format)
- [Changelog entry guidelines](#changelog-entry-guidelines)
- [PR Changelog check](#pr-changelog-check)
- [Unreleased section of CHANGELOG.md automatic update](#unreleased-section-of-CHANGELOG.md-automatic-update)

HashiCorp’s open-source projects have always maintained user-friendly, readable CHANGELOG.md that allow users to tell at a glance whether a release should have any effect on them, and to gauge the risk of an upgrade.

We use the [go-changelog](https://github.com/hashicorp/go-changelog) to generate and update the changelog automatically from files created in the `.changelog/` directory. It is important that when you raise your Pull Request, there is a changelog entry which describes the changes your contribution makes. Not all changes require an entry in the changelog, guidance follows on what changes do.

`@mongodb/docs-cloud-team` will be required reviewers for new changelog entry files contained in a Pull Request. We will wait up to 24 hours for a review, and after that proceed with the merge. Exceptions for merging ahead of the 24 hours may also apply.

## Script for creating changelog entry files

A script is defined to guide the creation of new entry files, simplifying the process and avoiding errors. You can invoke the script using the following make command:

```
make generate-changelog-entry
```

- The `subcategory` input prompt refers to the prefix of the changelog entry, used for specifying the relevant resource/data source when needed (e.g. data-source/mongodbatlas_project)

## Changelog format

The changelog format requires an entry in the following format, where HEADER corresponds to the changelog category, and the entry is the changelog entry itself. The entry should be included in a file in the `.changelog` directory with the naming convention `{PR-NUMBER}.txt`. For example, to create a changelog entry for pull request 1234, there should be a file named `.changelog/1234.txt`.

``````markdown
```release-note:{HEADER}
{ENTRY}
```
``````

If a pull request should contain multiple changelog entries, then multiple blocks can be added to the same changelog file. For example:

``````markdown
```release-note:note
resource/mongodbatlas_project: Deprecates `labels` attribute. All configurations using `labels` should be updated to use the new `tags` attribute instead.
```

```release-note:enhancement
resource/mongodbatlas_project: Adds `tags` attribute
```
``````


## Changelog entry guidelines

The CHANGELOG is intended to show operator-impacting changes to the codebase for a particular version. If every change or commit to the code resulted in an entry, the CHANGELOG would become less useful for operators. The lists below are general guidelines and examples for when a decision needs to be made to decide whether a change should have an entry.

### Header and entry values

``````markdown
```release-note:{HEADER}
{ENTRY}
```
``````

 `HEADER`
 - Must be one of the following values: `breaking-change`, `new-resource`, `new-datasource`, `new-guide`, `note`, `enhancement`, `bug`. Examples for each type can be seen below.

`ENTRY`
 - In the case of feature entry types (new-resource, new-datasource, new-guide) only the name of the new resource or guide is defined in the entry.
 - For other entry types:
    - Entry starts with the resource type, followed by its name (e.g. `resource/mongodbatlas_project: `). Use a `provider: ` prefix for provider-level changes.
    - For the description use a third person point of view, [active voice](https://www.mongodb.com/docs/meta/style-guide/writing/use-active-voice/#std-label-use-active-voice), and start with an uppercase character.
    - Surround attribute names with backticks (e.g. ```Adds `tags` attribute```). 

### Changes that should have a CHANGELOG entry

#### New resource

A new resource entry should only contain the name of the resource, and use the `release-note:new-resource` header.

``````markdown
```release-note:new-resource
mongodbatlas_stream_connection
```
``````

#### New data source

A new data source entry should only contain the name of the data source, and use the `release-note:new-datasource` header.

``````markdown
```release-note:new-datasource
mongodbatlas_stream_connection
```
``````

#### New full-length documentation guides (e.g., Upgrade Guide for a new major version)

A new full length documentation entry gives the title of the documentation added, using the `release-note:new-guide` header.

``````markdown
```release-note:new-guide
MongoDB Atlas Provider 1.15.0: Upgrade and Information Guide
```
``````

#### Resource and provider bug fixes

A new bug entry should use the `release-note:bug` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level fixes.

``````markdown
```release-note:bug
resource/mongodbatlas_database_user: Avoids sending database user password in update request if the value has not changed
```
``````

#### Resource and provider enhancements

A new enhancement entry should use the `release-note:enhancement` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level enhancements.

``````markdown
```release-note:enhancement
data-source/mongodbatlas_project: Adds `tags` attribute
```
``````

#### Deprecations

A deprecation entry should use the `release-note:note` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level changes.

``````markdown
```release-note:note
resource/mongodbatlas_project: Deprecates the `labels` attribute. All configurations using `labels` should be updated to use the new `tags` attribute instead.
```
``````

#### Breaking changes and removals

A breaking-change entry should use the `release-note:breaking-change` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level changes.

``````markdown
```release-note:breaking-change
data-source/mongodbatlas_search_indexes: Removes `page_num` and `items_per_page` attributes
```
``````

### Changes that should _not_ have a CHANGELOG entry

- Resource and provider documentation updates
- Testing updates
- Code refactoring
- Dependency updates

## PR Changelog check

A PR check is included to validate the changelog entry file. 
If a PR doesn't need a changelog entry its check can be skipped:
- Adding the label `skip-changelog-check` to the PR.
- Check in PRs with title `chore`, `test`, `doc` or `ci` is automatically skipped. However a changelog can still be added if needed.

## Unreleased section of CHANGELOG.md automatic update

After a PR is merged to `master` (or LTS branches like `v1-lts`) with a new entry in `.changelog` directory, [Generate Changelog workflow](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/generate-changelog.yml) will be triggered and it will update the `CHANGELOG.md` file with the new entry.
This workflow can also be triggered manually and it will update the `CHANGELOG.md` file with all entries in `.changelog` directory that are not present in the `CHANGELOG.md` file.

**Note:** After a backport release (e.g., from `v1-lts`), CHANGELOG.md changes from the LTS branch should be manually merged to `master` branch.
