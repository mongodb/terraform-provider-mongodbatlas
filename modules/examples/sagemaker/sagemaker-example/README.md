# Iris Flower Classification Example - Steps to Launch Solution

## Prerequisites

- Docker up and running
- AWS CLI
- Python 3

## Upload Model Artifacts to S3

**Note**: Update `profile_name` and `region_name` in `sagemaker-example/iris_classifier/upload.py` if you are not using AWS `default` profile.

```bash
$ cd sagemaker-example/iris_classifier
$ python3 upload.py

# Output
ModelDataS3URI: s3://sagemaker-ap-south-1-123456789012/model.tar.gz
ModelECRImageURI: 123456789034.dkr.ecr.ap-south-1.amazonaws.com/sagemaker-scikit-learn:0.23-1-cpu-py3
```

## Build and Push Lambda Functions

**Note**: Make sure Docker is up and running. Update `profile` and `region` in `sagemaker-example/lambda_functions/process_mdb_change_event/build.sh` and `sagemaker-example/lambda_functions/process_result/build.sh` if you are not using AWS `default` profile.

### Pull Lambda (Reads MongoDB Change Events)

```bash
$ cd sagemaker-example/lambda_functions/process_mdb_change_event
$ ./build.sh

# Output
Login Succeeded
sha256:822e5bf88cabe9dc1cb67258c022f99a3f549458a871b3a4f7b47686b4c20dd5
The push refers to repository [123456789012.dkr.ecr.ap-south-1.amazonaws.com/process-mdb-change-event]
5a2a24566d3d: Pushed
04c86612417d: Pushed
aa07e8319e17: Pushed
add489bc36b0: Pushed
2d244e0816c6: Pushed
1a3c18657fd6: Pushed
1a1430bb3d51: Pushed
b2c122fc6a0b: Pushed
latest: digest: sha256:299589ce12e3983177f1b6ab78324f241cbf496ae6adb4fd558db8d81280f2da size: 2002

PullLambdaECRImageURI: 123456789012.dkr.ecr.ap-south-1.amazonaws.com/process-mdb-change-event:latest
```

### Push Lambda (Writes Results to MongoDB)

```bash
$ cd sagemaker-example/lambda_functions/process_result
$ ./build.sh

# Output
Login Succeeded
sha256:04cbd21a8130e9bf4fc5d89c0cb8242dabaa370fba3abde35fc9d1d231340e8e
The push refers to repository [123456789012.dkr.ecr.ap-south-1.amazonaws.com/process-result]
a0c1c25f2660: Pushed
db3a9626d316: Pushed
aa07e8319e17: Pushed
add489bc36b0: Pushed
2d244e0816c6: Pushed
1a3c18657fd6: Pushed
1a1430bb3d51: Pushed
b2c122fc6a0b: Pushed
latest: digest: sha256:b8a0060e5358a7713ee18554d5543a68902171b39728ff131a14de9c4ab3f093 size: 2000

PushLambdaECRImageURI: 123456789012.dkr.ecr.ap-south-1.amazonaws.com/process-result:latest
```

## Create MongoDB Cluster

terraform module can be found [here](../../atlas-basic)

## Create Realm App and Service

Please check this [link](https://www.mongodb.com/docs/atlas/app-services/) to learn more about Atlas App Services.

You need an Atlas App Service to create triggers which respond to events like insert, update, read or delete on a collection.

Check the [Atlas App Services API (3.0)](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/) to create a Realm App and App Service.

Here are the direct links to API documentation for creating an App and a Service.

- [Get authentication tokens](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/#section/Get-Authentication-Tokens) (Use the access token as Bearer token in Authorization header in the further calls)
- [Create App](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/#tag/apps/operation/adminCreateApplication)
- [Create Service](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/#tag/services/operation/adminCreateService)

**Note**: Keep the App ID and Service ID handy for further steps.

## Launch the Solution

Finally execute the CloudFormation template `templates/mongodb-sagemaker-analytics-main.template.yaml` with the valid parameters.

### Testing

Insert below document in MongoDB collection `iris`.

```json
{
    "data": [
        [6.1, 2.8, 5.6, 2.2],
        [6.1, 2.8, 4.7, 1.2]
    ]
}
```

You should get the output in `predictions` collection as follows:

```json
{
    "_id": "63e3aa089506176cfb1d8cfc",
    "prediction": [
        "virginica",
        "versicolor"
    ],
    "inp_id":"63e3a9faf1f0658c92848094"
}
```
