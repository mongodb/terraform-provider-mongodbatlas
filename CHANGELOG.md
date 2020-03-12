## 0.4.2 (March 12, 2020)

**Fixed bugs:**

- mongodbatlas\_cluster fails to redeploy manually deleted cluster [\#159](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/159)

**Closed issues:**

- mongodbatlas\_alert\_configuration - not able to generate any alerts with event\_type = "OUTSIDE\_METRIC\_THRESHOLD" and matcher.fieldName != "TYPE\_NAME" [\#164](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/164)
- Can't create cluster ME\_SOUTH\_1 region [\#157](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/157)

**Merged pull requests:**

- chore: fix linting issues [\#169](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/169) ([marinsalinas](https://github.com/marinsalinas))
- Doc: Fix import for mongodbatlas\_project\_ip\_whitelist [\#166](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/166) ([haidaraM](https://github.com/haidaraM))
- chore: removed wrong validation for matchers.value [\#165](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/165) ([PacoDw](https://github.com/PacoDw))
- add default label to clusters [\#163](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/163) ([marinsalinas](https://github.com/marinsalinas))
- Cleaned Cluster state when it isn't found to allow create it again [\#162](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/162) ([PacoDw](https://github.com/PacoDw))
- cluster: removed array of regions due to they could be changed [\#158](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/158) ([PacoDw](https://github.com/PacoDw))

## 0.4.1 (February 26, 2020)

**Fixed bugs:**

- Add name argument in mongodbatlas\_project datasource [\#140](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/140)

**Closed issues:**

- Delete timeout for mongodbatlas\_private\_endpoint resource too short [\#151](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/151)
- mongodbatlas\_project name [\#150](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/150)
- mongodbatlas\_custom\_db\_role not waiting for resource creation [\#148](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/148)

**Merged pull requests:**

- Custom DB Roles: added refresh function to allow to create/remove multiple custom roles at the same time [\#155](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/155) ([PacoDw](https://github.com/PacoDw))
- chore: increase timeout when delete in private\_endpoint resource [\#154](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/154) ([marinsalinas](https://github.com/marinsalinas))
- add upgrade guide [\#149](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/149) ([themantissa](https://github.com/themantissa))
- Correct `mongodbatlas\_teams` resource name in docs [\#143](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/143) ([mattt416](https://github.com/mattt416))
- Project data source [\#142](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/142) ([PacoDw](https://github.com/PacoDw))

## 0.4.0 (February 18, 2020)

**Implemented enhancements to existing resources:**

- add mongodbatlas\_project datasource [\#137](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/137)
- Database Users: updated Read Function to avoid plugin error when it upgrades [\#133](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/133) ([PacoDw](https://github.com/PacoDw))
- Fix snapshot import with hyphened cluster\_name [\#131](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/131) ([marinsalinas](https://github.com/marinsalinas))
- chore: added database\_name as deprecated attribute [\#129](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/129) ([PacoDw](https://github.com/PacoDw))
- Encryption At Rest: fixed issues and added an enhancement [\#128](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/128) ([PacoDw](https://github.com/PacoDw))
- Add PIT enabled argumento to Cluster Resource and Data Source [\#126](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/126) ([marinsalinas](https://github.com/marinsalinas))
- Issue with import  mongodbatlas\_cloud\_provider\_snapshot\_restore\_job [\#109](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/109)
- Fix \#68: Added the ability to re-create the whitelist entry when it's remove manually [\#106](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/106) ([PacoDw](https://github.com/PacoDw))
- Add pitEnabled feature of mongodbatlas\_cluster resource [\#104](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/104)
- Updating `git clone` command to reference current repository [\#103](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/103) ([macintacos](https://github.com/macintacos))
- Cluster label and plugin attribute [\#102](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/102) ([PacoDw](https://github.com/PacoDw))
- Added functions to handle labels attribute in some resources [\#101](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/101) ([PacoDw](https://github.com/PacoDw))
- Added labels attr for Database User resource [\#100](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/100) ([PacoDw](https://github.com/PacoDw))
- Add provider\_name to containers data source [\#95](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/95) ([marinsalinas](https://github.com/marinsalinas))
- Whitelist [\#94](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/94) ([PacoDw](https://github.com/PacoDw))
- Network Peering RS: remove provider\_name=AWS as default, use Required=true instead i… [\#92](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/92) ([marinsalinas](https://github.com/marinsalinas))
- Added default Disk Size when it doesn't set up on cluster resource [\#77](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/77) ([PacoDw](https://github.com/PacoDw))
- mongodbatlas\_network\_containers datasource doesn't work with Azure [\#71](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/71)
- Support for AWS security groups in mongodbatlas\_project\_ip\_whitelist [\#67](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/67)

**Documentation improvements and fixes:**

- Breaking Change - [Upgrade Guide] (https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040) ([themantissa](https://github.com/themantissa))
- Readme: Updated env variables [\#135](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/135) ([PacoDw](https://github.com/PacoDw))
- Spelling and grammer [\#130](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/130) ([CMaylone](https://github.com/CMaylone))
- Shared tier doc edits [\#119](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/119) ([themantissa](https://github.com/themantissa))
- Update cluster doc [\#117](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/117) ([themantissa](https://github.com/themantissa))
- Update backup, add links [\#116](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/116) ([themantissa](https://github.com/themantissa))
- Update cluster.html.markdown [\#112](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/112) ([themantissa](https://github.com/themantissa))
Update database\_user.html.markdown [\#98](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/98) ([themantissa](https://github.com/themantissa))
- update containers and ip whitelist doc [\#96](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/96) ([themantissa](https://github.com/themantissa))
- Update project.html.markdown [\#91](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/91) ([themantissa](https://github.com/themantissa))
- website: collapse data sources sidebar by default [\#75](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/75) ([marinsalinas](https://github.com/marinsalinas))
- Improvements to Peering Resources [\#73](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/73) ([themantissa](https://github.com/themantissa))
- Remove dupe argument in docs [\#69](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/69) ([ktmorgan](https://github.com/ktmorgan))
- Clarify Azure Option in Doc [\#66](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/66) ([themantissa](https://github.com/themantissa))


**Fixed bugs:**

- Fix DiskSizeGB missing [\#111](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/111) ([marinsalinas](https://github.com/marinsalinas))
- Fix peering resource [\#107](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/107) ([PacoDw](https://github.com/PacoDw))
- fix: validate if mongo\_db\_major\_version is set in cluster resource [\#85](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/85) ([marinsalinas](https://github.com/marinsalinas))
- Cannot update GCP network peer [\#86](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/86)
- Cluster fails to build on 0.3.1 when mongo\_db\_major\_version is not specified [\#81](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/81)
- Crash \(panic, interface conversion error\) when creating mongodbatlas\_encryption\_at\_rest in Azure [\#74](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/74)
- Creating M2 cluster without specifying disk\_size\_gb results in 400 Bad Request [\#72](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues/72)

**New features:**

- X509 Authentication Database User [\#125](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/125) ([PacoDw](https://github.com/PacoDw))
- Database users: added x509\_type attribute [\#122](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/122) ([PacoDw](https://github.com/PacoDw))
- Private endpoints [\#118](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/118) ([PacoDw](https://github.com/PacoDw))
- Projects: adding teams attribute [\#113](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/113) ([PacoDw](https://github.com/PacoDw))
- Terraform resource for MongoDB Custom Roles [\#108](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/108) ([PacoDw](https://github.com/PacoDw))
- Alert configuration resource and data source [\#99](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/99) ([PacoDw](https://github.com/PacoDw))
- Feat: Global Cluster Configuration Resource and Data Source. [\#90](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/90) ([marinsalinas](https://github.com/marinsalinas))
- Auditing Resource and Data Source [\#82](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/82) ([PacoDw](https://github.com/PacoDw))
- Feat: Team Resource and Data Source [\#79](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/79) ([marinsalinas](https://github.com/marinsalinas))
- Maintenance window ds [\#78](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/78) ([PacoDw](https://github.com/PacoDw))
- Maintenance window rs [\#76](https://github.com/terraform-providers/terraform-provider-mongodbatlas/pull/76) ([PacoDw](https://github.com/PacoDw))

## 0.3.1 (November 11, 2019)

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
