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
# Usage: ./generate-doc.sh resource_name [feature_map_file]
#   resource_name: the terraform resource name. Example: search_deployment
#   feature_map_file (optional): mapping file to resolve feature folder. Default: scripts/docs/docs-feature-map.txt
#     Format blocks:
#       "Feature Path":
#       name1,
#       name2
#
#   Examples:
#     ./generate-doc.sh search_deployment
#     ./generate-doc.sh project
#     ./generate-doc.sh encryption_at_rest scripts/docs-feature-map.txt
#
# The scripts requires to install tfplugindocs and to create the resource templates in 
# templates/resources/${resource_name}.md.tmpl and 
# templates/data-sources/${resource_name}.md.tmpl
# templates/data-sources/${resource_name}s.md.tmpl

set -euo pipefail

TF_VERSION="${TF_VERSION:-"1.13.0"}" # TF version to use when running tfplugindocs. Default: 1.13.0
TEMPLATE_FOLDER_PATH="${TEMPLATE_FOLDER_PATH:-"templates"}" # PATH to the templates folder. Default: templates


# if [ -z "${resource_name}" ]; then
if [ $# -eq 0 ]; then
    echo "Error: Input param not found"
    echo "Usage: ./generate-doc.sh resource_name [feature_map_file]"
    echo "resource_name is the terraform resource and data source name."
    echo "feature_map_file (optional) maps names to feature folders (default: scripts/docs-feature-map.txt)."
    echo "Examples:"
    echo "  ./generate-doc.sh search_deployment"
    echo "  ./generate-doc.sh project"
    echo "  ./generate-doc.sh encryption_at_rest scripts/docs/docs-feature-map.txt"
    exit 1
fi

resource_name="$1"
FEATURE_MAP_FILE="${2:-scripts/docs/docs-feature-map.txt}"

# helper to trim whitespace
trim() { echo "$1" | sed -E 's/^\s+|\s+$//g'; }

# resolve feature folder from mapping file; format blocks:
#   "Feature Path":\nname1,\nname2
resolve_feature() {
    local key="$1"
    local map_file="$2"
    if [ ! -f "${map_file}" ]; then
        echo ""
        return 0
    fi
    local current_feature=""
    while IFS= read -r line || [ -n "$line" ]; do
        local trimmed
        trimmed=$(echo "$line" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')
        # skip empty or comment lines
        if [ -z "${trimmed}" ] || [[ "${trimmed}" =~ ^# ]]; then
            continue
        fi
        if [[ "${trimmed}" =~ ^\".+\"[[:space:]]*:\$ ]]; then
            # header like "Feature Path":
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
    done < "${map_file}"
    echo ""
}

# Determine destination directories based on mapping
feature_path="$(resolve_feature "${resource_name}" "${FEATURE_MAP_FILE}")"
if [ -z "${feature_path}" ] && [[ "${resource_name}" == *s ]]; then
    # try singular fallback for plural data sources
    singular_name="${resource_name%"s"}"
    feature_path="$(resolve_feature "${singular_name}" "${FEATURE_MAP_FILE}")"
fi

if [ -n "${feature_path}" ]; then
    feature_base_dir="docs/${feature_path}"
    resources_dest_dir="${feature_base_dir}/Resources"
    datasources_dest_dir="${feature_base_dir}/Data Sources"
else
    resources_dest_dir="docs/resources"
    datasources_dest_dir="docs/data-sources"
fi

# Ensure destination directories exist
mkdir -p "${resources_dest_dir}" "${datasources_dest_dir}"

if [ ! -f "${TEMPLATE_FOLDER_PATH}/resources/${resource_name}.md.tmpl" ]; then
    printf "Warning: we coudn't find the template for the %s resource. The default template templates/resources.md.tmpl will be used." "${resource_name}"
    printf "Please, make sure to include the resource template under %s.\n\n" "${TEMPLATE_FOLDER_PATH}/resources/${resource_name}.md.tmpl"
fi

if [ ! -f "${TEMPLATE_FOLDER_PATH}/data-sources/${resource_name}.md.tmpl" ]; then
    printf "Warning: we coudn't find the template for the %s data source. The default template templates/data-source.md.tmpl will be used." "${resource_name}"
    printf "Please, make sure to include the data source template under %s.\n\n" "${TEMPLATE_FOLDER_PATH}/data-sources/${resource_name}.md.tmpl"
fi

if [ ! -f "${TEMPLATE_FOLDER_PATH}/data-sources/${resource_name}s.md.tmpl" ]; then
    echo "Warning: we coudn't find the template for the ${resource_name}s data source"
    printf "Please, make sure to include the data source template under %s." "${TEMPLATE_FOLDER_PATH}/data-sources/${resource_name}.md.tmpl"
    printf "Skipping this check: We assume that the resource does not have a plural data source.\n\n"
fi

trap 'rm -R docs-out/' EXIT # temp dir cleanup when script exits

tfplugindocs generate --tf-version "${TF_VERSION}" --website-source-dir "${TEMPLATE_FOLDER_PATH}"  --rendered-website-dir "docs-out"

if [ ! -f "docs-out/resources/${resource_name}.md" ]; then
    echo "Error: We cannot find the documentation file for the resource ${resource_name}.md"
    echo "Please, make sure to include the resource template under templates/resources/${resource_name}.md.tmpl"
    printf "Skipping this step: We assume that only a data source is being generated.\n\n"
else
    printf "Moving the generated resource file %s.md to %s.\n" "${resource_name}" "${resources_dest_dir}"
    mv "docs-out/resources/${resource_name}.md" "${resources_dest_dir}/${resource_name}.md"
fi

if [ ! -f "docs-out/data-sources/${resource_name}.md" ]; then
    echo "Error: We cannot find the documentation file for the data source ${resource_name}.md"
    echo "Please, make sure to include the data source template under templates/data-sources/${resource_name}.md.tmpl"
    exit 1
else
    printf "Moving the generated data-source file %s.md to %s.\n" "${resource_name}" "${datasources_dest_dir}"
    mv "docs-out/data-sources/${resource_name}.md" "${datasources_dest_dir}/${resource_name}.md"
fi

if [ ! -f "docs-out/data-sources/${resource_name}s.md" ]; then
    echo "Warning: We cannot find the documentation file for the plural data source ${resource_name}s.md."
    echo "Please, make sure to include the data source template under templates/data-sources/${resource_name}s.md.tmpl"
    printf "Skipping this step: We assume that the resource does not have a plural data source.\n\n"
else
    printf "\nMoving the generated plural data-source file %s.md to %s.\n" "${resource_name}s" "${datasources_dest_dir}"
    mv "docs-out/data-sources/${resource_name}s.md" "${datasources_dest_dir}/${resource_name}s.md"
fi

printf "\nThe documentation for %s has been created.\n" "${resource_name}"
