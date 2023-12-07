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

delete_endpoint() {
    provider=$1
    count=$(atlas privateEndpoints "${provider}" list --projectId "${clean_project_id}" -o=go-template="{{len .}}")
    retVal=$?
    if [ $retVal -ne 0 ]; then
        count=0
    fi
    if [ "${count}" != "0" ]; then
        echo "Project ${clean_project_id} contains ${provider} endpoints, will start deleting"
        id=$(atlas privateEndpoints "${provider}" list --projectId "${clean_project_id}" -o=go-template="{{(index . 0).Id}}")
        atlas privateEndpoints "${provider}" delete "${id}" --force --projectId "${clean_project_id}"
    fi
}

projectToSkip="${PROJECT_TO_NOT_DELETE:-NONE}"

export MCLI_OPS_MANAGER_URL="${MONGODB_ATLAS_OPS_MANAGER_URL}"
export MCLI_PRIVATE_API_KEY="${MONGODB_ATLAS_PRIVATE_KEY}"
export MCLI_PUBLIC_API_KEY="${MONGODB_ATLAS_PUBLIC_KEY}"
export MCLI_ORG_ID="${MONGODB_ATLAS_ORG_ID}"
org_id="${MONGODB_ATLAS_ORG_ID}"

# Get all project Ids inside the organization
projects=$(atlas project ls --limit 500 --orgId "${org_id}"  -o json)

echo "${projects}" | jq -c '.results[].id' | while read -r id; do
    # Trim the quotes from the id
    clean_project_id=$(echo "$id" | tr -d '"')
    if [ "${clean_project_id}" = "${projectToSkip}" ]; then
        echo "Skipping project with ID ${projectToSkip}"
        continue
    fi

    clusters=$(atlas cluster ls --projectId "${clean_project_id}" -o=go-template="{{.TotalCount}}")
    if [ "${clusters}" != "0" ]; then
        echo "Project ${clean_project_id} contains clusters. Skipping..."
        continue
    fi

    set +e
    delete_endpoint "aws"
    delete_endpoint "gcp"
    delete_endpoint "azure"
    set -e

    echo "Deleting projectId ${clean_project_id}"
    # This command can fail if project has a cluster, a private endpoint, or general failure. The echo command always succeeds so the subshell will succeed and continue
    (
        atlas project delete "${clean_project_id}" --force || \
        echo "Failed to delete project with ID ${clean_project_id}"
    )
done
