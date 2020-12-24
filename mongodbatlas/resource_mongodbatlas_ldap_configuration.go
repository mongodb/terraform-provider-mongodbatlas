package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorLDAPConfigurationCreate  = "error creating MongoDB LDAPConfiguration (%s): %s"
	errorLDAPConfigurationUpdate  = "error updating MongoDB LDAPConfiguration (%s): %s"
	errorLDAPConfigurationRead    = "error reading MongoDB LDAPConfiguration (%s): %s"
	errorLDAPConfigurationDelete  = "error deleting MongoDB LDAPConfiguration (%s): %s"
	errorLDAPConfigurationSetting = "error setting `%s` for LDAPConfiguration(%s): %s"
)

func resourceMongoDBAtlasLDAPConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasLDAPConfigurationCreate,
		Read:   resourceMongoDBAtlasLDAPConfigurationRead,
		Update: resourceMongoDBAtlasLDAPConfigurationUpdate,
		Delete: resourceMongoDBAtlasLDAPConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ldap": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"authorization_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"bind_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bind_password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ca_certificate": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"authz_query_template": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"user_to_dn_mapping": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"substitution": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ldap_query": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasLDAPConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	ladpReq := &matlas.LDAPConfiguration{
		LDAP: expandLADP(d),
	}

	_, _, err := conn.LDAPConfigurations.Save(context.Background(), projectID, ladpReq)
	if err != nil {
		return fmt.Errorf(errorLDAPConfigurationCreate, projectID, err)
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasLDAPConfigurationRead(d, meta)
}

func resourceMongoDBAtlasLDAPConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ldapResp, _, err := conn.LDAPConfigurations.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorLDAPConfigurationRead, d.Id(), err)
	}

	if err := d.Set("ldap", flattenLDAP(ldapResp.LDAP)); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "ldap", d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasLDAPConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get the client connection.
	conn := meta.(*matlas.Client)

	ldapReq := &matlas.LDAPConfiguration{}

	if d.HasChange("ldap") {
		ldapReq.LDAP = expandLADP(d)
	}

	_, _, err := conn.LDAPConfigurations.Save(context.Background(), d.Id(), ldapReq)
	if err != nil {
		return fmt.Errorf(errorLDAPConfigurationUpdate, d.Id(), err)
	}

	return nil
}

func resourceMongoDBAtlasLDAPConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	// Get the client connection.
	conn := meta.(*matlas.Client)
	_, _, err := conn.LDAPConfigurations.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf(errorLDAPConfigurationDelete, d.Id(), err)
	}

	return nil
}

func expandLADP(d *schema.ResourceData) *matlas.LDAP {
	ldap := &matlas.LDAP{}

	if v, ok := d.GetOk("ldap.0.authentication_enabled"); ok {
		ldap.AuthenticationEnabled = v.(bool)
	}

	if v, ok := d.GetOk("ldap.0.authorization_enabled"); ok {
		ldap.AuthorizationEnabled = v.(bool)
	}

	if v, ok := d.GetOk("ldap.0.hostname"); ok {
		ldap.Hostname = v.(string)
	}

	if v, ok := d.GetOk("ldap.0.port"); ok {
		ldap.Port = v.(int)
	}

	if v, ok := d.GetOk("ldap.0.bind_username"); ok {
		ldap.BindUsername = v.(string)
	}

	if v, ok := d.GetOk("ldap.0.bind_password"); ok {
		ldap.BindPassword = v.(string)
	}

	if v, ok := d.GetOk("ldap.0.ca_certificate"); ok {
		ldap.CaCertificate = v.(string)
	}

	if v, ok := d.GetOk("ldap.0.authz_query_template"); ok {
		ldap.AuthzQueryTemplate = v.(string)
	}

	if v, ok := d.GetOk("ldap.0.user_to_dn_mapping"); ok {
		ldap.UserToDNMapping = expandDNMapping(v.([]interface{}))
	}

	return ldap
}

func expandDNMapping(p []interface{}) []*matlas.UserToDNMapping {
	mappings := make([]*matlas.UserToDNMapping, len(p))

	for k, v := range p {
		mapping := v.(map[string]interface{})
		mappings[k] = &matlas.UserToDNMapping{
			Match:        mapping["match"].(string),
			Substitution: mapping["substitution"].(string),
			LDAPQuery:    mapping["ldap_query"].(string),
		}
	}

	return mappings
}

func flattenLDAP(ldap *matlas.LDAP) []map[string]interface{} {
	policyItems := make([]map[string]interface{}, 1)
	policyItems = append(policyItems, map[string]interface{}{
		"authentication_enabled": ldap.AuthenticationEnabled,
		"authorization_enabled":  ldap.AuthorizationEnabled,
		"hostname":               ldap.Hostname,
		"port":                   ldap.Port,
		"bind_username":          ldap.BindUsername,
		"bind_password":          ldap.BindUsername,
		"ca_certificate":         ldap.CaCertificate,
		"user_to_dn_mapping":     flattenDNMapping(ldap.UserToDNMapping),
	})

	return policyItems
}

func flattenDNMapping(usersDNMappings []*matlas.UserToDNMapping) []map[string]interface{} {
	usersDN := make([]map[string]interface{}, 0)
	for _, v := range usersDNMappings {
		usersDN = append(usersDN, map[string]interface{}{
			"match":        v.Match,
			"substitution": v.Substitution,
			"ldap_query":   v.LDAPQuery,
		})
	}

	return usersDN
}
