---
name: monthly-test-suite-summary
description: Build the monthly test-stability summary for the OpEx document (Test stability → Overall metrics section) for the mongodb/terraform-provider-mongodbatlas repo. Aggregates the daily summarize-test-suite verdicts of one calendar month into a paste-ready markdown table plus a short narrative. Use when asked for the monthly OpEx test summary, monthly test stability metrics, or to update the Overall metrics section.
---

# Monthly Test Suite Summary (OpEx)

## Purpose

Produce the monthly replacement for the OpEx doc's *Test stability → Overall metrics* table. The old metric counted raw run conclusions (a run with only `OUT_OF_CAPACITY` flakes counted as "Fail"), which made the suite look permanently broken. This skill counts a run as a regression only when its daily summary verdict was `:red_circle:` — i.e., the daily `summarize-test-suite` classification found at least one code regression.

## Determinism contract

All data collection and aggregation is done by `scripts/monthly_summary.py`. The script's numbers are authoritative — **never alter, recompute, or "correct" them**. The agent's only reasoning work is Step 2 (sanity checks and narrative).

Per-run verdicts come from the daily `summarize-test-suite` skill's output (the `summary` artifact uploaded by every scheduled Test Suite run). This skill never re-classifies failures and never reads raw job logs.

## Inputs

Optional: an anchor date `YYYY-MM-DD`. "Last month" is the calendar month before the anchor's month. Default anchor is today; pass `--date` when the skill was not run on time (e.g. `--date 2026-08-01` on August 20 still yields July).

## Workflow

### Step 1: Run the aggregation script

```bash
python3 .agents/skills/monthly-test-suite-summary/scripts/monthly_summary.py [--date YYYY-MM-DD]
```

Requires `gh` authenticated against `mongodb/terraform-provider-mongodbatlas`. The script writes two artefacts (paths printed on stderr):

- `monthly-summary-<anchor-date>.md` in the current working directory — the metrics table plus exclusion footnote. This is the final output file.
- A JSON report (default `$TMPDIR/monthly-test-suite-summary/<month>-result.json`). Read it for the narrative step.

### Step 2: Sanity checks and narrative (the only agent reasoning)

1. **Sanity checks.** Flag in the narrative, with likely cause, when:
   - The month has far fewer runs than days (scheduled runs may have been skipped or the workflow renamed).
   - `no_summary_runs` is non-empty — the summary job failed or the artifact expired. These runs are excluded from the percentage; say so.
   - `unknown_verdict_runs` is non-empty — a summary artifact existed but no verdict emoji was found; suggest a manual look.
2. **Repeat offenders.** `recurring_tests` lists tests named in daily failing lists on ≥2 distinct days (a lower bound — daily summaries cap the list at 10 names plus `, and N more`). For each of the top recurring tests, read the cached daily summaries (`$TMPDIR/monthly-test-suite-summary/<run_id>/summary.md`) for those dates and determine whether it failed with the **same root cause every time** (persistent broken or flaky test — a ticket candidate; say which category, e.g. always `OUT_OF_CAPACITY` vs always a timeout) or **alternating causes** (general infra noise). Keep this to the handful of flagged tests.
3. **Red days.** For each run with verdict `red` in the JSON report, give one line with the run number, date, linked run URL, and the regression cause from its `regressions` entries. Group repeats (e.g. "3 of 4 red days were the same `config_server_type` regression").
4. **Category trends.** Summarize `category_totals` in one sentence (e.g. "capacity noise dominated: 340 failures across 22 days").

### Step 3: Append the narrative to the output file

Append a short `**Notes:**` bullet list to `monthly-summary-<anchor-date>.md`, below the script's table and footnote (never modify anything the script wrote): red days, repeat offenders with root-cause verdicts, category trend, and any sanity-check caveats. Keep it under ~10 bullets; this replaces the old table-plus-graph, not the whole section.

The finished file is **GitHub-flavored markdown** (standard `**bold**`, tables, `[label](url)` links — **not** Slack mrkdwn), paste-ready for the OpEx doc's *Test stability → Overall metrics* section. End your response by telling the user the file path; also show the file's full content so it can be reviewed without opening it.

Do not include internal ticket IDs or artefact names (same rule as the daily skill).

## Caveats

- **Artifact retention.** The daily `summary` artifact uses the repo's default retention. The skill is designed to run for last month only; older months will surface as "no summary available" runs.
- **Recurring-test counts are a lower bound.** Daily summaries truncate failing-test lists (first 10 + `, and N more`), so a test failing every day but ranked below the cut-off won't appear in `recurring_tests`.
- **Verdicts are only as good as the daily classification.** If a daily summary was wrong, this skill inherits the error by design — fix the daily skill's rules instead of overriding numbers here.
