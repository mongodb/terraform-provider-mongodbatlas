#!/usr/bin/env python3
"""Monthly test-suite stability summary for the OpEx doc.

Deterministic aggregator: it never classifies failures itself. It reuses the
per-run verdicts produced by the daily `summarize-test-suite` skill, which are
uploaded as a `summary` artifact (summary.md) on every scheduled Test Suite run.

For each scheduled run of the target month it:
  - downloads the summary artifact (cached under $TMPDIR),
  - parses the leading verdict emoji (red/yellow/green),
  - parses per-category failure counts,
  - parses named failing tests (lower bound: daily summaries cap the list),

then writes a paste-ready GitHub-flavored markdown table to
monthly-summary-<anchor-date>.md (in the current working directory) and a JSON
report for the agent's narrative step.

Usage:
    monthly_summary.py [--date YYYY-MM-DD] [--json-out PATH]

    --date      Anchor date; "last month" is computed relative to this date
                instead of today (use when the skill was not run on time).
    --json-out  Where to write the JSON report
                (default: $TMPDIR/monthly-test-suite-summary/<month>-result.json).
"""

import argparse
import datetime
import io
import json
import os
import re
import subprocess
import sys
import tempfile
import zipfile

REPO = "mongodb/terraform-provider-mongodbatlas"
WORKFLOW = "test-suite.yml"
CACHE_ROOT = os.path.join(tempfile.gettempdir(), "monthly-test-suite-summary")

# The leading emoji is the daily verdict; first occurrence wins.
VERDICT_RE = re.compile(r":(red|yellow|green)_circle:")

# Bullet form: "• Cloud capacity: 5 tests (e.g., ...)"
CATEGORY_BULLET_RE = re.compile(r"^\s*•\s*([A-Za-z /]+?):\s*(\d+)\s+tests?\b")
# Compressed form: "*Other failures*: Cloud capacity 5, Timeout 12, Cleanup 3"
CATEGORY_COMPRESSED_LINE_RE = re.compile(r"\*Other failures\*:\s*(.+)$")
CATEGORY_COMPRESSED_ITEM_RE = re.compile(r"([A-Za-z /]+?)\s+(\d+)(?:,|$)")
KNOWN_CATEGORIES = (
    "cloud capacity",
    "api errors",
    "api contract",
    "timeout",
    "cleanup",
)

FAILING_TESTS_LINE_RE = re.compile(r"\*Failing tests\*[^:]*:\s*(.+)$")
BACKTICK_RE = re.compile(r"`([^`]+)`")
SUBTEST_COUNT_SUFFIX_RE = re.compile(r"\s*\(×\d+\s+subtests\)\s*$")

REGRESSION_HEADER_RE = re.compile(r"^\*\d+ code regressions?\*")
REGRESSION_BULLET_RE = re.compile(r"^\s*•\s*(.+)$")


def gh_api(endpoint, raw=False):
    """Run `gh api`; return parsed JSON, or raw bytes when raw=True."""
    result = subprocess.run(["gh", "api", endpoint], capture_output=True)
    if result.returncode != 0:
        err = result.stderr.decode(errors="replace").strip()
        raise RuntimeError(f"gh api {endpoint} failed: {err}")
    return result.stdout if raw else json.loads(result.stdout)


def last_month_range(anchor):
    """Return (first_day, last_day) of the calendar month before the anchor's month."""
    first_of_anchor_month = anchor.replace(day=1)
    last_of_prev = first_of_anchor_month - datetime.timedelta(days=1)
    first_of_prev = last_of_prev.replace(day=1)
    return first_of_prev, last_of_prev


def list_scheduled_runs(start, end):
    created = f"{start.isoformat()}..{end.isoformat()}"
    data = gh_api(
        f"repos/{REPO}/actions/workflows/{WORKFLOW}/runs"
        f"?event=schedule&created={created}&per_page=100"
    )
    return data.get("workflow_runs", [])


def fetch_summary(run_id):
    """Return the summary.md text for a run, or None if unavailable. Cached."""
    cache_dir = os.path.join(CACHE_ROOT, str(run_id))
    cache_file = os.path.join(cache_dir, "summary.md")
    if os.path.exists(cache_file):
        with open(cache_file, encoding="utf-8") as f:
            return f.read()

    artifacts = gh_api(f"repos/{REPO}/actions/runs/{run_id}/artifacts").get(
        "artifacts", []
    )
    summary_artifact = next(
        (a for a in artifacts if a.get("name") == "summary" and not a.get("expired")),
        None,
    )
    if summary_artifact is None:
        return None

    blob = gh_api(
        f"repos/{REPO}/actions/artifacts/{summary_artifact['id']}/zip", raw=True
    )
    with zipfile.ZipFile(io.BytesIO(blob)) as zf:
        names = [n for n in zf.namelist() if n.endswith("summary.md")]
        if not names:
            return None
        text = zf.read(names[0]).decode("utf-8", errors="replace")

    os.makedirs(cache_dir, exist_ok=True)
    with open(cache_file, "w", encoding="utf-8") as f:
        f.write(text)
    return text


def parse_verdict(text):
    m = VERDICT_RE.search(text)
    return m.group(1) if m else None


def parse_categories(text):
    """Per-category failure counts from a daily summary. Returns {name: count}."""
    counts = {}
    for line in text.splitlines():
        m = CATEGORY_BULLET_RE.match(line)
        if m:
            name = m.group(1).strip().lower()
            if name in KNOWN_CATEGORIES:
                counts[name] = counts.get(name, 0) + int(m.group(2))
            continue
        m = CATEGORY_COMPRESSED_LINE_RE.search(line)
        if m:
            for item in CATEGORY_COMPRESSED_ITEM_RE.finditer(m.group(1)):
                name = item.group(1).strip().lower()
                if name in KNOWN_CATEGORIES:
                    counts[name] = counts.get(name, 0) + int(item.group(2))
    return counts


def normalize_test_name(name):
    return SUBTEST_COUNT_SUFFIX_RE.sub("", name).strip()


def parse_failing_tests(text):
    """Backticked test names from the '*Failing tests*' line. Lower bound: the
    daily summary caps this list (first 10, then ', and N more')."""
    for line in text.splitlines():
        m = FAILING_TESTS_LINE_RE.search(line)
        if m:
            names = BACKTICK_RE.findall(m.group(1))
            return [normalize_test_name(n) for n in names]
    return []


def parse_regressions(text):
    """Regression bullet lines from a red-run summary (for the narrative)."""
    regressions = []
    in_section = False
    for line in text.splitlines():
        if REGRESSION_HEADER_RE.match(line):
            in_section = True
            continue
        if in_section:
            bullet = REGRESSION_BULLET_RE.match(line)
            if bullet:
                regressions.append(bullet.group(1))
            elif line.strip() == "":
                in_section = False
    return regressions


def main():
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument(
        "--date",
        help="Anchor date YYYY-MM-DD; 'last month' is computed relative to it "
        "instead of today.",
    )
    parser.add_argument("--json-out", help="Path for the JSON report.")
    args = parser.parse_args()

    anchor = (
        datetime.date.fromisoformat(args.date) if args.date else datetime.date.today()
    )
    start, end = last_month_range(anchor)
    label = start.strftime("%b/%y")  # e.g. Jul/26
    month_id = start.strftime("%Y-%m")

    runs = list_scheduled_runs(start, end)
    runs.sort(key=lambda r: r.get("run_started_at") or r.get("created_at", ""))

    report_runs = []
    category_totals = {}
    test_days = {}  # test name -> set of ISO dates

    for run in runs:
        run_id = run["id"]
        date = (run.get("run_started_at") or run.get("created_at", ""))[:10]
        entry = {
            "run_id": run_id,
            "run_number": run.get("run_number"),
            "date": date,
            "url": run.get("html_url"),
            "conclusion": run.get("conclusion"),
            "verdict": None,
            "categories": {},
            "failing_tests": [],
            "regressions": [],
        }
        try:
            text = fetch_summary(run_id)
        except RuntimeError as e:
            print(f"warning: run {run_id}: {e}", file=sys.stderr)
            text = None

        if text is not None:
            entry["verdict"] = parse_verdict(text) or "unknown"
            entry["categories"] = parse_categories(text)
            entry["failing_tests"] = parse_failing_tests(text)
            entry["regressions"] = parse_regressions(text)
            for name, count in entry["categories"].items():
                category_totals[name] = category_totals.get(name, 0) + count
            for test in entry["failing_tests"]:
                test_days.setdefault(test, set()).add(date)
        report_runs.append(entry)

    red = sum(1 for r in report_runs if r["verdict"] == "red")
    no_summary = sum(1 for r in report_runs if r["verdict"] is None)
    unknown = sum(1 for r in report_runs if r["verdict"] == "unknown")
    without_regression = sum(
        1 for r in report_runs if r["verdict"] in ("yellow", "green")
    )
    classified = red + without_regression
    pct = f"{(without_regression / classified * 100):.2f}%" if classified else "n/a"

    recurring = [
        {"test": t, "days": len(dates), "dates": sorted(dates)}
        for t, dates in test_days.items()
        if len(dates) >= 2
    ]
    recurring.sort(key=lambda x: (-x["days"], x["test"]))

    no_summary_runs = [
        {"run_number": r["run_number"], "date": r["date"], "url": r["url"]}
        for r in report_runs
        if r["verdict"] is None
    ]
    unknown_runs = [
        {"run_number": r["run_number"], "date": r["date"], "url": r["url"]}
        for r in report_runs
        if r["verdict"] == "unknown"
    ]

    report = {
        "month": month_id,
        "label": label,
        "anchor_date": anchor.isoformat(),
        "totals": {
            "runs": len(report_runs),
            "red": red,
            "without_regression": without_regression,
            "no_summary": no_summary,
            "unknown_verdict": unknown,
            "pct_without_regression": pct,
        },
        "category_totals": category_totals,
        "recurring_tests": recurring,
        "no_summary_runs": no_summary_runs,
        "unknown_verdict_runs": unknown_runs,
        "runs": report_runs,
    }

    json_out = args.json_out or os.path.join(CACHE_ROOT, f"{month_id}-result.json")
    os.makedirs(os.path.dirname(json_out) or ".", exist_ok=True)
    with open(json_out, "w", encoding="utf-8") as f:
        json.dump(report, f, indent=2)

    # Markdown table (paste-ready for the OpEx doc).
    lines = [
        f"| TF Provider Test Suite Runs | {label} |",
        "| --- | --- |",
        f"| **Runs with code regression** | {red} |",
        f"| **Runs without code regression** (pass or infra noise only) | {without_regression} |",
        f"| **% without regression** | {pct} |",
    ]
    if no_summary_runs or unknown_runs:

        def refs(runs):
            return ", ".join(
                f"[#{r['run_number']} ({r['date']})]({r['url']})" for r in runs
            )

        exclusions = []
        if no_summary_runs:
            exclusions.append(f"no summary available: {refs(no_summary_runs)}")
        if unknown_runs:
            exclusions.append(f"unparseable summary verdict: {refs(unknown_runs)}")
        lines.append("")
        lines.append("*Excluded from the percentage: " + "; ".join(exclusions) + ".*")

    # Write the metrics section to monthly-summary-<anchor-date>.md in the cwd.
    # The agent appends the **Notes:** narrative section to this file (Step 3 of
    # the skill); it must not alter anything above that section.
    md_out = f"monthly-summary-{anchor.isoformat()}.md"
    with open(md_out, "w", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")

    print("\n".join(lines))
    print(f"\nMarkdown file: {md_out}", file=sys.stderr)
    print(f"JSON report: {json_out}", file=sys.stderr)


if __name__ == "__main__":
    main()
