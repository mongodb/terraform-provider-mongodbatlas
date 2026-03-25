# IaC Field Dispatcher Plugin — Design Spec

**Date:** 2026-03-25
**Author:** Marco Suma
**Status:** Approved

---

## Overview

The `iac-field-dispatcher` plugin automates the weekly discovery and implementation
of missing fields between the MongoDB Atlas Go SDK and the MongoDB Atlas Terraform
provider. It reads a delta spreadsheet, picks one unimplemented field per run,
creates a Jira Story, and triggers the existing `/add-field` workflow to implement
it end-to-end.

The goal is to steadily close the IaC/SDK gap without overwhelming the team with
simultaneous PRs. One field per week, one PR at a time, reviewed by a human before
merging.

---

## Problem Statement

An automated auditor maintains a spreadsheet tracking the delta between what the
MongoDB Atlas Go SDK exposes and what the Terraform provider surfaces. Rows where
`Used = Yes` and `Unpopulated` is non-empty represent fields that exist in the SDK
and are used by the provider but are not yet exposed as Terraform attributes. These
are the target fields.

Without automation, this delta is invisible week to week and requires manual
triage to turn into actionable work.

---

## Plugin Structure

```
~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/
├── .claude-plugin/
│   └── plugin.json
├── README.md
├── skills/
│   └── dispatch-next-field/
│       ├── SKILL.md
│       └── references/
│           └── spreadsheet-schema.md
└── agents/
    └── spreadsheet-reader.md
```

### Components

| Component | Type | Purpose |
|---|---|---|
| `dispatch-next-field` | Skill | Orchestrates the full weekly dispatch flow |
| `spreadsheet-reader` | Agent | Reads the Google Sheet via Glean, returns structured field candidates |

---

## Workflow

The `dispatch-next-field` skill executes this sequence on each run:

### Step 1 — Read the spreadsheet (via `spreadsheet-reader` agent)

Dispatches the `spreadsheet-reader` agent, which:

1. Calls `mcp__glean_default__read_document` with the spreadsheet URL:
   `https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew`
2. Parses the **"Latest - Terraform"** tab
3. Filters rows where `Used = Yes` AND `Unpopulated` column is non-empty
4. Returns a JSON array of candidates:

```json
[
  {
    "resource": "mongodbatlas_cluster",
    "field": "backupEnabled",
    "sdk_file": "path/to/sdk/file.go",
    "sdk_method": "CreateCluster",
    "input_class": "ClusterDescription"
  }
]
```

**Fallback:** If `read_document` returns nothing, the agent tries
`mcp__glean_default__search` with a descriptive query to locate the document
first, then reads it.

**Failure:** If the sheet is inaccessible via Glean, the agent returns a
structured error. The skill surfaces this to the user and exits without
creating any ticket.

### Step 2 — Query Jira for already-dispatched fields

Runs:
```bash
jira issue list --plain -q "project = CLOUDP AND labels = terraform-sdk-gaps" --paginate "0:50"
```

Paginates through all results (page size 50) until no more results are returned.
Extracts all field+resource combos that already have a ticket. Builds an
exclusion set. Jira is the sole source of truth for "what has already been
dispatched" — no separate state file is used.

### Step 3 — Pick one field

Takes the first candidate from the spreadsheet list that is NOT in the
exclusion set. Candidate order follows the spreadsheet row order.

If all fields are already ticketed, logs a clear message ("No new fields to
dispatch — all candidates already have tickets") and exits without error.

### Step 4 — Create Jira Story

Creates a CLOUDP Story with:

| Field | Value |
|---|---|
| Project | `CLOUDP` |
| Type | `Story` |
| Component | `Terraform` |
| Priority | `Major - P3` |
| Assigned Teams | `APIx DevOps Integrations` |
| Label | `terraform-sdk-gaps` |
| Documentation Changes | `Not Needed` |
| Summary | `Add \`<field>\` to \`mongodbatlas_<resource>\` resource` |
| Description | Structured body referencing SDK method, input class, SDK file path, and the specific unpopulated field(s) (Jira wiki markup) |

The ticket URL returned by the CLI is captured for use in Step 5.

### Step 5 — Invoke `/add-field`

Passes the Jira ticket URL to the existing `add-field` skill:

```
/add-field <ticket-url>
```

From this point, `add-field` owns the complete implementation: branch creation,
SDK research, code changes, unit tests, acceptance tests, and pushing the branch.

---

## Scheduling

Registered as a weekly remote trigger via the `schedule` skill:

- **Schedule:** Every Monday at 9:00 AM (local time)
- **Command:** `/dispatch-next-field`
- **Setup:** Run once during plugin installation (documented in README)

The plugin can also be triggered manually at any time by running
`/dispatch-next-field` in the provider repo directory.

---

## State Tracking

**Jira is the source of truth.** No separate state file exists.

Before picking a field, the skill queries all CLOUDP tickets with label
`terraform-sdk-gaps`. Any field+resource combo that already has a ticket is
excluded from consideration, regardless of that ticket's status.

This means:
- A field with an open PR is skipped (ticket exists)
- A field whose ticket was closed/resolved is also skipped (already done)
- A field whose implementation failed midway is skipped (ticket exists; the
  branch is available for manual inspection)

---

## Edge Cases

| Situation | Behaviour |
|---|---|
| All spreadsheet fields already ticketed | Skill exits cleanly with informational message |
| Glean cannot read the spreadsheet | Skill exits with error, no ticket created |
| Ticket created but `add-field` crashes | Next run skips that field (ticket exists); branch left for manual review |
| Same field appears in multiple spreadsheet rows | Treated as one candidate; deduplication is by field+resource combo |

---

## Jira Ticket Description Template

```
h2. Overview
Add the {{<field_name>}} attribute to the {{mongodbatlas_<resource>}} Terraform resource and data source.

h2. Background
The [IaC/SDK delta tracker|https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew]
identified this field as present in the Atlas Go SDK but not yet exposed in the Terraform provider.

* *SDK method:* {{<sdk_method>}}
* *Input class:* {{<input_class>}}
* *SDK file:* {{<sdk_file>}}

h2. Scope of Work
* Add {{<field_name>}} to the resource schema
* Add {{<field_name>}} to the data source schema
* Implement expander/flattener
* Wire into CRUD functions
* Write unit tests and acceptance tests
* Update docs and changelog

h2. Acceptance Criteria
* {{<field_name>}} is settable via Terraform configuration
* {{<field_name>}} is readable via the data source
* Unit tests pass ({{make test}})
* Acceptance test passes against cloud-dev.mongodb.com
```

---

## Dependencies

| Dependency | Purpose |
|---|---|
| `terraform-provider-new-field` plugin | Provides the `add-field` skill invoked in Step 5 |
| Glean MCP (`mcp__glean_default__*`) | Required for spreadsheet access |
| `jira` CLI | Required for Jira deduplication and ticket creation |
| Atlas API credentials in `~/.zshrc` | Required by `add-field` acceptance tests |

---

## Out of Scope

- Creating multiple tickets per run (intentional: one field per week to avoid PR overload)
- Handling `internal/serviceapi/` (autogenerated) resources — `add-field` already rejects these
- Prioritisation beyond spreadsheet row order — first unprocessed row wins
