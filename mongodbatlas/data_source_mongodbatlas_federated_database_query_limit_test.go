package mongodbatlas

// func TestAccDataSourceFederatedDatabaseQueryLimit_basic(t *testing.T) {
// 	SkipTestExtCred(t)
// 	var (
// 		queryLimit   matlas.DataFederationQueryLimit
// 		resourceName = "mongodbatlas_federated_database_query_limit.test"
// 		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		// projectName  = acctest.RandomWithPrefix("test-acc-project")
// 		// clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
// 		// hostname     = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
// 		// username     = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
// 		// password     = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
// 		// port         = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
// 	)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t); testCheckLDAP(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasLDAPConfigurationDestroy, // TODO: implement in resource_test
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasFederatedDatabaseQueryLimitConfig(projectName, orgID, clusterName, hostname, username, password, cast.ToInt(port)),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasLDAPVerifyExists(resourceName, &ldapVerify),

// 					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
// 					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
// 					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
// 					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
// 					resource.TestCheckResourceAttrSet(resourceName, "port"),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccMongoDBAtlasFederatedDatabaseQueryLimitConfig(projectName, orgID, clusterName1, clusterName2, hostname, username, password string, port int) string {
// 	return fmt.Sprintf(`
// 		resource "mongodbatlas_project" "test" {
// 			name   = "%[1]s"
// 			org_id = "%[2]s"
// 		}

// 		resource "mongodbatlas_cluster" "test" {
// 			project_id   = mongodbatlas_project.test.id
// 			name         = "%[3]s"

// 			// Provider Settings "block"
// 			provider_name               = "AWS"
// 			provider_region_name        = "US_EAST_2"
// 			provider_instance_size_name = "M10"
// 			provider_backup_enabled     = true //enable cloud provider snapshots
// 		}

// 		resource "mongodbatlas_cluster" "test" {
// 			project_id   = mongodbatlas_project.test.id
// 			name         = "%[4]s"

// 			// Provider Settings "block"
// 			provider_name               = "AWS"
// 			provider_region_name        = "US_EAST_2"
// 			provider_instance_size_name = "M10"
// 			provider_backup_enabled     = true //enable cloud provider snapshots
// 		}

// 		data "mongodbatlas_ldap_verify" "test" {
// 			project_id = mongodbatlas_ldap_verify.test.project_id
// 			request_id = mongodbatlas_ldap_verify.test.request_id
// 		}
// `, projectName, orgID, clusterName, hostname, username, password, port)
// }
