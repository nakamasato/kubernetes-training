import importlib.util
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location("targets", Path(__file__).with_name("select_targets.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class SelectionTest(unittest.TestCase):
    def test_unrelated_and_empty(self):
        self.assertEqual(module.select(["README.md", "go.mod", ""]), [])

    def test_single_target(self):
        self.assertEqual(module.select(["contents/argocd/setup/kustomization.yaml"]), ["argocd"])

    def test_case_only(self):
        self.assertEqual(module.select(["scripts/e2e/cases/argocd.sh"]), ["argocd"])

    def test_shared(self):
        for path in ["scripts/e2e/run.sh", "scripts/e2e/targets.json", ".github/workflows/e2e.yml"]:
            self.assertEqual(module.select([path]), sorted(module.TARGETS))

    def test_rename_or_delete_paths(self):
        self.assertEqual(module.select(["contents/argocd/deleted.yaml", "contents/prometheus-operator/new.yaml"]), ["argocd", "prometheus-operator"])

    def test_registered_cases_exist(self):
        for target in module.TARGETS:
            self.assertTrue(Path(__file__).with_name("cases").joinpath(target + ".sh").is_file())

    def test_monitoring_targets_are_independent(self):
        for target in ["prometheus", "prometheus-operator", "grafana"]:
            self.assertEqual(module.select([f"contents/{target}/kustomization.yaml"]), [target])

    def test_packaging_targets(self):
        for path in ["contents/helm/README.md", "contents/helm/hello-world/helloworld-chart/values.yaml"]:
            self.assertEqual(module.select([path]), ["helm"])
        for path in ["contents/kustomize/README.md", "contents/kustomize/example/overlays/prod/kustomization.yaml"]:
            self.assertEqual(module.select([path]), ["kustomize"])
        self.assertEqual(module.select(["contents/helm-vs-kustomize/README.md"]), [])
