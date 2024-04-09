#!/usr/bin/env bash
#
# Generate a changelog entry for a suggested PR by grabbing the next available
# auto incrementing ID in GitHub. User can choose to use another PR number

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
next_pr=$(($current_pr + 1))

echo "==> What is the new changelog entry? Suggested is $next_pr (next PR number)"
read next_pr

changelog-entry -pr $next_pr -dir ".changelog"

echo
echo "Successfully created $changelog_path. Don't forget to commit it and open the PR!"