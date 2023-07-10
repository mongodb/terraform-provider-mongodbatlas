package provider

// func TestAccProjectRSProject_CreateWithProjectOwner(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_CreateWithFalseDefaultSettings(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultSettings(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_CreateWithFalseDefaultAdvSettings(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_withUpdatedRole(t *testing.T) {
// 	var (
// 		resourceName    = "mongodbatlas_project.test"
// 		projectName     = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		roleName        = "GROUP_DATA_ACCESS_ADMIN"
// 		roleNameUpdated = "GROUP_READ_ONLY"
// 		clusterCount    = "0"
// 		teamsIds        = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleName),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
// 				),
// 			},
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleNameUpdated),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_importBasic(t *testing.T) {
// 	var (
// 		projectName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
// 		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		resourceName = "mongodbatlas_project.test"
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasic(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
// 					[]*matlas.ProjectTeam{},
// 					[]*apiKey{},
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportStateIdFunc:       testAccCheckMongoDBAtlasProjectImportStateIDFunc(resourceName),
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"with_default_alerts_settings"},
// 			},
// 		},
// 	})
// }
