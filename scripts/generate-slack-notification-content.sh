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
set -euo pipefail

# Check if the parameter is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <success|failure>"
    exit 1
fi

# Check if the parameter is either "success" or "failure"
if [ "$1" != "success" ] && [ "$1" != "failure" ]; then
    echo "Invalid parameter. Please provide either 'success' or 'failure'."
    exit 1
fi

oncall_tag=""
if [ "$1" == "success" ]; then
    text_value=":white_check_mark: HashiCorp Terraform Compatibility Matrix succeeded"
	action_text="Successful action"
else
    text_value=":red_circle: HashiCorp Terraform Compatibility Matrix failed"
	action_text="Failed action"
	oncall_tag="$2"
fi

server_url=$3
repository=$4
run_id=$5

json="{
        \"text\": \"$text_value\",
        \"blocks\": [
            {
                \"type\": \"section\",
                \"text\": {
                    \"type\": \"mrkdwn\",
                    \"text\": \"*$text_value* $oncall_tag\"
                }
            },
            {
                \"type\": \"actions\",
                \"elements\": [
                    {
                        \"type\": \"button\",
                        \"text\": {
                            \"type\": \"plain_text\",
                            \"text\": \":github: $action_text\"
                        },
                        \"url\": \"${server_url}/${repository}/actions/runs/${run_id}\"
                    }
                ]
            }
        ]
    }"

echo "$json" | jq -c .
