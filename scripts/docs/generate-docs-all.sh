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
# Usage: ./generate-docs-all.sh
#   This script automatically uses scripts/docs/docs-feature-map.txt if present to map
#   doc base names to feature paths (format blocks). If the map
#   file is missing or a name isn't found, files fall back to docs/resources and
#   docs/data-sources.
#
# The scripts requires to install tfplugindocs and to create the resource templates in 
# templates/resources/${resource_name}.md.tmpl and 
# templates/data-sources/${resource_name}.md.tmpl
# templates/data-sources/${resource_name}s.md.tmpl

set -euo pipefail

TF_VERSION="${TF_VERSION:-"1.13.0"}" # TF version to use when running tfplugindocs. Default: 1.13.0
TEMPLATE_FOLDER_PATH="${TEMPLATE_FOLDER_PATH:-"templates"}" # PATH to the templates folder. Default: templates
FEATURE_MAP_FILE="scripts/docs/docs-feature-map.txt"

trap 'rm -R docs-out/' EXIT # temp dir cleanup when script exits

tfplugindocs generate --tf-version "${TF_VERSION}" --website-source-dir "${TEMPLATE_FOLDER_PATH}"  --rendered-website-dir "docs-out" --provider-name "mongodbatlas"

printf "\nStarting file move\n\n"

if [ -f "${FEATURE_MAP_FILE}" ]; then
    printf "Using feature map file: %s\n\n" "${FEATURE_MAP_FILE}"
else
    printf "Feature map file not found at %s. Falling back to legacy directories.\n\n" "${FEATURE_MAP_FILE}"
fi

# helper to trim whitespace
trim() { echo "$1" | sed -E 's/^\s+|\s+$//g'; }

# helper to resolve feature path for a given base name using mapping file format blocks:
#   "Feature Path":\nname1,\nname2
feature_for() {
    local key="$1"
    if [ -n "${FEATURE_MAP_FILE}" ] && [ -f "${FEATURE_MAP_FILE}" ]; then
        local current_feature=""
        while IFS= read -r line || [ -n "$line" ]; do
            local trimmed
            trimmed=$(echo "$line" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')
            if [ -z "${trimmed}" ] || [[ "${trimmed}" =~ ^# ]]; then
                continue
            fi
            if [[ "${trimmed}" =~ ^\".+\"[[:space:]]*:\$ ]]; then
                current_feature=${trimmed%:}
                current_feature=${current_feature%\"}
                current_feature=${current_feature#\"}
                continue
            fi
            if [ -n "${current_feature}" ]; then
                IFS=',' read -r -a items <<< "${trimmed}"
                for item in "${items[@]}"; do
                    local name
                    name=$(echo "$item" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')
                    if [ "${name}" = "${key}" ]; then
                        echo "${current_feature}"
                        return 0
                    fi
                done
            fi
        done < "${FEATURE_MAP_FILE}"
    fi
    echo ""
}

# ensure destination dirs exist
ensure_dirs() {
    local feature="$1"
    if [ -n "${feature}" ]; then
        mkdir -p "docs/${feature}/Resources" "docs/${feature}/Data Sources"
    else
        mkdir -p "docs/resources" "docs/data-sources"
    fi
}

move_with_mapping() {
    local src_path="$1"   # full path to docs-out file
    local is_resource="$2" # "resource" or "data-source"
    local filename
    filename=$(basename -- "$src_path")
    local base="${filename%.md}"
    local feature
    feature=$(feature_for "${base}")
    # try singular for plurals when mapping missing and it's a data source
    if [ -z "${feature}" ] && [ "${is_resource}" = "data-source" ] && [[ "${base}" == *s ]]; then
        local singular="${base%"s"}"
        feature=$(feature_for "${singular}")
    fi
    ensure_dirs "${feature}"
    if [ -n "${feature}" ]; then
        if [ "${is_resource}" = "resource" ]; then
            printf "%s file moved to feature: %s -> %s\n" "Resource" "${filename}" "docs/${feature}/Resources/${filename}"
            mv "$src_path" "docs/${feature}/Resources/${filename}"
        else
            printf "%s file moved to feature: %s -> %s\n" "Data source" "${filename}" "docs/${feature}/Data Sources/${filename}"
            mv "$src_path" "docs/${feature}/Data Sources/${filename}"
        fi
    else
        if [ "${is_resource}" = "resource" ]; then
            printf "Resource file moved: %s\n" "${filename}"
            mv "$src_path" "docs/resources/${filename}"
        else
            printf "Data source file moved: %s\n" "${filename}"
            mv "$src_path" "docs/data-sources/${filename}"
        fi
    fi
}

# moving generated docs for resource and data sources that define custom docs templates
for file in "$TEMPLATE_FOLDER_PATH/resources"/*; do
    filenameTemplate=$(basename -- "$file")
    filename="${filenameTemplate%.*}"     
    if [ -f "docs-out/resources/${filename}.md" ]; then
        move_with_mapping "docs-out/resources/${filename}.md" "resource"
    fi
done

for file in "$TEMPLATE_FOLDER_PATH/data-sources"/*; do
    filenameTemplate=$(basename -- "$file")
    filename="${filenameTemplate%.*}"     
    if [ -f "docs-out/data-sources/${filename}.md" ]; then
        move_with_mapping "docs-out/data-sources/${filename}.md" "data-source"
    fi
done

# moving generated docs for autogenerated resources and data sources which use default template (ending with _api or _api_v*)
for file in docs-out/resources/*; do
    if [ -f "$file" ]; then
        filename=$(basename -- "$file")
        if [[ "$filename" == *_api.md || "$filename" == *_api_v*.md ]]; then
            move_with_mapping "$file" "resource"
        fi
    fi
done

for file in docs-out/data-sources/*; do
    if [ -f "$file" ]; then
        filename=$(basename -- "$file")
        if [[ "$filename" == *_api.md || "$filename" == *_api_v*.md ]]; then
            move_with_mapping "$file" "data-source"
        fi
    fi
done
