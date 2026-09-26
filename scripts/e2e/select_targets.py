#!/usr/bin/env python3
"""Select E2E targets from a NUL-separated git diff on stdin."""
import argparse
import fnmatch
import json
import sys
from pathlib import Path

TARGETS = json.loads(Path(__file__).with_name("targets.json").read_text())
REGISTRY = "scripts/e2e/targets.json"


def select(paths, targets=None, base_targets=None, head_targets=None):
    targets = TARGETS if targets is None else targets
    head_targets = targets if head_targets is None else head_targets
    selected = set()
    for path in paths:
        if path == REGISTRY:
            if base_targets is None:
                raise ValueError("registry changes require --base-registry")
            # Compare the PR's own snapshots, not additions from the merge ref.
            # Removed targets cannot run; ordering/formatting changes need no cluster.
            selected.update(
                target for target, patterns in head_targets.items()
                if target in targets and (
                    target not in base_targets
                    or set(patterns) != set(base_targets[target])
                )
            )
            continue
        if path == "scripts/e2e/tool_version.py":
            selected.update({"helm", "kustomize"} & targets.keys())
            continue
        if path.startswith("scripts/e2e/test_") and path.endswith(".py"):
            # The changes job always executes these tests without a cluster.
            continue
        for target, patterns in targets.items():
            if path == f"scripts/e2e/cases/{target}.sh" or any(
                fnmatch.fnmatchcase(path, pattern) for pattern in patterns
            ):
                selected.add(target)
    return sorted(selected)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-registry", type=Path)
    parser.add_argument("--head-registry", type=Path)
    args = parser.parse_args()
    base = json.loads(args.base_registry.read_text()) if args.base_registry else None
    head = json.loads(args.head_registry.read_text()) if args.head_registry else None
    paths = sys.stdin.buffer.read().decode().split("\0")
    print(json.dumps(select(paths, base_targets=base, head_targets=head)))
