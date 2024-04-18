#!/usr/bin/env bash
set -euo pipefail

: "${1?"Tag of new release must be provided"}"

RELEASE_TAG=$1
RELEASE_NUMBER=$(echo "${RELEASE_TAG}" | tr -d v)

IFS='.' read -r MAJOR MINOR PATCH <<< "$RELEASE_NUMBER"

# Check if it's a major release (patch version is 0)
if [ "$PATCH" -eq 0 ]; then
    UPGRADE_GUIDE_PATH="website/docs/guides/$MAJOR.$MINOR.$PATCH-upgrade-guide.html.markdown"
    echo "Checking for the presence of $UPGRADE_GUIDE_PATH"
    if [ ! -f "$UPGRADE_GUIDE_PATH" ]; then
        echo "Stopping release procees, upgrade guide $UPGRADE_GUIDE_PATH does not exist."
        exit 1
    else
        echo "Upgrade guide $UPGRADE_GUIDE_PATH exists."
    fi
else
    echo "No upgrade guide needed."
fi
