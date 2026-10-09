"""Offline checks for evidence completeness; no model or GitHub calls."""

import contextlib
import io
import json
import os
import runpy
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from prepare import analyze, excerpts, prepare

RUN = dict(id=123, run_number=7, run_attempt=2, head_sha="abcdef1234", html_url="https://example/run/123")


def job(job_id, conclusion: str | None = "failure", name="tests / network", status="completed"):
    return dict(id=job_id, name=name, status=status, conclusion=conclusion,
                steps=[dict(name="Acceptance Tests", status=status, conclusion=conclusion)])


class PrepareTests(unittest.TestCase):
    def setUp(self):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        self.output = Path(directory.name) / "evidence"

    def collect(self, jobs, logs):
        def api(endpoint, *args):
            if endpoint == "runs/123":
                return json.dumps(RUN)
            if endpoint.startswith("runs/123/jobs?"):
                self.assertIn("filter=latest", endpoint)
                self.assertEqual(args, ("--paginate", "--slurp"))
                return json.dumps([dict(jobs=jobs[:2]), dict(jobs=jobs[2:])])
            result = logs[int(endpoint.split("/")[1])]
            if isinstance(result, Exception):
                raise result
            return result

        with patch("prepare.github", side_effect=api) as calls:
            report = json.loads(prepare(123, self.output).read_text())
        downloads = [call.args[0] for call in calls.call_args_list if call.args[0].endswith("/logs")]
        self.assertEqual(len(downloads), len(set(downloads)))
        for call in calls.call_args_list:
            if call.args[0].endswith("/logs"):
                # gh rejects logs with terminal escape sequences unless we opt in.
                self.assertIn("--allow-escape-sequences", call.args)
        return report

    def test_index_keeps_verdicts_and_positions_without_matching_passing_errors(self):
        lines = [
            "=== RUN   TestPass", "[ERROR] expected rejection", "--- PASS: TestPass (0.01s)",
            "2026-09-01T01:02:03.123Z \x1b[31m    --- FAIL: TestBad/child (2.00s)\x1b[0m",
            "--- FAIL: TestBad (2.00s)", "FAIL\tgithub.com/example/pkg\t2.01s",
            "panic: test timed out after 5h0m0s", "##[error]Process completed with exit code 1.",
        ]
        index = excerpts("\n".join(lines)).split("\n\nCONTEXT")[0]
        self.assertNotIn("TestPass", index)
        self.assertNotIn("[ERROR]", index)
        self.assertIn("4:     --- FAIL: TestBad/child (2.00s)", index)
        self.assertIn("5: --- FAIL: TestBad (2.00s)", index)
        self.assertIn("6: FAIL\tgithub.com/example/pkg", index)
        self.assertIn("7: panic:", index)
        self.assertNotIn("\x1b", index)

    def test_collects_all_pages_and_nonstandard_failures_but_skips_healthy_and_self(self):
        jobs = [job(1), job(2, "success"), job(3, "timed_out"), job(4, "cancelled"),
                job(5, "skipped"), job(6, None, status="in_progress"),
                job(7, name="clean-after / cleanup-test-env-general"),
                job(8, name="caller / trigger-test-summary")]
        logs = {i: f"full log {i}\nFAIL\tpackage{i}\n" for i in (1, 3, 4, 7)}
        report = self.collect(jobs, logs)
        self.assertEqual(len(report["jobs"]), 8)
        self.assertEqual(report["run_attempt"], 2)
        for entry in report["jobs"]:
            if entry["id"] in logs:
                self.assertEqual(Path(entry["log_file"]).read_text(), logs[entry["id"]])
                self.assertIn("SIGNAL INDEX", Path(entry["evidence_file"]).read_text())
                analysis = json.loads(Path(entry["analysis_file"]).read_text())
                self.assertEqual(analysis["package_failures"], [f"package{entry['id']}"])
                self.assertNotIn("analysis", entry)  # Keep diagnostics out of the job inventory.
            else:
                self.assertNotIn("log_file", entry)
                self.assertNotIn("analysis_file", entry)
        self.assertEqual(report["jobs"][5]["status"], "in_progress")
        self.assertEqual(report["jobs"][0]["steps"], jobs[0]["steps"])

    def test_successful_job_preserves_whether_tests_ran_or_were_skipped(self):
        ran, skipped = job(1, "success"), job(2, "success")
        skipped["steps"][0]["conclusion"] = "skipped"
        report = self.collect([ran, skipped], {})
        self.assertEqual(report["jobs"][0]["steps"][0]["conclusion"], "success")
        self.assertEqual(report["jobs"][1]["steps"][0]["conclusion"], "skipped")

    def test_unavailable_logs_are_explicit_and_do_not_hide_other_jobs(self):
        report = self.collect([job(i) for i in range(1, 5)], {
            1: subprocess.CalledProcessError(1, "gh", stderr="gh: Not Found (HTTP 404)"),
            2: " \n", 3: subprocess.TimeoutExpired("gh", 120), 4: "no matching markers",
        })
        self.assertIn("HTTP 404", report["jobs"][0]["log_error"])
        self.assertEqual(report["jobs"][1]["log_error"], "empty log response")
        self.assertIn("timed out", report["jobs"][2]["log_error"])
        evidence = Path(report["jobs"][3]["evidence_file"]).read_text()
        self.assertIn("(no matching signals)", evidence)
        self.assertIn("1: no matching markers", evidence)

    def test_access_denial_stops_collection(self):
        with self.assertRaises(subprocess.CalledProcessError):
            self.collect([job(1)], {1: subprocess.CalledProcessError(1, "gh", stderr="HTTP 403")})
        self.assertFalse((self.output / "run.json").exists())

    def test_attempt_change_and_existing_output_cannot_publish_stale_evidence(self):
        with patch("prepare.github", side_effect=[json.dumps(RUN), '[{"jobs": []}]',
                                                 json.dumps({**RUN, "run_attempt": 3})]):
            with self.assertRaisesRegex(RuntimeError, "rerun"):
                prepare(123, self.output)
        self.assertFalse((self.output / "run.json").exists())
        with patch("prepare.github") as api:
            with self.assertRaises(FileExistsError):
                prepare(123, self.output)
            api.assert_not_called()

    def test_excerpts_bound_context_keep_all_signals_and_mark_gaps(self):
        lines = ["unrelated output"] * 200
        lines[40] = "--- FAIL: TestFirst (1.00s)"
        lines[120] = "--- FAIL: TestLast (1.00s)"
        result = excerpts("\n".join(lines))
        self.assertIn("41: --- FAIL: TestFirst", result)
        self.assertIn("121: --- FAIL: TestLast", result)
        self.assertIn("... omitted lines", result)
        self.assertIn("200: unrelated output", result)
        self.assertLess(len(result.splitlines()), 140)

    def test_distant_diagnostics_survive_among_passing_tests_without_implying_failure(self):
        lines = ["=== RUN   TestExpectedError", "Provider produced inconsistent result: expected error",
                 "--- PASS: TestExpectedError (0.01s)"]
        lines += [f"--- PASS: TestOther{i} (0.01s)" for i in range(100)]
        lines += ["=== RUN   TestRegression", "Error: Provider produced inconsistent result after apply"]
        lines += ["long diagnostic detail"] * 100
        lines += ["--- FAIL: TestRegression (1.00s)", "FAIL\tgithub.com/example/pkg"]
        result = excerpts("\n".join(lines))
        index = result.split("\n\nCONTEXT")[0]
        self.assertIn("105: Error: Provider produced inconsistent result after apply", index)
        self.assertIn("206: --- FAIL: TestRegression", index)
        self.assertIn("3: --- PASS: TestExpectedError", result)
        self.assertNotIn("--- FAIL: TestExpectedError", result)
        self.assertNotIn("TestOther50", result)

    def test_analysis_attributes_diagnostics_to_the_failing_test_only(self):
        lines = [
            "=== RUN   TestExpectedError",
            "Error: expected rejection from a passing negative test",
            "--- PASS: TestExpectedError (0.01s)",
            "=== RUN   TestRegression",
            "Error: Provider produced inconsistent result after apply",
            "    --- FAIL: TestRegression/child (1.00s)",
            "--- FAIL: TestRegression (1.00s)",
            "FAIL\tgithub.com/example/pkg\t2.01s",
        ]
        result = analyze("\n".join(lines))
        self.assertEqual([test["name"] for test in result["failed_tests"]], ["TestRegression"])
        self.assertEqual(result["failed_tests"][0]["subtests"], ["TestRegression/child"])
        diagnostics = "\n".join(result["failed_tests"][0]["diagnostics"])
        self.assertIn("Provider produced inconsistent result after apply", diagnostics)
        self.assertNotIn("expected rejection", diagnostics)
        self.assertEqual(result["package_failures"], ["github.com/example/pkg"])

    def test_analysis_separates_deadline_and_runtime_panics_and_build_failures(self):
        lines = [
            "=== RUN   TestA", "--- PASS: TestA (0.01s)",
            "panic: test timed out after 5h0m0s",
            "=== RUN   TestB", "panic: nil pointer dereference", "--- FAIL: TestB (0.02s)",
            "FAIL\tgithub.com/example/pkg\t1.0s",
        ]
        result = analyze("\n".join(lines))
        self.assertEqual(result["deadline_panics"], ["3: panic: test timed out after 5h0m0s"])
        self.assertEqual(result["panics"], ["5: panic: nil pointer dereference"])
        self.assertEqual([test["name"] for test in result["failed_tests"]], ["TestB"])

    def test_analysis_follows_parallel_output_markers(self):
        lines = [
            "=== RUN   TestRegression", "=== PAUSE TestRegression",
            "=== RUN   TestExpectedError", "=== PAUSE TestExpectedError",
            "=== CONT  TestRegression", "Error: Provider produced inconsistent result after apply",
            "=== CONT  TestExpectedError", "Error: expected rejection",
            "=== NAME  TestRegression", "Error: unexpected new value",
            "--- FAIL: TestRegression (1.00s)", "--- PASS: TestExpectedError (1.00s)",
        ]
        diagnostics = "\n".join(analyze("\n".join(lines))["failed_tests"][0]["diagnostics"])
        self.assertIn("Provider produced inconsistent result", diagnostics)
        self.assertIn("unexpected new value", diagnostics)
        self.assertNotIn("expected rejection", diagnostics)

    def test_analysis_does_not_attach_diagnostics_outside_active_test(self):
        for boundary in ("=== PAUSE TestFailure", "=== NAME  ", "--- FAIL: TestFailure (1.00s)",
                         "FAIL\tgithub.com/example/first\t1.00s", "ok\tgithub.com/example/first\t1.00s"):
            with self.subTest(boundary=boundary):
                result = analyze("\n".join([
                    "=== RUN   TestFailure", boundary, "Error: unrelated package teardown",
                    "--- FAIL: TestFailure (1.00s)",
                ]))
                self.assertNotIn("unrelated package teardown", "\n".join(result["failed_tests"][0]["diagnostics"]))

    def test_analysis_excludes_passing_subtest_diagnostics_from_failed_parent(self):
        result = analyze("\n".join([
            "=== RUN   TestParent", "=== RUN   TestParent/expected",
            "Error: Provider produced inconsistent result: expected rejection",
            "--- PASS: TestParent/expected (0.01s)", "=== RUN   TestParent/timeout",
            "Error: context deadline exceeded", "--- FAIL: TestParent/timeout (1.00s)",
            "--- FAIL: TestParent (1.01s)",
        ]))
        failure = result["failed_tests"][0]
        self.assertEqual(failure["name"], "TestParent")
        self.assertEqual(failure["subtests"], ["TestParent/timeout"])
        diagnostics = "\n".join(failure["diagnostics"])
        self.assertIn("context deadline exceeded", diagnostics)
        self.assertNotIn("Provider produced inconsistent result", diagnostics)

    def test_analysis_caps_diagnostics_and_reports_missing_markers(self):
        lines = ["=== RUN   TestMany"] + [f"Error: failure {i}" for i in range(60)]
        lines += ["Error: Provider produced inconsistent result after apply"]
        lines += ["--- FAIL: TestMany (0.02s)"]
        result = analyze("\n".join(lines))
        diagnostics = result["failed_tests"][0]["diagnostics"]
        self.assertEqual(len(diagnostics), 41)
        self.assertIn("failure 0", diagnostics[0])
        self.assertIn("22 diagnostic lines omitted", diagnostics[20])
        self.assertIn("Provider produced inconsistent result", diagnostics[-2])
        self.assertIn("--- FAIL: TestMany", diagnostics[-1])
        self.assertEqual(analyze("no test markers here"),
                         {"failed_tests": [], "package_failures": [], "panics": [], "deadline_panics": []})

    def test_monthly_summary_accepts_test_and_non_test_category_counts(self):
        monthly = Path(__file__).resolve().parents[2] / "monthly-test-suite-summary/scripts/monthly_summary.py"
        parse = runpy.run_path(str(monthly))["parse_categories"]
        self.assertEqual(parse("• API errors: 3 tests (approx.)\n• Cleanup: 8 failures (2 tests, 6 items)\n"
                               "• Timeout: 1 package"), {"api errors": 3, "cleanup": 8, "timeout": 1})
        self.assertEqual(parse("*Other failures*: API errors 3, Cleanup 8, Timeout 1"),
                         {"api errors": 3, "cleanup": 8, "timeout": 1})

    def test_monthly_summary_excludes_incomplete_runs_from_regression_free_percentage(self):
        monthly = Path(__file__).resolve().parents[2] / "monthly-test-suite-summary/scripts/monthly_summary.py"
        namespace = runpy.run_path(str(monthly))
        summaries = [
            ":red_circle: *Test Suite #1 — CODE REGRESSION DETECTED*",
            ":yellow_circle: *Test Suite #2 — Infrastructure noise only*",
            "\n:yellow_circle: *Test Suite #3 — Results incomplete*",
            ":green_circle: *Test Suite #4 — All tests passed*",
        ]
        self.assertEqual([namespace["parse_verdict"](text) for text in summaries],
                         ["red", "yellow", "incomplete", "green"])
        # Known regressions remain red even if other jobs have missing evidence.
        self.assertEqual(namespace["parse_verdict"](summaries[0] + "\nResults incomplete"), "red")
        runs = [dict(id=i, run_number=i, created_at=f"2026-09-0{i}T00:00:00Z", html_url=f"https://example/{i}")
                for i in range(1, 5)]
        self.output.mkdir()
        previous = Path.cwd()
        try:
            os.chdir(self.output)
            with patch.dict(namespace["main"].__globals__,
                            list_scheduled_runs=lambda *_: runs, fetch_summary=lambda i: summaries[i - 1]), \
                    patch("sys.argv", ["monthly_summary.py", "--date", "2026-10-01", "--json-out", "report.json"]), \
                    contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
                namespace["main"]()
            report = json.loads(Path("report.json").read_text())
            self.assertEqual(report["totals"]["without_regression"], 2)
            self.assertEqual(report["totals"]["pct_without_regression"], "66.67%")
            self.assertEqual(report["totals"]["incomplete"], 1)
            self.assertEqual(report["incomplete_runs"][0]["run_number"], 3)
            self.assertIn("incomplete test evidence", Path("monthly-summary-2026-10-01.md").read_text())
        finally:
            os.chdir(previous)


if __name__ == "__main__":
    unittest.main()
