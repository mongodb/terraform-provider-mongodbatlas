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

TOOL_VERSIONS_FILE=".tool-versions"

LATEST_TF_VERSION=$(echo "$TF_VERSIONS_ARR" | sed -E 's/^\["([^"]+).*/\1/')

TF_VERSION="${LATEST_TF_VERSION//x/0}"

# Update Terraform versions in test-suite.yml
sed -i.bak -E "/^ *terraform_matrix:/,/^ *provider_matrix:/ s|(default: ')[^']*(')|\1$TF_VERSIONS_ARR\2|" "$TEST_SUITE_YAML_FILE"

sed -i.bak -E "s|schedule_terraform_matrix: '.*'|schedule_terraform_matrix: '[\"$LATEST_TF_VERSION\"]'|" "$TEST_SUITE_YAML_FILE"

# Update Terraform versions in .tool-versions
sed -i.bak -E "s|^(terraform) [0-9]+\.[0-9]+\.[0-9]+|\1 $TF_VERSION|" "$TOOL_VERSIONS_FILE"
