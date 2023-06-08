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

for DIR in $(find ./examples -type f -name '*.tf' -exec dirname {} \; | sort -u); do
  [ ! -d "$DIR" ] && continue
  pushd "$DIR"
  echo; echo -e "\e[1;35m===> Initializing Example: $DIR <===\e[0m"; echo
  terraform init
  
  echo; echo -e "\e[1;35m===> Format Checking Example: $DIR <===\e[0m"; echo
  terraform fmt -check
  # Terraform syntax checks
  echo; echo -e "\e[1;35m===> Validating Example: $DIR <===\e[0m"; echo
   tflint \
     --enable-rule=terraform_deprecated_interpolation \
     --enable-rule=terraform_deprecated_index \
     --enable-rule=terraform_unused_declarations \
     --enable-rule=terraform_comment_syntax \
     --enable-rule=terraform_required_version
  popd
done
