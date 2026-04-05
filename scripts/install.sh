#!/usr/bin/env bash
set -euo pipefail

REPO="${FREEVIBE_REPO:-richer421/free-vibe}"
BINARY_NAME="${FREEVIBE_BINARY_NAME:-freevibe}"
INSTALL_DIR="${FREEVIBE_INSTALL_DIR:-/usr/local/bin}"
VERSION="${FREEVIBE_VERSION:-latest}"

usage() {
  cat <<USAGE
Usage: install.sh [options]

Install or update FreeVibe CLI from GitHub Release.

Options:
  --version <tag>          Version tag like v0.1.0 (default: latest)
  --install-dir <dir>      Install directory (default: /usr/local/bin)
  --repo <owner/name>      GitHub repo (default: richer421/free-vibe)
  -h, --help               Show this help
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
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

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
  echo "invalid version: $VERSION (expected vX.Y.Z)" >&2
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

current_binary="${INSTALL_DIR}/${BINARY_NAME}"
current_version=""
if [[ -x "$current_binary" ]]; then
  current_version="$("$current_binary" version 2>/dev/null | awk '{print $3}' || true)"
elif command -v "$BINARY_NAME" >/dev/null 2>&1; then
  current_binary="$(command -v "$BINARY_NAME")"
  current_version="$("$current_binary" version 2>/dev/null | awk '{print $3}' || true)"
fi

if [[ -n "$current_version" ]]; then
  echo "[freevibe] updating ${BINARY_NAME}: ${current_version} -> ${VERSION} (${os}/${arch})"
else
  echo "[freevibe] installing ${BINARY_NAME} ${VERSION} (${os}/${arch})"
fi

asset="${BINARY_NAME}_${VERSION}_${os}_${arch}.tar.gz"
download_url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"
echo "[freevibe] downloading: ${download_url}"
if ! curl -fsSL "$download_url" -o "$tmp_dir/$asset"; then
  echo "[freevibe] release asset not found: $asset" >&2
  exit 1
fi

tar -xzf "$tmp_dir/$asset" -C "$tmp_dir"
if [[ ! -d "$INSTALL_DIR" ]]; then
  parent_dir="$(dirname "$INSTALL_DIR")"
  if [[ -w "$parent_dir" ]]; then
    mkdir -p "$INSTALL_DIR"
  else
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

if [[ -w "$INSTALL_DIR" ]]; then
  install -m 0755 "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
  sudo install -m 0755 "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "[freevibe] ready: $INSTALL_DIR/$BINARY_NAME"
if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
  echo "[freevibe] hint: add to PATH -> export PATH=\"$INSTALL_DIR:\$PATH\""
fi

"$INSTALL_DIR/$BINARY_NAME" version
