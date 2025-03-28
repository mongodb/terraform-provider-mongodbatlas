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
# Usage: ./generate-doc.sh" ${resource_name}
#   resource_name is the terraform resource name. Example: search_deployment
#   echo "Examples:"
#   echo "  search_deployment"
#   echo "  project"
#   echo "  online_archive"
#   echo "  encryption_at_rest"
#
# The scripts requires to install tfplugindocs and to create the resource templates in 
# templates/resources/${resource_name}.md.tmpl and 
# templates/data-sources/${resource_name}.md.tmpl
# templates/data-sources/${resource_name}s.md.tmpl

set -euo pipefail

TF_VERSION="${TF_VERSION:-"1.12.5"}" # TF version to use when running tfplugindocs. Default: 1.12.5
TEMPLATE_FOLDER_PATH="${TEMPLATE_FOLDER_PATH:-"templates"}" # PATH to the templates folder. Default: templates


# if [ -z "${resource_name}" ]; then
if [ $# -eq 0 ]; then
    echo "Error: Input param not found"
    echo "Usage: ./generate-doc.sh resource_name"
    echo "resource_name is the terraform resource and data source name."
    echo "Examples:"
    echo "  search_deployment"
    echo "  project"
    echo "  online_archive"
    echo "  encryption_at_rest"
    exit 1
fi

resource_name="$1"

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

# ensure preview resource and data sources are also included during generation
export MONGODB_ATLAS_ENABLE_PREVIEW="true" 

trap 'rm -R docs-out/' EXIT # temp dir cleanup when script exits

tfplugindocs generate --tf-version "${TF_VERSION}" --website-source-dir "${TEMPLATE_FOLDER_PATH}"  --rendered-website-dir "docs-out"

if [ ! -f "docs-out/resources/${resource_name}.md" ]; then
    echo "Error: We cannot find the documentation file for the resource ${resource_name}.md"
    echo "Please, make sure to include the resource template under templates/resources/${resource_name}.md.tmpl"
    printf "Skipping this step: We assume that only a data source is being generated.\n\n"
else
    printf "Moving the generated resource file %s.md to the website folder.\n" "${resource_name}"
    mv "docs-out/resources/${resource_name}.md" "docs/resources/${resource_name}.md"
fi

if [ ! -f "docs-out/data-sources/${resource_name}.md" ]; then
    echo "Error: We cannot find the documentation file for the data source ${resource_name}.md"
    echo "Please, make sure to include the data source template under templates/data-sources/${resource_name}.md.tmpl"
    exit 1
else
    printf "Moving the generated data-source file %s.md to the website folder.\n" "${resource_name}"
    mv "docs-out/data-sources/${resource_name}.md" "docs/data-sources/${resource_name}.md"
fi

if [ ! -f "docs-out/data-sources/${resource_name}s.md" ]; then
    echo "Warning: We cannot find the documentation file for the plural data source ${resource_name}s.md."
    echo "Please, make sure to include the data source template under templates/data-sources/${resource_name}s.md.tmpl"
    printf "Skipping this step: We assume that the resource does not have a plural data source.\n\n"
else
    printf "\nMoving the generated plural data-source file %s.md to the website folder.\n" "${resource_name}s"
    mv "docs-out/data-sources/${resource_name}s.md" "docs/data-sources/${resource_name}s.md"
fi

printf "\nThe documentation for %s has been created.\n" "${resource_name}"
