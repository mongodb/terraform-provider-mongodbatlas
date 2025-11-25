#!/usr/bin/env bash
set -euo pipefail

echo "Generating SBOM..."
docker run --rm \
  -v "$PWD:/pwd" \
  "$SILKBOMB_IMG" \
  update \
  --purls /pwd/compliance/purls.txt \
  --sbom-out /pwd/compliance/sbom.json