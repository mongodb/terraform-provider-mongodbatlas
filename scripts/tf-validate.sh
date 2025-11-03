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

# Use local provider to validate examples
go build -o bin-examples/terraform-provider-mongodbatlas .

# Two TF CLI configs: one with local override, one without
TF_CLI_CONFIG_FILE_WITH="$PWD/bin-examples/tf-validate.local.tfrc"
TF_CLI_CONFIG_FILE_NO="$PWD/bin-examples/tf-validate.remote.tfrc"

cat << EOF > "$TF_CLI_CONFIG_FILE_WITH"
provider_installation {
  dev_overrides {
    "mongodb/mongodbatlas" = "$PWD/bin-examples"
  }
  direct {}
}
EOF

cat << EOF > "$TF_CLI_CONFIG_FILE_NO"
provider_installation {
  direct {}
}
EOF

for DIR in $(find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u); do
  [ ! -d "$DIR" ] && continue
  pushd "$DIR"
  # For directories named like v1.x.x, do NOT use local override
  if [[ "$(basename "$DIR")" == "v1.x.x" ]]; then
    export TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE_NO"
  else
    export TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE_WITH"
  fi
  echo; echo -e "\e[1;35m===> Example: $DIR <===\e[0m"; echo
  terraform init > /dev/null # suppress output as it's very verbose
  terraform fmt -check -recursive
  terraform validate

  rm -rf ".terraform"
  rm -rf ".terraform.lock.hcl"
  popd
done
