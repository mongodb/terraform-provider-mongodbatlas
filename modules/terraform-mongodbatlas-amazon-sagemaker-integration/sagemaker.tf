provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}


resource "mongodbatlas_event_trigger" "trigger" {
  project_id = var.atlas_project_id
  name = var.trigger_name
  type = "DATABASE"
  app_id = var.realm_app_id

  config_database= var.database_name
  config_collection = var.collection_name
  config_operation_types = ["INSERT"]
  config_service_id = var.service_id
  config_full_document = true

  event_processors {
      aws_eventbridge  {
          config_region = data.aws_region.current.name  
          config_account_id = data.aws_caller_identity.current.account_id
      }
    }
}

resource "aws_iam_role" "sage_maker_execution_role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = [
            "sagemaker.amazonaws.com"
          ]
        }
        Action = [
          "sts:AssumeRole"
        ]
      }
    ]
  })
  path = "/"
  managed_policy_arns = [
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonSageMakerFullAccess",
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonSageMakerCanvasFullAccess"
  ]

  inline_policy {
    name = "qs-sagemaker-execution-policy"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Effect = "Allow",
          Action = "s3:GetObject",
          Resource = "arn:${data.aws_partition.current.partition}:s3:::*"
        }
      ]
    })
  }
}

resource "aws_sagemaker_model" "model" {
  primary_container {
    image = var.model_ecr_image_uri
    model_data_url = var.model_data_s3_uri
    mode = "SingleModel"
    environment = {
      SAGEMAKER_PROGRAM = "inference.py"
      SAGEMAKER_SUBMIT_DIRECTORY = var.model_data_s3_uri
    }
  }
  execution_role_arn = aws_iam_role.sage_maker_execution_role.arn
}

resource "aws_sagemaker_endpoint_configuration" "endpoint_config" {
  production_variants  {
      initial_instance_count = 1
      initial_variant_weight = 1.0
      instance_type = "ml.c5.large"
      model_name = aws_sagemaker_model.model.name
      variant_name = aws_sagemaker_model.model.name
    }
}

resource "aws_sagemaker_endpoint" "endpoint" {
  endpoint_config_name = aws_sagemaker_endpoint_configuration.endpoint_config.name
}

resource "aws_cloudwatch_event_bus" "event_bus_for_capturing_mdb_events" {
  depends_on = [ mongodbatlas_event_trigger.trigger ]
  event_source_name = "aws.partner/mongodb.com/stitch.trigger/${mongodbatlas_event_trigger.trigger.trigger_id}"
  name = "aws.partner/mongodb.com/stitch.trigger/${mongodbatlas_event_trigger.trigger.trigger_id}"
}

resource "aws_cloudwatch_event_bus" "event_bus_for_sage_maker_results" {
  name = "qs-mongodb-sagemaker-results"
}

resource "aws_lambda_function" "lambda_function_to_read_mdb_events" {
  function_name = "pull-mdb-events"
  package_type = "Image"
  image_uri = var.pull_lambda_ecr_image_uri
  role = aws_iam_role.pull_lambda_function_role.arn
  environment {
    variables = {
      model_endpoint = aws_sagemaker_endpoint.endpoint.name
      region_name = data.aws_region.current.name
      eventbus_name = aws_cloudwatch_event_bus.event_bus_for_sage_maker_results.arn
    }
  }
  architectures = [
    "x86_64"
  ]
  memory_size = 1024
  timeout = 300
}

resource "aws_cloudwatch_event_rule" "event_rule_to_match_mdb_events" {
  description = "Event Rule to match MongoDB change events."
  event_bus_name = aws_cloudwatch_event_bus.event_bus_for_capturing_mdb_events.name
  event_pattern = jsonencode({
    account = [
      data.aws_caller_identity.current.account_id
    ]
  })
  is_enabled = true
  name = "pull-mdb-events"
}

resource "aws_cloudwatch_event_target" "read_mdb_event_target" {
  event_bus_name = aws_cloudwatch_event_bus.event_bus_for_capturing_mdb_events.name
  rule      = aws_cloudwatch_event_rule.event_rule_to_match_mdb_events.name
  target_id = "EventRuleToReadMatchMDBEventsID"
  arn =  aws_lambda_function.lambda_function_to_read_mdb_events.arn
}

resource "aws_iam_role" "pull_lambda_function_role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = [
            "lambda.amazonaws.com"
          ]
        }
        Action = [
          "sts:AssumeRole"
        ]
      }
    ]
  })
  path = "/"
  managed_policy_arns = [
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  ]
    inline_policy {
    name = "sagemaker-endpoint-invokation-policy"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Effect = "Allow",
          Action = "sagemaker:InvokeEndpoint",
          Resource = aws_sagemaker_endpoint.endpoint.arn
        },
        {
          Effect = "Allow",
          Action = "events:PutEvents",
          Resource = aws_cloudwatch_event_bus.event_bus_for_sage_maker_results.arn
        }
      ]
    })
  }
}

resource "aws_lambda_function" "lambda_function_to_write_to_mdb" {
  function_name = "push_lambda_function"
  package_type = "Image"
  role = aws_iam_role.push_lambda_function_role.arn
  image_uri = var.push_lambda_ecr_image_uri
  environment {
    variables = {
      mongo_endpoint = var.mongo_endpoint
      dbname = var.database_name
    }
  }
  architectures = [
    "x86_64"
  ]
  memory_size = 1024
  timeout = 300
}

resource "aws_iam_role" "push_lambda_function_role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = [
            "lambda.amazonaws.com"
          ]
        }
        Action = [
          "sts:AssumeRole"
        ]
      }
    ]
  })
  path = "/"
  managed_policy_arns = [
    "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  ]
  inline_policy {
    name = "sagemaker-endpoint-invokation-policy"
    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Effect = "Allow",
          Action = "sagemaker:InvokeEndpoint",
          Resource = aws_sagemaker_endpoint.endpoint.arn
        },
        {
          Effect = "Allow",
          Action = "events:PutEvents",
          Resource = aws_cloudwatch_event_rule.event_rule_to_match_mdb_events.arn
        }
      ]
    })
  }
}

resource "aws_cloudwatch_event_rule" "event_rule_to_capture_events_sent_from_lambda_function" {
  description = "Event Rule to match result events returned by pull Lambda."
  event_bus_name = aws_cloudwatch_event_bus.event_bus_for_sage_maker_results.name
  event_pattern = jsonencode({
    source = [
      "user-event"
    ]
    detail-type = [
      "user-preferences"
    ]
  })
  is_enabled = true
  name = "push-to-mongodb"
}

resource "aws_cloudwatch_event_target" "write_event_from_lambda_to_target" {
  event_bus_name = aws_cloudwatch_event_bus.event_bus_for_sage_maker_results.name
  rule      = aws_cloudwatch_event_rule.event_rule_to_capture_events_sent_from_lambda_function.name
  target_id = "EventRuleToCaptureEventsSentFromLambdaFunctionID"
  arn =  aws_lambda_function.lambda_function_to_write_to_mdb.arn
}

resource "aws_lambda_permission" "event_bridge_lambda_permission1" {
  function_name = aws_lambda_function.lambda_function_to_read_mdb_events.arn
  action = "lambda:InvokeFunction"
  principal = "events.amazonaws.com"
  source_arn = aws_cloudwatch_event_rule.event_rule_to_match_mdb_events.arn
}

resource "aws_lambda_permission" "event_bridge_lambda_permission2" {
  function_name = aws_lambda_function.lambda_function_to_write_to_mdb.arn
  action = "lambda:InvokeFunction"
  principal = "events.amazonaws.com"
  source_arn = aws_cloudwatch_event_rule.event_rule_to_capture_events_sent_from_lambda_function.arn
}