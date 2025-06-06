#!/usr/bin/env bash
set -euo pipefail

echo "Uploading SBOMs..."
docker run --rm \
  -v "$PWD:/pwd" \
  -e KONDUKTO_TOKEN \
  "$SILKBOMB_IMG" \
  upload \
  --sbom-in /pwd/compliance/sbom.json \
  --repo "$KONDUKTO_REPO"  \
  --branch "$KONDUKTO_BRANCH_PREFIX-linux-arm64"