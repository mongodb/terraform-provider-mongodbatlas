package x509authenticationdatabaseuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"customer_x509_cas": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	username := d.Get("username").(string)

	if username != "" {
		resp, _, err := connV2.X509AuthenticationApi.ListDatabaseUserCerts(ctx, projectID, username).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersRead, username, projectID, err))
		}
		if resp != nil && resp.Results != nil {
			if err := d.Set("certificates", flattenCertificates(*resp.Results)); err != nil {
				return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err))
			}
		}
	}

	resp, _, err := connV2.LDAPConfigurationApi.GetUserSecurity(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCustomerX509AuthDBUsersRead, projectID, err))
	}
	customerX509 := resp.GetCustomerX509()
	if err := d.Set("customer_x509_cas", customerX509.GetCas()); err != nil {
		return diag.FromErr(fmt.Errorf(errorX509AuthDBUsersSetting, "certificates", username, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":          projectID,
		"username":            username,
		"current_certificate": "",
	}))

	return nil
}
