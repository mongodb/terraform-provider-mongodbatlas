package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		CreateContext: resourceMongoDBAtlasLDAPConfigurationCreate,
		ReadContext:   resourceMongoDBAtlasLDAPConfigurationRead,
		UpdateContext: resourceMongoDBAtlasLDAPConfigurationUpdate,
		DeleteContext: resourceMongoDBAtlasLDAPConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Optional: true,
				Default:  636,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
	}
}

func resourceMongoDBAtlasLDAPConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	ldap := &matlas.LDAP{}

	if v, ok := d.GetOk("authentication_enabled"); ok {
		ldap.AuthenticationEnabled = v.(bool)
	}

	if v, ok := d.GetOk("authorization_enabled"); ok {
		ldap.AuthorizationEnabled = v.(bool)
	}

	if v, ok := d.GetOk("hostname"); ok {
		ldap.Hostname = v.(string)
	}

	if v, ok := d.GetOk("port"); ok {
		ldap.Port = v.(int)
	}

	if v, ok := d.GetOk("bind_username"); ok {
		ldap.BindUsername = v.(string)
	}

	if v, ok := d.GetOk("bind_password"); ok {
		ldap.BindPassword = v.(string)
	}

	if v, ok := d.GetOk("ca_certificate"); ok {
		ldap.CaCertificate = v.(string)
	}

	if v, ok := d.GetOk("authz_query_template"); ok {
		ldap.AuthzQueryTemplate = v.(string)
	}

	if v, ok := d.GetOk("user_to_dn_mapping"); ok {
		ldap.UserToDNMapping = expandDNMapping(v.([]interface{}))
	}

	ladpReq := &matlas.LDAPConfiguration{
		LDAP: ldap,
	}

	_, _, err := conn.LDAPConfigurations.Save(ctx, projectID, ladpReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationCreate, projectID, err))
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasLDAPConfigurationRead(ctx, d, meta)
}

func resourceMongoDBAtlasLDAPConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ldapResp, resp, err := conn.LDAPConfigurations.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationRead, d.Id(), err))
	}

	if err = d.Set("authentication_enabled", ldapResp.LDAP.AuthenticationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authentication_enabled", d.Id(), err))
	}
	if err = d.Set("authorization_enabled", ldapResp.LDAP.AuthorizationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authorization_enabled", d.Id(), err))
	}
	if err = d.Set("hostname", ldapResp.LDAP.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "hostname", d.Id(), err))
	}
	if err = d.Set("port", ldapResp.LDAP.Port); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "port", d.Id(), err))
	}
	if err = d.Set("bind_username", ldapResp.LDAP.BindUsername); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "bind_username", d.Id(), err))
	}
	if err = d.Set("ca_certificate", ldapResp.LDAP.CaCertificate); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "ca_certificate", d.Id(), err))
	}
	if err = d.Set("authz_query_template", ldapResp.LDAP.AuthzQueryTemplate); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authz_query_template", d.Id(), err))
	}
	if err = d.Set("user_to_dn_mapping", flattenDNMapping(ldapResp.LDAP.UserToDNMapping)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "user_to_dn_mapping", d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasLDAPConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection.
	conn := meta.(*MongoDBClient).Atlas

	ldap := &matlas.LDAP{}

	if d.HasChange("authentication_enabled") {
		ldap.AuthenticationEnabled = d.Get("").(bool)
	}

	if d.HasChange("authorization_enabled") {
		ldap.AuthorizationEnabled = d.Get("authorization_enabled").(bool)
	}

	if d.HasChange("hostname") {
		ldap.Hostname = d.Get("hostname").(string)
	}

	if d.HasChange("port") {
		ldap.Port = d.Get("port").(int)
	}

	if d.HasChange("bind_username") {
		ldap.BindUsername = d.Get("bind_username").(string)
	}

	if d.HasChange("bind_password") {
		ldap.BindPassword = d.Get("bind_password").(string)
	}

	if d.HasChange("ca_certificate") {
		ldap.CaCertificate = d.Get("ca_certificate").(string)
	}

	if d.HasChange("authz_query_template") {
		ldap.AuthzQueryTemplate = d.Get("authz_query_template").(string)
	}

	if d.HasChange("user_to_dn_mapping") {
		ldap.UserToDNMapping = expandDNMapping(d.Get("user_to_dn_mapping").([]interface{}))
	}

	ldapReq := &matlas.LDAPConfiguration{
		LDAP: ldap,
	}

	_, _, err := conn.LDAPConfigurations.Save(ctx, d.Id(), ldapReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationUpdate, d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasLDAPConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection.
	conn := meta.(*MongoDBClient).Atlas
	_, _, err := conn.LDAPConfigurations.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationDelete, d.Id(), err))
	}

	return nil
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
