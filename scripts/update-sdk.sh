#!/usr/bin/env bash

set -euo pipefail

LATEST_SDK_RELEASE=$(curl -sSfL -X GET  https://api.github.com/repos/mongodb/atlas-sdk-go/releases/latest | jq -r '.tag_name' | cut -d '.' -f 1)
echo  "==> Updating SDK to latest major version $LATEST_SDK_RELEASE"
gomajor get "go.mongodb.org/atlas-sdk/$LATEST_SDK_RELEASE@latest"
go mod tidy
echo "Done, remember to update build/ci/library_owners.json"
