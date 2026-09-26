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

    def test_operator_runner_and_workflow_changes_only_run_operator(self):
        for path in ["scripts/e2e/run.sh", ".github/workflows/e2e.yml"]:
            self.assertEqual(module.select([path]), ["operator"])

    def test_selector_changes_need_no_cluster(self):
        self.assertEqual(module.select(["scripts/e2e/select_targets.py"]), [])

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

    def test_registry_addition_runs_only_added_target(self):
        base = {"argocd": ["contents/argocd/**"]}
        head = {**base, "mysql": ["contents/mysql/**"]}
        self.assertEqual(module.select([module.REGISTRY, "scripts/e2e/cases/mysql.sh"],
                                       targets=head, base_targets=base), ["mysql"])

    def test_registry_path_change_runs_only_changed_target(self):
        base = {"a": ["old/**"], "b": ["b/**"]}
        head = {"a": ["new/**"], "b": ["b/**"]}
        self.assertEqual(module.select([module.REGISTRY], targets=head, base_targets=base), ["a"])

    def test_registry_removal_does_not_schedule_deleted_case(self):
        base = {"a": ["a/**"], "b": ["b/**"]}
        head = {"b": ["b/**"]}
        self.assertEqual(module.select([module.REGISTRY, "scripts/e2e/cases/a.sh", "a/deleted.yaml"],
                                       targets=head, base_targets=base), [])

    def test_registry_reordering_needs_no_cluster(self):
        base = {"a": ["a/**", "shared/**"], "b": ["b/**"]}
        head = {"b": ["b/**"], "a": ["shared/**", "a/**"]}
        self.assertEqual(module.select([module.REGISTRY], targets=head, base_targets=base), [])

    def test_registry_requires_base_snapshot(self):
        with self.assertRaisesRegex(ValueError, "base-registry"):
            module.select([module.REGISTRY])

    def test_merge_ref_does_not_schedule_unrelated_base_addition(self):
        base = {"a": ["a/**"]}
        head = {**base, "b": ["b/**"]}
        merged = {**head, "unrelated": ["other/**"]}
        self.assertEqual(module.select([module.REGISTRY], targets=merged,
                                       base_targets=base, head_targets=head), ["b"])
        self.assertEqual(module.select(["scripts/e2e/run.sh"], targets=merged), [])

    def test_unit_test_changes_need_no_cluster(self):
        self.assertEqual(module.select(["scripts/e2e/test_select.py", "scripts/e2e/test_isolation.py"]), [])

    def test_version_helper_only_runs_packaging(self):
        self.assertEqual(module.select(["scripts/e2e/tool_version.py"]), ["helm", "kustomize"])
