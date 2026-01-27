#!/usr/bin/env bash
# Fetch OpenAPI spec from an MMS branch
#
# Usage: fetch-mms-spec.sh <branch_name>
# Example: fetch-mms-spec.sh CLOUDP-375419

set -euo pipefail

BRANCH_NAME="${1:-}"
SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.json"

if [[ -z "$BRANCH_NAME" ]]; then
    echo "Error: Branch name is required"
    echo "Usage: $0 <branch_name>"
    echo "Example: $0 CLOUDP-375419"
    exit 1
fi

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

# Create directory if needed
mkdir -p "$(dirname "${SPEC_PATH}")"

# Get GitHub token
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed"
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub CLI. Run 'gh auth login' first."
    exit 1
fi

GITHUB_TOKEN=$(gh auth token)

echo "Fetching OpenAPI spec from 10gen/mms branch: ${BRANCH_NAME}..."

SPEC_URL="https://raw.githubusercontent.com/10gen/mms/refs/heads/${BRANCH_NAME}/server/openapi/services/openapi-mms.json"

if ! curl -fsSL \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    "${SPEC_URL}" \
    -o "${SPEC_PATH}"; then
    echo "Error: Failed to fetch spec from ${SPEC_URL}"
    echo "Make sure the branch exists and you have access to the 10gen/mms repository"
    exit 1
fi

echo "Successfully saved spec to ${SPEC_PATH}"
echo "File size: $(wc -c < "${SPEC_PATH}" | tr -d ' ') bytes"
