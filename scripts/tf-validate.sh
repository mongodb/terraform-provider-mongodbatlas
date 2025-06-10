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

export TF_CLI_CONFIG_FILE="$PWD/bin-examples/tf-validate.tfrc"

# Use local provider to validate examples
go build -o bin-examples/terraform-provider-mongodbatlas .

cat << EOF > "$TF_CLI_CONFIG_FILE"
provider_installation { 
  dev_overrides {
    "mongodb/mongodbatlas" = "$PWD/bin-examples"
  }
  direct {} 
}
EOF

# Function to check if directory is a V2 schema directory
is_v2_dir() {
  local parent_dir
  local grand_parent_dir
  parent_dir=$(basename "$1")
  grand_parent_dir=$(basename "$(dirname "$1")")
  local v2_parent_dirs=("mongodbatlas_backup_compliance_policy")
  local v2_grand_parent_dirs=("module_maintainer" "module_user" "migrate_cluster_to_advanced_cluster") # module_maintainer and module_user uses {PARENT_DIR}/vX/main.tf
  
  for dir in "${v2_parent_dirs[@]}"; do
    if [[ $parent_dir =~ $dir ]]; then
      return 0  # True
    fi
  done
  for dir in "${v2_grand_parent_dirs[@]}"; do
    if [[ $grand_parent_dir =~ $dir ]]; then
      return 0  # True
    fi
  done
  return 1  # False
}

for DIR in $(find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u); do
  [ ! -d "$DIR" ] && continue
  pushd "$DIR"
  echo; echo -e "\e[1;35m===> Example: $DIR <===\e[0m"; echo
  terraform init > /dev/null # suppress output as it's very verbose
  terraform fmt -check -recursive

  if is_v2_dir "$DIR"; then
    echo "v2 schema detected for $DIR"
    MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true terraform validate
  else
    echo "v1 schema detected for $DIR"
    terraform validate
  fi
  popd
done
