#!/usr/bin/env bash
set -euo pipefail

SPEC_URL="https://raw.githubusercontent.com/mongodb/openapi/main/openapi/v2/openapi-2025-03-12.yaml"
SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.yml"
FLATTENED_SPEC_PATH="tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"

mkdir -p "$(dirname "${SPEC_PATH}")"

echo "Downloading OpenAPI spec from ${SPEC_URL}..."
curl -fsSL "${SPEC_URL}" -o "${SPEC_PATH}"

echo "Saved spec to ${SPEC_PATH}"

mkdir -p "$(dirname "${FLATTENED_SPEC_PATH}")"

echo "Flattening OpenAPI spec:"
echo "  input : ${SPEC_PATH}"
echo "  output: ${FLATTENED_SPEC_PATH}"

npx -y --package github:mongodb/atlas-sdk-go#CLOUDP-356932 atlas-openapi-transformer flatten -i "${SPEC_PATH}" -o "${FLATTENED_SPEC_PATH}"

echo "Flattened spec saved to ${FLATTENED_SPEC_PATH}"
