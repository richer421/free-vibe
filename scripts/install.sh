#!/usr/bin/env bash
set -euo pipefail

REPO="${FREEVIBE_REPO:-richer421/free-vibe}"
BINARY_NAME="${FREEVIBE_BINARY_NAME:-freevibe}"
INSTALL_DIR="${FREEVIBE_INSTALL_DIR:-$HOME/.local/bin}"
MODE="install"
VERSION="${FREEVIBE_VERSION:-latest}"

usage() {
  cat <<USAGE
Usage: install.sh [options]

Install or update FreeVibe CLI from GitHub Release.

Options:
  --update                 Update mode (same install action with update message)
  --version <tag>          Version tag like v0.1.0 (default: latest)
  --install-dir <dir>      Install directory (default: ~/.local/bin)
  --repo <owner/name>      GitHub repo (default: richer421/free-vibe)
  -h, --help               Show this help
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --update)
      MODE="update"
      shift
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --repo)
      REPO="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin|linux) ;;
  *)
    echo "unsupported OS: $os" >&2
    exit 1
    ;;
esac

arch_raw="$(uname -m)"
case "$arch_raw" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "unsupported arch: $arch_raw" >&2
    exit 1
    ;;
esac

if [[ "$VERSION" == "latest" ]]; then
  api_url="https://api.github.com/repos/${REPO}/releases/latest"
  VERSION="$(curl -fsSL "$api_url" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
  if [[ -z "$VERSION" ]]; then
    echo "failed to resolve latest version from ${api_url}" >&2
    exit 1
  fi
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

echo "[freevibe] ${MODE} ${VERSION} (${os}/${arch})"
asset_candidates=(
  "${BINARY_NAME}_${VERSION}_${os}_${arch}.tar.gz"
)
# Backward compatibility for early releases using non-tag version text.
if [[ "$VERSION" == v* ]]; then
  asset_candidates+=("${BINARY_NAME}_${VERSION#v}_${os}_${arch}.tar.gz")
fi

downloaded_asset=""
for candidate in "${asset_candidates[@]}"; do
  download_url="https://github.com/${REPO}/releases/download/${VERSION}/${candidate}"
  echo "[freevibe] downloading: ${download_url}"
  if curl -fsSL "$download_url" -o "$tmp_dir/$candidate" 2>/dev/null; then
    downloaded_asset="$candidate"
    break
  fi
done

if [[ -z "$downloaded_asset" ]]; then
  echo "[freevibe] no matching release asset found for ${VERSION} (${os}/${arch})" >&2
  echo "[freevibe] tried:" >&2
  for candidate in "${asset_candidates[@]}"; do
    echo "  - ${candidate}" >&2
  done
  exit 1
fi

tar -xzf "$tmp_dir/$downloaded_asset" -C "$tmp_dir"
mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "[freevibe] installed to: $INSTALL_DIR/$BINARY_NAME"
if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
  echo "[freevibe] hint: add to PATH -> export PATH=\"$INSTALL_DIR:\$PATH\""
fi

"$INSTALL_DIR/$BINARY_NAME" version
