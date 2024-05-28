package federatedsettingsidentityprovider

import (
	"sort"
	"strings"

	admin20231115008 "go.mongodb.org/atlas-sdk/v20231115014/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func FlattenFederatedSettingsIdentityProvider(federatedSettingsIdentityProvider []admin20231115008.FederationIdentityProvider) []map[string]any {
	var federatedSettingsIdentityProviderMap []map[string]any
	if len(federatedSettingsIdentityProvider) > 0 {
		federatedSettingsIdentityProviderMap = make([]map[string]any, len(federatedSettingsIdentityProvider))

		for i := range federatedSettingsIdentityProvider {
			federatedSettingsIdentityProviderMap[i] = map[string]any{
				"acs_url":                      federatedSettingsIdentityProvider[i].AcsUrl,
				"associated_domains":           federatedSettingsIdentityProvider[i].AssociatedDomains,
				"associated_orgs":              FlattenAssociatedOrgs(federatedSettingsIdentityProvider[i].GetAssociatedOrgs()),
				"audience_uri":                 federatedSettingsIdentityProvider[i].AudienceUri,
				"display_name":                 federatedSettingsIdentityProvider[i].DisplayName,
				"issuer_uri":                   federatedSettingsIdentityProvider[i].IssuerUri,
				"okta_idp_id":                  federatedSettingsIdentityProvider[i].OktaIdpId,
				"pem_file_info":                FlattenPemFileInfo(federatedSettingsIdentityProvider[i].GetPemFileInfo()),
				"request_binding":              federatedSettingsIdentityProvider[i].RequestBinding,
				"response_signature_algorithm": federatedSettingsIdentityProvider[i].ResponseSignatureAlgorithm,
				"sso_debug_enabled":            federatedSettingsIdentityProvider[i].SsoDebugEnabled,
				"sso_url":                      federatedSettingsIdentityProvider[i].SsoUrl,
				"status":                       federatedSettingsIdentityProvider[i].Status,
				"idp_id":                       federatedSettingsIdentityProvider[i].Id,
				"protocol":                     federatedSettingsIdentityProvider[i].Protocol,
				"audience_claim":               federatedSettingsIdentityProvider[i].AudienceClaim,
				"client_id":                    federatedSettingsIdentityProvider[i].ClientId,
				"groups_claim":                 federatedSettingsIdentityProvider[i].GroupsClaim,
				"requested_scopes":             federatedSettingsIdentityProvider[i].RequestedScopes,
				"user_claim":                   federatedSettingsIdentityProvider[i].UserClaim,
			}
		}
	}

	return federatedSettingsIdentityProviderMap
}

func FlattenAssociatedOrgs(associatedOrgs []admin20231115008.ConnectedOrgConfig) []map[string]any {
	var associatedOrgsMap []map[string]any

	if len(associatedOrgs) == 0 {
		return nil
	}
	associatedOrgsMap = make([]map[string]any, len(associatedOrgs))

	for i := range associatedOrgs {
		if associatedOrgs[i].UserConflicts == nil {
			associatedOrgsMap[i] = map[string]any{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderId,
				"org_id":                     associatedOrgs[i].OrgId,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              FlattenAuthFederationRoleMapping(associatedOrgs[i].GetRoleMappings()),
				"user_conflicts":             nil,
			}
		} else {
			associatedOrgsMap[i] = map[string]any{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderId,
				"org_id":                     associatedOrgs[i].OrgId,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              FlattenAuthFederationRoleMapping(associatedOrgs[i].GetRoleMappings()),
				"user_conflicts":             FlattenFederatedUser(associatedOrgs[i].GetUserConflicts()),
			}
		}
	}

	return associatedOrgsMap
}

type mRoleAssignmentV2 []admin20231115008.RoleAssignment

func (ra mRoleAssignmentV2) Len() int      { return len(ra) }
func (ra mRoleAssignmentV2) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }
func (ra mRoleAssignmentV2) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].GetOrgId(), ra[j].GetOrgId())

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GetGroupId(), ra[j].GetGroupId())

	if compareVal != 0 {
		return compareVal < 0
	}

	return *ra[i].Role < *ra[j].Role
}

func FlattenRoleAssignments(roleAssignments []admin20231115008.RoleAssignment) []map[string]any {
	sort.Sort(mRoleAssignmentV2(roleAssignments))

	var roleAssignmentsMap []map[string]any

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]any{
				"group_id": roleAssignments[i].GroupId,
				"org_id":   roleAssignments[i].OrgId,
				"role":     roleAssignments[i].Role,
			}
		}
	}

	return roleAssignmentsMap
}

func FlattenFederatedUser(federatedUsers []admin20231115008.FederatedUser) []map[string]any {
	var userConflictsMap []map[string]any

	if len(federatedUsers) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]any, len(federatedUsers))

	for i := range federatedUsers {
		userConflictsMap[i] = map[string]any{
			"email_address":          federatedUsers[i].EmailAddress,
			"federation_settings_id": federatedUsers[i].FederationSettingsId,
			"first_name":             federatedUsers[i].FirstName,
			"last_name":              federatedUsers[i].LastName,
			"user_id":                federatedUsers[i].UserId,
		}
	}

	return userConflictsMap
}

type authFederationoleMappingsByGroupName []admin20231115008.AuthFederationRoleMapping

func (ra authFederationoleMappingsByGroupName) Len() int      { return len(ra) }
func (ra authFederationoleMappingsByGroupName) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra authFederationoleMappingsByGroupName) Less(i, j int) bool {
	return ra[i].ExternalGroupName < ra[j].ExternalGroupName
}

func FlattenAuthFederationRoleMapping(roleMappings []admin20231115008.AuthFederationRoleMapping) []map[string]any {
	sort.Sort(authFederationoleMappingsByGroupName(roleMappings))

	var roleMappingsMap []map[string]any

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]any, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]any{
				"external_group_name": roleMappings[i].ExternalGroupName,
				"id":                  roleMappings[i].Id,
				"role_assignments":    FlattenRoleAssignments(*roleMappings[i].RoleAssignments),
			}
		}
	}

	return roleMappingsMap
}

func FlattenPemFileInfo(pemFileInfo admin20231115008.PemFileInfo) []map[string]any {
	var pemFileInfoMap []map[string]any

	if certificates := pemFileInfo.GetCertificates(); len(certificates) > 0 {
		pemFileInfoMap = make([]map[string]any, 1)

		pemFileInfoMap[0] = map[string]any{
			"certificates": FlattenFederatedSettingsCertificates(certificates),
			"file_name":    pemFileInfo.FileName,
		}
	}

	return pemFileInfoMap
}

func FlattenFederatedSettingsCertificates(certificates []admin20231115008.X509Certificate) []map[string]any {
	var certificatesMap []map[string]any

	if len(certificates) > 0 {
		certificatesMap = make([]map[string]any, len(certificates))

		for i := range certificates {
			certificatesMap[i] = map[string]any{
				"not_after":  conversion.TimePtrToStringPtr(certificates[i].NotAfter),
				"not_before": conversion.TimePtrToStringPtr(certificates[i].NotBefore),
			}
		}
	}

	return certificatesMap
}
