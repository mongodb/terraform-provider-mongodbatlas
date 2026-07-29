#!/usr/bin/env bash
# Full workflow: Fetch, flatten, and generate a resource from any environment
#
# Usage: autogen-resource.sh <environment> <resource_name>
#
# Environments:
#   prod, production     - Production spec from mongodb/openapi (main branch)
#   dev, development     - Development spec from mongodb/openapi (dev branch)
#   mms:<branch>         - MMS branch (e.g., mms:CLOUDP-375419)
#   mms-pr:<number>      - MMS PR number (e.g., mms-pr:153849)
#
# Examples:
#   autogen-resource.sh prod log_integration
#   autogen-resource.sh dev alert_configuration_api
#   autogen-resource.sh mms:CLOUDP-375419 log_integration
#   autogen-resource.sh mms-pr:153849 log_integration

set -euo pipefail

ENVIRONMENT="${1:-}"
RESOURCE_NAME="${2:-}"

if [[ -z "$ENVIRONMENT" ]] || [[ -z "$RESOURCE_NAME" ]]; then
    echo "Error: Both environment and resource name are required"
    echo ""
    echo "Usage: $0 <environment> <resource_name>"
    echo ""
    echo "Environments:"
    echo "  prod, production     - Production spec"
    echo "  dev, development     - Development spec"
    echo "  mms:<branch>         - MMS branch (e.g., mms:CLOUDP-375419)"
    echo "  mms-pr:<number>      - MMS PR number (e.g., mms-pr:153849)"
    echo ""
    echo "Examples:"
    echo "  $0 prod log_integration"
    echo "  $0 mms:CLOUDP-375419 log_integration"
    echo "  $0 mms-pr:153849 log_integration"
    exit 1
fi

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo ""
echo "=========================================="
echo "Autogenerating Terraform Resource"
echo "=========================================="
echo "  Environment: ${ENVIRONMENT}"
echo "  Resource:    ${RESOURCE_NAME}"
echo "=========================================="
echo ""

# Step 1: Fetch
echo "Step 1/3: Fetching OpenAPI spec..."
"${SCRIPT_DIR}/fetch-spec.sh" "${ENVIRONMENT}"
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
