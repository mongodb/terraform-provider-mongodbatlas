#!/usr/bin/env bash
set -euo pipefail
: "${LINKER_FLAGS:=}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXTRACT_PURL_SCRIPT="${SCRIPT_DIR}/extract-purls.sh"

if [ ! -x "$EXTRACT_PURL_SCRIPT" ]; then
  echo "extract-purls.sh not found or not executable"
  exit 1
fi

echo "==> Generating purls"

# Define output and temp files
OUT_DIR="compliance"
LINUX_BIN="${OUT_DIR}/bin-linux"
DARWIN_BIN="${OUT_DIR}/bin-darwin"
WIN_BIN="${OUT_DIR}/bin-win.exe"
PURL_LINUX="${OUT_DIR}/purls-linux.txt"
PURL_DARWIN="${OUT_DIR}/purls-darwin.txt"
PURL_WIN="${OUT_DIR}/purls-win.txt"
PURL_ALL="${OUT_DIR}/purls.txt"

# Build and extract for Linux
GOOS=linux GOARCH=amd64 go build -ldflags "${LINKER_FLAGS}" -o "${LINUX_BIN}"
"$EXTRACT_PURL_SCRIPT" "${LINUX_BIN}" "${PURL_LINUX}"

# Build and extract for Darwin
GOOS=darwin GOARCH=amd64 go build -ldflags "${LINKER_FLAGS}" -o "${DARWIN_BIN}"
"$EXTRACT_PURL_SCRIPT" "${DARWIN_BIN}" "${PURL_DARWIN}"

# Build and extract for Windows
GOOS=windows GOARCH=amd64 go build -ldflags "${LINKER_FLAGS}" -o "${WIN_BIN}"
"$EXTRACT_PURL_SCRIPT" "${WIN_BIN}" "${PURL_WIN}"

# Combine, sort, and deduplicate
cat "${PURL_LINUX}" "${PURL_DARWIN}" "${PURL_WIN}" | LC_ALL=C sort | uniq > "${PURL_ALL}"

# Clean up temp files
rm -f "${LINUX_BIN}" "${DARWIN_BIN}" "${WIN_BIN}" "${PURL_LINUX}" "${PURL_DARWIN}" "${PURL_WIN}"