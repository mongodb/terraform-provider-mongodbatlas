#!/bin/bash

set -Eeou pipefail

# This script uses aws sts assume-role to generate temporary credentials
# and outputs them in $GITHUB_OUTPUT so those can be used in other workflow jobs.

# Define a function to convert a string to lowercase
function to_lowercase() {
  echo "$1" | tr '[:upper:]' '[:lower:]'
}
# Convert the input string to lowercase
aws_region=$(to_lowercase "$AWS_REGION")
# Replace all underscores with hyphens
aws_region=${aws_region//_/-}
# e.g. from US_EAST_1 to us-east-1

# Get the STS credentials
export AWS_REGION="$aws_region"
CREDENTIALS=$(aws sts assume-role --role-arn "$ASSUME_ROLE_ARN" --role-session-name newSession --output text --query 'Credentials.[AccessKeyId, SecretAccessKey, SessionToken]')

# Extract the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN from the STS credentials
AWS_ACCESS_KEY_ID=$(echo "$CREDENTIALS" | awk '{print $1}')
AWS_SECRET_ACCESS_KEY=$(echo "$CREDENTIALS" | awk '{print $2}')
AWS_SESSION_TOKEN=$(echo "$CREDENTIALS" | awk '{print $3}')

echo "::add-mask::${AWS_ACCESS_KEY_ID}"
echo "::add-mask::${AWS_SECRET_ACCESS_KEY}"
echo "::add-mask::${AWS_SESSION_TOKEN}"

{
  echo "aws_access_key_id=${AWS_ACCESS_KEY_ID}"
  echo "aws_secret_access_key=$AWS_SECRET_ACCESS_KEY"
  echo "aws_session_token=$AWS_SESSION_TOKEN"
} >> "$GITHUB_OUTPUT"
