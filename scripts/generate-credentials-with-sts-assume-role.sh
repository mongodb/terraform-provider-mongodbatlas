#!/bin/bash

# This script uses aws sts assume-role to generate temporary credentials
# and sets them as environment variables so those can be used at a later stage.
# role-arn = arn:aws:iam::358363220050:role/terraform-provider-mongodbatlas-acceptancetests

response=("$(aws sts assume-role \
    --role-session-name "newSession" \
    --role-arn "$ASSUME_ROLE_ARN" \
    --output text \
    --query 'Credentials.[AccessKeyId, SecretAccessKey, SessionToken]')")

# echo "${response[1]}"
# echo "${response[2]}"
# echo "${response[3]}"

export STS_AWS_ACCESS_KEY_ID=${response[1]}
export STS_AWS_SECRET_ACCESS_KEY=${response[2]}
export STS_AWS_SESSION_TOKEN=${response[3]}

echo "aws_access_key_id=${response[1]}"; echo "aws_secret_access_key=${response[2]}"; echo "aws_session_token=${response[3]}" >> "$GITHUB_OUTPUT"
# echo "aws_secret_access_key=${response[2]}" >> "$GITHUB_OUTPUT"
# echo "aws_session_token=${response[3]}" >> "$GITHUB_OUTPUT"