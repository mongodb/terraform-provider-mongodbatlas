package accesslistapikey_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectRSAccesslistAPIKey_SettingIPAddress(t *testing.T) {
	var (
		resourceName     = "mongodbatlas_access_list_api_key.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		ipAddress        = acc.RandomIP(179, 154, 226)
		description      = acc.RandomName()
		updatedIPAddress = acc.RandomIP(179, 154, 228)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithIPAddress(orgID, description, ipAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
				),
			},
			{
				Config: configWithIPAddress(orgID, description, updatedIPAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectRSAccessListAPIKey_SettingCIDRBlock(t *testing.T) {
	var (
		resourceName     = "mongodbatlas_access_list_api_key.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		cidrBlock        = acc.RandomIP(179, 154, 226) + "/32"
		description      = acc.RandomName()
		updatedCIDRBlock = acc.RandomIP(179, 154, 228) + "/32"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithCIDRBlock(orgID, description, cidrBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				Config: configWithCIDRBlock(orgID, description, updatedCIDRBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
				),
			},
		},
	})
}

func TestAccProjectRSAccessListAPIKey_SettingCIDRBlock_WideCIDR(t *testing.T) {
	var (
		resourceName     = "mongodbatlas_access_list_api_key.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		cidrBlock        = "100.10.0.0/16"
		description      = acc.RandomName()
		updatedCIDRBlock = "172.10.0.0/16"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithCIDRBlock(orgID, description, cidrBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				Config: configWithCIDRBlock(orgID, description, updatedCIDRBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
				),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.GetOrgAccessEntry(context.Background(), ids["org_id"], ids["entry"], ids["api_key_id"]).Execute()

		if err != nil {
			return fmt.Errorf("access list API Key (%s) does not exist", ids["api_key_id"])
		}

		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_access_list_api_key" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ProgrammaticAPIKeysApi.GetOrgAccessEntry(context.Background(), ids["org_id"], ids["entry"], ids["api_key_id"]).Execute()
		if err == nil {
			return fmt.Errorf("access list API Key (%s) still exists", ids["api_key_id"])
		}
	}

	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["org_id"], ids["api_key_id"], ids["entry"]), nil
	}
}

func configWithIPAddress(orgID, description, ipAddress string) string {
	return fmt.Sprintf(`

	   resource "mongodbatlas_api_key" "test" {
		  org_id = %[1]q
		  description = %[2]q
		  role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	    }

		resource "mongodbatlas_access_list_api_key" "test" {
			org_id = %[1]q
			api_key_id = mongodbatlas_api_key.test.api_key_id
			ip_address = %[3]q
		}
	`, orgID, description, ipAddress)
}

func configWithCIDRBlock(orgID, description, cidrBlock string) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_api_key" "test" {
		org_id = %[1]q
		description = %[2]q
		role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	  }

		resource "mongodbatlas_access_list_api_key" "test" {
		  org_id = %[1]q
		  api_key_id = mongodbatlas_api_key.test.api_key_id
		  cidr_block = %[3]q
		}
	`, orgID, description, cidrBlock)
}
