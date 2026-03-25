# IaC Field Dispatcher Plugin — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `iac-field-dispatcher` Claude Code plugin that weekly discovers one missing Terraform provider field from the Atlas SDK delta spreadsheet, creates a Jira Story, and invokes `/add-field` to implement it end-to-end.

**Architecture:** A skill (`dispatch-next-field`) orchestrates the flow: it dispatches a dedicated `spreadsheet-reader` agent to parse the Google Sheet via Glean, deduplicates against existing Jira tickets labelled `terraform-sdk-gaps`, creates one new CLOUDP Story, and hands off to the existing `/add-field` skill. Jira is the sole source of truth for dispatch state — no separate state file.

**Tech Stack:** Claude Code plugin (Markdown skills/agents), Glean MCP (`mcp__glean_default__*`), `jira` CLI (jira-cli v1.7.0+), existing `terraform-provider-new-field` plugin (`add-field` skill)

---

## File Map

| File | Action | Purpose |
|---|---|---|
| `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/.claude-plugin/plugin.json` | Create | Plugin manifest |
| `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/README.md` | Create | Plugin docs and setup instructions |
| `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/agents/spreadsheet-reader.md` | Create | Glean-based agent: reads sheet, returns JSON candidates |
| `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/SKILL.md` | Create | Orchestrator skill: dedup, pick, ticket, invoke add-field |
| `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/references/spreadsheet-schema.md` | Create | Documents Google Sheet column layout for the agent |

---

### Task 1: Plugin Scaffold

**Files:**
- Create: `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/.claude-plugin/plugin.json`

- [ ] **Step 1: Create the directory structure**

```bash
mkdir -p ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/.claude-plugin
mkdir -p ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/agents
mkdir -p ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/references
```

- [ ] **Step 2: Create plugin.json**

Write `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/.claude-plugin/plugin.json`:

```json
{
  "name": "iac-field-dispatcher",
  "version": "0.1.0",
  "description": "Weekly automation for discovering and implementing missing Terraform provider fields from the Atlas SDK delta spreadsheet",
  "author": {
    "name": "Marco Suma"
  }
}
```

- [ ] **Step 3: Verify plugin directory exists**

```bash
ls ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/
```

Expected: `.claude-plugin/`, `agents/`, `skills/` directories visible.

- [ ] **Step 4: Commit**

```bash
cd ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher
git init   # only if not already inside a git repo
git add .claude-plugin/plugin.json
git commit -m "feat: scaffold iac-field-dispatcher plugin"
```

---

### Task 2: Spreadsheet Schema Reference

**Files:**
- Create: `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/references/spreadsheet-schema.md`

This documents the Google Sheet column layout so the `spreadsheet-reader` agent
doesn't have to guess the structure.

- [ ] **Step 1: Create spreadsheet-schema.md**

Write `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/references/spreadsheet-schema.md`:

```markdown
# Spreadsheet Schema: IaC/SDK Delta Tracker

**URL:** https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew
**Tab:** Latest - Terraform

## Column Layout

| Column | Name | Description |
|--------|------|-------------|
| 1 | Method | Atlas Go SDK method name (e.g. `CreateCluster`) |
| 2 | Input Class | SDK input struct name (e.g. `ClusterDescription`) |
| 3 | SDK File | Path to the SDK file where this is defined |
| 4 | Used | `Yes` if the SDK method is used by the Terraform provider |
| 5 | Unpopulated | Comma-separated list of field names not yet exposed in the provider |

## Filtering Rules

Return only rows where ALL of the following are true:
1. `Used` column = `Yes` (case-insensitive)
2. `Unpopulated` column is non-empty

## Output Format

Split multi-field `Unpopulated` values into separate entries — one entry per
missing field:

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

**Deriving `resource`:** Convert the SDK method/class name to a Terraform
resource name using snake_case with `mongodbatlas_` prefix. If the mapping is
ambiguous, include the raw method name in the output; the dispatch skill will
resolve it when invoking `add-field`.
```

- [ ] **Step 2: Commit**

```bash
git add skills/dispatch-next-field/references/spreadsheet-schema.md
git commit -m "docs: add spreadsheet schema reference for iac-field-dispatcher"
```

---

### Task 3: Spreadsheet Reader Agent

**Files:**
- Create: `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/agents/spreadsheet-reader.md`

This agent is invoked by the dispatch skill. It returns a JSON array of candidate
fields. All Glean interaction is contained here.

- [ ] **Step 1: Create spreadsheet-reader.md**

Write `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/agents/spreadsheet-reader.md`:

```markdown
---
name: spreadsheet-reader
description: Use this agent to read the IaC/SDK delta spreadsheet and return a structured list of missing Terraform provider fields. Invoke when you need to find candidate fields for implementation.
color: blue
---

You are a spreadsheet reader agent. Your sole job is to read the IaC/SDK delta
tracking spreadsheet via Glean and return a structured JSON list of candidate
fields for the Terraform provider.

## Spreadsheet

URL: https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew
Tab: Latest - Terraform

## Instructions

1. Call `mcp__glean_default__read_document` with the spreadsheet URL above.

2. If the document is returned, parse it for the "Latest - Terraform" tab.
   Column layout and filtering rules are in
   `skills/dispatch-next-field/references/spreadsheet-schema.md`.

3. If `read_document` returns nothing or an error, try a fallback search:
   - Call `mcp__glean_default__search` with query:
     `"IaC SDK delta tracker Terraform Latest Unpopulated"`
   - Take the first spreadsheet URL from the results and call `read_document` on it.

4. If the spreadsheet is still inaccessible after the fallback, return:
   ```json
   {"error": "Spreadsheet not accessible via Glean. URL: https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew"}
   ```

## Output

Return ONLY a JSON array or error object. No preamble, no explanation.

- If candidates found: JSON array per `spreadsheet-schema.md` output format
- If no rows match the filter: `[]`
- If spreadsheet inaccessible: `{"error": "..."}`
```

- [ ] **Step 2: Verify file has valid frontmatter**

```bash
head -6 ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/agents/spreadsheet-reader.md
```

Expected: `---`, `name: spreadsheet-reader`, `description: ...`, `color: blue`, `---`

- [ ] **Step 3: Commit**

```bash
git add agents/spreadsheet-reader.md
git commit -m "feat: add spreadsheet-reader agent for iac-field-dispatcher"
```

---

### Task 4: Dispatch Skill

**Files:**
- Create: `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/SKILL.md`

The orchestrator. Reads candidates, deduplicates against Jira, picks one field,
creates the Story, hands off to `add-field`.

- [ ] **Step 1: Create SKILL.md**

Write `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/SKILL.md`:

```markdown
# Dispatch Next Field

Discovers one unimplemented Terraform provider field from the IaC/SDK delta
spreadsheet and dispatches the full implementation workflow.

Triggers when the user says: "dispatch next field", "run iac-field-dispatcher",
"pick next SDK gap", `/dispatch-next-field`, or when invoked by the weekly schedule.

## Step 1: Read the Spreadsheet

Invoke the `spreadsheet-reader` agent. It returns a JSON array of candidate
fields filtered from the "Latest - Terraform" tab.

**On error:** If the agent returns `{"error": "..."}`, stop and surface the
error to the user. Do not create any ticket.

**On empty:** If the agent returns `[]`, stop and report:
> "No candidate fields found — all `Used = Yes` rows have an empty `Unpopulated`
> column."

## Step 2: Build the Exclusion Set from Jira

Query all CLOUDP tickets labelled `terraform-sdk-gaps`, paginating until all
results are retrieved (stop when a page returns fewer than 50 results):

```bash
# Repeat with offset 0, 50, 100, ... until page returns < 50 rows
jira issue list --plain --no-headers \
  -q "project = CLOUDP AND labels = terraform-sdk-gaps" \
  --paginate "<offset>:50" \
  --columns KEY,SUMMARY
```

From the SUMMARY column, extract field and resource using this pattern:
`Add \`<field>\` to \`mongodbatlas_<resource>\` resource`

Build an exclusion set of `<field>@<resource>` strings.

## Step 3: Pick One Field

Iterate through the spreadsheet candidates in order (top row first). Pick the
FIRST candidate whose `<field>@<resource>` is NOT in the exclusion set.

If all candidates are in the exclusion set, stop and report:
> "Nothing new to dispatch — all candidates already have Jira tickets with
> label `terraform-sdk-gaps`."

## Step 4: Create the Jira Story

Create a CLOUDP Story via `jira` CLI. Pipe the description via stdin (never use
`-b "$(cat ..."` — it hangs):

```bash
cat <<'EOF' | jira issue create \
  -p CLOUDP \
  -t Story \
  -s "Add `<field>` to `mongodbatlas_<resource>` resource" \
  -C Terraform \
  -y "Major - P3" \
  --custom "assigned-teams=APIx DevOps Integrations" \
  --custom "documentation-changes=Not Needed" \
  --no-input \
  --raw
h2. Overview
Add the {{<field>}} attribute to the {{mongodbatlas_<resource>}} Terraform resource and data source.

h2. Background
The [IaC/SDK delta tracker|https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew]
identified this field as present in the Atlas Go SDK but not yet exposed in the Terraform provider.

* *SDK method:* {{<sdk_method>}}
* *Input class:* {{<input_class>}}
* *SDK file:* {{<sdk_file>}}

h2. Scope of Work
* Add {{<field>}} to the resource schema
* Add {{<field>}} to the data source schema (if a data source exists)
* Implement expander/flattener
* Wire into CRUD functions
* Write unit tests and acceptance tests
* Update docs and changelog

h2. Acceptance Criteria
* {{<field>}} is settable via Terraform configuration
* {{<field>}} is readable via the data source (if applicable)
* Unit tests pass ({{make test}})
* Acceptance test passes against cloud-dev.mongodb.com
EOF
```

After creation, add the `terraform-sdk-gaps` label:

```bash
jira issue edit <KEY> --label terraform-sdk-gaps --no-input
```

Construct the ticket URL: `https://jira.mongodb.org/browse/<KEY>`

Report to the user: "Created ticket <KEY>: <URL>"

## Step 5: Invoke add-field

Pass the Jira ticket URL to the `add-field` skill:

```
add field https://jira.mongodb.org/browse/<KEY>
```

The `add-field` skill owns the rest: branch creation, OAS research, code
changes, unit tests, acceptance tests, and push.
```

- [ ] **Step 2: Verify file structure**

```bash
grep "^## Step" ~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/skills/dispatch-next-field/SKILL.md
```

Expected: Steps 1 through 5 listed.

- [ ] **Step 3: Commit**

```bash
git add skills/dispatch-next-field/SKILL.md
git commit -m "feat: add dispatch-next-field skill for iac-field-dispatcher"
```

---

### Task 5: README

**Files:**
- Create: `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/README.md`

- [ ] **Step 1: Create README.md**

Write `~/.claude-work/plugins/marketplaces/local/plugins/iac-field-dispatcher/README.md`:

```markdown
# iac-field-dispatcher

Weekly automation for discovering and implementing missing Terraform provider
fields from the Atlas SDK/IaC delta spreadsheet.

## What It Does

Each run:
1. Reads the [IaC/SDK delta spreadsheet](https://docs.google.com/spreadsheets/d/1Rc2Xf9MdBkdWm3s4LdLUJsiJlyqTI8FvcXCBwh1M4Ew) ("Latest - Terraform" tab) via Glean
2. Queries CLOUDP Jira for tickets labelled `terraform-sdk-gaps`
3. Picks the first spreadsheet field not yet ticketed
4. Creates a CLOUDP Story (team: APIx DevOps Integrations, component: Terraform, label: terraform-sdk-gaps)
5. Invokes `add-field` to implement the field end-to-end (branch → code → tests → push)

One field per run. One PR per week. Human review before merge.

## Prerequisites

- Glean MCP configured and the spreadsheet indexed in Glean
- `jira` CLI configured (`~/.config/.jira/.config.yml`)
- `terraform-provider-new-field` plugin installed
- Atlas API credentials in `~/.zshrc`:
  ```bash
  export MONGODB_ATLAS_PUBLIC_KEY="..."
  export MONGODB_ATLAS_PRIVATE_KEY="..."
  export MONGODB_ATLAS_ORG_ID="..."
  export MONGODB_ATLAS_BASE_URL="https://cloud-dev.mongodb.com"
  export TF_ACC=1
  ```

## Usage

### Manual run

From inside the `terraform-provider-mongodbatlas` repo directory:

```
/dispatch-next-field
```

Or naturally: "dispatch next field" / "pick next SDK gap"

### Weekly schedule

Set up once (run this in any Claude Code session):

```
/schedule "every Monday at 9am, run /dispatch-next-field in ~/Dev/terraform-provider-mongodbatlas"
```

## State Tracking

Jira is the source of truth. A field is "dispatched" if any CLOUDP ticket with
label `terraform-sdk-gaps` has a summary matching
`Add \`<field>\` to \`mongodbatlas_<resource>\` resource`.

Closed/resolved tickets count as dispatched — those fields are done.

## Edge Cases

| Situation | Behaviour |
|---|---|
| All spreadsheet fields already ticketed | Exits cleanly with informational message |
| Glean cannot read the spreadsheet | Exits with error, no ticket created |
| Ticket created but `add-field` crashes | Next run skips that field (ticket exists); branch available for manual inspection |
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README for iac-field-dispatcher plugin"
```

---

### Task 6: Register the Weekly Schedule

- [ ] **Step 1: Open a Claude Code session in the provider repo**

```bash
cd ~/Dev/terraform-provider-mongodbatlas
```

- [ ] **Step 2: Register the schedule via the schedule skill**

Run in Claude Code:

```
/schedule "every Monday at 9am, run /dispatch-next-field in ~/Dev/terraform-provider-mongodbatlas"
```

- [ ] **Step 3: Verify the schedule was registered**

```
/schedule list
```

Expected: a Monday 9am entry targeting `dispatch-next-field` visible in the list.

---

### Task 7: Smoke Test

Verify the full flow works before relying on the schedule.

- [ ] **Step 1: Verify the spreadsheet-reader agent works**

In a Claude Code session in the provider repo:

```
Use the spreadsheet-reader agent to read the IaC delta spreadsheet and show me the first 3 candidates as JSON.
```

Expected: a JSON array with entries containing `field`, `resource`, `sdk_method`, `input_class`, `sdk_file`.

- [ ] **Step 2: Verify Jira deduplication query works**

```bash
jira issue list --plain --no-headers \
  -q "project = CLOUDP AND labels = terraform-sdk-gaps" \
  --paginate "0:50" --columns KEY,SUMMARY
```

Expected: command exits without error (empty list is fine if no tickets exist yet).

- [ ] **Step 3: Preview what would be picked (without creating a ticket)**

```
What would /dispatch-next-field pick next, without creating a ticket? Just tell me the field and resource.
```

Review the candidate. Confirm it looks like a real missing field.

- [ ] **Step 4: Full run**

```
/dispatch-next-field
```

Confirm:
- Jira ticket URL is printed
- `add-field` is invoked with that URL
- Branch is created in the provider repo
- Implementation proceeds per the `add-field` workflow
