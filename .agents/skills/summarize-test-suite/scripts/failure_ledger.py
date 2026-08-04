#!/usr/bin/env python3
"""Prepare, validate, and finalize GitHub Actions test-suite evidence."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from collections.abc import Iterable
from datetime import datetime
from pathlib import Path
from typing import Any, Callable

TEST_FAILURE_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\S+Z[^\S\r\n]+"
    r"--- FAIL:\s+(.+?)\s+\(\d+(?:\.\d+)?s\)(?:\s|$)"
)
PACKAGE_COMPLETION_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\S+Z "
    r"(ok|FAIL|\?)[^\S\r\n]*\t([^\t\s]+)(?:\t| \[|$)"
)
PACKAGE_TEARDOWN_RE = re.compile(r"\[ERROR\] Cleanup failed:")
GO_DEADLINE_RE = re.compile(r"panic: test timed out after")
GO_DEADLINE_EXECUTION_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\S+Z[^\S\r\n]+"
    r"panic: test timed out after\b"
)
PACKAGE_TEARDOWN_EXECUTION_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\S+Z[^\S\r\n]+"
    r"(?:\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} )?"
    r"\[ERROR\] Cleanup failed:"
)
COMPLETE_LOG_SUFFIX = b"Cleaning up orphan processes"
AUTH_PREFIX_RE = re.compile(r"-(pak|sa)$")
NESTED_SETUP_JOB_NAMES = frozenset({"get-provider-version", "change-detection"})
NESTED_REPORTING_JOB_NAMES = frozenset({"slack-notification-stream"})
REPOSITORY = "mongodb/terraform-provider-mongodbatlas"
SCHEMA_VERSION = 1
PREPARATION_STATUS_SCHEMA = "test-suite-preparation-status"
FAILURE_LEDGER_SCHEMA = "test-suite-failure-ledger"
EVIDENCE_SUMMARY_SCHEMA = "test-suite-evidence-summary"
ANALYSIS_INPUT_SCHEMA = "test-suite-analysis-input"
CLASSIFICATION_INPUT_SCHEMA = "test-suite-classification-input"
MODEL_DECISIONS_SCHEMA = "test-suite-model-decisions"
FINALIZATION_RESULT_SCHEMA = "test-suite-finalization-result"
FINALIZATION_STATUS_SCHEMA = "test-suite-finalization-status"
EXECUTION_METADATA_SCHEMA = "test-suite-summary-execution-metadata"
REPLAY_ARTIFACT_FILES = (
    "analysis-input.json",
    "model/classification-input.json",
    "execution-metadata.json",
    "failure-ledger.json",
    "model-decisions.json",
    "finalization-result.json",
    "finalization-status.json",
    "preparation-status.json",
    "evidence-summary.json",
)
SLACK_HARD_LIMIT = 2900
SLACK_SECTION_LIMIT = 3000
FAILING_TEST_LIMIT = 10
MAX_MODEL_DECISIONS_BYTES = 1_000_000
MAX_MODEL_DECISION_GROUPS = 1000
MAX_MODEL_DECISION_UNIT_IDS = 1000
MAX_MODEL_DECISION_EVIDENCE_REFS = 1000
MAX_MODEL_DECISION_CAUSE_CHARS = 240
MAX_MODEL_DECISION_AMBIGUITY_CHARS = 240
MAX_MODEL_DECISION_NOTE_CHARS = 120
MAX_MODEL_DECISION_OPTIONAL_TEXT_CHARS = 500
GITHUB_METADATA_TIMEOUT_SECONDS = 60
GITHUB_LOG_TIMEOUT_SECONDS = 120
MODEL_CATEGORIES = (
    "code_regression",
    "cloud_capacity",
    "api_contract",
    "timeout",
    "cleanup",
    "api_error",
)
DECISION_CATEGORIES = MODEL_CATEGORIES + ("unresolved",)
# Claude accepts only a subset of JSON Schema for structured output. Keep this
# projection compatible there and enforce stronger semantic and size limits in
# _validate_model_decisions.
MODEL_OUTPUT_TEXT_PATTERN = r"^[^\r\n]+$"
MODEL_DECISIONS_OUTPUT_SCHEMA = {
    "type": "object",
    "additionalProperties": False,
    "required": ["schema", "analysis_digest", "groups"],
    "properties": {
        "schema": {
            "type": "object",
            "additionalProperties": False,
            "required": ["name", "version"],
            "properties": {
                "name": {
                    "type": "string",
                    "const": MODEL_DECISIONS_SCHEMA,
                },
                "version": {
                    "type": "integer",
                    "const": SCHEMA_VERSION,
                },
            },
        },
        "analysis_digest": {
            "type": "string",
            "pattern": "^[0-9a-f]{64}$",
        },
        "groups": {
            "type": "array",
            "items": {
                "type": "object",
                "additionalProperties": False,
                "required": ["category", "cause", "evidence_refs"],
                "anyOf": [
                    {"required": ["unit_ids"]},
                    {"required": ["remaining"]},
                ],
                "properties": {
                    "unit_ids": {
                        "type": "array",
                        "minItems": 1,
                        "items": {"type": "string"},
                    },
                    "remaining": {
                        "type": "boolean",
                        "const": True,
                    },
                    "category": {
                        "type": "string",
                        "enum": list(DECISION_CATEGORIES),
                    },
                    "cause": {
                        "type": "string",
                        "pattern": MODEL_OUTPUT_TEXT_PATTERN,
                    },
                    "evidence_refs": {
                        "type": "array",
                        "minItems": 1,
                        "items": {"type": "string"},
                    },
                    "ambiguity": {
                        "type": "string",
                        "pattern": MODEL_OUTPUT_TEXT_PATTERN,
                    },
                    "note": {
                        "type": "string",
                        "enum": ["delete_on_timeout_unverified"],
                    },
                },
            },
        },
        "why": {
            "type": "string",
            "pattern": MODEL_OUTPUT_TEXT_PATTERN,
        },
        "action": {
            "type": "string",
            "pattern": MODEL_OUTPUT_TEXT_PATTERN,
        },
        "tldr": {
            "type": "string",
            "pattern": MODEL_OUTPUT_TEXT_PATTERN,
        },
    },
}
CATEGORY_LABELS = {
    "code_regression": "Code regression",
    "cloud_capacity": "Cloud capacity",
    "api_contract": "API contract",
    "timeout": "Timeout",
    "cleanup": "Cleanup",
    "api_error": "API errors",
}
PHASE_ANCHOR_RE = re.compile(r"\bStep\s+(\d+)/(\d+)\s+error:")
TEST_CONTROL_RE = re.compile(
    r"=== (?:RUN|CONT|PAUSE|NAME)\s+([^\s]+)"
)
TEST_EXECUTION_CONTROL_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\S+Z[^\S\r\n]+"
    r"=== (?:RUN|CONT|PAUSE|NAME)\s+[^\s]+"
)
GO_TEST_SOURCE_LINE_RE = re.compile(
    r"(?:^|\t)\d{4}-\d{2}-\d{2}T\S+Z[^\S\r\n]+"
    r"\S+\.go:\d+(?::\d+)?:[^\S\r\n]+\S"
)
TEST_HELPER_METADATA_RE = re.compile(r"\btest_name=\S+")
POST_TEST_DESTROY_RE = re.compile(
    r"Error running post-test destroy|error when destroying resource",
    re.IGNORECASE,
)
EXPLICIT_POST_TEST_DESTROY_RE = re.compile(
    r"Error running post-test destroy",
    re.IGNORECASE,
)
POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE = re.compile(
    r"\bCANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS\b",
    re.IGNORECASE,
)
ASSERTION_EVIDENCE_RE = re.compile(
    r"Error Trace:|Error:\s|Not equal:|Should be|expected:|actual:|"
    r"Received unexpected error:|Test:\s|Messages:\s|"
    r"diagnostic_(?:summary|detail)=|Error running",
    re.IGNORECASE,
)
HIGH_SIGNAL_EVIDENCE_RE = re.compile(
    r"panic:|OUT_OF_CAPACITY|NO_CAPACITY|HTTP (?:ERROR )?[45]\d\d|"
    r"MAX_GROUPS_PER_ORG_EXCEEDED|"
    r"INVALID_ATTRIBUTE|NOT_FOUND|does not exist|still exists|"
    r"DUPLICATE_|inconsistent result|unexpected new value|"
    r"Provider produced an invalid plan|bug in the provider|"
    r"(?:non-empty|unexpected) plan|plan was not empty|"
    r"Plugin did not respond|rpc error:|timed? out|timeout",
    re.IGNORECASE,
)
AUTHORIZATION_FAILURE_RE = re.compile(
    r"\bHTTP(?: ERROR)?\s+(?:401|403)\b",
    re.IGNORECASE,
)
AUTHORIZATION_REQUEST_RE = re.compile(
    r"\b(?P<url>https?://\S+)\s+"
    r"(?P<method>GET|POST|PUT|PATCH|DELETE):\s+"
    r"HTTP(?: ERROR)?\s+(?:401|403)\b",
    re.IGNORECASE,
)
MAX_GROUPS_PER_ORG_RE = re.compile(
    r"\bMAX_GROUPS_PER_ORG_EXCEEDED\b",
    re.IGNORECASE,
)
GROUPS_POST_HTTP_500_RE = re.compile(
    r"/api/atlas/v2/groups\s+POST:\s+HTTP 500(?:\s+Internal Server Error)?",
    re.IGNORECASE,
)
GROUPS_POST_HTTP_400_RE = re.compile(
    r"/api/atlas/v2/groups\s+POST:\s+HTTP 400(?:\s+Bad Request)?",
    re.IGNORECASE,
)
PROJECT_CREATION_FAILURE_RE = re.compile(
    r"error creating project|Project creation failed|"
    r"Error calling API in Create",
    re.IGNORECASE,
)
DETERMINISTIC_COHORT_WRAPPER_RE = re.compile(
    r"Error Trace:|"
    r"Error:\s+(?:Received unexpected error:|Error calling API in Create)|"
    r"Error running apply",
    re.IGNORECASE,
)
DETERMINISTIC_PRIMARY_FAILURE_RE = re.compile(
    r"Error Trace:|Error:\s|Received unexpected error:|"
    r"Provider produced an invalid plan|bug in the provider|"
    r"(?:non-empty|unexpected) plan|plan was not empty|"
    r"INVALID_ATTRIBUTE|OUT_OF_CAPACITY|NO_CAPACITY|"
    r"DUPLICATE_|still exists|already exists|"
    r"inconsistent result|unexpected new value|Plugin did not respond|"
    r"rpc error:|panic:|Not equal:|Should be|\bactual:",
    re.IGNORECASE,
)
DETERMINISTIC_HARD_CONFLICT_RE = re.compile(
    r"Provider produced an invalid plan|bug in the provider|"
    r"(?:non-empty|unexpected) plan|plan was not empty|"
    r"INVALID_ATTRIBUTE|OUT_OF_CAPACITY|NO_CAPACITY|"
    r"DUPLICATE_|still exists|already exists|"
    r"inconsistent result|unexpected new value|Plugin did not respond|"
    r"rpc error:|panic:|Not equal:|Should be|\bactual:|"
    r"timed? out|timeout",
    re.IGNORECASE,
)
DETERMINISTIC_TEST_POLICIES = {
    "max_groups_per_org_cleanup": {
        "category": "cleanup",
        "cause": "The test organization hit its project limit.",
    },
    "groups_post_http_500": {
        "category": "api_error",
        "cause": "Project creation failed because POST /groups returned HTTP 500.",
    },
    "post_destroy_private_endpoint_dependency_cleanup": {
        "category": "cleanup",
        "cause": "Post-test destroy left dependent resources behind.",
    },
}
SHAPE_B_TEST_RE = re.compile(
    r"create.*timeout.*delete.*on.*create|"
    r"delete.*on.*create.*create.*timeout",
    re.IGNORECASE,
)
SHAPE_B_EVIDENCE_RE = re.compile(
    r"Should be false|"
    r"(?:delete|deletion).{0,120}(?:timed? out|timeout)|"
    r"timeout while waiting for state to become ['\"]DELETED['\"]",
    re.IGNORECASE,
)
BUILD_FAILURE_EVIDENCE_RE = re.compile(
    r"\[build failed\]|undefined:|syntax error:|cannot use |"
    r"too many arguments|not enough arguments|does not implement|"
    r"no required module provides package",
    re.IGNORECASE,
)

CaptureGitHub = Callable[[str, bool], subprocess.CompletedProcess[bytes]]
DownloadGitHub = Callable[[str, Path, Path], int]


class PreparationError(RuntimeError):
    """Evidence preparation could not produce a trustworthy result."""


class FinalizationError(RuntimeError):
    """Model decisions could not be reconciled with prepared evidence."""


class WorkflowError(RuntimeError):
    """GitHub Actions orchestration could not complete safely."""


def _sorted_by_job_name(
    records: Iterable[dict[str, Any]],
) -> list[dict[str, Any]]:
    return sorted(records, key=lambda item: item["job_name"])


def _load_jobs(path: Path) -> dict[str, dict[str, Any]]:
    jobs: dict[str, dict[str, Any]] = {}
    with path.open(encoding="utf-8") as handle:
        for line_number, line in enumerate(handle, start=1):
            if not line.strip():
                continue
            try:
                job = json.loads(line)
            except json.JSONDecodeError as exc:
                raise ValueError(f"{path}:{line_number}: invalid JSON: {exc}") from exc
            job_id = str(job.get("id", ""))
            if not job_id:
                raise ValueError(f"{path}:{line_number}: job has no id")
            if job_id in jobs:
                raise ValueError(
                    f"{path}:{line_number}: duplicate job id {job_id}"
                )
            jobs[job_id] = job
    return jobs


def _job_role(name: str) -> str:
    parts = name.split(" / ")
    if parts[0] in {"clean-before", "clean-after"}:
        return "cleanup"
    if name == "variables":
        return "setup"

    nested_acceptance = (
        len(parts) == 3
        and AUTH_PREFIX_RE.search(parts[0]) is not None
        and parts[1].startswith("tests-")
    )
    leaf_name = parts[-1]
    if nested_acceptance and leaf_name in NESTED_SETUP_JOB_NAMES:
        return "setup"
    if name == "trigger-test-summary":
        return "reporting"
    if nested_acceptance and leaf_name in NESTED_REPORTING_JOB_NAMES:
        return "reporting"
    if nested_acceptance:
        return "test"
    return "unknown"


def _duration_minutes(job: dict[str, Any]) -> int | None:
    started_at = job.get("started_at")
    completed_at = job.get("completed_at")
    if not isinstance(started_at, str) or not isinstance(completed_at, str):
        return None
    try:
        started = datetime.fromisoformat(started_at.replace("Z", "+00:00"))
        completed = datetime.fromisoformat(completed_at.replace("Z", "+00:00"))
    except ValueError:
        return None
    seconds = (completed - started).total_seconds()
    if seconds < 0:
        return None
    return int(seconds // 60)


def _job_record(job_id: str, job: dict[str, Any]) -> dict[str, Any]:
    job_name = str(job.get("name", f"unknown job {job_id}"))
    return {
        "job_id": job_id,
        "job_name": job_name,
        "conclusion": str(job.get("conclusion", "unknown")),
        "role": _job_role(job_name),
        "started_at": job.get("started_at"),
        "completed_at": job.get("completed_at"),
        "duration_minutes": _duration_minutes(job),
    }


def _log_issue(log_path: Path | None, error_path: Path) -> str | None:
    if log_path is None:
        return "missing log"
    if error_path.is_file() and error_path.stat().st_size > 0:
        return "download command wrote an error"
    if log_path.stat().st_size == 0:
        return "empty log"
    with log_path.open("rb") as handle:
        handle.seek(max(0, log_path.stat().st_size - 4096))
        log_tail = handle.read()
    if COMPLETE_LOG_SUFFIX not in log_tail:
        return "log has no completion marker"
    return None


def _partial_log_proves_test_execution(log_path: Path) -> bool:
    with log_path.open(encoding="utf-8", errors="replace") as handle:
        return any(
            TEST_EXECUTION_CONTROL_RE.search(line)
            or TEST_FAILURE_RE.match(line)
            or PACKAGE_COMPLETION_RE.search(line)
            or GO_DEADLINE_EXECUTION_RE.search(line)
            or PACKAGE_TEARDOWN_EXECUTION_RE.search(line)
            for line in handle
        )


def _record_test_failure(
    pending_tests: dict[str, dict[str, Any]],
    full_name: str,
    line_number: int,
) -> None:
    parent_name, separator, _ = full_name.partition("/")
    entry = pending_tests.setdefault(
        parent_name,
        {
            "test": parent_name,
            "top_level_occurrences": 0,
            "top_level_line_numbers": [],
            "subtests": [],
            "subtest_line_numbers": [],
        },
    )
    if separator:
        if full_name not in entry["subtests"]:
            entry["subtests"].append(full_name)
        entry["subtest_line_numbers"].append(line_number)
    else:
        entry["top_level_occurrences"] += 1
        entry["top_level_line_numbers"].append(line_number)


def _flush_tests(
    tests: dict[tuple[str, str, str], dict[str, Any]],
    pending_tests: dict[str, dict[str, Any]],
    job_id: str,
    job_name: str,
    package: str | None,
    block_start_line: int,
    block_end_line: int,
) -> None:
    package_key = package or ""
    for parent_name, pending_entry in pending_tests.items():
        key = (job_id, package_key, parent_name)
        entry = tests.setdefault(
            key,
            {
                "job_id": job_id,
                "job_name": job_name,
                "package": package,
                "test": parent_name,
                "package_blocks": [],
                "top_level_occurrences": 0,
                "top_level_line_numbers": [],
                "subtests": [],
                "subtest_line_numbers": [],
            },
        )
        entry["top_level_occurrences"] += pending_entry["top_level_occurrences"]
        entry["top_level_line_numbers"].extend(
            pending_entry["top_level_line_numbers"]
        )
        entry["subtest_line_numbers"].extend(pending_entry["subtest_line_numbers"])
        entry["subtests"] = sorted(
            set(entry["subtests"]) | set(pending_entry["subtests"])
        )
        block = {
            "start_line": block_start_line,
            "end_line": block_end_line,
        }
        if block not in entry["package_blocks"]:
            entry["package_blocks"].append(block)
    pending_tests.clear()


def _record_package_teardown(
    package_teardowns: dict[tuple[str, str], dict[str, Any]],
    job_id: str,
    job_name: str,
    package: str,
    cleanup_line_numbers: list[int],
    package_failure_line_number: int,
    block_start_line: int,
) -> None:
    key = (job_id, package)
    entry = package_teardowns.setdefault(
        key,
        {
            "job_id": job_id,
            "job_name": job_name,
            "package": package,
            "occurrences": 0,
            "cleanup_line_numbers": [],
            "package_failure_line_numbers": [],
            "package_blocks": [],
        },
    )
    entry["occurrences"] += len(cleanup_line_numbers)
    entry["cleanup_line_numbers"].extend(cleanup_line_numbers)
    entry["package_failure_line_numbers"].append(package_failure_line_number)
    block = {
        "start_line": block_start_line,
        "end_line": package_failure_line_number,
    }
    if block not in entry["package_blocks"]:
        entry["package_blocks"].append(block)


def _unresolved_package_teardown(
    job_id: str,
    job_name: str,
    cleanup_line_numbers: list[int],
    block_start_line: int,
    block_end_line: int,
    reason: str,
    status: str | None = None,
    package: str | None = None,
) -> dict[str, Any]:
    return {
        "job_id": job_id,
        "job_name": job_name,
        "cleanup_line_numbers": cleanup_line_numbers.copy(),
        "package_block": {
            "start_line": block_start_line,
            "end_line": block_end_line,
        },
        "following_status": status,
        "following_package": package,
        "reason": reason,
    }


def _covered_package_blocks(
    test_entries: list[dict[str, Any]],
    teardown_entries: list[dict[str, Any]],
) -> dict[tuple[str, str], list[dict[str, int]]]:
    covered: dict[tuple[str, str], list[dict[str, int]]] = {}
    for entry in [*test_entries, *teardown_entries]:
        package = entry.get("package")
        if package is None:
            continue
        key = (str(entry["job_id"]), str(package))
        for block in entry.get("package_blocks", []):
            covered.setdefault(key, []).append(
                {
                    "start_line": int(block["start_line"]),
                    "end_line": int(block["end_line"]),
                }
            )
    return covered


def _uncovered_package_failure_entry(
    entry: dict[str, Any],
    covered: dict[tuple[str, str], list[dict[str, int]]],
) -> dict[str, Any] | None:
    key = (str(entry["job_id"]), str(entry["package"]))
    covered_blocks = covered.get(key, [])
    uncovered_blocks = [
        block
        for block in entry.get("package_blocks", [])
        if not any(
            candidate["start_line"] <= int(block["end_line"])
            <= candidate["end_line"]
            for candidate in covered_blocks
        )
    ]
    if not uncovered_blocks:
        return None
    return {
        **entry,
        "occurrences": len(uncovered_blocks),
        "line_numbers": sorted(
            {int(block["end_line"]) for block in uncovered_blocks}
        ),
        "package_blocks": uncovered_blocks,
        "build_failed": any(
            block.get("build_failed") is True for block in uncovered_blocks
        ),
    }


def build_ledger(logs_dir: Path, jobs_path: Path) -> dict[str, Any]:
    jobs = _load_jobs(jobs_path)
    log_paths = {path.stem: path for path in sorted(logs_dir.glob("*.log"))}

    tests: dict[tuple[str, str, str], dict[str, Any]] = {}
    package_failure_markers: dict[tuple[str, str], dict[str, Any]] = {}
    package_teardowns: dict[tuple[str, str], dict[str, Any]] = {}
    unresolved_package_teardown_markers: list[dict[str, Any]] = []
    parsed_test_execution_job_ids: set[str] = set()
    cleanup_jobs: list[dict[str, Any]] = []
    setup_jobs: list[dict[str, Any]] = []
    reporting_jobs: list[dict[str, Any]] = []
    failed_unclassified_jobs: list[dict[str, Any]] = []
    failed_test_jobs: list[dict[str, Any]] = []
    test_job_logs_for_sweep: list[dict[str, Any]] = []
    unsupported_test_jobs: list[dict[str, Any]] = []
    cancelled_jobs: list[dict[str, Any]] = []
    cancelled_cleanup_jobs: list[dict[str, Any]] = []
    timed_out_jobs: list[dict[str, Any]] = []
    timed_out_cleanup_jobs: list[dict[str, Any]] = []
    missing_log_jobs: list[dict[str, Any]] = []
    empty_log_jobs: list[dict[str, Any]] = []
    download_error_jobs: list[dict[str, Any]] = []
    partial_log_jobs: list[dict[str, Any]] = []
    cleanup_log_issues: list[dict[str, Any]] = []
    setup_log_issues: list[dict[str, Any]] = []
    reporting_log_issues: list[dict[str, Any]] = []
    unclassified_log_issues: list[dict[str, Any]] = []
    unresolved_failed_jobs: list[dict[str, Any]] = []
    deadline_panic_jobs: list[dict[str, Any]] = []
    job_records = {
        job_id: _job_record(job_id, job) for job_id, job in jobs.items()
    }
    unclassified_jobs = [
        record for record in job_records.values() if record["role"] == "unknown"
    ]
    active_jobs = [
        record
        for job_id, record in job_records.items()
        if jobs[job_id].get("conclusion") is None
    ]
    active_test_jobs = [
        record for record in active_jobs if record["role"] == "test"
    ]
    active_setup_jobs = [
        record for record in active_jobs if record["role"] == "setup"
    ]
    active_cleanup_jobs = [
        record for record in active_jobs if record["role"] == "cleanup"
    ]
    active_reporting_jobs = [
        record for record in active_jobs if record["role"] == "reporting"
    ]
    active_unclassified_jobs = [
        record for record in active_jobs if record["role"] == "unknown"
    ]

    failed_jobs = {
        job_id: job
        for job_id, job in jobs.items()
        if job.get("conclusion") not in {"success", "skipped", None}
    }

    for job_id, job in failed_jobs.items():
        name = str(job.get("name", f"unknown job {job_id}"))
        conclusion = str(job.get("conclusion", "unknown"))
        record = job_records[job_id]
        role = record["role"]

        log_path = log_paths.get(job_id)
        issue = _log_issue(
            log_path,
            logs_dir / f"{job_id}.err",
        )

        if role == "cleanup":
            if conclusion == "cancelled":
                cancelled_cleanup_jobs.append(record)
            else:
                cleanup_jobs.append(record)
                if conclusion == "timed_out":
                    timed_out_cleanup_jobs.append(record)
            if issue:
                cleanup_log_issues.append({**record, "reason": issue})
            continue

        if role == "setup":
            setup_jobs.append(record)
            if issue:
                setup_log_issues.append({**record, "reason": issue})
            continue

        if role == "reporting":
            reporting_jobs.append(record)
            if issue:
                reporting_log_issues.append({**record, "reason": issue})
            continue

        if role == "unknown":
            failed_unclassified_jobs.append(record)
            if issue:
                unclassified_log_issues.append({**record, "reason": issue})
            continue

        test_job_logs_for_sweep.append(record)
        if issue:
            if issue == "missing log":
                missing_log_jobs.append(record)
            elif issue == "empty log":
                empty_log_jobs.append(record)
            elif issue == "download command wrote an error":
                download_error_jobs.append(record)
            else:
                partial_log_jobs.append(record)
        if conclusion == "cancelled":
            cancelled_jobs.append(record)
            continue
        if conclusion not in {"failure", "timed_out"}:
            unsupported_test_jobs.append(record)
            continue
        failed_test_jobs.append(record)
        if conclusion == "timed_out":
            timed_out_jobs.append(record)

        parse_timed_out_partial_log = (
            conclusion == "timed_out"
            and issue == "log has no completion marker"
        )
        if (
            issue == "log has no completion marker"
            and log_path is not None
            and _partial_log_proves_test_execution(log_path)
        ):
            parsed_test_execution_job_ids.add(job_id)
        if issue and not parse_timed_out_partial_log:
            continue

        assert log_path is not None

        job_test_count_before = len(tests)
        job_package_count_before = len(package_failure_markers)
        job_teardown_count_before = len(package_teardowns)
        pending_tests: dict[str, dict[str, Any]] = {}
        pending_package_teardown_lines: list[int] = []
        deadline_lines: list[int] = []
        block_start_line = 1
        last_line_number = 0

        with log_path.open(encoding="utf-8", errors="replace") as handle:
            for line_number, line in enumerate(handle, start=1):
                last_line_number = line_number
                if TEST_EXECUTION_CONTROL_RE.search(line):
                    parsed_test_execution_job_ids.add(job_id)
                if GO_DEADLINE_RE.search(line):
                    deadline_lines.append(line_number)
                    parsed_test_execution_job_ids.add(job_id)
                if PACKAGE_TEARDOWN_RE.search(line):
                    pending_package_teardown_lines.append(line_number)
                    parsed_test_execution_job_ids.add(job_id)

                test_match = TEST_FAILURE_RE.match(line)
                if test_match:
                    parsed_test_execution_job_ids.add(job_id)
                    _record_test_failure(
                        pending_tests,
                        test_match.group(1),
                        line_number,
                    )

                package_match = PACKAGE_COMPLETION_RE.search(line)
                if package_match:
                    parsed_test_execution_job_ids.add(job_id)
                    status, package = package_match.groups()
                    _flush_tests(
                        tests,
                        pending_tests,
                        job_id,
                        name,
                        package,
                        block_start_line,
                        line_number,
                    )
                    if pending_package_teardown_lines:
                        if status == "FAIL":
                            _record_package_teardown(
                                package_teardowns,
                                job_id,
                                name,
                                package,
                                pending_package_teardown_lines,
                                line_number,
                                block_start_line,
                            )
                        else:
                            unresolved_package_teardown_markers.append(
                                _unresolved_package_teardown(
                                    job_id,
                                    name,
                                    pending_package_teardown_lines,
                                    block_start_line,
                                    line_number,
                                    "cleanup marker was not followed by a failed package",
                                    status,
                                    package,
                                )
                            )
                        pending_package_teardown_lines.clear()
                    if status == "FAIL":
                        build_failed = "[build failed]" in line.casefold()
                        key = (job_id, package)
                        entry = package_failure_markers.setdefault(
                            key,
                            {
                                "job_id": job_id,
                                "job_name": name,
                                "package": package,
                                "occurrences": 0,
                                "line_numbers": [],
                                "package_blocks": [],
                                "build_failed": False,
                            },
                        )
                        entry["occurrences"] += 1
                        entry["line_numbers"].append(line_number)
                        entry["build_failed"] = (
                            entry["build_failed"]
                            or build_failed
                        )
                        entry["package_blocks"].append(
                            {
                                "start_line": block_start_line,
                                "end_line": line_number,
                                "build_failed": build_failed,
                            }
                        )
                    block_start_line = line_number + 1

        _flush_tests(
            tests,
            pending_tests,
            job_id,
            name,
            None,
            block_start_line,
            last_line_number,
        )
        if pending_package_teardown_lines:
            unresolved_package_teardown_markers.append(
                _unresolved_package_teardown(
                    job_id,
                    name,
                    pending_package_teardown_lines,
                    block_start_line,
                    last_line_number,
                    "cleanup marker was not followed by a package completion",
                )
            )
        if deadline_lines:
            deadline_panic_jobs.append(
                {
                    **record,
                    "line_numbers": deadline_lines,
                }
            )

        if (
            len(tests) == job_test_count_before
            and len(package_failure_markers) == job_package_count_before
            and len(package_teardowns) == job_teardown_count_before
        ):
            unresolved_failed_jobs.append(record)

    test_entries = sorted(
        tests.values(),
        key=lambda item: (item["job_name"], item["package"] or "", item["test"]),
    )
    package_entries = sorted(
        package_failure_markers.values(),
        key=lambda item: (item["job_name"], item["package"]),
    )
    package_teardown_entries = sorted(
        package_teardowns.values(),
        key=lambda item: (item["job_name"], item["package"]),
    )
    for entry in test_entries:
        entry["package_blocks"].sort(key=lambda block: block["start_line"])
        entry["subtests"].sort()
        entry["derived_from_subtest_only"] = entry["top_level_occurrences"] == 0
    for entry in package_entries:
        entry["line_numbers"] = sorted(set(entry["line_numbers"]))
        blocks_by_range: dict[tuple[int, int], dict[str, Any]] = {}
        for block in entry["package_blocks"]:
            key = (int(block["start_line"]), int(block["end_line"]))
            normalized = blocks_by_range.setdefault(
                key,
                {
                    "start_line": key[0],
                    "end_line": key[1],
                    "build_failed": False,
                },
            )
            normalized["build_failed"] = (
                normalized["build_failed"]
                or block.get("build_failed") is True
            )
        entry["package_blocks"] = sorted(
            blocks_by_range.values(),
            key=lambda block: block["start_line"],
        )
    for entry in package_teardown_entries:
        entry["cleanup_line_numbers"] = sorted(set(entry["cleanup_line_numbers"]))
        entry["package_failure_line_numbers"] = sorted(
            set(entry["package_failure_line_numbers"])
        )
        entry["package_blocks"].sort(key=lambda block: block["start_line"])
    covered_package_ranges = _covered_package_blocks(
        test_entries,
        package_teardown_entries,
    )
    unresolved_package_failures = [
        uncovered
        for entry in package_entries
        if (
            uncovered := _uncovered_package_failure_entry(
                entry,
                covered_package_ranges,
            )
        )
        is not None
        and not uncovered["build_failed"]
    ]

    repeated_top_level_lines = sum(
        max(0, entry["top_level_occurrences"] - 1) for entry in test_entries
    )
    subtest_only_parents = [
        {
            "job_id": entry["job_id"],
            "job_name": entry["job_name"],
            "package": entry["package"],
            "test": entry["test"],
        }
        for entry in test_entries
        if entry["derived_from_subtest_only"]
    ]
    package_unattributed_tests = [
        {
            "job_id": entry["job_id"],
            "job_name": entry["job_name"],
            "test": entry["test"],
        }
        for entry in test_entries
        if entry["package"] is None
    ]

    unavailable_jobs = (
        missing_log_jobs
        + empty_log_jobs
        + download_error_jobs
        + partial_log_jobs
    )
    role_counts = {
        role: sum(
            1
            for job in jobs.values()
            if _job_role(str(job.get("name", ""))) == role
        )
        for role in ("test", "setup", "cleanup", "reporting", "unknown")
    }
    successful_test_job_ids = {
        job_id
        for job_id, job in jobs.items()
        if _job_role(str(job.get("name", ""))) == "test"
        and job.get("conclusion") == "success"
    }
    successful_test_jobs = len(successful_test_job_ids)
    test_jobs_ran = len(
        successful_test_job_ids | parsed_test_execution_job_ids
    )
    skipped_test_jobs = sum(
        1
        for job in jobs.values()
        if _job_role(str(job.get("name", ""))) == "test"
        and job.get("conclusion") == "skipped"
    )
    skipped_workflow_jobs = sum(
        1 for job in jobs.values() if job.get("conclusion") == "skipped"
    )
    test_jobs_without_execution_evidence = [
        record
        for job_id, record in job_records.items()
        if record["role"] == "test"
        and jobs[job_id].get("conclusion") not in {"success", "skipped"}
        and job_id not in parsed_test_execution_job_ids
    ]
    if test_jobs_ran:
        test_execution_state = "confirmed_ran"
    elif (
        test_jobs_without_execution_evidence
        or unclassified_jobs
    ):
        test_execution_state = "unknown"
    else:
        test_execution_state = "confirmed_none"
    cancelled_setup_jobs = [
        job for job in setup_jobs if job["conclusion"] == "cancelled"
    ]
    timed_out_setup_jobs = [
        job for job in setup_jobs if job["conclusion"] == "timed_out"
    ]
    ledger_complete = (
        not unavailable_jobs
        and not unresolved_failed_jobs
        and not cancelled_jobs
        and not timed_out_jobs
        and not unsupported_test_jobs
        and not package_unattributed_tests
        and not deadline_panic_jobs
        and not active_test_jobs
        and not active_unclassified_jobs
        and not unclassified_jobs
        and not unresolved_package_teardown_markers
        and not unresolved_package_failures
    )
    non_reporting_logs_complete = (
        ledger_complete
        and not cleanup_log_issues
        and not setup_log_issues
        and not unclassified_log_issues
        and not cancelled_cleanup_jobs
        and not timed_out_cleanup_jobs
        and not cancelled_setup_jobs
        and not timed_out_setup_jobs
        and not active_setup_jobs
        and not active_cleanup_jobs
    )

    warning_conditions = (
        (repeated_top_level_lines, "Repeated top-level failure lines were collapsed by job, package, and test name."),
        (subtest_only_parents, "Some parent tests were derived from subtest failures because no parent failure line was present."),
        (unavailable_jobs, "One or more failed test jobs had no readable log."),
        (unresolved_failed_jobs, "One or more failed test jobs had no test or package failure marker."),
        (cancelled_jobs, "One or more cancelled test jobs were unverifiable."),
        (cancelled_cleanup_jobs, "One or more cleanup jobs were cancelled."),
        (timed_out_jobs, "One or more timed-out test jobs may have ended before all failures were emitted."),
        (unsupported_test_jobs, "One or more test jobs had unsupported terminal conclusions and were excluded from test accounting."),
        (timed_out_cleanup_jobs, "One or more cleanup jobs timed out."),
        (setup_log_issues, "One or more failed setup jobs had no complete readable log."),
        (cancelled_setup_jobs, "One or more setup jobs were cancelled."),
        (timed_out_setup_jobs, "One or more setup jobs timed out."),
        (reporting_log_issues, "One or more failed reporting jobs had no complete readable log; reporting does not affect test-ledger completeness."),
        (unclassified_jobs, "One or more workflow jobs did not match a known job role and were excluded from test accounting."),
        (unresolved_package_teardown_markers, "One or more explicit package-cleanup markers could not be attributed to a failed package."),
        (unresolved_package_failures, "One or more packages failed without a canonical test, teardown, or build-failure marker."),
        (package_unattributed_tests, "One or more failed tests could not be attributed to a package."),
        (deadline_panic_jobs, "A Go deadline panic may identify running tests without emitting top-level failure lines."),
        (active_test_jobs or active_setup_jobs or active_cleanup_jobs, "One or more non-reporting workflow jobs are still active."),
        (active_reporting_jobs, "One or more reporting jobs are still active; reporting does not affect test-ledger completeness."),
        (active_unclassified_jobs, "One or more active jobs did not match a known job role."),
        (cleanup_log_issues, "One or more failed cleanup jobs had no complete readable log."),
    )
    warnings = [
        message for condition, message in warning_conditions if condition
    ]

    return {
        "schema": _schema(FAILURE_LEDGER_SCHEMA),
        "summary": {
            "unique_failed_tests": len(test_entries),
            "top_level_failure_lines": sum(
                entry["top_level_occurrences"] for entry in test_entries
            ),
            "subtest_failure_lines": sum(
                len(entry["subtest_line_numbers"]) for entry in test_entries
            ),
            "repeated_top_level_failure_lines_collapsed": repeated_top_level_lines,
            "subtest_only_parent_tests": len(subtest_only_parents),
            "package_unattributed_tests": len(package_unattributed_tests),
            "package_failure_markers": len(package_entries),
            "unresolved_package_failures": len(
                unresolved_package_failures
            ),
            "package_teardowns": len(package_teardown_entries),
            "package_teardown_marker_lines": sum(
                entry["occurrences"] for entry in package_teardown_entries
            ),
            "unresolved_package_teardown_markers": len(
                unresolved_package_teardown_markers
            ),
            "test_jobs": role_counts["test"],
            "setup_jobs": role_counts["setup"],
            "cleanup_jobs": role_counts["cleanup"],
            "reporting_jobs": role_counts["reporting"],
            "unclassified_jobs": role_counts["unknown"],
            "test_jobs_ran": test_jobs_ran,
            "test_execution_state": test_execution_state,
            "successful_test_jobs": successful_test_jobs,
            "skipped_test_jobs": skipped_test_jobs,
            "skipped_workflow_jobs": skipped_workflow_jobs,
            "failed_test_jobs": len(failed_test_jobs),
            "test_job_logs_for_sweep": len(test_job_logs_for_sweep),
            "cleanup_job_failures": len(cleanup_jobs),
            "setup_job_failures": len(setup_jobs),
            "cancelled_setup_job_failures": len(cancelled_setup_jobs),
            "timed_out_setup_job_failures": len(timed_out_setup_jobs),
            "reporting_job_failures": len(reporting_jobs),
            "unclassified_job_failures": len(failed_unclassified_jobs),
            "cancelled_job_failures": len(cancelled_jobs),
            "cancelled_cleanup_job_failures": len(cancelled_cleanup_jobs),
            "timed_out_job_failures": len(timed_out_jobs),
            "unsupported_test_job_conclusions": len(unsupported_test_jobs),
            "timed_out_cleanup_job_failures": len(timed_out_cleanup_jobs),
            "unavailable_test_job_logs": len(unavailable_jobs),
            "cleanup_jobs_with_unavailable_logs": len(cleanup_log_issues),
            "setup_jobs_with_unavailable_logs": len(setup_log_issues),
            "reporting_jobs_with_unavailable_logs": len(reporting_log_issues),
            "unresolved_failed_test_jobs": len(unresolved_failed_jobs),
            "deadline_panic_jobs": len(deadline_panic_jobs),
            "active_workflow_jobs": len(active_jobs),
            "active_test_jobs": len(active_test_jobs),
            "active_setup_jobs": len(active_setup_jobs),
            "active_cleanup_jobs": len(active_cleanup_jobs),
            "active_reporting_jobs": len(active_reporting_jobs),
            "active_unclassified_jobs": len(active_unclassified_jobs),
            "ledger_complete": ledger_complete,
            "non_reporting_logs_complete": non_reporting_logs_complete,
        },
        "tests": test_entries,
        "package_failure_markers": package_entries,
        "unresolved_package_failures": unresolved_package_failures,
        "package_teardowns": package_teardown_entries,
        "jobs": _sorted_by_job_name(job_records.values()),
        "unresolved_package_teardown_markers": sorted(
            unresolved_package_teardown_markers,
            key=lambda item: (
                item["job_name"],
                item["package_block"]["start_line"],
            ),
        ),
        "cleanup_jobs": _sorted_by_job_name(cleanup_jobs),
        "setup_jobs": _sorted_by_job_name(setup_jobs),
        "reporting_jobs": _sorted_by_job_name(reporting_jobs),
        "unclassified_jobs": _sorted_by_job_name(unclassified_jobs),
        "failed_unclassified_jobs": _sorted_by_job_name(
            failed_unclassified_jobs
        ),
        "failed_test_jobs": _sorted_by_job_name(failed_test_jobs),
        "test_job_logs_for_sweep": _sorted_by_job_name(
            test_job_logs_for_sweep
        ),
        "cancelled_jobs": _sorted_by_job_name(cancelled_jobs),
        "cancelled_cleanup_jobs": _sorted_by_job_name(
            cancelled_cleanup_jobs
        ),
        "timed_out_jobs": _sorted_by_job_name(timed_out_jobs),
        "unsupported_test_jobs": _sorted_by_job_name(
            unsupported_test_jobs
        ),
        "timed_out_cleanup_jobs": _sorted_by_job_name(
            timed_out_cleanup_jobs
        ),
        "missing_log_jobs": _sorted_by_job_name(missing_log_jobs),
        "empty_log_jobs": _sorted_by_job_name(empty_log_jobs),
        "download_error_jobs": _sorted_by_job_name(download_error_jobs),
        "partial_log_jobs": _sorted_by_job_name(partial_log_jobs),
        "cleanup_log_issues": _sorted_by_job_name(cleanup_log_issues),
        "setup_log_issues": _sorted_by_job_name(setup_log_issues),
        "reporting_log_issues": _sorted_by_job_name(
            reporting_log_issues
        ),
        "unclassified_log_issues": _sorted_by_job_name(
            unclassified_log_issues
        ),
        "unresolved_failed_jobs": _sorted_by_job_name(
            unresolved_failed_jobs
        ),
        "deadline_panic_jobs": _sorted_by_job_name(deadline_panic_jobs),
        "active_jobs": _sorted_by_job_name(active_jobs),
        "active_test_jobs": _sorted_by_job_name(active_test_jobs),
        "active_setup_jobs": _sorted_by_job_name(active_setup_jobs),
        "active_cleanup_jobs": _sorted_by_job_name(active_cleanup_jobs),
        "active_reporting_jobs": _sorted_by_job_name(
            active_reporting_jobs
        ),
        "active_unclassified_jobs": _sorted_by_job_name(
            active_unclassified_jobs
        ),
        "subtest_only_parents": subtest_only_parents,
        "package_unattributed_tests": package_unattributed_tests,
        "warnings": warnings,
    }


def _atomic_write(path: Path, data: bytes) -> None:
    temporary_path = path.with_name(f".{path.name}.tmp")
    try:
        with temporary_path.open("wb") as handle:
            handle.write(data)
        os.replace(temporary_path, path)
    finally:
        if temporary_path.exists():
            temporary_path.unlink()


def _atomic_write_json(path: Path, value: Any, *, compact: bool = False) -> None:
    if compact:
        text = json.dumps(value, separators=(",", ":"), sort_keys=True)
    else:
        text = json.dumps(value, indent=2, sort_keys=True)
    _atomic_write(path, f"{text}\n".encode())


def _invalidate_finalization_outputs(paths: list[Path]) -> None:
    failures: list[tuple[Path, Exception]] = []
    for path in paths:
        try:
            _atomic_write(path, b"")
        except Exception as exc:
            failures.append((path, exc))
    if failures:
        names = ", ".join(path.name for path, _ in failures)
        raise OSError(
            f"failed to invalidate finalization outputs: {names}"
        ) from failures[0][1]


def _atomic_write_analysis_input(path: Path, value: dict[str, Any]) -> None:
    preferred_order = (
        "schema",
        "provenance",
        "run_context",
        "gates",
        "unit_counts",
        "analysis_digest",
        "deterministic_test_groups",
        "deterministic_cohorts",
        "test_units",
        "package_failure_units",
        "package_teardown_units",
        "workflow_job_units",
        "issues",
        "evidence",
    )
    keys = [
        *[key for key in preferred_order if key in value],
        *sorted(set(value) - set(preferred_order)),
    ]
    lines = ["{"]
    for key_index, key in enumerate(keys):
        suffix = "," if key_index < len(keys) - 1 else ""
        encoded_key = json.dumps(key)
        item = value[key]
        if isinstance(item, list):
            lines.append(f"  {encoded_key}: [")
            for item_index, entry in enumerate(item):
                entry_suffix = "," if item_index < len(item) - 1 else ""
                encoded_entry = json.dumps(
                    entry,
                    separators=(",", ":"),
                    sort_keys=True,
                    ensure_ascii=False,
                )
                lines.append(f"    {encoded_entry}{entry_suffix}")
            lines.append(f"  ]{suffix}")
        else:
            encoded_item = json.dumps(
                item,
                separators=(",", ":"),
                sort_keys=True,
                ensure_ascii=False,
            )
            lines.append(f"  {encoded_key}: {encoded_item}{suffix}")
    lines.append("}")
    text = "\n".join(lines)
    _atomic_write(path, f"{text}\n".encode())


def _atomic_write_jsonl(path: Path, values: list[dict[str, Any]]) -> None:
    text = "".join(
        f"{json.dumps(value, separators=(',', ':'), sort_keys=True)}\n"
        for value in values
    )
    _atomic_write(path, text.encode())


def _schema(name: str) -> dict[str, Any]:
    return {"name": name, "version": SCHEMA_VERSION}


def _require_schema(value: Any, name: str, label: str) -> dict[str, Any]:
    if not isinstance(value, dict):
        raise FinalizationError(f"{label} must be a JSON object")
    if value.get("schema") != _schema(name):
        raise FinalizationError(
            f"{label} schema must be {name} version {SCHEMA_VERSION}"
        )
    return value


def _canonical_digest(value: Any, *, omit_key: str | None = None) -> str:
    if omit_key is not None and isinstance(value, dict):
        value = {key: item for key, item in value.items() if key != omit_key}
    encoded = json.dumps(
        value,
        separators=(",", ":"),
        sort_keys=True,
        ensure_ascii=False,
    ).encode()
    return hashlib.sha256(encoded).hexdigest()


def _required_workflow_env(name: str) -> str:
    value = os.environ.get(name)
    if value is None or not value.strip():
        raise WorkflowError(f"required environment variable is unavailable: {name}")
    if "\r" in value or "\n" in value:
        raise WorkflowError(f"environment variable contains a newline: {name}")
    return value


def _workflow_run_id() -> int:
    value = _required_workflow_env("GITHUB_RUN_ID")
    try:
        run_id = int(value)
    except ValueError as exc:
        raise WorkflowError("GITHUB_RUN_ID must be an integer") from exc
    if run_id <= 0:
        raise WorkflowError("GITHUB_RUN_ID must be positive")
    return run_id


def _workflow_evidence_dir(run_id: int) -> Path:
    runner_temp = Path(_required_workflow_env("RUNNER_TEMP"))
    if not runner_temp.is_absolute():
        raise WorkflowError("RUNNER_TEMP must be absolute")
    runner_temp = runner_temp.resolve()
    if runner_temp == Path("/"):
        raise WorkflowError("RUNNER_TEMP cannot be the filesystem root")
    return runner_temp / f"test-suite-summary-{run_id}"


def _workflow_unavailable_summary(run_id: int) -> str:
    return (
        ":red_circle: Test Suite summary unavailable. Preparation or "
        "finalization did not complete. Review "
        f"<https://github.com/{REPOSITORY}/actions/runs/{run_id}|"
        "the target run> manually.\n"
    )


def _model_schema_output() -> str:
    schema = json.dumps(
        MODEL_DECISIONS_OUTPUT_SCHEMA,
        separators=(",", ":"),
        sort_keys=True,
    )
    if any(character in schema for character in ("'", "\r", "\n")):
        raise WorkflowError("model output schema cannot be quoted safely")
    return schema


def _append_github_outputs(values: dict[str, str]) -> None:
    output_path = Path(_required_workflow_env("GITHUB_OUTPUT"))
    if not output_path.is_absolute():
        raise WorkflowError("GITHUB_OUTPUT must be absolute")
    with output_path.open("a", encoding="utf-8") as handle:
        handle.write(
            "".join(f"{name}={value}\n" for name, value in values.items())
        )


def _append_github_step_summary(summary: str) -> None:
    summary_path = Path(_required_workflow_env("GITHUB_STEP_SUMMARY"))
    if not summary_path.is_absolute():
        raise WorkflowError("GITHUB_STEP_SUMMARY must be absolute")
    with summary_path.open("a", encoding="utf-8") as handle:
        handle.write(summary.rstrip())
        handle.write("\n")


def _workflow_model_required(output_dir: Path, run_id: int) -> bool:
    status_path = output_dir / "preparation-status.json"
    classification_path = output_dir / "model" / "classification-input.json"
    try:
        status = json.loads(status_path.read_text(encoding="utf-8"))
        classification = json.loads(
            classification_path.read_text(encoding="utf-8")
        )
    except (OSError, json.JSONDecodeError) as exc:
        raise WorkflowError(
            "prepared evidence status or classification input is unavailable"
        ) from exc
    status = _require_schema(
        status,
        PREPARATION_STATUS_SCHEMA,
        "prepared evidence status",
    )
    classification = _require_schema(
        classification,
        CLASSIFICATION_INPUT_SCHEMA,
        "classification input",
    )
    if status.get("run_id") != run_id or status.get("status") != "ready":
        raise WorkflowError("prepared evidence is not ready")
    unit_counts = classification.get("unit_counts")
    model_tests = (
        unit_counts.get("model_tests")
        if isinstance(unit_counts, dict)
        else None
    )
    if type(model_tests) is not int or model_tests < 0:
        raise WorkflowError(
            "classification input has no valid model test count"
        )
    return model_tests > 0


def _stage_replay_artifact(output_dir: Path) -> tuple[str, ...]:
    artifact_dir = output_dir / "artifacts"
    replay_dir = artifact_dir / "summary-replay"
    artifact_dir.mkdir(parents=True, exist_ok=True)
    staging_dir = Path(tempfile.mkdtemp(prefix=".summary-replay-", dir=artifact_dir))
    copied_files: list[str] = []
    try:
        for relative_name in REPLAY_ARTIFACT_FILES:
            source = output_dir / relative_name
            if not source.is_file():
                continue
            target = staging_dir / relative_name
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(source, target)
            copied_files.append(relative_name)
        if not copied_files:
            raise WorkflowError("no replay files were available")

        previous_dir = staging_dir.with_name(f"{staging_dir.name}-previous")
        if replay_dir.exists():
            os.replace(replay_dir, previous_dir)
        try:
            os.replace(staging_dir, replay_dir)
        except Exception:
            if previous_dir.exists() and not replay_dir.exists():
                os.replace(previous_dir, replay_dir)
            raise
        if previous_dir.exists():
            shutil.rmtree(previous_dir)
        return tuple(copied_files)
    finally:
        if staging_dir.exists():
            shutil.rmtree(staging_dir)


def _stable_id(kind: str, *parts: object) -> str:
    identity = "\0".join(str(part) for part in parts)
    digest = hashlib.sha256(identity.encode()).hexdigest()[:16]
    return f"{kind}-{digest}"


def _read_log_lines(logs_dir: Path, job_id: str) -> list[str]:
    path = logs_dir / f"{job_id}.log"
    if not path.is_file():
        return []
    return path.read_text(encoding="utf-8", errors="replace").splitlines()


def _line_text(lines: list[str], line_number: int) -> str:
    if 1 <= line_number <= len(lines):
        return lines[line_number - 1].strip()
    return ""


def _job_context(job_name: str) -> tuple[str | None, str | None]:
    parts = job_name.split(" / ")
    if len(parts) != 3:
        return None, None
    auth_match = AUTH_PREFIX_RE.search(parts[0])
    auth = auth_match.group(1) if auth_match else None
    environment_match = re.search(
        r"-(dev|qa|stage|staging|prod|gov)$",
        parts[1],
        re.IGNORECASE,
    )
    environment = (
        environment_match.group(1).lower() if environment_match else None
    )
    return environment, auth


def _single_or_mixed(values: set[str]) -> str:
    if not values:
        return "unknown"
    if len(values) == 1:
        return next(iter(values))
    return "mixed"


def _test_phase_and_evidence(
    entry: dict[str, Any],
    lines: list[str],
    unit_id: str,
) -> tuple[
    dict[str, Any],
    list[dict[str, Any]],
    list[str],
    dict[str, Any],
    list[dict[str, Any]],
    list[dict[str, Any]],
    bool,
]:
    job_id = str(entry["job_id"])
    evidence: list[dict[str, Any]] = []
    evidence_refs: list[str] = []
    seen_refs: set[str] = set()

    def add_line(line_number: int, kind: str) -> None:
        if line_number < 1 or line_number > len(lines):
            return
        evidence_id = f"line:{unit_id}:{job_id}:{line_number}"
        if evidence_id in seen_refs:
            return
        seen_refs.add(evidence_id)
        evidence_refs.append(evidence_id)
        evidence.append(
            {
                "evidence_id": evidence_id,
                "unit_id": unit_id,
                "kind": kind,
                "job_id": job_id,
                "line_number": line_number,
                "text": _line_text(lines, line_number),
            }
        )

    failure_lines = sorted(
        set(
            entry.get("top_level_line_numbers", [])
            + entry.get("subtest_line_numbers", [])
        )
    )
    for line_number in failure_lines:
        add_line(line_number, "test_failure")

    package_blocks: list[dict[str, int]] = []
    for block in entry.get("package_blocks", []):
        start = int(block["start_line"])
        end = int(block["end_line"])
        package_blocks.append({"start_line": start, "end_line": end})

    test_name = str(entry["test"])
    failed_control_names = {test_name}
    failed_control_names.update(
        str(subtest) for subtest in entry.get("subtests", [])
    )
    scoped_ranges: list[tuple[int, int]] = []
    scoped_occurrences: list[list[tuple[int, int]]] = []
    scope_complete = True
    failure_blocks: dict[tuple[int, int], list[int]] = {}
    for failure_line in failure_lines:
        containing_blocks = [
            block
            for block in package_blocks
            if block["start_line"] <= failure_line <= block["end_line"]
        ]
        block = (
            containing_blocks[0]
            if containing_blocks
            else {"start_line": 1, "end_line": len(lines)}
        )
        key = (block["start_line"], block["end_line"])
        failure_blocks.setdefault(key, []).append(failure_line)

    for (block_start, _), block_failure_lines in failure_blocks.items():
        scope_end = max(block_failure_lines)
        controls: list[tuple[int, str]] = []
        for line_number in range(block_start, scope_end + 1):
            match = TEST_CONTROL_RE.search(lines[line_number - 1])
            if match:
                controls.append(
                    (
                        line_number,
                        match.group(1),
                    )
                )

        occurrence_ranges = []
        for index, (line_number, owner) in enumerate(controls):
            if owner not in failed_control_names:
                continue
            next_control = (
                controls[index + 1][0]
                if index + 1 < len(controls)
                else scope_end + 1
            )
            occurrence_ranges.append(
                (line_number, min(scope_end, next_control - 1))
            )
        if not occurrence_ranges:
            scope_complete = False
            occurrence_ranges = [
                (failure_line, failure_line)
                for failure_line in block_failure_lines
            ]
        scoped_ranges.extend(occurrence_ranges)
        scoped_occurrences.append(occurrence_ranges)

    scoped_line_numbers = {
        line_number
        for start, end in scoped_ranges
        for line_number in range(start, end + 1)
    } | set(failure_lines)
    semantic_context_candidates = {
        line_number
        for line_number in scoped_line_numbers
        if (
            PHASE_ANCHOR_RE.search(lines[line_number - 1])
            or POST_TEST_DESTROY_RE.search(lines[line_number - 1])
            or POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(
                lines[line_number - 1]
            )
            or SHAPE_B_EVIDENCE_RE.search(lines[line_number - 1])
            or ASSERTION_EVIDENCE_RE.search(lines[line_number - 1])
            or HIGH_SIGNAL_EVIDENCE_RE.search(lines[line_number - 1])
        )
    }
    source_context_candidates = set()
    for failure_line in failure_lines:
        preceding_line = failure_line - 1
        if (
            preceding_line in scoped_line_numbers
            and GO_TEST_SOURCE_LINE_RE.search(
                lines[preceding_line - 1]
            )
        ):
            source_context_candidates.add(preceding_line)
    context_candidates = (
        semantic_context_candidates | source_context_candidates
    )
    decision_context = [
        {
            "line_number": line_number,
            "text": _line_text(lines, line_number),
        }
        for line_number in sorted(semantic_context_candidates)
        if line_number not in failure_lines
        and TEST_CONTROL_RE.search(lines[line_number - 1]) is None
    ]
    full_context = [
        {
            "line_number": int(item["line_number"]),
            "text": str(item["text"]),
        }
        for item in decision_context
        if TEST_HELPER_METADATA_RE.search(str(item["text"])) is None
        and test_name not in str(item["text"])
    ]
    ordered_context_lines = sorted(
        context_candidates,
        key=lambda line_number: (
            min(abs(line_number - failure_line) for failure_line in failure_lines),
            line_number,
        ),
    )
    decision_context_by_line = {
        int(item["line_number"]): str(item["text"])
        for item in decision_context
    }
    full_post_destroy_lines = [
        line_number
        for line_number in sorted(context_candidates)
        if POST_TEST_DESTROY_RE.search(lines[line_number - 1])
    ]
    full_explicit_post_destroy_lines = [
        line_number
        for line_number in full_post_destroy_lines
        if EXPLICIT_POST_TEST_DESTROY_RE.search(lines[line_number - 1])
    ]
    full_step_anchors: list[dict[str, int]] = []
    for line_number in sorted(context_candidates):
        match = PHASE_ANCHOR_RE.search(lines[line_number - 1])
        if match:
            full_step_anchors.append(
                {
                    "line_number": line_number,
                    "step": int(match.group(1)),
                    "total_steps": int(match.group(2)),
                }
            )

    def deterministic_wrapper(line_number: int) -> bool:
        text = decision_context_by_line[line_number]
        return bool(
            MAX_GROUPS_PER_ORG_RE.search(text)
            or GROUPS_POST_HTTP_400_RE.search(text)
            or GROUPS_POST_HTTP_500_RE.search(text)
            or PROJECT_CREATION_FAILURE_RE.search(text)
            or DETERMINISTIC_COHORT_WRAPPER_RE.search(text)
            or PHASE_ANCHOR_RE.search(text)
            or POST_TEST_DESTROY_RE.search(text)
            or POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(text)
        )

    blocking_lines = [
        line_number
        for line_number in ordered_context_lines
        if line_number in decision_context_by_line
        and not deterministic_wrapper(line_number)
    ]
    primary_blocking_lines = [
        line_number
        for line_number in blocking_lines
        if DETERMINISTIC_PRIMARY_FAILURE_RE.search(
            decision_context_by_line[line_number]
        )
        is not None
    ]
    secondary_blocking_lines = [
        line_number
        for line_number in blocking_lines
        if line_number not in primary_blocking_lines
    ]
    dependency_cleanup_observed = any(
        POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(
            str(item["text"])
        )
        is not None
        for item in decision_context
    )
    deterministic_marker_lines = []
    for pattern in (
        MAX_GROUPS_PER_ORG_RE,
        GROUPS_POST_HTTP_500_RE,
        POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE,
    ):
        matching_line = next(
            (
                line_number
                for line_number in ordered_context_lines
                if pattern.search(lines[line_number - 1]) is not None
            ),
            None,
        )
        if matching_line is not None:
            deterministic_marker_lines.append(matching_line)
    if dependency_cleanup_observed:
        if full_explicit_post_destroy_lines:
            deterministic_marker_lines.append(
                full_explicit_post_destroy_lines[-1]
            )
        if full_post_destroy_lines:
            deterministic_marker_lines.append(full_post_destroy_lines[-1])
        if full_step_anchors:
            deterministic_marker_lines.append(
                full_step_anchors[-1]["line_number"]
            )
    required_lines = list(
        dict.fromkeys(
            deterministic_marker_lines
            + primary_blocking_lines
            + secondary_blocking_lines
        )
    )
    ordered_context_lines = required_lines + [
        line_number
        for line_number in ordered_context_lines
        if line_number not in required_lines
    ]
    context_line_numbers = ordered_context_lines[:60]
    for line_number in sorted(context_line_numbers):
        kind = (
            "failure_context"
            if line_number in semantic_context_candidates
            else "source_context"
        )
        add_line(line_number, kind)

    post_destroy_lines = [
        line_number
        for line_number in context_line_numbers
        if POST_TEST_DESTROY_RE.search(lines[line_number - 1])
    ]
    step_anchors: list[dict[str, int]] = []
    for line_number in context_line_numbers:
        match = PHASE_ANCHOR_RE.search(lines[line_number - 1])
        if match:
            step_anchors.append(
                {
                    "line_number": line_number,
                    "step": int(match.group(1)),
                    "total_steps": int(match.group(2)),
                }
            )
    occurrence_kinds: set[str] = set()
    for occurrence_ranges in scoped_occurrences:
        if any(
            start <= line_number <= end
            for start, end in occurrence_ranges
            for line_number in post_destroy_lines
        ):
            occurrence_kinds.add("post_test_destroy")
        elif any(
            start <= anchor["line_number"] <= end
            for start, end in occurrence_ranges
            for anchor in step_anchors
        ):
            occurrence_kinds.add("test_step")
        else:
            occurrence_kinds.add("test_execution")
    if len(occurrence_kinds) > 1:
        phase = {
            "kind": "mixed",
            "anchor_lines": failure_lines,
            "step_anchors": sorted(
                step_anchors,
                key=lambda item: item["line_number"],
            ),
        }
    elif post_destroy_lines:
        phase = {
            "kind": "post_test_destroy",
            "anchor_lines": sorted(post_destroy_lines),
            "step_anchors": sorted(
                step_anchors,
                key=lambda item: item["line_number"],
            ),
        }
    elif step_anchors:
        phase = {
            "kind": "test_step",
            "anchor_lines": [
                anchor["line_number"] for anchor in step_anchors
            ],
            "step_anchors": sorted(
                step_anchors,
                key=lambda item: item["line_number"],
            ),
        }
    else:
        phase = {
            "kind": "test_execution",
            "anchor_lines": failure_lines,
            "step_anchors": [],
        }

    identity_lines = set(failure_lines)
    context_text = "\n".join(
        lines[line_number - 1]
        for line_number in sorted(context_line_numbers)
        if line_number not in identity_lines
        and TEST_CONTROL_RE.search(lines[line_number - 1]) is None
    )
    shape_context_text = "\n".join(
        line
        for line in context_text.splitlines()
        if test_name not in line
    )
    shape_b_signal = (
        SHAPE_B_TEST_RE.search(test_name) is not None
        and SHAPE_B_EVIDENCE_RE.search(shape_context_text) is not None
    )
    leftover = (
        re.search(
            r"DUPLICATE_|still exists|already exists",
            context_text,
            re.IGNORECASE,
        )
        is not None
    )
    _, auth = _job_context(str(entry.get("job_name", "")))
    authorization_failure_observed = any(
        AUTHORIZATION_FAILURE_RE.search(str(item["text"])) is not None
        for item in decision_context
    )
    authorization_requests = {
        (match.group("method").upper(), match.group("url"))
        for item in decision_context
        if (match := AUTHORIZATION_REQUEST_RE.search(str(item["text"])))
        is not None
    }
    later_step_observed = any(
        anchor["step"] > 1 for anchor in full_step_anchors
    )
    machine_facts = {
        "leftover_indicator_observed": leftover,
        "shape_b_candidate": shape_b_signal and not leftover,
        "post_test_destroy": bool(post_destroy_lines),
        "evidence_scope_complete": scope_complete,
        "mixed_phase": len(occurrence_kinds) > 1,
        "pak_authorization_failure_candidate": (
            auth == "pak"
            and authorization_failure_observed
            and later_step_observed
        ),
        "authorization_failure_distinct_requests": len(
            authorization_requests
        ),
    }
    return (
        phase,
        evidence,
        evidence_refs,
        machine_facts,
        full_context,
        decision_context,
        len(failure_lines) == 1,
    )


def _deterministic_test_decision(
    evidence: list[dict[str, Any]],
    phase: dict[str, Any],
    machine_facts: dict[str, Any],
    test_name: str,
    full_context: list[dict[str, Any]],
    decision_context: list[dict[str, Any]],
    single_failure: bool,
    has_subtests: bool,
) -> dict[str, str] | None:
    # Go renders t.Log and t.Fatal with the same source-line shape. Keep a
    # terminal source-only diagnostic out of deterministic cohorts.
    if any(item.get("kind") == "source_context" for item in evidence):
        return None

    if (
        machine_facts.get("evidence_scope_complete") is True
        and machine_facts.get("shape_b_candidate") is not True
        and single_failure
        and not has_subtests
    ):
        marker_matches = [
            item
            for item in decision_context
            if POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(
                str(item["text"])
            )
            is not None
        ]
        explicit_post_destroy_matches = [
            item
            for item in decision_context
            if EXPLICIT_POST_TEST_DESTROY_RE.search(str(item["text"]))
            is not None
        ]
        full_step_anchor_observed = any(
            PHASE_ANCHOR_RE.search(str(item["text"])) is not None
            for item in decision_context
        )
        paired_cleanup_evidence = any(
            abs(
                int(marker["line_number"]) - int(anchor["line_number"])
            )
            <= 20
            for marker in marker_matches
            for anchor in explicit_post_destroy_matches
        )
        competing_primary = any(
            DETERMINISTIC_HARD_CONFLICT_RE.search(str(item["text"]))
            is not None
            or (
                DETERMINISTIC_PRIMARY_FAILURE_RE.search(str(item["text"]))
                is not None
                and POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(
                    str(item["text"])
                )
                is None
                and POST_TEST_DESTROY_RE.search(str(item["text"])) is None
            )
            for item in decision_context
        )
        evidence_match = next(
            (
                item
                for item in evidence
                if POST_DESTROY_PRIVATE_ENDPOINT_DEPENDENCY_RE.search(
                    str(item.get("text", ""))
                )
                is not None
            ),
            None,
        )
        if (
            paired_cleanup_evidence
            and not full_step_anchor_observed
            and not competing_primary
            and evidence_match
        ):
            signature = (
                "post_destroy_private_endpoint_dependency_cleanup"
            )
            policy = DETERMINISTIC_TEST_POLICIES[signature]
            return {
                "signature": signature,
                "category": policy["category"],
                "cause": policy["cause"],
                "evidence_ref": str(evidence_match["evidence_id"]),
            }

    if (
        machine_facts.get("evidence_scope_complete") is not True
        or machine_facts.get("mixed_phase") is True
        or machine_facts.get("post_test_destroy") is True
        or machine_facts.get("leftover_indicator_observed") is True
        or machine_facts.get("shape_b_candidate") is True
        or phase.get("kind") not in {"test_execution", "test_step"}
    ):
        return None

    context = [
        item
        for item in evidence
        if item.get("kind") != "test_failure"
        and test_name not in str(item.get("text", ""))
    ]
    last_primary_line = max(
        (
            int(item["line_number"])
            for item in full_context
            if PHASE_ANCHOR_RE.search(str(item["text"])) is not None
            or DETERMINISTIC_PRIMARY_FAILURE_RE.search(str(item["text"]))
            is not None
        ),
        default=0,
    )
    terminal_context = [
        item
        for item in full_context
        if int(item["line_number"]) >= last_primary_line
    ]

    def decision_for(
        signature: str,
        signal: re.Pattern[str],
        *,
        require_project_creation: bool,
        companion: re.Pattern[str] | None = None,
    ) -> dict[str, str] | None:
        full_matches = [
            item
            for item in terminal_context
            if signal.search(str(item["text"])) is not None
        ]
        if not full_matches:
            return None
        if require_project_creation and not any(
            PROJECT_CREATION_FAILURE_RE.search(str(item["text"]))
            is not None
            for item in terminal_context
        ):
            return None
        if any(
            signal.search(str(item["text"])) is None
            and (
                companion is None
                or companion.search(str(item["text"])) is None
            )
            and PROJECT_CREATION_FAILURE_RE.search(str(item["text"])) is None
            and DETERMINISTIC_COHORT_WRAPPER_RE.search(str(item["text"]))
            is None
            and PHASE_ANCHOR_RE.search(str(item["text"])) is None
            for item in terminal_context
        ):
            return None
        evidence_match = next(
            (
                item
                for item in context
                if signal.search(str(item.get("text", ""))) is not None
            ),
            None,
        )
        if evidence_match is None:
            return None
        policy = DETERMINISTIC_TEST_POLICIES[signature]
        return {
            "signature": signature,
            "category": policy["category"],
            "cause": policy["cause"],
            "evidence_ref": str(evidence_match["evidence_id"]),
        }

    max_groups_decision = decision_for(
        "max_groups_per_org_cleanup",
        MAX_GROUPS_PER_ORG_RE,
        require_project_creation=False,
        companion=GROUPS_POST_HTTP_400_RE,
    )
    if max_groups_decision is not None:
        return max_groups_decision
    return decision_for(
        "groups_post_http_500",
        GROUPS_POST_HTTP_500_RE,
        require_project_creation=True,
    )


def _allowed_categories(
    deterministic_decision: dict[str, str] | None,
    evidence: list[dict[str, Any]],
) -> list[str]:
    if deterministic_decision is not None:
        return [deterministic_decision["category"]]
    if not any(item.get("kind") == "failure_context" for item in evidence):
        return ["unresolved"]
    return list(DECISION_CATEGORIES)


def _workflow_units(ledger: dict[str, Any]) -> list[dict[str, Any]]:
    collections = (
        ("cleanup_jobs", "cleanup", "cleanup", "cleanup job failed"),
        (
            "cancelled_cleanup_jobs",
            "cleanup",
            "cleanup",
            "cleanup job was cancelled",
        ),
        ("setup_jobs", "setup", "code_regression", "setup job failed"),
        ("reporting_jobs", "reporting", None, "reporting job failed"),
        (
            "failed_unclassified_jobs",
            "unknown",
            None,
            "job role was not recognized",
        ),
        ("cancelled_jobs", "test_job", None, "test job was cancelled"),
        ("timed_out_jobs", "test_job", None, "test job timed out"),
        (
            "unsupported_test_jobs",
            "test_job",
            None,
            "test job had an unsupported conclusion",
        ),
        ("missing_log_jobs", "test_job", None, "test job log was missing"),
        ("empty_log_jobs", "test_job", None, "test job log was empty"),
        (
            "download_error_jobs",
            "test_job",
            None,
            "test job log download failed",
        ),
        (
            "partial_log_jobs",
            "test_job",
            None,
            "test job log was incomplete",
        ),
        (
            "unresolved_failed_jobs",
            "test_job",
            None,
            "no canonical test or package failure was recovered",
        ),
        (
            "deadline_panic_jobs",
            "test_job",
            None,
            "a deadline panic may have omitted failures",
        ),
        ("active_test_jobs", "active_test", None, "test job is active"),
        ("active_setup_jobs", "active_setup", None, "setup job is active"),
        ("active_cleanup_jobs", "active_cleanup", None, "cleanup job is active"),
        (
            "active_unclassified_jobs",
            "active_unknown",
            None,
            "unclassified job is active",
        ),
        (
            "active_reporting_jobs",
            "active_reporting",
            None,
            "reporting job is active",
        ),
    )
    by_job: dict[str, dict[str, Any]] = {}
    for collection, kind, category, default_reason in collections:
        for record in ledger.get(collection, []):
            job_id = str(record["job_id"])
            unit = by_job.setdefault(
                job_id,
                {
                    "unit_id": _stable_id("job", job_id),
                    "kind": kind,
                    "job_id": job_id,
                    "job_name": record["job_name"],
                    "role": record["role"],
                    "conclusion": record["conclusion"],
                    "duration_minutes": record.get("duration_minutes"),
                    "deterministic_category": category,
                    "sources": [],
                    "reasons": [],
                },
            )
            unit["sources"].append(collection)
            reason = str(record.get("reason") or default_reason)
            if reason not in unit["reasons"]:
                unit["reasons"].append(reason)
            if kind.startswith("active_"):
                unit["kind"] = kind
            if category is not None:
                unit["deterministic_category"] = category
    for unit in by_job.values():
        unit["sources"].sort()
        unit["reasons"].sort()
    return sorted(by_job.values(), key=lambda item: item["unit_id"])


def _package_failure_units(
    ledger: dict[str, Any],
    logs_dir: Path,
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    covered_blocks = _covered_package_blocks(
        ledger.get("tests", []),
        ledger.get("package_teardowns", []),
    )
    units: list[dict[str, Any]] = []
    evidence: list[dict[str, Any]] = []
    for entry in ledger.get("package_failure_markers", []):
        uncovered = _uncovered_package_failure_entry(entry, covered_blocks)
        if uncovered is None:
            continue
        key = (str(uncovered["job_id"]), str(uncovered["package"]))
        job_id, package = key
        unit_id = _stable_id("package-failure", job_id, package)
        lines = _read_log_lines(logs_dir, job_id)
        line_numbers = sorted(
            {
                int(line_number)
                for line_number in uncovered.get("line_numbers", [])
                if isinstance(line_number, int)
            }
        )
        context_candidates: set[int] = set(line_numbers)
        for marker_line in line_numbers:
            block = next(
                (
                    block
                    for block in uncovered.get("package_blocks", [])
                    if int(block["start_line"])
                    <= marker_line
                    <= int(block["end_line"])
                ),
                {
                    "start_line": max(1, marker_line - 30),
                    "end_line": marker_line,
                },
            )
            for line_number in range(
                max(int(block["start_line"]), marker_line - 30),
                min(marker_line, len(lines)) + 1,
            ):
                if line_number < 1:
                    continue
                if (
                    BUILD_FAILURE_EVIDENCE_RE.search(lines[line_number - 1])
                    or ASSERTION_EVIDENCE_RE.search(lines[line_number - 1])
                    or HIGH_SIGNAL_EVIDENCE_RE.search(lines[line_number - 1])
                ):
                    context_candidates.add(line_number)
        selected_lines = sorted(
            context_candidates,
            key=lambda line_number: (
                min(
                    abs(line_number - marker_line)
                    for marker_line in line_numbers
                ),
                line_number,
            ),
        )[:30]
        evidence_refs: list[str] = []
        build_failed = bool(uncovered.get("build_failed"))
        for line_number in sorted(selected_lines):
            evidence_id = f"line:{unit_id}:{job_id}:{line_number}"
            text = _line_text(lines, line_number)
            evidence_refs.append(evidence_id)
            evidence.append(
                {
                    "evidence_id": evidence_id,
                    "unit_id": unit_id,
                    "kind": "package_failure",
                    "job_id": job_id,
                    "line_number": line_number,
                    "text": text,
                }
            )
            build_failed = (
                build_failed
                or BUILD_FAILURE_EVIDENCE_RE.search(text) is not None
            )
        units.append(
            {
                "unit_id": unit_id,
                "kind": "package_failure",
                "job_id": job_id,
                "job_name": uncovered["job_name"],
                "conclusion": "failure",
                "package": package,
                "deterministic_category": (
                    "code_regression" if build_failed else "unresolved"
                ),
                "reason": (
                    "package build failed before canonical test output"
                    if build_failed
                    else "package failed without a canonical test identity"
                ),
                "evidence_refs": evidence_refs,
            }
        )
    return (
        sorted(units, key=lambda item: item["unit_id"]),
        sorted(evidence, key=lambda item: item["evidence_id"]),
    )


def _structured_issues(ledger: dict[str, Any]) -> list[dict[str, Any]]:
    issues = []
    for index, message in enumerate(ledger.get("warnings", []), start=1):
        issues.append(
            {
                "issue_id": _stable_id("issue", index, message),
                "message": message,
                "caps_confidence": not message.startswith(
                    "Repeated top-level failure lines"
                )
                and not message.startswith("Some parent tests"),
            }
        )
    return issues


def build_analysis_input(
    ledger: dict[str, Any],
    run: dict[str, Any],
    logs_dir: Path,
) -> dict[str, Any]:
    if ledger.get("schema") != _schema(FAILURE_LEDGER_SCHEMA):
        raise PreparationError(
            "failure ledger schema did not match the supported version"
        )
    if (
        ledger.get("run_id") != run.get("run_id")
        or ledger.get("run_attempt") != run.get("run_attempt")
    ):
        raise PreparationError(
            "failure ledger run identity did not match run metadata"
        )
    summary = ledger.get("summary")
    if not isinstance(summary, dict):
        raise PreparationError("failure ledger had no summary")
    tests = ledger.get("tests")
    if not isinstance(tests, list):
        raise PreparationError("failure ledger had no canonical tests")
    if summary.get("unique_failed_tests") != len(tests):
        raise PreparationError(
            "failure ledger test total did not match its canonical test rows"
        )

    evidence: list[dict[str, Any]] = []
    test_units: list[dict[str, Any]] = []
    deterministic_groups_by_signature: dict[str, dict[str, Any]] = {}
    environments: set[str] = set()
    auth_modes: set[str] = set()
    for record in ledger.get("jobs", []):
        environment, auth = _job_context(str(record.get("job_name", "")))
        if environment:
            environments.add(environment)
        if auth:
            auth_modes.add(auth)
    for entry in tests:
        unit_id = _stable_id(
            "test",
            entry["job_id"],
            entry.get("package") or "",
            entry["test"],
        )
        lines = _read_log_lines(logs_dir, str(entry["job_id"]))
        (
            phase,
            unit_evidence,
            evidence_refs,
            machine_facts,
            full_context,
            decision_context,
            single_failure,
        ) = _test_phase_and_evidence(entry, lines, unit_id)
        deterministic_decision = _deterministic_test_decision(
            unit_evidence,
            phase,
            machine_facts,
            str(entry["test"]),
            full_context,
            decision_context,
            single_failure,
            bool(entry.get("subtests")),
        )
        evidence.extend(unit_evidence)
        test_units.append(
            {
                "unit_id": unit_id,
                "kind": "test",
                "job_id": str(entry["job_id"]),
                "job_name": entry["job_name"],
                "package": entry.get("package"),
                "test": entry["test"],
                "subtests": entry.get("subtests", []),
                "phase": phase,
                "evidence_refs": evidence_refs,
                "machine_facts": machine_facts,
                "allowed_categories": _allowed_categories(
                    deterministic_decision,
                    unit_evidence,
                ),
            }
        )
        if deterministic_decision is not None:
            signature = deterministic_decision["signature"]
            group = deterministic_groups_by_signature.setdefault(
                signature,
                {
                    "group_id": _stable_id(
                        "deterministic-group",
                        signature,
                    ),
                    "signature": signature,
                    "category": deterministic_decision["category"],
                    "cause": deterministic_decision["cause"],
                    "unit_ids": [],
                    "evidence_refs": [],
                },
            )
            group["unit_ids"].append(unit_id)
            group["evidence_refs"].append(
                deterministic_decision["evidence_ref"]
            )

    package_teardown_units: list[dict[str, Any]] = []
    for entry in ledger.get("package_teardowns", []):
        unit_id = _stable_id(
            "teardown",
            entry["job_id"],
            entry["package"],
        )
        evidence_refs = []
        lines = _read_log_lines(logs_dir, str(entry["job_id"]))
        for line_number in entry.get("cleanup_line_numbers", []):
            evidence_id = (
                f"line:{unit_id}:{entry['job_id']}:{line_number}"
            )
            evidence_refs.append(evidence_id)
            evidence.append(
                {
                    "evidence_id": evidence_id,
                    "unit_id": unit_id,
                    "kind": "package_teardown",
                    "job_id": str(entry["job_id"]),
                    "line_number": line_number,
                    "text": _line_text(lines, line_number),
                }
            )
        package_teardown_units.append(
            {
                "unit_id": unit_id,
                "kind": "package_teardown",
                "job_id": str(entry["job_id"]),
                "job_name": entry["job_name"],
                "package": entry["package"],
                "evidence_refs": evidence_refs,
                "deterministic_category": "cleanup",
            }
        )

    workflow_units = _workflow_units(ledger)
    package_failure_units, package_failure_evidence = _package_failure_units(
        ledger,
        logs_dir,
    )
    evidence.extend(package_failure_evidence)
    deterministic_test_groups = []
    for group in deterministic_groups_by_signature.values():
        group["unit_ids"].sort()
        group["evidence_refs"].sort()
        group["test_count"] = len(group["unit_ids"])
        deterministic_test_groups.append(group)
    deterministic_test_groups.sort(key=lambda item: item["group_id"])
    deterministic_test_count = sum(
        group["test_count"] for group in deterministic_test_groups
    )
    active_non_reporting = any(
        unit["kind"]
        in {
            "active_test",
            "active_setup",
            "active_cleanup",
            "active_unknown",
        }
        for unit in workflow_units
    )
    active_units = [
        unit for unit in workflow_units if unit["kind"].startswith("active_")
    ]
    only_reporting_active = bool(active_units) and all(
        unit["role"] == "reporting" for unit in active_units
    )
    ledger_complete = summary.get("ledger_complete") is True
    non_reporting_complete = (
        summary.get("non_reporting_logs_complete") is True
    )
    execution_state = str(summary.get("test_execution_state", "unknown"))
    non_reporting_problem_jobs = [
        record
        for record in ledger.get("jobs", [])
        if record.get("role") != "reporting"
        and record.get("conclusion") not in {"success", "skipped"}
    ]
    issues = _structured_issues(ledger)
    reporting_only_active = only_reporting_active
    green_eligible = (
        (run.get("status") == "completed" or reporting_only_active)
        and ledger_complete
        and non_reporting_complete
        and execution_state == "confirmed_ran"
        and summary.get("successful_test_jobs", 0) > 0
        and not test_units
        and not package_failure_units
        and not package_teardown_units
        and not non_reporting_problem_jobs
    )
    run_incomplete = active_non_reporting or (
        run.get("status") != "completed"
        and not only_reporting_active
    )
    if run_incomplete:
        confidence_ceiling = "low"
    elif ledger_complete and non_reporting_complete:
        confidence_ceiling = "high"
    else:
        confidence_ceiling = "medium"
    confidence_reasons = [
        issue["message"] for issue in issues if issue["caps_confidence"]
    ]
    if any(
        unit["machine_facts"]["evidence_scope_complete"] is not True
        for unit in test_units
    ):
        confidence_ceiling = _cap_confidence(confidence_ceiling, "medium")
        confidence_reasons.append(
            "One or more failed tests lacked an unambiguous owned log region."
        )
    if any(
        unit["machine_facts"]["mixed_phase"] is True
        for unit in test_units
    ):
        confidence_ceiling = _cap_confidence(confidence_ceiling, "medium")
        confidence_reasons.append(
            "One or more canonical tests had heterogeneous failure phases."
        )

    analysis: dict[str, Any] = {
        "schema": _schema(ANALYSIS_INPUT_SCHEMA),
        "provenance": {
            "repository": REPOSITORY,
            "run_id": run["run_id"],
            "run_attempt": run["run_attempt"],
            "ledger_schema": ledger.get("schema"),
            "ledger_digest": _canonical_digest(ledger),
        },
        "run_context": {
            **run,
            "commit": str(run["head_sha"])[:7],
            "run_url": (
                f"https://github.com/{REPOSITORY}/actions/runs/{run['run_id']}"
            ),
            "environment": _single_or_mixed(environments),
            "auth": _single_or_mixed(auth_modes),
        },
        "gates": {
            "run_incomplete": run_incomplete,
            "reporting_only_active": reporting_only_active,
            "test_execution_state": execution_state,
            "successful_test_jobs": int(
                summary.get("successful_test_jobs", 0)
            ),
            "ledger_complete": ledger_complete,
            "non_reporting_logs_complete": non_reporting_complete,
            "green_eligible": green_eligible,
            "confidence_ceiling": confidence_ceiling,
            "confidence_reasons": confidence_reasons,
        },
        "unit_counts": {
            "tests": len(test_units),
            "deterministic_tests": deterministic_test_count,
            "model_tests": len(test_units) - deterministic_test_count,
            "deterministic_groups": len(deterministic_test_groups),
            "package_failures": len(package_failure_units),
            "package_teardowns": len(package_teardown_units),
            "workflow_jobs": len(workflow_units),
        },
        "deterministic_test_groups": deterministic_test_groups,
        "test_units": sorted(test_units, key=lambda item: item["unit_id"]),
        "package_failure_units": package_failure_units,
        "package_teardown_units": sorted(
            package_teardown_units,
            key=lambda item: item["unit_id"],
        ),
        "workflow_job_units": workflow_units,
        "issues": issues,
        "evidence": sorted(
            evidence,
            key=lambda item: item["evidence_id"],
        ),
    }
    analysis["analysis_digest"] = _canonical_digest(
        analysis,
        omit_key="analysis_digest",
    )
    return analysis


def build_classification_input(
    analysis_value: dict[str, Any],
) -> dict[str, Any]:
    analysis = _validate_analysis_input(analysis_value)
    deterministic_groups = analysis.get("deterministic_test_groups", [])
    deterministic_ids = {
        unit_id
        for group in deterministic_groups
        for unit_id in group["unit_ids"]
    }
    model_units = [
        unit
        for unit in analysis["test_units"]
        if unit["unit_id"] not in deterministic_ids
    ]
    model_ids = {unit["unit_id"] for unit in model_units}
    model_evidence = [
        item
        for item in analysis["evidence"]
        if item["unit_id"] in model_ids
    ]
    evidence_by_id = {
        item["evidence_id"]: item for item in analysis["evidence"]
    }
    units_by_id = {
        unit["unit_id"]: unit for unit in analysis["test_units"]
    }
    cohort_summaries = []
    for group in deterministic_groups:
        units = [units_by_id[unit_id] for unit_id in group["unit_ids"]]
        known_packages = {
            str(unit["package"])
            for unit in units
            if unit.get("package")
        }
        representative_evidence = []
        for evidence_ref in group["evidence_refs"]:
            text = str(evidence_by_id[evidence_ref]["text"])
            if text not in representative_evidence:
                representative_evidence.append(text)
            if len(representative_evidence) == 2:
                break
        cohort_summaries.append(
            {
                "group_id": group["group_id"],
                "signature": group["signature"],
                "category": group["category"],
                "cause": group["cause"],
                "unique_tests": group["test_count"],
                "jobs": len({unit["job_id"] for unit in units}),
                "packages": len(known_packages),
                "unattributed_tests": sum(
                    not unit.get("package") for unit in units
                ),
                "sample_tests": sorted(
                    str(unit["test"]) for unit in units
                )[:3],
                "representative_evidence": representative_evidence,
                "membership_digest": _canonical_digest(
                    sorted(group["unit_ids"])
                ),
            }
        )

    return {
        "schema": _schema(CLASSIFICATION_INPUT_SCHEMA),
        "analysis_digest": analysis["analysis_digest"],
        "run_context": analysis["run_context"],
        "gates": analysis["gates"],
        "unit_counts": {
            **analysis["unit_counts"],
            "deterministic_tests": len(deterministic_ids),
            "model_tests": len(model_units),
            "deterministic_groups": len(deterministic_groups),
            "model_units": len(model_units),
        },
        "deterministic_cohorts": cohort_summaries,
        "test_units": model_units,
        "issues": analysis["issues"],
        "evidence": model_evidence,
    }


def _short_model_text(
    value: Any,
    label: str,
    *,
    required: bool = False,
    limit: int = 500,
) -> str | None:
    if value is None and not required:
        return None
    if not isinstance(value, str) or not value.strip():
        raise FinalizationError(f"{label} must be a non-empty string")
    text = value.strip()
    if len(text) > limit:
        raise FinalizationError(f"{label} exceeds {limit} characters")
    if "\n" in text or "\r" in text or any(
        ord(character) < 32 and character != "\t" for character in text
    ):
        raise FinalizationError(f"{label} contains a control character")
    return text


def _validate_analysis_input(value: Any) -> dict[str, Any]:
    analysis = _require_schema(
        value,
        ANALYSIS_INPUT_SCHEMA,
        "analysis input",
    )
    digest = analysis.get("analysis_digest")
    if not isinstance(digest, str) or digest != _canonical_digest(
        analysis,
        omit_key="analysis_digest",
    ):
        raise FinalizationError("analysis input digest does not match its content")
    provenance = analysis.get("provenance")
    context = analysis.get("run_context")
    gates = analysis.get("gates")
    counts = analysis.get("unit_counts")
    if not all(
        isinstance(item, dict)
        for item in (provenance, context, gates, counts)
    ):
        raise FinalizationError("analysis input metadata is incomplete")
    if (
        provenance.get("repository") != REPOSITORY
        or provenance.get("ledger_schema") != _schema(FAILURE_LEDGER_SCHEMA)
        or re.fullmatch(
            r"[0-9a-f]{64}",
            str(provenance.get("ledger_digest", "")),
        )
        is None
    ):
        raise FinalizationError("analysis input ledger provenance is invalid")
    if (
        provenance.get("run_id") != context.get("run_id")
        or provenance.get("run_attempt") != context.get("run_attempt")
    ):
        raise FinalizationError("analysis input run provenance is inconsistent")
    required_context_strings = (
        "head_sha",
        "commit",
        "status",
        "run_url",
        "environment",
        "auth",
    )
    if (
        not isinstance(context.get("run_id"), int)
        or context["run_id"] <= 0
        or not isinstance(context.get("run_attempt"), int)
        or context["run_attempt"] <= 0
        or not isinstance(context.get("run_number"), int)
        or context["run_number"] <= 0
        or any(
            not isinstance(context.get(key), str)
            for key in required_context_strings
        )
        or (
            context.get("conclusion") is not None
            and not isinstance(context.get("conclusion"), str)
        )
        or context["commit"] != context["head_sha"][:7]
        or context["run_url"]
        != (
            f"https://github.com/{REPOSITORY}/actions/runs/"
            f"{context['run_id']}"
        )
    ):
        raise FinalizationError("analysis input run context is invalid")
    required_gate_booleans = (
        "run_incomplete",
        "reporting_only_active",
        "ledger_complete",
        "non_reporting_logs_complete",
        "green_eligible",
    )
    if (
        any(not isinstance(gates.get(key), bool) for key in required_gate_booleans)
        or gates.get("test_execution_state")
        not in {"confirmed_ran", "confirmed_none", "unknown"}
        or gates.get("confidence_ceiling") not in {"high", "medium", "low"}
        or not isinstance(gates.get("successful_test_jobs"), int)
        or isinstance(gates.get("successful_test_jobs"), bool)
        or gates["successful_test_jobs"] < 0
        or not isinstance(gates.get("confidence_reasons"), list)
        or any(
            not isinstance(reason, str)
            for reason in gates["confidence_reasons"]
        )
    ):
        raise FinalizationError("analysis input gates are invalid")
    for key in (
        "tests",
        "package_failures",
        "package_teardowns",
        "workflow_jobs",
    ):
        if (
            not isinstance(counts.get(key), int)
            or isinstance(counts.get(key), bool)
            or counts[key] < 0
        ):
            raise FinalizationError("analysis input unit counts are invalid")

    deterministic_groups = analysis.get("deterministic_test_groups", [])
    test_units = analysis.get("test_units")
    package_failure_units = analysis.get("package_failure_units")
    teardown_units = analysis.get("package_teardown_units")
    workflow_units = analysis.get("workflow_job_units")
    evidence = analysis.get("evidence")
    if not all(
        isinstance(item, list)
        for item in (
            deterministic_groups,
            test_units,
            package_failure_units,
            teardown_units,
            workflow_units,
            evidence,
        )
    ):
        raise FinalizationError("analysis input unit collections are invalid")
    if counts.get("tests") != len(test_units):
        raise FinalizationError("analysis input test total is unreconciled")
    if counts.get("package_failures") != len(package_failure_units):
        raise FinalizationError("analysis input package failure total is unreconciled")
    if counts.get("package_teardowns") != len(teardown_units):
        raise FinalizationError("analysis input package teardown total is unreconciled")
    if counts.get("workflow_jobs") != len(workflow_units):
        raise FinalizationError("analysis input workflow-job total is unreconciled")
    deterministic_count_keys = (
        "deterministic_tests",
        "model_tests",
        "deterministic_groups",
    )
    if "deterministic_test_groups" in analysis:
        if any(
            not isinstance(counts.get(key), int)
            or isinstance(counts.get(key), bool)
            or counts[key] < 0
            for key in deterministic_count_keys
        ):
            raise FinalizationError(
                "analysis input deterministic counts are invalid"
            )
    elif any(key in counts for key in deterministic_count_keys):
        raise FinalizationError(
            "analysis input deterministic counts lack their contract"
        )

    all_units = (
        test_units
        + package_failure_units
        + teardown_units
        + workflow_units
    )
    if any(not isinstance(unit, dict) for unit in all_units):
        raise FinalizationError("analysis input contains a non-object unit")
    if any(not isinstance(item, dict) for item in evidence):
        raise FinalizationError("analysis input contains non-object evidence")
    unit_ids = [unit.get("unit_id") for unit in all_units]
    if any(not isinstance(unit_id, str) or not unit_id for unit_id in unit_ids):
        raise FinalizationError("analysis input contains an invalid unit ID")
    if len(set(unit_ids)) != len(unit_ids):
        raise FinalizationError("analysis input contains duplicate unit IDs")

    evidence_ids = [item.get("evidence_id") for item in evidence]
    if any(
        not isinstance(evidence_id, str) or not evidence_id
        for evidence_id in evidence_ids
    ):
        raise FinalizationError("analysis input contains an invalid evidence ID")
    if len(set(evidence_ids)) != len(evidence_ids):
        raise FinalizationError("analysis input contains duplicate evidence IDs")
    known_unit_ids = set(unit_ids)
    if any(item.get("unit_id") not in known_unit_ids for item in evidence):
        raise FinalizationError("analysis input evidence references an unknown unit")
    evidence_owners = {
        item["evidence_id"]: item["unit_id"] for item in evidence
    }
    for item in evidence:
        if (
            not isinstance(item.get("kind"), str)
            or not isinstance(item.get("job_id"), str)
            or (
                "text" in item
                and not isinstance(item.get("text"), str)
            )
        ):
            raise FinalizationError("analysis input evidence is malformed")
    for unit in all_units:
        refs = unit.get("evidence_refs", [])
        if (
            not isinstance(refs, list)
            or any(not isinstance(ref, str) for ref in refs)
            or len(refs) != len(set(refs))
            or any(evidence_owners.get(ref) != unit["unit_id"] for ref in refs)
        ):
            raise FinalizationError(
                "analysis input unit evidence references are invalid"
            )
    for unit in test_units:
        if (
            unit.get("kind") != "test"
            or not isinstance(unit.get("job_id"), str)
            or not isinstance(unit.get("job_name"), str)
            or not isinstance(unit.get("test"), str)
            or not isinstance(unit.get("subtests"), list)
            or any(
                not isinstance(subtest, str)
                for subtest in unit["subtests"]
            )
            or not isinstance(unit.get("phase"), dict)
            or unit["phase"].get("kind")
            not in {
                "test_execution",
                "test_step",
                "post_test_destroy",
                "mixed",
            }
            or not isinstance(unit.get("machine_facts"), dict)
            or any(
                not isinstance(unit["machine_facts"].get(key), bool)
                for key in (
                    "leftover_indicator_observed",
                    "shape_b_candidate",
                    "post_test_destroy",
                    "evidence_scope_complete",
                    "mixed_phase",
                )
            )
            or (
                "pak_authorization_failure_candidate"
                in unit["machine_facts"]
                and not isinstance(
                    unit["machine_facts"][
                        "pak_authorization_failure_candidate"
                    ],
                    bool,
                )
            )
            or (
                "authorization_failure_distinct_requests"
                in unit["machine_facts"]
                and (
                    not isinstance(
                        unit["machine_facts"][
                            "authorization_failure_distinct_requests"
                        ],
                        int,
                    )
                    or isinstance(
                        unit["machine_facts"][
                            "authorization_failure_distinct_requests"
                        ],
                        bool,
                    )
                    or unit["machine_facts"][
                        "authorization_failure_distinct_requests"
                    ]
                    < 0
                )
            )
            or not isinstance(unit.get("allowed_categories"), list)
            or not unit["allowed_categories"]
            or len(unit["allowed_categories"])
            != len(set(unit["allowed_categories"]))
            or any(
                category not in DECISION_CATEGORIES
                for category in unit["allowed_categories"]
            )
            or not unit.get("evidence_refs")
        ):
            raise FinalizationError("analysis input test unit is malformed")
    test_units_by_id = {
        unit["unit_id"]: unit for unit in test_units
    }
    deterministic_ids: set[str] = set()
    allowed_deterministic_group_fields = {
        "group_id",
        "signature",
        "category",
        "cause",
        "test_count",
        "unit_ids",
        "evidence_refs",
    }
    for group in deterministic_groups:
        if not isinstance(group, dict):
            raise FinalizationError(
                "analysis input deterministic group is not an object"
            )
        if set(group) != allowed_deterministic_group_fields:
            raise FinalizationError(
                "analysis input deterministic group fields are invalid"
            )
        signature = group.get("signature")
        policy = DETERMINISTIC_TEST_POLICIES.get(str(signature))
        unit_ids_for_group = group.get("unit_ids")
        evidence_refs_for_group = group.get("evidence_refs")
        if (
            policy is None
            or group.get("group_id")
            != _stable_id("deterministic-group", signature)
            or group.get("category") != policy["category"]
            or group.get("cause") != policy["cause"]
            or not isinstance(unit_ids_for_group, list)
            or not unit_ids_for_group
            or any(
                not isinstance(unit_id, str)
                for unit_id in unit_ids_for_group
            )
            or len(unit_ids_for_group) != len(set(unit_ids_for_group))
            or unit_ids_for_group != sorted(unit_ids_for_group)
            or not isinstance(group.get("test_count"), int)
            or isinstance(group.get("test_count"), bool)
            or group["test_count"] != len(unit_ids_for_group)
            or not isinstance(evidence_refs_for_group, list)
            or len(evidence_refs_for_group) != len(unit_ids_for_group)
            or any(
                not isinstance(evidence_ref, str)
                for evidence_ref in evidence_refs_for_group
            )
            or len(evidence_refs_for_group)
            != len(set(evidence_refs_for_group))
            or evidence_refs_for_group != sorted(evidence_refs_for_group)
        ):
            raise FinalizationError(
                "analysis input deterministic group is malformed"
            )
        unknown_units = set(unit_ids_for_group) - set(test_units_by_id)
        if unknown_units:
            raise FinalizationError(
                "analysis input deterministic group references unknown tests"
            )
        overlap = deterministic_ids & set(unit_ids_for_group)
        if overlap:
            raise FinalizationError(
                "analysis input deterministic groups overlap"
            )
        if any(
            group["category"]
            not in test_units_by_id[unit_id]["allowed_categories"]
            for unit_id in unit_ids_for_group
        ):
            raise FinalizationError(
                "analysis input deterministic category is not allowed"
            )
        evidence_owners_for_group = {
            evidence_owners.get(evidence_ref)
            for evidence_ref in evidence_refs_for_group
        }
        if (
            None in evidence_owners_for_group
            or evidence_owners_for_group != set(unit_ids_for_group)
        ):
            raise FinalizationError(
                "analysis input deterministic evidence ownership is invalid"
            )
        deterministic_ids.update(unit_ids_for_group)
    if "deterministic_test_groups" in analysis:
        if (
            counts["deterministic_tests"] != len(deterministic_ids)
            or counts["model_tests"]
            != len(test_units) - len(deterministic_ids)
            or counts["deterministic_groups"] != len(deterministic_groups)
        ):
            raise FinalizationError(
                "analysis input deterministic counts are unreconciled"
            )
    for unit in package_failure_units:
        if (
            unit.get("kind") != "package_failure"
            or not isinstance(unit.get("job_id"), str)
            or not isinstance(unit.get("job_name"), str)
            or not isinstance(unit.get("package"), str)
            or unit.get("deterministic_category")
            not in {"code_regression", "unresolved"}
            or not isinstance(unit.get("reason"), str)
            or not unit.get("evidence_refs")
        ):
            raise FinalizationError("analysis input package failure unit is malformed")
    for unit in teardown_units:
        if (
            unit.get("kind") != "package_teardown"
            or not isinstance(unit.get("job_id"), str)
            or not isinstance(unit.get("job_name"), str)
            or not isinstance(unit.get("package"), str)
            or unit.get("deterministic_category") != "cleanup"
        ):
            raise FinalizationError("analysis input teardown unit is malformed")
    valid_roles = {"test", "setup", "cleanup", "reporting", "unknown"}
    for unit in workflow_units:
        if (
            not isinstance(unit.get("kind"), str)
            or not isinstance(unit.get("job_id"), str)
            or not isinstance(unit.get("job_name"), str)
            or unit.get("role") not in valid_roles
            or not isinstance(unit.get("conclusion"), str)
            or not isinstance(unit.get("sources"), list)
        ):
            raise FinalizationError("analysis input workflow unit is malformed")
    active_workflow_units = [
        unit
        for unit in workflow_units
        if unit["kind"].startswith("active_")
    ]
    reporting_only_active = bool(active_workflow_units) and all(
        unit["role"] == "reporting" for unit in active_workflow_units
    )
    if (
        gates["reporting_only_active"] != reporting_only_active
        or (
            context.get("conclusion") is None
            and context["status"] == "completed"
        )
        or (
            gates["reporting_only_active"]
            and (
                context["status"] == "completed"
                or context.get("conclusion") is not None
                or gates["run_incomplete"]
            )
        )
    ):
        raise FinalizationError(
            "analysis input live-reporting state is inconsistent"
        )
    if gates["green_eligible"] and (
        gates["run_incomplete"]
        or (
            context["status"] != "completed"
            and not gates["reporting_only_active"]
        )
        or gates["test_execution_state"] != "confirmed_ran"
        or gates["successful_test_jobs"] < 1
        or not gates["ledger_complete"]
        or not gates["non_reporting_logs_complete"]
        or test_units
        or package_failure_units
        or teardown_units
        or any(
            unit["role"] != "reporting"
            and unit["conclusion"] not in {"success", "skipped"}
            for unit in workflow_units
        )
    ):
        raise FinalizationError("analysis input green eligibility is inconsistent")
    return analysis


def _validate_model_decisions(
    analysis: dict[str, Any],
    value: Any,
) -> tuple[list[dict[str, Any]], dict[str, str | None]]:
    decisions = _require_schema(
        value,
        MODEL_DECISIONS_SCHEMA,
        "model decisions",
    )
    allowed_top_level = {
        "schema",
        "analysis_digest",
        "groups",
        "why",
        "action",
        "tldr",
    }
    unexpected = set(decisions) - allowed_top_level
    if unexpected:
        raise FinalizationError(
            "model decisions contain unsupported fields: "
            + ", ".join(sorted(unexpected))
        )
    if decisions.get("analysis_digest") != analysis["analysis_digest"]:
        raise FinalizationError(
            "model decisions do not match this analysis input digest"
        )
    groups = decisions.get("groups")
    if not isinstance(groups, list):
        raise FinalizationError("model decision groups must be a list")
    if len(groups) > MAX_MODEL_DECISION_GROUPS:
        raise FinalizationError("model decision groups exceed the limit")

    all_test_units = {
        unit["unit_id"]: unit for unit in analysis["test_units"]
    }
    deterministic_groups = analysis.get("deterministic_test_groups", [])
    deterministic_ids = {
        unit_id
        for group in deterministic_groups
        for unit_id in group["unit_ids"]
    }
    test_units = {
        unit_id: unit
        for unit_id, unit in all_test_units.items()
        if unit_id not in deterministic_ids
    }
    evidence_owners = {
        item["evidence_id"]: item["unit_id"] for item in analysis["evidence"]
    }
    explicit_ids: set[str] = set()
    remaining_group: dict[str, Any] | None = None
    validated_groups: list[dict[str, Any]] = []
    allowed_group_fields = {
        "unit_ids",
        "remaining",
        "category",
        "cause",
        "evidence_refs",
        "ambiguity",
        "note",
    }

    for index, group in enumerate(groups):
        if not isinstance(group, dict):
            raise FinalizationError(
                f"model decision group {index} must be an object"
            )
        group_unexpected = set(group) - allowed_group_fields
        if group_unexpected:
            raise FinalizationError(
                f"model decision group {index} contains unsupported fields: "
                + ", ".join(sorted(group_unexpected))
            )
        has_ids = "unit_ids" in group
        has_remaining = group.get("remaining") is True
        if has_ids == has_remaining:
            raise FinalizationError(
                f"model decision group {index} must have either unit_ids "
                "or remaining: true"
            )
        if "remaining" in group and group.get("remaining") is not True:
            raise FinalizationError(
                f"model decision group {index} has an invalid remaining selector"
            )
        if has_remaining:
            if remaining_group is not None:
                raise FinalizationError(
                    "model decisions may contain at most one remaining selector"
                )
            if index != len(groups) - 1:
                raise FinalizationError(
                    "the remaining selector must be the last decision group"
                )
            remaining_group = group
            resolved_ids: list[str] = []
        else:
            unit_ids = group["unit_ids"]
            if (
                not isinstance(unit_ids, list)
                or not unit_ids
                or len(unit_ids) > MAX_MODEL_DECISION_UNIT_IDS
                or any(not isinstance(unit_id, str) for unit_id in unit_ids)
            ):
                raise FinalizationError(
                    f"model decision group {index} has invalid unit IDs"
                )
            if len(set(unit_ids)) != len(unit_ids):
                raise FinalizationError(
                    f"model decision group {index} repeats a unit ID"
                )
            deterministic_references = set(unit_ids) & deterministic_ids
            if deterministic_references:
                raise FinalizationError(
                    "model decisions must not reclassify deterministic tests: "
                    + ", ".join(sorted(deterministic_references))
                )
            unknown_ids = set(unit_ids) - set(test_units)
            if unknown_ids:
                raise FinalizationError(
                    "model decisions reference unknown test units: "
                    + ", ".join(sorted(unknown_ids))
                )
            duplicates = explicit_ids & set(unit_ids)
            if duplicates:
                raise FinalizationError(
                    "model decisions classify test units more than once: "
                    + ", ".join(sorted(duplicates))
                )
            explicit_ids.update(unit_ids)
            resolved_ids = unit_ids.copy()

        category = group.get("category")
        if category not in DECISION_CATEGORIES:
            raise FinalizationError(
                f"model decision group {index} has an invalid category"
            )
        cause = _short_model_text(
            group.get("cause"),
            f"model decision group {index} cause",
            required=True,
            limit=MAX_MODEL_DECISION_CAUSE_CHARS,
        )
        evidence_refs = group.get("evidence_refs")
        if (
            not isinstance(evidence_refs, list)
            or not evidence_refs
            or len(evidence_refs) > MAX_MODEL_DECISION_EVIDENCE_REFS
            or any(not isinstance(ref, str) for ref in evidence_refs)
            or len(set(evidence_refs)) != len(evidence_refs)
        ):
            raise FinalizationError(
                f"model decision group {index} has invalid evidence refs"
            )
        ambiguity = _short_model_text(
            group.get("ambiguity"),
            f"model decision group {index} ambiguity",
            limit=MAX_MODEL_DECISION_AMBIGUITY_CHARS,
        )
        note = _short_model_text(
            group.get("note"),
            f"model decision group {index} note",
            limit=MAX_MODEL_DECISION_NOTE_CHARS,
        )
        validated_groups.append(
            {
                "unit_ids": resolved_ids,
                "remaining": has_remaining,
                "category": category,
                "cause": cause,
                "evidence_refs": evidence_refs,
                "ambiguity": ambiguity,
                "note": note,
            }
        )

    if remaining_group is not None:
        remaining_ids = sorted(set(test_units) - explicit_ids)
        if not remaining_ids:
            raise FinalizationError(
                "remaining selector did not match any test unit"
            )
        for group in validated_groups:
            if group["remaining"]:
                group["unit_ids"] = remaining_ids
                break
    elif explicit_ids != set(test_units):
        missing = set(test_units) - explicit_ids
        raise FinalizationError(
            "model decisions did not classify every test unit: "
            + ", ".join(sorted(missing))
        )

    assigned_ids = [
        unit_id
        for group in validated_groups
        for unit_id in group["unit_ids"]
    ]
    if len(assigned_ids) != len(set(assigned_ids)):
        raise FinalizationError(
            "model decision groups overlap after resolving remaining"
        )
    if set(assigned_ids) != set(test_units):
        raise FinalizationError(
            "model decision groups do not reconcile to the canonical test total"
        )

    for index, group in enumerate(validated_groups):
        disallowed_units = [
            unit_id
            for unit_id in group["unit_ids"]
            if group["category"]
            not in test_units[unit_id].get("allowed_categories", [])
        ]
        if disallowed_units:
            raise FinalizationError(
                f"model decision group {index} uses a category that is not "
                "allowed for: "
                + ", ".join(sorted(disallowed_units))
            )
        allowed_evidence = {
            ref
            for ref, owner in evidence_owners.items()
            if owner in group["unit_ids"]
        }
        unknown_refs = set(group["evidence_refs"]) - allowed_evidence
        if unknown_refs:
            raise FinalizationError(
                f"model decision group {index} references evidence not owned "
                "by its test units: "
                + ", ".join(sorted(unknown_refs))
            )
        shape_b_units = [
            unit_id
            for unit_id in group["unit_ids"]
            if test_units[unit_id]
            .get("machine_facts", {})
            .get("shape_b_candidate")
            is True
        ]
        if (
            group["category"] == "timeout"
            and shape_b_units
            and (
                group["note"] != "delete_on_timeout_unverified"
                or group["ambiguity"] is None
            )
        ):
            raise FinalizationError(
                "Shape-B timeout decisions require ambiguity and the "
                "delete_on_timeout_unverified note"
            )
        pak_authorization_units = [
            unit_id
            for unit_id in group["unit_ids"]
            if test_units[unit_id]
            .get("machine_facts", {})
            .get("pak_authorization_failure_candidate")
            is True
        ]
        if pak_authorization_units and group["ambiguity"] is None:
            raise FinalizationError(
                "PAK authorization-regression candidates require an "
                "ambiguity explanation"
            )
        if group["category"] == "unresolved" and group["ambiguity"] is None:
            raise FinalizationError(
                "unresolved decisions require an ambiguity explanation"
            )
        if group["note"] == "delete_on_timeout_unverified":
            if group["category"] != "timeout":
                raise FinalizationError(
                    "delete-on-timeout note requires the timeout category"
                )
            if not all(
                test_units[unit_id]
                .get("machine_facts", {})
                .get("shape_b_candidate")
                is True
                for unit_id in group["unit_ids"]
            ):
                raise FinalizationError(
                    "delete-on-timeout note is not supported by every selected unit"
                )
        elif group["note"] is not None:
            raise FinalizationError(
                f"model decision group {index} has an unsupported note"
            )

    validated_groups.extend(
        {
            "unit_ids": group["unit_ids"].copy(),
            "remaining": False,
            "category": group["category"],
            "cause": group["cause"],
            "evidence_refs": group["evidence_refs"].copy(),
            "ambiguity": None,
            "note": None,
        }
        for group in deterministic_groups
    )
    optional = {
        key: _short_model_text(
            decisions.get(key),
            f"model decisions {key}",
            limit=MAX_MODEL_DECISION_OPTIONAL_TEXT_CHARS,
        )
        for key in ("why", "action", "tldr")
    }
    return validated_groups, optional


def _confidence_rank(value: str) -> int:
    return {"low": 0, "medium": 1, "high": 2}.get(value, 0)


def _cap_confidence(value: str, cap: str) -> str:
    return value if _confidence_rank(value) <= _confidence_rank(cap) else cap


def _slack_text(value: str, limit: int | None = None) -> str:
    text = re.sub(r"\b(?:CLOUDP|HELP)-\d+\b", "internal ticket", value)
    text = text.replace("`", "'").replace("*", "")
    text = text.replace("<", "(").replace(">", ")")
    text = re.sub(r"\s+", " ", text).strip()
    if limit is not None and len(text) > limit:
        if limit <= 0:
            return ""
        if limit == 1:
            return "…"
        boundary = text.rfind(" ", 0, limit)
        cut = boundary if boundary > 0 else limit - 1
        return text[:cut].rstrip() + "…"
    return text


def _format_names(names: list[str], limit: int) -> str:
    shown = [f"`{_slack_text(name, 100)}`" for name in names[:limit]]
    if len(names) > limit:
        shown.append(f"and {len(names) - limit} more")
    return ", ".join(shown)


def _plural(count: int, singular: str, plural: str | None = None) -> str:
    return singular if count == 1 else (plural or f"{singular}s")


def _prioritized_causes(
    groups: list[dict[str, Any]],
    *,
    limit: int,
    text_limit: int,
) -> list[str]:
    causes: dict[str, dict[str, int]] = {}
    for index, group in enumerate(groups):
        cause = _slack_text(str(group["cause"]), text_limit)
        entry = causes.setdefault(
            cause,
            {"first_index": index, "test_count": 0},
        )
        entry["test_count"] += len(group["unit_ids"])
    if len(causes) <= limit:
        return list(causes)
    selected = sorted(
        causes,
        key=lambda cause: (
            -causes[cause]["test_count"],
            causes[cause]["first_index"],
        ),
    )[:limit]
    return sorted(
        selected,
        key=lambda cause: causes[cause]["first_index"],
    )


def _slack_bytes(value: str) -> int:
    return len(f"{value}\n".encode())


def _display_package(package: str) -> str:
    return package.removeprefix(f"github.com/{REPOSITORY}/")


def _urgency(verdict: str, confidence: str, requires_review: bool) -> str:
    if verdict == "red":
        return "investigate immediately"
    if verdict == "green":
        return {
            "high": "no on-call action",
            "medium": "please verify the run",
            "low": "please review the run manually",
        }[confidence]
    if confidence == "low":
        return "please review the run manually"
    if confidence == "medium" or requires_review:
        return "non-urgent review recommended"
    return "no immediate on-call action needed"


def _derive_verdict(
    analysis: dict[str, Any],
    category_counts: dict[str, int],
    unresolved_tests: int,
) -> tuple[str, str, str]:
    gates = analysis["gates"]
    execution_state = gates["test_execution_state"]
    setup_failures = [
        unit
        for unit in analysis["workflow_job_units"]
        if unit["role"] == "setup"
        and unit["kind"] != "active_setup"
    ]
    package_failures = analysis["package_failure_units"]
    if gates["run_incomplete"]:
        return "yellow", "Run incomplete", "run_incomplete"
    if (
        category_counts["code_regression"] > 0
        or any(
            unit["deterministic_category"] == "code_regression"
            for unit in package_failures
        )
        or (setup_failures and execution_state == "confirmed_ran")
    ):
        return "red", "CODE REGRESSION DETECTED", "code_regression"
    if execution_state != "confirmed_ran":
        return "red", "SUITE UNVERIFIED", "suite_unverified"
    if gates["green_eligible"]:
        return "green", "All tests passed", "green"
    if (
        unresolved_tests
        or any(
            unit["deterministic_category"] == "unresolved"
            for unit in package_failures
        )
        or not gates["ledger_complete"]
        or not gates["non_reporting_logs_complete"]
        or any(
            unit["role"] == "unknown"
            or unit["kind"] in {"test_job", "active_unknown"}
            for unit in analysis["workflow_job_units"]
        )
    ):
        return "yellow", "Automatic triage incomplete", "review"
    return "yellow", "Infrastructure noise only", "infrastructure"


def _render_summary(
    analysis: dict[str, Any],
    groups: list[dict[str, Any]],
    optional: dict[str, str | None],
    *,
    compact: bool = False,
) -> dict[str, Any]:
    test_units = {
        unit["unit_id"]: unit for unit in analysis["test_units"]
    }
    category_counts = {category: 0 for category in MODEL_CATEGORIES}
    by_category: dict[str, list[dict[str, Any]]] = {
        category: [] for category in DECISION_CATEGORIES
    }
    unresolved_tests = 0
    for group in groups:
        if group["category"] == "unresolved":
            unresolved_tests += len(group["unit_ids"])
        else:
            category_counts[group["category"]] += len(group["unit_ids"])
        by_category[group["category"]].append(group)
    if (
        sum(category_counts.values()) + unresolved_tests
        != analysis["unit_counts"]["tests"]
    ):
        raise FinalizationError(
            "derived category totals do not reconcile to the canonical test total"
        )

    verdict, headline, disposition = _derive_verdict(
        analysis,
        category_counts,
        unresolved_tests,
    )
    confidence = str(analysis["gates"]["confidence_ceiling"])
    if confidence not in {"high", "medium", "low"}:
        raise FinalizationError("analysis input has an invalid confidence ceiling")
    shape_b_groups = [
        group
        for group in groups
        if group["note"] == "delete_on_timeout_unverified"
    ]
    if shape_b_groups or any(group["ambiguity"] for group in groups):
        confidence = _cap_confidence(confidence, "medium")
    if disposition == "review":
        confidence = _cap_confidence(confidence, "medium")
    if disposition == "run_incomplete":
        confidence = "low"
    requires_review = bool(
        optional["action"] or shape_b_groups or disposition == "review"
    )
    urgency = _urgency(verdict, confidence, requires_review)

    context = analysis["run_context"]
    emoji = {
        "red": ":red_circle:",
        "yellow": ":yellow_circle:",
        "green": ":green_circle:",
    }[verdict]
    lines = [
        (
            f"{emoji} *Test Suite #{context['run_number']} — {headline}* "
            f"(`{_slack_text(str(context['commit']))}` on "
            f"`{_slack_text(str(context['environment']))}`, "
            f"`{_slack_text(str(context['auth']))}`)"
        ),
        (
            f"{confidence} confidence — {urgency} — "
            f"<{context['run_url']}|view run>"
        ),
    ]

    test_total = analysis["unit_counts"]["tests"]
    package_failure_total = analysis["unit_counts"]["package_failures"]
    teardown_total = analysis["unit_counts"]["package_teardowns"]
    workflow_units = analysis["workflow_job_units"]
    if disposition == "run_incomplete":
        lines.extend(
            [
                "",
                "The run still has active non-reporting jobs; retry after it completes.",
            ]
        )
    elif disposition == "suite_unverified":
        lines.extend(
            [
                "",
                "Provider behaviour was not verified because no completed test execution was confirmed.",
            ]
        )
    elif disposition == "review":
        if test_total == 0:
            lines.extend(
                [
                    "",
                    (
                        "No canonical failed-test identities were recoverable; "
                        "job or package evidence requires review."
                    ),
                ]
            )
        elif analysis["gates"]["ledger_complete"]:
            detail = (
                f"{unresolved_tests} "
                f"{_plural(unresolved_tests, 'test')} could not be "
                "classified automatically."
                if unresolved_tests
                else "Automatic triage could not fully classify the "
                "available failure evidence."
            )
            lines.extend(
                [
                    "",
                    (
                        f"{test_total} unique top-level "
                        f"{_plural(test_total, 'test')} failed; {detail}"
                    ),
                ]
            )
        else:
            lines.extend(
                [
                    "",
                    (
                        f"{test_total} unique top-level test failures observed; "
                        "the run total is unverified."
                    ),
                ]
            )
    elif test_total:
        lines.extend(
            [
                "",
                (
                    f"{test_total} unique top-level "
                    f"{_plural(test_total, 'test')} failed:"
                ),
            ]
        )

    if verdict == "red" and category_counts["code_regression"]:
        regression_total = category_counts["code_regression"]
        lines.extend(
            [
                "",
                (
                    f"*Code {_plural(regression_total, 'regression')} "
                    f"({regression_total} "
                    f"{_plural(regression_total, 'test')})*:"
                ),
            ]
        )
        regression_groups = by_category["code_regression"]
        for group in regression_groups[: (2 if compact else 5)]:
            names = [
                str(test_units[unit_id]["test"])
                for unit_id in group["unit_ids"]
            ]
            cause = _slack_text(
                str(group["cause"]),
                120 if compact else MAX_MODEL_DECISION_CAUSE_CHARS,
            )
            if len(names) == 1:
                lines.append(f"• `{_slack_text(names[0], 100)}` — {cause}")
            else:
                lines.append(
                    f"• {len(names)} {_plural(len(names), 'test')} "
                    f"({_format_names(names, 3)}) — {cause}"
                )
        if len(regression_groups) > (2 if compact else 5):
            omitted = sum(
                len(group["unit_ids"])
                for group in regression_groups[(2 if compact else 5) :]
            )
            lines.append(f"• and {omitted} more code-regression tests")

    other_categories = [
        category
        for category in MODEL_CATEGORIES
        if category != "code_regression" and category_counts[category]
    ]
    if other_categories:
        lines.append("")
        if compact:
            values = ", ".join(
                f"{CATEGORY_LABELS[category]} {category_counts[category]}"
                for category in other_categories
            )
            lines.append(f"*Other failures*: {values}")
        else:
            lines.append("*Other failures*:")
            for category in other_categories:
                category_groups = by_category[category]
                causes = _prioritized_causes(
                    category_groups,
                    limit=2,
                    text_limit=MAX_MODEL_DECISION_CAUSE_CHARS,
                )
                cause_suffix = (
                    f" — {'; '.join(causes)}" if causes else ""
                )
                lines.append(
                    f"• {CATEGORY_LABELS[category]}: "
                    f"{category_counts[category]} "
                    f"{_plural(category_counts[category], 'test')}"
                    f"{cause_suffix}"
                )

    if unresolved_tests:
        unresolved_causes = _prioritized_causes(
            by_category["unresolved"],
            limit=2,
            text_limit=(
                180 if compact else MAX_MODEL_DECISION_CAUSE_CHARS
            ),
        )
        suffix = (
            f" — {'; '.join(unresolved_causes)}"
            if unresolved_causes
            else ""
        )
        lines.append(
            f"• Unclassified: {unresolved_tests} "
            f"{_plural(unresolved_tests, 'test')}{suffix}"
        )

    shape_b_ids = [
        unit_id
        for group in shape_b_groups
        for unit_id in group["unit_ids"]
    ]
    if shape_b_ids:
        names = [str(test_units[unit_id]["test"]) for unit_id in shape_b_ids]
        causes = []
        for group in shape_b_groups:
            cause = _slack_text(
                str(group["cause"]),
                140 if compact else MAX_MODEL_DECISION_CAUSE_CHARS,
            )
            if cause not in causes:
                causes.append(cause)
        name_suffix = (
            ""
            if compact
            else f" — {_format_names(names, 3)}"
        )
        cause_suffix = f" — {causes[0]}" if causes else ""
        lines.append(
            "• Delete-on-timeout unverified: "
            f"{len(shape_b_ids)} {_plural(len(shape_b_ids), 'test')}"
            f"{name_suffix}{cause_suffix} (included in Timeout total)"
        )

    for unit in analysis["package_failure_units"][:3]:
        category = (
            "build failure"
            if unit["deterministic_category"] == "code_regression"
            else "unresolved package failure"
        )
        lines.append(
            f"• Package {category}: "
            f"`{_slack_text(_display_package(str(unit['package'])), 120)}`"
        )
    if package_failure_total > 3:
        lines.append(f"• and {package_failure_total - 3} more package failures")

    if teardown_total:
        packages = sorted(
            {
                _display_package(str(unit["package"]))
                for unit in analysis["package_teardown_units"]
            }
        )
        package_suffix = (
            f" — {_format_names(packages, 3)}" if not compact else ""
        )
        lines.append(
            f"• Package teardowns: {teardown_total}{package_suffix}"
        )
    cleanup_units = [
        unit for unit in workflow_units if unit["role"] == "cleanup"
    ]
    if cleanup_units:
        lines.append(f"• Cleanup jobs: {len(cleanup_units)}")
    setup_units = [
        unit
        for unit in workflow_units
        if unit["role"] == "setup" and not unit["kind"].startswith("active_")
    ]
    for unit in setup_units[:3]:
        lines.append(
            f"• Setup failure in `{_slack_text(unit['job_name'], 120)}` "
            f"— {unit['conclusion']}; affected tests were not verified"
        )
    reporting_units = [
        unit
        for unit in workflow_units
        if unit["role"] == "reporting"
        and not unit["kind"].startswith("active_")
    ]
    if reporting_units:
        lines.append(
            f"• Reporting jobs: {len(reporting_units)} failed "
            "(summary delivery only; excluded from test verdict)"
        )
    unknown_units = [
        unit for unit in workflow_units if unit["role"] == "unknown"
    ]
    for unit in unknown_units[:3]:
        lines.append(
            f"• Unclassified job: `{_slack_text(unit['job_name'], 120)}`"
        )
    test_job_units = [
        unit for unit in workflow_units if unit["kind"] == "test_job"
    ]
    for unit in test_job_units[:3]:
        reasons = unit.get("reasons") or ["test evidence was incomplete"]
        duration = (
            f", {unit['duration_minutes']}m"
            if unit.get("duration_minutes") is not None
            else ""
        )
        lines.append(
            f"• Unverified test job: "
            f"`{_slack_text(unit['job_name'], 120)}` — "
            f"{_slack_text(str(unit['conclusion']), 30)}{duration}; "
            f"{_slack_text(str(reasons[0]), 140)}"
        )
    if len(test_job_units) > 3:
        lines.append(f"• and {len(test_job_units) - 3} more unverified test jobs")

    if test_total and disposition not in {"run_incomplete", "green"}:
        names = sorted(str(unit["test"]) for unit in test_units.values())
        lines.extend(
            [
                "",
                f"*Failing tests*: {_format_names(names, 5 if compact else FAILING_TEST_LIMIT)}",
            ]
        )
    if verdict == "red" and optional["tldr"]:
        lines.extend(["", f"*TL;DR*: {_slack_text(optional['tldr'], 300)}"])
    if optional["action"]:
        lines.extend(["", f"*Action*: {_slack_text(optional['action'], 300)}"])
    if confidence != "high":
        why = optional["why"]
        if not why:
            ambiguities = [
                group["ambiguity"]
                for group in groups
                if group["ambiguity"]
            ]
            if shape_b_groups:
                why = (
                    "A delete wait expired without direct evidence distinguishing "
                    "provider cleanup failure from slow Atlas deletion."
                )
            elif ambiguities:
                why = str(ambiguities[0])
            elif analysis["gates"]["confidence_reasons"]:
                why = str(analysis["gates"]["confidence_reasons"][0])
            elif disposition == "run_incomplete":
                why = "The workflow is still running."
            else:
                why = "The available evidence does not fully verify every failure."
        lines.extend(
            [
                "",
                f"*Why {confidence} confidence*: {_slack_text(why, 300)}",
            ]
        )

    slack = "\n".join(lines).strip()
    return {
        "schema": _schema(FINALIZATION_RESULT_SCHEMA),
        "analysis_digest": analysis["analysis_digest"],
        "run_id": context["run_id"],
        "run_attempt": context["run_attempt"],
        "run_number": context["run_number"],
        "verdict": verdict,
        "headline": headline,
        "disposition": disposition,
        "confidence": confidence,
        "urgency": urgency,
        "test_total": test_total,
        "package_failures": package_failure_total,
        "category_counts": category_counts,
        "unresolved_tests": unresolved_tests,
        "package_teardowns": teardown_total,
        "workflow_jobs": len(workflow_units),
        "slack_mrkdwn": slack,
    }


def finalize_analysis(
    analysis_value: Any,
    decisions_value: Any,
) -> dict[str, Any]:
    analysis = _validate_analysis_input(analysis_value)
    groups, optional = _validate_model_decisions(analysis, decisions_value)
    result = _render_summary(analysis, groups, optional)
    if _slack_bytes(result["slack_mrkdwn"]) > SLACK_HARD_LIMIT:
        result = _render_summary(
            analysis,
            groups,
            optional,
            compact=True,
        )
    if _slack_bytes(result["slack_mrkdwn"]) > SLACK_HARD_LIMIT:
        context = analysis["run_context"]
        fallback_lines = [
            (
                f"{':red_circle:' if result['verdict'] == 'red' else ':yellow_circle:'} "
                f"*Test Suite #{context['run_number']} — {result['headline']}* "
                f"(`{_slack_text(str(context['commit']))}`)"
            ),
            (
                f"{result['confidence']} confidence — "
                f"{result['urgency']} — <{context['run_url']}|view run>"
            ),
            (
                f"{result['test_total']} failed tests; "
                "deterministic category counts are available in the "
                "finalization result."
            ),
        ]
        shape_b_count = sum(
            len(group["unit_ids"])
            for group in groups
            if group["note"] == "delete_on_timeout_unverified"
        )
        if shape_b_count:
            fallback_lines.append(
                f"Delete-on-timeout unverified: {shape_b_count} "
                f"{_plural(shape_b_count, 'test')} "
                "(included in Timeout total)."
            )
        if optional["action"]:
            fallback_lines.append(
                f"*Action*: {_slack_text(optional['action'], 300)}"
            )
        result["slack_mrkdwn"] = "\n".join(fallback_lines)
    if _slack_bytes(result["slack_mrkdwn"]) > SLACK_HARD_LIMIT:
        raise FinalizationError("deterministic Slack summary exceeded the hard limit")
    return result


def _capture_github(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
    command = ["gh", "api", endpoint]
    if paginate:
        command.extend(["--paginate", "--slurp"])
    try:
        return subprocess.run(
            command,
            capture_output=True,
            check=False,
            timeout=GITHUB_METADATA_TIMEOUT_SECONDS,
        )
    except subprocess.TimeoutExpired as exc:
        raise PreparationError(
            "GitHub metadata request timed out after "
            f"{GITHUB_METADATA_TIMEOUT_SECONDS} seconds"
        ) from exc


def _download_github(endpoint: str, log_path: Path, error_path: Path) -> int:
    with log_path.open("wb") as log_handle, error_path.open("wb") as error_handle:
        try:
            result = subprocess.run(
                ["gh", "api", endpoint],
                stdout=log_handle,
                stderr=error_handle,
                check=False,
                timeout=GITHUB_LOG_TIMEOUT_SECONDS,
            )
        except subprocess.TimeoutExpired as exc:
            raise PreparationError(
                "GitHub log request timed out after "
                f"{GITHUB_LOG_TIMEOUT_SECONDS} seconds"
            ) from exc
    return result.returncode


def _github_json(
    result: subprocess.CompletedProcess[bytes],
    label: str,
) -> Any:
    if result.returncode != 0:
        raise PreparationError(
            f"{label} request failed with exit code {result.returncode}"
        )
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError as exc:
        raise PreparationError(f"{label} response was not valid JSON") from exc


def _run_metadata(value: Any, run_id: int) -> dict[str, Any]:
    if not isinstance(value, dict) or value.get("id") != run_id:
        raise PreparationError("run metadata did not match the requested run ID")
    if not isinstance(value.get("head_sha"), str) or not value["head_sha"]:
        raise PreparationError("run metadata had no head SHA")
    if not isinstance(value.get("run_number"), int):
        raise PreparationError("run metadata had no numeric run number")
    if type(value.get("run_attempt")) is not int or value["run_attempt"] <= 0:
        raise PreparationError("run metadata had no valid run attempt")
    if not isinstance(value.get("status"), str):
        raise PreparationError("run metadata had no status")
    if "conclusion" not in value:
        raise PreparationError("run metadata had no conclusion field")
    if value["conclusion"] is not None and not isinstance(value["conclusion"], str):
        raise PreparationError("run metadata had an invalid conclusion")
    return {
        "run_id": run_id,
        "head_sha": value["head_sha"],
        "run_number": value["run_number"],
        "run_attempt": value["run_attempt"],
        "status": value["status"],
        "conclusion": value["conclusion"],
    }


def _job_metadata(
    value: Any,
    run_id: int,
    run_attempt: int,
) -> list[dict[str, Any]]:
    if not isinstance(value, list):
        raise PreparationError("paginated jobs response was not a list")

    jobs: list[dict[str, Any]] = []
    seen_job_ids: set[str] = set()
    expected_total: int | None = None
    for page in value:
        if not isinstance(page, dict) or not isinstance(page.get("jobs"), list):
            raise PreparationError("paginated jobs response had an invalid page")
        page_total = page.get("total_count")
        if type(page_total) is not int or page_total < 0:
            raise PreparationError(
                "paginated jobs response had no valid total count"
            )
        if expected_total is None:
            expected_total = page_total
        elif page_total != expected_total:
            raise PreparationError(
                "paginated jobs response had inconsistent total counts"
            )
        for job in page["jobs"]:
            if not isinstance(job, dict):
                raise PreparationError("jobs response contained a non-object job")
            job_id = job.get("id")
            name = job.get("name")
            conclusion = job.get("conclusion")
            if not isinstance(job_id, int) or not isinstance(name, str) or not name:
                raise PreparationError("jobs response contained an invalid job")
            if job.get("run_id") != run_id:
                raise PreparationError(
                    f"job {job_id} did not match requested run {run_id}"
                )
            if job.get("run_attempt") != run_attempt:
                raise PreparationError(
                    f"job {job_id} did not match run attempt {run_attempt}"
                )
            if conclusion is not None and not isinstance(conclusion, str):
                raise PreparationError("jobs response contained an invalid conclusion")
            job_id_text = str(job_id)
            if job_id_text in seen_job_ids:
                raise PreparationError(f"jobs response repeated job ID {job_id_text}")
            seen_job_ids.add(job_id_text)
            jobs.append(
                {
                    "id": job_id,
                    "name": name,
                    "conclusion": conclusion,
                    "run_id": run_id,
                    "run_attempt": run_attempt,
                    "started_at": job.get("started_at"),
                    "completed_at": job.get("completed_at"),
                }
            )
    if not jobs:
        raise PreparationError("jobs metadata was empty")
    if len(jobs) != expected_total:
        raise PreparationError(
            "paginated jobs response was incomplete: "
            f"expected {expected_total}, received {len(jobs)}"
        )
    return jobs


def prepare_evidence(
    run_id: int,
    output_dir: Path,
    *,
    capture_github: CaptureGitHub | None = None,
    download_github: DownloadGitHub | None = None,
) -> dict[str, Any]:
    if run_id <= 0:
        raise PreparationError("run ID must be a positive integer")
    if not output_dir.is_absolute():
        raise PreparationError("output directory must be absolute")
    output_dir = output_dir.resolve()
    if output_dir == Path("/"):
        raise PreparationError("output directory cannot be the filesystem root")

    output_dir.mkdir(parents=True, exist_ok=True)
    model_input_dir = output_dir / "model"
    model_input_dir.mkdir(parents=True, exist_ok=True)
    status_path = output_dir / "preparation-status.json"
    jobs_path = output_dir / "jobs.jsonl"

    run_attempt: int | None = None
    snapshots_verified = False

    def write_status(status: str) -> None:
        payload = {
            "schema": _schema(PREPARATION_STATUS_SCHEMA),
            "run_id": run_id,
            "status": status,
        }
        if run_attempt is not None:
            payload["run_attempt"] = run_attempt
        _atomic_write_json(
            status_path,
            payload,
            compact=True,
        )

    def invalidate_partial_evidence(status: str) -> None:
        nonlocal run_attempt
        run_attempt = None
        _atomic_write_json(
            output_dir / "run.json",
            {"run_id": run_id, "status": status},
        )
        _atomic_write_json(
            output_dir / "analysis-input.json",
            {
                "schema": _schema(ANALYSIS_INPUT_SCHEMA),
                "run_id": run_id,
                "status": status,
            },
            compact=True,
        )
        _atomic_write_json(
            model_input_dir / "classification-input.json",
            {
                "schema": _schema(CLASSIFICATION_INPUT_SCHEMA),
                "run_id": run_id,
                "status": status,
            },
            compact=True,
        )
        _atomic_write_jsonl(jobs_path, [])
        _atomic_write(output_dir / "failed-job-ids.txt", b"")

    write_status("preparing")
    invalidate_partial_evidence("preparing")
    try:
        if (capture_github is None or download_github is None) and shutil.which(
            "gh"
        ) is None:
            raise PreparationError("required command is unavailable: gh")
        capture = capture_github or _capture_github
        download = download_github or _download_github

        run_result = capture(
            f"repos/{REPOSITORY}/actions/runs/{run_id}",
            False,
        )
        _atomic_write(output_dir / "run.err", run_result.stderr)
        run = _run_metadata(_github_json(run_result, "run metadata"), run_id)
        run_attempt = run["run_attempt"]
        _atomic_write_json(output_dir / "run.json", run)
        write_status("preparing")

        jobs_result = capture(
            f"repos/{REPOSITORY}/actions/runs/{run_id}/jobs?per_page=100",
            True,
        )
        _atomic_write(output_dir / "jobs.err", jobs_result.stderr)
        jobs = _job_metadata(
            _github_json(jobs_result, "jobs metadata"),
            run_id,
            run_attempt,
        )
        _atomic_write_jsonl(jobs_path, jobs)

        failed_job_ids = [
            str(job["id"])
            for job in jobs
            if job.get("conclusion") not in {"success", "skipped", None}
        ]
        failed_job_ids_text = "".join(
            f"{job_id}\n" for job_id in failed_job_ids
        )
        _atomic_write(
            output_dir / "failed-job-ids.txt",
            failed_job_ids_text.encode(),
        )

        for job_id in failed_job_ids:
            error_path = output_dir / f"{job_id}.err"
            returncode = download(
                f"repos/{REPOSITORY}/actions/jobs/{job_id}/logs",
                output_dir / f"{job_id}.log",
                error_path,
            )
            if returncode != 0 and (
                not error_path.is_file() or error_path.stat().st_size == 0
            ):
                _atomic_write(
                    error_path,
                    f"gh api exited with status {returncode}\n".encode(),
                )

        try:
            final_jobs_result = capture(
                f"repos/{REPOSITORY}/actions/runs/{run_id}/jobs?per_page=100",
                True,
            )
            _atomic_write(output_dir / "jobs.err", final_jobs_result.stderr)
            final_jobs = _job_metadata(
                _github_json(final_jobs_result, "final jobs metadata"),
                run_id,
                run_attempt,
            )
            if sorted(final_jobs, key=lambda item: item["id"]) != sorted(
                jobs,
                key=lambda item: item["id"],
            ):
                raise PreparationError(
                    "job metadata changed while evidence was prepared"
                )

            final_run_result = capture(
                f"repos/{REPOSITORY}/actions/runs/{run_id}",
                False,
            )
            _atomic_write(output_dir / "run.err", final_run_result.stderr)
            final_run = _run_metadata(
                _github_json(final_run_result, "final run metadata"),
                run_id,
            )
            if final_run != run:
                raise PreparationError(
                    "run metadata changed while evidence was prepared"
                )
            run = final_run
            _atomic_write_json(output_dir / "run.json", run)
            snapshots_verified = True
        except Exception:
            invalidate_partial_evidence("invalidated")
            raise

        ledger = build_ledger(output_dir, jobs_path)
        ledger["run_id"] = run_id
        ledger["run_attempt"] = run_attempt
        _atomic_write_json(output_dir / "failure-ledger.json", ledger)

        compact_summary = {
            "schema": _schema(EVIDENCE_SUMMARY_SCHEMA),
            "run_id": run_id,
            "run_attempt": run_attempt,
            "summary": ledger["summary"],
        }
        _atomic_write_json(
            output_dir / "evidence-summary.json",
            compact_summary,
            compact=True,
        )
        analysis_input = build_analysis_input(ledger, run, output_dir)
        _atomic_write_analysis_input(
            output_dir / "analysis-input.json",
            analysis_input,
        )
        classification_input = build_classification_input(analysis_input)
        _atomic_write_analysis_input(
            model_input_dir / "classification-input.json",
            classification_input,
        )
        write_status("ready")
        return compact_summary
    except Exception:
        if not snapshots_verified:
            invalidate_partial_evidence("invalidated")
        write_status("failed")
        raise


def _workflow_prepare_main(argv: list[str]) -> None:
    if argv:
        raise WorkflowError("workflow-prepare does not accept arguments")

    run_id = _workflow_run_id()
    output_dir = _workflow_evidence_dir(run_id)
    output_dir.mkdir(parents=True, exist_ok=True)
    model_required = False
    failure: Exception | None = None
    try:
        _atomic_write(
            output_dir / "summary.md",
            _workflow_unavailable_summary(run_id).encode(),
        )
        _atomic_write_json(
            output_dir / "execution-metadata.json",
            {
                "schema": _schema(EXECUTION_METADATA_SCHEMA),
                "model": _required_workflow_env("SUMMARY_MODEL"),
                "summarizer_sha": _required_workflow_env("GITHUB_SHA"),
            },
            compact=True,
        )
        prepare_evidence(run_id, output_dir)
        model_required = _workflow_model_required(output_dir, run_id)
    except Exception as exc:
        failure = exc

    try:
        _stage_replay_artifact(output_dir)
    except Exception as exc:
        if failure is None:
            failure = exc

    _append_github_outputs(
        {
            "model_required": str(model_required).lower(),
            "model_schema": _model_schema_output(),
        }
    )
    if failure is not None:
        print(f"Workflow evidence preparation failed: {failure}", file=sys.stderr)
        raise SystemExit(1) from failure


def _write_failed_finalization_status(output_dir: Path) -> None:
    status_path = output_dir / "finalization-status.json"
    try:
        status = json.loads(status_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        status = None
    if (
        isinstance(status, dict)
        and status.get("schema") == _schema(FINALIZATION_STATUS_SCHEMA)
        and status.get("status") == "failed"
    ):
        return
    _atomic_write_json(
        status_path,
        {
            "schema": _schema(FINALIZATION_STATUS_SCHEMA),
            "status": "failed",
            "reason": "workflow finalization did not complete",
        },
        compact=True,
    )


def _workflow_finalize_main(argv: list[str]) -> None:
    if argv:
        raise WorkflowError("workflow-finalize does not accept arguments")

    run_id = _workflow_run_id()
    output_dir = _workflow_evidence_dir(run_id)
    output_dir.mkdir(parents=True, exist_ok=True)
    summary_path = output_dir / "summary.md"
    failure: BaseException | None = None
    finalized = False
    try:
        model_required = _workflow_model_required(output_dir, run_id)
        if not model_required:
            try:
                analysis = json.loads(
                    (output_dir / "analysis-input.json").read_text(
                        encoding="utf-8"
                    )
                )
            except (OSError, json.JSONDecodeError) as exc:
                raise WorkflowError("analysis input is unavailable") from exc
            analysis = _require_schema(
                analysis,
                ANALYSIS_INPUT_SCHEMA,
                "analysis input",
            )
            analysis_digest = analysis.get("analysis_digest")
            if not isinstance(analysis_digest, str):
                raise WorkflowError("analysis input has no digest")
            os.environ["MODEL_DECISIONS_JSON"] = json.dumps(
                {
                    "schema": _schema(MODEL_DECISIONS_SCHEMA),
                    "analysis_digest": analysis_digest,
                    "groups": [],
                },
                separators=(",", ":"),
                sort_keys=True,
            )
        elif not os.environ.get("MODEL_DECISIONS_JSON", "").strip():
            raise WorkflowError("model decisions are unavailable")

        _finalize_main(
            [
                "--analysis-input",
                str(output_dir / "analysis-input.json"),
                "--decisions-env",
                "MODEL_DECISIONS_JSON",
                "--output",
                str(summary_path),
                "--result-json",
                str(output_dir / "finalization-result.json"),
                "--slack-payload",
                str(output_dir / "slack-payload.json"),
                "--slack-prefix-env",
                "SLACK_ONCALL_TAG",
            ]
        )
        summary = summary_path.read_text(encoding="utf-8")
        finalized = True
        _append_github_step_summary(summary)
    except (Exception, SystemExit) as exc:
        failure = exc

    if not finalized:
        try:
            _write_failed_finalization_status(output_dir)
        except Exception as exc:
            if failure is None:
                failure = exc
        fallback = _workflow_unavailable_summary(run_id)
        try:
            _atomic_write(summary_path, fallback.encode())
        except Exception as exc:
            if failure is None:
                failure = exc
        try:
            _append_github_step_summary(fallback)
        except Exception as exc:
            if failure is None:
                failure = exc

    try:
        _stage_replay_artifact(output_dir)
    except Exception as exc:
        if failure is None:
            failure = exc

    if failure is not None:
        print(f"Workflow summary finalization failed: {failure}", file=sys.stderr)
        raise SystemExit(1) from failure


def _build_main(argv: list[str]) -> None:
    parser = argparse.ArgumentParser(
        description="Build a unique failed-test ledger from cached Actions logs."
    )
    parser.add_argument(
        "--logs-dir",
        required=True,
        type=Path,
        help="Directory containing one <job_id>.log file per downloaded failed job.",
    )
    parser.add_argument(
        "--jobs",
        required=True,
        type=Path,
        help="JSONL jobs metadata created in Step 1 of the skill.",
    )
    parser.add_argument(
        "--run-id",
        type=int,
        help="Optional workflow run ID to embed as ledger provenance.",
    )
    args = parser.parse_args(argv)

    if not args.logs_dir.is_dir():
        parser.error(f"logs directory does not exist: {args.logs_dir}")
    if not args.jobs.is_file():
        parser.error(f"jobs file does not exist: {args.jobs}")

    try:
        ledger = build_ledger(args.logs_dir, args.jobs)
    except ValueError as exc:
        parser.error(str(exc))
    if args.run_id is not None:
        ledger["run_id"] = args.run_id
    print(json.dumps(ledger, indent=2, sort_keys=True))


def _prepare_main(argv: list[str]) -> None:
    parser = argparse.ArgumentParser(
        description="Fetch Actions evidence and build a canonical failure ledger."
    )
    parser.add_argument("--run-id", required=True, type=int)
    parser.add_argument("--output-dir", required=True, type=Path)
    args = parser.parse_args(argv)

    try:
        compact_summary = prepare_evidence(args.run_id, args.output_dir)
    except (OSError, PreparationError, ValueError) as exc:
        print(f"Evidence preparation failed: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
    print(json.dumps(compact_summary, separators=(",", ":"), sort_keys=True))


def _finalize_main(argv: list[str]) -> None:
    parser = argparse.ArgumentParser(
        description=(
            "Validate model classifications and render a deterministic Slack summary."
        )
    )
    parser.add_argument("--analysis-input", required=True, type=Path)
    decisions_source = parser.add_mutually_exclusive_group(required=True)
    decisions_source.add_argument("--decisions", type=Path)
    decisions_source.add_argument(
        "--decisions-env",
        help="Environment variable containing the model decisions JSON.",
    )
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument(
        "--result-json",
        type=Path,
        help="Optional path for the structured deterministic finalization result.",
    )
    parser.add_argument(
        "--slack-payload",
        type=Path,
        help="Optional path for a Slack incoming-webhook JSON payload.",
    )
    parser.add_argument(
        "--slack-prefix-env",
        help=(
            "Optional environment variable prepended to the Slack summary, "
            "such as an on-call tag."
        ),
    )
    args = parser.parse_args(argv)

    if not args.analysis_input.is_file():
        parser.error(f"analysis input does not exist: {args.analysis_input}")
    if args.decisions is not None and not args.decisions.is_file():
        parser.error(f"model decisions do not exist: {args.decisions}")
    if args.decisions_env is not None and re.fullmatch(
        r"[A-Za-z_][A-Za-z0-9_]*",
        args.decisions_env,
    ) is None:
        parser.error("decisions environment variable name is invalid")
    if not args.output.is_absolute():
        parser.error("output path must be absolute")
    if args.result_json is not None and not args.result_json.is_absolute():
        parser.error("result JSON path must be absolute")
    if args.slack_payload is not None and not args.slack_payload.is_absolute():
        parser.error("Slack payload path must be absolute")
    if args.slack_prefix_env is not None and re.fullmatch(
        r"[A-Za-z_][A-Za-z0-9_]*",
        args.slack_prefix_env,
    ) is None:
        parser.error("Slack prefix environment variable name is invalid")
    if args.slack_prefix_env is not None and args.slack_payload is None:
        parser.error("--slack-prefix-env requires --slack-payload")

    status_path = args.output.parent / "finalization-status.json"
    decisions_audit_path = args.output.parent / "model-decisions.json"
    write_paths = {
        "summary output": args.output.resolve(),
        "finalization status": status_path.resolve(),
        "model decisions audit": decisions_audit_path.resolve(),
    }
    if args.result_json is not None:
        write_paths["result JSON"] = args.result_json.resolve()
    if args.slack_payload is not None:
        write_paths["Slack payload"] = args.slack_payload.resolve()
    by_path: dict[Path, list[str]] = {}
    for label, path in write_paths.items():
        by_path.setdefault(path, []).append(label)
    collisions = [
        labels for labels in by_path.values() if len(labels) > 1
    ]
    if collisions:
        parser.error(
            "finalization output paths must be distinct: "
            + "; ".join(", ".join(labels) for labels in collisions)
        )
    analysis_input_path = args.analysis_input.resolve()
    decisions_input_path = (
        args.decisions.resolve() if args.decisions is not None else None
    )
    for label, path in write_paths.items():
        if path == analysis_input_path:
            parser.error(f"{label} collides with the analysis input")
        if (
            decisions_input_path is not None
            and path == decisions_input_path
            and label != "model decisions audit"
        ):
            parser.error(f"{label} collides with the model decisions input")

    args.output.parent.mkdir(parents=True, exist_ok=True)
    finalization_outputs = [args.output]
    if args.result_json is not None:
        args.result_json.parent.mkdir(parents=True, exist_ok=True)
        finalization_outputs.append(args.result_json)
    if args.slack_payload is not None:
        args.slack_payload.parent.mkdir(parents=True, exist_ok=True)
        finalization_outputs.append(args.slack_payload)
    status: dict[str, Any] = {
        "schema": _schema(FINALIZATION_STATUS_SCHEMA),
        "status": "preparing",
    }
    _atomic_write_json(status_path, status, compact=True)

    try:
        _invalidate_finalization_outputs(finalization_outputs)
        analysis = json.loads(args.analysis_input.read_text(encoding="utf-8"))
        if args.decisions is not None:
            if args.decisions.stat().st_size > MAX_MODEL_DECISIONS_BYTES:
                raise FinalizationError("model decisions exceed the size limit")
            decisions_text = args.decisions.read_text(encoding="utf-8")
        else:
            decisions_text = os.environ.get(args.decisions_env or "")
            if decisions_text is None:
                raise FinalizationError(
                    "model decisions environment variable is not set"
                )
            if len(decisions_text.encode()) > MAX_MODEL_DECISIONS_BYTES:
                raise FinalizationError("model decisions exceed the size limit")
        _atomic_write(
            decisions_audit_path,
            f"{decisions_text.rstrip()}\n".encode(),
        )
        decisions = json.loads(decisions_text)
        result = finalize_analysis(analysis, decisions)
        summary = result["slack_mrkdwn"]
        slack_payload: dict[str, Any] | None = None
        if args.slack_payload is not None:
            prefix = (
                os.environ.get(args.slack_prefix_env, "").strip()
                if args.slack_prefix_env is not None
                else ""
            )
            slack_text = f"{prefix} {summary}".strip()
            if len(slack_text.encode()) > SLACK_SECTION_LIMIT:
                raise FinalizationError(
                    "Slack payload text exceeded the section limit"
                )
            slack_payload = {
                "text": f"Test Suite #{result['run_number']} summary",
                "blocks": [
                    {
                        "type": "section",
                        "text": {
                            "type": "mrkdwn",
                            "text": slack_text,
                        },
                    }
                ],
            }
        if args.result_json is not None:
            _atomic_write_json(args.result_json, result)
        if args.slack_payload is not None and slack_payload is not None:
            _atomic_write_json(args.slack_payload, slack_payload)
        _atomic_write(args.output, f"{summary}\n".encode())
        status = {
            "schema": _schema(FINALIZATION_STATUS_SCHEMA),
            "status": "ready",
            "run_id": result["run_id"],
            "run_attempt": result["run_attempt"],
            "analysis_digest": result["analysis_digest"],
            "decisions_digest": _canonical_digest(decisions),
        }
        _atomic_write_json(status_path, status, compact=True)
    except Exception as exc:
        invalidation_error: Exception | None = None
        try:
            _invalidate_finalization_outputs(finalization_outputs)
        except Exception as cleanup_exc:
            invalidation_error = cleanup_exc
        status["status"] = "failed"
        status["reason"] = str(exc)
        if invalidation_error is not None:
            status["invalidation_error"] = str(invalidation_error)
        _atomic_write_json(status_path, status, compact=True)
        print(f"Summary finalization failed: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
    print(summary)


def main() -> None:
    if len(sys.argv) > 1 and sys.argv[1] == "workflow-prepare":
        _workflow_prepare_main(sys.argv[2:])
        return
    if len(sys.argv) > 1 and sys.argv[1] == "workflow-finalize":
        _workflow_finalize_main(sys.argv[2:])
        return
    if len(sys.argv) > 1 and sys.argv[1] == "prepare":
        _prepare_main(sys.argv[2:])
        return
    if len(sys.argv) > 1 and sys.argv[1] == "finalize":
        _finalize_main(sys.argv[2:])
        return
    _build_main(sys.argv[1:])


if __name__ == "__main__":
    main()
