#!/usr/bin/env bash

# Copyright 2025 MongoDB Inc
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

# Generate OAuth2 access token from MongoDB Atlas service account credentials
# Required environment variables:
#   MONGODB_ATLAS_BASE_URL: Base URL for MongoDB Atlas API (e.g., https://cloud.mongodb.com/)
#   MONGODB_ATLAS_CLIENT_ID: Service account client ID
#   MONGODB_ATLAS_CLIENT_SECRET: Service account client secret

if [ -z "${MONGODB_ATLAS_BASE_URL:-}" ]; then
  echo "Error: MONGODB_ATLAS_BASE_URL environment variable is required" >&2
  exit 1
fi

if [ -z "${MONGODB_ATLAS_CLIENT_ID:-}" ]; then
  echo "Error: MONGODB_ATLAS_CLIENT_ID environment variable is required" >&2
  exit 1
fi

if [ -z "${MONGODB_ATLAS_CLIENT_SECRET:-}" ]; then
  echo "Error: MONGODB_ATLAS_CLIENT_SECRET environment variable is required" >&2
  exit 1
fi

# Create Basic Authentication header
AUTH_HEADER=$(echo -n "$MONGODB_ATLAS_CLIENT_ID:$MONGODB_ATLAS_CLIENT_SECRET" | base64)

# Remove trailing slash from base URL if present
BASE_URL="${MONGODB_ATLAS_BASE_URL%/}"
TOKEN_ENDPOINT="${BASE_URL}/api/oauth/token"

# Make the token request
RESPONSE=$(curl -s -w "\n%{http_code}" "$TOKEN_ENDPOINT" \
  -H "Accept: application/json" \
  -H "Authorization: Basic $AUTH_HEADER" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials")

# Extract HTTP status code and response body
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')

# Check if the request was successful
if [ "$HTTP_CODE" != "200" ]; then
  echo "Failed to generate OAuth2 token. HTTP Status: $HTTP_CODE" >&2
  echo "Response: $RESPONSE_BODY" >&2
  exit 1
fi

# Extract access token from JSON response
ACCESS_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
  echo "Failed to extract access token from response" >&2
  echo "Response: $RESPONSE_BODY" >&2
  exit 1
fi

echo "$ACCESS_TOKEN"
