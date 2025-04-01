package mongodbemployeeaccessgrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func NewTFModel(projectID, clusterName string, apiResp *admin.EmployeeAccessGrant) *TFModel {
	return &TFModel{
		ProjectID:      types.StringValue(projectID),
		ClusterName:    types.StringValue(clusterName),
		GrantType:      types.StringValue(apiResp.GetGrantType()),
		ExpirationTime: types.StringValue(conversion.TimeToString(apiResp.GetExpirationTime())),
	}
}

func NewAtlasReq(tfModel *TFModel) (*admin.EmployeeAccessGrant, error) {
	expirationTimeStr := tfModel.ExpirationTime.ValueString()
	expirationTime, ok := conversion.StringToTime(expirationTimeStr)
	if !ok {
		return nil, fmt.Errorf("expiration_time format is incorrect: %s", expirationTimeStr)
	}
	return &admin.EmployeeAccessGrant{
		GrantType:      tfModel.GrantType.ValueString(),
		ExpirationTime: expirationTime,
	}, nil
}
