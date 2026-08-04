"""Offline sanitized replays of historical test-suite evidence."""

import gzip
import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from typing import Any

from failure_ledger import (
    MODEL_DECISIONS_SCHEMA,
    SCHEMA_VERSION,
    build_analysis_input,
    build_classification_input,
    build_ledger,
    finalize_analysis,
)


FIXTURE_ROOT = (
    Path(__file__).resolve().parent / "testdata" / "replays"
)
RUN_NUMBERS = ("1093", "1098", "1102")


class HistoricalReplayTest(unittest.TestCase):
    @staticmethod
    def _read_evidence(run_number: str) -> dict[str, Any]:
        path = FIXTURE_ROOT / run_number / "evidence.json.gz"
        with gzip.open(path, "rt", encoding="utf-8") as handle:
            return json.load(handle)

    @staticmethod
    def _read_expected(run_number: str) -> dict[str, Any]:
        path = FIXTURE_ROOT / run_number / "expected.json"
        return json.loads(path.read_text(encoding="utf-8"))

    @staticmethod
    def _write_evidence(
        root: Path,
        evidence: dict[str, Any],
    ) -> Path:
        jobs_path = root / "jobs.jsonl"
        jobs_path.write_text(
            "".join(
                json.dumps(job, separators=(",", ":")) + "\n"
                for job in evidence["jobs"]
            ),
            encoding="utf-8",
        )
        for job_id, log in evidence["logs"].items():
            (root / f"{job_id}.log").write_text(log, encoding="utf-8")
        return jobs_path

    @staticmethod
    def _identity_digest(analysis: dict[str, Any]) -> str:
        identities = sorted(
            [unit["package"], unit["test"]]
            for unit in analysis["test_units"]
        )
        return hashlib.sha256(
            json.dumps(identities, separators=(",", ":")).encode()
        ).hexdigest()

    def _model_decisions(
        self,
        expected: dict[str, Any],
        analysis: dict[str, Any],
    ) -> dict[str, Any]:
        units_by_identity = {
            (unit["job_id"], unit["package"], unit["test"]): unit
            for unit in analysis["test_units"]
        }
        evidence_by_id = {
            item["evidence_id"]: item for item in analysis["evidence"]
        }
        deterministic_ids = {
            unit_id
            for group in analysis.get("deterministic_test_groups", [])
            for unit_id in group["unit_ids"]
        }
        groups = []
        for definition in expected["groups"]:
            units = []
            for selector in definition["selectors"]:
                identity = (
                    str(selector["job_id"]),
                    selector["package"],
                    selector["test"],
                )
                self.assertIn(
                    identity,
                    units_by_identity,
                    f"{definition['label']} selector was not replayed",
                )
                units.append(units_by_identity[identity])

            deterministic_units = [
                unit
                for unit in units
                if unit["unit_id"] in deterministic_ids
            ]
            if deterministic_units:
                self.assertEqual(
                    len(deterministic_units),
                    len(units),
                    f"{definition['label']} mixes deterministic and model units",
                )
                continue

            evidence_refs = []
            needle = definition["evidence_contains"].casefold()
            for unit in units:
                self.assertTrue(
                    unit["evidence_refs"],
                    (
                        f"{definition['label']} has no owned evidence for "
                        f"{unit['package']}::{unit['test']}"
                    ),
                )
                matching_refs = [
                    evidence_ref
                    for evidence_ref in unit["evidence_refs"]
                    if needle
                    in evidence_by_id[evidence_ref]["text"].casefold()
                ]
                self.assertTrue(
                    matching_refs,
                    (
                        f"{definition['label']} has no owned evidence matching "
                        f"{definition['evidence_contains']!r} for "
                        f"{unit['package']}::{unit['test']}"
                    ),
                )
                self.assertIn(
                    definition["category"],
                    unit["allowed_categories"],
                )
                evidence_refs.extend(matching_refs[:1])

            group = {
                "unit_ids": [unit["unit_id"] for unit in units],
                "category": definition["category"],
                "cause": definition["cause"],
                "evidence_refs": evidence_refs,
            }
            for key in ("ambiguity", "note"):
                if key in definition:
                    group[key] = definition[key]
            groups.append(group)

        classified_ids = {
            unit_id
            for group in groups
            for unit_id in group["unit_ids"]
        }
        self.assertEqual(
            classified_ids,
            {
                unit["unit_id"]
                for unit in analysis["test_units"]
                if unit["unit_id"] not in deterministic_ids
            },
        )
        return {
            "schema": {
                "name": MODEL_DECISIONS_SCHEMA,
                "version": SCHEMA_VERSION,
            },
            "analysis_digest": analysis["analysis_digest"],
            "groups": groups,
        }

    def _replay(
        self,
        run_number: str,
    ) -> tuple[
        dict[str, Any],
        dict[str, Any],
        dict[str, Any],
        dict[str, Any],
    ]:
        evidence = self._read_evidence(run_number)
        expected = self._read_expected(run_number)
        self.assertEqual(
            evidence["schema"],
            {
                "name": "test-suite-filtered-historical-evidence",
                "version": 1,
            },
        )
        self.assertEqual(
            expected["schema"],
            {
                "name": "test-suite-historical-expected-decisions",
                "version": 1,
            },
        )
        self.assertEqual(evidence["run"]["run_number"], int(run_number))
        self.assertEqual(
            evidence["source_run_id"],
            evidence["run"]["run_id"],
        )
        self.assertNotIn("groups", evidence)
        self.assertNotIn("result", evidence)

        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            jobs_path = self._write_evidence(root, evidence)
            ledger = build_ledger(root, jobs_path)
            ledger["run_id"] = evidence["run"]["run_id"]
            ledger["run_attempt"] = evidence["run"]["run_attempt"]
            analysis = build_analysis_input(ledger, evidence["run"], root)
            decisions = self._model_decisions(expected, analysis)
            result = finalize_analysis(analysis, decisions)
        return ledger, analysis, result, expected

    def test_historical_runs_reconcile_to_golden_results(self) -> None:
        for run_number in RUN_NUMBERS:
            with self.subTest(run_number=run_number):
                ledger, analysis, result, expected = self._replay(run_number)
                golden = expected["result"]

                self.assertTrue(ledger["summary"]["ledger_complete"])
                self.assertEqual(
                    self._identity_digest(analysis),
                    golden["identity_digest"],
                )
                self.assertEqual(
                    analysis["unit_counts"]["tests"],
                    golden["test_total"],
                )
                self.assertEqual(
                    result["category_counts"],
                    golden["category_counts"],
                )
                self.assertEqual(
                    result["package_teardowns"],
                    golden["package_teardowns"],
                )
                self.assertEqual(result["verdict"], golden["verdict"])
                self.assertEqual(result["headline"], golden["headline"])
                self.assertEqual(result["confidence"], golden["confidence"])
                self.assertEqual(
                    sorted(
                        unit["package"]
                        for unit in analysis["package_teardown_units"]
                    ),
                    golden["teardown_packages"],
                )
                self.assertEqual(
                    sorted(
                        unit["test"]
                        for unit in analysis["test_units"]
                        if unit["machine_facts"]["shape_b_candidate"]
                    ),
                    golden["shape_b_tests"],
                )

    def test_run_1093_preserves_http_error_401_evidence(self) -> None:
        _, analysis, _, expected = self._replay("1093")
        group = next(
            group
            for group in expected["groups"]
            if group["label"] == "test_utils_401"
        )
        self.assertEqual(len(group["selectors"]), 2)

        evidence_by_id = {
            item["evidence_id"]: item["text"]
            for item in analysis["evidence"]
        }
        units_by_identity = {
            (unit["job_id"], unit["package"], unit["test"]): unit
            for unit in analysis["test_units"]
        }
        for selector in group["selectors"]:
            unit = units_by_identity[
                (
                    str(selector["job_id"]),
                    selector["package"],
                    selector["test"],
                )
            ]
            self.assertTrue(
                any(
                    "HTTP ERROR 401" in evidence_by_id[evidence_ref]
                    for evidence_ref in unit["evidence_refs"]
                )
            )

    def test_run_1098_keeps_max_groups_separate_from_api_errors(self) -> None:
        _, analysis, result, expected = self._replay("1098")
        classification_input = build_classification_input(analysis)
        max_group = next(
            group
            for group in expected["groups"]
            if group["label"] == "max_groups_cleanup"
        )

        self.assertEqual(len(max_group["selectors"]), 159)
        self.assertEqual(analysis["unit_counts"]["tests"], 288)
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 282)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 6)
        self.assertEqual(
            classification_input["unit_counts"]["deterministic_tests"],
            282,
        )
        self.assertEqual(
            classification_input["unit_counts"]["model_tests"],
            6,
        )
        self.assertEqual(len(classification_input["test_units"]), 6)
        self.assertEqual(
            {
                (cohort["signature"], cohort["unique_tests"])
                for cohort in classification_input["deterministic_cohorts"]
            },
            {
                ("max_groups_per_org_cleanup", 159),
                ("groups_post_http_500", 122),
                (
                    "post_destroy_private_endpoint_dependency_cleanup",
                    1,
                ),
            },
        )
        self.assertEqual(result["category_counts"]["cleanup"], 160)
        self.assertEqual(result["category_counts"]["api_error"], 124)
        self.assertIn("Cleanup: 160 tests", result["slack_mrkdwn"])
        self.assertIn("API errors: 124 tests", result["slack_mrkdwn"])

        units_by_identity = {
            (unit["job_id"], unit["package"], unit["test"]): unit
            for unit in analysis["test_units"]
        }
        for selector in max_group["selectors"]:
            unit = units_by_identity[
                (
                    str(selector["job_id"]),
                    selector["package"],
                    selector["test"],
                )
            ]
            self.assertEqual(unit["allowed_categories"], ["cleanup"])
        self.assertEqual(
            hashlib.sha256(result["slack_mrkdwn"].encode()).hexdigest(),
            "2a923fb243b90acd957140720b418fe5b6b1fe3b7aa27c9eab8c162aef3d09af",
        )

    def test_classification_projection_compacts_only_supported_replays(
        self,
    ) -> None:
        expected_counts = {
            "1093": (131, 2),
            "1098": (282, 6),
            "1102": (0, 7),
        }
        for run_number, (
            deterministic_tests,
            model_tests,
        ) in expected_counts.items():
            with self.subTest(run_number=run_number):
                _, analysis, _, _ = self._replay(run_number)
                classification_input = build_classification_input(analysis)
                self.assertEqual(
                    classification_input["unit_counts"][
                        "deterministic_tests"
                    ],
                    deterministic_tests,
                )
                self.assertEqual(
                    classification_input["unit_counts"]["model_tests"],
                    model_tests,
                )
                self.assertEqual(
                    len(classification_input["test_units"]),
                    model_tests,
                )
                deterministic_ids = {
                    unit_id
                    for group in analysis["deterministic_test_groups"]
                    for unit_id in group["unit_ids"]
                }
                self.assertTrue(
                    all(
                        item["unit_id"] not in deterministic_ids
                        for item in classification_input["evidence"]
                    )
                )
                if deterministic_tests:
                    analysis_bytes = len(
                        json.dumps(
                            analysis,
                            separators=(",", ":"),
                            sort_keys=True,
                        ).encode()
                    )
                    classification_bytes = len(
                        json.dumps(
                            classification_input,
                            separators=(",", ":"),
                            sort_keys=True,
                        ).encode()
                    )
                    self.assertLess(
                        classification_bytes,
                        analysis_bytes // 5,
                    )

    def test_run_1102_preserves_shape_b_and_package_attribution(self) -> None:
        _, analysis, result, _ = self._replay("1102")
        packages = {
            unit["package"].rsplit("/", 1)[-1]
            for unit in analysis["package_teardown_units"]
        }

        self.assertEqual(
            packages,
            {
                "advancedcluster",
                "mongodbemployeeaccessgrant",
                "streamconnection",
                "streamprocessor",
            },
        )
        self.assertEqual(result["category_counts"]["timeout"], 7)
        self.assertEqual(
            result["slack_mrkdwn"].count(
                "Delete-on-timeout unverified: 1 test"
            ),
            1,
        )
        self.assertIn(
            "(included in Timeout total)",
            result["slack_mrkdwn"],
        )


if __name__ == "__main__":
    unittest.main()
