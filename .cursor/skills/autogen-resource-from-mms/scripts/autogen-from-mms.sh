#!/usr/bin/env bash
# Full workflow: Fetch, flatten, and generate a resource from an MMS branch
#
# Usage: autogen-from-mms.sh <branch_or_pr> <resource_name>
# Examples:
#   autogen-from-mms.sh CLOUDP-375419 log_integration
#   autogen-from-mms.sh 153849 log_integration  # PR number

set -euo pipefail

BRANCH_OR_PR="${1:-}"
RESOURCE_NAME="${2:-}"

if [[ -z "$BRANCH_OR_PR" ]] || [[ -z "$RESOURCE_NAME" ]]; then
    echo "Error: Both branch/PR and resource name are required"
    echo ""
    echo "Usage: $0 <branch_or_pr> <resource_name>"
    echo ""
    echo "Examples:"
    echo "  $0 CLOUDP-375419 log_integration      # Using branch name"
    echo "  $0 153849 log_integration             # Using PR number"
    exit 1
fi

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Determine if input is a PR number or branch name
BRANCH_NAME="$BRANCH_OR_PR"
if [[ "$BRANCH_OR_PR" =~ ^[0-9]+$ ]]; then
    echo "Resolving PR #${BRANCH_OR_PR} to branch name..."
    BRANCH_NAME=$(gh pr view "$BRANCH_OR_PR" --repo 10gen/mms --json headRefName -q '.headRefName')
    if [[ -z "$BRANCH_NAME" ]]; then
        echo "Error: Could not resolve PR #${BRANCH_OR_PR} to a branch name"
        exit 1
    fi
    echo "Resolved to branch: ${BRANCH_NAME}"
fi

echo ""
echo "=========================================="
echo "Autogenerating resource from MMS branch"
echo "=========================================="
echo "  Branch:   ${BRANCH_NAME}"
echo "  Resource: ${RESOURCE_NAME}"
echo "=========================================="
echo ""

# Step 1: Fetch
echo "Step 1/3: Fetching OpenAPI spec..."
"${SCRIPT_DIR}/fetch-mms-spec.sh" "${BRANCH_NAME}"
echo ""

# Step 2: Flatten
echo "Step 2/3: Flattening OpenAPI spec..."
"${SCRIPT_DIR}/flatten-spec.sh"
echo ""

# Step 3: Generate
echo "Step 3/3: Generating resource code..."
"${SCRIPT_DIR}/generate-resource.sh" "${RESOURCE_NAME}"
echo ""

echo "=========================================="
echo "Complete!"
echo "=========================================="
