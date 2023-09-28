# quickstart-mongodb-atlas-analytics-amazon-sagemaker-integration

## Overview

![simple-quickstart-arch](https://user-images.githubusercontent.com/5663078/229119386-0dbc6e30-a060-465e-86dd-f89712b0fc49.png)

This Partner Solutions template enables you to begin working with your machine learning models using MongoDB Atlas Cluster and Amazon SageMaker endpoints. With this template, you can utilize MongoDB as a data source and SageMaker for data analysis, streamlining the process of building and deploying machine learning models.


## MongoDB Atlas terraform Resources used by the templates

- [mongodbatlas_event_trigger](../../mongodbatlas/data_source_mongodbatlas_event_trigger.go)


## Environment Configured by the Partner Solutions template
The Partner Solutions template will generate and configure the following resources:
 - a [MongoDB Partner Event Bus](http://mongodb.com/docs/atlas/app-services/triggers/aws-eventbridge/#std-label-aws-eventbridge)
 - a [database trigger](https://www.mongodb.com/docs/atlas/app-services/triggers/database-triggers/) with your Atlas Cluster
 - lambda functions to run the machine learning model and send the classification results to your MongoDB Atlas Cluster. (See [iris_classifier](https://github.com/mongodb/mongodbatlas-cloudformation-resources/tree/master/examples/quickstart-mongodb-atlas-analytics-amazon-sagemaker-integration/sagemaker-example/iris_classifier) for an example of machine learning model to use with this template. See [lambda_functions](https://github.com/mongodb/mongodbatlas-cloudformation-resources/tree/master/examples/quickstart-mongodb-atlas-analytics-amazon-sagemaker-integration/sagemaker-example/lambda_functions) for an example of lambda functions to use to read and write data to your MongoDB Atlas cluster.)


