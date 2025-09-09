#!/usr/bin/env bash
set -euo pipefail

: "${1?"Tag of new release must be provided"}"

RELEASE_TAG=$1
RELEASE_NUMBER=$(echo "${RELEASE_TAG}" | tr -d v)

IFS='.' read -r MAJOR MINOR PATCH <<< "$RELEASE_NUMBER"

# Check if it's a major release (minor and patch versions are 0)
if [ "$PATCH" -eq 0 ] && [ "$MINOR" -eq 0 ]; then
    UPGRADE_GUIDE_PATH="docs/guides/$MAJOR.$MINOR.$PATCH-upgrade-guide.md"
    echo "Checking for the presence of $UPGRADE_GUIDE_PATH"
    if [ ! -f "$UPGRADE_GUIDE_PATH" ]; then
        echo "Stopping release process, upgrade guide $UPGRADE_GUIDE_PATH does not exist. Please visit our docs for more details: https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/RELEASING.md"
        exit 1
    else
        echo "Upgrade guide $UPGRADE_GUIDE_PATH exists."
    fi
else
    echo "No upgrade guide needed."
fi
