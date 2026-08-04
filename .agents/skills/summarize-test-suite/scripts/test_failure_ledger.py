import contextlib
import io
import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import failure_ledger
from failure_ledger import (
    ANALYSIS_INPUT_SCHEMA,
    CLASSIFICATION_INPUT_SCHEMA,
    EVIDENCE_SUMMARY_SCHEMA,
    FAILURE_LEDGER_SCHEMA,
    MODEL_DECISIONS_OUTPUT_SCHEMA,
    MODEL_DECISIONS_SCHEMA,
    PACKAGE_TEARDOWN_RE,
    PREPARATION_STATUS_SCHEMA,
    REPLAY_ARTIFACT_FILES,
    SCHEMA_VERSION,
    SLACK_HARD_LIMIT,
    FinalizationError,
    PreparationError,
    _canonical_digest,
    _capture_github,
    _download_github,
    _finalize_main,
    _job_context,
    _job_metadata,
    _job_role,
    _model_schema_output,
    _slack_text,
    _stage_replay_artifact,
    _workflow_finalize_main,
    _workflow_prepare_main,
    build_analysis_input,
    build_classification_input,
    build_ledger,
    finalize_analysis,
    prepare_evidence,
)

REPOSITORY_ROOT = Path(__file__).resolve().parents[4]


class FailureLedgerTest(unittest.TestCase):
    def setUp(self) -> None:
        self.temp_dir = tempfile.TemporaryDirectory()
        self.root = Path(self.temp_dir.name)
        self.jobs_path = self.root / "jobs.jsonl"

    def tearDown(self) -> None:
        self.temp_dir.cleanup()

    def workflow_env(
        self,
        *,
        run_id: int = 30228140283,
    ) -> tuple[dict[str, str], Path]:
        output_dir = self.root.resolve() / f"test-suite-summary-{run_id}"
        return (
            {
                "GITHUB_OUTPUT": str(self.root / "github-output.txt"),
                "GITHUB_RUN_ID": str(run_id),
                "GITHUB_SHA": "abcdef123456",
                "GITHUB_STEP_SUMMARY": str(
                    self.root / "github-step-summary.md"
                ),
                "RUNNER_TEMP": str(self.root),
                "SUMMARY_MODEL": "claude-opus-5",
                "SLACK_ONCALL_TAG": "<!subteam^oncall>",
            },
            output_dir,
        )

    @staticmethod
    def workflow_outputs(path: Path) -> dict[str, str]:
        return dict(
            line.split("=", 1)
            for line in path.read_text(encoding="utf-8").splitlines()
        )

    @staticmethod
    def write_json(path: Path, value: object) -> None:
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(json.dumps(value), encoding="utf-8")

    @staticmethod
    def read_json(path: Path) -> object:
        return json.loads(path.read_text(encoding="utf-8"))

    def write_workflow_preparation(
        self,
        output_dir: Path,
        *,
        run_id: int,
        model_tests: int,
    ) -> None:
        self.write_json(
            output_dir / "preparation-status.json",
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "run_attempt": 1,
                "status": "ready",
            },
        )
        self.write_json(
            output_dir / "model" / "classification-input.json",
            {
                "schema": {
                    "name": CLASSIFICATION_INPUT_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "unit_counts": {"model_tests": model_tests},
            },
        )

    def write_workflow_analysis(
        self,
        output_dir: Path,
        *,
        run_id: int,
        model_tests: int,
        analysis: dict[str, object],
    ) -> None:
        self.write_workflow_preparation(
            output_dir,
            run_id=run_id,
            model_tests=model_tests,
        )
        self.write_json(output_dir / "analysis-input.json", analysis)

    def run_workflow_finalize(
        self,
        env: dict[str, str],
        *,
        fails: bool = False,
    ) -> None:
        with mock.patch.dict(os.environ, env, clear=True):
            with contextlib.redirect_stdout(io.StringIO()):
                with contextlib.redirect_stderr(io.StringIO()):
                    if fails:
                        with self.assertRaises(SystemExit):
                            _workflow_finalize_main([])
                    else:
                        _workflow_finalize_main([])

    def write_jobs(self, *jobs: dict[str, object]) -> None:
        self.jobs_path.write_text(
            "".join(json.dumps(job) + "\n" for job in jobs),
            encoding="utf-8",
        )

    def write_log(self, job_id: int, text: str, *, complete: bool = True) -> None:
        suffix = (
            "\n2026-07-27T00:01:00Z Cleaning up orphan processes\n"
            if complete
            else "\n"
        )
        (self.root / f"{job_id}.log").write_text(
            text.rstrip() + suffix,
            encoding="utf-8",
        )

    @staticmethod
    def _test_job_name(leaf_name: str, *, auth: str = "pak") -> str:
        return (
            f"1.15.x-latest-{auth} / tests-1.15.x-latest-dev "
            f"/ {leaf_name}"
        )

    @staticmethod
    def _completed_json(
        value: object,
        *,
        returncode: int = 0,
        stderr: bytes = b"",
    ) -> subprocess.CompletedProcess[bytes]:
        return subprocess.CompletedProcess(
            ["gh", "api"],
            returncode,
            stdout=json.dumps(value).encode(),
            stderr=stderr,
        )

    def _analysis_from_failed_test_log(
        self,
        lines: list[str],
        *,
        leaf_name: str = "encryption",
        run_number: int = 1098,
        auth: str = "pak",
    ) -> dict[str, object]:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name(leaf_name, auth=auth),
                "conclusion": "failure",
            },
        )
        self.write_log(1, "\n".join(lines))
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": run_number,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1
        return build_analysis_input(ledger, run, self.root)

    def test_job_roles_scope_root_only_and_nested_names(self) -> None:
        self.assertEqual(_job_role("clean-before"), "cleanup")
        self.assertEqual(_job_role("clean-after / dev"), "cleanup")
        self.assertEqual(_job_role("variables"), "setup")
        self.assertEqual(_job_role(self._test_job_name("variables")), "test")
        self.assertEqual(_job_role("trigger-test-summary"), "reporting")
        self.assertEqual(
            _job_role(self._test_job_name("trigger-test-summary")),
            "test",
        )
        self.assertEqual(
            _job_context(self._test_job_name("variables")),
            ("dev", "pak"),
        )

    def test_package_teardown_marker_matches_its_go_producer(self) -> None:
        producer = (
            REPOSITORY_ROOT / "internal/testutil/acc/shared_resource.go"
        ).read_text(encoding="utf-8")
        self.assertIsNotNone(
            PACKAGE_TEARDOWN_RE.search(producer),
            "package teardown parser no longer matches the Go cleanup marker",
        )

    def test_workflow_model_schema_matches_claude_output_contract(self) -> None:
        workflow = (
            REPOSITORY_ROOT / ".github/workflows/test-suite.yml"
        ).read_text(encoding="utf-8")
        model_schema = _model_schema_output()

        self.assertEqual(
            json.loads(model_schema),
            MODEL_DECISIONS_OUTPUT_SCHEMA,
        )
        self.assertNotIn("'", model_schema)
        self.assertIn(
            "--json-schema '${{ steps.evidence.outputs.model_schema }}'",
            workflow,
        )

    def test_workflow_summary_job_keeps_production_contract(self) -> None:
        workflow = (
            REPOSITORY_ROOT / ".github/workflows/test-suite.yml"
        ).read_text(encoding="utf-8")
        job_start = workflow.index("\n  trigger-test-summary:\n")
        job = workflow[job_start:]
        job_header, separator, _ = job.partition("\n    steps:\n")
        clean_before_start = workflow.index("\n  clean-before:\n")
        clean_before_end = workflow.index("\n  tests:\n")
        clean_before = workflow[clean_before_start:clean_before_end]

        self.assertTrue(separator)
        self.assertNotIn("\n    if:", clean_before)
        self.assertIn(
            "\n  tests:\n"
            "    needs: [clean-before, variables]\n"
            "    if: ${{ !cancelled() }}\n",
            workflow,
        )
        self.assertIn(
            "\n  clean-after:\n"
            "    needs: tests\n"
            "    if: ${{ !cancelled() }}\n",
            workflow,
        )
        self.assertNotIn("\n    name:", job_header)
        self.assertIn(
            "\n    needs: [tests, variables, clean-after]\n",
            job_header,
        )
        self.assertIn(
            "\n    if: ${{ !cancelled() && "
            "needs.variables.outputs.trigger_test_summary == 'true' }}\n",
            job_header,
        )
        self.assertIn(
            "SUMMARY_MODEL: ${{ vars.TEST_SUITE_SUMMARY_MODEL || "
            "'claude-opus-5' }}",
            workflow,
        )
        self.assertIn(
            "persist-credentials: false",
            job,
        )
        self.assertIn(
            "run: python3 -B "
            ".agents/skills/summarize-test-suite/scripts/"
            "failure_ledger.py workflow-prepare",
            job,
        )
        self.assertIn(
            "run: python3 -B "
            ".agents/skills/summarize-test-suite/scripts/"
            "failure_ledger.py workflow-finalize",
            job,
        )
        run_lines = [
            line.strip()
            for line in job.splitlines()
            if line.strip().startswith("run:")
        ]
        self.assertEqual(len(run_lines), 2)
        self.assertIn(
            "steps.evidence.outcome == 'success'",
            job,
        )
        self.assertIn(
            "steps.evidence.outputs.model_required == 'true'",
            job,
        )
        self.assertIn(
            '--disallowedTools "mcp__*" "Read(/.git/**)"',
            job,
        )
        self.assertIn(
            '--add-dir "${{ runner.temp }}/test-suite-summary-'
            '${{ github.run_id }}/model"',
            job,
        )
        self.assertEqual(workflow.count("\n          name: summary\n"), 1)
        self.assertEqual(
            workflow.count("\n          name: summary-replay\n"),
            1,
        )
        self.assertIn(
            "if: ${{ steps.finalize.outcome == 'success' }}",
            job,
        )
        self.assertIn(
            "if: ${{ always() && (failure() || "
            "steps.finalize.outcome != 'success') }}",
            job,
        )
        self.assertIn(
            "${{ runner.temp }}/test-suite-summary-${{ github.run_id }}"
            "/artifacts/summary-replay",
            job,
        )
        for shell_text in (
            "steps.paths.outputs",
            "\n        run: |",
            "jq ",
            "mkdir ",
            "printf ",
            'cat "$SUMMARY_FILE"',
            "Persist summary to run page",
        ):
            with self.subTest(shell_text=shell_text):
                self.assertNotIn(shell_text, job)
        for harness_text in (
            "CLOUDP-424367_test_suite",
            "matrix.run_id",
            "matrix.run_date",
            "Initialize evaluation paths",
            "historical run",
        ):
            with self.subTest(harness_text=harness_text):
                self.assertNotIn(harness_text, workflow)

    def test_workflow_prepare_handles_model_and_no_model_runs(self) -> None:
        for index, (model_tests, expected) in enumerate(
            ((6, "true"), (0, "false")),
        ):
            with self.subTest(model_tests=model_tests):
                env, output_dir = self.workflow_env(
                    run_id=30228140283 + index,
                )
                env["GITHUB_OUTPUT"] = str(
                    self.root / f"github-output-{index}.txt"
                )
                run_id = int(env["GITHUB_RUN_ID"])

                def fake_prepare(
                    prepared_run_id: int,
                    prepared_dir: Path,
                ) -> dict:
                    self.assertEqual((prepared_run_id, prepared_dir), (run_id, output_dir))
                    self.write_workflow_preparation(
                        prepared_dir,
                        run_id=run_id,
                        model_tests=model_tests,
                    )
                    self.write_json(
                        prepared_dir / "analysis-input.json",
                        {"prepared": True},
                    )
                    return {}

                with mock.patch.dict(os.environ, env, clear=False):
                    with mock.patch(
                        "failure_ledger.prepare_evidence",
                        side_effect=fake_prepare,
                    ):
                        _workflow_prepare_main([])

                outputs = self.workflow_outputs(Path(env["GITHUB_OUTPUT"]))
                self.assertEqual(outputs["model_required"], expected)
                self.assertEqual(
                    json.loads(outputs["model_schema"]),
                    MODEL_DECISIONS_OUTPUT_SCHEMA,
                )
                if index == 0:
                    metadata = self.read_json(
                        output_dir / "execution-metadata.json"
                    )
                    self.assertEqual(
                        metadata["schema"],
                        failure_ledger._schema(
                            failure_ledger.EXECUTION_METADATA_SCHEMA
                        ),
                    )
                    self.assertEqual(metadata["model"], "claude-opus-5")
                    self.assertEqual(
                        metadata["summarizer_sha"],
                        "abcdef123456",
                    )
                    self.assertIn(
                        "summary unavailable",
                        (output_dir / "summary.md").read_text(
                            encoding="utf-8"
                        ),
                    )
                    self.assertTrue(
                        (
                            output_dir
                            / "artifacts"
                            / "summary-replay"
                            / "execution-metadata.json"
                        ).is_file()
                    )

    def test_workflow_prepare_failure_keeps_generic_diagnostic(self) -> None:
        env, output_dir = self.workflow_env()

        with mock.patch.dict(os.environ, env, clear=False):
            with mock.patch(
                "failure_ledger.prepare_evidence",
                side_effect=PreparationError("sensitive preparation detail"),
            ):
                with contextlib.redirect_stderr(io.StringIO()):
                    with self.assertRaises(SystemExit):
                        _workflow_prepare_main([])

        outputs = self.workflow_outputs(Path(env["GITHUB_OUTPUT"]))
        self.assertEqual(outputs["model_required"], "false")
        diagnostic = (output_dir / "summary.md").read_text(encoding="utf-8")
        self.assertIn("summary unavailable", diagnostic)
        self.assertNotIn("sensitive preparation detail", diagnostic)
        self.assertTrue(
            (
                output_dir
                / "artifacts"
                / "summary-replay"
                / "execution-metadata.json"
            ).is_file()
        )

    def test_replay_staging_exports_only_allowlisted_files(self) -> None:
        output_dir = self.root / "evidence"
        for index, relative_name in enumerate(REPLAY_ARTIFACT_FILES):
            path = output_dir / relative_name
            self.write_json(path, {"index": index})
        (output_dir / "123.log").write_text("raw log", encoding="utf-8")
        (output_dir / "run.err").write_text("raw error", encoding="utf-8")
        stale_replay_dir = output_dir / "artifacts" / "summary-replay"
        stale_replay_dir.mkdir(parents=True)
        (stale_replay_dir / "unexpected.txt").write_text(
            "stale",
            encoding="utf-8",
        )

        copied_files = _stage_replay_artifact(output_dir)

        replay_dir = output_dir / "artifacts" / "summary-replay"
        replay_files = {
            str(path.relative_to(replay_dir))
            for path in replay_dir.rglob("*")
            if path.is_file()
        }
        self.assertEqual(replay_files, set(REPLAY_ARTIFACT_FILES))
        self.assertEqual(set(copied_files), set(REPLAY_ARTIFACT_FILES))

        (output_dir / "analysis-input.json").unlink()
        copied_files = _stage_replay_artifact(output_dir)
        self.assertNotIn("analysis-input.json", copied_files)
        self.assertFalse((replay_dir / "analysis-input.json").exists())
        self.assertFalse((replay_dir / "unexpected.txt").exists())

        empty_dir = self.root / "empty-evidence"
        empty_dir.mkdir()

        with self.assertRaisesRegex(
            failure_ledger.WorkflowError,
            "no replay files were available",
        ):
            _stage_replay_artifact(empty_dir)

        self.assertFalse(
            (empty_dir / "artifacts" / "summary-replay").exists()
        )

    def test_replay_staging_restores_previous_snapshot_on_install_failure(
        self,
    ) -> None:
        output_dir = self.root / "evidence"
        artifact_dir = output_dir / "artifacts"
        replay_dir = artifact_dir / "summary-replay"
        source_path = output_dir / "analysis-input.json"
        source_path.parent.mkdir()
        source_path.write_text("previous", encoding="utf-8")
        _stage_replay_artifact(output_dir)
        source_path.write_text("replacement", encoding="utf-8")
        original_replace = failure_ledger.os.replace

        def fail_install(source: object, target: object) -> None:
            source_path_arg = Path(source)
            target_path_arg = Path(target)
            if (
                target_path_arg == replay_dir
                and source_path_arg.name.startswith(".summary-replay-")
                and not source_path_arg.name.endswith("-previous")
            ):
                raise OSError("simulated replay install failure")
            original_replace(source, target)

        with mock.patch(
            "failure_ledger.os.replace",
            side_effect=fail_install,
        ):
            with self.assertRaisesRegex(OSError, "install failure"):
                _stage_replay_artifact(output_dir)

        self.assertEqual(
            (replay_dir / "analysis-input.json").read_text(encoding="utf-8"),
            "previous",
        )
        self.assertFalse(
            any(
                path.name.endswith("-previous")
                for path in artifact_dir.iterdir()
            )
        )

    def test_job_metadata_rejects_a_different_run_attempt(self) -> None:
        with self.assertRaisesRegex(
            PreparationError,
            "did not match run attempt 2",
        ):
            _job_metadata(
                [
                    {
                        "total_count": 1,
                        "jobs": [
                            {
                                "id": 101,
                                "name": "variables",
                                "conclusion": "success",
                                "run_id": 30228140283,
                                "run_attempt": 1,
                            }
                        ],
                    }
                ],
                30228140283,
                2,
            )

    def test_github_subprocess_timeouts_fail_as_preparation_errors(self) -> None:
        timeout = subprocess.TimeoutExpired(["gh", "api"], 60)
        with mock.patch(
            "failure_ledger.subprocess.run",
            side_effect=timeout,
        ):
            with self.assertRaisesRegex(PreparationError, "metadata request timed out"):
                _capture_github("repos/example/actions/runs/1", False)
            with self.assertRaisesRegex(PreparationError, "log request timed out"):
                _download_github(
                    "repos/example/actions/jobs/1/logs",
                    self.root / "1.log",
                    self.root / "1.err",
                )

    def test_deduplicates_parent_tests_and_separates_non_test_failures(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("package-a"), "conclusion": "failure"},
            {"id": 2, "name": self._test_job_name("package-b"), "conclusion": "failure"},
            {"id": 3, "name": "clean-after / dev", "conclusion": "failure"},
            {
                "id": 4,
                "name": self._test_job_name("cancelled"),
                "conclusion": "cancelled",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestAccOne/sub-a (1.00s)",
                    "2026-07-27T00:00:01Z --- FAIL: TestAccOne/sub-b (1.00s)",
                    "2026-07-27T00:00:02Z --- FAIL: TestAccOne (2.00s)",
                    "2026-07-27T00:00:03Z OUT_OF_CAPACITY",
                    "2026-07-27T00:00:04Z OUT_OF_CAPACITY",
                    "2026-07-27T00:00:05Z --- FAIL: TestAccTwo (3.00s)",
                    "2026-07-27T00:00:06Z FAIL\tgithub.com/example/package-a\t3.000s",
                ]
            ),
        )
        self.write_log(
            2,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestAccOne (1.00s)",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/package-b\t1.000s",
                ]
            ),
        )
        self.write_log(
            3,
            "2026-07-27T00:00:00Z --- FAIL: TestCleanProject/id (1.00s)\n",
        )
        self.write_log(
            4,
            "2026-07-27T00:00:00Z --- FAIL: TestCancelled (1.00s)\n",
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["unique_failed_tests"], 3)
        self.assertEqual(ledger["summary"]["top_level_failure_lines"], 3)
        self.assertEqual(ledger["summary"]["subtest_failure_lines"], 2)
        self.assertEqual(ledger["summary"]["package_failure_markers"], 2)
        self.assertEqual(ledger["summary"]["cleanup_job_failures"], 1)
        self.assertEqual(ledger["summary"]["cancelled_job_failures"], 1)
        self.assertEqual(ledger["summary"]["failed_test_jobs"], 2)
        self.assertEqual(ledger["summary"]["test_job_logs_for_sweep"], 3)
        self.assertEqual(ledger["summary"]["test_jobs_ran"], 2)
        self.assertEqual(
            ledger["summary"]["test_execution_state"],
            "confirmed_ran",
        )
        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(
            {item["job_id"] for item in ledger["test_job_logs_for_sweep"]},
            {"1", "2", "4"},
        )
        self.assertEqual(
            [(item["job_id"], item["test"]) for item in ledger["tests"]],
            [("1", "TestAccOne"), ("1", "TestAccTwo"), ("2", "TestAccOne")],
        )
        self.assertEqual(
            ledger["tests"][0]["subtests"],
            ["TestAccOne/sub-a", "TestAccOne/sub-b"],
        )
        self.assertEqual(
            {item["package"] for item in ledger["tests"]},
            {
                "github.com/example/package-a",
                "github.com/example/package-b",
            },
        )

    def test_marks_missing_and_unresolved_failed_jobs_incomplete(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("missing"), "conclusion": "failure"},
            {
                "id": 2,
                "name": self._test_job_name("unresolved"),
                "conclusion": "failure",
            },
        )
        self.write_log(2, "2026-07-27T00:00:00Z command exited with status 1\n")

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unavailable_test_job_logs"], 1)
        self.assertEqual(ledger["summary"]["unresolved_failed_test_jobs"], 1)

    def test_cancelled_and_unsupported_log_issues_are_visible(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("cancelled"),
                "conclusion": "cancelled",
            },
            {
                "id": 2,
                "name": self._test_job_name("startup"),
                "conclusion": "startup_failure",
            },
        )
        self.write_log(
            2,
            "2026-07-27T00:00:00Z startup failed",
            complete=False,
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unavailable_test_job_logs"], 2)
        self.assertEqual(
            [item["job_id"] for item in ledger["missing_log_jobs"]],
            ["1"],
        )
        self.assertEqual(
            [item["job_id"] for item in ledger["partial_log_jobs"]],
            ["2"],
        )

    def test_rejects_duplicate_job_ids(self) -> None:
        duplicate = {
            "id": 1,
            "name": self._test_job_name("duplicate"),
            "conclusion": "failure",
        }
        self.jobs_path.write_text(
            f"{json.dumps(duplicate)}\n{json.dumps(duplicate)}\n",
            encoding="utf-8",
        )

        with self.assertRaisesRegex(ValueError, "duplicate job id 1"):
            build_ledger(self.root, self.jobs_path)

    def test_includes_subtest_only_parent_once(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("package-a"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestParent/sub-a (1.00s)",
                    "2026-07-27T00:00:01Z --- FAIL: TestParent/sub-b (1.00s)",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["unique_failed_tests"], 1)
        self.assertEqual(ledger["summary"]["subtest_only_parent_tests"], 1)
        self.assertEqual(ledger["summary"]["package_unattributed_tests"], 1)
        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertTrue(ledger["tests"][0]["derived_from_subtest_only"])

    def test_keeps_same_test_name_in_two_packages_distinct(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("multi-package"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestSameName (1.00s)",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/package-a\t1.000s",
                    "2026-07-27T00:00:02Z --- FAIL: TestSameName (1.00s)",
                    "2026-07-27T00:00:03Z FAIL\tgithub.com/example/package-b\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["unique_failed_tests"], 2)
        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(
            [item["package"] for item in ledger["tests"]],
            [
                "github.com/example/package-a",
                "github.com/example/package-b",
            ],
        )
        self.assertEqual(
            [item["package_blocks"][0] for item in ledger["tests"]],
            [
                {"start_line": 1, "end_line": 2},
                {"start_line": 3, "end_line": 4},
            ],
        )

    def test_ignores_prose_that_looks_like_a_package_completion(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("prose"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestProse (1.00s)",
                    "2026-07-27T00:00:01Z diagnostic: ok response",
                    "2026-07-27T00:00:02Z FAIL\tgithub.com/example/prose\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["tests"][0]["package"], "github.com/example/prose")
        self.assertEqual(
            ledger["tests"][0]["package_blocks"],
            [{"start_line": 1, "end_line": 3}],
        )

    def test_ignores_quoted_test_failure_text(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("quoted"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z diagnostic quoted "
                    "'--- FAIL: TestQuoted (1.00s)'",
                    "2026-07-27T00:00:01Z "
                    "FAIL\tgithub.com/example/quoted\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertEqual(ledger["tests"], [])

    def test_deadline_panic_and_partial_log_make_ledger_incomplete(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("deadline"), "conclusion": "failure"},
            {"id": 2, "name": self._test_job_name("partial"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z panic: test timed out after 5h0m0s",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/deadline\t18000.000s",
                ]
            ),
        )
        self.write_log(
            2,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestPartial (1.00s)",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/partial\t1.000s",
                ]
            ),
            complete=False,
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["deadline_panic_jobs"], 1)
        self.assertEqual(ledger["summary"]["unavailable_test_job_logs"], 1)
        self.assertEqual(len(ledger["partial_log_jobs"]), 1)

    def test_recognizes_build_failed_package_marker(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("build"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "2026-07-27T00:00:00Z FAIL\tgithub.com/example/build [build failed]",
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["package_failure_markers"], 1)
        self.assertEqual(
            ledger["package_failure_markers"][0]["package"],
            "github.com/example/build",
        )

    def test_preserves_parentheses_in_subtest_names(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("parentheses"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestParent/(east) (1.00s)",
                    "2026-07-27T00:00:01Z --- FAIL: TestParent/(west) (1.00s)",
                    "2026-07-27T00:00:02Z --- FAIL: TestParent (2.00s)",
                    "2026-07-27T00:00:03Z FAIL\tgithub.com/example/parentheses\t2.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(
            ledger["tests"][0]["subtests"],
            ["TestParent/(east)", "TestParent/(west)"],
        )

    def test_cleanup_log_issues_cancellations_and_active_jobs_are_visible(
        self,
    ) -> None:
        self.write_jobs(
            {"id": 1, "name": "clean-before / dev", "conclusion": "failure"},
            {"id": 2, "name": "clean-after / dev", "conclusion": "cancelled"},
            {"id": 3, "name": self._test_job_name("running"), "conclusion": None},
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["cleanup_job_failures"], 1)
        self.assertEqual(ledger["summary"]["cleanup_jobs_with_unavailable_logs"], 2)
        self.assertEqual(ledger["summary"]["cancelled_cleanup_job_failures"], 1)
        self.assertEqual(ledger["summary"]["active_workflow_jobs"], 1)

    def test_timed_out_job_is_observed_but_never_complete(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("timeout"),
                "conclusion": "timed_out",
                "started_at": "2026-07-27T00:00:00Z",
                "completed_at": "2026-07-27T01:01:59Z",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestObserved (1.00s)",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/timeout\t1.000s",
                ]
            ),
            complete=False,
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 1)
        self.assertEqual(ledger["summary"]["timed_out_job_failures"], 1)
        self.assertEqual(ledger["summary"]["unavailable_test_job_logs"], 1)
        self.assertEqual(ledger["timed_out_jobs"][0]["job_id"], "1")
        self.assertEqual(ledger["timed_out_jobs"][0]["duration_minutes"], 61)
        self.assertEqual(ledger["partial_log_jobs"][0]["job_id"], "1")

    def test_unsupported_test_conclusion_fails_closed_without_parsing(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("startup"),
                "conclusion": "startup_failure",
                "started_at": "2026-07-27T00:00:00Z",
                "completed_at": "2026-07-27T00:05:59Z",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestMustNotCount (1.00s)",
                    "2026-07-27T00:00:01Z "
                    "FAIL\tgithub.com/example/startup\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["test_jobs_ran"], 0)
        self.assertEqual(ledger["summary"]["test_execution_state"], "unknown")
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertEqual(ledger["summary"]["failed_test_jobs"], 0)
        self.assertEqual(ledger["summary"]["test_job_logs_for_sweep"], 1)
        self.assertEqual(
            ledger["summary"]["unsupported_test_job_conclusions"],
            1,
        )
        self.assertEqual(
            ledger["test_job_logs_for_sweep"][0]["job_id"],
            "1",
        )
        self.assertEqual(ledger["unsupported_test_jobs"][0]["job_id"], "1")
        self.assertEqual(
            ledger["unsupported_test_jobs"][0]["duration_minutes"],
            5,
        )

    def test_failed_test_job_before_test_execution_is_unknown(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("checkout"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "2026-07-27T00:00:00Z Error: checkout diagnostic mentioned "
            "=== RUN   TestStarted, but tests never started",
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["test_jobs_ran"], 0)
        self.assertEqual(ledger["summary"]["test_execution_state"], "unknown")
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["unresolved_failed_jobs"][0]["job_id"], "1")

        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1103,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = run["run_id"]
        ledger["run_attempt"] = run["run_attempt"]
        analysis = build_analysis_input(ledger, run, self.root)
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )

        self.assertEqual(result["verdict"], "red")
        self.assertEqual(result["headline"], "SUITE UNVERIFIED")

    def test_partial_failed_test_job_with_marker_confirms_execution(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("abrupt-failure"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z === RUN   TestStarted",
                    "2026-07-27T00:00:01Z process exited unexpectedly",
                ]
            ),
            complete=False,
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["test_jobs_ran"], 1)
        self.assertEqual(
            ledger["summary"]["test_execution_state"],
            "confirmed_ran",
        )
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertFalse(ledger["summary"]["ledger_complete"])

    def test_partial_deadline_or_teardown_proves_test_execution(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("deadline"),
                "conclusion": "failure",
            },
            {
                "id": 2,
                "name": self._test_job_name("teardown"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "2026-07-27T00:00:00Z panic: test timed out after 5h0m0s",
            complete=False,
        )
        self.write_log(
            2,
            "2026-07-27T00:00:00Z 2026/07/27 00:00:00 "
            "[ERROR] Cleanup failed: failed to delete shared resources",
            complete=False,
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertEqual(ledger["summary"]["test_jobs_ran"], 2)
        self.assertEqual(
            ledger["summary"]["test_execution_state"],
            "confirmed_ran",
        )
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertFalse(ledger["summary"]["ledger_complete"])

    def test_separates_setup_cleanup_reporting_and_test_roles(self) -> None:
        prefix = "1.15.x-latest-pak / tests-1.15.x-latest-dev"
        self.write_jobs(
            {"id": 1, "name": "variables", "conclusion": "success"},
            {
                "id": 2,
                "name": f"{prefix} / get-provider-version",
                "conclusion": "failure",
            },
            {
                "id": 3,
                "name": f"{prefix} / advanced_cluster",
                "conclusion": "skipped",
            },
            {
                "id": 4,
                "name": "clean-before / cleanup-test-env-general",
                "conclusion": "failure",
            },
            {"id": 5, "name": "trigger-test-summary", "conclusion": "failure"},
            {
                "id": 6,
                "name": f"{prefix} / slack-notification-stream",
                "conclusion": "skipped",
            },
        )
        self.write_log(2, "2026-07-27T00:00:00Z setup failed")
        self.write_log(4, "2026-07-27T00:00:00Z cleanup failed")
        self.write_log(5, "2026-07-27T00:00:00Z summary failed")

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertTrue(ledger["summary"]["non_reporting_logs_complete"])
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertEqual(ledger["summary"]["test_jobs"], 1)
        self.assertEqual(ledger["summary"]["setup_jobs"], 2)
        self.assertEqual(ledger["summary"]["cleanup_jobs"], 1)
        self.assertEqual(ledger["summary"]["reporting_jobs"], 2)
        self.assertEqual(ledger["summary"]["skipped_test_jobs"], 1)
        self.assertEqual(ledger["summary"]["setup_job_failures"], 1)
        self.assertEqual(ledger["summary"]["cleanup_job_failures"], 1)
        self.assertEqual(ledger["summary"]["reporting_job_failures"], 1)
        self.assertEqual(ledger["summary"]["failed_test_jobs"], 0)
        self.assertEqual(ledger["summary"]["unresolved_failed_test_jobs"], 0)
        self.assertEqual(ledger["setup_jobs"][0]["job_id"], "2")
        self.assertEqual(ledger["reporting_jobs"][0]["job_id"], "5")

    def test_unknown_failed_job_is_not_parsed_as_a_test(self) -> None:
        self.write_jobs(
            {"id": 1, "name": "unexpected infrastructure job", "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestMustNotCount (1.00s)",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/unknown\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertEqual(ledger["summary"]["unclassified_job_failures"], 1)
        self.assertEqual(ledger["unclassified_jobs"][0]["job_id"], "1")
        self.assertEqual(ledger["failed_unclassified_jobs"][0]["job_id"], "1")

    def test_successful_and_skipped_unknown_jobs_fail_closed(self) -> None:
        self.write_jobs(
            {"id": 1, "name": "renamed successful job", "conclusion": "success"},
            {"id": 2, "name": "renamed skipped job", "conclusion": "skipped"},
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unclassified_jobs"], 2)
        self.assertEqual(ledger["summary"]["unclassified_job_failures"], 0)
        self.assertEqual(
            [item["job_id"] for item in ledger["unclassified_jobs"]],
            ["2", "1"],
        )
        self.assertEqual(ledger["failed_unclassified_jobs"], [])

    def test_active_reporting_job_does_not_contaminate_test_completeness(
        self,
    ) -> None:
        self.write_jobs(
            {"id": 1, "name": "trigger-test-summary", "conclusion": None},
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["active_reporting_jobs"], 1)
        self.assertEqual(ledger["summary"]["active_test_jobs"], 0)

    def test_cancelled_setup_job_caps_non_reporting_completeness(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": (
                    "1.15.x-latest-pak / tests-1.15.x-latest-dev "
                    "/ get-provider-version"
                ),
                "conclusion": "cancelled",
            },
        )
        self.write_log(1, "2026-07-27T00:00:00Z setup was cancelled")

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertFalse(ledger["summary"]["non_reporting_logs_complete"])
        self.assertEqual(ledger["summary"]["cancelled_setup_job_failures"], 1)

    def test_records_explicit_package_teardown_without_changing_test_count(
        self,
    ) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("package-a"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z --- FAIL: TestAccOne (1.00s)",
                    "2026-07-27T00:00:01Z deleted project[ERROR] Cleanup failed:",
                    "2026-07-27T00:00:02Z FAIL\tgithub.com/example/package-a\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 1)
        self.assertEqual(ledger["summary"]["package_teardowns"], 1)
        self.assertEqual(ledger["summary"]["package_teardown_marker_lines"], 1)
        self.assertEqual(
            ledger["package_teardowns"][0],
            {
                "job_id": "1",
                "job_name": self._test_job_name("package-a"),
                "package": "github.com/example/package-a",
                "occurrences": 1,
                "cleanup_line_numbers": [2],
                "package_failure_line_numbers": [3],
                "package_blocks": [{"start_line": 1, "end_line": 3}],
            },
        )

    def test_collapses_repeated_teardowns_per_job_and_package(self) -> None:
        self.write_jobs(
            {"id": 1, "name": self._test_job_name("first"), "conclusion": "failure"},
            {"id": 2, "name": self._test_job_name("second"), "conclusion": "failure"},
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z [ERROR] Cleanup failed:",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/shared\t1.000s",
                    "2026-07-27T00:00:02Z [ERROR] Cleanup failed:",
                    "2026-07-27T00:00:03Z FAIL\tgithub.com/example/shared\t1.000s",
                ]
            ),
        )
        self.write_log(
            2,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z [ERROR] Cleanup failed:",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/shared\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["unique_failed_tests"], 0)
        self.assertEqual(ledger["summary"]["package_teardowns"], 2)
        self.assertEqual(ledger["summary"]["package_teardown_marker_lines"], 3)
        self.assertEqual(ledger["package_teardowns"][0]["occurrences"], 2)
        self.assertEqual(
            ledger["package_teardowns"][0]["package_failure_line_numbers"],
            [2, 4],
        )

    def test_unmatched_package_teardown_marker_is_incomplete(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("unmatched"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z FAIL\tgithub.com/example/plain\t1.000s",
                    "2026-07-27T00:00:01Z [ERROR] Cleanup failed:",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["package_teardowns"], 0)
        self.assertEqual(
            ledger["summary"]["unresolved_package_teardown_markers"],
            1,
        )
        self.assertEqual(
            ledger["unresolved_package_teardown_markers"][0]["reason"],
            "cleanup marker was not followed by a package completion",
        )

    def test_package_teardown_marker_followed_by_ok_is_unresolved(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("contradictory"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z [ERROR] Cleanup failed:",
                    "2026-07-27T00:00:01Z ok  \tgithub.com/example/contradictory\t1.000s",
                ]
            ),
        )

        ledger = build_ledger(self.root, self.jobs_path)

        self.assertFalse(ledger["summary"]["ledger_complete"])
        self.assertEqual(ledger["summary"]["package_teardowns"], 0)
        self.assertEqual(
            ledger["summary"]["unresolved_package_teardown_markers"],
            1,
        )
        self.assertEqual(
            ledger["unresolved_package_teardown_markers"][0]["following_status"],
            "ok",
        )

    def test_prepares_evidence_with_provenance_and_paginated_jobs(self) -> None:
        run_id = 30228140283
        calls: list[tuple[str, bool]] = []
        jobs_pages = [
            {
                "total_count": 2,
                "jobs": [
                    {
                        "id": 101,
                        "name": "variables",
                        "conclusion": "success",
                        "run_id": run_id,
                        "run_attempt": 1,
                        "started_at": "2026-07-27T00:00:00Z",
                        "completed_at": "2026-07-27T00:01:00Z",
                    }
                ]
            },
            {
                "total_count": 2,
                "jobs": [
                    {
                        "id": 202,
                        "name": self._test_job_name("package-a"),
                        "conclusion": "failure",
                        "run_id": run_id,
                        "run_attempt": 1,
                        "started_at": "2026-07-27T00:01:00Z",
                        "completed_at": "2026-07-27T00:02:00Z",
                    }
                ]
            },
        ]

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            calls.append((endpoint, paginate))
            if paginate:
                return self._completed_json(jobs_pages)
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": 1,
                    "status": "completed",
                    "conclusion": "failure",
                }
            )

        def download(endpoint: str, log_path: Path, error_path: Path) -> int:
            self.assertTrue(endpoint.endswith("/actions/jobs/202/logs"))
            log_path.write_text(
                "\n".join(
                    [
                        "2026-07-27T00:01:00Z --- FAIL: TestAccOne (1.00s)",
                        "2026-07-27T00:01:01Z "
                        "FAIL\tgithub.com/example/package-a\t1.000s",
                        "2026-07-27T00:01:02Z Cleaning up orphan processes",
                    ]
                )
                + "\n",
                encoding="utf-8",
            )
            error_path.write_bytes(b"")
            return 0

        compact = prepare_evidence(
            run_id,
            self.root,
            capture_github=capture,
            download_github=download,
        )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        run = json.loads((self.root / "run.json").read_text(encoding="utf-8"))
        ledger = json.loads(
            (self.root / "failure-ledger.json").read_text(encoding="utf-8")
        )
        stored_compact = json.loads(
            (self.root / "evidence-summary.json").read_text(encoding="utf-8")
        )
        analysis_input = json.loads(
            (self.root / "analysis-input.json").read_text(encoding="utf-8")
        )
        classification_input = json.loads(
            (self.root / "model" / "classification-input.json").read_text(
                encoding="utf-8"
            )
        )

        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "run_attempt": 1,
                "status": "ready",
            },
        )
        self.assertEqual(run["run_id"], run_id)
        self.assertEqual(run["run_attempt"], 1)
        self.assertEqual(ledger["run_id"], run_id)
        self.assertEqual(ledger["run_attempt"], 1)
        self.assertEqual(
            ledger["schema"],
            {"name": FAILURE_LEDGER_SCHEMA, "version": SCHEMA_VERSION},
        )
        self.assertEqual(compact["run_attempt"], 1)
        self.assertEqual(
            compact["schema"],
            {"name": EVIDENCE_SUMMARY_SCHEMA, "version": SCHEMA_VERSION},
        )
        self.assertEqual(stored_compact, compact)
        self.assertEqual(
            analysis_input["schema"],
            {"name": ANALYSIS_INPUT_SCHEMA, "version": SCHEMA_VERSION},
        )
        self.assertEqual(
            analysis_input["analysis_digest"],
            _canonical_digest(
                analysis_input,
                omit_key="analysis_digest",
            ),
        )
        self.assertEqual(
            classification_input["schema"],
            {
                "name": CLASSIFICATION_INPUT_SCHEMA,
                "version": SCHEMA_VERSION,
            },
        )
        self.assertEqual(
            classification_input["analysis_digest"],
            analysis_input["analysis_digest"],
        )
        self.assertEqual(classification_input["unit_counts"]["model_tests"], 1)
        self.assertEqual(compact["summary"]["unique_failed_tests"], 1)
        self.assertTrue(compact["summary"]["ledger_complete"])
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "202\n",
        )
        self.assertEqual(len(calls), 4)
        self.assertFalse(calls[0][1])
        self.assertTrue(calls[1][1])
        self.assertTrue(calls[2][1])
        self.assertFalse(calls[3][1])

    def test_incomplete_job_pagination_fails_preparation_closed(self) -> None:
        run_id = 30228140283
        (self.root / "jobs.jsonl").write_text(
            '{"id":999,"run_id":30228140283,"run_attempt":1}\n',
            encoding="utf-8",
        )
        (self.root / "failed-job-ids.txt").write_text("999\n", encoding="utf-8")

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            if paginate:
                return self._completed_json(
                    [
                        {
                            "total_count": 2,
                            "jobs": [
                                {
                                    "id": 101,
                                    "name": "variables",
                                    "conclusion": "success",
                                    "run_id": run_id,
                                    "run_attempt": 2,
                                }
                            ],
                        }
                    ]
                )
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": 2,
                    "status": "completed",
                    "conclusion": "failure",
                }
            )

        def unexpected_download(
            endpoint: str,
            log_path: Path,
            error_path: Path,
        ) -> int:
            self.fail(f"unexpected download: {endpoint} {log_path} {error_path}")

        with self.assertRaisesRegex(
            PreparationError,
            "expected 2, received 1",
        ):
            prepare_evidence(
                run_id,
                self.root,
                capture_github=capture,
                download_github=unexpected_download,
            )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "failed",
            },
        )
        self.assertEqual(
            json.loads((self.root / "run.json").read_text(encoding="utf-8")),
            {"run_id": run_id, "status": "invalidated"},
        )
        self.assertEqual(
            json.loads(
                (
                    self.root / "model" / "classification-input.json"
                ).read_text(encoding="utf-8")
            ),
            {
                "schema": {
                    "name": CLASSIFICATION_INPUT_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "invalidated",
            },
        )
        self.assertEqual(
            (self.root / "jobs.jsonl").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "",
        )

    def test_run_attempt_change_during_preparation_fails_closed(self) -> None:
        run_id = 30228140283
        run_fetches = 0

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            nonlocal run_fetches
            if paginate:
                return self._completed_json(
                    [
                        {
                            "total_count": 1,
                            "jobs": [
                                {
                                    "id": 101,
                                    "name": "variables",
                                    "conclusion": "success",
                                    "run_id": run_id,
                                    "run_attempt": 1,
                                }
                            ],
                        }
                    ]
                )
            run_fetches += 1
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": run_fetches,
                    "status": "completed",
                    "conclusion": "success",
                }
            )

        def unexpected_download(
            endpoint: str,
            log_path: Path,
            error_path: Path,
        ) -> int:
            self.fail(f"unexpected download: {endpoint} {log_path} {error_path}")

        with self.assertRaisesRegex(
            PreparationError,
            "run metadata changed",
        ):
            prepare_evidence(
                run_id,
                self.root,
                capture_github=capture,
                download_github=unexpected_download,
            )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "failed",
            },
        )
        self.assertEqual(
            json.loads((self.root / "run.json").read_text(encoding="utf-8")),
            {"run_id": run_id, "status": "invalidated"},
        )
        self.assertEqual(
            (self.root / "jobs.jsonl").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "",
        )

    def test_job_state_change_during_preparation_fails_closed(self) -> None:
        run_id = 30228140283
        job_fetches = 0

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            nonlocal job_fetches
            if paginate:
                job_fetches += 1
                return self._completed_json(
                    [
                        {
                            "total_count": 1,
                            "jobs": [
                                {
                                    "id": 101,
                                    "name": "variables",
                                    "conclusion": (
                                        "success" if job_fetches == 1 else "failure"
                                    ),
                                    "run_id": run_id,
                                    "run_attempt": 1,
                                }
                            ],
                        }
                    ]
                )
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": 1,
                    "status": "completed",
                    "conclusion": "success",
                }
            )

        def unexpected_download(
            endpoint: str,
            log_path: Path,
            error_path: Path,
        ) -> int:
            self.fail(f"unexpected download: {endpoint} {log_path} {error_path}")

        with self.assertRaisesRegex(
            PreparationError,
            "job metadata changed",
        ):
            prepare_evidence(
                run_id,
                self.root,
                capture_github=capture,
                download_github=unexpected_download,
            )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "failed",
            },
        )
        self.assertEqual(
            json.loads((self.root / "run.json").read_text(encoding="utf-8")),
            {"run_id": run_id, "status": "invalidated"},
        )
        self.assertEqual(
            (self.root / "jobs.jsonl").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "",
        )

    def test_failed_preparation_marks_stale_matching_evidence_unready(self) -> None:
        run_id = 30228140283
        stale_ledger = {
            "run_id": run_id,
            "summary": {"unique_failed_tests": 999},
        }
        (self.root / "failure-ledger.json").write_text(
            json.dumps(stale_ledger),
            encoding="utf-8",
        )
        (self.root / "evidence-summary.json").write_text(
            json.dumps(stale_ledger),
            encoding="utf-8",
        )
        (self.root / "run.json").write_text(
            json.dumps(
                {
                    "run_id": run_id,
                    "run_attempt": 1,
                    "status": "completed",
                }
            ),
            encoding="utf-8",
        )
        (self.root / "jobs.jsonl").write_text(
            '{"id":999,"run_id":30228140283,"run_attempt":1}\n',
            encoding="utf-8",
        )
        (self.root / "failed-job-ids.txt").write_text("999\n", encoding="utf-8")

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            self.assertFalse(paginate)
            self.assertTrue(endpoint.endswith(f"/actions/runs/{run_id}"))
            return subprocess.CompletedProcess(
                ["gh", "api"],
                1,
                stdout=b"",
                stderr=b"HTTP 503",
            )

        def unexpected_download(
            endpoint: str,
            log_path: Path,
            error_path: Path,
        ) -> int:
            self.fail(f"unexpected download: {endpoint} {log_path} {error_path}")

        with self.assertRaises(PreparationError):
            prepare_evidence(
                run_id,
                self.root,
                capture_github=capture,
                download_github=unexpected_download,
            )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "failed",
            },
        )
        self.assertEqual(
            json.loads(
                (self.root / "failure-ledger.json").read_text(encoding="utf-8")
            ),
            stale_ledger,
        )
        self.assertEqual(
            json.loads((self.root / "run.json").read_text(encoding="utf-8")),
            {"run_id": run_id, "status": "invalidated"},
        )
        self.assertEqual(
            (self.root / "jobs.jsonl").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual((self.root / "run.err").read_bytes(), b"HTTP 503")

    def test_download_exception_invalidates_unverified_snapshot(self) -> None:
        run_id = 30228140283
        failed_job = {
            "id": 202,
            "name": self._test_job_name("package-a"),
            "conclusion": "failure",
            "run_id": run_id,
            "run_attempt": 1,
        }

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            if paginate:
                return self._completed_json(
                    [{"total_count": 1, "jobs": [failed_job]}]
                )
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": 1,
                    "status": "completed",
                    "conclusion": "failure",
                }
            )

        def download(endpoint: str, log_path: Path, error_path: Path) -> int:
            raise OSError("simulated download failure")

        with self.assertRaisesRegex(OSError, "simulated download failure"):
            prepare_evidence(
                run_id,
                self.root,
                capture_github=capture,
                download_github=download,
            )

        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )
        self.assertEqual(
            status,
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": run_id,
                "status": "failed",
            },
        )
        self.assertEqual(
            json.loads((self.root / "run.json").read_text(encoding="utf-8")),
            {"run_id": run_id, "status": "invalidated"},
        )
        self.assertEqual(
            (self.root / "jobs.jsonl").read_text(encoding="utf-8"),
            "",
        )
        self.assertEqual(
            (self.root / "failed-job-ids.txt").read_text(encoding="utf-8"),
            "",
        )

    def test_log_download_failure_produces_ready_incomplete_evidence(self) -> None:
        run_id = 30228140283
        failed_job = {
            "id": 202,
            "name": self._test_job_name("package-a"),
            "conclusion": "failure",
            "run_id": run_id,
            "run_attempt": 1,
            "started_at": "2026-07-27T00:01:00Z",
            "completed_at": "2026-07-27T00:02:00Z",
        }

        def capture(endpoint: str, paginate: bool) -> subprocess.CompletedProcess[bytes]:
            if paginate:
                return self._completed_json(
                    [{"total_count": 1, "jobs": [failed_job]}]
                )
            return self._completed_json(
                {
                    "id": run_id,
                    "head_sha": "abc123def456",
                    "run_number": 1102,
                    "run_attempt": 1,
                    "status": "completed",
                    "conclusion": "failure",
                }
            )

        def download(endpoint: str, log_path: Path, error_path: Path) -> int:
            self.assertTrue(endpoint.endswith("/actions/jobs/202/logs"))
            log_path.write_bytes(b"")
            error_path.write_bytes(b"")
            return 1

        compact = prepare_evidence(
            run_id,
            self.root,
            capture_github=capture,
            download_github=download,
        )
        ledger = json.loads(
            (self.root / "failure-ledger.json").read_text(encoding="utf-8")
        )
        status = json.loads(
            (self.root / "preparation-status.json").read_text(encoding="utf-8")
        )

        self.assertEqual(status["status"], "ready")
        self.assertFalse(compact["summary"]["ledger_complete"])
        self.assertEqual(compact["summary"]["unavailable_test_job_logs"], 1)
        self.assertEqual(ledger["download_error_jobs"][0]["job_id"], "202")
        self.assertEqual(
            (self.root / "202.err").read_text(encoding="utf-8"),
            "gh api exited with status 1\n",
        )

    @staticmethod
    def synthetic_analysis(
        test_names: list[str],
        *,
        execution_state: str = "confirmed_ran",
        ledger_complete: bool = True,
        green_eligible: bool = False,
        workflow_units: list[dict[str, object]] | None = None,
        teardown_packages: list[str] | None = None,
        shape_b_unit: str | None = None,
    ) -> dict[str, object]:
        workflow_units = workflow_units or []
        teardown_packages = teardown_packages or []
        test_units = []
        evidence = []
        for index, name in enumerate(test_names):
            unit_id = f"test-{index}"
            evidence_id = f"evidence-{index}"
            test_units.append(
                {
                    "unit_id": unit_id,
                    "kind": "test",
                    "job_id": "101",
                    "job_name": "1.15.x-latest-pak / tests-1.15.x-latest-dev / package",
                    "package": "github.com/example/package",
                    "test": name,
                    "subtests": [],
                    "phase": {
                        "kind": "test_execution",
                        "anchor_lines": [index + 1],
                        "step_anchors": [],
                    },
                    "evidence_refs": [evidence_id],
                    "machine_facts": {
                        "leftover_indicator_observed": False,
                        "shape_b_candidate": unit_id == shape_b_unit,
                        "post_test_destroy": False,
                        "evidence_scope_complete": True,
                        "mixed_phase": False,
                    },
                    "allowed_categories": [
                        "code_regression",
                        "cloud_capacity",
                        "api_contract",
                        "timeout",
                        "cleanup",
                        "api_error",
                        "unresolved",
                    ],
                }
            )
            evidence.append(
                {
                    "evidence_id": evidence_id,
                    "unit_id": unit_id,
                    "kind": "test_failure",
                    "job_id": "101",
                    "line_number": index + 1,
                    "text": f"--- FAIL: {name}",
                }
            )
        teardown_units = [
            {
                "unit_id": f"teardown-{index}",
                "kind": "package_teardown",
                "job_id": str(200 + index),
                "job_name": "test package",
                "package": package,
                "evidence_refs": [],
                "deterministic_category": "cleanup",
            }
            for index, package in enumerate(teardown_packages)
        ]
        analysis: dict[str, object] = {
            "schema": {
                "name": ANALYSIS_INPUT_SCHEMA,
                "version": SCHEMA_VERSION,
            },
            "provenance": {
                "repository": "mongodb/terraform-provider-mongodbatlas",
                "run_id": 30228140283,
                "run_attempt": 1,
                "ledger_schema": {
                    "name": FAILURE_LEDGER_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "ledger_digest": "0" * 64,
            },
            "run_context": {
                "run_id": 30228140283,
                "run_attempt": 1,
                "run_number": 1102,
                "head_sha": "abcdef123456",
                "commit": "abcdef1",
                "status": "completed",
                "conclusion": "success" if green_eligible else "failure",
                "run_url": (
                    "https://github.com/mongodb/"
                    "terraform-provider-mongodbatlas/actions/runs/30228140283"
                ),
                "environment": "dev",
                "auth": "pak",
            },
            "gates": {
                "run_incomplete": False,
                "reporting_only_active": False,
                "test_execution_state": execution_state,
                "successful_test_jobs": 1 if green_eligible else 0,
                "ledger_complete": ledger_complete,
                "non_reporting_logs_complete": ledger_complete,
                "green_eligible": green_eligible,
                "confidence_ceiling": "high" if ledger_complete else "medium",
                "confidence_reasons": (
                    [] if ledger_complete else ["Evidence was incomplete."]
                ),
            },
            "unit_counts": {
                "tests": len(test_units),
                "package_failures": 0,
                "package_teardowns": len(teardown_units),
                "workflow_jobs": len(workflow_units),
            },
            "test_units": test_units,
            "package_failure_units": [],
            "package_teardown_units": teardown_units,
            "workflow_job_units": workflow_units,
            "issues": [],
            "evidence": evidence,
        }
        analysis["analysis_digest"] = _canonical_digest(
            analysis,
            omit_key="analysis_digest",
        )
        return analysis

    @staticmethod
    def model_decisions(
        analysis: dict[str, object],
        groups: list[dict[str, object]],
        **optional: str,
    ) -> dict[str, object]:
        return {
            "schema": {
                "name": MODEL_DECISIONS_SCHEMA,
                "version": SCHEMA_VERSION,
            },
            "analysis_digest": analysis["analysis_digest"],
            "groups": groups,
            **optional,
        }

    def test_analysis_input_marks_test_scoped_post_destroy_phase(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-23T00:57:00Z === RUN   TestPostDestroy",
                    "2026-07-23T00:57:01Z [ERROR] Error running post-test destroy "
                    "test_name=TestPostDestroy",
                    "2026-07-23T00:57:02Z --- FAIL: TestPostDestroy (3.00s)",
                    "2026-07-23T00:57:03Z FAIL\tgithub.com/example/package-a\t3.000s",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = run["run_id"]
        ledger["run_attempt"] = run["run_attempt"]

        analysis = build_analysis_input(ledger, run, self.root)

        self.assertEqual(
            analysis["test_units"][0]["phase"]["kind"],
            "post_test_destroy",
        )
        self.assertTrue(
            analysis["test_units"][0]["machine_facts"]["post_test_destroy"]
        )
        self.assertEqual(analysis["unit_counts"]["package_failures"], 0)
        self.assertEqual(
            analysis["provenance"]["ledger_digest"],
            _canonical_digest(ledger),
        )
        self.assertEqual(
            analysis["analysis_digest"],
            _canonical_digest(analysis, omit_key="analysis_digest"),
        )

    def test_interleaved_destroy_marker_is_not_attributed_to_another_test(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-23T00:57:00Z === RUN   TestOther",
                    "2026-07-23T00:57:01Z [ERROR] Error running post-test destroy",
                    "2026-07-23T00:57:02Z === NAME  TestTarget",
                    "2026-07-23T00:57:03Z --- FAIL: TestTarget (3.00s)",
                    "2026-07-23T00:57:04Z FAIL\tgithub.com/example/package-a\t3.000s",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = run["run_id"]
        ledger["run_attempt"] = run["run_attempt"]

        analysis = build_analysis_input(ledger, run, self.root)

        self.assertEqual(
            analysis["test_units"][0]["phase"]["kind"],
            "test_execution",
        )

    def test_post_destroy_dependency_is_deterministic_cleanup(self) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_AWS_basic"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:57:00Z === RUN   {test_name}",
                (
                    "2026-07-23T00:57:01Z warning issue performing "
                    "authorize: HTTP 400 CANNOT_ASSUME_ROLE"
                ),
                (
                    "2026-07-23T00:57:02Z [ERROR] error deleting "
                    "Encryption At Rest: HTTP 400 "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS diagnostic_summary="
                    '"error when destroying resource"'
                ),
                (
                    "2026-07-23T00:57:03Z [ERROR] "
                    "Error running post-test destroy"
                ),
                (
                    "2026-07-23T00:57:04Z "
                    "Error: error when destroying resource"
                ),
                f"2026-07-23T00:57:05Z --- FAIL: {test_name} (5.00s)",
                (
                    "2026-07-23T00:57:06Z "
                    "FAIL\tgithub.com/example/encryption\t5.000s"
                ),
            ]
        )
        unit = analysis["test_units"][0]
        cohort = build_classification_input(analysis)[
            "deterministic_cohorts"
        ][0]

        self.assertEqual(unit["phase"]["kind"], "post_test_destroy")
        self.assertEqual(unit["phase"]["step_anchors"], [])
        self.assertEqual(unit["allowed_categories"], ["cleanup"])
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 1)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 0)
        self.assertEqual(
            analysis["deterministic_test_groups"][0]["signature"],
            "post_destroy_private_endpoint_dependency_cleanup",
        )
        self.assertEqual(
            cohort["signature"],
            "post_destroy_private_endpoint_dependency_cleanup",
        )
        self.assertTrue(
            any(
                "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS"
                in line
                for line in cohort["representative_evidence"]
            )
        )
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )
        self.assertEqual(result["category_counts"]["cleanup"], 1)
        self.assertEqual(result["category_counts"]["api_contract"], 0)

    def test_post_destroy_dependency_named_conflict_stays_for_model(
        self,
    ) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_conflict"
        lines = [
            f"2026-07-23T00:57:00Z === RUN   {test_name}",
            (
                "2026-07-23T00:57:01Z "
                "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                "PRIVATE_ENDPOINTS"
            ),
            (
                "2026-07-23T00:57:02Z "
                "Error running post-test destroy"
            ),
            (
                "2026-07-23T00:57:03Z "
                "Error: Provider produced an invalid plan in "
                f"{test_name}"
            ),
        ]
        lines.extend(
            (
                f"2026-07-23T00:{index // 60 + 58:02d}:"
                f"{index % 60:02d}Z HTTP 500 secondary signal {index}"
            )
            for index in range(61)
        )
        lines.extend(
            [
                f"2026-07-23T02:00:00Z --- FAIL: {test_name} (63.00s)",
                (
                    "2026-07-23T02:00:01Z "
                    "FAIL\tgithub.com/example/encryption\t63.000s"
                ),
            ]
        )

        analysis = self._analysis_from_failed_test_log(lines)
        classification_input = build_classification_input(analysis)
        visible_evidence = "\n".join(
            item["text"] for item in classification_input["evidence"]
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn(
            f"Provider produced an invalid plan in {test_name}",
            visible_evidence,
        )
        self.assertIn(
            "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS",
            visible_evidence,
        )

    def test_post_destroy_dependency_requires_explicit_harness_anchor(
        self,
    ) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_testStep"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:57:00Z === RUN   {test_name}",
                (
                    "2026-07-23T00:57:01Z [ERROR] "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS diagnostic_summary="
                    '"error when destroying resource"'
                ),
                f"2026-07-23T00:57:02Z --- FAIL: {test_name} (2.00s)",
                (
                    "2026-07-23T00:57:03Z "
                    "FAIL\tgithub.com/example/encryption\t2.000s"
                ),
            ]
        )
        unit = analysis["test_units"][0]

        self.assertEqual(unit["phase"]["kind"], "post_test_destroy")
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)

    def test_post_destroy_dependency_same_line_conflict_stays_for_model(
        self,
    ) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_sameLine"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:57:00Z === RUN   {test_name}",
                (
                    "2026-07-23T00:57:01Z Error: "
                    "Provider produced inconsistent result; "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS; error when destroying resource"
                ),
                (
                    "2026-07-23T00:57:02Z "
                    "Error running post-test destroy"
                ),
                f"2026-07-23T00:57:03Z --- FAIL: {test_name} (3.00s)",
                (
                    "2026-07-23T00:57:04Z "
                    "FAIL\tgithub.com/example/encryption\t3.000s"
                ),
            ]
        )
        classification_input = build_classification_input(analysis)
        visible_evidence = "\n".join(
            item["text"] for item in classification_input["evidence"]
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn("Provider produced inconsistent result", visible_evidence)
        self.assertIn(
            "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS",
            visible_evidence,
        )

    def test_post_destroy_dependency_step_anchor_survives_cap(self) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_step"
        lines = [
            f"2026-07-23T00:57:00Z === RUN   {test_name}",
            "2026-07-23T00:57:01Z Step 1/1 error: apply failed",
        ]
        lines.extend(
            (
                f"2026-07-23T01:{index // 60:02d}:"
                f"{index % 60:02d}Z HTTP 500 secondary signal {index}"
            )
            for index in range(61)
        )
        lines.extend(
            [
                (
                    "2026-07-23T02:00:00Z "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS"
                ),
                (
                    "2026-07-23T02:00:01Z "
                    "Error running post-test destroy"
                ),
                f"2026-07-23T02:00:02Z --- FAIL: {test_name} (63.00s)",
                (
                    "2026-07-23T02:00:03Z "
                    "FAIL\tgithub.com/example/encryption\t63.000s"
                ),
            ]
        )

        analysis = self._analysis_from_failed_test_log(lines)
        unit = analysis["test_units"][0]
        visible_evidence = "\n".join(
            item["text"] for item in analysis["evidence"]
        )

        self.assertEqual(unit["phase"]["kind"], "post_test_destroy")
        self.assertEqual(len(unit["phase"]["step_anchors"]), 1)
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn("Step 1/1 error", visible_evidence)
        self.assertIn(
            "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS",
            visible_evidence,
        )

    def test_post_destroy_dependency_marker_survives_primary_cap(
        self,
    ) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_primaryCap"
        lines = [
            f"2026-07-23T00:57:00Z === RUN   {test_name}",
            (
                "2026-07-23T00:57:01Z "
                "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                "PRIVATE_ENDPOINTS"
            ),
            (
                "2026-07-23T00:57:02Z "
                "Error running post-test destroy"
            ),
        ]
        lines.extend(
            (
                f"2026-07-23T01:{index // 60:02d}:"
                f"{index % 60:02d}Z Error: "
                f"Provider produced an invalid plan {index}"
            )
            for index in range(61)
        )
        lines.extend(
            [
                f"2026-07-23T02:00:00Z --- FAIL: {test_name} (63.00s)",
                (
                    "2026-07-23T02:00:01Z "
                    "FAIL\tgithub.com/example/encryption\t63.000s"
                ),
            ]
        )

        analysis = self._analysis_from_failed_test_log(lines)
        classification_input = build_classification_input(analysis)
        visible_evidence = "\n".join(
            item["text"] for item in classification_input["evidence"]
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn(
            "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS",
            visible_evidence,
        )
        self.assertIn("Provider produced an invalid plan", visible_evidence)

    def test_adjacent_post_destroy_dependency_does_not_force_cleanup(
        self,
    ) -> None:
        analysis = self._analysis_from_failed_test_log(
            [
                "2026-07-23T00:57:00Z === RUN   TestNeighbor",
                (
                    "2026-07-23T00:57:01Z "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS"
                ),
                (
                    "2026-07-23T00:57:02Z "
                    "Error running post-test destroy"
                ),
                "2026-07-23T00:57:03Z === NAME  TestTarget",
                (
                    "2026-07-23T00:57:04Z "
                    "Error: Provider produced inconsistent result"
                ),
                "2026-07-23T00:57:05Z --- FAIL: TestTarget (5.00s)",
                (
                    "2026-07-23T00:57:06Z "
                    "FAIL\tgithub.com/example/encryption\t5.000s"
                ),
            ]
        )
        unit = analysis["test_units"][0]
        visible_evidence = "\n".join(
            item["text"]
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        )

        self.assertEqual(unit["phase"]["kind"], "test_execution")
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertNotIn(
            "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_PRIVATE_ENDPOINTS",
            visible_evidence,
        )

    def test_repeated_top_level_failures_do_not_force_cleanup(self) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_repeated"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:57:00Z === RUN   {test_name}",
                (
                    "2026-07-23T00:57:01Z "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS"
                ),
                (
                    "2026-07-23T00:57:02Z "
                    "Error running post-test destroy"
                ),
                f"2026-07-23T00:57:03Z --- FAIL: {test_name} (3.00s)",
                f"2026-07-23T00:57:04Z === RUN   {test_name}",
                (
                    "2026-07-23T00:57:05Z "
                    "Error: Received unexpected error"
                ),
                f"2026-07-23T00:57:06Z --- FAIL: {test_name} (2.00s)",
                (
                    "2026-07-23T00:57:07Z "
                    "FAIL\tgithub.com/example/encryption\t5.000s"
                ),
            ]
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)

    def test_split_subtest_cleanup_evidence_stays_for_model(self) -> None:
        test_name = "TestAccEncryptionAtRestPrivateEndpoint_subtests"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:57:00Z === RUN   {test_name}",
                f"2026-07-23T00:57:01Z === RUN   {test_name}/marker",
                (
                    "2026-07-23T00:57:02Z "
                    "CANNOT_DISABLE_ENCRYPTION_AT_REST_DUE_TO_"
                    "PRIVATE_ENDPOINTS"
                ),
                f"2026-07-23T00:57:03Z === NAME  {test_name}/marker",
                (
                    f"2026-07-23T00:57:04Z --- FAIL: "
                    f"{test_name}/marker (1.00s)"
                ),
                f"2026-07-23T00:57:05Z === RUN   {test_name}/anchor",
                (
                    "2026-07-23T00:57:06Z "
                    "Error running post-test destroy"
                ),
                f"2026-07-23T00:57:07Z === NAME  {test_name}/anchor",
                (
                    f"2026-07-23T00:57:08Z --- FAIL: "
                    f"{test_name}/anchor (1.00s)"
                ),
                f"2026-07-23T00:57:09Z --- FAIL: {test_name} (2.00s)",
                (
                    "2026-07-23T00:57:10Z "
                    "FAIL\tgithub.com/example/encryption\t2.000s"
                ),
            ]
        )
        unit = analysis["test_units"][0]

        self.assertEqual(len(unit["subtests"]), 2)
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)

    def test_analysis_scopes_parallel_evidence_and_shape_b(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        target = (
            "TestAccAdvancedCluster_"
            "createTimeoutWithDeleteOnCreateReplicaset"
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z === NAME  TestNeighbor",
                    "2026-07-27T00:00:01Z INVALID_ATTRIBUTE from neighbor",
                    f"2026-07-27T00:00:02Z === NAME  {target}",
                    "2026-07-27T00:00:03Z Error Trace: owned.go:10",
                    "2026-07-27T00:00:04Z Error: Should be false",
                    f"2026-07-27T00:00:05Z Test: {target}",
                    "2026-07-27T00:00:06Z === CONT  TestNeighborAfter",
                    "2026-07-27T00:00:07Z HTTP 409 from neighbor",
                    f"2026-07-27T00:00:08Z --- FAIL: {target} (3.00s)",
                    "2026-07-27T00:00:09Z FAIL\tgithub.com/example/package-a\t3.000s",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        unit = analysis["test_units"][0]
        owned_text = "\n".join(
            item["text"]
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        )

        self.assertIn("Error Trace: owned.go:10", owned_text)
        self.assertIn("Error: Should be false", owned_text)
        self.assertNotIn("INVALID_ATTRIBUTE from neighbor", owned_text)
        self.assertNotIn("HTTP 409 from neighbor", owned_text)
        self.assertTrue(unit["machine_facts"]["shape_b_candidate"])
        self.assertTrue(unit["machine_facts"]["evidence_scope_complete"])

    def test_analysis_marks_pak_later_step_authorization_candidate(self) -> None:
        lines = [
            "2026-07-23T00:00:00Z === RUN   TestAccProjectAccessList",
            "2026-07-23T00:00:01Z Step 2/2 error: apply failed",
            (
                "2026-07-23T00:00:02Z Error: "
                "https://cloud-dev.mongodb.com/api/atlas/v2/groups/one/"
                "accessList/10.0.0.1 DELETE: HTTP 401 Unauthorized"
            ),
            (
                "2026-07-23T00:00:03Z Error: "
                "https://cloud-dev.mongodb.com/api/atlas/v2/groups/one/"
                "accessList/10.0.0.2 DELETE: HTTP 401 Unauthorized"
            ),
            (
                "2026-07-23T00:00:04Z "
                "--- FAIL: TestAccProjectAccessList (2.00s)"
            ),
            (
                "2026-07-23T00:00:05Z "
                "FAIL\tgithub.com/example/project\t2.000s"
            ),
        ]

        pak_analysis = self._analysis_from_failed_test_log(
            lines,
            leaf_name="project",
            auth="pak",
        )
        sa_analysis = self._analysis_from_failed_test_log(
            lines,
            leaf_name="project",
            auth="sa",
        )

        pak_facts = pak_analysis["test_units"][0]["machine_facts"]
        sa_facts = sa_analysis["test_units"][0]["machine_facts"]
        self.assertTrue(pak_facts["pak_authorization_failure_candidate"])
        self.assertEqual(
            pak_facts["authorization_failure_distinct_requests"],
            2,
        )
        self.assertFalse(sa_facts["pak_authorization_failure_candidate"])
        self.assertEqual(
            sa_facts["authorization_failure_distinct_requests"],
            2,
        )

    def test_analysis_keeps_owned_terminal_go_source_context(self) -> None:
        analysis = self._analysis_from_failed_test_log(
            [
                "2026-07-23T00:00:00Z === RUN   TestNeighbor",
                (
                    "2026-07-23T00:00:01Z neighbor_test.go:12: "
                    "neighbor failure must not be attributed"
                ),
                "2026-07-23T00:00:02Z === RUN   TestAccNetworkLogging",
                (
                    "1.15.x / tests-dev / config\tUNKNOWN STEP\t"
                    "2026-07-23T00:00:03Z     transport_test.go:156: "
                    "temporary forced failure for test-suite summary validation"
                ),
                (
                    "2026-07-23T00:00:04Z "
                    "--- FAIL: TestAccNetworkLogging (0.00s)"
                ),
                (
                    "2026-07-23T00:00:05Z "
                    "FAIL\tgithub.com/example/config\t0.000s"
                ),
            ],
            leaf_name="config",
        )
        unit = analysis["test_units"][0]
        owned_evidence = [
            item
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        ]
        owned_text = "\n".join(
            item["text"]
            for item in owned_evidence
        )

        self.assertEqual(
            [item["kind"] for item in owned_evidence].count("source_context"),
            1,
        )
        self.assertIn("transport_test.go:156", owned_text)
        self.assertNotIn("neighbor_test.go:12", owned_text)
        self.assertTrue(unit["machine_facts"]["evidence_scope_complete"])
        self.assertEqual(unit["allowed_categories"], ["unresolved"])

    def test_source_context_stays_out_of_deterministic_cohort(self) -> None:
        test_name = "TestAccProject_sourceContext"
        analysis = self._analysis_from_failed_test_log(
            [
                f"2026-07-23T00:00:00Z === RUN   {test_name}",
                "2026-07-23T00:00:01Z Step 1/1 error: apply failed",
                (
                    "2026-07-23T00:00:02Z "
                    "POST /api/atlas/v2/groups: HTTP 400 "
                    "MAX_GROUPS_PER_ORG_EXCEEDED"
                ),
                (
                    "2026-07-23T00:00:03Z project_test.go:12: "
                    "terminal source context"
                ),
                f"2026-07-23T00:00:04Z --- FAIL: {test_name} (3.00s)",
                (
                    "2026-07-23T00:00:05Z "
                    "FAIL\tgithub.com/example/project\t3.000s"
                ),
            ],
            leaf_name="project",
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)

    def test_bare_failure_can_only_be_classified_as_unresolved(self) -> None:
        analysis = self._analysis_from_failed_test_log(
            [
                "2026-07-23T00:00:00Z === RUN   TestBare",
                "2026-07-23T00:00:01Z --- FAIL: TestBare (0.00s)",
                (
                    "2026-07-23T00:00:02Z "
                    "FAIL\tgithub.com/example/config\t0.000s"
                ),
            ],
            leaf_name="config",
        )
        unit = analysis["test_units"][0]

        self.assertEqual(unit["allowed_categories"], ["unresolved"])

    def test_analysis_keeps_evidence_from_earlier_failed_subtest(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("streamprocessor"),
                "conclusion": "failure",
            },
        )
        parent = "TestAccStreamProcessor_StateTransitionsUpdates"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {parent}",
                    f"2026-07-23T00:00:01Z === RUN   {parent}/first",
                    "2026-07-23T00:00:02Z HTTP 500 UNEXPECTED_ERROR",
                    f"2026-07-23T00:00:03Z === NAME  {parent}/first",
                    f"2026-07-23T00:00:04Z --- FAIL: {parent}/first (1.00s)",
                    "2026-07-23T00:00:05Z === CONT  TestNeighbor",
                    "2026-07-23T00:00:06Z HTTP 409 from neighbor",
                    f"2026-07-23T00:00:07Z === NAME  {parent}/second",
                    "2026-07-23T00:00:08Z Error: Should be true",
                    f"2026-07-23T00:00:09Z --- FAIL: {parent}/second (1.00s)",
                    f"2026-07-23T00:00:10Z --- FAIL: {parent} (2.00s)",
                    (
                        "2026-07-23T00:00:11Z "
                        "FAIL\tgithub.com/example/streamprocessor\t2.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1093,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        unit = analysis["test_units"][0]
        owned_text = "\n".join(
            item["text"]
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        )

        self.assertIn("HTTP 500 UNEXPECTED_ERROR", owned_text)
        self.assertIn("Error: Should be true", owned_text)
        self.assertNotIn("HTTP 409 from neighbor", owned_text)
        self.assertTrue(unit["machine_facts"]["evidence_scope_complete"])

    def test_analysis_excludes_evidence_from_passing_subtest(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        parent = "TestAccProject_subtests"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {parent}",
                    f"2026-07-23T00:00:01Z === RUN   {parent}/passing",
                    (
                        "2026-07-23T00:00:02Z "
                        "MAX_GROUPS_PER_ORG_EXCEEDED was expected"
                    ),
                    f"2026-07-23T00:00:03Z === NAME  {parent}/failed",
                    "2026-07-23T00:00:04Z INVALID_ATTRIBUTE regression",
                    f"2026-07-23T00:00:05Z --- FAIL: {parent}/failed (1.00s)",
                    f"2026-07-23T00:00:06Z --- FAIL: {parent} (2.00s)",
                    (
                        "2026-07-23T00:00:07Z "
                        "FAIL\tgithub.com/example/project\t2.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        unit = analysis["test_units"][0]
        owned_text = "\n".join(
            item["text"]
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        )

        self.assertNotIn("MAX_GROUPS_PER_ORG_EXCEEDED", owned_text)
        self.assertIn("INVALID_ATTRIBUTE regression", owned_text)
        self.assertNotEqual(unit["allowed_categories"], ["cleanup"])

    def test_max_groups_per_org_is_deterministic_cleanup(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_basic"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {test_name}",
                    "2026-07-23T00:00:01Z Step 1/1 error: apply failed",
                    (
                        "2026-07-23T00:00:02Z "
                        "POST /api/atlas/v2/groups: HTTP 400 "
                        "MAX_GROUPS_PER_ORG_EXCEEDED"
                    ),
                    f"2026-07-23T00:00:03Z --- FAIL: {test_name} (3.00s)",
                    (
                        "2026-07-23T00:00:04Z "
                        "FAIL\tgithub.com/example/project\t3.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        unit = analysis["test_units"][0]

        self.assertEqual(unit["allowed_categories"], ["cleanup"])
        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 1)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 0)
        self.assertEqual(
            analysis["deterministic_test_groups"][0]["signature"],
            "max_groups_per_org_cleanup",
        )
        decisions = self.model_decisions(analysis, [])
        result = finalize_analysis(analysis, decisions)
        self.assertEqual(result["category_counts"]["cleanup"], 1)
        self.assertEqual(result["category_counts"]["api_error"], 0)

        decisions["groups"] = [
            {
                "unit_ids": [unit["unit_id"]],
                "category": "api_error",
                "cause": "Incorrect reclassification.",
                "evidence_refs": [unit["evidence_refs"][0]],
            }
        ]
        with self.assertRaisesRegex(
            FinalizationError,
            "must not reclassify deterministic tests",
        ):
            finalize_analysis(analysis, decisions)

        tampered = json.loads(json.dumps(analysis))
        tampered["deterministic_test_groups"][0]["test_count"] = 2
        tampered["analysis_digest"] = _canonical_digest(
            tampered,
            omit_key="analysis_digest",
        )
        with self.assertRaisesRegex(
            FinalizationError,
            "deterministic group is malformed",
        ):
            finalize_analysis(
                tampered,
                self.model_decisions(tampered, []),
            )

    def test_max_groups_survives_the_context_evidence_cap(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_contextCap"
        lines = [
            f"2026-07-23T00:00:00Z === RUN   {test_name}",
            (
                "2026-07-23T00:00:01Z "
                "POST /api/atlas/v2/groups: HTTP 400 "
                "MAX_GROUPS_PER_ORG_EXCEEDED"
            ),
        ]
        lines.extend(
            (
                f"2026-07-23T00:{index // 60 + 1:02d}:{index % 60:02d}Z "
                "POST /api/atlas/v2/groups: HTTP 400 "
                f"MAX_GROUPS_PER_ORG_EXCEEDED duplicate {index}"
            )
            for index in range(61)
        )
        lines.extend(
            [
                f"2026-07-23T00:02:02Z --- FAIL: {test_name} (62.00s)",
                (
                    "2026-07-23T00:02:03Z "
                    "FAIL\tgithub.com/example/project\t62.000s"
                ),
            ]
        )
        self.write_log(1, "\n".join(lines))
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        unit = analysis["test_units"][0]
        owned_text = "\n".join(
            item["text"]
            for item in analysis["evidence"]
            if item["unit_id"] == unit["unit_id"]
        )

        self.assertEqual(unit["allowed_categories"], ["cleanup"])
        self.assertIn("MAX_GROUPS_PER_ORG_EXCEEDED", owned_text)

    def test_conflict_beyond_evidence_cap_stays_for_model(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_contextConflict"
        lines = [
            f"2026-07-23T00:00:00Z === RUN   {test_name}",
            (
                "2026-07-23T00:00:01Z "
                "POST /api/atlas/v2/groups: HTTP 400 "
                "MAX_GROUPS_PER_ORG_EXCEEDED"
            ),
            "2026-07-23T00:00:02Z INVALID_ATTRIBUTE provider regression",
        ]
        lines.extend(
            (
                f"2026-07-23T00:{index // 60 + 1:02d}:{index % 60:02d}Z "
                f"HTTP 500 nearer signal {index}"
            )
            for index in range(61)
        )
        lines.extend(
            [
                f"2026-07-23T00:02:02Z --- FAIL: {test_name} (62.00s)",
                (
                    "2026-07-23T00:02:03Z "
                    "FAIL\tgithub.com/example/project\t62.000s"
                ),
            ]
        )
        self.write_log(1, "\n".join(lines))
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        classification_input = build_classification_input(analysis)
        visible_evidence = "\n".join(
            item["text"] for item in classification_input["evidence"]
        )

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn("INVALID_ATTRIBUTE", visible_evidence)
        self.assertIn("MAX_GROUPS_PER_ORG_EXCEEDED", visible_evidence)

    def test_remaining_selector_covers_only_model_classification_units(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-23T00:00:00Z === RUN   TestDeterministic",
                    (
                        "2026-07-23T00:00:01Z "
                        "MAX_GROUPS_PER_ORG_EXCEEDED"
                    ),
                    (
                        "2026-07-23T00:00:02Z "
                        "--- FAIL: TestDeterministic (1.00s)"
                    ),
                    "2026-07-23T00:00:03Z === RUN   TestModel",
                    (
                        "2026-07-23T00:00:04Z "
                        "Error: Provider produced inconsistent result"
                    ),
                    "2026-07-23T00:00:05Z --- FAIL: TestModel (1.00s)",
                    (
                        "2026-07-23T00:00:06Z "
                        "FAIL\tgithub.com/example/project\t2.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1
        analysis = build_analysis_input(ledger, run, self.root)
        model_unit = next(
            unit
            for unit in analysis["test_units"]
            if unit["test"] == "TestModel"
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "code_regression",
                    "cause": "Provider produced inconsistent state.",
                    "evidence_refs": [model_unit["evidence_refs"][1]],
                }
            ],
        )

        result = finalize_analysis(analysis, decisions)

        self.assertEqual(result["test_total"], 2)
        self.assertEqual(result["category_counts"]["cleanup"], 1)
        self.assertEqual(result["category_counts"]["code_regression"], 1)

    def test_adjacent_max_groups_does_not_force_cleanup(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-23T00:00:00Z === RUN   TestNeighbor",
                    (
                        "2026-07-23T00:00:01Z "
                        "POST /api/atlas/v2/groups: HTTP 400 "
                        "MAX_GROUPS_PER_ORG_EXCEEDED"
                    ),
                    "2026-07-23T00:00:02Z === NAME  TestTarget",
                    "2026-07-23T00:00:03Z Step 1/1 error: apply failed",
                    "2026-07-23T00:00:04Z HTTP 500 UNEXPECTED_ERROR",
                    "2026-07-23T00:00:05Z --- FAIL: TestTarget (3.00s)",
                    (
                        "2026-07-23T00:00:06Z "
                        "FAIL\tgithub.com/example/project\t3.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)

        self.assertEqual(
            analysis["test_units"][0]["allowed_categories"],
            [
                "code_regression",
                "cloud_capacity",
                "api_contract",
                "timeout",
                "cleanup",
                "api_error",
                "unresolved",
            ],
        )

    def test_groups_post_http_500_is_a_visible_deterministic_cohort(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_create"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {test_name}",
                    "2026-07-23T00:00:01Z Step 1/1 error: apply failed",
                    (
                        "2026-07-23T00:00:02Z "
                        "https://cloud-dev.mongodb.com/api/atlas/v2/groups "
                        "POST: HTTP 500 Internal Server Error"
                    ),
                    (
                        "2026-07-23T00:00:03Z Messages: "
                        "Project creation failed: redacted"
                    ),
                    f"2026-07-23T00:00:04Z --- FAIL: {test_name} (3.00s)",
                    (
                        "2026-07-23T00:00:05Z "
                        "FAIL\tgithub.com/example/project\t3.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        classification_input = build_classification_input(analysis)

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 1)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 0)
        self.assertEqual(
            analysis["deterministic_test_groups"][0]["signature"],
            "groups_post_http_500",
        )
        self.assertEqual(classification_input["test_units"], [])
        cohort = classification_input["deterministic_cohorts"][0]
        self.assertEqual(cohort["unique_tests"], 1)
        self.assertEqual(cohort["jobs"], 1)
        self.assertEqual(cohort["packages"], 1)
        self.assertEqual(cohort["unattributed_tests"], 0)
        self.assertTrue(
            any(
                "/api/atlas/v2/groups POST: HTTP 500" in line
                for line in cohort["representative_evidence"]
            )
        )
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )
        self.assertEqual(result["category_counts"]["api_error"], 1)

    def test_groups_post_http_500_with_provider_assertion_stays_for_model(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_create"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {test_name}",
                    (
                        "2026-07-23T00:00:01Z "
                        "https://cloud-dev.mongodb.com/api/atlas/v2/groups "
                        "POST: HTTP 500 Internal Server Error"
                    ),
                    (
                        "2026-07-23T00:00:02Z Messages: "
                        "Project creation failed: redacted"
                    ),
                    "2026-07-23T00:00:03Z Error: Not equal:",
                    "2026-07-23T00:00:04Z actual: unexpected provider state",
                    f"2026-07-23T00:00:05Z --- FAIL: {test_name} (3.00s)",
                    (
                        "2026-07-23T00:00:06Z "
                        "FAIL\tgithub.com/example/project\t3.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertEqual(analysis["deterministic_test_groups"], [])

    def test_earlier_groups_500_does_not_hide_later_invalid_plan(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_create"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {test_name}",
                    (
                        "2026-07-23T00:00:01Z "
                        "https://cloud-dev.mongodb.com/api/atlas/v2/groups "
                        "POST: HTTP 500 Internal Server Error"
                    ),
                    (
                        "2026-07-23T00:00:02Z Messages: "
                        "Project creation failed during an expected error step"
                    ),
                    (
                        "2026-07-23T00:00:03Z "
                        "Step 2/2 error: unexpected plan"
                    ),
                    (
                        "2026-07-23T00:00:04Z Error: "
                        "Provider produced an invalid plan"
                    ),
                    f"2026-07-23T00:00:05Z --- FAIL: {test_name} (3.00s)",
                    (
                        "2026-07-23T00:00:06Z "
                        "FAIL\tgithub.com/example/project\t3.000s"
                    ),
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        classification_input = build_classification_input(analysis)

        self.assertEqual(analysis["unit_counts"]["deterministic_tests"], 0)
        self.assertEqual(analysis["unit_counts"]["model_tests"], 1)
        self.assertIn(
            "Provider produced an invalid plan",
            "\n".join(
                item["text"] for item in classification_input["evidence"]
            ),
        )

    def test_deterministic_cohort_counts_only_attributed_packages(
        self,
    ) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("project"),
                "conclusion": "failure",
            },
        )
        test_name = "TestAccProject_unattributed"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-23T00:00:00Z === RUN   {test_name}",
                    (
                        "2026-07-23T00:00:01Z "
                        "MAX_GROUPS_PER_ORG_EXCEEDED"
                    ),
                    f"2026-07-23T00:00:02Z --- FAIL: {test_name} (1.00s)",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1098,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        cohort = build_classification_input(analysis)[
            "deterministic_cohorts"
        ][0]

        self.assertEqual(cohort["packages"], 0)
        self.assertEqual(cohort["unattributed_tests"], 1)

    def test_shape_b_is_not_inferred_from_test_identity(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        name = "TestCreateTimeoutWithDeleteOnCreateTimeout"
        self.write_log(
            1,
            "\n".join(
                [
                    f"2026-07-27T00:00:00Z === RUN   {name}",
                    "2026-07-27T00:00:01Z Error: Received unexpected error",
                    f"2026-07-27T00:00:02Z Test: {name}",
                    f"2026-07-27T00:00:03Z --- FAIL: {name} (2.00s)",
                    "2026-07-27T00:00:04Z FAIL\tgithub.com/example/package-a\t2.000s",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)

        self.assertFalse(
            analysis["test_units"][0]["machine_facts"][
                "shape_b_candidate"
            ]
        )

    def test_package_only_failures_are_reportable_and_never_green(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("successful-package"),
                "conclusion": "success",
            },
            {
                "id": 2,
                "name": self._test_job_name("build-package"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            2,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z package.go:10: undefined: missing",
                    "2026-07-27T00:00:01Z FAIL\tgithub.com/example/build\t[build failed]",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )

        self.assertEqual(analysis["unit_counts"]["tests"], 0)
        self.assertEqual(analysis["unit_counts"]["package_failures"], 1)
        self.assertFalse(analysis["gates"]["green_eligible"])
        self.assertEqual(result["headline"], "CODE REGRESSION DETECTED")
        self.assertEqual(result["test_total"], 0)
        self.assertIn("Package build failure", result["slack_mrkdwn"])

    def test_separate_build_failure_block_is_not_hidden_by_a_test(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "\n".join(
                [
                    "2026-07-27T00:00:00Z === RUN   TestAPIError",
                    "2026-07-27T00:00:01Z Error: HTTP 500",
                    "2026-07-27T00:00:02Z --- FAIL: TestAPIError (2.00s)",
                    "2026-07-27T00:00:03Z FAIL\tgithub.com/example/package-a\t2.000s",
                    "2026-07-27T00:00:04Z package.go:10: undefined: missing",
                    "2026-07-27T00:00:05Z FAIL\tgithub.com/example/package-a [build failed]",
                ]
            ),
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned HTTP 500.",
                    "evidence_refs": [
                        analysis["test_units"][0]["evidence_refs"][0]
                    ],
                }
            ],
        )
        result = finalize_analysis(analysis, decisions)

        self.assertTrue(ledger["summary"]["ledger_complete"])
        self.assertEqual(analysis["unit_counts"]["tests"], 1)
        self.assertEqual(analysis["unit_counts"]["package_failures"], 1)
        self.assertEqual(
            analysis["package_failure_units"][0]["deterministic_category"],
            "code_regression",
        )
        self.assertEqual(result["headline"], "CODE REGRESSION DETECTED")

    def test_generic_package_only_failure_requires_review(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "failure",
            },
        )
        self.write_log(
            1,
            "2026-07-27T00:00:00Z FAIL\tgithub.com/example/package-a\t0.100s",
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "failure",
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )

        self.assertFalse(analysis["gates"]["ledger_complete"])
        self.assertEqual(result["headline"], "Automatic triage incomplete")
        self.assertIn("unresolved package failure", result["slack_mrkdwn"])

    def test_analysis_input_requires_matching_ledger_provenance(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "success",
            },
        )
        ledger = build_ledger(self.root, self.jobs_path)
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1
        run = {
            "run_id": 2,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "completed",
            "conclusion": "success",
        }

        with self.assertRaisesRegex(PreparationError, "run identity"):
            build_analysis_input(ledger, run, self.root)

    def test_finalize_groups_seven_timeouts_and_names_shape_b_once(self) -> None:
        analysis = self.synthetic_analysis(
            [f"TestTimeout{index}" for index in range(7)],
            shape_b_unit="test-0",
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "unit_ids": ["test-0"],
                    "category": "timeout",
                    "cause": "Delete wait expired without a leftover indicator.",
                    "evidence_refs": ["evidence-0"],
                    "ambiguity": (
                        "The log cannot distinguish provider cleanup from slow Atlas."
                    ),
                    "note": "delete_on_timeout_unverified",
                },
                {
                    "remaining": True,
                    "category": "timeout",
                    "cause": "The operation exceeded its polling deadline.",
                    "evidence_refs": [
                        f"evidence-{index}" for index in range(1, 7)
                    ],
                },
            ],
        )

        result = finalize_analysis(analysis, decisions)

        self.assertEqual(result["category_counts"]["timeout"], 7)
        self.assertEqual(result["test_total"], 7)
        self.assertEqual(result["headline"], "Infrastructure noise only")
        self.assertEqual(result["confidence"], "medium")
        self.assertIn(
            "(included in Timeout total)",
            result["slack_mrkdwn"],
        )
        self.assertNotIn("Timeout: 8", result["slack_mrkdwn"])

    def test_finalize_rejects_stale_digest_and_mutated_analysis(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned an API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )
        stale = dict(decisions)
        stale["analysis_digest"] = "0" * 64
        with self.assertRaisesRegex(FinalizationError, "do not match"):
            finalize_analysis(analysis, stale)

        mutated = json.loads(json.dumps(analysis))
        mutated["run_context"]["run_attempt"] = 2
        with self.assertRaisesRegex(FinalizationError, "digest"):
            finalize_analysis(mutated, decisions)

    def test_remaining_selector_is_last_and_reconciles_exactly(self) -> None:
        analysis = self.synthetic_analysis(["TestOne", "TestTwo", "TestThree"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "unit_ids": ["test-0"],
                    "category": "api_contract",
                    "cause": "The API contract changed.",
                    "evidence_refs": ["evidence-0"],
                },
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a generic API error.",
                    "evidence_refs": ["evidence-1", "evidence-2"],
                },
            ],
        )

        result = finalize_analysis(analysis, decisions)

        self.assertEqual(result["category_counts"]["api_contract"], 1)
        self.assertEqual(result["category_counts"]["api_error"], 2)
        self.assertIn(
            "API contract: 1 test — The API contract changed.",
            result["slack_mrkdwn"],
        )
        self.assertIn(
            "API errors: 2 tests — Atlas returned a generic API error.",
            result["slack_mrkdwn"],
        )

        decisions["groups"].reverse()
        with self.assertRaisesRegex(FinalizationError, "must be the last"):
            finalize_analysis(analysis, decisions)

    def test_summary_keeps_the_dominant_category_cause_visible(self) -> None:
        analysis = self.synthetic_analysis(
            ["TestOne", "TestTwo", "TestThree", "TestFour", "TestFive"]
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "unit_ids": ["test-0"],
                    "category": "api_error",
                    "cause": "Small API cause one.",
                    "evidence_refs": ["evidence-0"],
                },
                {
                    "unit_ids": ["test-1"],
                    "category": "api_error",
                    "cause": "Small API cause two.",
                    "evidence_refs": ["evidence-1"],
                },
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Dominant deterministic API cause.",
                    "evidence_refs": [
                        "evidence-2",
                        "evidence-3",
                        "evidence-4",
                    ],
                },
            ],
        )

        result = finalize_analysis(analysis, decisions)

        self.assertIn(
            "Dominant deterministic API cause.",
            result["slack_mrkdwn"],
        )
        self.assertIn("Small API cause one.", result["slack_mrkdwn"])
        self.assertNotIn("Small API cause two.", result["slack_mrkdwn"])

    def test_summary_preserves_complete_validated_causes(self) -> None:
        cause = (
            "Step 2/2 apply failed because Atlas rejected every project IP "
            "access list entry DELETE with HTTP 401 Unauthorized "
            '(UNEXPECTED_ERROR, "You are not authorized for this resource"), '
            "which also broke the post-test destroy."
        )
        self.assertGreater(len(cause), 160)
        self.assertLessEqual(len(cause), 240)
        cases = (
            ("code_regression", False),
            ("api_error", False),
            ("unresolved", False),
            ("timeout", True),
        )
        for category, shape_b in cases:
            with self.subTest(category=category, shape_b=shape_b):
                analysis = self.synthetic_analysis(
                    ["TestOne"],
                    shape_b_unit="test-0" if shape_b else None,
                )
                group = {
                    "remaining": True,
                    "category": category,
                    "cause": cause,
                    "evidence_refs": ["evidence-0"],
                }
                if category == "unresolved":
                    group["ambiguity"] = "The evidence is inconclusive."
                if shape_b:
                    group["ambiguity"] = "Deletion completion is ambiguous."
                    group["note"] = "delete_on_timeout_unverified"
                decisions = self.model_decisions(analysis, [group])

                result = finalize_analysis(analysis, decisions)

                self.assertIn(cause, result["slack_mrkdwn"])

    def test_finalize_rejects_category_not_allowed_by_unit(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        analysis["test_units"][0]["allowed_categories"] = ["timeout"]
        analysis["analysis_digest"] = _canonical_digest(
            analysis,
            omit_key="analysis_digest",
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a generic API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )

        with self.assertRaisesRegex(
            FinalizationError,
            "category that is not allowed",
        ):
            finalize_analysis(analysis, decisions)

    def test_finalize_accepts_representative_group_evidence(self) -> None:
        analysis = self.synthetic_analysis(["TestOne", "TestTwo"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned the same API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )

        result = finalize_analysis(analysis, decisions)
        self.assertEqual(result["category_counts"]["api_error"], 2)

    def test_finalize_rejects_evidence_owned_by_another_group(self) -> None:
        analysis = self.synthetic_analysis(["TestOne", "TestTwo"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "unit_ids": ["test-0"],
                    "category": "api_error",
                    "cause": "Atlas returned an API error.",
                    "evidence_refs": ["evidence-1"],
                },
                {
                    "remaining": True,
                    "category": "timeout",
                    "cause": "The operation timed out.",
                    "evidence_refs": ["evidence-1"],
                },
            ],
        )

        with self.assertRaisesRegex(
            FinalizationError,
            "evidence not owned",
        ):
            finalize_analysis(analysis, decisions)

    def test_finalize_requires_shape_b_note(self) -> None:
        analysis = self.synthetic_analysis(
            ["TestCreateTimeout"],
            shape_b_unit="test-0",
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "timeout",
                    "cause": "The delete wait expired.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )

        with self.assertRaisesRegex(
            FinalizationError,
            "Shape-B timeout decisions require",
        ):
            finalize_analysis(analysis, decisions)

    def test_finalize_rejects_malformed_nested_analysis(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        analysis["test_units"] = [None]
        analysis["analysis_digest"] = _canonical_digest(
            analysis,
            omit_key="analysis_digest",
        )

        with self.assertRaisesRegex(FinalizationError, "non-object unit"):
            finalize_analysis(
                analysis,
                self.model_decisions(analysis, []),
            )

    def test_finalize_cli_reads_decisions_from_environment(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a generic API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )
        analysis_path = self.root / "analysis-input.json"
        output_path = self.root / "summary.md"
        result_path = self.root / "finalization-result.json"
        slack_path = self.root / "slack-payload.json"
        analysis_path.write_text(json.dumps(analysis), encoding="utf-8")

        with mock.patch.dict(
            os.environ,
            {
                "MODEL_DECISIONS_JSON": json.dumps(decisions),
                "SLACK_PREFIX": "<!subteam^oncall>",
            },
        ):
            with contextlib.redirect_stdout(io.StringIO()):
                _finalize_main(
                    [
                        "--analysis-input",
                        str(analysis_path),
                        "--decisions-env",
                        "MODEL_DECISIONS_JSON",
                        "--output",
                        str(output_path),
                        "--result-json",
                        str(result_path),
                        "--slack-payload",
                        str(slack_path),
                        "--slack-prefix-env",
                        "SLACK_PREFIX",
                    ]
                )

        self.assertIn(
            "API errors: 1 test",
            output_path.read_text(encoding="utf-8"),
        )
        self.assertEqual(
            json.loads(
                (self.root / "finalization-status.json").read_text(
                    encoding="utf-8",
                )
            )["status"],
            "ready",
        )
        self.assertEqual(
            json.loads(
                (self.root / "model-decisions.json").read_text(
                    encoding="utf-8",
                )
            ),
            decisions,
        )
        self.assertEqual(
            json.loads(result_path.read_text(encoding="utf-8"))["test_total"],
            1,
        )
        slack_payload = json.loads(slack_path.read_text(encoding="utf-8"))
        self.assertEqual(slack_payload["text"], "Test Suite #1102 summary")
        self.assertTrue(
            slack_payload["blocks"][0]["text"]["text"].startswith(
                "<!subteam^oncall> :yellow_circle:"
            )
        )

    def test_workflow_finalize_model_decisions_and_stage_replay(self) -> None:
        env, output_dir = self.workflow_env()
        run_id = int(env["GITHUB_RUN_ID"])
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a generic API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )
        self.write_workflow_analysis(
            output_dir,
            run_id=run_id,
            model_tests=1,
            analysis=analysis,
        )
        env["MODEL_DECISIONS_JSON"] = json.dumps(decisions)

        self.run_workflow_finalize(env)

        summary = (output_dir / "summary.md").read_text(encoding="utf-8")
        self.assertIn("API errors: 1 test", summary)
        self.assertIn(
            "API errors: 1 test",
            Path(env["GITHUB_STEP_SUMMARY"]).read_text(encoding="utf-8"),
        )
        status = self.read_json(output_dir / "finalization-status.json")
        self.assertEqual(status["status"], "ready")
        replay_dir = output_dir / "artifacts" / "summary-replay"
        for relative_name in (
            "model-decisions.json",
            "finalization-result.json",
            "finalization-status.json",
        ):
            with self.subTest(relative_name=relative_name):
                self.assertTrue((replay_dir / relative_name).is_file())

    def test_workflow_finalize_synthesizes_zero_model_decisions(self) -> None:
        env, output_dir = self.workflow_env()
        run_id = int(env["GITHUB_RUN_ID"])
        analysis = self.synthetic_analysis([], green_eligible=True)
        self.write_workflow_analysis(
            output_dir,
            run_id=run_id,
            model_tests=0,
            analysis=analysis,
        )

        self.run_workflow_finalize(env)

        decisions = self.read_json(output_dir / "model-decisions.json")
        self.assertEqual(decisions["analysis_digest"], analysis["analysis_digest"])
        self.assertEqual(decisions["groups"], [])

    def test_workflow_finalize_invalid_model_output_is_diagnostic(self) -> None:
        env, output_dir = self.workflow_env()
        run_id = int(env["GITHUB_RUN_ID"])
        analysis = self.synthetic_analysis(["TestOne"])
        self.write_workflow_analysis(
            output_dir,
            run_id=run_id,
            model_tests=1,
            analysis=analysis,
        )
        env["MODEL_DECISIONS_JSON"] = "{invalid"

        self.run_workflow_finalize(env, fails=True)

        diagnostic = (output_dir / "summary.md").read_text(encoding="utf-8")
        self.assertIn("summary unavailable", diagnostic)
        self.assertNotIn("{invalid", diagnostic)
        self.assertIn(
            "summary unavailable",
            Path(env["GITHUB_STEP_SUMMARY"]).read_text(encoding="utf-8"),
        )
        status = self.read_json(output_dir / "finalization-status.json")
        self.assertEqual(status["status"], "failed")
        self.assertTrue(
            (
                output_dir
                / "artifacts"
                / "summary-replay"
                / "finalization-status.json"
            ).is_file()
        )

    def test_workflow_finalize_preparation_failure_is_diagnostic(self) -> None:
        env, output_dir = self.workflow_env()
        output_dir.mkdir(parents=True, exist_ok=True)
        self.write_json(
            output_dir / "preparation-status.json",
            {
                "schema": {
                    "name": PREPARATION_STATUS_SCHEMA,
                    "version": SCHEMA_VERSION,
                },
                "run_id": int(env["GITHUB_RUN_ID"]),
                "status": "failed",
            },
        )

        self.run_workflow_finalize(env, fails=True)

        self.assertIn(
            "summary unavailable",
            (output_dir / "summary.md").read_text(encoding="utf-8"),
        )
        status = self.read_json(output_dir / "finalization-status.json")
        self.assertEqual(status["status"], "failed")

    def test_finalize_cli_invalidates_outputs_when_ready_status_write_fails(
        self,
    ) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a generic API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )
        analysis_path = self.root / "analysis-input.json"
        output_path = self.root / "summary.md"
        result_path = self.root / "finalization-result.json"
        slack_path = self.root / "slack-payload.json"
        status_path = self.root / "finalization-status.json"
        analysis_path.write_text(json.dumps(analysis), encoding="utf-8")
        original_write_json = failure_ledger._atomic_write_json

        def fail_ready_status(
            path: Path,
            value: object,
            *,
            compact: bool = False,
        ) -> None:
            if (
                path == status_path
                and isinstance(value, dict)
                and value.get("status") == "ready"
            ):
                raise OSError("simulated ready-status write failure")
            original_write_json(path, value, compact=compact)

        with mock.patch.dict(
            os.environ,
            {"MODEL_DECISIONS_JSON": json.dumps(decisions)},
        ):
            with mock.patch(
                "failure_ledger._atomic_write_json",
                side_effect=fail_ready_status,
            ):
                with contextlib.redirect_stderr(io.StringIO()):
                    with self.assertRaises(SystemExit):
                        _finalize_main(
                            [
                                "--analysis-input",
                                str(analysis_path),
                                "--decisions-env",
                                "MODEL_DECISIONS_JSON",
                                "--output",
                                str(output_path),
                                "--result-json",
                                str(result_path),
                                "--slack-payload",
                                str(slack_path),
                            ]
                        )

        self.assertEqual(output_path.read_bytes(), b"")
        self.assertEqual(result_path.read_bytes(), b"")
        self.assertEqual(slack_path.read_bytes(), b"")
        status = json.loads(status_path.read_text(encoding="utf-8"))
        self.assertEqual(status["status"], "failed")
        self.assertIn("ready-status write failure", status["reason"])

    def test_finalize_cli_attempts_every_output_invalidation(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(analysis, [])
        analysis_path = self.root / "analysis-input.json"
        output_path = self.root / "summary.md"
        result_path = self.root / "finalization-result.json"
        slack_path = self.root / "slack-payload.json"
        status_path = self.root / "finalization-status.json"
        analysis_path.write_text(json.dumps(analysis), encoding="utf-8")
        output_path.write_text("stale summary", encoding="utf-8")
        result_path.write_text("stale result", encoding="utf-8")
        slack_path.write_text("stale payload", encoding="utf-8")
        original_atomic_write = failure_ledger._atomic_write

        def fail_result_invalidation(path: Path, data: bytes) -> None:
            if path == result_path and data == b"":
                raise OSError("simulated result invalidation failure")
            original_atomic_write(path, data)

        with mock.patch.dict(
            os.environ,
            {"MODEL_DECISIONS_JSON": json.dumps(decisions)},
        ):
            with mock.patch(
                "failure_ledger._atomic_write",
                side_effect=fail_result_invalidation,
            ):
                with contextlib.redirect_stderr(io.StringIO()):
                    with self.assertRaises(SystemExit):
                        _finalize_main(
                            [
                                "--analysis-input",
                                str(analysis_path),
                                "--decisions-env",
                                "MODEL_DECISIONS_JSON",
                                "--output",
                                str(output_path),
                                "--result-json",
                                str(result_path),
                                "--slack-payload",
                                str(slack_path),
                            ]
                        )

        self.assertEqual(output_path.read_bytes(), b"")
        self.assertEqual(
            result_path.read_text(encoding="utf-8"),
            "stale result",
        )
        self.assertEqual(slack_path.read_bytes(), b"")
        status = json.loads(status_path.read_text(encoding="utf-8"))
        self.assertEqual(status["status"], "failed")
        self.assertIn("failed to invalidate", status["reason"])
        self.assertIn("failed to invalidate", status["invalidation_error"])

    def test_finalize_cli_rejects_colliding_output_paths(self) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned an API error.",
                    "evidence_refs": ["evidence-0"],
                }
            ],
        )
        analysis_path = self.root / "analysis-input.json"
        decisions_path = self.root / "decisions.json"
        output_path = self.root / "summary.md"
        analysis_path.write_text(json.dumps(analysis), encoding="utf-8")
        decisions_path.write_text(json.dumps(decisions), encoding="utf-8")

        with contextlib.redirect_stderr(io.StringIO()):
            with self.assertRaises(SystemExit):
                _finalize_main(
                    [
                        "--analysis-input",
                        str(analysis_path),
                        "--decisions",
                        str(decisions_path),
                        "--output",
                        str(output_path),
                        "--result-json",
                        str(output_path),
                    ]
                )

        self.assertFalse(output_path.exists())

    def test_finalize_rejects_duplicate_unknown_missing_and_bad_counts(self) -> None:
        analysis = self.synthetic_analysis(["TestOne", "TestTwo"])
        cases = [
            (
                [
                    {
                        "unit_ids": ["test-0"],
                        "category": "api_error",
                        "cause": "API error.",
                        "evidence_refs": ["evidence-0"],
                    },
                    {
                        "unit_ids": ["test-0", "test-1"],
                        "category": "timeout",
                        "cause": "Timeout.",
                        "evidence_refs": ["evidence-1"],
                    },
                ],
                "more than once",
            ),
            (
                [
                    {
                        "unit_ids": ["unknown"],
                        "category": "api_error",
                        "cause": "API error.",
                        "evidence_refs": ["evidence-0"],
                    }
                ],
                "unknown test units",
            ),
            (
                [
                    {
                        "unit_ids": ["test-0"],
                        "category": "api_error",
                        "cause": "API error.",
                        "evidence_refs": ["evidence-0"],
                    }
                ],
                "did not classify every",
            ),
        ]
        for groups, message in cases:
            with self.subTest(message=message):
                with self.assertRaisesRegex(FinalizationError, message):
                    finalize_analysis(
                        analysis,
                        self.model_decisions(analysis, groups),
                    )

        bad_analysis = json.loads(json.dumps(analysis))
        bad_analysis["unit_counts"]["tests"] = 3
        bad_analysis["analysis_digest"] = _canonical_digest(
            bad_analysis,
            omit_key="analysis_digest",
        )
        with self.assertRaisesRegex(FinalizationError, "test total"):
            finalize_analysis(
                bad_analysis,
                self.model_decisions(
                    bad_analysis,
                    [
                        {
                            "remaining": True,
                            "category": "api_error",
                            "cause": "API error.",
                            "evidence_refs": ["evidence-0"],
                        }
                    ],
                ),
            )

    def test_finalizer_enforces_constraints_omitted_from_claude_schema(
        self,
    ) -> None:
        analysis = self.synthetic_analysis(["TestOne"])
        valid_group = {
            "unit_ids": ["test-0"],
            "category": "api_error",
            "cause": "Atlas returned an API error.",
            "evidence_refs": ["evidence-0"],
        }
        cases = [
            (
                self.model_decisions(analysis, [{}] * 1001),
                "groups exceed the limit",
            ),
            (
                self.model_decisions(
                    analysis,
                    [{**valid_group, "remaining": True}],
                ),
                "either unit_ids or remaining",
            ),
            (
                self.model_decisions(
                    analysis,
                    [{**valid_group, "unit_ids": ["test-0", "test-0"]}],
                ),
                "repeats a unit ID",
            ),
            (
                self.model_decisions(
                    analysis,
                    [
                        {
                            **valid_group,
                            "evidence_refs": ["evidence-0", "evidence-0"],
                        }
                    ],
                ),
                "invalid evidence refs",
            ),
            (
                self.model_decisions(
                    analysis,
                    [{**valid_group, "cause": "x" * 241}],
                ),
                "cause exceeds 240 characters",
            ),
            (
                self.model_decisions(
                    analysis,
                    [valid_group],
                    why="x" * 501,
                ),
                "why exceeds 500 characters",
            ),
        ]

        for decisions, message in cases:
            with self.subTest(message=message):
                with self.assertRaisesRegex(FinalizationError, message):
                    finalize_analysis(analysis, decisions)

    def test_unresolved_forces_review_but_code_regression_stays_red(self) -> None:
        analysis = self.synthetic_analysis(["TestOne", "TestTwo"])
        review = self.model_decisions(
            analysis,
            [
                {
                    "unit_ids": ["test-0"],
                    "category": "api_error",
                    "cause": "Atlas returned an API error.",
                    "evidence_refs": ["evidence-0"],
                },
                {
                    "remaining": True,
                    "category": "unresolved",
                    "cause": "The primary failure could not be isolated.",
                    "evidence_refs": ["evidence-1"],
                    "ambiguity": "Owned evidence has no primary failure anchor.",
                },
            ],
        )
        review_result = finalize_analysis(analysis, review)
        self.assertEqual(
            review_result["headline"],
            "Automatic triage incomplete",
        )
        self.assertEqual(review_result["confidence"], "medium")
        self.assertEqual(
            review_result["urgency"],
            "non-urgent review recommended",
        )
        self.assertEqual(review_result["unresolved_tests"], 1)
        self.assertIn(
            "2 unique top-level tests failed; "
            "1 test could not be classified automatically.",
            review_result["slack_mrkdwn"],
        )
        self.assertIn(
            "• Unclassified: 1 test",
            review_result["slack_mrkdwn"],
        )

        red = json.loads(json.dumps(review))
        red["groups"][0]["category"] = "code_regression"
        red["groups"][0]["cause"] = "Provider produced inconsistent state."
        red_result = finalize_analysis(analysis, red)
        self.assertEqual(red_result["headline"], "CODE REGRESSION DETECTED")
        self.assertEqual(red_result["verdict"], "red")

    def test_pak_authorization_candidate_requires_ambiguity(self) -> None:
        analysis = self.synthetic_analysis(["TestAuth", "TestTimeout"])
        analysis["test_units"][0]["machine_facts"][
            "pak_authorization_failure_candidate"
        ] = True
        analysis["test_units"][0]["machine_facts"][
            "authorization_failure_distinct_requests"
        ] = 20
        analysis["analysis_digest"] = _canonical_digest(
            analysis,
            omit_key="analysis_digest",
        )
        groups = [
            {
                "unit_ids": ["test-0"],
                "category": "code_regression",
                "cause": "PAK DELETE requests returned HTTP 401.",
                "evidence_refs": ["evidence-0"],
            },
            {
                "remaining": True,
                "category": "timeout",
                "cause": "The operation exceeded its polling deadline.",
                "evidence_refs": ["evidence-1"],
            },
        ]

        with self.assertRaisesRegex(
            FinalizationError,
            "PAK authorization-regression candidates require",
        ):
            finalize_analysis(
                analysis,
                self.model_decisions(analysis, groups),
            )

        groups[0]["ambiguity"] = (
            "The run cannot isolate AuthN from credentials or environment."
        )
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, groups),
        )

        self.assertEqual(result["verdict"], "red")
        self.assertEqual(result["confidence"], "medium")
        self.assertIn("*Code regression (1 test)*:", result["slack_mrkdwn"])
        self.assertIn("*Other failures*:", result["slack_mrkdwn"])
        self.assertIn("• Timeout: 1 test", result["slack_mrkdwn"])

    def test_finalize_green_and_suite_unverified(self) -> None:
        green = self.synthetic_analysis([], green_eligible=True)
        green_result = finalize_analysis(
            green,
            self.model_decisions(green, []),
        )
        self.assertEqual(green_result["verdict"], "green")
        self.assertEqual(green_result["headline"], "All tests passed")

        incomplete = self.synthetic_analysis([])
        incomplete["gates"]["run_incomplete"] = True
        incomplete["analysis_digest"] = _canonical_digest(
            incomplete,
            omit_key="analysis_digest",
        )
        incomplete_result = finalize_analysis(
            incomplete,
            self.model_decisions(incomplete, []),
        )
        self.assertEqual(incomplete_result["headline"], "Run incomplete")

        setup_job = {
            "unit_id": "job-setup",
            "kind": "setup",
            "job_id": "1",
            "job_name": "variables",
            "role": "setup",
            "conclusion": "failure",
            "duration_minutes": 1,
            "deterministic_category": "code_regression",
            "sources": ["setup_jobs"],
        }
        unverified = self.synthetic_analysis(
            [],
            execution_state="confirmed_none",
            workflow_units=[setup_job],
        )
        unverified_result = finalize_analysis(
            unverified,
            self.model_decisions(unverified, []),
        )
        self.assertEqual(unverified_result["verdict"], "red")
        self.assertEqual(unverified_result["headline"], "SUITE UNVERIFIED")

    def test_live_run_with_only_reporting_active_can_finalize_green(self) -> None:
        self.write_jobs(
            {
                "id": 1,
                "name": self._test_job_name("package-a"),
                "conclusion": "success",
            },
            {
                "id": 2,
                "name": "trigger-test-summary",
                "conclusion": None,
            },
        )
        ledger = build_ledger(self.root, self.jobs_path)
        run = {
            "run_id": 1,
            "run_attempt": 1,
            "run_number": 1102,
            "head_sha": "abcdef123456",
            "status": "in_progress",
            "conclusion": None,
        }
        ledger["run_id"] = 1
        ledger["run_attempt"] = 1

        analysis = build_analysis_input(ledger, run, self.root)
        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )

        self.assertTrue(analysis["gates"]["reporting_only_active"])
        self.assertFalse(analysis["gates"]["run_incomplete"])
        self.assertTrue(analysis["gates"]["green_eligible"])
        self.assertEqual(result["verdict"], "green")
        self.assertEqual(result["headline"], "All tests passed")

    def test_finalize_renders_unverified_test_job(self) -> None:
        test_job = {
            "unit_id": "job-timeout",
            "kind": "test_job",
            "job_id": "1",
            "job_name": self._test_job_name("package-a"),
            "role": "test",
            "conclusion": "timed_out",
            "duration_minutes": 360,
            "deterministic_category": None,
            "sources": ["timed_out_jobs"],
            "reasons": ["test job timed out"],
        }
        analysis = self.synthetic_analysis(
            [],
            ledger_complete=False,
            workflow_units=[test_job],
        )

        result = finalize_analysis(
            analysis,
            self.model_decisions(analysis, []),
        )

        self.assertEqual(result["headline"], "Automatic triage incomplete")
        self.assertIn("Unverified test job", result["slack_mrkdwn"])
        self.assertIn("timed_out, 360m", result["slack_mrkdwn"])
        self.assertNotIn(
            "0 unique top-level tests failed",
            result["slack_mrkdwn"],
        )

    def test_finalize_renders_teardown_packages_and_enforces_length(self) -> None:
        analysis = self.synthetic_analysis(
            [f"TestWithAnExtremelyLongName{index:03d}" for index in range(200)],
            teardown_packages=[
                "github.com/example/one",
                "github.com/example/two",
            ],
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "api_error",
                    "cause": "Atlas returned a repeated service-side error. " * 4,
                    "evidence_refs": [
                        f"evidence-{index}" for index in range(200)
                    ],
                }
            ],
            action="Review the affected Atlas service before rerunning.",
        )

        result = finalize_analysis(analysis, decisions)

        self.assertLessEqual(len(result["slack_mrkdwn"]), SLACK_HARD_LIMIT)
        self.assertEqual(result["package_teardowns"], 2)
        self.assertIn("Package teardowns: 2", result["slack_mrkdwn"])
        self.assertIn("github.com/example/one", result["slack_mrkdwn"])

    def test_slack_text_truncates_at_word_boundaries(self) -> None:
        self.assertEqual(
            _slack_text("alpha bravo charlie", 14),
            "alpha bravo…",
        )
        self.assertEqual(_slack_text("alpha bravo", 11), "alpha bravo")
        self.assertEqual(_slack_text("abcdefghijkl", 5), "abcd…")
        self.assertEqual(_slack_text("alpha", 1), "…")
        self.assertEqual(_slack_text("alpha", 0), "")

    def test_finalize_uses_utf8_byte_cap_and_aggregates_shape_b(self) -> None:
        analysis = self.synthetic_analysis(
            [f"TestCreateTimeout{index}" for index in range(100)],
        )
        for unit in analysis["test_units"]:
            unit["machine_facts"]["shape_b_candidate"] = True
        analysis["analysis_digest"] = _canonical_digest(
            analysis,
            omit_key="analysis_digest",
        )
        decisions = self.model_decisions(
            analysis,
            [
                {
                    "remaining": True,
                    "category": "timeout",
                    "cause": "🚨" * 200,
                    "evidence_refs": [
                        f"evidence-{index}" for index in range(100)
                    ],
                    "ambiguity": "Deletion completion is ambiguous.",
                    "note": "delete_on_timeout_unverified",
                }
            ],
            action="🔎" * 300,
        )

        result = finalize_analysis(analysis, decisions)

        self.assertLessEqual(
            len(f"{result['slack_mrkdwn']}\n".encode()),
            SLACK_HARD_LIMIT,
        )
        self.assertIn(
            "Delete-on-timeout unverified: 100 tests",
            result["slack_mrkdwn"],
        )
        self.assertIn(
            "included in Timeout total",
            result["slack_mrkdwn"],
        )


if __name__ == "__main__":
    unittest.main()
