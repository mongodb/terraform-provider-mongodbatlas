#!/usr/bin/env bash
export GH_TOKEN="${github_token}"

release_tag=$(gh release list --limit 1 --json tagName | jq -r '.[0].tagName')
echo "DEBUG: extracted release_tag: $release_tag"

if [[ -z "$release_tag" || "$release_tag" == "null" ]]; then
    echo "ERROR: Failed to extract valid release tag"
    exit 1
fi

mkdir -p release_artifacts

echo "Waiting 20 minutes before checking for release artifacts..."
sleep 1200 # 20 minutes initial wait

echo "Checking for terraform-provider .zip artifacts in GitHub release..."

max_attempts=5
attempt=1
artifact_found=false

while [ $attempt -le $max_attempts ]; do
    echo "Attempt $attempt: checking for artifacts..."
    gh release view "${release_tag}" --json assets --jq '.assets[].name' | grep -q 'terraform-provider-mongodbatlas.*\.zip'
    if [ $? -eq 0 ]; then
        echo "Artifacts found. Proceeding to download..."
        artifact_found=true
        break
    fi

    echo "Artifacts not available yet. Sleeping for 2 minutes before retry..."
    sleep 120
    attempt=$((attempt + 1))
done

if [ "$artifact_found" != true ]; then
    echo "ERROR: No terraform-provider .zip artifacts found after waiting. Exiting..."
    gh release view "${release_tag}" --json assets --jq '.assets[].name'
    exit 1
fi

mkdir -p release_artifacts
gh release download "${release_tag}" -p "terraform-provider-mongodbatlas*.zip" -D ./release_artifacts/

echo "Downloaded artifacts:"
ls -la release_artifacts/

echo "Removing any source code archives..."
rm -f release_artifacts/Source* release_artifacts/source* 2>/dev/null || true

echo "Final artifacts to trace:"
ls -la release_artifacts/

artifact_count=$(ls -1 release_artifacts/*.zip 2>/dev/null | wc -l)
if [ $artifact_count -eq 0 ]; then
    echo "ERROR: No terraform-provider .zip artifacts found for release ${release_tag}"
    echo "Available files in release:"
    gh release view "${release_tag}" --json assets --jq '.assets[].name'
    exit 1
fi

echo "Found $artifact_count terraform-provider artifacts to trace"

cat <<EOT >trace-expansions.yml
release_version: "$release_tag"
EOT
