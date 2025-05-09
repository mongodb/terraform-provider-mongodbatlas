resource_name="$1"

cd ./examples/mongodbatlas_$resource_name
terraform fmt -recursive
terraform validate
