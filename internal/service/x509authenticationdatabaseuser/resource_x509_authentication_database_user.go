package x509authenticationdatabaseuser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

const (
	errorX509AuthDBUsersCreate         = "error creating MongoDB X509 Authentication for DB User(%s) in the project(%s): %s"
	errorX509AuthDBUsersRead           = "error reading MongoDB X509 Authentication for DB Users(%s) in the project(%s): %s"
	errorX509AuthDBUsersSetting        = "error setting `%s` for MongoDB X509 Authentication DB User(%s): %s"
	errorCustomerX509AuthDBUsersCreate = "error creating Customer X509 Authentication in the project(%s): %s"
	errorCustomerX509AuthDBUsersRead   = "error reading Customer X509 Authentication in the project(%s): %s"
	errorCustomerX509AuthDBUsersDelete = "error deleting Customer X509 Authentication in the project(%s): %s"
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
				ForceNew: true,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"months_until_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 24 {
						errs = append(errs, fmt.Errorf("%q value should be between 1 and 24, got: %d", key, v))
					}
					return
				},
			},
			"current_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"customer_x509_cas": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{"months_until_expiration", "username"},
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"not_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)

	if expirationMonths, ok := d.GetOk("months_until_expiration"); ok {
		months := expirationMonths.(int)
		params := &admin.UserCert{
			MonthsUntilExpiration: &months,
		}
		certStr, _, err := connV2.X509AuthenticationApi.CreateDatabaseUserCertificate(ctx, projectID, username, params).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersCreate, username, projectID, err))
		}
		if err := d.Set("current_certificate", cast.ToString(certStr)); err != nil {
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "current_certificate", username, err))
		}
	} else {
		customerX509Cas := d.Get("customer_x509_cas").(string)
		userReq := &admin.UserSecurity{
			CustomerX509: &admin.DBUserTLSX509Settings{Cas: &customerX509Cas},
		}
		_, _, err := connV2.LDAPConfigurationApi.SaveLdapConfiguration(ctx, projectID, userReq).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorCustomerX509AuthDBUsersCreate, projectID, err))
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"username":      username,
		"serial_number": "", // not returned in create API, got later in Read
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	var (
		certificates []admin.UserCert
		serialNumber string
	)

	if username != "" {
		resp, _, err := connV2.X509AuthenticationApi.ListDatabaseUserCertificates(ctx, projectID, username).Execute()
		if err != nil {
			// new resource missing
			reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()
			if reset {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err))
		}
		if resp != nil && resp.Results != nil {
			certificates = *resp.Results
			if len(certificates) > 0 {
				serialNumber = cast.ToString(certificates[len(certificates)-1].GetId()) // Get SerialId from last user certificate
			}
		}
	}
	if err := d.Set("certificates", flattenCertificates(certificates)); err != nil {
		return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"username":      username,
		"serial_number": serialNumber,
	}))

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// We don't do anything because X.509 certificates can not be deleted or disassociated from a user.
	// More info: https://jira.mongodb.org/browse/HELP-53363
	d.SetId("")
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 1 && len(parts) != 2 {
		return nil, errors.New("import format error: to import a X509 Authentication, use the formats {project_id} or {project_id}-{username}")
	}
	var username string
	if len(parts) == 2 {
		username = parts[1]
	}
	projectID := parts[0]

	if username != "" {
		_, _, err := connV2.X509AuthenticationApi.ListDatabaseUserCertificates(ctx, projectID, username).Execute()
		if err != nil {
			return nil, fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err)
		}

		if err := d.Set("username", username); err != nil {
			return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "username", username, err)
		}
	}

	resp, _, err := connV2.LDAPConfigurationApi.GetLdapConfiguration(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf(errorCustomerX509AuthDBUsersRead, projectID, err)
	}
	customerX509 := resp.GetCustomerX509()
	if err := d.Set("customer_x509_cas", customerX509.GetCas()); err != nil {
		return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "project_id", username, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":          projectID,
		"username":            username,
		"current_certificate": "",
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenCertificates(userCertificates []admin.UserCert) []map[string]any {
	certificates := make([]map[string]any, len(userCertificates))
	for i, v := range userCertificates {
		certificates[i] = map[string]any{
			"id":         v.GetId(),
			"created_at": conversion.TimePtrToStringPtr(v.CreatedAt),
			"group_id":   v.GetGroupId(),
			"not_after":  conversion.TimePtrToStringPtr(v.NotAfter),
			"subject":    v.GetSubject(),
		}
	}
	return certificates
}
