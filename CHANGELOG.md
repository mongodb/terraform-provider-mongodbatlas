# Changelog

## [v1.12.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.2) (2023-10-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.1...v1.12.2)

**Enhancements**

- Supports `data_expiration_rule` parameter in 'mongodbatlas_online_archive` [\#1528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1528) ([AgustinBettati](https://github.com/AgustinBettati))
- Supports new `notifier_id` parameter in `mongodbatlas_alert_configuration` [\#1514](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1514) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- Fixes issue where Encryption at rest returns inconsistent plan when setting secret access key [\#1529](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1529) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes issue where alert configuration data source for third party notifications returns nil pointer [\#1513](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1513) ([AgustinBettati](https://github.com/AgustinBettati))
- Adjusts format of database user resource id as defined in previous versions [\#1506](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1506) ([AgustinBettati](https://github.com/AgustinBettati))
- Removes delete `partition_fields` statements [\#1499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1499) ([marcosuma](https://github.com/marcosuma))
- Changes validation of empty provider credentials from Error to Warning [\#1501](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1501) ([AgustinBettati](https://github.com/AgustinBettati))
- Uses `container_id` from created cluster in example [\#1475](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1475) ([mckornfield](https://github.com/mckornfield))
- Adjusts time for stale github issues to close after 1 week of inactivity [\#1512](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1512) ([AgustinBettati](https://github.com/AgustinBettati))
- Updates 1.10.0-upgrade-guide.html.markdown [\#1511](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1511) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Updates template issue with clearer guidelines [\#1510](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1510) ([marcosuma](https://github.com/marcosuma))
- Avoids including `provider_disk_type_name` property in cluster update request if attribute was removed [\#1508](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1508) ([AgustinBettati](https://github.com/AgustinBettati))

**Deprecations and Removals**

- Removes the data source `mongodbatlas_privatelink_endpoint_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_privatelink_endpoints_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes mongodbatlas_privatelink_endpoint_service_adl [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the resource `mongodbatlas_privatelink_endpoints_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_data_lake` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_data_lakes` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the resource `mongodbatlas_data_lake` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))

**Closed Issues**

- Error changing user [\#1509](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1509)
- Provider "registry.terraform.io/mongodb/mongodbatlas" planned an invalid value [\#1498](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1498)
- IP allowlist range force replacement on 1.12.0 [\#1495](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1495)
- Importing Online Archive resources is missing parameter partition\_fields in terraform state [\#1492](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1492)
- `mongodbatlas_network_container` faulty optional variable regions [\#1490](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1490)
- Matcher not allowing null [\#1489](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1489)
- Provider version 1.12.0 is breaking the resource mongodbatlas\_database\_user \(1.11.1 works correctly\) [\#1485](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1485)
- `mongodbatlas_project_ip_access_list` causes invalid plans [\#1484](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1484)
- Changes to oplog\_min\_retention\_hours not being applied when set to null [\#1481](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1481)
- Create alert with more than 1 notification [\#1473](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1473)

**Internal Improvements**

- Migrates online archive resource and data sources to new SDK [\#1523](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1523) ([AgustinBettati](https://github.com/AgustinBettati))
- Fixes cleanup-test-env script continues if delete of one project fails [\#1516](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1516) ([AgustinBettati](https://github.com/AgustinBettati))
- Updates atlas-sdk to v20231001001 [\#1515](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1515) ([AgustinBettati](https://github.com/AgustinBettati))
- Fixes module naming convention [\#1500](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1500) ([govardhanpagidi](https://github.com/govardhanpagidi))
- Updates cluster update handler to update advanced\_configuration first and make oplog\_min\_retention\_hours non-computed [\#1497](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1497) ([maastha](https://github.com/maastha))
- Adds coverage report to PRs [\#1496](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1496) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Adds sagemaker quickstart to repo [\#1494](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1494) ([govardhanpagidi](https://github.com/govardhanpagidi))
- Closes code block in "Resource: Cloud Provider Access Configuration Paths" documentation page [\#1487](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1487) ([mattt416](https://github.com/mattt416))
- Bump github.com/gruntwork-io/terratest from 0.43.13 to 0.44.0 [\#1482](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1482) ([dependabot[bot]](https://github.com/apps/dependabot))
- Uses retry.StateChangeConf for encryption-at-rest resource. [\#1477](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1477) ([marcosuma](https://github.com/marcosuma))
- Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.28.0 to 2.29.0, github.com/hashicorp/terraform-plugin-framework from 1.3.4 to 1.4.0, github.com/hashicorp/terraform-plugin-go from 0.18.0 to 0.19.0, github.com/hashicorp/terraform-plugin-mux from 0.11.2 to 0.12.0 [\#1468](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1468) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.10.0 to 0.12.0 [\#1466](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1466) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump golang.org/x/net from 0.13.0 to 0.17.0 [\#1524](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1524) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework from 1.4.0 to 1.4.1 [\#1521](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1521) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/zclconf/go-cty from 1.14.0 to 1.14.1 [\#1520](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1520) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.45.0 to 0.46.0 [\#1519](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1519) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.21 to 1.45.24 [\#1518](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1518) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/hcl/v2 from 2.18.0 to 2.18.1 [\#1517](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1517) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.8 to 1.45.21 [\#1505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1505) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.44.0 to 0.45.0 [\#1504](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1504) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.12.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.1-pre1) (2023-09-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.0...v1.12.1-pre1)

**Internal Improvements**

- Updates 1.12.0 release guide and Changelog [\#1488](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1488) ([maastha](https://github.com/maastha))
- Adjusts PR template so we ensure removals and deprecations are made in isolated PRs [\#1480](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1480) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- Adds missing DatabaseRegex field when creating FederatedDataSource [\#1486](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1486) ([marcosuma](https://github.com/marcosuma))

**Closed Issues**

- `tags` vs. `labels` usage in `mongodbatlas_cluster` resource[\#1370](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1370)

## [v1.12.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.0) (2023-09-20)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.11.1...v1.12.0)

**Enhancements**

- Support for `tags` attribute in `mongodbatlas_cluster`, `mongodbatlas_advanced_cluster`, and `mongodbatlas_serverless_instance`. See [Atlas Resource Tags](https://www.mongodb.com/docs/atlas/tags/) to learn more. [\#1461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1461) ([AgustinBettati](https://github.com/AgustinBettati))
- Support for new `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources [\#1432](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1432) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- **Breaking Change**: Fixes an issue where removing `collectionName` from user role doesn't work [\#1471](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1471) ([marcosuma](https://github.com/marcosuma)).
  - Note: As a result, `mongodbatlas_database_user` no longer requires `roles.collection_name` attribute and doesn't support an empty `collection_name`. You should remove any usage of `roles.collection_name = ""` in configurations for this resource when you upgrade to this version. For more details see:  https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.12.0-upgrade-guide. 
- Populates `total_count` in `mongodbatlas_alert_configurations` data source  [\#1476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1476) ([lantoli](https://github.com/lantoli))
- Improves error handling for `cloud_backup_schedule` resource. [\#1474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1474) ([marcosuma](https://github.com/marcosuma))
- Handles incorrect ids when importing `alert_configuration` or `project_ip_access_list` resources [\#1472](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1472) ([lantoli](https://github.com/lantoli))
- Changelog Spelling Fixes  [\#1457](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1457) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Adds `mongodbatlas_database_user` username parameter OIDC footnote in docs [\#1458](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1458) ([Zuhairahmed](https://github.com/Zuhairahmed))

**Deprecations and Removals**

- Deprecation of `labels` attribute in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` resources.
- Remove deprecated fields in `mongodbatlas_alert_configuration` resource [\#1385](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1385) ([AgustinBettati](https://github.com/AgustinBettati))
- Removal of `api_keys` attribute from `mongodbatlas_project` [\#1365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1365) ([maastha](https://github.com/maastha))
- Removal of attributes in `mongodbatlas_encryption_at_rest` resource: aws_kms, azure_key_vault, google_cloud_kms [\#1383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1383) ([maastha](https://github.com/maastha))
- Removal of MongoDB Atlas Terraform Provider v1.12.0 deprecated fields. [\#1418](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1418) ([marcosuma](https://github.com/marcosuma))
  - provider: duration_seconds
  - advanced_cluster resource: bi_connector
  - cloud_backup_snapshot_restore_job resource: delivery_type
  - cloud_provider_access_setup resource: aws
  - cluster resource: bi_connector, provider_backup_enabled, aws_private_link, aws_private_link_srv
  - database_user resource: provider_backup_enabled
  - project_api_key resource: role_names
  - cluster and clusters data sources: bi_connector
  - project_key and project_keys data sources: role_names

**Closed Issues**

- Alert notification interval\_min not working [\#1464](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1464)
- Changing DB user role from readwrite@DB.Col1 to readwrite@DB doesn't work [\#1462](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1462)
- Unable to deploy a DISK\_PARTITION\_UTILIZATION\_DATA AlertConfiguration [\#1410](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1410)
- Error: The terraform-provider-mongodbatlas\_v1.11.0 plugin crashed [\#1396](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1396)
- Trigger resource does not force replacement when app id changes [\#1310](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1310)

**Internal Improvements**

- Bump goreleaser/goreleaser-action from 4 to 5 [\#1470](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1470) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.4 to 1.45.8 [\#1469](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1469) ([dependabot[bot]](https://github.com/apps/dependabot))
- Merge feature branch into master [\#1460](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1460) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bump github.com/aws/aws-sdk-go from 1.45.2 to 1.45.4 [\#1459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1459) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.4.0 to 1.5.1 [\#1455](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1455) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/zclconf/go-cty from 1.13.3 to 1.14.0 [\#1454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1454) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/hcl/v2 from 2.17.0 to 2.18.0 [\#1453](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1453) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump actions/checkout from 3 to 4 [\#1452](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1452) ([dependabot[bot]](https://github.com/apps/dependabot))
- Fix docs and example fix for project\_api\_key resource after removing role\_names deprecated field [\#1441](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1441) ([AgustinBettati](https://github.com/AgustinBettati))
- Add breaking changes strategy for Terraform [\#1431](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1431) ([wtrocki](https://github.com/wtrocki))
- Extract Configure and Metadata framework functions into single implementation [\#1424](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1424) ([AgustinBettati](https://github.com/AgustinBettati))
- Fix INTMDB-1017 - Updated alert configuration schema with required params [\#1421](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1421) ([andreaangiolillo](https://github.com/andreaangiolillo))
- IP Access List doc updates for Terraform Resources, Data Sources [\#1414](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1414) ([zach-carr](https://github.com/zach-carr))
- Avoid diff in state after import for undefined optional attribute in alert config notification [\#1412](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1412) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate Resource: mongodbatlas\_project\_ip\_access\_list to Terraform Plugin Framework [\#1411](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1411) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Remove validation of empty public\_key and private\_key attributes in provider config to avoid breaking change [\#1402](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1402) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate DataSource mongodbatlas\_alert\_configuration to Terraform Plugin Framework [\#1397](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1397) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate DataSource: mongodbatlas\_project\_ip\_access\_list to Terraform Plugin Framework [\#1395](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1395) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrate Resource: mongodbatlas\_database\_user to Terraform Plugin Framework [\#1388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1388) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update project documentation to remove api\_keys references [\#1386](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1386) ([maastha](https://github.com/maastha))
- Migrates `mongodbatlas_alert_configuration` resource and removes deprecated fields [\#1385](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1385) ([AgustinBettati](https://github.com/AgustinBettati))
- Prepares migration to Terraform framework [\#1384](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1384) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrates `mongodbatlas_encryption_at_rest` resource to Terraform Plugin Framework [\#1383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1383) ([maastha](https://github.com/maastha))
- Adds new framework provider, main and acceptance tests to use mux server with existing sdk v2 provider [\#1366](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1366) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrates `mongodbatlas_project` resource to Terraform Plugin Framework and remove api\_keys attribute [\#1365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1365) ([maastha](https://github.com/maastha))

## [v1.11.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.11.1) (2023-09-06)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.11.0...v1.11.1)

**Enhancements**

- Adds Atlas OIDC Database User support to `mongodbatlas_database_user` [\#1382](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1382) ([andreaangiolillo](https://github.com/andreaangiolillo))
  - Note: This feature is only available starting in [MongoDB 7.0](https://www.mongodb.com/evolved#mdbsevenzero) or later. To learn more see https://www.mongodb.com/docs/atlas/security-oidc/  
- Adds Atlas `datasetNames` support in `mongodbatlas_federated_database_instance` [\#1439](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1439) ([maastha](https://github.com/maastha))
- Improves `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` documentation to highlight that `provider_volume_type=STANDARD` is not available for NVMe clusters [\#1430](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1430) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Adds a new example for `mongodbatlas_online_archive` [\#1372](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1372) ([rsolovev](https://github.com/rsolovev))
- Adds a new example for `mongodbatlas_cloud_backup_schedule` to create policies for multiple clusters [\#1403](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1403) ([maastha](https://github.com/maastha))


**Bug Fixes**

- Updates `tag_sets` to `storage_stores.read_preference` in `mongodbatlas_federated_database_instance` [\#1440](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1440) ([maastha](https://github.com/maastha))
- Updates cluster documentation about labels field [\#1425](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1425) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes null pointer error in `mongodbatlas_alert_configuration` [\#1419](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1419) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Updates `mongodbatlas_event_trigger` resource to force replacement when app id changes [\#1387](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1387) ([maastha](https://github.com/maastha))
- Updates deprecation message to 1.12.0 [\#1381](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1381) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes null pointer error in `mongodbatlas_project` data source [\#1377](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1377) ([andreaangiolillo](https://github.com/andreaangiolillo))


**Closed Issues**

- Provider registry registry.terraform.io does not have a provider named registry.terraform.io/hashicorp/mongodbatlas [\#1389](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1389)
- Event Trigger resource doesn't support wildcard collection name [\#1374](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1374)
- 1.11.0 - terraform provider `mongodbatlas_projects` access denied  [\#1371](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1371)
- Terraform-provider-mongodbatlas\_v1.10.2 plugin crashes when including backup schedule [\#1368](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1368)
- Error: Plugin did not respond - panic: interface conversion: interface is nil, not map[string]interface [\#1337](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1337)
- "Error: error deleting organization information" When importing organization [\#1327](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1327)
- `instance_size` for advance cluster marked as optional in the documentation [\#1311](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1311)
- `region_configs.#._specs.instance_size` in `mongodbatlas_advanced_cluster` is required [\#1288](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1288)

**Internal Improvements**

- Updates the release flow to remove the acceptance steps [\#1443](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1443) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bumps github.com/aws/aws-sdk-go from 1.44.334 to 1.45.2 [\#1442](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1442) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.28.0 [\#1429](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1429) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/zclconf/go-cty from 1.13.2 to 1.13.3 [\#1428](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1428) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.43.12 to 0.43.13 [\#1427](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1427) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.329 to 1.44.334 [\#1426](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1426) ([dependabot[bot]](https://github.com/apps/dependabot))
- Removes 3rd shard from 2 shard global cluster example [\#1423](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1423) ([BassT](https://github.com/BassT))
- Updates issue.yml to use issue number as Ticket title [\#1422](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1422) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrates to new Atlas SDK - `mongodbatlas_federated_database_instance` resource [\#1415](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1415) ([maastha](https://github.com/maastha))
- Updates broken links to the Atlas Admin API docs [\#1413](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1413) ([zach-carr](https://github.com/zach-carr))
- Self document make [\#1407](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1407) ([wtrocki](https://github.com/wtrocki))
- Adds instructions for updates of the Atlas SDK [\#1406](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1406) ([wtrocki](https://github.com/wtrocki))
- Bumps github.com/mongodb-forks/digest from 1.0.4 to 1.0.5 [\#1405](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1405) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.324 to 1.44.329 [\#1404](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1404) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps golangci/golangci-lint-action from 3.6.0 to 3.7.0 [\#1393](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1393) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.319 to 1.44.324 [\#1392](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1392) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.43.11 to 0.43.12 [\#1391](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1391) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps go.mongodb.org/atlas from 0.32.0 to 0.33.0 [\#1390](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1390) ([dependabot[bot]](https://github.com/apps/dependabot))
- Improves the release process [\#1380](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1380) ([andreaangiolillo](https://github.com/andreaangiolillo))
- clenaup-test-env.yml [\#1379](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1379) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Creates github action to delete projects in the test env [\#1378](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1378) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bumps github.com/aws/aws-sdk-go from 1.44.314 to 1.44.319 [\#1375](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1375) ([dependabot[bot]](https://github.com/apps/dependabot))
- Adds githooks [\#1373](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1373) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes cluster outage tests [\#1364](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1364) ([andreaangiolillo](https://github.com/andreaangiolillo))

## [v1.11.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.11.0) (2023-08-04)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.2...v1.11.0)

**Enhancements**

- [Azure Service Principles](https://learn.microsoft.com/en-us/azure/active-directory/develop/app-objects-and-service-principals?tabs=browser) support in `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization` [\#1343](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1343) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-545
- Atlas [Shared Cluster Backups](https://www.mongodb.com/docs/atlas/backup/cloud-backup/shared-cluster-backup/) support in `mongodbatlas_shared_tier_snapshot` and `mongodbatlas_shared_tier_restore_job` [\#1324](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1324) ([andreaangiolillo](https://github.com/andreaangiolillo)) and [\#1323](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1323) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-546
- Atlas Project `limits` support in `mongodbatlas_project` [\#1347](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1347) ([AgustinBettati](https://github.com/AgustinBettati)) - INTMDB-554
- New example for Encryption at Rest using Customer Key Management and multi-region cluster [\#1349](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1349) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-340

**Deprecations and Removals**   

- Marking `cloud_provider_access` resource and data source as deprecated [\#1355](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1355) ([AgustinBettati](https://github.com/AgustinBettati)) - INTMDB-967	

**Bug Fixes**

- Update `mongodbatlas_cloud_backup_schedule` to add the ID field to policyItems [\#1357](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1357) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `project_api_key` data source missing `project_assignment` attribute [\#1356](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1356) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support update of description for project api key resource [\#1354](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1354) ([AgustinBettati](https://github.com/AgustinBettati))
- Null pointer in `resource_mongodbatlas_cloud_backup_schedule` [\#1353](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1353) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Plugin did not respond - panic: interface conversion: interface is nil, not map[string]interface [\#1342](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1342) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Error deleting organization information when importing organization [\#1352](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1352) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Mark project api key resource as destroyed if not present [\#1351](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1351) ([AgustinBettati](https://github.com/AgustinBettati))
- `mongodbatlas_privatelink_endpoint_service` data source doc bug fix [\#1334](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1334) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Make region atributed optional computed in third-party-integration [\#1332](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1332) ([maastha](https://github.com/maastha))

**Closed Issues**

- json: cannot unmarshal number 4841168896 into Go struct field CloudProviderSnapshot.storageSizeBytes of type int [\#1333](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1333)
- Labels is not creating tags in the MongoAtlas UI [\#1319](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1319)
- `mongodbatlas_online_archive` `schedule` parameter update causing crashing in `terraform apply` [\#1318](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1318)
- Update Pager Duty integration fails with INTEGRATION\_FIELDS\_INVALID [\#1316](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1316)
- mongodbatlas\_event\_trigger is not updated if config\_match is added [\#1302](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1302)
- Updating the 'name' field of a 'mongodbatlas\_project' recreates a new Project [\#1296](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1296)
- mongodbatlas\_org\_invitation is missing ORG\_BILLING\_READ\_ONLY role support [\#1280](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1280)
- mongodbatlas\_alert\_configuration notification microsoft\_teams\_webhook\_url is always updated [\#1275](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1275)
- Provider not destroying API keys [\#1261](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1261)
- Has `project_assignment` of `mongodbatlas_api_key` not been implemented? [\#1249](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1249)
- Invalid attribute providerBackupEnabled specified. [\#1245](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1245)

**Internal Improvements**

- Fix documentation for `mongodbatlas_api_key` [\#1363](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1363) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Improve self-managed x509 database user docs [\#1336](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1336) ([andreaangiolillo](https://github.com/andreaangiolillo))
- add prefix to dependabot PR [\#1361](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1361) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update README.md with supported OS/Arch [\#1350](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1350) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add PR lint to repo [\#1348](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1348) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Mark `instance_size` in electable specs required in `advanced_cluster` documentation [\#1339](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1339) ([maastha](https://github.com/maastha))
- Update RELEASE.md Github issue [\#1331](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1331) ([maastha](https://github.com/maastha))
- Update privatelink endpoint service resources timeout config [\#1329](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1329) ([maastha](https://github.com/maastha))
- Use go-version-file in github actions [\#1315](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1315) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add autogenerated SDK to terraform [\#1309](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1309) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrate to terraform-plugin-testing [\#1301](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1301) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bump github.com/aws/aws-sdk-go from 1.44.308 to 1.44.314 [\#1360](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1360) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.43.10 to 0.43.11 [\#1358](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1358) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.3.0 to 1.4.0 [\#1346](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1346) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.43.8 to 0.43.10 [\#1345](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1345) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.304 to 1.44.308 [\#1344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1344) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.302 to 1.44.304 [\#1335](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1335) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.299 to 1.44.302 [\#1330](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1330) ([dependabot[bot]](https://github.com/apps/dependabot))


## [v1.10.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.2) (2023-07-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.1...v1.10.2)

**Bug Fixes:**

- `mongodbatlas_advanced_cluster` doc is not formatted correctly [\#1326](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1326) - INTMDB-941 ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_event_trigger` is not updated if `config_match` is added [\#1305](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1305) - INTMDB-922 ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_online_archive` `schedule` parameter update causing crashing in terraform apply - INTMDB-935 [\#1320](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1320) ([andreaangiolillo](https://github.com/andreaangiolillo))

**Internal Improvements:**

- Fix `mongodbatlas_online_archive` tests - INTMDB-938 [\#1321](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1321) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.298 to 1.44.299 [\#1312](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1312) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.10.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.1) (2023-7-13)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.0...v1.10.1)

**Enhancements:**

- Support for updating the name field of `mongodbatlas_project` without recreating a new Project - INTMDB-914 [\#1298](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1298) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `federation_settings_id` parameter in `mongodbatlas_organization` to enable linking to an existing federation upon Create - INTMDB-838 [\#1289](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1289) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `schedule` parameter in resource `mongodbatlas_online_archive` - INTMDB-828 [\#1272](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1272) ([andreaangiolillo](https://github.com/andreaangiolillo))
- New `mongodbatlas_advanced_cluster` doc examples for Multi-Cloud Clusters and Global Clusters - INTMDB-442 [\#1256](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1256) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `transactionLifetimeLimitSeconds` parameter in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluser` - INTMDB-874 [\#1252](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1252) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add Troubleshooting.md to include issue with using `dynamic` in Terraform - INTMDB-855 [\#1240](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1240) ([maastha](https://github.com/maastha))

**Bug Fixes:**
- Remove default value to [retainBackups parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Clusters/operation/deleteCluster) in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` - INTMDB-932 [\#1314](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1314) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_cloud_backup_snapshot_restore_job` extend guards for delivery type deletions - INTMDB-919 [\#1300](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1300) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_org_invitation` is missing `ORG_BILLING_READ_ONLY` role - INTMDB-904 [\#1287](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1287) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Change warning to error for org key delete - INTMDB-889 [\#1283](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1283) ([martinstibbe](https://github.com/martinstibbe))
- Add MicrosoftTeamsWebhookURL to values that are based on schema vs API - INTMDB-896 [\#1279](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1279) ([martinstibbe](https://github.com/martinstibbe))
- Update `group_id` -\> `project_id` for backup snapshots DOCSP-30798 [\#1273](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1273) ([jwilliams-mongo](https://github.com/jwilliams-mongo))
- Update example documentation for `mongodbatlas_project_api_key` - INTMDB-876 [\#1265](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1265) ([martinstibbe](https://github.com/martinstibbe))
- Make sure failed Terraform run rolls back properly - INTMDB-433 [\#1264](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1264) ([martinstibbe](https://github.com/martinstibbe))
- Fix "pause" during cluster and `mongodbatlas_advanced_cluster` update [\#1248](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1248) ([maastha](https://github.com/maastha))
- Add ForceNew to audit if the project id changes - INTMDB-435 [\#1247](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1247) ([martinstibbe](https://github.com/martinstibbe))

**Closed Issues:**

- `mongodbatlas_alert_configuration` failing to be created on apply due to `METRIC_TYPE_UNSUPPORTED` 400 error [\#1242](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1242)
- `app_id` property doesn't work in `mongodbatlas_event_trigger` resource [\#1224](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1224)
- Random change in `region_configs` order of `mongodbatlas_advanced_cluster` [\#1204](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1204)
- Problem returning some fields from `mongodbatlas_advanced_cluster` [\#1189](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1189)
- `node_count` documented as read-only for `mongodbatlas_advanced_cluster` [\#1187](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1187)
- In `mongodbatlas_third_party_integration` the `microsoft_teams_webhook_url` parameter keeps updating on every apply [\#1135](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1135)
- Point in Time Restore is not enabled when `should_copy_oplogs` is set to `true`, when copying backups to other regions [\#1134](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1134)
- Documentation: `analyzer` argument in Atlas search index is required [\#1132](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1132)
- Serverless Instance wants to do an in-place update on every run [\#1070](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1070)

**Internal Improvements:**

- INTMDB-912: Generate the CHANGELOG.md [\#1307](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1307) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Use GET one for cloud provider access to improve existing workflow - INTMDB-137 [\#1246](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1246) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.7 to 0.43.8 [\#1306](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1306) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.295 to 1.44.298 [\#1304](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1304) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.6 to 0.43.7 [\#1303](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1303) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-913: \[Terraform\] Enable fieldalignment linter [\#1297](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1297) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.294 to 1.44.295 [\#1293](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1293) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-911: \[Terraform\] Remove unused secret from code-health workflow [\#1291](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1291) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-910: \[Terraform\] Remove Automated Tests workflow [\#1290](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1290) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.291 to 1.44.294 [\#1286](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1286) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.26.1 to 2.27.0 [\#1285](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1285) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.5 to 0.43.6 [\#1284](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1284) ([dependabot[bot]](https://github.com/apps/dependabot))
- Remove slack key from repo [\#1282](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1282) ([maastha](https://github.com/maastha))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.3 to 0.43.5 [\#1277](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1277) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.287 to 1.44.291 [\#1276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1276) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-895: Third\_Party\_Integrations region field is required parameter in Terraform [\#1274](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1274) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Update RELEASING.md [\#1271](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1271) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-881: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - "Acceptance Tests" [\#1269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1269) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-892: \[Terraform\] Add APIx-Integration as a reviewer of dependabot PR [\#1268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1268) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.286 to 1.44.287 [\#1267](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1267) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.2 to 0.43.3 [\#1266](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1266) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-879: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Config [\#1263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1263) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-878: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Network [\#1260](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1260) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.284 to 1.44.286 [\#1259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1259) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.0 to 0.43.2 [\#1258](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1258) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-880: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Project [\#1257](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1257) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-883: Fix "Create JIRA ticket" Action [\#1255](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1255) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update README.md [\#1254](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1254) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update RELEASING.md [\#1253](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1253) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.283 to 1.44.284 [\#1251](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1251) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-871: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - backup [\#1250](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1250) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update acceptance-tests.yml [\#1244](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1244) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.282 to 1.44.283 [\#1243](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1243) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.272 to 1.44.282 [\#1237](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1237) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump golangci/golangci-lint-action from 3.4.0 to 3.6.0 [\#1235](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1235) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-843: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev. Clusters tests [\#1234](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1234) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-249: \[Terraform\] Lint and fix linting for examples [\#1221](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1221) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-808: Using vars instead of secrets for not sensitive info [\#1220](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1220) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump octokit/request-action from 2.1.7 to 2.1.9 [\#1211](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1211) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.2 to 2.17.0 [\#1206](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1206) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.10.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.0) (2023-6-15)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.9.0...v1.10.0)

**Enhancements:**

- New: [Organizations Management](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/createOrganization) including Create (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1176) - INTMDB-533  
- New: [Federated Database Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createFederatedDatabase) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1163) - INTMDB-801
- New: [Query Limit for Database Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createOneDataFederationQueryLimit) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1173) - INTMDB-802
- New: [Private Endpoint](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) resources and data sources for Federated Database Instance and Online Archive (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1182) - INTMDB-803
- New: [Data Lake Pipelines](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/createPipeline) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1174) - INTMDB-804
- New: [Data Lake Pipelines Run](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/getPipelineRun) data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1177) - INTMDB-805
- New: [Cluster Outage Simulation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation/operation/startOutageSimulation) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1188) - INTMDB-835
- Feature Add: [Cluster Enable Extended Storage](https://www.mongodb.com/docs/atlas/customize-storage/#minimum-disk-capacity-to-ram-ratios) Sizes in `mongodbatlas_project` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1128) - INTMDB-466
- Feature Add: [srvShardOptimizedConnectionString parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Multi-Cloud-Clusters/operation/createCluster) to `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` Data Sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1157) - INTMDB-694
- Feature Add: [retainBackups parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Multi-Cloud-Clusters/operation/createCluster) to `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1210) - INTMDB-781
- [Programmatic API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys) Resource Updates (`mongodbatlas_api_key`, `mongodbatlas_project_api_key` and `mongodbatlas_project_ip_access_list_key`) + Doc Cleanup (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1208) - INTMDB-655
- Release.md File Updates with Action Items (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1203) - INTMDB-690
- ChangeLog Generator (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1165) - INTMDB-720
- Upgrade to [Go 1.20](https://go.dev/blog/go1.20) (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1153) - INTMDB-658

**Deprecations and Removals:**

- `mongodbatlas_data_lake` and `mongodbatlas_privatelink_endpoint_service_adl` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1190) - INTMDB-806
-  Remove `mongodbatlas_private_ip_mode` and NEW_RELIC and FLOWDOCK in `mongodbatlas_third_party_integration` resources and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1159) - INTMDB-408
-  Remove mongodbatlas_cloud_provider (access, snapshot, snapshot_backup_policy, snapshot_restore_job) resources and datas sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1159) - INTMDB-408

**Bug Fixes:**

- `mongodbatlas_serverless_instance` wants to do an in-place update on every run (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1152) - INTMDB-710
- Documentation bug: analyzer argument in `mongodbatlas_search_index` is required (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1158) - INTMDB-780
- Point in Time Restore is not enabled when should_copy_oplogs is set to true, when copying backups to other regions (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1150) - INTMDB-783
- `mongodbatlas_third_party_integration` - microsoft_teams_webhook_url keeps updating on every apply (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1148) - INTMDB-784
- In `mongodbatlas_database_user` usernames with spaces breaks state due to URL encoding (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1164) - INTMDB-809
- `mongodbatlas_backup_compliance_policy` causing `mongodbatlas_cloud_backup_schedule` resource to fail (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1209) - INTMDB-827 
- `mongodbatlas_advanced_cluster` `node_count` parameter doc bug fix (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1193) - INTMDB-844
- Fix typos in docs for `mongodbatlas_network_peering` and  `mongodbatlas_network_container` resource imports (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1200) 

**Closed issues:**

- Online Archive: "Specified cloud provider is not supported in Data Lake Storage" but works in UI [\#1216](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1216)
- Asymmetric hardware is not supported by the provider [\#1214](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1214)
- `region` argument missing from `mongodbatlas_third_party_integration` for integration with PagerDuty [\#1180](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1180)
- Correct docs for importing network peering resources [\#1179](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1179)
- Terraform destroy produces a 500 \(UNEXPECTED\_ERROR\) on the underlying API call [\#1162](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1162)
- produced an unexpected new value: Root resource was present, but now │ absent [\#1160](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1160)
- Failed to respond to the plugin.\(\*GRPCProvider\).PlanResourceChange call [\#1136](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1136)
- Error: error creating MongoDB Cluster: unexpected EOF [\#674](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/674)

**Internal Improvements:**

- Chore\(deps\):  Bump actions/stale from 7 to 8 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1117)
- Chore\(deps\):  Bump github.com/zclconf/go-cty from 1.13.1 to 1.13.2 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1184)
- Chore\(deps\):  Bump github.com/aws/aws-sdk-go from 1.44.268 to 1.44.272 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1196)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.42.0 to 0.43.0 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1197)
- Chore(deps): Bump github.com/spf13/cast from 1.5.0 to 1.5.1 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1195)
- Chore(deps): Bump go.mongodb.org/atlas from 0.25.0 to 0.28.0 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1194)
- corrected documentation for advanced cluster and cluster (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1167) 
- Update component field to "Terraform" GitHub action (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1178)
- Add action to create JIRA ticket for a new Github issue (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1166)
- Terraform Provider Secrets Audit (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1202) 
- Add Code-health action (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1144)  
- Fix TestAccConfigDSSearchIndexes_basic (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1227) 

## [v1.9.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.9.0) (2023-4-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.2...v1.9.0)

**Enhancements:**

- New Feature: [Backup Compliance Policy](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/) Support  [\#1127](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1127) - INTMDB-587

**Bug Fixes:**

- Update resource [mongodbatlas_project](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project) to detect null boolean values in project settings [\#1145](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1145) - INTMDB-789
- Update on resource [mongodbatlas_search_index]([https://www.mongodb.com/docs/atlas/atlas-search/create-index/](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index)) resource docs [\#1137](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1137) - DOCSP-28948
- Removing resource [mongodbatlas_cluster](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) `disk_size_gb` examples [\#1133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1133)

**Closed issues:**

- mongodbatlas\_search\_index does not change name [\#1096](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1096)
- Unhelpful error when importing advanced cluster using mongodbatlas\_cluster resource [\#1089](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1089)
- Update Slack alert configuration fails with INTEGRATION\_FIELDS\_INVALID [\#1086](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1086)
- Upgrade to terraform-plugin-sdk v2.25.0 [\#1080](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1080)
- mongodbatlas\_project\_ip\_access\_list.comment should be optional [\#1079](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1079)
- Can't unset auto\_scaling [\#1072](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1072)
- mongodbatlas\_access\_list\_api\_key fails to import [\#1064](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1064)
- Terraform plan fails if API key created by `mongodbatlas_api_key` resource is deleted outside of Terraform [\#1057](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1057)
- mongodbatlas\_search\_index does not recreate when cluster\_name and project\_id fields change. [\#1053](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1053)
- ERROR: Pager Duty API key must consist of 32 hexadecimal digits [\#1049](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1049)
- mongodbatlas\_alert\_configuration not detecting drift [\#999](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/999)
- Provider insists on changing a sub-parameter even when no changes are necessary [\#997](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/997)
- Type: TEAM alert notification not saved properly [\#971](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/971)
- `app_id` In documentation is ambiguous for MongoDB Atlas Event Trigger [\#957](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/957)
- Provider panic with `authentication_enabled=true` input [\#873](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/873)
- Schema error when creating event\_trigger referring to Atlas App Services function  [\#858](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/858)

## [v1.8.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.2) (2023-3-30)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.1...v1.8.2)

**Enhancements:**

- Support for "TIMESERIES" Collection Type in [`mongodbatlas_online_archive`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/online_archive) [\#1114](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1114) - INTMDB-648
- Support for new "DATADOG" regions  in [`mongodbatlas_third_party_integration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/third_party_integration) [\#1105](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1105) - INTMDB-638

**Bug Fixes:**

- Error in unsetting auto_scaling in [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) [\#1112](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1112) - INTMDB-617
- Update Status Codes in  [`mongodbatlas_search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) [\#1104](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1104) - INTMDB-687
- [`mongodbatlas_project_ip_access_list`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_ip_access_list) comment should be optional [\#1103](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1103) - INTMDB-637
- Plan fails if API key created by [`mongodbatlas_api_key`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/api_key) resource is deleted outside of Terraform [\#1097](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1097) - INTMDB-581
- Google Cloud Terraform Provider Test Version Upgrade Refactoring [\#1098](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1098) - INTMDB-359


**Closed Issues:**

- mongodbatlas\_federated\_settings\_org\_role\_mapping INVALID\_ATTRIBUTE [\#1110](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1110)
- Errors when creating or importing timeseries online archive [\#1081](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1081)
- Unexpected EOF [\#1083](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1083)
- mongodbatlas\_access\_list\_api\_key creation fails after api\_key creation [\#1075](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1075)
- Panic when creating AWS privatelink endpoint [\#1067](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1067)
- error getting project 403 \(request IP\_ADDRESS\_NOT\_ON\_ACCESS\_LIST\) even if whitelisted IP [\#1048](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1048)
- Crash during the changes for `mongodbatlas_cloud_backup_schedule` interval with dynamic blocks. [\#1041](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1041)
- Regression in 1.8.0: mongodbatlas\_third\_party\_integration marks "type" attribute as deprecated [\#1032](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1032)
- Cannot reference `container_id` for `mongodbatlas_advanced_cluster` [\#1008](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1008)
- ERROR: Datadog API key must consist of 32 hexadecimal digits [\#1001](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1001)
- num\_shards value changed to 1 [\#970](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/970)
- Segmentation Fault in TerraForm Provider [\#969](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/969)
- Import statements are broken in documentation and help commands in the Terraform provider are outdated. [\#956](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/956)
- Invitation handling is not working after user accepted invitation. [\#945](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/945)
- bi\_connector settings don't work in mongodbatlas\_advanced\_cluster [\#893](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/893)
- Sensitive information exposed as resource Id - mongodbatlas\_x509\_authentication\_database\_user  [\#884](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/884)
- plugin crashes during apply: panic: runtime error: invalid memory address or nil pointer dereference [\#866](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/866)

**Internal Improvements:**

- Release staging v.1.8.2 [\#1115](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1115) ([martinstibbe](https://github.com/martinstibbe))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.221 to 1.44.226 [\#1109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1109) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.25.0 to 2.26.1 [\#1108](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1108) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/zclconf/go-cty from 1.12.1 to 1.13.1 [\#1107](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1107) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump actions/setup-go from 3 to 4 [\#1106](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1106) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.1 to 2.16.2 [\#1101](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1101) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.24.1 to 2.25.0 [\#1100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1100) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.216 to 1.44.221 [\#1099](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1099) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.40.56 to 1.44.216 [\#1094](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1094) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump go.mongodb.org/atlas from 0.21.0 to 0.23.1 [\#1092](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1092) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump golang.org/x/net from 0.1.0 to 0.7.0 [\#1071](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1071) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.8.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.1) (2023-3-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.0...v1.8.1)

**Enhancements:**

- Upgrade to [go1.19](https://go.dev/blog/go1.19) [\#1031](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1031) - INTMDB-390
- Add configurable timeouts to resources that wait for [`cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) to become IDLE [\#1047](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1047) - INTMDB-508
- Improve [`cloud_provider_access_authorization`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_provider_access) and [`encryption_at_rest`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/encryption_at_rest) (remove need for `time_sleep` arguments) [\#1045](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1045) - INTMDB-560
- Add [`search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) error handling [\#1077](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1077) - INTMDB-600
- New / Improved Upon [Resource Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples): 
  - [`cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) with NVMe Upgrade ([Non-Volatile Memory Express](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ssd-instance-store.html)) [\#1037](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1037) - INTMDB-32. See example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cluster)
  - [`privatelink_endpoint_serverless`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_serverless) Examples for AWS + Azure [\#1043](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1043) - INTMDB-424. See example for [AWS](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-atlas-privatelink-serverless) and [Azure](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/azure-atlas-privatelink-serverless)
  - Improvement for [`private_link_endpoint`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint) [\#1082](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1082) - INTMDB-410. see example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-privatelink-endpoint)
  - Improvement for [`encryption_at_rest`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/encryption_at_rest) [\#1060](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1060). see example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-encryption-at-rest/aws)

**Bug Fixes:**

- Resource [`ldap_configuration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/ldap_configuration) broken [\#1033](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1033) - INTMDB-440
- [`event_trigger`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/event_trigger) Import statements are broken in documentation [\#1046](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1046) - INTMDB-513
- [`event_trigger`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/event_trigger) Error Handler Update [\#1061](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1061) - INTMDB-517
- TEAM alert notification not saved properly [\#1029](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1029) - INTMDB-529
- [`alert_configuration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/alert_configuration) not detecting drift [\#1030](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1030) - INTMDB-542
- [`third_party_integration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/third_party_integration) marks "type" attribute as deprecated erroneously [\#1034](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1034) - INTMDB-556
- Error "Pager Duty API key must consist of 32 hexadecimal digits" [\#1054](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1054) - INTMDB-570
- Terraform provider stuck in changes for [`advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) configuration [\#1066](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1066) - INTMDB-572
- [`search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) does not recreate when `cluster_name` and `project_id` fields change [\#1078](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1078) - INTMDB-576
- POST Create Access List Entries for One Organization API Key endpoint supports list, but Terraform does not [\#1065](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1065) - INTMDB-579
- Typo in Readme [\#1073](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1073)
- Update project_api_key.html.markdown [\#1044](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1044)
- Doc Clean Up [`cloud_provider_access`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_provider_access) [\#1035](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1035)
- Update alert_configuration.html.markdown [\#1068](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1068)
- Doc Clean Up `provider_backup_enabled` deprecated to `cloud_backup` [\#1036](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1036)

**Closed Issues:**

- Unable to create third party integration of type Datadog with version 1.8.0 [\#1038](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1038)
- mongodbatlas\_third\_party\_integration - api\_token keeps updating on every apply [\#963](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/963)
- Unable to add mongodbatlas provider to CDK [\#952](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/952)
- \[Bug\] `update_snapshots` doesn't save at TF state with [`cloud_backup_schedule`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_schedule) resource [\#904](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/904)
- Cannot ignore changes for replication\_specs when [`auto_scaling`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#region_configs) is enabled [\#888](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/888)
- Warning: Deprecated attribute [\#1042](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1042)
- GCP Network Peering remains pending when created via terraform [\#917](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/917)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.0 to 2.16.1 [\#1062](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1062) ([dependabot[bot]](https://github.com/apps/dependabot))
- Update access\_list\_api\_key.html.markdown [\#1058](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1058) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.15.0 to 2.16.0 [\#1055](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1055) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.9 to 0.41.10 [\#1051](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1051) ([dependabot[bot]](https://github.com/apps/dependabot))
- Update CODEOWNERS to use APIx-Integrations [\#1050](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1050) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 [\#1026](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1026) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.8.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.0) (2023-1-26)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.7.0...v1.8.0)

**Enhancements:**

- Snapshot Distribution Support [\#979](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/979) - INTMDB-400
- Programmatically Create API Keys [\#974](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/974) - INTMDB-346
- Retrieve `org_id` from API Keys [\#973](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/973) - INTMDB-454
- `oplogMinRetentionHours` Parameter Support in `advanced_cluster` and `cluster` resources [\#1016](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1016) - INTMDB-397
- Analytics Node Tier New Features Support [\#994](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/994) - INTMDB-488
- Improve Default Alerts and Example Creation [\#993](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/993) - INTMDB-382
- Expand documentation for `cloud_backup_schedule` to include information about valid values for `frequency_interval` [\#1007](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1007) - INTMDB-547

**Depreciations:**

- `cloud_provider_snapshot`, `cloud_provider_snapshot_backup_policy`, `cloud_provider_snapshot_restore_job`, and `private_ip_mode` are now deprecated and will be removed from codebase as of v1.9 release [\#988](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/988) - INTMDB-409
- `NEW_RELIC` and `FLOWDOCK` in `third_party_integration` resource are now deprecated and will be removed from codebase as of v1.9 release [\#989](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/989) - INTMDB-482

**Bug Fixes:**

- Hide `current_certificate` when X.509 Authentication Database Users are Created [\#985](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/985) - INTMDB-468
- Import example added for `encryption_at_rest` resource [\#992](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/992) - INTMDB-530 
- Resource `cloud_backup_snapshot_export_job` variable name change [\#976](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/976) - INTMDB-523
- `update_snapshot` doesn't save at TF state with `cloud_backup_schedule` resource fix [\#974](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/974) - INTMDB-472
- Invitation handling after user accepts invitation fix [\#1012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1012) - INTMDB-511 
- `bi_connector` settings in `advanced_cluster` fix (breaking changes) [\#1010](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1010) - INTMDB-455
- `third_party_integration` api_token keeps updating on every apply fix [\#1011](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1011) - INTMDB-519
- `custom_db_role` error fix [\#1009](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1009) - INTMDB-448
- `ldap_configuration` and `ldap_verify` resources fix [\#1004](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1004) - INTMDB-543
- `cloud_backup_schedule` resource fix [\#968](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/968) - INTMDB-427
- `search_index_test` fix [\#964](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/964) - INTMDB-341
- Cannot ignore changes for replication_specs when autoscaling enabled fix [\#961](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/961) - INTMDB-464
- BI Connector documentation fix [\#1017](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1017) 
- `federated_settings_org_config` import example fix [\#996](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/996)

**Closed Issues:**

- Documentation: Expand documentation for mongodbatlas\_cloud\_backup\_schedule to include information about valid values for frequency\_interval  [\#1005](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1005)
- Serverless instance returns incorrect connection string [\#934](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/934)
- Terraform apply failed with Error: Provider produced inconsistent final plan This is a bug in the provider, which should be reported in the provider's own issue tracker. [\#926](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/926)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.7 to 0.41.9 [\#1013](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1013) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.7.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.7.0) (2023-1-16)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.6.1...v1.7.0)

**Enhancements:**

- AWS Secrets Manager (AWS SM) Authetication for Terraform Atlas Provider [\#975](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/975) - INTMDB-521 

**Bug Fixes:**

- Resource cloud_backup_snapshot_export_job variable name change [#976](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/976) - INTMDB-523
- Deprecate legacy mongodbatlas.erb given Terraform Registry autogeneration [#962](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/962) - INTMDB-477  

**Closed Issues:**

- Terraform plan fail: Asymmetric hardware is not supported by the v1.0 API [\#958](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/958)
- Error importing ressource mongodbatlas\_network\_peering.mongo\_peer [\#906](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/906)
- \[Bug\] `container_id` is unconfigurable  attribute at `mongodbatlas_advanced_cluster` resource [\#890](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/890)
- mongodbatlas\_alert\_configuration - api\_token keeps wanting to change [\#863](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/863)
- Docs - Example - Return a Connection String - Azure Private Endpoint [\#713](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/713)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.6 to 0.41.7 [\#978](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/978) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump actions/stale from 6 to 7 [\#977](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/977) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.4 to 0.41.6 [\#967](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/967) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/go-test/deep from 1.0.8 to 1.1.0 [\#966](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/966) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump goreleaser/goreleaser-action from 3 to 4 [\#965](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/965) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.6.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.6.1) (2022-12-6)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.6.0...v1.6.1)

**Enhancements:**

- Enable Adv Cluster and Cluster to have configurable timeouts [\#951](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/951) - INTMDB-503 
- Updated Prometheus Example [\#942](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/942) - INTMDB-498 
- Auto-Generate Changelog [\#944](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/944) - INTMDB-478 

**Bug Fixes:**

- Alert Configuration -- Api Token erroneous changes [\#941](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/941) - INTMDB-426
- Fix example private endpoint called out in issue 713 [\#907](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/907) - INTMDB-434 
- Cluster rename is inconsistently rejected by Terraform [\#929](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/929) - INTMDB-438 
- Terraform does not wait for cluster update when creating GCP private endpoints [\#943](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/943) - INTMDB-465
- mongodbatlas_federated_settings_connected_organization customer HELP + doc bug [\#924](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/924) - INTMDB-481  
- Serverless Private Endpoint Connection String Example Fix [\#940](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/940) - INTMDB-493 
- Fix regional mode endpoint test [\#946](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/946) - INTMDB-470 
- Skip tests for OPS GENIE and GOV [\#937](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/937) - INTMDB-484 
- Test Instability around privatelink tests fix [\#895](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/895) - INTMDB-384 
- Shorten test names that are too long to allow for targeting specific tests [\#932](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/932) - INTMDB-368 
- Remove container_id from configurable attribute in advanced_cluster [\#931](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/931) - INTMDB-463 

**Closed Issues:**

- No documented way to get config out of third party integration [\#939](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/939)
- Double checking Terraform Plan before I destroy Production [\#938](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/938)
- Issue: MongoDB Atlas Plugin Failure v 1.5.0 [\#928](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/928)
- MILLION\_RPU unit isn't supported by provider [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/854)
- MS Teams alert support in terraform provider is missing [\#827](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/827)

**Internal Improvements:**

- v1.6.1 - Conditionally ignore serverless connection string changes [\#953](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/953) ([evertsd](https://github.com/evertsd))
- Swap logic for variable substitution [\#950](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/950) ([martinstibbe](https://github.com/martinstibbe))
- Fix serverless endpoint tests [\#949](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/949) ([evertsd](https://github.com/evertsd))
- Release staging v1.6.1 [\#947](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/947) ([martinstibbe](https://github.com/martinstibbe))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.0 to 0.41.3 [\#936](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/936) ([dependabot[bot]](https://github.com/apps/dependabot))
- Serverless Endpoint Service Doc Bug [\#930](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/930) ([Zuhairahmed](https://github.com/Zuhairahmed))

## [v1.6.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.6.0) (2022-11-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.5.0...v1.6.0)

**Enhancements:** 

- Termination Protection for Advanced Cluster/Cluster/Serverless Instances [\#912](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/912) - INTMDB-444
- AWS/Azure Serverless Private Endpoints [\#913](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/913) - INTMDB-364

**Internal Improvements:**

- docs(website): fix federated_settings_org_config resource name by removing the misleading trailing s [\#908](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/908)
- chore(github): add link to contribution guidelines in PR template [\#910](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/910#issuecomment-1310007413)
- docs(resource/role_mapping): indent sub-elements of role_assignments [\#918](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/918)
- docs(resource/role_mapping): add link to reference of available role IDs [\#919](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/919)
- federated settings plural fix [\#914](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/914)  
- Chore(deps): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.24.0 to 2.24.1 [\#922](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/922)
- Chore(deps): Bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 [\#925](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/925)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.24 to 0.41.0 [\#923](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/923)
- Workaround to handle serverless endpoint tests failing due to provider name missing from API [\#927](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/927)
- Release staging v1.6.0 [\#921](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/921)

## [v1.5.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.5.0) (2022-11-07)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.6...v1.5.0)

**Enhancements:** 
- INTMDB-224 - Support AtlasGov with Terraform [\#865](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/865)
- INTMDB-314 - Feature add: Add ability to upgrade shared/TENANT tiers for clusters and advanced clusters [\#874](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/874)
- INTMDB-349 - New AtlasGov parameter to tell the provider to use the Atlas gov base URL [\#865](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/865)
- INTMDB-373 - Add new notification parameters to the mongodbatlas_alert_config resource [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)	
- INTMDB-378 - Document for users how to get a pre-existing container id [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)
- INTMDB-377 - Release 1.5 (both pre and then GA) [\#887](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/887)
- INTMDB-394 - MS Teams alert support [\#320](https://github.com/mongodb/go-client-mongodb-atlas/pull/320)

**Bug Fixes:** 
- INTMDB-326 - Review code/tests and docs for resource_mongodbatlas_search_index [\#891](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/891)
- INTMDB-334 - Determine best path forward for GCP PSC timeouts and implement [\#859](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/859)
- INTMDB-370 - Docs not complete for cloud_backup_snapshot_restore_job	[\#870](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/870)
- INTMDB-403 - Update third_party_integration.markdown [\#851](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/851)
- INTMDB-405 - Add cluster label to advanced clusters	[\#857](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/857)
- INTMDB-406 - MILLION_RPU unit isn't supported by provider #854 [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/854)
  
**Closed Issues:**
- MS Teams alert support in terraform provider is missing [\#827](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/827)
- MILLION_RPU unit isn't supported by provider not_stale  [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/854)

**Merged Pull Requests:**
- INTMDB-378: Add link for How To Guide for existing container ID [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)
- INTMDB-403: Update third_party_integration.markdown [\#851](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/851)
- INTMDB-404: mongodbatlas_advanced_cluster doc updates [\#852](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/852)
- INTMD-428: doc update to "mongodbatlas_projects" [\#869](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/869)
- Docs: fix custom_dns_configuration_cluster_aws [\#860](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/860)
- Relying on atlas api for unit validation on alert configuration [\#862](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/862)
- Adding a github actions to automatically close stale issues/PRs based on CLOUDP-79100 [\#872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/872)
- Encryption_at_rest M10+ limit doc update [\#886](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/886)
- Update cluster.html.markdown [\#878](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/878)
- Add APIx1 CODEOWNER [\#894](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/894)
- Chore(deps): Bump octokit/request-action from 2.1.6 to 2.1.7 [\#868](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/868)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.22 to 0.40.24 [\#875](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/875)
- Chore(deps): Bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 [\#897](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/897)

## [v1.4.6](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.6) (2022-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.5...v1.4.6)

**Enhancements and Bug Fixes:** 
- INTMDB-387 - Enable Azure NVME for Atlas Dedicated clusters [\#833](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/833)
- INTMDB-342 - Update TestAccDataSourceMongoDBAtlasPrivateEndpointRegionalMode_basic test to use a new project to prevent conflicts  [\#837](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/837)
- INTMDB-347 - Cloud_backup is not being correctly imported - issue [\#768](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/768)
- INTMDB-354 - Update docs around what requires an API key access list [\#834](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/834)
- INTMDB-363 - [Updated Feature] Add serverless backup to mongodbatlas_serverless_instance [\#830](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/830)
- INTMDB-379 - Release 1.4.6 (both pre and then GA)	
- INTMDB-381 - Customer is unable to disable backup auto export [\#823](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/823)
- INTMDB-383 - Update the BYOK/encryption at rest resource [\#805](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/805)
- INTMDB-385 - use_org_and_group_names_in_export_prefix is not working for a customer
- INTMDB-386 - Add new role types to invitation verification	[\#840](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/840)
- INTMDB-371 - Timeout when creating privatelink_endpoint [\#806](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/806)
- INTMDB-372 - Fix failing test for testAccMongoDBAtlasAlertConfigurationConfigWithMatchers	[\#836](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/836)
- INTMDB-358 - Upgrade to go1.18 [\#835](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/835)
- INTMDB-391 - Doc Fix for teams.html.markdown [\#838](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/838)

**Closed Issues:**
-  importing existing cluster does not populate backup status #768 [\#768](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/768)


**Merged Pull Requests:**
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.21 to 0.40.22 [\#842](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/842) ([dependabot[bot]](https://github.com/apps/dependabot))

- Rename team.html.markdown into teams.html.markdown [\#838](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/838) ([leo-ferlin-sutton](https://github.com/leo-ferlin-sutton))

- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.20 to 0.40.21 [\#825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/825) ([dependabot[bot]](https://github.com/apps/dependabot))

- Fix create index error msg[\#824](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/824) ([ebouther](https://github.com/ebouther))


## [v1.4.5](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.5)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.4...v1.4.5)

**Fixed**

- INTMDB-369: Fix parsing of `delivery_type_config` when using `point_in_time` for `cloud_backup_snapshot_restore_job`, in [\#813](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/813)
- INTMDB-322: Validated serverless alert_configurations and improved documentation on usage, addressing issue [\#722](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/722) in [\#819](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/819)

## [v1.4.4](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.4) (2022-08-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.3...v1.4.4)

**Fixed**

- INTMDB320 - Fix Global Cluster import documentation, in [\#796](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/796)
- INTMDB-331 - Update GCP documentation, issue [\#753](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/753), in [\#793](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/793)
- INTMDB-351 - Project data_source reads name, issue [\#788](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/788), in [\#795](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/795)
- INTMDB-362: Header Clarifications "Resource" vs "Data Source" in Documentation, in [\#803])(https://github.com/mongodb/terraform-provider-mongodbatlas/pull/803)
- INTMDB-343: Update go from 1.16 to 1.17 add windows arm64 build support, in [\#797](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/797)

## [v1.4.4-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.4-pre.1) (2022-08-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.3...v1.4.4-pre.1)

**Closed Issues:**

- Unable to update members in an existing "mongodbatlas\_teams" as the provider attempts to remove all users first [\#790](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/790)
- Please elaborate how to acquire PROJECTID and PEERINGID and PROVIDERNAME for import of network peering [\#789](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/789)
- error: error reading cloud provider access cloud provider access role not found in mongodbatlas, please create it first [\#781](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/781)

**Merged Pull Requests:**

- Update CONTRIBUTING.md [\#798](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/798) ([themantissa](https://github.com/themantissa))
- Fix federated\_settings\_identity\_provider attribute name [\#791](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/791) ([florenp](https://github.com/florenp))

## [v1.4.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.3) (2022-07-12)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.2...v1.4.3)

**Fixed:**

* INTMDB-335: Add option for multiple weekly monthly schedules @martinstibbe in [\#784](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/784)
* INTMDB-348: autoexport parameter not being set via provider [\#784](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/784)
* INTMDB-323: Removed the requirement to set `MONGODB_ATLAS_ENABLE_BETA` to use serverless and update the docs to match. [\#783](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/783)
* INTMDB-330 Fixed Serverless Instance Import Documentation. Closes [\#754](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/754)
* Fix typo in custom_db_role documentation [\#780](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/780)
* Fix typo in federated_settings_org_configs documentation [\#779](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/779)
 

## [v1.4.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.2) (2022-07-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.1...v1.4.2)

**Fixed:**

* INTMDB-313: Update project settings default flags by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/778

## [v1.4.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.1) (2022-07-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.0...v1.4.1)

**Fixed:**

* Update CHANGELOG.md by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/771
* INTMDB-313: Update project settings default flags by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/773


## [v1.4.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.0) (2022-07-5)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1...v1.4.0) 

**Closed Issues:**

Note: the binary executable for windows/arm64 is not available for this release.  Next release will include.
- Fix for Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/740)
- Feature Add: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/741)
- Fix for  Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/705)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/706)
- Fix for Correct import function for snapshot export bucket #714 [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/715)
- Fix for Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/717)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings
- Fix for Cannot import export bucket - bad state id encoding [\#708](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/708)
- Error missing expected { when updating the provider [\#697](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/697)

**Merged Pull Requests:**

- INTMDB-321: Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/740) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-313: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/741) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-301: Feature add: Add support for management of federated authentication configuration [\#742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/742) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-307: Add Regionalized Private Endpoint Settings [\#718](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/718) ([evertsd](https://github.com/evertsd))
- INTMDB-310: Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/705) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-311: Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/706) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-315: Correct import function for snapshot export bucket [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/715) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-309: Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/717) ([martinstibbe](https://github.com/martinstibbe))

## [v1.4.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.0-pre.1) (2022-06-29)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1...v1.4.0-pre.1) 

**Closed Issues:**

Note: the binary executable for windows/arm64 is not available for this release.  Next release will include.
- Fix for Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/740)
- Feature Add: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/741)
- Fix for  Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/705)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/706)
- Fix for Correct import function for snapshot export bucket #714 [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/715)
- Fix for Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/717)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings
- Fix for Cannot import export bucket - bad state id encoding [\#708](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/708)
- Error missing expected { when updating the provider [\#697](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/697)

**Merged Pull Requests:**

- INTMDB-321: Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/740) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-313: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/741) ([martinstibbe](https://github.com/martinstibbe)) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-301: Feature add: Add support for management of federated authentication configuration [\#742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/742) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-307: Add Regionalized Private Endpoint Settings [\#718](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/718) ([evertsd](https://github.com/evertsd))
- INTMDB-310: Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/705) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-311: Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/706) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-315: Correct import function for snapshot export bucket [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/715) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-309: Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/717) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1-pre.1) (2022-02-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.1...v1.3.1-pre.1)

**Closed issues:**
- Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686)
- Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) 
- Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687)
- Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695)

**Merged pull requests:**

- INTMDB-300: Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-302: Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) ([evertsd](https://github.com/evertsd))
-  INTMDB-303: Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-305: Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1) (2022-03-28)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.1...v1.3.1-pre.1)

- INTMDB-306: [Terraform] Release bug fix version 1.3.1
## [v1.3.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1-pre.1) (2022-02-23)

**Closed issues:**
- Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686)
- Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) 
- Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687)
- Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695)

**Merged pull requests:**

- INTMDB-300: Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-302: Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) ([evertsd](https://github.com/evertsd))
-  INTMDB-303: Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-305: Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.0) (2022-02-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.0-pre.1...v1.3.0)

**Merged pull requests:**

- Create 1.3.0-upgrade-guide.html.markdown [\#682](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/682) ([themantissa](https://github.com/themantissa))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.2 to 0.40.3 [\#681](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/681) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.3.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.0-pre.1) (2022-02-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.2.0...v1.3.0-pre.1)

**Closed issues:**

- Auto scaling of storage cannot be disabled for mongodbatlas\_advanced\_cluster via "disk\_gb\_enabled" [\#677](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/677)
- Can't create M0 free tier on Azure [\#675](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/675)
- Error: error creating MongoDB Cluster: unexpected EOF [\#674](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/674)
- attempting to assign api key to resource `mongodbatlas_project` results in `Error: Unsupported block type` [\#671](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/671)
- Error in documentation page : aws\_kms -\> aws\_kms\_config [\#666](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/666)
- Accepting organization invitation causes 404 [\#636](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/636)
- Alert configuration state is not stable [\#632](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/632)
- mongodbatlas\_cloud\_backup\_schedule with Azure results in "restore window days" mandatory [\#625](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/625)
- incorrect values from `mongodbatlas_cluster` and `mongodbatlas_clusters` datasources [\#618](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/618)
- mongodb atlas network container atlas cidr block value and real  value is not mached [\#617](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/617)
- config\_full\_document\_before not working in EventTriggers [\#616](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/616)
- Hitting timing issue when trying to integrate with `aws` provider's `aws_iam_access_key` resource [\#127](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/127)

**Merged pull requests:**

- INTMDB-291: pre-release v1.3.0 [\#680](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/680) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.1 to 0.40.2 [\#679](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/679) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.0 to 0.40.1 [\#678](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/678) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-299: Support Cloud Backup Export Jobs [\#673](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/673) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump octokit/request-action from 2.1.0 to 2.1.4 [\#672](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/672) ([dependabot[bot]](https://github.com/apps/dependabot))
- update the documentations for new changes [\#670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/670) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.9 to 0.40.0 [\#669](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/669) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-298: fixes a bug where you couldn't delete a team in team resource [\#668](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/668) ([coderGo93](https://github.com/coderGo93))
- INTMDB-297: set the container id to avoid null in state of data source cluster\(s\) [\#667](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/667) ([coderGo93](https://github.com/coderGo93))
- fixed typo in mongodbatlas\_teams sidebar [\#665](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/665) ([MartinCanovas](https://github.com/MartinCanovas))
- INTMDB-295: Fixes a bug about unauthorized error in project resource [\#664](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/664) ([coderGo93](https://github.com/coderGo93))
- INTMDB-293: Added container\_id in advanced cluster [\#663](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/663) ([coderGo93](https://github.com/coderGo93))
- INTMDB-294: \[Terraform\] Address security warnings from dependabot [\#661](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/661) ([thetonymaster](https://github.com/thetonymaster))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.8 to 0.38.9 [\#660](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/660) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-290: Added advanced configuration for datasource/resource of advanced cluster [\#658](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/658) ([coderGo93](https://github.com/coderGo93))
- Fix 1.2 upgrade/info guide formatting error [\#657](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/657) ([themantissa](https://github.com/themantissa))

## [v1.2.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.2.0) (2022-01-14)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.2.0-pre.1...v1.2.0)

**Merged pull requests:**

- INTMDB-268: Release v1.2.0 [\#656](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/656) ([abner-dou](https://github.com/abner-dou))

## [v1.2.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.2.0-pre.1) (2022-01-13)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.1...v1.2.0-pre.1)

**Closed issues:**

- mongodbatlas\_teams provides a Team resource [\#649](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/649)
- terraform-provider-mongodbatlas\_v0.8.0 plugin: panic: runtime error: invalid memory address or nil pointer dereference [\#644](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/644)
- backup snapshot and restore not working automated  [\#642](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/642)
- Delete default alerts [\#628](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/628)
- Cannot set multiple notifications for an alert [\#626](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/626)
- PRIVATE\_ENDPOINT\_SERVICE\_ALREADY\_EXISTS\_FOR\_REGION [\#590](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/590)
- INVALID\_CLUSTER\_CONFIGURATION when modifying a cluster to use replication\_specs \(eg for multi-region\) [\#588](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/588)

**Merged pull requests:**

- Re-branch from 651 due to conflict of docs fixes [\#654](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/654) ([themantissa](https://github.com/themantissa))
- INTMDB-268: Pre-release v1.2.0 [\#650](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/650) ([abner-dou](https://github.com/abner-dou))
- INTMDB-5: added parameter team name for alert configurations  [\#648](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/648) ([coderGo93](https://github.com/coderGo93))
- Fix markdown formatting in network\_container.html.markdown [\#647](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/647) ([pmacey](https://github.com/pmacey))
- INTMDB-15: Added parameter advanced conf for cluster datasource [\#646](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/646) ([coderGo93](https://github.com/coderGo93))
- INTMDB-284: Updated docs in alert configuration resource and datasource [\#645](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/645) ([coderGo93](https://github.com/coderGo93))
- INTMDB-285: Fix org\_invitations issue [\#643](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/643) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.9.0 to 2.10.1 [\#641](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/641) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-263: Create Resource and Datasource for private\_link\_endpoint\_adl [\#640](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/640) ([abner-dou](https://github.com/abner-dou))
- INTMDB-287: Fixes the issues in project api keys [\#639](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/639) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.5 to 0.38.8 [\#638](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/638) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-283: Fixes a bug about optional parameters in Cloud Backup Schedule [\#631](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/631) ([coderGo93](https://github.com/coderGo93))
- INTMDB-281: Fix realm event trigger issue [\#630](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/630) ([abner-dou](https://github.com/abner-dou))
- INTMDB-282: Updated test and docs for alert configuration using notifications [\#629](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/629) ([coderGo93](https://github.com/coderGo93))
- INTMDB-280: Fix cluster datasource scaling issue [\#627](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/627) ([abner-dou](https://github.com/abner-dou))
- INTMDB-272: Validate using interval\_min for PagerDuty, VictorOps, GenieOps [\#624](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/624) ([coderGo93](https://github.com/coderGo93))
- INTMDB-276: Added VersionReleaseSystem parameter for resource/datasource\(s\) of cluster [\#623](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/623) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.4 to 0.38.5 [\#622](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/622) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0 [\#621](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/621) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-24: Change computed paused to optional [\#620](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/620) ([coderGo93](https://github.com/coderGo93))
- INTMDB-257: Changed 'hcl' markdown tag to 'terraform' tag [\#619](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/619) ([abner-dou](https://github.com/abner-dou))
- Fix link in update guide and add version [\#615](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/615) ([themantissa](https://github.com/themantissa))
- INTMDB-279: Fixes a bug where it crashes when importing a trigger [\#614](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/614) ([coderGo93](https://github.com/coderGo93))
- mongodbatlas-project resource: add api\_keys attribute [\#504](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/504) ([tetianakravchenko](https://github.com/tetianakravchenko))

## [v1.1.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.1) (2021-11-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.0...v1.1.1)

**Closed issues:**

- Cannot update the default backup schedule policy without defining API Key access IPs. [\#610](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/610)

**Merged pull requests:**

- Release v1.1.1 [\#613](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/613) ([abner-dou](https://github.com/abner-dou))
- Fix documentation v1.1.0 [\#612](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/612) ([abner-dou](https://github.com/abner-dou))

## [v1.1.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.0) (2021-11-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.0-pre.1...v1.1.0)

**Merged pull requests:**

- INTMDB-264: Release v1.1.0 [\#611](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/611) ([abner-dou](https://github.com/abner-dou))
- Guide and minor main page changes for 1.1.0 [\#609](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/609) ([themantissa](https://github.com/themantissa))

## [v1.1.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.0-pre.1) (2021-11-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.2...v1.1.0-pre.1)

**Fixed bugs:**

- Plugin crash when changing value of iam\_assumed\_role\_arn for resource mongodbatlas\_cloud\_provider\_access\_authorization [\#565](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/565)
- terraform-provider-mongodbatlas\_v1.0.1 crashes after creating mongodbatlas\_cloud\_provider\_access\_authorization resource [\#554](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/554)
- INVALID\_CLUSTER\_CONFIGURATION when adding new regions [\#550](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/550)

**Closed issues:**

- update\_snapshots is not getting set on Atlas while using "mongodbatlas\_cloud\_backup\_schedule" tf resource  [\#594](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/594)
- Race condition when destroying cluster and disabling encryption at rest on the project-level [\#518](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/518)
- Scaling max/min is applied each time with disabled autoscaling [\#482](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/482)
- Documentation: Update contribution readme for developing the provider for terraform +14 newer [\#466](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/466)

**Merged pull requests:**

- INTMDB-264: pre-release v1.1.0 [\#608](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/608) ([abner-dou](https://github.com/abner-dou))
- INTMDB:273: Fix replication\_specs update error [\#607](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/607) ([abner-dou](https://github.com/abner-dou))
- Update cloud\_provider\_snapshots.html.markdown [\#605](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/605) ([nhuray](https://github.com/nhuray))
- Fix docs for third party data source. [\#604](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/604) ([jkodroff](https://github.com/jkodroff))
- Fix timeout in acctest [\#602](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/602) ([abner-dou](https://github.com/abner-dou))
- INTMDB-270: fix issue with project resource importer test [\#601](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/601) ([abner-dou](https://github.com/abner-dou))
- Update MDB version info [\#600](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/600) ([themantissa](https://github.com/themantissa))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.2 to 0.38.4 [\#599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/599) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-271: Fixing a bug and improving for custom zone mappings in global cluster config [\#597](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/597) ([coderGo93](https://github.com/coderGo93))
- INTMDB-275: Changed the pointer in some paremeters for custom db role [\#596](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/596) ([coderGo93](https://github.com/coderGo93))
- INTMDB-270: Added  'with\_default\_alerts\_settings' to project resource [\#595](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/595) ([abner-dou](https://github.com/abner-dou))
- Fix backup option from provider\_backup\_enabled to cloud\_backup [\#592](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/592) ([paikend](https://github.com/paikend))
- INTMDB-222: Added Synonyms to Search Index RS and DS [\#591](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/591) ([abner-dou](https://github.com/abner-dou))
- Fix typo: mongodbatlast =\> mongodbatlas [\#589](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/589) ([sonlexqt](https://github.com/sonlexqt))
- Chore\(deps\): Bump github.com/go-test/deep from 1.0.7 to 1.0.8 [\#587](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/587) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-260: Added GCP feature for Private Endpoint [\#586](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/586) ([coderGo93](https://github.com/coderGo93))
- INTMDB-227:Create new Resource and Datasource for Serverless Instance [\#585](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/585) ([abner-dou](https://github.com/abner-dou))
- INTMDB-269: Fix issue with default auto\_scaling\_disk\_gb\_enabled value [\#584](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/584) ([abner-dou](https://github.com/abner-dou))
- fixes failing snapshots because timeout is too short [\#583](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/583) ([ebostijancic](https://github.com/ebostijancic))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.1 to 0.38.2 [\#582](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/582) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-239: Added new resource/datasource and deprecate for cloud backup snapshot and restore job [\#581](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/581) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.12 to 0.38.1 [\#580](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/580) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-216: Added resource and datasource\(s\) of Advanced Cluster [\#570](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/570) ([coderGo93](https://github.com/coderGo93))
- Add Organisation and Project invitations [\#560](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/560) ([beergeek](https://github.com/beergeek))

## [v1.0.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.2) (2021-10-07)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.2-pre.1...v1.0.2)

**Closed issues:**

- gcp cluster doc issue [\#568](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/568)
- mongodbatlas\_auditing documentation mismatch [\#555](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/555)

**Merged pull requests:**

- INTMDB-246: Release v1.0.2 [\#579](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/579) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.11 to 0.37.12 [\#578](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/578) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.0.2-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.2-pre.1) (2021-10-04)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.1...v1.0.2-pre.1)

**Fixed bugs:**

- Error: error getting search index information: json: cannot unmarshal array into Go struct field IndexMapping.mappings.fields of type mongodbatlas.IndexField [\#545](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/545)
- Regression: Data Source for mongodbatlas\_cluster makes terraform hang indefinitely using version 1.0 [\#521](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/521)
- Enabling encryption at rest with any provider for a Cluster will throw error [\#517](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/517)
- realm trigger causes provider to authenticate with atlas even if realm triggers are not in use [\#512](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/512)
- Adding IP to access List failed when lot of entries [\#470](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/470)

**Closed issues:**

- Datalake configuration at creation time [\#561](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/561)
- Nested map variable works if defined in module consumer but not if defined in module itself [\#559](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/559)
- Getting blocked by IP when trying to create a project / cluster [\#557](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/557)

**Merged pull requests:**

- INTMDB-246: pre-release 1.0.2-pre.1 [\#577](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/577) ([abner-dou](https://github.com/abner-dou))
- Roll up for documentation fixes [\#576](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/576) ([themantissa](https://github.com/themantissa))
- INTMDB-259: Fix issue when create a tenant cluster without auto\_scaling\_disk\_gb [\#575](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/575) ([abner-dou](https://github.com/abner-dou))
- INTMDB-203: Fix IOPS restriction on NVME clusters [\#574](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/574) ([abner-dou](https://github.com/abner-dou))
- INTMDB-254: Fix replication\_specs behaviour when update cluster [\#573](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/573) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.10 to 0.37.11 [\#572](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/572) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.7.1 to 2.8.0 [\#571](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/571) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.8 to 0.37.10 [\#569](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/569) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-256: Fixes a bug for updated a role in cloud access authorization [\#567](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/567) ([coderGo93](https://github.com/coderGo93))
- INTMDB-245: Added an example for encryption at rest using azure with a cluster [\#566](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/566) ([coderGo93](https://github.com/coderGo93))
- INTMDB-221: Added projectOwnerID to project resource [\#564](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/564) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.7 to 0.37.8 [\#563](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/563) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-211: Add new advanced shard key options in global cluster resource [\#562](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/562) ([abner-dou](https://github.com/abner-dou))
- INTMDB-252: Added two parameters for cluster advanced configuration [\#558](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/558) ([coderGo93](https://github.com/coderGo93))
- Fix typo in import search index error [\#556](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/556) ([stefanosala](https://github.com/stefanosala))
- INTMDB-230: added  property to maintenance window rs ds [\#552](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/552) ([abner-dou](https://github.com/abner-dou))
- INTMDB-249: Lint and fix linting for examples [\#538](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/538) ([gssbzn](https://github.com/gssbzn))

## [v1.0.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.1) (2021-09-02)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.1-pre.1...v1.0.1)

**Merged pull requests:**

-  tag version 1.0.1 for release [\#553](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/553) ([coderGo93](https://github.com/coderGo93))

## [v1.0.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.1-pre.1) (2021-09-01)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.0...v1.0.1-pre.1)

**Fixed bugs:**

- Cannot define a mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy and enable provider\_backup\_enabled for an existing cluster in the same apply [\#350](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/350)

**Closed issues:**

- Not able to obtain CSRS hostnames from mongodbatlas\_cluster resource [\#543](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/543)
- Cloud Provider Access Setup for Datalake [\#486](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/486)
- Terraform: Unable to fetch connection strings when using 'data' resource for existing cluster [\#422](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/422)
- SDK framework update v2.0.0+ [\#408](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/408)
- Issues with resource and API design for Cloud Provider Snapshot Backup Policy [\#222](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/222)

**Merged pull requests:**

- Chore\(deps\): Bump go.mongodb.org/realm from 0.0.1 to 0.1.0 [\#551](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/551) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.7.0 to 2.7.1 [\#549](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/549) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-251: Update search rs and ds to use go-client v0.12.0 [\#548](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/548) ([abner-dou](https://github.com/abner-dou))
- tag version 1.0.1 for pre release [\#546](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/546) ([coderGo93](https://github.com/coderGo93))
- test: skip instead of fatal for team ids missing [\#544](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/544) ([gssbzn](https://github.com/gssbzn))
- Chore\(deps\): Bump github.com/spf13/cast from 1.3.1 to 1.4.1 [\#542](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/542) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.32.20 to 0.37.7 [\#541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/541) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/mongodb-forks/digest from 1.0.1 to 1.0.3 [\#540](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/540) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump octokit/request-action from 2.0.0 to 2.1.0 [\#539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/539) ([dependabot[bot]](https://github.com/apps/dependabot))
- feat: add dependabot [\#537](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/537) ([gssbzn](https://github.com/gssbzn))
- feat: mcli integration [\#536](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/536) ([gssbzn](https://github.com/gssbzn))
- docs: fix typo cluster.html.markdown [\#535](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/535) ([gssbzn](https://github.com/gssbzn))
- docs: update README [\#534](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/534) ([gssbzn](https://github.com/gssbzn))
- INTMDB-226 - Added forceNew to vpc\_id in network\_peering [\#533](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/533) ([abner-dou](https://github.com/abner-dou))
- Modified workflow to trigger the automated tests [\#532](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/532) ([coderGo93](https://github.com/coderGo93))
- INTMDB-247: Fixes a bug where it's taking 3 minutes to read a cluster [\#530](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/530) ([coderGo93](https://github.com/coderGo93))
- task: check examples are formatted correctly [\#529](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/529) ([gssbzn](https://github.com/gssbzn))
- feat: use golangci lint action [\#528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/528) ([gssbzn](https://github.com/gssbzn))
- task: remove misspell as a dependency [\#527](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/527) ([gssbzn](https://github.com/gssbzn))
- INTMDB-236: Updated the cluster configuration [\#526](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/526) ([coderGo93](https://github.com/coderGo93))
- INTMDB-244: add deprecation notes for cloud backup documentation [\#525](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/525) ([abner-dou](https://github.com/abner-dou))
- INTMDB-237: fix word in private\_endpoint resource documentation [\#523](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/523) ([abner-dou](https://github.com/abner-dou))
- INTMDB-243: Fixes a bug for encryption at rest with new parameters [\#522](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/522) ([coderGo93](https://github.com/coderGo93))
- INTMDB-235: Added example of ldap configuration docs [\#520](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/520) ([coderGo93](https://github.com/coderGo93))
- INTMDB-242: Fixes the bug when if you don set public/private key it would fail for getting realm client [\#519](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/519) ([coderGo93](https://github.com/coderGo93))
- Add stronger warning against attempting a shared tier upgrade [\#516](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/516) ([themantissa](https://github.com/themantissa))
- INTMDB-219: Fixed cluster scaling issue [\#515](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/515) ([abner-dou](https://github.com/abner-dou))
- INTMDB-218: fixes the bug when you try to add more than 100 ip whitelist [\#514](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/514) ([coderGo93](https://github.com/coderGo93))
- INTMDB-225: Fixed network peering resource for Azure [\#513](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/513) ([abner-dou](https://github.com/abner-dou))

## [v1.0.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.0) (2021-08-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.0-pre.1...v1.0.0)

**Merged pull requests:**

-  INTMDB-215: tag version 1.0.0 for release [\#511](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/511) ([coderGo93](https://github.com/coderGo93))
- Update 1.0.0-upgrade-guide.html.markdown [\#510](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/510) ([themantissa](https://github.com/themantissa))
- update the privatelink doc with Azure example [\#509](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/509) ([nikhil-mongo](https://github.com/nikhil-mongo))

## [v1.0.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.0-pre.1) (2021-08-10)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.1...v1.0.0-pre.1)

**Closed issues:**

- Multi cloud not supported? [\#497](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/497)
- Unsupported argument `bi_connector_config`  [\#491](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/491)
- Support for Mongo DB Cluster 4.4 [\#487](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/487)
- Backup policy ID requirement is a catch 22 [\#485](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/485)
- Updating from terraform 0.14.5 tp 0.15.0 or further version \(up to 1.0\) mongodbatlas started to add database\_name="" outsides roles part [\#480](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/480)
- resource mongodbatlas\_auditing audit\_filter param doesn't ignore whitespace changes in json [\#477](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/477)
- documentation: Azure private link example [\#469](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/469)
- Terraform not able to detect all the changes from the mongodb .tf files [\#465](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/465)
- mongodbatlas\_private\_endpoint and mongodbatlas\_private\_endpoint\_interface\_link gets re-created everytime [\#464](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/464)
- mongodbatlas\_database\_user lifecycle ignore\_changes is ignored [\#462](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/462)
- Unable to manage LDAP groups due to forced incorrect auth\_database\_name [\#447](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/447)
- Request for Docs - Migration from Ahkaryle-s provider [\#26](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/26)

**Merged pull requests:**

- INTMDB-215: tag version 1.0.0 for pre-release [\#508](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/508) ([coderGo93](https://github.com/coderGo93))
- INTDMB-223: Updated Cloud Backup to SDK v2 [\#507](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/507) ([abner-dou](https://github.com/abner-dou))
- INTMDB-233: Update linter version [\#506](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/506) ([MihaiBojin](https://github.com/MihaiBojin))
- INTMDB-232: Fix user agent version [\#505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/505) ([MihaiBojin](https://github.com/MihaiBojin))
- update resources documentation to address INTMDB-225 [\#503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/503) ([nikhil-mongo](https://github.com/nikhil-mongo))
- INTMDB-202: Changed to TypeSet for replication specs [\#502](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/502) ([coderGo93](https://github.com/coderGo93))
- INTDMB-223: update search index to sdk v2 [\#501](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/501) ([abner-dou](https://github.com/abner-dou))
- INTMDB-17: fixed import state method in search index resource [\#500](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/500) ([abner-dou](https://github.com/abner-dou))
- Adding autodefer parameter to automatically defer any maintenance [\#499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/499) ([vgarcia-te](https://github.com/vgarcia-te))
- docs: fix typo in GCP network\_peering.network\_name [\#498](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/498) ([daquinoaldo](https://github.com/daquinoaldo))
- INTMDB-180: file env variable spelling error [\#495](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/495) ([abner-dou](https://github.com/abner-dou))
- INTMDB-188: fixed issue related with read non-existing resource [\#494](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/494) ([abner-dou](https://github.com/abner-dou))
- add example for atlas-aws vpc peering [\#493](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/493) ([nikhil-mongo](https://github.com/nikhil-mongo))
- MongoDB Atlas-GCP VPC Peering [\#492](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/492) ([nikhil-mongo](https://github.com/nikhil-mongo))
- MongoDB Atlas - GCP VPC Peering example [\#490](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/490) ([nikhil-mongo](https://github.com/nikhil-mongo))
- INTMDB-183: Migrate to TF SDK 2 [\#489](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/489) ([coderGo93](https://github.com/coderGo93))
- INTMDB-17:  Resource/Data Source Atlas Search  [\#488](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/488) ([abner-dou](https://github.com/abner-dou))
- INTMDB-214: Deprecation of private endpoint [\#484](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/484) ([coderGo93](https://github.com/coderGo93))
- INTMDB-179: added more examples for connection strings [\#483](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/483) ([coderGo93](https://github.com/coderGo93))
- INTMDB-198: Fixes a bug where it appears empty private endpoint in cluster [\#481](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/481) ([coderGo93](https://github.com/coderGo93))
- INTMDB-201: Added to detect changes for name of cluster in update func [\#479](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/479) ([coderGo93](https://github.com/coderGo93))
- INTMDB-28: Added Event Triggers Realm [\#476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/476) ([coderGo93](https://github.com/coderGo93))
- INTMDB-145: Cloud backup schedule  [\#475](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/475) ([coderGo93](https://github.com/coderGo93))
- Starter example improvements and doc update [\#474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/474) ([themantissa](https://github.com/themantissa))
- INTMDB-212: Deprecation of Project IP Whitelist [\#473](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/473) ([coderGo93](https://github.com/coderGo93))
- INTMDB-18-Test for Online Archive and sync attribute for discussion [\#472](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/472) ([leofigy](https://github.com/leofigy))
- INTMDB-128: Modified design when you can get .id from various resources [\#471](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/471) ([coderGo93](https://github.com/coderGo93))
- update README - added plugin dev override [\#468](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/468) ([abner-dou](https://github.com/abner-dou))
- CLOUDP-90710: Expose BASE\_URL so that we can test terraform with a custom server [\#467](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/467) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-200: Fixes a bug about updating a region name with GCP [\#463](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/463) ([coderGo93](https://github.com/coderGo93))
- INTMDB-19: Added resource and datasource\(s\) of data lake [\#414](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/414) ([coderGo93](https://github.com/coderGo93))
- INTMDB-18 : DataSource and Resource support for Online Archive [\#413](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/413) ([leofigy](https://github.com/leofigy))

## [v0.9.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.1) (2021-05-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.1-pre.1...v0.9.1)

**Merged pull requests:**

- chore v0.9.1 changelog update [\#461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/461) ([leofigy](https://github.com/leofigy))
- Missing formatting backtick in the documentation [\#457](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/457) ([lescactus](https://github.com/lescactus))

## [v0.9.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.1-pre.1) (2021-05-14)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.0...v0.9.1-pre.1)

**Fixed bugs:**

- mongodbatlas\_cluster bi\_connector state changes on terraform CLI 0.14.2 even without any bi\_connector configuration - terraform 14  [\#364](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/364)
- Update the CA certificate with a os environment [\#442](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/442) ([pitakill](https://github.com/pitakill))

**Closed issues:**

- New single apply cloud provider access requires encryption\_at\_rest\_provider set in mongodbatlas\_cluster [\#452](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/452)
- Migration to mongodbatlas\_cloud\_provider\_access\_setup / authorization [\#451](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/451)
- GCP can't set region for cluster [\#450](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/450)
- Error verifying GPG signature for provider "mongodbatlas" [\#448](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/448)
- "CLUSTER\_DISK\_IOPS\_INVALID" related error/unexpected update-in-place [\#439](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/439)
- encryption\_at\_rest failing with UNEXPECTED ERROR \(and discussion of Cloud Provider Access possible improvement\) [\#409](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/409)
- Test update - Update test certificate  [\#407](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/407)
- mongodbatlas\_private\_endpoint and mongodbatlas\_private\_endpoint\_interface\_link not working as expected in version 0.7 [\#406](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/406)

**Merged pull requests:**

- INTMDB-207 chore: Doc update for changelog v0.9.1-pre.1 [\#460](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/460) ([leofigy](https://github.com/leofigy))
- INTMDB-208: Updated docs for upgrading private endpoints [\#458](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/458) ([coderGo93](https://github.com/coderGo93))
- INTMDB-205 fixing client update side effects [\#456](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/456) ([leofigy](https://github.com/leofigy))
-  INTMDB-205-client-update bumping the client version up [\#455](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/455) ([leofigy](https://github.com/leofigy))
- Test config update [\#454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/454) ([leofigy](https://github.com/leofigy))
- INTMDB-206 Documentation and example updates [\#453](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/453) ([leofigy](https://github.com/leofigy))
- updated cluster doc and examples  for the new IOPS change [\#446](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/446) ([nikhil-mongo](https://github.com/nikhil-mongo))
- fix page title and sidebar [\#445](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/445) ([themantissa](https://github.com/themantissa))
- chore v0.9.0 changelog [\#444](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/444) ([leofigy](https://github.com/leofigy))

## [v0.9.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.0) (2021-04-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.0-pre.1...v0.9.0)

**Implemented enhancements:**

- Test or TestAccResourceMongoDBAtlasDatabaseUser\_withAWSIAMType with [\#432](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/432) ([pitakill](https://github.com/pitakill))
- INTMDB 186 - Added authorization resource to split the cloud access provider config [\#420](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/420) ([leofigy](https://github.com/leofigy))

**Closed issues:**

- Outdated usage example about "mongodbatlas\_encryption\_at\_rest" [\#424](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/424)

**Merged pull requests:**

- Remove IOPS and adjust parameter description [\#443](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/443) ([themantissa](https://github.com/themantissa))
- remove unnecessary variables and use roles instead of keys for AWS Encryption-AtRest [\#441](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/441) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Update default IOPS and 0.9.0 info guide [\#440](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/440) ([themantissa](https://github.com/themantissa))
- \[Azure VNET Peering\] changed the incorrect parameter used for role assignment and role definition [\#438](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/438) ([nikhil-mongo](https://github.com/nikhil-mongo))
- chore changelog for v0.9.0 prerelease [\#437](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/437) ([leofigy](https://github.com/leofigy))
- Update release.yml [\#436](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/436) ([leofigy](https://github.com/leofigy))
- INTMDB-199: Fixes the error when updating an replication specs after removed one zone [\#434](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/434) ([coderGo93](https://github.com/coderGo93))
- Examples of terratest upgrade [\#431](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/431) ([coderGo93](https://github.com/coderGo93))
- Fix: small doc error [\#428](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/428) ([lescactus](https://github.com/lescactus))
- INTMDB-194: Added func to get db major version for testing [\#427](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/427) ([coderGo93](https://github.com/coderGo93))
- Add  examples creating user with aws\_iam\_type [\#426](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/426) ([alexsergeyev](https://github.com/alexsergeyev))
- INTMDB-155: Fixes a bug related to bi\_connector cluster by deprecating [\#423](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/423) ([coderGo93](https://github.com/coderGo93))
- INTMDB-168: updated docs the format of using dependencies [\#421](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/421) ([coderGo93](https://github.com/coderGo93))
- INTMDB-185: Added parameter regions for GCP network container [\#418](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/418) ([coderGo93](https://github.com/coderGo93))
- TeamsUpdate - fixing small bug, again missing update [\#417](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/417) ([leofigy](https://github.com/leofigy))
- Fixes test about ca certificate x509 [\#416](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/416) ([coderGo93](https://github.com/coderGo93))
- Working example for Atlas-encryptionAtRest-roles with a single tf apply [\#415](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/415) ([zohar-mongo](https://github.com/zohar-mongo))
- INTMDB-181: Detects unnecessary changes changes for azure/gcp encryption at rest [\#412](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/412) ([coderGo93](https://github.com/coderGo93))
- corrected the title by making the variable name plural [\#404](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/404) ([crosbymichael1](https://github.com/crosbymichael1))
- INTMDB-154: Deprecation for provider\_encrypt\_ebs\_volume  [\#403](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/403) ([leofigy](https://github.com/leofigy))
- INTMDB-133: Vendor removal to include terratest samples [\#395](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/395) ([leofigy](https://github.com/leofigy))
- INTMDB-114/115: Added resource, datasource and tests for LDAP configuration and verify [\#379](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/379) ([coderGo93](https://github.com/coderGo93))
- INTMDB-116: Added parameter ldap auth type for resource and datasource\(s\) of database user [\#376](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/376) ([coderGo93](https://github.com/coderGo93))
- INTMDB-16: Added resource and datasource for Custom DNS Configuration for Atlas Clusters on AWS [\#370](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/370) ([coderGo93](https://github.com/coderGo93))
- INTMDB-133: Examples for encryption at rest with roles  [\#369](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/369) ([coderGo93](https://github.com/coderGo93))

## [v0.9.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.0-pre.1) (2021-04-21)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.2...v0.9.0-pre.1)

**Closed issues:**

- TF support for creating api keys [\#433](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/433)
- Apple Silicon \(darwin/arm64\) support [\#430](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/430)
- Tenant Provider Acceptance tests are failing [\#419](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/419)
- 500 \(request "UNEXPECTED\_ERROR"\) [\#411](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/411)
- Error Creating cluster GCP - 500 UNEXPECTED\_ERROR [\#410](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/410)
- r/mongodbatlas\_third\_party\_integration fails on read after create [\#354](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/354)

## [v0.8.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.2) (2021-02-03)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.1...v0.8.2)

**Closed issues:**

- Issues with AWS/Azure Privatelink in v0.8.0 and v0.8.1 [\#401](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/401)
- mongodbatlas-0.4.2: panic: runtime error: invalid memory address or nil pointer dereference [\#399](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/399)
- Seemingly unable to set provider source to mongodb/mongodbatlas for terraform v0.14? [\#396](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/396)
- Missing connection\_strings for Azure Private Link in resource mongodbatlas\_cluster [\#390](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/390)
- Error in Docs [\#387](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/387)

**Merged pull requests:**

- INTMDB-177: chore: release changelog [\#402](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/402) ([leofigy](https://github.com/leofigy))
- INTMDB-174: updated an example for cluster [\#400](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/400) ([coderGo93](https://github.com/coderGo93))
- INTMDB-175: Added azure status and fixes the error about target state [\#398](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/398) ([coderGo93](https://github.com/coderGo93))

## [v0.8.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.1) (2021-01-28)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0...v0.8.1)

**Fixed bugs:**

- Removal of user scopes is not detected by the provider [\#363](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/363)

**Closed issues:**

- Parameter mismatch in mongodbatlas\_privatelink\_endpoint\_service [\#391](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/391)
- Can't add a team to a project [\#389](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/389)
- failed to create policy items while using mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#386](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/386)
- Unable to import resources with "terraform import" - 401 \(request "Unauthorized"\) [\#385](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/385)

**Merged pull requests:**

-  INTMDB-172: chore changelog update for v0.8.1 [\#397](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/397) ([leofigy](https://github.com/leofigy))
- INTMDB-169: delete encoding url path in private endpoint service [\#393](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/393) ([coderGo93](https://github.com/coderGo93))
- INTMDB-158: Added private endpoint in connection strings [\#392](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/392) ([coderGo93](https://github.com/coderGo93))
- INTMDB-163: Wrong order for PrivateLink Endpoint Service and detects unnecessary changes [\#388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/388) ([coderGo93](https://github.com/coderGo93))

## [v0.8.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0) (2021-01-20)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0-pre.2...v0.8.0)

**Closed issues:**

- Managing encryption at rest using iam roles fails [\#382](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/382)
- Unable to use mongodbatlas\_network\_peering data source [\#377](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/377)

**Merged pull requests:**

- INTMDB-153: Create 0.8.0-upgrade-guide.html.markdown [\#384](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/384) ([themantissa](https://github.com/themantissa))

## [v0.8.0-pre.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0-pre.2) (2021-01-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0-pre.1...v0.8.0-pre.2)

**Closed issues:**

- Bad Release Practice [\#381](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/381)

**Merged pull requests:**

- INTMDB-162: Fixes bug about detecting changes and make sensitive values [\#383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/383) ([coderGo93](https://github.com/coderGo93))

## [v0.8.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0-pre.1) (2021-01-15)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.7.0...v0.8.0-pre.1)

**Fixed bugs:**

- Unexpected behaviour from resource `mongodbatlas_teams` when adding username for user not yet part of/Pending to join Organisation [\#329](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/329)
- mongodbatlas\_alert\_configuration - metric\_threshold.threshold is not being passed when value is zero [\#311](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/311)
- After manually deletion of a billing alert, no more plan oder apply will succeed, becuase of an 404 during plan [\#305](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/305)

**Closed issues:**

- Deleting a user from mongo atlas results in a 404 error, not that user getting re-created [\#360](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/360)
- Replace "mongodbatlas\_project\_ip\_whitelist" resource/datasource/docs references with "mongodbatlast\_project\_ip\_accesslist" to reflect API/UI change. [\#358](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/358)
- Quick start for provider is not quick and comes with side effect about `replication_factor` field [\#356](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/356)
- mongodbatlas\_database\_user resource's id attribute does not have the username value [\#348](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/348)
- Mongodbatlas documentation issue with Data Sources [\#347](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/347)
- Please add support for Azure Private Link as a private endpoint [\#346](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/346)
- mongodbatlas\_maintenance\_window fails with BadRequest: Invalid Day of Week [\#289](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/289)

**Merged pull requests:**

- INTMDB-160: Resetting an encryption at rest [\#380](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/380) ([coderGo93](https://github.com/coderGo93))
- INTMDB-149: tag version 0.8.0 for pre-release [\#378](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/378) ([coderGo93](https://github.com/coderGo93))
- Fix typo "requirments" in the PR template [\#375](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/375) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Path escape import id of database user [\#373](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/373) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Update cluster to match Atlas doc [\#372](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/372) ([themantissa](https://github.com/themantissa))
- INTMDB-147: Changed to required in schema of roles for database users [\#371](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/371) ([coderGo93](https://github.com/coderGo93))
- INTMDB-144: Updated for scopes database users [\#368](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/368) ([coderGo93](https://github.com/coderGo93))
- Fix database user resource broken indent [\#367](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/367) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- INTMDB-142: Fixes the bug for alertconfiguration using data dog [\#366](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/366) ([coderGo93](https://github.com/coderGo93))
- INTMDB-133: Updated Encryption At Rest to work with IAM Roles [\#365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/365) ([coderGo93](https://github.com/coderGo93))
- INTMDB-141: Fixing 404 for existing database user [\#362](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/362) ([leofigy](https://github.com/leofigy))
- INTMDB-121: Prevents removing existing users [\#361](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/361) ([leofigy](https://github.com/leofigy))
- update wording [\#359](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/359) ([themantissa](https://github.com/themantissa))
- update the documentation and examples for adding the replication spec… [\#357](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/357) ([nikhil-mongo](https://github.com/nikhil-mongo))
- fix: code cleanup [\#355](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/355) ([AlphaWong](https://github.com/AlphaWong))
- Cloud Access Provider Datasources, Resources, and Documentation \(INTMDB 131\) [\#352](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/352) ([leofigy](https://github.com/leofigy))
- doc fix for db users data source [\#351](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/351) ([themantissa](https://github.com/themantissa))
- AWS and AZURE Private Endpoints \(INTMDB-123 & INTMDB-124\) [\#349](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/349) ([coderGo93](https://github.com/coderGo93))
- Basicexample for starting with Atlas [\#345](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/345) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Fix update function for DB users [\#341](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/341) ([EricZaporzan](https://github.com/EricZaporzan))

## [v0.7.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.7.0) (2020-10-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.5...v0.7.0)

**Fixed bugs:**

- X509 user creation and update throws error  [\#312](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/312)
- inherited\_roles are not correctly removed from custom\_db\_roles [\#280](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/280)
- alertConfigs fix field update in terraform state [\#334](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/334) ([leofigy](https://github.com/leofigy))

**Closed issues:**

- Warning when installing the provider on Terraform 0.13 [\#342](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/342)
- mongodbatals\_network\_container [\#336](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/336)
- Typo in documentation  [\#335](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/335)
- Crash when refreshing TF State for a single cluster [\#330](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/330)
- 500 response on /groups/GROUP-ID/peers [\#320](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/320)
- An invalid enumeration value M5 was specified. [\#318](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/318)
- Container ID on the cluster data source is always empty [\#314](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/314)
- email\_enabled always reported as a change for mongodbatlas\_alert\_configuration [\#306](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/306)

**Merged pull requests:**

- Quick docs for 3rd party [\#344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/344) ([themantissa](https://github.com/themantissa))
- chore: changelog v0.7.0 [\#343](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/343) ([leofigy](https://github.com/leofigy))
- documentation fix \#335 and examples added for the Azure VNET peering and AWS Private Link [\#340](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/340) ([nikhil-mongo](https://github.com/nikhil-mongo))
- mongodbatlas\_alert\_configuration - reset ID if was deleted and it's already in the plan [\#333](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/333) ([leofigy](https://github.com/leofigy))
- New resource and datasource for Project IP Access list [\#332](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/332) ([coderGo93](https://github.com/coderGo93))
- Client upgrade to fix metric threshold value set as 0 [\#331](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/331) ([leofigy](https://github.com/leofigy))
- docs: add mongo SLA link [\#328](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/328) ([gssbzn](https://github.com/gssbzn))
- Example added for database user scope [\#327](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/327) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Add "Sensitive: true" for securing sensitive data in state [\#325](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/325) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Create README and examples directory [\#324](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/324) ([themantissa](https://github.com/themantissa))
- fix: fixes a bug for issue 289 [\#323](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/323) ([coderGo93](https://github.com/coderGo93))
- Third party integrations  [\#321](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/321) ([leofigy](https://github.com/leofigy))
- changed from running on PR to manually trigger acceptance tests [\#319](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/319) ([coderGo93](https://github.com/coderGo93))
- Cluster docs [\#317](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/317) ([nikhil-mongo](https://github.com/nikhil-mongo))
- chore: changelog for v0.6.5 [\#316](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/316) ([leofigy](https://github.com/leofigy))
- Chore: Fix the ProviderVersion in the useragent string [\#309](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/309) ([MihaiBojin](https://github.com/MihaiBojin))

## [v0.6.5](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.5) (2020-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.6...v0.6.5)

## [v0.6.6](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.6) (2020-09-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.4...v0.6.6)

**Fixed bugs:**

- X509 is using the wrong authentication database when updating an existing user  [\#292](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/292)
- `mongodbatlas_cloud_provider_snapshot_backup_policy` `restore_window_days` \(optional value\) is being set even when omitted in resource config [\#290](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/290)

**Closed issues:**

- "mongodbatlas\_alert\_configuration" prints Slack API token in plain text [\#310](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/310)
- Can we create search index using terraform-provider-mongodbatlas? [\#308](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/308)
- Error: rpc error: code = Unavailable desc = transport is closing [\#302](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/302)
- Can't create alerts with using many event\_type  [\#232](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/232)

**Merged pull requests:**

- mongo atlas client update fix \#292 \#312 [\#315](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/315) ([leofigy](https://github.com/leofigy))
- DB user creation error because bad encoding in path [\#313](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/313) ([leofigy](https://github.com/leofigy))
- Database user scopes [\#307](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/307) ([coderGo93](https://github.com/coderGo93))
- Setting deterministic encoding id output, just sorting the keys [\#303](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/303) ([leofigy](https://github.com/leofigy))

## [v0.6.4](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.4) (2020-09-02)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.3...v0.6.4)

**Fixed bugs:**

- Unable to import $external auth users [\#285](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/285)
- custom\_db\_roles cannot be created with only inherited roles [\#279](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/279)
- mongodbatlas\_team data provider team\_id null after successful API query [\#277](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/277)

**Closed issues:**

- There is no parity between the Atlas API documentation and the provider doc in regards to alert event\_type values [\#295](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/295)
- Renaming a custom\_db\_role with attached users is not possible [\#284](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/284)
- changing cluster to \_NVME fails on commented-out IOPS [\#283](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/283)
- Error when assigning a custom db role to a database user.  [\#273](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/273)
- Error when creating `mongodbatlas_project_ip_whitelist` resource [\#266](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/266)
- Can't create a alert for Replication Oplog Window [\#227](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/227)

**Merged pull requests:**

- chore: add Changelog for 0.6.4 [\#301](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/301) ([marinsalinas](https://github.com/marinsalinas))
- fix: added a validation for restore\_window\_days [\#300](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/300) ([PacoDw](https://github.com/PacoDw))
- update auth\_database\_name [\#299](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/299) ([themantissa](https://github.com/themantissa))
- Fix \#227 \#232: Added a new Threshold attribute [\#298](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/298) ([PacoDw](https://github.com/PacoDw))
- Fix \#285: Unable to import $external auth users [\#297](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/297) ([PacoDw](https://github.com/PacoDw))
- Fixes many testacc [\#296](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/296) ([coderGo93](https://github.com/coderGo93))
- Fix \#279 \#280 [\#294](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/294) ([PacoDw](https://github.com/PacoDw))
- GitHub actions tests [\#293](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/293) ([coderGo93](https://github.com/coderGo93))
- Changed the harcoded links from hashicorp repo to mongodb repo [\#288](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/288) ([coderGo93](https://github.com/coderGo93))
- add note about container creation [\#287](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/287) ([themantissa](https://github.com/themantissa))
- Correct cluster labels documentation [\#286](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/286) ([themantissa](https://github.com/themantissa))
- Add templates to repo [\#282](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/282) ([themantissa](https://github.com/themantissa))
- Fix \#277: mongodbatlas\_team data provider team\_id null after successful API query [\#281](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/281) ([PacoDw](https://github.com/PacoDw))
- Create SECURITY.md [\#278](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/278) ([themantissa](https://github.com/themantissa))
- Update README.md [\#276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/276) ([themantissa](https://github.com/themantissa))
- Release configuration [\#275](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/275) ([MihaiBojin](https://github.com/MihaiBojin))
- Improvement for 503 error response while creating a cluster [\#274](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/274) ([coderGo93](https://github.com/coderGo93))
- Cleaned vendored deps [\#272](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/272) ([MihaiBojin](https://github.com/MihaiBojin))
- Replaced the digest auth library with one that supports SHA-256 [\#271](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/271) ([MihaiBojin](https://github.com/MihaiBojin))
- Updated changelog v0.6.3 [\#270](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/270) ([PacoDw](https://github.com/PacoDw))
- fix: fix golangci lint and travis [\#269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/269) ([gssbzn](https://github.com/gssbzn))
- feat: add a unique user agent [\#268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/268) ([gssbzn](https://github.com/gssbzn))
- fix: added validation for autoscaling compute enabled and when true a… [\#267](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/267) ([coderGo93](https://github.com/coderGo93))
- Added a field AwsIAMType for database user [\#264](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/264) ([coderGo93](https://github.com/coderGo93))
- Updated Backup Policy documentation [\#259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/259) ([PacoDw](https://github.com/PacoDw))

## [v0.6.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.3) (2020-07-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.2...v0.6.3)

**Fixed bugs:**

- Can't create a new cluster \(M2/M5\) after 0.6.2 version [\#265](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/265)

**Closed issues:**

- Creating cluster eventually returns 503 [\#256](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/256)

## [v0.6.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.2) (2020-07-16)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.1...v0.6.2)

**Fixed bugs:**

- Adding 16 whitelist entries at the same time causes an error [\#252](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/252)
- Error when create or import cluster - panic: runtime error: invalid memory address or nil pointer dereference [\#243](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/243)
- Cannot re-apply config when M2/M5 `disk_size_gb` is specified incorrectly [\#115](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/115)
- accepter\_region\_name not required for AWS on read/import/update [\#53](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/53)
- fix: resource/project\_ip\_whitelist - modify ip whitelist entry valida… [\#257](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/257) ([marinsalinas](https://github.com/marinsalinas))

**Closed issues:**

- In recommendations, prevent export of keys appearing in OS history [\#261](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/261)

**Merged pull requests:**

- Small change to recommendations [\#263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/263) ([themantissa](https://github.com/themantissa))
- Updated changelog v0.6.2 [\#262](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/262) ([PacoDw](https://github.com/PacoDw))
- Updated go version to v1.14 [\#260](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/260) ([PacoDw](https://github.com/PacoDw))
- Fix auto scaling attributes [\#255](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/255) ([PacoDw](https://github.com/PacoDw))
- add: project\_ip\_whitelist datasource [\#254](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/254) ([gmlp](https://github.com/gmlp))
- imp: team datasource add team name option [\#253](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/253) ([gmlp](https://github.com/gmlp))
- fix: fixes \#115  issue with disk size for shared tiers [\#251](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/251) ([gmlp](https://github.com/gmlp))
- Added golangci configuration and travis fix [\#248](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/248) ([PacoDw](https://github.com/PacoDw))
- Updated the name of module client mongodb atlas [\#244](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/244) ([coderGo93](https://github.com/coderGo93))
- fix: fixes \#53 accepter\_region\_name not required for AWS on read/import/update [\#242](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/242) ([gmlp](https://github.com/gmlp))

## [v0.6.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.1) (2020-06-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.0...v0.6.1)

**Fixed bugs:**

- Error when use provider\_name = TENANT on 0.6.0 mongodbatlas provider version. [\#246](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/246)

**Closed issues:**

- Add MongoDB Collection Data Source [\#250](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/250)

**Merged pull requests:**

- Updated changelog v0.6.1 [\#249](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/249) ([PacoDw](https://github.com/PacoDw))
- Fix \#246: Error when use provider\_name = TENANT on 0.6.0 mongodbatlas provider version [\#247](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/247) ([PacoDw](https://github.com/PacoDw))
- Fix \#243: Error when create or import cluster - panic: runtime error: invalid memory address or nil pointer dereference [\#245](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/245) ([PacoDw](https://github.com/PacoDw))

## [v0.6.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.0) (2020-06-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.5.1...v0.6.0)

**Implemented enhancements:**

- mongodbatlas\_database\_user can not be imported when they contain dashes "-" in the name [\#179](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/179)

**Fixed bugs:**

- Changes to mongodbatlas\_database\_user.role.collection\_name are ignored [\#228](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/228)
- Hour and minute properties don't update when they are zero for mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#211](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/211)
- Issues with advanced\_configuration section on mongodbatlas\_cluster [\#210](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/210)
- Changes are not detected when changing Team's role\_names attribute on mongodbatlas\_project [\#209](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/209)
- terraform plan and apply fails after upgrading this module to 0.5.0 [\#200](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/200)
- Issues upgrading cluster to an AWS NVME tier. [\#132](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/132)

**Closed issues:**

- Updating Snapshot Backup Policy: This resource requires access through a whitelist of ip ranges. [\#235](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/235)
- Cannot import mongodbatlas\_database\_user if username contains a hyphen [\#234](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/234)
- How to create a custom db role using built-in and connection action [\#226](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/226)
- connection\_strings returning empty private values [\#220](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/220)
- Documentation incorrect about accessing connection\_strings from clusters? [\#219](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/219)
- Incorrect description for atlas\_cidr\_block in mongodbatlas\_network\_peering documentation [\#215](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/215)
- RESOURSE or RESOURCE? Spelling change for readme.md [\#185](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/185)
- mongodbatlas\_encryption\_at\_rest key rotation impossible to perform with Azure KeyVault [\#80](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/80)

**Merged pull requests:**

- chore: updated changelog to v0.6.0 [\#241](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/241) ([PacoDw](https://github.com/PacoDw))
- Documentation Improvements and New Guide for 0.6.0 [\#240](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/240) ([themantissa](https://github.com/themantissa))
- fixes \#210: Issues with advanced\_configuration section on mongodbatlas\_cluster [\#238](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/238) ([gmlp](https://github.com/gmlp))
- New parameters about pagination for datasources [\#237](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/237) ([coderGo93](https://github.com/coderGo93))
- fix: fixes \#132 issues upgrading cluster to an AWS NVME tier [\#236](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/236) ([gmlp](https://github.com/gmlp))
- Cluster autoscaling [\#233](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/233) ([coderGo93](https://github.com/coderGo93))
- Fix \#228: Changes to mongodbatlas\_database\_user.role.collection\_name are ignored [\#231](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/231) ([PacoDw](https://github.com/PacoDw))
- fixes \#211: Hour and minute properties don't update when they are zero for mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#230](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/230) ([gmlp](https://github.com/gmlp))
- Terraform sdk [\#229](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/229) ([PacoDw](https://github.com/PacoDw))
- Fix \#209: Changes are not detected when changing Team's role\_names attribute on mongodbatlas\_project [\#225](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/225) ([PacoDw](https://github.com/PacoDw))
- New fields for snapshot restore jobs [\#224](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/224) ([coderGo93](https://github.com/coderGo93))
- Improve connection string doc [\#223](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/223) ([themantissa](https://github.com/themantissa))
- Update network\_peering.html.markdown [\#217](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/217) ([themantissa](https://github.com/themantissa))
- fix: fixed DatabaseUserID to allows names with multiple dashes [\#214](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/214) ([PacoDw](https://github.com/PacoDw))
- Fix \#80 - Update for GCP Encryption at rest [\#212](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/212) ([coderGo93](https://github.com/coderGo93))
- Added field container\_id in resource cluster [\#208](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/208) ([coderGo93](https://github.com/coderGo93))

## [v0.5.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.5.1) (2020-04-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.5.0...v0.5.1)

**Implemented enhancements:**

- Support new private and privateSrv connection strings [\#183](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/183)

**Closed issues:**

- Alert configuration roles array should not be required [\#201](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/201)
- Can't get PrivateLink-aware mongodb+srv address when using privatelink [\#147](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/147)

**Merged pull requests:**

- chore: updated changelog file for v0.5.1 [\#207](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/207) ([PacoDw](https://github.com/PacoDw))
- Fix travis, remove google cookie [\#204](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/204) ([marinsalinas](https://github.com/marinsalinas))
- Fix: improved validation to avoid error 404 [\#203](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/203) ([PacoDw](https://github.com/PacoDw))
- Changed roles to computed [\#202](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/202) ([PacoDw](https://github.com/PacoDw))
- Fixed the documetation menu [\#199](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/199) ([PacoDw](https://github.com/PacoDw))

## [v0.5.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.5.0) (2020-04-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.2...v0.5.0)

**Implemented enhancements:**

- mongodbatlas\_encryption\_at\_rest outputs IAM secrets to stdout  [\#93](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/93)
- Cloud Provider Snapshot Backup Policy [\#180](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/180) ([PacoDw](https://github.com/PacoDw))

**Fixed bugs:**

- TERRAFORM CRASH on importing mongodbatlas\_alert\_configuration [\#171](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/171)

**Closed issues:**

- Problem using Cross Region Replica Set in GCP [\#188](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/188)
- Delete this please. [\#187](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/187)
- Confusing output when modifying a cluster [\#186](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/186)
- Cluster auto-scaling [\#182](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/182)
- Docs with wrong resource type [\#175](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/175)
- On upgrade from 0.4.1 to 0.4.2 start getting errors [\#174](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/174)
- Can't create SYSTEM\_NORMALIZED\_CPU\_IOWAIT alert [\#172](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/172)
- mongodbatlas\_alert\_configuration - not able to specify ROLE for type\_name = "GROUP" [\#153](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/153)

**Merged pull requests:**

- chore: update Changelog file for v0.5.0 version [\#197](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/197) ([marinsalinas](https://github.com/marinsalinas))
- Add CONTRIBUTING file [\#196](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/196) ([themantissa](https://github.com/themantissa))
- Update MongoSDK to v0.2.0 [\#195](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/195) ([marinsalinas](https://github.com/marinsalinas))
- Doc update for private\_ip\_mode [\#194](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/194) ([themantissa](https://github.com/themantissa))
- Peering Container documentation fix [\#193](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/193) ([themantissa](https://github.com/themantissa))
- Update backup documenation [\#191](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/191) ([themantissa](https://github.com/themantissa))
- Fix documentation of roles block role\_name [\#184](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/184) ([fbreckle](https://github.com/fbreckle))
- Connection strings [\#181](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/181) ([coderGo93](https://github.com/coderGo93))
- Typo in `provider_disk_type_name` description [\#178](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/178) ([caitlinelfring](https://github.com/caitlinelfring))
- added roles in schema of notifications for alert configurations [\#177](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/177) ([coderGo93](https://github.com/coderGo93))
- fix-\#175 - missing word in resource name [\#176](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/176) ([themantissa](https://github.com/themantissa))
- Fix \#171: added validation to avoid nil type error [\#173](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/173) ([PacoDw](https://github.com/PacoDw))
- Fix Attributes Reference bullet points [\#168](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/168) ([brunopadz](https://github.com/brunopadz))

## [v0.4.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.2) (2020-03-12)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.1...v0.4.2)

**Fixed bugs:**

- mongodbatlas\_cluster fails to redeploy manually deleted cluster [\#159](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/159)

**Closed issues:**

- mongodbatlas\_alert\_configuration - not able to generate any alerts with event\_type = "OUTSIDE\_METRIC\_THRESHOLD" and matcher.fieldName != "TYPE\_NAME" [\#164](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/164)
- Still cannot create cluster in region ME\_SOUTH\_1 on plugin version 0.4.1 [\#161](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/161)
- mongoatlas\_cluster fails to create  - invalid enumeration value M2 was specified  [\#160](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/160)
- Can't create cluster ME\_SOUTH\_1 region [\#157](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/157)

**Merged pull requests:**

- chore: fix linting issues [\#169](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/169) ([marinsalinas](https://github.com/marinsalinas))
- chore: add changelog file for 0.4.2 version [\#167](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/167) ([marinsalinas](https://github.com/marinsalinas))
- Doc: Fix import for mongodbatlas\_project\_ip\_whitelist [\#166](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/166) ([haidaraM](https://github.com/haidaraM))
- chore: removed wrong validation for matchers.value [\#165](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/165) ([PacoDw](https://github.com/PacoDw))
- feature: add default label to clusters [\#163](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/163) ([marinsalinas](https://github.com/marinsalinas))
- Cleaned Cluster state when it isn't found to allow create it again [\#162](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/162) ([PacoDw](https://github.com/PacoDw))
- cluster: removed array of regions due to they could be changed [\#158](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/158) ([PacoDw](https://github.com/PacoDw))

## [v0.4.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.1) (2020-02-26)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.0...v0.4.1)

**Fixed bugs:**

- Add name argument in mongodbatlas\_project datasource [\#140](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/140)

**Closed issues:**

- Delete timeout for mongodbatlas\_private\_endpoint resource too short [\#151](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/151)
- mongodbatlas\_project name [\#150](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/150)
- mongodbatlas\_custom\_db\_role not waiting for resource creation [\#148](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/148)
- Cannot use mongodbatlas\_maintenance\_window - Error provider does not support [\#145](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/145)
- Error creating users with mongodbatlas\_database\_user \(following documentation examples\) [\#144](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/144)
- Auto Scale Cluster Tier Missing [\#141](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/141)

**Merged pull requests:**

- chore: add changelog file for 0.4.1 version [\#156](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/156) ([marinsalinas](https://github.com/marinsalinas))
- Custom DB Roles: added refresh function to allow to create/remove multiple custom roles at the same time [\#155](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/155) ([PacoDw](https://github.com/PacoDw))
- chore: increase timeout when delete in private\_endpoint resource [\#154](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/154) ([marinsalinas](https://github.com/marinsalinas))
- add upgrade guide [\#149](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/149) ([themantissa](https://github.com/themantissa))
- Correct `mongodbatlas_teams` resource name in docs [\#143](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/143) ([mattt416](https://github.com/mattt416))
- Project data source [\#142](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/142) ([PacoDw](https://github.com/PacoDw))

## [v0.4.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.0) (2020-02-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.3.1...v0.4.0)

**Implemented enhancements:**

- expose 'paused' as an argument for mongodbatlas\_cluster [\#105](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/105)
- Add pitEnabled feature of mongodbatlas\_cluster resource [\#104](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/104)
- Support for AWS security groups in mongodbatlas\_project\_ip\_whitelist [\#67](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/67)

**Fixed bugs:**

- Cannot update GCP network peer [\#86](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/86)
- Cluster fails to build on 0.3.1 when mongo\_db\_major\_version is not specified [\#81](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/81)
- Crash \(panic, interface conversion error\) when creating mongodbatlas\_encryption\_at\_rest in Azure [\#74](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/74)
- Creating M2 cluster without specifying disk\_size\_gb results in 400 Bad Request [\#72](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/72)

**Closed issues:**

- add mongodbatlas\_project datasource [\#137](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/137)
- Inconsistent documentation between GitHub repo and Terraform site [\#136](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/136)
- Cloud provider snapshot management [\#124](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/124)
- Add support in cluster-tier autoscaling [\#123](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/123)
- Continuous Backup is not supported for \(new\) AWS clusters [\#121](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/121)
- If don't specify IOPS when creating M10 or M20 cluster a 0 value is passed in causing failure [\#120](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/120)
- When is mongodbatlas\_project\_ip\_whitelist security group feature going to be released? [\#114](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/114)
- Error creating MongoDB Cluster: unexpected EOF [\#110](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/110)
- Issue with import  mongodbatlas\_cloud\_provider\_snapshot\_restore\_job [\#109](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/109)
- mongodbatlas\_network\_container Already exists [\#88](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/88)
- mongodbatlas\_network\_container doesn't form a valid json request [\#83](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/83)
- mongodbatlas\_network\_containers datasource doesn't work with Azure [\#71](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/71)
- missing schema for provider "mongodbatlas" resource type mongodbatlas\_ip\_whitelist [\#70](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/70)
- Whitelisted Project IPs when manually deleted causes failure at next plan/apply [\#68](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/68)
- Modifying project ip whitelist destroy and re-create all resources [\#51](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/51)

**Merged pull requests:**

- Changelog for v0.4.0 [\#138](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/138) ([marinsalinas](https://github.com/marinsalinas))
- Readme: Updated env variables [\#135](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/135) ([PacoDw](https://github.com/PacoDw))
- Database Users: updated Read Function to avoid plugin error when it upgrades [\#133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/133) ([PacoDw](https://github.com/PacoDw))
- Fix snapshot import with hyphened cluster\_name [\#131](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/131) ([marinsalinas](https://github.com/marinsalinas))
- Spelling and grammer [\#130](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/130) ([CMaylone](https://github.com/CMaylone))
- chore: added database\_name as deprecated attribute [\#129](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/129) ([PacoDw](https://github.com/PacoDw))
- Encryption At Rest: fixed issues and added an enhancement [\#128](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/128) ([PacoDw](https://github.com/PacoDw))
- Add PIT enabled argumento to Cluster Resource and Data Source [\#126](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/126) ([marinsalinas](https://github.com/marinsalinas))
- X509 Authentication Database User [\#125](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/125) ([PacoDw](https://github.com/PacoDw))
- Database users: added x509\_type attribute [\#122](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/122) ([PacoDw](https://github.com/PacoDw))
- Shared tier doc edits [\#119](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/119) ([themantissa](https://github.com/themantissa))
- Private endpoints [\#118](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/118) ([PacoDw](https://github.com/PacoDw))
- Update cluster doc [\#117](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/117) ([themantissa](https://github.com/themantissa))
- Update backup, add links [\#116](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/116) ([themantissa](https://github.com/themantissa))
- Projects: adding teams attribute [\#113](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/113) ([PacoDw](https://github.com/PacoDw))
- Update cluster.html.markdown [\#112](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/112) ([themantissa](https://github.com/themantissa))
- Fix DiskSizeGB missing [\#111](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/111) ([marinsalinas](https://github.com/marinsalinas))
- Terraform resource for MongoDB Custom Roles [\#108](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/108) ([PacoDw](https://github.com/PacoDw))
- Fix peering resource [\#107](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/107) ([PacoDw](https://github.com/PacoDw))
- Fix \#68: Added the ability to re-create the whitelist entry when it's remove manually [\#106](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/106) ([PacoDw](https://github.com/PacoDw))
- Updating `git clone` command to reference current repository [\#103](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/103) ([macintacos](https://github.com/macintacos))
- Cluster label and plugin attribute [\#102](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/102) ([PacoDw](https://github.com/PacoDw))
- Added functions to handle labels attribute in some resources [\#101](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/101) ([PacoDw](https://github.com/PacoDw))
- Added labels attr for Database User resource [\#100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/100) ([PacoDw](https://github.com/PacoDw))
- Alert configuration resource and data source [\#99](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/99) ([PacoDw](https://github.com/PacoDw))
- Update database\_user.html.markdown [\#98](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/98) ([themantissa](https://github.com/themantissa))
- update containers and ip whitelist doc [\#96](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/96) ([themantissa](https://github.com/themantissa))
- Add provider\_name to containers data source [\#95](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/95) ([marinsalinas](https://github.com/marinsalinas))
- Whitelist [\#94](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/94) ([PacoDw](https://github.com/PacoDw))
- Network Peering RS: remove provider\_name=AWS as default, use Required=true instead i… [\#92](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/92) ([marinsalinas](https://github.com/marinsalinas))
- Update project.html.markdown [\#91](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/91) ([themantissa](https://github.com/themantissa))
- Feat: Global Cluster Configuration Resource and Data Source. [\#90](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/90) ([marinsalinas](https://github.com/marinsalinas))
- fix: validate if mongo\_db\_major\_version is set in cluster resource [\#85](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/85) ([marinsalinas](https://github.com/marinsalinas))
- Auditing Resource and Data Source [\#82](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/82) ([PacoDw](https://github.com/PacoDw))
- Feat: Team Resource and Data Source [\#79](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/79) ([marinsalinas](https://github.com/marinsalinas))
- Maintenance window ds [\#78](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/78) ([PacoDw](https://github.com/PacoDw))
- Added default Disk Size when it doesn't set up on cluster resource [\#77](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/77) ([PacoDw](https://github.com/PacoDw))
- Maintenance window rs [\#76](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/76) ([PacoDw](https://github.com/PacoDw))
- website: collapse data sources sidebar by default [\#75](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/75) ([marinsalinas](https://github.com/marinsalinas))
- Improvements to Peering Resources [\#73](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/73) ([themantissa](https://github.com/themantissa))
- Remove dupe argument in docs [\#69](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/69) ([ktmorgan](https://github.com/ktmorgan))
- Clarify Azure Option in Doc [\#66](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/66) ([themantissa](https://github.com/themantissa))

## [v0.3.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.3.1) (2019-11-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.3.0...v0.3.1)

**Fixed bugs:**

- Confirmation on timelimit for a terraform apply [\#57](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/57)

**Closed issues:**

- Not able to create M0 clusters [\#64](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/64)
- No way to modify advanced configuration options for a cluster [\#61](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/61)
- mongodbatlas\_network\_peering outputting invalid json [\#59](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/59)
- Syntax are not mandatory and creates confusion [\#58](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/58)
- data source mongodbatlas\_network\_peering retrieves the same for id and connection\_id [\#56](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/56)
- Add resource for maintenance window [\#55](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/55)
- Error encryption\_at\_rest  rpc unavailable desc [\#54](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/54)
- specify oplog size? [\#52](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/52)
- Add resource for custom database roles [\#50](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/50)
- An invalid enumeration value US\_EAST\_1 was specified. [\#49](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/49)
- Version 0.3.0 [\#47](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/47)
- plugin.terraform-provider-mongodbatlas\_v0.2.0\_x4: panic: runtime error: index out of range [\#36](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/36)

**Merged pull requests:**

- chore: add changelog file for 0.3.1 version [\#65](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/65) ([marinsalinas](https://github.com/marinsalinas))
- Added format function to handle the mongo\_db\_major\_version attr [\#63](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/63) ([PacoDw](https://github.com/PacoDw))
- Added cast func to avoid panic by nil value [\#62](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/62) ([PacoDw](https://github.com/PacoDw))
- Cluster advanced configuration Options [\#60](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/60) ([PacoDw](https://github.com/PacoDw))

## [v0.3.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.3.0) (2019-10-08)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.2.0...v0.3.0)

**Closed issues:**

- Upgrade from M2 to M10 fails [\#42](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/42)
- GCP Peering endless terraform apply [\#41](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/41)
- AWS clusters default provider\_encrypt\_ebs\_volume to false [\#40](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/40)
- mongodbatlas\_network\_peering Internal Servier Error [\#35](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/35)
- Problem encryption\_at\_rest [\#33](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/33)
- Problem destroying network peering container [\#30](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/30)
- Bug VPC Peering between GCP and Atlas [\#29](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/29)

**Merged pull requests:**

- chore: add changelog file for 0.3.0 version [\#48](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/48) ([marinsalinas](https://github.com/marinsalinas))
- Clarify Doc Examples and Text [\#46](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/46) ([themantissa](https://github.com/themantissa))
- fix-\#40: added true value by defualt on provider\_encrypt\_ebs\_volume attr [\#45](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/45) ([PacoDw](https://github.com/PacoDw))
- make provider\_name forced new to avoid patch problems [\#44](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/44) ([marinsalinas](https://github.com/marinsalinas))
- Network peering [\#43](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/43) ([PacoDw](https://github.com/PacoDw))
- Update readme with more info [\#39](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/39) ([themantissa](https://github.com/themantissa))
- Fix: Network container [\#38](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/38) ([PacoDw](https://github.com/PacoDw))
- Doc updates [\#37](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/37) ([themantissa](https://github.com/themantissa))

## [v0.2.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.2.0) (2019-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.1.1...v0.2.0)

**Closed issues:**

- Unable to create project with peering only connections [\#24](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/24)
- importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#21](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/21)
- Support the vscode terraform extension [\#19](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/19)
- Bug: VPC Peering Atlas-GCP [\#17](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/17)
- PATCH network peering failed with no peer found [\#14](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/14)

**Merged pull requests:**

- chore: add changelog for new release [\#34](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/34) ([marinsalinas](https://github.com/marinsalinas))
- Add Private IP Mode Resource. [\#32](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/32) ([marinsalinas](https://github.com/marinsalinas))
- Moved provider\_name values to the correct section [\#31](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/31) ([kgriffiths](https://github.com/kgriffiths))
- website: add links to Atlas Region name reference. [\#28](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/28) ([themantissa](https://github.com/themantissa))
- Encryption at rest fix [\#27](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/27) ([marinsalinas](https://github.com/marinsalinas))
- website: make resources side nav expanded as default [\#25](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/25) ([marinsalinas](https://github.com/marinsalinas))
- fix: importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#23](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/23) ([PacoDw](https://github.com/PacoDw))
- Fix \#14: PATCH network peering failed with no peer found [\#22](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/22) ([PacoDw](https://github.com/PacoDw))
- fix: change the test configuration for AWS and GCP [\#20](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/20) ([PacoDw](https://github.com/PacoDw))

## [v0.1.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.1.1) (2019-09-05)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.1.0...v0.1.1)

**Fixed bugs:**

- panic: runtime error: index out of range [\#1](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1)

**Closed issues:**

- GCP peering problem [\#16](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/16)
- Cluster creation with Azure provider failed [\#15](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/15)
- Error creating MongoDB Cluster [\#9](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/9)
- Failed to create Atlas network peering container [\#7](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/7)
- Bug: Invalid attribute diskIOPS specified. [\#2](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2)

**Merged pull requests:**

- chore: update changelog [\#18](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/18) ([marinsalinas](https://github.com/marinsalinas))
- website: fix typo [\#13](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/13) ([heimweh](https://github.com/heimweh))
- fix: add the correct func to check the env variables on peering datasources [\#12](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/12) ([PacoDw](https://github.com/PacoDw))
- Fix diskIOPS attribute for GCP and Azure [\#11](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/11) ([PacoDw](https://github.com/PacoDw))
- website: fix data sources sidebar always collapsed [\#10](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/10) ([marinsalinas](https://github.com/marinsalinas))
- mongodbatlas\_network\_\(peering and container\): add more testing case [\#8](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/8) ([PacoDw](https://github.com/PacoDw))
- website: fix typo in MongoDB Atlas Services [\#5](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/5) ([marinsalinas](https://github.com/marinsalinas))
- Ip whitelist entries: removing all entries whitelist by terraform user [\#4](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/4) ([PacoDw](https://github.com/PacoDw))
- Refactored import function to get all ip\_addresses and cird\_blocks entries [\#3](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3) ([PacoDw](https://github.com/PacoDw))

## [v0.1.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.1.0) (2019-08-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/3e1c5c44b56aee2f153ec618c804dd170bbefbd4...v0.1.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
