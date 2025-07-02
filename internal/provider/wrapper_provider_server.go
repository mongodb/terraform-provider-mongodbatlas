package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func NewWrappedProviderServer(old func() tfprotov6.ProviderServer) func() tfprotov6.ProviderServer {
	return func() tfprotov6.ProviderServer {
		return &WrappedProviderServer{
			OldServer: old(),
		}
	}
}

type WrappedProviderServer struct {
	OldServer tfprotov6.ProviderServer
}

func (s *WrappedProviderServer) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	return s.OldServer.GetMetadata(ctx, req)
}

func (s *WrappedProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return s.OldServer.GetProviderSchema(ctx, req)
}

func (s *WrappedProviderServer) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	return s.OldServer.GetResourceIdentitySchemas(ctx, req)
}

func (s *WrappedProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	return s.OldServer.ValidateProviderConfig(ctx, req)
}

func (s *WrappedProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	return s.OldServer.ConfigureProvider(ctx, req)
}

func (s *WrappedProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return s.OldServer.StopProvider(ctx, req)
}

func (s *WrappedProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.ValidateResourceConfig(ctx, req)
}

func (s *WrappedProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "upgrade."+req.TypeName)
	return s.OldServer.UpgradeResourceState(ctx, req)
}

func (s *WrappedProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.ReadResource(ctx, req)
}

func (s *WrappedProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.PlanResourceChange(ctx, req)
}

func (s *WrappedProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.ApplyResourceChange(ctx, req)
}

func (s *WrappedProviderServer) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "import."+req.TypeName)
	return s.OldServer.ImportResourceState(ctx, req)
}

func (s *WrappedProviderServer) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "move."+req.TargetTypeName)
	return s.OldServer.MoveResourceState(ctx, req)
}

func (s *WrappedProviderServer) UpgradeResourceIdentity(ctx context.Context, req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.UpgradeResourceIdentity(ctx, req)
}

func (s *WrappedProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "data."+req.TypeName)
	return s.OldServer.ValidateDataResourceConfig(ctx, req)
}

func (s *WrappedProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "data."+req.TypeName)
	return s.OldServer.ReadDataSource(ctx, req)
}

func (s *WrappedProviderServer) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, "func."+req.Name)
	return s.OldServer.CallFunction(ctx, req)
}

func (s *WrappedProviderServer) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	return s.OldServer.GetFunctions(ctx, req)
}

func (s *WrappedProviderServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.ValidateEphemeralResourceConfig(ctx, req)
}

func (s *WrappedProviderServer) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.OpenEphemeralResource(ctx, req)
}

func (s *WrappedProviderServer) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.RenewEphemeralResource(ctx, req)
}

func (s *WrappedProviderServer) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	ctx = context.WithValue(ctx, config.ContextKeyTFSrc, req.TypeName)
	return s.OldServer.CloseEphemeralResource(ctx, req)
}
