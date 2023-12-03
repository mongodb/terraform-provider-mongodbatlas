# quickstart-mongodb-atlas



## Overview

![image](https://user-images.githubusercontent.com/5663078/229103723-4c6b9ab1-9492-47ba-b04d-7f29079e3817.png)

The Atlas Partner Solutions templates allow you to set up all you need to start using MongoDB Atlas. We provide four different templates:

- Deploy MongoDB Atlas without VPC peering. This option peers MongoDB Atlas with your existing VPC.
- Deploy MongoDB Atlas with VPC peering into a new VPC (end-to-end deployment). This option builds a complete MongoDB Atlas environment within AWS consisting of a project, cluster, and more.
- Deploy MongoDB Atlas with VPC peering into an existing VPC. This option peers MongoDB Atlas with a new VPC.
- Deploy MongoDB Atlas with Private Endpoint. This option connects MongoDB Atlas AWS VPC using Private Endpoint.

All the quickstart templates create an Atlas Project, Cluster, Database User and enable public access into your cluster.



## MongoDB Atlas CFN Resources used by the templates

- [MongoDB::Atlas::Cluster](../../mongodbatlas/resource_mongodbatlas_cluster.go)
- [MongoDB::Atlas::ProjectIpAccessList](../../mongodbatlas/fw_resource_mongodbatlas_project_ip_access_list.go)
- [MongoDB::Atlas::DatabaseUser](../../mongodbatlas/fw_resource_mongodbatlas_database_user.go)
- [MongoDB::Atlas::Project](../../mongodbatlas/fw_resource_mongodbatlas_project.go)
- [MongoDB::Atlas::NetworkPeering](../../mongodbatlas/resource_mongodbatlas_network_peering.go)
- [MongoDB::Atlas::NetworkContainer](../../mongodbatlas/resource_mongodbatlas_network_container.go)
- [MongoDB::Atlas::PrivateEndpoint](../../mongodbatlas/resource_mongodbatlas_privatelink_endpoint.go)


## Environment Configured by the Partner Solution templates
All Partner Solutions templates will generate the following resources:
- An Atlas Project in the organization that was provided as input.
- An Atlas Cluster with authentication and authorization enabled, which cannot be accessed through the public internet.
- A Database user that can access the cluster.
- The IP address range provided as input will be added to the Atlas access list, allowing the cluster to be accessed through the public internet.

The specific resources that will be created depend on which Partner Solutions template is used:

- A new AWS VPC (Virtual Private Cloud) will be created.
- A VPC peering connection will be established between the MongoDB Atlas VPC (where your cluster is located) and the VPC on AWS.

