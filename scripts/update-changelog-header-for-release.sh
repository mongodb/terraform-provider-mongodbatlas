#!/usr/bin/env bash
set -euo pipefail

: "${1?"Tag of new release must be provided"}"

CHANGELOG_FILE_PATH=CHANGELOG.md
RELEASE_TAG=$1
RELEASE_NUMBER=$(echo "${RELEASE_TAG}" | tr -d v)

# exit out if changelog already has the header updated with version number being released.
if grep -q "## $RELEASE_NUMBER (" "$CHANGELOG_FILE_PATH"; then
    echo "CHANGELOG already has a header defined for $RELEASE_NUMBER, no changes made to changelog."
    exit 0
fi

# Prepare the new version header
TODAYS_DATE=$(date "+%B %d, %Y")  # Format the date as "Month day, Year"
NEW_RELEASE_HEADER="## $RELEASE_NUMBER ($TODAYS_DATE)"

# Insert the new version header after the "(Unreleased)" line
sed -i "" -e "/(Unreleased)/a \\
\\
$NEW_RELEASE_HEADER" $CHANGELOG_FILE_PATH

echo "Changelog updated successfully defining header for new $RELEASE_TAG release."
