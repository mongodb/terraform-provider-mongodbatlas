package organization3

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312020/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func createOrganization(
	ctx context.Context,
	client *admin.APIClient,
	name, orgOwnerID string,
	secretExpiresAfterHours int,
) (*admin.CreateOrganizationResponse, *admin.ServiceAccountSecret, error) {
	req := admin.NewCreateOrganizationRequest(name)
	req.OrgOwnerId = new(orgOwnerID)
	skipDefaultAlerts := true
	req.SkipDefaultAlertsSettings = &skipDefaultAlerts
	req.SetServiceAccount(admin.OrgServiceAccountRequest{
		Name:                    name,
		Description:             fmt.Sprintf("organization3 SA for %s", name),
		Roles:                   []string{"ORG_OWNER"},
		SecretExpiresAfterHours: secretExpiresAfterHours,
	})
	resp, _, err := client.OrganizationsApi.CreateOrg(ctx, req).Execute()
	if err != nil {
		return nil, nil, err
	}
	sa, ok := resp.GetServiceAccountOk()
	if !ok {
		return resp, nil, fmt.Errorf("service account was not returned by CreateOrg")
	}
	secrets := sa.GetSecrets()
	if len(secrets) == 0 {
		return resp, nil, fmt.Errorf("service account secret was not returned by CreateOrg")
	}
	initial := secrets[0]
	return resp, &initial, nil
}

func getOrganization(ctx context.Context, client *admin.APIClient, orgID string) (*admin.AtlasOrganization, error) {
	org, resp, err := client.OrganizationsApi.GetOrg(ctx, orgID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil, nil
		}
		return nil, err
	}
	return org, nil
}

func updateOrganizationName(ctx context.Context, client *admin.APIClient, orgID, name string) error {
	_, _, err := client.OrganizationsApi.UpdateOrg(ctx, orgID, &admin.AtlasOrganization{Name: name}).Execute()
	return err
}

func deleteOrganization(ctx context.Context, client *admin.APIClient, orgID string) error {
	_, err := client.OrganizationsApi.DeleteOrg(ctx, orgID).Execute()
	return err
}

func getServiceAccount(ctx context.Context, client *admin.APIClient, orgID, clientID string) (*admin.OrgServiceAccount, error) {
	sa, resp, err := client.ServiceAccountsApi.GetOrgServiceAccount(ctx, orgID, clientID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}

func createServiceAccountSecret(
	ctx context.Context,
	client *admin.APIClient,
	orgID, clientID string,
	secretExpiresAfterHours int,
) (*admin.ServiceAccountSecret, error) {
	req := admin.NewServiceAccountSecretRequest(secretExpiresAfterHours)
	secret, _, err := client.ServiceAccountsApi.CreateOrgSecret(ctx, orgID, clientID, req).Execute()
	return secret, err
}

func deleteServiceAccountSecret(ctx context.Context, client *admin.APIClient, orgID, clientID, secretID string) error {
	_, err := client.ServiceAccountsApi.DeleteOrgSecret(ctx, clientID, secretID, orgID).Execute()
	return err
}
