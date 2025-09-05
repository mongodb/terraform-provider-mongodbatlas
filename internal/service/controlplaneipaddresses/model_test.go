package controlplaneipaddresses_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/controlplaneipaddresses"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.ControlPlaneIPAddresses
	expectedTFModel *controlplaneipaddresses.TFControlPlaneIpAddressesModel
}

func TestControlPlaneIpAddressesSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.ControlPlaneIPAddresses{
				Inbound: &admin.InboundControlPlaneCloudProviderIPAddresses{
					Aws: &map[string][]string{
						"some-region": {"inbound-aws-value"},
					},
					Azure: &map[string][]string{
						"some-region": {"inbound-azure-value"},
					},
					Gcp: &map[string][]string{
						"some-region": {"inbound-gcp-value"},
					},
				},
				Outbound: &admin.OutboundControlPlaneCloudProviderIPAddresses{
					Aws: &map[string][]string{
						"some-region": {"outbound-aws-value"},
					},
					Azure: &map[string][]string{
						"some-region": {"outbound-azure-value"},
					},
					Gcp: &map[string][]string{
						"some-region": {"outbound-gcp-value"},
					},
				},
			},
			expectedTFModel: &controlplaneipaddresses.TFControlPlaneIpAddressesModel{
				Inbound: controlplaneipaddresses.InboundValue{
					Aws: toTFMap(t, map[string][]string{
						"some-region": {"inbound-aws-value"},
					}),
					Azure: toTFMap(t, map[string][]string{
						"some-region": {"inbound-azure-value"},
					}),
					Gcp: toTFMap(t, map[string][]string{
						"some-region": {"inbound-gcp-value"},
					}),
				},
				Outbound: controlplaneipaddresses.OutboundValue{
					Aws: toTFMap(t, map[string][]string{
						"some-region": {"outbound-aws-value"},
					}),
					Azure: toTFMap(t, map[string][]string{
						"some-region": {"outbound-azure-value"},
					}),
					Gcp: toTFMap(t, map[string][]string{
						"some-region": {"outbound-gcp-value"},
					}),
				},
			},
		},
		"Null response in a specifc providers and root outbound property": {
			SDKResp: &admin.ControlPlaneIPAddresses{
				Inbound: &admin.InboundControlPlaneCloudProviderIPAddresses{},
			},
			expectedTFModel: &controlplaneipaddresses.TFControlPlaneIpAddressesModel{
				Inbound: controlplaneipaddresses.InboundValue{
					Aws:   types.MapNull(types.ListType{ElemType: types.StringType}),
					Azure: types.MapNull(types.ListType{ElemType: types.StringType}),
					Gcp:   types.MapNull(types.ListType{ElemType: types.StringType}),
				},
				Outbound: controlplaneipaddresses.OutboundValue{
					Aws:   types.MapNull(types.ListType{ElemType: types.StringType}),
					Azure: types.MapNull(types.ListType{ElemType: types.StringType}),
					Gcp:   types.MapNull(types.ListType{ElemType: types.StringType}),
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := controlplaneipaddresses.NewTFControlPlaneIPAddresses(t.Context(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func toTFMap(t *testing.T, values map[string][]string) basetypes.MapValue {
	t.Helper()
	result, diags := conversion.ToTFMapOfSlices(t.Context(), values)
	if diags.HasError() {
		t.Errorf("unexpected errors found when creating test cases: %s", diags.Errors()[0].Summary())
	}
	return result
}
