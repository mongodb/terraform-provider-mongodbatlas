"""
Uploads model artifacts to a default SageMaker S3 bucket
and prints S3 URI and SageMaker Docker image URI.
"""

import boto3
import sagemaker
import subprocess

profile_name = 'default'
region_name = 'us-east-1'

# Create sessions.
boto3.setup_default_session(
    profile_name=profile_name,
    region_name=region_name
)
boto_session = boto3.session.Session(
    profile_name=profile_name,
    region_name=region_name
)
s3 = boto_session.resource('s3')
sagemaker_session = sagemaker.Session()

# Build tar file with model data and inference code.
bash_cmd = "tar -cpzf model.tar.gz model.joblib inference.py"
process = subprocess.Popen(bash_cmd.split(), stdout=subprocess.PIPE)
process.communicate()

# S3 bucket for model artifacts.
default_bucket = sagemaker_session.default_bucket()

# Upload tar.gz to S3 bucket.
s3.meta.client.upload_file(
    'model.tar.gz',
    default_bucket,
    'model.tar.gz'
)

# Retrieve sklearn ECR image.
image_uri = sagemaker.image_uris.retrieve(
    framework="sklearn",
    region=boto_session.region_name,
    version="0.23-1",
    py_version="py3",
    instance_type="ml.t3.medium",
)

print('\nModelDataS3URI:', f"s3://{default_bucket}/model.tar.gz")
print('ModelECRImageURI:', image_uri, "\n")
