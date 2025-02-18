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

# TODO: remove this after releasing TPF
if git diff --quiet -- ./internal/config/preview_provider_v2.go; then
  V2_SCHEMA_DISABLED=true
else
  V2_SCHEMA_DISABLED=false
fi

if $V2_SCHEMA_DISABLED; then
  echo "enabling Advanced Cluster V2 Schema"
  make enable-preview-provider-v2
fi
# end TODO

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
  local parent_dir="$1"
  local v2_dirs=("module_maintainer" "module_user" "migrate_cluster_to_advanced_cluster")
  
  for dir in "${v2_dirs[@]}"; do
    if [[ $parent_dir =~ $dir ]]; then
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

  PARENT_DIR=$(basename "$(dirname "$DIR")") # module_maintainer and module_user uses {PARENT_DIR}/vX/main.tf

  if is_v2_dir "$PARENT_DIR"; then
    echo "v2 schema detected for $DIR"
    MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true terraform validate
  else
    terraform validate
  fi
  popd
done

# TODO: remove this after releasing TPF
if $V2_SCHEMA_DISABLED; then
  echo "restoring Advanced Cluster V2 Schema"
  git restore ./internal/config/preview_provider_v2.go
fi
# end TODO
