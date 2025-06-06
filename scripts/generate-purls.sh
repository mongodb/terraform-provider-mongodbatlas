#!/usr/bin/env bash
set -euo pipefail
: "${LINKER_FLAGS:=}"

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
go version -m "${LINUX_BIN}" | awk '$1 == "dep" || $1 == "=>" { print "pkg:golang/" $2 "@" $3 }' | LC_ALL=C sort > "${PURL_LINUX}"

# Build and extract for Darwin
GOOS=darwin GOARCH=amd64 go build -ldflags "${LINKER_FLAGS}" -o "${DARWIN_BIN}"
go version -m "${DARWIN_BIN}" | awk '$1 == "dep" || $1 == "=>" { print "pkg:golang/" $2 "@" $3 }' | LC_ALL=C sort > "${PURL_DARWIN}"

# Build and extract for Windows
GOOS=windows GOARCH=amd64 go build -ldflags "${LINKER_FLAGS}" -o "${WIN_BIN}"
go version -m "${WIN_BIN}" | awk '$1 == "dep" || $1 == "=>" { print "pkg:golang/" $2 "@" $3 }' | LC_ALL=C sort > "${PURL_WIN}"

# Combine, sort, and deduplicate
cat "${PURL_LINUX}" "${PURL_DARWIN}" "${PURL_WIN}" | LC_ALL=C sort | uniq > "${PURL_ALL}"

# Clean up temp files
rm -f "${LINUX_BIN}" "${DARWIN_BIN}" "${WIN_BIN}" "${PURL_LINUX}" "${PURL_DARWIN}" "${PURL_WIN}"