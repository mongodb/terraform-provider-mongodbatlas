#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <binary_path> <output_file>"
  exit 1
fi

BINARY_PATH="$1"
OUTPUT_FILE="$2"

go version -m "$BINARY_PATH" | \
  awk '$1 == "dep" || $1 == "=>" { print "pkg:golang/" $2 "@" $3 }' | \
  LC_ALL=C sort > "$OUTPUT_FILE" 