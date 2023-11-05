#!/usr/bin/env bash

set -euo pipefail

LATEST_SDK_RELEASE=$(curl -sSfL -X GET  https://api.github.com/repos/mongodb/atlas-sdk-go/releases/latest | jq -r '.tag_name' | cut -d '.' -f 1)
echo  "==> Updating SDK to latest PREVIEW major version $LATEST_SDK_RELEASE"

gomajor get "go.mongodb.org/atlas-sdk/$LATEST_SDK_RELEASE@dev-latest"
ret=$?
if [ $ret -ne 0 ]; then
	go mod tidy
	echo "Finished update - This is SDK Preview, ***** DONT MERGE *******"
else
	echo "Failed to update - You may need to update to the latest SDK first using the 'Update SDK' GitHub Action"
	exit 1
fi	
