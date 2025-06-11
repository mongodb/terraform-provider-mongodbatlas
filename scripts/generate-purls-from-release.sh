#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <release_version>"
  exit 1
fi

RELEASE_VERSION="$1"
OUT_DIR="compliance"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXTRACT_PURL_SCRIPT="${SCRIPT_DIR}/extract-purls.sh"

if [ ! -x "$EXTRACT_PURL_SCRIPT" ]; then
  echo "extract-purls.sh not found or not executable"
  exit 1
fi

mkdir -p "$OUT_DIR"

# Define platforms and temp files
PLATFORMS=(
  "linux_amd64"
  "darwin_amd64"
  "windows_amd64"
)
BIN_PATHS=()
PURL_FILES=()

for PLATFORM in "${PLATFORMS[@]}"; do
  ZIP_FILE="release-${PLATFORM}.zip"
  case "$PLATFORM" in
    linux_amd64)
      BIN_PATH="./terraform-provider-mongodbatlas_v${RELEASE_VERSION}"
      ;;
    darwin_amd64)
      BIN_PATH="./terraform-provider-mongodbatlas_v${RELEASE_VERSION}"
      ;;
    windows_amd64)
      BIN_PATH="./terraform-provider-mongodbatlas_v${RELEASE_VERSION}.exe"
      ;;
  esac
  PURL_FILE="${OUT_DIR}/purls-${PLATFORM}.txt"
  BIN_PATHS+=("$BIN_PATH")
  PURL_FILES+=("$PURL_FILE")

  # Download
  curl -L "https://github.com/mongodb/terraform-provider-mongodbatlas/releases/download/v${RELEASE_VERSION}/terraform-provider-mongodbatlas_${RELEASE_VERSION}_${PLATFORM}.zip" \
    -o "$ZIP_FILE"
  # Extract
  unzip -o "$ZIP_FILE"
  # Extract PURLs
  "$EXTRACT_PURL_SCRIPT" "$BIN_PATH" "$PURL_FILE"
  # Clean up zip and extracted bin after use
  rm -f "$ZIP_FILE"
  rm -f "$BIN_PATH"
done

# Combine, sort, and deduplicate
cat "${PURL_FILES[@]}" | LC_ALL=C sort | uniq > "${OUT_DIR}/purls.txt"
cat "${OUT_DIR}/purls.txt"

# Clean up temp purl files
rm -f "${PURL_FILES[@]}" 