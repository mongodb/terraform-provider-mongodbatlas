# Changelog

## [v0.8.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.1)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0...v0.8.1)

**Fixed bugs:**

- Removal of user scopes is not detected by the provider [\#363](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/363)

**Closed issues:**

- Parameter mismatch in mongodbatlas\_privatelink\_endpoint\_service [\#391](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/391)
- Can't add a team to a project [\#389](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/389)
- failed to create policy items while using mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#386](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/386)
- Unable to import resources with "terraform import" - 401 \(request "Unauthorized"\) [\#385](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/385)

**Merged pull requests:**

- INTMDB-158: Added private endpoint in connection strings [\#392](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/392) ([coderGo93](https://github.com/coderGo93))
- INTMDB-163: Wrong order for PrivateLink Endpoint Service and detects unnecessary changes [\#388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/388) ([coderGo93](https://github.com/coderGo93))
- INTMDB-169: delete encoding url path in private endpoint service [\#393](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/393) ([coderGo93](https://github.com/coderGo93))

## [v0.8.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.7.0...v0.8.0)

**Fixed bugs:**

- Unexpected behaviour from resource `mongodbatlas\_teams` when adding username for user not yet part of/Pending to join Organisation [\#329](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/329)
- mongodbatlas\_alert\_configuration - metric\_threshold.threshold is not being passed when value is zero [\#311](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/311)
- After manually deletion of a billing alert, no more plan oder apply will succeed, becuase of an 404 during plan [\#305](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/305)

**Closed issues:**

- Deleting a user from mongo atlas results in a 404 error, not that user getting re-created [\#360](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/360)
- Replace "mongodbatlas\_project\_ip\_whitelist" resource/datasource/docs references with "mongodbatlast\_project\_ip\_accesslist" to reflect API/UI change. [\#358](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/358)
- Quick start for provider is not quick and comes with side effect about `replication\_factor` field [\#356](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/356)
- mongodbatlas\_database\_user resource's id attribute does not have the username value [\#348](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/348)
- Mongodbatlas documentation issue with Data Sources [\#347](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/347)
- Please add support for Azure Private Link as a private endpoint [\#346](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/346)
- mongodbatlas\_maintenance\_window fails with BadRequest: Invalid Day of Week [\#289](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/289)

**Merged pull requests:**

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
- chore: changelog v0.7.0 [\#343](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/343) ([leofigy](https://github.com/leofigy))
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
- `mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy` `restore\_window\_days` \(optional value\) is being set even when omitted in resource config [\#290](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/290)

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
- Error when creating `mongodbatlas\_project\_ip\_whitelist` resource [\#266](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/266)
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
- Added a field AwsIAMType for database user [\#264](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/264) ([coderGo93](https://github.com/coderGo93))
- Updated Backup Policy documentation [\#259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/259) ([PacoDw](https://github.com/PacoDw))

## [v0.6.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.3) (2020-07-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.2...v0.6.3)

**Fixed bugs:**

- Can't create a new cluster \(M2/M5\) after 0.6.2 version [\#265](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/265)

**Closed issues:**

- Creating cluster eventually returns 503 [\#256](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/256)

**Merged pull requests:**

- Updated changelog v0.6.3 [\#270](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/270) ([PacoDw](https://github.com/PacoDw))
- fix: fix golangci lint and travis [\#269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/269) ([gssbzn](https://github.com/gssbzn))
- feat: add a unique user agent [\#268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/268) ([gssbzn](https://github.com/gssbzn))
- fix: added validation for autoscaling compute enabled and when true a… [\#267](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/267) ([coderGo93](https://github.com/coderGo93))

## [v0.6.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.2) (2020-07-16)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.1...v0.6.2)

**Fixed bugs:**

- Adding 16 whitelist entries at the same time causes an error [\#252](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/252)
- Error when create or import cluster - panic: runtime error: invalid memory address or nil pointer dereference [\#243](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/243)
- Cannot re-apply config when M2/M5 `disk\_size\_gb` is specified incorrectly [\#115](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/115)
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

- fix: fixes \#53 accepter\_region\_name not required for AWS on read/import/update [\#242](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/242) ([gmlp](https://github.com/gmlp))
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
- Typo in `provider\_disk\_type\_name` description [\#178](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/178) ([caitlinelfring](https://github.com/caitlinelfring))
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
- Correct `mongodbatlas\_teams` resource name in docs [\#143](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/143) ([mattt416](https://github.com/mattt416))
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
- Feat: Global Cluster Configuration Resource and Data Source. [\#90](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/90) ([marinsalinas](https://github.com/marinsalinas))
- Auditing Resource and Data Source [\#82](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/82) ([PacoDw](https://github.com/PacoDw))
- Feat: Team Resource and Data Source [\#79](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/79) ([marinsalinas](https://github.com/marinsalinas))
- Maintenance window ds [\#78](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/78) ([PacoDw](https://github.com/PacoDw))
- Added default Disk Size when it doesn't set up on cluster resource [\#77](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/77) ([PacoDw](https://github.com/PacoDw))
- Maintenance window rs [\#76](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/76) ([PacoDw](https://github.com/PacoDw))

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
