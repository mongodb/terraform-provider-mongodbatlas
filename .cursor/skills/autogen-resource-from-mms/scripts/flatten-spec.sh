#!/usr/bin/env bash
# Flatten the OpenAPI spec for code generation
#
# Usage: flatten-spec.sh [input_file] [output_file]
# Defaults:
#   input:  tools/codegen/atlasapispec/raw-multi-version-api-spec.json
#   output: tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml

set -euo pipefail

INPUT_SPEC="${1:-tools/codegen/atlasapispec/raw-multi-version-api-spec.json}"
OUTPUT_SPEC="${2:-tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml}"

# Ensure we're in the repo root
if [[ ! -f "tools/codegen/config.yml" ]]; then
    echo "Error: Must be run from the terraform-provider-mongodbatlas repository root"
    exit 1
fi

if [[ ! -f "$INPUT_SPEC" ]]; then
    echo "Error: Input spec not found at ${INPUT_SPEC}"
    echo "Run fetch-mms-spec.sh first to download the spec"
    exit 1
fi

# Create output directory if needed
mkdir -p "$(dirname "${OUTPUT_SPEC}")"

echo "Flattening OpenAPI spec:"
echo "  Input:  ${INPUT_SPEC}"
echo "  Output: ${OUTPUT_SPEC}"

if ! npx -y --package github:mongodb/atlas-sdk-go atlas-openapi-transformer flatten \
    -i "${INPUT_SPEC}" \
    -o "${OUTPUT_SPEC}"; then
    echo "Error: Failed to flatten the OpenAPI spec"
    exit 1
fi

echo "Successfully flattened spec to ${OUTPUT_SPEC}"
echo "File size: $(wc -c < "${OUTPUT_SPEC}" | tr -d ' ') bytes"
