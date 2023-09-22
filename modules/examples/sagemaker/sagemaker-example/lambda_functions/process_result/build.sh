#! /bin/bash

set -x
# Change these variables according to your configuration.
profile=default
region=us-east-1
repo_name=process-result

# Get AWS account ID.
account_id=$(aws sts get-caller-identity --profile ${profile} --query Account --output text)

# If the repository doesn't exist in ECR, create it.
if ! aws ecr describe-repositories \
    --profile "${profile}" \
    --region "${region}" \
    --repository-names "${repo_name}" > /dev/null 2>&1
then
    aws ecr create-repository \
        --profile "${profile}" \
        --region "${region}" \
        --repository-name "${repo_name}" > /dev/null
fi

# Authenticate Docker to the Amazon ECR private registry.
aws ecr get-login-password \
    --profile "${profile}" \
    --region "${region}" \
| docker login \
    --username AWS \
    --password-stdin "${account_id}.dkr.ecr.${region}.amazonaws.com"

# Build the docker image locally with the image name and then push it to ECR
# with the full name.
full_name="${account_id}.dkr.ecr.${region}.amazonaws.com/${repo_name}:latest"
docker build -q -t "${repo_name}" .
docker tag "${repo_name}" "${full_name}"
docker push "${full_name}"

printf "\nPushLambdaECRImageURI: %s\n\n" "${full_name}"
