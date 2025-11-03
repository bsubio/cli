#!/usr/bin/env python3
"""
Extract release notes from CHANGES.yaml for a specific version.
Usage: python extract-changelog.py <version>
Example: python extract-changelog.py 0.1.0
"""

import sys
import yaml
from pathlib import Path


def extract_changelog(version: str, changes_file: Path = Path("CHANGES.yaml")) -> str:
    """Extract formatted changelog for a specific version."""

    if not changes_file.exists():
        return f"No CHANGES.yaml found. Using version {version}."

    with open(changes_file) as f:
        data = yaml.safe_load(f)

    releases = data.get("releases", {})

    if version not in releases:
        return f"Version {version} not found in CHANGES.yaml."

    release = releases[version]

    # Build markdown changelog
    lines = []

    # Add summary
    if summary := release.get("release_summary", "").strip():
        lines.append(summary)
        lines.append("")

    # Add breaking changes (if any)
    if breaking := release.get("breaking_changes", []):
        lines.append("## ⚠️ Breaking Changes")
        for change in breaking:
            lines.append(f"- {change}")
        lines.append("")

    # Add major changes
    if major := release.get("major_changes", []):
        lines.append("## 🎉 Major Changes")
        for change in major:
            lines.append(f"- {change}")
        lines.append("")

    # Add minor changes
    if minor := release.get("minor_changes", []):
        lines.append("## ✨ Minor Changes")
        for change in minor:
            lines.append(f"- {change}")
        lines.append("")

    # Add bugfixes
    if bugfixes := release.get("bugfixes", []):
        lines.append("## 🐛 Bug Fixes")
        for fix in bugfixes:
            lines.append(f"- {fix}")
        lines.append("")

    # Add security fixes
    if security := release.get("security_fixes", []):
        lines.append("## 🔒 Security Fixes")
        for fix in security:
            lines.append(f"- {fix}")
        lines.append("")

    # Add known issues
    if issues := release.get("known_issues", []):
        lines.append("## ⚠️ Known Issues")
        for issue in issues:
            lines.append(f"- {issue}")
        lines.append("")

    # Add metadata footer
    lines.append("---")
    if release_date := release.get("release_date"):
        lines.append(f"*Released: {release_date}*")

    return "\n".join(lines)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python extract-changelog.py <version>", file=sys.stderr)
        print("Example: python extract-changelog.py 0.1.0", file=sys.stderr)
        sys.exit(1)

    version = sys.argv[1]
    # Remove 'v' prefix if present
    if version.startswith("v"):
        version = version[1:]

    changelog = extract_changelog(version)
    print(changelog)
