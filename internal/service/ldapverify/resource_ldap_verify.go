package ldapverify

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

const (
	errorCreate   = "error creating MongoDB LDAPVerify (%s): %s"
	errorRead     = "error reading MongoDB LDAPVerify (%s): %s"
	errorSettings = "error setting `%s` for LDAPVerify(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"authz_query_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	params := new(admin.LDAPVerifyConnectivityJobRequestParams)

	if v, ok := d.GetOk("hostname"); ok {
		params.Hostname = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		params.Port = v.(int)
	}
	if v, ok := d.GetOk("bind_username"); ok {
		params.BindUsername = v.(string)
	}
	if v, ok := d.GetOk("bind_password"); ok {
		params.BindPassword = v.(string)
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		params.CaCertificate = conversion.Pointer(v.(string))
	}
	if v, ok := d.GetOk("authz_query_template"); ok {
		params.AuthzQueryTemplate = conversion.Pointer(v.(string))
	}

	ldap, _, err := connV2.LDAPConfigurationApi.VerifyLdapConfiguration(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, projectID, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"SUCCESS", "FAILED"},
		Refresh:    resourceRefreshFunc(ctx, projectID, ldap.GetRequestId(), connV2),
		Timeout:    3 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, projectID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldap.GetRequestId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	requestID := ids["request_id"]
	ldapResp, resp, err := connV2.LDAPConfigurationApi.GetLdapConfigurationStatus(context.Background(), projectID, requestID).Execute()
	if err != nil || ldapResp == nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorRead, d.Id(), err))
	}

	if err := d.Set("hostname", ldapResp.Request.GetHostname()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "hostname", d.Id(), err))
	}
	if err := d.Set("port", ldapResp.Request.Port); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "port", d.Id(), err))
	}
	if err := d.Set("bind_username", ldapResp.Request.GetBindUsername()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "bind_username", d.Id(), err))
	}
	if err := d.Set("links", conversion.FlattenLinks(ldapResp.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "links", d.Id(), err))
	}
	if err := d.Set("validations", flattenValidations(ldapResp.GetValidations())); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "validations", d.Id(), err))
	}
	if err := d.Set("request_id", ldapResp.GetRequestId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "request_id", d.Id(), err))
	}
	if err := d.Set("status", ldapResp.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "status", d.Id(), err))
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}

func flattenValidations(validations []admin.LDAPVerifyConnectivityJobRequestValidation) []map[string]string {
	ret := make([]map[string]string, len(validations))
	for i := range validations {
		validation := &validations[i]
		ret[i] = map[string]string{
			"status":          validation.GetStatus(),
			"validation_type": validation.GetValidationType(),
		}
	}
	return ret
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a LDAP Verify use the format {project_id}-{request_id}")
	}
	projectID := parts[0]
	requestID := parts[1]

	_, _, err := connV2.LDAPConfigurationApi.GetLdapConfigurationStatus(ctx, projectID, requestID).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorRead, requestID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorSettings, "project_id", requestID, err)
	}

	if err := d.Set("request_id", requestID); err != nil {
		return nil, fmt.Errorf(errorSettings, "request_id", requestID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": requestID,
	}))
	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, projectID, requestID string, connV2 *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		ldap, resp, err := connV2.LDAPConfigurationApi.GetLdapConfigurationStatus(ctx, projectID, requestID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		status := ldap.GetStatus()
		return ldap, status, nil
	}
}
