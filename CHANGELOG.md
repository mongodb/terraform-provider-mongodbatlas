## 0.3.1 (Unreleased)

**Fixed bugs:**

- Confirmation on timelimit for a terraform apply [\#57](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/57)

**Closed issues:**

- Not able to create M0 clusters [\#64](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/64)
- No way to modify advanced configuration options for a cluster [\#61](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/61)
- mongodbatlas\_network\_peering outputting invalid json [\#59](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/59)
- Syntax are not mandatory and creates confusion [\#58](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/58)
- data source mongodbatlas\_network\_peering retrieves the same for id and connection\_id [\#56](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/56)
- Add resource for maintenance window [\#55](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/55)
- Error encryption\_at\_rest  rpc unavailable desc [\#54](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/54)
- specify oplog size? [\#52](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/52)
- Modifying project ip whitelist destroy and re-create all resources [\#51](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/51)
- Add resource for custom database roles [\#50](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/50)
- An invalid enumeration value US\_EAST\_1 was specified. [\#49](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/49)
- Version 0.3.0 [\#47](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/47)
- plugin.terraform-provider-mongodbatlas\_v0.2.0\_x4: panic: runtime error: index out of range [\#36](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/36)

**Merged pull requests:**

- Added format function to handle the mongo\_db\_major\_version attr [\#63](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/63) ([PacoDw](https://github.com/PacoDw))
- Added cast func to avoid panic by nil value [\#62](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/62) ([PacoDw](https://github.com/PacoDw))
- Cluster advanced configuration Options [\#60](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/60) ([PacoDw](https://github.com/PacoDw))

## 0.3.0 (October 08, 2019)

**Closed issues:**

- Upgrade from M2 to M10 fails [\#42](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/42)
- GCP Peering endless terraform apply [\#41](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/41)
- AWS clusters default provider\_encrypt\_ebs\_volume to false [\#40](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/40)
- mongodbatlas\_network\_peering Internal Servier Error [\#35](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/35)
- Problem encryption\_at\_rest [\#33](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/33)
- Problem destroying network peering container [\#30](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/30)
- Bug VPC Peering between GCP and Atlas [\#29](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/29)

**Merged pull requests:**

- Clarify Doc Examples and Text [\#46](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/46) ([themantissa](https://github.com/themantissa))
- fix-\#40: added true value by defualt on provider\_encrypt\_ebs\_volume attr [\#45](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/45) ([PacoDw](https://github.com/PacoDw))
- make provider\_name forced new to avoid patch problems [\#44](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/44) ([marinsalinas](https://github.com/marinsalinas))
- Network peering [\#43](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/43) ([PacoDw](https://github.com/PacoDw))
- Update readme with more info [\#39](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/39) ([themantissa](https://github.com/themantissa))
- Fix: Network container [\#38](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/38) ([PacoDw](https://github.com/PacoDw))
- Doc updates [\#37](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/37) ([themantissa](https://github.com/themantissa))

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
