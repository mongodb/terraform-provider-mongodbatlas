#!/usr/bin/env bash
# Fetch OpenAPI spec from different environments
#
# Usage: fetch-spec.sh <environment>
#
# Environments:
#   prod, production     - Production spec from mongodb/openapi (main branch)
#   dev, development     - Development spec from mongodb/openapi (dev branch)
#   mms:<branch>         - MMS branch (e.g., mms:CLOUDP-375419)
#   mms-pr:<number>      - MMS PR number (e.g., mms-pr:153849)
#
# Examples:
#   fetch-spec.sh prod
#   fetch-spec.sh dev
#   fetch-spec.sh mms:CLOUDP-375419
#   fetch-spec.sh mms-pr:153849

set -euo pipefail

ENVIRONMENT="${1:-}"

if [[ -z "$ENVIRONMENT" ]]; then
    echo "Error: Environment is required"
    echo ""
    echo "Usage: $0 <environment>"
    echo ""
    echo "Environments:"
    echo "  prod, production     - Production spec from mongodb/openapi (main branch)"
    echo "  dev, development     - Development spec from mongodb/openapi (dev branch)"
    echo "  mms:<branch>         - MMS branch (e.g., mms:CLOUDP-375419)"
    echo "  mms-pr:<number>      - MMS PR number (e.g., mms-pr:153849)"
    exit 1
fi

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

# Determine spec URL and output path based on environment
SPEC_URL=""
SPEC_PATH=""
NEEDS_AUTH=false

case "$ENVIRONMENT" in
    prod|production)
        echo "Environment: Production (mongodb/openapi main branch)"
        SPEC_URL="https://raw.githubusercontent.com/mongodb/openapi/main/openapi/v2.yaml"
        SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.yml"
        ;;
    dev|development)
        echo "Environment: Development (mongodb/openapi dev branch)"
        SPEC_URL="https://raw.githubusercontent.com/mongodb/openapi/dev/openapi/v2.yaml"
        SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.yml"
        ;;
    mms:*)
        BRANCH_NAME="${ENVIRONMENT#mms:}"
        echo "Environment: MMS branch '${BRANCH_NAME}'"
        SPEC_URL="https://raw.githubusercontent.com/10gen/mms/refs/heads/${BRANCH_NAME}/server/openapi/services/openapi-mms.json"
        SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.json"
        NEEDS_AUTH=true
        ;;
    mms-pr:*)
        PR_NUMBER="${ENVIRONMENT#mms-pr:}"
        echo "Environment: MMS PR #${PR_NUMBER}"
        
        if ! command -v gh &> /dev/null; then
            echo "Error: GitHub CLI (gh) is not installed"
            exit 1
        fi
        
        BRANCH_NAME=$(gh pr view "$PR_NUMBER" --repo 10gen/mms --json headRefName -q '.headRefName')
        if [[ -z "$BRANCH_NAME" ]]; then
            echo "Error: Could not resolve PR #${PR_NUMBER} to a branch name"
            exit 1
        fi
        echo "Resolved PR #${PR_NUMBER} to branch: ${BRANCH_NAME}"
        
        SPEC_URL="https://raw.githubusercontent.com/10gen/mms/refs/heads/${BRANCH_NAME}/server/openapi/services/openapi-mms.json"
        SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.json"
        NEEDS_AUTH=true
        ;;
    *)
        echo "Error: Unknown environment '${ENVIRONMENT}'"
        echo ""
        echo "Valid environments: prod, dev, mms:<branch>, mms-pr:<number>"
        exit 1
        ;;
esac

mkdir -p "$(dirname "${SPEC_PATH}")"

echo "Fetching OpenAPI spec..."
echo "  URL: ${SPEC_URL}"
echo "  Output: ${SPEC_PATH}"

CURL_OPTS=("-fsSL")

if [[ "$NEEDS_AUTH" == "true" ]]; then
    if ! gh auth status &> /dev/null; then
        echo "Error: Not authenticated with GitHub CLI. Run 'gh auth login' first."
        exit 1
    fi
    GITHUB_TOKEN=$(gh auth token)
    CURL_OPTS+=("-H" "Authorization: token ${GITHUB_TOKEN}")
fi

if ! curl "${CURL_OPTS[@]}" "${SPEC_URL}" -o "${SPEC_PATH}"; then
    echo "Error: Failed to fetch spec from ${SPEC_URL}"
    exit 1
fi

echo "Successfully saved spec to ${SPEC_PATH}"
echo "File size: $(wc -c < "${SPEC_PATH}" | tr -d ' ') bytes"
