"""Offline checks for evidence completeness; no model or GitHub calls."""

import json
import runpy
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from prepare import excerpts, prepare

RUN = dict(id=123, run_number=7, run_attempt=2, head_sha="abcdef1234", html_url="https://example/run/123")


def job(job_id, conclusion="failure", name="tests / network", status="completed"):
    return dict(id=job_id, name=name, status=status, conclusion=conclusion,
                steps=[dict(name="Acceptance Tests", conclusion=conclusion)])


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
            else:
                self.assertNotIn("log_file", entry)
        self.assertEqual(report["jobs"][5]["status"], "in_progress")
        self.assertEqual(report["jobs"][0]["failed_steps"], ["Acceptance Tests"])

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

    def test_monthly_summary_accepts_test_and_non_test_category_counts(self):
        monthly = Path(__file__).resolve().parents[2] / "monthly-test-suite-summary/scripts/monthly_summary.py"
        parse = runpy.run_path(str(monthly))["parse_categories"]
        self.assertEqual(parse("• API errors: 3 tests (approx.)\n• Cleanup: 8 failures (2 tests, 6 items)\n"
                               "• Timeout: 1 package"), {"api errors": 3, "cleanup": 8, "timeout": 1})
        self.assertEqual(parse("*Other failures*: API errors 3, Cleanup 8, Timeout 1"),
                         {"api errors": 3, "cleanup": 8, "timeout": 1})


if __name__ == "__main__":
    unittest.main()
