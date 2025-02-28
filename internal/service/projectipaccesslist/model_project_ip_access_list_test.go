package projectipaccesslist_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaccesslist"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

var (
	AWSSecurityGroup = "AWSSecurityGroup"
	CIDRBlock        = "CIDRBlock"
	IPAddress        = "IPAddress"
	Comment          = "Comment"
	GroupID          = "GroupID"
	IPAddressID      = conversion.EncodeStateID(map[string]string{
		"entry":      IPAddress,
		"project_id": GroupID,
	})
	CidrBlockID = conversion.EncodeStateID(map[string]string{
		"entry":      CIDRBlock,
		"project_id": GroupID,
	})
	AwsSecurityGroupID = conversion.EncodeStateID(map[string]string{
		"entry":      AWSSecurityGroup,
		"project_id": GroupID,
	})
)

func TestNewMongoDBProjectIPAccessList(t *testing.T) {
	testCases := []struct {
		tfModel        *projectipaccesslist.TfProjectIPAccessListModel
		expectedResult *[]admin.NetworkPermissionEntry
		name           string
	}{
		{
			name: "NewMongoDBProjectIPAccessList",
			tfModel: &projectipaccesslist.TfProjectIPAccessListModel{
				AWSSecurityGroup: types.StringValue(AWSSecurityGroup),
				CIDRBlock:        types.StringValue(CIDRBlock),
				IPAddress:        types.StringValue(IPAddress),
				Comment:          types.StringValue(Comment),
			},
			expectedResult: &[]admin.NetworkPermissionEntry{
				{
					AwsSecurityGroup: &AWSSecurityGroup,
					CidrBlock:        &CIDRBlock,
					IpAddress:        &IPAddress,
					Comment:          &Comment,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := projectipaccesslist.NewMongoDBProjectIPAccessList(tc.tfModel)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfProjectIPAccessListModel(t *testing.T) {
	testCases := []struct {
		tfModel        *projectipaccesslist.TfProjectIPAccessListModel
		sdkModel       *admin.NetworkPermissionEntry
		expectedResult *projectipaccesslist.TfProjectIPAccessListModel
		name           string
	}{
		{
			name: "NewTfProjectIPAccessListModel with IpAddress",
			tfModel: &projectipaccesslist.TfProjectIPAccessListModel{
				Timeouts: timeouts.Value{},
			},
			sdkModel: &admin.NetworkPermissionEntry{
				IpAddress: &IPAddress,
				Comment:   &Comment,
				GroupId:   &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListModel{
				ID:               types.StringValue(IPAddressID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(""),
				CIDRBlock:        types.StringValue(""),
				IPAddress:        types.StringValue(IPAddress),
				Comment:          types.StringValue(Comment),
				Timeouts:         timeouts.Value{},
			},
		},
		{
			name: "NewTfProjectIPAccessListModel with CidrBlock",
			tfModel: &projectipaccesslist.TfProjectIPAccessListModel{
				Timeouts: timeouts.Value{},
			},
			sdkModel: &admin.NetworkPermissionEntry{
				CidrBlock: &CIDRBlock,
				IpAddress: &IPAddress,
				Comment:   &Comment,
				GroupId:   &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListModel{
				ID:               types.StringValue(CidrBlockID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(""),
				CIDRBlock:        types.StringValue(CIDRBlock),
				IPAddress:        types.StringValue(IPAddress),
				Comment:          types.StringValue(Comment),
				Timeouts:         timeouts.Value{},
			},
		},
		{
			name: "NewTfProjectIPAccessListModel with AwsSecurityGroup",
			tfModel: &projectipaccesslist.TfProjectIPAccessListModel{
				Timeouts: timeouts.Value{},
			},
			sdkModel: &admin.NetworkPermissionEntry{
				AwsSecurityGroup: &AWSSecurityGroup,
				IpAddress:        &IPAddress,
				Comment:          &Comment,
				GroupId:          &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListModel{
				ID:               types.StringValue(AwsSecurityGroupID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(AWSSecurityGroup),
				CIDRBlock:        types.StringValue(""),
				IPAddress:        types.StringValue(IPAddress),
				Comment:          types.StringValue(Comment),
				Timeouts:         timeouts.Value{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := projectipaccesslist.NewTfProjectIPAccessListModel(tc.tfModel, tc.sdkModel)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTfProjectIPAccessListDSModel(t *testing.T) {
	testCases := []struct {
		sdkModel       *admin.NetworkPermissionEntry
		expectedResult *projectipaccesslist.TfProjectIPAccessListDSModel
		name           string
	}{
		{
			name: "NewTfProjectIPAccessListDSModel with IpAddress",
			sdkModel: &admin.NetworkPermissionEntry{
				AwsSecurityGroup: &AWSSecurityGroup,
				CidrBlock:        &CIDRBlock,
				IpAddress:        &IPAddress,
				Comment:          &Comment,
				GroupId:          &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListDSModel{
				ID:               types.StringValue(IPAddressID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(AWSSecurityGroup),
				CIDRBlock:        types.StringValue(CIDRBlock),
				IPAddress:        types.StringValue(IPAddress),
				Comment:          types.StringValue(Comment),
			},
		},
		{
			name: "NewTfProjectIPAccessListDSModel with CidrBlock",
			sdkModel: &admin.NetworkPermissionEntry{
				CidrBlock: &CIDRBlock,
				Comment:   &Comment,
				GroupId:   &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListDSModel{
				ID:               types.StringValue(CidrBlockID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(""),
				CIDRBlock:        types.StringValue(CIDRBlock),
				IPAddress:        types.StringValue(""),
				Comment:          types.StringValue(Comment),
			},
		},
		{
			name: "NewTfProjectIPAccessListDSModel with AwsSecurityGroup",
			sdkModel: &admin.NetworkPermissionEntry{
				AwsSecurityGroup: &AWSSecurityGroup,
				CidrBlock:        &CIDRBlock,
				Comment:          &Comment,
				GroupId:          &GroupID,
			},
			expectedResult: &projectipaccesslist.TfProjectIPAccessListDSModel{
				ID:               types.StringValue(AwsSecurityGroupID),
				ProjectID:        types.StringValue(GroupID),
				AWSSecurityGroup: types.StringValue(AWSSecurityGroup),
				CIDRBlock:        types.StringValue(CIDRBlock),
				IPAddress:        types.StringValue(""),
				Comment:          types.StringValue(Comment),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, _ := projectipaccesslist.NewTfProjectIPAccessListDSModel(context.Background(), tc.sdkModel)

			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}
