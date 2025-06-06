#!/usr/bin/env bash
set -euo pipefail

: "${RELEASE_VERSION:?RELEASE_VERSION environment variable not set}"
DATE=$(date +'%Y-%m-%d')

echo "Augmenting SBOM..."
docker run \
	--pull=always \
	--platform="linux/amd64" \
	--rm \
	-v "${PWD}:/pwd" \
	-e KONDUKTO_TOKEN \
	"$SILKBOMB_IMG" \
	augment \
	--sbom-in "/pwd/compliance/sbom.json" \
	--repo "$KONDUKTO_REPO" \
	--branch "$KONDUKTO_BRANCH_PREFIX-linux-arm64" \
	--sbom-out "/pwd/compliance/augmented-sbom-v${RELEASE_VERSION}-${DATE}.json"
