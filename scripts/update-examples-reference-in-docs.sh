#!/bin/bash

set -euo pipefail

: "${1?"Tag of new release must be provided"}"

RELEASE_TAG=$1

# Define the regex pattern for the old and new versioned GitHub URLs (e.g., for /examples or /troubleshooting directories)
OLD_URL_PATTERN="https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/[a-zA-Z0-9._-]*\/"
NEW_URL="https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/$RELEASE_TAG\/"

FILES=()

# 1) docs/index.md
FILES+=("./docs/index.md")

# 2) collect all *.md and *.md.tmpl under target directories.
TARGET_DIRS=(
  "./docs/resources"
  "./templates/resources"
  "./docs/data-sources"
  "./templates/data-sources"
  "./docs/ephemeral-resources"
  "./templates/ephemeral-resources"
  "./docs/guides"
)

for DIR in "${TARGET_DIRS[@]}"; do
  if [ -d "$DIR" ]; then
    while IFS= read -r -d '' f; do
      FILES+=("$f")
    done < <(find "$DIR" -type f \( -name "*.md" -o -name "*.md.tmpl" \) -print0)
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
