#!/usr/bin/env bash
set -euo pipefail

: "${1?"Name of resource or data source must be provided."}"

SDK_BRANCH="${SDK_BRANCH:-"main"}"
# URL to download Atlas Admin API Spec
atlas_admin_api_spec="https://raw.githubusercontent.com/mongodb/atlas-sdk-go/${SDK_BRANCH}/openapi/atlas-api-transformed.yaml"

echo "Downloading api spec"
curl -L "$atlas_admin_api_spec" -o "./api-spec.yml"

resource_name=$1
resource_name_lower_case="$(echo "$resource_name" | awk '{print tolower($0)}')"
resource_name_snake_case="$(echo "$resource_name" | perl -pe 's/([a-z0-9])([A-Z])/$1_\L$2/g')"

pushd "./internal/service/$resource_name_lower_case" || exit

# Running HashiCorp code generation tools

echo "Generating provider code specification"
# Generate provider code specification using api spec and generator config
tfplugingen-openapi generate --config ./tfplugingen/generator_config.yml --output provider-code-spec.json ../../../api-spec.yml

echo "Generating resource and data source schemas and models"
# Generate resource and data sources schemas using provider code specification
tfplugingen-framework generate data-sources --input provider-code-spec.json --output ./ --package "$resource_name_lower_case"
tfplugingen-framework generate resources --input provider-code-spec.json --output ./ --package "$resource_name_lower_case"


rm ../../../api-spec.yml
rm provider-code-spec.json


rename_file() {
    local old_name=$1
    local new_name=$2

    # Check if the original file exists
    if [ -e "$old_name" ]; then
        # If the target new name exists, use the alternative name with _gen
        if [ -e "$new_name" ]; then
            echo "File $new_name already exists, writing content in ${new_name%.*}_gen.go"
            new_name="${new_name%.*}_gen.go"
        fi

        echo "Created file in $new_name"
        # Rename the file
        mv "$old_name" "$new_name"
    fi
}

rename_file "${resource_name_snake_case}_data_source_gen.go" "data_source_schema.go"
rename_file "${resource_name_snake_case}s_data_source_gen.go" "pural_data_source_schema.go"
rename_file "${resource_name_snake_case}_resource_gen.go" "resource_schema.go"
