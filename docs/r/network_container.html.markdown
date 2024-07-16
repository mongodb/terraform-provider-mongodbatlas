# Resource: mongodbatlas_network_container

`mongodbatlas_network_container` provides a Network Peering Container resource. The resource lets you create, edit and delete network peering containers. You must delete network peering containers before creating clusters in your project. You can't delete a network peering container if your project contains clusters. The resource requires your Project ID.  Each cloud provider requires slightly different attributes so read the argument reference carefully.

 Network peering container is a general term used to describe any cloud providers' VPC/VNet concept.  Containers only need to be created if the peering connection to the cloud provider will be created before the first cluster that requires the container.  If the cluster has been/will be created first Atlas automatically creates the required container per the "containers per cloud provider" information that follows (in this case you can obtain the container id from the cluster resource attribute `container_id`).

The following is the maximum number of Network Peering containers per cloud provider:
<br> &#8226;  GCP -  One container per project.
<br> &#8226;  AWS and Azure - One container per cloud provider region.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Example with AWS

```terraform
  resource "mongodbatlas_network_container" "test" {
    project_id       = "<YOUR-PROJECT-ID>"
    atlas_cidr_block = "10.8.0.0/21"
    provider_name    = "AWS"
    region_name      = "US_EAST_1"
  }

```

### Example with GCP

```terraform
resource "mongodbatlas_network_container" "test" {
  project_id       = "<YOUR-PROJECT-ID>"
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "GCP"
  regions = ["US_EAST_4", "US_WEST_3"]
}
```

### Example with Azure

```terraform
resource "mongodbatlas_network_container" "test" {
  project_id       = "<YOUR-PROJECT-ID>"
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AZURE"
  region           = "US_EAST_2"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the Atlas project for this Network Peering Container.
* `atlas_cidr_block` - (Required) CIDR block that Atlas uses for the Network Peering containers in your project.  Atlas uses the specified CIDR block for all other Network Peering connections created in the project. The Atlas CIDR block must be at least a /24 and at most a /21 in one of the following [private networks](https://tools.ietf.org/html/rfc1918.html#section-3):
  * Lower bound: 10.0.0.0 -	Upper bound: 10.255.255.255 -	Prefix: 10/8
  * Lower bound: 172.16.0.0 -	Upper bound:172.31.255.255 -	Prefix:	172.16/12
  * Lower bound: 192.168.0.0 -	Upper bound:192.168.255.255 -	Prefix:	192.168/16

    **Atlas locks this value** if an M10+ cluster or a Network Peering connection already exists. To modify the CIDR block, ensure there are no M10+ clusters in the project and no other Network Peering connections in the project.

    **Important**: Atlas limits the number of MongoDB nodes per Network Peering connection based on the CIDR block and the region selected for the project. Contact [MongoDB Support](https://www.mongodb.com/contact?tck=docs_atlas) for any questions on Atlas limits of MongoDB nodes per Network Peering connection.

* `provider_name`  - (Required GCP and AZURE, Optional but recommended for AWS) Cloud provider for this Network Peering connection.  Accepted values are GCP, AWS, AZURE. If omitted, Atlas sets this parameter to AWS.
* `region_name` - (Required AWS only) The Atlas AWS region name for where this container will exist, see the reference list for Atlas AWS region names [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `region` - (Required AZURE only) Atlas region where the container resides, see the reference list for Atlas Azure region names [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
* `regions` - (Optional GCP only) Atlas regions where the container resides. Provide this field only if you provide an `atlas_cidr_block` smaller than `/18`. [GCP Regions values](https://docs.atlas.mongodb.com/reference/api/vpc-create-container/#request-body-parameters).



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `container_id` - The Network Peering Container ID.
* `id` - Terraform's unique identifier used internally for state management.
* `provisioned` - Indicates whether the project has Network Peering connections deployed in the container.

**AWS ONLY:**

* `region_name` - Atlas name for AWS region where the Atlas container resides.
* `vpc_id` - Unique identifier of Atlas' AWS VPC.

**GCP ONLY:**

* `gcp_project_id` - Unique identifier of the GCP project in which the network peer resides. Returns null. This value is populated once you create a new network peering connection with the network peering resource.
* `network_name` - Unique identifier of the Network Peering connection in the Atlas project. Returns null. This value is populated once you create a new network peering connection with the network peering resource.

**AZURE ONLY:**

* `region` - Azure region where the Atlas container resides.
* `azure_subscription_id` - Unique identifier of the Azure subscription in which the VNet resides.
* `vnet_name` - 	The name of the Azure VNet. Returns null. This value is populated once you create a new network peering connection with the network peering resource.


## Import

Network Peering Containers can be imported using project ID and network peering container id, in the format `PROJECTID-CONTAINER-ID`, e.g.

```
$ terraform import mongodbatlas_network_container.my_container 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a
```

See detailed information for arguments and attributes: [MongoDB API Network Peering Container](https://docs.atlas.mongodb.com/reference/api/vpc-create-container/)

-> **NOTE:** If you need to get an existing container ID see the [How-To Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/howto-guide.html).