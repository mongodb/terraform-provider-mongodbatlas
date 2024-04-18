package pushbasedlogexport_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigPushBasedLogExport_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0")

	var (
		orgID                = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName          = acc.RandomProjectName()
		s3BucketNamePrefix   = fmt.Sprintf("tf-%s", acc.RandomName())
		s3BucketName1        = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2        = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		config               = configBasic(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, nonEmptyPrefixPath, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckAwsEnvBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(commonChecks(s3BucketName1, nonEmptyPrefixPath)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigPushBasedLogExport_noPrefixPath(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0")

	var (
		orgID                = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName          = acc.RandomProjectName()
		s3BucketNamePrefix   = fmt.Sprintf("tf-%s", acc.RandomName())
		s3BucketName1        = fmt.Sprintf("%s-1", s3BucketNamePrefix)
		s3BucketName2        = fmt.Sprintf("%s-2", s3BucketNamePrefix)
		s3BucketPolicyName   = fmt.Sprintf("%s-s3-policy", s3BucketNamePrefix)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		config               = configBasic(projectName, orgID, s3BucketName1, s3BucketName2, s3BucketPolicyName, awsIAMRoleName, awsIAMRolePolicyName, defaultPrefixPath, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckAwsEnvBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(commonChecks(s3BucketName1, nonEmptyPrefixPath)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
