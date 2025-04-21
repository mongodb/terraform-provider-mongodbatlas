package controlplaneipaddresses

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func NewTFControlPlaneIPAddresses(ctx context.Context, apiResp *admin.ControlPlaneIPAddresses) (*TFControlPlaneIpAddressesModel, diag.Diagnostics) {
	inbound := apiResp.GetInbound()
	inboundAwsTfMap, inAwsDiags := conversion.ToTFMapOfSlices(ctx, inbound.GetAws())
	inboundGcpTfMap, inGcpDiags := conversion.ToTFMapOfSlices(ctx, inbound.GetGcp())
	inboundAzureTfMap, inAzureDiags := conversion.ToTFMapOfSlices(ctx, inbound.GetAzure())

	outbound := apiResp.GetOutbound()
	outboundAwsTfMap, outAwsDiags := conversion.ToTFMapOfSlices(ctx, outbound.GetAws())
	outboundGcpTfMap, outGcpDiags := conversion.ToTFMapOfSlices(ctx, outbound.GetGcp())
	outboundAzureTfMap, outAzureDiags := conversion.ToTFMapOfSlices(ctx, outbound.GetAzure())

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
