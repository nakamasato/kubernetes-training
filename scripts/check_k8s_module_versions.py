#!/usr/bin/env python3
"""Ensure documented client-go/controller-runtime versions match go.mod."""

from __future__ import annotations

import re
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
MODULES = {
    "client-go": "k8s.io/client-go",
    "controller-runtime": "sigs.k8s.io/controller-runtime",
}
# Covers pkg.go.dev module@version URLs, GitHub versioned source/release URLs,
# and prose such as "client-go v0.37.1".
REFERENCE_PATTERNS = {
    "client-go": re.compile(
        r"(?:k8s\.io/client-go@|github\.com/kubernetes/client-go/(?:blob|tree|releases/tag)/|\bclient-go\s+)"
        r"(v\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?)"
    ),
    "controller-runtime": re.compile(
        r"(?:sigs\.k8s\.io/controller-runtime@|github\.com/kubernetes-sigs/controller-runtime/(?:blob|tree|releases/tag)/|\bcontroller-runtime\s+)"
        r"(v\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?)"
    ),
}


def module_versions() -> dict[str, str]:
    versions: dict[str, str] = {}
    for line in (ROOT / "go.mod").read_text(encoding="utf-8").splitlines():
        fields = line.split()
        if len(fields) >= 2 and fields[0] in MODULES.values():
            versions[fields[0]] = fields[1]
    missing = set(MODULES.values()) - versions.keys()
    if missing:
        raise RuntimeError(f"go.mod is missing required module(s): {', '.join(sorted(missing))}")
    return {name: versions[module] for name, module in MODULES.items()}


def tracked_files() -> list[Path]:
    output = subprocess.check_output(
        ["git", "-C", str(ROOT), "ls-files", "--cached", "--others", "--exclude-standard", "-z"]
    )
    paths = [ROOT / path.decode() for path in output.split(b"\0") if path]
    return [path for path in paths if path.is_file() and path.name not in {"go.mod", "go.sum"}]


def main() -> int:
    expected = module_versions()
    errors: list[str] = []
    for path in tracked_files():
        try:
            content = path.read_text(encoding="utf-8")
        except (UnicodeDecodeError, OSError):
            continue
        for module, pattern in REFERENCE_PATTERNS.items():
            for match in pattern.finditer(content):
                found = match.group(1)
                if found != expected[module]:
                    line = content.count("\n", 0, match.start()) + 1
                    errors.append(
                        f"{path.relative_to(ROOT)}:{line}: {module} {found} does not match go.mod ({expected[module]})"
                    )
    if errors:
        print("Documented Kubernetes module versions are out of sync with go.mod:", file=sys.stderr)
        print("\n".join(errors), file=sys.stderr)
        return 1
    print("Documented client-go and controller-runtime versions match go.mod.")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (OSError, subprocess.CalledProcessError, RuntimeError) as error:
        print(f"version check failed: {error}", file=sys.stderr)
        raise SystemExit(2)
