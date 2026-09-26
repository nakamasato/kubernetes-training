#!/usr/bin/env python3
"""Read the version link in a tool's README, which is updated by Renovate."""
import re
import sys
from pathlib import Path


def version(tool):
    if tool not in {"helm", "kustomize"}:
        raise ValueError(f"Unsupported tool: {tool}")
    readme = Path(__file__).resolve().parents[2] / "contents" / tool / "README.md"
    match = re.search(r"^\[(v\d+\.\d+\.\d+)\]\(https://github\.com/", readme.read_text(), re.MULTILINE)
    if not match:
        raise ValueError(f"No version link in {readme}")
    return match.group(1)


if __name__ == "__main__":
    print(version(sys.argv[1]))
