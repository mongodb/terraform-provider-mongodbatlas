
output "sage_maker_endpoint_arn" {
  description = "SageMaker endpoint ARN"
  value = aws_sagemaker_endpoint.endpoint.arn
}

output "event_bus_name" {
  description = "Event Bus Name"
  value = aws_cloudwatch_event_bus.event_bus_for_capturing_mdb_events.arn
}
