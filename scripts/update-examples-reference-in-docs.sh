#!/bin/bash

set -euo pipefail

: "${1?"Tag of new release must be provided"}"

RELEASE_TAG=$1

# Define the old URL pattern and new URL
OLD_URL_PATTERN="https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/[a-zA-Z0-9._-]*\/examples"
NEW_URL="https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/$RELEASE_TAG\/examples"

FILES=()

# 1) docs/index.md
FILES+=("./docs/index.md")

# 2) collect docs/resources (.md), templates/resources (.md.tmpl),
#    docs/data-sources (.md), templates/data-sources (.md.tmpl)
TARGETS=( "./docs/resources|*.md" "./templates/resources|*.md.tmpl" "./docs/data-sources|*.md" "./templates/data-sources|*.md.tmpl" )

for TARGET in "${TARGETS[@]}"; do
  IFS='|' read -r DIR PATTERN <<< "$TARGET"
  if [ -d "$DIR" ]; then
    while IFS= read -r -d '' f; do
      FILES+=("$f")
    done < <(find "$DIR" -type f -name "$PATTERN" -print0)
  fi
done

# Update links in each target file
for FILE_PATH in "${FILES[@]}"; do
  TMP_FILE_NAME="${FILE_PATH}.tmp"
  rm -f "$TMP_FILE_NAME"

  # Use sed to update the URL and write to temporary file
  sed "s|$OLD_URL_PATTERN|$NEW_URL|g" "$FILE_PATH" > "$TMP_FILE_NAME"

  # Move temporary file to original file
  mv "$TMP_FILE_NAME" "$FILE_PATH"

  echo "Link updated successfully in $FILE_PATH"
done
