package controlplaneipaddresses_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/controlplaneipaddresses"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115014/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.ControlPlaneIPAddresses
	expectedTFModel *controlplaneipaddresses.TFControlPlaneIpAddressesModel
	name            string
}

func TestControlPlaneIpAddressesSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFModelTestCase{
		{
			name:    "Complete SDK response",
			SDKResp: &admin.ControlPlaneIPAddresses{},
			expectedTFModel: &controlplaneipaddresses.TFControlPlaneIpAddressesModel{
				Inbound: controlplaneipaddresses.InboundValue{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := controlplaneipaddresses.NewTFControlPlaneIPAddresses(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
