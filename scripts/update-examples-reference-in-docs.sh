#!/bin/bash

set -euo pipefail

: "${1?"Tag of new release must be provided"}"

FILE_PATH="./docs/index.md"
RELEASE_TAG=$1

# Define the old URL pattern and new URL
OLD_URL_PATTERN="\[example configurations\](https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/[a-zA-Z0-9._-]*\/examples)"
NEW_URL="\[example configurations\](https:\/\/github.com\/mongodb\/terraform-provider-mongodbatlas\/tree\/$RELEASE_TAG\/examples)"


TMP_FILE_NAME="docs.tmp"
rm -f $TMP_FILE_NAME

# Use sed to update the URL and write to temporary file
sed "s|$OLD_URL_PATTERN|$NEW_URL|g" "$FILE_PATH" > "$TMP_FILE_NAME"

# Move temporary file to original file
mv "$TMP_FILE_NAME" "$FILE_PATH"

echo "Link updated successfully in $FILE_PATH"
