#!/usr/bin/env bash

# Copyright 2021 MongoDB Inc
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

arch_name=$(uname -m)

for DIR in $(find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u); do
  [ ! -d "$DIR" ] && continue
  
  
  # Skip directories with "v08" or "v09" in their name for ARM64
  if [[ "$arch_name" == "arm64" ]] && echo "$DIR" | grep -qE "v08|v09"; then
      echo "Skip directories with \"v08\" or \"v09\" in their name for ARM64"
      echo "TF provider does not have a package available for ARM64 for version < 1.0"
      echo "Skipping directory: $DIR"
      continue
  fi

  pushd "$DIR"

  echo; echo -e "\e[1;35m===> Initializing Example: $DIR <===\e[0m"; echo
  terraform init
  
  echo; echo -e "\e[1;35m===> Format Checking Example: $DIR <===\e[0m"; echo
  terraform fmt -check

  echo; echo -e "\e[1;35m===> Validating Example: $DIR <===\e[0m"; echo
  terraform validate
  
  echo; echo -e "\e[1;35m===> Validating Syntax Example: $DIR <===\e[0m"; echo
  # Terraform syntax checks
  tflint \
    --enable-rule=terraform_deprecated_interpolation \
    --enable-rule=terraform_deprecated_index \
    --enable-rule=terraform_unused_declarations \
    --enable-rule=terraform_comment_syntax \
    --enable-rule=terraform_required_version \
    --minimum-failure-severity=warning
  popd
done
