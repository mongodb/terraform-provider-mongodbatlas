## 0.2.0 (September 19, 2019)

**Closed issues:**

- Unable to create project with peering only connections [\#24](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/24)
- importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#21](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/21)
- Support the vscode terraform extension [\#19](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/19)
- Bug: VPC Peering Atlas-GCP [\#17](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/17)
- PATCH network peering failed with no peer found [\#14](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/14)

**Merged pull requests:**

- Add Private IP Mode Resource. [\#32](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/32) ([marinsalinas](https://github.com/marinsalinas))
- Moved provider\_name values to the correct section [\#31](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/31) ([kgriffiths](https://github.com/kgriffiths))
- website: add links to Atlas Region name reference. [\#28](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/28) ([themantissa](https://github.com/themantissa))
- Encryption at rest fix [\#27](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/27) ([marinsalinas](https://github.com/marinsalinas))
- website: make resources side nav expanded as default [\#25](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/25) ([marinsalinas](https://github.com/marinsalinas))
- fix: importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#23](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/23) ([PacoDw](https://github.com/PacoDw))
- Fix \#14: PATCH network peering failed with no peer found [\#22](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/22) ([PacoDw](https://github.com/PacoDw))
- fix: change the test configuration for AWS and GCP [\#20](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/20) ([PacoDw](https://github.com/PacoDw))


## 0.1.1 (September 05, 2019)

**Fixed bugs:**

- panic: runtime error: index out of range [\#1](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/1)

**Closed issues:**

- GCP peering problem [\#16](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/16)
- Cluster creation with Azure provider failed [\#15](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/15)
- Error creating MongoDB Cluster [\#9](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/9)
- Failed to create Atlas network peering container [\#7](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/7)
- Bug: Invalid attribute diskIOPS specified. [\#2](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/2)

**Merged pull requests:**

- website: fix typo [\#13](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/13) ([heimweh](https://github.com/heimweh))
- fix: add the correct func to check the env variables on peering datasources [\#12](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/12) ([PacoDw](https://github.com/PacoDw))
- Fix diskIOPS attribute for GCP and Azure [\#11](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/11) ([PacoDw](https://github.com/PacoDw))
- website: fix data sources sidebar always collapsed [\#10](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/10) ([marinsalinas](https://github.com/marinsalinas))
- mongodbatlas\_network\_\(peering and container\): add more testing case [\#8](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/8) ([PacoDw](https://github.com/PacoDw))
- website: fix typo in MongoDB Atlas Services [\#5](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/5) ([marinsalinas](https://github.com/marinsalinas))
- Ip whitelist entries: removing all entries whitelist by terraform user [\#4](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/4) ([PacoDw](https://github.com/PacoDw))
- Refactored import function to get all ip\_addresses and cird\_blocks entries [\#3](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/3) ([PacoDw](https://github.com/PacoDw))


## 0.1.0 (August 19, 2019)

Initial Release
