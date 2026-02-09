#!/usr/bin/env sh
set -eu

ALPINE_VERSION="3.23.3"
ARCH="x86_64"
ROOTFS_BASENAME="alpine-minirootfs-${ALPINE_VERSION}-${ARCH}"
ARCHIVE_NAME="${ROOTFS_BASENAME}.tar.gz"
ROOTFS_DIR="chroot/${ROOTFS_BASENAME}"
BASE_URL="https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/${ARCH}"
ARCHIVE_URL="${BASE_URL}/${ARCHIVE_NAME}"

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
  echo "Usage: ./scripts/setup-chroot.sh [--force]"
  echo "  --force   Re-download and re-extract even if target directory already exists"
  exit 0
fi

FORCE=0
if [ "${1:-}" = "--force" ]; then
  FORCE=1
fi

if [ "$FORCE" -eq 0 ] && [ -d "$ROOTFS_DIR" ]; then
  echo "Rootfs already exists at: $ROOTFS_DIR"
  echo "Nothing to do. Use --force to re-download and re-extract."
  exit 0
fi

mkdir -p chroot

TMP_ARCHIVE="$(mktemp /tmp/alpine-minirootfs.XXXXXX.tar.gz)"
cleanup() {
  rm -f "$TMP_ARCHIVE"
}
trap cleanup EXIT

echo "Downloading: $ARCHIVE_URL"
if command -v curl >/dev/null 2>&1; then
  curl -fL "$ARCHIVE_URL" -o "$TMP_ARCHIVE"
elif command -v wget >/dev/null 2>&1; then
  wget -O "$TMP_ARCHIVE" "$ARCHIVE_URL"
else
  echo "Error: neither curl nor wget is installed." >&2
  exit 1
fi

if [ -d "$ROOTFS_DIR" ]; then
  rm -rf "$ROOTFS_DIR"
fi

echo "Extracting into chroot/"
tar -xzf "$TMP_ARCHIVE" -C chroot

echo "Done. Rootfs available at: $ROOTFS_DIR"
