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
projects=$(atlas project ls -o json)

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

    echo "Deleting projectId ${clean_string_id}"
    # This command will fail if the project has a cluster inside
    atlas project delete "${clean_project_id}" --force
done
