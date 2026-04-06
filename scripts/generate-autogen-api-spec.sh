#!/usr/bin/env bash
set -euo pipefail

# Accepts a URL or a local file path to an OpenAPI spec.
# If a local file is provided, it is copied directly; otherwise, it is downloaded from the given URL.
# TODO: Use openapi dev branch, don't merge to master, remove in CLOUDP-372674
SPEC_SOURCE="${1:-https://raw.githubusercontent.com/mongodb/openapi/dev/openapi/v2.yaml}"
SPEC_PATH="tools/codegen/atlasapispec/raw-multi-version-api-spec.yml"
FLATTENED_SPEC_PATH="tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml"

mkdir -p "$(dirname "${SPEC_PATH}")"

if [[ -f "${SPEC_SOURCE}" ]]; then
  echo "Copying local spec from ${SPEC_SOURCE}..."
  cp "${SPEC_SOURCE}" "${SPEC_PATH}"
else
  echo "Downloading OpenAPI spec from ${SPEC_SOURCE}..."
  curl -fsSL "${SPEC_SOURCE}" -o "${SPEC_PATH}"
fi

echo "Saved spec to ${SPEC_PATH}"

mkdir -p "$(dirname "${FLATTENED_SPEC_PATH}")"

echo "Flattening OpenAPI spec:"
echo "  input : ${SPEC_PATH}"
echo "  output: ${FLATTENED_SPEC_PATH}"

npx -y --package github:mongodb/atlas-sdk-go atlas-openapi-transformer flatten -i "${SPEC_PATH}" -o "${FLATTENED_SPEC_PATH}"

echo "Flattened spec saved to ${FLATTENED_SPEC_PATH}"
