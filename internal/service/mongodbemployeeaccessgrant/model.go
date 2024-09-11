package mongodbemployeeaccessgrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func NewTFModel(projectID, clusterName string, apiResp *admin.EmployeeAccessGrant) (*TFModel, error) {
	id, err := conversion.IDWithProjectIDClusterName(projectID, clusterName)
	if err != nil {
		return nil, err
	}
	return &TFModel{
		ID:             types.StringValue(id),
		ProjectID:      types.StringValue(projectID),
		ClusterName:    types.StringValue(clusterName),
		GrantType:      types.StringValue(apiResp.GetGrantType()),
		ExpirationTime: types.StringValue(conversion.TimeToString(apiResp.GetExpirationTime())),
	}, nil
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
