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

usage() {
	echo "Usage: $0 [true|false]"
	echo "  true:  Returns details of supported Terraform versions"
	echo "  false: Returns only supported Terraform version numbers"
	exit 1
}

fetch_terraform_releases_page() {
	local page="$1"
	local api_version="2022-11-28"
	curl -s \
		--request GET \
		--url "https://api.github.com/repos/hashicorp/terraform/releases?per_page=100&page=$page" \
		--header "Authorization: Bearer $GITHUB_TOKEN" \
		--header "X-GitHub-Api-Version: $api_version"
}

get_last_day_of_month() {
	last_day_of_month=0
	case $1 in
	01 | 03 | 05 | 07 | 08 | 10 | 12)
		last_day_of_month=31
		;;
	04 | 06 | 09 | 11)
		last_day_of_month=30
		;;
	02)
		# February: check if it's a leap year
		if ((year % 4 == 0 && (year % 100 != 0 || year % 400 == 0))); then
			last_day_of_month=29
		else
			last_day_of_month=28
		fi
		;;
	esac
	echo $last_day_of_month
}

add_end_support_date() {
	new_json_list=$(echo "$1" | jq -c '.[]' | while IFS= read -r obj; do
		input_date=$(echo "$obj" | jq -r '.published_at')

		# Extract the year, month, day, hour, minute, and second from the input date
		year="${input_date:0:4}"
		month="${input_date:5:2}"
		hour="${input_date:11:2}"
		minute="${input_date:14:2}"
		second="${input_date:17:2}"
		last_day_of_month=$(get_last_day_of_month "$month")
		new_year=$((year + 2))

		new_date="${new_year}-${month}-${last_day_of_month}T${hour}:${minute}:${second}Z"

		echo "$obj" | jq --arg new_date "${new_date}" '.end_support_date = $new_date'
	done | jq -s '.')

	echo "$new_json_list"
}

get_terraform_supported_versions_details() {
	page=1

	while true; do
		response=$(fetch_terraform_releases_page "$page")
		if [[ "$(printf '%s\n' "$response" | jq -e 'length == 0')" == "true" ]]; then
			break
		else
			versions=$(echo "$response" | jq -r '.[] | {version: .tag_name, published_at: .published_at}')
			filtered_versions_json=$(printf '%s\n' "${versions}" | jq -s '.')
			updated_date_versions=$(add_end_support_date "$filtered_versions_json")
			filtered_out_deprecated_versions=$(echo "$updated_date_versions" | jq 'map(select((.version | test("alpha|beta|rc"; "i") | not) and ((.end_support_date | fromdate) >= now))) | map(select(.version | endswith(".0")))')
			# echo "filtered_out_deprecated_versions: ${filtered_out_deprecated_versions}"
			if [ -z ${json_array+x} ]; then
				# If json_array is empty, assign filtered_versions directly
				json_array=$(jq -n --argjson filtered_versions "${filtered_out_deprecated_versions}" '$filtered_versions')
			else
				# If json_array is not empty, append filtered_versions to it
				json_array=$(echo "$json_array" | jq --argjson filtered_versions "${filtered_out_deprecated_versions}" '. + $filtered_versions')
			fi

			((page++))
		fi
	done

	echo "$json_array"
}

get_terraform_supported_versions() {
	json_array=$(get_terraform_supported_versions_details)

	versions_array=$(printf '%s\n' "${json_array}" | jq -r '.[] | .version')

	formatted_output=$(echo "$versions_array" | awk 'BEGIN { ORS = "" } {gsub(/^v/,"",$1); gsub(/\.0$/,".x",$1); printf("%s\"%s\"", (NR==1?"":", "), $1)}' | sed 's/,"/," /g')

	echo "[$formatted_output]" | jq -c .
}

get_latest_terraform_version() {
	api_version="2022-11-28"
	latest_version=""

	for page in 1 2; do
		response=$(fetch_terraform_releases_page "$page")

		if [[ "$(printf '%s\n' "$response" | jq -e 'length == 0')" == "true" ]]; then
			break
		fi

		versions=$(echo "$response" | jq -r '.[] | select(.tag_name | test("alpha|beta|rc"; "i") | not) | .tag_name')

		for version in $versions; do
			version_cleaned="${version#v}"
			if [[ -z "$latest_version" || "$(echo -e "$latest_version\n$version_cleaned" | sort -V | tail -n 1)" == "$version_cleaned" ]]; then
				latest_version="$version_cleaned"
			fi
		done
	done

	echo "$latest_version"
}

if [ $# -ne 1 ]; then
	usage
fi

get_details=$1
if [ "$get_details" = "true" ]; then
	get_terraform_supported_versions_details
elif [ "$get_details" = "false" ]; then
	get_terraform_supported_versions
elif [ "$get_details" = "latest" ]; then
	get_latest_terraform_version
else
	echo "Invalid parameter."
	usage
fi
