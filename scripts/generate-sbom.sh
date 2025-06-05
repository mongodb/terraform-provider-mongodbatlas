#!/usr/bin/env bash
set -euo pipefail

# Authenticate Docker to AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 901841024863.dkr.ecr.us-east-1.amazonaws.com

echo "Generating SBOMs..."
docker run --rm \
  -v "$PWD:/pwd" \
  901841024863.dkr.ecr.us-east-1.amazonaws.com/release-infrastructure/silkbomb:2.0 \
  update \
  --purls /pwd/compliance/purls.txt \
  --sbom-out /pwd/compliance/sbom.json