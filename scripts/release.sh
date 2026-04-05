#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "usage: scripts/release.sh vX.Y.Z" >&2
  exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
  echo "invalid semver tag: $VERSION" >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "working tree is dirty, commit first" >&2
  exit 1
fi

git fetch origin --tags
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "tag already exists: $VERSION" >&2
  exit 1
fi

git tag -a "$VERSION" -m "release $VERSION"
git push origin "$VERSION"

echo "tag pushed: $VERSION"
echo "GitHub Actions will build and publish release assets automatically."
