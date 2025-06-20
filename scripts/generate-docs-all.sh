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

TF_VERSION="${TF_VERSION:-"1.12.2"}" # TF version to use when running tfplugindocs. Default: 1.12.2
TEMPLATE_FOLDER_PATH="${TEMPLATE_FOLDER_PATH:-"templates"}" # PATH to the templates folder. Default: templates

trap 'rm -R docs-out/' EXIT # temp dir cleanup when script exits

tfplugindocs generate --tf-version "${TF_VERSION}" --website-source-dir "${TEMPLATE_FOLDER_PATH}"  --rendered-website-dir "docs-out" --provider-name "mongodbatlas"

printf "\nStarting file move\n\n"

# moving generated docs for resource and data sources that define custom docs templates
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

# moving generated docs for autogenerated resources and data sources which use default template (ending with _api or _api_v*)
for file in docs-out/resources/*; do
    if [ -f "$file" ]; then
        filename=$(basename -- "$file")
        if [[ "$filename" == *_api.md || "$filename" == *_api_v*.md ]]; then
            printf "Autogenerated resource file moved: %s\n" "${filename}"
            mv "$file" "docs/resources/${filename}"
        fi
    fi
done

for file in docs-out/data-sources/*; do
    if [ -f "$file" ]; then
        filename=$(basename -- "$file")
        if [[ "$filename" == *_api.md || "$filename" == *_api_v*.md ]]; then
            printf "Autogenerated data source file moved: %s\n" "${filename}"
            mv "$file" "docs/data-sources/${filename}"
        fi
    fi
done
