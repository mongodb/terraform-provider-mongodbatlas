#!/usr/bin/env bash
set -euo pipefail

release_date=${DATE:-$(date -u '+%Y-%m-%d')}

export DATE="${release_date}"

if [ -z "${AUTHOR:-}" ]; then
  AUTHOR=$(git config user.name)
fi

if [ -z "${VERSION:-}" ]; then
  VERSION=$(git tag --list 'v*' --sort=-taggerdate | head -1 | cut -d 'v' -f 2)
fi

if [ "${AUGMENTED_REPORT:-false}" = "true" ]; then
  target_dir="."
  file_name="ssdlc-compliance-${VERSION}-${DATE}.md"
  SBOM_TEXT="  - See Augmented SBOM manifests (CycloneDX in JSON format):
      - This file has been provided along with this report under the name 'linux_amd64_augmented_sbom_v${VERSION}.json'
      - Please note that this file was generated on ${DATE} and may not reflect the latest security information of all third party dependencies."

else # If not augmented, generate the standard report
  target_dir="compliance/v${VERSION}"
  file_name="ssdlc-compliance-${VERSION}.md"
  SBOM_TEXT="  - See SBOM Lite manifests (CycloneDX in JSON format):
      - https://github.com/mongodb/terraform-provider-mongodbatlas/releases/download/v${VERSION}/sbom.json"
  # Ensure terraform-provider-mongodbatlas version directory exists
  mkdir -p "${target_dir}"
fi

export AUTHOR
export VERSION
export SBOM_TEXT

echo "Generating SSDLC report for Terraform Provider for MongoDB Atlas version ${VERSION}, author ${AUTHOR} and release date ${DATE}..."

envsubst < templates/releases/ssdlc-compliance.template.md \
  > "${target_dir}/${file_name}"

echo "SSDLC compliance report ready. Files in ${target_dir}/:"
ls -l "${target_dir}/"

echo "Printing the generated report:"
cat "${target_dir}/${file_name}"