# Database users can be imported using project ID and username, in the format `project_id-username-auth_database_name`, e.g.
terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-my_user-admin

# NOTE: Terraform will want to change the password after importing the user if a password argument is specified.
