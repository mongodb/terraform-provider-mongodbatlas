# Atlas Terraform Provider Example: Encryption at rest - AWS

This example sets up encryption at rest using AWS KMS for your Atlas Project. It creates the encryption key in AWS KMS, as well as an IAM role and policy so that Atlas can access the key, and finally enables encryption at rest for the Atlas Project.

## Usage

1. Set your Atlas public & private API key via environment variables:

        $ export MONGODB_ATLAS_PUBLIC_KEY="xxxxxxxx"
        $ export MONGODB_ATLAS_PRIVATE_KEY="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

2. Set your AWS access key & secret via environment variables:

        $ export AWS_ACCESS_KEY_ID="xxxxxxxxxxxxxxxxxxxx"
        $ export AWS_SECRET_ACCESS_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

3. Initialize terraform: `terraform init`
4. Run terraform apply and supply values for any of the variables when
prompted: `terraform apply`

## (Optional) Enabling encryption at rest for an existing cluster

1. Import the cluster using the Project ID and cluster name (e.g. `5beae24579358e0ae95492af-MyCluster`):

        $ terraform import mongodbatlas_cluster.my_cluster ProjectId-ClusterName

2. Add any non-default values to the cluster resource *mongodbatlas_cluster.my_cluster* in *main.tf*. And add the following attribute: `encryption_at_rest_provider = "AWS"`
3. Run terraform apply to enable encryption at rest for the cluster: `terraform apply`
4. (Optional) To remove the cluster from TF state, in case you want to disable project-level encryption and delete the role and key without deleting the imported cluster:
    1. First disable encryption on the cluster by changing the attribute `encryption_at_rest_provider = "NONE"` for the cluster resource *mongodbatlas_cluster.my_cluster* in *main.tf*. If you skip this and the next step, you won't be able to disable encryption on the project-level
    2. Run terraform apply to disable encryption for the cluster: `terraform apply`
    3. Finally, remove the cluster from TF state:

            terraform state rm mongodbatlas_cluster.my_cluster

    4. You should now be able to run terraform destroy without deleting the cluster: `terraform destroy`
