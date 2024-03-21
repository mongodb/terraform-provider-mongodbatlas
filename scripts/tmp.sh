#!/bin/bash

input_array='[{
  "version": "v1.7.0",
  "published_at": "2024-01-17T19:59:32Z",
  "end_support_date": "2026-01-31T19:59:32Z"
},
{
  "version": "v1.6.0",
  "published_at": "2023-10-04T17:11:35Z",
  "end_support_date": "2025-10-31T17:11:35Z"
}]'

indexFile="../website/docs/index.html.markdown"

transform_array() {
    local arr="$1"
    local updated_arr="["
    local isFirstElement=true

    for ((i=0; i<$(jq length <<< "$arr"); i++)); do
        version=$(jq -r ".[$i].version" <<< "$arr" | sed 's/^v//;s/\.0/.x/')
        published_at=$(jq -r ".[$i].published_at" <<< "$arr" | cut -dT -f1)
        end_support_date=$(jq -r ".[$i].end_support_date" <<< "$arr" | cut -dT -f1)

        if [ "$isFirstElement" = false ]; then
            updated_arr+=","
        fi
        updated_arr+="{\"version\": \"$version\", \"published_at\": \"$published_at\", \"end_support_date\": \"$end_support_date\"}"
        isFirstElement=false
    done

    updated_arr+="]"

    echo "$updated_arr"
}

generate_matrix_markup() {
    local output_array="$1"

    table="| HashiCorp Terraform Release | HashiCorp Terraform Release Date  | HashiCorp Terraform Full Support End Date  | MongoDB Atlas Support End Date |\n"
    table+="|:-------:|:------------:|:------------:|:------------:|\n"

    for ((i=0; i<$(jq length <<< "$output_array"); i++)); do
        version=$(jq -r ".[$i].version" <<< "$output_array")
        published_at=$(jq -r ".[$i].published_at" <<< "$output_array")
        end_support_date=$(jq -r ".[$i].end_support_date" <<< "$output_array")

        table+="| $version | $published_at | $end_support_date | $end_support_date |\n"
    done

    echo -e "$table"
}

update_index_markdown_file() {
    local markup="$1"
    local tempFile="$indexFile.tmp"
    local placeholderStart="<!-- MATRIX_PLACEHOLDER_START -->"
    local placeholderEnd="<!-- MATRIX_PLACEHOLDER_END -->"
    local inPlaceholder=0

    # Ensure the temp file is empty or does not exist
    > "$tempFile"

    while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ "$line" == "$placeholderStart" ]]; then
            inPlaceholder=1
            echo "$line" >> "$tempFile"
            echo "$markup" >> "$tempFile"
            continue
        fi

        if [[ "$line" == "$placeholderEnd" ]]; then
            inPlaceholder=0
            echo "$line" >> "$tempFile"
            continue
        fi

        if [[ $inPlaceholder -eq 0 ]]; then
            echo "$line" >> "$tempFile"
        fi
    done < "$indexFile"

    mv "$tempFile" "$indexFile"

    echo "Updated content between placeholders in $indexFile"
}

# Transform the array data
updated_array=$(transform_array "$input_array")

# Generate markup for the compatibility matrix
markup=$(generate_matrix_markup "$updated_array")

# Update the Markdown file with the generated markup
update_index_markdown_file "$markup"
