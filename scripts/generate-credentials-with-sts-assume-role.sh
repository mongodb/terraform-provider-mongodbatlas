#!/bin/bash

# This script uses aws sts assume-role to generate temporary credentials
# and sets them as environment variables so those can be used at a later stage.
# role-arn = arn:aws:iam::358363220050:role/terraform-provider-mongodbatlas-acceptancetests

# Define a function to convert a string to lowercase
function to_lowercase() {
  echo "$1" | tr '[:upper:]' '[:lower:]'
}

# Convert the input string to lowercase
aws_region=$(to_lowercase "$AWS_REGION")
# Replace all underscores with hyphens
aws_region=${aws_region//_/-}

# Print the output string
echo "$aws_region"
export AWS_REGION="$aws_region"

# Get the STS credentials
CREDENTIALS=$(aws sts assume-role --role-arn ${ASSUME_ROLE_ARN} --role-session-name newSession --output text)

# Extract the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN from the STS credentials
AWS_ACCESS_KEY_ID=$(echo $CREDENTIALS | awk '{print $1}')
AWS_SECRET_ACCESS_KEY=$(echo $CREDENTIALS | awk '{print $2}')
AWS_SESSION_TOKEN=$(echo $CREDENTIALS | awk '{print $3}')

# Export the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN environment variables
export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN

# Print the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN to the console
echo "AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}"
echo "AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}"
echo "AWS_SESSION_TOKEN: ${AWS_SESSION_TOKEN}"

# echo "aws_access_key_id=${AWS_ACCESS_KEY_ID}"; echo "aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}"; echo "aws_session_token=${AWS_SESSION_TOKEN}" >> "$GITHUB_OUTPUT"


# response=("$(aws sts assume-role \
#     --role-session-name "newSession" \
#     --role-arn "$ASSUME_ROLE_ARN" \
#     --output text \
#     --query 'Credentials.[AccessKeyId, SecretAccessKey, SessionToken]')")

# echo "${response[0]}"
# echo "${response[1]}"
# echo "${response[2]}"

# export STS_AWS_ACCESS_KEY_ID=${response[0]}
# export STS_AWS_SECRET_ACCESS_KEY=${response[1]}
# export STS_AWS_SESSION_TOKEN=${response[2]}

# echo "aws_access_key_id=${response[0]}"; echo "aws_secret_access_key=${response[1]}"; echo "aws_session_token=${response[2]}" >> "$GITHUB_OUTPUT"
# # echo "aws_secret_access_key=${response[2]}" >> "$GITHUB_OUTPUT"
# # echo "aws_session_token=${response[3]}" >> "$GITHUB_OUTPUT"