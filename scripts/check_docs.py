#!/usr/bin/env python3
"""Check local Markdown file links only; this script makes no network requests."""
from pathlib import Path
import re

root = Path(__file__).resolve().parents[1]
missing = []
for doc in root.rglob("*.md"):
    if ".git" in doc.parts:
        continue
    for target in re.findall(r"\[[^\]]*\]\(([^)]+)\)", doc.read_text()):
        if "://" in target or target.startswith(("#", "mailto:")):
            continue
        file = target.split("#", 1)[0]
        if file and not (doc.parent / file).exists():
            missing.append(f"{doc.relative_to(root)} -> {target}")
if missing:
    raise SystemExit("Missing local links:\n" + "\n".join(missing))
print("Local Markdown file links passed")
