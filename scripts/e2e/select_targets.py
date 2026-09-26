#!/usr/bin/env python3
"""Select E2E targets from a NUL-separated git diff on stdin."""
import fnmatch
import json
import sys
from pathlib import Path

TARGETS = json.loads(Path(__file__).with_name("targets.json").read_text())


def select(paths):
    selected = set()
    for path in paths:
        if path == ".github/workflows/e2e.yml" or (
            path.startswith("scripts/e2e/")
            and not path.startswith("scripts/e2e/cases/")
        ):
            return sorted(TARGETS)
        for target, patterns in TARGETS.items():
            if path == f"scripts/e2e/cases/{target}.sh" or any(
                fnmatch.fnmatchcase(path, pattern) for pattern in patterns
            ):
                selected.add(target)
    return sorted(selected)


if __name__ == "__main__":
    paths = sys.stdin.buffer.read().decode().split("\0")
    print(json.dumps(select(paths)))
