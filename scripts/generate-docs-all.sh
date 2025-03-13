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

#
# Shell script to generate the Terraform documentation for the resource and data sources.
#
# Usage: ./generate-docs-all.sh"
#
# The scripts requires to install tfplugindocs and to create the resource templates in 
# templates/resources/${resource_name}.md.tmpl and 
# templates/data-sources/${resource_name}.md.tmpl
# templates/data-sources/${resource_name}s.md.tmpl

set -euo pipefail

TF_VERSION="${TF_VERSION:-"1.11.2"}" # TF version to use when running tfplugindocs. Default: 1.11.2
TEMPLATE_FOLDER_PATH="${TEMPLATE_FOLDER_PATH:-"templates"}" # PATH to the templates folder. Default: templates

# ensure preview resource and data sources are also included during generation
export MONGODB_ATLAS_ENABLE_PREVIEW="true" 

trap 'rm -R docs-out/' EXIT # temp dir cleanup when script exits

tfplugindocs generate --tf-version "${TF_VERSION}" --website-source-dir "${TEMPLATE_FOLDER_PATH}"  --rendered-website-dir "docs-out" --provider-name "mongodbatlas"

printf "\nStarting file move\n\n"

for file in "$TEMPLATE_FOLDER_PATH/resources"/*; do
    filenameTemplate=$(basename -- "$file")
    filename="${filenameTemplate%.*}"     
    if [ -f "docs-out/resources/${filename}" ]; then
        printf "Resource file moved: %s\n" "${filename}"
        mv "docs-out/resources/${filename}" "docs/resources/${filename}"
    fi
done

for file in "$TEMPLATE_FOLDER_PATH/data-sources"/*; do
    filenameTemplate=$(basename -- "$file")
    filename="${filenameTemplate%.*}"     
    if [ -f "docs-out/data-sources/${filename}" ]; then
        printf "Data source file moved: %s\n" "${filename}"
        mv "docs-out/data-sources/${filename}" "docs/data-sources/${filename}"
    fi
done
