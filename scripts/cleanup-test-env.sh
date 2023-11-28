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

projectToSkip="${PROJECT_TO_NOT_DELETE:-NONE}"

# Get all project Ids inside the organization
projects=$(atlas project ls --limit 500 -o json)

echo "${projects}" | jq -c '.results[].id' | while read -r id; do
    # Trim the quotes from the id
    clean_project_id=$(echo "$id" | tr -d '"')
    if [ "${clean_project_id}" = "${projectToSkip}" ]; then
        echo "Skipping project with ID ${projectToSkip}"
        continue
    fi

    countAWS=$(atlas privateEndpoints aws list --projectId "${clean_project_id}" -o=go-template="{{len .}}")
    if [ "${countAWS}" != "0" ]; then
        echo "Project ${clean_project_id} contains AWS endpoints, will start deleting it now and will try to delete the project in the next execution"
        idAWS=$(atlas privateEndpoints aws list --projectId "${clean_project_id}" -o=go-template="{{(index . 0).Id}}")
        atlas privateEndpoints aws delete "${idAWS}" --force --projectId "${clean_project_id}" || \
        echo "Failed to delete AWS private endpoint with project ID ${clean_project_id}, endpoint ID: ${idAWS}"
        continue
    fi

    countGCP=$(atlas privateEndpoints gcp list --projectId "${clean_project_id}" -o=go-template="{{len .}}")
    if [ "${countGCP}" != "0" ]; then
        echo "Project ${clean_project_id} contains GCP endpoints, will start deleting it now and will try to delete the project in the next execution"
        idGCP=$(atlas privateEndpoints gcp list --projectId "${clean_project_id}" -o=go-template="{{(index . 0).Id}}")
        atlas privateEndpoints gcp delete "${idGCP}" --force --projectId "${clean_project_id}" || \
        echo "Failed to delete GCP private endpoint with project ID ${clean_project_id}, endpoint ID: ${idGCP}"
        continue
    fi

    countAzure=$(atlas privateEndpoints azure list --projectId "${clean_project_id}" -o=go-template="{{len .}}")
    if [ "${countAzure}" != "0" ]; then
        echo "Project ${clean_project_id} contains Azure endpoints, will start deleting it now and will try to delete the project in the next execution"
        idAzure=$(atlas privateEndpoints azure list --projectId "${clean_project_id}" -o=go-template="{{(index . 0).Id}}")
        atlas privateEndpoints azure delete "${idAzure}" --force --projectId "${clean_project_id}" || \
        echo "Failed to delete Azure private endpoint with project ID ${clean_project_id}, endpoint ID: ${idAzure}"
        continue
    fi

    clusters=$(atlas cluster ls --projectId "${clean_project_id}" -o=go-template="{{.TotalCount}}")
    if [ "${clusters}" != "0" ]; then
        echo "Project ${clean_project_id} contains clusters. Skipping..."
        continue
    fi

    echo "Deleting projectId ${clean_project_id}"
    # This command can fail if project has a cluster, a private endpoint, or general failure. The echo command always succeeds so the subshell will succeed and continue
    (
        atlas project delete "${clean_project_id}" --force || \
        echo "Failed to delete project with ID ${clean_project_id}"
    )
done
