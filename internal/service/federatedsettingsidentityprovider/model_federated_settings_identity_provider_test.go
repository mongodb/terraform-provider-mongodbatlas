package federatedsettingsidentityprovider_test

import (
	"testing"
	"time"

	admin20231115008 "go.mongodb.org/atlas-sdk/v20231115014/admin"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
)

var (
	groupID                    = "group_id"
	organizationID             = "organization_id"
	role                       = "role"
	emailAddress               = "emailAddress"
	federationSettingsID       = "federationSettingsId"
	firstName                  = "firstName"
	lastName                   = "lastName"
	userID                     = "userId"
	roleAssignmentsID          = "roleAssignmentsID"
	externalGroupName          = "externalGroupName"
	notBefore                  = time.Now()
	notAfter                   = notBefore.Add(time.Hour * time.Duration(1))
	fileName                   = "fileName"
	domainAllowList            = []string{"domainAllowList"}
	domainRestrictionEnabled   = true
	identityProviderID         = "identityProviderId"
	postAuthRoleGrants         = []string{"postAuthRoleGrants"}
	acsURL                     = "acsURL"
	associatedDomains          = []string{"associatedDomains"}
	audienceURI                = "audienceURI"
	displayName                = "displayName"
	issuerURI                  = "issuerURI"
	oktaIdpID                  = "oktaIdpID"
	requestBinding             = "requestBinding"
	responseSignatureAlgorithm = "responseSignatureAlgorithm"
	ssoDebugEnabled            = true
	ssoURL                     = "ssoUrl"
	status                     = "ACTIVE"
	samlProtocol               = "SAML"
	oidcProtocol               = "OIDC"
	audienceClaim              = []string{"audienceClaim"}
	clientID                   = "clientId"
	groupsClaim                = "groupsClaim"
	requestedScopes            = []string{"requestedScopes"}
	userClaim                  = "userClaim"

	roleAssignments = []admin20231115008.RoleAssignment{
		{
			GroupId: &groupID,
			OrgId:   &organizationID,
			Role:    &role,
		},
	}
	flattenedRoleAssignments = []map[string]any{
		{
			"group_id": &groupID,
			"org_id":   &organizationID,
			"role":     &role,
		},
	}
	pemCertificates = []admin20231115008.X509Certificate{
		{
			NotAfter:  &notAfter,
			NotBefore: &notBefore,
		},
	}
	flattenedPemCertificates = []map[string]any{
		{
			"not_after":  conversion.TimePtrToStringPtr(&notAfter),
			"not_before": conversion.TimePtrToStringPtr(&notBefore),
		},
	}
	federationRoleMapping = []admin20231115008.AuthFederationRoleMapping{
		{
			ExternalGroupName: externalGroupName,
			Id:                &roleAssignmentsID,
			RoleAssignments:   &roleAssignments,
		},
	}
	flattenedFederationRoleMapping = []map[string]any{
		{
			"external_group_name": externalGroupName,
			"id":                  &roleAssignmentsID,
			"role_assignments":    flattenedRoleAssignments,
		},
	}
	federatedUser = []admin20231115008.FederatedUser{
		{
			EmailAddress:         emailAddress,
			FederationSettingsId: federationSettingsID,
			FirstName:            firstName,
			LastName:             lastName,
			UserId:               &userID,
		},
	}
	flattenedFederatedUser = []map[string]any{
		{
			"email_address":          emailAddress,
			"federation_settings_id": federationSettingsID,
			"first_name":             firstName,
			"last_name":              lastName,
			"user_id":                &userID,
		},
	}
	associatedOrgs = []admin20231115008.ConnectedOrgConfig{
		{
			DomainAllowList:          &domainAllowList,
			DomainRestrictionEnabled: domainRestrictionEnabled,
			IdentityProviderId:       identityProviderID,
			OrgId:                    organizationID,
			PostAuthRoleGrants:       &postAuthRoleGrants,
			RoleMappings:             &federationRoleMapping,
			UserConflicts:            nil,
		},
	}
	flattenedAssociatedOrgs = []map[string]any{
		{
			"domain_allow_list":          &domainAllowList,
			"domain_restriction_enabled": domainRestrictionEnabled,
			"identity_provider_id":       identityProviderID,
			"org_id":                     organizationID,
			"post_auth_role_grants":      &postAuthRoleGrants,
			"role_mappings":              flattenedFederationRoleMapping,
			"user_conflicts":             nil,
		},
	}
	pemFileInfo = admin20231115008.PemFileInfo{
		FileName:     &fileName,
		Certificates: &pemCertificates,
	}
	flattenedPemFileInfo = []map[string]any{
		{
			"certificates": flattenedPemCertificates,
			"file_name":    &fileName,
		},
	}
)

func TestFlattenRoleAssignments(t *testing.T) {
	testCases := []struct {
		name   string
		input  []admin20231115008.RoleAssignment
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenRoleAssignments",
			input:  roleAssignments,
			output: flattenedRoleAssignments,
		},
		{
			name:   "Empty FlattenRoleAssignments",
			input:  []admin20231115008.RoleAssignment{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenRoleAssignments(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenFederatedUser(t *testing.T) {
	testCases := []struct {
		name   string
		input  []admin20231115008.FederatedUser
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenFederatedUser",
			input:  federatedUser,
			output: flattenedFederatedUser,
		},
		{
			name:   "Empty FlattenFederatedUser",
			input:  []admin20231115008.FederatedUser{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenFederatedUser(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenAuthFederationRoleMapping(t *testing.T) {
	testCases := []struct {
		name   string
		input  []admin20231115008.AuthFederationRoleMapping
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenAuthFederationRoleMapping",
			input:  federationRoleMapping,
			output: flattenedFederationRoleMapping,
		},
		{
			name:   "Empty FlattenAuthFederationRoleMapping",
			input:  []admin20231115008.AuthFederationRoleMapping{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenAuthFederationRoleMapping(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenFederatedSettingsCertificates(t *testing.T) {
	testCases := []struct {
		name   string
		input  []admin20231115008.X509Certificate
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenFederatedSettingsCertificates",
			input:  pemCertificates,
			output: flattenedPemCertificates,
		},
		{
			name:   "Empty FlattenFederatedSettingsCertificates",
			input:  []admin20231115008.X509Certificate{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenFederatedSettingsCertificates(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenPemFileInfo(t *testing.T) {
	testCases := []struct {
		name   string
		input  admin20231115008.PemFileInfo
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenPemFileInfo",
			input:  pemFileInfo,
			output: flattenedPemFileInfo,
		},
		{
			name:   "Empty FlattenPemFileInfo",
			input:  admin20231115008.PemFileInfo{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenPemFileInfo(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenAssociatedOrgs(t *testing.T) {
	testCases := []struct {
		name   string
		input  []admin20231115008.ConnectedOrgConfig
		output []map[string]any
	}{
		{
			name:   "Non empty FlattenAssociatedOrgs without UserConflics",
			input:  associatedOrgs,
			output: flattenedAssociatedOrgs,
		},
		{
			name: "Non empty FlattenAssociatedOrgs with UserConflics",
			input: []admin20231115008.ConnectedOrgConfig{
				{
					DomainAllowList:          &domainAllowList,
					DomainRestrictionEnabled: domainRestrictionEnabled,
					IdentityProviderId:       identityProviderID,
					OrgId:                    organizationID,
					PostAuthRoleGrants:       &postAuthRoleGrants,
					RoleMappings:             &federationRoleMapping,
					UserConflicts:            &federatedUser,
				},
			},
			output: []map[string]any{
				{
					"domain_allow_list":          &domainAllowList,
					"domain_restriction_enabled": domainRestrictionEnabled,
					"identity_provider_id":       identityProviderID,
					"org_id":                     organizationID,
					"post_auth_role_grants":      &postAuthRoleGrants,
					"role_mappings":              flattenedFederationRoleMapping,
					"user_conflicts":             flattenedFederatedUser,
				},
			},
		},
		{
			name:   "Empty FlattenAssociatedOrgs",
			input:  []admin20231115008.ConnectedOrgConfig{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenAssociatedOrgs(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}

func TestFlattenFederatedSettingsIdentityProvider(t *testing.T) {
	var nilStringPtr *string
	var nilStringSlicePtr *[]string
	var nilString string
	var nilMap []map[string]any
	var nilBoolPtr *bool
	testCases := []struct {
		name   string
		input  []admin20231115008.FederationIdentityProvider
		output []map[string]any
	}{
		{
			name: "Non empty SAML FlattenFederatedSettingsIdentityProvider",
			input: []admin20231115008.FederationIdentityProvider{
				{
					AcsUrl:                     &acsURL,
					AssociatedDomains:          &associatedDomains,
					AssociatedOrgs:             &associatedOrgs,
					AudienceUri:                &audienceURI,
					DisplayName:                &displayName,
					IssuerUri:                  &issuerURI,
					OktaIdpId:                  oktaIdpID,
					PemFileInfo:                &pemFileInfo,
					RequestBinding:             &requestBinding,
					ResponseSignatureAlgorithm: &responseSignatureAlgorithm,
					SsoDebugEnabled:            &ssoDebugEnabled,
					SsoUrl:                     &ssoURL,
					Status:                     &status,
					Id:                         identityProviderID,
					Protocol:                   &samlProtocol,
				},
			},
			output: []map[string]any{
				{
					"acs_url":                      &acsURL,
					"associated_domains":           &associatedDomains,
					"associated_orgs":              flattenedAssociatedOrgs,
					"audience_uri":                 &audienceURI,
					"display_name":                 &displayName,
					"issuer_uri":                   &issuerURI,
					"okta_idp_id":                  oktaIdpID,
					"pem_file_info":                flattenedPemFileInfo,
					"request_binding":              &requestBinding,
					"response_signature_algorithm": &responseSignatureAlgorithm,
					"sso_debug_enabled":            &ssoDebugEnabled,
					"sso_url":                      &ssoURL,
					"status":                       &status,
					"idp_id":                       identityProviderID,
					"protocol":                     &samlProtocol,
					"audience_claim":               nilStringSlicePtr,
					"client_id":                    nilStringPtr,
					"groups_claim":                 nilStringPtr,
					"requested_scopes":             nilStringSlicePtr,
					"user_claim":                   nilStringPtr,
				},
			},
		},
		{
			name: "Non empty OIDC FlattenFederatedSettingsIdentityProvider",
			input: []admin20231115008.FederationIdentityProvider{
				{
					AssociatedDomains: &associatedDomains,
					AssociatedOrgs:    &associatedOrgs,
					DisplayName:       &displayName,
					IssuerUri:         &issuerURI,
					Id:                identityProviderID,
					Protocol:          &oidcProtocol,
					AudienceClaim:     &audienceClaim,
					ClientId:          &clientID,
					GroupsClaim:       &groupsClaim,
					RequestedScopes:   &requestedScopes,
					UserClaim:         &userClaim,
				},
			},
			output: []map[string]any{
				{
					"acs_url":                      nilStringPtr,
					"associated_domains":           &associatedDomains,
					"associated_orgs":              flattenedAssociatedOrgs,
					"audience_uri":                 nilStringPtr,
					"display_name":                 &displayName,
					"issuer_uri":                   &issuerURI,
					"okta_idp_id":                  nilString,
					"pem_file_info":                nilMap,
					"request_binding":              nilStringPtr,
					"response_signature_algorithm": nilStringPtr,
					"sso_debug_enabled":            nilBoolPtr,
					"sso_url":                      nilStringPtr,
					"status":                       nilStringPtr,
					"idp_id":                       identityProviderID,
					"protocol":                     &oidcProtocol,
					"audience_claim":               &audienceClaim,
					"client_id":                    &clientID,
					"groups_claim":                 &groupsClaim,
					"requested_scopes":             &requestedScopes,
					"user_claim":                   &userClaim,
				},
			},
		},
		{
			name:   "Empty FlattenFederatedSettingsIdentityProvider",
			input:  []admin20231115008.FederationIdentityProvider{},
			output: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := federatedsettingsidentityprovider.FlattenFederatedSettingsIdentityProvider(tc.input)
			assert.Equal(t, tc.output, resultModel)
		})
	}
}
