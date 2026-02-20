#!/usr/bin/env bash
set -euo pipefail

# Rebase commits on the current branch (that are ahead of master) and re-sign them with GPG.
# Usage: ./scripts/resign-commits.sh [base-branch]
# Default base: master

BASE_BRANCH="${1:-master}"

echo "==> Checking for uncommitted changes..."
if ! git diff-index --quiet HEAD --; then
  echo "Error: You have uncommitted changes. Please commit or stash them first."
  exit 1
fi

echo "==> Checking if base branch '$BASE_BRANCH' exists..."
if ! git rev-parse --verify "$BASE_BRANCH" >/dev/null 2>&1; then
  echo "Error: Base branch '$BASE_BRANCH' not found."
  exit 1
fi

echo "==> Rebasing onto $BASE_BRANCH and re-signing commits..."
git rebase --exec 'git commit --amend --no-edit -S' "$BASE_BRANCH"

echo "==> Done. All commits have been re-signed."
