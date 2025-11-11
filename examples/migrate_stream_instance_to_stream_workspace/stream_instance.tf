# BEFORE: Original stream instance resource (comment out after migration)
# resource "mongodbatlas_stream_instance" "example" {
#   project_id = var.project_id
#   instance_name = var.workspace_name
#   data_process_region = {
#     region = "VIRGINIA_USA"
#     cloud_provider = "AWS"
#   }
#   stream_config = {
#     tier = "SP30"
#   }
# }
