#!/usr/bin/env bash

set -euo pipefail

echo "Script will update to dev preview. Note: you need to have latest major version in order to be able to use it"

LATEST_SDK_RELEASE=$(curl -sSfL -X GET  https://api.github.com/repos/mongodb/atlas-sdk-go/releases/latest | jq -r '.tag_name' | cut -d '.' -f 1)
echo  "==> Updating SDK to latest PREVIEW major version $LATEST_SDK_RELEASE"

go get "go.mongodb.org/atlas-sdk/$LATEST_SDK_RELEASE@dev-latest"
go mod tidy
echo
echo "Finished update - This is SDK Preview, ***** DONT MERGE THESE CHANGES *******"
echo "If no files were changed you may need first to update to the latest SDK using the 'Update SDK' GitHub Action"
