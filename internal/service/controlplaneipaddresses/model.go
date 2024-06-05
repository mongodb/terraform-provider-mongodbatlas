package controlplaneipaddresses

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20231115014/admin"
)

func NewTFControlPlaneIPAddresses(ctx context.Context, apiResp *admin.ControlPlaneIPAddresses) (*TFControlPlaneIpAddressesModel, diag.Diagnostics) {
	inbound := apiResp.GetInbound()
	inboundAwsTfMap, inAwsDiags := toTFMap(ctx, inbound.GetAws())
	inboundGcpTfMap, inGcpDiags := toTFMap(ctx, inbound.GetGcp())
	inboundAzureTfMap, inAzureDiags := toTFMap(ctx, inbound.GetAzure())

	outbound := apiResp.GetOutbound()
	outboundAwsTfMap, outAwsDiags := toTFMap(ctx, outbound.GetAws())
	outboundGcpTfMap, outGcpDiags := toTFMap(ctx, outbound.GetGcp())
	outboundAzureTfMap, outAzureDiags := toTFMap(ctx, outbound.GetAzure())

	allDiags := slices.Concat(inAwsDiags, inGcpDiags, inAzureDiags, outAwsDiags, outGcpDiags, outAzureDiags)
	if allDiags.HasError() {
		return nil, allDiags
	}

	return &TFControlPlaneIpAddressesModel{
		Inbound: InboundValue{
			Aws:   inboundAwsTfMap,
			Gcp:   inboundGcpTfMap,
			Azure: inboundAzureTfMap,
		},
		Outbound: OutboundValue{
			Aws:   outboundAwsTfMap,
			Gcp:   outboundGcpTfMap,
			Azure: outboundAzureTfMap,
		},
	}, nil
}

func toTFMap(ctx context.Context, values map[string][]string) (basetypes.MapValue, diag.Diagnostics) {
	return types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, values)
}
