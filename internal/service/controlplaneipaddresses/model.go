package controlplaneipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20231115014/admin"
)

func NewTFControlPlaneIPAddresses(ctx context.Context, apiResp *admin.ControlPlaneIPAddresses) (*TFControlPlaneIpAddressesModel, diag.Diagnostics) {
	return &TFControlPlaneIpAddressesModel{}, nil
}
