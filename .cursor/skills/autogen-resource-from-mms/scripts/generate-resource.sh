#!/usr/bin/env bash
# Generate a Terraform resource from the flattened OpenAPI spec
#
# Usage: generate-resource.sh <resource_name>
# Example: generate-resource.sh log_integration

set -euo pipefail

RESOURCE_NAME="${1:-}"

if [[ -z "$RESOURCE_NAME" ]]; then
    echo "Error: Resource name is required"
    echo "Usage: $0 <resource_name>"
    echo "Example: $0 log_integration"
    exit 1
fi

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

# Check if resource is configured
if ! grep -q "^  ${RESOURCE_NAME}:" tools/codegen/config.yml 2>/dev/null; then
    echo "Warning: Resource '${RESOURCE_NAME}' not found in tools/codegen/config.yml"
    echo "Make sure the resource is configured before running generation"
fi

# Check if flattened spec exists
FLATTENED_SPEC="tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"
if [[ ! -f "$FLATTENED_SPEC" ]]; then
    echo "Error: Flattened spec not found at ${FLATTENED_SPEC}"
    echo "Run flatten-spec.sh first"
    exit 1
fi

echo "Generating resource: ${RESOURCE_NAME}..."

if ! go run ./tools/codegen/main.go "${RESOURCE_NAME}"; then
    echo "Error: Code generation failed"
    exit 1
fi

echo ""
echo "Successfully generated resource: ${RESOURCE_NAME}"
echo ""
echo "Generated files:"

# Find and list generated files
MODEL_FILE="tools/codegen/models/${RESOURCE_NAME}.yaml"
if [[ -f "$MODEL_FILE" ]]; then
    echo "  - ${MODEL_FILE}"
fi

# Try to determine package name from the model
if [[ -f "$MODEL_FILE" ]]; then
    PACKAGE_NAME=$(grep "^packageName:" "$MODEL_FILE" 2>/dev/null | cut -d' ' -f2 || echo "")
    if [[ -n "$PACKAGE_NAME" ]]; then
        PACKAGE_DIR="internal/serviceapi/${PACKAGE_NAME}"
        if [[ -d "$PACKAGE_DIR" ]]; then
            for f in "${PACKAGE_DIR}"/*.go; do
                if [[ -f "$f" ]]; then
                    echo "  - ${f}"
                fi
            done
        fi
    fi
fi

echo ""
echo "Next steps:"
echo "  1. Review generated schema for correctness"
echo "  2. Register the resource with the provider (if new)"
echo "  3. Add acceptance tests"
echo "  4. Add documentation"
