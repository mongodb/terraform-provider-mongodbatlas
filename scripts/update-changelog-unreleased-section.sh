#!/bin/bash

# Updates Unreleased section of CHANGELOG.md by generating content with all commited changelog entry files defined after last release.
# Content of existing unreleased header and previous releases is not modified.

set -o errexit
set -o nounset

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__parent="$(dirname "$__dir")"

CHANGELOG_FILE_NAME="CHANGELOG.md"
CHANGELOG_TMP_FILE_NAME="CHANGELOG.tmp"
TARGET_SHA=$(git rev-parse HEAD)

PREVIOUS_RELEASE_TAG=$(git describe --abbrev=0 --match='v*.*.*' --tags)
PREVIOUS_RELEASE_SHA=$(git rev-list -n 1 "$PREVIOUS_RELEASE_TAG")

if [ "$TARGET_SHA" == "$PREVIOUS_RELEASE_SHA" ]; then
  echo "Nothing to do"
  exit 0
fi

# contains all content of CHANGELOG.md starting from the last release (excludes Unreleased section)
PREVIOUS_CHANGELOG=$(sed -n -e "/## $(echo "${PREVIOUS_RELEASE_TAG}" | tr -d v)/,\$p" "$__parent"/$CHANGELOG_FILE_NAME)

# this if then clause is only defined to handle legacy format, can be removed after coming release (succeeding 1.15.3)
if [ -z "$PREVIOUS_CHANGELOG" ]
then
    PREVIOUS_CHANGELOG=$(sed -n -e "/## \[${PREVIOUS_RELEASE_TAG}\]/,\$p" "$__parent"/$CHANGELOG_FILE_NAME)
fi

if [ -z "$PREVIOUS_CHANGELOG" ]
then
    echo "Unable to locate previous changelog contents."
    exit 1
fi

# changelog-build -local-fs performs internal git checkouts; restore caller's HEAD afterward.
ORIGINAL_HEAD_REF=$(git -C "$__parent" symbolic-ref -q --short HEAD || true)
ORIGINAL_HEAD_SHA=$(git -C "$__parent" rev-parse HEAD)

restore_head() {
  if [ -n "$ORIGINAL_HEAD_REF" ]; then
    git -C "$__parent" checkout --quiet "$ORIGINAL_HEAD_REF"
  else
    git -C "$__parent" checkout --quiet "$ORIGINAL_HEAD_SHA"
  fi
}

trap restore_head EXIT

CHANGELOG=$("$(go env GOPATH)"/bin/changelog-build -this-release "$TARGET_SHA" \
                      -last-release "$PREVIOUS_RELEASE_SHA" \
                      -git-dir "$__parent" \
                      -entries-dir .changelog \
                      -changelog-template "$__dir/changelog/changelog.tmpl" \
                      -note-template "$__dir/changelog/release-note.tmpl" \
                      -local-fs)

restore_head
trap - EXIT

if [ -z "$CHANGELOG" ]
then
    echo "No changelog generated."
    exit 0
fi

rm -f $CHANGELOG_TMP_FILE_NAME

sed -n -e "1{/## /p;}" "$__parent"/$CHANGELOG_FILE_NAME > $CHANGELOG_TMP_FILE_NAME # places header (first line) of current changelog into new one

{
    echo "$CHANGELOG" 
    echo 
    echo "$PREVIOUS_CHANGELOG"
} >> $CHANGELOG_TMP_FILE_NAME

cp $CHANGELOG_TMP_FILE_NAME $CHANGELOG_FILE_NAME
rm $CHANGELOG_TMP_FILE_NAME

echo "Successfully generated changelog."
exit 0
