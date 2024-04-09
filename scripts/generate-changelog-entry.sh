#!/usr/bin/env bash
#
# Generate a changelog entry for a suggested PR by grabbing the next available
# auto incrementing ID in GitHub. User can choose to use another PR number
set -euo pipefail

if ! command -v curl &> /dev/null
then
  echo "curl not be found"
  exit 1
fi

if ! command -v jq &> /dev/null
then
  echo "jq not be found"
  exit 1
fi

if ! command -v changelog-entry &> /dev/null
then
  echo "changelog-entry not be found"
  exit 1
fi

current_pr=$(curl -s "https://api.github.com/repos/mongodb/terraform-provider-mongodbatlas/issues?state=all&per_page=1" | jq -r ".[].number")
next_pr=$((current_pr + 1))

echo "==> What PR number should be used for this changelog entry? Leave emtpy to use $next_pr (next PR number)"
read -r changelog_entry

if [ -n "$changelog_entry" ]; then
    pr="$changelog_entry"
else
    pr="$next_pr"
fi

changelog-entry -pr "$pr" -dir ".changelog" -allowed-types-file="./scripts/changelog/allowed-types.txt"

echo
echo "Successfully created $next_pr. Don't forget to commit it and open the PR!"
