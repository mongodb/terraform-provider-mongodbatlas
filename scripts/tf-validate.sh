#!/usr/bin/env bash

# Copyright 2023 MongoDB Inc
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

set -Eeou pipefail

# Delete Terraform execution files so the script can be run multiple times
find ./examples -type d -name ".terraform" -exec rm -rf {} +
find ./examples -type f -name ".terraform.lock.hcl" -exec rm -f {} +

export TF_CLI_CONFIG_FILE="$PWD/examples-bin/tf-validate.tfrc"

export TF_PLUGIN_CACHE_DIR="$PWD/examples-cache"
rm -rf "$TF_PLUGIN_CACHE_DIR"
mkdir -p "$TF_PLUGIN_CACHE_DIR"

# Use local provider to validate examples
go build -o examples-bin/terraform-provider-mongodbatlas .

cat << EOF > "$TF_CLI_CONFIG_FILE"
provider_installation { 
  dev_overrides {
    "mongodb/mongodbatlas" = "$PWD/examples-bin"
  }
  direct {} 
}
EOF

# Function to validate a single directory
validate_dir() {
  local dir=$1
  local tempfile=$(mktemp)

  # Capture all output to a temp file to keep it grouped
  {
    [ ! -d "$dir" ] && return 0
    cd "$dir"
    echo
    echo -e "\e[1;35m===> Example: $dir <===\e[0m"
    echo
    TF_LOG=TRACE terraform init 2>&1
    terraform fmt -check -recursive 2>&1
    terraform validate 2>&1
  } &> "$tempfile"

  # Output the grouped logs
  cat "$tempfile"
  local exit_code=$?
  rm -f "$tempfile"
  return $exit_code
}

export -f validate_dir
export TF_CLI_CONFIG_FILE
export TF_PLUGIN_CACHE_DIR

# Find all directories and run validation in parallel with 10 workers
find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u | \
  xargs -P 10 -I {} bash -c 'validate_dir "$@"' _ {}
