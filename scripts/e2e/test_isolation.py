"""Exercise cleanup and explicit contexts without needing a real cluster."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
MOCK = r'''#!/usr/bin/env python3
import json, os, pathlib, sys
args = sys.argv[1:]
name = pathlib.Path(sys.argv[0]).name
with open(os.environ["CALLS"], "a") as f:
    f.write(json.dumps([name, args, os.environ.get("KUBECONFIG")]) + "\n")
if name == "kind" and args[:2] == ["create", "cluster"]:
    pathlib.Path(args[args.index("--kubeconfig")+1]).write_text("private kubeconfig")
    if os.environ.get("FAIL_CREATE"):
        sys.exit(1)
if name == "kubectl":
    if os.environ.get("FAIL_APPLY") and "apply" in args:
        sys.exit(1)
    if "--raw" in args:
        print("ok")
'''


class IsolationTest(unittest.TestCase):
    def run_case(self, fail=None):
        with tempfile.TemporaryDirectory() as tmp:
            tmp = Path(tmp)
            tools = tmp / "bin"
            tools.mkdir()
            for name in ["kind", "kubectl"]:
                executable = tools / name
                executable.write_text(MOCK)
                executable.chmod(0o755)
            original = tmp / "original-config"
            original.write_text("current-context: do-not-touch\n")
            calls = tmp / "calls"
            env = dict(os.environ, PATH=str(tools) + os.pathsep + os.environ["PATH"],
                       KUBECONFIG=str(original), CALLS=str(calls),
                       E2E_ARTIFACTS=str(tmp / "artifacts"))
            if fail:
                env[fail] = "1"
            result = subprocess.run(["bash", "scripts/e2e/run.sh", "argocd"],
                                    cwd=ROOT, env=env, capture_output=True, text=True)
            self.assertEqual(result.returncode, 1 if fail else 0, result.stderr)
            self.assertEqual(original.read_text(), "current-context: do-not-touch\n")
            entries = [json.loads(line) for line in calls.read_text().splitlines()]
            create = next(args for tool, args, _ in entries if tool == "kind" and args[:2] == ["create", "cluster"])
            cluster = create[create.index("--name")+1]
            config = create[create.index("--kubeconfig")+1]
            self.assertNotEqual(config, str(original))
            self.assertFalse(Path(config).exists(), "temporary config must be removed")
            self.assertTrue(any(tool == "kind" and args[:2] == ["delete", "cluster"] for tool, args, _ in entries))
            for tool, args, kubeconfig in entries:
                self.assertEqual(kubeconfig, config)
                if tool == "kubectl":
                    self.assertEqual(args[:4], ["--kubeconfig", config, "--context", "kind-" + cluster])
                if tool == "kind" and args[:2] == ["delete", "cluster"]:
                    self.assertEqual(args[args.index("--kubeconfig")+1], config)
                    self.assertEqual(args[args.index("--name")+1], cluster)

    def test_success(self):
        self.run_case()

    def test_apply_failure(self):
        self.run_case("FAIL_APPLY")

    def test_creation_failure(self):
        self.run_case("FAIL_CREATE")
