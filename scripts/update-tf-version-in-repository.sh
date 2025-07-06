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
DEV_BRANCH_YAML_FILE=".github/workflows/update-dev-branches.yml"
EXAMPLES_YAML_FILE=".github/workflows/examples.yml"
TOOL_VERSIONS_FILE=".tool-versions"
DOC_SCRIPT="scripts/generate-doc.sh"
DOC_ALL_SCRIPT="scripts/generate-docs-all.sh"

LATEST_TF_VERSION=$(echo "$TF_VERSIONS_ARR" | sed -E 's/^\["([^"]+).*/\1/')

# Update Terraform versions in test-suite.yml
sed -i.bak -E "/^ *terraform_matrix:/,/^ *provider_matrix:/ s|(default: ')[^']*(')|\1$TF_VERSIONS_ARR\2|" "$TEST_SUITE_YAML_FILE"
sed -i.bak -E "s|schedule_terraform_matrix: '.*'|schedule_terraform_matrix: '[\"$LATEST_TF_VERSION\"]'|" "$TEST_SUITE_YAML_FILE"

# Update Terraform version in examples.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$EXAMPLES_YAML_FILE"

# Update Terraform version in code-health.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$CODE_HEALTH_YAML_FILE"

# Update Terraform version in update-dev-branches.yml
sed -i.bak -E "s|terraform_version: '.*'|terraform_version: '$LATEST_TF_VERSION'|" "$DEV_BRANCH_YAML_FILE"

# Update Terraform version in acceptance-tests.yml
sed -i.bak -E "s~terraform_version \|\| '[0-9]+\.[0-9]+\.x'~terraform_version \|\| '$LATEST_TF_VERSION'~" "$ACCEPTANCE_TESTS_YAML_FILE"

# Update patch version occurrences
LATEST_TF_PATCH_VERSION=$(./scripts/get-terraform-supported-versions.sh "latest")

# Update Terraform versions in .tool-versions
sed -i.bak -E "s|^(terraform) [0-9]+\.[0-9]+\.[0-9]+|\1 $LATEST_TF_PATCH_VERSION|" "$TOOL_VERSIONS_FILE"

# Update Terraform versions in generate-doc scripts
sed -i.bak -E "/TF_VERSION=/ s/[0-9]+\.[0-9]+\.[0-9]+/$LATEST_TF_PATCH_VERSION/g" "$DOC_SCRIPT"
sed -i.bak -E "/TF_VERSION=/ s/[0-9]+\.[0-9]+\.[0-9]+/$LATEST_TF_PATCH_VERSION/g" "$DOC_ALL_SCRIPT"
