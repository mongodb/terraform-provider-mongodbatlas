#!/usr/bin/env bash

# Copyright 2024 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
set -euo pipefail

# example value of TF_VERSIONS_ARR='["1.9.x", "1.8.x", "1.7.x"]'
TF_VERSIONS_ARR=$(./scripts/get-terraform-supported-versions.sh "false")

TEST_SUITE_YAML_FILE=".github/workflows/test-suite.yml"
ACCEPTANCE_TESTS_YAML_FILE=".github/workflows/acceptance-tests.yml"
CODE_HEALTH_YAML_FILE=".github/workflows/code-health.yml"
EXAMPLES_YAML_FILE=".github/workflows/examples.yml"
UPDATE_DEV_BRANCHES_YAML_FILE=".github/workflows/update-dev-branches.yml"
TOOL_VERSIONS_FILE=".tool-versions"

LATEST_TF_VERSION=$(echo "$TF_VERSIONS_ARR" | sed -E 's/^\["([^"]+).*/\1/')

# Update Terraform versions in test-suite.yml
sed -i.bak -E "/^ *terraform_matrix:/,/^ *provider_matrix:/ s|(default: ')[^']*(')|\1$TF_VERSIONS_ARR\2|" "$TEST_SUITE_YAML_FILE"
sed -i.bak -E "s|schedule_terraform_matrix: '.*'|schedule_terraform_matrix: '[\"$LATEST_TF_VERSION\"]'|" "$TEST_SUITE_YAML_FILE"

# Update Terraform version in examples.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$EXAMPLES_YAML_FILE"

# Update Terraform version in code-health.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$CODE_HEALTH_YAML_FILE"

# Update Terraform version in update-dev-branches.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$UPDATE_DEV_BRANCHES_YAML_FILE"

# Update Terraform version in acceptance-tests.yml
sed -i.bak -E "s~terraform_version \|\| '[0-9]+\.[0-9]+\.x'~terraform_version \|\| '$LATEST_TF_VERSION'~" "$ACCEPTANCE_TESTS_YAML_FILE"

# Update patch version occurrences
LATEST_TF_PATCH_VERSION=$(./scripts/get-terraform-supported-versions.sh "latest")

# Update Terraform versions in .tool-versions
sed -i.bak -E "s|^(terraform) [0-9]+\.[0-9]+\.[0-9]+|\1 $LATEST_TF_PATCH_VERSION|" "$TOOL_VERSIONS_FILE"

MIN_TF_VERSION=$(echo "$TF_VERSIONS_ARR" | jq -r 'last' | sed 's/\.x$//')

# Only bump TF version if min supported version is greater than the one set in the file.
# Skips files that require a higher version for specific features (e.g. write-only attributes).
version_lt() {
	local major1="${1%%.*}" minor1="${1#*.}"
	local major2="${2%%.*}" minor2="${2#*.}"
	[ "$major1" -lt "$major2" ] && return 0
	[ "$major1" -eq "$major2" ] && [ "$minor1" -lt "$minor2" ] && return 0
	return 1
}

# Update required_version field in examples versions.tf files.
find examples -name "versions.tf" -not -path "*/.terraform/*" | while read -r file; do
	current=$(grep -oE 'required_version = ">= [0-9]+\.[0-9]+"' "$file" | grep -oE '[0-9]+\.[0-9]+' | head -1 || true)
	if [ -n "$current" ] && version_lt "$current" "$MIN_TF_VERSION"; then
		sed -i.bak -E "s/required_version = \">= [0-9]+\.[0-9]+\"/required_version = \">= $MIN_TF_VERSION\"/" "$file"
	fi
done

