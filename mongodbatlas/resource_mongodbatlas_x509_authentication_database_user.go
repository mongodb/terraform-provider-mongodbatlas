package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorX509AuthDBUsersCreate         = "error creating MongoDB X509 Authentication for DB User(%s) in the project(%s): %s"
	errorX509AuthDBUsersRead           = "error reading MongoDB X509 Authentication for DB Users(%s) in the project(%s): %s"
	errorX509AuthDBUsersSetting        = "error setting `%s` for MongoDB X509 Authentication DB User(%s): %s"
	errorCustomerX509AuthDBUsersCreate = "error creating Customer X509 Authentication in the project(%s): %s"
	errorCustomerX509AuthDBUsersRead   = "error reading Customer X509 Authentication in the project(%s): %s"
	errorCustomerX509AuthDBUsersDelete = "error deleting Customer X509 Authentication in the project(%s): %s"
)

func resourceMongoDBAtlasX509AuthDBUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasX509AuthDBUserCreate,
		ReadContext:   resourceMongoDBAtlasX509AuthDBUserRead,
		DeleteContext: resourceMongoDBAtlasX509AuthDBUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasX509AuthDBUserImportState,
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
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

func resourceMongoDBAtlasX509AuthDBUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)

	var currentCertificate string

	if expirationMonths, ok := d.GetOk("months_until_expiration"); ok {
		res, _, err := conn.X509AuthDBUsers.CreateUserCertificate(ctx, projectID, username, expirationMonths.(int))
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersCreate, username, projectID, err))
		}

		currentCertificate = res.Certificate
	} else {
		customerX509Cas := d.Get("customer_x509_cas").(string)
		_, _, err := conn.X509AuthDBUsers.SaveConfiguration(ctx, projectID, &matlas.CustomerX509{Cas: customerX509Cas})
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorCustomerX509AuthDBUsersCreate, projectID, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"username":            username,
		"current_certificate": currentCertificate,
	}))

	return resourceMongoDBAtlasX509AuthDBUserRead(ctx, d, meta)
}

func resourceMongoDBAtlasX509AuthDBUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	username := ids["username"]
	currentCertificate := ids["current_certificate"]

	var (
		certificates []matlas.UserCertificate
		err          error
	)

	if username != "" {
		certificates, _, err = conn.X509AuthDBUsers.GetUserCertificates(ctx, projectID, username)
		if err != nil {
			// new resource missing
			reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()
			if reset {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err))
		}
	}

	if err := d.Set("current_certificate", cast.ToString(currentCertificate)); err != nil {
		return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "current_certificate", username, err))
	}

	if err := d.Set("certificates", flattenCertificates(certificates)); err != nil {
		return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err))
	}

	return nil
}

func resourceMongoDBAtlasX509AuthDBUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	ids := decodeStateID(d.Id())
	currentCertificate := ids["current_certificate"]
	projectID := ids["project_id"]

	if currentCertificate == "" {
		_, err := conn.X509AuthDBUsers.DisableCustomerX509(ctx, projectID)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorCustomerX509AuthDBUsersDelete, projectID, err))
		}
	}

	d.SetId("")

	return nil
}

func resourceMongoDBAtlasX509AuthDBUserImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

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
		_, _, err := conn.X509AuthDBUsers.GetUserCertificates(ctx, projectID, username)
		if err != nil {
			return nil, fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err)
		}

		if err := d.Set("username", username); err != nil {
			return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "username", username, err)
		}
	}

	customerX509, _, err := conn.X509AuthDBUsers.GetCurrentX509Conf(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(errorCustomerX509AuthDBUsersRead, projectID, err)
	}

	if err := d.Set("customer_x509_cas", customerX509.Cas); err != nil {
		return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorX509AuthDBUsersSetting, "project_id", username, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"username":            username,
		"current_certificate": "",
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenCertificates(userCertificates []matlas.UserCertificate) []map[string]interface{} {
	certificates := make([]map[string]interface{}, len(userCertificates))
	for i, v := range userCertificates {
		certificates[i] = map[string]interface{}{
			"id":         v.ID,
			"created_at": v.CreatedAt,
			"group_id":   v.GroupID,
			"not_after":  v.NotAfter,
			"subject":    v.Subject,
		}
	}

	return certificates
}
