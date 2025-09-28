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

for DIR in $(find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u); do
  [ ! -d "$DIR" ] && continue
  pushd "$DIR"
  echo; echo -e "\e[1;35m===> Example: $DIR <===\e[0m"; echo
  TF_LOG=TRACE terraform init
  terraform fmt -check -recursive
  terraform validate

  rm -rf ".terraform"
  rm -f ".terraform.lock.hcl"
  popd
done
