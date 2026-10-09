# /// script
# requires-python = ">=3.10"
# dependencies = []
# ///
"""Download each affected job once; attribute failures without classifying them."""

import argparse
import json
import re
import subprocess
from pathlib import Path

REPO = "mongodb/terraform-provider-mongodbatlas"
ANSI = re.compile(r"\x1b\[[0-9;]*[A-Za-z]")
TIMESTAMP = re.compile(r"^\d{4}-\d{2}-\d{2}T\S+\s")
SIGNAL = re.compile(r"^\s*--- FAIL: |^FAIL(?:\s|$)|^panic:|^fatal error:|^##\[error\]")
# Locate diagnostics even when a long plan/trace separates them from FAIL.
DIAGNOSTIC = re.compile(
    r"provider produced (?:inconsistent|an unexpected)|this is a bug in the provider|"
    r"plugin did not respond|rpc error:", re.IGNORECASE,
)
RUN = re.compile(r"^=== RUN\s+(Test\S+)")
VERDICT = re.compile(r"^\s*--- (PASS|FAIL|SKIP):\s+(Test\S+)")
PACKAGE_FAIL = re.compile(r"^FAIL\s+(\S+)")
DEADLINE_PANIC = re.compile(r"^panic: test timed out after")
# Broad on purpose: the window already excludes passing tests, so a failed test's
# capacity / timeout / API diagnostics stay visible for the model to classify.
ERROR_ISH = re.compile(
    r"Error:|error:|panic:|unexpected|invalid|not found|does not exist|"
    r"timed out|deadline exceeded|OUT_OF_CAPACITY|NO_CAPACITY|No Capacity",
    re.IGNORECASE,
)
DIAGNOSTIC_CAP = 40


def github(endpoint, *args):
    return subprocess.run(
        ["gh", "api", f"repos/{REPO}/actions/{endpoint}", *args],
        check=True, capture_output=True, text=True, encoding="utf-8", errors="replace",
        timeout=120,
    ).stdout


def sanitized(log):
    return [TIMESTAMP.sub("", ANSI.sub("", line)) for line in log.splitlines()]


def analyze(log):
    """Attribute each diagnostic to the test that produced it; never classify.

    Each line belongs to the most recent `=== RUN` test, so a passing test's
    expected error cannot be attributed to a later failed test. Parallel tests
    can interleave, so treat this as a heuristic and expand the full log when it
    would change a verdict.
    """
    lines = sanitized(log)
    owner: list[str | None] = [None] * len(lines)
    current: str | None = None
    failures = {}
    packages = set()
    panics = []
    deadline_panics = []
    for i, line in enumerate(lines):
        run = RUN.match(line)
        if run:
            current = run.group(1).split("/")[0]
            owner[i] = current
            continue
        verdict = VERDICT.match(line)
        if verdict:
            status, name = verdict.group(1), verdict.group(2)
            top = name.split("/")[0]
            owner[i] = top
            if status == "FAIL":
                failure = failures.setdefault(top, {"subtests": [], "fail_line": i})
                if name != top and name not in failure["subtests"]:
                    failure["subtests"].append(name)
            continue
        owner[i] = current
        package = PACKAGE_FAIL.match(line)
        if package:
            packages.add(package.group(1))
        elif DEADLINE_PANIC.match(line):
            deadline_panics.append(f"{i + 1}: {line}")
        elif line.startswith("panic:"):
            panics.append(f"{i + 1}: {line}")

    failed = []
    for name in sorted(failures):
        window = [i for i, line_owner in enumerate(owner) if line_owner == name]
        diagnostics = [
            f"{i + 1}: {lines[i]}" for i in window
            if DIAGNOSTIC.search(lines[i]) or SIGNAL.search(lines[i]) or ERROR_ISH.search(lines[i])
        ]
        if len(diagnostics) > DIAGNOSTIC_CAP:
            omitted = len(diagnostics) - DIAGNOSTIC_CAP
            diagnostics = diagnostics[:DIAGNOSTIC_CAP] + [f"... {omitted} more diagnostic lines omitted"]
        failed.append({
            "name": name,
            "subtests": sorted(failures[name]["subtests"]),
            "fail_line": failures[name]["fail_line"] + 1,
            "diagnostics": diagnostics,
        })
    return {
        "failed_tests": failed,
        "package_failures": sorted(packages),
        "panics": panics,
        "deadline_panics": deadline_panics,
    }


def excerpts(log):
    """Keep every signal and nearby lines, with original 1-based log positions."""
    lines = sanitized(log)
    hits = [i for i, line in enumerate(lines) if SIGNAL.search(line) or DIAGNOSTIC.search(line)]
    selected = set(range(max(0, len(lines) - 40), len(lines)))
    for i in hits:
        selected.update(range(max(0, i - 30), min(len(lines), i + 6)))
    index = "\n".join(f"{i + 1}: {lines[i]}" for i in hits)
    context = []
    previous = -1
    for i in sorted(selected):
        if i != previous + 1:
            context.append("... omitted lines; read the full log to expand context ...")
        context.append(f"{i + 1}: {lines[i]}")
        previous = i
    return f"SIGNAL INDEX (not classifications)\n{index or '(no matching signals)'}\n\nCONTEXT\n" + "\n".join(context)


def prepare(run_id, output):
    # A fresh directory prevents a rerun from reusing stale logs or partial output.
    output = Path(output).resolve()
    output.mkdir(parents=True, exist_ok=False)
    run = json.loads(github(f"runs/{run_id}"))
    pages = json.loads(github(f"runs/{run_id}/jobs?filter=latest&per_page=100", "--paginate", "--slurp"))
    jobs = []
    for page in pages:
        for job in page["jobs"]:
            entry = {key: job.get(key) for key in (
                "id", "name", "status", "conclusion", "html_url", "started_at", "completed_at",
            )}
            jobs.append(entry)
            entry["failed_steps"] = [step["name"] for step in job.get("steps", [])
                                     if step.get("conclusion") not in (None, "success", "skipped", "neutral")]
            # Include the summarizer in the inventory, but do not analyze itself.
            if job["name"].split(" / ")[-1] == "trigger-test-summary":
                continue
            if job["status"] != "completed" or job["conclusion"] in ("success", "skipped", "neutral"):
                continue
            try:
                log = github(f"jobs/{job['id']}/logs")
            except subprocess.CalledProcessError as error:
                if "403" in error.stderr or "SSO" in error.stderr:
                    raise  # Do not work around an access denial.
                entry["log_error"] = error.stderr.strip() or f"gh exited {error.returncode}"
                continue
            except subprocess.TimeoutExpired:
                entry["log_error"] = "log download timed out after 120 seconds"
                continue
            if not log.strip():
                entry["log_error"] = "empty log response"
                continue
            raw = output / f"{job['id']}.log"
            context = output / f"{job['id']}.txt"
            raw.write_text(log, encoding="utf-8")
            context.write_text(excerpts(log), encoding="utf-8")
            entry.update(log_file=str(raw), evidence_file=str(context), analysis=analyze(log))
    if json.loads(github(f"runs/{run_id}"))["run_attempt"] != run["run_attempt"]:
        raise RuntimeError("Run was rerun during collection; prepare again in a fresh directory")
    report = {key: run[key] for key in ("id", "run_number", "run_attempt", "head_sha", "html_url")}
    report["jobs"] = jobs
    destination = output / "run.json"
    destination.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    return destination


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("run_id", type=int)
    parser.add_argument("output", type=Path, help="New directory for run.json and job logs")
    args = parser.parse_args()
    print(prepare(args.run_id, args.output))
