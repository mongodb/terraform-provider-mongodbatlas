package encryptionatrest

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

type EarService interface {
	UpdateEncryptionAtRest(ctx context.Context, groupID string, encryptionAtRest *admin.EncryptionAtRest) (*admin.EncryptionAtRest, *http.Response, error)
}

type EarServiceFromClient struct {
	client *admin.APIClient
}

func (a *EarServiceFromClient) UpdateEncryptionAtRest(ctx context.Context, groupID string, encryptionAtRest *admin.EncryptionAtRest) (*admin.EncryptionAtRest, *http.Response, error) {
	return a.client.EncryptionAtRestUsingCustomerKeyManagementApi.UpdateEncryptionAtRest(ctx, groupID, encryptionAtRest).Execute()
}

func ServiceFromClient(client *admin.APIClient) EarService {
	return &EarServiceFromClient{
		client: client,
	}
}
