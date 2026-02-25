---
name: summarize-test-suite
description: Summarize GitHub Actions test suite execution results from a workflow run URL. Fetches job logs, identifies failing tests, and categorizes root causes. Use when the user shares a GitHub Actions run URL and asks about test failures, flaky tests, or CI results.
---

# Summarize Test Suite Execution

## When to Use

The user provides a GitHub Actions workflow run or job URL, e.g.:
- `https://github.com/mongodb/terraform-provider-mongodbatlas/actions/runs/<run_id>`
- `https://github.com/mongodb/terraform-provider-mongodbatlas/actions/runs/<run_id>/job/<job_id>`

## Workflow

### Step 1: Identify Failed Jobs

Use `gh api` with the `--paginate` flag and `"all"` permissions (needed for network access).

```bash
gh api repos/{owner}/{repo}/actions/runs/{run_id}/jobs --paginate
```

Parse the JSON to list jobs with `conclusion: "failure"` and their failed steps. Record each failed job's `id`.

### Step 2: Fetch Logs for Each Failed Job

For each failed job ID, fetch and filter logs in parallel:

```bash
gh api repos/{owner}/{repo}/actions/jobs/{job_id}/logs 2>&1 | grep -i -E "(FAIL|panic|Error|---)" | tail -60
```

Then extract the test-level results to identify which tests actually failed vs passed:

```bash
gh api repos/{owner}/{repo}/actions/jobs/{job_id}/logs 2>&1 | grep -E "(=== RUN|--- PASS|--- FAIL|FAIL\t)" | grep -v "    ---"
```

**Important**: Many `[ERROR]` log lines come from tests that intentionally exercise error paths and ultimately PASS. Always cross-reference error logs against the `--- FAIL` / `--- PASS` verdicts.

### Step 3: Get Error Details for Each Failing Test

For each `--- FAIL: TestName` found, fetch surrounding context:

```bash
gh api repos/{owner}/{repo}/actions/jobs/{job_id}/logs 2>&1 | grep -B10 "FAIL: TestName"
```

### Step 4: Categorize and Summarize

Group failures into root cause categories. Common categories for this repo:

| Category | Indicators |
|----------|-----------|
| **Cloud capacity** | `OUT_OF_CAPACITY`, `NO_CAPACITY`, `No Capacity` |
| **API errors** | `HTTP 400`, `HTTP 409`, `HTTP 500` with specific error codes |
| **Timeout / flake** | `timeout while waiting for state`, eventual consistency assertions |
| **Cleanup race** | `still exists` after destroy, `CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS` |
| **Code regression** | Assertion failures on attribute values, unexpected plan diffs |

### Step 5: Produce Output

Provide two versions:

1. **Detailed summary** — table with test names, error messages, and root cause categories.
2. **Slack-ready summary** — concise format the user can paste directly, structured as:

```
**Test Suite #N — Failure Summary**

**X jobs failed, Y tests total across Z packages**

**Category 1 (N tests)** — brief explanation:
- `TestName1` — short error
- `TestName2` — short error

**Category 2 (N tests):**
- `TestName3` — short error

**TL;DR**: One-sentence overall assessment.
```

## Key Learnings

- Always request `"all"` permissions for `gh api` calls — TLS certificate verification fails in the default sandbox.
- Log output is large; always filter with `grep` rather than reading raw logs.
- The `--- PASS` / `--- FAIL` verdict is the source of truth — many ERROR-level log lines are from tests that pass (they test error handling paths).
- Package-level `FAIL` without a corresponding `--- FAIL` test line can indicate cleanup failures or build errors rather than individual test failures.
