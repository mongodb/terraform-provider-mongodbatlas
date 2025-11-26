#!/usr/bin/env bash
set -euo pipefail

SOURCE_SPEC_URL="${SOURCE_SPEC_URL:-https://raw.githubusercontent.com/mongodb/openapi/dev/openapi/v2.yaml}" # TODO: change to main once API is in production
SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.yml"
FLATTENED_SPEC_PATH="tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"

mkdir -p "$(dirname "${SPEC_PATH}")"

echo "Downloading OpenAPI spec from ${SOURCE_SPEC_URL}..."
curl -fsSL "${SOURCE_SPEC_URL}" -o "${SPEC_PATH}"

echo "Saved spec to ${SPEC_PATH}"

mkdir -p "$(dirname "${FLATTENED_SPEC_PATH}")"

echo "Flattening OpenAPI spec:"
echo "  input : ${SPEC_PATH}"
echo "  output: ${FLATTENED_SPEC_PATH}"

npx -y --package github:mongodb/atlas-sdk-go atlas-openapi-transformer flatten -i "${SPEC_PATH}" -o "${FLATTENED_SPEC_PATH}"

echo "Flattened spec saved to ${FLATTENED_SPEC_PATH}"
